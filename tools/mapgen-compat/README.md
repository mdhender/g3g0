# mapgen differential test

`run.sh` checks that the Go map generator (`cmd/mapgen`) is byte-compatible with G3's C `mapgen-g3`.
G3 is the reference.
Expected results always come from running the pinned C build; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `dcdb868d44707f4cff0eef5702224948fc7fd3b7`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.

The port follows `mapgen/mapgen.c` and `mapgen/z.c` and the linked support code: `lib/rnd.c` (MD5 RNG), `lib/ilist.c`, `lib/tiles.c` and `lib/roads.c` (lists and Knuth shuffles), and `lib/checked_alloc.c` (zeroed allocation).

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use an independent clone:

  ```sh
  git clone https://github.com/playbymail/olympia-g3 /path/to/olympia-g3
  G3_REPO=/path/to/olympia-g3 tools/mapgen-compat/run.sh
  ```

- What G3 needs to build `mapgen-g3` and `genesis-g3`: CMake 3.25 or later, Ninja and a C11 compiler.
- Go (the version in `go.mod`), Bash, Perl, `xxd`, `shasum`, `cmp` and `diff`.

## Running

```sh
tools/mapgen-compat/run.sh
```

The script:

1. Exports the pinned revision with `git archive` into `tmp/mapgen-compat/g3`.
   It uses neither the working tree nor an existing build of the clone.
   It builds `mapgen-g3` and `genesis-g3` there with the `debug` preset.
2. Builds the Go port.
3. Creates each input set in `inputs.txt`.
   They come from the upstream fixtures (`tests/mapgen/fixtures`), from the inputs of upstream regressions (`tests/mapgen/regress/<name>`), or from the pinned `genesis-g3`.
4. Runs every input set with every seed in `seeds.txt`.
   Each run gets a fresh directory with copies of `Map`, `Regions`, `Land`, `Cities` and the original seed bytes.
   The original seed is never the file that a run advances.
   The C and Go programs each run twice.
5. Requires that the C program exits 0, and that both C runs and both Go runs agree.
   It also requires the C and Go runs to agree byte for byte on `loc`, `gate`, `road`, the final `randseed`, stdout, stderr and the exit status.
   On a mismatch it prints both SHA-256 hashes, the first differing byte, and the start of a unified diff.
6. For each input set, requires that every full 16-byte seed gives a different reference map.
   A port that ignored the seed would fail this check.
7. Runs G3's own mapgen regression scripts (`tests/mapgen/regress/*/check.sh`) from the export, unchanged, against the C build and against Go.
   They run the C build by default, and the Go binary through their `MAPGEN_BIN` setting (`tests/mapgen/mapgen.sh`).

The script exits 0 only when every check passes.
Each run's outputs are left under `tmp/mapgen-compat/runs/<input>/<seed>/{c1,c2,go1,go2}`, and the summary is written to `tmp/mapgen-compat/results.tsv`.

`go test ./internal/mapgen` runs the `secret-sea-route` and `save-seed-unwritable` regressions in-process when `G3_ROOT` names a G3 checkout.
They read the seed and `EXPECT.sha256` from that checkout:

```sh
G3_ROOT=tmp/mapgen-compat/g3 go test -count=1 ./internal/mapgen
```

## Seed corpus

`seeds.txt` records the exact seed bytes, in hex:

| label | bytes | why |
|---|---|---|
| fixture | `736c6172746962617274666173743f0a` | upstream fixture seed (`slartibartfast?\n`) |
| secret-sea-route | `53303030303030303030303030303032` | upstream `secret-sea-route` reproducer (`S000000000000002`) |
| zero | `00000000000000000000000000000000` | all-zero state |
| ones | `ffffffffffffffffffffffffffffffff` | all-ones state |
| random-1 … random-4 | see `seeds.txt` | arbitrary states from `openssl rand -hex 16` |
| short-5 | `68656c6c6f` | short file: `load_seed()` reads 5 bytes, rest stay zero |
| missing | (no file) | `load_seed()` prints `randseed could not be opened.` to stdout and starts from zero |

