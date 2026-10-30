# Strict whole-tree turn comparison (#56, #134)

Pinned reference: [playbymail/olympia-g3 at
9a03d126b329ed466c4af48425553c025eaca26f](https://github.com/playbymail/olympia-g3/tree/9a03d126b329ed466c4af48425553c025eaca26f).
No turn engine code is implemented here. The reference is real C, never Go.

The #291 merged-pin proof, current capture commands/provenance and explicit
remaining original-fixture NYIs are in
[merged-pin proof](../README.md#merged-pin-proof-291). The older snapshot and
stage-exit evidence below retain their historical scope.

## Repair and improvement original-fixture proof (#223)

At **0.145.0-alpha**, immutable **9a03d126**, the complete direct REPAIR/
IMPROVE owners make all three unchanged original scripts below pass. The
anticipated STOP/USE boundary was **not reached**: USE601's real sailing
interrupt was already accepted in219, so stop-running now finishes too.
This is not acceptance of other USE callbacks240, RAZE224, combat235, a
construction family exit or a stage exit. No upstream bug/blocker was
demonstrated. Supplied #289 ATTACK235/fog264/missing-argument ownership and
the stage5 immediate/console deferrals are unchanged.

Before reusing `tmp/_issue215/canonical`, verify every **1237 original
blob/mode** against isolated bare pin9a03d126 and the canonical executable:
SHA-256 `2ea19140ebd00e1ca672bd022ff958c139efb3aeca59f5f91b3cb8ebb93b5212`.
Debug/Ninja and original
`-ffile-prefix-map=/home/user/workspace/repo/tmp/_issue215/canonical/=`
retain the trailing slash. `tmp/issue223-canonical-provenance.log` records
this verification; no source, fixture, pin or canonical build was changed.
Protected `tmp/issue291-roundtrip` stayed idle and was not used for engine
runs, input copies or a repeat serializer corpus. Localized callbacks add
no new serializer dependency and imply no new full207-corpus proof.

Read complete original READMEs/check.sh/Lua, construction scenario.patch
orders, and stop-running's actual frozen scenario.tgz command records, not
just its backlog names. Source reading is retained in
`tmp/issue223-original-fixture-source.log` / `tmp/issue223-stop-frozen.log`;
the parent supplied `tmp/parent223-fixture-source.txt` as additional reading.
No Lua authoring, regeneration, original --update or input overlay was used.

Recreated native pre-cleanup launcher: preserve argv0/stdin, capture streams
**O_RDWR**, invoke bash sourcing timeout.sh with
`run_limited 120 "$0" "$S4_ENGINE" "$@"`, seek/replay stdout then stderr,
close both, write real status/newline, then `cp -a .` to the absent capture
cwd before original script cleanup. Each invocation has fresh sequential
001/002/etc, never overwritten. Launcher source/binary SHA-256:
`6fa176d41814de8a9e24423b198cf1b9da9b3efe3ef6a37b8d768d20caf72c48` /
`fa788748b34ba75651013eec04e69e39bd62feea2d55ebd0053651602345c12a`.
Source/commands/hashes remain in `tmp/issue223-capture-provenance.log`.

```sh
. tools/lib/timeout.sh
export G3_REPO="$PWD/tmp/upstream.git" S4_ROOT="$PWD"
# Recreate the launcher described above and build Go freshly under tmp/_issue223.
# For EACH fixture below run C1, C2, THEN Go, using fresh absent capture bases.
export S4_CAPTURE="$PWD/tmp/_issue223/$fixture-$run"
export S4_ENGINE="$PWD/tmp/_issue215/canonical/build/canonical/olympia-g3"
OLYMPIA_ENGINE_BIN="$PWD/tmp/_issue223/capture" \
  run_limited 150 bash /bin/bash \
  "$PWD/tmp/_issue215/canonical/tests/olympia/regress/$fixture/check.sh"
# Reject C markers before accepting status, compare complete invocation inventories,
# then for every001/002/etc use the EXISTING strict comparator and mode check:
run_limited 120 compare "$PWD/tmp/_issue223/compare" --manifest "$capture"
run_limited 120 compare "$PWD/tmp/_issue223/compare" "$c1" "$c2"
run_limited 120 compare "$PWD/tmp/_issue223/compare" "$c1" "$go"
python3 -I tools/lib/tree-modes.py "$c1/cwd" "$c2/cwd"
python3 -I tools/lib/tree-modes.py "$c1/cwd" "$go/cwd"
```

**Nine native engine captures**, three actual invocation branches. For each
fixture C1/C2 precede Go; all inventories are exactly001, every engine status
is **0**, every original script returns **YES/0**. C markers are rejected
before accepting statuses. All raw streams/status and every cwd path/type/
mode/empty directory/byte match, with no exclusions or normalization:

| Unchanged original script / invocation | C/C/Go cwd entries | Result |
| --- | --- | --- |
| construction-build001 | 242 | strict exact; original coverage/hashes/JSON pass |
| construction-refused001 | 245 | strict exact; original coverage/hashes/JSON pass |
| stop-running001 | 259 | strict exact; original coverage/hashes/JSON pass |

STOP exercises saved repairs with damage+defense, defense-only, damage-only
and no progress, alongside the existing collect/make/sail/USE601 owners and
their real successor orders. The direct command oracle separately checks
IMPROVE's NULL interrupt and saved STOP without inventing a callback.
Evidence: `tmp/issue223-fixtures.log`, all
`tmp/issue223-*-{inventory,manifest}.txt`, `tmp/issue223-*-{cc,cgo}.log` and
original-script logs. Matching captures and scratch binaries/helpers/WORKs
are removed after proof; decisive logs/manifests and the old-Go dispatch NYI
log remain. See [command proof](../command-compat/README.md#repair-and-improve-223-stage-4-n22a)
and [full spec no-loss](../olympia-spec/README.md#repair-and-improvement-no-loss-223).

## Daily BUILD original-fixture proof (#222)

**0.144.0-alpha**, unchanged immutable **9a03d126**. Before reusing
`tmp/_issue215/canonical`, all **1237 original blobs and modes** were verified
against the isolated bare pin, together with Debug/Ninja and
`-ffile-prefix-map=/home/user/workspace/repo/tmp/_issue215/canonical/=`
(trailing slash retained). Canonical binary SHA-256 remains
`2ea19140ebd00e1ca672bd022ff958c139efb3aeca59f5f91b3cb8ebb93b5212`.
Audit: `tmp/issue222-canonical-provenance.log`. No original C/source/fixture,
pin or protected full207 input cache was changed; no engine ran in that
cache or on linked inputs. No new full-corpus claim.

Read unchanged construction-build/refused READMEs, check.sh, Lua and frozen
scenario orders, raft-sail's three cases and **both actual raft-loc-depth
branches**. No authoring/regeneration, upstream --update or fixture edits.
The native pre-cleanup launcher was recreated under unique `tmp/_issue222`:
preserve argv0/stdin, open raw streams **O_RDWR**, invoke bash sourcing
timeout.sh and `run_limited 120 "$0" "$S4_ENGINE" "$@"`, seek/replay stdout
then stderr, close, write actual status/newline, then `cp -a .` to an absent
capture/cwd **after replay and before script cleanup**. Every invocation
gets distinct sequential001/002/003, never overwritten. Launcher source /
binary SHA-256:
`aa20eeb688fde9d0e6de417df1e242acd17b1fc28d8191dc6cfb043956fa8ed0` /
`49af958e2316c2eff170060406aab9f2fd67418c8316b9eb38556a318dee2d16`.
Provenance: `tmp/issue222-capture-provenance.log`.

```sh
. tools/lib/timeout.sh
export G3_REPO="$PWD/tmp/upstream.git" S4_ROOT="$PWD"
# Recreate the native sequential launcher as described above.
# For EACH fixture, run C1, C2, then Go (fresh tmp/_issue222/olympia).
export S4_CAPTURE="$PWD/tmp/_issue222/$fixture-$run" # fresh absent run base
export S4_ENGINE="$PWD/tmp/_issue215/canonical/build/canonical/olympia-g3"
OLYMPIA_ENGINE_BIN="$PWD/tmp/_issue222/capture" \
  run_limited 150 bash /bin/bash \
  "$PWD/tmp/_issue215/canonical/tests/olympia/regress/$fixture/check.sh"
# Compare complete invocation inventories before comparing each001/002/003:
run_limited 120 compare "$PWD/tmp/_issue222/compare" --manifest "$capture"
run_limited 120 compare "$PWD/tmp/_issue222/compare" "$c1" "$c2"
run_limited 120 compare "$PWD/tmp/_issue222/compare" "$c1" "$go"
python3 -I tools/lib/tree-modes.py "$c1/cwd" "$c2/cwd"
python3 -I tools/lib/tree-modes.py "$c1/cwd" "$go/cwd"
```

**21 native engine captures**, seven actual invocation branches. For each
script fresh **C1,C2 precede Go**, C marker rejection precedes accepting0,
and inventories match across all three runs. All seven C/C comparisons
match raw status/stdout/stderr plus **every cwd path/type/mode/empty
directory/byte**, using existing turn-compat and tree-modes.py without
exclusions or normalization. Native scripts have C engine0/scriptYES0.

| Original script / invocation | Entries | Go engine / script | Strict result |
| --- | --- | --- | --- |
| construction-build001 | 242 | 2 / NO1 | C/C exact; now actual **command improve start NYI**, #223 |
| construction-refused001 | 245 | 2 / NO1 | C/C exact; actual **command repair start NYI**, #223 |
| raft-sail001 sail | 235 | 0 / YES0 | C/C/Go exact |
| raft-sail002 lose | 235 | 0 / YES0 | C/C/Go exact |
| raft-sail003 lose-ship | 235 | 0 / YES0 | C/C/Go exact |
| raft-loc-depth001 load | 227 | 0 / NO1 (whole script) | This branch C/C/Go exact |
| raft-loc-depth002 attack | 235 | 2 / NO1 | C/C exact; **regular_combat NYI**, #235 |

Direct BUILD no longer stops at start, but unfinished construction scripts
are **not** accepted as passes or a family exit. Existing raft-sail all
three cases remain passing. Stage4 #289's settled ATTACK235/fog264/original
missing-argument obligations and stage5 console remainder are unchanged.

Logs/manifests/inventories: `tmp/issue222-fixtures.log`,
`tmp/issue222-*-{cc,cgo,*-script}.log`,
`tmp/issue222-*-{manifest,inventory}.txt`. Failed Go captures retained in
`tmp/issue222-construction-{build,refused}-go/001` and
`tmp/issue222-raft-loc-depth-go/002`. Matching captures, temporary native/
Go binaries/helpers and disposable WORK/RUN_DIRs are removed. Completion's
bounded real-C/C/Go [command oracle](../command-compat/README.md#daily-build-and-completion-222-stage-4-n21)
discriminates the complete callbacks independently of those later NYIs.

## Construction start prerequisite (#221)

**0.143.0-alpha**, unchanged immutable **9a03d126**. All **1237 original
source blobs/modes** in `tmp/_issue215/canonical` were verified against the
isolated bare pin before native reuse. The Debug/Ninja build retains
`-ffile-prefix-map=/home/user/workspace/repo/tmp/_issue215/canonical/=`;
binary SHA-256 remains
`2ea19140ebd00e1ca672bd022ff958c139efb3aeca59f5f91b3cb8ebb93b5212`.
Provenance: `tmp/issue221-canonical-provenance.log`. The protected full207
input cache is idle/unmodified; no in-place engine/input links, shared C
checkout, new pin or full-corpus rerun.

Original construction-build/refused scripts, READMEs, Lua and frozen
scenarios were read unchanged. Neither authoring/regeneration nor upstream
fixture edits are used. Both need complete BUILD222 and repair/improve223;
the helper oracle is deliberately not partial direct dispatch.

The native pre-cleanup capture method below was recreated under unique
`tmp/_issue221`: O_RDWR raw streams, preserved argv0/stdin, bash sourcing
timeout.sh and `run_limited 120 "$0" "$S4_ENGINE" "$@"`, seek/replay stdout
then stderr, close, actual status/newline, and `cp -a .` to an absent cwd
after replay and before script cleanup. Every invocation gets its own
001/002/... directory; both scripts actually invoke the engine once.
Launcher source/binary hashes are
`d67b22c6f2f51114edc47907566689870cb3283b924d9d9fdcf15903943012cc` /
`9b40f6ebb9dd7abec781b5db29dd69a2f2d38aa405c7ffe2fbd47891601a3932`.

```sh
. tools/lib/timeout.sh
export G3_REPO="$PWD/tmp/upstream.git" S4_ROOT="$PWD"
# Recreate native sequential capture as described above, then C1/C2/Go:
export S4_CAPTURE="$PWD/tmp/_issue221/$fixture-$run"
export S4_ENGINE="$PWD/tmp/_issue215/canonical/build/canonical/olympia-g3"
# Go run only after both C repeats; use fresh tmp/_issue221/olympia instead.
OLYMPIA_ENGINE_BIN="$PWD/tmp/_issue221/capture" \
  run_limited 150 bash /bin/bash \
  "$PWD/tmp/_issue215/canonical/tests/olympia/regress/$fixture/check.sh"
run_limited 120 compare "$PWD/tmp/_issue221/compare" "$c1/001" "$c2/001"
python3 -I tools/lib/tree-modes.py "$c1/001/cwd" "$c2/001/cwd"
```

**Six native captures**, actual invocation inventories exactly001. Real C1
and C2 run before Go for each script. Marker rejection precedes C status
acceptance. Strict comparator plus tree-modes.py match raw status/streams
and every cwd path/type/mode/empty directory/byte: construction-build
**242 entries**, construction-refused **245 entries**, C engine0/scriptYES0.
Go honestly differs: **engine2/scriptNO1**, respectively **command build
start NYI** and **command repair start NYI**. Completion stays NYI too;
no construction script or family exit is newly accepted.

Logs/manifests: `tmp/issue221-construction-*-{cc,cgo,*-script}.log`,
`tmp/issue221-construction-*-manifest.txt`, `tmp/issue221-fixtures.log`,
source/capture provenance. Failed Go captures remain in
`tmp/issue221-construction-{build,refused}-go`; matching C captures,
temporary launcher/compiler/Go binaries and disposable WORKs are removed.
See [bounded helper proof](../command-compat/README.md#construction-predicates-table-and-start-221-stage-4-n20).

## Ferry originals (#220)

**0.142.0-alpha**, unchanged immutable **9a03d126**. All **1237 original
source blobs and 644/755 modes** in `tmp/_issue215/canonical` were verified
against the isolated bare pin before reuse; canonical C SHA-256 is still
`2ea19140ebd00e1ca672bd022ff958c139efb3aeca59f5f91b3cb8ebb93b5212`.
The pristine source and protected `tmp/issue291-roundtrip` cache remain
untouched. No cached input is executed in place or linked, and no shared
additional C checkout is used. Original scripts/READMEs/frozen scenarios and
EXPECT files are unchanged; no `--update`, Lua authoring or build-scenario
prerequisite is introduced.

The [bounded native pre-cleanup method](#recruitment-raise-and-refused-215)
was recreated under unique `tmp/_issue220`. The native launcher retains
argv0/stdin, opens raw captures **O_RDWR**, invokes bash sourcing timeout.sh
and `run_limited 120 "$0" "$S4_ENGINE" "$@"`, then seeks/replays stdout
before stderr, closes both, writes actual status/newline and finally
`cp -a . "$capture/cwd"` **after replay, before original script cleanup**.
Script-owned redirected out/err are therefore included. Capture destinations
are absent, not existing directories that would nest a second copy. Launcher
executable SHA-256:
`750c45d79a4bdc0127b84cecdae500534d268d6e18275735eea5cfc1730208b8`;
source SHA-256:
`f4bb57495e5271a934369866a0643f063ab136c03e5ebe9984b19116f232554f`.

Each fixture runs real **C1, C2, then Go**, with every engine invocation
allocated a distinct sequential **001/002/003** capture directory. Inspection
of these six scripts finds exactly **one invocation each**; all complete
invocation inventories are exactly `001`, compared across C1/C2/Go. Thus
**18 native engine captures**, six C/C comparisons and five strict C/Go
matches, not an overwritten last-invocation-only observation. C marker
rejection precedes accepting status; raw stdout/stderr/status and **every cwd
path/type/mode/empty directory/byte**, including the root, are compared with
the existing comparator and tree-modes.py, without normalization/exclusions.

```sh
. tools/lib/timeout.sh
export G3_REPO="$PWD/tmp/upstream.git" S4_ROOT="$PWD"
# Recreate/build the native sequential-capture launcher described above.
export S4_CAPTURE="$PWD/tmp/_issue220/$fixture-$run" # fresh absent run base
export S4_ENGINE="$PWD/tmp/_issue215/canonical/build/canonical/olympia-g3"
# Go uses a freshly built cmd/olympia, only after both C runs.
OLYMPIA_ENGINE_BIN="$PWD/tmp/_issue220/capture" \
  run_limited 150 bash /bin/bash \
  "$PWD/tmp/_issue215/canonical/tests/olympia/regress/$fixture/check.sh"
run_limited 120 compare "$PWD/tmp/_issue220/compare" --manifest "$capture"
run_limited 120 compare "$PWD/tmp/_issue220/compare" "$c1" "$c2"
run_limited 120 compare "$PWD/tmp/_issue220/compare" "$c1" "$go"
python3 -I tools/lib/tree-modes.py "$c1/cwd" "$go/cwd"
```

| Original fixture | Complete C cwd entries | Go engine / script | Strict result |
| --- | --- | --- | --- |
| move-ship | 239 | 0 / YES0 | C/C/Go exact, newly passing |
| movement-refused | 239 | 0 / YES0 | C/C/Go exact, newly passing |
| trade-refused | 237 | 2 / NO1 | C/C exact; actual `command buy start` NYI, #225 |
| wait-ferry-range | 235 | 0 / YES0 | C/C/Go exact, retained |
| wait-gm-stack | 237 | 0 / YES0 | C/C/Go exact, retained |
| wait-not-weather | 245 | 0 / YES0 | C/C/Go exact, retained |

All C engines/scripts exit0. move-ship proves the original fee/horn/wait/board/
out-and-back sail/unload turn, including original coverage, reports, saved
ship and JSON assertions. movement-refused proves every original verb's fail
and rejected cells. trade-refused remains an honest later BUY/SELL boundary:
no successful callback or rejection-only overlay is added. Family #186 still
has later combat/storm dependencies; this is not family/stage exit.

Strict manifest SHA-256s (C1/C2/Go identical except the honest trade boundary):

- move-ship: `3ef71115ed831ec414a0a9ae48616ba2532c91921a3deaca19a47f8f453e7f50`
- movement-refused: `f5d0e3b032eaaf50e8954a7fd8d598beae85507e5dc9a0edc1a6e86597e7f2c3`
- trade-refused C/C: `28511102be321d13de314f844ba833ea41b27a81623294bc2f4545ad8a7c623e`
- wait-ferry-range: `46c6090941300054159d3cca3693131baa37e4ed4c0e307fcc4b7a05ad719929`
- wait-gm-stack: `1d88382f463499fd192e1cb53d8898acfe033740edf4e921cfce2735f5a8fdd0`
- wait-not-weather: `b3addb02a6b6f2975e8c1230fde63ac83075a769be4fbbf96191f86cf3771885`

Logs/inventories/manifests: `tmp/issue220-fixtures.log`,
`tmp/issue220-FIXTURE-RUN.log`, `tmp/issue220-FIXTURE-RUN-invocations.txt`,
`tmp/issue220-FIXTURE-RUN-001-manifest.txt`,
`tmp/issue220-{canonical-provenance,source-definitions}.log`,
`tmp/issue220-{launcher,fixture-manifests}.sha256`. The decisive trade C/Go
captures remain under `tmp/issue220-trade-refused-captures/{c1-001,go-001}`;
full strict mismatch list: `tmp/issue220-trade-refused-go-001.diff`.
Matching captures, disposable WORK/RUN_DIR and task helpers/binaries are
removed; pristine C/cache and decisive evidence remain. Linux amd64 only,
no sanitizer/platform or full-corpus claim. See
[command proof](../command-compat/README.md#ferry-boarding-unload-and-fee-220-stage-4-n19)
and [spec no-loss](../olympia-spec/README.md#ferry-no-loss-220).

## SAIL originals (#219)

**0.141.0-alpha**, unchanged immutable **9a03d126**. Before native reuse,
all **1237 original blobs/modes** in `tmp/_issue215/canonical` were verified
against the bare pin. Canonical C retains SHA-256
`2ea19140ebd00e1ca672bd022ff958c139efb3aeca59f5f91b3cb8ebb93b5212`.
The shared additional checkout is never used. The protected
`tmp/issue291-roundtrip` cache remains idle and untouched; no in-place or
linked cache input is executed. Original scripts, README/scenarios and
expected bytes are read without edits, regeneration or `--update`.

The [bounded native pre-cleanup method](#recruitment-raise-and-refused-215)
is recreated under unique `tmp/_issue219`. The launcher opens **O_RDWR** raw
streams, preserves original argv0/stdin, sources timeout.sh in bash and calls
`run_limited 120 "$0" "$S4_ENGINE" "$@"`. It seeks/replays stdout then
stderr, closes them, writes actual status plus newline, then `cp -a .` to an
absent `cwd` **after replay and before original cleanup**. Thus script-owned
redirected `out`/`err` bytes are captured too. Scripts are bounded at **150**
seconds; native/comparator/test children at **120**. All invocations explicitly
set `G3_REPO="$PWD/tmp/upstream.git"`.

**Every engine invocation is captured separately**, not overwritten: each
fixture/run has an exclusive base and deterministic `001`, `002`, `003`
subdirectories, allocated in actual execution order. `raft-sail` runs
**sail/lose/lose-ship**; `raft-loc-depth` runs **load/attack**. Compare complete
per-run invocation inventories as well as each corresponding capture. No
`cp -a` destination already exists. Real C1/C2 run before Go for each script;
raw dependency/NYI/stub/sanitizer markers are rejected before accepting any C
status. Existing strict comparator plus `tree-modes.py` compare **all** cwd
paths/types/modes/empty directories/bytes and raw status/stdout/stderr.

```sh
. tools/lib/timeout.sh
export G3_REPO="$PWD/tmp/upstream.git" S4_ROOT="$PWD"
# Recreate/build the native sequential-capture launcher as described above.
export S4_CAPTURE="$PWD/tmp/_issue219/$fixture-$run" # fresh absent run base
export S4_ENGINE="$PWD/tmp/_issue215/canonical/build/canonical/olympia-g3"
# Go run only after both C repeats; use a freshly built cmd/olympia instead.
OLYMPIA_ENGINE_BIN="$PWD/tmp/_issue219/capture" \
  run_limited 150 bash /bin/bash \
  "$PWD/tmp/_issue215/canonical/tests/olympia/regress/$fixture/check.sh"
run_limited 120 compare "$PWD/tmp/_issue219/compare" --manifest "$capture"
run_limited 120 compare "$PWD/tmp/_issue219/compare" "$c1" "$c2"
run_limited 120 compare "$PWD/tmp/_issue219/compare" "$c1" "$go"
python3 -I tools/lib/tree-modes.py "$c1/cwd" "$go/cwd"
```

**30 native engine runs**, ten invocations for each C1/C2/Go set. All ten
C/C captures match, engine0 and original script YES0. Five corresponding Go
captures match strictly; five remain honest later-owner boundaries:

| Fixture / invocation | C cwd entries | Go engine / script | Strict C/Go / actual next stop |
| --- | --- | --- | --- |
| sail-out-first / 001 | 240 | 0 / YES0 | exact |
| raft-sail / 001 sail | 235 | 0 / YES0 | exact |
| raft-sail / 002 lose | 235 | 0 / YES0 | exact |
| raft-sail / 003 lose-ship | 235 | 0 / YES0 | exact |
| raft-loc-depth / 001 load | 227 | 0 / script1 overall | exact load case only |
| raft-loc-depth / 002 attack | 235 | 2 / 1 | `regular_combat` NYI, canonical combat owners through #235 |
| move-ship / 001 | 239 | 2 / 1 | `command fee start` NYI, #220 |
| sail-stale-storm / 001 | 237 | 2 / 1 | `USE v_dissipate callback` NYI, #265 |
| stop-running / 001 | 259 | 2 / 1 | `command repair interrupt` NYI, #223/#240 |
| movement-refused / 001 | 239 | 2 / 1 | `command board start` NYI, #220 |

Only **sail-out-first** and **raft-sail** are newly passing complete scripts.
raft-loc-depth's load-case parity is not whole-script acceptance. STOP's real
next NYI is now **repair interrupt**, not SAIL; both direct/USE SAIL interrupts
are genuinely owned and covered by the command oracle, not overlaid in this
fixture. Stale/live/moved adjacent bound-storm completion is independently
real-C/Go exact there; no future dissipate command is accelerated here.
No movement-family/stage exit or later command/mode implementation is claimed.

Matching capture manifest SHA-256s (three equal each):

- sail-out-first: `51f85264c5314f8bb771536f00b423141a2b9c4b655ec9623fdc42a9b96b0007`
- raft-sail sail: `7bfc2cdeec022e47eb30bf0859ef9cf2fa83e506e03023cf647802458359f750`
- raft-sail lose: `6b39bc9bdc7bd5180f442dafd8b5450bfe28bb63cd50113df17fec9d9058d19b`
- raft-sail lose-ship: `f8ba27099aaf8d47d03e673017ed1007eec33b8ca2b9ee911cac4260ae0276e8`
- raft-loc-depth load: `9edeb2e639d6e4228f1b1f55497eec42454e681ecdef612bf09665f39a899bb4`

Logs/manifests: `tmp/issue219-fixtures.log`,
`tmp/issue219-FIXTURE-{c1,c2,go}.log`,
`tmp/issue219-FIXTURE-RUN-ID-manifest.txt`,
`tmp/issue219-{fixture-manifests,launcher}.sha256` and
`tmp/issue219-canonical-provenance.log`. The five decisive C/Go failures are
retained under `tmp/issue219-FIXTURE-captures/{c1-ID,go-ID}`, with complete
mismatch lists `tmp/issue219-FIXTURE-ID.diff`. Matching captures, task helpers/
binaries and disposable WORKs are removed; pristine C/cache remain.
See [complete lifecycle proof](../command-compat/README.md#sail-219-stage-4-n18)
and [spec no-loss](../olympia-spec/README.md#sail-no-loss-219).

## FLY originals (#218)

**0.140.0-alpha**, immutable **9a03d126b329ed466c4af48425553c025eaca26f**.
Before reuse, all **1237 original source blobs/modes** in
`tmp/_issue215/canonical` match the bare pin; canonical C SHA-256 remains
`2ea19140ebd00e1ca672bd022ff958c139efb3aeca59f5f91b3cb8ebb93b5212`.
The original move-land/fly-out-first/movement-refused/save-fly-cap scripts,
READMEs, frozen scenarios and expected files are unchanged. No protected
`tmp/issue291-roundtrip` input/cache is executed or modified.

The [native pre-cleanup method](#recruitment-raise-and-refused-215) is
recreated under `tmp/_issue218`. The launcher opens **O_RDWR** raw streams,
retains original argv[0]/stdin, sources timeout.sh and calls
`run_limited 120 "$0" "$S4_ENGINE" "$@"`. It records the actual status,
seeks/replays stdout then stderr, closes both streams, writes status/newline,
then `cp -a . "$S4_CAPTURE/cwd"` after replay and before script cleanup.
Script-owned redirected streams are included. Launcher executable SHA-256:
`8ec21d42117b9310fcefc7c5ed7b4cc15e4f6030f77427a80e3589d56389114f`;
source SHA-256 `7a8c93161fa1117ef39c8876dea9160a3fe26f1cc91d437aaa37b5b816417acc`.

```sh
. tools/lib/timeout.sh
export G3_REPO="$PWD/tmp/upstream.git" S4_ROOT="$PWD"
# For each fixture below, run c1, c2, then go, never concurrently:
export S4_CAPTURE="$PWD/tmp/_issue218/$fixture-$run"
# c1/c2 use the verified canonical binary; go uses freshly built cmd/olympia.
export S4_ENGINE="$PWD/tmp/_issue215/canonical/build/canonical/olympia-g3"
OLYMPIA_ENGINE_BIN="$PWD/tmp/_issue218/capture" \
  run_limited 150 bash "$(command -v bash)" \
  "$PWD/tmp/_issue215/canonical/tests/olympia/regress/$fixture/check.sh"
# Reject raw dependency/NYI/stub/sanitizer/runtime-error markers in C BEFORE
# accepting its status. Manifest and compare with the existing comparator:
run_limited 120 compare "$PWD/tmp/_issue218/compare" --manifest "$S4_CAPTURE"
run_limited 120 compare "$PWD/tmp/_issue218/compare" "$c1" "$c2"
run_limited 120 compare "$PWD/tmp/_issue218/compare" "$c1" "$go"
python3 tools/lib/tree-modes.py "$c1/cwd" "$go/cwd"
```

| Original fixture | C1/C2/Go engine | C / Go scripts | Strict result |
| --- | --- | --- | --- |
| move-land | 0/0/0 | YES0 / YES0 | C/C and C/Go, **241 entries**, newly passing |
| fly-out-first | 0/0/0 | YES0 / YES0 | C/C and C/Go, **240 entries**, newly passing |
| movement-refused | 0/0/2 | YES0 / NO1 | C/C, **239 entries**; first Go stop is **command board start NYI** |
| save-fly-cap | 0/0/0 | YES0 / YES0 | C/C and C/Go, **226 entries**, retained |

Comparisons include raw status/streams and **all** cwd paths/types/modes/
empty directories/bytes, without exclusions or normalization. Both new
successes also pass the originals' own report/hash/state-export checks;
move-land passes every original movement-verb coverage assertion.
movement-refused still requires BOARD/FERRY/UNLOAD (#220) and SAIL (#219),
not another FLY overlay. Its Go result remains honestly red.

The three equal manifests per passing fixture have SHA-256:

- move-land: `a1a639c8730e7b197fb36ad1425efe79f8ff44088b9fd11697b5d72a2707519e`
- fly-out-first: `0d8c3e0f806cd49520677deead7f973f4a3ad8471b53395ac426c003f967f7a2`
- save-fly-cap: `5f9a1067b2b9505ba5fe06195397958749dc0fba9cb14abcb81e2344be16b69d`
- Both C movement-refused manifests: `f5d0e3b032eaaf50e8954a7fd8d598beae85507e5dc9a0edc1a6e86597e7f2c3`

Evidence: `tmp/issue218-fixtures.log`,
`tmp/issue218-FIXTURE-{c1,c2,go}.log`, corresponding `-manifest.txt`,
`tmp/issue218-{canonical-provenance,source-definitions}.log`, and
`tmp/issue218-{launcher,fixture-manifests}.sha256`.
Matching captures, task helpers/binaries and disposable WORKs are removed;
the decisive C/Go pair remains
`tmp/issue218-movement-refused-captures/{c1,go}`, with the full mismatch list
at `tmp/issue218-movement-refused.diff` and first stop at
`tmp/issue218-first-nyi.log`.
No pin change, new source dependency, full-corpus rerun, family/stage exit,
later gameplay or stage5 acceptance is implied.
See [callback proof](../command-compat/README.md#fly-218-stage-4-n17)
and [spec no-loss](../olympia-spec/README.md#fly-no-loss-218).

## TERRORIZE and INCITE start (#217)

**0.139.0-alpha**, unchanged **9a03d126b329ed466c4af48425553c025eaca26f**.
All **1237 original source blobs/modes** in `tmp/_issue215/canonical` are
verified before reuse; canonical C retains SHA-256
`2ea19140ebd00e1ca672bd022ff958c139efb3aeca59f5f91b3cb8ebb93b5212`.
Original social scripts/READMEs/scenarios/expected data are unchanged.

Recreate the [native pre-cleanup capture method](#recruitment-raise-and-refused-215)
under `tmp/_issue217`: O_RDWR streams, original argv0/stdin, bounded native
engine120/script150, real status and stdout-then-stderr replay before copying
the entire cwd. Real C runs twice before Go. Reject raw dependency/NYI/stub/
sanitizer/runtime-error markers before C status acceptance. Every invocation
sets `G3_REPO="$PWD/tmp/upstream.git"`; no protected cached input is executed.

| Original fixture | C1/C2/Go engine | Scripts C / Go | Strict result |
| --- | --- | --- | --- |
| social-switched | 0/0/0 | YES0 / YES0 | C/C and C/Go, **243 entries** |
| social-refused | 0/0/0 | YES0 / YES0 | C/C and C/Go, **241 entries**, newly passing |
| social-bribe | 0/0/2 | YES0 / NO1 | C/C, **245 entries**; Go reaches **USE d_incite finish NYI** |

Existing turn-compat comparator plus tree-modes.py compare raw status/streams
and **all** cwd paths/types/modes/empty directories/bytes, with no exclusions
or normalization. Original script coverage/hash/JSON checks pass on both
strict C/Go successes. social-refused's three capture manifests have SHA-256
`04726e1cb1e3371f15218a0fb1a754cbdea18360e84ac657bf8c4e2aacb38a48`;
social-switched retains
`9571120d8a7e2f100813cc53507da7eb4242033b89257da88dad7f28175d3ac8`.
Both C social-bribe manifests:
`b6e5024b460e695fd327327ead5ffa1e8be48dc5bc198233c2ee2e78daf36ef3`.

Logs/manifests: `tmp/issue217-fixtures.log`,
`tmp/issue217-social-FIXTURE-{c1,c2,go}.log` and `-manifest.txt`,
`tmp/issue217-canonical-provenance.log`, `tmp/issue217-launcher.sha256`.
Matching captures, task helpers/binaries and disposable exports are removed;
the decisive failed Go capture remains `tmp/_issue217/social-bribe-go`.
Its start is now real; #242 still owns complete USE675 finish/queued ATTACK
after combat. No success overlay, family #189 exit, stage exit, full serializer
rerun, pin change, later-owner acceleration or stage5 work is implied.
See [callback proof](../command-compat/README.md#terrorize-and-incite-start-217-stage-4-n16)
and [spec no-loss](../olympia-spec/README.md#terrorize-and-incite-no-loss-217).

## BRIBE and social fixtures (#216)

At **0.138.0-alpha**, immutable **9a03d126b329ed466c4af48425553c025eaca26f**
is unchanged. The original social README/Lua/scenario/check.sh/EXPECT files
are read, never modified. All **1237 original blobs and modes** in the
separate canonical export `tmp/_issue215/canonical` are reverified before
reusing its C binary, SHA-256
`2ea19140ebd00e1ca672bd022ff958c139efb3aeca59f5f91b3cb8ebb93b5212`.
No protected roundtrip input/cache is executed or needed for this proof.

The native pre-cleanup method documented below is recreated under
`tmp/_issue216`: O_RDWR stream captures, original argv[0]/stdin, bash sourcing
timeout.sh and `run_limited 120 "$0" "$S4_ENGINE" "$@"`; record actual status,
seek/replay stdout then stderr, close both, write status/newline, then
`cp -a . "$S4_CAPTURE/cwd"` **after replay, before original script cleanup**.
The unchanged scripts themselves run under `run_limited 150`. Real C runs
twice before Go; raw dependency/NYI/stub/sanitizer/runtime-error rejection
precedes C status acceptance. Launcher SHA-256:
`e61415bad36b03243802bb6119846b5bb80e624bf32100ea91c8981dcac948f9`.

| Original fixture | C1/C2/Go engine status | C scripts / Go script | Strict result |
| --- | --- | --- | --- |
| social-switched | 0/0/0 | YES 0 / YES 0 | C/C and C/Go, all **243** cwd entries identical |
| social-bribe | 0/0/2 | YES 0 / NO 1 | C/C, all **245** cwd entries identical; Go stops at INCITE start NYI |
| social-refused | 0/0/2 | YES 0 / NO 1 | C/C, all **241** cwd entries identical; Go stops at INCITE start NYI |

The existing strict comparator (`tools/turn-compat`), supplemented by
`tools/lib/tree-modes.py`, covers raw status/stdout/stderr and **every** cwd
path/type/mode/empty directory/file byte. There are no exclusions or
normalization, including script-owned stderr. All three social-switched
capture manifests have SHA-256
`9571120d8a7e2f100813cc53507da7eb4242033b89257da88dad7f28175d3ac8`;
its original coverage/report/state-export assertions pass too. Its delayed
contract-500 SWITCH, rival's finish refusal and latecomer's start refusal
are actual complete-engine behavior, not a callback overlay.

social-bribe/refused still require later INCITE/TERRORIZE and, for actual
USE675 success, combat. Their original scripts remain honestly red; no
partial-oracle success closes family #189. Matching captures and temporary
launchers/binaries/sources are removed; failed-Go captures and decisive logs/
manifests remain. Logs: `tmp/issue216-fixtures.log`,
`tmp/issue216-social-FIXTURE-{c1,c2,go}.log`, corresponding `-manifest.txt`,
`tmp/issue216-{canonical-provenance,source-definitions}.log`; failed-Go trees
are `tmp/_issue216/social-{bribe,refused}-go`.

No pin move, original fixture edit, full serializer/207-roundtrip rerun,
new dependency, family/stage-exit or remote-CI acceptance is implied.
See [callback/source proof](../command-compat/README.md#bribe-216-stage-4-n15)
and [spec no-loss](../olympia-spec/README.md#bribe-no-loss-gate-216).

## Recruitment raise and refused (#215)

At **0.137.0-alpha**, pin **9a03d126b329ed466c4af48425553c025eaca26f** stays
unchanged. Original recruitment-raise/refused README, Lua authoring source,
scenario patches, expected manifests/JSON and complete check.sh scripts were
read without modifying them. The unchanged scripts run real C twice, then
fresh Go through native OLYMPIA_ENGINE_BIN capture before script cleanup:

| Fixture | C1 / C2 / Go engine | Original scripts | Strict C/C and C/Go |
| --- | --- | --- | --- |
| recruitment-raise | 0 / 0 / 0 | YES / 0 each | 241 cwd entries identical |
| recruitment-refused | 0 / 0 / 0 | YES / 0 each | 239 cwd entries identical |

Both comparisons include raw status/stdout/stderr and **all** cwd paths,
types, modes, empty directories and file bytes through the existing
`tools/turn-compat` comparator and `tools/lib/tree-modes.py`. Nothing is
excluded or normalized. The original scripts' own coverage, manifest and
state-export assertions also pass. Matching captures are deleted only after
all three manifest hashes agree:

- recruitment-raise: `82008ab2fd77b7b21a318866ed012b96f04cbc49f2b5abef4b89f9e7578b6a9c`
- recruitment-refused: `7357ab6ae4ca0a3771763ece6028a56b2e6172eb6f50209abdd55e6fb1d43185`

Reference provenance: verify all **1237** original source blobs/modes and
the protected 207-input cache manifest before reuse. The accepted cache
`tmp/issue291-roundtrip` remains idle and unchanged: C SHA-256
`d8967f5700225d1f8c3714d203858f201508b7349c36c6daead8cd37b5cda634`,
input-manifest SHA-256
`af1453a48658ac26e97ece4d46df3345484897e11010b04314821668f01d5b53`.
No cached input is executed in place, linked or written. The original
check.sh unpacks fresh physical disposable inputs itself.

A separate pristine archive builds the canonical reference as follows; the
**build/canonical** directory is significant for reproducible DWARF bytes:

```sh
src="$PWD/tmp/_issue215/canonical"
# Export the immutable source with git archive into $src first.
cmake -S "$src" -B "$src/build/canonical" -G Ninja \
  -DCMAKE_BUILD_TYPE=Debug -DCMAKE_C_FLAGS="-ffile-prefix-map=$src/="
cmake --build "$src/build/canonical" --target olympia-g3 -j4
go build -o tmp/_issue215/compare ./tools/turn-compat
```

Canonical C SHA-256 is the accepted
`2ea19140ebd00e1ca672bd022ff958c139efb3aeca59f5f91b3cb8ebb93b5212`.
The trailing slash in the compiler map is retained; no assertion text is
rewritten later. The temporary native Go launcher lives under tmp/_issue215
(not a discoverable scratch test package); SHA-256
`740125c50805f59fbf761a25ea837da680802b7806bcbf7475789206ebec5032`.
It opens O_RDWR stdout/stderr captures, sources timeout.sh in bash and invokes
`run_limited 120 "$0" "$S4_ENGINE" "$@"`, preserving launcher argv[0]/stdin.
It records the real status, seeks/replays stdout then stderr, closes streams,
writes status plus newline, then `cp -a . "$S4_CAPTURE/cwd"` **after replay,
before script cleanup**, and exits that status. Script-owned redirected
stderr is therefore included too. The scripts themselves are bounded at
150 seconds. All compatibility calls explicitly use
`G3_REPO="$PWD/tmp/upstream.git"`.

Raw marker rejection precedes accepting C status or comparisons; the command
oracle separately exercises 24 compiled native 0/134 default/update rejection
controls. There are no temporary successful callbacks in the full engine.
The task's launcher, Go/compare binaries, marker controls, scratch Go sources
and matching captures are removed; verified original C/cache and decisive
failed-Go evidence, logs and manifests remain.

Logs: `tmp/issue215-fixtures.log`, `tmp/issue215-FIXTURE-{c1,c2,go}.log`,
`tmp/issue215-FIXTURE-{c1,c2,go}-manifest.txt`,
`tmp/issue215-fixture-manifests.sha256`, `tmp/issue215-cache-provenance.log`
and `tmp/issue215-canonical-build.log`. Full spec adds exactly the two
recruitment names, retaining all prior 56 and stage-3 36 names. See
[callback/source evidence](../command-compat/README.md#raiserally-215-stage-4-n14b)
and [spec gate](../olympia-spec/README.md#raise-rally-no-loss-gate-215).
This localized callback task does not rerun the already-accepted #291 four
full serializer/207-roundtrip corpora and makes no family/stage-exit claim.

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

## Merged-pin turn (#291)

At **0.136.3-alpha**, the immutable merged reference above passes the complete
existing gate, real C twice before Go/recording:

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue291-roundtrip" \
  WORK="$PWD/tmp/issue291-turn" RUN_TIMEOUT=120 \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/turn-compat/run.sh --full --record --command-coverage \
  --use-coverage --phase-coverage --state-export --dump-integration
```

The ordinary comparison preserves **233 cwd entries / status 0**. All
**13 lifecycle / 10 command / 15 USE / 23 phase / 13 export** controls and
**84 integration** controls (**74 status 0 / 10 status 1**) pass, plus **20
flag-free aftermath** comparisons. Every primary comparison includes all raw
streams/status and cwd paths/types/modes/empty directories/bytes. Full four-dump
hashes remain exactly those in [#177's table](#four-dump-integration-closure-177).
Only the pin and actual binary provenance change in c-determinism.txt; its
entire manifest body and small snapshots remain unchanged. New metadata SHA-256:
`7800d8c36851a21e2d819428f0eae17bd81e02a358c0b2030e8bb921a7a757e0`.

Logs and summaries: `tmp/issue291-turn-full.log`,
`tmp/issue291-turn/{c-determinism,lifecycle,command-coverage,use-coverage,
phase-coverage,state-coverage,integration-coverage}.txt`. Native original-script
captures, canonical compiler provenance and explicit pending ATTACK/fog,
console, sail-interrupt and quick-potion owners are in
[merged-pin proof](../README.md#merged-pin-proof-291). This does not accept a
whole gameplay family or stage exit.

## Follow-C PR-head turn and fixtures (#288)

**0.136.1-alpha**, immutable inspected **PR #365 snapshot 7064d6c**, OPEN at
initial recheck, not its later advanced/merged head (see the linked drift
notice below):

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue288-roundtrip" \
  WORK="$PWD/tmp/issue288-turn" \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/turn-compat/run.sh --full --record --command-coverage \
  --use-coverage --phase-coverage --state-export --dump-integration
```

Two unchanged C world turns precede recording/Go; the ordinary turn matches
**233 entries / status 0**. All **13 lifecycle / 10 command / 15 USE / 23
phase / 13 export** and **84 integration** controls pass, including **74 status
0 / 10 status 1** integration cases and **20 flag-free** integration aftermath
comparisons. No exclusions or normalization: raw status/streams and all cwd
types/modes/empty directories/bytes. Small recordings and the entire manifest
body are unchanged; **only pin and binary provenance** advance in
testdata/c-determinism.txt, and its test requires the reviewed pin.

Original scripts use the unchanged trusted native `tmp/stage4-capture` via
OLYMPIA_ENGINE_BIN, capturing before script cleanup. A fresh disposable C
archive is compiled with `-ffile-prefix-map=ABS_SOURCE/=`; assertion paths are
canonical at compilation, never rewritten in output. C SHA-256
`f48d888089c3a0b192abc4d90c84b5fda4c2a1e0eab82ac6b0d3da42236bd294`.
Each original script runs unchanged C twice then fresh Go:

| Fixture | C/C | Go / strict C/Go |
| --- | --- | --- |
| gm-refused (changed bare-argument coverage) | YES, 237 entries | YES, exact |
| gm-credit | YES, 237 entries | YES, exact |
| production-harvest | YES, 249 entries | YES, exact |
| production-refused | YES, 237 entries | YES, exact |
| player-settings / player-refused | YES, 237 entries each | YES, exact |
| stop-running (new frozen turn-1 input) | YES, 259 entries | NO: engine 2, command sail interrupt NYI |
| use-quick-cast (new) | YES, 246 entries | NO: engine 2, USE item key 6 NYI |

Both new scripts return **1** on Go; these are explicit pending gameplay,
not accepted statuses or reference recordings. Matching captures are removed;
decisive failure captures and manifests remain under `tmp/issue288-fixtures`.
The unchanged-Go gm-refused panic (engine 2/script 1) is separately retained
under `tmp/issue288-before-gm-refused`. Callback owners and full source review:
[follow-C gate](../README.md#follow-c-pr-head-gate-288).
Evidence: `tmp/issue288-{turn-full,fixtures,before,source-proof}.log` and
`tmp/issue288-turn/{c-determinism,lifecycle,*coverage}.txt`.

## Final-tier turn and original fixtures (#178)

At unchanged **0.115.1-alpha / f73e48118805ffe8ff5e763bfa4f20ae00c3a705**:

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue179-roundtrip" \
  WORK="$PWD/tmp/issue178-turn" \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  tools/turn-compat/run.sh --full --command-coverage --use-coverage \
  --phase-coverage --state-export --dump-integration
```

The ordinary turn passes **233 cwd entries / status 0**. All **13 lifecycle /
10 command / 15 USE / 23 phase / 13 export** controls pass, as do all **84**
integration controls (**74 status 0 / 10 status 1**) and their **20** flag-free
aftermath comparisons. Unchanged real C repeats before Go; unfiltered raw
streams/status and every cwd entry/type/mode/empty directory/regular byte
match, without normalization or exclusions. The four entire reference dumps
are byte-identical to accepted parent #177, retaining the **2116 / 5100 /
2523 / 16422-byte** sizes and hashes in the [four-dump table below](#four-dump-integration-closure-177).

All **six original phase scripts** also run themselves, real C twice before
Go, from a separate pristine bare-pin export, using the [native pre-cleanup
capture method](#drum-fixture-rerun-182): **phase-day0 252 / phase-daily 245 /
phase-month-open 242 / phase-month-gold 241 / phase-month-decay 246 /
phase-month-close 242** cwd entries. Original **savage-cap** likewise passes
C/C/Go, **238 entries**. Every script prints YES and every script/engine exits
**0**. Raw captures pass the same strict comparator and tree-modes.py before
deletion; no upstream script, fixture, EXPECT.json or expected manifest is
edited. Savage's script-owned stderr file is captured after raw stderr replay.
Native helper sources/launchers and matching captures are removed.

For each fixture in the order above and engine in `c1 c2 go` order, the
temporary launcher invocation was:

```sh
G3_REPO="$PWD/tmp/upstream.git" I178_CAPTURE="$capture" I178_ENGINE="$bin" \
  OLYMPIA_ENGINE_BIN="$PWD/tmp/issue178-fixture-launch" \
  bash "$source/tests/olympia/regress/$fixture/check.sh"
```

`source` was `tmp/issue178-fixture-source`, exported with git archive at the
pin; `capture` was the absolute `tmp/issue178-fixture-captures/$fixture/$engine`.
`bin` selected the verified cached C executable for c1/c2 and the turn work's
Go executable for go. The native helper preserves raw streams/status, uses
existing run_limited with **120s**, replays streams to the script's original
descriptors, then copies the complete cwd before returning. The source export
is removed after all seven comparisons.

Evidence: `tmp/issue178-turn-full.log`, `tmp/issue178-turn-summary.log`,
`tmp/issue178-turn/{c-determinism,lifecycle,command-coverage,use-coverage,phase-coverage,state-coverage,integration-coverage}.txt`,
the four full integration dump files, `tmp/issue178-fixtures.log` and
`tmp/issue178-fixture-captures/*/{*-script.log,*-manifest.txt,*.diff}`.
Linux amd64/GCC 12.2.0/Go 1.26.4 only; no native macOS/Windows claim.
These gates pass, and parent completed the initially withheld full-box gate
with verified physical copies via the existing GAMES override, without a
harness fix or exemption. Parent independently repeated the strict ordinary
turn, every coverage/export/lifecycle control, all **84** integration
controls and **20** aftermath comparisons; all four complete dump hashes
are unchanged. Evidence: `tmp/parent178-turn-full.log`,
`tmp/parent178-turn-summary.log`, `tmp/parent178-turn/`. See [accepted
stage-3 evidence](../README.md#final-tier-validation-178). Stage 3 is accepted
and paused; no stage-4 work is authorized.

## Four-dump integration closure (#177)

**0.115.1-alpha**, unchanged immutable reference
**f73e48118805ffe8ff5e763bfa4f20ae00c3a705**. This is harness work only:
`--dump-integration` extends the existing full-engine family loop, runners,
strict comparator and all-entry mode checks. It requires `--full`; default
checks remain small. No engine gameplay, generic helper layer, upstream
harness/fixture, expected byte, pin or inherited recording changes.

```sh
G3_REPO="$PWD/tmp/upstream.git" C_CACHE="$PWD/tmp/issue179-roundtrip" \
  WORK="$PWD/tmp/issue177-turn" tools/turn-compat/run.sh --full \
  --command-coverage --use-coverage --phase-coverage --state-export --dump-integration
```

### Source contracts and matrix

The complete original coverage.c (**453 lines / 12597 bytes**, SHA-256
`5730dc2b64b817a6ee1a1d9c250856a400f6aaf54b6ac50251b5f37e5e5ded5a`)
and state_export.c (**290 lines / 6496 bytes**, SHA-256
`2f710e5ec96c810996206c87efc7908fe50f75732b28f7119237e9ef1fb8210b`)
remain the writer definitions. All four return without opening for NULL;
otherwise `fopen(path, "w")` truncates an existing file, or prints exactly
`Can't write PATH!\n` and returns nonfatally. Empty writer paths therefore
fail to open, unlike NULL. Existing unchanged command-oracle cases retain
direct NULL/empty/repeated-write accounting, immediate suppression and byte
escaping coverage; CLI empty/missing values instead set errflag and return 1
before load. main.c **802–853** removes exact long options before BSD getopt,
consumes any nonempty next argument (including dash-prefixed paths), and the
last repeated assignment wins. An omitted option remains NULL.

main.c **1024–1045** calls command → USE → phase writers, then check_db,
optional save_db, optional save_logdir, state export, and finally Times
masthead/close. check_db and its complete checks were read; load_db also
checks at io.c **2906**. io.c **3045–3099** renames log into save/turn,
renames players into that saved directory and recreates log; save_db first
saves randseed, then routes/system/boxes/master/orders/admissions. rnd.c
**286–325** reads the seed and exits 1 on failed seed save. c2.c **339–472**
opens/flushes press and rumor streams early, writes times_0 late, and closes
the original FILEs. These real owners, not a modeled lifecycle, determine
overwrite/rename/cache behavior. All complete source-file hashes are in
`tmp/issue177-source-audit.log`. affected.sh already treats turn-compat as a
whole-engine/source consumer; no dependency rule or pin update is needed.

The **84 new controls** use pristine original **phase-day0** input, unpacked
from the isolated pin via original scenario_unpack. Each C reference runs
twice deterministically **before** Go. Final enabled path assignments use
reverse writer order to distinguish parse order from execution order.
Bits/indexes **0/1/2/3 = command/USE/phase/state**:

| Cases | Count | Contract exercised |
|---|---:|---|
| `subset-0` … `subset-15` | 16 | Every enabled/NULL combination, including completely flag-free |
| `shared-3,5,6,7,9,10,11,12,13,14,15` | 11 | Every pair/triple/four-way shared destination; later actual writer wins |
| `load`, `load-save`, `turn-nosave`, `inside` | 4 | All four at load/turn and save/no-save boundaries; destinations outside/inside lib |
| `overwrite`, `repeat` | 2 | Existing-file truncation; each option repeated with missing/unused prior destinations, only last used |
| `all-fail`, `save-fail`, `times-fail` | 3 | Four nonfatal open failures in writer order; fatal save stops before export; late Times failure after export |
| `{system,seed,report,log,saved,times,press,rumor}-{0..3}` | 32 | Each writer collides with saved DB/RNG, players, log/300, save/1/300, times_0/press/rumor |
| `{failure,empty,missing,dash}-{0..3}` | 16 | One failed writer with three successful writers; CLI empty/missing values; literal dash-prefixed destinations |

All **84** match C/C/Go status, unfiltered streams, **whole raw TSV/JSON**,
and every cwd entry/type/mode/empty directory/regular byte. **74** exit 0;
**10** exit 1 (eight parse failures plus seed-save and late Times failures).
Complete C manifests, including raw stream/status hashes, are retained in
`tmp/issue177-turn/integration-coverage.txt`. No exclusions or normalization.
For **20** controls (all subsets and the four load/save/inside controls),
flag-free C also repeats before flag-free Go. Only after the strict primary
comparison are requested dumps removed from disposable captures; every
remaining state/RNG/report/file/stream observation equals flag-free C.

The combined case is not an empty-root/zero-counter comparison: **8 regular
factions / 8 units / 8 faction-location records**, unit distribution
**1,0,1,1,2,1,1,1**; actual command `use` is **0/0/1/0**, actual USE item
**7/drum ok=1**, savage spawn and day-1 trade both **ran=1/acted=1**.
Supplemental cell/JSON assertions guard against vacuous input; they do not
replace the primary complete byte comparisons. Complete C dump evidence:

| `tmp/issue177-turn/` file | Bytes / rows | SHA-256 |
|---|---:|---|
| `integration-commands.tsv` | 2116 / 148 data rows | `1dade8bdfb65dc21476fa211be3573d7bafcad9c4d8280ffafa119fd6439696d` |
| `integration-use.tsv` | 5100 / 151 data rows | `b54a151a39fe67db0c4027c273496170c42a76083fa570d5f160083c6b245a9e` |
| `integration-phase.tsv` | 2523 / 74 data rows | `109a9151092d5a55ad68f3282558b0cff3d87f719fef077473f69276cb8eeeb1` |
| `integration-state.json` | 16422 | `455c499a16b469eff43f6aec998a8f6e6df0c80f4b911f647bbc9676174543db` |

### Rejection controls and independent fixtures

The unchanged strict ordinary turn still matches **233 entries**, and all
**13 lifecycle / 10 command / 15 USE / 23 phase / 13 export** controls pass.
Log: `tmp/issue177-turn-full.log`; existing family/lifecycle summaries remain
under `tmp/issue177-turn/`. Matching captures are deleted promptly.

`tmp/issue177-negative.log` exercises the actual runner and comparator, not
replacement successful gameplay. A temporary native pre-exec diagnostic
injector execs the **real unchanged C binary**: six marker families in either
stream on actual statuses **0 and 1** give **24 rejections (harness 2)** before
expected-status acceptance. Status 0 includes a real failed dump open; status
1 includes a real failed randseed save. Missing dependency, not implemented,
NYI, stub, AddressSanitizer and runtime-error markers are all rejected.
Six real-capture mutations independently reject JSON/TSV whitespace and a
directory mode (**1**), or missing stream/cwd and a symlink (**2**, capture
error rather than mismatch). Existing fast tests retain the broader negative
type/mode/error controls. No recording refresh is performed.

Three decisive failure reruns also repeat real C before Go and retain raw
streams/status and complete manifests in `tmp/issue177-negative/`:
`all-fail` **0 / 251 entries**, `save-fail` **1 / 251**, `times-fail`
**1 / 254**. all-fail stderr orders command/use/phase open diagnostics before
`check_db()`/`save_db()`, then the state diagnostic afterwards. The failed
seed save has all three TSVs but no JSON; late Times failure retains all four.

All **six original phase scripts** themselves run C twice then Go, using the
temporary native pre-cleanup capture method documented for #182.
All script/engine statuses are **0**, scripts print **YES**, complete raw
streams/status and cwd entries/types/modes/bytes match:
**day0 252 / daily 245 / month-open 242 / month-gold 241 / month-decay 246 /
month-close 242**. Logs: `tmp/issue177-fixtures.log`; original script logs,
complete manifests and strict diffs: `tmp/issue177-fixture-captures/`.
Temporary launchers, probes, fixture exports and matching cwd captures are
removed. No first-time gameplay prerequisite or product defect was found.

Bounded vet, fresh bounded full-enabled tests and the 14 relevant default
tools pass; unchanged full spec retains **36 pass / 89 fail / 0 error** (127s).
Golden remains **NO**, initialization **2** at unsupported -i. The entire
**578307 / 2134398-byte** inherited command input/output remains unchanged,
SHA-256 `3481dcd40549f16ff6fc67201a70bab1e00c08a599799151e7706734ff594659` /
`c4cc3c8ee0cd5f0b40ac6eb239aa7bf51235aaf1448fa6d9aad81f5a821f1bac`.
The read-only cache retains **1212** original files, executable SHA-256
`151c7355890eba7f8bfda1f4f16f4a490947aaa1db1ffcc3ca859cd659c01158`,
**201 pristine worlds / 41938 input hashes**. Every compat invocation uses
explicit isolated G3_REPO; every engine/test uses run_limited. Full tests
enable opt-in tests but do not supply corpus child-entry variables; the final
201-world/four-corpus and peripheral tier stays **#178**. Evidence here is
Linux amd64 only; no native Darwin/Windows verification is claimed.
This is **not stage exit**, whole-spec parity, or stage-4 authorization.

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
