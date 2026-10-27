# entity allocation and codes differential test

`run.sh` checks that the Go port's box allocation, kind and subkind chains, entity codes and names (`internal/olympia`: `code.go`, `loop.go`, `io.go`'s `ifName`) match G3's `olympia/code.c`, the entity loops of `olympia/loop.h` and `if_name()` in `olympia/io.c`.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `2567e6a5dd0ba0d95516d8bed9e1bafb87a51985`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.

The port follows, from `olympia/code.c`: `letter_val`, `int_to_code`, `code_to_int`, `print_box_usage`, `scode`, `name`, `set_name`, `set_banner`, `box_code`, `box_code_less`, the chains (`add_next_chain`, `remove_next_chain`, `add_sub_chain`, `remove_sub_chain` and their cache), `delete_box`, `change_box_kind`, `change_box_subkind` and `alloc_box`.
From `olympia/oly.h`: `kind_first`, `kind_next`, `sub_first`, `sub_next` and `loc`.
From `olympia/loop.h`: `loop_boxes`, `loop_kind`, `loop_subkind` and `loop_loc_or_ship`.
`new_ent` and the mint allocator are left for later tasks. `display_name` and the name helpers built on it that `check_db` uses (`box_name`, `just_name`, `plural_item_name`) are checked by [`helper-compat`](../helper-compat); `display_kind` and the rest are left for later tasks.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another.
- A C11 compiler (`CC`, default `cc`), Go, Bash, Python 3, `git`, `cmp` and `diff`.

## Running

```sh
tools/code-compat/run.sh            # compare
tools/code-compat/run.sh --update   # rewrite internal/olympia/testdata/code-c-output.txt from C
```

The script:

1. Exports `lib/` and `olympia/` at the pinned revision with `git archive` into `tmp/code-compat/g3`.
2. Generates stubs for the command routines `cmd_tbl` names, `sout()` and the three routines `code.c` calls only from `display_name()` and its users, and copies `if_name()` (static in `io.c`) out of the export.
3. Writes a C driver to `tmp/code-compat/driver.c` and links it with the exported `olympia/code.c`, `olympia/glob.c`, `olympia/z.c`, `lib/checked_alloc.c`, `lib/ilist.c`, `lib/rng.c` and `lib/rnd.c`. It compiles from the export's root, so `asfail()` prints `olympia/code.c` as the file, as the Go port does.
4. Runs each section of [`internal/olympia/testdata/code-cases.txt`](../../internal/olympia/testdata/code-cases.txt) in a fresh process, since a section may end in a failed assertion, which aborts. stderr goes to stdout, in order, and `exit STATUS` follows each section. It does this twice and requires the same output both times.
5. Requires that output to equal [`internal/olympia/testdata/code-c-output.txt`](../../internal/olympia/testdata/code-c-output.txt), or with `--update`, writes it there.
6. Runs `go test -run TestCodeMatchesC ./internal/olympia`, which runs the same cases on the Go port and requires the recorded C output line for line.

`go test ./...` runs `TestCodeMatchesC` too, so the port is held to the recorded C output without a G3 clone.

## What the cases cover

| section | covers |
|---|---|
| codes | `int_to_code` at the edges of every range; `code_to_int` and `scode` on valid, malformed and high-bit input; `letter_val`; and checksums of `int_to_code` and back over every box number from -10 to `MAX_BOXES`+10, of `code_to_int` and `scode` over every string of up to four bytes from a 19-byte alphabet, and of `code_to_int` over every three- and four-byte code shape in both cases |
| allocation, deletion, changes | chain order after allocation in increasing, decreasing and mixed order across kinds and subkinds (which exercises the chain caches), after `delete_box` of the head, middle and tail and allocation after it, and after `change_box_kind` and `change_box_subkind`, including no-op changes and kind 0 |
| loops | `loop_kind` and `loop_subkind` while the body deletes the current box or allocates one before or after it; `loop_boxes`; `loop_loc_or_ship` with and without ships |
| names | `name`, `set_name` (brackets become braces), `set_banner` (cut to 50 bytes, cleared by an empty string), `if_name` with and without `pretty_data_files`, `box_code` and `box_code_less` |
| usage | `print_box_usage` on an empty world and with boxes at the ends of each range |
| assertions | `alloc_box` of a box in use (deleted or not) or out of range, `delete_box` twice or of a missing box, changing the kind or subkind of a deleted box, and `name`/`set_name` of a missing or deleted box: the stderr text and exit status 134 |

## Results

Last run: 2026-10-27, Linux amd64, GCC 12.2.0, Go 1.26.4.

- 20 sections, 652 lines. The C driver is deterministic, its output matches `code-c-output.txt`, and `TestCodeMatchesC` passes.

## Upstream defects reproduced on purpose

`print_box_usage()` (`olympia/code.c` lines 178-209): its ranges disagree with the coding system at the top of the file ("locs" is 20000-59999, leaving out the provinces at 10000-19999 and overlapping "play", "lloc" and "subl"; "rest" starts at 79900, not 79000), and it prints and divides by `high - low`, one less than the size of each range. Not yet filed in G3.
