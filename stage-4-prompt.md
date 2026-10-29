Start stage 4 of the mdhender/g3g0 engine port: “Commands, by family” in #19. Plan the stage, create its issues, then implement its tasks sequentially using the repository workflow. Stage 3 is accepted and closed in #169/#178. Do not repeat its implementation or start stage 5.

This file is a handoff for a future stage-4 run. Writing, committing or pushing it does not start stage 4, create its issues or authorize implementation in the drafting thread. Execute it only when the owner asks to start the stage. Inspect the actual checkout, origin/main and current issues then; the snapshots below are not assertions about a future checkout.

## Read first

- AGENTS.md, WORKFLOW.md and tools/README.md, especially the current **Final-tier validation (#178)** section, check tiers, source dependencies and stage-exit commands. Historical task notes retain earlier temporary boundaries; do not mistake them for the accepted final status.
- #19 and the final acceptance comments in #169/#178. Stage-3 completion thread if useful: https://ampcode.com/threads/T-01a184c1-d416-7744-b33b-39cf0c9b7689. Durable requirements and evidence are in the repository and issues; do not depend on that thread's temporary files or agents.
- tools/olympia-spec/README.md and passing.txt, tools/olympia-roundtrip/README.md, tools/command-compat/README.md and tools/turn-compat/README.md.
- At the chosen immutable C pin: doc/claude/testing.md, doc/claude/ultron-fixtures.md, doc/claude/rng-streams.md, doc/turn-execution-order.md, tests/olympia/fixtures/README.md, tests/olympia/engine.sh, and the actual regress check scripts and scenarios.
- The complete C command and USE tables, item dispatch and callback definitions; tests/ultron/backlog.tsv, use-backlog.tsv and phase-backlog.tsv. Treat backlogs as coverage maps, not substitutes for source or evidence that every path is exercised.
- Existing Go owners, especially internal/olympia/glob.go, input.go, use.go, the scheduler, command-family files and current real-C oracles. Read callers and shared dependencies before deciding what remains to port.

## Accepted starting point

Last accepted stage/task commit on origin/main:
https://github.com/mdhender/g3g0/commit/8255aee778c0736ab8c4d442a7bae098b902b7f8

Version: **0.115.1-alpha**. This documentation-only handoff does not bump it. At handoff drafting, main and origin/main agree at that acceptance commit, the worktree is clean, no implementing agents remain, and the only open local issue is tracker #19. No stage-4 issue or task plan exists. Recheck ALL open issues, including unlabeled ones, before planning and before each task; do not duplicate later work.

Accepted C pin: **f73e48118805ffe8ff5e763bfa4f20ae00c3a705**, including merged G3 #360 followed in #179. Unaffected tools deliberately retain independent older pins; do not repin them without a real source dependency.

Final acceptance was Linux amd64, GCC 12.2.0, Go 1.26.4:

- Bounded vet and fresh bounded full-enabled tests pass across **11 test packages**. All **15 engine/library defaults**, genesis and mapgen pass at their own pins. Whole-game suites actually execute through their full compatibility tools, not merely G3G0_FULL=1 without game lists.
- Command oracle: unchanged real C twice and Go match **1380 sections / 107214 output lines / 2134398 output bytes**. Input is **578307 bytes**, SHA-256 `3481dcd40549f16ff6fc67201a70bab1e00c08a599799151e7706734ff594659`; output SHA-256 `c4cc3c8ee0cd5f0b40ac6eb239aa7bf51235aaf1448fa6d9aad81f5a821f1bac`. Preserve both entire inherited byte prefixes when appending cases, including intentional C trailing spaces.
- **201 deterministic roundtrips**, zero failures/errors, plus **one separate EXPECTED_ABORT**. Supplementary fresh C/C/Go replay compares every pristine input and cached C, with complete raw status/streams and cwd types/modes/empty directories/bytes. The separate C/C abort is **134**, exact `swear.c (180): FALSE`; no successful `release-swear-np+turn` snapshot exists or may be manufactured.
- Full **box/boxfile/check/loadsave each pass 202 worlds**: 201 actual cached inputs and the tool's own base world. Box plain/pretty, all three boxfile variants and check cover **8,146,800 boxes**; loadsave covers **8,146,791**. Check/loadsave emit zero check_db messages.
- Strict ordinary turn matches real C twice and Go on raw status **0**, stdout/stderr and **all 233 cwd entries**, reports/Times/auxiliary files included. All **13 lifecycle / 10 command / 15 USE / 23 phase / 13 export** controls pass. All **84** four-dump integration controls pass (**74 status 0 / 10 status 1**), plus **20** flag-free aftermath comparisons.
- The four complete dump modes are implemented and byte-exact: --coverage, --use-coverage, --phase-coverage and --state-export, including combinations, failures, collisions and state/RNG/file aftermath. Preserve main's command → USE → phase writers → check_db → save_db/save_logdir → state export → Times order. Existing tests and READMEs record the complete dump hashes.
- All six unchanged phase scripts and savage-cap pass real C twice then Go, including complete pre-cleanup raw/cwd captures: **252 / 245 / 242 / 241 / 246 / 242 / 238 entries** for day0, daily, month-open, month-gold, month-decay, month-close and savage-cap.
- Spec is **36 pass / 89 fail / 0 error** across **125 fixtures**, preserving all 21 stage-2 passes. The **876-byte** passing.txt has SHA-256 `6d2080c29961395d8387364983b6543a6d05c7277b63bdfd16f9025a91736b67`. Preserve every one of the 36 names, not just the original 21 or the numeric count. Golden is **NO**, initialization exits **2** at unsupported -i; this is not whole-spec parity.

## Scope and family ownership

Port the remaining actual gameplay command families in #19, not a second dispatcher or a fixture-specific simulation. Inventory what already works first: stages 1–3 implemented substantial movement, stacking/death/transfer helpers, NPC behavior, reporting, scheduling and coverage, plus GUARD/STACK, hinder-meditation and the drum callback. An existing helper or a passing refusal path does not prove the whole command is complete.

Use #19's family list as the starting ownership map:

1. Players and GM: player-*, gm-*, governance-*.
2. Attitude, loyalty, stacking and inventory: attitude-*, loyalty-*, stacking-*, inventory-*.
3. Movement: move-*, movement-refused, sail-*, fly-*, coast-*; include raft-loc-depth and raft-sail rather than losing names not listed in #19's shorthand.
4. Production, construction, recruitment and trade: production-*, construction-*, recruitment-*, trade-*, raze-*, mint-*.
5. Exploration, stealth and skills: exploration-*, stealth-*, skills-*.
6. Social: social-* and their real shared ownership/loyalty prerequisites.
7. Combat and pillage: combat-*, guard-pillage, and the related failure-path fixtures assigned by actual ownership, not just their names.
8. Magic and USE: magic-* and every actual use-* fixture; subdivide by school/item/shared closure, not one giant USE implementation task. Re-enumerate rather than assuming #19's historical “32” remains the count.
9. Assign order-*, wait-*, token-*, rng-repeat*, release-swear-np, savage-cap, inherited loader/save/CLI/phase checks and every other actual regress fixture exactly once to its owning family or explicit later-stage boundary. Already passing fixtures are retained regression obligations, not reasons to reimplement their code.

For every command-table verb, USE skill and item key, map start/completion/interruption callbacks, shared owners, implementation status, fixtures and missing outcome coverage. Port complete in-scope callback contracts, including paths the frozen fixtures miss; fixture success is the family exit gate, not permission to leave other reachable callbacks as NYI. Preserve genuinely NULL/unimplemented C entries as C defines them; do not invent functionality.

World initialization/immediate mode and the golden flow remain stage 5. Separate CLI modes (-a, -e/-E, -s, -M, -A, -h, -T, -w, -t and PWD handling) remain stage 6. Lua bindings/scenario authoring remain stage 7. Normal in-game command behavior required by a family is in scope; an unrelated bulk CLI mode or world-creation pipeline is not. Establish and file any genuine shared prerequisite separately before implementation; do not silently pull entire later stages into a family.

At the accepted pin **mail-libdir-space** explicitly invokes -r -S -M -h, -M, and -E -M with external-tool stubs; it is a stage-6 mode fixture, not an unfinished command family. Recheck its script at the selected pin and record that boundary in the initial inventory. Do not demand 125/125 spec success by implementing deferred modes, and do not label the remaining suite green. Any other proposed deferral needs the same source-backed later-stage dependency; large gameplay work or a failing fixture is not by itself a deferral. Ask the owner if actual requirements conflict with #19's stage boundaries.

## Planning and pin gate

1. Audit ALL open local issues, including unlabeled ones; bugs/follow-C take priority. Recheck current upstream issues and main before relying on the historical snapshot. Label/attach applicable follow-C issues and handle them before new family work.
2. Select and record an explicit immutable stage-4 pin. If moving from f73e481, run tools/affected.sh OLD NEW with isolated G3_REPO, inspect complete relevant source/fixture changes, handle follow-C work and affected recordings explicitly, and preserve independent unaffected pins. Never silently follow a moving branch or blanket-rebaseline tools.
3. Independently re-establish the inherited baseline at that pin before new gameplay: bounded vet/fresh tests, required defaults, all roundtrips plus the separate abort, four full corpora, strict ordinary turn and all coverage/export/integration controls, all phase scripts and spec preserving the accepted 36 names. Explain fixture/corpus/count changes; do not mistake a pin regression for a new command gap. A pin change affecting the abort requires source-backed review of its new contract.
4. Enumerate every pinned regress check.sh and every command/USE table entry; publish a complete fixture/callback ownership and dependency matrix in the stage issue. Record current passes, genuine missing implementations, shared prerequisites, exact C owners and byte-exact acceptance per family. Reconcile #19's shorthand with the actual inventory. Use dependency order, not a promise to follow the list above blindly.
5. Create a stage issue under #19 and one child issue per family, assigned to @me. Split oversized families into bounded implementation/prerequisite child tasks before coding, with explicit dependency order. Create and attach all initial tasks before implementation. No family-sized task may quietly require a context-filling combat or magic dependency graph; file additional discovered closures before expanding scope.
6. Give each task complete C definition/caller/dependency references, existing Go owners, concrete asymmetric success/refusal/failure cases, multi-day/interruption cases where applicable, and a path to the family's original fixture checks. Extend the existing oracles. Where fixture coverage is missing, request it upstream; do not edit their harness or invent expected behavior here.

## Required per-task process

1. Re-audit all open local issues. Work bugs/follow-C first. File local defects here and upstream defects in playbymail/olympia-g3, assigned to @me. Follow C, including known nonblocking bugs with source-linked comments. If an upstream defect blocks acceptance, stop affected work and ask the owner; do not mask it with guards, recovery, allowed failures or a convenient repin.
2. Launch exactly **one implementing subagent per bounded task**, sequentially. It reads the task, stage plan, guidance and complete pinned definitions, implements in the existing owners and leaves changes **UNCOMMITTED**. Do not hand the whole stage or an oversized family to one subagent. Parent owns source/diff review, independent verification, delivery and tracker updates.
3. Require unchanged real C twice before comparing Go or refreshing C-only recordings. Compare complete raw output/status, relevant persistent state, lists/allocation/cache effects and RNG aftermath; ordinary fixture checks alone may not expose those. Expected bytes never come from Go, a rewritten C model, fixture-specific success overlays or replacement gameplay callbacks.
4. Parent independently reviews complete pinned C and the diff, runs bounded vet, fresh bounded tests, relevant compatibility defaults and exact task/family checks, and verifies the version bump. New features below v1 use minor bumps; fixes/harness changes use patch bumps; alpha remains. Documentation-only commits need no bump.
5. Commit directly to main with Closes #TASK for that accepted task only; push and tick #19, the stage and family checklists after acceptance. Close a family only when all its assigned fixtures pass and its in-scope implementation is complete. Use “G3 #N” for upstream references. No branches, PRs, tags or GitHub Actions workflows. Subscribe this thread to CI for the exact pushed SHA/branch after every push; no checks found is not a remote green result.

## Implementation invariants

- Reuse existing dispatch, command parsing/reparsing, priority queues, scheduler, USE metadata, inventory/ownership/death, topology, output and RNG owners. Change the source of truth instead of layering parallel helpers or duplicate side effects. Separate behavior-preserving refactors from gameplay changes and verify each phase.
- Preserve C widths/sign/wrap, NULL versus empty, validation/assertion/message order, saved command arguments, costs/charges and when they occur relative to revalidation, running/waiting/completion/interrupt semantics, zero-time inhibition, poll/timer behavior and retained order queues. A start that succeeds is not proof its completion/interrupt path is correct.
- Preserve original-list versus copied sorts, allocation timing, lazy backlinks, cached exit/FILE ownership, exactly-once transfer/death/reclamation, report routing and all RNG stream tags/keys/draw order. Exercise real dispatch and cross-family interactions, not only direct helper calls.
- Keep genuinely missing Go callbacks explicitly NYI until implemented. Never turn a non-NULL C start/completion/interrupt into nil or return success to satisfy coverage. Keep existing accounting hooks at their original points, including rejected/failed attempts and saved pre-reparse command indices.
- Preserve the completed report, coverage and state-export contracts while gameplay changes their legitimate contents. New implemented commands may add counters/state; explain changes from actual C, not from adjusted thresholds or normalized JSON/TSV. Accepted flag-free turns must not gain dump files or extra effects.
- Test likely mistakes with discriminating inputs: signed/boundary values, asymmetric ownership/stacks, hidden/contact/admission differences, insufficient resources, stale targets after delays, interrupted actions and keyed versus legacy RNG. Keep permanent tests for public behavior, not mock call order or private implementation layout.

## C source, cache and validation safety

NEVER touch or build in the shared additional olympia-g3 checkout or ~/Software/playbymail/olympia-g3. Read git show at the pin or git archive into tmp. Use isolated bare tmp/upstream.git; create/fetch there if missing. Explicitly set `G3_REPO="$PWD/tmp/upstream.git"` for EVERY compatibility/affected invocation. Keep per-tool pins independent; do not globally override G3_REV.

The completion orb had a verified exact-pin cache at tmp/issue179-roundtrip, preserved older inputs at tmp/stage3-roundtrip, and evidence under tmp/issue178-* and tmp/parent178-*. A new thread may have none of these. Verify revision, original archive/modes, executable provenance and input names/types/modes/bytes before reuse; rebuild absent inputs from pristine pinned fixtures, not prior Go output. Historical cache: **1212** original files; C executable SHA-256 `151c7355890eba7f8bfda1f4f16f4a490947aaa1db1ffcc3ca859cd659c01158`; **201 inputs / 41938 files / 5228 directories**.

Never run engines on protected cached inputs or use symlinks/hardlinks as input copies. Make verified physical disposable copies, outside any WORK the tool empties. The full-box read-only serializers also use the supplied GAMES directory directly: use the **existing GAMES override** to disposable `cases/*/input/lib` copies. Stage 3 resolved this without a harness change or policy exemption; do not recreate that false blocker. Use the existing copy workflow in tools/README.md for all four full corpus tools.

Wrap EVERY C/Go engine/test executable with tools/lib/timeout.sh's run_limited or the existing run-limited.sh test adapter, including alternate-engine fixture launchers and expected aborts. Fresh test command:

```sh
. tools/lib/timeout.sh
run_limited 120 go "$(command -v go)" vet ./...
G3G0_FULL=1 go test -count=1 -exec "bash $PWD/tools/lib/run-limited.sh" ./...
```

Check current tool supervision rather than assuming it exists everywhere; the lists default historically needed temporary bounded compiler/native/test launchers, documented in tools/README.md. Run shared-scratch/cache tools sequentially, keep defaults small and full-world work behind --full / G3G0_FULL=1. Delete matching captures and temporary helpers promptly; retain decisive failure/provenance evidence.

Use unchanged upstream check.sh scripts with OLYMPIA_ENGINE_BIN, preserving their own expected-status, coverage, manifest and JSON checks. Capture real C twice then Go before original script cleanup when proving strict fixture equality, using turn-compat's documented bounded native capture method. Require raw status/streams and complete cwd entries/types/modes/empty directories/bytes, with tree-modes.py and the existing strict comparator; no exclusions or output normalization. Scenario authoring/build-scenario.sh/check-lua.sh remain C-only, not prerequisites for replaying frozen scenarios.

Reject dependency/NYI/stub/sanitizer markers BEFORE accepting expected statuses or replacing recordings, including status 134 or intentionally failing cases. Exercise actual negative rejection controls. Append to inherited recording bytes, never rewrite prefixes to match Go. Extend affected.sh/provenance only for actual new source dependencies. Let the existing spec tool record genuine newly passing names; do not edit passing.txt by hand or remove inherited names to manufacture no-loss success.

## Known limitations and stopping point

At handoff, the only open upstream issues are **G3 #309**, invalid RNG ranges/overflow (won't fix), and **G3 #114**, a wider state-export enhancement. No accepted stage-3 input reached an invalid RNG range; stage-4 commands may expose new cases, so reassess against the actual failing input rather than assuming nonblocking forever. G3 #358/#359 are closed, fixed by merged #360 and followed in accepted #179; their old stage-2 red findings are historical, not current-pin defects. Recheck all statuses/source at execution time.

Known platform exceptions remain red: scanf **78 deterministic glibc probes / 11 macOS numeric-width differences**, exit 1; island **114/119**, five differences/seven failed checks, exit 1, all **17** upstream gates passing for both engines. Do not rebaseline, suppress or call them passing. Prior acceptance was Linux-only, not native macOS/Windows or sanitizer validation. Unsupported initialization/golden and genuinely deferred modes remain explicit, not successful stubs.

Stage-4 exit requires every in-scope command family's assigned original fixtures to pass, complete in-scope callback contracts, and a complete ownership inventory accounting for all other regress fixtures without unexplained gaps or lost accepted names. Record exact new passes and any source-backed later-stage mode deferrals. An expected-failure fixture passing its original check is valid; a plain failed/aborted turn is not passing parity. Do not confuse the separate release-swear-np abort control with a successful roundtrip or successful turn snapshot.

Run the final tier sequentially in tools/README.md order: bounded vet; fresh bounded full-enabled tests; required pinned defaults and peripheral checks with known red exceptions disclosed; all inherited and meaningful new deterministic roundtrips plus the exact separate abort; full box/boxfile/check/loadsave corpora on disposable inputs; strict whole-cwd ordinary turn and all lifecycle/command/USE/phase/export/integration/aftermath controls; all six phase scripts, relevant original family scripts and the full spec with **no lost names against the accepted 36-pass baseline**. Require unchanged real C/C/Go full-byte evidence for new family paths and retain flag-free stage-3 acceptance. Do not claim whole-spec/golden success while deferred-mode failures remain.

Only after stage exit and independent final acceptance pass, commit durable commands/counts/hashes/provenance/limitations in existing tool docs, close the stage issue and family/task issues, update #19, summarize every open bug and any prerequisite for stage 5, then **PAUSE**. Do not plan or implement stage 5 until the owner asks.
