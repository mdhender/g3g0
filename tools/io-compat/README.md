# line reader and io.c field helper differential test

`run.sh` checks that the Go line readers and field helpers (`internal/olympia`: `z.go`, `io.go`, `u.go`) match G3's `olympia/z.c` and the top of `olympia/io.c`.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `e733acb63bacb21167a1f2589d8dcdc023364896`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.

The port follows:

- `olympia/z.c`: `getlin`, `getlin_ew`, `eat_leading_trailing_whitespace`, `readfile`, `closefile`, `readlin` and `readlin_ew`, and the `iswhite` and `isdigit` macros of `olympia/z.h`.
- `olympia/io.c`, lines 1 to 817 (up to `print_loc_info`): the `linehash`, `t_string` and `advance` macros, the `line` and `ext_boxnum` statics, `zero_check`, `if_name`, `box_scan`/`box_print`, `boxlist0_scan`/`boxlist0_print`, `boxlist_scan`/`boxlist_print`, `admit_print_sup`/`admit_print`/`admit_scan`, `ilist_print`/`ilist_scan`, `known_print`/`known_scan`, `scan_ints`, the `skill_list`, `item_list`, `trade_list` and `req_list` pairs, and `olytime_scan`/`olytime_print`.
- `olympia/u.c`: `set_bit` and `test_bit`. `known_print` sorts with `int_comp` (`olympia/seed.c`) through `loop_known` (`olympia/loop.h`).

## Prerequisites

G3 #345 resolves the former Linux blocker (G3 #342): `scan_ints` now
reads unlimited decimal tokens with `strtoll`, saturates at signed 64 bits,
then keeps the low 32 bits (low 16 for olytime shorts). A failed conversion
leaves that field and the following fields unchanged. The recording was
regenerated from this real pinned C; generic `clib.Scanner` is unchanged.

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
| olytime | partial and overflowing decimal input, which leaves unread fields as they were; printing all-zero and signed values |
| scanf number width | historical section name: `scan_ints` consumes the entire numeric token, including leading zeros beyond 512 bytes |
| continuation at end of file | a continued field on a file's last line, with and without a final newline, in every list scan (`boxlist`, `boxlist0`, `admit`, `ilist`, `known`, skill, item, trade and req lists): `readlin_ew` returns NULL and the scan keeps what it has read; the next `readlin` is NULL |
| empty continuation line | a blank or white-space line after a trailing backslash in the skill, item, trade and req lists, and an item or trade tag with no data: the empty entry gets its `bad ... 0` diagnostic and the list ends there (G3 [#318](https://github.com/playbymail/olympia-g3/issues/318), fixed upstream); the next `readlin` is the line after it |
| short entries | skill and req entries with too few numbers, on the first line and on continuation lines: each missing field is 0 for that entry (G3 [#311](https://github.com/playbymail/olympia-g3/issues/311), fixed upstream), not the previous entry's value |
| bare tags | a bare `sl`, `il`, `tl` and `rq` tag, with and without a trailing tab: one all-zero entry, reported (`bad skill 0`, `bad item 0`) and dropped, and not a continuation (G3 [#312](https://github.com/playbymail/olympia-g3/issues/312)) |
| zero_check | a `loc_info` that is all zero, has a field set, or has an empty but allocated list |

Out of scope, because C is undefined there: `readlin` before any `readfile` (C dereferences a NULL `point`).

## Results

Last run: 2026-10-27, Linux amd64, GCC 12.2.0, Go 1.26.4.

- 49 sections, 1598 lines. The C driver is deterministic, its output matches `io-c-output.txt`, and `TestIOMatchesC` passes.

## Upstream defects reproduced on purpose

None at the pinned revision. Upstream has fixed the ones this port used to reproduce: a short skill or req entry read uninitialized locals (G3 [#311](https://github.com/playbymail/olympia-g3/issues/311)), and the continuation test read `s[strlen(s)-1]` before an empty entry (G3 [#312](https://github.com/playbymail/olympia-g3/issues/312) and [#318](https://github.com/playbymail/olympia-g3/issues/318)).
