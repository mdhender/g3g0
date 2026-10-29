# Suggested course correction: gameplay fidelity, not exact C equivalence

## Status and authorization

This is a proposed prompt for the owner to review and launch. Saving it does
**not** adopt a new compatibility policy, revise existing acceptance, authorize
engine changes, or start another stage. Until the owner approves the change,
WORKFLOW.md and the current issue requirements remain authoritative. Do not
use this proposal alone to waive an active task's checks.

## Prompt

Recommend and, after explicit owner approval of the contract below, plan the
course correction for mdhender/g3g0 from exact C equivalence to **gameplay
fidelity to Olympia-G3**. The goal is to finish a maintainable, playable Go
engine implementing the same game, not to reproduce every C implementation
artifact. This is a reduction in the compatibility promise, not proof that
the remaining engine work is already done.

Do not restart the port, discard accepted exact work, introduce new game
rules, or substitute a fixture-specific simulation. Do not create another
thread for this work. Separate policy migration from engine changes; propose
the concrete policy and plan edits first. Obtain owner approval before
applying conflicting workflow/tracker changes or any deliberate gameplay
deviation. A recommendation to relax proof is not authorization to relax
behavior arbitrarily.

### 1. Recover the live state and preserve what works

Read AGENTS.md, WORKFLOW.md, tools/README.md, stage-4-prompt.md,
stage-4-continue-prompt.md, #19, #183, family trackers #184–193, their current
task issues/comments, and the existing Go owners. Audit all open bugs and
follow-C issues, including unlabeled issues. Inspect main, origin/main and
the worktree; preserve concurrent changes. Treat all numbers below as a
dated snapshot, not a substitute for that audit.

At drafting on 2026-10-29, the recorded accepted baseline is engine
0.136.1-alpha, tasks #194–214 plus bounded #288, with affected tools pinned
to 7064d6c1767ba09772dc5e12af57719806184258. Other tools have independent
pins. The original spec records 56 pass / 71 fail / 0 error at that pin;
golden initialization remains unsupported. #289–291 track later merged-C
work; #215 is the next canonical gameplay task after priority work. Recheck
all of this before making a plan; do not repeat or overwrite work delivered
since this snapshot.

Keep existing engine behavior, exact tests, recordings, immutable provenance,
serializer support and accepted historical evidence. Keep cheap exact tests
running by default unless a separately reviewed change requires migration.
Do not delete tests, weaken assertions, regenerate goldens from Go, or rewrite
accepted byte prefixes merely to make the policy change appear green. Retain
strict tools as an optional diagnostic/legacy-compatibility lane even where
their expensive proof is no longer a release gate.

### 2. Define the compatibility that remains required

Adopt a written contract with these minimum requirements:

- **Game rules and observable outcomes:** preserve order syntax and meaning,
  validation/refusal conditions, priorities, turn/day sequencing, action
  durations, interruption/revalidation, costs and consumption timing, skill
  progression, ownership, inventory, loyalty, knowledge/visibility, movement,
  trade, combat, magic, NPC behavior and world evolution. Complete reachable
  start/completion/interruption paths; passing one refusal is not completion.
- **Player information:** reports must give the same gameplay facts to the
  correct players, without information leaks. Preserve meaningful ordering,
  entity references, quantities, dates, notices, failures and actionable order
  templates. Cosmetic wrapping/spacing may differ; missing events, altered
  quantities or disclosure to the wrong faction may not.
- **Supported inputs and integration:** enumerate supported command/CLI modes,
  encodings, saved-order syntax, configuration, file paths, report consumers
  and export schemas. Preserve externally consumed fields, IDs and meaningful
  ordering. Byte details remain required wherever a documented consumer
  actually depends on them; inventory that dependency before relaxing it.
- **Persistence:** existing supported C saves must load with the same game
  meaning and existing IDs/references. Go save/reload must preserve world state,
  queued/running orders and reproducibility. If historical quirks are valid
  deployed input, retain their reader support. Unless separately approved,
  keep the existing C-readable writer format; demonstrate C loading a Go save
  and semantic state equality on representative worlds. This is format
  interoperability, not a promise that alternating C/Go turns replay identically.
- **Randomness:** preserve game-defined probabilities, bounds, weighting and
  independence needed for fairness. Same Go version + seed + supported inputs
  must replay deterministically, including across save/reload. Record the RNG
  algorithm/version and seed persistence contract; cross-version replay changes
  require explicit migration/release notes. Preserve current C-exact RNG while
  it works; relaxing the gate does not itself justify replacing it.
- **Reliability:** valid gameplay must not panic, hang, corrupt saves or lose
  orders. Invalid input must produce a bounded, understandable failure without
  committing a corrupt world. Preserve operationally significant success/error
  status and required output locations, but not C's signal or assertion text.

