Start stage 3 of the mdhender/g3g0 engine port: “Fixture dumps” in #19. Plan this stage, create its tasks, then implement them sequentially using the repository workflow. Stage 2 is accepted and closed in #132/#166. Do not repeat its implementation or start stage 4.

## Read first

- AGENTS.md, WORKFLOW.md, and tools/README.md, especially check tiers, source dependencies and accepted stage-2 evidence.
- #19 and #132/#166, including their final acceptance comments. Earlier red-gate and baseline-in-progress notes are historical, not the final state.
- tools/olympia-spec/README.md and passing.txt, tools/olympia-roundtrip/README.md, and tools/turn-compat/README.md.
- The chosen pinned C coverage.c/state_export.c, their declarations, all callers and instrumented execution paths, and the upstream fixture scripts and harness documentation that consume their outputs.
- Stage-2 completion thread if useful: https://ampcode.com/threads/T-01a182fa-25df-733d-85b0-f9ded807fbb0. Durable requirements and evidence are in the repository and issues; do not depend on that thread's temporary files or agents.

This file is a handoff for a future stage-3 run. Writing or committing it does not itself start stage implementation. When executing it, inspect the actual checkout, origin/main and current issue state rather than assuming nothing has changed.

## Accepted starting point

Last accepted engine/task commit on origin/main:
https://github.com/mdhender/g3g0/commit/6b381b7fdfebaacd593ce6ecebce1119d9f38b63

Version: **0.105.1-alpha**. This handoff is documentation-only and does not bump it. At stage-2 closure the worktree was clean, no implementing agents remained, and the only open local issue was tracker #19. No stage-3 issue or task plan had been created. Recheck ALL open issues, including unlabeled ones, before planning and before each task; do not duplicate issues created since this handoff.

Accepted stage-2 C pin: **e733acb63bacb21167a1f2589d8dcdc023364896**, including G3 #357. Some unaffected tools deliberately retain older pins. Do not repin them without a real source dependency.

Independent final acceptance on Linux amd64, GCC 12.2.0, Go 1.26.4:

- Vet and fresh bounded full-enabled tests pass; the four whole-game suites actually run through their full compatibility tools.
- All 15 engine/library default tools pass, as do genesis and mapgen at their own pins.
- Command oracle: unchanged real C twice and Go match **1234 sections / 82688 lines / 1635428 output bytes**. Input is **26815 lines / 513156 bytes**. Output SHA-256: `d08f0a35fbcd021dc44470da263bf483302cd070c31a8ea3bcfe54321b6b4336`; input SHA-256: `6a711d79ca5d5a65ca2184b430766037e89df591c6caf1c662bdd7498c174c27`.
- **201 deterministic roundtrips pass**, zero failures/errors, retaining all 196 inherited names. One separate **EXPECTED_ABORT** validates actual C twice exiting **134** at `swear.c (180): FALSE`, with matching raw captures/partial trees and no dependency/sanitizer markers. No successful `release-swear-np+turn` snapshot exists or may be manufactured.
- Full box/boxfile/check/loadsave pass **202 worlds**: 201 actual cached inputs plus each tool's separate base world. Box plain/pretty, all three boxfile variants and check cover **8,146,800 boxes**; loadsave covers **8,146,791 saved boxes**. Check/loadsave have zero check_db diagnostics.
- Strict full turn passes two identical C runs and exact C/Go status **0**, raw stdout/stderr and **all 233 cwd entries**, including reports, Times, auxiliary files, empty directories, all-entry types/modes and regular bytes. All **13 lifecycle controls** pass under that same artifact-inclusive contract. No exclusions or normalization.
- Spec: **21 pass / 102 fail / 0 error**, preserving the earlier 19 passes and adding `mint-top-id`/`token-vanish`. Preserve all **21** now recorded in passing.txt. Golden initialization still exits **2** at unsupported `-i`; golden is **NO**, not passing whole-spec parity.

The old stage-2 eight-report red gate is resolved. Do not restore artifact exclusions or describe the accepted strict gate as still red. Reports absent from the exit world's regular-player paths have separate asymmetric real-C coverage.

## Scope and boundaries

Implement the existing long-form options and their actual C contracts:

- `--coverage`: command-start outcome accounting and deterministic verb dump from coverage.c.
- `--use-coverage`: USE skill/item/invalid-target accounting and dump from coverage.c.
- `--phase-coverage`: real turn-phase accounting and dump from coverage.c.
- `--state-export`: the deterministic faction snapshot from state_export.c.

The Go CLI already parses these path options but rejects their modes before load_db. Inspect and reuse existing options, command tables, dispatch, USE metadata, turn phases, entity traversal, topology, output and file owners. Remove only the relevant NYI gates as complete implementations land; keep later modes explicitly unsupported.

