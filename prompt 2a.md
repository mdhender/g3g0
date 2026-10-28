Continue stage 2 of mdhender/g3g0 #19: “Output and reports,” tracked in #132. Resume at #140, route rendering. The previous thread stopped after completing and pushing #139 at my request so this work could continue in a new thread. Do not restart stage planning or begin stage 3.

Read first:
- AGENTS.md, WORKFLOW.md and tools/README.md, especially check tiers and C-source dependencies.
- #19 and #132, including comments and acceptance evidence. The stage issue's original baseline-in-progress wording is historical; its later acceptance comment and task checklist record the accepted baseline and completed work.
- #140 and each subsequent task issue before implementing it. All remaining tasks #140–#166 are already created, assigned and attached to #132.
- stage-2-prompt.md for the original requirements, interpreted with the accepted progress below; do not repeat its already-completed pin move and planning.
- Previous thread if useful: https://ampcode.com/threads/T-01a18264-2704-773c-b2b2-22efe6f30cec. Durable requirements and evidence are also in issues and repository documentation; do not depend on access to its temporary files or agents.

## Starting state

Last accepted engine/task commit on origin/main:
https://github.com/mdhender/g3g0/commit/22669a8f598d234ea9519cfe119c75d02ba54041

Version: 0.79.0-alpha. This handoff prompt is a later documentation-only change and does not bump that version. Inspect the actual checkout and origin/main rather than assuming its HEAD is still the task commit. The previous thread left a clean working tree, no implementing agents running and no #140 implementation started. Do not create a replacement stage or duplicate task issues. Inspect ALL current open issues again; none were open bugs/follow-C at the last task audit, but that may change.

Immutable stage C pin:
e733acb63bacb21167a1f2589d8dcdc023364896, including G3 #357.
Previous stage-1 pin: 6a25f2d252ce45bd5adbe3baaccb3be3eed56279.
Do not silently follow current upstream main. Some unaffected tools deliberately retain their own older pins; code-compat's complete upstream dependency set is identical at the stage pin. Do not repin unrelated tools.

Accepted and pushed tasks, already ticked in #19/#132:
- #131: follow G3 #357 and move affected tool pins/provenance. USE 886 remains unported, explicitly NYI; its follow-C requirements are recorded for stage 4.
- #133: expected-abort corpus accounting, without manufacturing a successful release-swear-np+turn snapshot.
- #134: strict whole-tree turn comparator, with NO report exclusions.
- #135: bounded raw file-byte observation in the existing command oracle.
- #136: MASTER routing, styles/html, vectors, section primitives and exactly-once formatting through the existing output/cache owner.
- #167: separately filed/fixed local oracle parser bug; empty name/plural text must not reuse stale prior text. Engine unchanged.
- #137: names/codes audit. No engine gap found; existing code.go retained, direct exact-C coverage extended.
- #138: complete location headers and barrier rendering.
- #139: ship/location/road/storm descriptions, displayOwner and real shipCap/shipWeight. Existing character descriptions unchanged.

## Accepted evidence and remaining red gate

Before #131, the parent independently re-established the inherited baseline at the new immutable pin: all 196 inherited deterministic roundtrips, 201 meaningful cases overall, full box/boxfile/check/loadsave across 202 worlds, the real no-orders turn's 190 state files/raw streams/status and all 13 lifecycle controls. Full check/loadsave had zero check_db diagnostics. Spec was 19 pass / 104 fail / 0 error, with the exact same 19 passing tests; three added upstream fixtures explain growth from the old 19/101 baseline. Golden initialization still stops at unsupported -i and is not this stage's exit.

#133 now explicitly validates the missing post-turn snapshot: release-swear-np's actual C turn must exit 134 at swear.c:180, without dependency/sanitizer markers, twice deterministically. Full roundtrip acceptance is 201 PASS / 0 FAIL / 0 ERROR / 1 separate EXPECTED_ABORT. Full-corpus tools enumerate those 201 saved inputs plus their separate base world: 202 worlds. Other failed turns remain errors.

#134 removed every artifact exclusion. The Go comparator and unchanged tools/lib/tree-modes.py enforce paths, types, empty directories and all-entry permission/special bits, including directory/root modes (a stronger policy than regular-file modes alone), regular-file bytes, and raw streams/status. Unsupported nonregular entries or failed captures/listings are harness errors (2), distinct from mismatches (1). Defaults remain small and run no world turn.

