# olympia load/save round trip

`run.sh` checks that the Go engine (`cmd/olympia`) loads and saves a game exactly as G3's C `olympia-g3` does, with no turn in between.
G3 is the reference.
Expected results always come from running the pinned C build on the same input; nothing here is re-baselined from Go output, and `--update` is refused.

This is the exit check for stage 0 (#20): load and save must match on every case before the turn is ported.
Since `load_db` and `save_db` were ported (#28), every case passes.

## Cases

A case is a saved game.
The games are:

- `world`: G3's shared base world, `tests/olympia/fixtures/world.tgz`.
- Every regress fixture with a saved game, named after its directory under `tests/olympia/regress`: `scenario.patch` over `world.tgz`, unpacked with G3's own `scenario_unpack` (`tests/olympia/fixtures/scenario.sh`), or the fixture's `scenario.tgz` (`trade-routes`).
  `mint-top-id`, `rng-selftest` and `save-seed-unwritable` have no saved game of their own.
- Each of those after one C turn, named `<game>+turn`.
  The turn is `olympia-g3 -r -l ./lib -S test-use-const-report-date`, as the fixtures run it.

At the pinned revision that is 98 games and 196 cases.

Each case runs, with both engines, on identical copies of the game:

```sh
olympia-g3 -l ./lib -S test-use-const-report-date
```

with `argv[0]` set to `olympia-g3` and stdin from `/dev/null`.
`test-use-const-report-date` fixes the masthead date that `save_db` writes into `times_0`; without it, the C results would depend on the day they were made.

A case passes when the exit status, stdout and stderr match byte for byte, and so does the whole run directory: every file under it, and no file on only one side.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `47535543d83fbc5649d54b63b49e578f2b11ab1c`.
The revision is pinned in `run.sh` (`PINNED_REV`); `G3_REV` tests against another.
Change the pin only on purpose, together with [`tools/olympia-spec`](../olympia-spec).

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another. Its working tree and build are not used.
- What G3 needs to build `olympia-g3`: CMake 3.25 or later, Ninja and a C11 compiler.
- Go (the version in `go.mod`), Bash (3.2 is enough), Perl, `patch`, `tar`, `shasum`, `cmp` and `diff`.
- macOS, for the reasons given in [`tools/olympia-cli-compat`](../olympia-cli-compat/README.md#prerequisites).

## Running

```sh
tools/olympia-roundtrip/run.sh                    # every case
tools/olympia-roundtrip/run.sh world move-go      # just these cases
tools/olympia-roundtrip/run.sh 'move-*' '*+turn'  # shell patterns on case names
tools/olympia-roundtrip/run.sh --list             # print the case names
tools/olympia-roundtrip/run.sh --fresh            # drop the cached C side first
tools/olympia-roundtrip/run.sh --determinism-all  # check C's determinism on every case
```

The script:

1. Exports the revision with `git archive` into `tmp/olympia-roundtrip/g3` and builds `olympia-g3` there with the `debug` preset.
2. Builds the Go port into `tmp/olympia-roundtrip/bin/olympia-g3`.
3. For each selected case, makes the input and the C result unless they are cached.
   On a fixed sample of cases (`world`, `move-go`, `trade-routes` and each `+turn`; `DETERMINISM_SAMPLE` in `run.sh`), C runs the round trip twice, and the case is an `ERROR` if its two runs differ.
   `--determinism-all` checks every selected case, including those already cached without the check.
   A case is also an `ERROR` if the game cannot be unpacked or the C turn fails.
4. Runs Go on a fresh copy of the input and compares.
   It prints one line per case: `PASS`, or `FAIL` with what differs (the exit statuses, `stdout`, `stderr`, and the number of files that differ with the first five), then a `pass N fail N error N` line.

The export, the C build, the inputs and the C results are kept and reused while the revision, the C binary and the way they are made stay the same.
Only Go is rebuilt and rerun, so a later run takes about 25 seconds for every case, or under a second for a few.
The first run, or one after `--fresh`, takes about 3 minutes.

Each case's files are under `tmp/olympia-roundtrip/cases/<case>`:

| Path | What it is |
|---|---|
| `input/` | the saved game both engines start from |
| `turn/` | for a `+turn` case, the C turn's `stdout`, `stderr` and `status` |
| `c/`, `go/` | each engine's run directory (`cwd/`) and its `stdout`, `stderr` and `status` |
| `diff.txt` | the statuses, then `diff -u` of stdout and stderr, then `diff -ru` of the run directories |
| `error` | why the case could not be prepared |
| `deterministic` | C's round trip was run twice and gave the same result |

`tmp/olympia-roundtrip/results.tsv` summarizes the last run, one row per case it ran.
It exits 0 only when every selected case passes.

## Results

Last run: 2026-10-27, macOS 26.6.2 arm64, Apple clang 21.0.0, CMake 4.3.3, Go 1.26.4, with `--fresh --determinism-all`.

pass 196, fail 0, error 0.
C gave the same result twice on every case, and Go matches it on every case: exit status, stdout, stderr and every file of the run directory.
