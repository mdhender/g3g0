# check_db differential test

`run.sh` checks that the Go port of `check_db()` (`internal/olympia`: `check.go`) matches G3's `olympia/check.c`: every check, in C's order, with its repairs, its stderr messages and its assertions.
It shows the repairs the way the engine would: it loads a database with the box-file passes, runs `check_db()` on it, and writes it back out.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `9a03d126b329ed466c4af48425553c025eaca26f`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.
The script exports that revision with `git archive`, so the clone's checkout does not matter.

## Shared USE assertion provenance (#290 amended scope)

At **0.136.2-alpha**, the parent's narrow #290 amendment advances check from
**f73e48118805ffe8ff5e763bfa4f20ae00c3a705** to the immutable merged pin above,
alongside command/helper. Command's new C oracle requires the shared Go
`skill_school` assertion at **use.c:1189**, incompatible with check's old-pin
**1182**. Parent independently reproduced this source-line-only conflict and
authorized these two extra pin/recording refreshes in #290; #291 retains the
remaining eight affected pins and full-corpus proof. No check_db gameplay,
extractor, assertion remapping or test weakening is added here.

Complete relevant old/new source diff contains only the seven-line USE guard.
All **46 extracted helper definitions** are body-identical at old/new pins and
verbatim in the new generated source under original `#line` directives.
`skill_school` is **157 bytes**, SHA-256
`3f855279d65348a713e156f3e9621bed2c62033d59ef7962f3fe42c207465f67`;
its return-type line moves **1174→1181**, assertion **1182→1189**. The entire
check.c, io.c and world.tgz are unchanged from the old pin, hashes respectively
`e60cef6a324e0b2e4708b8b78bb3769243295e380e4ec04f42bfda67f9143f27`,
`8ce3a1ccbb044a973004f87e48ae51ed82590a5dd614c96237d4424652f7824c`,
`5c59fc014ea3d2f9a14e4756bfa7ccc8c7c11b23bf6611320713ec19760d8ce9`.
All **101** pristine archive files match the merged pin. Generated helpers.c:
**17233 bytes**, SHA-256
`e6a4c5187a06a6eefdb1a267c371aeedfb7a8a51f4e69acfc34a565dc08a3f0a`.
Preflight check-driver: **206936 bytes**, SHA-256
`f357580f82202f744f09ce7f6fa4ba2c497ff38db713461b8d910add8720d58a`.

