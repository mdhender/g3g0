# genesis differential test

`run.sh` checks that the Go world generator (`cmd/genesis`) is byte-compatible with G3's C `genesis-g3`.
G3 is the reference.
Expected results always come from running the pinned C build; nothing here is re-baselined from Go output.
Two sources of expectations are used:

- G3's own genesis regression gates (`tests/genesis/regress`, added for [G3 #290](https://github.com/playbymail/olympia-g3/issues/290)). They run unchanged from the pinned export, with G3's committed expectations.
- A downstream differential matrix, [`cases.txt`](cases.txt), for what the gates don't cover. Its expectation for each case is the pinned C build's output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `2322a0474de5fb38308d795f6cd72beb811652e9`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.

The port follows `genesis/genesis.c`, `genesis/pcg.c` (PCG32 and splitmix64), `genesis/names.h` (the name pools, copied verbatim into `internal/genesis/names.go`) and `lib/map_image.c` (the `--image` BMP writer).
`doc/genesis.md` describes the design; where it and the source disagree, the source wins.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use an independent clone:

  ```sh
  git clone https://github.com/playbymail/olympia-g3 /path/to/olympia-g3
  G3_REPO=/path/to/olympia-g3 tools/genesis-compat/run.sh
  ```

- What G3 needs to build `genesis-g3` and `mapgen-g3`: CMake 3.25 or later, Ninja and a C11 compiler.
- Go (the version in `go.mod`), Bash (3.2 is enough), Perl, `xxd`, `shasum`, `cmp` and `diff`.
- macOS. genesis's argument diagnostics come from the C library's `getopt_long`, and the port reproduces the BSD libc that macOS ships.
  glibc prints different messages, so on Linux the argument-error cases would fail.

## Running

```sh
tools/genesis-compat/run.sh
```

The script:

1. Exports the pinned revision with `git archive` into `tmp/genesis-compat/g3`.
   It uses neither the working tree nor an existing build of the clone.
   It builds `genesis-g3` and `mapgen-g3` there with the `debug` preset.
2. Builds the Go ports of both.
3. Runs every case in `cases.txt`.
   Each run gets a fresh, empty directory.
   Both programs get the same argument vector, including `argv[0]` (`genesis-g3`), because the usage message prints `argv[0]` and the `getopt_long` messages print its last element.
   The C and Go programs each run twice.
4. Checks the reference's exit status against the case's expectation: `ok` must exit 0 and write all five files, `help` must exit 0 and write nothing, and `fail` must exit nonzero.
5. Requires that both C runs and both Go runs agree, and that C and Go agree byte for byte on every file left in the directory (`Map`, `Regions`, `Land`, `Cities`, `randseed`, any BMP, and which files exist), on stdout, on stderr and on the exit status.
   On a mismatch it prints both SHA-256 hashes, the first differing byte, and the start of a unified diff.
6. Requires that cases in the same group, which differ only in their seed, give different reference worlds.
   A port that ignored the seed would fail this check.
7. Feeds every world the Go port generated to the C `mapgen-g3` and to the Go `mapgen`.
   The C program must exit 0, and the two must agree byte for byte on `loc`, `gate`, `road`, the final `randseed`, stdout, stderr and the exit status.
   This is a smoke check that the worlds are usable; the comparison against C genesis is the gate.
8. Runs G3's own genesis regression gates (`tests/genesis/regress/*/check.sh`) from the export, unchanged, against the C build and against Go.
   They run the C build by default, and the Go binary through their `GENESIS_BIN` setting (`tests/genesis/genesis.sh`).
   Each must print `YES`.

The script exits 0 only when every check passes.
Each run's outputs are left under `tmp/genesis-compat/runs/<case>/{c1,c2,go1,go2}`, the mapgen runs under `tmp/genesis-compat/mapgen/<case>/{c,go}`, and the summaries are written to `tmp/genesis-compat/results.tsv` (cases) and `tmp/genesis-compat/gates.tsv` (G3 gates).

`go test ./internal/genesis` checks the C library emulation (`atoi`, `strtoull`, `getopt_long`), the argument diagnostics, determinism, and seed sensitivity.
Its expected messages were observed from the pinned C build.

## Case matrix

`cases.txt` records each case's exact argument vector. It covers:

