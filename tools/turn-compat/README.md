# Strict whole-tree turn comparison (#56, #134)

Pinned reference: [playbymail/olympia-g3 at
e733acb63bacb21167a1f2589d8dcdc023364896](https://github.com/playbymail/olympia-g3/tree/e733acb63bacb21167a1f2589d8dcdc023364896).
No turn engine code is implemented here. The reference is real C, never Go.

There are no artifact exclusions. Both C/C determinism and C/Go comparison
cover every cwd entry, including empty directories, types, permission/special
bits and regular-file bytes, plus unfiltered stdout/stderr/status. The strict
world gate and all thirteen artifact-inclusive lifecycle controls pass at #166;
historical stage-1 acceptance below had a narrower contract.

```sh
tools/turn-compat/run.sh                 # fast recorded checks; no world turn or C build
tools/turn-compat/run.sh --full          # real C/C/Go world turns and lifecycle controls
tools/turn-compat/run.sh --full --record # explicitly refresh C-only recordings, then compare Go
```

## Default checks

The Go tests exercise the comparator's observable contract on a tiny recorded
slice of the real C world turn: `system`, `randseed`, the empty `ship`, raw
stdout, raw stderr, status and the small real `times_0` report. Negative
controls check changed report/database bytes, missing files (including empty
ones), extra files, orders, factions, trade routes, unknown paths outside lib,
whitespace, status, deleted empty directories, regular/directory/root modes,
special permission bits and file/directory replacement in both directions.
Symlinks, FIFOs and sockets are rejected explicitly, including in formerly
excluded report paths; required capture files and cwd cannot be missing.
Synthetic mutations test the comparator, not expected engine behavior.
The shared mode helper's own controls inject lstat/scandir failures on either
or both sides; shell controls ensure capture/listing failure stays exit 2,
distinct from an ordinary mismatch (1). C streams are checked for missing
dependency/stub/NYI markers before accepting any result or recording, including
expected failures; controls cover markers in either stream and missing streams.

The harness also runs the real Go engine with the turn arguments below plus
`-?`, comparing its rejected-option diagnostic and unchanged input file with
two recorded C runs. It never treats that diagnostic as a successful turn.
`go test ./...` needs no C checkout and never runs a whole-world turn here.

`testdata/c-determinism.txt` records the C pin, binary and world archive SHA-256,
arguments, successful repeat check, and an all-entry mode/type/digest manifest. Tests
check each small snapshot against that manifest. This is recorded provenance,
not a claim to have rerun C during default tests. The recordings total about
25 KB; `--record` writes them only after real pinned C repeats match, before
looking at Go. It never accepts NYI as a successful turn expectation.

## Full comparison

Both engines get identical copies of the pinned `world.tgz` and run:

```sh
olympia-g3 -r -l ./lib -S test-use-const-report-date
```

`argv[0]` is `olympia-g3`, stdin is `/dev/null`. Every engine invocation uses
`run_limited` from `tools/lib/timeout.sh` (`RUN_TIMEOUT`, default 120 seconds).
Test executables also use its `run-limited.sh` adapter through `bash`.
Timeout is status 124; signal death remains 128 + signal, without shell noise.
No stdout/stderr filtering, NYI special case, or status normalization occurs.

The C turn runs twice on every `--full` invocation. Both statuses must be zero
and the *entire* captures must match, including all report files. A failed
or nondeterministic C reference stops the harness with exit 2 and retained
evidence, rather than blessing a Go comparison. Only then is Go compared to C.

After the ordinary turn matches, full mode checks thirteen lifecycle controls
on copies of that same input: no `-S`; empty, dot-only, partial and non-directory
cwd-relative `act`; and EISDIR at `times_press`, `times_rumor`, `times_0`,
`players`, `email`, `totimes`, `forward` and `factions`. Missing `act` is covered
by the ordinary world. Each control repeats C (entire captures, artifacts
included), then compares Go's raw streams/status and every cwd entry.
Times failures require C status 1, while writer errors and no-admission paths
require C status 0. The masthead failure occurs after saving; press/rumor
failures precede the turn. Valid admissions and external account charging are
explicit unsupported branches, not successful lifecycle controls.
Only one control's captures are retained at a time; matching captures are
deleted immediately. A mismatch exits 1 with evidence under
`$WORK/lifecycle/CASE`; a C failure/nondeterminism exits 2. The concise passing
summary is `$WORK/lifecycle.txt`. Defaults still run no whole-world turn.

### Whole-tree contract

At the pin, `io.c:save_db` calls `save_seed`, `save_trade_routes`, `save_system`,
`write_all_boxes`, `write_master`, `save_orders` and `rename_act_join_files`.
These own `randseed`, optional `trade_routes`, `system`, `loc`, `item`, `skill`,
`gate`, `road`, `ship`, `unform`, `misc`, `fact/*`, `master`, `orders/*` and
act/join state. Comparing only an allowlist would miss new state or extra files.

Instead, the comparator captures **every entry in the working directory**,
including the cwd root itself, `lib/log`, `lib/save`, Times and auxiliary
reports, lore input, `startloc`, unknown paths and everything outside `lib`.
Missing/extra entries and type replacements fail; empty files and directories
count. File contents are raw bytes, with no whitespace or routing normalization.

`tools/lib/tree-modes.py` actually compares lstat types and permission/special
bits on **all entries**, including directories and the root, not only regular
files. The Go capture deliberately enforces the same stronger mode contract;
the shared helper remains unchanged. The shell runs both comparators for C/C
and C/Go. Symlinks and other nonregular/nondirectory entries are unsupported
capture errors (exit 2), not silently omitted, followed or treated as ordinary
differences. The helper alone can list such types; capture validation rejects
them first. Failed/partial listings are never accepted as equal.

Capture-container permissions and stdout/stderr/status *file modes* are harness
plumbing, not engine output: only their exact bytes are compared. In the
manifest these observations have the sentinel mode `----------`; cwd entries
carry their real modes. Directory rows have `-` instead of a byte digest.

## Command coverage CLI controls (#170)

`--full --command-coverage` uses this same full engine runner and strict
comparator, not an alternate engine harness. Ten cases run unchanged pinned C
twice before Go: load and load/save, turn with and without save, repeated
`--coverage` options (last path wins), ENOENT/ENOTDIR/EISDIR nonfatal opens,
and collisions with `lib/system` and the post-turn `lib/players` report.
Dump paths outside lib, all raw streams/status and every entry/type/mode/byte
are compared. `lib/system` is written by coverage before save overwrites it;
`lib/players` retains its coverage bytes through the subsequent lifecycle.

The first four cases independently repeat flag-free C and compare flag-free
Go. Only after strict C/C/Go comparison of the coverage captures is the one
requested new `cov.tsv` removed from those disposable captures, then all
remaining state/RNG/report/filesystem/stream observations must equal the
flag-free reference. This secondary side-effect control is not an exclusion
from the strict primary comparison. Matching captures are removed promptly;
`command-coverage.txt` retains C dump hashes and the passing summary.
The ordinary world turn and all 13 inherited lifecycle controls still run.
Defaults remain fast and do not run any world engine.

## USE coverage CLI controls (#172)

`--full --use-coverage` extends the same runner and all-entry comparator used
by `--command-coverage`. The original ten controls also run for USE (load,
load/save, turn, no-save, last-option-wins, three open failures and two path
collisions). Four compare independent flag-free C/C/Go aftermath, **only after**
strict primary C/C/Go equality, removing only the requested new dump from
disposable captures for that secondary check. Three further cases combine
command/USE dumps, share their destination (USE writes last), or fail command
coverage's open while USE still writes. No normalization/artifact exclusions.

Two meaningful initialized-world controls unpack original pinned patches
freshly via original `scenario_unpack`, never cached/post-engine game inputs:

- `shared-checks`: `use-checks/scenario.patch`, SHA-256
  `3c33e158e35ad2a78743980ed8413a9cd915da73a977056bd00e16bad5925002`.
  Remove only `orders/300` (the ten unported implicit callbacks); keep original
  check orders for 301/302, all factions, world state and real shared checks.
- `hide`: `use-stealth/scenario.patch`, SHA-256
  `0c57dcd7fe9f3c07e0f03c9ae119cda3b5f63747825181a361171bdd2a28343a`.
  Keep the whole world/factions; select only implemented shadow noble 6947's
  `use 638 1` then `use 638 0` orders. This exercises actual start/completion
  and inhibited zero-time completion without fake thief/sneak/torture/form.

These are explicitly downstream **selected-order** controls, not claims that
the unchanged upstream full fixtures or remaining USE gameplay pass. Both
also request command coverage. Each entire unchanged C engine runs twice and
matches raw status/streams and every cwd entry/type/mode/byte before Go is
compared. The ordinary world and all thirteen lifecycle controls still run.
Matching captures are deleted; `use-coverage.txt` retains summaries/dump hashes.
Defaults remain fast, with no world engines or new duplicate harness.

At 0.108.0-alpha, the ordinary **233-entry** turn, all **10 command** and
**15 USE** controls and all **13 lifecycle** controls match. The four ordinary
aftermath cases per dump and both selected-order aftermath cases also match.
Each selected USE dump is **5,100 bytes**. Original `use_cells` rejects vacuous
selected controls: `shared-checks` has invalid checks **5**, skill checks **1**
each for **603/613/631/801**, and item checks **1** each for **8/9**; `hide`
has skill **638 ok=2**. All other selected cells are zero.

Actual C dump SHA-256 values (also exactly matched by C repeat and Go):

| Control | USE dump | Command dump |
| --- | --- | --- |
| Ordinary zero-count world | `b01e1c9960eeeb0d076bf033000139f00e53b98957f476bdbb4e2417213549e0` | See command controls above |
| Selected shared checks | `2083eb2692db35b59f52c679cd8f9e7668a571b47bc39806bf48cb4fe7f083a5` | `263cdc40552e1f654d3487ad51d9a7bf420f3ab6d25100112ef7220f7e9ce252` |
| Selected hide | `c30833c3ee3cbaccf2be0ced5f02035a040799b8e5ee2b4340a7e6e3732adf87` | `a7cc4922749a7aeb516feee75b547b3e33f7acb752439e30b3641951727802bb` |

Provenance: exact-pin C binary SHA-256
`1f03220bac19c2b3fd3a75c3b3aac34cc0df56210472148b11fba82d8e05c8bf`,
pristine world archive SHA-256
`5c59fc014ea3d2f9a14e4756bfa7ccc8c7c11b23bf6611320713ec19760d8ce9`.
The binary is reused read-only from `tmp/stage3-roundtrip`; every engine input
is a disposable copy. Logs: `tmp/issue172-turn-full.log` and
`tmp/issue172-turn/{c-determinism,command-coverage,use-coverage,lifecycle}.txt`.

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/stage3-roundtrip" \
  WORK="$PWD/tmp/issue172-turn" tools/turn-compat/run.sh --full --command-coverage --use-coverage
```

## Cache, prerequisites and evidence

Requires Go, Bash, Python 3, Perl, `mkfifo`, `diff`, `tar`, `shasum`; full mode also needs a
clone containing the pin and CMake/Ninja/C11 when no reusable C build exists.

* `WORK`: defaults to `tmp/turn-compat`; all builds and captures live here.
* `G3_REPO`: defaults to `~/Software/playbymail/olympia-g3`; only `git archive`
  reads the pinned commit. A bare cache such as `tmp/upstream.git` works.
* `C_CACHE`: defaults to `tmp/olympia-roundtrip`. An exported cache with an
  exact `g3.rev` and `g3/build/debug/olympia-g3` is reused read-only. Otherwise
  the script exports the exact commit into `$WORK/reference/g3` and builds
  with upstream's unchanged debug preset. It never builds in the source clone.
* Input is cached by pin and archive digest. C is repeated each full run,
  so an old deterministic marker cannot hide changed C behavior.
  Engines run only in actual directory copies, never against cached inputs
  or using a symlink as the input-copy source.

Example using the accepted existing export without rebuilding:

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/stage2-roundtrip" \
  tools/turn-compat/run.sh --full
```

On mismatch, exit 1; `$WORK/{c,go}/{cwd,stdout,stderr,status}` retain raw
evidence and `diff.txt` lists all cwd path/type/mode/byte differences and unfiltered
stdout/stderr diffs. `c-determinism.txt` retains repeat metadata. The redundant
matching C repeat is deleted immediately; matching Go/C game outputs are
deleted after success. The input and C build remain reusable.

### #166 final strict integration evidence (2026-10-28)

Linux amd64, GCC 12.2.0, Go 1.26.4, unchanged exact stage pin above:

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/stage2-roundtrip" \
  WORK="$PWD/tmp/issue166-turn" tools/turn-compat/run.sh --full
```

Two fresh pinned C turns match, then Go matches status **0**, raw streams and
all **233 cwd entries**, including saved reports, Times, auxiliary files,
empty directories and all-entry types/permission/special bits/regular bytes.
All **13** controls repeat C and match Go under the identical strict contract:
no-save **230** entries; empty-act/act-file **234**; dot-act/partial-act **236**;
press failure **219**, rumor failure **220**, and the other six controls **233**.
Press/rumor/masthead failures have status **1**; the other controls status **0**.
No report exclusion, normalization, source pin or recording changed.

Evidence: `tmp/issue166-turn-full.log` and
`tmp/issue166-turn/{c-determinism,lifecycle}.txt`; the manifest records fresh
binary/world digests. Matching whole-world/repeat/control captures and input
are removed. The parent's independent repeat also passes the whole-world
comparison and all thirteen controls: `tmp/parent166-turn-full.log` and
`tmp/parent-stage2-turn/{c-determinism,lifecycle}.txt`.
The exit world has no regular players, posts or email rows: the unchanged
**1234-section** command oracle supplies asymmetric real-C/C/Go coverage for
those branches, rather than inferring them from this world.

### #134 strict-gate evidence (2026-10-28)

With `G3_REPO="$PWD/tmp/upstream.git"`, `C_CACHE="$PWD/tmp/stage2-roundtrip"`
and `WORK="$PWD/tmp/issue134-turn"`, `--full --record` ran two fresh unchanged
pinned C turns, both status 0. All 233 cwd entries, modes/types/bytes and raw
streams/status matched. The provenance manifest now includes directories,
modes and every previously excluded report digest. Original small recording
bytes are unchanged; only the tiny real-C `times_0` report is added.

C/Go both exit 0 with identical raw stdout/stderr. Strict comparison exits 1
on exactly eight report-byte differences: `lib/save/1/{100,200,202,203,204,205,
210,players}`. No paths/types/modes or other bytes differ. These are pending
stage-2 reports, not normalized successes. Lifecycle controls are not reached
after the ordinary strict turn fails; their comparator is strengthened too,
but their historical passing status is not new strict acceptance.
Evidence: `tmp/issue134-turn-full-record.log`,
`tmp/issue134-turn/{determinism.diff,c-determinism.txt,diff.txt,c,go}`.

### Historical stage-1 evidence (narrower artifact-excluding contract)

Initial Linux amd64 recording: C statuses 0/0, identical stdout, stderr and
entire post-turn directory; Go status 2, `olympia-g3: load_db (-r -S): not yet
implemented`. The full check fails on status, stdout, stderr, 15 changed state
files and missing `orders/204`. That is the required red integration gate,
not a newly discovered bug and not a passing turn claim.

After #71, the opt-in diagnostic reaches the real NPC dispatcher and stops at
`olympia-g3: ping_garrisons: not yet implemented`. Two pinned C runs still exit
0 and match including reports; Go exits 2. The raw/state comparison remains
red (15 changed database files, missing `orders/204`, stdout/stderr/status).
Garrison/token and later daily/month phases are not skipped or accepted.

After #73 (0.50.0-alpha), real garrison setup and token reconciliation run
at their existing day-zero positions. The opt-in diagnostic using
`C_CACHE="$PWD/tmp/issue104-roundtrip"` repeats pinned C with statuses 0/0
and identical entire captures, including reports. Go exits 2 at
`olympia-g3: match_all_trades: not yet implemented`, after day-zero interrupted
and player orders and lore scanning, at day 1 before daily command/events.
Comparison still honestly fails on the same 15 database paths, missing
`orders/204`, raw stdout/stderr and status. This is diagnostic progress,
not full-turn or stage-exit acceptance; default checks remain recorded/fast.

After #81 (0.59.0-alpha), the opt-in diagnostic completes all 30 days and
enters `post_month()`, stopping at `olympia-g3: move_city_gold: not yet
implemented`. Two pinned C turns still exit 0 and match in their entire
captures, including reports; Go exits 2. The comparison remains red on
15 changed database paths, missing `orders/204`, raw stdout/stderr and status.
This establishes the next real integration boundary, not stage-exit acceptance.

After #123 (0.76.0-alpha), the actual main tail is wired in pinned C order.
With `G3_REPO="$PWD/tmp/upstream.git"` and
`C_CACHE="$PWD/tmp/issue104-roundtrip"`, full comparison exits 0: C statuses
0/0 with identical entire captures, then Go status 0 with identical raw
stdout/stderr and all 190 compared state paths/bytes (no missing/extra files).
All thirteen lifecycle controls also match, with deterministic C repeats:
no-save compares 189 files; dot/partial/non-directory act compare 191;
press/rumor failures compare 189; the remaining controls compare 190.
No output filtering, expectation re-recording or comparator exclusions changed.
Reports/artifact bytes remain deferred; #90 retains final stage-exit acceptance.

Final #90 evidence rerun, 2026-10-28 (Linux amd64, GCC 12.2.0, Go 1.26.4):
`G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue104-roundtrip"
tools/turn-compat/run.sh --full` exits 0, independently repeating the complete
C/C and C/Go checks above: all 190 state files, raw streams/status and all
13 lifecycle controls match. Logs: `tmp/issue90-evidence/turn-full.log`;
metadata and control summary: `tmp/turn-compat/{c-determinism,lifecycle}.txt`.
The stage-1 turn/state integration is verified, not pending. The separate
full boxfile/check/loadsave gates also pass after independently accepted #130
(0.76.1-alpha): 197 games each, boxfile in all three variants. Evidence is in
`tmp/parent130-{boxfile,check,loadsave}-full.log`; the initial Linux mode-listing
failure is resolved, not a turn mismatch or an outstanding blocker. Final #90
acceptance passed the parent's independent rerun, including the full turn and
all lifecycle controls (`tmp/parent90-turn-full.log`).
Report artifact equality remains outside this stage, and no recordings changed.
