Continue stage 4 of mdhender/g3g0 in this new thread, following its existing
bounded-task plan. Resume at **#227 N25 RESEARCH/FORGET**, after a fresh live
audit and environment baseline. Do not restart planning, recreate issues,
repeat accepted implementations, or start stage 5. Finish stage 4, independently
accept it, update its trackers, then PAUSE.

This handoff is for the owner to launch in a new thread. Saving/pushing it does
not authorize the drafting thread to resume work. The owner explicitly asked
that thread to finish its current task, push this prompt, and stop. It completed
and delivered #226; **#227 has NOT been started**, and no implementer is running.
Source thread: https://ampcode.com/threads/T-01a18983-9e79-77eb-a846-70ffa87221b6.

## Recover the live state before implementation

- Read AGENTS.md, WORKFLOW.md, stage-4-prompt.md and tools/README.md. The
  original stage prompt's starting version/pin/planning state is historical;
  current accepted task comments and this handoff supersede it.
- Inspect local main, origin/main, worktree and version; fetch missing history
  if shallow before using git log/blame. Preserve others' changes. Use actual
  remote main, not the source thread's old local checkout.
- Audit ALL open local issues, including unlabeled ones, and upstream issues
  and main. Bug/follow-C work takes priority. At handoff, only the settled
  remainder of #289 is open local bug/follow-C work. Upstream main is still
  **9a03d126b329ed466c4af48425553c025eaca26f**; open G3 #363 (ctype headers),
  #309 (invalid RNG ranges/overflow, wontfix) and #114 (wider export) have no
  demonstrated new blocker. Recheck, rather than assuming they are harmless.
- Read #19, stage #183, family trackers #184–193 and their comments. All 94
  initial bounded tasks #194–287 exist. Canonical task/dependency index:
  https://github.com/mdhender/g3g0/issues/183#issuecomment-6493036887.
  Earlier #183 comments contain full command/USE/item/fixture/dependency
  matrices. Recover them from GitHub, not absent tmp files; do not recreate
  the plan. Fresh-read tracker bodies before updating and preserve others' edits.
- Read #288's complete comments, #289, accepted #290/#291 and latest #226
  acceptance. Read tools/command-compat/README.md, tools/turn-compat/README.md,
  tools/olympia-roundtrip/README.md and tools/olympia-spec/{README.md,passing.txt}.
  Historical sections describe their own pins and scope, not current parity.
- Read pinned original testing/fixture/RNG/turn-execution documents, tables,
  complete callback definitions/callers and actual frozen scripts/scenarios.
  Backlog names and historical line ranges are not complete source boundaries.

## Delivered state and selected C pin

