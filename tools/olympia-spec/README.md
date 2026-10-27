# olympia spec against the Go engine

`run.sh` runs G3's olympia spec, unchanged, against the Go engine (`cmd/olympia`), and records which gates pass in [`passing.txt`](passing.txt).
G3 is the reference.
Expected output comes only from G3's committed files, which its C engine wrote; nothing here re-baselines.

The spec is the three gates that G3 documents for another engine in `doc/claude/testing.md` ("Running the spec against another engine"):

| Gate | What it checks |
|---|---|
| `tests/olympia/regress/run-all.sh` | every regress fixture's `check.sh` (not `check-lua.sh`); one `PASS`/`FAIL`/`ERROR` line per fixture |
| `run/olympia-g3.sh` | the full golden turn (`-i`, then `-r -S`) from `lib.tgz`, written to `OLYMPIA_RUN_DIR` |
| `tests/olympia/golden_check.sh` | that turn's `lib` against `tests/olympia/golden/manifest.sha256`; prints `YES` |

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `ab24d451ee25d61364e08cde1fbfc223e220fbd1`.
The revision is pinned in `run.sh` (`PINNED_REV`).
`G3_REV` runs another revision, but then `passing.txt` is neither compared nor written.
Change the pin only on purpose, and rerun this and [`tools/olympia-roundtrip`](../olympia-roundtrip) when you do.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another. Its working tree and build are not used.
- Go (the version in `go.mod`), Bash (3.2 is enough), Perl, Python 3 (`rng-selftest` runs `-R` on a pseudo-terminal), `patch`, `shasum` and `diff`.
- macOS, for the reasons given in [`tools/olympia-cli-compat`](../olympia-cli-compat/README.md#prerequisites).

No C build is needed: with `OLYMPIA_ENGINE_BIN` set, the gates run only that engine.

## Running

```sh
tools/olympia-spec/run.sh              # every fixture, the golden turn, and passing.txt
tools/olympia-spec/run.sh move-go ...  # just those fixtures; passing.txt is left alone
```

The script:

1. Exports the pinned revision with `git archive` into `tmp/olympia-spec/g3`.
2. Builds the Go port into `tmp/olympia-spec/bin/olympia-g3` and sets `OLYMPIA_ENGINE_BIN` to it.
3. Sets `OLYMPIA_RUN_DIR` to `${TMPDIR:-/tmp}/g3g0-olympia-spec-run` (or `RUN_DIR`), which must be outside this repository and the G3 clone, and empties it.
4. Runs `run-all.sh` from the export.
   With fixture names, only those run, and the golden turn is skipped.
5. On a full run, runs `run/olympia-g3.sh` and then `golden_check.sh`.
   `golden` passes when the turn exits 0 and the check prints `YES`.
6. Compares the passing gates with `passing.txt`, prints those newly passing and those no longer passing, and rewrites the file.
   If a recorded gate no longer passes, it exits 1 and leaves the file unchanged; `FORCE_RECORD=1` records the loss anyway.

Logs are kept in `tmp/olympia-spec`: `run-all.log` (each failing fixture's own output, indented), `olympia-g3.log` and `golden.log`.
The golden turn's directory is kept in the run dir until the next run.

`--update` is refused.
`run-all.sh` never re-baselines, and every `check.sh --update` and `golden_check.sh --update` refuses to run while `OLYMPIA_ENGINE_BIN` is set.
The gates also run inside the export, so they cannot write to the clone.

## Results

Last run: 2026-10-26, macOS 26.6.2 arm64, Go 1.26.4. About 45 seconds.

- `run-all.sh`: pass 1, fail 100, error 0. Only `rng-selftest` passes; every other fixture needs `load_db` and the turn.
- `run/olympia-g3.sh` exits 2 (the Go engine stops at `-i` with "not yet implemented"), so `golden_check.sh` reports that the run diverges from the manifest.
- The same `run-all.sh` with the pinned C build passes all 101 fixtures.
