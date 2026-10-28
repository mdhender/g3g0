# island differential test

`run.sh` checks that the Go island generator (`cmd/island`) is byte-compatible with G3's C `island-g3`.
G3 is the reference.
Expected results always come from running the pinned C build; nothing here is re-baselined from Go output.
Two sources of expectations are used:

- G3's own island regression gates (`tests/island/regress`, added for [G3 #292](https://github.com/playbymail/olympia-g3/issues/292)). They run unchanged from the pinned export, with G3's committed expectations.
- A downstream differential matrix, [`cases.txt`](cases.txt), for what the gates don't cover. Its expectation for each case is the pinned C build's output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `6bbcd63055d051ade6f745f368281628441a2c27`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.

The port follows `island/island.c` and the RNG in `lib/rnd.c`, which mapgen shares (`internal/rnd`).
The C library behavior it depends on (`atoi`, `getopt`, `fgets` line splitting, `perror`) is in `internal/clib`.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use an independent clone:

  ```sh
  git clone https://github.com/playbymail/olympia-g3 /path/to/olympia-g3
  G3_REPO=/path/to/olympia-g3 tools/island-compat/run.sh
  ```

- What G3 needs to build `island-g3` and `genesis-g3`: CMake 3.25 or later, Ninja and a C11 compiler.
- Go (the version in `go.mod`), Bash (3.2 is enough), Perl, `xxd`, `shasum`, `cmp` and `diff`.
- macOS. island parses its flags with the C library's `getopt()`, and the port reproduces the BSD libc that macOS ships, which does not permute arguments.
  glibc's `getopt()` permutes, so on Linux the cases with an operand before an option would fail.

## Running

```sh
tools/island-compat/run.sh
```

The script:

1. Exports the pinned revision with `git archive` into `tmp/island-compat/g3`.
   It uses neither the working tree nor an existing build of the clone.
   It builds `island-g3` and `genesis-g3` there with the `debug` preset.
2. Builds the Go port.
3. Creates each input map in `cases.txt`.
   They come from the upstream fixtures (`tests/island/fixtures`), from the small maps in [`maps/`](maps), from maps the script generates (one byte repeated over a rectangle, optionally with CRLF line ends), or from the `Map` that the pinned `genesis-g3` writes.
4. Runs every case.
   A case is a pipeline of one or more island runs from one map and one starting `randseed`.
   Each run gets a fresh directory with its own copy of the starting seed, so the original seed is never the file a run advances.
   The stages run one after another, each reading the stage before's stdout; that is what a pipe does, because island reads stdin to EOF before it loads the seed and saves the seed before it prints.
   Both programs run with `argv[0]` set to `island-g3`, and each runs every case twice.
5. Checks the reference's exit status against the case's expectation: `ok` must exit 0 at every stage, and `fail` must exit nonzero at the last stage.
6. Requires that both C runs and both Go runs agree, and that C and Go agree byte for byte on every stage's stdout, stderr and exit status, and on the final `randseed` (or its absence).
   On a mismatch it prints both SHA-256 hashes, the first differing byte, and the start of a unified diff.
7. Requires that cases in the same group, which differ only in their seed, give different reference output.
   A port that ignored the seed would fail this check.
8. Runs G3's own island regression gates (`tests/island/regress/*/check.sh`) from the export, unchanged, against the C build and against Go.
   `ISLAND_BIN` selects a native launcher for either binary, preserving the
   gate's `exec -a island-g3` name through the shared `run_limited` supervisor.
   Each must print `YES`.

All engine calls, including genesis input generation and both engines in the
upstream gates, use `tools/lib/timeout.sh` (`RUN_TIMEOUT`, default 60 seconds).
Timeout is status 124; ordinary signal exits are unchanged. No upstream source,
fixture or expected byte is modified by the launcher.

The script exits 0 only when every check passes.
Each run's outputs are left under `tmp/island-compat/runs/<case>/{c1,c2,go1,go2}`, and the summaries are written to `tmp/island-compat/results.tsv` (cases) and `tmp/island-compat/gates.tsv` (G3 gates).

`go test ./internal/island` checks a few cases against values observed from the pinned C build, and the error paths.
`go test ./internal/clib` checks the `getopt()` emulation.

## Case matrix

`cases.txt` records each case's map, exact seed bytes and argument vectors. It covers:

| area | cases |
|---|---|
| seeds | the four upstream fixture seeds and four more 16-byte seeds, which must all give different maps on the same input; a short seed file, an empty one, one longer than 16 bytes, a missing one (alone and in a chain), a `randseed` directory, and a read-only `randseed` |
| `-s` | 1, 2, 50 and 300; 0, -5 and `INT_MIN`, which are the random size; `atoi` with leading blanks, trailing junk, no digits, 32-bit wrap and `long` saturation; more than the map has room for; `INT_MAX - 1`, which allocates and runs out of room, and `INT_MAX`, whose allocation fails |
| `-b`, `-c` | border 0, 1, 10, 19, 20, 40, 41 and 1000 on a 40×40 map, and 7, 8, 15 and 16 on a 30×15 map; a negative border; shelf 0, 1, 6, 8, 13, 20 and 29 around land, then 30, 100000000 and `INT_MAX`, which cover all the water; a negative shelf, and 20 with no land |
| land and water | upstream's `land-40.map` with three seeds; no water left for an island, from an all-land map, a 1×1 map, a border reaching the middle of the map, and a shelf covering all the water; worlds generated by `genesis-g3` at sizes 30, 40, 60 and 99, including `--defer-oceans`: with no shelf or border, since a genesis world's water is mostly within the default shelf of land, and once with the defaults, which leave no room |
| shapes and glyphs | 1×1, 1×30, 30×1, 8×8, 30×15 and 80×80 maps; ragged rows; each ocean glyph mapgen accepts; a byte above 0x7f; a mixed map of ocean glyphs, terrain letters, `#`, a tab and a high byte; CRLF lines; blank LF and CRLF lines; `\r` alone as the line end; lines starting with NUL at the start, middle and end, and after a skipped row; NUL lines only; no final newline |
| input bounds | a 119×120 map, the largest accepted; a 120-character row; 121 rows; 250-character lines, which `fgets` splits; an empty file; blank lines only |
| chains | four stages; the same sizes in increasing and decreasing order; a chain that fills a small map; a chain whose last stage fails |
| `getopt()` | attached arguments, an argument in a cluster (`-b1c1`), repeated options, an argument that looks like an option, an operand before and after the options (the first stops option parsing), `--`, a lone `-`, an empty operand, unknown options before and after valid ones, `-:`, `--s`, missing arguments, and `-h` |

