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
Its calendar/schedule suites use real pinned caln/wthr streams and preserve
process-lifetime picks/list/cursor separately from the streams' turn guards.
Only test overlays trace daily gameplay callbacks: scheduling acceptance does
not accept the unfinished weather/damage/healing/faery/event implementations.
Run rng/command/move/helper/loadsave/default turn comparisons alongside it;
ordinary tests use the small direct-stream recording, not a full world.
Natural weather creation and BIND (#75) extend command-compat rather than
duplicating its real mint, placement, garrison output and save-box lifecycle
in a weather-only tool. The new cases call actual storm.c definitions, not
callback traces. Run command/move/helper/rng/code/loadsave/turn-driver/default
turn comparisons sequentially. All four scheduled weather days together
remain opt-in integration; unrelated weather magic remains explicitly
unported. Additional source dependencies cloud.c, tunnel.c and
quest.c supply only original greater-region globals (tracked by affected.sh).
Storm movement/decay (#76) extends the same command comparison with complete
day.c storm_decay/storm_move/storm_owner_touch_loc and storm.c move_storm,
move_bound_storms and dissipate_storm. Monthly relocation uses set_where
directly; ship-bound motion emits weather and touches owners. Shared day-zero
location initialization and month-end call positions are preserved. Run the
same sequential comparisons above. No new source dependency/pin or SAIL
handler, weather-magic command or full-world turn success is implied.
Faery hunts (#77) also extend command-compat, extracting complete faery.c
create_elven_hunt/warn_human/auto_faery_sup/auto_faery definitions alongside
the existing original faer prefix. This reuses #111 entry hooks, actual
creation/mint, NPC/exdr routing, queue/dispatch, logs and save-box/order writers,
not a separate modeled NPC harness. Cases cover the 15-unit cap (all faction
units), slot/retry ocean rejection, quantities/order, shared warning memory,
stone/stack filtering, and real spawned WAIT → eight-day MOVE → WAIT.
Daily call position is unchanged. Run command/move/helper/rng/code/loadsave/
turn-driver/default turn sequentially. MESSAGE/ATTACK and unrelated commands
remain explicitly unported; no new source dependency/pin or full-world success.
Structural destruction (#120), the bounded prerequisite split from #78,
extends command-compat with complete u.c add_structure_damage, sink_ship,
building_collapses, get_rid_of_collapsed_mine and kill_stack_ocean. It reuses
real death/survival/inheritance/extraction, arrival/nearest-land routing,
output/VECT and bound-storm state, not successful gameplay substitutes.
Run the same eight comparisons sequentially before vet/default tests. No new
source dependency or pin: u.c was already tracked. Daily coastal/mine/inn
events remain #78, character damage/healing/ocean polling #79, and monthly
mine decay #88. Full worlds remain opt-in; no stage-exit success is implied.
For command parsing, dispatch, handlers or completion, also run
`tools/command-compat/run.sh` (small real-C dispatch/scheduler recordings, no world turn).
MOVE/GO and all six movement synonyms use command-compat's actual C callbacks
through dispatch/evenings, including queued multi-day NPC MOVE → WAIT. Also run
move-compat for the shared travel, message and raft crew-loss contracts; full
SAIL/FLY and the eventual whole-engine raft fixture (#90) are separate work.
Stack-arrival knowledge, contacts, weather views, ocean bookkeeping and city
matching also use command-compat, reusing its actual Hades/Faery hooks, keyed
streams, entity creation, order queues and real database writer. Candidate
settlement stops explicitly at `attempt_trade` pending #81; empty and
unmatched city paths execute real `match_trades`, not a skipped callback.
For exit topology, distances, randomized exit selection, weights or travel
calculations, movement parsing/admission or movement bookkeeping, and character
descriptions/stack display (including rank/aura/equipment), departure/arrival
messages or garrison notice/spotting, also run
`tools/move-compat/run.sh` (small real-C exit-field/order/RNG/weight/travel and
parsing/admission/stack/message recordings, no world turn). Its message helpers
reuse the already listed `u.c`, `garr.c`, `move.c`, `glob.c`, `display.c` and
`sout.c` dependencies; output includes real faction logs and serialized state.
For mint allocation or character/unique-item lifecycle changes, run
`tools/code-compat/run.sh` and `tools/command-compat/run.sh`: the former checks
the allocator and shared mint draws, the latter real creation/deletion,
stack placement, ownership and serialized state (no full combat or death).
Movement prisoner escapes and nearest-land routing use command-compat's real
stack release, entity RNG, arrival, swimming, logging and serialization closure;
also run move-compat and rng-compat for their shared movement/stream contracts.
Hades population spawning and NPC route queueing also extend command-compat:
actual seeded zero/existing/capped populations, direct-target filtering,
asymmetric topology, direction persistence, shared exit shuffles and generated
MOVE → WAIT execution. Run move/helper/rng/code comparisons alongside it for
the shared topology, creation and stream closure. Autonomous NPC dispatch
also runs in command-compat; it does not substitute successful combat.

Savage fort spawning and savage autonomous order generation likewise extend
command-compat, including initial population caps, per-fort keyed checks,
recursive human filtering, ordered targets and real spawned MOVE execution.
Run command/move/helper/rng/code/loadsave/turn-driver/default turn comparisons.
Queued USE/ATTACK/RAZE/UNSTACK/DIE are not successful command execution; their
unported handlers retain explicit stops. The command tool additionally reads
`olympia/savage.c` (also tracked by `affected.sh`).

Unsworn/mob/undead autonomous programs also extend command-compat, including
the 4/5-turn mob boundary, summoner clearing, location-keyed RNG, saved-order
append and generated MOVE → WAIT. Run the same command/move/helper/rng/code/
loadsave/turn-driver/default-turn checks. ATTACK/PILLAGE/RECLAIM execution
remains explicitly unported. The command
tool additionally reads `olympia/necro.c` (also tracked by `affected.sh`).

Bandit programs and day-0 NPC dispatch extend the same command comparison:
keyed victim identity in saved target order, special-program safety filters,
all child UNSTACK orders, fallback directions and actual MOVE → WAIT;
summon/prisoner/command/queued-order skips and every program branch; real
savage-before-Hades creation and shared RNG/queue/serialization state. The
existing process_orders callback now invokes this dispatcher. Garrison guard
setup, daily default-garrison pay and token reconciliation are real; calendar
and later daily and monthly phases remain explicit stops. After #73 the
opt-in full-world diagnostic reaches `match_all_trades` and still fails;
it is not stage acceptance.

Garrison guard setup and default-garrison payroll extend command-compat, using
the real inventory/strength/maintenance tariff, subtype iteration, descriptions,
ownership and location/garrison output routes. Cases cover the 9/10 boundary,
every guard reset, nested stack leaders, captive/default/NPC skips, unequal
balances/pay, zero-pay generation/output and invalid negative-pay assertions.
Month-end garrison tax forwarding remains #116 after real maintenance #87.
Run command/move/helper/code/loadsave/turn-driver/default-turn comparisons.

Inventory transfer and token acquisition (#117) extend command-compat with
the complete real `consume_item`, mutually recursive `move_item`/`drop_item`,
and compiled `move_token` acquisition closure. Real lord interruption/NP,
creation/shared mint, inventory/lore/trade, naming and output dependencies
are reused. Run command/move/helper/rng/code/loadsave/turn-driver/default-turn
comparisons. No new upstream source dependency or pin change is needed:
`code.c` is already linked whole and `art.c`/`u.c` are already tracked.
Independent melting and `check_token_units` remain #73 after #118/#119;
the acquisition path contains no successful death/reclamation substitute.
GIVE/DROP/USE handlers and `create_npc_token` are outside this helper task.

Inheritance, prisoner release and desertion (#118) extend that same command
comparison with complete `stackmate_inheritor`, `take_unit_items`,
`put_back_cookie`, `drop_stack`, `extract_stacked_unit`, `move_prisoner`,
`unit_deserts`, `plural_item_box` and `box_name_qty`. The natural disappearance
path runs real desertion/extraction, cookie return and dead-character
conversion, without calling unfinished death/reclamation helpers. Existing
`code.c` is linked whole; `u.c`, `stack.c`, `swear.c` and `gm.c` are already
tracked dependencies (including original combat gold scalars). No pin or
dependency-table/affected.sh change is needed. Run command/move/helper/rng/
code/loadsave/turn-driver/default-turn comparisons. Full shared death remains
#119 and token reconciliation remains #73; no UNSTACK/DIE/RECLAIM handler or
whole-turn/stage-exit acceptance is implied.

Shared character death/reclamation (#119) extends command-compat with complete
`survive_fatal`, `char_reclaim`, `dead_char_body`, `kill_char`,
`hack_unique_item`, `forget_skill`, `random_hades_loc` and `char_rep_location`.
The last helper is a **new report.c dependency**: its location strings feed
the real death log, not deferred stage-2 report assembly. Original use.c
`skill_school` is also extracted, reusing its existing Go implementation.
Run command/move/helper/rng/code/loadsave/turn-driver/default-turn comparisons.
Actual survival, corpses, transcendence/arrival, token removal/destruction,
lord/NP, inheritance, extraction and interruption paths are exercised; no
RECLAIM/FORGET/DIE command, sibling damage/combat/resurrection helper or token
reconciliation is implied by #119 alone (reconciliation follows in #73).
G3 #355's missing token-vanishing notification is deliberately
reproduced and source-linked, not independently fixed.

Token melting/reconciliation (#73) extends command-compat with complete
art.c `melt_token_units` 1054–1080 and `check_token_units` 1197–1259, reusing
the full accepted death/reclamation, creation, lord/NP, inventory, output,
orders and shared RNG closure. The real process_orders definition and existing
clock/progress helpers are also extracted for a day-zero owner assertion
after actual NPC/guard setup. Only phase instrumentation is omitted; later
unported turn dependencies fail explicitly and cannot be recorded as success.
Existing processOrders/postMonth callback positions are unchanged. Run command/
move/helper/rng/code/loadsave/turn-driver/default-turn comparisons. No new
source dependency or pin change, unrelated command handlers or artifact
generation. Full-world turn success and stage exit remain unfinished.

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
tools/turn-compat/run.sh --full      # stage 1 exit; currently fails at match_all_trades
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
`move-compat` uses the same supervision for both C recordings and its Go test.
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
| [`turn-driver-compat`](turn-driver-compat) | `internal/olympia/testdata/turn-{driver,phase}-c-output.txt`, `{calendar,schedule}-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `u.c`, `input.c`, `day.c`, `z.c`, `storm.c`, `seed.c`, `glob.c`, `rnd.c`; `lib/` `ilist.c`, `rng.c` |
| [`move-compat`](move-compat) | `internal/olympia/testdata/move-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `dir.c`, `loc.c`, `u.c`, `swear.c`, `move.c`, `input.c`, `perm.c`, `stack.c`, `storm.c`, `garr.c`, `art.c`, `combat.c`, `day.c`, `code.c`, `glob.c`, `z.c`, `sout.c`, `display.c`, `use.c`, `io.c`, `seed.c`; `lib/` `checked_alloc.c`, `ilist.c`, `exit_views.c`, `rng.c`, `rnd.c`, `item_ents.c`, `admits.c`, `cstrings.c`, `skill_ents.c`, `trades.c`, `req_ents.c` |
| [`command-compat`](command-compat) | `internal/olympia/testdata/command-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `code.c`, `glob.c`, `z.c`, `sout.c`, `input.c`, `combat.c`, `art.c`, `display.c`, `order.c`, `swear.c`, `stack.c`, `loc.c`, `garr.c`, `u.c`, `dir.c`, `perm.c`, `storm.c`, `use.c`, `c1.c`, `produce.c`, `c2.c`, `immed.c`, `stealth.c`, `move.c`, `day.c`, `gm.c`, `seed.c`, `io.c`, `buy.c`, `lore.c`, `npc.c`, `hades.c`, `faery.c`, `savage.c`, `necro.c`, `report.c`, `cloud.c`, `tunnel.c`, `quest.c`; `lib/` `checked_alloc.c`, `ilist.c`, `exit_views.c`, `rng.c`, `rnd.c`, `cstrings.c`, `item_ents.c`, `skill_ents.c`, `orders.c`, `wait_args.c`, `flag_ents.c`, `admits.c`, `trades.c`, `req_ents.c` |

`genesis-compat`, `island-compat` and `mapgen-compat` have pins of their own and are not covered here.
