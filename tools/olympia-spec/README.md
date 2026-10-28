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

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `e733acb63bacb21167a1f2589d8dcdc023364896`.
The revision is pinned in `run.sh` (`PINNED_REV`).
`G3_REV` runs another revision, but then `passing.txt` is neither compared nor written.
Change the pin only on purpose, and rerun this and [`tools/olympia-roundtrip`](../olympia-roundtrip) when you do.

At #131, the exact recorded 19-pass set is retained: 19 pass / 104 fail /
0 error at the new pin. The generated passing.txt changes provenance only.
Golden initialization still stops at explicitly unsupported `-i`; this
no-loss progress check does not claim later commands/modes or whole-spec parity.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another. Its working tree and build are not used.
- Go (the version in `go.mod`), Bash (3.2 is enough), Perl, Python 3 (`rng-selftest` runs `-R` on a pseudo-terminal), `patch`, `shasum` and `diff`.
- Linux or macOS (the pinned C engine owns its BSD short-option parser).

No C build is needed: with `OLYMPIA_ENGINE_BIN` set, the gates run only that engine.

## Running

```sh
tools/olympia-spec/run.sh              # every fixture, the golden turn, and passing.txt
tools/olympia-spec/run.sh move-go ...  # just those fixtures; passing.txt is left alone
```

The script:

1. Exports the pinned revision with `git archive` into `tmp/olympia-spec/g3`.
2. Builds the Go port into `tmp/olympia-spec/bin/olympia-g3` and sets `OLYMPIA_ENGINE_BIN` to a native launcher that sends it through `run_limited`, preserving the fixture's `exec -a` `argv[0]` (#99). A shell-script launcher discards that name before it can reach the engine.
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

Final #166 run: 2026-10-28, Linux amd64, Go 1.26.4, pinned revision above.
`G3_REPO="$PWD/tmp/upstream.git" tools/olympia-spec/run.sh`:

- Regress: **pass 21, fail 102, error 0** (70 seconds, 123 fixtures).
  All **19** baseline passes remain. The unchanged upstream checks now also
  pass `mint-top-id` and `token-vanish`; the tool records just those two new
  names in `passing.txt`. No fixture, manifest or expected bytes changed.
- The separate strict whole-world C/C/Go comparison and all 13 lifecycle
  controls pass, including report artifacts. Remaining spec failures still
  request unsupported coverage/state exports, commands and modes; they are
  not accepted successful gameplay or a claim of whole-spec parity.
- Golden initialization still exits 2 at `load_db (-i): not yet implemented`;
  golden_check reports `NO`. The spec tool exits 0 because no recorded pass was
  lost, **not** because all regressions or the golden flow pass.
- Logs: `tmp/issue166-spec.log`, `tmp/issue166-spec-audit.log`, plus
  `tmp/olympia-spec/{run-all,olympia-g3,golden}.log`. Full C spec was not rerun.
  Every Go engine invocation uses the existing native `run_limited` launcher.

Historical foundation run (before stage-1 turn integration): 2026-10-27,
Linux amd64, Go 1.26.4. About 54 seconds. The limitations below describe that
run, not the current implementation.

- `run-all.sh`: pass 15, fail 100, error 1; no previously passing gate lost.
  - G3 #345's `cli-options` (20 cases) and `load-number-width` (13 cases)
    now pass. Both unchanged gates also print `YES` for the actual pinned C
    engine. Generic Scanner and the Go BSD option parser are unchanged.
  - Seven newly passing G3 #340/#341 regressions: `assert-flush`, `load-bad-master`, `load-box-usage`, `load-here-duplicate`, `load-orphan-unit`, `load-stray-orders`, `trade-route-seed-parse`. These also pass on the actual pinned C engine.
  - Passing: `rng-selftest`, and, since `load_db`, `save_db` and `-l`/`-S` were ported (#52), the fixtures that only load and save: `load-comment-lines`, `load-continuation-blank`, `load-continuation-eof`, `load-short-entries` and `save-fly-cap`.
  - 97 fixtures fail because their turn (`-r`) is not ported; `save-seed-unwritable` is one of them (it wants the turn to exit 1 on an unwritable `randseed`, and Go stops at `-r` before loading). `trade-routes` fails for the same reason, through `--use-coverage`, which its check reads back. `mail-libdir-space` needs `-M`, `-h` and the turn.
  - `save-libdir-space` (G3 #320) is an `ERROR`, not a `FAIL`: it runs the engine for its reference first, a `-r -S` turn, and an engine that cannot run the turn is an error to it.
- `run/olympia-g3.sh` exits 2 (the Go engine stops at `-i` with "not yet implemented"), so `golden_check.sh` reports that the run diverges from the manifest.
- Full C spec not rerun at this pin; the two G3 #345 regressions pass on C.
