Start stage 1 of the engine port ("A turn with no orders" in #19). Stage 0 is done (#20, closed). You have my permission to plan the stage and then start its work without asking again.

Read first:
- #19 for the stage's scope and exit check, and #20 with its closing comment as an example of a finished stage.
- AGENTS.md, WORKFLOW.md ("Port stages") and tools/README.md (the compat tools and check tiers).
- Your memory notes, especially port-stage-process.

Plan:
1. Pin the stage to the current olympia-g3 main. Stage 0 ended at 47535543d83fbc5649d54b63b49e578f2b11ab1c. Run tools/affected.sh to see which compat tools a pin move affects.
2. Create the stage issue as a sub-issue of #19. Give it a detailed plan that follows the order the C turn driver runs in (input.c: the clock change, command loading, the priority queues, the STATE_LOAD → RUN → DONE machine, day 0, the daily loop, daily_events, post_month), then the world's own work each turn on the bare map.
3. Create one task sub-issue per piece of work, all before work starts. Size each so one sub-agent can finish it without running out of context, and split anything bigger. Each task names its byte-for-byte check against C: a new or extended compat tool driven by the real C files, or a step toward the exit check. The exit check is that `-r -S` on world.tgz writes the same database files as C, with the same stdout, stderr and exit status. Reports are stage 2: the stage exit compares only the database files.
4. Plan the branches stage 0 left as "not yet implemented" that a turn reaches. The ones I know of: set_lord between two real players, load_command falling through to do_command, and add_item's lore path.

Then, for each task, in order:
1. Work any open follow-c and bug issues first. Label any follow-c issue G3 files without the label, and make it a sub-issue of the stage. Related follow-c issues can share one run.
2. Launch one sub-agent per task. It reads the issue, the stage notes, AGENTS.md and WORKFLOW.md, and leaves its changes uncommitted. It must not touch the working tree of ~/Software/playbymail/olympia-g3 (someone else uses it). It reads C at the pinned revision with `git show` or `git archive` into tmp/.
3. Verify the work yourself:
   - `go vet ./...` and `go test ./...` pass.
   - The task's check against C matches byte for byte.
   - The diff makes no independent fixes. Upstream bugs are reproduced, each with a comment linking the upstream issue.
   - The version is bumped.
4. Commit with "Closes #N" for that task only, and write "G3 #N" for upstream issues. Push to main and tick the task in #19 and the stage issue.
5. File every bug, ours in this repo and upstream ones in playbymail/olympia-g3, assigned to @me. Fix our own bugs as tasks in this stage.
6. If an upstream bug blocks progress, stop and ask me.

Keep tests fast:
- Default runs use the recorded cases. Whole-game comparisons go behind --full / G3G0_FULL=1, as in stage 0.
- Wrap every C or Go engine run in tools/lib/timeout.sh's run_limited.
- Write as little to disk as you can.

Carry over from stage 0:
- On macOS, C abort() flushes stdio and exits 134. SIGSEGV flushes nothing and exits 139. Go assertions panic with assertion{file,line,cond}.
- The RNG streams are a contract. Each stream's tag, keys and draw order must match doc/claude/rng-streams.md and C exactly.

Known open bugs, none blocking: ours #53, and upstream G3 #310, #321, #322, #328–#334.

When the stage exit check passes, close the stage issue, summarize all open bugs, and call out any that must be fixed before stage 2. Then pause. Don't plan stage 2 until I say so.