The last independent strict full run was after #136, not after #139: two fresh C turns matched all 233 cwd entries. C and Go both exited 0, with identical raw stdout/stderr and no path/type/mode differences. Exactly eight saved-report byte differences remained:
lib/save/1/{100,200,202,203,204,205,210,players}.
This is an honest pending-renderer gate, NOT passing stage-2 acceptance. The strict harness stops there, so the strengthened 13 lifecycle controls have NOT yet received new strict artifact-inclusive acceptance. Their previous passes used the stage-1 artifact-excluding contract. Recheck the current full result when appropriate; do not assume these exact eight differences persist forever or normalize them away.

Latest #139 parent acceptance: command C/C/Go 903 sections / 62031 lines / 1160307 output bytes, prior recording prefixes byte-exact. Recording SHA-256:
9d4c3cc9f920c22d921dfe865c875c13b0e2d73a5c74d8a74488aa8818ce8f40.
Independent command/move/helper/code/loadsave/turn-driver/default-turn/RNG, go vet and fresh bounded default tests all passed. These routine task checks are not a rerun of the full corpus or stage exit. Task issue closing comments and tool READMEs retain precise evidence.

## Next task and remaining order

Start #140 only after the current open-issue audit. Its owner is existing dir.go. Port complete unchanged dir.c list_exit_extras 839–864, list_exits_sup 879–930, list_road_sup 933–970, list_exits/list_sailable_routes 973–1072 and their real dependencies. Reuse actual exitsFromLoc. Cover hidden/see_all 0/1/2, IN filtering, roads/transit, barriers/Hades toll, empty exits, ocean routes and Uldim/summerbridge text. Paired aftermath must detect extra exit shuffles/RNG draws. No SAIL or other movement handlers.

Continue sequentially using the detailed issue bodies:
- #141 markets; #142 visible characters/owner stacks; #143 inner locations/ship sightings; #144 ancillary location sections/posts/weather/rulers; #145 complete location-report root.
- #146 inventory/unclaimed; #147 capacity/item skills; #148 admissions/attitudes/player entity rendering; #149 known/partial skills; #150 character detail sections; #151 complete character root.
- #152 templates; #153 banners/account footer; #154 unit summaries; #155 stack/capacity tables; #156 player/storm/garrison reports.
- #157 summary totals/rankings; #158 complete summary; #159 lore sheets; #160 skill lore/queues/GM listing.
- #161 GM economy/items; #162 GM skills/geography/loyalty; #163 complete GM; #164 include assembly; #165 auxiliary files/saved reports/Times; #166 strict integration and final acceptance.

## Required per-task process

1. Inspect ALL open local issues, including unlabeled ones, before every task. Work bugs/follow-C first. Label and attach applicable follow-C issues to the stage. File every newly discovered local defect here and upstream defect in playbymail/olympia-g3, assigned to @me. The local oracle parser defect #167 is an example of a separately accepted prerequisite. If a new prerequisite or oversized closure appears, file/split it before implementation, not silently expand the active task. Stop and ask me on a blocking upstream defect.
2. Launch exactly one implementing subagent for that task. It reads its issue, stage notes and guidance, implements and leaves changes UNCOMMITTED. Parent owns independent review/acceptance and delivery. Do not hand the whole remaining stage to one subagent or create replacement threads for ordinary task phases.
3. Independently read the diff and complete pinned C yourself. Require exact real-C byte/state checks, go vet ./..., fresh bounded go test ./..., relevant compatibility tools and a version bump. New features below v1 use minor bumps; fixes/harness work use patch bumps. No expected output from Go, rewritten C models or successful missing-dependency stubs.
4. Commit directly to main with “Closes #N” for that task only, then push and tick #19/#132. Write upstream references as “G3 #N.” No branches, PRs, tags or GitHub Actions workflows. Subscribe this thread to CI for the exact pushed head after each push. Previous pushes had no configured CI checks; local acceptance is the evidence, not a remote green check.

## Existing owners and invariants

Reuse display.go, code.go, dir.go, buy.go, weights.go, sout.go/world.go and report.go/lore.go/olympia.go projections. Replace projections with real rendering in the existing owners; never run both effect paths. summary.go and gm.go may be created as source-corresponding owners when their tasks arrive. The complete location root and posts still stop explicitly; command-compat now extracts real show_loc_header and noncharacter descriptions, not successful stand-ins.

