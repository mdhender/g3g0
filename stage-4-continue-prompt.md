Continue stage 4 of mdhender/g3g0 on this larger orb. Do the work in this
thread, following the existing bounded-task plan. Do not restart stage-4
planning, recreate its issues, repeat accepted implementations, or start
stage 5. This prompt is for the owner to launch in a new thread; saving it
does not authorize its drafting thread to continue implementation.

## Read first and recover the live state

- Read AGENTS.md, WORKFLOW.md, stage-4-prompt.md and tools/README.md, especially
  **Follow-C PR-head gate (#288)**, the source-dependency/check tiers and the
  full-corpus disposable-copy workflow. The original stage prompt's starting
  version, open-issue list and planning state are historical, superseded below.
- Inspect local main, origin/main, the worktree and current version. Fetch
  history if shallow before relying on it. Preserve changes made by others.
  The source thread delivers the accepted #288 fix and this prompt to main;
  use the actual remote tip, not the old source thread's local checkout.
- Audit ALL open local issues, including unlabeled ones, and current upstream
  issues/main. Bug/follow-C work takes priority. Read #288's complete comments
  and #289 before any new gameplay task. Recheck changes since this handoff.
- Read #19, stage #183 and family trackers #184–193, with their comments.
  All 94 initial bounded tasks #194–287 already exist. The canonical index is
  https://github.com/mdhender/g3g0/issues/183#issuecomment-6493036887.
  The stage's earlier comments contain complete command, USE, fixture and
  source/dependency matrices. Recover those from GitHub; do not recreate them
  because old tmp files are absent. Current amendments and accepted task
  comments supersede historical pin/implementation-state descriptions.
- Read tools/command-compat/README.md, tools/turn-compat/README.md,
  tools/olympia-roundtrip/README.md and tools/olympia-spec/{README.md,passing.txt}.
  At the selected C pin, read the original testing/fixture/RNG/turn-execution
  documents, command/USE/item tables, complete callback definitions/callers,
  and actual frozen scripts/scenarios, not only backlogs or fixture names.

## Delivered state and source boundary

The source thread completed #194–214 sequentially, then the bounded original
scope of priority follow-C #288. Delivered engine version: **0.136.1-alpha**.
No #215 implementation was started. No stage/family exit or stage-5 acceptance
was claimed. Task/family/root progress is on GitHub, not inferred from this
historical snapshot.

The five affected tools (command, turn, CLI, roundtrip, spec) now default to
**7064d6c1767ba09772dc5e12af57719806184258**. This is the deliberately inspected
immutable PR #365 snapshot, not its later merged head. Unaffected tools retain
independent pins; do not globally override G3_REV or blanket-rebaseline.
The accepted #288 engine change is only the two missing-argument guards in
normal-turn CREDIT/RELORE; other behavior remains C's original contract.

Upstream PR playbymail/olympia-g3#365 advanced during validation and merged at
2026-10-29T20:09:31Z. Later head:
**3434e8ecd75f52d3df6d1e57809b215785992940**; merge:
**9a03d126b329ed466c4af48425553c025eaca26f**. Those additional changes are
**not incorporated**. They are explicitly split into priority bug/follow-C
**#289**, attached to #183, from the later owner comment on #288:
https://github.com/mdhender/g3g0/issues/288#issuecomment-6507018145.

## First priority: #289's source and scope gate

Assess #289 before resuming #215. Review the complete delta from 7064d6c to
the later head/merge and run tools/affected.sh with isolated G3_REPO before
choosing a new immutable pin. It adds ATTACK/USE missing-argument guards,
a fog naming guard, five immediate-GM guards, and original missing-argument /
gm-immediate-missing fixtures. A saved command's CO/li text is reparsed on
load; do not assume its argument slots began empty.

The existing USE dispatcher can be updated in its owner, but ATTACK and fog
have pending complete-contract tasks #235 and #264. The five console commands
need immediate mode, which #19 assigns to stage 5. Do not implement a
missing-argument-only ATTACK overlay, expand into combat/weather closures,
invent successful callbacks, or silently implement stage 5 to satisfy a
composite fixture. Inventory actual dependencies and split bounded work and
later-owner obligations before implementation. If immediate-mode parity now
conflicts with the authorized stage boundary, ask the owner for that specific
scope decision; neither this prompt nor the upstream fixture authorizes
stage 5. Do not claim #289 complete merely because the pin moved.

After the priority follow-C scope is resolved and its required bounded work
accepted, resume canonical gameplay at **#215 N14b RAISE/RALLY**, family #187.
Its prerequisite #214 already owns complete cookie allocation in npc.go.
Use the existing swear.go/npc.go/glob.go/USE owners; do not duplicate that
machinery. Read #215 and its complete pinned source before briefing its one
implementer. Then proceed sequentially through the existing index, splitting
newly discovered oversized prerequisites before coding.

## Accepted evidence to retain, not whole-stage parity

Linux amd64, GCC 12.2.0, Go 1.26.4. The #288 implementer and parent independently
passed bounded vet, fresh full-enabled tests, all eight sequential defaults
(command/move/helper/rng/code/loadsave/turn-driver/turn), affected CLI, the
205-roundtrip gate, strict ordinary turn/all dump integration controls and
the original-fixture no-loss gates. Exact commands, counts, hashes, source and
fixture ownership are durable in the tool READMEs.

- Command oracle: **1781 sections / 148098 lines**, input **742721 bytes**,
  SHA-256 `2a7bbb12574ea36135181b4b8249c8324306cb63d0ef2119aacfe6177fd9c6c9`;
  output **2880391 bytes**,
  SHA-256 `f71b3833e6e44b3571261e078bac705aca4a9edaeca32547b04c96527e9841cf`.
  Entire accepted prefixes, including **21 input / 67 output high-bit bytes**
  and authentic C trailing spaces, remain unchanged. Preserve them when
  appending, unless exact reviewed follow-C behavior explains specific changes.
  Stage-3 input/output prefixes are **578307 / 2134398 bytes**, hashes
  `3481dcd40549f16ff6fc67201a70bab1e00c08a599799151e7706734ff594659` /
  `c4cc3c8ee0cd5f0b40ac6eb239aa7bf51235aaf1448fa6d9aad81f5a821f1bac`.
- **205 deterministic load/save roundtrips**, zero failures/errors, plus
  **one separate EXPECTED_ABORT**. Real C release-swear-np is exactly **134 /
  swear.c (180): FALSE** at the current reviewed pin; no successful +turn
  snapshot exists. Never manufacture one. The abort exception is exact-pin,
  not a generic nonzero-status allowance; review source and entire fixture
  again if changing its pin.
- Strict ordinary turn: **233 cwd entries**, raw status/streams and every
  type/mode/empty directory/byte. All **13 lifecycle / 10 command / 15 USE /
  23 phase / 13 export / 84 integration** controls and **20 flag-free aftermath**
  comparisons pass. Preserve all four complete dump contracts and their order.
- Strict unchanged scripts gm-refused, gm-credit, production-harvest/refused
  and player-settings/refused match C/C/Go. New stop-running and use-quick-cast
  match C/C but honestly fail Go at **sail interrupt NYI** and **USE item key 6
  NYI**, engine 2/script 1. Owners: #219/#223/#240 for STOP closures;
  #256/#278/#281 for quick-cast's aura-view/Trance/potion closures. They are
  frozen inputs and require no Lua authoring to replay.
- Spec at 7064d6c: **56 pass / 71 fail / 0 error**, 127 fixtures; exact prior
  56 names and stage-3's 36 retained. Golden remains **NO**, initialization
  **2** at unsupported -i. Read passing.txt for names, not just a count.
  New merged fixtures are beyond this snapshot, not accepted or forgotten.
- Full serializer corpora were not rerun as #288 task acceptance. Earlier
  stage-3/final baseline evidence is historical, not a new-pin full-tier claim.
  Run the required pin-change/stage-exit tier on this orb; never substitute
  G3G0_FULL=1 without actual game lists for executing full compatibility tools.

Known platform exceptions remain red: scanf 78 glibc probes/11 macOS numeric
differences, exit 1; island 114/119 with five differences/seven failed checks,
exit 1. Do not rebaseline or call these passing. Upstream open issues at final
audit included G3 #363 (ctype include/shadowing), #309 (invalid RNG ranges,
won't fix) and #114 (wider export); recheck live status and relevance rather
than declaring any harmless or blocking without an actual failing case.

## Clean larger-orb setup and safety

Use an **a1.large (8 CPUs / 16 GB RAM)** orb. The new orb starts from remote
main with its own filesystem. Do not assume the old thread's untracked tmp
trees, caches, binaries, captures, logs or agents transferred. Tracked setup
installs Go from go.mod and CMake 4.1.3; inspect installed tools first.

Never build/write the shared additional C checkout or ~/Software reference.
Create an isolated bare clone if absent, and fetch missing immutable pins there:

```sh
mkdir -p tmp
git clone --bare https://github.com/playbymail/olympia-g3 tmp/upstream.git
```

For EVERY compatibility/affected invocation explicitly set
`G3_REPO="$PWD/tmp/upstream.git"`. Use git show or pristine git archive exports,
not the shared working tree. Rebuild absent C sources/corpora from unchanged
pinned fixtures, never Go output. Reuse a cache only after verifying its pin,
all original source/mode bytes, binary provenance and input manifest.

Use physical disposable input copies, no symlinks/hardlinks, and never run an
engine on protected cached input. Full box/boxfile/check/loadsave serializers
use the existing **GAMES override** pointing to disposable cases/*/input/lib
copies OUTSIDE WORK; see tools/README.md. Do not replace this solved workflow
with a harness change or policy exemption. Run shared-scratch tools sequentially.

Bound EVERY C/Go engine/test, alternate fixture launcher and expected abort
with tools/lib/timeout.sh run_limited or the existing run-limited.sh adapter:

```sh
. tools/lib/timeout.sh
run_limited 120 go "$(command -v go)" vet ./...
G3G0_FULL=1 RUN_TIMEOUT=120 go test -count=1 \
  -exec "bash $PWD/tools/lib/run-limited.sh" ./...
G3_REPO="$PWD/tmp/upstream.git" \
  GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
  bash tools/command-compat/run.sh
```

Original fixture proof uses unchanged check.sh with bounded native
OLYMPIA_ENGINE_BIN capture before script cleanup. Recreate the documented
capture method if its old temporary helper is absent. Require real C twice
before Go, raw status/stdout/stderr plus ALL cwd paths/types/modes/empty
directories/bytes via the existing strict comparator and tree-modes.py.
Reject dependency/NYI/stub/sanitizer markers BEFORE accepting expected statuses
or replacing recordings, including genuine 134 cases; exercise actual negative
controls. Expected bytes come only from unchanged C, never a rewritten model,
new upstream --update, normalized output or a fixture-specific overlay.

## Per-task process and stopping rule

Re-audit all open issues before each task. Exactly one implementing subagent
per bounded coding task, sequentially, leaves changes UNCOMMITTED. Parent
reads complete pinned definitions/diff, independently verifies, checks version
(minor for features, patch for fixes/harness; retain alpha), accepts and commits
directly to main with Closes #TASK only. Push and subscribe exact SHA/main CI;
no checks found is not remote green. Fresh-read/update root/stage/family
checklists without overwriting others' edits. No branches, PRs, tags or CI
workflows. Keep durable evidence/provenance/limits in existing docs; clean
task helpers, binaries, matching captures and generated coverage/use dumps.

Follow the original stage-exit final tier only after all in-scope contracts
and assigned original family fixtures are complete. Do not infer a family
exit from one start callback or partial oracle. Report every lost name,
unexplained difference and true blocker honestly. Finish stage 4, independently
accept it and update #19/#183/families, then PAUSE. Do not plan or implement
stage 5 without the owner's explicit authorization. No further thread is
needed merely for a phase change; keep this thread's implementation/proof here.