| area | cases |
|---|---|
| seeds | the default (seed 0); seeds 1, 12345, `0xdeadbeef`, `0777`, 2⁶⁴−1 and 9876543210123456789, which must all differ; `strtoull` base-0 spellings (`16`, `0x10`, `0X10`, `020`), `-1`, overflow, leading blanks and `+`, `0x` with no digits, `08`, trailing junk, and no digits |
| sizes | 10, 11, 12, 20, 36, 50, 98 and 99; an `atoi` value that wraps to 10; leading blanks and trailing junk |
| oceans | 1 and 12 oceans; a request clamped to the cap; seeds where 4-coloring fails and the count is reduced (sizes 13, 20 and 40) |
| islands | 2, 3 and 10 iterations, including 2 on a 10×10 map |
| naming | a world big enough that the region pool wraps and names get a numeric suffix |
| `--defer-oceans` | frontier growth, explicit and default; `bfs` growth; 1 and 12 oceans; sizes 10, 30, 60 and 99; lakes, and a pool name "Great Sea" recast as "Great Lake"; tiny landlocked water filled with land |
| start cities | 0, 1, 6 and 16; a 10×10 map where spacing reduces the count; 16 under `--defer-oceans` |
| `--image` | the defaults; `--pixels` 1, 2, 3, 5 and 64; `--grid` 1, 2 and 4; rows that need padding; a repeated `--image`; `--pixels` and `--grid` out of range without `--image`, which is not an error |
| `getopt_long` | attached short arguments, clusters, `--name=value`, unique abbreviations, a long option inside a short cluster (`-D-size`), permuted operands, a lone `-`, `--` ending the options, repeated options, an option argument of `--`, and `POSIXLY_CORRECT` |
| help | `-h`, `--help`, `--he`, `-hx`, and `--help` after an invalid `--size` |
| errors | every range check; an invalid `--ocean-growth`; `bfs` without `--defer-oceans`; unknown, ambiguous and malformed options; missing and unexpected arguments |
| output failures | a read-only directory; `Map`, `Regions`, `Land`, `Cities`, `randseed` or the BMP path is a directory; an empty `--image` path; a BMP in a missing directory |

Seeds for the rare paths (coloring retries, start-city spacing retries, lakes, "Great Lake" and the pool wrap) were found by scanning seeds with the pinned C build.

Two `--defer-oceans` retry paths are not exercised: "couldn't place N ocean origins" and "N deferred oceans not 4-colorable".
Scans of up to 400 seeds per configuration, at sizes 10 to 99 with up to 12 oceans and both growth modes, never reached them in the C build.
No case reaches the three `assert()`s, "no land for start cities" or "no open water left for oceans" either.
G3's gates don't reach any of these paths either; its `tests/genesis/README.md` lists them as unreached.

The G3 gates cover a subset of this matrix with G3's own expectations: defaults, seed parsing, sizes and the ocean cap, coloring retries, island iterations, the region pool wrap, deferred oceans, `bfs` growth, start cities, `--image`, validation errors and write failures.
They pass no unknown or malformed options, because `getopt_long`'s messages differ between C libraries, so the `getopt_long` cases here have no upstream counterpart.

## Results

Last run: 2026-10-26, macOS 26.6.2 arm64, Apple clang 21.0.0, CMake 4.3.3, Go 1.26.4.

- 131 cases: 84 `ok`, 5 `help` and 42 `fail`.
  All 131 PASS: Go matches C byte for byte on every file, stdout, stderr and exit status, and both programs are deterministic.
- The seed groups gave 7, 4 and 3 distinct reference worlds for 7, 4 and 3 cases.
- All 84 Go-generated worlds are accepted by the C `mapgen-g3`, and the Go `mapgen` matches it on each.
- G3 gates: all 12 (`coloring-retry`, `defaults`, `defer-oceans`, `image`, `island-iterations`, `ocean-growth-bfs`, `region-pool-wrap`, `seed-parsing`, `sizes`, `start-cities`, `validation-errors`, `write-failures`) print `YES` for the C build and for Go.

The full table of hashes is in [`results.tsv`](results.tsv), and the gate results are in [`gates.tsv`](gates.tsv). Both are copied from the run, not edited by hand.

## Known differences

These are limited to failure paths the C build cannot reproduce exactly. None is reached by any case:

- **Assertion failures.** `assert()` prints the compiler's `__FILE__`, which is the absolute path of the source on the build machine.
  The port prints `genesis/genesis.c` with the C line number, and exits 134 as `abort()` does.
- **Undefined behavior.** More than 100 water bodies under `--defer-oceans` would overrun fixed-size C arrays before the assertion that guards them.
  The port panics with an index error there instead.
- **Byte order.** `randseed` is two `uint64_t` values written with `fwrite`, so its bytes are in the host's byte order.
  The port writes them in Go's native byte order, which matches on every platform; the recorded results are from a little-endian machine.
  G3's gate manifests pin `randseed`, so, as its `tests/genesis/README.md` notes, they hold only on little-endian hosts.

## Upstream defects reproduced on purpose

None are known at the pinned revision.
When G3 has an unfixed defect, `WORKFLOW.md` requires the port to reproduce it, with a comment linking the G3 issue, and to list it here.

## Upstream changes followed

Each change moves the pinned revision to the G3 commit that merged it:

| G3 change | G3 commit | g3g0 issue | change |
|---|---|---|---|
| [#291](https://github.com/playbymail/olympia-g3/pull/291), for G3 [#290](https://github.com/playbymail/olympia-g3/issues/290) | `2322a04` | [#10](https://github.com/mdhender/g3g0/issues/10) | genesis regression gates with `GENESIS_BIN` selection; tests only, so the port is unchanged and `run.sh` now runs the gates against C and Go |