Preserve allocation timing, open-before-pPlayer in bottomOut, bounded cached FILE/inode ownership, original-list in-place sorts versus copied skill sorts, knowledge mutation, retained duplicate lore queues, lazy all-character pledge backlinks, input reads, assertion order, output bits, scratch fields and persistent process/world state. Format once before fanout; prisoners are filtered before MASTER. STYLE_PREV is one saved previous style, not a stack. vector save/restore transfers ownership. Do not reopen log files per line or close caches merely to observe them.

The command oracle's readfile PATH observes regular files up to 65536 bytes with length-framed raw bytes and no implicit production flush/close/reset. Use life flush explicitly, or legacy readlog's intentional global flush. Preserve exact blank lines, trailing spaces and missing final newlines. Extend existing oracle cases as byte-exact prefixes and compare two real C runs before accepting new expectations. Reject missing-C-dependency markers even for expected failures. Literal trailing spaces in C recordings are intentional, not lint cleanup.

At this pin wout is out; wrap_set has no definition/calls. Do not invent wrapping. Raw html()/style lines in ordinary reports ARE in scope; separate -h generation/password/public-site lifecycle is not. Keep later-stage commands, coverage/state export, world creation, valid new-player admissions, external -A charging and -M mail sending deferred. Writing lib/email is in scope. Times open/masthead/close and saveLogdir already have real owners; verify/extend them rather than implement duplicates.

The exit world has NPC/system/silent players but NO regular players, posts or email rows. Full-turn equality does not replace separate asymmetric C-driven cases for those branches. Preserve source bugs; do not make independent game fixes. Known scanf Linux/macOS and peripheral island discrepancies remain disclosed, not converted to passing expectations.

## C source safety, caches and practical checks

NEVER touch/build the shared additional olympia-g3 checkout or ~/Software/playbymail/olympia-g3. Read only through git show at the immutable pin or git archive into tmp/. Use an isolated bare tmp/upstream.git; clone it if absent, fetch into it rather than the shared checkout, and pass G3_REPO="$PWD/tmp/upstream.git" to EVERY compatibility/affected tool. Do not assume a new orb has this thread's untracked caches, binaries or logs.

The prior orb has exact-pin tmp/stage2-roundtrip with g3.rev, an exported C source/build, and accepted cases. Validate pin/build/cache provenance before reuse. C_CACHE="$PWD/tmp/stage2-roundtrip" is the reusable turn-tool cache; WORK="$PWD/tmp/stage2-roundtrip" is the roundtrip cache. Inspect scripts for full-corpus game-list variables. Rebuild missing inputs from pristine pinned fixtures, not prior engine outputs.

Run every C/Go engine with tools/lib/timeout.sh run_limited. Test executables use:
go test -count=1 -exec "bash $PWD/tools/lib/run-limited.sh" ./...
Run shared-scratch/cache tools sequentially. Use tools/lib/tree-modes.py for portable failure-aware type/mode comparisons. Keep defaults small; whole-world checks stay behind --full / G3G0_FULL=1. Remove matching captures promptly, retain decisive logs. For shared-output changes, use the command/move/helper/code/loadsave/turn-driver/default-turn tier; route/weight/RNG changes also require the relevant move/RNG gates. Follow tools/README.md for full ordering.

Never run an engine against cached input or pass a symlink as an input-copy source. A temporary #133 negative-control script once did that and contaminated one input after an otherwise passing run; the parent rebuilt the case from the pristine scenario and reran all 201 successfully. That mistake was validation-script contamination, not a reason to normalize comparisons or reuse partial failed worlds.

## Stage exit and stopping point

On identical copies of pinned world.tgz, run C twice and Go with:
olympia-g3 -r -l ./lib -S test-use-const-report-date

Require identical complete working directories, including empty directories, entry paths/types/modes, every regular-file byte, reports, Times, auxiliary files and paths outside lib, plus raw stdout/stderr/status. No exclusions or whitespace/routing/order/diagnostic normalization. Then pass all 13 lifecycle controls under that same strict contract.

Final gates: go vet ./..., fresh bounded G3G0_FULL=1 go test ./..., every required default/full tool in tools/README.md order, all inherited/meaningful new deterministic roundtrips and expected-abort control, full box/boxfile/check/loadsave corpus, strict world/lifecycle comparison, and spec with no lost passes against the 19-pass baseline. Disclose documented platform/peripheral exceptions honestly.

Only when the stage exit and final tier pass, close #132, update #19 and acceptance evidence, summarize every open bug and anything required before stage 3, then PAUSE. Do not plan or implement stage 3 until I say so.
