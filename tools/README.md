# Compat tools

Each directory here holds a `run.sh` that checks part of the Go port against G3's C code at a pinned revision, and a README that says what it checks.
This page says which checks to run when, and which upstream files each tool reads, so that moving the pin re-records only the tools whose sources changed.

## Check tiers

**Routine task acceptance.**
Seconds to a few minutes; writes little to `tmp/`.

```sh
go vet ./...
go test ./...                        # recorded cases only; the whole-game tests skip
tools/<tool>/run.sh                  # each tool whose area the task touched, default mode
```

For load and save work, add `tools/olympia-roundtrip/run.sh` (about 25 seconds once its C cache is built).
`tools/box-compat/run.sh --games 10`, `tools/boxfile-compat/run.sh --games 10`, `tools/check-compat/run.sh --games 10` and `tools/loadsave-compat/run.sh --games 10` compare a fixed sample of saved games when a change could show only on real data.
For turn work, add `tools/turn-compat/run.sh`: small C recordings, comparator
contracts and a byte-identical diagnostic run, never a whole-world turn by default.
For clock/calendar or progress output work, also run `tools/turn-driver-compat/run.sh`.
For command parsing, dispatch, handlers or completion, also run
`tools/command-compat/run.sh` (small real-C dispatch/scheduler recordings, no world turn).
For mint allocation or character/unique-item lifecycle changes, run
`tools/code-compat/run.sh` and `tools/command-compat/run.sh`: the former checks
the allocator and shared mint draws, the latter real creation/deletion,
stack placement, ownership and serialized state (no full combat or death).

**Pin move.**
Run `tools/affected.sh OLD NEW` and, for each tool it lists, move the pin in its `run.sh` and README, then rerun it (with `--update` where the tool records C output, and review the diff of the recorded file).
Then run the routine checks.
`olympia-roundtrip` rebuilds its C cache for the new revision by itself, which takes about 3 minutes.
`turn-compat` uses `--full --record` to refresh its small C recordings and
determinism metadata; review those bytes. Its full Go turn remains an honest
failure until stage 1 integration, not an expected-success NYI baseline.

**Stage exit.**
Everything, with the whole-game comparisons:

```sh
go vet ./...
G3G0_FULL=1 go test ./...
tools/<tool>/run.sh                  # every tool pinned to the olympia revision
tools/olympia-roundtrip/run.sh --determinism-all
tools/box-compat/run.sh --full       # after olympia-roundtrip, whose games it reads
tools/boxfile-compat/run.sh --full
tools/check-compat/run.sh --full
tools/loadsave-compat/run.sh --full
tools/turn-compat/run.sh --full      # stage 1 exit; currently fails on Go's unported -r
```

Run these one at a time: `box-compat`, `boxfile-compat`, `check-compat` and `loadsave-compat` read the games `olympia-roundtrip` caches, so do not run them while `olympia-roundtrip --fresh` rebuilds that cache.

## The opt-in switches

| switch | effect |
|---|---|
| `G3G0_FULL=1` | lets the whole-game Go tests (`TestBoxWorlds`, `TestBoxfileWorlds`, `TestCheckWorlds`, `TestLoadSaveWorlds`) run; `-short` always skips them. They compare nothing on their own: they need the game list that `box-compat`, `boxfile-compat`, `check-compat` or `loadsave-compat` hands them, and skip without it. |
| `--full` | `box-compat`, `boxfile-compat`, `check-compat`, `loadsave-compat`: also compare every saved game (`world.tgz` and each game cached by `olympia-roundtrip`). |
| `--full` | `turn-compat`: run two identical pinned C world turns to establish determinism, then compare a real Go turn's state paths/bytes and raw stdout/stderr/status. Report artifacts alone are deferred. |
| `--games N` | `box-compat`, `boxfile-compat`, `check-compat`, `loadsave-compat`: compare `N` saved games, always the same ones: `world.tgz` and cached games spread evenly over the sorted list. |
| `--determinism-all` | `olympia-roundtrip`: check that C gives the same result twice on every selected case, not just the fixed sample. |

By default `box-compat`, `boxfile-compat`, `check-compat` and `loadsave-compat` run only their recorded cases, and `olympia-roundtrip` checks C's determinism only on its fixed sample of cases when it builds its cache.
The game comparisons delete each game's output as soon as it matches, and keep only failing games' output under the tool's `tmp/` directory.

## Time limits