All canonical tasks **#194–226** are accepted and pushed sequentially. Priority
follow-C #288, #290 and #291 are also accepted. Current engine version is
**0.148.0-alpha**. Latest code commit is
**26553f5e5392bc0869e05c48fd55b60b67d304bb** (#226 STUDY). The following
documentation-only handoff commit does not change that version.

The eleven affected defaults now use immutable merged G3
**9a03d126b329ed466c4af48425553c025eaca26f**: rng, helper, check, loadsave,
olympia-cli, olympia-roundtrip, olympia-spec, turn, turn-driver, command and move.
Other tools retain independent pins. Do not globally override G3_REV,
blanket-repin/rebaseline, or redo #290/#291's accepted pin work.

Recent deliveries: #218 FLY, #219 SAIL including real USE601 interruption,
#220 FEE/BOARD/UNLOAD/FERRY, #221 construction predicates/table/start helpers,
#222 BUILD completion/dispatch, #223 REPAIR/IMPROVE, #224 RAZE/rebuild,
#225 BUY/SELL and #226 STUDY. Existing shared settlement, skill allocation and
grant owners are complete; reuse them, do not duplicate them. No family exit,
stage exit, initialization/golden acceptance or stage-5 work is claimed.
Exact-SHA/main CI subscriptions have reported no checks found, not remote green.

## Settled #289 scope — do not ask again

The owner approved retaining immediate mode in stage 5 while completing all
stage-4 follow-C work. Decision:
https://ampcode.com/threads/T-01a1896d-7b58-73b5-8e7d-cdd7f139aafc.

USE's missing-argument guard and the merged-pin proof are accepted #290/#291.
Remaining mandatory stage-4 obligations are:

- Complete ATTACK and its missing-argument guard in canonical **#235**.
- Complete weather summons and fog's naming guard in canonical **#264**.
- Strict unchanged original **missing-argument** C/C/Go fixture acceptance
  after those complete owners, before stage-4 exit. Saved CO/li text reparses
  on load; test stale/nonempty stored argument slots, not only empty slots.

Keep their canonical order. No guard-only ATTACK overlay, accelerated combat/
weather closure, successful callback stub or new scope deferral. Until complete
owners land, retain honest C/C proof and actual Go NYI evidence. Full merged-pin
load/save corpus acceptance does NOT mean this gameplay fixture passes.

Immediate -i, the five console commands/guards (be, additem, poof, know,
makeloc), world creation/golden flow and original **gm-immediate-missing** are
explicit stage-5 obligations. Do not implement them in stage 4 or treat that
deferred remainder alone as blocking canonical gameplay. Preserve their honest
unsupported results. Update #289 as complete owners land; do not use Closes #289
or tick it wholly complete. At stage-4 exit leave it OPEN, stating that only
the owner-approved immediate/console remainder remains. No other deferral is
authorized, and no stage-5 implementation is authorized by this prompt.

## Next bounded task: #227 N25 RESEARCH/FORGET, family #188

Read its live issue and complete current pinned use.c/glob.c owners before
briefing exactly one implementer. Read source, not skill.c: the owner is
**use.go/use.c**. Source-thread discovery (not implementation) identified:

- v_forget, use.c:1077–1149; research_notknown:1945–1977;
  v_research:1980–2062; d_research:2065–2118. Verify actual complete declaration/
  body boundaries; issue ranges are historical and can cut off definitions.
- Real glob rows: FORGET start only, wait0/poll0/priority1, NULL finish and
  interrupt; RESEARCH start/finish, wait7/poll0/priority3, NULL interrupt.
- Existing rpSkillEnt/pSkillEnt/forgetSkill/learnSkill/skillSchool are shared
  owners from P4; #226 adds complete STUDY. Reuse existing addNP, payment,
  maxEffAura, location/building ownership and streams.SkilResearch/Pick.
  The entire original skills RNG prefix/static helpers is already extracted
  by command-compat; no new RNG model or adapter is needed.
- Preserve learning-versus-know=dont NP refunds, ordered category cascades,
  duplicate rows, aura reduction and exact magician/weather flag effects.
  Do not introduce a missing-argument guard for FORGET that C does not have.
- RESEARCH checks missing/invalid/unknown skill, religion temple versus other
  tower and first-character ownership, then magician effective-aura>30 and
  province civilization>1, candidate availability, then ability to pay 25 gold.
  Completion does not repeat those location/owner/start checks: validate skill,
  charge 25 gold, allocate char, draw success, then select the still-available skill.
  Preserve money lost on failure, candidate evaporation diagnostics and
  discovery as a real unknown skill row, not a learned skill.
- research_notknown filters on absent target row and PRESENT prerequisite row,
  not learned prerequisite; preserve original research list order/duplicates,
  selection bounds and keyed per-day draws. Discriminate same-day stable versus
  next-day fresh research, changed candidates and failed-payment ordering.
- Original **skills-learn** and **skills-refused** require STUDY/RESEARCH/FORGET;
  read full README/check.sh/Lua/scenario.patch and require strict C/C/Go once
  #227 completes those prerequisites. They currently stop at RESEARCH start.
  **rng-repeat-leaves** combines research with USE 615 archery, whose complete
  owner is **#246**; retain its honest boundary now and accept original combined
  integration when that owner lands (family #192), not a silent deferral.

Then continue #228 onward in the existing index. Split newly proven oversized
prerequisites before coding; do not expand a bounded task into later gameplay.

## Accepted evidence to retain, not whole-stage parity

Latest parent independently passed bounded vet, fresh full-enabled tests
(11 packages; Olympia 2.438s), all eight sequential defaults, real-C oracle and
full unchanged spec. Durable commands/provenance are in existing tool READMEs.

- Current command oracle: **2024 sections /182841 output lines**.
  Input **856415 bytes**, SHA256
  `a79e7e2fbf55bca76ef11c8cd345387f6fda0751db7eb4dcc594beb9a562f5a8`;
  output **3513820 bytes**, SHA256
  `c2c17c88f50d7b82f2f4ecef68131d155de093176de7356b2bb126f558dc3db0`.
  Preserve BOTH ENTIRE accepted prefixes when appending, including high-bit
  bytes/trailing spaces. The inherited #225 prefixes contain 47 input/93 output
  high-bit bytes. Do not undo accepted #290's genuine assertion-location shifts.
- Spec: **129 fixtures /73 pass /56 fail /0 error**. Preserve every exact
  current name, not merely the count. passing.txt remains **1430 bytes**, SHA256
  `5e82f52fcab18c59f92abde108db03f4b57d0253213474e325880f531e38975f`.
  Stage 3's exact 36-name list can be recovered from
  `git show 9deba6e850c6079c9088fe2a14959e839c363eed:tools/olympia-spec/passing.txt`;
  its 876 bytes hash to
  `6d2080c29961395d8387364983b6543a6d05c7277b63bdfd16f9025a91736b67`.
  Golden remains NO, initialization exits 2 at unsupported -i. Only the existing tool
  writes passing.txt; never rewrite original expected bytes or run upstream --update.
- Latest skills-learn/refused replay: unchanged C1/C2 scripts YES/0, engines 0,
  all 241 cwd entries each and raw streams/status match. Fresh Go engine 2/script 1
  at RESEARCH start NYI; neither whole script is accepted yet.
- Accepted native script passes include recruitment/social successes, all three
  raft-sail branches, sail-out-first, move-ship/movement-refused, construction-
  build/refused, stop-running, raze-rebuild, trade-market/opium/refused and all 14
  retained seed-parser branches. STOP-running really passes since #223; do not
  reintroduce its obsolete sail-interrupt boundary. trade-routes still reaches
  USE v_find_buy NYI (#244); raft-loc-depth's load passes but attack needs #235;
  sail-stale-storm needs #265; quick-cast retains its actual later magic/item owners.
  Read current task-specific READMEs for complete scopes, not this compact list.
- #291 independently accepted **207 deterministic load/save roundtrips**, plus
  one separate exact EXPECTED_ABORT; strict supplementary C/C/Go raw/whole-cwd
  replay on all 207 inputs; actual full box/boxfile/check/loadsave corpora over 208
  cases; strict ordinary 233-entry turn and all lifecycle/dump/integration/
  flag-free aftermath controls. The separate release-swear-np abort is native 134,
  swear.c:180 FALSE at this exact pin; no successful +turn snapshot exists.
  Never manufacture one or generalize its narrow exception. Read the full
  merged-pin proof in tools/README.md; its commands/source/cache manifests
  supersede older 205-case evidence. Scoped callbacks since #291 did not rerun
  full 207 corpora and do not claim new full-tier acceptance.
- Disclosed platform exceptions remain red: scanf 78 glibc probes/11 macOS
  numeric differences, exit 1; island 114/119, five differences/seven failed checks,
  exit 1. Do not normalize/rebaseline them or call them passing.

## New-orb setup and safety

The new thread/orb does not inherit untracked tmp trees, caches, binaries,
captures, logs or agent state. Do not rely on those existing. Use the larger
a1.large (8 CPU/16 GB) environment for this continuation; inspect actual resources
and installed tools. Tracked setup installs Go from go.mod and CMake 4.1.3.
Re-establish the bounded routine baseline before #227; historical full-tier proof
is not a claim that new-orb full corpora have run. Run full tiers on an actual
source/pin dependency change and stage exit as tools/README.md requires.

Never build/write the shared additional C checkout or ~/Software reference.
Create isolated tmp/upstream.git if absent, fetch immutable pins there, and use
git show/pristine git archive exports. EVERY compatibility/affected invocation
must explicitly set `G3_REPO="$PWD/tmp/upstream.git"`. Rebuild absent corpora
from unchanged C fixtures, never Go saves. Verify source blobs/modes, pin,
binary/build provenance and input manifests before reusing any cache.

The source orb's canonical C was an isolated pristine archive under
tmp/_issue215/canonical, CMake Debug/Ninja, using
`-DCMAKE_C_FLAGS=-ffile-prefix-map=ABS_SOURCE/=` (trailing slash essential).
All 1237 original blobs/modes were verified before each reuse. Its executable
SHA256 is `2ea19140ebd00e1ca672bd022ff958c139efb3aeca59f5f91b3cb8ebb93b5212`.
If rebuilding in a different environment, establish and record that build's
actual provenance/hash; do not normalize assertion paths after execution.

Use physical disposable input copies, no symlinks/hardlinks. Never execute an
engine on protected cached input. Full serializers use the existing GAMES
override pointing to disposable cases/*/input/lib OUTSIDE WORK. Run shared-
scratch tools sequentially; do not replace the solved workflow with a harness
change or policy exemption. Localized callback work need not repeat idle full
corpora unless source dependencies require it; final stage acceptance still must.

Bound EVERY C/Go engine/test, alternate fixture launcher and expected abort:

```sh
. tools/lib/timeout.sh
run_limited 120 go "$(command -v go)" vet ./...
G3_REPO="$PWD/tmp/upstream.git" G3G0_FULL=1 RUN_TIMEOUT=120 \
  go test -count=1 -exec "bash $PWD/tools/lib/run-limited.sh" ./...
for tool in command move helper rng code loadsave turn-driver turn; do
  G3_REPO="$PWD/tmp/upstream.git" RUN_TIMEOUT=120 \
    GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
    bash "tools/$tool-compat/run.sh"
done
```

Original fixtures use unchanged check.sh with bounded NATIVE OLYMPIA_ENGINE_BIN
capture before script cleanup. Recreate the helper when absent: preserve argv0
and stdin, create a distinct numbered directory per invocation, open stdout/
stderr O_RDWR, invoke bash sourcing timeout.sh with
`run_limited 120 "$0" "$S4_ENGINE" "$@"`, seek/replay stdout then stderr,
close streams, write actual status/newline, and cp-a the ENTIRE cwd into the
absent capture/cwd before returning the engine status. Replaying only the last
invocation loses cases (seed-parser has 14; trade-routes has 5).

For each original script run unchanged C1, C2, THEN fresh Go. Compare invocation
inventories first, then raw status/stdout/stderr and ALL cwd paths/types/modes/
empty directories/bytes using existing tools/turn-compat and tree-modes.py,
no exclusions. Original scripts may redirect streams into cwd files; compare
those too. Reject dependency/NYI/stub/sanitizer/runtime-error markers BEFORE
expected statuses or recording replacement, including genuine 134 cases;
exercise actual native 0/134 default/update controls, not shell-signal guesses.
Expected bytes come only from unchanged C, never a model, normalized output,
fixture overlay, rewritten upstream harness or Go-produced reference.

## Per-task delivery and stopping rule

Re-audit all open issues before each task. Exactly ONE implementing subagent
per bounded coding task, sequentially, leaves changes UNCOMMITTED. Parent
reads complete pinned definitions/callers and full diff, independently repeats
bounded vet/fresh tests, relevant compatibility and unchanged real-C-twice/Go
acceptance, checks version (minor features, patch fixes/harness, retain alpha),
then commits directly to main with Closes only that accepted task. Push and
subscribe exact SHA/main CI; no checks found is not remote green. Update root/
stage/family/index checklists from fresh content. No branches, PRs, tags or
CI workflows. Upstream behavior is authoritative; no independent fixes. File
discovered bugs under WORKFLOW policy; stop affected work on a proven blocker.

Keep durable commands/provenance/limits/results in existing tool docs. Retain
decisive logs/manifests/failed-Go evidence, but remove temporary helpers,
binaries, matching captures, disposable WORK/RUN_DIRs and generated root or
internal/olympia coverage.tsv/use.tsv. A partial oracle is not a family exit.

Only after every in-scope contract and assigned original family fixture is
complete, run tools/README.md's entire final tier sequentially, preserving
all 73 current passes, exact 36 stage-3 names and flag-free/dump contracts. Include
strict original missing-argument; preserve the explicit gm-immediate-missing
stage-5 exception with #289 open. Report every lost name, unexplained difference
and open bug honestly. Independently accept stage 4, update #19/#183/families,
then STOP. Do not plan or implement stage 5 without a new owner request.