The pinned C engine remains the primary rules reference, not an infallible
specification of undefined behavior. Do not silently fix observable C gameplay
bugs: document/reproduce the discrepancy, file upstream, and seek an owner
decision. Propose a workflow amendment permitting safe handling of undefined
behavior, crashes and hangs without waiting for a C fix, provided it preserves
valid game semantics and is recorded as a divergence. Do not infer a replacement
game rule from a crash; unresolved gameplay consequences remain blockers.

### 3. Retire proof obligations, not correctness obligations

Propose removing these universal per-task/stage gates after policy approval:

| Old proof obligation | New required evidence |
| --- | --- |
| Raw byte equality for every report, diagnostic, JSON/TSV dump and save | Parsed semantic equality plus targeted exact tests for real format/consumer contracts; retain raw differences for review |
| Identical signal, abort code, source line, diagnostic order or undefined C arithmetic/scanf behavior | Safe Go failure, unchanged valid-input rules, integrity/recovery tests and an explicit divergence record |
| Every cwd path/type/mode/empty directory/byte and cached FILE/inode/allocation side effect | Required files and permissions, correct destinations, safe writes/cleanup and no unauthorized writes; no C-internal lifecycle imitation |
| Identical C RNG streams/tags/keys/draw counts, allocation order and random aftermath | Go replay, persistence, probability/fairness tests and rule-level comparisons; preserve IDs already in an imported world |
| Two fresh C runs and full-world serializer/corpus replays for every bounded task or pin move | Targeted real-C comparisons for affected rules; broader representative integration at family/stage/release boundaries |
| Reproducing an expected C abort such as release-swear-np as Go acceptance | Retain historical abort evidence; test the owner-approved safe Go behavior separately, never invent a successful C snapshot |

Do not normalize away gameplay facts, merge distinct errors indiscriminately,
ignore all ordering, renumber imported entities, or whitelist unexplained state
differences. Freshly allocated IDs may be compared by an explicit one-to-one
entity correspondence only if all references and player-visible relationships
remain correct; existing IDs must remain unchanged. Label every normalization
and test that deliberately wrong state is rejected. A C crash, missing dependency,
NYI, stub, sanitizer failure or timeout is not a usable successful oracle.

Retain immutable C pins and source review for chosen comparisons, isolated
upstream exports, bounded processes, physical disposable world copies,
protected caches, local vet/tests and honest platform limits. Dropping filesystem
equality does not waive file safety. Existing exact failures remain visible in
their original lane; do not relabel them exact passes or silently alter
tools/olympia-spec/passing.txt.

### 4. Build a representative, auditable C comparison suite

Reuse pinned upstream frozen scenarios and existing real engine dispatch.
Select concrete fixture names and inputs in a committed manifest before
implementation. Start with the families below, inspect their actual scripts
and dependencies, and supplement missing outcomes with downstream semantic
tests based on inspected C rules. Do not edit upstream fixtures or manufacture
C expectations. Propose explicitly amending the current upstream-only coverage
policy so local semantic tests are allowed; send shared coverage requests upstream.

| Slice | Representative scenarios and assertions |
| --- | --- |
| Baseline turn and phases | Pinned world no-orders turn; all six phase-* fixtures; calendar, upkeep, production, weather/NPC and report facts |
| Orders and ownership | attitude-set/refused, stacking-stack/refused, inventory-give/refused, loyalty-bind, token-swear/vanish; asymmetric factions, quantities, privacy and references |
| Movement | move-go, move-out-first, coast-rocks, wait-ferry-range, plus actual sail/fly/raft fixtures; duration, capacity, arrival, denial and interruption |
| Economy and skills | production-harvest/refused plus inspected construction, recruitment, trade and skills fixtures; resources, partial work, completion and delayed-target invalidation |
| Combat and magic | Actual combat-*, guard-pillage and magic-/use-* representatives after source inventory; casualties/capture/loot, aura/charges, visibility, invalid target and insufficient resources |
| Cross-family lifecycles | stop-running and use-quick-cast when their real owners land; next-order execution, partial progress, timing and potion consumption; never replace NYI with success |
| Operational contracts | GM/player refusal and missing-argument cases, valid C-save import, C reading a Go save, multi-turn Go reload/replay, report/export consumers and safe unwritable-output handling |

For each family include success, refusal and a boundary/failure case, plus a
multi-day/interruption path where applicable. Add at least one representative
multi-turn world with several interacting factions/orders, not just isolated
commands and the mostly idle baseline. Record fixture pin, input/seed, commands,
supported modes, compared fields, tolerances, expected outcomes and owner.

For deterministic cases compare C and Go semantic state and report facts
exactly. For random paths first exploit existing matching RNG to compare exact
outcomes where practical. If a reviewed change legitimately breaks lockstep,
compare deterministic costs/eligibility/invariants exactly and evaluate random
outcomes over a fixed, predeclared seed set using C distributions and
rule-derived probabilities. Set sample sizes and acceptance bounds before
seeing Go results; choose power to reject meaningful balance errors, and report
uncertainty. Do not accept a single lucky result or equality of averages alone.
Treat divergence needing an impractically large sample as a testing/design
problem, not permission to waive the check.