`missing` intentionally produces the same map as `zero`, so the distinctness check covers only the 16-byte seeds.

## Results

Last run: 2026-10-26, macOS 26.6.2 arm64, Apple clang 21.0.0, CMake 4.3.3, Go 1.26.4.

- 9 input sets × 10 seeds = 90 cases.
  All 90 PASS: Go matches C byte for byte on every output, and both programs are deterministic.
- Each input set gives 8 distinct reference maps for its 8 full seeds.
- The `fixtures`/`secret-sea-route` case reproduces upstream's `secret-sea-route/EXPECT.sha256` for `gate`, `loc` and `road`.
- G3 regressions: every `tests/mapgen/regress/*/check.sh` at the pinned revision prints `YES` for the C build and for Go.

The full table of hashes and final seeds is in [`results.tsv`](results.tsv). It is copied from the run, not edited by hand.

`#` map holes, the Uldim and Summerbridge glyphs (`^ v { } [ ]`), sea-lane glyphs (`; : ~ "`) and `*` city squares come from upstream's `map-glyphs` regression inputs, added for [G3 #276](https://github.com/playbymail/olympia-g3/issues/276), which run with every seed in the corpus.

## Known differences

These are limited to failure paths that the C build cannot reproduce exactly:

- **Assertion failures.** `asfail()` prints the compiler's `__FILE__`, which is the absolute path of the source on the build machine.
  The port prints `mapgen/mapgen.c` with the C line number, and exits 134 as `abort()` does.
- **Crashes.** Where the C code dereferences NULL, the port exits 139 and prints nothing, which is what the shell sees for SIGSEGV.
  The process is not actually killed by a signal.
  Out-of-range input, such as a map wider or taller than 100 squares, is undefined behavior in C; the port panics there instead.

## Upstream defects reproduced on purpose

As `WORKFLOW.md` requires, the port reproduces these and links each one from `internal/mapgen/mapgen.go`:

- [G3 #285](https://github.com/playbymail/olympia-g3/issues/285): `bridge_map_holes()` skips holes in the last row and column.
- [G3 #287](https://github.com/playbymail/olympia-g3/issues/287): a hole at (`max_row`, `max_col`) crashes the final "highest province" line.

## Upstream fixes followed

Each fix moves the pinned revision to the G3 commit that merged it:

| G3 fix | G3 commit | g3g0 issue | change |
|---|---|---|---|
| [#277](https://github.com/playbymail/olympia-g3/pull/277) (fixes #274) | `06a5959` | [#2](https://github.com/mdhender/g3g0/issues/2) | `bridge_map_hole_sup()` checks each neighbor for NULL before reading `mark` |
| [#282](https://github.com/playbymail/olympia-g3/pull/282) (fixes #272) | `b086a16` | [#3](https://github.com/mdhender/g3g0/issues/3) | `random_province_gates()` gates to the second province, `(r2,c2)`; changes `gate` for every seed |
| [#283](https://github.com/playbymail/olympia-g3/pull/283) (fixes #273, #278) | `8854d66` | [#4](https://github.com/mdhender/g3g0/issues/4) | unnamed regions are allowed (stone-circle line prints `??`); `Regions` and `Land` entries with no province are skipped with a diagnostic |
| [#284](https://github.com/playbymail/olympia-g3/pull/284) (fixes #275) | `635bb93` | [#5](https://github.com/mdhender/g3g0/issues/5) | `random_province()` and `mark_bad_locs()` scan the last map column (the port has no `not_random_province()`, which nothing calls) |
| [#286](https://github.com/playbymail/olympia-g3/pull/286) (fixes #281) | `dcdb868` | [#6](https://github.com/mdhender/g3g0/issues/6) | `flood_land_inside()` treats another region's boundary square as a wall instead of asserting |
