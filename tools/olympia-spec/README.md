# olympia spec against the Go engine

`run.sh` runs G3's olympia spec, unchanged, against the Go engine (`cmd/olympia`), and records which gates pass in [`passing.txt`](passing.txt).
G3 is the reference.
Expected output comes only from G3's committed files, which its C engine wrote; nothing here re-baselines.

The current merged-pin no-loss proof and remaining stage-4/stage-5 boundaries
are recorded in [merged-pin proof](../README.md#merged-pin-proof-291).
Historical task results below describe their own pins, not whole-spec parity.

## Direct trade-entry no-loss (#225)

At **0.147.0-alpha**, unchanged immutable **9a03d126**, full unchanged spec
exits0: **129 fixtures /73 pass /56 fail /0 error (132s)**. Every **70 prior
name** and stage3's **36 names** is retained as an exact set. The only newly
passing original names are **`trade-market`**, **`trade-opium`** and
**`trade-refused`**. Only the existing tool writes passing.txt; no original
fixture, scenario, expected byte or original --update is used.

The complete original gold_each/v_buy/v_sell owner removes actual old-Go
BUY and SELL dispatch NYIs, without repeating accepted220's v_fee or changing
shared settlement, arrival, inventory, reports or monthly restock. Separate
[native original-script proof](../turn-compat/README.md#direct-trade-entry-original-fixture-proof-225)
runs unchanged C1,C2 **before Go** on all three plus retained seed parsing.
All **51 captures /17 branches** match raw streams/status and every cwd
path/type/mode/empty directory/regular byte, with no exclusions: market239,
refused237, opium235, seed-parse228 for001–013 and229 for014. The script's
entire14-invocation inventory is preserved;001 is not overwritten.

Original trade-routes remains **FAIL**, honestly **engine2 / script1** at
**command USE v_find_buy callback: not yet implemented** (#244). All five
branches repeat strict C/C successfully and fail strict C/Go; raw streams,
status, full manifests, inventories and failed-Go/C1 captures are retained.
This is not a new blocker or authorization for trade-route discovery, a
trade-family/stage exit, whole-spec parity or repeated full207 serializers.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/_issue225/spec" \
  RUN_DIR=/tmp/g3g0-issue225-spec-run RUN_TIMEOUT=120 \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

Tool-written passing.txt: **1430 bytes**, SHA-256
`5e82f52fcab18c59f92abde108db03f4b57d0253213474e325880f531e38975f`.
Stage3 source list: **876 bytes**, SHA-256
`6d2080c29961395d8387364983b6543a6d05c7277b63bdfd16f9025a91736b67`.
Golden remains **NO**, initialization **2**, unsupported `load_db (-i)`.
ATTACK235/fog264/full original missing-argument remain stage4; immediate-i,
five console guards/original gm-immediate-missing stay the explicit stage5
remainder. Supplied audit has no demonstrated new blocker. Evidence:
`tmp/issue225-spec.log`, `tmp/issue225-spec-names.log`,
`tmp/issue225-spec-logs/` and linked native proof. Matching captures,
temporary helpers/binaries/disposable WORKs/external RUN_DIR/generated
coverage/use dumps are removed; decisive logs/manifests/failed-Go evidence
and pristine canonical/cache are retained. Changes stay uncommitted for
parent independent acceptance/delivery; no GitHub writes or stage5 work.

## Raze no-loss (#224)

At **0.146.0-alpha**, unchanged immutable **9a03d126**, full unchanged spec:
**129 fixtures /70 pass /59 fail /0 error (133s)**, tool exit0. All **69 prior
names** and stage3's **36 names** are retained, checked as exact name sets.
The **only** newly passing original name is **`raze-rebuild`**. Only the existing
tool writes passing.txt; no expected fixture byte or original script is edited.

The complete RAZE owner removes the actual old Go **command raze start NYI**.
Original raze-rebuild now builds1375, razes it, builds higher3743, and satisfies
the real T_loc-chain/lib/loc check and whole-lib C manifest. Separate bounded
[native original-script proof](../turn-compat/README.md#raze-and-retained-construction-proof-224)
runs unchanged C1,C2 before Go for raze-rebuild plus retained construction-build,
construction-refused and stop-running. All **12 captures /four invocation
branches** match raw streams/status and every cwd entry/type/mode/empty
directory/byte, **236/242/245/259 entries**, with no exclusions or normalization.
This is localized acceptance evidence, not a construction-family/stage exit,
whole-spec parity or acceptance of unrelated command/USE callbacks.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/_issue224/spec" \
  RUN_DIR=/tmp/g3g0-issue224-spec RUN_TIMEOUT=120 \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

Tool-written passing.txt: **1391 bytes**, SHA-256
`c01eb3b256617f9840d4df4caa8b418c17751479d08a57d43350597c657cd71b`.
Stage3 source list: **876 bytes**, SHA-256
`6d2080c29961395d8387364983b6543a6d05c7277b63bdfd16f9025a91736b67`.
Golden remains honestly **NO**, initialization **2**, unsupported
`load_db (-i)`. ATTACK235/fog264/full original missing-argument remain stage4;
immediate mode/five console guards/original gm-immediate-missing remain the
explicit owner-approved stage5 remainder. No new blocker was demonstrated,
no console implementation/fixture overlay/original --update/new pin or repeat
full207 serializer corpus is introduced. Protected roundtrip cache stays idle.

Evidence: `tmp/issue224-spec.log`, `tmp/issue224-spec-names.log`,
`tmp/issue224-spec-logs/{run-all,olympia-g3,golden,build}.log`, plus native
inventories/manifests/comparison logs. Disposable WORK/external RUN_DIR,
temporary helpers/binaries, matching captures and generated coverage/use dumps
are removed; decisive failed-Go evidence and pristine canonical/cache remain.
Everything is left uncommitted for the parent's independent acceptance.

## Repair and improvement no-loss (#223)

At **0.145.0-alpha**, unchanged immutable **9a03d126**, full original spec:
**129 fixtures /69 pass /60 fail /0 error (132s)**, exit0. All **66 prior
names** and stage3's **36 names** remain, checked as sets rather than counts.
Exactly three newly passing original names:

- `construction-build`
- `construction-refused`
- `stop-running`

The anticipated stop-running USE boundary is not a failure in this run:
USE601's actual sailing callbacks already work after219, and complete
REPAIR223 removes the last boundary in this particular frozen scenario.
Other USE callbacks/STOP240 and RAZE224 remain their own obligations.
The separate [native original-script proof](../turn-compat/README.md#repair-and-improvement-original-fixture-proof-223)
runs C1,C2,Go, with all **nine** captures and all three complete cwd/raw
comparisons exact (**242/245/259 entries**). None of this infers a family or
stage exit, whole-spec parity or acceptance of other USE commands.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/_issue223/spec" \
  RUN_DIR=/tmp/g3g0-issue223-spec RUN_TIMEOUT=120 \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

Only the existing tool writes passing.txt: **1378 bytes**, SHA-256
`bb7120a4ded88e68bc4c3eca917aadd10d295db9ff7e6f344dd7a63b4e2d7735`.
Stage3's source list remains **876 bytes**, SHA-256
`6d2080c29961395d8387364983b6543a6d05c7277b63bdfd16f9025a91736b67`.
Golden stays honestly **NO**, initialization **2**, unsupported `load_db (-i)`.
No original fixture or oracle edits, upstream --update, fixture overlays,
new pin, stage5 implementation or full207 serializer-corpus rerun. #289's
ATTACK235/fog264/full missing-argument obligations and immediate/console
stage5 remainder are unchanged; no new upstream blocker was demonstrated.

Evidence: `tmp/issue223-spec.log`, `tmp/issue223-spec-names.log`,
`tmp/issue223-spec-logs/{run-all,olympia-g3,golden,build}.log`.
The disposable WORK/external RUN_DIR, temporary binaries/helpers, matching
captures and generated coverage/use dumps are removed after checks. The
pristine verified canonical source and protected idle roundtrip cache remain.

## Daily BUILD no-loss (#222)

At **0.144.0-alpha**, unchanged immutable **9a03d126**, full original spec:
**129 fixtures /66 pass /63 fail /0 error (128s)**, exit0. Exactly every
prior **66 name** and stage3 **36 names** remains, checked by name, not count.
No new passing name: construction scripts still need repair/improve223,
raft-loc-depth's attack branch needs regular_combat235. The separate strict
[original-fixture replay](../turn-compat/README.md#daily-build-original-fixture-proof-222)
ran all seven native invocation branches C1,C2,Go; raft-sail all three
branches stay exact. This is not a family/stage exit or whole-spec parity.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/_issue222/spec" \
  RUN_DIR=/tmp/g3g0-issue222-spec RUN_TIMEOUT=120 \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

Only the tool writes passing.txt, unchanged **1325 bytes**, SHA-256
`326cf4c3827bb0f8be7ab11094f900cda2ae50ceb81b55187b569198738c9f8d`.
Stage3 source list remains **876 bytes**, SHA-256
`6d2080c29961395d8387364983b6543a6d05c7277b63bdfd16f9025a91736b67`.
Golden remains **NO**, initialization **2**, unsupported `load_db (-i)`.
No fixture changes, new pin, console/stage5 implementation or full207
serializer/roundtrip rerun. #289's settled remaining owners are unchanged.

Evidence: `tmp/issue222-spec.log`, `tmp/issue222-spec-names.log`,
`tmp/issue222-spec-logs/{run-all,olympia-g3,golden,build}.log`.
Disposable spec WORK/external RUN_DIR, temporary helpers/binaries,
matching captures and generated coverage.tsv/use.tsv are removed. Decisive
logs/manifests and failed Go captures remain. Linux amd64, GCC12.2.0,
Go1.26.4; all engines use existing bounded run_limited/native launchers.

## Construction start no-loss (#221)

At **0.143.0-alpha**, unchanged immutable **9a03d126**, full original spec:
**129 fixtures /66 pass /63 fail /0 error (130s)**. No new construction
pass is expected: direct BUILD start/finish remain genuinely NYI until222,
repair/improve until223. Exactly all prior **66 names** and stage-3's
**36 names** are retained, checked by name, not count. The existing tool
alone writes passing.txt, unchanged **1325 bytes**, SHA-256
`326cf4c3827bb0f8be7ab11094f900cda2ae50ceb81b55187b569198738c9f8d`.
Stage-3 source list remains **876 bytes**, SHA-256
`6d2080c29961395d8387364983b6543a6d05c7277b63bdfd16f9025a91736b67`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/_issue221/spec" \
  RUN_DIR=/tmp/g3g0-issue221-spec RUN_TIMEOUT=120 \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

Golden remains **NO**, initialization exits **2** at unsupported `load_db
(-i)`. #289's settled canonical ATTACK235/fog264/full missing-argument
stage4 obligations and deferred stage5 console remainder are unchanged.
No pin/fixture edit, family/stage exit, stage5, GitHub write or delivery
acceptance. No new full207 serializer/roundtrip corpus claim; its cache
remains idle/unmodified. Bounded C/C/Go original construction scripts
remain [honest later-owner gates](../turn-compat/README.md#construction-start-prerequisite-221).

Evidence: `tmp/issue221-spec.log`, `tmp/issue221-spec-names.log`,
`tmp/issue221-spec-logs/{run-all,olympia-g3,golden,build}.log`.
Spec WORK/RUN_DIR, temporary helpers/binaries, matching captures and
generated coverage/use dumps are removed; decisive logs/manifests and
failed Go construction captures remain. Linux amd64, GCC12.2.0/Go1.26.4.

## Ferry no-loss (#220)

At **0.142.0-alpha**, unchanged immutable
**9a03d126b329ed466c4af48425553c025eaca26f**, the full original spec reports
**129 fixtures: 66 pass /63 fail /0 error (129s)**. Exactly **move-ship** and
**movement-refused** are newly passing; every prior **64 name** and all
stage-3 **36 names** are retained, compared by name rather than count.
The stage-3 source list remains **876 bytes**, SHA-256
`6d2080c29961395d8387364983b6543a6d05c7277b63bdfd16f9025a91736b67`.
Only the existing tool writes passing.txt: **1325 bytes**, SHA-256
`326cf4c3827bb0f8be7ab11094f900cda2ae50ceb81b55187b569198738c9f8d`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/_issue220/spec" \
  RUN_DIR=/tmp/g3g0-issue220-spec RUN_TIMEOUT=120 \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

Separate unchanged original scripts run real C twice before Go through the
[bounded native pre-cleanup capture](../turn-compat/README.md#ferry-originals-220),
including complete invocation inventories and raw status/streams/all cwd
paths/types/modes/empty directories/bytes. All six scripts invoke their engine
once: **18 captures**, six C/C matches, five strict C/Go matches. move-ship
and movement-refused match **239 cwd entries each**, engine0/scriptYES0;
the inherited wait-ferry-range/wait-gm-stack/wait-not-weather matches retain
**235/237/245** entries. trade-refused matches C/C **237** entries, but Go
honestly stops **engine2/scriptNO1** at actual **command buy start NYI**;
the complete BUY/SELL owner remains #225, not an overlay in #220.

Golden stays **NO**, initialization exits **2** at unsupported `load_db (-i)`.
The #289 canonical ATTACK235/fog264/full original missing-argument stage4
obligations and owner-deferred stage5 console remainder remain unchanged.
There is no pin/fixture/expected-data edit, movement-family/stage exit or
stage5 implementation. Linux amd64, GCC **12.2.0**, Go **1.26.4** only;
accepted #291 full207 serializer/roundtrip corpora are not rerun for this
localized callback task, and their protected cache remains untouched.

Evidence: `tmp/issue220-spec.log`, `tmp/issue220-spec-names.log`,
`tmp/issue220-spec-logs/{run-all,olympia-g3,golden,build}.log` and linked
native provenance/manifests. Matching captures, spec WORK/RUN_DIR, temporary
helpers/binaries and generated root/internal coverage/use dumps are removed;
pristine C/cache and decisive failed trade C/Go captures/logs are retained.
See [complete callback/source proof](../command-compat/README.md#ferry-boarding-unload-and-fee-220-stage-4-n19).

## SAIL no-loss (#219)

At **0.141.0-alpha**, unchanged immutable
**9a03d126b329ed466c4af48425553c025eaca26f**, the full original spec reports
**129 fixtures: 64 pass /65 fail /0 error (131s)**. Exactly **raft-sail** and
**sail-out-first** are newly passing. Every prior **62 name** and stage-3's
full **36-name** set is retained, checked by name rather than count. The
stage-3 source list remains **876 bytes**, SHA-256
`6d2080c29961395d8387364983b6543a6d05c7277b63bdfd16f9025a91736b67`.
The existing tool alone writes passing.txt: **1298 bytes**, SHA-256
`96ff0be54197af07a5ce93d31032d5038b11ceda8ac3e3e29c3489cca656e0aa`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/_issue219/spec" \
  RUN_DIR=/tmp/g3g0-issue219-spec RUN_TIMEOUT=120 \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

Separate unchanged original scripts use real C twice before Go and the
[bounded native pre-cleanup capture](../turn-compat/README.md#sail-originals-219).
All **30 engine invocations** are individually captured: ten per C1/C2/Go
set, including raft-sail's **three** cases and raft-loc-depth's **two**.
sail-out-first matches raw status/streams and **240 complete cwd entries**;
raft-sail matches all three **235-entry** cases. raft-loc-depth's load case
matches **227 entries**, but its attack case remains engine2/script1 at real
`regular_combat` NYI. It is **not** a passing complete script.

Honest later-owner failures remain: move-ship **FEE start** (#220),
movement-refused **BOARD start** (#220), sail-stale-storm **USE v_dissipate**
(#265), and stop-running now **REPAIR interrupt** (#223/#240), rather than
the formerly observed SAIL interrupt. These match C/C and remain red for
Go, not fixture-specific successful replacements or reasons to accelerate
later owners. Complete direct/USE601 SAIL and live/moved/deleted adjacent
bound-storm lifecycle have independent command-oracle parity.

Golden remains **NO**, initialization **2** at unsupported `load_db (-i)`.
#289's canonical ATTACK/fog/full original missing-argument obligations and
explicit stage5 console remainder are unchanged. No pin, original fixture or
expected bytes are edited, and no movement-family/stage exit is claimed.
Linux amd64, GCC **12.2.0**, Go **1.26.4** only. Accepted #291 full serializer/
207-roundtrip corpora are not rerun for this bounded callback task.

Logs: `tmp/issue219-spec.log`, `tmp/issue219-spec-names.log`,
`tmp/issue219-spec-logs/{run-all,olympia-g3,golden,build}.log`. Disposable spec
WORK/RUN_DIR, helpers/binaries and matching captures are removed; pristine C
and decisive failed captures/manifests/logs remain. See
[complete lifecycle/source proof](../command-compat/README.md#sail-219-stage-4-n18).

## FLY no-loss (#218)

At **0.140.0-alpha**, unchanged immutable
**9a03d126b329ed466c4af48425553c025eaca26f**, the full unchanged spec reports
**129 fixtures: 62 pass /67 fail /0 error (212s)**. Exactly **fly-out-first**
and **move-land** are newly passing. All prior **60 names** and the full
stage-3 **36-name** set remain; the latter source list verifies against its
accepted **876-byte** SHA-256
`6d2080c29961395d8387364983b6543a6d05c7277b63bdfd16f9025a91736b67`.
The existing tool alone writes passing.txt: **1273 bytes**, SHA-256
`52f944165625d4fc88c24e6da0ed498129660242dd5b9ac8b1e425154a9ac1fc`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/_issue218/spec" \
  RUN_DIR=/tmp/g3g0-issue218-spec RUN_TIMEOUT=120 \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

Separate unchanged original scripts execute real C twice before Go through
the [bounded native pre-cleanup capture](../turn-compat/README.md#fly-originals-218).
move-land/fly-out-first match raw status/streams and **241/240 complete cwd
entries**, engine0/scriptYES0. save-fly-cap retains its strict **226-entry**
success. movement-refused matches C/C (**239 entries**) but Go honestly stops
at **command board start NYI**, engine2/script1; #219/#220 still own the
remaining sailing/boarding/ferry/unload contracts. No callbacks are overlaid.

Golden remains **NO**, initialization **2** at unsupported `load_db (-i)`.
#289's canonical ATTACK/fog/full original missing-argument obligations and
owner-approved stage5 console remainder are unchanged. No fixture, expected
bytes or pin is edited; no movement-family or stage exit is claimed.
Linux amd64 only, GCC **12.2.0**, Go **1.26.4**; full serializer/207-roundtrip
corpora remain the accepted #291 evidence, not rerun for this local callback.

Logs: `tmp/issue218-spec.log`, `tmp/issue218-spec-names.log`,
`tmp/issue218-spec-logs/{run-all,olympia-g3,golden,build}.log`.
Task helpers/binaries, disposable exports and matching captures are removed;
decisive failed C/Go evidence is retained in the linked native-capture proof.
See [complete callback/source proof](../command-compat/README.md#fly-218-stage-4-n17).

## TERRORIZE and INCITE no-loss (#217)

At **0.139.0-alpha**, unchanged immutable **9a03d126b329ed466c4af48425553c025eaca26f**,
the full original spec reports **129 fixtures: 60 pass /69 fail /0 error
(195s)**. Exactly **social-refused** is newly passing. All prior **59 names**
and stage-3's **36 names** remain; the latter source list is verified against
its accepted 876-byte SHA-256. The tool alone writes passing.txt:
**1249 bytes**, SHA-256
`7d384e9540c11cbe94fc9cf71c3db4565ded05b0dacacb7ff9c4186551c92e63`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue217-spec" \
  RUN_TIMEOUT=120 GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

Separate unchanged original scripts run real C twice before Go through the
[bounded native pre-cleanup method](../turn-compat/README.md#terrorize-and-incite-start-217).
social-refused matches raw status/streams and **241 complete cwd entries**,
engine/script **0/YES0**; social-switched retains its strict **243-entry** pass.
social-bribe matches C/C (**245 entries**) but Go now reaches **USE d_incite
finish NYI**, engine2/script1. #242 still owns that completion and queued
ATTACK after combat. Direct INCITE's timer-only finish is genuinely NULL,
not a successful replacement for USE completion.

Golden remains **NO**, initialization **2** at unsupported `load_db (-i)`.
#289's ATTACK/fog and authorized stage5 console remainder are unchanged.
No upstream fixtures/expected bytes/pins change; family #189 and stage4 stay
open. Logs: `tmp/issue217-spec.log`, `tmp/issue217-spec-names.log`,
`tmp/issue217-spec/{run-all,olympia-g3,golden}.log`. Disposable exports/binaries
are removed after the run; decisive logs and failed-Go fixture evidence remain.

## BRIBE no-loss gate (#216)

At **0.138.0-alpha**, unchanged immutable **9a03d126b329ed466c4af48425553c025eaca26f**,
the full unchanged spec reports **129 fixtures: 59 pass / 70 fail / 0 error
(179s)**. Only **social-switched** is newly passing. Exact set checks retain
all prior **58 names** and every stage-3 **36** name (the stage-3 source list
is verified against its accepted 876-byte SHA-256). The tool alone writes
passing.txt: **1234 bytes**, SHA-256
`7a4fd07d17e2822c50d996912f1504d0bb2ae436424f191c73576a4aba5b633e`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/_issue216/spec" \
  RUN_DIR=/tmp/g3g0-issue216-spec RUN_TIMEOUT=120 \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

Separate original-script proof runs real C twice before Go; social-switched
matches raw status/streams and all **243 cwd entries**, paths/types/modes/
empty directories/bytes, with original script **YES / 0**, engine **0**.
social-bribe/refused each pass C/C but Go honestly stops at **INCITE start
NYI**, engine **2** / script **1**; their later social/combat owners remain.
Golden remains **NO**, initialization exits **2** at unsupported `load_db (-i)`.
#289's pending ATTACK/fog and explicit stage-5 console remainder are unchanged.
No fixture/expected output is edited and no family/stage-exit is claimed.

Logs: `tmp/issue216-spec.log`, `tmp/_issue216/spec/{run-all,olympia-g3,golden}.log`,
`tmp/issue216-spec-names.log`. Temporary export/launcher sources and binaries
are removed; logs remain. See [native capture evidence](../turn-compat/README.md#bribe-and-social-fixtures-216)
and [complete callback evidence](../command-compat/README.md#bribe-216-stage-4-n15).

## RAISE RALLY no-loss gate (#215)

At **0.137.0-alpha**, unchanged immutable **9a03d126b329ed466c4af48425553c025eaca26f**,
the existing full original spec reports **129 fixtures: 58 pass / 71 fail /
0 error (156s)**. The only new names are **recruitment-raise** and
**recruitment-refused**. Every prior **56-name** and stage-3 **36-name** entry
is retained, checked as exact sets, not just counts. The tool alone writes
passing.txt: **1218 bytes**, SHA-256
`4e0c70e4454892cc8036db90eed178c33aea53d51389cde5acf4164ce46b6aa0`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue215-spec" \
  RUN_DIR=/tmp/g3g0-issue215-spec RUN_TIMEOUT=120 \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

The generated alternate-engine launcher invokes the existing run_limited
supervisor with the original argv[0]/stdin. Its scratch Go source/binaries
and disposable export are removed after the run; use an underscore-prefixed
WORK for future invocations to avoid even a transient extra scratch Go package.
Separate unchanged native-script proof ran real C twice before Go for both
newly passing fixtures: script **YES / 0**, engine **0**, strict **241 / 239
cwd entries**, raw streams/status and every path/type/mode/empty directory/byte.
See [capture provenance](../turn-compat/README.md#recruitment-raise-and-refused-215)
and [complete callback evidence](../command-compat/README.md#raiserally-215-stage-4-n14b).

Golden remains **NO**, initialization exits **2** at unsupported `load_db (-i)`.
Pending #289 ATTACK/fog owners and owner-approved console remainder remain
unchanged. This is a no-loss/local callback gate, not full-spec, family or
stage-exit acceptance. No original fixture, expected bytes or pin is edited.
Logs: `tmp/issue215-spec.log`, `tmp/issue215-spec/{run-all,olympia-g3,golden}.log`,
and `tmp/issue215-spec-names.log`.

## Merged-pin no-loss gate (#291)

At **0.136.3-alpha**, immutable **9a03d126b329ed466c4af48425553c025eaca26f**,
the existing full tool reports **129 fixtures: 56 pass / 73 fail / 0 error
(194s)**. All **56 accepted names** and the full stage-3 **36-name** set remain;
only the provenance line changes in passing.txt, written by the tool, never
hand-edited. Its SHA-256 is
`d070fda317b046dc1bedb06a73d0b63ad874ed877917d9eb3a36de241d89c8c3`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue291-spec" \
  RUN_DIR=/tmp/g3g0-issue291-spec RUN_TIMEOUT=120 \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

The two new originals add failures, not passes: missing-argument first reaches
`command USE v_summon_fog callback` NYI; complete #264/#235 and strict fixture
acceptance remain mandatory before stage-4 exit. gm-immediate-missing reaches
unsupported `load_db (-i)` and remains explicitly stage 5. Both have strict
real C/C proof before the Go run; neither C/Go comparison passes. STOP sail interruption
and quick-cast item6 remain NYI too. No fixture, expected byte or callback is
changed. Golden stays **NO**, initialization **2**. Tool exit 0 means no
accepted name was lost, not whole-spec, family or stage completion. Full C
spec was not rerun; separate C/C/native original-script proof covers the 18
selected fixtures documented in the linked proof.

Logs: `tmp/issue291-spec.log` and
`tmp/issue291-spec/{run-all,olympia-g3,golden}.log`.

## Follow-C PR-head no-loss gate (#288)

At **0.136.1-alpha**, immutable **7064d6c1767ba09772dc5e12af57719806184258**
(inspected **PR #365 snapshot**, OPEN at initial recheck, not its later
advanced/merged head; see the linked drift notice below), the full spec has
**127 fixtures: 56 pass / 71 fail / 0 error (139s)**. All **56 prior names**
and the complete stage-3 **36-name** set remain. The existing tool alone writes
passing.txt: only its provenance pin changes, SHA-256
`ef62cb05244e9d80e6c12d82c3f55f0222633d60ce6aaeffc9bdf13beeda75b8`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/_issue288/spec" \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

The changed gm-refused script now genuinely passes after following C's two
guards. Both new original scripts pass unchanged C twice, but **fail Go**:
stop-running reaches `command sail interrupt: not yet implemented`; use-quick-cast
reaches `USE item key 6: not yet implemented` (engine 2 / script 1 each).
These add no passing name. Full fixture/callback ownership, strict captures
and retained production/player controls are documented in
[the pin gate](../README.md#follow-c-pr-head-gate-288).

Golden remains **NO**, initialization **2** at unsupported `load_db (-i)`.
This is a no-loss gate, not whole-spec parity, later-stage completion, or a
claim that a new fixture passes because its source pin moved.
Logs: `tmp/issue288-spec.log` and `tmp/issue288-spec-{run-all,olympia-g3,golden}.log`.

## Stack command progress (#201, stage 4 N04)

At **0.123.0-alpha**, unchanged pin **f73e48118805ffe8ff5e763bfa4f20ae00c3a705**,
the full original spec records **41 pass / 84 fail / 0 error**, **125 fixtures /
122s**, preserving every previous **38** pass and adding only
**stacking-stack, stacking-refused, release-swear-np**. The existing tool
generated passing.txt; no upstream scenario, script or expected bytes change.
Its SHA-256 is
`6771ff03ab00bd8b9f3b6b85482b1891f04800c18be96fceefb2108895cb464f`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue201-spec" \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

Parent independently repeats both stacking scripts, move-go, move-out-first
and release-swear-np with unchanged real C twice then fresh Go. Every original
script prints **YES / exits 0**. Complete pre-cleanup raw status/streams and
cwd types/modes/empty directories/bytes match **241/239/237/240/229 entries**.
The last script expects engine **134**, exact `olympia/swear.c (180): FALSE`;
its success does **not** create a successful turn/roundtrip snapshot. C source
paths are canonicalized at compilation, not by filtering capture output.
See [stack acceptance](../command-compat/README.md#stack-commands-201-stage-4-n04).

Logs: `tmp/issue201-spec.log`, `tmp/issue201-spec/{run-all,olympia-g3,golden}.log`,
`tmp/parent201-FIXTURE-{c1,c2,go}.log`; strict manifests remain under
`tmp/parent201-fixtures/`. Golden remains **NO**, unsupported initialization
exits **2**. The **84 remaining failures**, mixed family #185 and stage 4
remain unaccepted.

## Attitude command progress (#199, stage 4 N02)

At **0.121.0-alpha**, unchanged pin **f73e48118805ffe8ff5e763bfa4f20ae00c3a705**,
the full original spec reports **38 pass / 87 fail / 0 error**, **125 fixtures /
121s**. The tool preserves all **36** accepted names and records only the two
genuine new passes: **attitude-set** and **attitude-refused**. No upstream
scenario, script, manifest, expected state or expected byte changes.
passing.txt SHA-256:
`ecdcdc06667ba5c2aea3b1d930ef5c01391e384f379cd2f3ca23c4c9a44c9923`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/parent199-spec" \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/olympia-spec/run.sh
```

Parent also runs both original scripts against unchanged real C twice then
fresh Go, all **YES / status 0**. Strict pre-cleanup raw status/streams and
whole-cwd comparison matches **239 entries** for attitude-set and **237** for
attitude-refused, including types, modes, empty directories and bytes.
The first scenario deliberately uses a safe haven; these fixtures do not
require the still-pending general combat implementation. See
[command acceptance evidence](../command-compat/README.md#admission-attitudes-and-contact-199-stage-4-n02).

Logs: `tmp/parent199-spec.log`, `tmp/parent199-spec/{run-all,olympia-g3,golden}.log`,
`tmp/parent199-attitude-{set,refused}-{c1,c2,go}.log`; retained strict manifests
under `tmp/parent199-fixtures/`. Matching captures and temporary fresh binary
are removed. Golden remains **NO**, initialization exits **2** at unsupported
`-i`; **87 failing fixtures** remain unaccepted. Neither mixed family #185
nor stage 4 is complete.

## Final-tier no-loss repeat (#178)

At unchanged **0.115.1-alpha**, immutable
**f73e48118805ffe8ff5e763bfa4f20ae00c3a705**, the original full spec repeats
**36 pass / 89 fail / 0 error**, **125 fixtures / 119s**. The no-loss tool
exits **0**, with **no recorded-list change**; this is not whole-spec parity.
The complete **876-byte** passing.txt remains SHA-256
`6d2080c29961395d8387364983b6543a6d05c7277b63bdfd16f9025a91736b67`.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue178-spec" \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  tools/olympia-spec/run.sh
```

All six phase fixtures and savage-cap also independently pass their original
scripts with strict pre-cleanup C/C/Go capture; see [the current turn evidence](../turn-compat#final-tier-turn-and-original-fixtures-178).
Golden remains **NO**, initialization **2** at unsupported `-i`, explicitly
not a stage-3 blocker. No expected bytes or harnesses are changed. Logs:
`tmp/issue178-spec.log`, `tmp/issue178-spec/{run-all,olympia-g3,golden}.log`.
The initially withheld full-box gate now passes using parent's verified
physical copies and the existing GAMES override; no fix or exemption was
needed. Parent's independent spec repeat again gives **36 pass / 89 fail /
0 error**, **121s**, no passing-list change and golden **NO**. See [accepted
stage-3 evidence](../README.md#final-tier-validation-178). Stage 3 is accepted
and paused; the **89 failures and golden NO remain disclosed**.

## Drum prerequisite progress (#182)

At **0.115.0-alpha**, the unchanged complete spec records **36 pass / 89 fail /
0 error** (**119s**, **125 fixtures**), preserving every **34** accepted pass
and adding only **phase-day0** and **savage-cap**. The existing tool adds those
two actual passing check-script names to passing.txt; no expected byte,
scenario, manifest, EXPECT.json or upstream harness is edited. savage-cap's
original turn reaches the real drum callback after day-zero spawning; it does
not require implementing SUMMON-SAVAGE or combat.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue182-spec" \
  tools/olympia-spec/run.sh
```

All six original phase scripts separately pass unchanged real C twice then
Go, with exact raw status/streams and every cwd entry/type/mode/byte captured
before cleanup. phase-day0 newly matches **252 entries**; all five previously
passing phase fixtures remain identical. [Strict fixture evidence](../turn-compat#drum-fixture-rerun-182).
The final native-capture repeat additionally runs unchanged **savage-cap**
against C twice then Go, all **238 cwd entries** and raw streams/status exact.
No additional missing gameplay boundary is exposed. Golden remains **NO**:
initialization exits **2** at unsupported `load_db (-i)`. Remaining **89**
failures are not accepted parity. #177/#178 remain open; neither the six
phase passes nor this no-loss spec progress establishes stage exit or
whole-spec parity. The 201-world/four-corpus final tier remains #178.

Logs: `tmp/issue182-spec.log`,
`tmp/issue182-spec/{run-all,olympia-g3,golden}.log`,
`tmp/issue182-fixtures{,-final}.log`, `tmp/issue182-fixture-captures/`.

## Hinder-meditation prerequisite progress (#181)

At **0.114.0-alpha**, the unchanged complete spec records **34 pass / 91 fail /
0 error** (**119s**), preserving all **33** accepted passes and adding only
**phase-month-decay**. The existing tool writes that single new progress-list
name; no upstream fixture, script, manifest, EXPECT.json or expected byte is
edited. All six phase scripts separately run real C twice then Go under the
strict raw-stream/status/all-cwd-entry comparison: month-decay newly matches
all **246 entries**, and the other four previously passing phase fixtures
remain identical. [Fixture evidence](../turn-compat#hinder-meditation-fixture-rerun-181).

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue181-spec" \
  tools/olympia-spec/run.sh
```

phase-day0 still stops at separately filed **#182**, USE item key 7/drum;
it is not accepted or implemented here. Five of six phase fixtures now pass,
not six. Golden remains **NO**, initialization exits **2** at unsupported
`load_db (-i)`. No whole-spec or stage-exit parity is claimed. Full 201-world/
four-corpus acceptance is deferred to #178 as instructed.

Logs: `tmp/issue181-spec.log`,
`tmp/issue181-spec/{run-all,olympia-g3,golden}.log`,
`tmp/issue181-fixtures.log`, `tmp/issue181-fixture-captures/`.
Temporary fixture launchers/probes and matching captures are removed; the
decisive phase-day0 C/Go captures are retained. No newly exposed callback
beyond the already filed drum prerequisite was encountered.

## GUARD/STACK prerequisite progress (#180)

At **0.113.0-alpha**, the unchanged complete spec records **33 pass / 92 fail /
0 error** (**119s**), preserving every previously accepted **32** pass,
including all inherited 23. The existing tool adds only **phase-daily** to
passing.txt; no manifest, EXPECT.json, scenario, upstream script or expected
byte is edited. That fixture separately passes real C twice and Go with
identical raw status/streams and all **245 cwd entries**, including modes.

```sh
G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue180-spec" \
  tools/olympia-spec/run.sh
```

phase-day0 still fails: GUARD now executes, but Go later stops at **USE item
key 7**, original savage.c `v_use_drum` **70–146**. Real C twice passes its
unchanged script. That separately sized prerequisite is documented for the
parent to file before expansion, not implemented or accepted as success.
phase-month-decay's #181 hinder-meditation boundary is unchanged. The three
accepted month fixtures and all other recorded passes remain green; only
four of the six phase fixtures currently pass. Golden remains **NO**, with
initialization **exit 2** at unsupported `load_db (-i)`.

Logs: `tmp/issue180-spec.log`, `tmp/issue180-spec/{run-all,olympia-g3,golden}.log`;
strict affected-fixture evidence and the additional source/sizing boundary:
[turn comparison](../turn-compat#guardstack-fixture-rerun-180),
`tmp/issue180-{fixtures,prerequisite}.log`. No whole-spec, all-six fixture,
full-corpus, native-platform or stage-exit claim is made.

The spec is the three gates that G3 documents for another engine in `doc/claude/testing.md` ("Running the spec against another engine"):

| Gate | What it checks |
|---|---|
| `tests/olympia/regress/run-all.sh` | every regress fixture's `check.sh` (not `check-lua.sh`); one `PASS`/`FAIL`/`ERROR` line per fixture |
| `run/olympia-g3.sh` | the full golden turn (`-i`, then `-r -S`) from `lib.tgz`, written to `OLYMPIA_RUN_DIR` |
| `tests/olympia/golden_check.sh` | that turn's `lib` against `tests/olympia/golden/manifest.sha256`; prints `YES` |

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `9a03d126b329ed466c4af48425553c025eaca26f`.
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

## Location export integration (#176)

At **0.112.0-alpha**, the unchanged full spec runs **125 fixtures** and exits
**0**: **32 pass / 93 fail / 0 error** in **118 seconds**. Every inherited
**23** pass remains. The tool itself adds exactly nine genuine unchanged
script passes to `passing.txt`: `coast-rocks`, `move-go`, `move-out-first`,
`phase-month-close`, `phase-month-gold`, `phase-month-open`, `player-refused`,
`wait-gm-stack`, `wait-not-weather`. Expected files/harnesses are not edited.
Command: `G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue176-spec"
tools/olympia-spec/run.sh`; log: `tmp/issue176-spec.log`, full raw script
results under `tmp/issue176-spec/`. Golden initialization still exits **2**
at unsupported `-i`; golden is **NO**, not a passing spec.

All six phase fixtures separately ran their original scripts against real C
twice then Go with bounded engines and pre-cleanup raw status/stream/complete
cwd capture. C/C matches for all six; C/Go matches fully for **month-open,
month-gold and month-close**, including all JSON bytes and **242/241/242** cwd
entries/types/modes/bytes. **day0, daily and month-decay remain NO**: their
Go engines stop at **guard**, **stack** and **USE v_hinder_med** respectively
(engine **2**, script **1**; both C scripts are YES/**0**). These first
observed missing-gameplay prerequisites remain outside export, not rewritten
expectations or successful substitutes. See [exact outcomes and prerequisite
scope](../turn-compat#state-export-save-lifecycle-176),
`tmp/issue176-fixtures.log` and `tmp/issue176-prerequisites.log`.
No six-fixture, full-C-spec, native-platform or stage-exit acceptance claimed.

## Results

Follow-C #179 run: 2026-10-29, Linux amd64, Go 1.26.4, immutable
`f73e48118805ffe8ff5e763bfa4f20ae00c3a705`:

- `G3_REPO="$PWD/tmp/upstream.git" tools/olympia-spec/run.sh` exits 0 with
  **23 pass / 102 fail / 0 error** (80 seconds, **125 fixtures**). All **21**
  inherited passes remain. Only the genuinely passing new unchanged checks
  `load-trade-no-price` and `summary-trunc-name` are added to `passing.txt`,
  along with the new provenance pin. No upstream expected bytes are rewritten.
- Those two check-only fixtures add no saved-game corpus names. Both also
  pass unchanged on real C twice; all eight real-turn cases match exact raw
  C/C/Go status, streams and complete **239-entry** cwd trees. See
  `tmp/issue179-fixtures.log` and [the pin gate](../README.md#follow-c-pin-gate-179).
- Golden remains **NO**: initialization exits **2** at explicitly unsupported
  `load_db (-i)`. The 102 existing failures, remaining USE handlers, phase
  coverage, state export and later commands/modes remain unimplemented; this
  is a no-lost-passes check, not whole-spec parity or stage-3 completion.
- Logs: `tmp/issue179-spec.log` and
  `tmp/olympia-spec/{run-all,olympia-g3,golden}.log`. Full C spec and macOS
  acceptance were not rerun. Every Go engine uses the native bounded launcher.

Final #166 run: 2026-10-28, Linux amd64, Go 1.26.4, then-pin
`e733acb63bacb21167a1f2589d8dcdc023364896`.
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