Keep original upstream script results separate from semantic suite results:
an unchanged hash-based script can fail while its semantic case passes.
Record both statuses and the reviewed reason. Test the semantic comparator
with wrong quantity, wrong owner, missing event and leaked-information controls.

### 5. Migrate the remaining plan explicitly

1. Publish an owner-reviewable contract, non-goals, consumer/save compatibility
   inventory, divergence ledger and representative-suite manifest. For each
   existing check mark **retained required**, **diagnostic exact lane**,
   **replaced with named semantic evidence**, or **deferred by stage scope**.
   Every changed requirement needs a rationale and replacement coverage.
2. After approval, update WORKFLOW.md, #19/#183, affected family/task issues,
   tools/README.md and active continuation prompts together so they no longer
   give contradictory instructions. Preserve accepted history and annotate
   superseded future requirements rather than rewriting past evidence. Read
   live issue bodies/comments immediately before editing; preserve others' edits.
3. Reassess #289–291 against the new contract. Keep real missing-argument safety
   fixes and complete ATTACK/fog owners; propose replacing disproportionate
   pin/full-corpus proof with affected-rule evidence. Do not close a proof task
   as implemented merely because its requirement was retired: record the owner
   decision, retained work and replacement task. Preserve the stage-5 console
   exception and keep its remaining obligations visible.
4. Run the semantic baseline and retained cheap exact checks on the unchanged
   accepted engine. List every unsupported path, failure and coverage gap.
   Prove comparator negative controls. Do not call a changed gate a fixed engine
   or infer semantic acceptance from historical exact counts alone.
5. Reuse #215 onward and the existing ownership/dependency plan. Remove work
   that exists only to imitate C artifacts, not complete gameplay callbacks.
   Keep bounded tasks and bug priority; split new real prerequisites before
   coding. Do not create duplicate stages or silently accelerate later-stage modes.
6. Pilot the contract on the next bounded gameplay task and one cross-family
   integration case before using it as a stage-exit gate. Report correctness
   findings and validation cost; adjust inadequate coverage through review,
   not by weakening a failing threshold after seeing results.

No stage-5 world creation/immediate-mode, stage-6 bulk modes or stage-7 scripting
implementation is authorized by this prompt. Revise their future acceptance
wording to the approved contract without calling them complete. Scenario replay
does not require porting Lua authoring. Golden-flow semantic acceptance and
the original byte-exact golden YES are distinct claims.

### 6. Acceptance, risks and reporting

A bounded task is accepted only with reviewed complete in-scope gameplay,
passing go vet ./... and fresh bounded go test ./..., retained affected exact
checks, named semantic cases with discriminating expected results, and no
unexplained regression in previously supported gameplay. Use the repository's
version/commit workflow; policy-only documentation need not bump the engine.

A family/stage exit additionally requires all in-scope callbacks and assigned
semantic scenarios, representative multi-turn interaction, save interoperability
and Go replay, with no unexplained semantic failures. Enumerate all remaining
NYIs, original-script failures, approved divergences and source-backed later-stage
deferrals. Require explicit owner approval for gameplay changes or new deferrals;
a red test is not a deferral. Run the agreed broader suite at release boundaries,
with supported platforms stated, rather than claiming untested portability.

Address these risks in the plan and evidence:

- Semantic projections can hide lost state, privacy leaks or meaningful order:
  review compared fields and use negative controls and cross-family worlds.
- Relaxed RNG parity can change balance and historical campaign replay:
  preserve the current algorithm unless necessary, test distributions and
  version replay expectations. Equal seeds need not imply equal C/Go futures.
- Relaxed serialization/report bytes can break tools or operators:
  inventory real consumers, test import/export and document any migration.
  Do not assume bidirectional campaign interchange from a single save test.
- Fixing C crashes can accidentally change rules: isolate invalid-input safety
  from gameplay decisions and record upstream issues/owner-approved deviations.
- A smaller suite can miss rare interactions: keep exact tools runnable,
  promote real regressions into semantic coverage, and broaden release sampling.
- Changing metrics can inflate progress: retain historical exact results and
  separately report semantic passes, exact-script passes, failures and NYIs.

Deliver the proposed/approved policy status, concrete plan edits, selected suite,
commands and results, open decisions and remaining work. Never report this as
“drop-in compatible,” “byte-compatible,” “all upstream tests pass” or “C replay
compatible” without separately satisfying those specific contracts. Use wording
such as **“gameplay-faithful Go port for the listed supported features; semantic
comparison against C at [pin]; documented differences and unsupported modes.”**
State separately which save formats, report consumers, modes, replay versions
and platforms were actually tested. Completion of this migration changes the
engineering target; it does not by itself finish stage 4 or the port.
