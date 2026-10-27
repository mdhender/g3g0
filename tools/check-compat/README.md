# check_db differential test

`run.sh` checks that the Go port of `check_db()` (`internal/olympia`: `check.go`) matches G3's `olympia/check.c`: every check, in C's order, with its repairs, its stderr messages and its assertions.
It shows the repairs the way the engine would: it loads a database with the box-file passes, runs `check_db()` on it, and writes it back out.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `47535543d83fbc5649d54b63b49e578f2b11ab1c`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.
The script exports that revision with `git archive`, so the clone's checkout does not matter.

The port follows `olympia/check.c` (all of it), in `check.go`:

| upstream | Go |
|---|---|
| `check_db` | `checkDb` |
| `check_here`, `check_swear` | `checkHere`, `checkSwear` |
| `check_indep`, `check_gm`, `check_skill_player`, `check_eat_player`, `check_npc_player`, `check_garr_player` | `checkIndep`, `checkGm`, `checkSkillPlayer`, `checkEatPlayer`, `checkNpcPlayer`, `checkGarrPlayer`, which share `checkSpecialPlayer` |
| `check_glob`, `check_nowhere`, `check_skills`, `check_item_counts`, `check_loc_name_lengths`, `check_moving`, `check_prisoner` | `checkGlob`, `checkNowhere`, `checkSkills`, `checkItemCounts`, `checkLocNameLengths`, `checkMoving`, `checkPrisoner` |

The helpers it calls were ported in [#48](https://github.com/mdhender/g3g0/issues/48), and are checked by [`helper-compat`](../helper-compat).
`loop_here` and `loop_inv` assert `valid_box()` of their box at the caller's line (`check.c` lines 49 and 400); both run inside `loop_boxes`, which yields only valid boxes, and nothing in those checks deletes a box, so the assertions cannot fail and the port leaves them out.

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
| check_here | units, a stacked unit, a location and a ship missing from their location's here list (added at the end); here lists naming boxes that are elsewhere or nowhere (removed), a box that lists itself, duplicates of a box that belongs; a stale box listed twice, which `remove_from_here_list()` removes both copies of at once, so the second visit fails its assertion (G3 [#331](https://github.com/playbymail/olympia-g3/issues/331)) |
| check_swear | units missing from their player's list, directly and through a lord (added, and the list sorted), a unit of another player; lists naming another player's units, a location, and duplicates (each reported and removed, a duplicate reported twice); a lord loop (`player()`'s assertion) |
| check_indep | units with no lord, sworn to a unit with no player, with no CH, sworn to a location, with a command that is loading and one that is running (each sworn to 100: `prev_lord` and loyalty cleared, CH and command created); the independent player created first; a unit whose command is in error (`init_load_sup()`'s assertion); a unit that is nowhere, and one in a region (`loc_depth()`'s assertion in `touch_loc()`, G3 [#330](https://github.com/playbymail/olympia-g3/issues/330)) |
| check_nowhere | a unit, locations and ships below a region that are nowhere (ships first), a region that is not reported; a raft (`loc_depth()`'s assertion) |
| check_skills | an orphaned subskill, skills with no learn time and no SK, a skill two schools offer, offered and researched skills in another school, an offered item, a skill no one offers; a loop of required skills (`skill_school()`'s assertion) |
| check_item_counts | a unit in an inventory, a unique item held by the wrong unit and by a location (who_has repaired to each in turn), a quantity of 3, 0 and -1 of a unique item, a unique item no one holds, a unique item whose who_has is an item |
| check_loc_name_lengths | names of 25 and 26 bytes, with high-bit bytes, an unnamed location, a ship (not checked) |
| check_moving | a moving stack whose leader has no command (stopped, with every unit under it), a command that is loading (stopped) and running (kept), followers that disagree with their leader both ways, one with no CH until `check_indep()` gives it one; a stack loop (`stack_leader()`'s assertion) |
| check_prisoner | prisoners in a location (freed), stacked (kept) and nowhere (freed) |
| garrison_magic allocated | the warning |

## Results

Last run: 2026-10-27, macOS 26.6.2 arm64, Apple clang 21.0.0, Go 1.26.4.

- Check 1: 30 sections, 2071 lines. The C driver is deterministic, its output matches `check-c-output.txt`, and `TestCheckMatchesC` passes.
- Check 2 (`--full`): 197 games (`world.tgz` and the 196 cached round-trip inputs; 7,947,575 boxes in their written `master` files, 36,999 in `world.tgz`): C and Go `lib` trees (files, modes and contents), stdout, stderr and exit status identical. No game gives any message but `check_db()`'s stage line: no repair branch, warning or assertion is reached by a real saved game.

## Upstream defects reproduced on purpose

- `check_indep()` swears a unit with no player and no location to the independent player, and `set_lord()`'s `touch_loc()` then fails `loc_depth()`'s assertion, so `check_db()` aborts instead of repairing the database: G3 [#330](https://github.com/playbymail/olympia-g3/issues/330) (reproduced in `touchLoc`, `day.go`). A unit directly in a region fails the same way, since `viewloc()` steps out of the region to 0.
- `check_here()` walks a copy of each here list and removes a stale entry with `remove_from_here_list()`, which removes every copy of it. A stale box listed twice is reached again, and `remove_from_here_list()`'s assertion that it is in the list fails, so `check_db()` aborts (the section "check_here removes from here lists"): G3 [#331](https://github.com/playbymail/olympia-g3/issues/331).
