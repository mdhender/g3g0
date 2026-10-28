# check_db helpers differential test

`run.sh` checks that the Go port of the helpers `check_db()` (`olympia/check.c`) calls matches G3's C code.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `e733acb63bacb21167a1f2589d8dcdc023364896`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.

At #131's pin move, two identical C runs retain all 37 sections / 1487 lines;
the only recording delta is is_unit's swear.c assertion offset 297→294.

The port follows:

| upstream | functions | Go |
|---|---|---|
| `olympia/code.c` | `display_name`, `box_name`, `just_name`, `plural_item_name` | `code.go` |
| `olympia/u.c` | `loc_depth`, `clear_temps`, `cap` | `u.go` (`cap` is `capitalize`) |
| `olympia/loc.c` | `add_char_here`, `all_stack`, `subloc`, `viewloc`, `add_to_here_list`, `remove_from_here_list`, `in_here_list`, `first_character` | `loc.go` |
| `olympia/stack.c` | `stack_parent`, `stack_leader` | `stack.go` |
| `olympia/swear.c` | `set_loyal`, `set_lord` (for a unit with no player), `is_unit` | `swear.go` |
| `olympia/use.c` | `skill_school` | `use.go` |
| `olympia/move.c` | `restore_stack_actions` | `move.go` |
| `olympia/day.c` | `touch_loc`, `touch_loc_pl` | `day.go` |
| `olympia/garr.c` | `garrison_here` | `garr.go` |
| `olympia/order.c` | `rp_order_head`, `top_order`, `pop_order` | `order.go` |
| `olympia/input.c` | `parse_arg`, `oly_parse`, `get_command`, `set_state`, `load_command` (for an order that parses), `init_load_sup` | `input.go` |
| `olympia/seed.c`, `olympia/z.c` | `int_comp`, `safe_qsort` | `seed.go`, `z.go` |
| `olympia/oly.h` | `noble_item`, `char_moving`, `is_prisoner`, `item_unique`, `learn_time`, `req_skill` | `world.go` |
| `olympia/loop.h` | `loop_here`, `loop_stack`, `loop_inv` (`loop_char` and the other `loop_kind` loops, and `loop_units`, were already ported) | `loop.go` |

Left out, because `check_db()` does not reach them:

- `display_kind` and the `is_port_city` it calls, `nice_num`, `just_name_qty`, `plural_item_box`, `box_name_qty` and `box_name_kind`.
- In `set_lord`, moving a unit from one player to another (`flush_unit_orders`, `interrupt_order`, `clear_all_att`, `reswear_all_sworn`, the pledge, and the noble point exchange). `check_db()` calls `set_lord` only for units with no player. The port panics if it is reached.
- In `load_command`, an order that does not parse, which goes to `do_command()`. The port panics if it is reached.

The C driver stubs those routines with a stub that prints `stub NAME` and aborts, so a case that reaches one shows.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another.
- A C11 compiler (`CC`, default `cc`), Go, Bash, Python 3, `git`, `cmp` and `diff`.

## Running

```sh
tools/helper-compat/run.sh            # compare
tools/helper-compat/run.sh --update   # rewrite internal/olympia/testdata/helper-c-output.txt from C
```

The script:

1. Exports `lib/` and `olympia/` at the pinned revision with `git archive` into `tmp/helper-compat/g3`.
2. Copies the helpers above, and `sout()`, `min()`, `test_bit()`, `set_bit()`, `lord()`, `player()`, `parse_line()`, `remove_comment()`, `remove_ctrl_chars()`, `find_command()`, `oly_parse_cmd()` and `p_order_head()` that they call, out of the export into one file, each under a `#line` directive that names its own file and line, so that `asfail()` prints them as the engine's build does. It generates stubs for the command routines `cmd_tbl` names and for the routines left out above.
3. Writes a C driver to `tmp/helper-compat/driver.c` and links it with the exported `olympia/code.c`, `olympia/glob.c`, `olympia/z.c`, `lib/checked_alloc.c`, `lib/ilist.c`, `lib/rng.c`, `lib/rnd.c`, `lib/cstrings.c`, `lib/item_ents.c` and `lib/orders.c`.
4. Runs each section of [`internal/olympia/testdata/helper-cases.txt`](../../internal/olympia/testdata/helper-cases.txt) in a fresh process, since a section may end in a failed assertion, which aborts. stderr goes to stdout, in order, and `exit STATUS` follows each section. It does this twice and requires the same output both times.
5. Requires that output to equal [`internal/olympia/testdata/helper-c-output.txt`](../../internal/olympia/testdata/helper-c-output.txt), or with `--update`, writes it there.
6. Runs `go test -run TestHelperMatchesC ./internal/olympia`, which runs the same cases on the Go port and requires the recorded C output line for line.

`go test ./...` runs `TestHelperMatchesC` too, so the port is held to the recorded C output without a G3 clone.

## What the cases cover

| section | covers |
|---|---|
| names | `display_name`, `box_name` and `just_name` for every kind and for subkinds of each, named and unnamed boxes, brackets in names, a name longer than `sout()`'s buffer, high-bit bytes, `garrison_magic`, deleted, missing and out-of-range boxes, and nobles made of an item (named, unnamed, missing, deleted, a noble made of a noble) |
| cap, plural_item_name | `cap` on empty, lower-case, upper-case, punctuation, high-bit and long strings; `plural_item_name` for qty 1 and others, with no `entity_item`, a NULL and an empty plural name, and the stderr warning |
| loc_depth | every location subkind, including completed and unfinished rafts at building depth and their subloc/viewloc walks; the stderr message and assertion for subkinds that are not locations and for box 0 |
| here lists, stacks | `in_here_list`, `add_to_here_list` and `remove_from_here_list` (duplicates, and the assertions), `stack_parent`, `stack_leader` (and its loop and non-unit assertions), `loop_stack`, `loop_here` and `restore_stack_actions` over stacks with non-units, missing and deleted boxes in the here lists |
| skill_school, macros | parent chains, a loop (the assertion), and the `oly.h` attribute macros with and without their substructures |
| is_unit and player | members and non-members, a player with no `entity_player`, `loop_units`, a non-player (the assertion), and a lord loop in `player` (the assertion) |
| clear_temps and loops, safe_qsort | `clear_temps`, `loop_kind`, `loop_inv` skipping missing items and non-positive quantities, and `safe_qsort` with `int_comp` |
| subloc and viewloc | units, stacks, ships, cities, castles, graveyards, sewers, faery hills, regions and nowhere, and `viewloc` of a unit (the assertion) |
| parse_arg | numbers, codes, negative numbers, and `garrison` and its near misses with and without a garrison or a location; `first_character` and `garrison_here` |
| set_lord and init_load_sup | swearing units with no player to a player and to a unit, with and without loyalty changes, `prev_lord`, unsorted unit lists, commands that are loading, running, in error (the assertion) or another unit's (the assertion); queued orders that parse (conditions, comments, control bytes, fuzzy command names, arguments, an empty order, the longest order and one too long, the assertion), dead units and units with no player; and a unit that is nowhere, whose `touch_loc` fails `loc_depth`'s assertion |
| touch_loc | units in cities, castles, provinces, ships, inner locations, stacks, and a region (the assertion) |

## Results

Last run: 2026-10-27, Linux amd64, GCC 12.2.0, Go 1.26.4.

- 37 sections, 1487 lines. The C driver is deterministic, its output matches `helper-c-output.txt`, and `TestHelperMatchesC` passes.
