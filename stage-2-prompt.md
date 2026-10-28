Start stage 2 of the engine port: “Output and reports” in mdhender/g3g0 #19. Stage 1 is complete (#55 and #90, closed). You have my permission to plan the stage and then start its work without asking again.

Read first:
- #19 for stage 2’s scope and exit check.
- #55 and #90, including their closing comments and acceptance evidence.
- AGENTS.md, WORKFLOW.md (“Port stages”), and tools/README.md (compat tools and check tiers).
- Your memory notes, especially port-stage-process, if available.
- The stage-1 handoff thread: https://ampcode.com/threads/T-01a17dcc-e13a-7252-9bd9-befaa68a78bb

Baseline:
- Stage 1 ended at bd654a7f5557e3ebaa6d203f3affa23f1421a639, version 0.76.2-alpha, pushed to main.
- Its C pin was 6a25f2d252ce45bd5adbe3baaccb3be3eed56279.
- The real no-orders world turn matches C on all 190 state files, raw stdout/stderr and exit status. All 13 lifecycle controls pass.
- Deterministic roundtrips pass all 196 cases. Full box/boxfile/check/loadsave checks pass across 197 worlds.
- Spec progress is 19 pass / 101 fail / 0 error, with no lost passes. Golden initialization still stops at unsupported -i; that is not this stage’s exit.
- Report artifact equality was deliberately deferred. The existing turn comparator excludes those artifacts.

Stage 2 scope:
- sout.c output routing, out/wout, indentation, styles, wrapping, logs and per-player report assembly.
- Names, codes and display helpers.
- Turn/location/character/player reports, summaries, lore sheets, order templates and GM reports.
- Saved reports, Times files, and player-list/email/forward/faction files.

Do not assume those areas are all unported. Stage 1 already implemented substantial display/output functionality, cached log ownership, Times and save_logdir lifecycle, lore input reads, and report-path persistent-state effects. Inspect and reuse what exists.

Plan:
1. Pin the stage to current olympia-g3 main. Run tools/affected.sh from the stage-1 pin to the proposed pin and handle affected recordings and follow-C work explicitly. Do not silently follow a moving branch.
2. Create a stage issue as a sub-issue of #19. Plan from actual pinned C ownership and call order, including output prerequisites and main’s report tail.
3. Create every task sub-issue before implementation begins. Size each so one sub-agent can finish without exhausting context; split oversized work before coding.
4. Give each task a byte-for-byte check against complete, unchanged real C definitions and dependencies, or a concrete step toward the stage exit. Expected output must never come from Go or a rewritten model of C.
5. Extend the existing turn harness to compare reports as well as database state. Stage-1 artifact exclusions must not hide stage-2 mismatches. Preserve the existing fast/default checks.

Exit:
On identical copies of pinned world.tgz, run:
  olympia-g3 -r -l ./lib -S test-use-const-report-date
Require identical whole lib/ trees, including log/save reports, Times and auxiliary output files, plus identical raw stdout, stderr and exit status. Check paths, file types and modes where applicable. Repeat C to establish determinism. Do not normalize away whitespace, routing, ordering or real diagnostics.

For each task, in order:
1. Inspect ALL open local issues, including unlabeled ones. Work bugs and follow-C issues first. Label applicable upstream-filed follow-C issues and attach them to the stage.
2. Launch one implementing sub-agent. It reads its issue, stage notes, AGENTS.md and WORKFLOW.md and leaves changes uncommitted.
3. Independently review the diff and pinned C yourself. Verify no independent fixes, the task’s exact C comparison, go vet ./..., fresh go test ./..., and a version bump.
4. Commit with “Closes #N” for that task only. Write “G3 #N” for upstream references. Push to main and tick the task in #19 and the stage issue.
5. File every discovered bug, ours here and upstream defects in playbymail/olympia-g3, assigned to @me. Fix our bugs as stage tasks. Reproduce upstream bugs with comments linking their issues.
6. If an upstream bug blocks progress, stop and ask me.

Important implementation boundaries:
- Replace stage-1 report-state projections with real reporting through the existing owners; do not layer a second reporting system over them.
- Preserve their allocation timing, saved-list ordering, knowledge changes, input reads and assertion order. Do not perform side effects twice when adding rendering.
- Preserve the accepted log cache and file lifecycle rather than reopening files per line.
- No successful stubs for unported rendering dependencies.
- Keep later-stage commands, coverage/state-export modes, world creation, valid new-player admissions, external account charging, mailing and HTML out of scope unless a genuine stage-2 prerequisite is established.

C source safety:
Never touch the working tree or build under ~/Software/playbymail/olympia-g3; someone else uses it. The additional orb clone may refer to that same checkout. Read C using git show at the pin or git archive into tmp/. Stage 1 used bare tmp/upstream.git. Reuse exact-pin caches only if present and validated; do not assume a new thread inherits them.

Keep checks practical:
- Defaults use recorded small cases. Whole-game checks stay behind --full / G3G0_FULL=1.
- Wrap every C/Go engine run in tools/lib/timeout.sh’s run_limited.
- Run tools sequentially when they recreate scratch directories or consume the same game cache.
- Write little to disk; promptly remove matching per-game captures.
- Preserve existing recordings as exact prefixes when extending cases.
- Reject missing-C-dependency markers before accepting expected failures or updating recordings.
- Use tools/lib/tree-modes.py for portable mode/type comparisons; #130 fixed BSD stat usage that falsely compared filesystem statistics on Linux.

Carry forward:
- macOS abort() flushes stdio and exits 134; SIGSEGV flushes nothing and exits 139. Go assertions use assertion{file,line,cond}.
- RNG stream tags, keys and draw order must match C and doc/claude/rng-streams.md exactly.
- At stage-1 closure there were no local bugs/follow-C issues open. Known upstream bugs were G3 #309 (wontfix) and #352–#356; recheck current status.
- Known scanf Linux/macOS and peripheral island differences are documented in tools/README.md. Do not claim those checks pass or silently change their expectations.

When the stage exit and final check tier pass, close the stage issue, summarize every open bug and any that must be fixed before stage 3, update #19 and the acceptance evidence, then pause. Do not plan stage 3 until I say so.
