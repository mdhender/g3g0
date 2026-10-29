# olympia load/save round trip

`run.sh` checks that the Go engine (`cmd/olympia`) loads and saves a game exactly as G3's C `olympia-g3` does, with no turn in between.
G3 is the reference.
Expected results always come from running the pinned C build on the same input; nothing here is re-baselined from Go output, and `--update` is refused.

This is the exit check for stage 0 (#20): load and save must match on every case before the turn is ported.
Since `load_db` and `save_db` were ported (#28), every case passes.

## Final-tier replay (#178)

At **0.115.1-alpha**, the unchanged tool reuses the verified exact-pin
`tmp/issue179-roundtrip` cache, without --fresh, and passes **201 / 0 fail /
0 error / 1 separate EXPECTED_ABORT**, **144s**:

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue179-roundtrip" \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  tools/olympia-roundtrip/run.sh --determinism-all
```

Before reuse, all **1212** original archived files equal
**f73e48118805ffe8ff5e763bfa4f20ae00c3a705**, C binary SHA-256
`151c7355890eba7f8bfda1f4f16f4a490947aaa1db1ffcc3ca859cd659c01158`.
All **201** input names/types/modes/bytes equal the preserved old stage3
corpus: **41938 files / 5228 directories**; no symlinks or manufactured
`release-swear-np+turn`. Full input-manifest SHA-256:
`a1207a76ba2a7374d339a22eec21d9ade625fc1e309c0d540c393d1cbc8f1f91`.

Supplementary replay additionally runs **fresh unchanged C twice then Go on
all 201 inputs**, each on a real disposable copy, through
`run_limited 120 olympia-g3 BIN -l ./lib -S test-use-const-report-date`.
The existing strict turn comparator plus tree-modes.py compare C/C, C/Go
and cached-C/fresh-C, including unfiltered status/streams and every cwd
entry/type/mode/empty directory/byte. Every comparison matches. The separate
fresh C/C `-r` control again gives exact **134 / swear.c (180): FALSE**;
it supplies no successful post-turn input. **36** marker-injection controls
using actual status-0/134 captures reject invalid references before status
acceptance. Matching captures and temporary scripts/probes are removed.

Evidence: `tmp/issue178-{integrity-before,roundtrip,roundtrip-strict,marker-controls}.log`,
`tmp/issue178-input-manifest.txt`, `tmp/issue178-roundtrip-results.tsv`
(SHA-256 `96cafa5021e93976c9e1c0de550fb0a100aa7003251c5570e57595f090adb2e2`),
`tmp/issue178-roundtrip-strict/*/{manifest.txt,*.diff}` and separate abort
raw streams/status/manifest. Parent also completed the initially withheld
full-box gate using physical disposable input copies via the existing GAMES
override: **202 worlds / 8,146,800 boxes, plain and pretty, PASS**. No harness
fix or exemption was needed. Parent independently verified all **949**
evidence digests and protected-input/source integrity, and repeated bounded
vet/tests, the strict turn and spec no-loss checks. See [accepted stage-3
evidence and the copy workflow](../README.md#final-tier-validation-178).
Stage 3 is accepted; paused without starting stage 4.

## Cases

A case is a saved game.
The games are:

- `world`: G3's shared base world, `tests/olympia/fixtures/world.tgz`.
- Every regress fixture with a saved game, named after its directory under `tests/olympia/regress`: `scenario.patch` over `world.tgz`, unpacked with G3's own `scenario_unpack` (`tests/olympia/fixtures/scenario.sh`), or the fixture's `scenario.tgz` (`trade-routes`).
  `mint-top-id`, `rng-selftest`, `save-seed-unwritable`, `load-trade-no-price`
  and `summary-trunc-name` have no saved game of their own. The last two are
  G3 #360 check-only fixtures: they mutate fresh `trade-market` copies for
  four real turns each. They do not add ordinary or `+turn` corpus names.
- Each of those after one successful C turn, named `<game>+turn`.
  The turn is `olympia-g3 -r -l ./lib -S test-use-const-report-date`, as the fixtures run it.

At the pinned revision that is 101 games and **201 saved-game cases**:
101 ordinary games and 100 successful post-turn snapshots. All 196 inherited
cases remain, plus `savage-cap` and `token-vanish` before/after a turn and
the ordinary `release-swear-np` load/save case.

`release-swear-np+turn` is **unavailable**, not a saved-game case. The unchanged
pinned [check.sh](https://github.com/playbymail/olympia-g3/blob/f73e48118805ffe8ff5e763bfa4f20ae00c3a705/tests/olympia/regress/release-swear-np/check.sh)
and [README.md](https://github.com/playbymail/olympia-g3/blob/f73e48118805ffe8ff5e763bfa4f20ae00c3a705/tests/olympia/regress/release-swear-np/README.md)
require its real C turn to abort with status 134 and the `swear.c` FALSE
assertion, without sanitizer errors (G3 #357). No successful save exists.
The harness explicitly checks that contract twice, requiring the exact pinned
line 180, rejecting missing-dependency/unsupported markers, and comparing raw
streams, status and partial directories without normalization. It reports
`EXPECTED_ABORT` separately from load/save `PASS`, never moves failed output
into an `input/`, and does not run a Go turn on that later-command fixture.
Any other failed C turn is still an error, including an unrelated status 134.

Each case runs, with both engines, on identical copies of the game:

```sh
olympia-g3 -l ./lib -S test-use-const-report-date
```

with `argv[0]` set to `olympia-g3` and stdin from `/dev/null`.
`test-use-const-report-date` fixes the masthead date that `save_db` writes into `times_0`; without it, the C results would depend on the day they were made.

A case passes when the exit status, stdout and stderr match byte for byte, and so does the whole run directory: every file under it, and no file on only one side.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `f73e48118805ffe8ff5e763bfa4f20ae00c3a705`.
The revision is pinned in `run.sh` (`PINNED_REV`); `G3_REV` tests against another.
Change the pin only on purpose, together with [`tools/olympia-spec`](../olympia-spec).

The expected-abort exception applies **only to this exact resolved revision**.
`G3_REV` overrides at any other revision retain all ordinary `+turn` candidates
and require status 0. In particular, valid older snapshots are not silently
dropped. A future pin needs its own source review; arbitrary nonzero turns
are never accepted by this contract.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another. Its working tree and build are not used.
- What G3 needs to build `olympia-g3`: CMake 3.25 or later, Ninja and a C11 compiler.
- Go (the version in `go.mod`), Bash (3.2 is enough), Perl, `patch`, `tar`, `shasum`, `cmp` and `diff`.
- Linux or macOS (the pinned C engine owns its BSD short-option parser).

## Running

```sh
tools/olympia-roundtrip/run.sh                    # every case
tools/olympia-roundtrip/run.sh world move-go      # just these cases
tools/olympia-roundtrip/run.sh 'move-*' '*+turn'  # shell patterns on case names
tools/olympia-roundtrip/run.sh --list             # print available case names; no builds/runs
tools/olympia-roundtrip/run.sh --fresh            # drop the cached C side first
tools/olympia-roundtrip/run.sh --determinism-all  # check C's determinism on every case
```

Selecting `release-swear-np`, or a pattern matching it or its unavailable
`+turn`, also selects the C-only abort control. `--list` prints only available
saved-game names and describes the unavailable snapshot on stderr, without
building or running either engine. Selecting only `release-swear-np+turn`
fails explicitly with no available case, rather than succeeding on an empty
selection. Broad patterns such as `*+turn` run the available snapshots and
also validate/account for the expected abort.

The script:

1. Exports the revision with `git archive` into `tmp/olympia-roundtrip/g3` and builds `olympia-g3` there with the `debug` preset.
2. Builds the Go port into `tmp/olympia-roundtrip/bin/olympia-g3`.
3. For each selected case, makes the input and the C result unless they are cached.
   On a fixed sample of cases (`world`, `move-go`, `trade-routes` and each `+turn`; `DETERMINISM_SAMPLE` in `run.sh`), C runs the round trip twice, and the case is an `ERROR` if its two runs differ.
   `--determinism-all` checks every selected case, including those already cached without the check.
   A case is also an `ERROR` if the game cannot be unpacked or an ordinary C turn fails.
   When selected, the exact-pin expected-abort control always runs twice,
   even with a warm cache; a wrong status/assertion, timeout, dependency marker,
   sanitizer report or nondeterministic result is an `ERROR`.
4. Runs Go on a fresh copy of the input and compares.
   It prints one line per case: `PASS`, or `FAIL` with what differs (the exit statuses, `stdout`, `stderr`, and the number of files that differ with the first five), then a `pass N fail N error N expected-abort N` line.
   Matching Go captures are deleted promptly; failures remain for inspection.

The export, the C build, the inputs and the C results are kept and reused while the revision, the C binary and the way they are made stay the same.
Cache format remains 1: successful inputs and C load/save results are unchanged,
so exact-pin builds/results remain valid. Expected-abort runs are separate,
never treated as cached successful snapshots. A legacy failed `+turn` directory
may remain as historical diagnostics; without `input/lib` it is not a game
for the full-corpus tools.
Only Go is rebuilt and rerun, so a later run takes about 25 seconds for every case, or under a second for a few.
The first run, or one after `--fresh`, takes about 3 minutes.

Each case's files are under `tmp/olympia-roundtrip/cases/<case>`:

| Path | What it is |
|---|---|
| `input/` | the saved game both engines start from |
| `turn/` | for a `+turn` case, the C turn's `stdout`, `stderr` and `status` |
| `c/`, `go/` | each engine's run directory (`cwd/`) and its `stdout`, `stderr` and `status`; matching `go/` is removed |
| `diff.txt` | the statuses, then `diff -u` of stdout and stderr, then `diff -ru` of the run directories |
| `error` | why the case could not be prepared |
| `deterministic` | C's round trip was run twice and gave the same result |

`tmp/olympia-roundtrip/results.tsv` summarizes the last run, one row per case it ran.
The selected abort control has a separate `EXPECTED_ABORT` (or `ERROR`) row,
with C status and no Go status. Its raw captures remain under
`tmp/olympia-roundtrip/expected-aborts/release-swear-np+turn/{turn,turn-again}`;
matching partial `cwd/` trees are removed after exact comparison, and the
`deterministic` marker records the repeat. It exits 0 only when every selected
case passes and every selected control meets its contract.

## Results

For #179 (**0.108.1-alpha**, Linux amd64, GCC 12.2.0, Go 1.26.4), rebuild
from pristine `f73e48118805ffe8ff5e763bfa4f20ae00c3a705` under
`tmp/issue179-roundtrip`, leaving `tmp/stage3-roundtrip` intact. All **201**
inherited names and the bytes/digests of all **41938 input files** are unchanged from
the old-pin corpus. All **201 C/C determinism checks** were freshly made.
The finished-script warm repeat passes **201 / 0 fail / 0 error / 1 EXPECTED_ABORT**
in **145s**, freshly repeating the real C abort twice at `swear.c (180)` /
status **134** and every Go load/save. No successful abort snapshot exists.
The final-source repeat after the recording/full-tool checks has the same
counts in **147s**. Its final audit also confirms all **5228 input directories**
and every input file retain the old corpus's exact bytes/types/modes.
The exact new-pin abort source, fixture and README are byte-identical to the
accepted old-pin versions; the pin exception follows that review, not a
general allowance. Marker rejection now precedes all C result/status handling,
including cached references and expected aborts, with negative controls for
either stream and absent captures. Engines require real directory inputs.

All **1212** archived regular source/fixture files match the immutable export;
C binary SHA-256 `151c7355890eba7f8bfda1f4f16f4a490947aaa1db1ffcc3ca859cd659c01158`.
Logs: `tmp/issue179-roundtrip-{final,final2}.log`,
`tmp/issue179-{source-audit,corpus-audit,final-audit}.log`,
`tmp/issue179-inputs.sha256`, `tmp/issue179-roundtrip/results.tsv`.
The initial pristine run completed every case and determinism marker but its
shell footer failed after this script was edited while it was running. The
145s and 147s complete reruns use the finished script and verified cache;
that initial shell failure is not represented as a green run.
The two unchanged new upstream check scripts separately pass real C twice
and Go, including exact raw status/streams/all-entry trees on all eight turns:
`tmp/issue179-fixtures.log`. See the [pin gate](../README.md#follow-c-pin-gate-179).

Final #166 run, 2026-10-28, Linux amd64, GCC 12.2.0, Go 1.26.4,
unchanged pin `e733acb63bacb21167a1f2589d8dcdc023364896`:
**201 PASS / 0 FAIL / 0 ERROR / 1 separate EXPECTED_ABORT**, in **619s**.
The old cache was absent: source, build, all pristine inputs, all C results
and all 201 C/C determinism checks were rebuilt, not inherited from old
markers. Go ran on fresh directory copies. All 196 inherited case names pass;
the five new saved cases pass too. The actual unavailable
`release-swear-np+turn` C control ran twice with status **134**, exact
`swear.c (180): FALSE`, identical raw captures and partial trees, and no
dependency/sanitizer marker. No successful snapshot or Go turn is invented.
Matching Go/repeat/abort cwd captures were deleted by the harness.

Command: `G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/stage2-roundtrip"
tools/olympia-roundtrip/run.sh --determinism-all`. Logs/TSV:
`tmp/issue166-roundtrip-{all.log,results.tsv}`. Provenance:
`tmp/issue166-{provenance,export-bytes}.log`; all **1208** exported regular
source files equal the pinned archive bytes. `g3.rev` and format-1 `c.stamp`
agree with the fresh binary SHA-256
`0e3dee47344b8e66abc0ef91691cc27366adfa87a52b4b4d83ee14aa9dda1002`.
Inputs/results remain in `tmp/stage2-roundtrip`; full tools enumerate its
201 actual `input/lib` directories plus their own base world = **202 worlds**.

Historical #133 run at `e733acb63bacb21167a1f2589d8dcdc023364896`,
2026-10-28, Linux amd64, GCC 12.2.0,
Go 1.26.4: `--determinism-all` passes **201**, fail 0, error 0, with **1
separate expected abort**, in 153 seconds. Exact-pin format-1 inputs/C
results and their 201 determinism checks were reused; every Go load/save
ran freshly, and the expected C abort ran twice freshly with exact raw
streams/status and partial-tree equality. All 196 inherited names remain
and pass. The full tools' input enumeration is unchanged: 201 cached
`input/lib` directories plus their separate base world = 202 worlds.

Command: `G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/stage2-roundtrip"
tools/olympia-roundtrip/run.sh --determinism-all`. Evidence:
`tmp/issue133-roundtrip-all.log`, `tmp/issue133-roundtrip-results.tsv` and
`tmp/issue133-enumeration.log`. Negative controls against the actual harness
functions reject status 0/1/124/139, wrong assertion line/file, dependency
markers alongside the correct assertion, and sanitizer reports. An actual
bounded C turn on the aborting fixture, presented to ordinary `+turn`
preparation, is rejected with status 134 and creates no saved input.
Evidence: `tmp/issue133-controls.log`. Fresh `--list` exports build/run no
engines; an older `G3_REV=6a25f2d` override lists exactly the 196 inherited
cases, while selecting only the pinned unavailable snapshot exits 2.
Vet, fresh bounded default tests, default turn and command comparisons pass;
logs: `tmp/issue133-{vet,tests,turn-default,command-default}.log`.

Parent acceptance caught contamination from a temporary negative-control script:
its symlink input let C write aborted-turn logs into one cached input after the
implementer's successful run. The parent discarded that generated case and
rebuilt it from the pristine pinned scenario, then repeated the full corpus and
fresh checks (`tmp/parent133-*.log`). The corrected independent control uses a
real directory copy and rejects the ordinary failed turn without changing its
source input. No checked-in harness path uses that scratch symlink.

Historical stage-1 pin (`6a25f2d`): 2026-10-28, Linux amd64, GCC 12.2.0, Go 1.26.4:
`--determinism-all` passes all 196 cases (98 games before and after a C turn),
fail 0, error 0, in 608 seconds. Every case has a C/C determinism check;
Go matches status, raw stdout/stderr and every run-directory file, with no
differing files. Command: `G3_REPO="$PWD/tmp/upstream.git"
WORK="$PWD/tmp/issue104-roundtrip" tools/olympia-roundtrip/run.sh --determinism-all`.
Evidence: `tmp/issue90-evidence/roundtrip-all.log` and `roundtrip-results.tsv`;
reusable C results and inputs remain in `tmp/issue104-roundtrip/cases`.

Historical run: 2026-10-27, macOS 26.6.2 arm64, Apple clang 21.0.0, CMake 4.3.3, Go 1.26.4, with `--fresh --determinism-all`.

pass 196, fail 0, error 0.
C gave the same result twice on every case, and Go matches it on every case: exit status, stdout, stderr and every file of the run directory.