The affected engine tools (`rng`, `entity`, `code`, `helper`, `io`, `box`, `boxfile`,
`check`, `loadsave`, `mapgen` and `olympia-cli` compat, plus `olympia-roundtrip`
and `olympia-spec`, `turn-compat`, `turn-driver-compat` and `command-compat`) run their C/Go processes under `run_limited` from
[`lib/timeout.sh`](lib/timeout.sh), with a limit of `RUN_TIMEOUT` seconds.
[`lib/run-limited.sh`](lib/run-limited.sh) adapts it for `go test -exec` and
unchanged upstream fixtures. `genesis-compat` and `island-compat` also use it.
At the limit it sends SIGTERM, then SIGKILL 2 seconds later, so a hung Go binary is stopped as surely as a C one ([#47](https://github.com/mdhender/g3g0/issues/47): the Go runtime ignores the SIGALRM the scripts used to rely on).
A run that times out has status 124, for both engines; any other status is the program's own, with 128 + N for death by signal N.
The scripts report a status of 124 as a timeout.

## Upstream files each tool reads

Paths are in the G3 repository.
A tool needs rerunning after a pin move if any file it reads changed; [`affected.sh`](affected.sh) applies this table to `git diff --name-only OLD NEW`.

| tool | records | reads |
|---|---|---|
| [`rng-compat`](rng-compat) | `internal/rng/testdata/c-output.txt` | `lib/rng.[ch]`, `lib/rnd.[ch]`, `lib/ilist.c`, `lib/lists.h`; line ranges of `olympia/` `buy.c`, `combat.c`, `npc.c`, `quest.c`, `cloud.c`, `faery.c`, `hades.c`, `tunnel.c`, `basic.c`, `c1.c`, `code.c`, `dir.c`, `day.c`, `seed.c`, `storm.c`, `swear.c`, `use.c`, `u.c` (any change can move a range) |
| [`lists-compat`](lists-compat) | `internal/lists/testdata/c-output.txt` | `lib/*.[ch]` |
| [`entity-compat`](entity-compat) | `internal/olympia/testdata/entity-c-output.txt` | `olympia/*.h`, `lib/*.h`, `olympia/glob.c`, `lib/checked_alloc.c` |
| [`code-compat`](code-compat) | `internal/olympia/testdata/code-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `code.c`, `glob.c`, `z.c`, `io.c`, `sout.c`; `lib/` `checked_alloc.c`, `ilist.c`, `rng.c`, `rnd.c` |
| [`helper-compat`](helper-compat) | `internal/olympia/testdata/helper-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `code.c`, `glob.c`, `z.c`, `sout.c`, `u.c`, `seed.c`, `loc.c`, `stack.c`, `swear.c`, `use.c`, `move.c`, `day.c`, `garr.c`, `order.c`, `input.c`; `lib/` `checked_alloc.c`, `ilist.c`, `rng.c`, `rnd.c`, `cstrings.c`, `item_ents.c`, `orders.c` |
| [`io-compat`](io-compat) | `internal/olympia/testdata/io-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `code.c`, `glob.c`, `z.c`, `io.c`, `sout.c`, `u.c`, `seed.c`; `lib/` `checked_alloc.c`, `ilist.c`, `rng.c`, `rnd.c`, `admits.c`, `item_ents.c`, `req_ents.c`, `skill_ents.c`, `trades.c` |
| [`box-compat`](box-compat) | `internal/olympia/testdata/box-c-output.txt` | what `io-compat` reads, and `olympia/input.c`, `lib/cstrings.c`, `tests/olympia/fixtures/world.tgz`; with `--full` or `--games`, the games `olympia-roundtrip` caches |
| [`boxfile-compat`](boxfile-compat) | `internal/olympia/testdata/boxfile-c-output.txt` | what `box-compat` reads, and `olympia/swear.c` |
| [`check-compat`](check-compat) | `internal/olympia/testdata/check-c-output.txt` | what `helper-compat` reads, and `olympia/check.c`, `olympia/io.c`, `lib/admits.c`, `lib/req_ents.c`, `lib/skill_ents.c`, `lib/trades.c`, `tests/olympia/fixtures/world.tgz`; with `--full` or `--games`, the games `olympia-roundtrip` caches |
| [`loadsave-compat`](loadsave-compat) | `internal/olympia/testdata/loadsave-c-output.txt` | what `check-compat` reads, and `olympia/buy.c`, `olympia/dir.c`, `olympia/add.c`, `olympia/quest.c`; with `--full` or `--games`, the games `olympia-roundtrip` caches |
| [`scanf-compat`](scanf-compat) | `internal/clib/testdata/scanf-c-output.txt` | nothing upstream: it tests the host's C library |
| [`olympia-cli-compat`](olympia-cli-compat) | `results.tsv` | the whole tree: it builds `olympia-g3` |
| [`olympia-roundtrip`](olympia-roundtrip) | nothing | the whole tree: it builds `olympia-g3` and runs the regress fixtures' saved games |
| [`olympia-spec`](olympia-spec) | `passing.txt` | the whole tree: it runs G3's spec |
| [`turn-compat`](turn-compat) | `testdata/` small real C snapshots and determinism manifest | the whole tree: it builds `olympia-g3` and turns `tests/olympia/fixtures/world.tgz` |
| [`turn-driver-compat`](turn-driver-compat) | `internal/olympia/testdata/turn-driver-c-output.txt` | `olympia/*.h`, `lib/*.h`, `olympia/u.c` |
| [`command-compat`](command-compat) | `internal/olympia/testdata/command-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `code.c`, `glob.c`, `z.c`, `sout.c`, `input.c`, `combat.c`, `order.c`, `swear.c`, `stack.c`, `loc.c`, `garr.c`, `u.c`, `perm.c`, `storm.c`, `use.c`, `c1.c`, `produce.c`, `c2.c`, `immed.c`, `stealth.c`, `move.c`, `day.c`, `seed.c`, `io.c`, `buy.c`; `lib/` `checked_alloc.c`, `ilist.c`, `rng.c`, `rnd.c`, `cstrings.c`, `item_ents.c`, `skill_ents.c`, `orders.c`, `wait_args.c`, `flag_ents.c`, `admits.c`, `trades.c`, `req_ents.c` |

`genesis-compat`, `island-compat` and `mapgen-compat` have pins of their own and are not covered here.
