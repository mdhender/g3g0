# box-file pass and writer differential test

`run.sh` checks that the Go port's box-file passes and writers (`internal/olympia`: `io.go`) match G3's `olympia/io.c`: pass 1 (`fast_scan`, or `scan_boxes`, `scan_chars` and `scan_all_boxes`), pass 2 (`read_boxes`, `read_chars` and `read_all_boxes`), and the writers (`open_write_fp`, `write_kind`, `write_player`, `write_chars`, `write_leftovers`, `write_all_boxes` and `write_master`), with the files they read and write, which files exist afterwards, their modes, and the diagnostics.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `e733acb63bacb21167a1f2589d8dcdc023364896`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.
The script exports that revision with `git archive`, so the clone's checkout does not matter.

The port follows:

- `olympia/io.c`, lines 2126 to 2640 at the pinned revision (`open_write_fp` to `write_master`).
- `olympia/u.c`: `lookup()` and `stage()`; `olympia/swear.c`: `lord()` and `player()`, which `write_master()` uses; `olympia/sout.c`: `sout()`'s cut to `LEN` bytes.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another.
- A C11 compiler (`CC`, default `cc`), Go, Bash, Python 3, `git`, `tar`, `cmp`, `diff`, `find`, `stat` and `sort`.
- Optional: the inputs `tools/olympia-roundtrip/run.sh` caches under `tmp/olympia-roundtrip/cases` (`GAMES`). Without them only `world.tgz` is used for check 2.

## Running

```sh
tools/boxfile-compat/run.sh             # compare the recorded cases (check 1)
tools/boxfile-compat/run.sh --games 10  # and 10 saved games, always the same ones (check 2)
tools/boxfile-compat/run.sh --full      # and every saved game (check 2)
tools/boxfile-compat/run.sh --update    # rewrite internal/olympia/testdata/boxfile-c-output.txt from C
```

Check 2 runs only with `--full` or `--games N`; `N` games are `world.tgz` and cached games spread evenly over the sorted list.
[`tools/README.md`](../README.md) says when to run which.

The script:

