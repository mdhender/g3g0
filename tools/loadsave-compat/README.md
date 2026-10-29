# load_db/save_db parts 1 and 2 differential test

`run.sh` checks that the Go port of the routines `load_db()` and `save_db()` call to read or write files of their own, or to set up globals ([#50](https://github.com/mdhender/g3g0/issues/50)), and the last steps of `load_db()` before `print_box_usage()` ([#51](https://github.com/mdhender/g3g0/issues/51)), matches G3: the files they write (names, modes and contents), their stderr, and the globals they set.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `9a03d126b329ed466c4af48425553c025eaca26f`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.
The script exports that revision with `git archive`, so the clone's checkout does not matter.

For #291, all extracted definitions remain byte-identical to f73e481;
skill_school's assertion naturally moves seven source lines with the USE
guard. The deterministic **61-section / 3411-line** default recording is
unchanged. Full-corpus commands, disposable-input provenance, counts and
limits are recorded in [merged-pin proof](../README.md#merged-pin-proof-291).

| upstream | Go |
|---|---|
| `load_system`, `save_system`, `cleanup_posts` (`io.c`) | `loadSystem`, `saveSystem`, `cleanupPosts` (`io.go`) |
| `load_orders`, `load_player_orders`, `save_orders`, `save_player_orders`, `queue_order`, `p_order_head` (`order.c`) | `loadOrders`, `loadPlayerOrders`, `saveOrders`, `savePlayerOrders`, `queueOrder`, `pOrderHead` (`order.go`) |
| `init_trade_routes`, `save_trade_routes`, `apply_trade_route_seed`, `splitmix64` (`buy.c`) | `initTradeRoutes`, `saveTradeRoutes`, `applyTradeRouteSeed`, `splitmix64` (`buy.go`) |
| `determine_map_edges`, `region_row`, `region_col` (`dir.c`) | `determineMapEdges`, `regionRow`, `regionCol` (`dir.go`) |
| `rename_act_join_files` (`add.c`) | `renameActJoinFiles` (`add.go`) |
| `region`, `set_where` (`loc.c`), `loop_province` (`loop.h`) | `region`, `setWhere` (`loc.go`), `loopProvince` (`loop.go`) |
| `create_relics`, `create_a_relic` (`quest.c`) | `createRelics`, `createARelic` (`quest.go`) |
| `create_unique_item_alloc`, `add_item` (`u.c`) | `createUniqueItemAlloc`, `addItem` (`u.go`) |
| `investigate_possible_trade`, `trades_to_check` (`buy.c`), `loop_trade` (`loop.h`) | `investigatePossibleTrade`, `world.tradesToCheck` (`buy.go`), `loopTrade` (`loop.go`) |
| `init_ocean_chars`, `ocean_chars` (`move.c`) | `initOceanChars`, `world.oceanChars` (`move.go`) |
| `delete_deadchars` (`io.c`) | `deleteDeadchars` (`io.go`) |

`read_trade_route_seed()` in `init_trade_routes()` uses a dedicated unlimited
decimal reader: ASCII whitespace, optional sign, then digits. Overflow saturates
to 2^64−1 even with a minus sign; otherwise minus negates as uint64. Generic
`clib.Scanner` is unchanged. G3's `trade-route-seed-parse` fixture covers this.

## Running

```sh
tools/loadsave-compat/run.sh             # compare the recorded cases (check 1)
tools/loadsave-compat/run.sh --games 10  # and 10 saved games, always the same ones (check 2)
tools/loadsave-compat/run.sh --full      # and every saved game (check 2)
tools/loadsave-compat/run.sh --update    # rewrite internal/olympia/testdata/loadsave-c-output.txt from C
```

The prerequisites, and the way the driver is built, are those of [`check-compat`](../check-compat), whose generator `run.sh` extends: it also copies `region()` and `set_where()` from `loc.c`, the order-file routines from `order.c`, `determine_map_edges()` from `dir.c`, `rename_act_join_files()` from `add.c`, `create_relics()` from `quest.c`, `add_item()` and `create_unique_item_alloc()` from `u.c`, `init_ocean_chars()` from `move.c` and `investigate_possible_trade()` from `buy.c` into the helpers; `io.c` up to `delete_deadchars()` (so through `save_system()`), and its `delete_deadchars()`, `cleanup_posts()` and `remove_dir()`; and `buy.c`'s trade-route block, which the driver `#include`s so that it can print `buy.c`'s statics. It defines the globals of other files they use (`max_map_*`, `new_players`, `ocean_chars`, `trades_to_check`, and the special regions and players), and stubs `test_known()` and `queue_lore()`, which `add_item()` calls only for an item that has lore.

1. Check 1 runs each section of [`internal/olympia/testdata/loadsave-cases.txt`](../../internal/olympia/testdata/loadsave-cases.txt) in a fresh process and directory, with umask 022, twice, and lists the directory's tree after each. Its output must equal [`loadsave-c-output.txt`](../../internal/olympia/testdata/loadsave-c-output.txt); then `go test -run TestLoadSaveMatchesC` runs the same cases on the Go port.
2. Check 2, with `--full` or `--games N`: for `world.tgz` and every cached game input (`tmp/olympia-roundtrip/cases/*/input/lib`), from a copy of the game with `libdir` `lib`, the driver (`-g`) runs what `load_db()` runs of these routines, in its order (`init_trade_routes()`, `load_system()`, the box-file passes, `load_orders()`, `check_db()`, `determine_map_edges()`, `create_relics()`, `init_ocean_chars()`, `delete_deadchars()`), prints the globals, and runs what `save_db()` runs of them (`cleanup_posts()`, `save_trade_routes()`, `save_system()`, `write_all_boxes()`, `write_master()`, `save_orders()`, `rename_act_join_files()`). `TestLoadSaveWorlds` (run with `LOADSAVE_COMPAT_GAMES` and `G3G0_FULL=1`) does the same in Go. The `lib` trees, stdout, stderr and exit status must be identical.

Every C driver and Go test process runs under `tools/lib/timeout.sh`'s
`run_limited`, with `RUN_TIMEOUT` seconds (default 10). A timeout fails the
check; signal exits remain 134 for an assertion and 139 for SIGSEGV.

The invalid order-file names `999` and `0` are independent recorded sections
([#54](https://github.com/mdhender/g3g0/issues/54)). In the combined case,
real pinned C on Linux printed their errors in the reverse of the macOS
recording: `load_orders()` follows host `readdir()` order. Each isolated
section keeps its raw diagnostic, exit status and post-save tree comparison;
neither the engine's directory enumeration nor its output is sorted.

`load_seed()`/`save_seed()` and `print_box_usage()` are not run here, so `lib/randseed` is not written; `load_db()` and `save_db()` as a whole, with them, are checked by [`olympia-roundtrip`](../olympia-roundtrip) ([#52](https://github.com/mdhender/g3g0/issues/52)). `load_db()`'s new-world branches (between `determine_map_edges()` and `create_relics()`) are not ported yet.

None of the routines of #51 draws from an RNG: `create_a_relic()` allocates fixed box numbers (401 to 403) with `alloc_box()`, not `new_ent()`.

## What the cases cover

| sections | cover |
|---|---|
| a whole load and save | a small game with `lib/system`, orders and a Nowhere location but no relics, through load, print and save |
| load_system | every recognized line, comments, blank lines, white space, control and high-bit bytes, NUL bytes, unrecognized lines, partial and repeated `sysclock` lines, `atoi()` of junk and overflow, a value that is empty, a long line and a last line with no newline; no file, a directory, an unreadable file; `save_system()` with no directory |
| trade routes | no seed; `lib/trade_routes` (winning over `-G`, with the warning, and over `lib/trade-route-seed`); `-G`, including 0; `lib/trade-route-seed`; files that are not numbers, empty, or a directory; signs, overflow, leading zeros past scanf's 512-character limit, hexadecimal, no newline; `save_trade_routes()` that cannot write |
| load_orders, save_orders | lines with and without `:`, white space, unit 0, units with no box, file names that are not numbers or name no box, hidden files, the 250-order limit, a long line; no orders directory, orders a file, a player's file a directory or unreadable; duplicate orders and missing player sections diagnosed and skipped; dead characters skipped and stale files removed; a save that cannot write |
| cleanup_posts | posts in a location and nowhere, deleted before the boxes are written |
| determine_map_edges | provinces in ordinary and special regions; a province in no region is included even when a special region is 0 |
| create_relics | none of the relics, all three, one and a box of another kind; `nowhere_loc` 0 or naming no box (the assertion), a character, a location with items and trades, and 401 itself |
| init_ocean_chars | characters at sea, stacked at sea, on a ship at sea, on an island, on land and nowhere; run twice |
| delete_deadchars | dead characters in a location, nowhere and leading a stack; one missing from its location's here list (the assertion) |
| rename_act_join_files | a rename that works, a missing file, a target that is a non-empty directory, and one that is a file, and a player listed twice |

## Results

Final #166 run, 2026-10-28, Linux amd64, GCC 12.2.0, Go 1.26.4,
at then-pin `e733acb63bacb21167a1f2589d8dcdc023364896`: default **61 sections / 3411 lines** match
deterministic C and the unchanged recording. `--full` passes **202 games /
8,146,791 saved boxes**, with **zero check_db messages**. C/Go `lib` trees
(types/modes/bytes), raw stdout/stderr and statuses match. This count is after
the real load/save lifecycle, not the box-only 8,146,800 count.

Command: `G3_REPO="$PWD/tmp/upstream.git"
GAMES="$PWD/tmp/stage2-roundtrip/cases" tools/loadsave-compat/run.sh --full`.
The corpus is the fresh #166 roundtrip's 201 actual saved inputs plus the
tool's separate base world; no unavailable abort snapshot is included.
Evidence: `tmp/issue166-loadsave-{default,full}.log`. Shared `run_limited`
bounds all engine/test processes; matching per-game copies are removed.

Historical foundation evidence:

2026-10-27, Linux amd64 at the current pin: all 61 sections (3411 lines)
match deterministic C byte for byte after C `--update`. The trade-route
numeric cases now follow G3 #341; no generic Scanner behavior changed.
The whole-engine `world` load/save roundtrip builds and passes too (G3 #336
is fixed upstream). G3's new `trade-route-seed-parse` regression passes on
both actual engines, including its 14 numeric cases.

Previous run, before the split: 2026-10-27, macOS arm64: check 1, 59 sections, matching; check 2 (`--full`), 197 saved games, C and Go identical.
Of those games, every one already has the three relics, so `create_relics()` makes none; one has characters at sea; seven have dead characters for `delete_deadchars()`.
`add_item()`'s branch for an item already held and `investigate_possible_trade()`'s match cannot be reached from a loaded game: the box files' inventories and trades keep only items that exist, and `create_a_relic()` makes only items that do not.
