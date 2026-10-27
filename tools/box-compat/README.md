# entity section, load_box and save_box differential test

`run.sh` checks that the Go port's entity sections, `loadBox` and `saveBox` (`internal/olympia`: `io.go`), and the command parser they use (`input.go`, and `iStrcmp`/`fuzzyStrcmp` in `z.go`), match G3's `olympia/io.c` and `olympia/input.c`.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `66b03fb09b541560268e66b553179d06ea30de80`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.

The port follows:

- `olympia/io.c`, lines 773 to 2081: every `print_*`/`scan_*` pair (`loc_info`, `magic`, `char`, `loc`, `subloc`, `item`, `item_magic`, `player`, `skill`, `command`, `gate`, `misc`), `load_box` and `save_box`.
- `olympia/input.c`: `parse_line`, `remove_comment`, `remove_ctrl_chars`, `find_command` (and `fuzzy_find`) and `oly_parse_cmd`, which `scan_command` calls.
- `olympia/z.c`: `i_strcmp` and the `fuzzy_*` matchers.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another.
- A C11 compiler (`CC`, default `cc`), Go, Bash, Python 3, `git`, `tar`, `cmp` and `diff`.
- Optional: the inputs `tools/olympia-roundtrip/run.sh` caches under `tmp/olympia-roundtrip/cases` (`GAMES`). Without them only `world.tgz` is used for check 2.

## Running

```sh
tools/box-compat/run.sh            # compare
tools/box-compat/run.sh --update   # rewrite internal/olympia/testdata/box-c-output.txt from C
```

The script:

1. Exports `lib/`, `olympia/` and `tests/olympia/fixtures/world.tgz` at the pinned revision with `git archive` into `tmp/box-compat/g3`.
2. Generates the stubs that `tools/io-compat` uses, and copies out of the export, verbatim: `sout()`, `set_bit()`, `test_bit()`, `int_comp()`, the command parser of `olympia/input.c`, and `olympia/io.c` from its first line to just before `open_write_fp()` (just after `save_box()`), under `#line 1 "olympia/io.c"` so that assertions name `olympia/io.c` and its lines.
3. Builds a C driver (`tmp/box-compat/driver.c`) that `#include`s that copy of io.c, linked with the exported `code.c`, `glob.c`, `z.c` and the `lib/` modules.
4. Check 1: runs each section of [`internal/olympia/testdata/box-cases.txt`](../../internal/olympia/testdata/box-cases.txt) in a fresh process and directory, twice, requires the same output both times and requires it to equal [`box-c-output.txt`](../../internal/olympia/testdata/box-c-output.txt). Then it runs `go test -run TestBoxMatchesC`, which runs the same cases on the Go port.
5. Check 2: for `world.tgz` and every cached game input (`tmp/olympia-roundtrip/cases/*/input/lib`), plain and with `pretty_data_files` set, the driver (`-w`) allocates every box from the box files' `n kind subkind` lines, loads `loc`, `item`, `skill`, `gate`, `road`, `ship`, `unform`, `misc` and `fact/*` (sorted) with `load_box`, and saves every box with `save_box`. `TestBoxWorlds` (run with `BOX_COMPAT_GAMES`) does the same in Go, and stdout and stderr must be identical for every game.

`go test ./...` runs `TestBoxMatchesC` too, so the port is held to the recorded C output without a G3 clone. `TestBoxWorlds` is skipped unless `run.sh` runs it.

The driver's load loop is `read_boxes()`'s, except that a line that is not a box number ends the file (printing `stop`), where `read_boxes()` skips a comment line silently and reports and skips any other. Pass 1 (`scan_boxes`) and pass 2 (`read_boxes`) themselves are tested by [`tools/boxfile-compat`](../boxfile-compat/README.md).

## What the cases cover

| section | covers |
|---|---|
| every field and section | every tag of every section, in an order unlike save_box's; continuation lines in `il`, `tl`, `sl`, `rq`; plain and pretty output |
| names and pretty | `na` lines and `if_name` |
| zero_check | the six sections C suppresses when all zero, including fields io.c never saves (`next`, `mage_worked`, `one_turn_use`, `opium_double`) and an allocated empty list; the sections that always print; `CO` with no command; deleted and missing boxes |
| diagnostics | an unknown or short tag in every section, high-byte tags, box references that are not boxes, a comment line, an unknown top-level tag and the lines it skips |
| line endings | boxes without a blank line between them, two blank lines, a comment inside a section, end of file without a blank line |
| unforseen box | `load_box` of a box pass 1 did not allocate: the diagnostic with `libdir`, and the assertion |
| continuation at end of file | a continued field on the file's last line, at top level and in a section: the list ends there (G3 [#313](https://github.com/playbymail/olympia-g3/issues/313), fixed upstream), and `save_box` writes what was read |
| empty continuation line | a continued `il`, `sl`, `tl` and `rq` followed by a blank or white-space line, at top level and in a section: the list ends there (G3 [#318](https://github.com/playbymail/olympia-g3/issues/318), fixed upstream), after a `bad item 0`/`bad skill 0` diagnostic for the empty entry; the blank line was the box separator, so `load_box` reads on into the next box's lines as bad lines and its later fields |
| fly_cap and ride_cap | `fc` without `rc`, `rc` without `fc`, `rc` with `fc 0`, both set, a negative `fc`, and `fc 0` alone (an all-zero `IT`, not printed): `print_item` saves `fc` when `fly_cap` is set, whatever `ride_cap` is (G3 [#314](https://github.com/playbymail/olympia-g3/issues/314), fixed upstream) |
| narrowing | every `(schar)`, `(short)`, `(uchar)` and `int` field past its range, `sscanf` of a short `ar` line, `%c` of a high byte |
| commands | `oly_parse_cmd` and `find_command`: conditions, comments, quotes, control and high bytes, exact and fuzzy matches, unknown commands; `li` and `us` in a box, and the `bad cmd` diagnostic; `remove_comment`'s quotes (G3 [#315](https://github.com/playbymail/olympia-g3/issues/315), fixed upstream): a `#` inside a double- or single-quoted token is kept, an unterminated quote runs to the end of the line, a quote inside a word (`Joe's`, `a"#b"`) is ordinary text, a quote right after a closing quote starts another, and a quote after a `#` is ignored; in an `li` line too |
| real boxes from world.tgz | 25 boxes chosen so that every section and list tag appears, with every box they reference allocated |

## Results

Last run: 2026-10-27, macOS 26.6.2 arm64, Apple clang 21.0.0, Go 1.26.4.

- Check 1: 20 sections, 3625 lines. The C driver is deterministic, its output matches `box-c-output.txt`, and `TestBoxMatchesC` passes.
- Check 2: 197 games (`world.tgz` and the 196 cached round-trip inputs), 7,947,575 saved boxes (36,999 in `world.tgz`), plain and pretty: C and Go output identical.

## Upstream defects reproduced on purpose

None at the pinned revision. (`print_item()` used to print `fc` (`fly_cap`) only when `ride_cap` was set, G3 [#314](https://github.com/playbymail/olympia-g3/issues/314); upstream fixed it at the pinned revision.)