Plan from complete pinned definitions and their actual dependencies, not file names alone. Preserve opt-in behavior, path parsing, NULL versus empty paths, allocation timing, integer accounting, row membership/order, escaping, exact whitespace/newlines, file-open failures, diagnostics/status and main's call order. Check table sentinel/index rules, unrecognized commands, zero-count rows, outcome boundaries, immediate-mode suppression, repeated calls and independent combinations of flags.

Instrument real execution where C does. Do not infer counters from saved state or reports, manufacture expected totals, or reuse test overlays that replace gameplay callbacks. Distinguish a C command genuinely marked unimplemented from a command whose C implementation exists but whose Go handler is still NYI. Do not turn the latter into successful dispatch or a fabricated C coverage outcome.

USE accounting does not authorize porting the remaining USE handlers. Reuse or port only the metadata/accounting prerequisites required by this stage, with complete real-C checks. Exercise reachable implemented paths and direct accounting contracts without successful substitutes for missing gameplay. If an acceptance fixture proves a genuine missing prerequisite, file and size it separately before implementation rather than silently expanding a dump task.

State export is not a new save format or a replacement report layer. Preserve exact JSON bytes and byte-string escaping, including NULL/empty names, controls, quotes/backslashes and non-ASCII bytes. Match real faction/unit/location visibility, skills/items, prisoner/ownership rules, ordering and exit selection from the source. Verify its state/RNG/allocation effects against actual C rather than assuming a comment proves it observational. Preserve copied versus original-list sorts and output/cache ownership; do not mutate persistent data merely to sort an export. Keep its real position after the save lifecycle.

Ordinary flag-free execution must remain identical to accepted stage 2 and must not create dump files or alter RNG draws, saved data, reports, diagnostics or output ordering. Compare combinations and failure paths, including dump destinations outside lib and paths colliding with existing entries, according to actual C behavior.

Later-stage commands, world creation, immediate mode, valid admissions, external charging, mail sending and separate HTML conversion remain out of scope unless a separately established prerequisite is required. Do not implement upstream enhancement G3 #114's proposed wider state export instead of the format at the chosen pin. No successful stubs for missing dependencies.

## Planning and pin gate

1. Audit ALL open local issues, including unlabeled ones. Work bugs/follow-C first; label and attach applicable follow-C issues to the stage. Recheck current upstream defects and source before treating an old finding as resolved or blocking.
2. Inspect current upstream main and choose an explicit immutable stage-3 revision. Record why it is selected. Never silently follow a moving branch. If moving from the accepted stage-2 pin, run `tools/affected.sh OLD NEW` with isolated G3_REPO, inspect every relevant source change and fixture addition/removal, and handle follow-C work and affected recordings explicitly. Do not blanket-rebaseline unrelated tools.
3. Before new dump implementation, independently re-establish the inherited baseline at that pin: strict whole-world turn and all 13 artifact-inclusive lifecycle controls, all inherited deterministic roundtrips and the separately accounted abort control, full box/boxfile/check/loadsave corpus, default tools, vet/fresh tests and spec with no lost accepted passes. Explain corpus or fixture-count changes; do not mistake a pin-move regression for a dump gap. File/split prerequisite work before proceeding.
4. Create a stage issue as a sub-issue of #19, assigned to @me, with pinned C ownership/call sites, dependencies, task order, byte-exact checks and exit criteria. Create and attach every bounded task issue before implementation. Do not assume one task per source file; coverage families, real hook integration and export traversal may need separate closures. A final acceptance task must not hide first-time implementation of a large dependency.
5. Inventory all `phase-*` fixtures at the chosen pin and map each to the actual phase hooks and dump/export contracts it needs. At the accepted stage-2 pin these are `phase-day0`, `phase-daily`, `phase-month-open`, `phase-month-gold`, `phase-month-decay` and `phase-month-close`. Re-enumerate rather than assuming this list never changes. Record direct C-driven coverage for contracts those fixtures do not reach.

## Required per-task process

1. Inspect all open local issues again; bugs/follow-C take priority. Split oversized closures or new prerequisites before coding. File every discovered local defect here and every upstream defect in playbymail/olympia-g3, assigned to @me. Stop and ask the owner on a blocking upstream defect; preserve nonblocking C bugs with source-linked comments rather than independent fixes.
2. Launch exactly one implementing subagent for the task. It reads its issue, stage notes and guidance, implements in the existing owners, and leaves changes **UNCOMMITTED**. Parent owns independent source/diff review, acceptance and delivery. Do not hand the whole stage to one subagent or create replacement threads for ordinary phases.
3. Independently read complete pinned C and the diff. Require exact output/state checks against unchanged real C, bounded vet, fresh bounded tests, relevant compatibility tools and a version bump. New features below v1 use minor bumps; fixes/harness work use patch bumps. Expected output never comes from Go, rewritten C models or a successful missing-dependency stub.
4. Commit directly to main with `Closes #TASK` for that task only, push, and tick both #19 and the stage tracker after acceptance. Write upstream references as “G3 #N.” No branches, PRs, tags or GitHub Actions workflows. Subscribe the current thread to CI for the exact pushed head after each push. Prior pushes had no configured CI checks; local acceptance, not a fictitious remote green check, is the evidence.

