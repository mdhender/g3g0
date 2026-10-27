# load_db/save_db part 1 differential test

`run.sh` checks that the Go port of the routines `load_db()` and `save_db()` call to read or write files of their own, or to set up globals ([#50](https://github.com/mdhender/g3g0/issues/50)), matches G3: the files they write (names, modes and contents), their stderr, and the globals they set.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `47535543d83fbc5649d54b63b49e578f2b11ab1c`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.
The script exports that revision with `git archive`, so the clone's checkout does not matter.

| upstream | Go |
|---|---|
| `load_system`, `save_system`, `cleanup_posts` (`io.c`) | `loadSystem`, `saveSystem`, `cleanupPosts` (`io.go`) |
| `load_orders`, `load_player_orders`, `save_orders`, `save_player_orders`, `queue_order`, `p_order_head` (`order.c`) | `loadOrders`, `loadPlayerOrders`, `saveOrders`, `savePlayerOrders`, `queueOrder`, `pOrderHead` (`order.go`) |
| `init_trade_routes`, `save_trade_routes`, `apply_trade_route_seed`, `splitmix64` (`buy.c`) | `initTradeRoutes`, `saveTradeRoutes`, `applyTradeRouteSeed`, `splitmix64` (`buy.go`) |
| `determine_map_edges`, `region_row`, `region_col` (`dir.c`) | `determineMapEdges`, `regionRow`, `regionCol` (`dir.go`) |
| `rename_act_join_files` (`add.c`) | `renameActJoinFiles` (`add.go`) |
| `region`, `set_where` (`loc.c`), `loop_province` (`loop.h`) | `region`, `setWhere` (`loc.go`), `loopProvince` (`loop.go`) |

`fscanf(fp, "%llu", ...)` in `init_trade_routes()` is `clib.Scanner.Uint64`, which reads as macOS's scanf does (at most 512 characters, sign and leading zeros included).

## Running

```sh
tools/loadsave-compat/run.sh             # compare the recorded cases (check 1)
tools/loadsave-compat/run.sh --games 10  # and 10 saved games, always the same ones (check 2)
tools/loadsave-compat/run.sh --full      # and every saved game (check 2)
tools/loadsave-compat/run.sh --update    # rewrite internal/olympia/testdata/loadsave-c-output.txt from C
```

The prerequisites, and the way the driver is built, are those of [`check-compat`](../check-compat), whose generator `run.sh` extends: it also copies `region()` and `set_where()` from `loc.c`, the order-file routines from `order.c`, `determine_map_edges()` from `dir.c` and `rename_act_join_files()` from `add.c` into the helpers; `io.c` up to `delete_deadchars()` (so through `save_system()`), and its `cleanup_posts()` and `remove_dir()`; and `buy.c`'s trade-route block, which the driver `#include`s so that it can print `buy.c`'s statics. It defines the globals of other files they use (`max_map_*`, `new_players`, and the special regions and players).

1. Check 1 runs each section of [`internal/olympia/testdata/loadsave-cases.txt`](../../internal/olympia/testdata/loadsave-cases.txt) in a fresh process and directory, with umask 022, twice, and lists the directory's tree after each. Its output must equal [`loadsave-c-output.txt`](../../internal/olympia/testdata/loadsave-c-output.txt); then `go test -run TestLoadSaveMatchesC` runs the same cases on the Go port.
2. Check 2, with `--full` or `--games N`: for `world.tgz` and every cached game input (`tmp/olympia-roundtrip/cases/*/input/lib`), from a copy of the game with `libdir` `lib`, the driver (`-g`) runs what `load_db()` runs of these routines, in its order (`init_trade_routes()`, `load_system()`, the box-file passes, `load_orders()`, `check_db()`, `determine_map_edges()`), prints the globals, and runs what `save_db()` runs of them (`cleanup_posts()`, `save_trade_routes()`, `save_system()`, `write_all_boxes()`, `write_master()`, `save_orders()`, `rename_act_join_files()`). `TestLoadSaveWorlds` (run with `LOADSAVE_COMPAT_GAMES` and `G3G0_FULL=1`) does the same in Go. The `lib` trees, stdout, stderr and exit status must be identical.

`load_seed()`/`save_seed()`, `load_db()`'s new-world branches and the steps after `determine_map_edges()` are left to [#51](https://github.com/mdhender/g3g0/issues/51) and [#52](https://github.com/mdhender/g3g0/issues/52), so `lib/randseed` is not written.

## What the cases cover

| sections | cover |
|---|---|
| a whole load and save | a small game with `lib/system` and orders, through load, print and save |
| load_system | every recognized line, comments, blank lines, white space, control and high-bit bytes, NUL bytes, unrecognized lines, partial and repeated `sysclock` lines, `atoi()` of junk and overflow, a value that is empty, a long line and a last line with no newline; no file, a directory, an unreadable file; `save_system()` with no directory |
| trade routes | no seed; `lib/trade_routes` (winning over `-G`, with the warning, and over `lib/trade-route-seed`); `-G`, including 0; `lib/trade-route-seed`; files that are not numbers, empty, or a directory; signs, overflow, leading zeros past scanf's 512-character limit, hexadecimal, no newline; `save_trade_routes()` that cannot write |
| load_orders, save_orders | lines with and without `:`, white space, unit 0, units with no box, file names that are not numbers or name no box, hidden files, the 250-order limit, a long line; no orders directory, orders a file, a player's file a directory or unreadable; one player named twice (the assertion); a player with no `PL` and a box that is not a player (NULL dereference, exit 139); dead characters skipped and stale files removed; a save that cannot write |
| cleanup_posts | posts in a location and nowhere, deleted before the boxes are written |
| determine_map_edges | provinces in ordinary and special regions; a province in no region, skipped while any special region is 0 |
| rename_act_join_files | a rename that works, a missing file, a target that is a non-empty directory, and one that is a file, and a player listed twice |

## Results

Last run: 2026-10-27, macOS arm64: check 1, 48 sections, matching; check 2 (`--full`), 197 saved games, C and Go identical.
