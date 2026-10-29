# Strict whole-tree turn comparison (#56, #134)

Pinned reference: [playbymail/olympia-g3 at
f73e48118805ffe8ff5e763bfa4f20ae00c3a705](https://github.com/playbymail/olympia-g3/tree/f73e48118805ffe8ff5e763bfa4f20ae00c3a705).
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

## Drum fixture rerun (#182)

At **0.115.0-alpha**, only the complete original drum USE callback and real
item-dispatch wiring are newly implemented. The unchanged immutable pin is
**f73e48118805ffe8ff5e763bfa4f20ae00c3a705**. No upstream fixture/script,
expected bytes, comparator, artifact policy or normalization is changed.

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue179-roundtrip" \
  WORK="$PWD/tmp/issue182-turn" tools/turn-compat/run.sh \
  --full --command-coverage --use-coverage --phase-coverage --state-export
```

The strict **233-entry** ordinary turn and all **13 lifecycle / 10 command /
15 USE / 23 phase / 13 export** controls pass, including all four nonempty
export flag-free aftermath comparisons. Each unchanged real C run repeats
before Go; raw streams/status and every cwd entry/type/mode/empty directory/
regular byte match. The three expected Times status-1 failures remain exact.
Matching captures are removed. Evidence: `tmp/issue182-turn-full.log` and
`tmp/issue182-turn/{c-determinism,lifecycle,command-coverage,use-coverage,phase-coverage,state-coverage}.txt`.

All **six original phase scripts** run real C twice before Go, unpacking their
own pristine fixtures from a separate immutable export, never cached engine
inputs. A temporary native launcher preserves argv[0] through run_limited and
captures raw streams/status plus the **complete cwd before original cleanup**.
The existing strict comparator and tree-modes.py compare all entries without
exclusions. Every script prints **YES**; scripts and engines exit **0**:

| fixture | strict C/C and C/Go result |
|---|---|
| phase-day0 | **252 entries, newly passing** |
| phase-daily | **245 entries**, retained |
| phase-month-open | **242 entries**, retained |
| phase-month-gold | **241 entries**, retained |
| phase-month-decay | **246 entries**, retained |
| phase-month-close | **242 entries**, retained |

The final repeat also runs the newly passing original **savage-cap** script
against C twice then Go: **238 entries**, exact raw status/streams and every
entry/type/mode/byte. Its native capture includes the script-owned stderr
file after raw engine stderr has been replayed to that original destination.
All seven final script comparisons pass in `tmp/issue182-fixtures-final.log`.

Dependency/NYI/stub/sanitizer markers are rejected before status acceptance.
No further missing gameplay boundary appears. Raw per-script logs, complete
mode/type/byte digest manifests and strict diffs remain under
`tmp/issue182-fixture-captures/`; summary: `tmp/issue182-fixtures.log`.
Matching captures and temporary launchers/source exports are removed. An
initial launcher used an environment variable also assigned by the original
script, failing before capture/engine execution; this is not a gameplay
failure or accepted parity. It was corrected using a distinct variable and
the whole six-fixture run restarted. Initial evidence:
`tmp/issue182-fixtures-initial.log`.

The read-only C cache's **1212 originals** and executable SHA-256
`151c7355890eba7f8bfda1f4f16f4a490947aaa1db1ffcc3ca859cd659c01158`
are verified against the pin (`tmp/issue182-source-audit.log`). Engines run
only on real disposable copies. Final integrity also verifies all **41938**
cached input file hashes unchanged (`tmp/issue182-final-integrity.log`).
This accepts the drum prerequisite and its
fixture progress, **not stage exit**: #177/#178 remain open, and the full
201-world/four-corpus final tier is intentionally deferred to #178.

## Hinder-meditation fixture rerun (#181)

At **0.114.0-alpha**, complete original hinder-meditation start/omen/finish,
reset_cast_where and check_char_where run through real USE dispatch. No
upstream source, fixture, script, expected manifest/JSON, comparator, exclusion
or normalization is changed. Only this prerequisite is implemented, not #182.

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue179-roundtrip" \
  WORK="$PWD/tmp/issue181-turn" tools/turn-compat/run.sh \
  --full --command-coverage --use-coverage --phase-coverage --state-export
```

Strict ordinary turn passes **233 entries**, plus all **13 lifecycle /
10 command / 15 USE / 23 phase / 13 export** controls, including the four
nonempty-export flag-free aftermath comparisons. Every C run repeats before
Go; raw status/stdout/stderr and ALL cwd entry types/modes/empty directories/
regular bytes match. Matching turn/control captures are removed. Logs and
provenance: `tmp/issue181-turn-full.log`,
`tmp/issue181-turn/{c-determinism,lifecycle,command-coverage,use-coverage,phase-coverage,state-coverage}.txt`.

All **six** original phase scripts are rerun against unchanged C twice before
Go. A temporary native launcher preserves argv[0], invokes the existing
run_limited adapter, captures raw engine streams/status and copies its entire
cwd **before the original script's cleanup trap**. Scripts unpack their own
pristine scenario.patch/world.tgz from a separate immutable archive; no engine
runs on cache inputs or symlink-copy sources. The unchanged strict comparator
and tree-modes.py compare every captured entry without exceptions.

| fixture | original C twice (script / engine) | Go script / engine | strict result |
|---|---|---|---|
| phase-day0 | YES / 0 | NO / 2 (script 1) | C/C all **252 entries** identical; Go stops at #182 drum |
| phase-daily | YES / 0 | YES / 0 | C/C/Go all **245 entries** identical |
| phase-month-open | YES / 0 | YES / 0 | C/C/Go all **242 entries** identical |
| phase-month-gold | YES / 0 | YES / 0 | C/C/Go all **241 entries** identical |
| phase-month-decay | YES / 0 | **YES / 0, newly passing** | C/C/Go all **246 entries** identical |
| phase-month-close | YES / 0 | YES / 0 | C/C/Go all **242 entries** identical |

Dependency/NYI/stub/sanitizer markers are rejected before accepting status,
including expected failures. phase-day0 still stops at **USE item key 7**;
this is already filed/sized separately as #182, not a newly exposed callback.
No expansion into drum or other gameplay occurs here. Its decisive raw C/Go
captures remain in `tmp/issue181-fixture-captures/phase-day0/{c1,go}`; all
matching captures, native launcher/probe sources/binaries and temporary
fixture source archives are removed. Per-run script logs, complete mode/byte
digest manifests and strict diffs remain in `tmp/issue181-fixture-captures/`;
summary: `tmp/issue181-fixtures.log`.

Verified read-only C cache: **1212 original files** equal immutable
**f73e48118805ffe8ff5e763bfa4f20ae00c3a705**; executable SHA-256
`151c7355890eba7f8bfda1f4f16f4a490947aaa1db1ffcc3ca859cd659c01158`.
Audit: `tmp/issue181-source-audit.log`, with final integrity recheck in
`tmp/issue181-final-integrity.log`. Spec now preserves all **33** accepted
passes and adds only **phase-month-decay** (**34/91/0**). Five of six phase
fixtures pass; neither all-six parity nor stage exit is claimed. Final
201-world/four-corpus acceptance remains #178, intentionally not repeated.

## GUARD/STACK fixture rerun (#180)

At **0.113.0-alpha**, complete original GUARD/STACK starts and their small
gone/visibility validator execute in the real engine; no fixture, upstream
harness, expected bytes, turn comparator or tree-modes policy changes.

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue179-roundtrip" \
  WORK="$PWD/tmp/issue180-turn" tools/turn-compat/run.sh \
  --full --command-coverage --use-coverage --phase-coverage --state-export
```

This passes the strict **233-entry** ordinary turn, all **13 lifecycle /
10 command / 15 USE / 23 phase / 13 nonempty-export** controls and four
nonempty-export flag-free aftermath comparisons. Each real C run repeats
before Go; raw statuses/streams and every cwd entry/type/mode/empty directory/
regular byte match without normalization or exclusions. Summary files:
`tmp/issue180-turn/{c-determinism,lifecycle,command-coverage,use-coverage,phase-coverage,state-coverage}.txt`;
raw log: `tmp/issue180-turn-full.log`.

Both affected unchanged phase scripts run C twice then Go. A temporary native
launcher invokes each engine through run_limited, tees its raw streams/status
and copies its complete cwd **before the script's cleanup trap**. The scripts
still unpack their pristine original scenario.patch/world.tgz themselves.
The launcher and exported source are removed after comparison. The existing
strict turn comparator and tree-modes.py compare captures with no changes:

| fixture | original C twice | Go script / engine | strict result |
|---|---|---|---|
| phase-daily | YES / 0 | **YES / 0** | C/C/Go, all **245 cwd entries**, identical |
| phase-day0 | YES / 0 | **NO / 2**, script status **1** | C/C all **252 entries** identical; Go reaches later USE drum NYI |

The day0 diagnostic is exactly `USE item key 7: not yet implemented`, after
GUARD succeeds. Original use.c **869–872** maps key **7** (`use_drum`,
oly.h:411), moves c->b into c->a and calls complete savage.c `v_use_drum`
**70–146**. Its definition is **1483 bytes**, SHA-256
`9dc335d781d225c678f8f50fcfdf094baaa965c639e62d169d324d5fa580103d`.
Actual C's generated savage **97167** executes `use 98 1` on day **8**
(`lib/save/1/100`); this is neither an export mismatch nor a successful turn.
Missing-dependency/NYI markers are rejected before accepting engine statuses.

Separately sized next prerequisite, **not implemented here**: complete
v_use_drum plus original USE-item dispatch wiring, reusing existing
callSavage/createSavage, numSavages/cap, shared real randomized land exits,
location messages, creation, queues, initLoadSup and RNG owners; asymmetric
C/C/Go start/lifecycle/aftermath cases. No summon command, combat or ATTACK.
Parent must file/attach it before expansion. Source/sizing evidence:
`tmp/issue180-prerequisite.log`. #181 hinder meditation is untouched.

Logs/manifests: `tmp/issue180-fixtures.log`,
`tmp/issue180-fixture-captures/*-{script.log,manifest.txt,diff}`.
Matching captures (including the second deterministic day0 C copy) and all
temporary launchers/source are deleted. Decisive day0 C/Go raw captures and
complete cwd remain at `tmp/issue180-fixture-captures/phase-day0/{c1,go}`.
Read-only cache verified before reuse: all **1212** originals equal pin
**f73e481**, C executable SHA-256
`151c7355890eba7f8bfda1f4f16f4a490947aaa1db1ffcc3ca859cd659c01158`
(`tmp/issue180-cache-audit.log`). No engine runs on cached input or symlink
sources. No six-phase-fixture, final-corpus or stage-exit acceptance is claimed.

## State export save lifecycle (#176)

**0.112.0-alpha** adds `--full --state-export`: **13** actual-engine controls
using the unchanged pinned **phase-month-open** scenario, with four nonempty
regular factions, units and asymmetric routes. Each case unpacks pristine
fixture input; engines only run in fresh real copies. No upstream scripts,
scenarios, expected bytes, comparator exclusions or normalization change.

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue179-roundtrip" \
  WORK="$PWD/tmp/issue176-state-turn" tools/turn-compat/run.sh \
  --full --state-export
```

Cases: load only, load+save, turn+save, turn without save, duplicate path
(last wins), ENOENT, ENOTDIR, the pre-save `lib/log` path, collisions with
saved system/post-writer players, existing saved `lib/fact` directory,
saved randseed and Times output. In particular, state export is **after**
save_logdir: `lib/log` is no longer a directory then, unlike the earlier
coverage writers. Times follows export; its collision is compared as real
file/cache behavior, not a fabricated failure. All C statuses are **0**;
failed opens are nonfatal with original diagnostics. C runs twice before
Go; raw status/streams and every cwd entry/type/mode/byte match.

Four load/turn/save combinations also run flag-free C twice and Go. Only
after the strict complete comparison is the requested dump removed from
disposable captures; every remaining observation matches flag-free C.
The nonempty complete turn JSON is **3520 bytes**, SHA-256
`c7bd7db18c75fad5844a7470e2512535c32afe7e87d814664c57bd5a3ee1a729`,
retained at `tmp/issue176-state-turn/state-turn.json`. Load/load-save JSON
hash is `161677dae01c4c2b30589c7bbb9efbc7b8a1a7577edef48e43dcde707449c9c9`.
No selected-field JSON check substitutes for raw bytes. Log/summary:
`tmp/issue176-state-turn.log`, `tmp/issue176-state-turn/state-coverage.txt`.

The earlier complete run also passes strict **233-entry** ordinary turn,
all **13 lifecycle / 10 command / 15 USE / 23 phase** controls and an empty
root export family (`tmp/issue176-turn-full.log` and `tmp/issue176-turn/`).
The nonempty repeat above strengthens the new export controls; it does not
replace inherited cases. All **15** default tools, bounded vet and fresh
bounded full-enabled tests pass. Read-only C cache provenance is independently
verified: all **1212** archived originals equal the pin, binary SHA-256
`151c7355890eba7f8bfda1f4f16f4a490947aaa1db1ffcc3ca859cd659c01158`.

All six unchanged phase fixture scripts also ran themselves against bounded
real C twice and Go, with a temporary capture launcher preserving complete
raw engine status/streams/cwd **before each script's cleanup**. Each C/C
capture passes the strict turn comparator and tree-modes.py. Results:

| fixture | C twice | Go script / engine | complete C/Go cwd |
|---|---|---|---|
| phase-day0 | YES / 0 | NO / 2: command guard start NYI | differs; stopped before dumps |
| phase-daily | YES / 0 | NO / 2: command stack start NYI | differs; stopped before dumps |
| phase-month-open | YES / 0 | YES / 0 | **242 entries**, exact |
| phase-month-gold | YES / 0 | YES / 0 | **241 entries**, exact |
| phase-month-decay | YES / 0 | NO / 2: USE v_hinder_med callback NYI | differs; stopped before dumps |
| phase-month-close | YES / 0 | YES / 0 | **242 entries**, exact |

All three NO scripts exit **1**, not engine status 2. Successful JSON hashes:
open above; gold **2111 bytes**,
`03e6bdf8e326bca1e97195a2d116a45d143419da79c71a2e43c6bd94454f1d2c`;
close **3520 bytes**,
`4005271ab9a072c0cae4b3a7983f6718c08d51331c788522469fb4353b4a9a72`.
Logs/manifests: `tmp/issue176-fixtures{,-discovery}.log`,
`tmp/issue176-fixture-captures/`. Matching captures and temporary launchers/
scripts/source export are removed; decisive failed C/Go captures remain.

These three **first observed** missing-gameplay boundaries are separate
prerequisites, not export implementation work or an exhaustive inventory:
combat.c `v_guard` **2947–2964** (small callback), stack.c `v_stack`
**703–741** (plus real char/visibility checks), basic.c `v_hinder_med`,
`hinder_med_omen`, `d_hinder_med` **331–420** (cast-location/visibility and
real keyed magic RNG closure). Scope/source hashes: `tmp/issue176-prerequisites.log`.
No gameplay callback is replaced with success, fixture orders removed, or
USE handler ported. Per the no-issue-edits instruction no GitHub issue is
created/edited: parent must separately file/size these before expanding.
There is **no six-fixture/stage-exit acceptance**; #177 still owns broader
four-dump combinations/failure closure and #178 the final full rerun.

## Monthly phase CLI controls (#174)

At **0.110.0-alpha**, `--full --phase-coverage` adds **23** actual-world
C/C/Go controls to the existing strict comparator. No upstream harness or
expected fixture is edited. The ten original command controls and fifteen
USE controls remain independently selectable and unchanged in behavior.

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue179-roundtrip" \
  WORK="$PWD/tmp/issue174-turn" tools/turn-compat/run.sh \
  --full --command-coverage --use-coverage --phase-coverage
```

Phase cases cover load only (all 74 rows zero, including production), load
and save, turn with/without save, last duplicate wins, ENOENT/ENOTDIR/EISDIR,
destinations colliding with saved system and post-writer files, command+phase,
USE+phase and all three counters, each pair/all three sharing a path, and
prior command/USE write failures not preventing the later phase writer.
They check cwd-relative paths outside lib, including `-R` consumed as a path;
missing/empty phase paths are usage/status **1**. `--phase-coverage=cov.tsv`
is deliberately **not** long-option syntax: original getopt scans its
characters, reports illegal options, then the `t` in `cov.tsv` runs ilist_test
and exits **0** before errflag. The initial expected-1 harness assumption was
rejected by real C, corrected from source/streams, and the complete suite
rerun (`tmp/issue174-turn-discovery.log` preserves that discovery).

Every case runs unchanged original C twice and compares the raw status and
streams and **every cwd entry/type/mode/byte** before comparing Go, without
normalization or exclusions. Four phase cases additionally repeat flag-free
C/C/Go and, only after the strict full comparison, remove the requested dump
from disposable captures to prove the complete flag-free aftermath is unchanged.
All use real directory copies, never cached inputs or symlink sources.

The complete ordinary-turn **74-row / 2523-byte phase TSV** is retained at
`tmp/issue174-turn/phase-turn.tsv`, SHA-256
`9fb3a1f3d900bfc3d5b971a84966fbf7328d1f3d10838110ac43c5d04e16ab16`.
It is compared as raw bytes, not selected totals; notably token ran is **2**
on both grouped rows, production ran/acted is **1/1**, and the real turn has
**432** city transfers, **9** unformed grants, **6472** storm weakenings and
**4** link openings. Its tax/NP/claim/garrison/rank actions are zero, not
invented fixture effects. Direct meaningful nonzero coverage for all 43
monthly rows is in [command-compat](../command-compat#monthly-phase-hooks-and-cli-integration-174).

The final run passes the strict **233-entry** ordinary turn, all **13 lifecycle,
10 command, 15 USE and 23 phase** controls, including expected diagnostic
statuses. Summaries: `tmp/issue174-turn/{c-determinism,lifecycle,command-coverage,
use-coverage,phase-coverage}.txt`; raw run log: `tmp/issue174-turn-full.log`.
Matching captures/inputs are removed. No small turn recording or manifest
expectation is refreshed. The read-only cache's **1212** original archived
files and C binary SHA-256
`151c7355890eba7f8bfda1f4f16f4a490947aaa1db1ffcc3ca859cd659c01158`
are verified before reuse (`tmp/issue174-cache-audit.log`).

Marker rejection is exercised in both streams by the existing default tests;
command-compat also performs four actual linker probes, rejecting invalid
references before expected assertions or recording updates. All engines use
run_limited; every compat invocation explicitly sets isolated G3_REPO.
No full 201-world rerun, state export, six phase-fixture pass, golden success,
or full export combinations are claimed: those integrations remain #176/#177.

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

### #179 follow-C pin gate (2026-10-29)

Linux amd64, GCC 12.2.0, Go 1.26.4, merged immutable pin
`f73e48118805ffe8ff5e763bfa4f20ae00c3a705`:

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue179-roundtrip" \
  WORK="$PWD/tmp/issue179-turn" tools/turn-compat/run.sh --full --command-coverage --use-coverage
```

The freshly built unchanged C engine repeats exactly; Go matches status 0,
raw stdout/stderr and **all 233 cwd entries**, including report/Times/auxiliary
artifacts, empty directories, entry types/modes and bytes. All **13 lifecycle
controls**, **10 command controls** and **15 USE controls** pass under the
same strict contract, including the three expected status-1 Times failures
and flag-free aftermath comparisons. No exclusions or normalization.

A real-C-only recording run updates only the pin and C binary SHA-256 in
`testdata/c-determinism.txt`; its entire manifest body and all small snapshot
bytes remain unchanged. The provenance assertion advances with that recording.
The finished full rerun above passes after the initial recording run exposed
the old literal pin in that assertion. Matching captures are cleaned; retain
`tmp/issue179-turn/{c-determinism,lifecycle,command-coverage,use-coverage}.txt`
and `tmp/issue179-turn-final.log`. New fixture raw C/C/Go checks and the rest
of the inherited acceptance are in [the pin gate](../README.md#follow-c-pin-gate-179).
Remaining USE gameplay, phase coverage and state export are still gated.

### #166 final strict integration evidence (2026-10-28)

Linux amd64, GCC 12.2.0, Go 1.26.4, then-stage pin
`e733acb63bacb21167a1f2589d8dcdc023364896`:

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