Default preflight runs unchanged real C twice before replacement/Go. The
separate byte audit rejects dependency/NYI/stub/sanitizer markers before
accepting the results; raw statuses remain **21 × 0 / 10 × 134**. The reviewed
diff is **only output line 1634's assertion 1182→1189**; paths/types/modes,
serialized bytes and every other output byte are unchanged. Normal --update
independently repeats unchanged C twice before replacement/Go. Both pairs
and the recording match: **31 sections / 2347 lines / 47399 bytes**, SHA-256
`dda3edf28f247c271a430dc258d3059a9e93bede685fdce054144760059b62bc`.
Old output SHA-256:
`1888acfb30da506493e9cecd896241ab3cc133c832317c17b4e770aa6f30b43c`.
Input remains unchanged, SHA-256
`5e256d6dd7e158e78f5971f274d51144e639b1b6e8488b30565f097f2ff78423`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue290-check-preflight" \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/check-compat/run.sh # expected 1: exact reviewed old-pin line diff
# Audit both C outputs and source extraction before replacement.
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue290-check-update" \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/check-compat/run.sh --update
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue290-check-default" \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/check-compat/run.sh # fresh C/C/Go match, no recording update
```

Every C/Go engine/test is bounded by shared run_limited (**10s** per process).
Inputs/source are disposable physical git archives from isolated bare
`tmp/upstream.git`, never shared upstream checkouts or protected caches.
Evidence: `tmp/issue290-{check-preflight,check-update,check-default,amend-source-proof}.log`;
post-amendment routine gate results and exact commands are in
[command-compat's #290 record](../command-compat/README.md#missing-use-argument-290-follow-c).
Matching captures, generated source and task binaries are removed after proof.
**No new-pin full-corpus result is claimed here**; that proof remains #291.

The port follows `olympia/check.c` (all of it), in `check.go`:

| upstream | Go |
|---|---|
| `check_db` | `checkDb` |
| `check_here`, `check_swear` | `checkHere`, `checkSwear` |
| `check_indep`, `check_gm`, `check_skill_player`, `check_eat_player`, `check_npc_player`, `check_garr_player` | `checkIndep`, `checkGm`, `checkSkillPlayer`, `checkEatPlayer`, `checkNpcPlayer`, `checkGarrPlayer`, which share `checkSpecialPlayer` |
| `check_glob`, `check_nowhere`, `check_skills`, `check_item_counts`, `check_loc_name_lengths`, `check_moving`, `check_prisoner` | `checkGlob`, `checkNowhere`, `checkSkills`, `checkItemCounts`, `checkLocNameLengths`, `checkMoving`, `checkPrisoner` |

The helpers it calls were ported in [#48](https://github.com/mdhender/g3g0/issues/48), and are checked by [`helper-compat`](../helper-compat).
`loop_here` and `loop_inv` assert `valid_box()` of their box at the caller's line (`check.c` lines 49 and 414); both run inside `loop_boxes`, which yields only valid boxes, and nothing in those checks deletes a box, so the assertions cannot fail and the port leaves them out.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another.
- A C11 compiler (`CC`, default `cc`), Go, Bash, Python 3, `git`, `tar`, `cmp`, `diff`, `find`, `stat` and `sort`.
- Optional: the inputs `tools/olympia-roundtrip/run.sh` caches under `tmp/olympia-roundtrip/cases` (`GAMES`). Without them only `world.tgz` is used for check 2.

## Running

```sh
tools/check-compat/run.sh             # compare the recorded cases (check 1)
tools/check-compat/run.sh --games 10  # and 10 saved games, always the same ones (check 2)
tools/check-compat/run.sh --full      # and every saved game (check 2)
tools/check-compat/run.sh --update    # rewrite internal/olympia/testdata/check-c-output.txt from C
```

Check 2 runs only with `--full` or `--games N`; `N` games are `world.tgz` and cached games spread evenly over the sorted list.
[`tools/README.md`](../README.md) says when to run which.

The script:

1. Exports `lib/`, `olympia/` and `tests/olympia/fixtures/world.tgz` at the pinned revision with `git archive` into `tmp/check-compat/g3`.
2. Copies out of the export, verbatim, the helpers `check_db()` calls, as `helper-compat` does (each under a `#line` directive that names its own file and line, so that `asfail()` prints them as the engine's build does), and `lookup()` and `stage()` from `olympia/u.c`; and `olympia/io.c` from its first line to just before `load_system()`, and its `remove_dir()`, as `boxfile-compat` does. It generates stubs for the command routines `cmd_tbl` names, and for the routines the helpers call only on paths not ported (`set_lord()` moving a unit between players, and `load_command()` of an order that does not parse), each of which prints `stub NAME` and aborts.
3. Builds a C driver (`tmp/check-compat/driver.c`) that `#include`s that copy of io.c, linked with the helpers and the exported `olympia/check.c`, `code.c`, `glob.c`, `z.c` and the `lib/` modules.
4. Check 1: runs each section of [`internal/olympia/testdata/check-cases.txt`](../../internal/olympia/testdata/check-cases.txt) in a fresh process and directory, with umask 022, twice, and lists the directory's tree after each (every file's mode and contents). Each section writes a small database, loads it, runs `check_db()` and writes it out again. It requires the same output both times and requires it to equal [`check-c-output.txt`](../../internal/olympia/testdata/check-c-output.txt). Then it runs `go test -run TestCheckMatchesC`, which runs the same cases on the Go port.
5. Check 2, with `--full` or `--games N`: for `world.tgz` and every cached game input (`tmp/olympia-roundtrip/cases/*/input/lib`), the driver (`-g`) runs, from a copy of the game with `libdir` `lib`: `fast_scan()` or else `scan_all_boxes()`, `read_all_boxes()`, `check_db()`, `write_all_boxes()` and `write_master()`. `TestCheckWorlds` (run with `CHECK_COMPAT_GAMES`) does the same in Go on its own copy. The two `lib` trees must be identical (`diff -r`, and every path's mode), and so must stdout, stderr and the exit status. It then lists every message `check_db()` gave over all the games, with the number of games that gave it.

`go test ./...` runs `TestCheckMatchesC` too, so the port is held to the recorded C output without a G3 clone. `TestCheckWorlds` is skipped unless `G3G0_FULL=1` is set and `CHECK_COMPAT_GAMES` names a list of games, as `run.sh` does. `-short` always skips it.

### What the driver leaves out of load_db()

In the engine, `load_db()` runs more between the box-file passes and `check_db()`, and after it: `load_seed()`, `init_trade_routes()` and `load_system()` before the passes, and `load_orders()` between `read_all_boxes()` and `check_db()`; then `determine_map_edges()` and the rest. The driver runs only the passes, `check_db()` and the writers. Two things follow:

- With no orders loaded, `check_indep()`'s `set_lord()` loads no queued order into a unit it swears to the independent player. In the engine it would, and an order that does not parse goes to `do_command()`, which the port does not have yet ([#28](https://github.com/mdhender/g3g0/issues/28) wires `check_db()` into `load_db()`).
- What `check_db()` changes that the box files do not keep (the load and run queues, players' `locs`, the `temp` fields) is not compared on real games. The writers zero every `temp` before they use it.

## What the cases cover

| section | covers |
|---|---|
| a consistent game | every special player, a player, a noble, a school and a unique item where they belong: only the stage line |
| the special players | all missing (created, with their subkinds, and named), present with no names (named), and each one not a player (the assertion of each check, in order) |
| check_here | units, a stacked unit, a location and a ship missing from their location's here list (added at the end); here lists naming boxes that are elsewhere or nowhere (removed), a box that lists itself, duplicates of a box that belongs; stale duplicates removed and reported once (G3 #331) |
| check_swear | units missing from their player's list, directly and through a lord (added, and the list sorted), a unit of another player; lists naming another player's units, a location, and duplicates (each reported and removed, a duplicate reported twice); a lord loop (`player()`'s assertion) |
| check_indep | units with no lord, sworn to a unit with no player, with no CH, sworn to a location, with a command that is loading and one that is running (each sworn to 100: `prev_lord` and loyalty cleared, CH and command created); the independent player created first; a unit whose command is in error (`init_load_sup()`'s assertion); orphan units nowhere or directly in a region repaired without touching location 0 (G3 #330) |
| check_nowhere | a unit, locations and ships below a region that are nowhere (ships first), a region that is not reported; completed and unfinished rafts load and save; invalid nowhere_loc (zero and an absent box) only warns; a valid destination relocates orphan units in character iteration order, appends after a resident and preserves located units; a repeated check does not move them again. The new relocation case captures raw stderr separately in `diagnostics`, whose exact bytes and saved locations/here lists are recorded. |
| check_skills | an orphaned subskill, skills with no learn time and no SK, a skill two schools offer, offered and researched skills in another school, an offered item, a skill no one offers; a loop of required skills (`skill_school()`'s assertion) |
| check_item_counts | a unit in an inventory, a unique item held by the wrong unit and by a location (who_has repaired to each in turn), a quantity of 3, 0 and -1 of a unique item, a unique item no one holds, a unique item whose who_has is an item |
| check_loc_name_lengths | names of 25 and 26 bytes, with high-bit bytes, an unnamed location, a ship (not checked) |
| check_moving | a moving stack whose leader has no command (stopped, with every unit under it), a command that is loading (stopped) and running (kept), followers that disagree with their leader both ways, one with no CH until `check_indep()` gives it one; a stack loop (`stack_leader()`'s assertion) |
| check_prisoner | prisoners in a location (freed), stacked (kept) and nowhere (freed) |
| garrison_magic allocated | the warning |

## Results

Final #166 run: 2026-10-28, Linux amd64, GCC 12.2.0, Go 1.26.4,
at then-pin `e733acb63bacb21167a1f2589d8dcdc023364896`.

- Check 1: **31 sections / 2347 lines**, deterministic C, unchanged
  recording and matching Go.
- Check 2 (`--full`): **202 games / 8,146,800 boxes**, including the
  separate 36,999-box base world and all 201 pristine cached roundtrip inputs.
  C/Go `lib` trees (types/modes/bytes), raw stdout/stderr and status match.
  **Zero check_db messages** beyond its stage line; the corpus does not
  exercise repair branches merely by loading these consistent worlds.

Command: `G3_REPO="$PWD/tmp/upstream.git"
GAMES="$PWD/tmp/stage2-roundtrip/cases" tools/check-compat/run.sh --full`.
Evidence: `tmp/issue166-check-{default,full}.log`. All engine/test processes
use shared `run_limited`; matching per-game captures are deleted promptly.

## Upstream fixes followed

G3 #340 fixes orphan-location touches (#330) and duplicate stale here entries
(#331). The recorded cases follow those repairs at the pinned revision.
G3 #351 adds nowhere relocation (#347) and raft depths (#348). The C driver
extracts the real `set_where` dependency; `nowhere N` only sets the global,
and `stderr NAME` redirects the real stderr to the recorded file.