## C source safety and practical validation

NEVER touch or build in the shared additional olympia-g3 checkout or ~/Software/playbymail/olympia-g3. Read with git show at the immutable pin or git archive into tmp. Use isolated bare `tmp/upstream.git`; create it if absent and fetch there, not into the shared checkout. Explicitly set `G3_REPO="$PWD/tmp/upstream.git"` for EVERY compatibility/affected tool invocation.

The completion orb had a freshly built exact-stage-2-pin cache at `tmp/stage2-roundtrip` and strict metadata at `tmp/parent-stage2-turn`. A future thread may have neither. Validate revision, archive/build provenance and input integrity before reuse. Never run an engine against cached input or use a symlink as an input-copy source. Rebuild missing inputs from pristine pinned fixtures, not a previous engine's output. Full corpus tools use `GAMES="$PWD/tmp/stage2-roundtrip/cases"` when that validated cache is selected; inspect current scripts for actual cache variables.

Wrap EVERY C/Go engine in tools/lib/timeout.sh's run_limited. Test executables use `go test -count=1 -exec "bash $PWD/tools/lib/run-limited.sh" ./...`. Run shared-scratch/cache tools sequentially. Keep default checks small; whole-world work stays behind `--full` / `G3G0_FULL=1`. Delete matching captures promptly and retain decisive failure/provenance logs. Use tools/lib/tree-modes.py for portable, failure-aware all-entry mode/type comparisons.

Extend existing oracles/fixtures where they own the behavior. Require two deterministic real-C runs before Go comparison or recording refresh; verify complete extracted definitions/dependencies unchanged. Preserve inherited recording input/output byte prefixes when appending cases, including intentional trailing spaces. Do not normalize JSON/TSV, whitespace, routing, diagnostics or filesystem comparisons to hide differences.

Reject missing-C-dependency/NYI/stub/sanitizer markers BEFORE accepting expected statuses or replacing recordings, including expected failures. Exercise negative controls proving the harness rejects invalid references. Preserve the exact-pin expected-abort contract; ordinary failed turns remain errors, and pin changes require source-backed review rather than carrying an arbitrary allowed failure forward.

Keep stage-2 invariants: real dispatch/phase order, allocation timing, input reads/assertions, retained queues, in-place versus copied sorts, lazy backlinks, exactly-once effects, bounded cached FILE/inode ownership, and RNG stream tags/keys/draw order. Extend affected.sh/provenance only for real new dependencies.

## Known limitations to carry forward

At stage-2 closure the open upstream defects were G3 #309 (invalid RNG ranges/overflow, won't fix), G3 #358 (positive-quantity zero-price BUY division by zero), and G3 #359 (NULL relational pointer comparison during long-name truncation). Recheck their current status and the chosen source pin. G3 #352–#356 were closed and followed in #131; G3 #114 was an enhancement, not an accepted export expansion.

The owner designated G3 #358 nonblocking for stage 2. Its malformed case is **C136/Go2, NOT passing parity**; no guard, recovery or repin was introduced to conceal it. Reassess if stage-3 acceptance reaches it or promises broader malformed-input coverage. G3 #359's observed deterministic C truncation was preserved; reassess for compiler/sanitizer portability. Neither disclosure authorizes silently changing behavior.

Known scanf Linux/macOS differences remain red: 78 deterministic glibc probes, 11 numeric-width differences, exit 1. Island remains 114/119 matching cases, five differences/seven failed checks, exit 1; all 17 upstream gates pass for both engines. Do not describe these tools as passing or rewrite their expectations. Stage-2 validation was Linux-only; macOS and the full C spec were not rerun.

## Stage exit and stopping point

Require every `phase-*` fixture at the chosen pin to pass through the unchanged upstream check scripts against Go, exercising actual turn phases rather than modeled substitutes. Also require byte-exact deterministic C/C/Go checks for all four dump modes and their integration, including nonempty/zero output, filters/order/escaping, combined flags, failure paths, and state/RNG/file-lifecycle aftermath. Parsing JSON/TSV or passing phase fixtures alone is not byte parity for the complete contracts.

Run the final tier sequentially in tools/README.md order: bounded vet, fresh bounded `G3G0_FULL=1` tests, every required pinned default, all inherited and meaningful new deterministic roundtrips plus the separate expected-abort control, full box/boxfile/check/loadsave corpus, strict whole-working-directory turn and all 13 lifecycle controls, and the spec with **no lost passes against the 21-pass stage-2 baseline**. Record which additional fixtures pass; do not claim unrelated commands or golden initialization are implemented. Keep flag-free stage-2 acceptance green as well as the new opt-in paths.

Only after stage exit and final acceptance pass, close the stage issue, update #19 and committed acceptance evidence, summarize every open bug and any prerequisite for the next stage, then **PAUSE**. Do not plan or implement stage 4 until the owner asks.
