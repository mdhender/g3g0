# Turn comparison (stage 1, #56)

Pinned reference: [playbymail/olympia-g3 at
e733acb63bacb21167a1f2589d8dcdc023364896](https://github.com/playbymail/olympia-g3/tree/e733acb63bacb21167a1f2589d8dcdc023364896).
No turn engine code is implemented here. The reference is real C, never Go.

For #131, `--full --record` refreshes provenance from two actual turns at
this pin: only the pin and C binary hash change; all recorded state/raw
streams, fixture hash and entire-capture manifest entries remain identical.
The hardcoded provenance test follows the new pin. The stage-1 comparator
exclusions below deliberately remain until #134; this is not report-byte
acceptance. Historical evidence below retains its original pin and hashes.

```sh
tools/turn-compat/run.sh                 # fast recorded checks; no world turn or C build
tools/turn-compat/run.sh --full          # real C/C/Go world turns and lifecycle controls
tools/turn-compat/run.sh --full --record # explicitly refresh C-only recordings, then compare Go
```

## Default checks

The Go tests exercise the comparator's observable contract on a tiny recorded
slice of the real C world turn: `system`, `randseed`, the empty `ship`, raw
stdout, raw stderr and status. Negative controls check changed bytes, missing
files (including empty ones), extra files, orders, factions, trade routes,
unknown paths, output whitespace, exit status, and the precise report exclusions.
These mutations test the comparator, not expected engine behavior.

The harness also runs the real Go engine with the turn arguments below plus
`-?`, comparing its rejected-option diagnostic and unchanged input file with
two recorded C runs. It never treats that diagnostic as a successful turn.
`go test ./...` needs no C checkout and never runs a whole-world turn here.

`testdata/c-determinism.txt` records the C pin, binary and world archive SHA-256,
arguments, successful repeat check, and a state/output digest manifest. Tests
check each small snapshot against that manifest. This is recorded provenance,
not a claim to have rerun C during default tests. The recordings total about
20 KB; `--record` writes them only after real pinned C repeats match, before
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
and the *entire* captures must match, including deferred report files. A failed
or nondeterministic C reference stops the harness with exit 2 and retained
evidence, rather than blessing a Go comparison. Only then is Go compared to C.

After the ordinary turn matches, full mode checks thirteen lifecycle controls
on copies of that same input: no `-S`; empty, dot-only, partial and non-directory
cwd-relative `act`; and EISDIR at `times_press`, `times_rumor`, `times_0`,
`players`, `email`, `totimes`, `forward` and `factions`. Missing `act` is covered
by the ordinary world. Each control repeats C (entire captures, artifacts
included), then compares Go's raw streams/status and every non-artifact file.
Times failures require C status 1, while writer errors and no-admission paths
require C status 0. The masthead failure occurs after saving; press/rumor
failures precede the turn. Valid admissions and external account charging are
explicit unsupported branches, not successful lifecycle controls.
Only one control's captures are retained at a time; matching captures are
deleted immediately. A mismatch exits 1 with evidence under
`$WORK/lifecycle/CASE`; a C failure/nondeterminism exits 2. The concise passing
summary is `$WORK/lifecycle.txt`. Defaults still run no whole-world turn.

### Database ownership and report boundary

At the pin, `io.c:save_db` calls `save_seed`, `save_trade_routes`, `save_system`,
`write_all_boxes`, `write_master`, `save_orders` and `rename_act_join_files`.
These own `randseed`, optional `trade_routes`, `system`, `loc`, `item`, `skill`,
`gate`, `road`, `ship`, `unform`, `misc`, `fact/*`, `master`, `orders/*` and
act/join state. Comparing only an allowlist would miss new state or extra files.

Instead, the comparator compares **every file in the working directory**, with
only these exact stage-2 output artifacts excluded:

* `lib/log/*`, `lib/save/*`: `sout.c` and `io.c:save_logdir`.
* `lib/players`, `lib/email`, `lib/forward`, `lib/factions`, `lib/totimes`:
  the writers in `main.c`.
* `lib/times_0`, `lib/times_press`, `lib/times_rumor`: `c2.c`.

Report routines still run in C. Any changes they make to database state
(for example known/lore state in faction records) remain compared. Lore input,
`startloc`, unknown files and files outside `lib` remain compared too. No broad
`times_*` or faction filename filtering is used. Missing and extra files fail;
empty files count. Nonregular state files are harness errors, not followed
symlinks. Empty directories and file permissions are not serialized DB bytes.

## Cache, prerequisites and evidence

Requires Go, Bash, Perl, `cmp`/`diff`, `tar`, `shasum`; full mode also needs a
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

Example using the accepted existing export without rebuilding:

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue97-roundtrip" \
  tools/turn-compat/run.sh --full
```

On mismatch, exit 1; `$WORK/{c,go}/{cwd,stdout,stderr,status}` retain raw
evidence and `diff.txt` lists all state path/byte differences and unfiltered
stdout/stderr diffs. `c-determinism.txt` retains repeat metadata. The redundant
matching C repeat is deleted immediately; matching Go/C game outputs are
deleted after success. The input and C build remain reusable.

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
