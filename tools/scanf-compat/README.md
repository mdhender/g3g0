# sscanf %d differential test

`run.sh` checks that `clib.Scanner` (`internal/clib/stdio.go`) matches the C library's `sscanf("%d")` on macOS, the reference platform.
The C library is the reference.
Expected results always come from running a C driver on macOS; nothing here is re-baselined from Go output.

`Scanner.Int` is the `%d` used by `internal/mapgen` for `sscanf(buf, "%d,%d", …)` (`mapgen/mapgen.c:1764`) and `sscanf(buf, "%d,%d %c", …)` (`mapgen/mapgen.c:1940`), and by `sscanD` in `internal/olympia` for the `%d` and `%hd` scans of `olympia/io.c`.

## Prerequisites

- macOS. Other C libraries may differ (the 512-character limit, for one, comes from the BSD `vfscanf`), so record expectations only on macOS; the script warns elsewhere.
- A C11 compiler (`CC`, default `cc`), Go, Bash, `cmp` and `diff`. No G3 clone is needed.

## Running

```sh
tools/scanf-compat/run.sh            # compare
tools/scanf-compat/run.sh --update   # rewrite internal/clib/testdata/scanf-c-output.txt from C
```

The script:

1. Writes a C driver to `tmp/scanf-compat/driver.c` and builds it with `cc -std=c11 -Wall`.
2. Runs it twice over [`internal/clib/testdata/scanf-cases.txt`](../../internal/clib/testdata/scanf-cases.txt) and requires the same output both times.
   Each case is a format (`d`, `dd`, `d,d` or `d,d c`, each followed by `%n`) and an input; the driver prints `sscanf`'s return value, the targets and where scanning stopped.
3. Requires that output to equal [`internal/clib/testdata/scanf-c-output.txt`](../../internal/clib/testdata/scanf-c-output.txt), or with `--update`, writes it there.
4. Runs `go test -run TestScannerMatchesC ./internal/clib`, which runs the same cases through `Scanner` and requires the recorded C output line for line.

`go test ./...` runs `TestScannerMatchesC` too.

## What the cases cover

| area | covers |
|---|---|
| values and signs | `+`, `-`, `-0`, leading zeros, trailing junk, `0x10` (only the `0` is read) |
| white space | every `isspace()` byte before the number; none allowed between the sign and the digits |
| failures | empty and all-blank input (input failure, `EOF`), a lone sign, a sign followed by a non-digit (matching failure) |
| overflow | values around the 32- and 64-bit limits: `strtoimax()` saturates to the 64-bit range and the `int` keeps the low 32 bits, so `99999999999999999999` scans as -1 and `-99999999999999999999` as 0 |
| 512-character limit | one conversion reads at most 512 characters, sign and leading zeros included; the rest of the digits stay in the input and the next `%d` reads them |
| mapgen shapes | `%d,%d` and `%d,%d %c` with blanks, missing fields, overflow and over-long numbers |

## Results

Last run: 2026-10-26, macOS 26.6.2 arm64, Apple clang 21.0.0, Go 1.26.4.

- 78 cases. The C driver is deterministic, its output matches `scanf-c-output.txt`, and `TestScannerMatchesC` passes.
