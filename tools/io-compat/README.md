# line reader and io.c field helper differential test

`run.sh` checks that the Go line readers and field helpers (`internal/olympia`: `z.go`, `io.go`, `u.go`) match G3's `olympia/z.c` and the top of `olympia/io.c`.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `6bbcd63055d051ade6f745f368281628441a2c27`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.

The port follows:

- `olympia/z.c`: `getlin`, `getlin_ew`, `eat_leading_trailing_whitespace`, `readfile`, `closefile`, `readlin` and `readlin_ew`, and the `iswhite` and `isdigit` macros of `olympia/z.h`.
- `olympia/io.c`, lines 1 to 735 (up to `print_loc_info`): the `linehash`, `t_string` and `advance` macros, the `line` and `ext_boxnum` statics, `zero_check`, `if_name`, `box_scan`/`box_print`, `boxlist0_scan`/`boxlist0_print`, `boxlist_scan`/`boxlist_print`, `admit_print_sup`/`admit_print`/`admit_scan`, `ilist_print`/`ilist_scan`, `known_print`/`known_scan`, the `skill_list`, `item_list`, `trade_list` and `req_list` pairs, and `olytime_scan`/`olytime_print`.
- `olympia/u.c`: `set_bit` and `test_bit`. `known_print` sorts with `int_comp` (`olympia/seed.c`) through `loop_known` (`olympia/loop.h`).

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another.
- A C11 compiler (`CC`, default `cc`), Go, Bash, Python 3, `git`, `cmp` and `diff`.

## Running

```sh
tools/io-compat/run.sh            # compare
tools/io-compat/run.sh --update   # rewrite internal/olympia/testdata/io-c-output.txt from C
```

The script:

1. Exports `lib/` and `olympia/` at the pinned revision with `git archive` into `tmp/io-compat/g3`.
2. Generates stubs for the command routines `cmd_tbl` names and the three routines `code.c` calls only from `display_name()`, and copies out of the export, verbatim: `sout()` (`olympia/sout.c`, the real one, so `if_name` is cut to `BUF_LEN` as in the engine), `set_bit()` and `test_bit()` (`olympia/u.c`), `int_comp()` (`olympia/seed.c`), and `olympia/io.c` from its first line to just before `print_loc_info()`.
   The io.c helpers are `static`, so the driver `#include`s that copy of the top of io.c, which also gives it the `line` and `ext_boxnum` statics.
3. Writes a C driver to `tmp/io-compat/driver.c` and links it with the exported `olympia/code.c`, `olympia/glob.c`, `olympia/z.c`, `lib/checked_alloc.c`, `lib/ilist.c`, `lib/rng.c`, `lib/rnd.c` and the list modules the helpers use (`admits`, `item_ents`, `req_ents`, `skill_ents`, `trades`).
4. Runs each section of [`internal/olympia/testdata/io-cases.txt`](../../internal/olympia/testdata/io-cases.txt) in a fresh process and a fresh, empty directory, since a section may crash and writes its own input files. stderr goes to stdout, in order, and `exit STATUS` follows each section. It does this twice and requires the same output both times.
5. Requires that output to equal [`internal/olympia/testdata/io-c-output.txt`](../../internal/olympia/testdata/io-c-output.txt), or with `--update`, writes it there.
6. Runs `go test -run TestIOMatchesC ./internal/olympia`, which runs the same cases on the Go port and requires the recorded C output line for line.

`go test ./...` runs `TestIOMatchesC` too, so the port is held to the recorded C output without a G3 clone.

## What the cases cover

| section | covers |
|---|---|
| readlin, readlin_ew | empty lines, white space, CR, NUL and high bytes, the last line without a newline, end of file; `readlin_ew`'s clean-up (bytes from 0x80 up become spaces, as `char` is signed; a leading control byte is not skipped) |
| line length | lines of 2046 to 9000 bytes, cut to 2047 (`LEN`-1), with and without a newline; lines across the 8192-byte `read()` buffers |
| files that cannot be read | a missing file (the `perror` message), a directory (read as empty), an empty file |
| closefile and failures | `closefile` and a failed `readfile`, which leave C's `read()` buffer in place: `readlin` returns what is left in it, a partial line at its end included, and then NULL |
| getlin, getlin_ew | the same clean-up; newline-terminated lines are not cut, the last line without a newline is cut to 1023 bytes; a directory |
| eat_leading_trailing_whitespace, linehash, t_string, next | short lines, high-byte tags (negative hashes), the scanners' skip of the first byte |
| box_scan, box_print | valid, invalid, zero, negative, overflowing and malformed references and their diagnostic; `if_name` with and without `-p`, unnamed boxes, names past `sout`'s 2048 bytes |
| boxlist, boxlist0, ilist, known | continuation lines (with white space, empty, CR, a backslash mid-line), invalid references, zero, signs, overflow; printing with every wrap point; `known_print`'s sort of the caller's list |
| admit | scanning targets, senses and lists across continuation lines, bad targets and boxes; printing, including the cases that print nothing and the `count` that starts at 2 |
| skill, item, trade, req lists | multi-entry lines, invalid entries, narrowing to `char`/`short`/`schar`, numbers past the `long` range; the entries print leaves out |
| olytime | partial and overflowing `%hd` input, which leaves fields as they were; printing all-zero and signed values |
| scanf number width | `sscanf`'s 512-byte limit on one `%d` |
| continuation at end of file | `readlin_ew` returns NULL and C crashes (status 139); the Go port's nil dereference is its counterpart |
| zero_check | a `loc_info` that is all zero, has a field set, or has an empty but allocated list |

Out of scope, because C is undefined there: a skill or req line with too few numbers (the C reads uninitialized locals), an empty skill, item, trade or req entry (`s[strlen(s)-1]` reads before the string), and `readlin` before any `readfile` (C dereferences a NULL `point`).

## Results

Last run: 2026-10-26, macOS 26.6.2 arm64, Apple clang 21.0.0, Go 1.26.4.

- 35 sections, 1196 lines. The C driver is deterministic, its output matches `io-c-output.txt`, and `TestIOMatchesC` passes.

## Upstream defects reproduced on purpose

- `skill_list_scan()` and `req_list_scan()` (`olympia/io.c` lines 489-524 and 687-714) read `know`, `experience` and `consume` from locals that `sscanf()` may leave uninitialized, and every list scan with a continuation test reads `s[strlen(s)-1]`, before the string when it is empty. Both are undefined; the Go port uses 0 for the locals (keeping their value from entry to entry) and treats an empty entry as the last. Not yet filed in G3.
- A continuation backslash on the last line of a file makes `readlin_ew()` return NULL, which every list scan dereferences. Not yet filed in G3.