1. Exports `lib/`, `olympia/` and `tests/olympia/fixtures/world.tgz` at the pinned revision with `git archive` into `tmp/boxfile-compat/g3`.
2. Generates the stubs that `tools/box-compat` uses, and copies out of the export, verbatim: `sout()`, `set_bit()`, `test_bit()`, `lookup()`, `stage()`, `int_comp()`, `lord()` and `player()` (under `#line`, so that `player()`'s assertion names `olympia/swear.c` and its line), the command parser of `olympia/input.c`, and `olympia/io.c` from its first line to just before `load_system()` (just after `write_master()`), under `#line 1 "olympia/io.c"`, followed by io.c's `remove_dir()` (which `write_all_boxes()` calls), under `#line` with its own line number.
3. Builds a C driver (`tmp/boxfile-compat/driver.c`) that `#include`s that copy of io.c, linked with the exported `code.c`, `glob.c`, `z.c` and the `lib/` modules.
4. Check 1: runs each section of [`internal/olympia/testdata/boxfile-cases.txt`](../../internal/olympia/testdata/boxfile-cases.txt) in a fresh process and directory, with umask 022, twice, and lists the directory's tree after each (every file's mode and contents). It requires the same output both times and requires it to equal [`boxfile-c-output.txt`](../../internal/olympia/testdata/boxfile-c-output.txt). Then it runs `go test -run TestBoxfileMatchesC`, which runs the same cases on the Go port.
5. Check 2, with `--full` or `--games N`: for `world.tgz` and every cached game input (`tmp/olympia-roundtrip/cases/*/input/lib`), three ways (with `master`; with `master` and `pretty_data_files`; and with `master` removed, so that pass 1 is `scan_all_boxes()`), the driver (`-g`) runs, from a copy of the game with `libdir` `lib`, what `load_db()` and `save_db()` run of these: `fast_scan()` or else `scan_all_boxes()`, `read_all_boxes()`, `write_all_boxes()` and `write_master()`. `TestBoxfileWorlds` (run with `BOXFILE_COMPAT_GAMES`) does the same in Go on its own copy. The two `lib` trees must be identical (`diff -r`, and every path's mode), and so must stdout and stderr.

`go test ./...` runs `TestBoxfileMatchesC` too, so the port is held to the recorded C output without a G3 clone. `TestBoxfileWorlds` is skipped unless `G3G0_FULL=1` is set and `BOXFILE_COMPAT_GAMES` names a list of games, as `run.sh` does: it builds the package's test binary once and runs it on one game at a time, and deletes each game's copies and output as soon as it matches C's, keeping only failing games'. `-short` always skips it.

### Lines between boxes

`read_boxes()` skips a comment line silently, and reports any other line that is not a box number (such as a second blank line between boxes) with `read_boxes: unexpected line` and skips it.
Before upstream fixed G3 [#316](https://github.com/playbymail/olympia-g3/issues/316), it never moved past either and looped forever; `run.sh` and the Go test no longer need to stop such a section.
`run.sh` stops any section that runs for more than 10 seconds, and fails.

### What the abort and crash statuses mean for the files

`asfail()` explicitly flushes every open stdio stream before `abort()`, so a section that fails an assertion while writing a file (`write_player()` of a unit that is not a character, `player()` of a lord loop in `write_master()`) leaves the file with everything written to it so far; the port's assertion panic is recovered with `clib.FlushAll()`. A crash (status 139) flushes nothing.

## What the cases cover

| section | covers |
|---|---|
| fast_scan | `master` lines with extra white space, no `.`, no subkind, junk, a subkind past `schar`, no newline at the end; no `master`, an empty one; blank lines and comments rejected before allocation; duplicate box (`alloc_box()` assertion); another `libdir` |
| scan_boxes | comments before and inside a box, a line of white space ending a box, kinds and subkinds in other cases, subkind `0` and `0x`, a subkind with a space, bad kinds and subkinds (the `scan_boxes(N)` diagnostics), short and tab-separated box lines, kind `deleted`, no blank line at the end; a missing file; a blank line between boxes and a box number 0 (assertion); a duplicate box across files |
| scan_chars, scan_all_boxes | fact files whose names do and do not start with a digit, a directory among them; no `fact`, `fact` a file, an unreadable `fact`; every box file, and none of them (the `stage()` line and every diagnostic) |
| read_boxes, read_chars, read_all_boxes | box lines with leading white space, the end of the file, a continued field on the file's last line; a missing file; comment lines first, between boxes, last and alone, one and three extra blank lines between boxes, blank lines first and last, and junk, indented and box-number-0 lines between boxes and last, each skipped (G3 [#316](https://github.com/playbymail/olympia-g3/issues/316), fixed upstream); a line that is not a box and a negative box number; a box pass 1 did not allocate (`load_box()`'s assertion); `fact` with other names and a dot directory; no `fact`; a box in two files |
| write_kind | chain order, deleted boxes, an empty kind, pretty names; no `libdir`; a read-only, a write-only and a directory in the file's place |
| write_player, write_chars, write_leftovers | units in list order, a unit twice, a deleted unit, a unit with no box; no units; no `fact`; a unit that is a location or a player (assertion, with the file as `abort()` leaves it); what each leaves to `misc`, and `temp` |
| write_all_boxes | every kind, an old `fact` file and stale box files; no `fact`, no `libdir`, a read-only `libdir`; `fact` cleared by `remove_dir()`, not the shell (G3 [#320](https://github.com/playbymail/olympia-g3/issues/320), fixed upstream): files, a dot file, nested and empty directories, and symbolic links to a directory, a file and nowhere (removed, not followed), `fact` a file, a symbolic link to a directory (its contents removed, the link kept), unreadable, and holding an entry it cannot remove (no message for any of them); a `libdir` with a space, which clears its own `fact` and nothing else, with and without `win_flag`; and `win_flag` (no effect) with and without `fact` |
| write_master | every kind, roads indexed as road, names with tabs, the lord chain to a player, no lord and a lord that is not a player (`fact/0`), pretty (no effect); a lord loop (assertion), a lord with no box (crash); no `libdir`; no boxes |
| round trip | `scan_all_boxes`, `read_all_boxes`, `write_all_boxes` and `write_master` on a small game |

## Results

Last run: 2026-10-27, Linux amd64, GCC 12.2.0, Go 1.26.4.

- Check 1: 70 sections, 1771 lines. The C driver is deterministic, its output matches `boxfile-c-output.txt`, and `TestBoxfileMatchesC` passes.
- Check 2 was not rerun at this pin. Historical macOS `--full`: 197 games (`world.tgz` and the 196 cached round-trip inputs; 7,947,575 boxes in their written `master` files, 36,999 in `world.tgz`), each with `master`, with `master` and pretty, and without `master`: C and Go `lib` trees (files, modes and contents), stdout and stderr identical.

## Upstream defects reproduced on purpose

None at the pinned revision. (`read_boxes()` used to loop forever on a comment line or a line that is not a box number, G3 [#316](https://github.com/playbymail/olympia-g3/issues/316); and `write_all_boxes()` cleared `fact` with `system("rm -rf LIBDIR/fact")`, libdir unquoted, G3 [#320](https://github.com/playbymail/olympia-g3/issues/320); upstream fixed both at the pinned revision.)
