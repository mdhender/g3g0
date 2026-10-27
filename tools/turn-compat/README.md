# Turn comparison (stage 1, #56)

Pinned reference: [playbymail/olympia-g3 at
2567e6a5dd0ba0d95516d8bed9e1bafb87a51985](https://github.com/playbymail/olympia-g3/tree/2567e6a5dd0ba0d95516d8bed9e1bafb87a51985).
No turn engine code is implemented here. The reference is real C, never Go.

```sh
tools/turn-compat/run.sh                 # fast recorded checks; no world turn or C build
tools/turn-compat/run.sh --full          # real C/C/Go world turns; currently exits 1
tools/turn-compat/run.sh --full --record # explicitly refresh C-only recordings, then compare Go
```

## Default checks

The Go tests exercise the comparator's observable contract on a tiny recorded
slice of the real C world turn: `system`, `randseed`, the empty `ship`, raw
stdout, raw stderr and status. Negative controls check changed bytes, missing
files (including empty ones), extra files, orders, factions, trade routes,
unknown paths, output whitespace, exit status, and the precise report exclusions.
These mutations test the comparator, not expected engine behavior.

The harness also runs the real Go engine with the turn arguments below plus
`-?`, comparing its rejected-option diagnostic and unchanged input file with
two recorded C runs. It never treats that diagnostic as a successful turn.
`go test ./...` needs no C checkout and never runs a whole-world turn here.

`testdata/c-determinism.txt` records the C pin, binary and world archive SHA-256,
arguments, successful repeat check, and a state/output digest manifest. Tests
check each small snapshot against that manifest. This is recorded provenance,
not a claim to have rerun C during default tests. The recordings total about
20 KB; `--record` writes them only after real pinned C repeats match, before
looking at Go. It does not accept Go's current NYI as an expectation.

## Full comparison

Both engines get identical copies of the pinned `world.tgz` and run:

```sh
olympia-g3 -r -l ./lib -S test-use-const-report-date
```

`argv[0]` is `olympia-g3`, stdin is `/dev/null`. Every engine invocation uses
`run_limited` from `tools/lib/timeout.sh` (`RUN_TIMEOUT`, default 120 seconds).
Test executables also use its `run-limited.sh` adapter through `bash`.
Timeout is status 124; signal death remains 128 + signal, without shell noise.
No stdout/stderr filtering, NYI special case, or status normalization occurs.

The C turn runs twice on every `--full` invocation. Both statuses must be zero
and the *entire* captures must match, including deferred report files. A failed
or nondeterministic C reference stops the harness with exit 2 and retained
evidence, rather than blessing a Go comparison. Only then is Go compared to C.

### Database ownership and report boundary

At the pin, `io.c:save_db` calls `save_seed`, `save_trade_routes`, `save_system`,
`write_all_boxes`, `write_master`, `save_orders` and `rename_act_join_files`.
These own `randseed`, optional `trade_routes`, `system`, `loc`, `item`, `skill`,
`gate`, `road`, `ship`, `unform`, `misc`, `fact/*`, `master`, `orders/*` and
act/join state. Comparing only an allowlist would miss new state or extra files.

Instead, the comparator compares **every file in the working directory**, with
only these exact stage-2 output artifacts excluded:

* `lib/log/*`, `lib/save/*`: `sout.c` and `io.c:save_logdir`.
* `lib/players`, `lib/email`, `lib/forward`, `lib/factions`, `lib/totimes`:
  the writers in `main.c`.
* `lib/times_0`, `lib/times_press`, `lib/times_rumor`: `c2.c`.

Report routines still run in C. Any changes they make to database state
(for example known/lore state in faction records) remain compared. Lore input,
`startloc`, unknown files and files outside `lib` remain compared too. No broad
`times_*` or faction filename filtering is used. Missing and extra files fail;
empty files count. Nonregular state files are harness errors, not followed
symlinks. Empty directories and file permissions are not serialized DB bytes.

## Cache, prerequisites and evidence

Requires Go, Bash, Perl, `cmp`/`diff`, `tar`, `shasum`; full mode also needs a
clone containing the pin and CMake/Ninja/C11 when no reusable C build exists.

* `WORK`: defaults to `tmp/turn-compat`; all builds and captures live here.
* `G3_REPO`: defaults to `~/Software/playbymail/olympia-g3`; only `git archive`
  reads the pinned commit. A bare cache such as `tmp/upstream.git` works.
* `C_CACHE`: defaults to `tmp/olympia-roundtrip`. An exported cache with an
  exact `g3.rev` and `g3/build/debug/olympia-g3` is reused read-only. Otherwise
  the script exports the exact commit into `$WORK/reference/g3` and builds
  with upstream's unchanged debug preset. It never builds in the source clone.
* Input is cached by pin and archive digest. C is repeated each full run,
  so an old deterministic marker cannot hide changed C behavior.

Example using the accepted existing export without rebuilding:

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue97-roundtrip" \
  tools/turn-compat/run.sh --full
```

On mismatch, exit 1; `$WORK/{c,go}/{cwd,stdout,stderr,status}` retain raw
evidence and `diff.txt` lists all state path/byte differences and unfiltered
stdout/stderr diffs. `c-determinism.txt` retains repeat metadata. The redundant
matching C repeat is deleted immediately; matching Go/C game outputs are
deleted after success. The input and C build remain reusable.

Initial Linux amd64 recording: C statuses 0/0, identical stdout, stderr and
entire post-turn directory; Go status 2, `olympia-g3: load_db (-r -S): not yet
implemented`. The full check fails on status, stdout, stderr, 15 changed state
files and missing `orders/204`. That is the required red integration gate,
not a newly discovered bug and not a passing turn claim.
