# olympia command-line differential test

`run.sh` checks that the Go engine's command line (`cmd/olympia`) matches G3's C `olympia-g3`.
G3 is the reference.
Expected results always come from running the pinned C build; nothing here is re-baselined from Go output.
Two sources of expectations are used:

- G3's `rng-selftest` gate (`tests/olympia/regress/rng-selftest`), which runs unchanged from the pinned export with G3's committed `EXPECT.txt`.
- A downstream matrix of argument vectors, [`cases.txt`](cases.txt). Its expectation for each vector is the pinned C build's output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `e733acb63bacb21167a1f2589d8dcdc023364896`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.

The port follows `main()` in `olympia/main.c` (the long-option pass, the BSD `oly_getopt()` loop, the usage text and what `main()` prints before it dispatches), `test_random()` in `olympia/z.c` and `ilist_test()` in `lib/ilist.c`, and from #52 its load and save (`load_db()`, `check_db()`, `save_db()` with `-S`, `times_masthead()` and the closing output), which [`olympia-roundtrip`](../olympia-roundtrip) checks on saved games.
Stage 1 also ports the no-orders turn on an existing world; report artifacts remain deferred. Unsupported modes stop with `<argv[0]>: <what>: not yet implemented` on stderr and exit status 2:

- A vector that asks for a mode that runs after `load_db()` and is not ported (`-e`, `-E`, `-s`, `-a`, `-i`, `-M`, `-h`, `--use-coverage`, `--phase-coverage` or `--state-export`) stops before `load_db()`, so the message names `load_db` and the modes that were asked for. `--coverage` is implemented in #170; its empty-directory vectors still stop at the genuine world-initialization NYI, not at a coverage gate. Actual initialized-world dump C/C/Go and flag-free comparisons live in `turn-compat --full --command-coverage`.
- Any other vector (no arguments, `-r`, `-r -S`, `-S`, `-l`, `-T`, `-w`, `-p`, `-f`, `-G`, `-A`, `test-use-const-report-date`, `--`, `-`, an operand) loads the game. Each runs in an empty directory, so there is no `lib`, and `load_db()` takes the branch that builds a new world: Go prints what C prints up to `INIT: post_production()` (on stdout, `lib/randseed could not be opened.`; on stderr, the missing files and the players `check_db()` creates) and stops there with `post_production: not yet implemented`, where C goes on and aborts in `create_faery()`. These vectors stay `nyi`. On the pinned initialized bare world, `-r -S` completes and `turn-compat --full` verifies identical database files, raw stdout/stderr and exit status against C. Unported gameplay branches still stop explicitly rather than succeeding silently.

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
- Linux or macOS. G3 #345 makes the C engine own its BSD `oly_getopt()`:
  it does not permute, says `illegal option`, and names `argv[0]` as given.
  The Go port's existing BSD short-option parser is unchanged.

## Running

```sh
tools/olympia-cli-compat/run.sh
```

Set `C_CACHE` to an `olympia-roundtrip` cache to reuse its C build. The tool
requires the cache's `g3.rev` to equal the exact pin and an executable debug
engine; otherwise it rejects the cache. The default still builds a clean archive.

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
Each run's outputs are left under `tmp/olympia-cli-compat/runs/<number>-<label>/{c,go}`, and the summary is written to `tmp/olympia-cli-compat/results.tsv`.
The number starts at 1 and increases for each noncomment case, keeping labels such as `a` and `A` separate even on case-insensitive filesystems.
The summary keeps the original labels, without the numeric prefix.

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

Last run: 2026-10-27, Linux amd64, GCC 12.2.0, Go 1.26.4.

- 92 vectors: 62 `same` and 30 `nyi`. All 92 PASS. For the 15 `nyi` vectors that now load (#52), Go's stdout is all of C's and its stderr is C's up to the "not yet implemented" line.
- G3 gate: `rng-selftest` prints `YES` for the C build and for Go.
- G3's unchanged `cli-options` gate also prints `YES` for both engines
  (20 cases, including fixed `argv[0]` values), through the native timeout
  launcher of `olympia-spec` (#99). No Go option-parser change was needed.

The full table of hashes and exit statuses is in [`results.tsv`](results.tsv), copied from the run, not edited by hand.

The C stderr and merged hashes of the 30 `nyi` vectors, where C aborts (134), changed with [#47](https://github.com/mdhender/g3g0/issues/47).
Before it, the program ran as the shell's own child, and bash wrote its `Abort trap: 6` report, with the process ID and the script's line number and text, into the captured stderr, so those hashes differed on every run.
Now the time limit (`tools/lib/timeout.sh`) runs it under a supervisor that exits with the same status, and the captured stderr is only what C wrote.

## Known differences

- **Buffer size.** C flushes a full stdout buffer early, at a size that depends on the file; the port holds stdout until exit (or line by line on a terminal, or not at all after `-f`).
  Nothing ported yet writes enough to stdout for that to change how it interleaves with stderr.

## Upstream defects reproduced on purpose

None are known at the pinned revision.
