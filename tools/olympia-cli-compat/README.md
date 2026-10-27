# olympia command-line differential test

`run.sh` checks that the Go engine's command line (`cmd/olympia`) matches G3's C `olympia-g3`.
G3 is the reference.
Expected results always come from running the pinned C build; nothing here is re-baselined from Go output.
Two sources of expectations are used:

- G3's `rng-selftest` gate (`tests/olympia/regress/rng-selftest`), which runs unchanged from the pinned export with G3's committed `EXPECT.txt`.
- A downstream matrix of argument vectors, [`cases.txt`](cases.txt). Its expectation for each vector is the pinned C build's output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `6bbcd63055d051ade6f745f368281628441a2c27`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.

The port follows `main()` in `olympia/main.c` (the long-option pass, the `getopt()` loop, the usage text and what `main()` prints before it dispatches), `test_random()` in `olympia/z.c` and `ilist_test()` in `lib/ilist.c`.
Every other mode stops with `<argv[0]>: <what>: not yet implemented` on stderr and exit status 2.
Every such mode reaches `load_db()` first, so the message names `load_db` and the modes that were asked for.

### `sizeof(struct box)`

`main()` prints `sizeof(struct box) = N`, and the Go engine has to print the bytes of the C build it is compared with.
The Go value is `sizeofBox` in [`internal/olympia/box.go`](../../internal/olympia/box.go): `cSizeof` walks the Go `box` type by reflection and lays it out as C does (`int8`/`int16`/`int32` as themselves, every pointer, slice, string and func as a C pointer, with C's padding).
So it gives the size for the host it is built on, as a C build on that host does: 144 on the 64-bit reference hosts.
The `sizeof(int) != sizeof(int *)` warning before it is printed on the same condition.
The entity types have to change when the pinned revision changes `struct box`, and nothing that isn't a C field may be added to them; the `rng-selftest` gate and the matrix catch a mismatch.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another.
- What G3 needs to build `olympia-g3`: CMake 3.25 or later, Ninja and a C11 compiler.
- Go (the version in `go.mod`), Bash, Perl, Python 3 (the gate runs `-R` on a pseudo-terminal), `shasum`, `cmp` and `diff`.
- macOS. The engine parses its flags with the C library's `getopt()`, and the port reproduces the BSD libc that macOS ships: it does not permute, says `illegal option`, and names `argv[0]` as given.
  glibc's `getopt()` differs on all three.

## Running

```sh
tools/olympia-cli-compat/run.sh
```

The script:

1. Exports the pinned revision with `git archive` into `tmp/olympia-cli-compat/g3` and builds `olympia-g3` there with the `debug` preset.
   It uses neither the working tree nor an existing build of the clone.
2. Builds the Go port.
3. Runs every vector in `cases.txt` with both engines, each in an empty directory with stdin from `/dev/null` and `argv[0]` set to `olympia-g3`.
   Each vector runs twice per engine: with stdout and stderr in separate files, and with both in one file (`2>&1`), which shows how C's buffered stdout and unbuffered stderr interleave.
4. For a `same` vector, requires Go to match C byte for byte on stdout, stderr, the merged file, and both exit statuses.
   For an `nyi` vector, requires Go to exit nonzero, print `not yet implemented` on stderr, and print the start of C's stdout.
5. Runs G3's `rng-selftest` gate from the export, unchanged, against the C build and against Go (through `OLYMPIA_ENGINE_BIN`).
   It checks `-R` on a pipe against `EXPECT.txt`, and that the 10×10 grid printed on a terminal is the same draws in draw order.
   Each must print `YES`.

The script exits 0 only when every check passes.
Each run's outputs are left under `tmp/olympia-cli-compat/runs/<case>/{c,go}`, and the summary is written to `tmp/olympia-cli-compat/results.tsv`.

The gate can also be run by hand from the pinned clone:

```sh
go build -o tmp/olympia ./cmd/olympia
cd ~/Software/playbymail/olympia-g3
OLYMPIA_ENGINE_BIN=/path/to/g3g0/tmp/olympia tests/olympia/regress/rng-selftest/check.sh
```

`go test ./internal/olympia` checks `-R` (on a pipe and on a terminal), `-t`, the usage errors and the unported modes against output recorded from the pinned C build.
`go test ./internal/clib` checks the `getopt()` emulation.

## Case matrix

| area | cases |
|---|---|
| `-R` | alone, twice, with an operand, unknown options before and after it and clustered with it (`-xR`, `-Rx`), a missing argument after it, every flag before it, `-l` and `-G` attached and apart, `-l ''`, `-f` before and after it, `-?`, `--` and `-` after it |
| long options | `test-use-const-report-date` and each `--` path option before and after `-R`; a path that is missing or empty, which is an error but still lets `-R` win; a path that looks like an option (`--coverage -R`, `--coverage --use-coverage`), which is taken as the path; `--coverage=x`, which is not a long option and reaches `getopt()` as a cluster |
| `-t` | alone, twice, before and after an unknown option and `-R`, clustered with an unknown option (`-tx`, `-xt`), after `-f`, `-l dir`, mode flags and a long option, and with an operand after it |
| usage errors | unknown options alone, in a cluster and around valid ones; `-?` and `-:`; `--foo`; `-l` and `-G` without an argument; every flag before an unknown option; `-f` before and after one, which changes how stdout and stderr interleave; an operand after one |
| not yet implemented | no arguments, each mode flag, `-r -S`, each long option, `--`, `-- -R`, a lone `-`, an operand before `-R`, and `-l -R` |

## Results

Last run: 2026-10-26, macOS 26.6.2 arm64, Apple clang 21.0.0, CMake 4.3.3, Go 1.26.4.

- 92 vectors: 62 `same` and 30 `nyi`. All 92 PASS.
- G3 gate: `rng-selftest` prints `YES` for the C build and for Go.

The full table of hashes and exit statuses is in [`results.tsv`](results.tsv), copied from the run, not edited by hand.

## Known differences

- **Buffer size.** C flushes a full stdout buffer early, at a size that depends on the file; the port holds stdout until exit (or line by line on a terminal, or not at all after `-f`).
  Nothing ported yet writes enough to stdout for that to change how it interleaves with stderr.

## Upstream defects reproduced on purpose

None are known at the pinned revision.