## Results

Final #166 Linux run, 2026-10-28, GCC 12.2.0, Go 1.26.4, unchanged pin:
**114/119 matching cases**, five case differences / seven failed checks,
exit **1**, not passing parity. `size-saturate` and `size-int-max-1` expose
the disclosed huge-allocation policy; `opt-operand-first`, `opt-lone-dash`
and `opt-empty-arg` expose glibc/BSD option permutation. All **17 upstream
gates pass for both engines**. Expectations are unchanged.

The #166 bounded-launch integration covers the previously direct genesis
input generation and gate launches. Repeating the complete tool with that
launcher gives byte-identical before/after `results.tsv` and `gates.tsv`,
including all five known differences. Independent launcher controls preserve
`argv[0] = island-g3`, statuses **134/139**, timeout **124**, and missing-target
failure **127**, without fabricated successful output.
Command: `G3_REPO="$PWD/tmp/upstream.git" tools/island-compat/run.sh`.
Evidence: `tmp/issue166-island-{before,default}.log`,
`tmp/island-compat/{results,gates}.tsv` and
`tmp/issue166-launcher-controls.log`. Matching captures were removed;
the five failing cases remain for inspection.

Historical reference-platform run: 2026-10-26, macOS 26.6.2 arm64,
Apple clang 21.0.0, CMake 4.3.3, Go 1.26.4.

- 119 cases: 87 `ok` and 32 `fail`.
  All 119 PASS: Go matches C byte for byte on every stage's stdout, stderr and exit status and on the final `randseed`, and both programs are deterministic.
- The seed groups gave 8, 4 and 3 distinct reference outputs for 8, 4 and 3 cases.
- G3 gates: all 17 (`border-shelf`, `chain`, `defaults`, `empty-input`, `existing-land`, `map-bounds`, `map-shape`, `missing-randseed`, `no-room`, `no-water`, `nul-lines`, `ocean-glyph-output`, `ocean-glyphs`, `shelf-large`, `size`, `usage-errors`, `write-failure`) print `YES` for the C build and for Go.

The full table of hashes and final seeds is in [`results.tsv`](results.tsv), and the gate results are in [`gates.tsv`](gates.tsv). Both are copied from the run, not edited by hand.

## Known differences

The port models one C library behavior instead of reproducing it:

- **Allocation.** `malloc()` for the island's square list fails only when `target_size + 1` overflows; on macOS any smaller request, up to 16 GiB, is granted without being touched.
  The port reports the failure for `-s 2147483647` only, and allocates no more squares than the map has.
  The `size-int-max` and `size-int-max-1` cases check both sides of that limit against C.

## Upstream defects reproduced on purpose

None are known at the pinned revision.
When G3 has an unfixed defect, `WORKFLOW.md` requires the port to reproduce it, with a comment linking the G3 issue, and to list it here.

## Upstream changes followed

The port started at `6cbf9a3`, after the first six of these changes.

| G3 change | G3 commit | g3g0 issue | change |
|---|---|---|---|
| [#293](https://github.com/playbymail/olympia-g3/pull/293), for G3 [#292](https://github.com/playbymail/olympia-g3/issues/292) | `a5f7825` | [#11](https://github.com/mdhender/g3g0/issues/11) | island regression gates with `ISLAND_BIN` selection |
| [#296](https://github.com/playbymail/olympia-g3/pull/296) (fixes #294) | `0bc2460` | [#13](https://github.com/mdhender/g3g0/issues/13) | any size below 1 is the random size, counted from 0 |
| [#297](https://github.com/playbymail/olympia-g3/pull/297) (fixes #295) | `c105f0b` | [#14](https://github.com/mdhender/g3g0/issues/14) | an empty map is an error: `No map on stdin!`, exit 1 |
| [#301](https://github.com/playbymail/olympia-g3/pull/301) (fixes #298, #299, #300) | `12383aa` | [#15](https://github.com/mdhender/g3g0/issues/15) | a missing `randseed` is reported on stderr; maps over 119 columns or 120 rows are rejected; the usage line names `island-g3` |
| [#304](https://github.com/playbymail/olympia-g3/pull/304) (fixes #302, #303) | `d33b3c6` | [#16](https://github.com/mdhender/g3g0/issues/16) | a border wider than the map stops at its edge; a line starting with NUL is skipped |
| [#306](https://github.com/playbymail/olympia-g3/pull/306) (fixes #305) | `6cbf9a3` | [#17](https://github.com/mdhender/g3g0/issues/17) | with no water left outside the shelf and border, island prints `No room for an island!` and exits 1 instead of overwriting a land square |
| [#308](https://github.com/playbymail/olympia-g3/pull/308) (fixes #307) | `6bbcd63` | [#18](https://github.com/mdhender/g3g0/issues/18) | the continental shelf is marked in one breadth-first pass instead of a 4^shelf recursion; output is unchanged, and a large `-c` now finishes |
