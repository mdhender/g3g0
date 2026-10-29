# Command dispatch, scheduler, waits, interruptions and ownership transfers

Pinned reference: [Olympia G3 at
f73e48118805ffe8ff5e763bfa4f20ae00c3a705](https://github.com/playbymail/olympia-g3/tree/f73e48118805ffe8ff5e763bfa4f20ae00c3a705).

## Form nobles (#200, stage 4 N03)

**0.122.0-alpha** adds the three complete c1.c definitions
`form_new_noble` **947–974**, `v_form` **977–1002**, and `d_form`
**1003–1067** (the return-type line is 1003, not 1005). Their unchanged
extractions occur exactly once, **1975 source bytes** total; individual
SHA-256 proofs are in `tmp/issue200-source-proof.log`. Actual callers are
glob.c's FORM start/completion slots and d_form's private formation call.
The genuine NULL interrupt remains NULL. Existing c1.go/glob.go own the port;
setName/changeBoxKind/setWhere/setLord/joinStack, NP, command initialization,
orders, report routing and RNG owners are reused. No new source dependency,
pin, package or affected.sh/table change. `noble_cost(pl)` is constant **1**.

Sixteen appended status-0 cases distinguish selected versus first unformed,
invalid foreign/nonexistent/character targets and diagnostic-before-fallback,
duplicate removal, empty/exhausted and later-populated queues, default/empty/
trimmed names, exact NP charging after formation, city-before-NP start checks,
NP depletion/replenishment and changed owner/location/behind during the delay.
Fresh hire, 100 health, 80/80 attack/defense, 50 break, retained unrelated
character fields, and new noble command initialization/retained orders are
observed directly. Actual queued dispatch and daily_command_loop run days
1–7; successful starts alone are not the test. Command-only persistence uses
the original print_command/read_boxes/scan_command through a thin io.c
adapter, discards process-local scheduler queues, then calls actual
initial_command_load. It covers successful relocation outside the city,
completion NP failure before target validation, and genuine NULL interruption
with no charge/formation and subsequent WAIT orders intact. This is a real
serialized command reload, not a claim of a new full-world corpus run.

Raw streams/status, save_box bytes, allocations, unsaved fresh-hire state,
lists, queues and RNG aftermath match C. Exhaustion and no-FORM controls both
mint **77542**, then yield entity **563/3/4/2** and legacy **38266**. Formation
also borrows the existing cached FILE after log/200 is renamed: both completion
messages reach the **155-byte** held inode, log/200 stays absent before/after
close. No parallel formation, allocation fallback or output-cache owner.

Unchanged real C ran twice before every C-only update/Go comparison. Final:
**1474 sections / 115418 output lines / 2285364 output bytes**; input
**32334 lines / 614892 bytes**. Additions: **16 sections / 407 input lines /
6219 input bytes / 1527 output lines / 28420 output bytes**. Entire HEAD
prefixes **608673 input / 2256944 output bytes** and stage-3 prefixes
**578307 / 2134398** remain byte-exact, including all **13 non-UTF8 bytes**
and authentic trailing spaces. `tmp/issue200-prefix-proof.log` records hashes;
`tmp/issue200-c-cache-final.log` records final C/C/update/Go. The pre-port
Go test fails explicitly at `command form start` (`tmp/issue200-go-before.log`).

| file | SHA-256 |
|---|---|
| command-cases.txt | `d6a60edbada701d15f81a694b26b79dc0cca1e267ceba4ff382074e1c8342011` |
| command-c-output.txt | `ce9ae9ece4b418600f9ec27d53acdaddd851be612310b643c9305ff84555bf94` |

Four temporary linker controls emit a dependency marker from actual v_form
and exit **0/134**, in default/update modes. All reject section **1458** with
harness status **2** before status acceptance/record replacement; the
then-current C recording stays unchanged. Bounded direct calls verify both
statuses. Logs: `tmp/issue200-negative-{0,134}-{default,update}.log` and
`tmp/issue200-negative-statuses.log`. Probes/build helpers and matching
captures are removed. Bounded vet and fresh full-enabled tests pass across
**11 packages / 10 explicit opt-in/child-entry skips**; eight sequential
defaults pass at their existing pins. Logs: `tmp/issue200-{vet,tests}-clean.log`
and `tmp/issue200-TOOL-final.log`. Reproduce with N01's bounded command block
below (explicit isolated G3_REPO on every compatibility invocation).

Both unchanged original fixture scripts ran real C twice, then fresh Go.
C/C match raw status/streams and **241 inventory-give / 239 recruitment-refused
cwd entries**, including types/modes/bytes; original scripts print YES.
Go really stops at **ACCEPT start** and **RECRUIT start**, respectively
(engine status **2**, script status **1**), before those original turns can
finish. These are observed paths, not speculative combat/later-stage gates;
do not suppress NYI or claim family exit. Logs:
`tmp/issue200-{inventory-give,recruitment-refused}-{c1,c2,go}.log`, retained
Go raw stderr/stdout logs and `tmp/issue200-fixtures/*/c-manifest.txt`.
The full existing spec tool, not manual passing.txt edits, reruns at the pin:
**38 pass / 87 fail / 0 error**, no changed names (all 36 stage-3 names retained),
golden **NO**, initialization status **2**. Command:
`G3_REPO="$PWD/tmp/upstream.git" WORK="$PWD/tmp/issue200-spec" GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" bash tools/olympia-spec/run.sh`
(`tmp/issue200-spec.log`). No new verified engine bug or oversized FORM
dependency was found. Remaining mixed family callbacks and stage exit are
separate; no full-world, golden, native non-Linux or sanitizer acceptance.
Parent independently reviewed all three complete pinned definitions and the
full engine/adapter/case diff, verified complete inherited byte prefixes, and
reran bounded vet, fresh full-enabled tests and all eight sequential defaults
successfully (`tmp/parent200-{vet,tests,command,move,helper,rng,code,loadsave,turn-driver,turn}.log`).

## Admission, attitudes and contact (#199, stage 4 N02)

**0.121.0-alpha** adds complete perm.c `p_admit` **27–47**, `v_admit`
**92–138**, `set_att` **241–275**, `is_defend` **319–367**, `v_set_att`
**370–395** and its four wrappers **398–423**, plus stealth.c `add_contact`
**301–308**, `v_contact` **311–332** and the `cloak_lord` **617–622** leaf.
The existing perm.go/stealth.go/glob.go owners bind the six actual starts;
their C NULL completion/interruption slots remain NULL. Existing parser,
cmdShift, scheduler, admission/hostility/clear helpers, report routing,
hasSkill and RNG remain authoritative. CONTACT appends directly as C does;
it does not also call addContact or implement SEEK.

**Source correction:** the task handoff's first-occurrence removal premise
is not this pin. `lib/ilist.c:235–241` removes **all** matching values,
backwards. setAtt therefore uses existing RemValue, not a new list variant.
Only the selected destination sorts; the other lists retain their order.
The original definitions are extracted unchanged exactly once, with hashes
in `tmp/issue199-source-proof.log`. No new source dependency, pin or
affected.sh/table change is needed (perm.c/stealth.c were already tracked).

Fourteen appended status-0 cases distinguish first-match admission rows,
first-use reset versus repeat append, ALL unit/faction exceptions, unformed
members, duplicates/order, empty targets and arguments, absent allocations,
invalid arguments followed by valid updates (including beyond argument 8),
same-faction hostility refusal with TRUE return and the independent exemption.
Raw disposition lists establish complete duplicate removal and selective
sorting. Defence covers hostile-before-default, unit-before-faction,
target-before-target-faction, learning/known cloak fallback and independence.
Contact checks direction, duplicate retention, character/player recipient
output, hidden-unit visibility, shared-helper allocation and actual
clearContacts. Real queued dispatch composes all six verbs with foreign
STACK and WAIT, and records actual permission rejections/coverage. Original
save_box bytes, raw lists/flags/allocations, streams/status, entity and legacy
RNG aftermath are compared; no model supplies expected bytes.

Real C ran twice before each C-only recording update/Go comparison. Final:
**1458 sections / 113891 output lines / 2256944 output bytes**; input
**31927 lines / 608673 bytes**. Additions: **14 sections / 335 input lines /
6444 input bytes / 1601 output lines / 29798 output bytes**. Entire inherited
HEAD prefixes **602229 input / 2227146 output bytes** and stage-3 prefixes
**578307 / 2134398 bytes** are unchanged, including non-UTF8/trailing spaces.
Proof: `tmp/issue199-prefix-proof.log`; C observations:
`tmp/issue199-discrimination.log`; C/C/update/Go:
`tmp/issue199-after-final.log`.

| file | SHA-256 |
|---|---|
| command-cases.txt | `86a170ab7e8c771ac83bf89033bb2d87e2f0eb609ef95c26bac983d550a2b508` |
| command-c-output.txt | `fff1197f03348079b4cf703b8c113d33d5762a9c722aadca7139955a982944db` |

Four temporary linker controls emit a dependency marker from actual v_admit
and exit **0/134**, in both default/update modes. All reject section **1444**
with harness status **2** before accepting status or replacing recordings;
the inherited recording hash stays unchanged. Bounded direct executions
confirm both injected statuses. Logs: `tmp/issue199-negative-STATUS-MODE.log`
and `tmp/issue199-negative-statuses.log`. Probes/build helpers are removed.

Bounded vet and fresh full-enabled tests pass (**11 packages**, **10 explicit
opt-in/child-entry skips**); standalone whole-game skips are not corpus
acceptance. Logs: `tmp/issue199-{vet,tests}-clean.log`. Eight sequential bounded
defaults pass using the N01 reproduction command below, with explicit isolated
G3_REPO: command **1458/113891**, move **135/7158**, helper **37/1487**, RNG
**321**, code **43/969**, loadsave **61/3411**, turn-driver **7/152**, phase
**5/484**, calendar **2/119**, schedule **4/1801**, default turn raw streams/
status and **2 cwd entries**. Logs: `tmp/issue199-TOOL-final.log`.
An interim fresh suite, started before
the final C-only refresh added the cwd-localizing lore_root line, correctly
rejected the one-line stale recording; final fresh checks use the refreshed
real-C bytes (`tmp/issue199-tests-before-final-record.log`).

Parent independently reviewed the complete pinned definitions, shared list
semantics and full source/adapter/case diff, confirmed both inherited byte
prefixes, and reran bounded vet, fresh full-enabled tests and all eight
sequential defaults (`tmp/parent199-{vet,tests,command,move,helper,rng,code,loadsave,turn-driver,turn}.log`).

**Original attitude fixture gate passed:** parent rejected the initial
assumption that these fixtures required N28/combat. attitude-set explicitly
places the nobles in a safe haven to prevent fighting; attitude-refused only
exercises failed/rejected orders. Both unchanged original check.sh scripts
ran against unchanged real C twice, then the fresh Go engine, all **YES /
status 0**. Complete pre-cleanup raw streams/status and all cwd entries,
types, modes, empty directories and bytes match: **239** entries for
attitude-set, **237** for attitude-refused. The existing bounded native
capture launcher, strict turn comparator and tree-modes.py were used, with
the verified C executable SHA-256
`565ca7fb46a07d733aaff3c46bde5d2b7d124d8c44a98d3842d82ff886d55ab0`.
Logs: `tmp/parent199-attitude-{set,refused}-{c1,c2,go}.log`; retained manifests:
`tmp/parent199-fixtures/attitude-{set,refused}/manifest.txt`.
Matching captures and the temporary fresh Go binary are removed afterwards.
Actual hostile auto-attacks outside safe havens still require N28/combat;
no original script, expected byte, NYI or battle behavior is altered.
This accepts N02 and its two original fixtures, not the remaining mixed
attitude/inventory/loyalty family tasks or stage exit.

## Validation, acquisition and wounds prerequisite (#198, stage 4 N01)

**0.120.0-alpha** adds the complete pinned u.c `check_char_here` **924–930**,
`check_still_here` **965–996**, `add_char_damage` **541–579**,
`count_man_items` **1178–1197**, `diff_region` **2627–2631** and swear.c
`np_to_acquire` **187–196** / `enough_np_to_acquire` **199–212** to existing
u.go/swear.go. The oracle extracts each original definition unchanged exactly
once, including return types. `tmp/issue198-source-proof.log` records their
individual sizes/hashes and exact-pin checks. Complete direct callers were
reviewed, including disabled SWEAR code and commented necro region checks;
none is newly wired. Existing visibility/contact, greater-region, death,
survival, inheritance, interruption, NP totals and entity RNG remain owners.
No source dependency, pin, upstream harness or unrelated callback changes.

Twenty-one appended status-0 sections compare raw streams/status, real
save_box bytes and raw export allocation/list observations. They distinguish:
999 before kind/visibility checks (even with a real box at 999); invisible
versus contacted targets, contact direction, faction and same-stack visibility,
fog, relocation and named revalidation diagnostics; moving targets accepted by
still-here but rejected by check_char_gone; ordinary continents collapsed by
greater_region versus all six special realms, including nested real ships;
noble 1 / garrison 0, nonfighting man items versus nonman fighters, signed sums;
nonpositive damage on absent records/boxes; health -1 at 49/50, health 1/0/-2,
clamped INT_MAX wounds, exact message order, new/already-sick state, repeated
post-health keyed draws and next-turn/schar boundaries. Real dispatched WAIT
survives fatal-wound skill consumption, while actual death interrupts WAIT,
flushes queued orders and transfers inheritance through accepted owners.

Independent no-call versus nonpositive/already-sick/sentinel controls replace
the master seed before the first entity observation: these detect an unwanted
sickness draw's otherwise-hidden turn-guard initialization. Entity and legacy
RNG aftermath are compared without replacing C draws. Acquisition cases cover
base/oath/all-skill NP totals, short/exact/repeated checks without charging,
the independent previous-owner exemption versus unrelated acquisitions,
zero-cost NPCs, absent CH/player records, conditional CH allocation and C int/
short wrapping. Thin adapters only call the original helpers/leaf and print
results; no replacement gameplay or success adapters are added.

Four temporary linker controls wrap check_char_here, emit a dependency marker
and exit **0/134** in both default/update modes. Each is rejected at section
**1423** with harness status **2**, before status acceptance or replacement;
the recording hash stays unchanged. Direct bounded executions verify both
injected statuses (`tmp/issue198-negative-statuses.log`), and rejection logs
are `tmp/issue198-negative-{0,134}-{default,update}.log`. Probes/compiler
wrapper and matching captures are removed before the final checks.

Unchanged real C ran twice before each C-only recording update and Go
comparison. Final corrected cases (explicit colocation for visibility checks)
match **1444 sections / 112290 output lines / 2227146 output bytes**. Inputs:
**31592 lines / 602229 bytes**. Additions: **21 sections, 423 input lines /
8689 bytes, 1912 output lines / 34477 bytes**. Both entire inherited HEAD
prefixes (**593540 input / 2192669 output bytes**) and stage-3 prefixes
(**578307 / 2134398 bytes**) remain byte-exact, including non-UTF8 and real
C trailing spaces. Prefix proof: `tmp/issue198-prefix-proof.log`.

| file | SHA-256 |
|---|---|
| command-cases.txt | `c19193becde431c53dcddf2def85e56dd36e29cac6bc734a7754a69b628ce0d8` |
| command-c-output.txt | `6d2ba0f66cc973df8660cb471a7118c4102eedcd4d1f949dcfe4391f2290317e` |

The initial unknown-operation red test is `tmp/issue198-before.log`; final
C/C/update/Go evidence is `tmp/issue198-c-corrected.log`. The extracted
recording checks in `tmp/issue198-discrimination.log` show, in particular,
moving **here=1 / gone=0 / still=1**, ordinary-region **0** versus special
**1**, man totals **1/0/11/10/-2147483645**, no-draw control **599/1** versus
initialized entity aftermath **267/3**, and unrelated independent acquisition
**3 NP** versus the previous owner's **0**. These are observed original C
results, not invented expectations.

Final sequential defaults pass with explicit isolated G3_REPO and bounded
test execution: command **1444/112290**, move **135/7158**, helper **37/1487**,
RNG **321**, code **43/969**, loadsave **61/3411**, turn-driver **7/152**,
phase **5/484**, calendar **2/119**, schedule **4/1801**, default turn raw
status/streams and **2 cwd entries**. Logs: `tmp/issue198-TOOL-final.log`.
Then bounded vet and fresh full-enabled tests pass, **11 test packages**,
Olympia **1.942s**, no package exclusions (`tmp/issue198-{vet,tests}-clean.log`).
Whole-game tests still require supplied corpora; their standalone skips are
not full-corpus acceptance. Reproduction (from repository root):

```sh
for tool in command move helper rng code loadsave turn-driver turn; do
  G3_REPO="$PWD/tmp/upstream.git" GOFLAGS="'-exec=bash $PWD/tools/lib/run-limited.sh'" \
    bash "tools/$tool-compat/run.sh" || break
done
. tools/lib/timeout.sh
run_limited 120 go "$(command -v go)" vet ./...
G3_REPO="$PWD/tmp/upstream.git" G3G0_FULL=1 RUN_TIMEOUT=120 \
  run_limited 180 go "$(command -v go)" test -count=1 -v \
    -exec "bash $PWD/tools/lib/run-limited.sh" ./...
```

Source/docs diff-check passes. Full diff-check reports only **20 authentic
C recording trailing-space lines**, intentionally untrimmed. Temporary
matching command captures/build helpers are removed after validation.

This is prerequisite-helper acceptance only. Original social/inventory/USE
family scripts await their later callback tasks; no full family, full-world,
whole-spec, native non-Linux or sanitizer acceptance is implied. The separate
release-swear-np expected abort is unchanged, not a successful turn. The
all-open local audit finds **102 issues / no bug or follow-C**; no new verified
engine bug or oversized prerequisite was found. Parent independently reviewed
all seven complete pinned definitions, existing shared owners and the full
source/adapter/case diff, and confirmed both complete inherited byte prefixes.
Parent bounded vet, fresh full-enabled tests (Olympia **2.074s**) and all eight
sequential defaults pass (`tmp/parent198-{vet,tests,command,move,helper,rng,code,loadsave,turn-driver,turn}.log`).

## Province line-of-sight prerequisite (#197, stage 4 P5)

**0.119.0-alpha** adds complete dir.c **128–169**
`los_province_distance` to existing dir.go. The oracle extracts the unchanged
705-byte definition once, including its return-type line (SHA-256
`fd65deb03d868bc1bf03baf7d9c663afa6ce549163ad251d7c792440abf1d5b4`).
dir.c and the source callers (projection, FIND BUY and immediate distance)
were already tracked: no source dependency, pin or callback wiring changes.

Both endpoints normalize through province. Every province's temp/next resets
to -1/0, lazily allocating location records, even for distance zero. The BFS
uses those original next fields, full prov_dest length and insertion order;
only zero/already-visited neighbors skip. The target is tested on dequeue,
before expansion, preserving all queue/frontier aftermath. No coordinate
approximation, cardinal limit, separate queue, endpoint validation or guards.

Thirteen appended ordinary sections require status **0** after marker
rejection: same endpoint with global reset/allocation and untouched
nonprovinces; directed/disconnected reverse paths; graph versus wrapped/nearby
coordinates; nested character/stack/city/building/ship normalization; unequal
alternate lengths; swapped equal-path neighbor order; ten-entry adjacency
with zeros/duplicates/self-links/cycles; early success versus exhaustion;
changed topology/repeated queries; independent no-call versus call aftermath.
Parent review corrected the normalization fixture's original `5002 4 30`
(T_item) to **`5002 9 30` (T_ship/sub_galley)**, as pinned oly.h defines.
The actual ship links to province 19999; `los 1001 5002` normalizes the nested
stack/castle/city source to 10101 and the ship destination to 19999 (distance
1). `los 5001 1000` retains building/character same-province coverage (0).
The recording now explicitly saves **`5002 ship galley`**, not `item galley`.
Thin adapters expose the helper, append one graph edge, seed next, and observe
every province's raw temp/next/record presence/neighbor list. Existing export
state/save_box observers compare every box's saved bytes, raw lists and
allocations. The independent controls' keyed/entity/exitdir/legacy RNG and
shared mint/create/save_new tail is identical; complete final states retain
the authentic C temp differences and newly allocated bare LO section.

Cases preceded implementation; bounded preimplementation tests rejected the
new operations. Real pinned C ran twice before the C-only recording update,
then the final comparison repeated C/C and matched Go: **1423 sections /
110378 lines / 2192669 output bytes**. Input: **31169 lines / 593540 bytes**.
Additions: **227 input lines / 4540 bytes**, **1391 output lines / 28237 bytes**.
Both entire inherited HEAD prefixes (**589000 / 2164432 bytes**) and stage-3
prefixes (**578307 / 2134398 bytes**) are byte-identical, including non-UTF8
and intentional trailing spaces. No expected bytes originate from Go.

| file | SHA-256 |
|---|---|
| command-cases.txt | `5f831281ac27b0fb8466927bcbf7e3a80caa0d8b6fe22e67576d0c9a6cb1126c` |
| command-c-output.txt | `cb222a7c00a01ed7be22b0166e938a497eeef3e578d24668563c13ceab6cbb5a` |

Initial red evidence: `tmp/issue197-{before,c-first}.log`. After the ship
fixture correction, unchanged C ran twice before refreshing the unaccepted
P5 recording suffix, then Go matched, including live native traversal.
Current evidence: `tmp/issue197-corrected-{command,source-proof,prefix-proof}.log`.
Command/move/helper/rng/code/loadsave/turn-driver/default-turn pass sequentially
at their unchanged independent pins with explicit isolated G3_REPO
(`tmp/issue197-corrected-{command,move,helper,rng,code,loadsave,turn-driver,turn}.log`).
Default turn checks the diagnostic's raw status/streams and **2 cwd entries**,
not a world turn. Literal bounded `go vet ./...` and fresh full-enabled
`go test -count=1 -v -exec "bash $PWD/tools/lib/run-limited.sh" ./...` pass
after cleanup, all **11 test packages** (Olympia **2.223s**;
`tmp/issue197-corrected-{vet,tests}-clean.log`), with no package exclusions.

Four temporary linker controls wrap the real LOS entry, emit a rejected
dependency marker and exit **0/134**, each in default/update mode. All stop
with harness status **2** at section **1410**, before status acceptance or
recording replacement; the recording hash stays unchanged
(`tmp/issue197-negative-{0,134}-{default,update}.log`). Temporary helpers,
executables and matching captures are removed. Source/docs diff-check is
clean; the full check flags only authentic new C recording trailing spaces,
which are deliberately retained. No bug/follow-C appeared in the all-open
local audit. The parent-identified fixture-kind defect is corrected; no
engine defect or blocker was discovered, and the engine code is unchanged
by this correction. Version remains **0.119.0-alpha**.
Parent independently reviewed the complete pinned definition and adapters,
confirmed the corrected ship record and both accepted byte prefixes, and
reran bounded vet, fresh full-enabled tests and all eight sequential defaults
successfully (`tmp/parent197-{vet,tests,command,move,helper,rng,code,loadsave,turn-driver,turn}.log`).
No family/projection/trade/whole-world/stage-exit acceptance is implied.

## Skill allocation/grant prerequisite (#196, stage 4 P4)

**0.118.0-alpha** adds complete use.c **1016–1035** `p_skill_ent` and
**1857–1892** `learn_skill` to existing use.go. The oracle extracts both
unchanged definitions once (SHA-256, including return-type line:
`ffc33e5389857f0268b8ffe15c02c557ff77a035fd3ea1c4a3887a8442fe928c` and
`d80087723d13b92515bfe79651ff6fcf6654c906552b54991c1ab659ecd4e431`).
use.c was already tracked; no new source dependency, pin or dispatch wiring.

Allocation searches the original list for the **first** matching skill,
regardless of knowledge, otherwise appending one zeroed row. Grant always
prints Learned, sets only know in that row, adds **50** to missile only below
50 (not set-to-50), and unconditionally allocates magic, even for nonmagic.
Magic-school ancestry adds current/max aura and prints the maximum on **every
grant**, including already-known skills; magician becomes 1, and weather's
exact school ID sets knows_weather only inside that magic branch. No cost,
validation, idempotence repair or sorting is introduced.

Thirteen appended ordinary sections require status **0**, after rejecting
dependency/NYI/stub/sanitizer markers. They distinguish missing/partial/known
and duplicate-order rows, append order, retained signed study/experience/
exp_this_month and skill-use counters, absent character/magic records, plain
school/subskill, nested magic ancestry and weather subskill versus school,
nonmagic weather ID, archery -73 repeats and 0/49/50/51/-32768/32767, negative
aura and signed 32-bit wrap/repeated grants. Raw row/record observers supplement
actual save_box and exact logs; player NP is unchanged. Independent no-call
and grant controls have identical magic/entity/legacy/shared-mint aftermath.

Cases preceded implementation; the preimplementation Go run rejects the new
operations explicitly. Real pinned C ran twice before each C-only recording
update, then Go matched **1410 sections / 108987 lines / 2164432 bytes**.
Input is **30942 lines / 589000 bytes**. The entire inherited **584230-byte
input / 2153045-byte output** prefixes remain identical, including all original
stage-3 bytes and intentional trailing spaces. New additions are **215 input
lines / 4770 bytes**, **615 output lines / 11387 bytes**.

| file | SHA-256 |
|---|---|
| command-cases.txt | `631ccfa582ec096f232a36b8e3a8d8dfec3d06cfd511dad5f9798a1de454f262` |
| command-c-output.txt | `5af45da01ac26f14906ad94078bf0d0adfbc9c7fd2623210a676713eb801f57e` |

Evidence: `tmp/issue196-{before,c-first,after-final,prefix-proof,source-proof}.log`.
Final command/move/helper/rng/code/loadsave/turn-driver/default-turn comparisons
pass **sequentially**, each at its independent existing pin with explicit
isolated G3_REPO (`tmp/issue196-*-final.log`). Final command repeats unchanged
real C twice, then matches Go and the live native traversal control. Default
turn compares the diagnostic's status, raw streams and **2 cwd entries**, not
a world turn. Literal bounded vet and fresh full-enabled bounded tests pass
all **11 test packages**, again after cleanup (Olympia **1.759s**;
`tmp/issue196-{vet,tests}-clean.log`). Full `git diff --check` passes; raw C
recording bytes are not normalized.

Parent independently reviewed both complete C definitions, the magic-school
macro and source/adapter/case changes, confirmed both inherited byte prefixes,
and repeated all eight compatibility checks plus literal bounded vet/fresh
full-enabled tests successfully (`tmp/parent196-*.log`; Olympia **2.380s**).

Four temporary linker controls wrap p_skill_ent, emit a dependency marker and
exit **0/134**, each in default/update mode. All reject at section **1397**
with harness status **2**, before expected-status acceptance or recording
replacement; both recording hashes remain unchanged. Direct bounded launches
confirm both injected engine statuses. See `tmp/issue196-negative-*.log`.
Probe sources, executables and matching captures/build trees were removed.
ALL **104** open local issues, including unlabeled issues, were audited;
none was bug/follow-C, and no new bug/blocker was found in this bounded task.

This is only the shared skill-grant prerequisite for later STUDY/auraculum/
eat-dead owners, not STUDY/research/USE-family or stage-exit completion.

## Potion allocation prerequisite (#195, stage 4 P3)

**0.117.0-alpha** ports only complete alchem.c **8–42** `new_potion` into
the natural alchem.go owner. The unchanged C definition is extracted once
(SHA-256 `263519d6384cd65f112e4f62c998366af46bd5ed2eb81b97ba9c5ef0a9320a3a`).
It reuses real unique-item allocation/inventory, naming, item-magic records,
province traversal, output and persistent `MagcPotion(who,new,0,1,2)`.
`alchem.c` is the sole new command-tool source dependency, also tracked by
affected.sh. No pin changes, RNG adapters, brewing callbacks or item wiring.

Eight appended sections cover both actual keyed names, exact IDs, ownership,
inventory order and complete save_box bytes; nested actor → character → ruins
→ province versus unlocated province zero; absent actor magic records;
repeated same-turn allocations, repeat keyed probes, day and turn changes,
relocation, and preferred-range exhaustion with real fallback. Independent
no-call and equal-count plain-mint controls distinguish shared mint advancement
from unchanged sibling magic/entity/legacy RNG. Produced-one logs retain the
real item code/name and routing.

**Exhaustion is an assertion, not a returned -1, at this pin.** Complete
code.c `new_ent` calls `alloc_box(-1,...)` after both ranges are exhausted;
the actual C/C/Go cases and plain-mint control preserve status **134** and
`olympia/code.c (690): n > 0 && n < MAX_BOXES`. Neither the potion draw nor
its naming/records/output can run after that assertion. The literal negative
return branch in new_potion remains ported, but cannot be reached through
this allocator; no mock failure, dummy ID, caught C assertion or invented
successful failure-aftermath recording is used. All six ordinary sections
require **0**; both exhaustion controls require **134**, after marker rejection.

Cases preceded implementation. The old Go adapter failed explicitly on
`potion`/`entity_range`; unchanged real C then ran twice before the C-only
recording update and Go implementation. Totals: **1397 sections / 108372
output lines / 2153045 output bytes**, input **30727 lines / 584230 bytes**.
Entire inherited **582083-byte input / 2147299-byte output** prefixes remain
byte-identical, including intentional C trailing spaces (the new `hl 1000 `
line also deliberately retains its C space).

| file | SHA-256 |
|---|---|
| command-cases.txt | `d680b5533b1335fe4b9ffead00943983090a9c5aadf0f38c2ea75e10ea103f7a` |
| command-c-output.txt | `252a2e9e46ba4d781278959796c5918497add61ac795fceb65f74382e87eb801` |

Evidence: `tmp/issue195-{before,c-first,after,prefix-proof,source-proof}.log`.
Sequential command/move/helper/rng/code/loadsave/turn-driver/default-turn
comparisons pass at their existing pins, each with explicit isolated G3_REPO
and existing run_limited supervision (`tmp/issue195-*-final.log`). The final
command rerun again requires unchanged real C twice before Go, including live
native traversal. Literal bounded `go vet ./...` and fresh full-enabled
bounded `go test -count=1 ./...` pass all **11 test packages**, repeated after
scratch cleanup (Olympia **2.163s**; `tmp/issue195-{vet,tests}-clean.log`).

Parent independently reviewed the primitive and allocator definitions,
confirmed both entire inherited byte prefixes, and repeated literal bounded
vet/fresh full-enabled tests plus all eight compatibility checks successfully
(`tmp/parent195-*.log`; Olympia **1.925s**). The only full diff whitespace
warning is the intentional new C recording line `hl 1000 `; source checks pass.

Four temporary real linker controls wrap new_potion, emit a dependency marker
and exit **0/134** in default/update runs. Every tool run rejects at section
**1389** with status **2**, before expected-status acceptance or recording
replacement; both recording hashes remain unchanged. Direct bounded launches
also verify the injected engine statuses. Evidence:
`tmp/issue195-negative-{0,134}-{default,update}.log` and
`tmp/issue195-negative-statuses.log`. Probes and matching captures/build trees
are removed. ALL **105** open local issues were audited, with no bug/follow-C;
upstream G3 #309/#114/#361/#362 do not block this primitive.
This is a primitive prerequisite, not alchemy/item-family or stage acceptance.

## Current aura cap prerequisite (#194, stage 4 P2)

**0.116.0-alpha** ports complete art.c **54–75** `max_current_aura` and
`limit_cur_aura` into the existing art.go, reusing maxEffAura/hasAuraculum.
The C definitions are extracted literally unchanged, once each (SHA-256
`31dc267a7e18dd8f822ed28de91efaa1801a725eca46534e290f0d7764095a9f`,
`0104eb8b11ec9e47a599a042f13eaa08a6020ebdce606d8a0d442931b433956d`).
No new source dependency, pin move, gameplay callback wiring or RNG adapter.

Multiply effective aura by five at C int width **before** the negative clamp.
Only a zero result gets the one-aura fallback, and that checks the **raw**
recorded auraculum even when unheld, absent or invalid. Limiting only lowers
strictly excessive current aura, evaluates the cap again in that branch,
and leaves negative current aura and absent magic records alone.

**Nine appended sections** cover asymmetric held artifact/per-entry bonuses
(quantities 4/19/3 and a zero entry), below/equal/above/repeated limiting,
negative innate versus negative effective aura, held/unheld/missing/invalid
recorded artifacts and exact zero fallback, missing character/item magic
records, signed addition/multiplication boundaries and a no-call RNG/mint
control. Boundaries include INT_MIN/INT_MAX, 429496729/429496730 and
858993460: the last gives cap **4**, not a wide-integer product. Raw record
flags/inventory/list observers and complete real save_box bytes accompany
magic/entity/legacy RNG and shared allocation aftermath. All nine require
status **0**, after marker rejection.

Cases were appended before implementation; unchanged real C ran twice before
the C-only recording update. The then-unextended Go adapter failed explicitly
on the new operations; subsequent real C/C/Go comparisons pass. Final totals:
**1389 sections / 107961 output lines / 2147299 output bytes**, input
**30604 lines / 582083 bytes**. Entire inherited **578307-byte input** and
**2134398-byte output** prefixes equal HEAD byte-for-byte, including C spaces.

| file | SHA-256 |
|---|---|
| command-cases.txt | `da5596ea0a46567db3f26e0ff6738c06938f13e13ea7749180eb730fc9ab3401` |
| command-c-output.txt | `3e20877343c1ca2c0261fb8af315dc8e06076c62daa807e1c6bb3d18966232ae` |

Temporary linker controls wrapped max_current_aura, emitted a dependency
marker and exited **0/134** in both ordinary and `--update` runs. All four
reject at section **1380** with tool status **2**, before status acceptance
or recording replacement; expected bytes remain unchanged. Controls removed.
Final normal command/move/helper/rng/code/loadsave/turn-driver/default-turn
checks pass sequentially, explicitly using `G3_REPO=$PWD/tmp/upstream.git`
and each tool's existing pin and run_limited supervision. Evidence:
`tmp/issue194-*-final.log`, `tmp/issue194-negative-{0,134}-{default,update}.log`.
Parent removed its completed baseline launcher sources after an initial
scratch-package setup failure. Literal bounded `go vet ./...` and fresh
full-enabled bounded `go test -count=1 ./...` then pass all **11 test packages**
(Olympia **2.225s**), without excluding package trees. Parent independently
repeated all eight compatibility checks and confirmed both inherited byte
prefixes. Evidence: `tmp/parent194-{vet,tests,command,move,helper,rng,code,loadsave,turn-driver,turn}.log`.
No whole-world repeat, magic/USE family exit or stage completion is claimed.

## Drum USE item prerequisite (#182)

**0.115.0-alpha** ports complete original savage.c `v_use_drum` **70–146**
(**77 lines / 1483 bytes**, SHA-256
`9dc335d781d225c678f8f50fcfdf094baaa965c639e62d169d324d5fa580103d`),
and actual use.c **869–872** wiring: copy c.b into c.a, then call the drum
callback. The original shared item tail **911–920** counts the callback's
outcome and inhibits zero-time completion. Other item callbacks still stop
explicitly, never returning fabricated success/failure. The boundary test
excludes only this now-implemented key; all remaining NYI contracts stay.
No summon-savage command, ATTACK/combat, independent fix or stage 4 work.

Raw speeds **0/1** can summon; **2/negative/256** still emit sounds and shuffle
LAND exits in a province. A building/ship instead emits to its immediate
parent and returns before the shuffle/response branch. The exact original
sound order/capitalization and opposite-direction table are retained. A failed
controlled-human candidate permits a later attempt, only the first success
spawns, and **every** remaining neighbor still hears the sound. The **200**
cap counts the actual successful spawn. Existing callSavage/createSavage,
queues, initLoadSup, NPC/exitdir/mint/legacy RNG and allocation owners are reused.

All **100** archived command originals equal the immutable pin; all **583**
extractor-listed complete definitions occur literally unchanged exactly once.
The **1212** original files in read-only `tmp/issue179-roundtrip` equal the
pin; its executable retains SHA-256
`151c7355890eba7f8bfda1f4f16f4a490947aaa1db1ffcc3ca859cd659c01158`.
Audited original dependency SHA-256 values:

| definition (source lines) | SHA-256 |
|---|---|
| create_savage (savage.c 13–29) | `eff24412752d543e04e19847fc3d51b595d0a28a59de6092df72218c26913b8a` |
| call_savage (savage.c 32–67) | `08a50be096ac37b33d04a9b0d2dce52afd248ba2284ccd1a6bb87f812b243032` |
| init_load_sup (input.c 703–741) | `6e644e702057e833cddd719c3e94dd897531ad453355c27a688131ec51a7905a` |
| exits_from_loc_nsew_select (dir.c 776–803) | `734a0dc4901db925078a615a0848bd6ad71982c5f168d99a2ba2f7fd1be848b4` |
| controlled_humans_here (npc.c 65–84) | `95f1831398fec32d9be6a223c09fb179a83c97cd126d6471b92780da7ac1bcbd` |
| v_use_item (use.c 785–920) | `76a9a4db221fb71559256a6fb36b2c431afa28140bc81b0f2715b49098acf144` |

No new source dependency or pin move: savage.c/use.c and all actual
dependencies are already tracked, so affected.sh is unchanged. The real
original item/command tables and actual queues execute, without success overlays.
The small `savage drum` adapter calls the complete original callback with raw
arguments; parsed USE cases separately check the original item dispatch.

**14 appended sections** cover raw speed branches, nested building/city and
ship/province, no/ocean/controlled/hidden-barrier exits, failed candidate then
later success, cap **199→200 / 200**, four-neighbor output after success,
real queued USE/zero-time completion inhibition, command/USE accounting,
saved command arguments, exact inventory/loyalty/order/load-queue state,
shared shuffle/NPC/entity/legacy RNG aftermath and subsequent actual item mint.
The seed-2 four-land-neighbor call tries **9004** (controlled) before spawning
**129744** in **9003**, with **16** savages, summon loyalty **3**, loaded MOVE
and retained STACK. Its sounds still reach **9002/9001** afterwards.

Unchanged C runs twice before Go or recording replacement, matching **1380
sections / 107214 output lines / 2134398 bytes**, SHA-256
`c4cc3c8ee0cd5f0b40ac6eb239aa7bf51235aaf1448fa6d9aad81f5a821f1bac`.
Input: **30431 lines / 578307 bytes**, SHA-256
`3481dcd40549f16ff6fc67201a70bab1e00c08a599799151e7706734ff594659`.
The ENTIRE inherited **571691-byte input / 2105533-byte output** remain exact
prefixes, SHA-256 respectively
`c5bbd288cc49cadc94ff68dde1315f8a67f359e271a265422e8be5be2a1e06b3` /
`fe4b9575b0630d05abb40a63e18724f78266be65b41817a22e859d49c1194cd4`.
Original C trailing spaces in four appended saved here-lists are intentional;
do not normalize them to appease diff whitespace checks. New sections require
status **0**, after dependency/NYI/stub/sanitizer rejection.

Six actual linker controls wrap the direct drum call at section **1366**
(status **0/134**, each default/update) and original **asfail** at inherited
expected-abort section **27** (status **134**, default/update). Markers in
stderr or stdout cause harness **2** **before** status acceptance/recording
replacement; the recording hash remains unchanged in every control. Probe
sources/launchers/captures are removed. Logs:
`tmp/issue182-negative-{v_use_drum,asfail}-*.log`,
`tmp/issue182-negatives.log`. The subsequent normal unchanged-C/C/Go default
is `tmp/issue182-command-default-final.log`.

Logs: `tmp/issue182-source-audit.log`,
`tmp/issue182-command-update-complete.log`, `tmp/issue182-command-default.log`.
An initial fixture setup left a supposedly controlled unit unsworn, provoking
an invalid test-world C spawn/status **139**; no recording was accepted from
that run. The fixture was corrected to real controlled loyalty. The first
successful C recording exposed a Go constant spelling error, then the obsolete
drum NYI assertion; both were corrected rather than accepting failed Go output.
These initial failures remain in `tmp/issue182-command-{update,tests}.log` and
`tmp/issue182-command-update-final.log`; the final normal comparison passes.

## Hinder-meditation USE prerequisite (#181)

**0.114.0-alpha** ports complete original basic.c `v_hinder_med` **331–357**,
`hinder_med_omen` **360–390**, `d_hinder_med` **393–420**, scry.c
`reset_cast_where` **37–51** and u.c `check_char_where` **898–921**. Only
the two actual hinder-meditation USE callbacks are newly wired. Existing
chargeAura/checkAura, charWhere/contact/visibility and Streams.MagcOmen remain
the real owners. No other magic, drum support (#182), independent fix or
stage-4 work is included.

The original order is retained: clamp b to **1..3** before checking aura;
only then consume a valid projected location/ship and save it in c.d, even
when target validation fails. Completion charges aura before validating the
target at that saved location. The signed-char addition is stored **before**
the **5** cap, including **127 + 1 → -128**. Moving targets are not rejected
merely for moving; visibility/location still govern. Invalid projected IDs
remain stored. Omen's stream uses turn-root `turn` → `magc`, keyed target/
caster and tag `omen` (**0x6f6d656e**), range **1..4**, without a day child
or sequential draw. Same-turn master changes do not reseed it; next turn does.

Every archived original (**100 files**) equals the immutable pin; all **582**
extractor-listed complete definitions occur literally unchanged exactly once.
The entire original magic RNG prefix is extracted unchanged (**4969 bytes**,
SHA-256 `2d88d0f9ea8b08fdebc81aa062b489f6b7c84ace8ca851158dfb2667e6cf4134`).
New complete-definition SHA-256 values:

| definition | SHA-256 |
|---|---|
| v_hinder_med | `84ce74cad7fc370358a2f700857fca80d766cdf535bc0188039331f2dfc73a4c` |
| hinder_med_omen | `768cd27170b038e09be9ccdde131e83375c02db60ab7b90d797743e7dbb4144a` |
| d_hinder_med | `34ac1871a1aa81ad1849d8bc291d50342588006c0a1b0b5a0d2d7f64277a59e2` |
| reset_cast_where | `0ed474847e6e200997192f5579de826296d9da33194ee10634c71337d4170320` |
| check_char_where | `49cf90d1a614b2d8d1dc98850d2f772619f431785d5b74227b890dba20bb7000` |

**basic.c is a new command-oracle dependency**, tracked in affected.sh and
tools/README.md; the pin is unchanged. The original use table and real queues
execute, not callback-success overlays. `usemeta callback_start/callback_finish`
selects the actual table callback with raw arguments, alongside real parsed
USE/start/evening/finish/interrupt paths. `usemeta rolls` observes real C magic
leaves and Go's existing streams, not a replacement RNG model.

**16 appended asymmetric sections** cover clamp extremes, absent/exact/lost
aura, zero/noncharacter/missing/garrison targets, projected locations/ships,
reset versus retention, saved destination after caster movement/new projection,
target movement/kind/hidden completion revalidation, same-faction hidden and
foreign hidden/fog contacts, all **four omen branches**, signed-char wrap/cap,
actual ten-evening completion, interruption, command/USE accounting, target
allocation, skill experience/state, logs and independent magic/legacy/entity
RNG aftermath. No-call, interrupted and successful controls retain legacy
**28797**, entity **57/2/4/1**, and magic **1/44/83** at seed 2. The silent
fourth omen is independently observed as **4/13/9**, with no target log line.

Unchanged C twice and Go match **1366 sections / 105616 lines / 2105533
output bytes**, SHA-256
`fe4b9575b0630d05abb40a63e18724f78266be65b41817a22e859d49c1194cd4`.
Input is **30049 lines / 571691 bytes**, SHA-256
`c5bbd288cc49cadc94ff68dde1315f8a67f359e271a265422e8be5be2a1e06b3`.
The **entire** inherited **564972-byte input / 2068804-byte output** prefixes
retain respectively `b778a5b247549f69da5dc81b3e29abb331c2175a5cdb0eb4940df3b0a6122975`
and `ab33212cfdbf47c158a014c6c3d4ba25113e3c2f9f6d0234723f793a5c09146d`.
Literal C whitespace is preserved. Every new section requires status **0**
after dependency/NYI/stub/sanitizer rejection, not just equal failed runs.

Six real linker controls wrap the direct actual-table entry adapter (status
**0/134**, each default/update), then original **asfail** at inherited expected-
abort section **27** (default/update). All reject with harness **2** at the
marker check, **before** status acceptance or recording replacement; the
recording is unchanged in all six. An initial attempt to wrap v_hinder_med
could not interpose its same-translation-unit table reference and did not
trigger a marker; it was not counted as a successful negative control.
Probe sources/launchers/captures are deleted; final normal comparison passes.

Logs: `tmp/issue181-{source-audit,final-integrity}.log`,
`tmp/issue181-command-{update,default}-final.log`,
`tmp/issue181-negative-{drv_use_callback,asfail}-*.log`.
Bounded vet/fresh full-enabled tests, all fifteen defaults, strict full turn
and all controls pass. Five of six unchanged phase fixtures now pass C/C/Go;
only phase-month-decay is newly passing. phase-day0 still reaches separately
filed #182's drum callback. [Full fixture evidence](../turn-compat#hinder-meditation-fixture-rerun-181).
The final 201-world/four-corpus tier remains #178, not rerun here.

## GUARD and STACK command prerequisites (#180)

**0.113.0-alpha** adds complete original combat.c `v_guard` **2947–2964**,
stack.c `v_stack` **703–741**, and u.c `check_char_gone` **933–962**, wiring
only the two original command starts in glob.go. Existing charWhere/charHere/
contacted, admission and joinStack remain the authoritative owners, unchanged.
Validation retains garrison/character/visibility/gone, then self/cycle/prisoner/
admission order. The active `char_gone` macro tests nonzero moving, even on
the same epoch, not elapsed days. GUARD normalizes any nonzero raw argument
to TRUE; normal parsing still clamps negative arguments to zero. No added
prisoner filter, guard clearing, combat, UNSTACK, USE or independent fix.

The real C table calls the complete extracted definitions unchanged; there
are no callback overlays. `callback_start [raw-a]` calls that actual table
entry directly, optionally setting its first argument to check signed/raw
boundaries independently of the real parser. Real load/do/start/evening
paths separately exercise queued commands, completion, counters and outputs.
All **100** archived source files equal the immutable pin, and **577** complete
extractor-listed definitions occur literally unchanged exactly once. New
complete-definition SHA-256 values:

- `v_guard`: `0a51deac1d45d9e77ce8f309d4b08acc410b376de40a5ae12642d0acb38c75ec`
- `v_stack`: `10a0a105fd5536cee532596d6d1757c5e01a91062f404a369ccee59e02b4deec`
- `check_char_gone`: `d54ff88d3fb7c9522e82ed7826658dc419ea1caccf6e39c38f4ddb782592ffdf`

No new upstream source dependency, affected.sh change or pin move is needed.
Full source/closure audit: `tmp/issue180-source-audit.log`.

**18 appended asymmetric sections** check raw negative/256/zero flags, omitted
arguments, prisoner GUARD and immediate suppression; missing garrison, zero,
noncharacter and missing targets; self, cycle before prisoner, prisoner before
admission, absent before moving/hidden, same-epoch/negative moving; hidden and
fog visibility, character versus faction contacts, real garrison visibility,
same-faction/same-stack exemptions, building depth, admission inclusion and
exclusion, both refusal recipients, nested reparenting/children/prisoner
insertion, retained guard, real command/phase dumps and saved state. Legacy
next draw **62497** and entity aftermath **701/2/1/1** match the independent
no-command control. Every new C section requires status **0**, after marker
rejection, not merely equal failed recordings.

Unchanged C twice and Go match **1350 sections / 103698 output lines /
2068804 bytes**, SHA-256
`ab33212cfdbf47c158a014c6c3d4ba25113e3c2f9f6d0234723f793a5c09146d`.
Input: **29695 lines / 564972 bytes**, SHA-256
`b778a5b247549f69da5dc81b3e29abb331c2175a5cdb0eb4940df3b0a6122975`.
The entire inherited **559758-byte input / 2035552-byte output** prefixes
retain respectively `35c5e55721fe82b4481ddc2dcb2e1414f82041937c8556fd3801e85830b91615`
and `57ad48790379b5df3ef83e155320c0ef6c6d78026477c29e81c22626f98d28f0`.
Literal C whitespace is retained. Byte audit: `tmp/issue180-final-recording.log`.

Four temporary actual linker probes wrap only v_guard at the first new
section (**1332**), emit a dependency marker on stdout/status **0** or
stderr/status **134**, in both default/update modes. All reject with harness
status **2** before status acceptance or recording replacement; the recording
hash stays unchanged during the probes. Probes/launchers/captures are removed;
the final real command default passes without them. Logs:
`tmp/issue180-negative-{0,134}-{default,update}.log`,
`tmp/issue180-command-{update,default}.log`.

Bounded vet, fresh bounded full-enabled tests and all **15** engine/library
defaults pass. Strict ordinary turn and all **13 lifecycle / 10 command /
15 USE / 23 phase / 13 export** controls pass, including four nonempty-export
flag-free aftermath controls. The unchanged phase-daily script now passes
C/C/Go, including all **245 cwd entries**. phase-day0 passes C twice but Go
now stops later, at **USE item key 7 / savage.c v_use_drum**, not GUARD.
That separately sized prerequisite must be filed by the parent before any
expansion; no drum handler or #181 hinder-meditation work is included.
See [strict fixture evidence](../turn-compat#guardstack-fixture-rerun-180).
No all-six, full-corpus, native-platform or stage-exit claim is made.

## Complete location export (#176)

**0.112.0-alpha** completes original `export_location` **133–189** in the
existing state_export.go. The temporary location-boundary test is replaced
by real C/C/Go contract sections in `TestCommandMatchesC`. Main now accepts
the path and calls the root writer after save_db/save_logdir and before
Times, at the original position. Other mode gates and missing USE/gameplay
callbacks stay NYI. No G3 #114 format expansion or invented visibility filter.

The wrapper still includes **all unchanged state_export.c** and exposes one
additional static entry point, not a replacement implementation. All **100**
archived files equal the pin and all **574** extracted complete definitions
occur literally unchanged exactly once. Source audit includes all real
dir.c exit constructors/cache, loc.c all_char_here/add_char_here/province/
subloc, u.c loc_depth/loc_hidden/test_known, swear.c lord/player, original
macros and glob.c tables. No new dependency or affected.sh/pin change.
Whole state_export.c remains **6496 bytes**, SHA-256
`2f710e5ec96c810996206c87efc7908fe50f75732b28f7119237e9ef1fb8210b`;
complete export_location SHA-256
`e110d3b2d8868e0e2d80813e068170c1b4b2d891fc8cf1c5bb669f398f66e2e5`.
Full per-definition ranges/hashes: `tmp/issue176-source-audit.log`.

Three appended asymmetric sections cover actual nonempty location/faction/
root bytes, lowest-ID viewer versus a higher-ID listed unit with a different
lord/knowledge, undiscovered exits versus known hidden destinations, unsorted
roads and duplicated directions/destinations, port insertion, distance,
barriers, collapsed mines, cloud links, open links, hidden outward exits,
province/subloc/building/ship names and depths. Original loc_depth_s has
**NULL at build depth**, so a building's depth is **empty**, not an invented
"build" label; `<no depth` also retains the original missing bracket.
Ships have no exported exits even though the shared constructor can build
them. `all_char_here` includes foreign units, prisoners and nested character
stacks, but prunes noncharacter branches and excludes units inside buildings
or ships from the outer location. No extra visibility/prisoner filter.

Exit order is the real constructor's order, never sorted. C owns/reuses its
static exit scratch and frees old entries on the next construction; export
borrows it immediately and never sorts/reclaims it. Go uses the existing
fresh scratch constructor without retaining references. Only copied here IDs
sort. Raw/saved/allocation/list snapshots before and after export are exactly
identical: **3050 / 1647 / 795 bytes**, SHA-256 respectively
`043c619e22245dae42fb1432c93654c3d26d583095777af55ca12efe80cabdf3`,
`f640a8a19a5a26b9bceaefbfcd68dc9d944d9321aacb9da52c820a03ab19c1ae`,
`efbc06099e4cfd80225d30174215a2a412c260ca221c49acd02f3a8b6441a650`.
Explicit later knowledge changes are not mislabeled read-only. Legacy next
draw **62497** and entity draws **701/2/1/1** match the inherited independent
no-export control; a borrowed buffered output stream survives location/root
export and directory rename without reopening or losing its inode.
Evidence: `tmp/issue176-bytes-aftermath.log` and complete raw recording.

Unchanged real C twice and Go match **1332 sections / 101686 lines /
2035552 output bytes**, SHA-256
`57ad48790379b5df3ef83e155320c0ef6c6d78026477c29e81c22626f98d28f0`.
Input: **29344 lines / 559758 bytes**, SHA-256
`35c5e55721fe82b4481ddc2dcb2e1414f82041937c8556fd3801e85830b91615`.
The entire inherited **556686-byte input / 2001138-byte output** prefixes
retain their accepted hashes `cf2cfca1a439847fdb69d72003a9bc992071b7dc1bd2fc226f5d41857e4136a3`
and `22c85b5fb8598701ecd92b7fd304f8f01f4b62f7143c53fcd0388be4765bb540`.
Literal trailing spaces in C save/output bytes remain deliberately unchanged.

Four actual linker probes wrap the new direct location entry point: marker
on stdout/status **0** and stderr/status **134**, both default and update.
All reject at section **1329** with harness status **2**, before expected
statuses or recording replacement; the complete output hash stays unchanged.
Probe sources, compiler launcher and probe captures are removed. Logs:
`tmp/issue176-negative-{0,134}-{default,update}.log`. The normal command
default passes again after their removal (`tmp/issue176-command-default.log`).

Bounded vet/fresh full-enabled tests (**11 packages**, Olympia **1.934s**) and
all **15 engine/library defaults** pass. Strict ordinary turn and all **13
lifecycle / 10 command / 15 USE / 23 phase** controls pass, plus **13** new
nonempty export controls and four flag-free aftermath comparisons; see
[turn evidence](../turn-compat#state-export-save-lifecycle-176).
All six unchanged phase scripts ran real C twice then Go: **three pass and
three expose separate unported gameplay prerequisites**, not export byte
differences. Spec is **32 pass / 93 fail / 0 error**, retaining all 23 prior
passes. Exact fixture outcomes and remaining boundaries are documented there
and in [spec](../olympia-spec#location-export-integration-176). No six-fixture
or stage-exit success is claimed; broader combinations remain #177.

## State export byte strings, units and faction roots (#175)

**0.111.0-alpha** adds `internal/olympia/state_export.go`, corresponding to
the complete original `json_str` (49–67), `sort_ids` (70–75), `export_unit`
(78–130), `export_faction` (192–253) and `write_state_export` (256–290).
`exportLocation` is an explicit NYI boundary for #176, with a Go test proving
that a faction with a location stops there. **The CLI gate and main integration
are untouched.** No six-phase-fixture success, new format or G3 #114 expansion
is claimed. There is no gameplay implementation or successful dependency substitute.

The C wrapper includes the **entire unchanged original state_export.c**, then
exposes its actual static functions. It does not extract/rewrite export bodies,
replace `export_location`, or alter any callback. Direct unit cases may have a
real location/ship; faction/root cases have actual `subloc == 0` for every
listed unit and therefore never call the real location exporter. All required
real dependencies were already compiled by this oracle. The only new source
dependency is `olympia/state_export.c`, now tracked in affected.sh and the tools
dependency notes. No tool pin changes.

The whole file is **6496 bytes**, SHA-256
`2f710e5ec96c810996206c87efc7908fe50f75732b28f7119237e9ef1fb8210b`.
All **100** original archived command-tool files equal the immutable pin.
All **574** extracted full definitions occur literally unchanged exactly once,
including the complete subloc/stack_parent/has_item/int_comp/lord/player
dependencies. The full original code/glob/z/sout/io and lib list/RNG modules
remain directly compiled. Audit: `tmp/issue175-source-audit.log` (individual
state-export definition ranges/hashes, declarations and callers included).
The declarations remain oly.h:1286–1287/proto.h:275–276; main's only assignment
is 835–840 and its call is 1036, after save_db/save_logdir. That integration
belongs to #176, not this primitive closure.

Nine appended sections cover NULL/empty strings, NUL at the start/in the
middle, **every nonzero byte 0x01–0xff**, quote/backslash, all controls and DEL.
UTF-8-looking bytes remain individual lowercase Latin-1 escapes, not decoded
runes. They also cover empty/nonempty regular faction roots, absent player
records, system/silent exclusion, signed-short turn endpoints, repeated writes,
NULL paths and nonfatal empty/ENOENT/EISDIR destinations. Exact raw JSON bytes
are compared, not parsed/normalized JSON.

Unit tests retain insertion-order inventories/skills while sorting temporary
IDs, including duplicate items (has_item returns the **first** quantity),
zero/negative quantities, invalid inventory references, known versus learning/
unknown skills, and a known valid non-skill ID (no extra kind filter). Direct
sort_ids tests mutate the supplied original list, including empty/singleton/
duplicate lists. Faction tests preserve unsorted original unit lists and
duplicates: C walks the actual player units list, **not** a fresh lord scan,
and adds no prisoner or character-kind filter. Thus a listed prisoner,
wrong-lord unit and unformed entry appear; an unlisted owned character does
not. Stacked_under follows the immediate character parent, while where walks
through that parent to the real ship (or zero in root tests).

The C “read-only” comment is not used as evidence. Before/after observations
compare **all saved boxes plus raw inventory/skill/unit/here/output lists,
record allocation presence, temp and output-order fields** byte for byte.
The observer restores only save_box's own temp mark; it never resets exporter
effects. Snapshot bytes/SHA-256:

- unit/ship case: **1712**,
  `e60e1c55dcf832e44ed8222a8821f4c55089f941303926b998ed2f476d15983b`;
- faction/root case: **1381**,
  `35992be9e3c1f4160aa597c78ea70cb224400a03542284e8bcf4b07f28b5b39a`.

Absent player/character/magic records stay absent; no new entity IDs or lazy
backlinks are created. The list source actually allocates/reclaims temporary
copies using its ordinary heap allocator, not entity allocations; Go keeps
the corresponding temporary copies, without reordering persistent lists.
This is persistent-state/allocation-ownership evidence, not a cross-language
total heap-allocation count. Independent no-export controls preserve the next
legacy draw **62497**, then all four keyed entity draws: **701/2/1/1** at
turn zero and **103/1/3/2** at turn 32767. A borrowed buffered log stream
survives export and directory rename: both writes remain in the old **46-byte**
inode, the replacement path stays absent until explicit close, and the later
reopen produces a **23-byte** file. Full evidence:
`tmp/issue175-bytes-aftermath.log`.

Unchanged real C twice and Go match **1329 sections / 100457 lines /
2001138 output bytes**, SHA-256
`22c85b5fb8598701ecd92b7fd304f8f01f4b62f7143c53fcd0388be4765bb540`.
Input is **29184 lines / 556686 bytes**, SHA-256
`cf2cfca1a439847fdb69d72003a9bc992071b7dc1bd2fc226f5d41857e4136a3`.
The complete accepted **1320-section / 99685-line** baseline remains a byte
prefix: input **552971** bytes
(`deeb7928b106cd8bcdfa36102516dfe087f59c735b51290920e85a68d054ffca`),
output **1983914** bytes
(`7d7315f4335b7cf7ffa5c1fc710846204a9b4863d2d203974fceef131a8593e4`).
Literal trailing spaces/high bytes from the original writers are retained.

Six temporary linker controls prove marker rejection before status acceptance
or recording replacement: root write section **1323**, status **0/134**, each
default/update; enabled negative-index expected-assertion section **1284**,
status **134**, default/update. Both streams are exercised. Every run rejects
with harness status **2**, preserving the output hash above. Probe sources,
compiler wrapper and matching probe captures are removed. Logs:
`tmp/issue175-negative-{root-{0,134},failure-134}-{default,update}.log`.

Validation: bounded vet and fresh bounded full-enabled tests pass (Olympia
**1.877s**), as do all **15** engine/library defaults, including a clean normal
command run after probe removal. Strict full ordinary turn (**233 entries**),
all **13 lifecycle / 10 command / 15 USE / 23 phase** controls and flag-free
aftermath comparisons pass raw status/streams and complete cwd types/modes/
bytes. The read-only `tmp/issue179-roundtrip` C cache is independently checked
against **1212** original archived files and accepted binary SHA-256
`151c7355890eba7f8bfda1f4f16f4a490947aaa1db1ffcc3ca859cd659c01158`.
Spec retains **23/102/0**; initialization exits **2**, golden **NO**.
All tools run sequentially with explicit isolated G3_REPO and bounded engines;
the unchanged lists C driver uses a temporary bounded compiler launcher, then
removes it. Logs: `tmp/issue175-command-{update,default}.log`,
`tmp/issue175-{vet,tests,cache-audit,turn-full,spec}.log`,
`tmp/issue175-*-default.log`; strict summaries: `tmp/issue175-turn/`.
No full roundtrip/four-corpus repeat, peripheral genesis/mapgen/island/scanf,
full C spec, native Darwin/Windows, export CLI or six-phase-fixture success
is claimed. No new prerequisite or upstream blocker was found. Changes are
uncommitted; no push, issue edit or checklist tick.

## Monthly phase hooks and CLI integration (#174)

**0.110.0-alpha** adds the remaining **84 hook sites** (32 ran / 52 acted)
in the existing monthly owners: **43 monthly event rows**, with no new
gameplay. Together with #173, every original caller is present: **135 sites,
50 ran / 85 acted, in 52 complete original hook-containing definitions**.
The actual `process_orders` and `daily_events` scheduling definitions are
also extracted unchanged. `--phase-coverage` now reaches the existing world
field and main writes it after command/USE coverage, before check_db/save_db.
State export and the later unsupported modes remain gated.

The complete audit maps **each of the 135 individual C file:line callers to
its actual Go file:line and method**, including all repeated branch positions:
`tmp/issue174-source-audit.log`. All **100** original archived command-tool
files equal the immutable pin; every one of those 52 full definitions occurs
unchanged exactly once in the compiled extraction. Complete coverage.c remains
**12,597 bytes**, SHA-256
`5730dc2b64b817a6ee1a1d9c250856a400f6aaf54b6ac50251b5f37e5e5ded5a`.
There is no source dependency/pin change, C body edit, phase-macro substitution,
gameplay callback overlay, new RNG adapter or successful missing dependency.

Monthly source ownership (inclusive definition ranges, not broad task ranges):

| Pinned C definition(s) | Existing Go owner, under internal/olympia |
|---|---|
| day.c add_claim_gold 635–653; add_noble_points 655–676; add_unformed 699–719 | day.go addClaimGold/addNoblePoints/addUnformed |
| day.c increment_current_aura 722–774 | day_monthly.go incrementCurrentAura |
| day.c decrement_ability_shroud 806–825; decrement_loc_barrier 828–853; decrement_region_shroud 856–876; decrement_meditation_hinder 879–898 | day_monthly.go decrementAbilityShroud/decrementLocBarrier/decrementRegionShroud/decrementMeditationHinder |
| day.c loyalty_decay 1007–1075; relic_decay 1078–1108; pillage_decay 1111–1130 | day_monthly.go loyaltyDecay/relicDecay/pillageDecay |
| day.c auto_drop 1133–1167; charge_maint_costs 1352–1376; animal_deaths 1379–1418 | day.go autoDrop/chargeMaintCosts/animalDeaths |
| day.c inn_income 1421–1489; temple_income 1492–1522 | day_income.go innIncome/templeIncome |
| day.c collapsed_mine_decay 1525–1549; ghost_warrior_decay 1556–1584; corpse_decay 1591–1619; dead_body_rot 1626–1652 | day_decay.go collapsedMineDecay/ghostWarriorDecay/corpseDecay/deadBodyRot |
| day.c storm_decay 1655–1678; storm_move 1705–1734; quest_decay 1737–1753; collect_taxes 1756–1789 | day.go stormDecay/stormMove/questDecay/collectTaxes |
| day.c hide_mage_decay 1945–1963 | day_monthly.go hideMageDecay |
| day.c link_decay 1966–1997 | day.go linkDecay |
| day.c move_city_gold 2005–2024 | day_production.go moveCityGold |
| day.c special_locs_open 2046–2085; special_locs_close 2088–2127; post_month 2147–2203 | day.go specialLocsOpen/specialLocsClose/postMonth |
| garr.c garrison_gold 446–549; determine_noble_ranks 718–761 | garr.go garrisonGold/determineNobleRanks |

The inherited 51 sites were re-audited against their actual Go locations,
not just #173's intended list: art.go, buy.go, day.go, day_damage.go,
day_entities.go, faery.go, garr.go, hades.go, input.go, move.go, npc.go,
savage.go and storm.go. Their complete pinned definitions/ranges are in the
#173 table below and the new all-caller map. Token reconciliation still counts
twice per turn in its day0-labeled rows. Production hooks occur **only in
postMonth immediately before postProduction**, never its load initialization.

The **12 appended sections** exercise every monthly event with nonzero actual
gameplay, verified separately in `tmp/issue174-hook-coverage.log` without
including primitive sweeps or seeded counters. They distinguish seasonal
else-if precedence/all four branches; aura caps, positive/negative shrouds,
barrier decrement+dissipation and pillage wrap; summon release versus contract
and fear desertion/decay multiplicity; relic decrement/vanish/short wrap;
per-building income/traveller and temple zero/nonzero offerings; maintenance
paid/starved/no-debt; animal **inventory entries with deaths, not individuals**;
ghosts/corpses/body and deletion-safe collapsed mines; storm weakening,
dissipation and relocation; actual repeated post_month tax/castle/owner,
NP/claim/unformed/rank/production/autoquit/link/quest composition; positive
versus zero city transfers and attempted zero-tax collection. Empty calls,
immediate gameplay with suppressed counters, repeated grouped rows, and a
real maintenance call wrapping signed int32 counters are included. Original
saved boxes, queues, allocation/mint effects, scalars and logs remain observed.

Unchanged real C twice and Go match **1320 sections / 99685 lines /
1983914 output bytes**, SHA-256
`7d7315f4335b7cf7ffa5c1fc710846204a9b4863d2d203974fceef131a8593e4`.
Input: **29025 lines / 552971 bytes**, SHA-256
`deeb7928b106cd8bcdfa36102516dfe087f59c735b51290920e85a68d054ffca`.
The complete accepted #173 prefixes remain exact: **544674 input bytes**
(`c0d27f47c17ef533051667df155fa4b96c6377b1d87da67c2614cc8932214572`)
and **1902639 output bytes**
(`a93ea34d069efd3a47ada8bf71f4ee75bfe0b9636335b589793b6dfcbe15c90e`).
Native traversal also matches raw C/C/Go with unchanged shared input.

Four temporary linker probes emit a dependency marker for the enabled
negative-index expected-assertion section **1284**, exiting **0** or **134**.
Both normal and `--update` runs reject with harness status **2** before
expected-status checks or recording replacement; the recording hash above
stays unchanged. Probe/wrapper sources and matching probe captures are removed.
Logs: `tmp/issue174-negative-{0,134}-{default,update}.log`.

Validation logs: `tmp/issue174-command-{update,default}.log`,
`tmp/issue174-{vet,tests,cache-audit,turn-full,spec}.log` and
`tmp/issue174-*-default.log`. All 15 defaults, bounded vet/fresh bounded tests,
strict ordinary turn, 13 lifecycle/10 command/15 USE/**23 phase** controls and
spec no-loss check pass. See [turn evidence](../turn-compat#monthly-phase-cli-controls-174).
No full 201-world/corpus rerun, native Darwin/Windows, golden initialization,
six phase-fixture success or state-export integration is claimed. Export is
still NYI, as required; unchanged phase fixture success belongs to #176/#177.

## Phase primitives and real day-zero/daily hooks (#173)

At **0.109.0-alpha**, coverage.go adds the complete **74-row** PH enum,
fixed row metadata, world-owned int32 ran/acted arrays and exact TSV writer.
The phase CLI gate and main's phase writer remain unwired until #174;
state export remains gated. No monthly-only hook is ported here. Shared
`checkTokenUnits` counts wherever its existing real caller invokes it, including
the month-end call; its rows remain labeled `day0`, as in C. No new gameplay,
phase callback replacement, RNG adapter, upstream dependency or pin move.

The four phase-macro substitutions are removed from the main extracted prefix
and the separate Hades/Faery/NPC translation units. The driver links the
**entire unchanged coverage.c**, not a phase model or edited function body:
**12,597 bytes**, SHA-256
`5730dc2b64b817a6ee1a1d9c250856a400f6aaf54b6ac50251b5f37e5e5ded5a`.
All **100** archived original files match the immutable bare-repository pin.
The source audit checks these **22 complete definitions** unchanged exactly
once, including the two actual scheduling owners:

| Original C source | Complete definitions / inclusive lines |
|---|---|
| art.c | check_token_units 1200–1262 |
| buy.c | match_all_trades 698–729 |
| day.c | ship_coastal_damage 261–327; random_loc_damage 454–486; dogs_bark_at_hidden_chars 489–516; heal_char_sup 519–578; default_garrison_pay 581–613; heal_characters 616–632; increment_stone_ring_aura 777–803; noncreator_curse_erode 901–963; daily_events 2206–2271 |
| faery.c / hades.c | auto_faery 607–634 / auto_hades 499–526 |
| garr.c / input.c | ping_garrisons 1064–1093 / process_interrupted_units 1021–1052; process_orders 1057–1102 |
| move.c / npc.c | check_ocean_chars 852–898 / queue_npc_orders 613–691 |
| savage.c / stack.c | init_savage_attacks 353–409 / weekly_prisoner_escape_check 345–374 |
| storm.c | create_some_storms 1525–1553; natural_weather 1556–1613 |

These contain **51 hook sites**: 18 ran calls and 33 acted calls, for **31
day0/day1/daily event rows** across 18 routines. Each has nonzero actual-C
gameplay evidence in the appended cases. Group rows reuse the first contiguous
routine's ran index, even if direct accounting increments a later row's ran.
Trade actions count reduced buyers, not settlements; token melt/add, hidden
dog checks, and NPC creation/unit attempts count after the original call even
when it changes nothing. C assertion locations and branch multiplicity remain.

The **25 appended sections** include 18 real-gameplay sections, one direct
accounting section and six expected assertions (status **134**). Direct tests
compare all original C enum names/evaluated values (generated from oly.h),
all zero and asymmetric rows, every contiguous group, int32 overflow, repeats,
NULL versus empty paths, suppression before invalid-index assertions, immediate
writing and ENOENT/ENOTDIR/EISDIR/empty-path failures. Assertions retain
coverage.c **405/417**, token owner **art.c:1218**, and real gen_item failure.

Real cases exercise both stop branches with retained successors; multi-seller
settlement; token successful creation and zero-effect attempts; nine/ten guards,
zero/positive payroll and immediate gameplay; healing/recovery/sick deaths;
reef/rock damage and both omen branches; mine thresholds and repeated inn
triggers; building/outdoor prisoner escapes; ring cap, curse inventory-entry
multiplicity and hidden dogs without hounds; natural weather; repeated Faery
and Hades population/unit passes; every independent NPC dispatch branch; savage
spawn before Hades mint; ocean plunge/cleanup; and the **actual process_orders
day-zero/day-one/daily pipeline**. Saved boxes, real logs, retained orders,
mint IDs and independent RNG probes are compared as applicable. No enabled
case calls post_month or claims complete monthly phase parity. Schedule-only
callback overlays remain scheduling evidence, not gameplay evidence.

Real C twice and Go match **1308 sections / 96599 lines / 1902639 output
bytes**, SHA-256
`a93ea34d069efd3a47ada8bf71f4ee75bfe0b9636335b589793b6dfcbe15c90e`.
Input: **28560 lines / 544674 bytes**, SHA-256
`c0d27f47c17ef533051667df155fa4b96c6377b1d87da67c2614cc8932214572`.
All inherited **1283 sections / 530718 input bytes** (SHA-256
`3184d7795bd11f7d527f2c4e5ac7212fa4d379798015f5c41c2398ae2244127a`)
and **1769914 output bytes** (SHA-256
`5a3d7d10790885c7045e4d52f46e3f7754341c1fb63a6273a66dac5190b303dc`)
remain exact prefixes, including original C whitespace. Native traversal also
matches raw C/C/Go. No expectation is computed from Go or modeled totals.

Negative linker probes wrap phase_ran only for the enabled negative-index
input in expected-assertion section **1284**, emitting a dependency marker and
exiting either **0** or **134**. Normal and `--update` modes reject all four
with harness status **2**, before expected-status acceptance or any recording
replacement. The recording checksum is unchanged. Temporary probe sources and
compiler wrapper are removed before the final ordinary oracle run.
Logs: `tmp/issue173-negative-{0,134}-{default,update}.log`,
`tmp/issue173-command-{update,default}.log`, `tmp/issue173-source-audit.log`
and `tmp/issue173-hook-coverage.log`.

All 15 engine/library defaults, bounded vet and fresh bounded tests pass.
Strict ordinary full turn, 13 lifecycle, 10 command and 15 USE controls pass
at the same pin; no phase CLI is enabled. The 1212 original cached C files
and accepted C executable hash are verified before read-only cache reuse;
engines run only on actual disposable copies. Spec preserves **23/102/0**;
golden initialization remains unsupported, not passing. No full-corpus or
roundtrip rerun is claimed here; the inherited #179 pin gate remains evidence.
All compat invocations explicitly set `G3_REPO="$PWD/tmp/upstream.git"`;
engines/tests use shared run_limited, with tools run sequentially. Logs:
`tmp/issue173-{vet,tests,cache-audit,turn-full,spec}.log`,
`tmp/issue173-*-default.log`; strict summaries in `tmp/issue173-turn/`.
Phase fixture success remains for #174/#176, not this partial integration.

## Follow-C load validation and summary traversal (#179)

At **0.108.1-alpha**, the affected pin moves to merged G3 #360. The complete
original io.c and report.c definitions are still used unchanged by C. Go's
summary loop now checks its NULL sentinel before calling prevWord; its old
comment preserving the nonportable comparison is obsolete and removed.
One appended section checks the exact four new upstream summary names:
unbroken 20-byte name, repeated stop words, word break and short name.
Existing broader name boundaries and stop-word-table cases remain unchanged.
Loaded-trade boundary coverage extends [io-compat](../io-compat), not a second
loader harness.

Real C twice and Go pass **1283 sections / 90478 lines / 1769914 output
bytes**, SHA-256 `5a3d7d10790885c7045e4d52f46e3f7754341c1fb63a6273a66dac5190b303dc`.
Input: **27735 lines / 530718 bytes**, SHA-256
`3184d7795bd11f7d527f2c4e5ac7212fa4d379798015f5c41c2398ae2244127a`.
ALL inherited **530432 input bytes** and **1769491 output bytes** remain
exact prefixes, retaining their #172 hashes below. The only added output is
the 17-line successful summary section. Live native traversal also still
matches raw C/C/Go. No unrelated output, source dependency or oracle changes.
Logs: `tmp/issue179-command-{update,default}.log`,
`tmp/issue179-recording-audit.log`; full new-fixture/whole-world evidence is
in [the pin gate](../README.md#follow-c-pin-gate-179).

## USE accounting and real shared dispatch (#172)

At **0.108.0-alpha**, coverage.go completes the USE counters and writer;
use.go wires the real shared `v_use`, `d_use`, `i_use`, `v_use_item` and
`d_use_item` through cmdTbl. The main writer follows command coverage and
precedes check/save. Only the USE CLI gate is removed; phase/export remain
explicitly gated. Existing metadata/helpers and hide/bind-storm callbacks are
reused. **Every other non-NULL skill callback remains explicit NYI**, including
implicit/shipbuild/SAIL. All item gameplay callbacks remain NYI before return
accounting or fabricated outcomes. Go-only boundary tests reach every absent
skill start through real command dispatch and every item key through vUse,
checking pre-callback mutation and absence of either fabricated counter family.

The existing C extractor now retains the entire unchanged **6,122-byte**
coverage.c prefix through the USE writer (before phase coverage), SHA-256
`827715623143788543f71c4c37f76590f4f500602ec095849dc27421d4d150ce`.
Complete original use.c definitions are extracted unchanged exactly once:
`v_use` **524–653**, `d_use` **685–762**, `i_use` **765–782**,
`v_use_item` **785–920**, `d_use_item` **923–954**. Missing references are
abort-only rejected dependencies, never successful gameplay substitutions.
The source audit compares all **100** archived original files against the
immutable pin; coverage.c is **12,597 bytes**, SHA-256
`5730dc2b64b817a6ee1a1d9c250856a400f6aaf54b6ac50251b5f37e5e5ded5a`.
No upstream source dependency or pin changes; use.c/coverage.c were tracked.

Twenty-two appended sections cover direct outcomes, placeholder/sentinel/high
indices and keys, signed overflow, NULL/empty paths, immediate suppression,
lazy skill-only allocation, numeric skill ordering and real database names,
zero rows, repeated/truncating writes and ENOENT/ENOTDIR/EISDIR opens. Real
shared paths cover invalid/wrong-kind/school/knowledge checks, no table entry,
all nine magical-item safe-haven keys, pre-shift magic rejection, post-shift
requirements, hide success/failure/finish/interrupt, bind-storm scroll correction
and speedup/consumption, and the original NULL-start production branch.
Restored `use_ent=0` is looked up locally, not written back. Failed finishes
increment use_count before callbacks but never start counters. A hide scroll
passes shared access yet fails hide's own knowledge check: it counts **fail**,
not checks, and is not consumed. Successful production consumes requirements
and grants experience in original order; a missing product asserts, not NYI
or success. Item finish lost/zero/invalid-key paths invoke real C, not stubs.

The previous **1,260 sections / 523,069 input bytes** and **1,703,330 output
bytes** remain exact prefixes (SHA-256 respectively
`c0099469c48a560b66f8dbcb4e3d6131b9b94ed182011e15a5397ef4a521f0e5`
and `f77a1f488eb9ae4482777766d79a4ed2e4ebd36d4d574d29d4825e485ea25975`).
Real C is run twice before recording/comparing Go; marker rejection precedes
expected-status acceptance and `--update`, including expected assertion exits.
See also [the actual CLI controls](../turn-compat#use-coverage-cli-controls-172).

Accepted #172 recording at the same immutable pin: **1,282 sections / 27,725
input lines / 530,432 input bytes**, SHA-256
`9466b4111cb9732208571b30bbe7bbab5adafedf4ff9e87b12759f6cbe90d207`;
**90,461 output lines / 1,769,491 output bytes**, SHA-256
`927a74512cd4171cc161874a886601e69dcd2348d38cc269249981c681aca305`.
The 22 additions contain 15 successful and 7 expected-assertion sections.
`tmp/issue172-source-audit.log` records byte identity of 100 archived source
files, the complete unchanged 6,122-byte coverage prefix (SHA-256
`827715623143788543f71c4c37f76590f4f500602ec095849dc27421d4d150ce`),
and each of the five unchanged shared USE definitions extracted exactly once.

Negative acceptance probes wrapped `use_coverage_skill` only for the enabled
negative-outcome input in expected-assertion section 1264: a missing-dependency
marker followed by either exit 0 or exit 134 was rejected with harness exit 2,
both ordinarily and with `--update`, before status acceptance or recording
writes. The recording checksum stayed unchanged. Logs are
`tmp/issue172-negative-{0,134}-{default,update}.log`; the temporary wrappers
were removed and the ordinary unwrapped oracle passed again in
`tmp/issue172-command-default.log`.

Independent parent review of the pinned coverage/shared USE definitions and
complete diff passes. Fresh bounded vet/tests, the eight relevant default
comparisons and all 92 CLI vectors pass. The parent full repeat also passes
the strict 233-entry world turn, 10 command controls, 15 USE controls and all
13 lifecycle controls, including nonzero shared checks/HIDE and flag-free
aftermath. Both inherited recording prefixes remain independently byte-exact.
Logs: `tmp/parent172-{vet,tests,cli,turn-full}.log` and
`tmp/parent172-*-default.log`; summaries in `tmp/parent172-turn/`.
During acceptance G3 #360 merged at `f73e48118805ffe8ff5e763bfa4f20ae00c3a705`.
Follow-C #179 takes priority before the next feature task; this task's
acceptance remains at its original immutable pin, not an implicit repin.

## USE metadata and shared prerequisites (#171)

At **0.107.0-alpha**, existing use.go owns the complete original `use_tbl`
**84–229**: **135 skill rows**, initial placeholder at **0**, separate zero
sentinel at **136**, exact allow/time/poll and callback NULL bits. The older
stage prose's 88–237 range includes part of `v_implicit`; actual table data
ends at 229. `v_implicit` and `v_shipbuild` remain NYI, as do all other absent
non-NULL gameplay callbacks. Only already implemented `v_hide`/`d_hide` and
`v_bind_storm`/`d_bind_storm` are reused. Full SAIL is NOT implemented, so its
three callbacks remain explicit NYI, not NULLs. A Go boundary test checks all
non-NULL callbacks, the four reused function identities and all three still
gated command-level USE callbacks. `--use-coverage` stays gated too.

Complete unchanged pinned definitions newly extracted by this existing oracle:

| use.c definition | Lines |
|---|---|
| `find_use_entry` | 251–261 |
| `use_tbl_skill` | 269–274 |
| `may_use_skill` | 287–323 |
| `magically_speed_casting` | 326–359 |
| `correct_use_item` | 369–385 |
| `meets_requirements` | 388–433 |
| `consume_requirements` | 436–482 |
| `consume_scroll` | 485–499 |
| `experience_use_speedup` | 502–521 |
| `add_skill_experience` | 660–682 |

The C oracle additionally retains the **18,205-byte unchanged original struct,
callback declarations and table (9–229)**, SHA-256
`0cc15e0524f80c711fb50f47243859b3fb73449a1b7e78653f1f1d1e74e222c3`.
Full use.c is **53,389 bytes / 2,134 lines**, SHA-256
`7582e5f81f26e6fe43c47c1e3444b4d4a2487edebb2dcb8908600b2a6a7bcd06`.
Every extracted definition is byte-equal to the complete pinned definition,
with original `#line` locations; all **100** archived lib/olympia regular
sources equal the immutable archive. Symbol/callback-name/order checks
supplement execution of all original table rows and numeric IDs. Absent C
gameplay references retain abort-only dependency markers; no helper calls
them, no successful callback overlay is used and no table row is replaced.

Dependencies read/verified include complete shared USE callers, STUDY's
scroll caller, coverage.c's sentinel walker, existing `rp_skill_ent`,
`has_skill`, `exp_level`, `skill_school`, oly.h's actual field widths,
REQ constants and magic/quick-cast/use-key macros, loop.h's copied/filtered
inventory traversal, u.c `has_item`, `consume_item`/`sub_item` and
`destroy_unique_item`, and existing item/skill/list allocation, output,
deletion, serialization and shared mint owners. No new upstream source
dependency or pin: all of these files were already tracked by affected.sh.

Twelve appended sections exercise original metadata/lookup/sentinel; last
positive inventory artifact before last scroll (opposing ID/list order),
known-before-item grants and filtered invalid/zero inventory entries; first
grant correction versus a real item use key, including non-scroll grants;
absent/empty record nonallocation; first sufficient OR alternative and the
terminal entry's consumption policy; repeated consumption and ignored failure
before a later item; raw inventory/list order, skill and scroll/deleted-box
effects; unique creation and legacy RNG aftermath; **91** experience/duration
boundary pairs; magical speed zero/one-day/less/equal/greater/negative waits and
quick-cast values plus nonmagic/save-quick/trance exemptions; and experience
4→5, 11→12, 20→21, 34→35, unknown/learning/known, no_exp, school versus
child and once-per-month flags. Advancing the calendar alone does NOT reset
`exp_this_month`; this helper reads the flag and increments only the given
entry, not its parent, despite C's comment. Actual unmet trailing-OR assertions
at **412** and **463** compare status **134** and raw diagnostics. Original
index-before-bound OR-loop order is preserved; successful trailing-OR
out-of-bounds reads are not turned into UB expectations.

Real C twice and Go match **1260 sections / 86711 output lines / 1703330 output
bytes**, from **27307 input lines / 523069 input bytes**. All **1248** inherited
sections and both inherited byte prefixes remain unchanged:

| File | Prefix bytes | Prefix SHA-256 |
|---|---:|---|
| command-cases.txt | 517452 | `17250c35be656f2f61bce19409fe32e2cffb936f9b77fe2363c7b7a71b88469e` |
| command-c-output.txt | 1682820 | `b211c52e66534520a987e79411e6916f8465e94412270f2c1462a0f158b3da95` |

Full input SHA-256:
`c0099469c48a560b66f8dbcb4e3d6131b9b94ed182011e15a5397ef4a521f0e5`.
Full output SHA-256:
`f77a1f488eb9ae4482777766d79a4ed2e4ebd36d4d574d29d4825e485ea25975`.
The new literal C output preserves two intentional trailing-space `mu` list
lines; unrestricted `git diff --check` reports those raw bytes. Source-only
whitespace, gofmt and shell syntax checks pass.

Before recording refresh, four temporary linker rejection controls wrap only
the direct requirement helper observer at expected-abort section **1254**:
emit the missing-C-dependency marker then either return **0** or abort **134**,
in ordinary and `--update` modes. All reject with harness **2** BEFORE expected
status/recording acceptance; recording hashes stay unchanged. Temporary probes
and matching captures are removed. An exploratory nonexistent-owner call
reached C's invalid `rp_char` dereference (**139**, violates the caller's
valid-owner precondition), not a defined helper assertion; it was rejected
before recording and removed rather than manufacturing a UB expectation.
No valid-input upstream blocker or gameplay prerequisite was found.

Validation: bounded vet and fresh bounded default tests pass (Olympia
**1.780s** on final rerun). Sequential command, move, helper, RNG, code, loadsave,
turn-driver and default-turn checks pass; CLI **92** vectors and both
`rng-selftest` gates pass. EVERY compat call explicitly sets
`G3_REPO=$PWD/tmp/upstream.git`; engines/tests use run_limited. CLI reuses only
the pristine exact-pin C executable in `tmp/stage3-roundtrip/g3/build/debug`
(SHA-256 `1f03220bac19c2b3fd3a75c3b3aac34cc0df56210472148b11fba82d8e05c8bf`),
never its cached inputs. Logs: `tmp/issue171-source-audit.log`,
`tmp/issue171-recording-integrity.log`, `tmp/issue171-{vet,tests}-final.log`,
`tmp/issue171-command-{update,default}.log`,
`tmp/issue171-negative-{return,abort}-{default,update}.log`,
`tmp/issue171-*-default.log` and `tmp/issue171-cli.log`.

Only metadata and shared helpers are implemented here. #172 owns v_use,
d_use, i_use, shared item dispatch and USE accounting/CLI integration;
remaining gameplay handlers, full-corpus reruns and stage exit are NOT claimed.

Independent parent acceptance repeats bounded vet, fresh bounded tests,
command/move/helper/RNG/code/loadsave/turn-driver/default-turn and CLI checks;
all pass. The command oracle repeats real C twice before comparing Go:
**1260 sections / 86711 lines**, with both inherited recording prefixes
independently verified unchanged. Parent reviewed the complete pinned use.c,
table, shared helpers and diff; missing gameplay remains explicitly NYI.
Logs: `tmp/parent171-{vet,tests,command,move,helper,rng,code,loadsave,turn-driver,turn,cli}.log`.
The all-open audit found new follow-C #179, attached to the stage: G3 #360
is still unmerged and upstream main remains the immutable stage pin. Its
trade-loader/name-truncation changes are pending upstream, not silently
applied to this unrelated USE prerequisite. Recheck before the next task.

## Command coverage (#170)

At **0.106.0-alpha**, coverage.go implements complete pinned coverage.c
**28–117**, with world-owned optional path and lazy `[COV_MAX]int32` rows.
The oracle compiles the **2627-byte unchanged prefix (lines 1–119)**, including
globals, cov_num_cmds, coverage_count, cov_name_comp and write_coverage. Its
SHA-256 is `ded3432cf53cda5786bbc7d0913227d6cf91b24b9e3aa853abd56daefea472b8`.
Existing complete input.c check_allow and do_command remain unchanged in the
C reference; Go inserts hooks at the identical branches and retains the
original command index across reparsing. Definitions/dependencies verified:
oly.h COV constants, proto.h declarations, whole glob.c table, z.c
str_save/safe_qsort/asfail, checked_alloc.c zero-initializing allocator and
main.c long option parsing plus write-before-check/save lifecycle. The full
coverage.c USE/phase definitions were read but remain outside this task.
All 100 exported lib/olympia regular source files match the immutable archive.

The actual table has **148** named rows (empty-name placeholder included)
and a terminal NULL-name/allow/start row at index **148**. The sentinel is
excluded from counters/TSV membership, not replaced with another verb. Direct
dispatch of this real C sentinel supplies the genuine NULL-start path:
check_allow(NULL) succeeds, start is NULL, and coverage maps index 148 to row
0's unimplemented cell. No C table or callback is changed for the test.
Go preserves that implicit terminal-row dispatch, including fuzzy NULL-name
output. The mail-only NULL-start rows have empty (not NULL) allow strings and
are rejected earlier, as the real `begin` case proves. Go NYI callbacks still
panic before any C start outcome is recorded; a dedicated test protects that
boundary without pretending to be a C success/failure expectation.

Fourteen appended cases cover disabled/immediate guards before invalid-outcome
assertions/allocation, NULL versus empty paths, zero-inclusive byte-sorted
rows, first/last/sentinel/out-of-range indexes and all outcomes, signed 32-bit
boundary values (no C signed-overflow UB expectation), repeated/truncating
writes, retained counters across path changes, nonfatal ENOENT/ENOTDIR/EISDIR,
real unknown/empty/allow rejection, start success/failure, start-only timed
accounting (not finishes/interrupts), successful/failed cardinal reparsing and
the actual sentinel NULL-start. Two invalid outcomes abort at the unchanged
coverage.c **62** assertion. USE/phase instrumentation and gameplay callbacks
are untouched; other dumps/later modes remain explicitly gated.

Pinned real C twice and Go match **1248 sections / 85822 lines / 1682820 output
bytes**. Inputs are **27052 lines / 517452 bytes**. Exact inherited prefixes:

| File | Unchanged prefix bytes | Prefix SHA-256 |
|---|---:|---|
| command-cases.txt | 513156 | `6a711d79ca5d5a65ca2184b430766037e89df591c6caf1c662bdd7498c174c27` |
| command-c-output.txt | 1635428 | `d08f0a35fbcd021dc44470da263bf483302cd070c31a8ea3bcfe54321b6b4336` |

Full input SHA-256: `17250c35be656f2f61bce19409fe32e2cffb936f9b77fe2363c7b7a71b88469e`.
Full C-only output SHA-256: `b211c52e66534520a987e79411e6916f8465e94412270f2c1462a0f158b3da95`.
The initial timed fixture attempted WAIT interruption one day too early and
correctly failed with real C SIGSEGV; it was not recorded or accepted. Adding
the real third HIDE finish reaches the intended WAIT callback in both engines.

Before any recording refresh, four temporary linker controls wrapped
coverage_count only at the enabled negative-outcome assertion section **1240**:
emit a missing-dependency marker and either return **0** or abort **134**, in
both ordinary and `--update` modes. Every run rejects with harness **2** before
expected-status/recording acceptance, preserving the recording hash. Temporary
probe source/compiler wrapper and captures were removed. Logs:
`tmp/issue170-negative-{return,abort}-{default,update}.log`.
Both command and whole-engine harnesses reject dependency/NYI/stub/sanitizer
markers before accepting results; whole-engine tests exercise either stream
and missing streams too.

Validation, Linux amd64, unchanged pin, all changes left uncommitted:

- Bounded `go vet ./...` and fresh `go test -count=1 -exec "bash
  $PWD/tools/lib/run-limited.sh" ./...` pass (Olympia **1.757s**).
- Sequential command, move (**135/7158**), helper (**37/1487**), RNG (**321**),
  code (**43/969**), loadsave (**61/3411**), turn-driver (**7/152**, phase
  **5/484**, calendar **2/119**, schedule **4/1801**) and default-turn pass.
- CLI matrix: all **92** vectors pass and both real rng-selftest gates say YES.
- `turn-compat --full --command-coverage`: unchanged full C twice and Go
  match status **0**, raw streams and all **233** ordinary cwd entries;
  all ten coverage CLI cases and four independent flag-free controls pass,
  followed by all **13** inherited strict lifecycle controls. Zero/nonzero
  C dump hashes are respectively
  `be13a1dbbe7c01a22b3a822960e3e2fa53794211ba1d9829ebc61b168fc0118f`
  and `5fc30138c641ccedbb73dd510ae33cac4a6f230ab17c13b203a4285a3a73af77`.

Every compat invocation explicitly sets `G3_REPO=$PWD/tmp/upstream.git`;
whole-engine tools reuse only the pristine read-only `tmp/stage3-roundtrip`
C binary, running on real directory copies. Engines/tests use run_limited.
No shared upstream checkout, cached input, C source pin, unrelated recording
or gameplay callback is changed. Source-only whitespace, gofmt and shell
syntax checks pass; the new literal C recording intentionally retains one
trailing-space command echo (`> `), so unrestricted `git diff --check` reports
that raw reference byte. Logs: `tmp/issue170-{source-prefix,vet,tests}.log`,
`tmp/issue170-command-{update,default}.log`, the seven other
`tmp/issue170-*-default.log`, `tmp/issue170-cli.log`, and
`tmp/issue170-turn-full.log` / `tmp/issue170-turn/command-coverage.txt`.
No new bug, blocking upstream defect or oversized prerequisite was found.
No native Darwin/Windows run or stage-3 final-corpus/spec acceptance is claimed.

Independent parent complete-source/diff review confirms the command prefix and
both inherited recording byte prefixes unchanged. Bounded vet, fresh bounded
tests and all nine routine tools above pass again. The parent also repeats
`turn-compat --full --command-coverage`: strict 233-entry flag-free turn,
all ten coverage CLI cases, four independent flag-free aftermath controls and
all thirteen lifecycle controls pass. No C cache or input changes. Evidence:
`tmp/parent170-{vet,tests}.log`, `tmp/parent170-*-default.log`,
`tmp/parent170-turn-full.log` and `tmp/parent170-turn/command-coverage.txt`.

## Auxiliary files and saved report lifecycle (#165)

Existing olympia.go replaces its five writer projections directly in main
and the oracle with complete main.c bodies, keeping opens before traversal
and the original players/email/totimes/forward/factions order. Shared real
ownership, garrison traversal, names and output-cache owners are reused;
write_faction_sup and write_forward_sup have the same Go implementation,
matching their identical C bodies. There is no second effects pass/cache.

| Complete original definition | Inclusive pinned lines |
|---|---|
| main.c `write_totimes` | 50–79 |
| main.c `write_email` | 82–107 |
| main.c `list_a_player` | 109–149 |
| main.c `write_player_list` | 152–189 |
| main.c `write_forward_sup` | 192–204 |
| main.c `write_forwards` | 207–248 |
| main.c `write_faction_sup` | 251–263 |
| main.c `write_factions` | 266–290 |
| c2.c `open_times` | 338–377 |
| c2.c `times_masthead` | 380–454 |
| c2.c `close_times` | 457–472 |
| io.c `save_logdir` (unchanged entire IO translation unit) | 3031–3069 |

All eleven extracted definitions are checked unchanged exactly once; the
entire real io.c remains unchanged. `times_masthead` and its original months
declaration now join the existing extraction. The fixture enables the real
constant-report-date flag; no Times implementation is added. Existing c2.go,
io.go and sout.go are unchanged. No source dependency or pin changes.

Sixteen appended sections (**12 success / 4 fatal status 1**) compare real
bodies, raw diagnostics/status and saved state:

- Private versus visible email; empty versus absent email/visible/full-name
  pointers; regular/system/silent/missing records; CIS exclusion only from
  totimes. Direct original static list helper cases cover public-only,
  name-only, empty-name and absent-record allocation branches unreachable
  through the root's email eligibility gate. Duplicate preseeded new-player
  IDs produce stars and one footer, without admitting anyone.
- Two garrisons and a three-noble pledge chain deduplicate factions, not
  equal addresses: each garrison has its ordinary character route plus two
  identical-address ruling routes (three matching rows each). Empty addresses
  produce email/totimes blank lines but no forward/faction rows.
- Failure after all five successful writers preserves renamed old bodies,
  reports all five EISDIR failures and leaves missing player records absent.
  Inherited ownership-cycle/castle failure controls still establish that
  validation happens only after successful forward open.
- Actual banners/account/include logs and players bytes match before/after
  save, with an explicitly empty recreated log; no-save retains live reports.
  Buffered cached output reaches the renamed inode twice (**406 bytes**),
  then close/reset creates a new **203-byte** live log, leaving saved bytes
  unchanged. No duplicate report pass or invented output membership.
- Times has zero/singular/plural **0/1/3** regular players, including missing
  player records without allocating them; turns **0/8/9/17**, month-zero clamp,
  year rollover, custom From/Subject/URL, fixed **January 1, 2000** date,
  exact **201-byte press / 116-byte rumor** headers, repeated renamed streams
  (**402/232 bytes**) and close/reset. Masthead EISDIR/ENOENT/ENOTDIR and
  rumor failure retain fatal status **1** with no perror/newline. The inherited
  press-before-rumor flush and press/rumor failure controls remain real.

Constants audited against oly.h: player **1**, char **2**, system/regular/
silent **48/49/62**, garrison **64**, **8** months, raw player_email and
player_compuserve accessors. The recorded body rows, counts and lengths above
are actually reached C branches, not fixture intentions.

Pinned C twice deterministically and Go match **1234 sections / 82688 lines /
1635428 output bytes**. Inherited **1218 sections / 81791 lines / 1615381
output bytes**, and **26569 lines / 507858 input bytes**, remain exact prefixes.
New input is **26815 lines / 513156 bytes**. SHA-256:

- Input: `6a711d79ca5d5a65ca2184b430766037e89df591c6caf1c662bdd7498c174c27`
- Output: `d08f0a35fbcd021dc44470da263bf483302cd070c31a8ea3bcfe54321b6b4336`
- Inherited input: `03264d856def78a7aca8ed6b2592f63bd9915c83303b5b86c819020f39c07681`
- Inherited output: `237f5b00267c9c5fe26739bf7b9ef972b5f651961928423dd35453ef46de30a4`

Four temporary missing-dependency probes wrap write_forwards at inherited
ownership-cycle section **818**, which ordinarily requires **134**. Both
return **0** and abort **134**, in default and --update modes, reject with
harness **2** at the marker check before expected-status/update acceptance.
Both recording hashes remain unchanged. An initial probe also rejected the
prior failed-open section **817**; it was narrowed to section 818 before all
four final controls. Probe source/compiler wrapper and matching captures are
removed. Logs: `tmp/issue165-negative-*.log`.

Sequential command refresh/default/post-probe, move **135/7158**, helper
**37/1487**, RNG **321 results**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**,
and default turn pass. Default turn compares raw status/stdout/stderr and two
cwd entries under the strict contract, not a world turn. Bounded vet and fresh
bounded default tests pass. Every compat call explicitly sets isolated
`G3_REPO=$PWD/tmp/upstream.git`; engines/tests use shared run_limited.
Logs: `tmp/issue165-{command-*,validation,source-audit,vet,tests}.log`.

Version **0.105.0-alpha**. Parent's all-four-open-issue
audit is the task prerequisite: #165/#166 and trackers #19/#132, no bug/follow-C.
No newly verified defect/prerequisite/oversized closure or independent fix.
G3 #358/#359 are unchanged/nonblocking; malformed **C136/Go2** is not parity.
No valid admissions, later commands, charge/mail/HTML conversion or new owner.

Independent parent review verifies all eleven extracted definitions complete,
unchanged and exactly once, plus exact inherited input/output prefixes. All
eight routine tools, bounded vet and fresh bounded tests pass (Olympia
**1.726s**); logs `tmp/parent165-*.log`. Fresh strict full comparison passes
two identical C turns and exact C/Go status **0**, raw streams and all **233
cwd entries**, reports/types/modes/empty directories/bytes included. All **13
lifecycle controls** pass under that same contract. Evidence:
`tmp/parent165-strict.log`, `tmp/parent-stage2-turn/lifecycle.txt`.
The broader final tier and stage acceptance remain #166.

## Include section assembly (#164)

Existing report.go completes all six pinned report.c definitions:

| Definition | Inclusive lines |
|---|---|
| `rep_player` | 1353–1363 |
| `rep_char` | 1366–1388 |
| `rep_loc` | 1391–1410 |
| `inc` | 1413–1420 |
| `gen_include_sup` | 1423–1566 |
| `gen_include_section` | 1569–1586 |

The requested 1568 start includes a blank line before the complete final
definition. The existing oracle already extracts all six complete unchanged
definitions, each verified exactly once against the immutable git object,
with original `#line` attribution. Every OUT/LOG/EAT constant, including
OUT_SUMMARY's banner alias and the unsupported OUT_HTML_INDEX, is verified.
No new C dependency, source extraction, pin or successful substitute is added.
Adapters only set a corpse's saved name and the scanner recipient for fixtures.

Both include projections are removed. Main and oracle call the real root
directly, once; main still calls it LAST immediately before closeLogfile.
Real names, lines/style output, showLocHeader, output sorting and the existing
bounded FILE cache are reused. World-owned character/location scratch lists
retain their allocation across clear/repeated calls, matching the C statics.
Output memberships come exclusively from real output, not a projection pass.
Initial validation precedes output; corpse misc allocation occurs only later
in repChar. Original output knowledge sorts in place; unit order is unchanged.

Thirteen appended cases (**8 status 0 / 5 status 134**) actually reach:

- Every recognized special-code branch, using actual sout output to populate
  membership. Exact include order is **0,24,14,10,11,13,12,20,21,7,22,200,6,4,
  2,3,5,25,23**. EAT_OKAY/OUT_LOC lack headings; banner/include inputs do not
  recurse. Repeating gives **4,058 → 7,782 bytes**, not a second effects pass.
- Empty outputs contain only the real include bit **1** and **28 bytes**.
  Missing regular/system/NPC/savage records allocate; silent records remain
  absent and produce no file. Routing resets to **0/0**. Scanner recipient
  suppresses its own player heading but still includes **0,24,23**.
- Living/dead characters, named/empty/NULL corpse saved names, and real
  plain/city/ship headers. Output-order and ID ties produce **0,3001,1001,
  3000,3002,1000,3003,9100,9000,9001,9200**, **3,004 → 5,778 bytes** on repeat.
  Absent corpse misc becomes allocated **0 → 1**; original units stay
  **1000,1001**. NULL uses the actual libc `(null)` spelling; the saved name,
  rather than the body's current name, supplies the heading.
- Invalid **149999/9999** and absent index **8** skip inclusion. Complete
  real player/character/banner/template calls produce only actual memberships
  **1000,1001,0,5**, then include adds **1**; no player bit is manufactured
  when empty player_ent_info emits nothing. The resulting report is **2,205
  bytes**. Sequential recipients clear prior character/location scratch.
- Buffered output keeps a renamed cached inode: **277 → 529 bytes**, with
  the original path absent. Closing permits a fresh **252-byte** original
  file, without changing the renamed **529-byte** file.
- Foreign-player **report.c:1495**, ordinary-item/present-index-8
  **1503**, unsupported-skill **1507**, and output-open **sout.c:158**
  assertions. Validation/open failures precede corpse misc allocation.

At **0.104.0-alpha**, unchanged C twice and Go match **1218 sections /
81791 lines / 1615381 output bytes**. Inherited **1205 sections / 80495 lines /
1579486 output bytes** remain an exact prefix, SHA-256
`3e2694726523765da39111b0479262cca4ebe42145f03ee85b316594512b9c90`.
Inherited **26294 lines / 501764 input bytes** remain an exact prefix,
SHA-256 `9b6227ece07be605c4abf80cfa8b293910489f64bd4471752b8757a7dc79c2be`.
New input: **26569 lines / 507858 bytes**, SHA-256
`03264d856def78a7aca8ed6b2592f63bd9915c83303b5b86c819020f39c07681`.
New output SHA-256:
`237f5b00267c9c5fe26739bf7b9ef972b5f651961928423dd35453ef46de30a4`.
Constants, definitions, prefixes and reached branches are checked in
`tmp/issue164-{source-evidence,evidence}.log`.

Four temporary linker probes wrap the real include root only on foreign-player
assertion section **1213**, emit the rejected missing-dependency marker, and
return **0** or abort **134** (both child statuses separately verified).
Default/update runs all reject with harness **2** at the marker check before
expected-status/update acceptance. Both recording hashes remain unchanged.
Probe source/compiler wrapper and matching captures are deleted. Logs:
`tmp/issue164-negative-{return,abort}-{default,update}.log`,
`tmp/issue164-probe-{evidence.log,before.sha256}`.

Sequential command refresh/default/post-probe-default, move **135/7158**,
helper **37/1487**, RNG **321**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**,
and default turn (raw streams/status and **2 cwd entries** equal) pass.
Bounded vet and fresh bounded default tests pass (Olympia **1.854s**).
Every compat invocation explicitly uses isolated `G3_REPO=$PWD/tmp/upstream.git`;
engines/tests use shared run_limited supervision. Source-only whitespace and
gofmt checks pass; the single new trailing-space save_box line is kept raw.
Logs: `tmp/issue164-*.log`. An initial new composition fixture omitted item
10 required by real character weights; adding the existing prerequisite box
fixed that fixture, without a renderer/engine change or recording acceptance
of its C139 failure.

The parent supplied the all-open-issue audit: only #164–166 and trackers
#19/#132, no bug/follow-C. No new defect, missing renderer prerequisite or
oversized closure found. No later writers/commands/HTML conversion, independent
fix or second effects pass. Independent parent complete-C/diff review verifies
all six definitions unchanged exactly once and both inherited prefixes
byte-exact. All eight routine tools, bounded vet and fresh bounded tests pass
(Olympia **1.805s**); logs `tmp/parent164-*.log`.

Fresh strict post-#164 has only **one** saved-file byte difference:
`lib/save/1/players`. C repeats deterministically, both engines exit **0**,
and raw streams/paths/types/modes match. Exit **1** still precedes all 13
lifecycle controls, NOT stage acceptance. Log `tmp/parent164-strict.log`;
captures `tmp/parent-stage2-turn`. G3 #358/#359 remain nonblocking;
malformed **C136/Go2** is not parity.

## Complete GM report and faction tables (#163)

Existing gm.go completes seven pinned gm.c definitions:

| Definition | Inclusive lines |
|---|---|
| `wealth_list_comp` | 739–748 |
| `gm_faction_wealth` | 751–791 |
| `nobles_list_comp` | 794–803 |
| `gm_nobles_list` | 806–846 |
| `gm_player_details` | 849–961 |
| `list_all_notices` | 964–992 |
| `gm_report` | 1025–1047 |

The existing command oracle already extracted the complete real definitions.
All **23** gm.c definitions, including these seven, are verified unchanged
exactly once against the immutable git object, with original `#line` locations.
Adapters expose faction/detail/notice leaves and observe player/skill scratch;
they do not rewrite C reporting. No C source dependency or pin changes.
Go removes gmReportState from report.go and replaces its main/oracle calls
with the real gmReport root, preserving every call once in source order.
Final listAllItems still targets **skillPlayer**, independently of the GM.
Real showLocPosts, cached output ownership, player/unit/skill loops, ownership,
formatters and scratch sorting are reused. Two static faction lists retain
their allocation across calls. No later include/writer/command work or
parallel projection/rendering pass is added; this closure reads no lore files.

Twelve appended cases (**8 status 0 / 4 status 134**) actually reach:

- Wealth **1,234 / 1,234 / 1,234** and noble-count ties in ascending ID order
  despite reversed allocation, regular-only ranking rows, and scratch counting
  all owners including silent/NPC/system (**9,999 / 8,888 / 7,777** gold).
  Changed lords and repeated/empty reports prove scratch and list clearing.
- Details exclude system but include silent/NPC players, preserve missing
  player records with **???** age, and show both **14** and **-2** ages.
  Real ownership produces **4,001 gold / 6 units / 1 building / 1 sublocation /
  2 ships / 3 unique known skills**, including finished/unfinished ships,
  foreign ownership, a prisoner and duplicate/partial/unexpected skill rows.
  Repeats preserve the original rows; final player/skill scratch is observed.
- Actual creator/anonymous, multiline/empty-line/raw-inline posts in a
  province and city, full location names and indentation; a ship post is
  counted in ship scratch but omitted by the original location-only loop.
  Invalid post text asserts at **display.c:912** before a separator/output.
- Complete root output order, all-player allocation and stored knowledge
  sorting, with **2,804-byte** GM output and independent **189-byte** item
  output. Repeated zero-gate roots with the same recipient retain a renamed
  cached inode (**5,132 bytes**); the original path remains absent.
  Root failure controls preserve output-open **sout.c:158**, loyalty
  **gm.c:659** and post **display.c:912** assertion ordering.

At **0.103.0-alpha**, unchanged C twice and Go match **1205 sections /
80495 lines / 1579486 output bytes**. Inherited **1193 sections / 79511 lines /
1556344 output bytes** remain exact, SHA-256
`1ce0bf7f4c6aaf6772cc23a8e857cd8f17ab972805b79474ad5bbb7bb873d1a0`.
Inherited **25968 lines / 495251 input bytes** also remain an exact prefix,
SHA-256 `68727937f0b5f22793da9651202a84d5f945f9e442409ad75c6b222d37a9ec6b`.
New input: **26294 lines / 501764 bytes**, SHA-256
`9b6227ece07be605c4abf80cfa8b293910489f64bd4471752b8757a7dc79c2be`.
New output SHA-256:
`3e2694726523765da39111b0479262cca4ebe42145f03ee85b316594512b9c90`.
Constants, complete definitions, exact prefixes and reached branches are
checked in `tmp/issue163-{source-evidence,evidence}.log`.

Four temporary linker probes wrap only the invalid-notice adapter call,
emit the rejected missing-dependency marker, and return **0** or abort **134**.
Default/update runs all reject **section 1199**, harness **2**, before expected
status or update acceptance; both recording hashes remain unchanged.
Probe source/compiler wrapper and matching probe captures are removed.
Logs: `tmp/issue163-negative-{return,abort}-{default,update}.log`,
`tmp/issue163-probe-evidence.log` and `tmp/issue163-probe-before.sha256`.

Command refresh/default, move/helper/RNG/code/loadsave/turn-driver/default
turn, bounded vet and fresh bounded defaults pass sequentially (Olympia
**1.491s**); logs are `tmp/issue163-*.log`. Every compatibility invocation
explicitly uses `G3_REPO=$PWD/tmp/upstream.git`; engines/tests use shared
run_limited. Matching command captures are removed.
Source-only whitespace/gofmt checks pass; full diff whitespace checking flags
exactly eight intentional trailing-space C-recording lines (unit lists and
region totals/separators). These bytes are not normalized.

The supplied parent audit covers all six open issues, including unlabeled:
only #163–166 and trackers #19/#132, no bug/follow-C. No new verified defect,
missing renderer prerequisite or oversized closure was found. Independent
parent complete-C/diff review verifies all 23 gm.c definitions unchanged
exactly once and both inherited prefixes byte-exact. All eight routine tools,
bounded vet and fresh bounded tests pass (Olympia **1.728s**);
logs `tmp/parent163-*.log`.

Fresh strict post-#163 remains red: deterministic C repeats, both engines
**0**, equal raw streams/paths/types/modes, eight report-byte differences,
exit **1** before 13 lifecycle controls. Log `tmp/parent163-strict.log`;
captures `tmp/parent-stage2-turn`. No stage-exit acceptance is claimed.
G3 #358/#359 remain nonblocking; malformed **C136/Go2** is not parity.

## GM skills geography and loyalty (#162)

Existing gm.go completes nine pinned gm.c definitions, verified unchanged
exactly once in the existing real-C extraction with original `#line` locations:

| Definition | Inclusive lines |
|---|---|
| `skill_use_comp` | 40–54 |
| `gm_show_skill_use_counts` | 57–104 |
| `skills_known_comp` | 107–116 |
| `gm_show_skills_known` | 118–179 |
| `gm_show_gate_stats` | 476–536 |
| `gm_show_locs_visited` | 539–634 |
| `gm_loyalty_stats` | 637–678 |
| `region_occupy_comp` | 681–690 |
| `gm_land_stats` | 693–736 |

Go reuses actual entity loops, copied known-skill rows, player/skill/gate
record access, in-place knowledge sorting, location topology, region lookup,
names/codes and the existing cached output owner. The three original static
scratch lists have world-owned counterparts, cleared without discarding
their allocation; entity temps retain the original reset/count side effects.
There is no new C dependency or pin, rewritten C model, successful stub,
second render/projection pass, later root/include/writer or command. Existing
gmReportState and main remain unchanged until #163.

Thirteen appended cases (**12 status 0 / 1 status 134**) actually reach:

- Use-count descending order, **1,234** ties with lower ID first despite
  reversed allocation order, category/zero exclusion, negative counts, owner
  lookup and list clearing after a formerly used skill becomes unused.
- Known counts **3 / 2 / 2 / 0**, lower-ID ties, duplicate known rows,
  learning/don't-know/unexpected states excluded, zero-use blank versus
  comma-formatted use counts, school exclusion, all characters including
  silent-owned prisoners/garrisons/unowned units, and absent character records
  remaining absent. Original skill rows/order are unchanged on repeats.
- Zero gates: all fixture players still allocate and knowledge still sorts,
  but routing deliberately remains **MASTER (-2) / OUT_LORE (2)**. The next
  ordinary output demonstrates that retained route, not just an observer.
  Nonzero controls report **2/3 (66%)** found, **1/3** sealed, **2/3** notify
  jump and **1/3** notify unseal; a separate control reports **0/2** found.
  Missing gate records stay absent and stale gate temps are cleared.
- Empty and nonempty location denominators; **1/2** provinces and **3/5**
  sublocations visited, **33% visible / 66% hidden**, and **50%** of the
  visited province's immediate sublocations found, split **50% / 50%**.
  Buildings/ships/nonlocations and sublocations in unvisited provinces do not
  contaminate the immediate-sublocation denominator. All-player allocation,
  stored knowledge sorting and location scratch resets remain observable.
- Loyalty **6** total: **2** oath, **1** fear, **1** contract, **1** unsworn,
  plus **1** NPC in the denominator; silent owners/prisoners count, nonzero
  character subtypes skip invalid loyalty, missing records read as unsworn,
  and the empty ordinary-character denominator is safe. Invalid loyalty **5**
  emits `loy_kind = 5` then the exact **gm.c:659** assertion before output.
- Region counts **2 / 1 / 1**, ID-ordered ties despite reversed allocation,
  all-character inclusion, empty-region exclusion, stale-temp clearing and
  changed counts after relocation; total remains **4** on both reports.
- An **838-byte** renamed cached inode across repeated reports and a
  zero-gate/ordinary-output transition; the original path remains missing.
  Silent output recipients allocate and record output through the real owner.

At **0.102.0-alpha**, complete unchanged C twice and Go match **1193
sections / 79511 lines / 1556344 output bytes**. Inherited **1180 sections /
78766 lines / 1540387 output bytes** remain an exact prefix, SHA-256
`3956cf60d3cb631b65c5e3b362b978e1b686d8bc3d26cee815698d5a0a2c91f6`.
Inherited **25661 lines / 488967 input bytes** also remain exact, SHA-256
`35e452b48a6a3580e30b2334d978781712a6320c37a07a12d6de7f7829f37f36`.
New input: **25968 lines / 495251 bytes**, SHA-256
`68727937f0b5f22793da9651202a84d5f945f9e442409ad75c6b222d37a9ec6b`.
New output SHA-256:
`1ce0bf7f4c6aaf6772cc23a8e857cd8f17ab972805b79474ad5bbb7bb873d1a0`.
Constants, complete definitions, prefix bytes and reached output are checked
in `tmp/issue162-evidence.log`.

Four temporary linker probes wrap only the adapter's invalid-loyalty call,
emitting the rejected dependency marker and returning **0** or aborting
**134**. Default/update runs all reject **section 1192**, harness **2**,
before expected-status or update acceptance; both recording hashes stay
unchanged. Probe source/compiler wrapper and matching captures are removed.
Logs: `tmp/issue162-negative-{return,abort}-{default,update}.log`,
`tmp/issue162-probe-evidence.log` and `tmp/issue162-probe-before.sha256`.

Command refresh/default and post-probe default, then move/helper/RNG/code/
loadsave/turn-driver/default turn pass sequentially, followed by bounded vet
and fresh bounded default tests (Olympia **1.737s**). Every compatibility invocation explicitly
uses `G3_REPO=$PWD/tmp/upstream.git`; engines/tests use shared run_limited.
Logs: `tmp/issue162-*.log`. Source-only whitespace/gofmt checks pass; full
`git diff --check` flags exactly four intentional trailing-space C-recording
lines from the region totals/separators. These bytes are not normalized.

The supplied parent audit covers all seven open issues, including unlabeled:
only #162–166 and trackers #19/#132, no bug/follow-C. No new verified defect,
prerequisite or oversized closure was found. No issue/checklist changes,
commit/push, full-world/lifecycle or stage-exit claim. Latest strict evidence
remains parent's post-#160 red gate: deterministic C, both engines **0**, equal
raw streams/paths/types/modes, eight report-byte differences, exit **1**
before 13 lifecycle controls. `tmp/parent-stage2-turn` is untouched.
G3 #358/#359 remain nonblocking; malformed **C136/Go2** is not parity.

Independent parent complete-C/diff review verifies all nine definitions
unchanged exactly once and both inherited prefixes byte-exact. All eight
routine tools, bounded vet and fresh bounded tests pass (Olympia **1.828s**);
logs `tmp/parent162-*.log`. Strict evidence remains post-#160, not acceptance.

## GM item economy and attributes (#161)

New source-corresponding gm.go completes these unchanged pinned gm.c
definitions, each verified byte-identical exactly once in the existing
extraction (original `#line` attribution retained):

| Definition | Inclusive lines |
|---|---|
| `pct` | 31–37 |
| `gm_show_interesting_attributes` | 183–280 |
| `gm_list_animate_items` | 283–329 |
| `gm_show_gold` | 332–377 |
| `gm_show_control_arts` | 380–406 |
| `gm_count_stuff` | 409–473 |
| `list_all_items` | 995–1022 |

The oracle exposes those complete real definitions, not rewritten C models.
It reuses complete original `beast_capturable` (u.c **2547–2561**),
`token_player` (art.c **1019–1032**), `player` (swear.c **21–37**), names,
codes, loops, output/cache and real save writers. No new C source dependency,
pin, successful stub or second output layer is introduced. All 16 original
C gold declarations are retained unchanged. Go reuses the ten existing actual
accounting fields and adds only the six missing pinned scalars: common magic,
lead to gold, Times, petty thief, pillage and ferry. No gameplay is added to
populate them. Existing gmReportState and main's call sites are unchanged;
the later skills/geography/loyalty and complete GM root remain #162/#163.

Twelve appended cases (**11 status 0 / 1 status 134**) reach and compare:

- Empty categories, zero regular-player denominator, no galley line, and
  signed/truncating `pct`, including zero denominators.
- Attribute counts **2/1/3/1**, location counts **2/1/2**, format **1/3**,
  and ram **1/3 (33%)**. Negative nonzero flags count once; ship location
  attributes, nongalley rams and nonregular format flags are excluded.
  Missing player/character/location/sublocation records stay missing.
- Every animate OR predicate, separate attack/defense/missile selection,
  man-only/animal-only rows, all flags, and ghost **18**'s zero-stat exception.
  Capturability alone does not select an item. Raw fixtures with actual
  noble-item records exercise both outcomes of the original character-based
  beast predicate even on the item chain; no substituted animal predicate.
- All 16 unequal gold counters. Generated gold stays **12,083** when player
  combat changes **999,999 → -111111**; the latter's original comma formatting
  is retained. Zero/negative/cancelling generated sums suppress rows, not
  headers; a positive sum of one reaches **-300% / 400%**.
- Two captured regular-owner tokens, ownership changes and independent,
  NPC/silent/noncharacter/absent-owner exclusions; unrelated items are skipped.
- Every finished/unfinished structure category, unequal totals, real subtype
  chain transitions, collapsed/other subtypes and wrong-kind exclusions.
- Item listing with missing records, long names, unequal signed/nonzero
  weights and all capacities; real item-chain order and zero defaults.
- Route reset, a **927-byte** renamed cached inode across repeated sections,
  silent-recipient allocation/output, and the original output-open assertion
  at **sout.c:158**, before recipient allocation.

At **0.101.0-alpha**, two deterministic real C runs and Go match **1180
sections / 78766 lines / 1540387 output bytes**. The inherited **1168
sections / 78078 lines / 1522550 output bytes** are an exact prefix (SHA-256
`dbe620447bec3e071112f3e342ba0fbe70d309f4489265f24a0d36a4fcc6a52d`).
The inherited **25396 lines / 483251 input bytes** also remain an exact
prefix (SHA-256
`73b6858bc1bb9182873f3ed9fdb5c8a81f806eafb4ec2d5bd9e9096e1067a62b`).
New input totals **25661 lines / 488967 bytes**, SHA-256
`35e452b48a6a3580e30b2334d978781712a6320c37a07a12d6de7f7829f37f36`;
new output SHA-256
`3956cf60d3cb631b65c5e3b362b978e1b686d8bc3d26cee815698d5a0a2c91f6`.
Constants (kinds, every structure subtype, token/NI/regular subtypes, ghost,
MASTER **-2**, OUT_LORE **2**) and the rendered branches above are verified
in `tmp/issue161-evidence.log`, which also records complete-definition hashes.

Four temporary linker controls wrap only the oracle entry at expected-abort
section **1179**: a missing-dependency marker followed by return **0** or
abort **134**, in default and update modes. All reject with harness **2**
at the marker check, before expected-status/update acceptance; both recording
hashes stay unchanged. Probe source/compiler wrapper and matching captures
are removed. Logs: `tmp/issue161-negative-{return,abort}-{default,update}.log`.

Sequential command refresh/default (and post-probe default), move
**135/7158**, helper **37/1487**, RNG **321**, code **43/969**, loadsave
**61/3411**, turn-driver **7/152**, phase **5/484**, calendar **2/119**,
schedule **4/1801**, and default turn (raw streams/status and **2 cwd entries**)
pass. Bounded `go vet ./...` and fresh bounded default tests pass; Olympia
**2.302s**. Every compatibility call explicitly uses isolated
`G3_REPO=$PWD/tmp/upstream.git`; engines/tests use shared run_limited.
Full source/recording whitespace checks pass. Logs: `tmp/issue161-*.log`.

No new verified defect, prerequisite or oversized closure. A development-only
raw beast fixture lacked its noble-item record and C exited 139; the fixture
was corrected before recording, not accepted as malformed-input parity.
The supplied parent audit of all eight open issues found only #161–166 and
#19/#132, no bug/follow-C. No strict whole-world/lifecycle rerun or stage
acceptance is claimed: parent post-#160 remains red on eight report-byte
differences, engines **0**, deterministic C repeats and equal raw streams,
paths/types/modes, exit **1** before 13 lifecycle controls. G3 #358/#359 remain
unchanged/nonblocking; malformed **C136/Go2** is not parity.
No later GM root, includes, writers, commands, issue/checklist edits,
commit or push; all changes are left uncommitted for parent review.

Independent parent complete-C/diff review verifies the seven definitions
unchanged exactly once and both inherited prefixes byte-exact. The eight
routine tools, bounded vet and fresh bounded tests pass (Olympia **2.586s**);
logs `tmp/parent161-*.log`. Strict evidence remains post-#160, not acceptance.

## Skill lore queues and GM skill listing (#160)

Existing lore.go completes unchanged lore.c `out_skill_line` **137–146**,
`deliver_skill_lore` **149–199**, `deliver_lore` **202–221**,
`show_lore_sheets` **258–294** and `gm_show_all_skills` **297–369**.
Queue/comparator **230–255** and character skill/item scans **372–411**
retain their existing real owners. Existing u.go adds complete `weeks`,
u.c **1616–1630**, using the accepted nice-number/plural helpers. All thirteen
complete C definitions (the twelve lore functions and weeks) are verified
unchanged exactly once in the existing extraction. No new source/pin or
successful dependency substitute is needed.

Main and the existing oracle replace showLoreSheetsState/gmShowAllSkillsState
with the real roots. The diagnostic-only deliverLoreState/loreSheetInput and
optional-emitter plumbing are removed, as is the temporary
TestDeliverLoreSkillBoundary. One real sheet read/render pass uses the existing
commandOut/FILE cache; no second knowledge, allocation, scan or reader pass.
Queues sort in place using the original C-int subtraction comparator; adjacent
duplicates suppress delivery only and remain queued. Both roots reset routing
to zero, not the caller's previous route. GM output allocates through the first
real output/open before listing validation, lists all schools and their nested
descendants before any sheet reads, then delivers schools before nonschools.

Ten appended cases (**7 status 0 / 3 status 134**) compare raw bytes,
streams/status, actual saved records and list/record observers:

- Ordinary delivery reaches neither/offered-only/research-only/both lists.
  Offered **702,701,702** retains order and duplicates; hidden school research
  is not rendered as a skill row, while the same skill is subsequently shown
  when offered directly by a leaf. Exact headers, blanks, **one~day**,
  **one~week**, **two~weeks**, NP **-1/0/1/3**, and absent skill records match.
- GM listing reaches nested school membership, offered/research combinations,
  **0~days**, **eight~days**, **three~weeks**, and **-1~weeks**. Raw **2444**
  bytes prove all listing text precedes sheets in school order
  **638,700,900**, followed by nonschools **500,600,701**. GM-only research
  rows are visible without granting knowledge or queuing anything.
- Retained queue **700,700,701,701** and known insertion order **701,700**
  survive repeats. Modified input is reread while a renamed cached inode
  retains **1011** exact bytes; the original log path stays missing. MASTER
  **-2**, OUT_LORE **2**, CHAR_FIELD **6**, style **5**, initial indent **3**
  and secondary indent **8** are verified against actual reached bytes.
- Skill states **0/1/2/3**, duplicate/shared faction rows, prisoners, silent
  owners, and absent owner/character records reuse the existing scan/queue
  effects. Silent GM delivery really renders **293** bytes and allocates its
  recipient; silent players are not filtered by these roots. An unknown
  function still aborts for a silent owner at **lore.c:43**.
- Immediate GM listing reaches **use.c:1182** before a deliberately invalid
  sheet function can be read. Blocked queue output retains **sout.c:158**
  before any sheet validation. Empty queues reset routing without allocating
  absent player records. Inherited diagnostics/read-error/assertion cases
  remain unchanged and continue exercising the real C definitions.

At **0.100.0-alpha**, complete unchanged real C repeats deterministically and
Go matches **1168 sections / 78078 lines / 1522550 output bytes**. Entire
inherited **1158 sections / 77449 lines / 1508521 output bytes** and
**25169 input lines / 478453 input bytes** are exact byte prefixes.
New input: **25396 lines / 483251 bytes**. Output SHA-256:
`dbe620447bec3e071112f3e342ba0fbe70d309f4489265f24a0d36a4fcc6a52d`.
Input SHA-256:
`73b6858bc1bb9182873f3ed9fdb5c8a81f806eafb4ec2d5bd9e9096e1067a62b`.
Source-only whitespace checks pass; eight intentional C-recording trailing
spaces remain exact. Closure, constants, prefixes and reached-branch evidence:
`tmp/issue160-evidence.log`.

Four temporary linker controls wrap only gm_show_all_skills at the inherited
expected listing-cycle abort (**section 766**). Missing-dependency markers
with return **0** / abort **134**, in default/update modes, all reject with
harness **2** before expected-status/update acceptance and preserve the
recording hash above. Probe source/compiler wrapper and matching captures
are removed. Logs: `tmp/issue160-negative-{0,134}-{default,update}.log`.

Sequential command refresh/default **1168/78078**, move **135/7158**, helper
**37/1487**, RNG **321**, code **43/969**, loadsave **61/3411**, turn-driver
**7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**, and default
turn (equal status/raw streams and **2** cwd entries) pass, followed by bounded
vet and fresh bounded default tests (Olympia **1.981s**). Every compat call
explicitly uses isolated `G3_REPO=$PWD/tmp/upstream.git`; engines use shared
run_limited and test executables use tools/lib/run-limited.sh. Logs:
`tmp/issue160-*.log`. Temporary evidence assertions were corrected after
inspecting the actual C output (the leaf legitimately offers the otherwise
hidden research, and the silent GM file has 293 bytes); no engine/expected
output change was made to satisfy those inspection assertions.

Parent audited all nine open issues including unlabeled ones: only #160–166
and trackers #19/#132, no bug/follow-C. No new verified defect, blocking
prerequisite or oversized closure found. No later GM reports/includes/writers
or commands. Independent parent review verifies all twelve lore definitions
plus weeks unchanged exactly once and both inherited prefixes byte-exact.
All eight routine tools, bounded vet and fresh bounded tests pass (Olympia
**2.118s**); logs `tmp/parent160-*.log`.
Fresh strict post-#160 comparison repeats C deterministically and remains red
with eight saved-report byte differences, engines **0**, matching raw streams/
paths/types/modes, exit **1** before 13 lifecycle controls. Captures:
`tmp/parent-stage2-turn`; log `tmp/parent160-strict.log`. No stage-exit
acceptance. G3 #358/#359 remain unchanged/nonblocking; malformed **C136/Go2**
is not parity.

## Lore sheet rendering (#159)

Existing lore.go completes unchanged lore.c `lore_function` **9–44** and
`deliver_lore_sheet` **47–117**, plus `deliver_lore` **202–221** with one
explicit dependency boundary: its skill branch stops NYI until #160 ports
`deliver_skill_lore` and its offered/research output. A dedicated Go test
requires that stop before output or recipient allocation. The original C
skill renderer remains fully extracted and unchanged, not stubbed or modeled.
Direct sheets, item-to-skill sheets and non-skill delivery execute real C.

The old diagnostic-only sheet implementation is replaced by one shared
reader/parser. Real delivery emits through the accepted commandOut/cache
owner; existing stage-1 queue/GM projections use the same reader without
artifact emission, retaining all their input diagnostics and assertions until
#160 replaces those roots. They neither call the NYI real skill entry nor
perform an additional read. No queue/GM/main root is newly wired, no later
skill output is silently implemented, and no successful dependency stub,
new source dependency/pin or second output owner is added. Existing clib
ReadFile/getlin retain successful-open read errors, NUL and final-line rules.

Sixteen appended cases (**11 status 0 / 5 status 134**) compare raw file
bytes, streams/status and saved records:

- Item lore **777** versus display **300**, item-to-leaf **701** via school
  **700**, direct skill sheets with display **1001**, nonitem **9000**, and
  absent item magic mapping to **0**. Headers and substitution use the
  display ID, never the sheet ID. Adjacent **@/@@/@@@@** collapse separately;
  comments, leading-space #/$ literals, blank lines, actual embedded NUL,
  retained CR and final nonnewline text match exact C bytes.
- Actual **T_item=4** fixtures reach both legal dollar functions (not merely
  similarly populated location records). Capturable-only and negative
  capturable values, attack-only **12**, defense-only **-3**, missile-only
  **7**, nonanimal exclusion, neither flag, absent item records and
  **item_ghost_warrior=18** with zero combat stats reach the real predicates.
  The output contains the ghost's **(0,0,0)** exception and all three other
  combat tuples in actual item-chain order. Empty expansions are covered.
- ENOENT/ENOTDIR emit the header, unavailable placeholder and exact perror;
  directory fopen succeeds then read fails, producing only the header and
  delivery's trailing blank, no missing-file diagnostic/placeholder.
- Invalid unknown/empty/space-suffixed function names retain **lore.c:43**.
  Immediate output proves the header and earlier body precede the abort.
  A skill cycle emits its header before **use.c:1182**; blocked output fails
  first at **sout.c:158**, before even that path validation.
- Repeated delivery rereads changed input while borrowing the renamed cached
  log inode. **227** bytes retain MASTER **-2**, OUT_LORE **2**, style **5**,
  indent **3/8**, exact blanks and substitutions; `log/200` stays absent.
  These numeric constants and reached bytes, and complete extracted C
  definition equality, are checked in `tmp/issue159-evidence.log`.

At **0.99.0-alpha**, two unchanged real pinned C runs are deterministic and
Go matches **1158 sections / 77449 lines / 1508521 output bytes**. Entire
inherited **1142 sections / 77030 lines / 1499775 output bytes** and
**24967 input lines / 473763 input bytes** remain exact byte prefixes.
New input: **25169 lines / 478453 bytes**. Output SHA-256:
`b312ad0755b50233903f58ed7a07039c5e70e3c2cb1170362646f6e7919e232f`.
Input SHA-256:
`ee32b252987eb9af7bdae83a8ee1b743fd030df5890b431f93a01a8bd73a181b`.
Intentional C trailing spaces and the whitespace-suffixed input are retained.

Four temporary linker controls wrap only the direct sheet adapter at the
expected cycle failure (**section 1156**). A missing-dependency marker with
either return **0** or abort **134** makes default and update modes reject
with harness **2**, before expected status/update acceptance. All recording
hashes remain unchanged; probe source/compiler wrapper and matching captures
are removed. Logs: `tmp/issue159-negative-{0,134}-{default,update}.log`.

Sequential checks pass with explicit `G3_REPO=$PWD/tmp/upstream.git` and
shared run_limited: command refresh/default **1158/77449**, move **135/7158**,
helper **37/1487**, RNG **321**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**,
default turn (equal status/raw streams and **2** cwd entries), bounded vet
and fresh bounded default tests (Olympia **1.812s**). The first fresh test run
found a missing player box in the new NYI-boundary test fixture, not an engine
failure; that fixture was corrected and vet/fresh tests rerun successfully.
Logs: `tmp/issue159-*.log`; matching captures removed.

All ten open issues, including unlabeled ones, are only #159–166 and trackers
#19/#132: no bug/follow-C priority work. No new verified defect or blocking
prerequisite found. #160's explicit real-skill dependency remains deferred as
authorized. Parent strict post-#158 remains red: eight saved-report byte
differences, engines **0**, matching raw streams/paths/types/modes, exit **1**
before 13 lifecycle controls. No new whole-world run or stage-exit acceptance;
G3 #358/#359 remain unchanged/nonblocking, malformed **C136/Go2** is not parity.
No commit, push, issue/checklist edit or parent strict capture change.

Independent parent complete-C review verifies the three task definitions
unchanged exactly once and both recording prefixes byte-exact. The eight
routine tools, bounded vet and fresh bounded tests pass (Olympia **1.744s**);
logs `tmp/parent159-*.log`. Strict whole-world evidence remains post-#158.

## Complete skills, provinces and summary root (#158)

Existing summary.go adds complete unchanged summary.c `summary_skills`
**179–228**, `summary_spells` **231–280**, `summary_provinces` **348–389**
and `summary_report` **438–502**. Main and the existing oracle replace
summaryReportState/summaryProvincesState with the actual renderers; both
copied report.go projections are removed. The root calls the accepted #157
totals, units, men, gold and land helpers exactly once in original order,
then skills, spells and provinces. Original copied unit/skill traversal,
stored knowledge sorting, accounting, ranking and cached FILE/output owners
are reused. SHOW_TOP remains FALSE and `summary_sublocs` remains excluded
by `#if 0`. No later root, new game command, second output layer, source
dependency/pin change, or successful C substitute is introduced. The existing
C extraction already contains every complete definition and its dependencies;
only direct skills/spells adapter dispatch is added.

Nine appended sections (seven **0**, two **134**) compare full raw report
bytes, streams/status, saved records and scratch observations:

- Two unequal regular factions plus system/NPC/silent owners: totals
  **2 players / 4 controlled rows / 3 other characters**, actual characters
  **2/1**, men **13/9**, gold **1144/222**, skills **2/4 versus 0/4**, spells
  **2/2 versus 1/2**, provinces **2/3 versus 1/3**. A prisoner is included;
  duplicate unit and skill rows count only once per distinct skill. Partial,
  don't-know and malformed-known rows do not qualify. Full header, underline,
  all seven rows and trailing blanks retain exact order and inline markup.
- Magic follows the prerequisite school's subtype, not the leaf subtype.
  Direct skills/spells passes reverse the two factions' rankings, then switch
  to all nonmagic and all magic, exercising both zero denominators without
  a slash. Stored descending/duplicate skill rows remain byte-exact.
- All player subtypes allocate in the province pass, including absent regular,
  system, NPC and silent records. Original stored knowledge sorts from
  **9002,9010,99,9000** to **99,9000,9002,9010**; lair/item knowledge is
  retained but not counted as provinces. Repeats clear location scratch.
  Separate root cases cover absent records, zero regular factions, no output
  files, route reset and retained initial indent. A valid noncharacter unit
  row with no CH and a skill with no SK stay unallocated through direct
  skills/spells and the root; this is a raw-list control, not a normal unit.
- Two complete roots borrow the renamed cached inode; `log/200` stays absent
  and the original contains **1400** exact bytes, including changed province
  totals, retained style and the original one-shot second-indent consumption.
  Output-open failure remains **sout.c:158**, and a skill-school cycle remains
  **use.c:1182**, both status **134**. No independent guards/fixes are added.

At **0.98.0-alpha**, complete unchanged real C repeats deterministically and
Go matches **1142 sections / 77030 lines / 1499775 output bytes**. Inherited
**1133 sections / 76217 lines / 1479826 output bytes** and **24697 input
lines / 468406 input bytes** remain exact unchanged HEAD prefixes.
New input: **24967 lines / 473763 bytes**. Output SHA-256:
`3815d51f0ed58c53301d31813b12f749d7fc1d770bc49a140b6e6b08ca79c89d`.
Input SHA-256:
`cc04ab6923b598184344e221572b1e6aede1c139916fe15a048b913a77e83960`.
Complete-definition equality, fixture constants and actually reached bytes
are checked in `tmp/issue158-evidence.log`.

Four temporary linker probes wrap only the direct summary-parts adapter at
the new expected-cycle failure, section **1141**. Emitting a missing-dependency
marker and either returning **0** or aborting **134** makes both default and
update modes reject with harness status **2** before expected-failure
acceptance or recording replacement. Every before/after hash is unchanged.
Logs: `tmp/issue158-negative-{0,134}-{default,update}.log`. An initial probe
attempt itself read req_skill(700) when that box was absent in an earlier
control, and was rejected at section **1138**, status **139** (harness **2**),
without changing the recording. The probe was corrected to check valid_box
first, then all four intended controls reached the marker. This is a
temporary probe error, not an engine defect or accepted negative result;
log `tmp/issue158-probe-initial-invalid-read.log`. Both probe files were
removed before the final real-C refresh/default passes.

Sequential command refresh/default **1142/77030**, move **135/7158**, helper
**37/1487**, RNG **321**, code **43/969** (unchanged older pin), loadsave
**61/3411**, turn-driver **7/152**, phase **5/484**, calendar **2/119**,
schedule **4/1801** and default turn (**2 cwd entries**, raw streams/status)
all pass. Bounded `go vet ./...` and fresh bounded
`go test -count=1 -exec "bash $PWD/tools/lib/run-limited.sh" ./...` pass
(Olympia **1.727s**). Every compat invocation uses explicit isolated G3_REPO;
engines use the existing run_limited owner. Logs:
`tmp/issue158-command{,-update}.log`,
`tmp/issue158-{move,helper,rng,code,loadsave,turn-driver,turn,vet,test}.log`.
Matching captures and temporary probes are removed. Source-only whitespace
checks pass; literal C recording trailing spaces remain unnormalized.
All eleven open local issues, including unlabeled ones, are only #158–166
and the two trackers, with no bug/follow-C. No new defect, prerequisite or
oversized closure was found. No commit, push, issue edit or checklist tick.

Independent parent review verifies all fourteen extracted summary definitions
unchanged exactly once; the same eight-tool/vet/fresh-test tier passes
(Olympia **1.601s**), logs `tmp/parent158-*.log`.
These are routine task checks, not strict whole-world/lifecycle or stage-exit
acceptance. Fresh parent post-#158 strict result has eight saved-report byte differences,
exit **1** before lifecycle, with identical C repeats, both engines **0** and
equal raw streams/paths/types/modes; log `tmp/parent158-strict.log`, captures
`tmp/parent-stage2-turn`. G3 #358 remains nonblocking/unchanged
(malformed **C136/Go2 is not parity**); G3 #359's source-linked native
truncation comment and behavior remain unchanged.

## Summary rankings and totals (#157)

Source-corresponding summary.go implements the active complete summary.c
`rank_comp` **15–22**, `ranking` **33–41**, `out_rank_mine_top` **72–89**,
`out_ranking` **92–119**, `collect_game_totals` **127–159**,
`summary_units` **162–176**, `summary_men` **283–303**, `summary_gold`
**306–320** and `summary_land_owned` **323–345**. Existing u.go supplies
complete `ordinal` **1711–1724**, retaining C's literal **10–19** check
(thus **111st**, not an independent suffix correction). world.go owns the
original ranks/nplayers/ncontrolled/nother scratch fields. Existing accounting,
inventory, copied unit loops, garrison strength/unique ruler lists, sorting
and sout.go's cached output owner are reused. SHOW_TOP remains FALSE;
the dormant `top_rank` display is not invoked or newly implemented.

The summary root, skills/spells/provinces and projection replacement remain
**#158**. Neither main nor summaryReportState/summaryProvincesState is changed
or additionally wired, so the existing persistent effects still run once.
The C oracle already extracts every complete definition and dependency above.
The adapter only dispatches direct calls and observes original scratch fields;
no new source dependency, pin, rewritten C model or output layer is added.
Exact complete-definition equality, fixture constants and actual reached
branches are checked in `tmp/issue157-evidence.log`.

Nine appended sections (eight status **0**, one **134**) compare raw report
bytes, streams/status, saved records and original scratch state:

- Five asymmetric regular factions rank **4000,1000,1000,0,-5** as
  **1st,2nd,2nd,4th,5th**, excluding system/NPC/silent values even when
  larger. Zero/positive/negative denominators, all-five ties and repeated
  membership changes exercise scratch-list clearing without evaluating top.
- Controlled totals use actual unit-list rows, not character ownership:
  duplicated and foreign listed units give **2 players / 4 controlled /
  3 other**, while actual owner counts are **2,1,1,1,1**. Repeated calls
  clear character marks and counters, retaining unrelated player temps.
- Zero regular factions create no rank list or output, and leave absent
  player records absent. One regular faction with no player record is
  counted without allocation; only real output subsequently allocates it.
- Men totals **1231/999/4000/20/0**, gold **2000/2000/9000/7/0** include
  prisoners, nested/foreign and nonregular characters, nonfighter men and
  signed nonzero man flags, but exclude player/location holdings and
  nonman/missing-item records. Player temps reset each pass; character
  temp **41** survives. Gold ties produce equal first ranks.
- Actual garrisons reach strengths **9,9,10,10**, then **10** after one
  added soldier. Shared pledged rulers count each faction only once;
  repeated same-faction rulers are deduplicated. An empty castle contributes
  nothing. Regular and nonregular ruler counts rise **1 → 2**, with an
  unrelated regular faction remaining zero.
- Output after renaming the cached log appends to the original inode;
  `log/200` remains absent and `log/rank-original` contains **174** exact
  bytes. The failed output open retains **sout.c:158**, not a new assertion.

At **0.97.0-alpha**, complete unchanged real C repeats deterministically and
Go matches **1133 sections / 76217 lines / 1479826 output bytes**. All inherited
**1124 sections / 75560 lines / 1466335 output bytes** and **24431 input
lines / 463271 input bytes** remain exact unchanged HEAD prefixes.
New input: **24697 lines / 468406 bytes**. Output SHA-256:
`efd25aa486381603c682cbb58f9bb6da12d975e78dfcf173d0848c001ac6fff9`.
Input SHA-256:
`453a5347380170032686c5424fc83349e0a93fb0029589b93980053c63d79a34`.

Four temporary linker probes wrap only the direct summary adapter at expected
failure section **1132**, emit the rejected dependency marker and exit **0**
or abort **134**. Both default and update modes reject with harness status
**2** before expected-failure acceptance or recording replacement. Every
before/after recording hash stays unchanged. Logs:
`tmp/issue157-negative-{0,134}-{default,update}.log`. Probe source/compiler
wrapper and matching/provisional captures are removed before final refresh.

Final real-C command refresh/default pass **1133/76217**, including native
traversal on unchanged shared input. Sequential move **135/7158**, helper
**37/1487**, RNG **321**, code **43/969** (its unchanged older pin), loadsave
**61/3411**, turn-driver **7/152**, phase **5/484**, calendar **2/119**,
schedule **4/1801** and default turn (**2 cwd entries**, raw streams/status)
all pass. Bounded `go vet ./...` and fresh bounded
`go test -count=1 -exec "bash $PWD/tools/lib/run-limited.sh" ./...` pass
(Olympia **1.532s**). All use explicit isolated G3_REPO and the existing
run_limited owner. Source-only whitespace checks pass; literal C recording
trailing spaces remain unnormalized. Matching captures are removed. Logs:
`tmp/issue157-command{,-update}.log`,
`tmp/issue157-{move,helper,rng,code,loadsave,turn-driver,turn,vet,test}.log`.
Independent parent review verifies all ten complete C definitions unchanged
exactly once; the same eight-tool/vet/fresh-test tier passes (Olympia
**1.618s**), logs `tmp/parent157-*.log`.
These are task checks, not a whole-world/lifecycle or stage-exit acceptance.
The latest parent strict post-#156 result remains identical C repeats,
both engines **0**, equal raw streams/paths/types/modes, eight saved-report
byte differences, exit **1** before lifecycle; parent captures are untouched.
No new defect, prerequisite or oversized closure was found. The parent supplied
the all-twelve-open-issue audit, including unlabeled issues and no bug/follow-C.
G3 #358 remains unchanged/nonblocking (**C136/Go2 is not parity**), and
G3 #359's source-linked native truncation remains unchanged. No commit, push,
issue change or checklist tick is part of this task.

## Player-report composition, storms and garrisons (#156)

Existing report.go now implements complete `player_report_sup` **1188–1212**
and `player_report` **1323–1350**. The root replaces `playerReportState` in
both main and the existing oracle: support, accepted unit summary (#154),
location-stack and capacity tables (#155), storms, garrisons and final blank
run exactly once in original order. MASTER/banner routing resets to zero;
the ordinary cached FILE owner, loops, weights and output sorting are reused.
There is no copied state projection or parallel renderer.

Existing c1.go adds complete `print_hiring_status` **904–921** and
`print_unformed` **924–944**, storm.go adds `storm_report` **326–362**, and
garr.go adds `garr_own_s` **813–857**, `garrison_summary` **860–906**, plus
its original ownership predicate `may_rule_here` **142–165**. The C oracle
already extracted all these complete definitions and their real dependencies.
The adapter adds only direct report calls and NP-counter/fast-study fixture
fields. No new game commands, C source dependency, pin, modeled definition
or successful rendering substitute is introduced. Complete-definition equality
against isolated immutable source, constants and actually reached branches
are checked in `tmp/issue156-evidence.log`.

Eleven appended sections (eight status **0**, three **134**) compare complete
raw report bytes/streams/status and real state/save observers:

- Regular, NPC, system and silent roots include absent-record controls.
  System skips support but still runs tables; silent root skips everything.
  Direct system support does not allocate, nor do absent-PL unformed/NPC
  hiring checks. A direct silent support call still renders, as C does.
- NP gained/spent, fast-study **1/2/negative** and six unsorted unformed IDs
  preserve original wording, spacing and the first-five truncation. Hiring
  covers turns **0/7/8/15/32767/-32768**, including zero countdown becoming
  eight. No independent grammar, scheduling or signed-boundary fix is made.
- Owned fog/wind/rain, foreign/silent ownership, absent summoners, nested
  province lookup, negative and comma-formatted strengths retain actual storm
  iteration. A storm without a province prints location **0** successfully;
  there is no invented location assertion.
- Garrison output order **1102,1100,1101** follows shared sorting, including
  ID ties. Forwarding **-1/0/123**, recursive **21** figures versus the
  garrison's own **26** maintenance cost, tax **57**, absent-castle filtering
  and zero-strength ownership are exercised. Seven rulers display
  `1000 1001 1002  ... 1006`; five, three and one retain their full IDs.
- Full asymmetric faction composition includes listed foreign/nested units,
  actual unit/stack/capacity tables, NP/unformed/fast-study, owned storm and
  shared garrison rulers, system/NPC/silent controls and saved player,
  character, inventory and storm records. Repeating after renaming the cached
  file appends to the same inode: **3856 bytes** in `log/faction-original`,
  with `log/200` still absent. Style bytes and the foreign-faction footnote
  are compared literally. Failures retain **c1.c:909** and **sout.c:158**.

At **0.96.0-alpha**, the extended recording is **1124 sections / 75560 lines /
1466335 output bytes**. The inherited **1113 sections / 74785 lines /
1445511 output bytes** and **24146 input lines / 457547 input bytes** are
exact unchanged HEAD prefixes. New input: **24431 lines / 463271 bytes**.
Output SHA-256:
`0bb63affc1deda8e8a109ebdeab22773ac4c031016f72cd45d7b2410e78201a5`.
Input SHA-256:
`1631629b551d0979b5c34c67eb7bbe4aaa3487760a1d5dd26716ee23b323568a`.

Four temporary linker probes wrap only `print_hiring_status` at expected
assertion section **1120**, emit the rejected dependency marker and exit **0**
or abort **134**. Both default and update modes reject with harness status
**2** before expected-status acceptance or recording replacement. Every
before/after recording hash remains the output SHA-256 above. Logs:
`tmp/issue156-negative-{0,134}-{default,update}.log`. Temporary compiler
wrapper/probe source are removed before the final real-C refresh/default.

Final complete-real-C refresh and default both repeat C deterministically
and match Go **1124/75560**, including live native traversal on unchanged
shared input. Sequential move **135/7158**, helper **37/1487**, RNG **321**,
code **43/969** (its deliberately unchanged older pin), loadsave **61/3411**,
turn-driver **7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**
and default turn pass. The default turn compares raw streams/status and
**2 cwd entries**, not a whole-world turn. All tools receive explicit isolated
G3_REPO; engines and fresh test executables use the shared run_limited owner.
Bounded `go vet ./...` and fresh bounded
`go test -count=1 -exec "bash $PWD/tools/lib/run-limited.sh" ./...` pass
(Olympia **2.072s**). Logs: `tmp/issue156-command{,-update}.log` and
`tmp/issue156-{move,helper,rng,code,loadsave,turn-driver,turn,vet,test}.log`.
Matching captures are removed; original C whitespace is retained.

No new defect, prerequisite or oversized closure was found. The parent's
all-thirteen-open-issue audit includes unlabeled issues and no bug/follow-C
work. G3 #358 remains nonblocking and unchanged: malformed **C136/Go2** is
not passing parity. G3 #359's source-linked comment and observed fifteen-byte
truncation are unchanged. No later summary/GM/lore/include/writer roots,
whole-world/lifecycle or stage-exit acceptance are claimed. Independent parent
review confirms twelve complete original C definitions unchanged exactly once;
the eight-tool tier, bounded vet and fresh tests pass (Olympia **1.461s**),
logs `tmp/parent156-*.log`. Fresh parent strict post-#156 evidence has identical
C repeats, both engines **0**, equal raw streams and paths/types/modes, but
eight saved-report byte differences and exit **1** before lifecycle. Log:
`tmp/parent156-strict.log`; captures: `tmp/parent-stage2-turn`.

## Location stack and capacity tables (#155)

Existing report.go implements complete unchanged report.c `loc_ind_s`
**1062–1073**, `loc_stack_catchup` **1076–1086**, `loc_stack_rep_sup`
**1092–1115**, `loc_stack_report` **1118–1185** and `stack_capacity_report`
**1215–1320**, retaining the original foreign-faction explanation scratch.
All five C definitions and their real dependencies were already extracted;
the existing adapter adds only direct calls, trivial static-helper wrappers
and raw temp observers. Complete-definition equality against the immutable
bare-repository source is verified in `tmp/issue155-evidence.log`. No new C
source dependency/pin, rewritten definition, successful missing dependency,
second output owner, later player-report root or main wiring is introduced.
Actual `determineStackWeights`, names/codes, unit loops, output sorting and
cached FILE routing are reused, as is accepted #154 `unitSummary`.

Eighteen appended sections (fourteen status **0**, four **134**) compare full
raw logs/streams/status, real saved records, original lists and scratch:

- Empty/missing-PL and prisoner-only roots clear location, ship and character
  temps before returning without output/allocation; item temp **61** survives.
  Capacity does not clear temps or exclude prisoners. Direct system/silent
  output allocates through the existing output owner; immediate output leaves
  the absent PL absent. Neither root sorts the stored faction units.
- Region/plain/city/castle and nested ship paths retain exact **0/2/4/6**-space
  indentation, ancestor-first catch-up and **-1** markers. Nonzero ancestor
  temp stops recursion, including a positive unit-link marker. Reverse list
  building leaves **9200 -> 1000 -> 1001 -> 1002 -> 0**; prisoners and unlisted
  descendants remain zero. Unit and location output-order ties use IDs.
  Foreign listed descendants have the star/footnote; unlisted descendants do
  not gain rows. Relocation/repeats rebuild links and reset explanation state.
- Capacity follows stored roots including duplicates and recursively weighs
  own, foreign and prisoner descendants. Distinct weights/caps actually print
  walking/riding/flying **33/78/55%**, then **186/365/311%** overloaded; the
  recursive leader totals **158**, with **32/98/131%** ratios. Zero capacity
  emits blank cells, including after nonblank rows. Exact fullness prints
  **0 100%** in all three columns, negative load prints **241 -141%**, and
  negative caps retain the real self-carry/no-extra-capacity semantics.
- All three **999%** cells print `(899) 999%`; weight **1000** over capacity
  **100** reaches all three greater-than-999 branches and prints `(900)  -- `.
  A raw negative quantity stays stored but is skipped by the real inventory
  iterator. Custom noble items and native C signed product/ratio overflow
  retain original comma formatting; no independent arithmetic fix is made.
- Accepted summary plus both new roots share a renamed cached inode:
  **892 bytes** reach `log/original`, with `log/200` still absent. The new
  renderers do not flush, close or reset the cache. Original assertions remain
  at **sout.c:158** before linked-list/weight work, **u.c:1480** for an invalid
  weight root, and **u.c:388** for invalid ancestor depth after the header.

Pinned constants and actually reached bytes are checked in the evidence log:
player/char/loc/item/ship **1/2/3/4/9**, plain/city/galley/roundship/castle/region
**3/14/30/31/32/50**, system/silent **48/62**, CHAR_FIELD **6**. No new defect,
prerequisite or oversized closure was found; the all-open local audit includes
all fourteen issues, including unlabeled ones, with no bug/follow-C work.

At **0.95.0-alpha**, complete unchanged real C runs twice deterministically
and Go matches **1113 sections / 74785 lines / 1445511 output bytes**. The
entire inherited **1095 sections / 73930 lines / 1423874 output bytes** and
**23766 input lines / 449169 input bytes** remain exact HEAD prefixes.
New input: **24146 lines / 457547 bytes**. Output SHA-256:
`d0a1306d9a95ac8fc458e197e035936bbe74dc8b4e6c095c760aacc450a31559`.
Input SHA-256:
`508d320c7627f422d5749da882623b7bd01b47484f6e108cd17b02ce4832c907`.
Live native traversal also matches C/C/Go with shared input unchanged. Logs:
`tmp/issue155-command{,-update}.log` and `tmp/issue155-evidence.log`.

Four temporary linker probes wrap only `loc_stack_report` at expected-failure
section **1108**, emitting the rejected dependency marker and returning **0**
or aborting **134**. Both default and update modes reject with harness status
**2**, before expected-status acceptance or recording replacement. Each log
records unchanged before/after SHA-256
`80cbf285bd16513b363bb94c708e704a054d221c5f7d2a0de6cd960a9924c040`
(the initial 1112-section recording, before appending the exact-full control).
Logs: `tmp/issue155-negative-{0,134}-{default,update}.log`. Temporary source
and compiler wrapper were removed before the final real-C refresh/default.

Sequential command refresh/default, move/helper/RNG/code/loadsave/turn-driver/
default-turn, bounded vet and fresh bounded default tests pass with explicit
isolated G3_REPO and run_limited for engine/test executables. Logs:
`tmp/issue155-{move,helper,rng,code,loadsave,turn-driver,turn,vet,test}.log`.
Fresh Olympia defaults: **1.846s**. Matching captures removed; literal C
whitespace retained. Independent parent review verifies all five complete C
definitions unchanged exactly once; the same eight-tool/vet/fresh-test tier
passes (Olympia **1.586s**), logs `tmp/parent155-*.log`. No whole-world/lifecycle or
stage-exit acceptance; latest strict evidence remains the post-#153 eight
report-byte differences, exit **1** before lifecycle. Parent strict captures
are untouched. G3 #358 stays nonblocking/unchanged: original division/Go panic,
malformed **C136/Go2**, is not passing parity. G3 #359's observed fifteen-byte
truncation and existing source-linked comment remain unchanged.

## Complete faction unit summaries (#154)

Existing report.go adds complete unchanged report.c `sum_fighters`
**801–822**, TRUNC_NAME and the original seventeen-entry `stupid_words`
table **825–848**, `strip_leading_stupid_word`, `stupid_word`, `prev_word`,
`summary_trunc_name` **851–925**, and `unit_summary_sup` / `unit_summary`
**937–1059**. Required u.c `knum` **1694–1708** is ported in existing u.go.
All these complete C definitions, tables and statics were already extracted;
only trivial direct-call wrappers and raw scratch observers are added to the
existing adapter. No C definitions are rewritten and no successful missing
dependency substitute is added. Existing inventory/unit loops, names/codes,
items/skills, stack/subloc, output/cache and allocation owners are reused.
The later player-report root and its projection are deliberately unchanged;
this task does not wire an additional effect pass into main.

Twelve appended sections (ten status **0**, two **134**) compare real raw
logs/streams/status, stored records/lists, scratch totals and allocation flags:

- Empty roots clear **all character temps** before output and allocate PL
  after the first blank, even for system/silent players and immediate output.
  One unit emits no totals; a missing CH/CM remains missing. Multiple units
  follow the original stored list, including duplicates, without sorting or
  filtering prisoners. Repeated roots reset totals; direct helpers accumulate.
- All six loyalties produce original `ucofns` prefixes. Health **100/-1**
  overrides sick; **37** healthy/sick, **0/-2**, behind **127/-128**, magician
  **1/-1**, nonmagician with stored aura and aura **12345/-99** exercise
  minimum-width (not truncation) columns. Captives show both ` ?? ` fields;
  stacked noncaptives show their real parent and enclosing sublocation.
- Original knum boundaries **0/9998/9999/999999/1000000/2147483647** and
  negative values are reached through real loyalty/inventory/totals. The
  **9999 -> 9k** discontinuity and signed 32-bit totals are preserved. Seven
  ordered rows, including a duplicate, total **29959/1000000/999999/10000/0**.
  Overflow controls yield gold **-2147473650**, fighters **-2146473651**.
- Real fighter predicates distinguish man/nonman, all-zero, all-one,
  attack-only/defense-only/missile-only **2**, negative stats, missing item
  record, duplicate quantities, invalid item IDs and zero/negative quantities.
  Fighters total **102**, then **139** after giving the man-marked ghost
  warrior missile **2**; its ghost identity alone does not pass max-stat >1.
  Inventory order, invalid/raw entries and records are unchanged.
- Names of **14/15/16** bytes and spaces immediately before/at width **15**,
  all seventeen leading and trailing table entries, case-insensitive matching,
  a nonword prefix and multiple spaces execute the original pointer/byte
  algorithm. Short `The Short Name` stays untouched. Only one leading word
  is stripped: `The the Extraordinary Warrior` becomes **the**. Trailing
  words strip repeatedly, but `Alpha  of Extraordinary` retains **Alpha **,
  including its space; an all-space tail is not normalized. No convenient
  split/join/trim algorithm replaces these quirks.
- Borrowed cached output survives rename: **181** bytes reach `log/original`,
  `log/200` remains absent, and the renderer does not flush/close/reset it.
  Open failure asserts at original **sout.c:158** before explicit PL access;
  a deliberately absent gold entity asserts at **u.c:1832**, not a missing
  C-definition marker. The guard still rejects dependency markers first.

Numeric fixture constants were checked in pinned oly.h: player/char/loc/item
**1/2/3/4**, plain **3**, system/silent **48/62**, gold/peasant/worker/sailor
**1/10/11/19**, ghost warrior **18**, loyalties **0–5**, CHAR_FIELD **6** and
TRUNC_NAME **15**. Above branches and bytes were inspected in actual C output.
Source portability caveat retained rather than fixed: the original truncator
can call `prev_word(s, NULL)` for an unbroken long first word, whose `t > s`
pointer comparison is not portable C. Native pinned C repeats deterministically
and returns the first fifteen bytes in these cases; Go preserves that observed
result. Likewise the signed-overflow probes record this native C build's
behavior, not a claim about every optimizing C compiler. No observed blocking
defect or new prerequisite was found. Parent confirmed the NULL comparison in
current upstream and filed [G3 #359](https://github.com/playbymail/olympia-g3/issues/359),
assigned to the owner; the port links it without changing observed behavior.
No upstream code was changed.

At **0.94.0-alpha**, complete unchanged real C runs twice deterministically
and Go matches **1095 sections / 73930 lines / 1423874 output bytes**. Entire
inherited **1083 sections / 73273 lines / 1406675 output bytes** and **23438
input lines / 441216 input bytes** are exact HEAD prefixes. New input totals
**23766 lines / 449169 bytes**. Output SHA-256:
`cd3e5c556e5a2e71d5f8a6eb051e149a266e96e6007001e54c8513b708614b5f`.
Input SHA-256:
`a003ba328dfdd2b7ab5718a1eeed84b5c07e3ae6f2930e27419983b0c81d4edd`.
Live native traversal also matches C/C/Go with shared input unchanged. Logs:
`tmp/issue154-command{,-update}.log`.

Four temporary linker probes wrap unit_summary only when player **200**
has no PL and its log path is a directory, at new expected-failure section
**1093**. Marker plus return **0** or abort **134**, default and update modes,
all reject with harness status **2** before expected-status acceptance or
recording replacement. Each log records the identical before/after output
hash above: `tmp/issue154-negative-{0,134}-{default,update}.log`. Probe source,
object and compiler wrapper are removed before the final normal comparison.

Sequential command refresh/default, move/helper/RNG/code/loadsave/turn-driver/
default-turn, bounded vet and fresh bounded default tests pass with explicit
isolated G3_REPO and run_limited for engine/test executables. Logs:
`tmp/issue154-{move,helper,rng,code,loadsave,turn-driver,turn,vet,test}.log`.
Olympia fresh defaults: **1.896s**. Matching captures are removed. Source-only
whitespace validation passes; literal raw C recording and the all-space-tail
input control intentionally retain trailing spaces. Independent parent review
and the same eight-tool/vet/fresh-test tier pass (Olympia **1.559s**); logs
`tmp/parent154-*.log`. No full-world, lifecycle, later report-root or stage-exit
acceptance is claimed. Latest strict evidence remains the post-#153 eight
report-byte differences, exit **1** before lifecycle; the deferred root is not
wired by this task. `tmp/parent-stage2-turn` is untouched. G3 #358
remains unchanged/nonblocking; malformed **C136/Go2** is not passing parity.

## Complete player banner and account footer (#153)

Existing report.go implements complete unchanged report.c `player_banner`
**1645–1722**, `report_account_sup` **1725–1791** and `report_account`
**1794–1815**. Main and the existing command oracle replace the two state
projections with those renderers, once in the original order. No second
effect pass remains. The existing player/calendar/name, html/style and
sout.go cached FILE/output owners supply the entire closure. The -w option
is retained in world state for the banner's image path; oracle life open/save
updates that same flag as its existing C adapter does. No new source
dependency, pin, later report root or handler is added.

The existing fixture adapter adds direct `account_sup`, banner win/title and
rules/Times URL setters only. Complete original C definitions were already
extracted; extraction and the dependency-marker guard are unchanged. No
successful rendering substitutes or modeled C expectations are introduced.
The account command, `/tmp/oly-acct.tmp` read/unlink and account-summary text
are all **#if 0** at the pin: no account input read is active. External
`charge_account` **1818–1856** / -A remains explicitly unsupported.

Eleven appended sections (nine status **0**, two **134**) compare full raw
logs/streams/status, saved records and routing/style/cache observations:

- Banner emits for regular **49**, system **48**, savage **51**, NPC **52**
  and plain subtype **0**; silent **62** has neither output nor allocation.
  Missing nonsilent PL records allocate, including system players. The
  compuserve warning retains its original missing space at `you.To`, plus
  indent **+3** and the following blank line. Names/title are not escaped.
- Turns **0**, every month **1–8**, **9/16/17**, signed-short endpoints
  **32767/-32768**, all eight real season names and year rollover are
  observed. Turn zero clamps the month to **0**, yielding year **1**;
  negative remainder clamps too, without changing the negative year.
  The **report.c:1706** month assertion is retained; remainder plus clamp
  makes it unreachable for these valid stored turn values. Both literal
  `../head.gif` and -w `head.gif` image tags match, including attributes.
- Footer root emits/allocates only regular players, not NPC/system/silent
  or plain players. Quotes, angle brackets, ampersands, spaces and empty
  URLs are emitted raw; the three spaces between links are preserved.
  HTML lines carry STYLE_HTML **2**, while the trailing blank and later
  style probes use **0** even after STYLE_PREV. With no regular players,
  the loop never resets style; the existing HTML style **2** survives.
- MASTER **-2** / OUT_BANNER **0**, prior nondefault route, indent/second
  indent/day and style are observed exactly. Both roots reset routing to
  **0/0**, not the previous route. Direct account_sup preserves its caller's
  route **-2/44**, performs no subtype filtering, and can allocate PL even
  on a character recipient through the existing output owner.
- Immediate footer output does not allocate missing PL; banner still does
  at its explicit p_player after initial HTML. No log file opens in either
  immediate path. First HTML on a blocked ordinary output fails at original
  **sout.c:158** before the explicit banner allocation or later footer
  output. Source retains grab_fp/open-before-p_player ordering; failure
  controls compare the actual diagnostic/assertion, not post-abort state.
- Buffered banner then footer share the renamed cached inode: `log/200`
  stays absent and the original file receives **486** raw bytes, while
  close/reset permits a new **203-byte** footer log. The other faction's
  retained log contains **689** bytes. Observers flush explicitly; neither
  renderer closes, flushes, resets or duplicates the accepted cache.

Numeric fixtures were verified against pinned oly.h/sout.h: player/character
**1/2**, subtypes above, NUM_MONTHS **8**, STYLE_TEXT/HTML/PREV **1/2/-1**,
MASTER/OUT_BANNER **-2/0**, and signed-short stored turn arithmetic. Actual
recording branches and bytes were inspected, not inferred from case names.

At **0.93.0-alpha**, unchanged real C runs twice deterministically and Go
matches **1083 sections / 73273 lines / 1406675 output bytes**. Entire
inherited **1072 sections / 72518 lines / 1388090 output bytes** and **23243
input lines / 437538 input bytes** remain exact byte prefixes against HEAD.
New input totals **23438 lines / 441216 bytes**. Output SHA-256:
`6b2a5504ef21fc7dd2eb956e2a98eb9dd4725859ef7c68b4e9471f08ea93a14a`.
Input SHA-256:
`63e67a732891bdb21774efc3894488e0e5f5a3d42aeeb492eb6c848df373d935`.
#168 live native traversal still matches raw C/C/Go with shared input
unchanged. Logs: `tmp/issue153-command{,-update}.log`.

Four temporary linker probes wrap report_account only when player **200**
has no PL and its log path is a directory, at new expected-failure section
**1081**. Marker plus return **0** or abort **134**, normal and update modes,
all reject with harness status **2** at the dependency guard BEFORE expected
failure/status acceptance or recording replacement. Each negative log records
the identical before/after output hash above. Probe source/compiler wrapper
are removed before the final ordinary comparison. Logs:
`tmp/issue153-negative-{0,134}-{default,update}.log`.

Sequential command refresh/default, move/helper/RNG/code/loadsave/turn-driver/
default-turn, bounded vet and fresh bounded defaults pass with explicit
isolated G3_REPO for every tool and run_limited for engine/test executables.
Counts: move **135/7158**, helper **37/1487**, RNG **321**, code **43/969**,
loadsave **61/3411**, turn-driver **7/152**, phase **5/484**, calendar **2/119**,
schedule **4/1801** (the latter scheduling checks are not gameplay acceptance).
Default turn compares status/raw streams and **2 cwd entries**, not a world
turn. Logs: `tmp/issue153-{move,helper,rng,code,loadsave,turn-driver,turn,vet,test}.log`.
Matching captures and temporary probes are removed; literal C whitespace is
retained. No new defect/prerequisite or oversized closure was found.

Independent parent review confirms all three complete C definitions unchanged
and present exactly once. The eight-tool tier, bounded vet and fresh tests
pass (Olympia **1.846s**); logs `tmp/parent153-*.log`. No later report root,
full-corpus or stage-exit acceptance. Fresh parent post-#153 strict result
remains a red gate: identical C repeats, C/Go **0**,
equal raw streams and paths/types/modes, eight saved-report byte differences,
exit **1** before lifecycle; log `tmp/parent153-strict.log`, captures
`tmp/parent-stage2-turn`. G3 #358
stays nonblocking and unchanged: original division/Go panic, malformed
**C136/Go2**, is NOT passing parity.

## Complete order templates (#152)

Existing order.go adds complete unchanged order.c `orders_template_sup`
**370–454**, `orders_other` **457–505**, `orders_template` **508–547** and
`list_order_templates` **550–570**. Main and the existing command oracle call
the renderer once instead of listOrderTemplatesState; that copied projection
is removed from report.go. The real player/unit loops, order queues, commands,
names, loyalty/number formatting and sout.go cached output owner are reused.
No handler, later report root, new source dependency or pin change is added.

The existing report fixture adapter adds a direct template call and `saved`
input action. A trivial wrapper includes complete unchanged io.c and exposes
its existing static `read_boxes`; Go uses existing readBoxes. Existing
`lore_file` creates actual saved PL/CO fixtures. Passwords, broken-mailer
flags and command lines/states are therefore read by the actual saved-box
parsers, not successful loader substitutes or new gameplay handlers. The
original complete order template definitions were already extracted by this
oracle; extraction, real loops/queue owners and missing-dependency guard stay
unchanged. No compared C function is rewritten.

Eleven appended sections (ten status **0**, one **134**) compare full raw
template logs, output/frame state, saved records and retained queue contents:

- Missing PL and absent/allocated-empty/nonempty passwords, embedded quotes
  and spaces, broken mailer **0/1/-1**, repeated templates and allocation
  after the initial password read. No escaping or independent mail fix.
- Autoquit **5/4/3/1/0/-1**, ages **1/2** and immediate missing-record age
  **0**. Integer half thresholds, singular/plural and disabled access match.
- Own contract **49/50**, fear **1/2**, prisoners' name-without-location,
  ordinary name/location, faction queues, multiple orders and raw comments.
  Direct character recipients reach the separate contract warning at **49**,
  suppress it at **50** and for fear/foreign recipients; prisoner output is
  suppressed by the existing output owner, not by skipping the template.
- Saved RUN wait **-1/0/1/12001** reaches still~executing, zero/plural,
  one/singular and comma-number output; LOAD prints the stored line and DONE
  omits it. The stored command and queued follow-up remain unchanged.
- Foreign queued orders print without names/current command/loyalty warnings.
  Empty foreign LOAD/RUN/**ERROR** queues still cause empty templates; absent
  and DONE commands do not. The foreign heading is emitted only once, and
  only when a template survives. Faction-self, owned-unit, invalid-box and
  **T_deadchar 8** entries are filtered at their original positions.
- System **48** and silent **62** missing records remain absent; regular
  **49**, savage **51**, NPC **52** and plain subtype **0** receive templates.
  MASTER uses OUT_TEMPLATE **5**, then resets routing to **0/0**, preserving
  indent/second-indent/style/unit fields. Immediate mode allocates a missing
  faction only when autoquit explicitly accesses it.
- Buffered repeats keep the renamed FILE/inode: `log/200` stays missing while
  `log/original` receives **250 raw bytes**, then close/reset permits a new
  **140-byte** `log/200`. Observers explicitly flush; rendering itself does
  not close or reset cache ownership. A direct nonplayer PL argument reaches
  the original **order.c:464** assertion after earlier output.

Numeric fixtures above were verified against pinned oly.h/sout.h, including
LOY_contract/fear **1/3**, DONE/LOAD/RUN/ERROR **0/1/2/3**, player/character
**1/2**, deadchar **8**, and all subtype/routing constants. The actual new
recording was inspected to verify reached branches, not just case names.

At **0.92.0-alpha**, unchanged real C runs twice deterministically and Go
matches **1072 sections / 72518 lines / 1388090 output bytes**. All inherited
**1061 sections / 71543 lines / 1368616 output bytes** and **23007 input
lines / 432470 input bytes** remain exact byte prefixes against HEAD.
New input totals **23243 lines / 437538 bytes**. Output SHA-256:
`c76cb77668157edfb55f984dd5682ab252af3dae9fda14c113967641ef84a99a`.
Input SHA-256:
`9ff4678da1296821d37b4bd1c81a7d98ef56cd917e9ffed4b5f2bfd1fd2e13cc`.
#168 live native traversal still matches raw C/C/Go with shared input
unchanged. Logs: `tmp/issue152-command{,-update}.log`.

Four temporary linker probes wrapped only the direct orders_template call
with nonplayer pl **1000** at expected-failure section **1071**. Marker plus
return **0** or abort **134**, in normal and update modes, all rejected with
harness status **2** at the dependency guard, BEFORE expected-failure/status
acceptance or recording replacement. The recording hash stayed unchanged in
all four runs. Probe source/compiler wrapper were removed before the final
ordinary comparison. Logs: `tmp/issue152-negative-{0,134}-{default,update}.log`.

Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded defaults pass, using explicit isolated G3_REPO
for every tool and run_limited for engines and test executables. Counts:
move **135/7158**, helper **37/1487**, RNG **321**, code **43/969**, loadsave
**61/3411**, turn-driver **7/152**, phase **5/484**, calendar **2/119**,
schedule **4/1801**. The default turn compares raw status/streams and two cwd
entries; it runs no world turn. Olympia's fresh default suite took **1.384s**.
Logs: `tmp/issue152-{move,helper,rng,code,loadsave,turn-driver,turn,vet,test}.log`.
Source-only diff whitespace check passes; literal C saved-record trailing
whitespace is retained. Matching captures and temporary probes are removed.

Independent parent review verified all four complete C definitions unchanged
and present exactly once. The same eight-tool tier, bounded vet and fresh
tests pass (Olympia **1.485s**); logs `tmp/parent152-*.log`.
No new defect/prerequisite or upstream blocker found. No later root,
full-corpus or stage-exit acceptance. Fresh parent strict post-#152
result has deterministic C repeats, C/Go **0**, equal raw streams and
paths/types/modes, eight saved-report byte differences, exit **1** before
lifecycle; log `tmp/parent152-strict.log`, captures `tmp/parent-stage2-turn`.
G3 #358 is nonblocking and unchanged; its division/Go panic and
malformed **C136/Go2** are NOT passing parity.

## Complete character-report root (#151)

Existing report.go implements complete unchanged report.c `char_rep_sup`
**681–710** and `character_report` **713–743**. Main and the command oracle
call the renderer instead of characterReportState. The copied projection
and now-unused reportAttitudeState are removed; no double effect pass.
Existing buy.go adds complete `list_pending_trades`, buy.c **1123–1158**,
already extracted unchanged in the real-C oracle. No new C dependency or pin.

The root retains the actual subordinate order: location/loyalty, stack,
health/combat/misc, optional banner, lazy pledges, optional magic, prisoners,
attitudes, known/partial skills, inventory, capacity, item skills, pending
trades, final blank. Every subordinate reuses its accepted owner and normal
out/wiout/cache routing. Silent owners skip before any output or collection.
Prisoners temporarily clear their flag only for the original two output
lines, then restore it; they never run char_rep_sup. Only show_pledged's first
eligible call collects backlinks across all characters. CH allocations,
retained allocated-empty CM, original inventory/attitude sorts, copied skill
sorts, owner record allocation and repeats remain owned by those routines.

Eleven appended sections compare complete actual C twice with Go: eight
status **0**, three **134**. The asymmetric root prints every subordinate,
including aura **11 (9+2)**, prisoner health **42**, Oath-17, shroud/vision,
concealment, banner, all three attitudes, partial **0/3**, inventory total
**71**, distinct land/ride/fly ratios **15/20/11%**, item study/use and stored
SELL-before-BUY offers. It observes actual faction logs, raw lists, skill
rows and save_box bytes before/after repeated roots; changing a forward
pledge does not recollect the existing **1001/1002/1003** backlinks. Separate
silent and prisoner-only controls verify ready **0**, unsorted inventory and
attitudes, untouched missing owner records for silent players, original
prisoner flags, then ready **1** only after the first eligible root.

Garrison output actually reaches both castle owner and pledge ruler logs
using OUT_GARR **6** and the garrison unit field. A missing-CH character
acquires CH at lazy collection, while a preallocated backlink list retains
bare CM after clearing. The MASTER fixture opens character **1000**'s FILE
before renaming it; after repeated complete roots the observer reads **1003
raw bytes**, including the prefix, from that same cached inode. Prisoner
**1001** is observed separately. A
different-recipient char_rep_sup control preserves who versus num semantics.

Pending trades preserve stored order and records, skip PRODUCE/CONSUME,
nonpositive quantities and invalid items, and print **2147483647** quantity,
**-123456789** price and positive-quantity zero price without market_report's
affordability division. This is the unchanged pending renderer, NOT a fix
or claimed parity for malformed market data. Invalid owner asserts at
buy.c **1129** before output; complete-root controls preserve report.c
**504** and use.c **1358** assertions in original call order. Constants were
verified against pinned oly.h: player/character/location/item/skill **1–5**,
plain **3**, castle **32**, magic **43**, regular **49**, region **50**,
silent **62**, garrison **64**, BUY/SELL **1/2**, hide_self **638**.

At **0.91.0-alpha**, final unchanged-C refresh/default comparison passes
**1061 sections / 71543 lines / 1368616 output bytes**. Accepted **1050
sections / 70245 lines / 1339950 output bytes** and **426753 input bytes**
remain exact unchanged prefixes against HEAD. New input: **23007 lines /
432470 bytes**. Output SHA-256:
`cb7ee86183cc0070d30a6fbd7c2fc9c4e9fc5474a858003d4149482c9c6f9cee`.
Input SHA-256:
`b2da676cbdbd92c6129887b35629ab48f7b8159b1a1baf743ca239d384bc9d2c`.
#168 live native traversal still matches raw C/C/Go over unchanged shared
input. Logs: `tmp/issue151-command{,-update}.log`.

Four temporary linker probes wrap only the direct list_pending_trades call
with invalid num **9999**, new section **1058**: marker plus return **0** or
abort **134**, ordinary and update mode. All reject with harness status **2**
at the missing-dependency guard BEFORE expected-failure acceptance/update;
the recording hash remains unchanged. Probe source/compiler wrapper removed;
the final default run rebuilds unchanged C. Logs:
`tmp/issue151-negative-{0,134}-{default,update}.log`.

Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded defaults pass with explicit isolated G3_REPO
and run_limited engine supervision. Logs: `tmp/issue151-*.log`.
The final fresh Olympia suite took **1.444s**. `git diff --check` reports
32 added literal trailing spaces, all unchanged real-C saved-list/report
bytes; the rest of the diff passes. Those bytes must not be stripped.
No later root/command, full-world/full-corpus/strict-lifecycle or stage-exit
acceptance. Independent parent acceptance passes the same tool tier, vet
and fresh tests (Olympia **1.380s**); logs `tmp/parent151-*.log`.
The parent's fresh post-#151 strict run has identical C repeats, C/Go 0,
matching streams/paths/types/modes, and eight report-byte differences,
exit **1** before lifecycle; log `tmp/parent151-strict.log`, captures in
`tmp/parent-stage2-turn`. G3 #358
remains nonblocking, with the original division/Go panic unchanged;
malformed C136/Go2 is NOT parity.

## Character detail sections (#150)

Existing report.go adds complete pinned report.c `char_rep_stack_info`
**410–433**, `show_pledged` **466–496**, `prisoner_health` and
`char_rep_prisoners` **499–537**, and health/combat/misc/magic **540–678**.
It reuses the existing `char_rep_location` **357–407**, two-pass
`collect_pledge_backlinks` **438–463** and world-owned lazy flag, not a second
collector. The character-report root remains its existing state projection
until #151; these direct renderers do not prematurely replace it or run both
effect paths. The original first eligible root call still owns collection
timing; the new direct show_pledged entry collects only on its own first call.

The complete unchanged C definitions and real closure were already linked:
stack_parent/loop_here snapshots, player/ownership, has_skill and
char_really_hidden/count_stack_any, noble-item and character/item accessors,
has_auraculum/max_eff_aura and real inventory filtering, existing names and
char_rep_location topology, wiout/out and the accepted FILE cache. Only
direct oracle entry points and numeric setters are new; no new dependency,
pin, command handler, rewritten C rendering or successful substitute.

Fifteen appended sections compare real unchanged C twice and Go: ten status
**0**, five status **134**. Health covers absent CH without allocation,
**-1** n/a versus **-2**, **0/1/99/100/127/-128**, signed sickness and recovery.
Combat distinguishes character **32767/-12/71** from noble-item **-7/91/13**
and recordless zero stats, front/negative/wrapped behind and break **0/57/100**;
the two spaces before the behind explanation are intentional. Stack/prisoner
fixtures use descending character IDs with interleaved noncharacters,
nested/direct children, prisoner filtering, positive/negative health and
repeat calls. Assertions preserve report.c **504** and **519**, including
first-invalid before opening a blocked output and invalid after a valid entry.

Pledge observers verify no eager collection from health/misc/stack, ALL
characters (including silent and prisoner), original iteration order, both
passes, missing CH allocation and retained allocated-empty bare CM, and no
recollection after changing a forward pledge. Hidden-skill fixtures separately
use learning and known rows (no duplicate-first-row masking), and actually
print all three concealment phrases: alone, other character/man inventory/
stack parent, and switched off. Misc shroud prints negative nonzero values;
magic shroud requires positive values. Vision, negative/zero/positive quick
and hide-mage gates, **32768 → -32768**, turn **32767+127 = 32894**, negative
innate clamp, per-entry rather than quantity aura bonuses, negative bonuses,
held/unheld auraculum, and location/ship versus character/invalid/zero project
targets are reached. The cached renamed inode receives **326 raw bytes**
under the actual MASTER frame without reopening. Constants were checked
against pinned oly.h: player/character/location/item/ship **1/2/3/4/9**,
plain **3**, galley **30**, silent **62**, hide_self **638**.

At **0.90.0-alpha**, command comparison passes **1050 sections / 70245 lines /
1339950 output bytes**. Accepted **1035 sections / 69581 lines / 1324865
output bytes** and **22446 lines / 420081 input bytes** remain byte-exact
prefixes against HEAD. New input: **22741 lines / 426753 bytes**.
Output SHA-256:
`97ef87c83a29d2ce7d7e20ed73a535f7734f9b507157dcc12ebae703f0639565`.
Input SHA-256:
`7cdc4a7fcd75b09184e2a4e25e2b79afd1532244e6903503fe6e88781696640f`.
#168 live native traversal still matches raw C/C/Go over unchanged shared input.
Logs: `tmp/issue150-command{,-update}.log`.

Four temporary linker controls wrap only the direct prisoner-health entry at
the new invalid section **1045**: marker plus return **0** or abort **134**,
ordinary and update mode. All reject with harness status **2** at the missing
dependency check, before expected-failure acceptance or recording replacement;
the recording hash remains unchanged during those probes. Probe sources and
compiler wrapper were removed, and final refresh/default runs rebuilt unchanged
C. Logs: `tmp/issue150-negative-{0,134}-{default,update}.log`.

Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded default tests pass with explicit isolated
`G3_REPO` throughout and engines under run_limited. Evidence:
`tmp/issue150-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.
The fresh Olympia suite took **1.572s**. `git diff --check` reports one new
literal trailing space in C's saved here-list row; it is intentional exact
oracle output, not whitespace to strip.

No full-world/full-corpus/strict-lifecycle/stage-exit acceptance is implied.
Latest strict evidence is still the parent's post-#145 diagnostic: deterministic
C repeats, C/Go 0 and equal raw streams/paths/types/modes, eight saved-report
byte differences, strict exit **1** before lifecycle. This task does not
refresh that gate. G3 #358 remains nonblocking; original division/Go panic
are unchanged, and malformed C136/Go2 is NOT passing parity.

Independent parent review verified all ten complete unchanged C definitions
exactly once, inherited byte prefixes and both hashes above. Parent sequential
compatibility checks, bounded vet and fresh bounded default tests pass
(Olympia **1.513s**). Logs:
`tmp/parent150-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.

## Known and partial skills (#149)

Existing use.go implements complete pinned use.c `exp_s` **977–994**,
`rep_skill_comp`/`flat_skill_comp` **1199–1238**, active `list_skill_sup`
**1262–1271**, `list_skills` **1276–1335**, `fractional_skill_qualifier`
**1345–1361** and `list_partial_skills` **1364–1406**. Existing lore.go adds
complete `np_req_s`, lore.c **120–134**. The real closure reuses use.c
`exp_level` **957–974** and `skill_school` **1174–1189**, original skill
accessor macros/records, skill_ents copy/reclaim/safe_qsort, box_name and
its existing code/display-name closure, sout/wiout and the accepted FILE
cache. These complete unchanged C definitions were already extracted/linked;
only direct oracle entry points/fixture setters/observers were added. No
new upstream source dependency or pin, modeled C renderer or successful stub.

Both lists sort COPIES only. Raw rows and actual save_box bytes before and
after repeated calls retain stored order, including duplicate skills and
mixed learning states. The original comparator puts nonknown entries FIRST
despite its comment; it avoids school reads for mixed known/nonknown pairs.
Known output skips those entries, groups by root school then ID, and indents
by the immediate prerequisite's presence (not recursive depth). Skill **708**
belongs to school **701** and prints before school **704**, distinguishing
school grouping from a flat numeric sort. No-exp nonzero values suppress
experience and short-circuit school lookup. A singleton cyclic no-exp skill
succeeds, while a two-entry comparator still visits that cycle and asserts.
The later character root remains its existing state projection until #151:
these direct renderers do not also run a projection or wire later roots.

Eighteen appended sections compare actual C twice and Go: eight status **0**
and ten status **134**. Cases cover absent/allocated-empty/all-known lists,
all **0/1/2** know states, unknown-only output, duplicate rows, nested
prerequisites, root/magic/recordless schools, positive and negative no-exp,
signed study days, exact blank lines and indentation/secondary-indent resets.
Experience inputs **-32768/-1/0/4/5/11/12/20/21/34/35/32767** exercise every
band/boundary and all five literal names. NP cases include absent records,
zero/negative costs, **1/2/2147483647**, zero/negative/maximum learn times,
and learning rows that must omit NP text. The renamed cached inode receives
**87 raw bytes**, including actual MASTER-routed skill lines, with no reopen.

Constants were checked against pinned oly.h/sout.h: character/skill **2/5**,
magic subtype **43**, SKILL_dont/learning/know **0/1/2**, experience levels
**1–5**, MASTER **-2**; **701–708** are synthetic valid entities and **9999**
is deliberately absent. Malformed cases preserve source assertions at
use.c **1283/1371** (invalid owner before output/open or missing-record
checks), **1182** (school cycle), **1349/1358** (qualifier state before skill
reads), and **990** (invalid experience level 0/6). The original assertion
and buffering behavior is compared, not normalized into successful reports.

At **0.89.0-alpha**, final unchanged-C refresh and default command check pass:
**1035 sections / 69581 lines / 1324865 output bytes**. The inherited
**1017 sections / 68918 lines / 1310660 output bytes** and **22223 lines /
414715 input bytes** remain exact unchanged prefixes against HEAD. New input:
**22446 lines / 420081 bytes**. Output SHA-256:
`6a6f77b0297bda957882e962dbef9b69a083846768bb50af5787bb785a64cac8`.
Input SHA-256:
`2f4b2838af26d55dfaaf3b8d22b394a58de189cd9f38cbb5c98aad50ec048c11`.
#168 live native traversal still matches C/C/Go with its shared input unchanged.
Logs: `tmp/issue149-command{,-update}.log`.

Four temporary linker controls wrapped only list_partial_skills at the new
malformed-state section **1029**: marker plus return **0** or abort **134**,
ordinary and update mode. All reject with harness status **2** at the missing
dependency check, before expected-failure acceptance/recording replacement;
the recording hash remains unchanged. Probe source/wrapper were removed and
the final default command run rebuilt unchanged C. Logs:
`tmp/issue149-negative-{0,134}-{default,update}.log`.

Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded default tests pass, with explicit isolated
`G3_REPO` throughout and engines under run_limited. Evidence:
`tmp/issue149-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.
No new production defect/prerequisite, later root/command, full-world,
full-corpus, strict-lifecycle or stage-exit acceptance is claimed. Latest
strict evidence remains the parent's post-#145 run: identical C repeats,
C/Go 0, matching streams and paths/types/modes, eight saved-report byte
differences, exit **1** before lifecycle. This task does not refresh it.
G3 #358 remains nonblocking: original division/Go panic untouched;
malformed C136/Go2 is NOT passing parity.

Independent parent review verified all eight complete unchanged C definitions
exactly once, inherited byte prefixes and both hashes above. Parent sequential
compatibility checks, bounded vet and fresh bounded default tests pass
(Olympia **1.558s**). Logs:
`tmp/parent149-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.

## Admission permissions, attitudes and player entity reports (#148)

Existing perm.go implements complete unchanged perm.c `admit_comp`,
`print_admit_sup`, `print_admit` **141–223** and `print_att_sup`/`print_att`
**426–488**, including both complete closing definitions.
Existing report.go replaces playerEntInfoState with complete report.c
`player_ent_info` **783–798** and main calls that renderer exactly once.
The closure uses accepted #146 `show_unclaimed` **746–780**, inventory
comparison/formatting, existing box-code/validity and player/disposition
accessors, int_comp/safe_qsort, and the original sout.c output/cache owner.
All C definitions were already extracted/linked; no new dependency or pin.
The character root remains deferred: reportAttitudeState is still its
state-only projection until #151, not an extra rendering/sorting pass.

Eleven appended sections compare original C/C/Go bytes, raw lists, record
allocation and save_box output. Coverage includes missing/allocated-empty
player records, missing/empty attitude lists, invalid-only lists and targets,
valid empty admissions (heading but no command row), sense 0/1/negative,
equal admission targets, duplicates, unsorted admission members, retained
invalid references, all three in-place attitude sorts and inventory sorting.
Explicit member counts **11/12/13** run both without and with `all` (which
itself consumes one slot). Attitude valid counts are **11/12/13** across
hostile/neutral/defend; invalid rows do not consume slots. Real continuation
lines retain ten admission padding spaces and seven attitude header spaces.
Neutral-first and defend-first controls cover the lazy heading and literal
trailing space in `defend `. Primary/secondary indentation and day framing
are preserved, including a repeated root using the renamed cached inode
(**358 raw bytes** under `log/original`). Original nonplayer assertion is
perm.c **197**, status **134**; the other ten new sections exit **0**.

Constants checked against pinned oly.h: player/character/item **1/2/4**;
system/regular/savage/NPC/silent subtypes **48/49/51/52/62**. Root cases
actually emit all three sections for each nonsilent subtype, including
subtype 0, and compare raw files **202–205/207**; silent **206** has no file,
no missing-player allocation and retains unsorted original lists/inventory.
Items **501/502** are synthetic and **9999/8888** deliberately absent.
The subtype raw-file fixture uses its actual cwd-relative output root;
missing-file observations are not used as evidence of rendered bytes.

At **0.88.0-alpha**, final refresh/default checks match complete real C twice
and Go: **1017 sections / 68918 lines / 1310660 output bytes**. The inherited
**1006 sections / 68035 lines / 1290206 output bytes** and **21934 lines /
408111 input bytes** remain exact unchanged prefixes against HEAD. New input:
**22223 lines / 414715 bytes**. Output SHA-256:
`0f4c2412955269f5667d56926ab1003a9f98b1cd5920704bc7b163370cdd724d`.
Input SHA-256:
`79f6954b243f34ec6086307632a659c5eb6eea4ab554b0ebb1ad752e7c4bb5c8`.
#168 live native traversal C/C/Go still matches, shared input unchanged.
Logs: `tmp/issue148-command{,-update}.log`.

Temporary linker controls wrapped only print_admit at expected-abort section
**1016**, emitted the missing-dependency marker and returned **0** or aborted
**134**. Normal/update runs all rejected with harness status **2** at the
marker check, before expected-failure acceptance/recording replacement;
each probe preserved the recording hash. Probe sources/wrapper were removed
before the final unchanged-C refresh. Logs: `tmp/issue148-negative-*.log`.

Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded default tests pass, all compatibility tools
with explicit isolated `G3_REPO` and engines under run_limited. See
`tmp/issue148-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.
No new production defect or prerequisite was found. These are permission
reports, not valid new-player admissions or new commands. No later report
root, full-world/full-corpus/strict-lifecycle or stage-exit acceptance.
Latest strict evidence remains the parent's post-#145 run: identical C
repeats, C/Go 0, matching streams and paths/types/modes, eight saved-report
byte differences, exit **1** before lifecycle. This task does not refresh it.
G3 #358 remains nonblocking: original division/Go panic untouched and
malformed C136/Go2 is NOT passing parity. Independent parent review verified
the six complete unchanged C definitions exactly once, inherited prefixes
and both hashes above. Parent sequential compatibility checks, bounded vet
and fresh bounded default tests pass (Olympia **1.595s**). Logs:
`tmp/parent148-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.

## Capacity and item-granted skills (#147)

Existing report.go implements complete unchanged pinned report.c
`show_carry_capacity` **73–105**, `show_item_skills_sup` **108–193**
(active code, retaining its two distinct loops) and `show_item_skills`
**196–210**. The existing oracle already extracts all three definitions;
direct `report capacity` and `report item_skills` expose them. Fixture-only
`may_study`/`may_use` fields append real stored lists on both sides.
No new dependency, pin, successful C substitute or output/cache owner.

The complete real closure reuses u.c `determine_unit_weights`/
`add_item_weight`, `max`, `comma_num`; use.c `skill_school`, `has_skill`,
`rp_skill_ent`, `exp_level`; original inventory/access macros, code.c
`box_name`/`just_name` and naming dependencies; and complete sout.c
`out`/`wiout`/routing/cached FILE ownership. Go reuses weights.go, use.go,
code.go and sout.go. `wout` is the pinned `out` macro, not invented wrapping.
Later character/player roots remain projections: neither these sections nor
a projection-plus-render double pass are wired into those roots yet.

Fourteen appended sections compare real C/C/Go, including raw/saved state:

- Zero weight/capacity, loaded zero capacity, separate walking/riding/flying
  self-weight exclusions, unequal positive capacities, overfull loads,
  negative capacity's self-carry/no-extra-capacity semantics, custom noble
  item, signed-int multiplication/ratio overflow and literal comma output.
- Empty and missing magic records, study-only/use-only and separate mixed
  lists, original inventory/list order and duplicates, nonpositive/invalid
  inventory rows filtered before skill validation, repeated rendering,
  schools/subskills and multi-level prerequisites. Recipient (not owner)
  knowledge determines `???`; a learning school is insufficient while a
  known school with negative experience still has novice level.
- Primary indentation restoration, secondary indent reset only on `wiout`,
  unchanged secondary indent with empty lists, day fields, MASTER routing
  and actual renamed cached inode bytes (328 bytes under `log/old`).
- Original assertions at report.c **126/162**, use.c **1182** (prerequisite
  cycle) and u.c **1480** (noncharacter weight owner). Five new sections
  exit **134**; the other nine exit **0**. Buffered output on assertion is
  compared as emitted by actual C, not artificially flushed for inspection.

Constants verified against pinned oly.h: char/item/skill **2/4/5**, default
noble item/peasant **10**, learning/known **1/2**. Items 501–504 and skills
701–704 are synthetic, 9999 deliberately absent. Capacity fixtures allocate
the actual default noble item before weight lookup (the initially incomplete
synthetic fixture lacked it and C exited 139; that fixture was corrected,
not accepted or masked). Observed asymmetric output is **100/300 land (33%),
125/160 ride (78%), 100/180 fly (55%)**, becoming **560/300 (186%),
585/160 (365%), 560/180 (311%)** under load. Negative-capacity control prints
**0/0 land, 0/30 ride (0%), 0/22 fly (0%)**. Stored inventory/skill/grant
lists remain unchanged by these section renderers, with no new records on
the empty-list path. No production defect or new prerequisite was found.

At **0.87.0-alpha**, refresh and fresh default command checks match complete
pinned C twice and Go: **1006 sections / 68035 lines / 1290206 output bytes**.
The inherited **992 sections / 67520 lines / 1278600 output bytes** and
**21714 lines / 403193 input bytes** remain exact unchanged prefixes, checked
against HEAD. New input: **21934 lines / 408111 bytes**.
Output SHA-256:
`62d0443f0b7d7122c2b2a28c8b2556cb4c9995a20782f7428a63b5a1afe62e86`.
Input SHA-256:
`6ae3de40cb526f9b9c1d473959fa4ac5e64fcdcb8c564129fc3a0df5a54da5b3`.
#168 live native traversal still matches **28 lines / 426 bytes**, shared
input unchanged. Logs: `tmp/issue147-command{,-update}.log`.

Temporary linker controls wrapped only `show_item_skills` at expected-abort
section **1002**, emitting the missing-dependency marker then exiting **0**
or aborting **134**. Ordinary and `--update` runs all rejected with harness
status **2** at the marker check before status/recording acceptance, and
recording SHA-256 was unchanged in all four runs. Probe source/compiler
wrapper removed; final refresh/default checks use only unchanged actual C.
Logs: `tmp/issue147-negative-{return,abort}-{normal,update}.log`.

Sequential checks pass, each with explicit `G3_REPO="$PWD/tmp/upstream.git"`
and bounded engine execution: command **1006/68035**, move **135/7158**,
helper **37/1487**, RNG **321**, code **43/969** (its own unchanged older
pin), loadsave **61/3411**, turn-driver **7/152**, phase **5/484**, calendar
**2/119**, schedule **4/1801**, default turn (raw streams/status and both cwd
entries match), bounded `go vet ./...`, fresh bounded
`go test -count=1 -exec "bash $PWD/tools/lib/run-limited.sh" ./...`
(Olympia **1.742s**). Logs: `tmp/issue147-{move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.
Source/docs diff whitespace check passes; original C recording trailing
spaces remain unmodified. Matching command captures removed after checking.

No whole-world/full-corpus/strict-lifecycle or stage-exit acceptance claimed.
Latest strict integration evidence is still the parent's post-#145 run:
C repeats identical, C/Go statuses zero, streams and entry paths/types/modes
match, eight saved-report byte differences, exit **1** before lifecycle.
This task does not refresh that gate. G3 #358 remains nonblocking; original
division and Go panic are unchanged, malformed C136/Go2 is NOT passing
parity. Independent parent review verified all three complete unchanged C
definitions exactly once, the inherited prefixes and hashes above. Parent
sequential compatibility checks, bounded vet and fresh bounded default tests
pass (Olympia **1.729s**). Logs:
`tmp/parent147-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.

## Inventories and unclaimed items (#146)

Existing report.go implements complete unchanged pinned report.c
`inv_item_comp` **213–220**, `extra_item_info` **246–290**,
`show_char_inventory` **294–344** and `show_unclaimed` **746–780**.
These C definitions were already extracted in full by this oracle; two direct
report operations now expose the sections without implementing later roots.
Their complete real closure is the existing `loop_inv`/`safe_qsort`, original
item/magic access macros, u.c `comma_num`/`nice_num`, code.c
`plural_item_box`/`plural_item_name`/`box_name` and naming dependencies, and
complete sout.c output/cache/routing. No C substitute, source dependency,
pin, inventory storage, magic record or output owner was added.

Both renderers sort the original inventory before traversal/output, retaining
invalid/nonpositive rows which traversal and save_box may filter. Products,
sum and comparator retain C int widths. Inventory prints its original total
and empty message; unclaimed prints neither and always passes viewer zero
to the aura formatter. Capacity precedence is fly, ride, land, then nothing;
fighter tuples include the zero-stat ghost exception; signed bonus strings
concatenate literally, including `+-4`, without invented spaces. Aura checks
the recipient's magician flag, not inventory ownership or effective aura.
Later character/player roots still execute only their existing projections;
they do not render these sections or run both paths. Their original sort
positions are unchanged and share the comparator. #147's capacity table and
item skills, #148's player root and #151's character root remain separate.

Ten appended cases compare raw/saved original lists, bytes, allocations and
warnings: empty/fully filtered inventories, duplicate IDs, invalid positive
and nonpositive rows, repeats, independent unclaimed sorting, unequal and
negative weights, all capacity precedence branches, independent attack/
defense/missile fighter tests, missing item records, ghost zero stats,
positive/negative equipment bonuses, magician recipient versus nonmagician
owner, unclaimed aura suppression, null/empty plural warnings, signed-int
product/total overflow, prisoner-filtered output and cached renamed inodes.
Actual C output confirms `fly 77`, `ride 69`, `cap 62`, `(0,0,0)`,
`+-11 aura` only for the magician recipient and unclaimed `cap 39`/weight
21 with no aura/total. All ten new sections exit zero. Constants checked
against pinned oly.h: char/item **2/4**, ghost warrior **18**. IDs 501–504
are synthetic ordinary items, not assumed built-in equipment IDs; 9999 is
deliberately absent. Fixtures set the real typed fields; no threshold or
branch is inferred merely from the case name.

At **0.86.0-alpha**, refresh and fresh default command runs match two complete
real C runs and Go: **992 sections / 67520 lines / 1278600 output bytes**.
Inherited **982 sections / 66938 lines / 1262604 output bytes** and
**21537 lines / 399178 input bytes** remain exact byte prefixes (verified
against HEAD). New input: **21714 lines / 403193 bytes**.
Output SHA-256:
`96fe2f2e43b4b63c26cd466fffd40a7f64e1dbda485170dd119a349799e8d743`.
Input SHA-256:
`18f9256047a84b48c1975fd0aed190fad42ac7b2d5f77b00761b75475afa01f6`.
#168's separate live native C/C/Go check matches **28 lines / 426 bytes**,
with shared input unchanged. The existing missing-dependency marker check
still precedes expected-failure acceptance and recording replacement.
Logs: `tmp/issue146-command-update.log`, `tmp/issue146-command.log`.

The required shared-output tier passed sequentially with explicit
`G3_REPO="$PWD/tmp/upstream.git"`: command **992/67520**, move **135/7158**,
helper **37/1487**, RNG **321 results**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**,
and default turn (raw streams/status and both cwd entries match).
Bounded `go vet ./...` and fresh
`go test -count=1 -exec "bash $PWD/tools/lib/run-limited.sh" ./...` pass;
Olympia **1.513s**. Logs: `tmp/issue146-{move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.
These defaults do not run saved-world/full-corpus comparisons. Source/docs
diff whitespace checks pass; original C recording trailing spaces are kept.

Independent parent review verified all four unchanged complete C definitions
exactly once, inherited byte prefixes and both hashes above. Parent sequential
compatibility checks, bounded vet and fresh bounded default tests pass
(Olympia **1.402s**). Logs:
`tmp/parent146-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.

No whole-world, strict lifecycle, full corpus or stage-exit acceptance is
claimed here. The parent's post-#145 strict full run remains the latest:
deterministic C/C, C/Go zero statuses, matching raw streams and entry
paths/types/modes, but eight saved-report byte differences; exit 1 before
lifecycle controls. G3 #358 remains nonblocking; original division and Go
panic are untouched, and malformed C136/Go2 is not passing parity.

## Complete location report root (#145)

Existing display.go implements complete pinned display.c `show_loc`
**1056–1101**. Existing report.go replaces its allocation projection with
complete report.c `turn_end_loc_reports` **1589–1642**, and the existing main
caller now calls that renderer exactly once. The command oracle extracts the
complete unchanged `show_loc` definition and original `show_loc_no_header`
declaration; the complete original turn-end root was already extracted.
The last rejected `show_loc` dependency is removed only now that its entire
rendering closure is real. No new source dependency, pin or output owner.

The original subordinate sequence is header/blank, barrier/blank, pillage/blank,
ruler, sailable routes, exits, stats, nearby cities, taught skills, city-only
market, inner locations, posts, weather, characters, ships. It reuses all
accepted #138–#144 sections, names, knowledge, traversal, ship weights,
ownership, original macros and the existing sout.go cached FILE owner.
Turn-end uses real `lines` for separate-format headings and MASTER/OUT_LOC
versus per-location routing, not a second report layer. It reads format before
`p_player`, allocates every nonsilent record, sorts the original touched list,
skips but retains invalid/deleted IDs, and resets routing/header globals to
zero/false rather than restoring incoming values. `show_loc` sets display
strings true and then false; an empty turn-end call never changes that flag.
Assertion failures are not given deferred cleanup that C does not perform.

Nine appended sections compare actual complete-C small-world report bytes,
raw globals/lists/record-presence, output bits, real save_box bytes and
market sort state. The city composition proves barrier → pillage → routes →
damage → nearby cities → skills → market → inner owner stack → post → rain →
characters → docked ships order. A province composition proves ruler/pledge
before routes and posts before privileged fog descriptions/character suppression.
A ship composition actually renders ocean routes, exits, damaged capacity and
captain display text; buildings suppress the outer weather as in C.
Direct calls cover headers/suppression, signed pillage, singular/plural/recent
text, nonzero indentation/day framing, and display-flag reset (not restoration).
The turn-end cases cover regular/system/savage/NPC and unknown-subtype players,
both formats including negative nonzero format, silent skips, missing records
for every named nonsilent subtype, unsorted/stale/deleted touched IDs, repeated
calls and format switching. Actual C output confirms all intended branches.
Shared shuffle/NPC/entity aftermath and renamed cached-inode/reopen controls
reuse existing oracle operations. Invalid direct IDs and malformed posts
reached through the full root retain source assertions **1061/912**.

Fixture constants were checked against the pin's oly.h: player/char/loc/item/
skill/road/ship/post/storm **1/2/3/4/5/7/9/10/11**, ocean/forest/plain/mountain/
city **1/2/3/4/14**, galley/roundship/castle/tower **30/31/32/46**, system/
regular/region/savage/NPC **48/49/50/51/52**, fog/wind/rain **57/58/59**,
silent/garrison **62/64**, gold/peasant **1/10**, BUY/SELL **1/2**; sout.h
OUT_LOC is **4**. No positive-quantity zero-price character BUY is fabricated.

At **0.85.0-alpha**, refresh and a fresh default command run pass real C/C/Go:
**982 sections / 66938 lines / 1262604 output bytes**. Both complete inherited
prefixes are byte-exact: **973 sections / 66015 lines / 1241870 output bytes**
and **21238 lines / 393356 input bytes**. New input: **21537 lines / 399178
bytes**. Recording SHA-256:
`cef1f88175fc1a5cd42973a4428ab27cd9f5bd6a136f13e49e5bc0f74072b87b`.
Input SHA-256:
`0e36ad84c75cfdda2a50b659510408c1c2f975822225f17368d6f5f8f3b9edd4`.
#168's separate live raw native-order C/C/Go check still matches **28 lines /
426 bytes**, with shared input unchanged. Expectations come only from two
deterministic real C runs, never Go or callback traces.

A temporary compiler/driver overlay injected a missing-dependency marker and
status **134** in expected-invalid root section **980**. `--update` rejected it
before accepting the expected failure or replacing the recording: status **2**,
recording hash unchanged. Temporary overlay files/captures were removed.
Logs: `tmp/issue145-{source-check,branches,prefix,marker}.log`.

Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded default tests pass. Evidence:
`tmp/issue145-{command-update,command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.
All compatibility invocations explicitly use isolated `tmp/upstream.git`;
engine/test processes are bounded by run_limited. code-compat deliberately
retains its unchanged older pin. Intentional C save_box trailing spaces remain.
Independent parent review verified both complete C definitions unchanged
exactly once, the inherited byte prefixes and hashes above. Parent sequential
compatibility checks, bounded vet and fresh bounded tests pass (Olympia
**1.373s**); logs `tmp/parent145-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.

Parent also reran strict `turn-compat --full` with the isolated exact-pin C
cache. Two fresh C turns matched all artifacts. C/Go both exited **0**, raw
stdout/stderr/status and all entry paths/types/modes matched, but eight saved
report files still differ: `lib/save/1/{100,200,202,203,204,205,210,players}`.
The strict gate exits **1** before lifecycle controls; no full-world/lifecycle
or stage-exit acceptance is claimed. Evidence: `tmp/parent145-strict.log` and
`tmp/parent-stage2-turn/{c-determinism.txt,diff.txt,c,go}`. Later report
roots/commands remain deferred. G3 #358 is owner-declared nonblocking: existing division and Go panic
are unchanged, and malformed C **136** versus Go **2** is not passing parity.

## Ancillary location sections and posts (#144)

Existing display.go now implements complete stage-pin display.c
`show_loc_stats` **92–122**, `show_nearby_cities` **842–869**,
`show_loc_skills` **872–892**, `show_loc_posts` **895–963**,
`show_weather` **966–1012** and `show_loc_ruler` **1015–1046**.
The oracle extracts these complete unchanged definitions, replacing the
rejected `show_loc_posts` dependency, and adds complete storm.c
`can_see_weather_here` **1661–1676** in the existing storm.go owner.
The existing location-report root remains deferred to #145; no parallel
projection/output path, POST handler or weather-magic command is introduced.

The real reused closure is display.c `liner_desc`/`liner_desc_storm`,
`safe_haven_s`; storm.c `weather_here`; garr.c `garrison_here`,
`top_ruler`, `province_admin`, `rank_s`; loc.c `first_character`,
`province`, `building_owner`, `all_char_here`/`add_char_here`; swear.c
`player`/`lord`; u.c `ship_cap`, `ship_weight`, `determine_unit_weights`,
`add_item_weight`, `comma_num`; original oly.h `noble_item`, `safe_haven`,
`loc_damage`, `garrison_castle`, `item_creator` and storm-field macros;
original loop.h direct-here and taught-skill copied loops;
and complete code.c/sout.c/io.c for names, comma_append, raw output,
cached FILE/inode ownership and actual save_box. All source dependencies were
already tracked: no pin or dependency-table change.

Nineteen appended sections exercise direct post order, nonpost and nested-post
skips, single/multiple/empty text lines, first/last quotes, authors and absent
item-magic records, both full-location names, retained text on repeat calls,
MASTER/style/day output and renamed cached inodes. Missing/allocated-empty
post text asserts at display.c **912**, including a malformed second post
after valid output but before its separator. Invalid post/skills containers
retain macro assertions **904/878**. Nearby cities retain duplicates/order
and safe-haven suffixes; taught skills retain stored order and duplicates.
Weather controls cover all eight rain/wind/fog combinations, signed strengths
and cancellation, ordinary versus privileged output, stale province views,
output allocating a missing player record before visibility lookup, and the
distinction between province weather aggregation and the queried location's
direct storm descriptions.
Buildings suppress weather; nested storms do not contribute to province sums.
Rulers retain first-character garrison selection, castle/no-owner/no-castle
branches, top pledge/rank, and ocean/sublocation skips. Stats retain damage
byte wrapping, zero damaged capacity, nested noble/item weights, repeat reads,
absent records and signed C-int percentage arithmetic. Raw record flags,
lists, output bits, indentation and save_box bytes supplement report bytes.

Fixture values are checked against pinned oly.h, not descriptive labels:
T_player **1**, T_char **2**, T_loc **3**, T_item **4**, T_skill **5**,
T_ship **9**, T_post **10**, T_storm **11** (not T_deadchar **8**),
ocean **1**, plain **3**, mountain **4**, city **14**, galley **30**,
roundship **31**, castle **32**, tower **46**, fog/wind/rain **57/58/59**,
garrison **64**, item_peasant **10**. The synthetic fallback-weight control
includes that real item record; no invalid fixture crash is blessed as parity.
The overflow fixture builds **30,030,000** weight using a short-safe item
weight **30,000** and inventory quantity **1000**, rather than silently
truncating a large item weight. Actual C renders **-1291%** at capacity
**1,000,000**, distinguishing signed 32-bit multiplication from a widened ratio.

At **0.84.0-alpha**, recording refresh and a fresh default command run both
pass real C/C/Go: **973 sections / 66015 lines / 1241870 bytes**. All inherited
**954 sections / 65118 lines / 1224413 output bytes** and **20885 lines /
385519 input bytes** remain exact prefixes. New input is **21238 lines /
393356 bytes**. Recording SHA-256:
`4a83dabbf7bf4a812d7d75423a2e012c787eb95c90b972f8998560130bbb4f1f`.
#168's separate live native-order check still matches C/C/Go (**28 lines /
426 bytes**) without mutating its shared input. A temporary compiler/driver
overlay injected a missing-dependency marker and status **134** in expected
malformed-post section **958**. `--update` rejected it before status acceptance
or recording replacement (status **2**, then-current recording hash unchanged).
Temporary overlay scripts and captures were removed; the diagnostic is
`tmp/issue144-marker.log`.

Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded `go vet ./...` and fresh bounded `go test -count=1` defaults pass.
Evidence logs: `tmp/issue144-{command-update,command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.
Every compatibility invocation explicitly uses isolated `tmp/upstream.git`;
every engine and Go test executable is bounded by run_limited. code-compat
deliberately retains its unchanged older pin. Default turn is the small
diagnostic comparator, not a world turn; loadsave's saved worlds were not run.
Intentional C save_box trailing spaces in the recording are preserved.

Independent parent review verified all seven complete unchanged C definitions
exactly once, both inherited byte prefixes and the recording/input hashes.
Parent sequential compatibility checks, bounded vet and fresh bounded default
tests pass (Olympia **1.469s**). Logs:
`tmp/parent144-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.

No full-world or stage-exit acceptance is implied. G3 #358 remains owner-declared
nonblocking: original division and Go panic are unchanged; malformed C **136**
versus Go **2** is not passing parity.

## Inner locations and ship sightings (#143)

Existing display.go adds complete stage-pin display.c `show_inner_locs`
**685–717**, `show_sublocs_here` **720–797** and `show_ships_here`
**800–839**. The oracle extracts each complete unchanged definition exactly
once, with original source locations and loop-macro assertions. Three trivial
wrappers expose the static functions; expected output is produced only by two
deterministic executions of real C, never by a replacement renderer.

The reused real closure is display.c `liner_desc` and its location/ship/character
formatters, `display_owner`, `show_owner_stack`, `show_chars_below_highlight`,
`highlight_units` and `show_loc_barrier`; u.c `loc_hidden`/`test_known`, swear.c
`player`, original oly.h `is_loc_or_ship`/`see_all`/`loc_link_open` macros,
loop.h's copied direct here-list iteration, dir.c `is_port_city` and topology,
loc.c `province`/`all_char_here`, actual ship capacity/weight dependencies,
and complete original code.c/sout.c/io.c (names, wiout, cached FILE ownership,
raw output and save_box). All were already linked/extracted. There is no new
source dependency, pin, gameplay command, full show_loc root, posts or parallel
output layer. Previously accepted display.go bytes remain unchanged.

Fourteen appended sections distinguish unknown hidden parents from known
descendants (the parent still prunes), known hidden containers/ships, see_all
**0/1/2**, character versus player recipients, and no renderer-added knowledge.
Direct cities omit owners/recursion but retain barriers; cities reached by
inner recursion do render owners/descendants. Buildings recurse through both
locations and ships. Open crosslinks follow direct locations in original list
order, retain duplicates and direct-target repeats, and omit owners, barriers
and recursion; closed links remain absent even when known or see_all is set.
Signed nonzero link_open is exercised. Empty/wrong-kind/all-hidden inputs
create no header/file and preserve both indent fields and absent records.
Nonempty cases compare original save_box bytes, raw known-list order before
save, output bits and exact primary/secondary indentation, including ship
wiout's fixed secondary **3**. Ocean and port headings, multiple ship order,
owner/descendant/barrier order, deliberate absence of ship-inner recursion,
MASTER/style/day routing and renamed cached-inode output are byte-observed.

Fixture numeric constants were checked against pinned oly.h: T_loc **3**,
T_ship **9**, ocean **1** (not desert **5**), plain **3**, cave **13**, city
**14**, ruins **17**, temple **29**, galley **30**, roundship **31**, castle
**32**, tower **46**. Actual C output contains `Ships sighted:` on the ocean
and `Ships docked at port:` on the city control. The first recording attempt
correctly rejected status **139** from a synthetic orphan linked city: its
is_port_city dependency needs a province. The fixture now has a real province;
no crash was converted to passing parity or bypassed in production.

Four expected-abort cases retain original assertions: invalid containers at
display.c **691/727/806**, and invalid knowledge recipient at u.c **2334**
even with see_all=2 (knowledge is evaluated first). A temporary in-memory
negative-control wrapper emitted a missing-dependency marker and status 134
only in that last expected-abort case. `--update` rejected it **before** status
acceptance/recording replacement (section 953); the recording hash remained
unchanged. Temporary negative-control captures were removed; its diagnostic
log remains `tmp/issue143-marker.log`.

At **0.83.0-alpha**, real C/C/Go matches **954 sections / 65118 lines /
1224413 bytes**. The entire inherited **940 sections / 64201 lines /
1206911 output bytes** and **20552 lines / 378933 input bytes** remain exact
unchanged prefixes. New input: **20885 lines / 385519 bytes**. Recording
SHA-256: `0f9c882fd130b4613fa6a32a9cc58d920e8f35e837049dfa5a1273ffc3728d43`.
Both recording refresh and fresh default command runs passed, including
#168's separate live raw native C/C/Go traversal over unchanged shared input.
Matching command captures were removed after the checks.

Sequential bounded command/move/helper/RNG/code/loadsave/turn-driver/default
turn checks passed with explicit `G3_REPO=$PWD/tmp/upstream.git`: move
**135/7158**, helper **37/1487**, RNG **321**, code **43/969**, loadsave
**61/3411**, turn-driver **7/152**, phases **5/484**, calendar **2/119**,
schedule **4/1801**. Default turn checks raw streams/status and two cwd entries,
not a whole world. Bounded `go vet ./...` and fresh bounded
`go test -count=1 -exec ... ./...` pass (Olympia **1.430s**).
Logs: `tmp/issue143-command{,-update}.log` and
`tmp/issue143-{move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.
Source/document diff whitespace checks pass; full diff-check flags only 25
literal trailing-space save_box list lines in real C expectations, preserved
without normalization. No newly verified production defect or prerequisite.

G3 #358 remains owner-declared nonblocking. Existing market division/Go panic
are untouched, and malformed C **136** / Go **2** remains **not passing** raw
failure parity. No full-world/full-corpus/strict lifecycle or stage-exit
acceptance is implied. Independent parent review verified all three complete
unchanged C definitions exactly once, both inherited byte prefixes and the
recording hash above. Parent sequential compatibility checks, bounded vet and
fresh bounded tests passed (Olympia **1.724s**); logs are
`tmp/parent143-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.

## Visible characters and owner stacks (#142)

Existing display.go adds complete stage-pin display.c `highlight_units`
**547–558**, `show_chars_below_highlight` **579–593**, `show_owner_stack`
**596–626** and `show_chars_here` **629–682**. Existing `show_chars_below`
**561–576** and accepted character descriptions are unchanged. The oracle
extracts all five complete unchanged C definitions, retaining their original
source locations and macro assertions. New static-function wrappers only
expose the real functions; they do not render or implement expected logic.

The reused real closure includes display.c liner_desc/liner_desc_char and
display_with; swear.c player; loc.c all_char_here/add_char_here and province;
stack.c stack_parent; u.c count_stack_any/count_any and loc_depth; storm.c
weather_here; original loop.h snapshots and oly.h char_really_hidden macros;
and whole original code.c/sout.c/io.c, including init_spaces, wiout, routing,
the accepted FILE cache and save_box. All these dependencies were already
extracted or linked. No new source dependency, pin, location-report root,
gameplay command or parallel output/description owner is added.

Fifteen appended sections distinguish the first hidden owner (shown) from
later hidden owners (skipped), while noncharacters do not consume the first
owner exception. Ordinary sightings hide even own truly hidden units. Nested
characters retain depth-first snapshot order and the same descendant depth;
prisoners are described and noncharacter branches pruned. Highlight bytes
depend on a player recipient, not a character recipient or MASTER routing;
literal spacing is checked through depth 150, the original spaces buffer's
defined range. Ocean `, flying` applies only to direct sightings, not their
descendants or owner-stack rendering. Fog short-circuits province rendering
before indent changes (even with negative primary indent), including empty
provinces; signed storm cancellation, city/building and owner-stack controls
distinguish that branch. Empty/all-hidden cases emit no header or file and
preserve both indentation fields; successful wiout calls reset secondary
indent while restoring primary indent. Raw save_box/record observers,
MASTER/style/day/output bits and renamed cached-inode bytes remain compared.

Five expected-abort controls compare assertion order and source locations:
highlight depth **display.c:551** precedes indent **552**, invalid owner
container fails its original loop macro at **606**, invalid sighting container
fails earlier in **u.c:388** loc_depth, and negative owner depth fails at
**display.c:551**. Missing-C-dependency markers are rejected before expected
statuses and before recording replacement, including these failure cases.

At **0.82.0-alpha**, the real pinned C repeats deterministically, then Go
matches **940 sections / 64201 lines / 1206911 bytes**. The inherited **925
sections / 63547 lines / 1195280 recording bytes** and **20305 input lines /
373845 input bytes** remain exact unchanged prefixes. New input totals are
**20552 lines / 378933 bytes**; recording SHA-256:
`b4f5f85b6135273b3c09dc125208a7d9e022e599d863bf905428f9b16fe9b660`.
#168's separate live native C/C/Go comparison remains enabled and passes over
the same unchanged input directory. Recording refresh evidence:
`tmp/issue142-command-update.log`. Matching C captures are removed after checks.

The default command comparison reran both real C executions and Go without
updating expectations, and passed again. Sequential bounded command/move/
helper/RNG/code/loadsave/turn-driver/default-turn tools passed with explicit
`G3_REPO=$PWD/tmp/upstream.git`: move **135/7158**, helper **37/1487**, RNG
**321**, code **43/969**, loadsave **61/3411**, turn-driver **7/152**, phase
**5/484**, calendar **2/119**, schedule **4/1801**. Default turn matches raw
streams/status and two cwd entries; it is not a world-turn test. Bounded
`go vet ./...` and fresh `go test -count=1 -exec ... ./...` passed (Olympia
**1.519s**). Logs: `tmp/issue142-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`.
Independent byte inspection confirms each of the five extracted C definitions
appears complete and unchanged exactly once, and the existing Go
showCharsBelow body is unchanged. Source/documentation diff whitespace checks
pass; full diff-check reports only twelve literal trailing-space here-list
lines emitted by the real save_box writer in the new recording. Those bytes
are deliberately preserved, not normalized.

Parent review corrected the ocean fixture's mistaken desert subkind before
acceptance, then regenerated only new real-C expectations. The corrected
recording has exactly two `, flying` suffixes on direct sightings and none
on descendants or owner-stack lines. Parent independently reviewed the full
pinned C/Go diff, verified all five extracted definitions and inherited
prefixes, and reran all eight tools, bounded vet and fresh bounded defaults.
Evidence: `tmp/parent142-*.log`; these supersede the pre-review byte/hash
measurements in the original `issue142` logs.

G3 #358 remains owner-declared nonblocking and unrelated to this task. The
original market division and Go panic are untouched; malformed C **136** /
Go **2** is still **not** passing raw/status failure parity. No full-world,
full-corpus, strict lifecycle or stage-exit acceptance is implied by #142.
No newly verified local/upstream defect or blocking prerequisite was found.

## City market rendering (#141)

The existing buy.go now implements complete unchanged buy.c
`list_market_items` **1027–1066** and `market_report` **1069–1120**. The
oracle extracts both original definitions, reusing its actual buyer_list,
seller_list, reduce_qty, names/plurals, output/cache and save_box closure.
No new upstream dependency or pin, settlement, BUY/SELL handler, movement
handler or full location-report root is added. Production notices use the
original loop_trade filter and order, not the offer visibility filter.

Nine appended sections compare empty/repeated reports; nested native buyer
order before all sellers; stable equal-price seller order and city offers
last; real original seller sort-field mutation (including unavailable rows),
unchanged buyer/city sort fields, inventories, quantities, queues and all
other trade fields; city reserve bypass and character balance/reserve limits;
cloak 0/1/2; zero availability; positive zero-price city offers and character
sellers; cleared and invalid-item trade filtering; hidden zero-price buyers;
duplicate/hidden production notices, zero quantity/month and non-PRODUCE
controls; moving/prisoner filters with a captive's eligible descendant;
absent item metadata, its real plural warning and nonallocation; exact signed
wide columns, minimum-width-not-truncation, markup and trailing spaces.
`economy raw` supplements real save_box bytes with sort and every trade field.
MASTER/style/day/indent output bits and a renamed cached inode are observed
with existing readfile/life flush primitives; no observation closes the cache.

At **0.81.0-alpha**, independently accepted by the parent, actual pinned C/C/Go
matches **925 sections / 63547 lines / 1195280 bytes**, including raw streams,
statuses, report bytes and state. The inherited **916 sections / 62819 lines /
1177391 recording bytes** and **20078 input lines / 368722 input bytes** are
exact unchanged prefixes. New inputs total **20305 lines / 373845 bytes**.
Recording SHA-256:
`a9a0ea8d3d58c13dcdf8d808fbec9355950928a6c3ad2bad38f8d59c81b97b2d`.
The independent #168 live native-order check also matches raw C/C/Go
(**28 lines / 426 bytes**) over the same unchanged directory.

A temporary linker negative control returned success after emitting a
missing-C-dependency marker from market_report. `--update` rejected it with
status **2**, `section 916 reached a C stub`, before replacing any recording;
the recording hash stayed unchanged. Its sources/binary tree were removed.
Log: `tmp/issue141-negative.log`.

### G3 #358: deliberately not a passing failure-parity case

[G3 #358](https://github.com/playbymail/olympia-g3/issues/358), assigned to
the owner, remains open and explicitly **nonblocking by owner decision**.
The investigation executed real io.c read_boxes/load_box/trade_list_scan,
then the real market_report → buyer_list → list_market_items → reduce_qty
path, not a model or an investigate_possible_trade inference. Loaded positive
quantity/zero-price character BUY data is accepted and divides by zero even
with no gold. On this Linux amd64/GCC 12.2.0 host C exits **136/SIGFPE**, twice
deterministically. Real v_buy zero-quantity cancellation and a loaded price-1
control instead exit **0**, twice each. Normal validated BUY commands have
not been shown to produce the malformed positive-quantity zero-price record;
no corpus prevalence is claimed. Prior evidence and exact source/build/run
provenance are retained in `tmp/issue141-probe/` and the upstream issue.

After implementation, an isolated malformed-input probe again ran the actual
C report twice (both **136**, raw streams/status/log identical). An equivalent
temporary Go test executable invoked actual marketReport and exited **2**
with `runtime error: integer divide by zero`. **Those raw bytes and statuses
do not match; this is not recorded or counted as a passing C/Go case.** The
original division remains in reduceQty, with a source-linked comment; there
is no guard, row suppression, successful substitute or new panic recovery.
olympia.go's existing host-dependent-failure comment is clarified only;
engine failure handling is unchanged. Temporary test source/executable were
removed. Evidence: `tmp/issue141-malformed-c{1,2}/`,
`tmp/issue141-malformed-cases.txt`, `tmp/issue141-malformed-go.{stdout,stderr,status}`.

Routine command/move/helper/RNG/code/loadsave/turn-driver/default-turn tools
passed sequentially with explicit `G3_REPO=$PWD/tmp/upstream.git`, bounded
engines, then bounded vet and fresh bounded default Go tests. Code-compat
retains its deliberately unaffected older pin. Logs are
`tmp/issue141-{command,move,helper,rng,code,loadsave,turn-driver,turn,vet,tests}.log`;
recording refresh evidence is `tmp/issue141-command-update.log`. These are
bounded small-case checks, not a new full-world/full-corpus or stage-exit
acceptance claim. Parent independently reviewed the complete pinned C and Go
diff, verified both extracted definitions unchanged exactly once and every
inherited input/recording byte unchanged, then reran all eight tools, bounded
vet and fresh bounded default tests. Parent logs: `tmp/parent141-*.log`.

## Portable admission recording and live native order (#168)

This is a local harness correction, prerequisite to #140, not an engine or
upstream game regression change. On this Linux filesystem, unchanged pinned
C twice and Go traverse `q2,p4,z0,b3,a1`; the inherited recording baked in
`b3,z0,p4,a1,q2`. Native readdir order is correct and must not be sorted or
normalized into portable saved expectations.

Only section **793** (zero-based, the 794th section) changes. Its five
zero-through-four-nonblank-field inputs are unchanged. It now calls real
admission traversal with one candidate at a time, using the existing
`life remove` fixture operation between candidates. All five partial reads,
blank-line cleaning, saved records, and the no-admission tail still execute.
The other **902 sections** are byte-identical in both input and output:
the **329029-byte input prefix / 32759-byte suffix** and **1007603-byte
recording prefix / 151856-byte suffix** are unchanged. No other saved byte
is replaced, including diagnostics, whitespace and source assertions.

`run.sh` also creates one small five-account live fixture, then runs two
fresh C processes and Go over that **same unchanged input directory**.
The existing driver invokes the complete unchanged extracted pinned
`make_new_players` through `admission probe`; no new C source closure or
driver action is needed. C/C compares stdout, stderr and status separately
and byte-for-byte. The Go test compares the existing raw stream/status
envelope solely against those real C bytes, never modeled expected output.
All engines use `run_limited`; relative paths after the same `lore_root`
give both C processes and Go identical actual cwd paths. Recursive snapshots
check file bytes, modes, paths, directory membership and native entry order
after each engine; neither admissions nor fixture setup may mutate the
shared tree between engines. Dependency markers on either live C stream
are rejected before recording updates; C must exit zero.

`TestCommandAdmissionNativeOrder` compares raw C/Go output regardless of
the filesystem's native order; even naturally alphabetical order is valid.
On this filesystem the observed order differs from alphabetical controls
`a1,b3,p4,q2,z0` and `z0,q2,p4,b3,a1`.
Temporary Go build overlays that alphabetized traversal or reversed native
traversal both failed at the raw C/Go comparison (test exit 1); overlays were
removed and the unchanged engine passed again. No production engine file
was edited. Plain `go test` explicitly skips only this live test when
`COMMAND_NATIVE_DIR` is absent; all portable recorded probes always run.
The tool sets that variable to its generated `WORK/native` directory.

Validation at **0.79.1-alpha**:

- `G3_REPO="$PWD/tmp/upstream.git"` explicitly for both `--update` and
  ordinary command comparison; pin remains `e733acb63bacb21167a1f2589d8dcdc023364896`.
  Each passed **903 recorded sections / 62047 lines / 1160604 bytes**, plus
  the **28-line / 426-byte** live C/C/Go comparison with unchanged input.
  New input is **19745 lines / 362373 bytes**. Recording SHA-256:
  `e8fedf667b0a699da4bc9d0d1bb542fe7f9f23f8ddc9a4ae005d63faa2cafb5a`.
- Bounded `go vet ./...` and fresh bounded `go test -count=1 ./...` passed,
  as did shell syntax and `git diff --check`. Update/default/native logs and
  intentional sorted/reversed failure logs are under `tmp/issue168-*.log`.

No engine/report behavior, upstream harness, source pin/dependency set or
shared upstream clone is changed. Route rendering remains separate in #140.
This does not claim Darwin/native Windows, full-world or stage-exit parity.

## Route rendering (#140)

dir.go adds complete dir.c `list_exit_extras` **839–864**,
`list_exits_sup` **879–930**, `list_road_sup` **933–970**, `list_exits`
**973–1031** and `list_sailable_routes` **1034–1072**. It reuses real
exitsFromLoc and the existing names, topology, knowledge, seasonal, output/cache,
capitalization and plural helpers. No movement handler or location-report root
is added. All C definitions are extracted unchanged exactly once (verified
against git show at the pin); no source dependency or pin changes. The direct
fixture wrapper calls the real static helpers with supplied exit_view fields,
including transit/extras combinations not constructed by extra_routes, and
observes nullable/empty/consumed header state and both indentation fields.

Thirteen appended sections cover hidden/known routes, see_all 0/1/2 and IN
filtering, named/unnamed destinations and regions, first-road blank lines,
singular/plural/zero/negative distances, roads, transit, barrier and Hades toll,
empty exits, nonship sailing no-op, ocean/port/mountain/IN filtering, road-only
sailable routes, Uldim/summerbridge precedence and seasonal boundaries. The
invalid summerbridge retains dir.c:1029's assertion/status 134; other added
sections require status 0. Actual save_box state, MASTER/style raw readfile
bytes and paired no-call aftermath (20 matching shuffle/NPC observations)
are retained. Extras intentionally reset second_indent through real wiout;
ordinary exits include transit, while road transit renders impassable.

The initial recording-prefix audit found the inherited filesystem-dependent
admission diagnostic order, reproduced on untouched main and fixed separately
in #168. The accepted #168 input/recording prefixes are preserved exactly;
only the 13 route sections and their 772 real-C recording lines are appended.
Independent parent acceptance at **0.80.0-alpha** passes C/C/Go at
**916 sections / 62819 lines / 1177391 bytes**, plus the live native traversal
comparison. The **1160604-byte / 62047-line** recording prefix and
**362373-byte / 19745-line** input prefix are byte-exact. Recording SHA-256:
`c3262a550b22427ee79da758c25cc9b5a042cab82d70be85c64b6d669d2e6703`.
Parent independently reviewed the complete pinned C and Go diff, verified all
five extracted definitions unchanged exactly once, and passed sequential
command/move/helper/RNG/code/loadsave/turn-driver/default-turn comparisons,
bounded vet and fresh bounded default tests. Logs: `tmp/parent140-*.log`.
No full corpus/full turn, lifecycle or strict stage-2 exit is claimed.

## Noncharacter descriptions (#139)

display.go now implements complete stage-pin display.c `ship_cap_s` **77–89**,
`display_owner` / `incomplete_string` / `liner_desc_ship` **208–272**,
`liner_desc_loc` **288–333**, road/storm descriptions **480–521**, and the
noncharacter branches of `liner_desc` **528–544**. Character descriptions and
accepted headers are unchanged. Existing safe-haven, hidden, mine-depth,
road-distance, storm, names and output/cache helpers are reused. No existing
ship-capacity/ship-weight helper existed: weights.go adds complete u.c
`ship_cap` **2599–2607** and `ship_weight` **1525–1553**, reusing the real
unit weights and character-only recursive traversal. Macro fields retain their
C widths, including unsigned damage, signed ram/level/strength and 32-bit
percentage intermediates. display_owner retains first_char_here's direct-list
semantics and its original macro assertion at u.c:2395.

The command extractor removes four rejected noncharacter dependencies and
extracts the complete unchanged definitions, plus u.c `first_char_here`
**2390–2406**. An independent comparison against git show at the immutable pin
checks that all eleven definitions occur unchanged exactly once. No rewritten
liner_desc or successful rendering stub is used; show_loc and show_loc_posts
still abort with dependency markers, rejected before status checks/recording.
All source files were already dependencies: no pin/dependency-set change.

Thirteen small appended cases distinguish finished/unbuilt ships; absent,
zero, empty, negative and overloaded capacity; recursive noble/inventory
weight versus excluded building contents; completion/load arithmetic widths;
damage-byte wrapping, defense/fee int limits and ram signed-byte limits;
building depth/level/safe haven/completion; city and Hades-pit suffixes;
distinct ship/location banner markup and display gating; hidden roads with
same/unequal province, singular/plural/zero distances; and absent/owned/unowned
storms, signed strength and all eight directions. Unsupported kind and invalid
owner controls require original assertions/status 134; all other new cases
require status 0. Actual save_box observations retain original records and
lists before/after reads. MASTER/style/output-bit and independent raw readfile
observations include **103-byte** banner-visible and **62-byte** banner-hidden
logs, using the accepted owner rather than a second output implementation.

Two actual unchanged C runs match Go at **903 sections / 62031 lines /
1160307 bytes**. A fresh non-update command comparison also passes. The entire
inherited **890-section / 61457-line / 1149073-byte** recording prefix and
**19471-line / 356567-byte** input prefix are byte-exact. Current output SHA-256:
`9d4c3cc9f920c22d921dfe865c875c13b0e2d73a5c74d8a74488aa8818ce8f40`.
Seven added save_box here-list lines retain C's literal trailing spaces.
Initial fixture construction incorrectly omitted the noble base item; C's
139 was rejected before recording, then the fixture was corrected to allocate
the actual item_peasant (10). No engine defect or blocking upstream defect
was found; existing signed comma_num formatting is preserved, not normalized.

At **0.79.0-alpha**, sequential command/move/helper/code/loadsave/turn-driver/
default-turn/RNG comparisons, bounded vet and fresh bounded defaults pass.
All tools use explicit `G3_REPO="$PWD/tmp/upstream.git"` and engine run_limited
supervision. The nine usual opt-in/child-entry test skips are not full-world
checks. Logs: `tmp/issue139-command-{update,default}.log`,
`tmp/issue139-{move,helper,code,loadsave,turn-driver,turn,rng}-default.log`,
`tmp/issue139-{vet,go-test,prefix-check,extract-check}.log`.
Matching command captures are removed after comparison. No movement handler,
location-report root/posts, full corpus/full turn or stage-exit acceptance is
included; strict stage-2 report-tree parity remains deferred.

## G3 #357 follow-C and pin move (#131)

Pin-move recording: **859 sections / 59841 lines**, from two identical real
C runs. All original input sections remain byte-exact; their old descriptive
names are retained even where the bug they named is now fixed. The appended
savage-cap case probes the current NPC stream before an explicit source-
location reseed, distinguishing the pre-begin_npc cap from a post-draw cap.
Existing cases already distinguish unpaid NP aborts, unowned knowledge
assertions, the 199→200 savage tally and regular-holder token notices.
The only inherited recording changes are those four fixes (including their
state/output-bit/mint aftermath) and art.c's +3 assertion offsets.

USE 886 / v_curse_noncreat was never ported: glob.go's USE callbacks remain
explicit NYI, and the C oracle does not extract this validator. Its stage-4
implementation must retain the minimum c->b=1 and possession/destroyable-item
checks, then print `The noncreator possession curse has no effect yet.` and
return FALSE, never reaching completion or charging aura. The turn's existing
noncreator_curse_erode NOTYET log path is unchanged, not successful casting.

The implementation/evidence sections below are historical at their stated
pins. Their old hashes and reproduced-bug descriptions are retained, not
claims that G3 #352–#356 remain unfixed at the current pin. See also
[the affected-tool pin-move notes](../README.md#stage-2-pin-move-and-follow-c-131).

```sh
tools/command-compat/run.sh
tools/command-compat/run.sh --update # ONLY real C generates the recording
```

Requires Bash, git, Python 3, a C11 compiler (`CC`, default `cc`, with GNU
section GC or Darwin dead stripping), Go, Perl (timeout supervision), cmp and diff.
`G3_REPO` defaults to `~/Software/playbymail/olympia-g3`; a bare cache containing
the pin also works. `G3_REV` deliberately overrides the pin.
`WORK` defaults to `tmp/command-compat` and is emptied first.
Only archives under `tmp/` are compiled; no upstream working tree is modified
or built. Each C driver and Go test binary runs under shared `run_limited`,
via `tools/lib/timeout.sh` / `run-limited.sh`, with `RUN_TIMEOUT` seconds.

## Empty name/plural oracle inputs (#167)

The local C driver's `name N [TEXT]` and `plural N [TEXT]` operations clear
their reused text buffer before scanning and require an integer ID conversion.
Omitted text now reaches the complete unchanged real `set_name` or `str_save`
as an empty string, rather than retaining the previous operation's text.
This corrects the oracle parser only; no engine, source pin or C dependency
changes. The separate #137 names audit is not part of this prerequisite.

Three appended sections use existing operations and actual `save_box`,
`plural_item_box` and `box_name_qty` observations. They distinguish clearing
a nonempty name, clearing a nonempty plural, and an unrelated `field` token
before an empty plural. Before the fix, C retained the old saved name/plural
or installed `weight` as the plural. Afterward, cleared saved fields disappear,
plural quantities fall back to the singular name with exactly two warning
lines per plural-quantity observation, and singular quantities stay silent.

Two actual pinned C runs are identical at **877 sections / 60622 lines**;
the Go comparison passes. The inherited **874 sections / 60548 lines** and
entire input prefix are byte-identical. Only C generated the 74 appended
recording lines; dependency-marker checks precede updates. Four temporary
negative controls (missing/nonnumeric ID for each operation) returned status
0 before and are rejected with status 2 after the parser correction.
Fresh default command and turn comparisons, bounded `go vet ./...` and
fresh bounded `go test -count=1 ./...` pass. Logs:
`tmp/issue167-command-{before,update,default}.log`,
`tmp/issue167-parser-negative-{before,after}.log`,
`tmp/issue167-prefix-check.log`, `tmp/issue167-turn-default.log`,
`tmp/issue167-vet.log` and `tmp/issue167-go-test.log`.

## Location headers (#138)

The existing display.go owner now ports complete unchanged display.c
`loc_inside_string` **8–49**, `show_loc_barrier` **52–63**, `safe_haven_s`
**66–74**, and `loc_civ_s` / `show_loc_header` **125–166** at the stage pin.
The C extractor adds those exact definitions, exposing the static barrier
through a fixture-only wrapper. Only the rejected `show_loc_header` dependency
is replaced: `show_loc` and `show_loc_posts`, and noncharacter `liner_desc`
dependencies still abort explicitly. No all-location root or duplicate state
projection is introduced. Existing mineDepth, safeHaven, locHidden, topology,
names and sout/cache owners are reused. display.c was already tracked, so
neither the dependency set nor any pin changes.

Seven appended asymmetric sections distinguish direct-province, ocean and
nested-building containment; sublocations in land/ocean; raw named/empty
regions; signed mine depth with truncation; depth → safe haven → hidden → civ
suffix order; wilderness/signed civ; Faery/Hades/ocean/nonprovince suppression;
and tunnel's direct civ string versus header suppression. Raw record-presence
observations before/after rendering and actual save_box output check absence
of unintended allocations or state changes. MASTER/header/banner/include
section routing uses the actual output owner, with style, both indents and day
fields. Real nonzero/zero barrier returns, prisoner suppression and output bits
are observed alongside independent readfile bytes without flushing production
streams. The two invalid controls retain exact code.c:228 and u.c:388
assertions, in original order.

Two actual deterministic pinned C runs, followed by Go, match **890 sections /
61457 lines / 1149073 bytes**; default non-update comparison also passes.
The entire inherited **883 sections / 61012 lines / 1140518 bytes** output
prefix and **19305 lines / 353182 bytes** input prefix remain byte-exact.
Inputs now total **19471 lines / 356567 bytes**. Output SHA-256:
`4109632dbaf4950a4f15b27515c51a534e002909d5db0b6a30e9bf029a9140d9`.
For example, the two real MASTER log observations each contain **132 bytes**,
including `0:7::3/8:17:` versus `1:0::3/8:17:` headers and barrier lines;
the normal character log remains **43 bytes** after a prisoner barrier call.
Missing-dependency rejection still runs before expected-status checks and any
recording update. All newly expected failure statuses are explicitly 134.
Matching command captures were removed after comparison. Source/documentation
diff whitespace checks pass; five new save_box here-list lines retain their
literal C trailing spaces, not normalized output.

At **0.78.0-alpha**, command/move/helper/code/loadsave/turn-driver/default-turn
and RNG comparisons pass sequentially, each with explicit bare G3_REPO and
run_limited engine supervision. Code retains its existing pin; its complete
source/header dependency set is identical at the stage pin (checked with git
diff). Bounded vet and fresh bounded default tests pass (Olympia 1.780s);
the ordinary nine opt-in/child-entry skips are not extra whole-world checks.
Logs: `tmp/issue138-{command,command-default,move,helper,code,loadsave,turn-driver,turn,rng}.log`,
`tmp/issue138-{vet,go-test,prefix-check,code-source-compat,extract-check}.log`.
No upstream/port defect or blocker was found. Full strict report-tree parity
remains pending later stage-2 tasks and is not claimed or normalized here;
no full corpus, full turn or stage-exit acceptance was run for this task.

## Bounded names/codes audit (#137)

Audited complete unchanged pinned `code.c` definitions at **47–174,
212–343 and 350–466** against the existing `code.go` owner. No engine gap
was demonstrated; no naming layer, engine implementation or source dependency
was added. The parser prerequisite is separately accepted #167 above.

Existing oracles already cover the following, so their inputs and recordings
are retained rather than duplicated:

| oracle | retained discriminating coverage |
|---|---|
| code | Every numeric code-range transition, signed-int extremes, zero and invalid shapes, scode's single opening-bracket removal, exhaustive box-space conversion checksum and designated string/shape checksums; missing/deleted name assertions; bracket substitutions and replacement; absent/empty/replaced banners and 50/51/60-byte truncation. |
| helper | Kind/subkind fallback, invalid/deleted/unnamed/empty names, noble-item recursion and missing items, capitalization/high bytes/long names, absent/null/empty plural names and warning bytes. |
| move | Real port topology; character noble/item fallback and warnings; banner absent/empty/truncation and display flags. |
| command | Ordinary city versus port city and mountain suppression, all three storm kinds, default subkinds and quantity helpers through real inheritance. |

Six appended sections fill the remaining direct-observation/composition gaps:

1. Zero, negative/out-of-range codes and the garrison sentinel, both missing
   and explicitly named. `box_code`, `box_name` and `just_name` return
   `Garrison`; `box_code_less` still returns `999`, and `display_name` is
   independent of that sentinel shortcut.
2. An unnamed, named, then explicitly emptied port city. Display fallback is
   `City`, never `Port city`; kind remains `port city` independent of name.
3. Noble-item singular fallback and capitalization, explicit noble-name
   priority, empty item-name fallback and missing-item code-only fallback.
4. Absent item record, null plural pointer and explicit empty plural: all
   four plural/quantity primitives emit their real warning bytes for each
   nonsingular call (12 new warnings). Singular calls stay silent. set_name
   substitutes brackets; directly stored plural names preserve brackets.
5. Exact tilde separators and quantity formatting at 0/1/10/11/1001/-2 and
   both signed-int extremes; negative formatting is not normalized.
6. 2040/2048/2049-byte names through nested formatting. At 2049 bytes the
   direct display/just/singular-plural names remain 2049 bytes; formatted
   box/quantity results truncate to 2048. Noble-item capitalization itself
   truncates the fallback to 2048 bytes before outer formatting.

`names show N Q` observes each actual primitive separately; `names repeat N
LENGTH` constructs a bounded raw input and calls real set_name. Neither
models a C naming routine. Whole unchanged code.c remains linked, with real
cap, nice_num and port dependencies. Names sections require status 0;
dependency markers are rejected before recording updates.

Actual C/C/Go passes at **883 sections / 61012 lines**. The entire accepted
877-section input prefix (**351278 bytes / 19225 lines**) and output prefix
(**1075000 bytes / 60622 lines**) are unchanged. C alone generated the new
**65518 bytes / 390 lines** after two deterministic runs. Recording SHA-256:
`57aa1797f37841a3c4d46aa31a0cd398f95cba7e457bdfac9ff95b4875661566`.

Sequential code **43/969**, helper **37/1487**, move **135/7158**, command
**883/61012**, default turn, bounded vet and fresh bounded default Go tests
pass. The code tool retains its independent `6a25f2d` pin: all its upstream
dependencies are unchanged at the stage pin, including code.c blob
`1a9996f1c7ea3da603a35a3c534a00549c4375a7`. No full-world/stage-exit or
later report task acceptance is implied. Evidence:
`tmp/issue137-source-audit.log`, `tmp/issue137-prefix-and-byte-audit.log`,
`tmp/issue137-command-audit-{update,default}.log`,
`tmp/issue137-{code,helper,move,turn-default,vet,go-test}.log`.

## Bounded raw file observations (#135)

`readfile PATH` observes a whitespace-free cwd-relative (or absolute) path,
without flushing, closing, resetting or reopening any production output owner.
It opens and closes only its own independent read handle. Regular files are
limited to **65536 bytes**; nonregular files, oversized reads and other I/O
errors fail the harness, never truncate or masquerade as an empty file.
ENOENT/ENOTDIR print `file bytes PATH missing`. Existing files print
`file bytes PATH LENGTH:` and a newline, exactly LENGTH raw bytes, then the
framing `\nend file bytes\n`. The extra framing newline is not file content;
the explicit length distinguishes empty files, blank lines, trailing spaces
and absent final newlines. There is no whitespace or diagnostic normalization.

Use existing `life flush` for intentional post-flush observations. `readlog N`
retains its historical explicit global fflush followed by unframed log bytes;
the Go interpreter now explicitly mirrors that flush even with buffering on.
Neither operation implements production rendering. Existing `report` actions
already invoke complete pinned C roots; Go's deferred projections remain in
their existing owners. No speculative report-leaf catalog is added.

Three appended cases cover raw blank/trailing-whitespace bytes, empty/missing
and ENOTDIR, buffered report output before/after explicit flush, repeated
observation without cache reset, renamed cached-inode writes, actual saved-log
bytes, and the intentionally flushing legacy readlog. The sole new leaf action,
`report location N`, routes the existing `char_rep_location(N)` result through
real `out`, reusing Go's accepted charRepLocation/commandOut owners. Complete
unchanged report.c:357–407 was already extracted for death logging; no C
dependency or rendering implementation is added. Cases distinguish `nowhere`
from a named province after relocation. Ordinary routing also preserves its
two trailing spaces. Future auxiliary and Times body cases use the same observer;
their deferred rendering is not implemented or claimed here.

New `section file observer` cases require C status 0 in addition to the existing
missing-dependency rejection, including in `--update` mode; an assertion cannot
silently become a successful observer recording.

At **0.76.6-alpha**, two actual pinned C runs before recording and Go match
**862 sections / 59951 lines / 1040469 output bytes**. The previous **859
sections / 59841 lines / 1038509 output bytes**, and **18885 input lines /
344932 input bytes**, remain exact unchanged byte prefixes. New input totals
**18934 lines / 345845 bytes**. Recorded output SHA-256:
`e9d1a9c2bf62020280ff8c19109fe84906c31ae5a84e44f6d6c38b606140e423`.

A temporary one-byte mutation of C-recorded report text (`nowhere` → `Nowhere`)
made the bounded fresh Go comparator fail at line **59877**, with the exact
different bytes quoted. The original C bytes and hash were restored. A separate
temporary C driver probe emitted the rejected missing-dependency marker and
returned success in section **859**; `--update` rejected it with status **2**
before replacing the recording (hash unchanged). Neither probe remains.
Logs: `tmp/issue135-command-update.log`, `tmp/issue135-negative-report-byte.log`,
`tmp/issue135-negative-dependency.log`, `tmp/issue135-validation.log`.

Fresh bounded `go vet ./...` and `go test -count=1 ./...` pass (Olympia
1.767s); test binaries use `tools/lib/run-limited.sh`. Sequential default
command comparison repeats C deterministically and passes all 862 sections;
default turn passes its comparator contracts and exact raw streams/status
and two cwd entries (types/modes/bytes). All tool calls explicitly use
`G3_REPO="$PWD/tmp/upstream.git"`. No full corpus or strict whole-world turn
is rerun for this harness-only task; the stage-2 report differences remain.
Engine code, source pin, dependencies and existing recordings' prefixes are
unchanged. Parent acceptance independently repeated command C/C/Go, vet,
fresh bounded tests and default turn checks; evidence is in
`tmp/parent135-{command,vet,tests,turn}.log`.

## MASTER, styles and section primitives (#136)

At **0.77.0-alpha**, twelve appended sections exercise complete unchanged
real sout.c:273–283,448–582,620–719, already linked whole. No extraction,
pin, dependency or successful stub changes are needed. Go extends the same
sout.go/world.go owner: format once into the LEN-1-byte text before routing,
prisoners before MASTER, one previous style (not a stack), vector ownership
transfer, and literal section/header bytes. bottomOut/openPlayerOutput still
open before p_player and retain the existing buffered cached inode lifetime.
The admission projection now calls vectorPlayers; the account projection
retains the three original HTML/previous-style transitions, without adding
its deferred body. wout is out, with no active wrapping engine at this pin.
Raw html() styles are implemented; separate -h generation is not.

Fixture-only `sout` actions select route/frame/style, vector clear/add/players/
save/restore and state observation. `sout text WHO KIND IND COUNT TEXT` repeats
bounded fixture text (at most 16384 bytes, `|` means newline) and calls actual
out/wiout/html/lines/match_lines/log_write. It does not supply expected output
or replace a C dependency. Existing readfile's independent 64-KiB/no-flush
observer, life flush and readlog are reused unchanged.

Cases distinguish normal/MASTER/vector, duplicates, prisoner suppression,
location allocations and garrison fanout, silent players and scanner/skill
exclusions, all eight
frame variants (unit/day/secondary), negative indents, wiout's secondary reset,
repeated PREV and nested HTML, comma/empty/long headings, exactly 72 dashes,
2046/2047/2048/3000-byte inputs, unwrapped lines, raw colon/tilde/newlines,
output bits, transfer/restore including empty vectors, log_write restoration
of both ordinary and MASTER fields, prisoner GM, and buffered renamed-inode
writes. Invalid vectors and log categories retain original assertion locations
457/708. Every new successful case requires status 0; only the two named
invalid controls require 134, after missing-dependency rejection.

Two complete real C runs are compared before every recording replacement.
Final C/C/Go: **874 sections / 60548 lines / 1073778 output bytes**; inputs
**19198 lines / 350721 bytes**. The entire previous **59951 lines / 1040469
output bytes** and **18934 lines / 345845 input bytes** remain unchanged byte
prefixes. Output SHA-256:
`260a5fea60c89dffdd78241e681bf41b0a989f6320b21e579ce043c08cd101f0`.
A Go-only Stringer probe guards the single-format-before-fanout boundary,
including duplicate, prisoner-only and empty vectors.

Negative controls: a temporary missing-dependency marker with successful C
exit is rejected in --update at section 0, status 2, with recording hash
unchanged (`tmp/issue136-negative-dependency.log`). A one-byte Title→title
mutation of the recorded HTML name fails the Go comparator at line 60197
(`tmp/issue136-negative-byte.log`); original bytes/hash were restored. Neither
probe remains. Initial fresh tests exposed unallocated recipients in the new
Go-only test setup; allocating those fixtures fixed the test, not engine code.

Sequential command/move/helper/code/loadsave/turn-driver/default-turn checks
pass, all with explicit `G3_REPO="$PWD/tmp/upstream.git"` and run_limited engine
supervision. Bounded vet and fresh bounded defaults pass. Evidence:
`tmp/issue136-command-{preupdate,update}.log`,
`tmp/issue136-{command,move,helper,code,loadsave,turn-driver,turn,vet,test}.log`.
Source diff whitespace checks pass; trailing spaces in the C recording are
literal preserved data. Parent acceptance independently repeated these checks
(`tmp/parent136-*.log`) and the strict full turn. Two fresh C turns match;
C/Go both exit 0 with identical raw streams and no path/type/mode differences.
Exactly the same eight saved-report byte differences remain deferred renderer
work, not a claimed green stage-2 gate (`tmp/parent136-turn-full.log`).
Full corpus was not rerun. No later renderer or external lifecycle feature is
implemented.

## Reference and stub boundary

The script extracts complete functions from **real pinned input.c**, retaining
original `#line` locations: parser, `cmd_shift`, `check_allow`, `get_command`,
`set_state`, `load_command`, `command_done`, `finish_command`, `do_command`,
and `init_load_sup`, precedence/comparison and both queue sorts,
`initial_command_load`, `min_pri_ready`, `start_phase`, `evening_phase`,
`process_player_orders`, and `daily_command_loop`. It also extracts the real
`init_wait_list` and complete `check_all_waits`, plus produce.c's
`init_collect_list` and original `collectors` storage, plus combat.c's
complete `auto_attack`, `check_auto_attack_sup`, `check_all_auto_attacks` and
`clear_second_waits`, plus `fort_covers`. Wrappers expose static helpers to the
driver; no compared routine is rewritten or stubbed. It links the real
`glob.c` table (including original handler pointers), `code.c`, `z.c`,
**sout.c**, and list/allocation libraries. The raw command output is therefore
real `out`/`wout` → `vout` → `out_sup` → `bottom_out`, not a printf stub.
The driver observes output files through the `log` operation and at the actual
combat boundary, after flushing C stdio. stdout, stderr and status are recorded
separately, without normalization. Two fresh C runs must be byte-identical before recording or
comparing. Default Go tests use the recording without C or saved games.

Real complete handlers: `v_stop`, `v_xyzzy`, `v_plugh`, `v_format`, `v_notab`,
`v_hide`, `d_hide`, `v_wait`, `d_wait`, `i_wait`, `v_move`, `d_move`,
`v_exit`, `v_enter`, `v_north`, `v_south`, `v_east`, `v_west`.
Real dependencies include order-head/top/pop operations,
player/lord and stack lookup, `all_stack`, `count_any`, `count_stack_any`,
`rp_skill_ent`, `has_skill`, `exp_level`, `check_skill`, parser garrison lookup,
`here_pos` and sparse output-bit operations. Target predicates are extracted
whole from perm.c (`is_hostile`), u.c (`contacted`, `char_where`, `char_here`,
`test_known`, `loc_hidden`, `has_item`, `loc_depth`), loc.c (`in_safe_now`,
`province`, `building_owner`), and storm.c (`weather_here`), using
the actual oly.h macros. The driver only constructs state and calls these
functions; it does not emulate their decisions.

Wait comparisons extract complete c1.c `clear_wait_parse`, `parse_wait_args`,
`wait_leader`, `check_wait_conditions` and the small real `flag_raised`
dependency, plus the original `wait_tags`, flag structure/storage and u.c
number-name table. Real u.c `lookup`, `nice_num`, `comma_num`, `gold_s` and
code.c `just_name_qty` provide the exact diagnostics; existing visibility,
inventory, weather and location helpers are reused. The `signal` driver
operation constructs flag records only; it is not a replacement `v_flag`.
`wait` in the case language clears/reparses and checks conditions directly;
it does **not** execute `v_wait`, `d_wait` or `i_wait`. `wait_check` reevaluates
the cached parse, while `wait_parse`/`wait_dump` observe even partial parses.
The handler cases use real queued `wait` orders through load/start/evening
dispatch, or construct persisted commands and call `init_wait`/`poll_wait`.
`interrupt` invokes the table's actual interrupt callback; `generic_interrupt`
calls the complete real `interrupt_order`. `lists` observes wait/collector/second-wait
storage and the selected command's second-wait flag. `wait_add` only constructs
list state, including stale entries and duplicates; it models no decisions.
Collector cases construct saved COLLECT commands and call `init_collect`;
they never execute an aborting COLLECT handler or claim harvesting works.

All other command routines named by the C table are **aborting** stubs,
never successful substitutes. The remaining aborting dependencies are
unrelated noncharacter descriptions and unported command/combat handlers.
`is_port_city` and trade settlement/dispatch are real extracted definitions.
Reaching any
stub fails the script, including `--update`. Only instrumentation
(`coverage_count` and phase markers) is a no-op: stage 3, not gameplay state. No compared
scheduler, auto-attack scan, wait predicate or wait handler is a successful stub.

Positive selection controls are named `section combat boundary ...`. The only
accepted combat boundary is an **aborting** `regular_combat`: actual extracted
`auto_attack` has already written the target's code to the real output log and
set the player's output bit. The boundary dumps those observations, emits a
harness diagnostic and aborts (134); it never returns fake combat success.
Both missing expected combat and unexpected combat fail, including `--update`.
Go exercises its production `regular_combat` not-implemented panic, then the
test observer reads its real log/output bit in the same format. The recorded
134 is the C harness abort, not a claim about the full C engine's combat result
or Go's CLI exit code. No production test hooks replace combat.

## Production boundary

Go preserves each saved command-table index and the C layout of
`cmd_tbl_ent` (three function pointers). NULL stays nil; all unported
non-NULL handlers raise the engine's explicit not-implemented error, including
collect. Tests compare presence of all three pointers at **every** table index
against actual C, independently of the existing entity layout/metadata test.
The only implemented handlers are the complete real handlers listed above.
The tiny no-op handlers really return TRUE in C. Skill lookup and counting
are real dependencies, not successful placeholders.

The `commandOut` route implements default character/player and location/garrison
output, prisoner suppression, real raw C line framing and LEN-1 truncation, append-only
log writes and the player's output-bit mutation. The small `logWrite` MASTER
route shares that same bottom writer; code alerts preserve GM output bits
and raw log bytes, not just diagnostics. Immediate output goes to
stdout with C framing, but without command echo; it does **not** bypass all
allow checks. Indentation, secondary indentation and garrison unit prefixes
are preserved. Report assembly and the broader sout interface remain stage 2.
No report-side mutations on these paths are omitted. Full turn execution still
stops at explicit unfinished subsystem boundaries. Wait callbacks, polling, generic order interruption and
the full ownership-transfer closure are ported. COLLECT execution and
validated-trade settlement still stop explicitly. SAIL and COLLECT
interrupt callbacks remain explicit unported boundaries, not successful stubs.
MOVE's interrupt pointer is genuinely NULL in C, so generic interruption can
complete a MOVE without calling its start/finish callbacks. Direction and
ENTER/EXIT table rows have only a start callback: their real reparsing switches
the command to MOVE before completion, not an invented alias finish pointer.
The real auto-attack
scan excludes stacked followers and leaders whose health is neither 100 nor
-1, then completes the local and outer-ship target scans. Healthy nonhostile
leaders return normally. Only a selected attack emits real auto-attack output
and stops explicitly at `regular_combat`; combat resolution is not implemented.
Garrison fan-out now uses real `players_who_rule_here` and `province_admin`;
selected garrison attacks emit actual routed output and stop at
`regular_combat`, just like ordinary units. Prisoners retain C's output suppression.
Clearing second waits is the actual list-based combat.c operation, not a scan
of all characters. Attitude command families remain unported.

## Discriminating cases

- All 148 table entries: actual NULL/start/finish/interrupt presence and
  serialized indexes, in addition to existing entity layout checks.
- Unknown → unknown → valid versus unknown queue exhaustion: recursive
  reporting, `load_command`'s actual return value, status, next-order
  consumption, `last_order_turn` and exact load/run queue membership.
- A comment-only/empty command is valid parse index 0, then **disallowed**;
  it is not reported as unrecognized and does not increment command count.
  A NULL line uses actual C `%s` formatting, not a Go pointer dereference.
- `fomrat 258` distinguishes fuzzy correction and its echo; formatting is
  stored as signed char **2**, while raw output still prints **258**.
  `notab 256` stores zero but prints the nonzero-argument message. Command
  count crosses signed-short 32767. Next orders lower current priority.
- Both condition markers execute despite the opposite previous status:
  pinned C parses/saves `conditional` but never tests it in dispatch.
- NULL/empty allow, player/character, restricted `r`/`g`, GM override, and
  immediate `i` exemption versus still-disallowed non-`i` orders.
- A real failing hide start increments command count, unlike allow rejection;
  learning a skill is not knowing it. A real timed hide start/finish changes
  hidden state, then hide 0 sets wait zero and inhibits its finish routine.
- Persisted hide with wait **3**, poll **0 versus 1**: the poll case mutates
  hidden state before timeout and emits two real finish messages. These
  deliberately constructed saved fields test the generic dispatcher; the
  table's default hide poll remains 0. Failure while stacked completes early
  and loads the priority-0 next order. Inhibition preserves status/state.
- Negative wait remains negative/running; double same-day calls print C's
  diagnostic; day 130 stores a signed-char debug value, not an int.
- A moving follower suspends even with failed status; the wait-ID exception
  decrements wait without calling any wait handler (`inhibit_finish=1`).
  A dead character completes without decrementing wait or consuming orders.
- Shift exposes ninth and tenth numeric arguments instead of always zeroing
  h. Player queues, prisoner output suppression versus real mutations,
  immediate completion preserving queued orders, and all dispatcher NULL /
  invalid-unit / days-executing assertion locations are observed.

- Character-then-player initial loading, day-0 player-only dispatch, then
  character start: real queued orders, report-setting mutations and log bytes.
- Load sorting distinguishes player-before-character, reverse here-list
  position versus ID, equal precedence across locations versus insertion
  order, and stack depth. Run sorting distinguishes priority before stack
  depth and ID ties; next-load insertion exposes the real completion order.
- Ready priority changes independently for prisoners, movement and second
  waits, returning 99 only when all candidates are excluded. Mixed blocked
  and ready units at the same priority also test each snapshot exclusion.
  Bad saved state, priority and missing here-list membership preserve exact
  C assertions; malformed location zero preserves SIGSEGV status 139 and
  unflushed stdout rather than inventing an assertion. This raw-state case
  bypasses `check_db`; G3 #351 now repairs nowhere units at load/save.
- A failed priority-3 hide loads priority-1 format, which executes before the
  second priority-3 command: observable output distinguishes snapshot
  re-evaluation from draining the old queue. Timed hide spans three evenings,
  mutates hidden state and lowers current priority when loading its successor.
- Evening snapshots do not execute newly loaded orders; second waits skip
  both starts and aging, then daily clearing makes them eligible next day.
  A running prisoner still ages; a moving follower ages without finishing.
- Auto gate consumption at priority 99, priority 0/1 then 3, literal zero-char
  world, healthy follower/unhealthy leader exclusions, and no second scan
  after making a leader healthy later that day. Health 100/-1 nonhostile
  leaders complete priority-3 commands. Hostile controls stop after actual
  low-priority state/output mutations but before the priority-3 handler.

- Unit and player hostile lists versus neutral/defend lists, same-faction
  short-circuit, exact target versus target-faction IDs, automatic sublocation
  monsters, NPC subkind and both NPC loyalties, vulnerability exceptions,
  and garrison-only hostile lists.
- Hidden-alone versus inventory/stack visibility, contact direction and
  faction contact, mismatched locations, fog versus same faction/shared stack,
  first-character garrison exception versus a later garrison. Sublocations do
  not acquire province fog's visibility rule. Buildings and ships suppress
  weather; province storms sum in signed C widths, excluding nested storms.
- Safe self/ancestors through nested stacks and ships, unchanged by movement;
  hidden-location flags and faction knowledge independently; exact signed
  first inventory quantity (including duplicate entries); building/ship owners
  selected by direct here-list order, not faction or nested ownership.
- `contacted` really allocates an absent character record. Invalid known/item
  operands and nonbuilding ownership retain exact assertion file/line/condition.

- Complete integrated scans: safe/moving short-circuits even before an invalid
  here-list, prisoners, hidden/fogged units and target contact, exact nonzero
  signed vulnerability-item quantities, unsupported target kinds, same-faction
  exclusions, direct here-list order rather than recursive/ID order, and
  stack-leader conversion (including leader vulnerability and own-stack skip).
  Two deliberately inconsistent saved here-list/parent constructions expose
  conversion and own-stack exclusion without changing the C scan.
- Hidden/known locations, empty/same-faction/enemy owners, building versus
  nonbuilding first-character lookup, nested occupants, and ships bypassing
  location-hidden checks. Outer scans exclude self, unfinished/ghost ships,
  nonships, empty and neutral ships; unlike local scans they do not exclude
  same-faction owners. Local targets precede outer targets. Multiple eligible
  targets prove the first selection and exact precombat output, not combat
  resolution or a second successful attack.
- `is_ship` is the actual subkind macro, not a T_ship kind check. Following
  G3 #351, completed rafts reach the outer combat boundary and unfinished
  rafts can be selected locally. Both have building depth and zero fort cover.
  Every other `fort_covers` subkind and its invalid-subkind assertion are
  real-C controls. A raft initiator with no local target can reach the outer scan.

- All wait predicate families and weather aliases, case-insensitive tags,
  complete and partial parse contents, clear/reparse and cached reevaluation.
  Missing operands for every argument-taking tag, unknown/empty conditions,
  negative operands clamped by the real `parse_arg`, zero/default quantities,
  large signed-int quantities and independent signed-short clocks are covered.
- Ordered implicit OR returns the first true condition, not a conjunction.
  Parsing still validates the entire input before testing the first condition;
  an early true term cannot hide a malformed later term. NOT lasts one term,
  repeated NOT does not toggle, and a trailing NOT alone keeps waiting.
  Literal `OR`, `AND` and `hour` are **unknown conditions at this pin**, not
  extensions we implement. Empty predicates assert in this low-level checker;
  the start handler instead rejects an empty order before predicate checking.
- Unit existence versus hidden/contact/fog visibility; direct `loc` versus
  outer `shiploc`/`ship`, nested stacks on ships, unfinished ships and rafts,
  and kind versus ship-subkind validation. Ferry checks require both local
  presence and a nonzero signed-char horn, not merely a ship's existence.
- G3 #349's formerly crashing `ferry 2147483647`, one-past-table and highest
  in-range missing IDs, zero, missing operand, deleted ship-subkind box,
  ordinary nonship and negative IDs, with same-state unit/loc/ship controls.
  Real `parse_arg` clamps negative IDs to zero. Completed ships/rafts with
  raised horns succeed; unfinished rafts remain invalid ferries.
- Gold/item positive messages show actual holdings, while negative messages
  show requested holdings. Duplicate inventories are not summed; invalid
  item kinds keep waiting (or produce C's negative message), rather than
  receiving an invented item-validation error. Item's optional quantity
  greedily consumes the following token even when it is another condition.
- Stack, top and first-character ownership are distinct; the GM player is
  its own wait leader. Flag lookup distinguishes default faction, explicit
  unit/faction and global zero, case-insensitive first-match order, quoted
  flags and digit-leading zero versus negative optional sender tokens.
- Rain/fog/wind check the **province**, while clear checks the **direct
  sublocation**: a building or ship can be clear while its province is raining.
  Negative storm strength is not positive weather but still prevents clear.

- Real wait starts immediately satisfy, or append and remain indefinitely
  running; same-day polling after another command can finish and load/execute
  a successor. Unresolved polling does not age a command; evenings do.
  Resumed waits preserve elapsed days and lazily reconstruct the predicate
  parse. Saved poll zero does not suppress an indefinite wait's finish callback.
- Both initializers append without clearing or deduplication, in `loop_char`
  order, excluding players and absent/nonrunning/wrong-command records.
  Collector initialization does not execute COLLECT or mutate its command.
  Polling skips absent/nonrunning/wrong-command records before asserting the
  selected wait is -1, preserving the actual input.c:821 assertion.
- The poll loop indexes the **live shrinking list**, not a snapshot or Go
  range. Adjacent satisfied waits leave the shifted neighbor for a later poll.
  A satisfied/unresolved pair distinguishes an incorrect fixed-length range
  by the unresolved command's debug field (0 on the first poll, 2 on the next).
  A three-unit asymmetric case also observes later satisfied entries.
- Empty starts fail without clearing an earlier cached parse. Malformed
  arguments instead emit `Wait finished: <parse diagnostic>` and succeed,
  exactly as C does; no independent validation fix. A deliberately constructed
  resumed empty command bypasses start validation and retains c1.c:1363's
  assertion; this raw state is not evidence of a publicly reachable engine bug.
  The interrupt callback removes **all** matching list entries, but changes
  neither the command nor its cached parse; repeated removal is harmless.
- Second waits suppress evening aging, not the final wait-list poll, so a day
  predicate can finish despite second_wait. The unchanged `clear_second_waits`
  clears only listed commands, including duplicate entries, leaving an unlisted
  flag intact. Daily clearing allows a time wait to age on the following day.

New real-source functions are c1.c:1639–1692 (`v_wait`, `d_wait`, `i_wait`),
input.c:793–825 (`init_wait_list`, `check_all_waits`), and produce.c:525–543
(`init_collect_list`, with its original storage at line 522). The previously
extracted combat.c `clear_second_waits` is reused unchanged.

## Interruption and transfer closure (#62)

Complete real functions added: input.c `interrupt_order` and
`process_interrupted_units`; order.c `is_stop_order`, `stop_order`,
`flush_unit_orders`; swear.c `set_loyal`, `set_lord`, `reswear_all_sworn`;
perm.c `clear_all_att`; use.c `char_np_total`; u.c `add_np`, `deduct_np`;
move.c `restore_stack_actions`; day.c `touch_loc`, `touch_loc_pl`;
loc.c `viewloc`; seed.c `int_comp`. The full real io.c writer provides
`save_box` and its serialization dependencies, linked with admits/trades/req
lists. Only `phase_ran`/`phase_acted` are instrumentation no-ops (stage 3).
No transfer effects or command callbacks are replaced by successful stubs.

`transfer`, `stops`, `flush`, `stop_check` and `np_total` call those functions.
`save` calls **actual C save_box**, compared against Go's existing saveBox;
`pstate` observes unsaved NP gained/spent, output bits (including LOG_CODE),
touched locations, units and exact remaining order strings. `readlog` reads
the real specified player's raw sout log, including GM code alerts.
Additional field/order operations only construct state.

- Two ordinary players transfer a running WAIT unit with two skill entries
  (learning and unknown, each np_req 5). Refund is **14**, charge **19**:
  oath rate 3 is read before loyalty changes to rate 8. Old queued hide/move
  orders must be flushed before interruption, never executed or loaded;
  the new owner's format order is loaded. WAIT list removal, retained cached
  parse, both priority queues, serialized command fields and logs are observed.
- A dependent is resworn to the **old direct character lord**, not its faction;
  its own loyalty, pledge and attitudes stay intact. Transfer clears who's
  attitudes and pledge, invalidates backlinks, restores stack movement,
  removes old-player membership, sorts new membership and touches both a city
  and its province. Same-player reparenting preserves these effects instead.
- Loaded and absent commands, ordinary → independent → different ordinary
  player, same-previous-player reacquisition, zero and nonplayer prevLord,
  detachment to lord zero, NPC loyalty transitions, negative costs and signed
  short overflow are recorded. Deliberately raw nonplayer prevLord and
  negative/overflow fields are low-level contracts, not public-engine defects.
- Insufficient funds leaves the new player's NP/spent unchanged but still
  transfers, refunds, emits **Paid**, and logs actual LOG_CODE (10) through
  MASTER to the GM. Immediate mode preserves day-0 framing without output-bit
  mutation. The port does not independently assert or roll back this C path.
- Generic interruption calls real WAIT i_wait, removes duplicate wait-list
  entries and loads the successor. NULL callback, loaded command, absent
  command, stack leader and follower behavior are distinct controls.
- Day-0 STOP consumes only the top queued order, distinguishes running WAIT
  from loaded/idle commands, and preserves C's sometimes surprising output:
  loaded hide says **Interrupt current order** but stays loaded. Case/leading
  spaces/trailing arguments match, whereas attached comments and misspellings
  do not; the pinned fuzzy comparison is of the remainder, not the token.

The recording contains 165 sections / 5953 lines on Linux x86-64.
All 24 added sections exit 0. The earlier 141 sections / 4852 lines are
byte-for-byte unchanged; no callback reachability exception was needed.
Helper-compat retains its aborting transfer-dependency stubs; these full
transfer cases are in command-compat and never reach them.
Whole-game checks remain opt-in in `turn-compat`. The #62 cases consume no RNG.

## Mint and entity lifecycle closure (#63)

Added whole real u.c `new_char` (282–317), `create_unique_item`
(1948–1961), `destroy_unique_item` (1978–2007), `beast_capturable`
(2547–2560) and `sub_item` (639–658). The existing real
`create_unique_item_alloc`/`add_item` are included unchanged, with actual
buy.c `investigate_possible_trade` (764–806). **buy.c is a new input** in
tools/README.md and affected.sh. At #63, new items had no lore; #64 below
completes that branch with the real lore routines. No token, movement command, death conversion,
kill_char or combat routine is replaced by a successful stub.

`new_char` uses real `set_lord` and `set_where`. Its character-target branch
uses the complete small stack.c join/leave closure: `stacked_beneath`,
`first_prisoner_pos`, `promote`, `unstack`, `leave_stack`, `stack`,
`join_stack` (55–73, 116–254), with loc.c `all_char_here` (87–104), here-list
mutation and set_where (228–274, 299–311). Actual sout.c `vector_char_here` and VECT
fanout remain linked; Go now preserves their persistent recipient list and
output bits/logs. No report body is substituted for a database mutation.

The new driver operations construct arguments, call these exact routines
and observe `save_box`, kind/subkind/uniqueness, NP counters and raw logs:

- Character defaults (60/60 attack/defense, 50 breakpoint), named/unnamed
  entities, bracket replacement, signed-char health, ordinary ownership,
  touched locations and sorted faction units. Location and ship placement;
  actual stacking before prisoners with moving-day inheritance and logs.
- Capturable/noncapturable NI, undead, NPC/summon loyalty and zero breakpoint.
- Unique items held by characters or locations, consecutive allocations,
  retained zero-quantity inventory entries, and deleted boxes whose unique
  owner/subkind remain but whose serialization is omitted by real save_box.
- Dead-body NP refund/logging to old lord, independent-owner previous-lord
  fallback and absent-player control. No full death pipeline is claimed.
- Wrong-kind/nonunique/absent inventory assertions; sub_item missing,
  insufficient, zero and last quantity plus argument assertions.
- Nested join/leave with descendants, persistent here order, release flag
  cleanup, summon-to-NPC loyalty, recursive movement restoration, exact
  VECT logs, and descendant/captive/different-location assertion controls.

Two real C runs produce **186 sections / 6585 lines** and the Go test matches.
The previous **165 sections / 5953 lines** are byte-identical. Full allocator
branch/range, exhaustion and shared-stream discrimination live in
[`code-compat`](../code-compat); these lifecycle cases consume that same
actual C mint stream, not a driver RNG or assigned IDs.

Pending-feature review adds two cases after the original 184 sections /
6527 lines, preserving that entire recording prefix:

- Joining the default unlocated characters succeeds in both actual C and
  the unchanged Go: 1000 is placed under 1001, its here-list contains 1000,
  both ordinary stack messages are logged, and status is 0. No missing
  NULL guard was reproduced. `all_char_here` (loc.c:87–104) asserts
  `valid_box(who)` first; oly.h:477 embeds `x_loc_info` in every box and
  oly.h:947 defines `rp_loc_info(n)` as its address. The C NULL branch is
  unreachable for a valid box, including one whose `where` is zero. Go's
  existing translation remains unchanged apart from an explanatory comment.
- An unnamed noble can leave `dead_char_body`'s saved name NULL
  (u.c:625). Omitting the optional `body_new` name now constructs this state
  without changing any previously named case. Real C destruction logs
  `(null)~[x542] has passed on.  Gained 7 NPs.`, refunds 7 NP to player 201,
  deletes item 77542 and omits it from subsequent `save_box` output, status 0.
  Before correction the same Go case panicked, observed as status 139 with
  no flushed stdout. `destroyUniqueItem` now uses the existing `clib.Str`
  convention; the same recording matches without changing expected bytes.

## Item lore state and day-0 initialization (#64)

The complete real pinned functions now include lore.c `queue_lore`
(230–245), `scan_char_item_lore` (390–411), u.c `set_known`
(2347–2367), day.c `storm_owner_touch_loc` (1681–1702) and
`init_locs_touched` (2306–2322), storm.c `update_weather_view_loc_sup`
(1616–1625) and `init_weather_views` (1647–1658). The existing real
`add_item`/`test_known`, player/lord, inventory loops, sparse-bit helpers,
viewloc and province close the dependency chain. The `queue_lore` aborting
stub is removed. **lore.c is a new input** in tools/README.md and affected.sh.
The new driver operations only construct or observe state or call these
actual functions. Report bodies and lore-sheet delivery remain stage 2;
no knowledge, queue, touched-location or weather-view mutation is skipped.

18 new cases distinguish unknown from known, known-but-not-queued from
prequeued-but-unknown, ordinary duplicate suppression from `anyway`
duplicates, item IDs from their lore-sheet IDs, player and unit ownership,
zero-quantity acquisition, repeat acquisition and faction changes. Scanning
uses actual character/inventory order, rejects invalid and nonpositive
inventory entries, includes prisoners and excludes player inventories.
State dumps include unsaved queue/known insertion order, touched locations
and weather views; real `save_box` checks the persistent/unsaved boundary.

Day-0 cases cover nested city/building/stack locations, prisoner exclusion
from character touching but inclusion for weather mages, player and unit
storm owners, unowned/no-owner storm controls, repeated append-only init,
preservation of existing bits, unlocated/region locations, nonmages and
signed mage flags. A mage with no province records bit 0 exactly as C.
Ownership assertions retain their exact C diagnostics and status 134.
The unowned lore scan reproduces [G3 #353](https://github.com/playbymail/olympia-g3/issues/353):
`test_known` dereferences `bx[0]`, status 139 with no flushed stdout/stderr.
No NULL guard or independent fix was added. G3 #352's separate NP-transfer
behavior remains unchanged.

Before #109, two actual C runs produced **204 sections / 7002 lines**, byte-identical;
Go matches their raw stdout, stderr, status and observed/serialized state.
The prior **186 sections / 6585 lines** remain byte-identical. No new RNG
calls or streams are introduced. Whole-game turns remain opt-in and full
turn orchestration is still #65 and later work, not enabled by this task.

## Location and garrison output (#109)

The former aborting `players_who_rule_here` dependency is replaced by its
complete pinned definition (garr.c 168–201) and `province_admin` (73–98),
using original structs/macros and `#line` locations. The existing complete
`garrison_here`, `building_owner`, player/lord walks and sparse bits close
the dependency chain; original sout.c is still linked unchanged.

Three appended cases compare actual hidden-location logs for direct/outer
touching and silent factions; garrison versus ordinary silent-unit routing;
and garrison output before the real, still-aborting regular-combat boundary.
Before/after `save_box` records show persistent knowledge, pledge and castle
fields unchanged, while explicit state captures retain unsaved touched/output
mutations. Report assembly and MOVE/SAIL callback integration are not included.

Two fresh bounded C runs produce **207 sections / 7193 lines / 137909 bytes**,
and the Go recording test matches. All prior **204 sections / 7002 lines /
135225 bytes** remain byte-exact, with no rewritten prefix. Literal trailing
spaces in the new recordings are actual C output. The corresponding move
recording has 83 sections / 4067 lines / 81524 bytes and broader routing controls.

## Hades/Faery entry ambushes (#111)

Complete pinned definitions: npc.c `create_hades_bandit` (202–249),
`create_faery_bandit` (252–290), `hades_attack_check` (293–319), and
`faery_attack_check` (322–348); u.c `gen_item` (1885–1891), `can_pay`
(2135–2140), `charge` (2143–2147), `stack_has_item` (2150–2167),
`has_use_key` (2170–2190), `stack_has_use_key` (2193–2212),
`stack_sub_item` (2223–2286), and `autocharge` (2289–2294); order.c
`queue_order` (160–171), `queue` (199–208), and `save_player_orders`
(211–256). Original hades.c/faery.c prefixes supply actual tags, static
streams, turn guards, day children and keyed leaves in separate translation
units, with original `#line` locations. **npc.c, hades.c and faery.c are new
source dependencies** in tools/README.md and affected.sh.

The existing `w.streams` instance supplies both region leaves and sequential
mint draws. There is no new RNG, global test reset, assigned success ID or
hardcoded successful roll. `master` writes and loads a real 16-byte seed in
both engines without resetting streams. Region leaves are **keyed**, not
counter-advancing: repeated identical day/actor/location keys replay their
roll. Early-return tests reload a different master in the same turn after
calling the hook, then observe original-root region rolls; the following
turn reseeds. This distinguishes drawing before safety/NPC predicates from
an incorrect early predicate return. Subsequent real unique-item allocation
observes the shared mint cursor independently.

21 appended cases cover every Hades monster and both Faery kinds; actual
minted IDs, quantities, gold and creation order; seeded ambush failure and
success; retaliation/no retaliation and same-day repeated entry; hidden,
NPC loyalty/subkind, noncharacter, transcend-alone versus troops/stack,
safe-haven nonexclusion, friendly/foreign prisoner stones and nonpositive
inventories. Payment cases preserve actor-first then friendly stack order,
signed 32-bit totals, insufficient-funds nonmutation, signed use-key width,
inventory validity/positive-quantity filters and exact negative-quantity
assertions. These raw negative/overflow constructions are helper contracts,
not claims about publicly reachable engine defects. Generation keeps real
lore/trade bookkeeping and the unique-item assertion.

`bandits` serializes actual entities with C `save_box`, observes output bits
and ready/run queues. `pstate` observes exact remaining orders;
`orders_file` calls actual `save_player_orders` and reads its raw database
file. Ordinary spectators' touched-location state makes appearance messages
and output bits observable through actual sout routing. Real queued WAIT
is loaded by `init_load_sup` at priority 1, executes, completes, and loads
the queued ATTACK at priority 3. ATTACK **execution is not performed** by
these C contracts and remains explicit NYI in Go; no successful combat or
unfinished-handler replacement was added. Autonomous programs, daily Faery
hunting (#77), MOVE integration (#67) and whole-turn acceptance remain
separate tasks. #77 should reuse these entry hooks and helpers.

Two fresh bounded C runs produce **228 sections / 8574 lines / 156526 bytes**,
matched byte-for-byte by Go, including raw stdout/stderr/status, logs,
serialized entities/orders, output bits and queue/RNG observations. All prior
**207 sections / 7193 lines / 137909 bytes** and the complete prior command
input are exact prefixes. Movement's full **135 sections / 7150 lines /
142755 bytes**, including its complete input, is unchanged. Recording
whitespace remains raw C output. Default tests read only recordings; world
turns remain opt-in. Every reached unfinished C dependency marker still
fails before recording, even for expected-abort cases and `--update`.

## Stack arrival state (#112)

Added complete pinned definitions, with original headers/macros and `#line`
locations: loc.c `mark_loc_stack_known` (281–295), move.c
`touch_loc_after_move` (692–707) and `move_stack` (710–773), storm.c
`update_weather_view_locs` (1628–1644), stealth.c `clear_contacts`
(285–299), buy.c `market_here` (342–355), `seller_comp` (373–383),
`seller_list` (386–427), `buyer_list` (430–464), `scan_trades`
(635–653), and `match_trades` (656–695). These are **288 C lines** across
11 complete definitions, plus original move.c `ocean_chars` storage (8).
All source files were already tracked by tools/README.md and affected.sh;
there are no new upstream source dependencies or pin changes.

The real accepted `set_where`, recursive stack walks, knowledge/touch/bit
helpers, sout routing/logging, io serializer and #111 Hades/Faery hooks are
reused. Go uses the existing persistent `w.streams`, including the shared
mint cursor, with no modeled RNG, reset or success-specific behavior. The
city closure executes real matching, including actual offer scans and seller
sort-key mutation, for empty, filtered and unmatched paths. Candidate pairs
stop explicitly at **unported `attempt_trade` (#81)**, not at an invented
successful settlement. Its C boundary aborts with the rejected unfinished
marker; it cannot enter an accepted recording, including `--update`.

22 appended cases discriminate:

- Hidden leaders, captive roots, captive descendants and their noncaptive
  children; different factions; ship passengers and noncharacter-pruned
  branches. Knowledge/weather/touch exclude captive descendants, while
  root characters are included even when captive. Contacts clear the entire
  character stack, including captives, but not ship roots or unrelated units.
- Actual location touching versus the unused destination parameter, nested
  city/building views and outer touched bits; weather's province versus
  destination; absent-player allocation, signed mage flags, nonmages,
  duplicate calls, preserved bits and province zero. Contact allocation and
  retained allocated-empty lists are separately observable.
- Ocean entry/exit, repeated same-ocean arrival, nonzero preserved timers,
  signed-char extremes, unrelated/follower timers, duplicate ocean entries,
  inconsistent preexisting timer/list state, and a city inside an ocean.
  Here-list order and all saved versus unsaved state are observed separately.
- City/noncity/ship market walks, actual ship-subkind rather than kind,
  filtered and repeated BUY/SELL/PRODUCE entries, captive matching, unmatched
  items, self exclusion, offer traversal/cost/tie order, moving/captive seller
  exclusions and city offers last. Raw mismatched-owner fixtures assert at
  the original C line; these are helper contracts, not public-command bugs.
- Real successful Faery entry (master 0/day 17) and Hades entry (master
  13/day 7), and a deliberately shared region wrapper with **both** successful
  hooks (master 13/day 7): actual Faery unit 136063 precedes Hades unit 93970
  in the arrival here list, priority queue and appearance log. The real keyed
  C leaves selected these inputs, never a modeled RNG. Touching occurs before
  appearances; Faery's special GM log occurs before relocation, only on region
  entry, while repeated same-day arrival still runs the real ambush hook.
  Ocean/city/building/outside-region, failed-draw, hidden-lone and NPC controls
  do not ambush. Master reload observations preserve the existing stream root.

Temporary BUY and SELL candidate-pair probes exercised the real Go arrival
path to its NYI settlement dependency: location, knowledge, touch, weather
and contacts had changed; trade quantities and ocean-exit state had not.
They were removed rather than coupling permanent tests to a temporary NYI
name. A temporary actual-C candidate-pair case was rejected in both normal
and `--update` modes, with the expected recording unchanged, then removed.
No C settlement success or Go settlement equivalence is claimed.

Two fresh bounded deterministic actual C runs produce **250 sections /
10405 lines / 186888 bytes**, compared exactly with Go's raw stdout,
stderr, status, real serialized boxes, faction logs, output bits and unsaved
state. The complete prior **228-section input (57628 bytes)** and **8574
lines / 156526 output bytes** are exact prefixes. Movement's complete input
and **135 sections / 7150 lines / 142755 bytes** are unchanged. Literal C
serialization whitespace is retained. New arrival cases require status 0
except four explicit assertion contracts (134); a crashing fixture cannot
be silently accepted as an ordinary arrival.

Full MOVE callbacks remain #67, settlement #81 and
report assembly stage 2. Report-path state is included here, not silently
excluded. This is helper/state acceptance, not whole-turn acceptance.

## Movement prisoner escapes and nearest land (#113)

The five new production definitions are complete pinned stack.c
`check_prisoner_escape` (255–286), `prisoner_escapes` (289–328),
`prisoner_movement_escape_check` (331–342), u.c `find_nearest_land`
(2010–2092), and sout.c `vector_stack` (585–602): **185 C definition
lines**. The comparator extracts the first four with original `#line`
locations; sout.c is already linked whole. Existing real release/unstack,
moving restoration, touch, arrival, Hades/Faery, output and serialization
dependencies are reused. dir.c `location_direction` (182–194) is now also
extracted whole: **dir.c is a new read-source dependency**, recorded in
tools/README.md and affected.sh. No pin change.

The comparator retains the original u.c entity-RNG declarations/tags and
all eleven complete functions (64–181; 77 definition lines) from its prefix,
including the turn guard, per-day child, prisoner roll and direction/pick
leaves. The new closure totals 275 definition lines including those eleven
functions and `location_direction`; all other gameplay dependencies were
already real. Go reuses **persistent `w.streams`**, not a fresh RNG or a
per-call reseed. `entity_rolls` observes actual keyed leaves. Reloading a
different master within the turn tests guard persistence; changing day
changes only the prison leaf, changing turn reseeds. Even a non-ocean land
lookup seeds the entity stream before returning, just as C does.

Important source distinction: pinned `loop_char_here` (loop.h:60–69)
constructs a depth-first **snapshot** via `all_char_here` before executing
escape bodies. It is not a live shrinking outer traversal. Adjacent
successful escapes and a successful parent followed by its nested prisoner
all execute; replacing this with a live index walk would skip prisoners.
The real here lists still shrink/reorder during unstack/promote, and nested
restoration/arrival walks keep their existing C ordering. The swimming
parent case proves that its real arrival clears the nested prisoner's
contacts before the nested escape runs on shore; the adjacent sibling then
swims independently. No alleged live-list bug was invented or fixed.

23 appended cases discriminate:

- Actual C-chosen rolls **20/21** at threshold 20 and **10/11** with hounds
  at threshold 10; failed/no-hound, singular/plural barking, whole-stack
  faction ownership (foreign captor hounds do not count), nested vector
  recipients, release flags, summoned loyalty conversion, moving restoration
  and exact promotion/here-list order. Province early return retains contact
  state; nested temple-to-city arrival clears contacts and updates weather.
- Real completed galley, roundship and raft paths: port exit without a swim,
  unfinished-ship exit into ocean with real ocean bookkeeping, first-island
  choice in here-list order, and an actual Hades shore ambush producing unit
  136063 and its real queued WAIT. Swimming emits real faction and special
  logs, updates output bits, and executes real `move_stack` arrival effects.
- Independent asymmetric west/north/south/east land outcomes **9004 / 9001 /
  9003 / 9002** from the actual keyed direction leaves, not fixed/modelled
  landing choices. Missing-exit rotation restarts from the origin, while
  the next outer retry preserves the current ocean cursor. Plan B logging,
  full bounded Plan C fallback, region/depth/ocean filtering and entity-chain
  candidate order are compared. Two real fallback picks select **9020** then
  **9005**; an empty eligible set returns zero.
- Deliberately malformed direct-helper fixtures (an ocean with no exits,
  and an item masquerading as an island) retain exact assertions at u.c
  **2053** and **2034**, status 134 and flushed output. These are not claims
  of public-engine defects. Existing SIGSEGV139/no-flush contracts remain.

Two fresh bounded deterministic C runs produce **273 sections / 11966
lines / 215667 bytes**, exactly matched by Go's raw stdout/stderr/status,
real serialized boxes, faction logs, output bits, vector/contact allocation,
moving flags, ocean lists, known/touched/weather state and RNG observations.
The complete prior **250-section input (67815 bytes)** and **10405 lines /
186888 output bytes** are exact prefixes. Movement's entire input (53589
bytes) and **135 sections / 7150 lines / 142755 bytes** remain unchanged.
Literal C whitespace is retained. New cases require status 0 except the two
explicit malformed-helper assertions; unfinished C dependency markers are
rejected before status handling and before any `--update` recording.
An ephemeral expected-abort case followed real escape → city arrival →
candidate trade settlement: both normal and `--update` modes rejected its
actual C dependency marker without changing the recording. A temporary Go
probe confirmed the real production boundary after release, moving restore,
contact clear and arrival knowledge/touch, without changing trade quantities.
Both probes were removed; no permanent test depends on the temporary NYI name.

Default tests remain small recorded comparisons; full worlds stay opt-in.
This adds no successful gameplay stubs, replacement RNG, modeled landing or
independent C fix. MOVE integration (#67), daily/weekly prisoner events
(#80), settlement (#81), reporting and stage/whole-turn acceptance remain
outside this task.

## Complete MOVE/GO and synonym integration (#67)

New production definitions at the unchanged pin (261 C definition lines):

| source | complete functions and original lines |
|---|---|
| move.c | `discover_road` 140–162; `v_move` 584–689; `d_move` 776–832 |
| move.c | `v_exit` 1096–1105; `v_enter` 1108–1124; `v_north` 1127–1136; `v_south` 1139–1148; `v_east` 1151–1160; `v_west` 1163–1172 |
| order.c | `prepend_order` 174–181 (unlike append, no 250-order limit) |

The real prerequisite departure/arrival, stack, knowledge, payment, hooks,
escape and output helpers are reused without new gameplay substitutes.
Command-compat imports 1646 already-compared definition lines from move-compat's
travel/topology, admission and character-description closure, plus the 261
new lines above: 1907 additional extracted definition lines total. Original
combat `fight`/`wield` structures (combat.c:110 onward), display storage
(display.c:402), exitdir tags/static stream/turn guard (dir.c:22 onward),
headers, macros and complete function `#line` locations are retained.
New tracked source inputs are **art.c, display.c and lib/exit_views.c**, in
tools/README.md and affected.sh. No upstream source or build flags were repaired.

18 appended cases observe actual callbacks through the original command table:

- An NPC queues MOVE → WAIT; real start and four evenings change wait
  4 → 3 → 2 → 1 → loaded WAIT. WAIT starts the next day and finishes on day 6.
  Daily loops, queues, serialized command arrays, log framing, output bits,
  contacts, knowledge, touched locations and weather views all match. A second
  exit shuffle and real unique-item allocation observe the shared exdr/mint
  streams after movement, rather than resetting or modeling them.
- Zero-distance building movement detaches the old stack, suspends/clears
  guards then immediately restores the full moving stack, including prisoners,
  before loading the successor. Outer-level fallback reparses and reinserts
  the **original complete line** ahead of WAIT at each level, even after
  argument shifting. A building exit, city exit and subsequent province trip
  expose the distinct zero/one/four-day paths.
- Hades computes cost before leaving the old friendly stack and charges
  actor-first then friendly stack inventory. Captives count toward cost but
  foreign gold is untouched. Insufficient payment keeps location, old stack,
  guards and inventories unchanged. Special-log order is compared.
- A deleted destination fails at actual evening completion before restoring
  moving flags. Hidden roads reveal every reverse road for the complete
  descendant snapshot, including foreign prisoners; roads to other locations
  are not revealed. Without a reverse road the forward road remains unknown
  to the foreign follower. Raw known-list insertion order is observed.
- Every cardinal long/short synonym and ENTER/IN/EXIT/OUT reparses to MOVE;
  explicit ENTER uses only its first operand and empty ENTER assumes IN.
  Generic interruption demonstrates MOVE's actual NULL interrupt pointer.
  No arguments, already-there, unknown/alternative destinations, overload,
  water rejection and refused building entry preserve exact failure order.
- Actual fog arrival broadcasts to the stack but suppresses public arrival.
  Real Hades/Faery arrival hooks allocate bandits, queue their WAIT/ATTACK,
  touch/output/log and consume the shared mint stream in C order. These
  focused saved-state completion cases set the remaining wait to one; the
  separate NPC journey exercises every real elapsed evening without shortening.
  ATTACK is not executed or accepted as successful combat.
- Actual start-time prisoner escapes run before departure output, using C's
  **snapshot** `loop_char_here`, not live shrinking traversal. Adjacent escaped
  prisoners are all processed; remaining leader movement and here ordering
  are observed before serialization and departure messages.

Two fresh bounded deterministic C runs generate **291 sections / 14137 lines /
259462 bytes**, exactly matched by Go raw stdout/stderr/status and the actual
serialized/unsaved state, queue/order files, faction logs and output bits.
All 18 new sections must exit zero; a crashing fixture cannot be recorded as
successful movement. The entire old command input **4545 lines / 77598 bytes**
and old output **273 sections / 11966 lines / 215667 bytes** are exact prefixes:
**zero old bytes changed**, with no newly-reachable-boundary exception.
The entire movement input **3150 lines / 53589 bytes** and recording
**135 sections / 7150 lines / 142755 bytes** are unchanged, including #108's
actual running raft/roundship/galley crew-loss controls and base-distance cases.
All literal trailing whitespace remains real C serialization/output.

An ephemeral expected-abort probe executes actual queued MOVE → evening →
city arrival → candidate trade settlement. Both normal and `--update` runs
reject its unfinished C dependency marker (status 2), before expected-134
handling or recording. The accepted recording's SHA-256 remains
`0853cefeca98e716f700fd80260e66693ab97d33ef79a9b24e9ce247553728f2`.
The probe is removed; no private permanent test locks an NYI name. Rejection
still recognizes only the explicit harness marker, not ordinary game text.

Defaults remain fast recorded checks, whole-world runs opt-in. No new RNG
stream/reset, successful combat/trade substitute, full SAIL/FLY implementation
or full-turn success is claimed. Full upstream raft-sail fixture execution
remains an explicit #90 engine-integration obligation; full trade settlement
is #81 and the stage exit is still unfinished. G3 #309/#352/#353 remain the
known intentionally reproduced nonblocking upstream concerns.

## Hades population and NPC route generation (#68)

Complete new production definitions at the unchanged pin:

| source | definitions and original lines |
|---|---|
| hades.c | `create_hades_nasty` 404–470; `auto_hades_sup` 473–496; `auto_hades` 499–526 |
| npc.c | `get_exit_dir` 87–97; `choose_npc_direction` 100–122; `npc_move` 125–145 |

These six definitions total **174 C lines**. The real Hades prefix already
supplies `hads_nasty`, the tags, stream and turn guard. Its translation unit
now also contains the original complete spawn/auto functions. A separate
original npc.c prefix retains the per-location stream and `begin_npc`,
`npc_spawn`, `npc_qty`, `npc_behavior` (34–62; 26 definition lines). All
headers, macros and function `#line` locations are original. Only phase
instrumentation is a no-op; no gameplay dependency is replaced by success.
Both upstream files were already tracked inputs; no new dependency or pin.

Go reuses the persistent **`w.streams`**, existing keyed `HadsNasty`, shared
sequential mint/exdr, real `newChar`, `genItem`, queue/loading, topology,
serialization and output closure. The route selector shuffles **before**
testing direction persistence, including when that direction wins. NPC
behavior reseeds by location on each nonempty selection, unlike Hades' and
exdr's once-per-turn guards. NPC day-0 dispatcher wiring remains **#71**;
the existing unported dispatcher stop is unchanged, not prematurely skipped.

14 appended cases compare actual C state, not a model:

- Zero population creates all 25 slots and every monster kind with its real
  minted ID, name, health, attack/defense, spirit quantity, loyalty, breakpoint,
  player membership, location/here order, loaded WAIT and subsequent ATTACK
  or MOVE. With 24 existing units only slot **24** spawns (actual unit **90729**).
  Existing units outside Hades, followers, captives and negative-health units
  still count. Populations 25 and 26 neither allocate replacements nor delete
  excess units; all faction members still receive the real per-unit handling.
- Region children are indexed in saved here order, without hidden/terrain
  filtering: the asymmetric list includes underground, plain and city children.
  Lore, touched locations, player serialization, exact order files and queue
  insertion order remain observable. Repeated slot 3 replays its kind/quantity
  across day/master changes within one turn, while minted IDs advance; the next
  turn reseeds. A later real unique-item allocation observes the shared mint.
- Target order is direct here-list order, not recursive/ID order. Hidden,
  NPC-loyal regular-player characters and captives remain targets; nonregular
  owners and transcend-alone characters are filtered. Transcend with troops
  or descendants is not alone. No eligible direct target queues movement even
  when a nested ordinary-player character exists. Existing orders are retained.
- Building/city/ship paths queue literal `move out` without changing npc_dir.
  No exits/all-water exits preserve absent misc and signed-char direction state.
  Land selection excludes water but **not** hidden/barrier/impassable exits;
  original NSEW topology includes UP/DOWN but not roads or inner cities.
  Actual queued strings include all cardinal directions, `move up`, `move down`.
- Actual NPC roll **9** keeps its direction; roll **10** takes the shuffled first
  exit (actor 1003 north instead chooses east). Missing/zero/invalid previous
  directions use the shuffled list. Multiple actors, day/master/turn changes
  and later explicit shuffles discriminate persistent exdr order from keyed
  per-location behavior. Generated `npc_move` → four real MOVE evenings → WAIT
  exercises actual accepted command dispatch, arrival, log and output-bit state.

Two fresh bounded deterministic C runs produce **305 sections / 17256 lines /
299103 bytes**, exactly matched by Go raw stdout/stderr/status and real
serialized state, order files, logs, output bits, queues and RNG observations.
All 14 new sections require status 0. The old command input **5232 lines /
87733 bytes** and output **291 sections / 14137 lines / 259462 bytes** are exact
prefixes: **zero old bytes changed**, no newly reachable boundary exception.
The entire move input **3150 lines / 53589 bytes** and output **135 sections /
7150 lines / 142755 bytes** remain unchanged. Literal C whitespace is retained.

An ephemeral expected-134 probe runs real Hades targeting → queue/load →
ATTACK dispatch. Both normal and `--update` reject the reached unfinished
C dependency marker (status 2) **before** expected-abort handling/recording.
The recording SHA-256 stays
`4c0f5a13dcd2ef2d86b13a24a0b63b8d23001be2c8a4843baf2c890cf6a0cd56`.
The probe is removed; no retained test depends on a temporary private NYI name.
Marker rejection remains specific, not a broad ban on ordinary game text.
Combat, full SAIL/FLY, autonomous dispatcher integration and whole-turn/stage
acceptance are not claimed. Default checks use recordings; full worlds remain
opt-in. No independent bug fix, second RNG or successful missing-gameplay stub.

## Savage forts and autonomous savage queues (#69)

Sized before implementation: **245 C definition lines** at the unchanged pin,
with all remaining dependencies already in the accepted command closure:

| source | complete definitions and original lines |
|---|---|
| savage.c | `create_savage` 13–29; `call_savage` 32–65; `savage_hates` 219–241; `savage_hate_here` 244–263; `auto_savage` 266–343; `init_savage_attacks` 351–403 |
| npc.c | `controlled_humans_here` 65–84 |

**savage.c is a new source dependency**, recorded in `tools/README.md` and
`tools/affected.sh`. The C harness extracts these complete original definitions,
plus the original MAX_SAVAGES macro and static population storage, with original
headers, macros and `#line` locations. Only phase instrumentation remains a
no-op. The existing original npc.c prefix provides all NPC tags/stream/leaves;
Go uses the same persistent `w.streams` and shared mint/exdr, real character
creation, inventory/lore, recursive traversal, queue/loading, movement,
serialization and output. There is no new RNG or successful gameplay stub.

19 appended C-recorded cases cover:

- Master 0 / turn 1 / fort **9409** really rolls **1**, fort **9100** rolls
  **10**. Ship **9498** also rolls 1 when queried, but is not visited by
  `loop_loc`. Spawning considers every building-depth location, including a
  mine, not only the hated-structure predicate. Fort iteration is entity order,
  not saved here order. The population counts all savage noble-items regardless
  of subkind, faction, health, prisoner or stack state.
- Initial **199** allows two real spawns (110158 and 140449), retaining tally
  199 throughout; the next initialization recounts **201** and stops. Initial
  200 and 201 allocate nothing, preserve excess units and leave exit/mint
  streams untouched. No increment-on-spawn cap has been introduced.
- Successful keyed checks still reject nested forts; zero land exits retain
  exact stderr `init_savage_attacks: no exits?`. All-human routes spawn nothing
  after the shuffle. Asymmetric fallback rejects controlled sources 9003 and
  9001 before actually spawning in hidden/barrier/unsworn source **9004**;
  ocean exits are excluded. Source-keyed quantity and minted-actor-keyed WAIT
  are recorded independently, along with day/master/turn and sibling controls.
- Human detection recursively visits cities, ships and stacks, excludes self,
  unsworn loyalty and nonzero character subkind, but does not exclude hidden,
  prisoner, negative-health, NPC-loyal or NPC-faction ordinary characters.
- All twelve hated complete/unfinished building/ship subkinds and negative
  mine/ghost-ship/raft/city controls; follower no-op; in-structure RAZE; empty
  first target ENTER → RAZE; direct here order rather than ID or recursive
  target order. A nonhuman-owned first target is not replaced by the next
  human-owned one. Direct children queue UNSTACK in saved order (not just
  savage children), then DIE or movement as real human detection dictates.
- Human-owned targets add exactly one missing drum, preserve existing quantity,
  and append USE → keyed WAIT → ATTACK after saved orders. Actual actor 1000
  gets WAIT 48, actor 1002 gets 38; day replay preserves the draw, master and
  turn changes reseed it. Autonomous route MOVE executes four real evenings
  through WAIT. A genuinely spawned unit's numeric MOVE executes seven real
  evenings, arrives, and **loads** USE (priority 3) without executing it.
  Serialized command/DB, log, output-bit, queue and order-file state are real.

Two fresh bounded C runs are deterministic: **324 sections / 18602 lines /
320880 bytes**, exactly matched by Go raw stdout/stderr/status. All 19 new
sections require status 0. The previous entire command input **5739 lines /
97486 bytes** and output **305 sections / 17256 lines / 299103 bytes** remain
exact byte prefixes; **zero old bytes or reachable boundaries changed**.
The entire move input **3150 lines / 53589 bytes** and output **135 sections /
7150 lines / 142755 bytes** are unchanged. Raw C trailing whitespace and the
fixture's missing-plural-name diagnostics are retained, not normalized.

An ephemeral expected-abort probe queues real savage RAZE, loads it and
dispatches it. Both normal and `--update` reject the reached missing-dependency
marker with status 2 before expected-abort handling or recording. The recording
SHA-256 remains `e6d7515ec89b8b5219e441ed7c40867dadba8093109a77d2b2d1f9911b64f8a9`.
The probe was removed; no permanent test depends on a temporary private NYI name.
Generated RAZE/USE/ATTACK/UNSTACK/DIE are queued, **not** claimed as successful
execution; unported handlers remain explicit stops. The unrelated drum/summon/
keep-savage command handlers are not reached or ported. Autonomous dispatcher
wiring remains #71; full combat and whole-turn/stage acceptance remain pending.
Default checks are small recordings, with full worlds still opt-in. Parent
review reported the reproduced cap overrun as [G3 #354](https://github.com/playbymail/olympia-g3/issues/354);
the Go source links it and deliberately preserves C's behavior. It is
nonblocking; known G3 #309/#352/#353 are unchanged.

## Unsworn, mob and undead autonomous queues (#70)

Sized before implementation: **94 C definition lines**, with the remaining
closure already accepted. Complete original pinned definitions are extracted
with original headers/macros and `#line` locations:

| source | definitions and original lines |
|---|---|
| npc.c | `auto_unsworn` 149–165; `auto_mob` 168–199 |
| necro.c | `auto_undead` 363–387 |
| loc.c | `subloc_here` 332–351, used by original loc.h:22 `city_here` macro |

**necro.c is a new input**, tracked in tools/README.md and affected.sh.
The existing original npc.c prefix supplies its static stream, exact tags,
`begin_npc` and keyed leaves. Go reuses persistent `w.streams`, shared mint
and sequential exdr, real topology/movement, saved-order append, serialization,
logging and output bits. No new stream/reset or successful gameplay stub.

20 appended real-C cases cover:

- Seeded unsworn skip, first direct city in saved here order (9101 before
  9100), second-coin route fallback, absent-city/landless paths and nested-city
  exclusion. Master 0/turn 1/location 9000 gives actor 1000 first coin **2**,
  actor 1001 coins **1,1** and actor 1004 **1,2**. Building skip preserves a
  previously selected location stream even after master reload. City movement
  goes out; city lookup filters kind, not hidden/barrier status. Raw wrong-kind
  city and negative-health/subkind constructions are helper-contract controls,
  not claims about public-engine defects.
- Exact missing-misc mob warning with no allocation/reseed; age **4** no-op
  versus age **5** stay/disperse (actors 1000/1001, turn 10). Away-from-home
  disperses regardless of future creation age, follower/prisoner/ordinary
  loyalty. Signed-short turn promotion is tested around 32767/-32768 without
  narrowing the age difference. WAIT precedes literal `reclaim "disperses."`,
  appended after existing work. Day replay, master/location/turn changes and
  wide keyed probes distinguish persistent stream selection from reseeding.
- Undead without misc/summoner, summoner elsewhere, same province but distinct
  sublocations, and same-sublocation stacked/captive summoner cases. ATTACK
  appends before `summoned_by` clears; an immediate second call runs normal
  behavior, proving clearing happened before return. The attack branch leaves
  the previously selected NPC stream untouched. Other branches seed by actual
  sublocation, including building/city/ship short circuits; province coins
  distinguish MOVE from `pillage 1`. No extra loyalty, health or prisoner filter
  is introduced. Inventories and loyalty serialize unchanged.
- Actually generated undead cardinal MOVE runs all four real evenings, arrives
  and executes its following WAIT; actually generated unsworn numeric city MOVE
  runs the real one-day callback and following WAIT. Ready/run queues, command
  arrays, location/here order, knowledge/touch state, real faction logs, output
  bits and exact order files are compared. Subsequent shuffle and unique-item
  allocation observe shared exdr/mint state. Asymmetric multi-exit/day/master/
  turn cases retain existing guarded exdr behavior, not a test reset.

Repeated direct `orders_file` calls can retain an earlier file when the queue
becomes empty: real `save_player_orders` opens the file lazily. This is a
harness-only call sequence, not evidence of an engine save defect: production
`save_orders` (order.c:355–367, called by save_db) removes/recreates the orders
directory before calling the per-player helper. The raw helper output is kept.

Two fresh bounded deterministic C runs produce **344 sections / 19994 lines /
344574 bytes**, matched exactly by Go's raw stdout/stderr/status and observed
state. All 20 new sections require status 0. The complete previous command
input **6272 lines / 106631 bytes** and output **324 sections / 18602 lines /
320880 bytes** are exact prefixes: zero old bytes changed. The entire move
input **3150 lines / 53589 bytes** and recording **135 sections / 7150 lines /
142755 bytes** are unchanged. Literal C whitespace remains untouched.

An ephemeral expected-134 probe executes real undead queue → load → PILLAGE
dispatch. Normal and `--update` runs both reject its missing-C-dependency marker
with status 2, before abort acceptance/recording. The recording hash stays
`541ad99deb0b9f29a2ef7a6be4bd5235147da9986a536ca85efc5c19b1a1ea52`.
An ephemeral Go probe follows generated PILLAGE, ATTACK, and mob WAIT → RECLAIM
to their actual explicit NYI stops. Both probes are removed; no permanent test
depends on private temporary NYI names. Marker rejection remains specific to
unfinished C dependencies, not ordinary game text.

NPC dispatcher wiring remains #71. ATTACK/PILLAGE/RECLAIM are queued, not
accepted as successful execution. Default checks use small recordings; whole
worlds remain opt-in and no full-turn/stage success is claimed. No new upstream
defect was identified; known intentionally reproduced G3 #309/#352/#353/#354
remain unchanged and nonblocking.

## Bandit programs and day-0 NPC dispatch (#71)

Sized before implementation: **132 C definition lines**, all gameplay
dependencies already accepted. Complete pinned definitions are npc.c
`auto_bandit` **351–403** and `queue_npc_orders` **613–691**. The comparison
additionally extracts complete u.c `stage` **2563–2594** (32 lines), reusing
Go's existing real stage output. Original `#line` locations, macros, headers,
NPC stream/tags and keyed leaves remain unchanged. `auto_unsworn` and
`auto_mob` move into the original NPC-prefix translation unit because they
are static C functions called by the real dispatcher. Only phase instrumentation
is omitted; there are no new upstream source inputs or pin changes.

Production uses the existing `processOrders` callback, persistent `w.streams`,
`NPCVictim(who, len(targets))`, shared mint/exdr, creation, Hades/savage/
unsworn/undead/mob/movement, queues, output and serialization. There is no
replacement RNG, independent behavior fix or successful unported handler.
The sublocation-monster program's empty branch is an intentional C no-op;
unknown programs retain **npc.c:687 FALSE**, status 134 with flushed output.

13 appended real-C cases distinguish:

- All direct child UNSTACK orders (including noncharacters) before ATTACK,
  saved target order rather than ID/sorted/recursive order, and keyed victim
  indices rather than target shuffling. Actors **1000** and **1011** select
  different victims (**1002** and **1003**) from the same three targets.
  Moving the first target to the end changes 1000's victim to **1003**.
  Day replay, master reload, next turn and subsequent real mint/exdr calls
  observe the exact shared/keyed stream contracts.
- Noncharacters, NPC subkind/loyalty and really-hidden-alone exclusions;
  hidden with descendants, prisoners, negative health and safe havens are
  not independently excluded. Faery friendly-descendant stones versus foreign
  captives, zero stones, Hades learned transcend-alone versus troop/descendant
  controls, and ordinary-bandit controls are compared through actual queues.
- Landless/nested-target and building fallback paths, preserved prior stream
  when no victim and no reseeding route occurs, asymmetric land/water routes,
  persisted directions and shared exit shuffles. Generated bandit MOVE runs
  four real evenings, arrives and executes its following WAIT; logs, output
  bits, knowledge, commands and order-file state are observed, not modeled.
- Summon, prisoner, LOAD/RUN/ERROR and queued-order skips before invalid-program
  assertions; DONE and absent-command eligibility; every ordinary/undead/NI
  savage/peasant/angry-peasant, all bandit, monster-no-op and token branch.
  Other NI items and unmatched subkinds intentionally do nothing. Independent
  saved unit order differs from ID order and is preserved by the existing
  `loopUnits` snapshot, taken after savage initialization and Hades processing.
- An actual newly minted savage precedes the 25 real Hades allocations and
  keeps its loaded MOVE plus USE/WAIT/ATTACK queue without a spurious autonomous
  rerun. Serialized state, ready queues and minted identities expose the
  initializer order. Actual dispatcher stderr is `queue_npc_orders()`.

Two fresh bounded deterministic C runs produce **357 sections / 22259 lines /
375381 bytes**, exactly matched by Go raw stdout/stderr/status and state.
The 13 new sections require status 0 except the explicit invalid-program
assertion (134). The complete prior input **6808 lines / 115505 bytes** and
output **344 sections / 19994 lines / 344574 bytes** are exact byte prefixes:
zero old bytes changed, no reached-boundary expectation exception. The full
move input **3150 lines / 53589 bytes** and output **135 sections / 7150 lines /
142755 bytes** remain unchanged. Literal C whitespace is retained.

An ephemeral expected-134 probe executes generated bandit ATTACK through the
actual queue/load/dispatch path. Both normal and `--update` runs reject its
unfinished-C-dependency marker with status 2 before abort acceptance/recording.
The recording hash stays
`c769c52ccec35e1c42d33d50be06150343b7af9f7af8fd95050277671fed7534`.
The probe was removed; marker rejection is specific, not a ban on game text.
No permanent test names a temporary private NYI boundary.

The opt-in real whole-world diagnostic uses the accepted issue104-roundtrip
cache: two C turns exit 0 with identical complete captures; Go exits 2 at
**ping_garrisons**, after the real NPC dispatcher. It still fails state/raw
output/status comparison. No early full-turn enablement, skipped garrison/
token/daily/month work or stage-exit acceptance is claimed. Defaults remain
small recordings. No new upstream bug was found; known nonblocking
G3 #309/#352/#353/#354 remain unchanged.

## Garrison guard setup and daily default payroll (#72)

The narrowed closure adds four complete definitions: garr.c
`garrison_strength` (54–71) and `ping_garrisons` (1064–1093), day.c
`default_garrison_pay` (581–613) and `maint_cost` (1171–1199), approximately
110 source lines. The existing real inventory, description, ownership, stack
and output closure is reused. `gm.c`'s original `gold_garrison` definition is
extracted with its source line; `affected.sh` and the dependency table track it.
A thin wrapper exposes only the static payroll entry point. Actual `oly.h`
macros and `loop.h` loops remain unchanged in the C reference.

The 13 added sections discriminate:

- Every maintenance tariff and unsupported/signed item IDs; man versus
  non-man, zero upkeep, independently insufficient attack/defense, equality,
  missing optional soldier attributes, negative stats and wide quantities.
  `loop_inv` excludes absent/deleted items and zero/negative inventory entries.
- Strength 9/10, repeated threshold crossing, initially guarded empty/weak
  garrisons, unrelated ordinary guards, multiple garrisons behind another
  first character, and exact original missing-character assertion.
- Real `liner_desc` and province output; castle ownership through a nested
  unit lord and pledged ruler; touched versus untouched viewers; missing,
  invalid and unowned castles. Guard is cleared before description/output.
  `show_to_garrison` stays true through routing, then becomes false even if
  it was true on entry. A subsequent location emit exposes that reset.
- Topmost (not immediate) stack leaders, a default middle ancestor beneath
  a nondefault root, self-default skips, captive skips, NPC subtype and both
  NPC/summon loyalty filters. Ownership and stack ancestry differ deliberately.
  Initial balances 13/2/91/3/17, pay 7 and two payroll calls expose generation
  versus transfer: the leader's 3 gold stays unchanged; non-NPC accounting
  moves 19 → 47 → 75, including the unowned non-NPC recipient while omitting
  three paid NPCs. Orders, raw logs, output bits and serialized boxes match.
- Zero pay still runs generation/lore and emits `Received 0~gold`; a negative
  configured amount retains the actual `qty >= 0` assertion, not successful
  negative generation. This and a missing character structure are deliberately
  malformed harness states, not claims of public-engine defects.

Two fresh bounded C runs produce **370 sections / 23109 lines / 392577 bytes**,
matched exactly by Go. The entire previous input **7325 lines / 124357 bytes**
and output **357 sections / 22259 lines / 375381 bytes** are unchanged exact
prefixes. Move input **3150 lines / 53589 bytes** and output **135 sections /
7150 lines / 142755 bytes** are wholly unchanged. Raw C trailing whitespace
is intentionally retained.

An ephemeral expected-134 garrison section queued and actually dispatched
ATTACK. Both normal and `--update` rejected `out-of-scope C dependency:
v_attack` with status 2 before accepting the abort or recording. The recording
hash stayed `fc9869bc2aaa8d0c1ef2bf844f30bafb801cd87016ed5ce896c6a2a60a431125`;
the probe was removed. No permanent test binds a temporary NYI phase name.

Command (normal and update), move, helper, code, loadsave, turn-driver/default
turn comparisons, bounded vet and fresh bounded default tests pass. C tools
were run sequentially before vet/all-tests; defaults remain small recordings.
The optional full diagnostic using `C_CACHE=tmp/issue104-roundtrip` repeats
actual C successfully/deterministically, then honestly fails comparison with
Go status 2 at **check_token_units**, after real guard setup. The Go faction
207 guard log is **424 lines / 33034 bytes**, byte-identical to C's corresponding
guard lines and an exact prefix of its saved log. The same 15 DB paths and
missing `orders/204`, plus raw streams/status, remain red. No full-turn or
stage-exit success is claimed; calendar/token and later world work stay
explicit, and month-end `garrison_gold` remains #116 after maintenance #87.
No new upstream bug/blocker was found; G3 #309/#352/#353/#354 remain unchanged.

## Inventory transfer and token acquisition (#117)

Ten complete production definitions, **231 C lines**, at the unchanged pin:

| source | definitions and original lines |
|---|---|
| u.c | `consume_item` 1894–1900; `move_item` 1908–1935; `drop_item` 2102–2132 |
| art.c | `token_player` 1016–1029; `swear_token_units` 1032–1051; `add_token_unit_sup` 1083–1120; `add_token_units` 1123–1141; `move_token` 1153–1194 |
| code.c | `display_kind` 318–343; `box_name_kind` 461–466 |

The comparator extracts complete original u.c/art.c definitions with original
`#line` locations, headers and macros. code.c is already linked whole, so its
two newly ported naming helpers require no extraction. No new source input
or pin change: all three files are already tracked in `affected.sh` and the
dependency table. The compiled `move_token` branch does **not** melt units;
the melting calls are in `#if 0`. Go preserves that distinction and reuses
real `setLord`, interruption/NP/membership, character/unique creation, shared
persistent mint/entity streams, inventory, lore, pending trades and output.
There were no existing Go `consumeItem`, `moveItem` or `dropItem` definitions
to wrap; their source of truth is now u.go.

22 appended actual-C cases discriminate:

- Transfer's nonpositive early return precedes even invalid box/quantity
  validation, unlike direct consumption/discard. Insufficient stock leaves
  destination untouched. Ordinary/self/unique moves, zero entries and lore
  plus pending BUY/SELL queues retain exact inventory order and quantities.
  Deliberately malformed direct calls retain unique-consumption, negative
  quantity and excessive-unique assertions (134), not successful transfers.
- Unique drops choose the province, first island in saved here order, or
  ocean fallback after real bounded Plan B/C. Subsequent entity leaves and
  unique allocation observe real stream state. A nowhere body is destroyed
  even by direct zero-quantity drop and refunds its old/previous player;
  a located body moves to the province without deletion or NP refund. These
  helper boundaries do not assert a public DROP command defect.
- Regular, independent, other NPC, player-as-holder and location holders;
  same versus different faction; empty-independent reset, absent optional
  magic, signed/zero capacity and stale noncharacter list rows. Stale rows
  are skipped by swearing but still count toward creation; they are explicit
  malformed saved-state probes, not evidence of a public-engine defect.
- A running real WAIT is interrupted and queued old orders are flushed.
  The transferred oath noble refunds/charges **8 NP** (20/30 → 28/22),
  retains loyalty, loses pledge/attitudes and changes membership. Real missing
  units are minted as **110158**, **140449**, with restricted `r`, token
  backlinks, health, beast/nonbeast breakpoints, touched locations and output.
  Land creation is in the province rather than the owner's temple; ocean
  creation is aboard the owner's ship. Independent transfer defers melting;
  reacquisition and later unique item **64258** distinguish shared mint order.
- City, coastal port city, mountain city, all three storm kinds and default
  names use actual `display_kind` and `box_name_kind` strings. Raw serialized
  boxes, command/ready queues, order files, faction/GM logs and **all player
  output bits** are compared, not synthetic success values.

Two fresh bounded deterministic C runs produce **392 sections / 24501 lines /
415416 bytes**, exactly matched by Go's raw stdout/stderr/status and state.
The previous entire command input **7677 lines / 130850 bytes** and output
**370 sections / 23109 lines / 392577 bytes** remain exact prefixes:
**zero old bytes changed**, no newly reached old-C boundary exception.
Move input **3150 lines / 53589 bytes** and output **135 sections / 7150 lines /
142755 bytes** remain wholly unchanged. Literal C trailing whitespace stays.

An ephemeral expected-134 inventory probe queued and dispatched ATTACK.
Both normal and `--update` rejected the actual unfinished dependency marker
with status 2 **before** abort handling or recording. SHA-256 stayed
`9dedf1d276843c1f6f7946dd614ddc5f8dbdbab402ee44dbbff48d90c595cfc3`.
The probe was removed; no permanent test binds a private temporary NYI name.

The local creation count advances after a failed `add_token_unit_sup`, exactly
as C does. That return path cannot be driven through the real allocator at
this pin: `new_ent` calls `alloc_box(n,...)` even after all searches return -1,
so exhaustion asserts before `new_char` can return a negative ID. No allocator
replacement or fabricated failed-allocation success was added to reach it.

Defaults remain small recordings and full worlds opt-in. `checkTokenUnits`
is still explicitly unported; inheritance/release/desertion (#118), full
death/reclamation (#119), then melting/reconciliation (#73) remain separate.
No unrelated GIVE/DROP/USE handlers or token-artifact generation were added.
No whole-turn or stage-exit acceptance is claimed. Known nonblocking
G3 #309/#352/#353/#354 remain unchanged; no new public-engine defect was
reproduced by this task's bounded cases.

## Inheritance, prisoner release and desertion (#118)

Nine complete definitions, **345 C lines**, at the unchanged pin:

| source | complete definitions and original lines |
|---|---|
| u.c | `stackmate_inheritor` 399–419; `take_unit_items` 422–538; `put_back_cookie` 582–593 |
| stack.c | `drop_stack` 377–442; `extract_stacked_unit` 466–520; `move_prisoner` 635–649 |
| swear.c | `unit_deserts` 303–343 |
| code.c | `plural_item_box` 431–442; `box_name_qty` 453–458 |

The script extracts whole actual definitions with original headers, macros
and `#line` locations; code.c remains linked whole. Original gm.c
`gold_combat` and `gold_combat_indep` definitions are extracted alongside
`gold_garrison`. All source dependencies were already tracked by
tools/README.md and affected.sh; there is no pin or source-list change.
The only new wrapper observes original static entity RNG leaves. Production
uses the persistent `w.streams`, existing keyed EntTake/EntDrop, snapshot
inventory/here loops, real token acquisition, lord/NP/interruption/dependent
handling, movement/prisoner escape, logging and serialization. There is no
replacement RNG, successful arrival stub or independent gameplay fix.

23 appended real-C sections distinguish:

- None, explicit, MATES and MATES_SILENT recipients; first eligible direct
  child **1003** beats parent **1001**, lower-ID nested child **1002** and
  captive **1005**. No eligible child falls back to the parent; no parent
  means discard/release. Silent modes suppress the heading and item/prisoner
  descriptions without suppressing actual inventory and prisoner changes.
- TAKE_ALL, TAKE_SOME and TAKE_NI (including the added noble item). At master
  0/turn 1/from 1000, gold **101** selects **86**, peasants **17** select **7**;
  both original copied remainders are actually removed. Item **12**'s gate
  is 2, so its independently probed zero amount must NOT be used; item **15**
  really selects zero through gate 1. Unique token **60001** also selects
  zero, then its copied remainder moves to the province and its units swear
  to the independent faction, with no bogus recipient heading. Day replay,
  master reload within a guarded turn, next turn and return to the old turn
  expose exact persistent stream behavior, keys, subkeys and draw ranges.
- Foreign regular gold **137** increments ordinary combat accounting; same
  faction gold does not. Independent transfers/discards accumulate
  **93 → 110 → 196**, the last increment counting the **86 selected**, not
  all 101 discarded. Ordinary discarded gold adds **31**. Live inventories
  end empty, not with an accidentally retained discarded remainder.
- Inherited token **60000** moves a running WAIT noble, flushes old queued
  work, clears pledge, reassigns its dependent to the old direct lord and
  exchanges **8 NP** (20/30 → 28/22). Real logs retain the inheritance's
  three-space indentation through token/lord processing. Shared mint creates
  the subsequent unique item **65748**, with no private stream/reset.
- Prisoner-only first-heading output, same/foreign-faction transfer, saved
  here-list order, moving restoration and release-swear normalization. The
  same-faction prisoner is cleared only after its captive description; the
  foreign captive stays captive. With no recipient, release-swear is cleared
  before drop so neither fealty nor disappearance happens.
- Direct release's actual keyed success (prisoner **1002**, gate **4**) runs
  real lord/NP/order handling; gate **5** prisoners **1001**, **1007**, **1015**
  exercise all three disappearance messages. Natural drop → unit_deserts →
  extract_stacked_unit runs nested/moving mixed children, leaves noncharacters
  behind, returns the NPC cookie and changes the released unit to deadchar.
  A prisoner-only descendant stack actually escapes, touches and moves using
  the existing arrival closure. Ordinary drop's vector and summon restoration
  are compared too, not replaced with a successful return.
- Extraction chooses the first eligible saved-order child, reparents other
  characters (including captives), promotes it immediately after the extracted
  unit rather than appending after a sibling, and preserves noncharacter and
  nested-child state. Prisoner-only building extraction executes real escape
  and arrival. Captive desertion into the captor's faction clears the captive
  bit without extraction; foreign desertion stays captive; noncaptives extract.
  Exact new_lord, show_day, membership, dependent and cookie state is saved.
- Singular/plural/zero/negative/large quantities and raw missing-plural
  diagnostics. Two malformed helper-contract controls retain exact original
  assertions (drop with wrong parent; move_prisoner on a nonprisoner), status
  134, not successful operations or public command-defect claims.

Two fresh bounded deterministic C runs produce **415 sections / 26512 lines /
448122 bytes**, exactly matched by Go's raw stdout/stderr/status and observed
DB, log, all-player output bits, accounting, orders/queues and RNG state.
All 23 new sections require status 0 except the two explicit assertions.
The complete previous command input **8086 lines / 138905 bytes** and output
**392 sections / 24501 lines / 415416 bytes** are exact prefixes: zero old
bytes changed. New input is **8731 lines / 150387 bytes**. Entire move input
**3150 lines / 53589 bytes** and output **135 sections / 7150 lines / 142755
bytes** are unchanged. Raw C trailing whitespace is intentionally preserved.

An ephemeral expected-134 section queued and dispatched ATTACK. Normal and
`--update` both rejected `out-of-scope C dependency: v_attack` with status 2
before accepting the abort or writing the recording. SHA-256 stayed
`9b4e7ded457fee992ff85b8268f7e78bd21dfeea9006995dbd8d0b13911d5c2c`.
The probe was removed; no permanent test binds a temporary private NYI name.
Rejection remains specific to source dependency markers, not ordinary text.

Command (normal/update), move, helper, rng, code, loadsave, turn-driver/default
turn comparisons pass; C tools run sequentially before vet/fresh bounded
default tests. Full worlds remain opt-in. No full turn or stage exit is
claimed; attempt_trade remains the explicit #81 boundary. Shared death and
reclamation remain #119; reconciliation remains #73. No UNSTACK/DIE/RECLAIM
handler was added. No new public-engine bug/blocker was reproduced; known
nonblocking G3 #309/#352/#353/#354 remain deliberately unchanged.

## Complete shared character death and reclamation (#119)

Eight complete production definitions, **296 C lines**, at the unchanged pin
`6a25f2d252ce45bd5adbe3baaccb3be3eed56279`:

| source | complete definitions and original lines |
|---|---|
| u.c | `survive_fatal` 220–249; `char_reclaim` 252–262; `dead_char_body` 596–635; `kill_char` 734–833; `hack_unique_item` 1938–1944 |
| use.c | `forget_skill` 1038–1067 |
| hades.c | `random_hades_loc` 529–555 |
| report.c | `char_rep_location` 357–407 |

The C tool extracts complete actual definitions with original headers/macros
and `#line` locations. It additionally extracts use.c `skill_school`
1174–1189 (16 lines, already implemented in Go) for the original `magic_skill` macro.
**report.c is a new source dependency**, tracked in tools/README.md and
affected.sh: this location formatter feeds the death log, not deferred report
assembly. Original Hades/entity prefixes and persistent `w.streams` preserve
all tags, keys, draw order, turn guards and shared mint/exdr state. There is
no second RNG, reset, successful gameplay stub or replacement C definition.

28 appended actual-C sections discriminate:

- Ordinary, captive and stacked deaths, with the dying unit's captive bit
  cleared for its notification and parent notification, then set to suppress
  subsequent output. Child **1003** precedes **1002** as MATES inheritor;
  gold **101 → 86** and accounting, captive transfer/release, child promotion,
  dependent lord changes, attitudes/pledge removal, actual running WAIT
  interruption, saved orders and ready queues follow the existing real closure.
  MATES_SILENT suppresses headings without suppressing transfers.
- Survival consumes the known skill, restores health **100** and clears
  sickness, leaving stack/items/orders/lord intact. Visible and captive-output
  variants and a second fatal attempt distinguish survival from real death.
  Absent/learning skills do not survive. Forgetting removes learning as well
  as known entries, but returns false for missing entry/character structure.
  Magic-school ancestry, ordinary skills, equipment aura bonuses, excess
  **100 → 99** rather than a wholesale clamp, zero and negative aura are real.
- Reclamation sets melt and bypasses both survival and transcendence, preserves
  their skill entries, zeros magician aura, restores the NPC cookie, discards
  inventory rather than giving it to mates, and flushes/interrupts real orders.
- Ordinary corpses retain their entity ID and original nullable name, old lord,
  peasant weight, unique-item inventory ownership and full death timestamp.
  NPC subkind/loyalty, melt and no-grave cases become retained T_deadchar IDs
  with zero subtype instead. Ocean cases select the first actual island in
  saved order or execute bounded Plan B/C and leave no corpse when no land
  exists. Later real entity leaves and unique creation observe shared RNG.
- Multiple eligible Hades locations are chosen in entity-chain order, not
  region-here order. Actor **1000 → 9101**, **1001 → 9100**; actual death
  moves 1000 to 9101, updates knowledge/touch/output and restores health/sickness
  without relinquishing lord/skills. Day replay, master reload during a
  guarded turn, next turn and return to the old turn expose exact keyed state;
  the subsequent real unique mint is **76084**. No eligible candidate returns
  zero without a pick; trying to transcend in that deliberately incomplete
  world reaches the existing actual arrival `loc_depth(0)` assertion (134),
  not a fabricated successful nowhere arrival.
- Token unit membership removal, positive capacity decrement, exhausted unique
  destruction, positive melt capacity preservation and zero-capacity melt
  destruction. The actual `player(who_has) == sub_pl_regular` expression is
  preserved. The reproduced missing regular-holder vanishing notification is
  filed/assigned upstream as [G3 #355](https://github.com/playbymail/olympia-g3/issues/355)
  and linked in production, not silently described as harmless behavior.
- Exact land, nested building/stack, named/unnamed region, ocean and nowhere
  location strings. Missing-region, noncharacter kill and surviving noncharacter
  child assertions retain original file/line/condition and flushed status 134.
  These and the absent-Hades world are explicitly malformed harness controls,
  not reproduced public-engine defects. Initial fixture-construction mistakes
  (unallocated optional-structure target, raw subtype assignment bypassing
  chains, invalid unique holder and missing ocean exits) were corrected in
  the harness, never worked around in either engine.

Two fresh bounded deterministic C runs produce **443 sections / 27958 lines /
471158 bytes**, exactly matched by Go raw stdout/stderr/status and serialized
death fields, names, kind chains, logs, all-player output bits, orders/queues,
NP/accounting and RNG observations. New sections require 0 except four explicit
malformed-state assertions (134). The complete previous command input
**8731 lines / 150387 bytes** and output **415 sections / 26512 lines / 448122
bytes** remain exact byte prefixes: **zero old bytes or reached boundaries
changed**. New input is **9325 lines / 161068 bytes**. Entire move input
**3150 lines / 53589 bytes** and output **135 sections / 7150 lines / 142755
bytes** are unchanged. Raw C trailing whitespace remains intentional.

An ephemeral expected-134 section queued, loaded and dispatched actual ATTACK.
Both normal and `--update` rejected `out-of-scope C dependency: v_attack` with
status 2 before accepting the abort or writing the recording. SHA-256 stayed
`67b12101d5b61ed8e095cfa3e189d833e6243263f617baa4a3ebebe64c7a8798`.
The probe was removed; no permanent test names a temporary private NYI boundary,
and marker rejection remains specific rather than banning ordinary game text.

Defaults remain small recordings and whole worlds opt-in. This does not add
RECLAIM/FORGET/DIE handlers, damage/ocean-stack death/resurrection/combat
callers, token melting/reconciliation or stage-2 report assembly. Existing
`checkTokenUnits` remains explicitly unported until #73; `attempt_trade`
remains #81. No full-turn or stage-exit success is claimed. Known nonblocking
G3 #309/#352/#353/#354 remain unchanged; G3 #355 is also nonblocking and is
deliberately reproduced rather than independently fixed.

Command (update and normal), move, helper, rng, code, loadsave, turn-driver
and default turn comparisons pass. C tools ran sequentially before bounded
`go vet ./...` and fresh bounded default `go test -count=1` with the shared
`-exec` timeout wrapper; both pass. Full-world and stage-exit tiers were not run.

### Pre-acceptance location buffer refinement

One additional real-C section exercises literal long nested location and
region names, first below individual-join boundaries and then lengthening
the inner name to **1169 bytes**. The outer and region names are **962** and
**889 bytes**. Twelve actual location calls cover land versus ocean wording,
nested joins that fill the buffer before province/region append, and joins
that first reach the limit at the province or region. C returns **2048-byte**
location strings at ten boundaries, with **1926/1916-byte** below-limit controls.

Before changing production, `--update` ran actual pinned C twice
deterministically and recorded this section; the targeted Go comparison
failed at line **27977**: the location was **2900 bytes**, versus C's **2048**.
All four concatenations in `charRepLocation` now call existing `sout` using
the exact C formats, including `%s, in province %s`; no wrapper or final-only
truncation was added. The same targeted comparison passes after the change.

The comparison is now **444 sections / 28020 lines / 500621 bytes**; input is
**9377 lines / 166322 bytes**. The entire previous **443-section / 27958-line /
471158-byte** output and **9325-line / 161068-byte** input are exact prefixes.
Only actual C generated the new expected bytes. Version remains
**0.49.0-alpha**, and G3 #355's intentional reproduction is unchanged.

After this refinement, command `--update` and normal comparisons pass
sequentially, followed by bounded vet and fresh bounded default all-tests.

## Token melting and reconciliation (#73)

The complete new production definitions are **art.c**, not use.c:
`melt_token_units` **1054–1080** (27 lines) and `check_token_units`
**1197–1259** (63 lines), **90 C lines** at the unchanged pin. Original
headers, macros and `#line` locations are retained. The only omitted effects
are phase counters; no death, ownership, creation or output dependency is
replaced with successful behavior. No new upstream source dependency or pin.

Go reuses the existing `loopUnits` snapshot, complete `charReclaim`/`killChar`,
`addTokenUnits`, `tokenPlayer`, lord/NP/interruption, inventory/prisoner/cookie,
lore, output/logging, real serialization and persistent `w.streams` closure.
The existing processOrders and postMonth callback positions are unchanged;
only the phase NYI definition is removed. G3 #355's player-ID/subtype mix-up
remains deliberately reproduced and source-linked in killChar.

Seventeen appended real-C sections cover:

- No tokens; regular empty, deficient, full and excess lists; repeated checks;
  zero/negative capacities; independent, other NPC, player-as-holder and
  location owners. An absent owner retains **art.c:1215 `owner`**, status 134.
- Allocated deleted/dead/noncharacter rows count toward regular capacity and
  retain the exact LOG_CODE diagnostic, including the apparently mistaken
  **duplicated `box_code(unit)`**, not an independently corrected token ID.
  Such malformed helper-state controls are not public-engine defect claims.
- Adjacent **1000 and 1001** both melt while real death removes them from the
  live list, proving copied traversal. Full running WAIT interruption and
  queued-order flushing, child extraction/promotion, captive release, discarded
  inventory, two returned cookies, magician aura zeroing, preserved survival/
  transcend skills, retained deadchar IDs/death times, all player output bits
  and raw logs/order files are compared. The positive token ends at capacity
  **1**, with no live units; a second check does not recreate independent units.
- Capacity 1 is retained on melting; exhausted zero/negative capacities
  destroy the actual unique token. Destroying two adjacent subtype entries
  does not skip the second. A genuine non-melting death decrements capacity
  **3 → 2** before regular reconciliation actually mints the two replacements.
- Actual NP requirement **8**: balance **7** preserves the foreign running
  unit; balance **8** transfers it through real set_lord, interruption,
  pledge/attitude clearing, dependent re-lording and **8 NP** refund/charge.
  A subsequent requirement **1**, balance **2** transfer proves the above
  threshold too. Unowned units omit the old-faction notification. Repeated
  calls cannot charge twice after ownership matches.
- Multiple tokens are visited in actual subtype-chain order, not allocation
  order; exact minted IDs, later real unique-item allocation, turn/day/master
  replay observe the existing shared mint/turn guards without a private RNG.
- The complete original **input.c `process_orders` 1057–1102** (46 lines)
  is additionally extracted with already-ported u.c `olytime_increment`
  **1365–1371**, `olytime_turn_change` **1379–1385**, `print_dot` **2371–2385**
  and original dot counter. A malformed-token case runs the actual day-zero
  NPC initialization and guard setup before its own owner assertion, before
  interrupted/player orders. Later unported turn callbacks are explicit
  missing-dependency failures, never accepted as a full successful turn.
  The separate turn-driver phase trace verifies both callback positions.

Two fresh bounded deterministic C runs produce **461 sections / 29818 lines /
527434 bytes**, exactly matched by Go raw stdout/stderr/status, serialized
state, death/kind chains, token counts/lists, logs, output bits, accounting,
orders/queues and shared mint observations. New sections require status 0
except the two explicit missing-owner assertions (134). The entire previous
command input **9377 lines / 166322 bytes** and output **444 sections /
28020 lines / 500621 bytes** remain exact prefixes: **zero old bytes changed**,
no newly reachable old-boundary exception. New input is **9861 lines /
175204 bytes**. Entire move input **3150 lines / 53589 bytes** and output
**135 sections / 7150 lines / 142755 bytes** remain unchanged. Literal C
trailing whitespace is retained; only real C generated expected bytes.

An ephemeral expected-134 section queued, loaded and dispatched actual ATTACK.
Normal and `--update` both rejected `out-of-scope C dependency: v_attack`
with status 2 before accepting the abort or writing the recording. SHA-256
remained `ab403e3f469edd5c5c74b5425346b9acd0e781418012001d0cbd628eb6d49753`.
The probe is removed; no permanent test binds a private temporary NYI name.
Marker rejection remains specific, not a broad ban on ordinary game text.

Command (normal/update), move, helper, rng, code, loadsave, turn-driver and
default turn comparisons pass sequentially before bounded vet/default tests.
The optional real-world diagnostic repeats C successfully/deterministically,
then honestly fails with Go status 2 at **match_all_trades** (day 1), the
same 15 DB differences and missing `orders/204`, plus raw streams/status.
No full-turn/stage-exit success, unrelated RECLAIM/USE/combat handler or token
artifact generation is claimed. Version **0.50.0-alpha**; defaults remain
recorded/fast and worlds opt-in. No new verified public-engine bug or blocker
was found; nonblocking G3 #309/#352/#353/#354/#355 remain unchanged.

## Natural weather creation and binding (#75)

This extends the existing real-source command harness instead of introducing
weather-compat: creation/binding needs its already accepted allocator, shared
mint stream, set_where, aura/command machinery, garrison/location output and
save_box. No parallel lifecycle, RNG adapter or successful gameplay stub is
introduced. Pin remains **6a25f2d252ce45bd5adbe3baaccb3be3eed56279**.

Complete original storm.c definitions are extracted with their original
locations: RNG prefix **40–104**, BIND **107–187**, new_storm **279–323**,
create_some_storms **1525–1554** and natural_weather **1558–1613**. Existing
weather_here and view helpers remain unchanged. Only phase counters are
omitted. u.c's complete nprovinces **2511–2526**, greater_region **2609–2624**,
deduct_aura **1773–1784**, charge_aura **1788–1799** and check_aura
**1803–1814** close the bounded helper dependencies. Original greater-region
global definitions additionally come from **cloud.c**, **tunnel.c** and
**quest.c**; these dependencies are listed in tools/README.md and affected.sh.

Go uses the existing world-owned WthrShuffle/WthrStorm with the exact
province, day*256+kind key, original candidate iteration order and sequential
shuffle counter. Strength leaves do not consume the shuffle counter. The
nonzero nprovinces cache retains C static lifetime on the world. The scratch
candidate slice itself has no observable persistence: C clears it on every
call. BIND's start/finish callbacks are installed in the real command table.

**21 appended sections** cover all eight seasons, month/year rollover,
0/1/31/32/63/64 provinces, reversed map construction with ocean/forest/plain/
mountain terrain, all six special-region exclusions, zero/positive/negative
and nested current storms, candidate exhaustion and repeated generation on
the same and later days. Below-threshold active seasons still shuffle;
Blossom/Harvest do not. Rain/fog/rain order is retained in Waning Days.
Real storm saves, province here-lists, signed-short strength wrapping,
explicit-ID relocation, first-weather message suppression, output flags,
ruler/pledge/touched versus untouched logs and resetting show_to_garrison
are compared. A mint collision skips occupied **110158 → 110159**; real
unique-item allocation interleaved with storms proves the shared allocator
order. Pre/post shuffle probes observe persistent wthr state and keyed
strength independently of mint.

BIND is queued/loaded/dispatched/completed through real callbacks, then
rebound between two ships; old/new ship lists, the storm backlink and aura
**12 → 9 → 6 → 3** are serialized. Ownership, ship, province and aura failures
are checked at both start and finish, including aura lost after validation.
Invalid kind/depth/month retain exact source assertions, stderr and status
134. All other new sections require status 0.

Two actual bounded C executions yield **482 sections / 32263 lines / 558979
bytes**, byte-identical to Go. Previous **461 sections / 29818 lines / 527434
bytes** of output and **9861 lines / 175204 bytes** of input remain exact
unchanged prefixes; new input totals **10197 lines / 181494 bytes**. Raw C
save trailing whitespace is deliberately retained. The entire previous move
recording and all calendar/schedule recordings remain unchanged.

A removed ephemeral expected-134 case dispatched the actual unported ATTACK
handler. Both normal and --update runs rejected its missing-dependency
marker with status 2 **before** status acceptance or recording updates.
The expected-output SHA-256 remained
`4ba88adc2281065fa898669e94c43b5bb876ab020c42b6a5ca2672499598ef1b`.

Command normal/update, move (135/7150), helper (37/1487), rng (321 results),
code (43/969), loadsave (61/3411), turn-driver (7/152), ordering-only phase
trace (5/484), calendar (2/119), schedule (4/1801) and default turn comparisons
pass sequentially. The schedule/phase traces are **not gameplay acceptance**.
No full-world or four-weather-day integration success is claimed: it remains
opt-in and the known full turn boundary is day-1 match_all_trades (#81).
Movement/decay (#76) and unrelated weather magic remain explicitly unported;
the existing dailyEvents call position is unchanged. No newly verified bug
or upstream blocker was found; known nonblocking upstream issues are retained.
Bounded go vet and fresh bounded default all-tests also pass. Full-world
tests were not run for this task.
Version **0.52.0-alpha**; no commit or push by this task.

## Storm movement, owner visibility and month-end decay (#76)

This continues command-compat rather than adding a duplicate weather-compat.
The complete pinned closure is **221 C definition lines**:

| source | complete definitions and original lines |
|---|---|
| day.c | storm_decay 1655–1678; storm_owner_touch_loc 1681–1702; storm_move 1705–1734 |
| storm.c | move_storm 191–250; move_bound_storms 253–276; dissipate_storm 365–425 |

The owner-touch helper and its real init_locs_touched call were already
implemented; they are retained and exercised through the shared initializer,
not bypassed with a weather-only replacement. Existing postMonth calls now
reach real decay and movement. Original C definitions, headers, macros and
assertion locations are extracted; only phase instrumentation is omitted.
All source inputs are already tracked, so there is no dependency-table,
affected.sh or pin change. Production reuses shared placement, kind iteration,
knowledge/touch state, inventory/cookie return, reclamation, output and saves.

**20 appended real-C sections** distinguish:

- Strength 0/1/2/negative and signed-short -32768/32767 decrement behavior;
  adjacent deletion during kind-chain traversal; coexisting same-kind storms;
  positive-strength direct dissipation checks weather *before* removal, whereas
  decay zeros strength before calling it. No independently corrected notice.
- Show false/true, all three weather messages, character/unowned/invalid/player
  owners, real ruler/pledge routing and show_to_garrison reset versus preservation.
  Dissipation does not enable that flag; monthly relocation does not alter it.
- Unequal deferred destinations, nonzero positive/negative npc_dir gates,
  an invalid deferred destination behind a zero gate, exact gate/destination
  reset and repeated calls. Monthly movement uses set_where directly, emits
  no weather notification, sets knowledge but does not touch a location.
  Player 201 gains known 9001 with an empty touched list.
- Real BIND followed by ship-bound movement, province selection from a city
  destination, adjacent stale/nonstorm bound-list rows, repeated same-province
  motion and saved storm/ship backlinks. Ship-bound motion touches player 201's
  province 9002 without adding knowledge. Cookie/home pairs return two items;
  incomplete pairs do not return an item. Missing optional ship state is safe.
- Shared init_locs_touched includes the captive storm owner's locations despite
  skipping that character directly. A nested temple/city helper-state control
  yields touched city 9100 and province 9001, without adding knowledge; ordinary
  player 201 touches only its character's province 9000.
- Paired RNG controls give identical two wthr shuffles and strength probes,
  and the same subsequent shared-mint unique item **67592**, with and without
  real movement/decay/touch/dissipation. Natural minted storms also run decay
  and deferred relocation through the actual helpers. Lifecycle has no draws.
- Four malformed helper-state controls retain exact C kind/depth/subtype
  assertions and status 134. Raw unallocated optional-structure targets were
  corrected to allocated deleted boxes before recording; these constructions
  are not claims about public-engine defects.

Two bounded C executions are deterministic: **502 sections / 33420 lines /
576837 bytes**, exactly matched by Go raw stdout/stderr/status and real saves.
The previous **482 sections / 32263 lines / 558979 bytes** of output and
**10197 lines / 181494 bytes** of input remain exact unchanged prefixes.
New input totals **10585 lines / 189472 bytes**. Entire move input and output
are unchanged; raw C trailing whitespace is retained.

A removed ephemeral expected-134 weather section dispatched actual ATTACK.
Both normal and --update rejected its missing-dependency marker with status 2
before expected-abort handling or recording, leaving the recording unchanged.
No permanent test names a temporary private NYI boundary.

Command normal/update, move, helper, rng, code, loadsave, turn-driver and
default turn comparisons run sequentially before bounded vet and fresh default
tests. Whole worlds remain opt-in; no full-world or stage-exit success is
claimed. SAIL and unrelated weather magic handlers remain explicit boundaries.
Report assembly is not accepted; the real touched/known/output DB mutations
remain in scope. No new verified bug/blocker was found; known nonblocking
G3 #309/#352/#353/#354/#355 remain unchanged. Version **0.53.0-alpha**;
changes remain uncommitted for parent review.

## Faery hunts (#77)

The same real-source harness now extracts **faery.c:529–634** in full:
`create_elven_hunt`, `warn_human`, `auto_faery_sup`, and `auto_faery`.
Original faer tag/key/turn-guard definitions were already included for #111;
its accepted entry hooks and creations are reused, not duplicated. Only phase
instrumentation is omitted. No new C source dependency or pin change is needed.
The production daily call position is unchanged; its former NYI is replaced.

Thirteen new small recorded cases exercise:

- Zero population fills slots 0–14 in creation order; 14 existing faction units
  (including a captive, an outside-region unit and an artificial noncharacter
  with ownership state) produce exactly one hunt at slot 14. Populations 15 and
  16 neither allocate nor remove anything. C counts the whole faction list.
- Asymmetric region child order and mixed ocean/land maps, retry-keyed rejection,
  actual elf quantities, repeated slot draws, day changes, master reload under
  the same turn guard, and reseeding on a new turn. Interleaved real bandit,
  unique-item mint and exdr shuffles preserve keyed hunt placement/quantity
  while sharing actual allocation and saved list order.
- First encounters queue actual MESSAGE plus body lines and log to LOG_SPECIAL;
  later encounters queue ATTACK. Memory is on the owning faction, shared across
  hunts, persisted in real save-box bytes. Direct saved here order changes both
  warning and attack order; noncharacters/nonregular factions and nested humans
  are skipped, while hidden, captive, unsworn and safe-haven targets are not.
- Same-faction nested stone holders protect their stack leader; other-faction
  holders do not. Removing protection warns once, then attacks. Protected-only
  paths queue real MOVE OUT; no land exits produce no invented MOVE or WAIT.
- An empty-current existing hunt executes real MOVE → WAIT. An actually spawned
  hunt **110158**, with **25 elves**, first executes its initialized WAIT, then
  an eight-day generated MOVE, arrives at 9001, and executes the next WAIT.
  Its real output includes the 25-men extra travel day, departure/arrival,
  known/touched mutations, empty saved orders, and completed command state.
- A valid faction-owned but unlocated hunt reaches the original `loop_here`
  call-site assertion, not a null-owner crash: exact stderr
  `assertion failure: olympia/faery.c (581): valid_box(where)` and status 134.
  Go asserts immediately after `pMisc`, before entering `loopHere`. Only this
  named Faery section accepts 134; missing-dependency rejection still precedes
  status acceptance and --update recording.

Two fresh C executions are deterministic: **515 sections / 36108 lines /
615332 bytes**, exactly matched by Go stdout/stderr/status, save-box state and
order files. The previous **502 sections / 33420 lines / 576837 bytes** of
output and **10585 lines / 189472 bytes** of input are exact unchanged prefixes.
The original #77 **514 sections / 36094 lines / 615031 bytes** also remain an
exact unchanged prefix after the assertion review. Input now totals **11065
lines / 197681 bytes**; entire move recordings remain unchanged. No expectations
came from Go.

An allocated C box cannot lack `loc_info`: pinned `oly.h:477` embeds
`x_loc_info` in `struct box`, and `rp_loc_info` at line 947 takes its address.
Thus no valid allocated-region case can reach `create_elven_hunt`'s `p`
assertion at faery.c:538. The source assertion remains ported; no fabricated
pointer state or unallocated-box crash is used as its expected behavior.

The expected-abort missing-dependency negative control is checked in both normal
and --update modes before status acceptance or recording. Temporary probes are
removed; no tests bless unported MESSAGE/ATTACK execution. Discarded artificial
empty/unallocated-region constructions were not public-engine bug reproductions
and prompted no engine workaround. Ocean-only maps are not routine tests.

Command, move, helper, rng, code, loadsave, turn-driver and default turn checks
run sequentially before bounded vet/fresh default tests. Whole-world checks
remain opt-in; no full-world turn or stage-exit success is implied. Reports and
unrelated command handlers remain explicit boundaries. No new verified public
bug or blocker; known G3 #309/#352/#353/#354/#355 are unchanged. Version
**0.54.0-alpha**, left uncommitted for parent acceptance.

## Structural damage, sinking and collapse (#120)

This bounded prerequisite was split from #78 before implementation. Five
complete **u.c definitions / 191 C lines** are extracted at the unchanged pin:

| definition | original definition lines, including return type and closing brace |
|---|---|
| kill_stack_ocean | 184–217 (34) |
| sink_ship | 1013–1068 (56) |
| get_rid_of_collapsed_mine | 1071–1094 (24) |
| building_collapses | 1097–1140 (44) |
| add_structure_damage | 1143–1175 (33) |

The issue's 1069/1176 endpoints include following blank lines; its mine
1072 endpoint is the signature, preceded by `void` at 1071. No source input
or pin changed, so the dependency table and affected.sh remain unchanged.
Original headers/macros and assertion locations remain intact; thin wrappers
only expose the two static helpers. Production reuses full killChar,
surviveFatal, inheritance/extraction, dead bodies, lord/NP/interruption,
arrival/nearest-land routing, output/VECT and bound-storm bookkeeping. It
adds no replacement RNG or successful unported gameplay boundary.

The **21 appended sections** distinguish:

- Zero, 98/99/100 and over-100 damage, exact-cap versus over-cap branches,
  canDestroy false's reset to 99, zero damage at 99/100, loaded byte 255,
  and negative damage's original u.c:1148 assertion. The unsupported direct
  removal of an uncollapsed mine retains u.c:1077, status 134. These two
  malformed helper-contract controls are not public-engine bug claims.
- Mine collapse keeps its location, damage 100, collapsed subtype and delay
  8; final removal again evacuates occupants added after collapse. Nested
  captive/moving characters and direct noncharacters preserve here
  order. Ordinary buildings are deleted. Castle collapse clears both
  matching garrison backlinks, leaving the other castle and ordinary-unit
  backlinks alone. Raw field observations additionally expose references that
  save_box would omit after the referenced castle is deleted.
- Real VECT output reaches nested foreign occupants and touched province
  viewers in order. Land/city sinking moves whole stacks, both adjacent roots
  and direct noncharacters; the no-optional-ship-record path also executes.
  Bound lists with two adjacent stale/nonstorm rows followed by two real
  storms are cleared in full; both real storm backlinks are zeroed, not the
  nonstorm field. Raw deleted-ship observations prove the list is empty.
- Ocean sinking kills descendants before ancestors and visits direct roots
  in saved order. Actual bodies, discarded inventory, NP, running WAIT
  interruption and queued-order removal are observed. Real skill survival
  consumes the skill, restores health/sickness, extracts and washes ashore;
  a surviving running WAIT remains running. A surviving child washes ashore
  before its parent dies. The three-level snapshot visits **1002, 1003, 1001,
  1000**, not sorted order or a live mutating traversal.
- Unequal ocean routes yield **Mountain 9003**; swapping east/west produces
  **Plain 9001** with the same guarded entity leaves. An island case chooses
  **9101** before 9100 in saved here order. No-land deaths run bounded Plan B/C
  without inventing a survivor. Day/master reload during a guarded turn,
  next-turn reseeding, later real entity probes and shared unique mint
  **65748** / **60040** remain exact. Paired zero-damage/no-call controls
  preserve the same entity probes and unique allocation.

Direct noncharacter ocean contents are intentionally ignored by sink_ship's
ocean branch, unlike land evacuation. This raw here-list construction tests
the helper branch, not a claim that public commands can put ordinary items
there. Early fixture mistakes (storm/garrison kind constants, an unowned
nested character, a missing peasant item structure and missing region wrappers)
were corrected in the fixture, never worked around in production or filed as
public-engine defects.

Two fresh bounded deterministic C runs produce **536 sections / 37892 lines /
644189 bytes**, exactly matched by Go raw stdout/stderr/status, real save_box
bytes, all-player output bits, logs, inventory, orders/ready queues, ownership,
stack/location order and additional raw retained fields. Only the two explicit
new assertions accept 134; all other new sections require status 0. The whole
previous input **11065 lines / 197681 bytes** and output **515 sections /
36108 lines / 615332 bytes** are exact prefixes: **zero old bytes changed**.
New input totals **11650 lines / 207996 bytes**. Entire move input **3150 lines /
53589 bytes** and output **135 sections / 7150 lines / 142755 bytes** remain
unchanged. Literal C trailing whitespace is retained.

An ephemeral expected-134 section queued/loaded/dispatched actual ATTACK.
Normal and --update both rejected `out-of-scope C dependency: v_attack` with
status 2 before abort acceptance or recording. The expected-output SHA-256
stayed `d05242fe101fd50f66f5e39f2bfa5b27a41e7e458bb306dfb275594e68278ac4`.
The probe is removed; no permanent test names a temporary private NYI boundary.

Command normal/update, move, helper, rng, code, loadsave, turn-driver and default
turn comparisons pass sequentially before bounded vet and fresh default tests.
Defaults remain small recordings; whole worlds remain opt-in. Daily coastal,
mine and inn events remain #78; character damage/healing/ocean polling remain
#79; monthly mine decay remains #88. Unrelated commands/combat and report
assembly stay explicit. No full-world turn or stage-exit success is claimed;
the known day-1 full-turn boundary remains match_all_trades (#81). No newly
verified public-engine bug or blocker was found; known nonblocking
G3 #309/#352/#353/#354/#355 remain unchanged. Version **0.55.0-alpha**; changes
are left uncommitted for parent review, not accepted/closed by this task.

## Daily coastal and location damage (#78)

The bounded daily closure reuses accepted #120 rather than introducing a
successful destruction substitute or duplicate weather tool. Six complete
definitions, **255 C lines**, are extracted at the unchanged pin
`6a25f2d252ce45bd5adbe3baaccb3be3eed56279`:

| source | definition | original lines (including return type and closing brace) |
|---|---|---|
| day.c | near_rocky_coast | 237–258 (22) |
| day.c | ship_coastal_damage | 261–327 (67) |
| day.c | mine_calamity | 330–399 (70) |
| day.c | inn_calamity | 403–451 (49) |
| day.c | random_loc_damage | 454–486 (33) |
| u.c | count_loc_char_item | 1336–1349 (14) |

The original `mine_depth` macro is oly.h:1102 (`shaft_depth / 3`, with a
signed-short shaft). Some issue endpoints included adjacent blank/body lines;
the table above identifies the actual complete definitions. Original headers,
macros and assertion locations are preserved; only phase counters are omitted.
No new upstream source dependency or pin change is needed. Go calls the
world's `streams.WthrWreck` directly: tag `wthr`, keyed `wreck` leaves,
`(entity, day*16+slot)`, without advancing the shared weather shuffle counter.
Ship slots are 0/1 for reef, 2/3 for rocks and 4 for omens; mine uses 0/1/2,
inn 0/1/2/3; their random-event gates use slots 8 and 9. The real turn guard,
shared master and mint stream remain untouched. Production dailyEvents order
is unchanged; only the two former NYI definitions were removed.

**17 appended sections** use small actual-C-selected cases, not a retained
random search or Go-derived expected output:

- Land/ocean/coast/rocky/open topology, rocky precedence in either exit order,
  all three unfinished ship subtype skips with otherwise actual reef-hit
  keys, and completed galley/roundship/raft handling. Real reef hits for IDs
  **9203/9124/9422** yield **3/4/5**; rock hits **9133/9135/9344/9518** yield
  **6/7/8/9**. Independent C probes verify gate 1, adjacent false gate 2 and
  maximum false gates 75/50. Damage **99** survives, exact **100** really sinks;
  adjacent deletion preserves kind-chain traversal. A real fatal-survival
  skill is consumed and its occupant washes ashore at 9001. Logs retain actual
  reef/rocks text, sinking, survival and output bits.
- Open-sea keys **9473/9230/9143** give **1/2/3**, reaching both birds arms
  and sharks; ID 9100 gives **92**, with no omen. Ship messages route through
  real touched-location output, not an invented per-occupant route.
- All four mine messages, zero workers, clamping to one worker, independent
  one/five victim draws, damage endpoints **1/15**, and uneven nested worker
  allocation. Saved depth-first order is **1001,1003,1002,1004,1000**, not ID
  order; losses **1,2,2** preserve later inventories. Hidden/captive/safe-haven
  characters are not independently excluded; captive output is suppressed by
  real output routing but its workers still die. The branch under a direct
  noncharacter is pruned. Count **27 → 22** excludes its 100 workers.
- Mine gate 9 with depths **8/9/10** proves `<`, equality and `>`; signed
  shaft **-4 / 3 → -1**, absent optional state and zero depth do not trigger.
  Independent C endpoint gates **1/90** pair shafts **2/3** and **269/270**.
  Two adjacent actual collapses retain subtype 53, damage 100, delay 8 and
  evacuate occupants in exact order; later random passes skip both mines.
- All **6 subjects / 5 actions / 7 consequences** of inn wording and damage
  **5/15** are emitted by actual C. Rare gate **9128 → 1** really damages
  an ownerless inn without output, then **88 + 12 → 100** collapses it with
  an owner and real messages. False gates **2/100**, unfinished inn/mine
  subtypes and ordinary nontrigger controls do nothing. Owner selection
  skips a direct noncharacter but does not skip a captive first character:
  damage still occurs, without incorrectly notifying the next character.
- Per-day replay for ships, mines and inns distinguishes subslot aliasing.
  Already-seeded weather plus master reload during the same guarded turn
  retains ship damage **3 → 6**, day 2 does not hit, and replaying day 1
  yields **9**. Next turn and return to the old turn reseed correctly. Paired
  no-call controls preserve all four shuffle observations and real shared
  unique allocations **77018 / 60845**.
- Two malformed helper-contract controls retain real loc.c:92
  `valid_box(who)` and loc.c:395 `loc_depth(where) == LOC_build`, with flushed
  status 134. All other new sections require status 0. These assertions, a
  corrected initially missing worker item, and a corrected duplicate raw
  inventory construction are not reproduced public-engine defects. The
  existing `loopHere` intentionally has no internal assertion; callers keep
  their original macro call-site assertions. No iteration helper was changed.

Two fresh bounded deterministic C runs produce **553 sections / 39215 lines /
669382 bytes**, exactly matched by Go's raw stdout/stderr/status and real
serialized state, inventories, ownership, topology/here order, output bits,
logs and RNG/mint observations. The complete previous input **11650 lines /
207996 bytes** and output **536 sections / 37892 lines / 644189 bytes** remain
exact byte prefixes: **zero old bytes changed**. New input is **12248 lines /
218785 bytes**. Entire move input **3150 lines / 53589 bytes** and output
**135 sections / 7150 lines / 142755 bytes** remain identical. Literal C
trailing whitespace is deliberately retained, including save_box here lists;
`git diff --check` consequently reports those recording bytes, not source
whitespace errors.

A removed ephemeral expected-134 section queued/loaded/dispatched actual
ATTACK. Both normal and `--update` rejected
`out-of-scope C dependency: v_attack` with status 2 **before** expected-status
acceptance or recording. Output SHA-256 stayed
`19bc0cfdeb87f5ff1c2c09294cb830caf88c1e84ef029f0a12574a8b248de7ba`.
No permanent case binds a temporary private missing-dependency name.

Command (update and normal), move, helper, rng, code, loadsave, turn-driver and
default turn checks pass sequentially before bounded vet and fresh default
all-tests. Every C/Go engine process uses the shared timeout supervisor.
World turns remain opt-in and were not run; the known day-1 full-world
boundary remains match_all_trades (#81). No full-turn or stage-exit acceptance
is claimed. Character damage/healing/ocean polling (#79), monthly mine decay
(#88), reports and unrelated commands/combat remain explicit boundaries.
No newly verified public-engine bug or upstream blocker was found; both
open-issue arrival checks found only enhancements, no local bug/follow-c.
Known nonblocking G3 #309/#352/#353/#354/#355 remain deliberately unchanged.
Version **0.56.0-alpha**, left uncommitted for parent review, not accepted,
closed, staged, committed or pushed by this task.

## Ocean characters and weekly healing (#79)

The actual issue names ocean polling and healing, not `add_char_damage`.
This extends the existing lifecycle/stream harness rather than introducing
an upkeep-only tool or widening the task to unrelated damage/commands.
Complete original definitions at unchanged pin
`6a25f2d252ce45bd5adbe3baaccb3be3eed56279` are extracted with original headers,
macros and `#line` locations (only phase instrumentation is omitted):

| source | definition | exact complete lines | count |
|---|---|---|---|
| move.c | init_ocean_chars | 835–849 | 15 |
| move.c | check_ocean_chars | 852–898 | 47 |
| day.c | heal_char_sup | 519–578 | 60 |
| day.c | heal_characters | 616–632 | 17 |

These endpoints are checked against the pinned source, not older boundary
notes: day.c:581 starts default_garrison_pay, not more healing. The four
definitions total **139 C lines**, of which **124** implement the three new
Go definitions. The existing Go initOceanChars body and loadDb call remain
unchanged. The C harness also extracts the original day.c:50–107 upkeep
prefix, including static state, purpose tags, begin_upkeep, up_today and
up_heal. No new upstream file is read, so dependency metadata and pin stay
unchanged. Go reuses the world's existing `streams.UpHeal` directly:
`upkp` turn guard, per-day `heal` child, `(who, 0)` recovery in 1–100 and
independent `(who, 1)` amount in 3–15. No adapter, new RNG or manual reseed.

**11 small appended sections** compare actual gameplay, not callbacks:

- Real init_ocean_chars populates the list in character-chain order, then
  check_ocean_chars runs repeatedly. Two adjacent leaders go **14 → 15 → 16**,
  really drown, become serialized bodies, then have timers/list entries
  cleaned on the following check. Invalid/no-longer-ocean cleanup allocates
  char state before the branch as C does and removes every duplicate value.
- Nested followers start at **16**, increment to **17**, but the parent gate
  stops independent plunges. A duplicate leader reaches 16 later in that
  same snapshot: real reverse stack deaths remove its children, while the
  leader's real known survive-fatal skill **611** is consumed. It washes
  ashore at **9001** through actual extraction/movement/arrival output.
  Its live-list removal does not skip the snapshot's remaining entries;
  dead follower cleanup happens on the next poll. Plunging uses real VECT,
  deaths, logs, output bits, ownership, body saves and touched/known state.
- Separate nonocean duplicates prove removal cannot skip the next snapshot
  entry. Signed-char timers **127 → -128**, and a follower **126 → 127**,
  preserve C widths before later land relocation clears both.
- Actual recovery keys **1103/1020/1161/1040** yield **5/6/10/11** on day 7.
  On plain only 5 recovers; inside an inn 5, 6 and 10 recover while 11
  remains sick. Recovery is followed by a gain in the same invocation;
  continued illness loses the independent amount. Amount endpoint keys
  **1023/1001** yield **3/15**. An exact **85 + 15 → 100** and a clamped
  **99 → 100** emit fully-healed messages. Sick loss clamps to remaining
  health, including zero, before the real MATES death/survival path.
- Outer traversal handles **0, -1, 99, 100, 101**: zero healthy health gains,
  zero sick health dies, negative and 100+ health skip. Two adjacent sick
  deaths transfer **19 → 42** items through real MATES inheritance and leave
  the next wounded character reachable. Actual loaded WAIT orders, queue
  cleanup, serialized commands/order files, real logs and output are observed.
  Another fatal illness really consumes survival skill 611 and restores 100.
- Initially seeded upkeep survives master reload in the guarded turn:
  health **50 → 56 → 62** replays the same day's amount; day 14 yields 10,
  turn 2 yields 4, return to turn 1 under the new master yields 5.
  Paired no-call controls retain both weather shuffles and unique-item
  mint allocations **65748 / 68643**, proving subsystem isolation.

Two deterministic bounded pinned C runs produce **564 sections / 40140 lines /
684313 bytes**, exactly matched by Go's raw stdout/stderr/status and real
state/log/order observations. Prior command input **12248 lines / 218785 bytes**
and output **553 sections / 39215 lines / 669382 bytes** remain exact byte
prefixes: **zero old bytes changed**. New input is **12580 lines / 224361 bytes**.
Entire move input **3150 lines / 53589 bytes** and output **135 sections /
7150 lines / 142755 bytes** remain identical. Literal C trailing whitespace
in recorded here lists is retained, not normalized.

A removed ephemeral expected-134 case queued, loaded and dispatched ATTACK.
Normal and `--update` both rejected `out-of-scope C dependency: v_attack`
with status **2**, before abort-status acceptance and before recording.
Expected-output SHA-256 stayed
`ed21ae2739e9eb4ba189ee09b83bf7f307282fb3eb3766c46c24ebb1f129fdf1`.
The probe and scratch key-search data are removed; no permanent test depends
on that private missing-handler name. All new permanent sections require
status zero; no successful unimplemented handler was introduced.

Command (update and normal), move, helper, rng, code, loadsave, turn-driver
and default turn comparisons pass sequentially before bounded vet and fresh
default all-tests. All engine processes retain shared run_limited supervision.
DailyEvents positions and weekly gating are unchanged. Real-C schedule-suite
weekly boundaries pass, but callback traces accept only scheduling, not this
healing gameplay. Whole-world turns remain opt-in and were not run; the known
day-1 full-turn boundary is still match_all_trades (#81). add_char_damage,
unrelated commands/combat, later daily/monthly phases and reports remain
explicit boundaries. No full-world/stage-exit success or task acceptance is
claimed. No newly verified public-engine bug or blocker was found; arrival
checks still show only enhancements, no local bug/follow-c. Known nonblocking
G3 #309/#352/#353/#354/#355 remain unchanged. Version **0.57.0-alpha**; changes
remain uncommitted, unstaged and unaccepted for parent review.

## Daily prisoner, ring-aura, curse and dog closure (#80)

This extends the existing command/lifecycle harness, rather than introducing
the duplicate `upkeep-compat` named in the original issue acceptance wording.
The pin remains `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`. No new source file
is read, so the dependency table/affected.sh metadata is unchanged.

| source | complete definition | exact lines | count |
|---|---|---|---|
| stack.c | weekly_prisoner_escape_check | 345–375 | 31 |
| day.c | dogs_bark_at_hidden_chars | 489–516 | 28 |
| day.c | increment_stone_ring_aura | 777–803 | 27 |
| day.c | noncreator_curse_erode | 901–963 | 63 |
| u.c | bark_dogs | 2633–2665 | 33 |

These **182 C lines** are extracted unchanged, including disabled preprocessor
branches; only phase instrumentation is omitted. The production port implements
the active branch, not disabled design alternatives. The real **daily_events
2206–2271 (66 lines)** and original **day.c 189–235** calendar prefix are also
extracted. No scheduler model or successful gameplay stub replaces a callback.
The existing persistent world Streams supplies EntPrisoner/EntBark/CalDay;
the original C entity prefix, tags, guarded turn/day streams and leaf keys
remain intact. No production RNG wrapper or adapter was introduced.

**16 appended small sections** compare actual C and Go raw stdout, stderr,
status, saveBox bytes and real logs/knowledge/touched/output/VECT state:

- Outdoor weekly roll **12** escapes at chance 2, but fails in an inn at
  chance 1; roll **1** escapes from that inn through real outward movement.
  Captive captors, ocean flyers and release-sworn prisoners skip. Direct
  sibling snapshot ordering and nested prisoners retain the real live outer
  character chain. Same-owner follower hounds halve the threshold and bark
  plurally; changing their owner permits the identical prisoner roll to escape.
- A ship prisoner over ocean is not mistaken for a flying prisoner: it really
  swims to the first island, clears moving/contact/ocean state, marks knowledge
  and weather/touched state, and emits both captor and captive messages/logs.
- Ring aura uses the existing max_eff_aura contract: below/at/above maximum,
  nonmagicians, nested buildings, innate negative clamp, quantity-independent
  artifact bonuses, auraculum, zero-quantity and deleted-item skips. Repeated
  calls stop exactly at the effective cap; missing magic stays unallocated.
- Curse skips missing item magic, zero strength, creator, oath, zero quantity
  and deleted items. Contract/fear/unsworn holders log once per qualifying
  inventory entry, not per quantity, on repeated calls. Loyalty rate **37**
  stays unchanged. The active C `#if 1` is literally log-only NOTYET.
- Dogs require a direct really-hidden character, not merely hide_self or a
  nested hidden descendant. Real stack item sums include same-owner nested
  hounds; every direct character supplies its own `(who, hound ordinal)` key.
  Seeded zero/one/multiple barks distinguish no output, **The hound**, **A
  hound**, and **The hounds**; VECT includes descendants and the location.
  Day changes, guarded same-turn master reload, next turn and returning to
  the old turn retain the actual entity-stream contract.
- Real daily_events on a small map observes actual calendar picks **curse 8,
  dog 15, faery 25**. Calls on days **6, 7, 8, 14, 15, 16** demonstrate
  nonfiring days, weekly escape, log-only curse, then ring aura before dog
  output. This is actual gameplay composition, NOT phase-trace acceptance.
  Day 25/faery and full-world integration are not claimed by this case.
- Paired no-call controls preserve both weather shuffles and unique-item
  mint IDs **65748 / 68643** after real dog/weekly draws.

The C driver runs twice deterministically: **580 sections / 41331 lines /
705663 bytes**, matched exactly by Go. Previous input **12580 lines / 224361
bytes** and output **564 sections / 40140 lines / 684313 bytes** remain exact
prefixes, zero old bytes changed. New input: **13058 lines / 232235 bytes**.
Move cases/recordings remain byte-identical (**3150 / 7150 lines**, **53589 /
142755 bytes**). Literal C trailing whitespace in saves is retained.

A removed temporary expected-134 case queued/loaded/dispatched ATTACK.
Normal and `--update` both rejected `out-of-scope C dependency: v_attack`
with **status 2**, before expected-status handling or recording. Output hash
remained `ea769dd8b92f50ae7a8b47bc93b2a219ed2f78c1657f3fdd3d29f8c68d823d4a`.
All permanent new sections require status zero; every engine invocation still
uses shared run_limited supervision. Only the four owned production NYIs and
the harness daily_events stub are removed; command/combat boundaries remain.

**Upstream evidence:** closed G3 #233 (under #227) explicitly tracks coverage
of the historical log-only curse path, not implementation of its public effect.
Pinned use.c:134 exposes USE 886; art.c:664–722 validates/casts the curse,
charges aura and advertises strength. day.c:948–958 logs NOTYET with the actual
loyalty effect disabled. Filed nonblocking [G3 #356](https://github.com/playbymail/olympia-g3/issues/356),
assigned to @me, for that source-verified supported-command limitation, not a
new regression claim or malformed-state crash. Production links the intentional
reproduction. No full public-command turn reproduction is claimed. Existing
G3 #309/#352/#353/#354/#355 remain deliberately unchanged and source-linked.

The eight comparison tools pass sequentially. Full-world/stage-exit checks
are opt-in and not run here; the actual turn still stops at day-1
match_all_trades (#81). Production dailyEvents order/weekly gates are unchanged.
Version **0.58.0-alpha**. Changes remain uncommitted, unstaged and unaccepted
for parent review; no issue closure, push or stage tick is implied.

## Day-1 trade settlement and live validation (#81)

This extends **command-compat**, rather than the duplicate economy-compat
named in the original issue. This harness already owns real inventory/lore/
token transfer, arrival, command completion, market lists, output and saves.
The pin stays `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`; no new source
file is read, so affected.sh and the source dependency table are unchanged.
Earlier historical #81 boundary notes above are superseded by this section.

The script extracts these complete, unchanged definitions from real buy.c:

| definition | exact source lines | count |
|---|---|---|
| reduce_qty | 475–503 | 29 |
| attempt_trade | 506–632 | 127 |
| match_all_trades | 698–729 | 32 |
| check_validated_trades | 735–746 | 12 |

These **200 C lines** reuse the accepted real seller/buyer lists, scan_trades,
match_trades, market_here and investigate_possible_trade. Only phase counters
are omitted. Original gm.c scalar definitions at **10, 11, 22** supply
gold_pot_basket, gold_trade and gold_opium; the production port mutates int32
fields, not observer-derived totals. There are **no settlement RNG draws or
new RNG adapters**. Shared transfer machinery keeps its existing entity/mint
streams, tags, guards and draw order per pinned rng-streams.md.

Production removes only the owned attemptTrade/matchAllTrades stops and the
two conditional completion stops in input.go. Both callers retain their C
nonempty-list guards. The validator walks the **live growing list** and clears
it **after** the loop; temporary market lists retain pointers to real trades
while mutations remain visible. Unported BUY/SELL/FIND and monthly market
restock (#83) are not implemented or blessed by tests.

**26 appended small sections** compare raw C/Go stdout, stderr and status,
actual saveBox/trade/inventory bytes, faction/special logs, lore/knowledge,
ownership/NP mutations and gold metrics:

- Cost mismatch, unequal quantities, buyer gold reserve at the **seller's**
  price, seller inventory reserve, and reserves above stock/balance.
- Equal-price seller here-order and character-first/city-last ordering even
  when the city is cheaper; BUY and SELL entry scans, mismatched items and
  self-offer skips. City ordinary sales generate stock without consuming its
  insufficient inventory, and city purchases ignore gold/reserve inventory.
- City ordinary accounting; pots/baskets **22 gold**, opium **33 gold**, and
  actual LOG_SPECIAL output. None/seller/buyer/both cloaks and city-side cloak
  display preserve exact C precedence and messages.
- Real city/character unique-item ownership/lore transfers. An NPC-token sale
  transfers its actual controlled unit, refunds **1 NP** then charges **1 NP**,
  updates faction membership and emits the real token/ownership logs.
- A live cascade starts with only **1000** queued by real add_item. Settlement
  appends more work; **1002** then buys the third good from **1003**, who ends
  with **10 gold**, before the queue clears. A snapshot/fixed-length validator
  cannot complete this chain. Real do_command and evening_phase separately
  drain real inventory-triggered queues; real move_stack arrival settles too.
- Actual process_orders on a small map executes day-1 global matching, then
  all 30 days through real daily callbacks. Valid tiny Hades/Faery regions and
  factions allow their real autonomous setup, not successful callback stubs.
  This is not the opt-in full-world turn or post_month/stage-exit acceptance.
- Two deliberate raw invalid contracts retain exact reduce_qty/item-pair
  assertions/status 134. A deliberately inconsistent missing-unique-inventory
  case reproduces C's no-replacement behavior; this and malformed raw trades
  alone are **not public-command defect proof**. No new bug was verified or
  filed; known nonblocking G3 #309/#352/#353/#354/#355/#356 remain unchanged.

Two deterministic C runs produce **606 sections / 42545 lines / 728447 bytes**,
matched exactly by Go. Previous **580 sections / 41331 lines / 705663 bytes**
and input **13058 lines / 232235 bytes** are exact prefixes, zero old bytes
changed. New input: **13560 lines / 242037 bytes**. Move recordings remain
unchanged. Literal C whitespace is retained.

A removed temporary expected-134 section dispatched ATTACK. Normal and
`--update` both rejected `out-of-scope C dependency: v_attack` with **status 2**,
before abort-status acceptance or recording. The recording stayed SHA-256
`33aaf710e8cc6871ebed0b62072d0ce6688ea4930185a5b798068a8b4d42faaa`.
Every engine process still uses shared run_limited supervision. Only actual
definitions replace the three former C trade stubs; unrelated missing
dependencies remain explicit failures.

Run command/move/helper/rng/code/loadsave/turn-driver/default-turn sequentially,
then bounded vet/default tests. Whole-world/stage-exit checks remain opt-in
and are left for parent acceptance; no new full-world boundary is inferred.
Version **0.59.0-alpha**. Changes remain uncommitted and unaccepted; no commit,
push, issue closure, stage tick or stage-exit success is implied.

## #82 — city gold and existing-world production

Acceptance clarification: [#82 comment](https://github.com/mdhender/g3g0/issues/82#issuecomment-6461526153).
This extends this harness, not a duplicate economy tool. Complete original
definitions at the same pin are extracted, retaining `#line` and the real
inventory, traversal, topology, trade queue and save-box dependencies:

| Pinned C source | Definitions/ranges (including return declaration) |
| --- | --- |
| `day.c` | `compute_civ_levels` 1792–1929; `post_production` 1932–1942; `move_city_gold` 2005–2024 |
| `produce.c` | original macros and complete `terr_prod`/`mine_prod` tables 10–88; `replenish` 91–99; `mine_production` 102–117; `location_production` 120–169 |
| `seed.c` | `seed_taxes` 640–704 |
| `u.c` | `clear_temps` 1352–1362 (existing Go helper reused) |

Ten new small sections independently run real production helpers:

- Asymmetric civ propagation reaches temps **9 → 18 → 36**, civ **2/4/9**,
  with no diffusion into an ocean. A second case distinguishes an ocean's
  own crown bonus as a diffusion source. First nested castle/city/tower win
  over later duplicates; major/minor cities and per-province flag resets are
  distinct. Repeated calculation clears old location scratch/civ state.
- Every terrain row and mine depth **1–20**, including empty rows and sparse
  mithril yields; depth uses signed-short shaft/3 truncation and clamps above
  20. Real replenish checks below/equal/above floors and preserves excess stock.
- Coastal and island flotsam, inland controls, poppy cap-before-double,
  consumed one-shot tending flag, next-month recapping, no-misc poppies and
  swamp tending left untouched. A shared exdr shuffle follows production.
- City tax transfer is additive and repeatable with empty cities; reseeding
  resets tax/menial/thief stocks, skips oceans/buildings/regions, preserves
  opium-before-pillage rounding (**46**, **30**, **10** tax cookies) and real
  serialized inventory/LO/SL state.
- Raw zero/negative mine depths and opium-10 tax state preserve exact C
  assertions/status 134. These are malformed helper/saved-state contracts,
  **not publicly reachable defect claims**. No new C bug was verified.

`postProduction` runs real civ then stops honestly at `location_trades` (#83),
before location production, tax reseeding or the post flag. A Go boundary
test proves that ordering, not successful composition. A removed temporary
expected-134 case called real C `post_production`; both normal and `--update`
rejected its explicit missing `location_trades` dependency with **status 2**,
before assertion acceptance/recording, leaving the recording unchanged.
No successful C gameplay substitute is introduced.

Two deterministic C runs match Go: **616 sections / 43868 lines / 746528 bytes**.
All previous **606 sections / 42545 lines / 728447 bytes** and input
**13560 lines / 242037 bytes** remain exact byte prefixes (zero old bytes
changed). New input: **13839 lines / 247379 bytes**. Recording SHA-256:
`a3c184638f69745d3a00ef5b17c5be60d524182d681d1ffd176fbf02eef801f0`.
Literal C trailing whitespace is intentional and is never normalized.

Command/move/helper/rng/code/loadsave/turn-driver/default-turn passed
sequentially, followed by `go vet ./...` and fresh bounded default tests.
Every engine process used shared `run_limited`, with `G3_REPO` pointing to the
bare cache. No new pin/source dependency, local RNG adapter, BUY/SELL or
new-world command. New-world load initialization remains explicit stage 5;
postMonth call positions are unchanged. No full-world run or stage-exit
success is claimed. Version **0.60.0-alpha**, all changes **uncommitted**.

## #83 — real monthly markets and successful post-production composition

Scope clarification: [#83 source-verified comment](https://github.com/mdhender/g3g0/issues/83#issuecomment-6461652386).
This supersedes #82's temporary `location_trades` boundary above. The task
fits one bounded closure: no child tasks or duplicate economy tool are needed.
Pinned C remains `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`; only bare-cache
`git show`/`git archive` into `tmp/` were used, never an upstream working tree.

| Pinned source | Complete extracted definitions/data (inclusive ranges) |
| --- | --- |
| `buy.c` | original `economy_rng` declaration 34; original economy macros and `begin_economy`/all `econ_*` helpers 215–312 |
| `buy.c` | `find_trade` 809–826; `new_trade` 851–870; `add_city_trade` 1161–1172 |
| `buy.c` | original `opium_data`/`MAX_OPIUM_ECON` 1195–1212; `opium_market_delta` 1215–1269; `expire_trades` 1271–1309 |
| `buy.c` | `loc_trade_sup` 1318–1364; `trade_suffuse_ring` 1367–1395; `location_trades` 1398–1414; `is_tradegood` 1500–1516 |
| `art.c` | `new_suffuse_ring` 1437–1488 |

The Go helpers reuse persistent `w.streams.BeginEconomy/EconStock/EconRing/
EconCost` exactly, plus the accepted real mint, settlement, inventory/lore
and save-box lifecycle. No new stream or RNG adapter. `loop_trade` is the
filtered pointer snapshot at `loop.h:332–345`, **not** a live append loop.
The separate raw opium scan includes sold-out BUY entries and takes the
first duplicate; raw expiration ages zero-quantity entries and deletes
adjacent expired entries without skipping them. The already-live validated
settlement queue remains unchanged and its earlier recording is preserved.

Ten new recorded sections distinguish:

- Sold-out, partial, no-sale and absent opium orders, bounds **0/7**, first
  duplicate precedence, both cloaks and actual quantities/prices including
  levels 1 and 6. Lower demand does not discard excess outstanding BUY stock.
- Zero-quantity offers are not reused by `new_trade`; annual production,
  **turns 1/2/8/9/10**, override-only production, consume ignoring seasons,
  stock floors/surpluses and copied cost/cloak/lifetime metadata.
- Tradegood aging/deletion, inventory held by either character retaining BUY
  entries (including zero-quantity/negative-expiry entries), city inventory
  not retaining them, copied lifetimes and deletion after holders consume
  their goods. These are exact helper contracts, not new public bug claims.
- Asymmetric monthly cities, a non-city untouched, gates **1/2/3**, distinct
  market/turn keys and ring prices **500/600/950/1050**. All five actual
  `new_suffuse_ring` variants serialize use keys **10–14**, lore **938–942**,
  Golden ring name, weight and unique owner, with real shared mint IDs.
- Existing SELL stock prevents a second mint; actual `match_trades` purchases
  clear SELL quantity, move unique ownership, debit **950** each and queue each
  unique ring for lore (two rings remain two entries). Same-turn sold-out
  replacement works; turn 9 rejects restock,
  turn 10 creates a different ring at **600**, followed by a shared-mint control.
- Actual `post_production` executes civ → markets → production → taxes → flag,
  twice across turns 8/9: civ **1**, **25** poppy opium, opium-adjusted city
  taxes **80 then 70**, `post 1`, Faery/Cloud ring creation and ordinary-market
  exclusion. The obsolete `TestPostProductionStopsAtLocationTrades` is removed;
  the successful C composition owns that ordering contract now.

Missing-C-dependency rejection still precedes expected assertion/status
handling and any `--update`. All new market sections must return **0**;
malformed fixtures cannot silently become accepted assertion recordings.
During construction an incomplete tiny map lacked required resource entities
and province routes; it was corrected, not accepted as gameplay evidence.

Two C runs and Go match **626 sections / 45421 lines / 770256 bytes**.
Previous **616 sections / 43868 lines / 746528 bytes** and input
**13839 lines / 247379 bytes** remain exact prefixes, with zero old bytes
changed. New input: **14174 lines / 253908 bytes**. Recording SHA-256:
`e3e311370c0273aabce56934ce9ece9e4bdcd5c330c79e13ba01e6ca0273995b`.
Raw C trailing whitespace remains intentional.

`seed_city_trade` (`seed.c:483–579`) is INIT-only through `seed_city` and
world/region builders, not called by monthly markets. Trade-route init/save
already exists. Neither that INIT driver, FIND/BUY/SELL commands nor new-world
load initialization is included here: `loadDb` retains its explicit stage-5
boundary. No postMonth ordering change, new source dependency/pin, independent
fix or full-world success. Version **0.61.0-alpha**; all work **uncommitted**.

Final checks passed in order: command **626/45421**, move **135/7150**,
helper **37/1487**, RNG **321 results**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152** plus phase **5/484**, calendar **2/119**, schedule
**4/1801**, then default turn (raw status/stdout/stderr and one state file).
Bounded `go vet ./...` and fresh `go test -count=1 ./...` passed afterward;
test binaries used the shared `run-limited.sh` executor. Every compat call
used `G3_REPO=$PWD/tmp/upstream.git`. No full-world run; parent owns it.
No newly verified upstream defect or blocker was found. `git diff --check`
reports only literal C recording trailing whitespace, intentionally retained.

## #84 — faction accounts, castle taxes and seasonal announcements

Pin remains **6a25f2d252ce45bd5adbe3baaccb3be3eed56279**. The existing
command harness extracts these complete definitions, not modeled callbacks:

| File | Complete definitions / pinned lines |
| --- | --- |
| `day.c` | `add_claim_gold` 635–653; `add_noble_points` 655–676; `add_unformed_sup` 678–696; `add_unformed` 699–719 |
| `day.c` | `collect_taxes` 1756–1789; `special_locs_open` 2046–2085; `special_locs_close` 2088–2127; `clear_orders_sent` 2130–2144 |
| `c1.c` | `next_np_turn` 881–901 |
| `gm.c` | original `gold_taxes` scalar at 13 |

Only phase instrumentation is omitted. No source dependency/pin change:
all four files were already tracked. Go reuses its existing `clearOrdersSent`,
real `newEnt`, `genItem`/`consumeItem`, `addNP`, owner/NP/lord helpers, output
and database writer. It adds no alternate allocator, account model or separate
upkeep/garrison harness. All ten new `accounts` sections require status **0**.
Missing-dependency rejection precedes this status contract and `--update`.
A temporary negative driver emitted the missing-dependency marker with status
134 at the first accounts section: both ordinary and `--update` runs rejected
it with status **2**, and the recording hash stayed unchanged. The probe was
removed afterward.

Recorded cases cover:

- Unequal claim balances **9/71 → 34/96 → 59/121**, while system/NPC/silent
  players retain their balances. A regular player lacking a player record
  receives claim inventory without allocating that record; NP querying later
  allocates it. Clearing orders covers positive/negative values and missing
  records without allocation.
- `next_np_turn` ignores asymmetric `first_turn` values **1/777**. Turns
  **0/1/7/8/9/15/16/32767/-1** retain the C formula, including a turn-zero
  grant, **7 → 1**, **8 → 0**, **9 → 7**, and signed-short NP wrap. Missing
  records allocate even on a non-grant turn; nonregular players are skipped.
- Unformed list lengths **0/4/5/6**: replenish exactly to five, leave five/six
  untouched, include allocated NPC/system/silent records, and skip missing
  regular/NPC records. Repeated calls do not allocate extras. A separate
  malformed-list control verifies that stale/duplicate entries count verbatim;
  it is not public bug evidence. Raw list state and real saved entities are
  both compared. The shared mint produces a later unique item **62565** after
  five unformed allocations versus **70368** after a full-list no-op; mixed
  faction replenishment produces **60539**.
- Castle iteration/ownership uses real `building_owner`, including its first
  occupant behavior, unfinished-castle exclusion and ownerless castles. An
  odd **11** cookies becomes **5** gold for the first owned castle, zero for
  the next; unowned provinces retain **7/1** cookies until ownership changes.
  One cookie is consumed despite producing zero gold. A reversed-construction,
  asymmetric-owner control consumes **2147483647** cookies into **1073741823**
  gold, preserving real prisoner output suppression. Repeated tax collection
  emits ordinary zero-tax logs; totals, inventories, output bits and saved
  state are compared.
- Both seasonal functions execute all four real output branches, with
  Summerbridge precedence over colliding Uldim flags. Uldim **1/2**, invalid
  Summerbridge-only flags, cities and regions do not announce. Real faction
  logs include `LOG_CODE` framing and province messages; saves show unchanged
  route flags. The existing real-C turn-driver phase suite owns post_month's
  month-2/month-6 calendar gates; it proves scheduling only.

Two pinned C runs and Go match **636 sections / 46887 lines / 791460 bytes**.
Previous **626 sections / 45421 lines / 770256 bytes** and input
**14174 lines / 253908 bytes** remain exact prefixes, with no old byte changed.
New input: **14488 lines / 259432 bytes**. Recording SHA-256:
`f987d718a9b667f3687b97280debe46fd2e01418a678ae3571bcdb7db8186772`.
Literal C trailing whitespace is retained.

`garrison_gold` remains explicit at its existing post_month position; #116
follows maintenance #87. No full post_month integration or full-world success
is claimed, and no maintenance/account-command/report-body work is included.
Version **0.62.0-alpha**; changes are **uncommitted**. Scope clarification:
https://github.com/mdhender/g3g0/issues/84#issuecomment-6461896006.

Final checks passed sequentially: command **636/46887**, move **135/7150**,
helper **37/1487**, RNG **321 results**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152** plus phase **5/484**, calendar **2/119**, schedule
**4/1801**, then default turn (raw status/stdout/stderr and one state file).
Bounded `go vet ./...` and fresh `go test -count=1 ./...` passed afterward;
test binaries used the shared `run-limited.sh` executor. Every compat call
used `G3_REPO=$PWD/tmp/upstream.git`. No full-world run, new verified bug or
blocker. `git diff --check` findings are only literal C trailing whitespace.

## Monthly aura, shrouds and loyalty (#85)

This extends the same real-C command harness, not a separate upkeep tool.
Pin: `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`. Complete definitions newly
extracted (inclusive original source ranges):

| Source | Definitions and ranges |
| --- | --- |
| `day.c` | `increment_current_aura` 722–774; `decrement_ability_shroud` 806–825; `decrement_loc_barrier` 828–853; `decrement_region_shroud` 856–876; `decrement_meditation_hinder` 879–898 |
| `day.c` | `loyalty_decay` 1007–1075; `relic_decay` 1078–1108; `pillage_decay` 1111–1130; `hide_mage_decay` 1945–1963 |
| `scry.c` | `notify_loc_shroud` 333–359 |
| `u.c` | `loyal_s` 1569–1605 |

The original guarded upkeep prefix now spans **day.c 50–115**, including the
whole `up_loyal` definition (109–114). Its tags, persistent turn guard and
keyed leaves are C source, not harness replacements. The original
`quest.c:62` `nowhere_loc` global is also extracted alongside `nowhere_region`.
`scry.c` is the only newly tracked source dependency; the pin is unchanged.

Production reuses effective aura/auraculum, the world's persistent
`streams.UpLoyal`, real `unitDeserts`/`setLord`/interruption, NP accounting,
extraction, inventory transfer/lore, output and save-box implementations.
Nothing returns fabricated successful gameplay. `postMonth` is unchanged:
the nine new phases remain in exact C order after unformed settlement.

The twelve new `monthly` sections all require status **0**. Small cases prove:

- Aura recovery per positive inventory **entry**, not quantity or bonus:
  quantities **19/3** and bonuses **7/-2** give **0 → 6 → 12** with the
  auraculum. The effective **23** cap is reached from **22** without overshoot;
  other cases cover one-below cap, already-over-max, negative innate aura,
  an unheld/held auraculum, nonmagicians and absent magic/character records.
- Shroud/barrier/hinder/hide timers **2/1/0/negative** and absent records,
  exact barrier disappearance logs, nested skilled recipients versus
  learning/unskilled characters and noncharacter subtree pruning. Pillage
  respects recent-loot, decrements negative values and wraps signed-byte
  **-128 → 127**, unlike the positive-only magical timers.
- Every loyalty kind, positive/negative fresh-hire skips and absent records.
  Contract **49/50**, zero/negative rates and **501/510** distinguish the
  pre-decrement desertion gate and the **max(50, rate/10)** charge. Fear
  **2/1/0/-1** tests after decrement, then clamps misses at zero. Keyed hit
  and miss controls replay across days and same-turn master changes, and
  reseed on a new turn. A real fear miss survives twice, then deserts on
  the next-turn hit. Summon **1 → 0** stays stacked until the following call;
  zero/negative release unstacks, restores actions and retains its owner.
- Real fear desertion interrupts running WAIT, clears run/wait queues and
  queued orders, removes attitudes/pledge, resets backlinks, extracts the
  stack and reswears dependents to their old lord. Its independent-faction
  destination correctly leaves **11 NP** unchanged (no regular-player NP
  refund/charge). Saved unit/owner/here order, output and order-file bytes
  are compared, not just a callback trace.
- Owned relic timers **1/2/0/-1**, an absent magic record, nonrelic exclusion
  and short **-32768 → 32767**. Real acquisition queues lore **410/411**;
  transfer moves the actual inventory and unique owner to nowhere, retains
  known/lore state and emits exact vanish messages. Repeated negative decay
  remains negative and transfers again, as C does. These raw signed-field
  controls alone are not evidence of a publicly reachable upstream defect.

Missing-dependency rejection still precedes monthly status checks and
`--update`: a temporary driver emitted the marker with status 134 at section
636; ordinary and update modes both rejected it with **status 2**, leaving
the recording hash unchanged. The probe was removed after verification.

Two pinned C runs and Go match **648 sections / 48420 lines / 815585 bytes**.
Previous output **636 sections / 46887 lines / 791460 bytes** and input
**14488 lines / 259432 bytes** remain exact byte prefixes. New input:
**14923 lines / 267413 bytes**. Recording SHA-256:
`490f2f26175f108f7c5f7082dfc00f52ee947dd369f42a04abca496fde747673`.
Expected bytes are solely C-generated, including trailing whitespace.

Sequential final checks passed: command **648/48420**, move **135/7150**,
helper **37/1487**, RNG **321 results**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**,
default turn (raw status/stdout/stderr and one state file), then bounded vet
and fresh bounded default tests with the shared engine executor. Every tool
used `G3_REPO=$PWD/tmp/upstream.git`. No new verified bug/blocker, full-world
run, independent fix, commit, push, task closure/tick or stage-2 work.

Version **0.63.0-alpha**; all changes remain **uncommitted** for acceptance.
Scope clarification: https://github.com/mdhender/g3g0/issues/85#issuecomment-6462040567.
`garrison_gold` still stops at its unchanged #116 boundary; no successful
full `postMonth`, maintenance/income, command-family or report-body scope.

## Monthly inn and temple income (#86)

This extends the same real-C harness, not a duplicate upkeep tool. Pin remains
`6a25f2d252ce45bd5adbe3baaccb3be3eed56279`. New complete definitions:

| pinned source | definition / range |
|---|---|
| `olympia/day.c` | `inn_income`, **1421–1489** |
| `olympia/day.c` | `temple_income`, **1492–1522** |
| `olympia/loc.c` | `count_loc_structures`, **352–366** |
| `olympia/day.c` | original upkeep stream/tags and complete guarded leaf prefix, **50–157**, including `up_income` **150–155** |
| `olympia/gm.c` | original `int gold_inn = 0;` **12**, `int gold_temple = 0;` **19** |

Only phase instrumentation is omitted. Go reuses persistent
`w.streams.UpIncome`, real `buildingOwner`/`hasSkill`, inventory/lore,
output routing and `saveBox`; both accounting scalars and intermediate amounts
are C-width signed ints (`int32`). Nothing substitutes a modeled owner,
income stream, inventory or successful unported phase. No new source
dependency/pin/affected.sh entry is needed. `postMonth` call order is unchanged;
the earlier `garrisonGold` boundary remains #116 after maintenance #87.

Source-backed acceptance clarification:
https://github.com/mdhender/g3g0/issues/86#issuecomment-6462255396.
The issue's original temple-sharing/occupancy wording is not the C contract:

- Inn sharing uses **subloc(inn)**, not province, and counts direct **T_loc**
  children matching **sub_inn OR 0**, including unowned structures. The count
  helper checks either subtype once and does not descend or count items/chars.
  Actual unequal incomes are **12/72/13** with five sharing structures: the
  middle inn adds its **60** traveller bonus **after** sharing. Empty and
  unfinished inns are skipped; nested inns share only their immediate parent.
- Inn base endpoints **50/75**, traveller gate **1/2/8**, and actual bonus
  endpoints **50/130** are C-recorded. Pillage suppresses an otherwise-hit
  traveller gate, divides before sharing and selects all three exact flavor
  messages. Pillage **127** yields zero without creating inventory entries.
  Out-of-order fixture allocation and placement still emit subtype-chain order.
- Temples have **no sharing or occupant/capacity scaling**. One- and
  three-occupant priest-owned temples each pay **100**, while an empty temple,
  a first nonpriest followed by a priest, learning religion, absent character
  records/skills and unfinished temples do not pay. Local pillage pays **33**;
  an unpillaged nested temple still pays **100**, not its province's **33**.
  A first captive priest remains the owner: inventory/metrics receive the
  offering even when real prisoner-output routing suppresses its message.
- Same-turn day/master changes replay **120/120**; the next turn pays **61**;
  returning to the old turn after the master change reseeds and pays **57**.
  Scalars initialized at **2147483000/-2147483000** accumulate real income
  twice without reset or signed overflow. Saved inventories, record presence,
  here order and raw logs are compared, not just totals.
- Real-income/no-call pairs retain identical weather shuffles/strengths,
  exit shuffles and minted IDs **66288/61234** across the turn boundary.

The twelve ordinary `income` sections require status **0**. The explicitly
named invalid-count section alone requires **134**: `income count 9999 37 0`
records `assertion failure: olympia/loc.c (358): valid_box(where)` from the
original `loop_here` macro expansion. Go preserves this assertion at its
call-site before `loopHere`, rather than silently returning or producing a
nil-dereference/runtime panic. This is the parent's task-review correction,
not an independent gameplay change or publicly reachable upstream defect.
A temporary compiled
driver wrapper emitted `out-of-scope C dependency: income_probe` and status
**134** at section **648**. Both ordinary and `--update` runs rejected it
with status **2** before income status handling and left the recording hash
unchanged. The temporary probe was removed. Raw malformed-field tests alone
are not public upstream-defect evidence; no new verified bug/blocker was found.

Two real pinned C runs and Go match **661 sections / 49624 lines / 836364
bytes**. Previous output **660 sections / 49617 lines / 836187 bytes**, and
its input **15298 lines / 274081 bytes**, remain exact byte prefixes.
The accepted **648-section / 48420-line / 815585-byte** output and
**14923-line / 267413-byte** input prefixes are also byte-verified unchanged,
retaining the earlier 636-section prefix transitively. New input:
**15301 lines / 274174 bytes**. Recording SHA-256:
`2ffa80dd6dd61f28cb3b10a57f7efae146139376ed179ec0d49f07e834f0f91f`.
Expected bytes come only from real C, including its trailing whitespace.

Sequential final checks passed: command **661/49624**, move **135/7150**,
helper **37/1487**, RNG **321 results**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**,
default turn (raw status/stdout/stderr and one state file), then bounded vet
and fresh bounded default tests with the shared `run_limited` engine executor.
Every compat call used `G3_REPO=$PWD/tmp/upstream.git`. No full-world run,
independent fix, commit, push, issue closure/tick or stage-2 work.

Version **0.64.0-alpha**; changes remain **uncommitted** for parent acceptance.

## Monthly maintenance, starvation and animal deaths (#87)

This extends command-compat rather than duplicating its stack/inventory,
command, output and save closure in upkeep-compat; scope clarification is on
[#87](https://github.com/mdhender/g3g0/issues/87#issuecomment-6462462502).
Pin remains `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`.
Complete original definitions are extracted, with original `#line` locations:

| source | definitions / inclusive ranges |
|---|---|
| `day.c` | accepted `maint_cost` **1171–1199**; `men_starve` **1200–1301**; `unit_maint_cost` **1304–1318** |
| `day.c` | `charge_maint_sup` **1325–1349**; `charge_maint_costs` **1352–1376**; `animal_deaths` **1379–1418** |
| `day.c` | existing original guarded upkeep prefix **50–157**, including `up_starve` **117–122** and `up_animal` **129–135** |
| `move.c` | complete `check_captain_loses_sailors` **1175–1247**, reusing the accepted Go helper |

Only phase instrumentation is omitted, as before. No source dependency or
pin change is needed. Production reuses `maintCost`, `stackHasItem`, actor-first
`autocharge`, `consumeItem`, formatters and persistent `w.streams.UpStarve` /
`UpAnimal`; no successful command/captain or RNG replacement is supplied.
Existing postMonth call positions remain unchanged, with the earlier
`garrisonGold` stop still #116. Full SAIL and later decay remain unported.

The **22 new sections** cover:

- Asymmetric expensive-first versus cheap-first inventory: five gold funds
  one elite guard in the former but four peasants in the latter. Payment
  repeatedly cycles items and skips unaffordable tariffs; casualties likewise
  cycle inventory, not cost order. Seven unfunded men lose three, not two.
- Exact affordability, one short, no funds, partially affordable cheap troops
  with unusable gold left, empty/free units and the full signed-int maximum
  cost **2147483647**, without an overflowing C expression.
- Actor-first friendly stack borrowing, foreign rich stack exclusion, friendly
  prisoner funds, regular-prisoner/unsworn charging and nonregular filtering.
  Noble items are excluded by unit cost only: starvation and animal mortality
  still include them. Invalid/nonpositive inventory is filtered by loop_inv.
- Real moving roundship captain, parsed SAIL in STATE_RUN with positive wait:
  loss of two sailors and one pirate changes wait **4 → 7**, emits both real
  crew-loss messages and LOG_SPECIAL lines, and saves real command state.
  No SAIL callback is invoked or substituted.
- Peasant `starved`, both keyed nonpeasant flavors, same-turn master/day
  changes, repetition and next-turn reseeding; real repeated herd losses,
  item-specific leaves, hit/miss, noble/prisoner animals and faction filtering.
- Independent original-C leaf probes: `(1698,51,1)` is **9** and dies;
  `(4403,51,1)` is **10** and survives; `(1094,51,0)` is **5**, but its actual
  first individual uses index **1**, rolls **614** and survives. Thus neither
  zero-based indexing nor a `<=10` comparison can match the recording.
- Paired call/no-call controls prime the same upkeep guard, preserving sibling
  income leaves and downstream weather, exit shuffle and mint bytes.
- Raw invalid helper calls reproduce loop_inv assertions at day.c **1310**
  and **1223**, and fully paid direct starvation asserts at **1259**.

All ordinary upkeep sections require status **0**; named invalid sections
require **134**. A temporary harness wrapper emitted an out-of-scope dependency
marker and aborted at section **676**, which otherwise expects 134. Both normal
and `--update` runs rejected it with status **2**, before assertion acceptance
or recording replacement; the recording hash stayed unchanged. The probe was
removed. No new verified upstream defect or blocker was discovered.

Two real C runs and Go match **683 sections / 50572 lines / 853708 bytes**.
The accepted **661-section / 49624-line / 836364-byte** output and
**15301-line / 274174-byte** input remain exact byte prefixes. New input is
**15712 lines / 281967 bytes**. Recording SHA-256:
`ac0900cfbee0f46d55a0b88053e29457429354c090ec55c5cc0281f7045c05c7`.
Expected bytes come only from actual pinned C, including trailing whitespace.

Sequential checks passed: command **683/50572**, move **135/7150**, helper
**37/1487**, RNG **321 results**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**,
default turn (identical raw status/stdout/stderr and one state file), then
bounded vet and fresh bounded default tests using shared `run_limited`.
Every compat call used `G3_REPO=$PWD/tmp/upstream.git`. No full-world run,
independent fix, commit, push, issue closure/tick or stage-2 work.

Version **0.65.0-alpha**; changes remain **uncommitted** for parent acceptance.

## Month-end garrison tax forwarding (#116)

Pin remains `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`. The harness extracts
complete original definitions, including all assertions and their source lines:

| Source | Definition |
|---|---|
| `garr.c` | `garrison_gold` **446–549** |
| `build.c` | `province_subloc` **160–177** (new tracked source dependency) |
| `gm.c` | `int gold_provinces = 0;` **14** |

The requested 446–558 span crosses into the unrelated
`count_castle_garrisons` definition beginning at 552; the complete tax
settlement ends at 549 at this exact pin. Existing real `clear_temps`
**u.c:1352–1362**, `charge_maint_sup` **day.c:1325–1349**, strength, ownership,
inventory, actor-first stack payment, output and save-box closure are reused.
There is no successful maintenance/gameplay substitute. Only existing phase
instrumentation is omitted. Production `postMonth`'s position is unchanged;
its next explicit unfinished boundary is `ghost_warrior_decay` (#88).

The **23 new sections** cover:

- Tax **53**, maintenance **20**, own balance **7**, owner balance **13**:
  consume **33**, forward **16**, leave the original **7**, pay owner to **29**.
  A repeated zero-base call resets location temps and tax, retains forward
  **16**, and really starves **10 → 7** soldiers. Raw output bits and logs match.
- Full-tax spend and own funds beyond the tax base; friendly-stack money
  supplies **6** of **20** upkeep after the actor's **14**. The implementation
  measures spent against the actor's before/after gold, not stack totals.
- Paid strength **9/10**, plus real starvation **10 → 9** and **11 → 10**
  with tax **46** and tariff **5**: the **one** unspendable gold is consumed
  before strength testing; only the surviving ten restore it locally.
- Direct province castle and first-city castle restoration, full odd
  remainders and forward **-1**; province-child precedence, first-city-only
  search and no deeper recursion. Odd one-gold remainder truncates to zero
  without resetting an old forwarding value.
- Multiple garrisons produce **16 + 23 = 39** and **29** for two castles;
  actual owner payouts, faction routing, balances and signed-int GM sum
  **2147483500 → 2147483568** match. A shared province's tax is consumed by
  the first subtype-chain garrison, not duplicated for the second.
- Missing castle, absent misc, unowned castle, default garrison with stale
  fields, and no garrisons; temps clear even with no active garrison. A noble
  item with zero unit upkeep forwards **27** from tax **55**, then retains
  forward **27** on a zero-base call without inventing a maintenance payment.
- Paired real-tax/no-call controls preserve keyed income, weather, exit
  shuffle and mint bytes across same-turn master changes and the next turn.
- Nonprovince garrison aborts at **garr.c:471**; unlocated
  `province_subloc` reaches `subloc_here`'s **loc.c:338** `loop_here` assertion.
  Both require **134**, while every ordinary new section requires **0**.

Fixtures use real subtype-chain updates (`change_box_subkind`), not just
assignment to `skind`, and create the soldier box required by the original
strength macros. A temporary linker wrapper emitted an out-of-scope marker
and aborted on section **701**, which ordinarily requires status 134. Both
normal and `--update` modes rejected it with status **2** before assertion
acceptance or recording replacement; the recording hash was unchanged.
The probe was removed. No new verified upstream defect or blocker was found.

Two pinned C runs and Go match **706 sections / 52160 lines / 881837 bytes**.
Prior **683 sections / 50572 lines / 853708 bytes** of output and
**15712 lines / 281967 bytes** of input are exact byte prefixes. New input is
**16225 lines / 291764 bytes**. Output SHA-256:
`eca91b35e8bee3a1c75866ae378b803f2822f106542b86a8066fce565068c800`.
The C writer's trailing whitespace remains unmodified.

Sequential checks passed: command **706/52160**, move **135/7150**, helper
**37/1487**, RNG **321 results**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**,
default turn (identical raw status/stdout/stderr and one state file), followed
by bounded vet and fresh bounded default tests using shared `run_limited`.
Every compat call used `G3_REPO=$PWD/tmp/upstream.git`. No full-world run,
independent fix, commit, push, closure/tick or stage-2 work.

Version **0.66.0-alpha**; changes remain **uncommitted** for parent acceptance.

## Monthly body, ghost and mine decay (#88)

At pin `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`, extract the complete
original day.c definitions: `collapsed_mine_decay` **1525–1549**,
`ghost_warrior_decay` **1556–1584**, `corpse_decay` **1591–1619** and
`dead_body_rot` **1626–1652**. The task's broader 1525–1553, 1556–1588 and
1591–1623 bands include the comments after the first three definitions.
The real existing upkeep prefix contains `up_corpse` **138–143**; Go reuses
persistent `streams.UpCorpse`, not a reset or adapter. No new source file,
pin or affected.sh dependency is needed. Source-backed clarification:
[#88 comment](https://github.com/mdhender/g3g0/issues/88#issuecomment-6462926473).

The **17 new sections** compare exact saved boxes, kind/subtype chains,
inventories, faction output bits, raw logs, status and RNG aftermath:

- Empty passes; ghost quantities **0/1/3** (and negative raw control),
  loss of exactly one per call, full `is_npc` filtering for nonzero subtype,
  `LOY_npc` and `LOY_summon`, including regular-subtype NPCs.
- Seed-0 turn-9 corpse keys **1000 → 2**, **1001 → 1**, **1100 → 0**;
  a one-corpse stock clamps the two draw to one. Repeats keep the same keyed
  draws even after day/master changes; next-turn keys change independently.
  Zero/negative inventory and all three NPC branches skip decay.
- Body ages **11/12/13** and future timestamps; missing character record
  becomes an allocated empty `CH` even at age 11, then disappears at 12.
  Character owners receive the decomposition line; location owners do not.
  Adjacent current-box deletions retain the remaining subtype chain.
- Real unique creation/destruction refunds **7 NP** to the old lord and
  **2 NP** to the independent body's previous lord. Owner inventories and
  body boxes disappear; another decay call does not refund twice. Signed
  short clock/death wrapping occurs before promoted integer subtraction.
- Mine delays **0/1/2**, missing misc, **-128/-127/127**, plus a noncollapsed
  subtype control: **-128 → 127**, **-127 → 0 → delete**, **127 → 126**.
  Real `get_rid_of_collapsed_mine` evacuates nested characters, prisoners,
  loose items and locations, keeps the stack together and marks touched
  locations. A collapsed mine nested inside the current mine is relocated
  then deleted later in the same subtype traversal. Saved `here_list` order
  captures snapshot mutation, not just a mocked deletion result.
- Paired decay/no-call controls preserve sibling upkeep-income, weather,
  exit-shuffle and mint output, including same-turn master changes and the
  next turn. Mint IDs **66288** then **61234** match both controls.

The three malformed raw controls require **134**: owner **day.c:1638**,
missing owner inventory **u.c:2004**, and unlocated occupied mine reaching
`loc_depth` **u.c:388**. The last is **not** a `loop_here` assertion. Source
inspection confirms the actual mine-removal `loop_here` at **u.c:1083**
expands `valid_box(fort)` there; loop.h **242–252** does not change that
caller's file/line. loop_subkind **153–161** prefetches next before deletion.
These raw invalid-state controls alone establish no public-gameplay defect.
No newly verified C bug or blocker was found; none was independently fixed.

A temporary linker wrapper emitted an out-of-scope dependency marker in
section **714**, which ordinarily expects the owner assertion/status 134.
Both normal and `--update` modes rejected it with status **2** before
assertion/status acceptance or recording replacement; the recording hash
was unchanged. The temporary probe was removed, then the normal tool passed.

Two real C runs and Go match **723 sections / 53123 lines / 896467 bytes**.
Prior **706 sections / 52160 lines / 881837 bytes** of output and
**16225 lines / 291764 bytes** of input are exact byte prefixes. New input:
**16648 lines / 299004 bytes**. Output SHA-256:
`bff609ace017f709a32b9649bbe5e995bfd088e875f5f9df7a348f8add3657fc`.
Original C trailing whitespace is retained, not normalized.

Sequential checks passed: command **723/53123**, move **135/7150**, helper
**37/1487**, RNG **321 results**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**,
default turn (raw status/stdout/stderr and one state file identical), then
bounded vet and fresh bounded default tests. All engine runs used shared
`run_limited`; every compat invocation used `G3_REPO=$PWD/tmp/upstream.git`.
The normal command comparison also passed again after dependency probes.
No full-world turn, new tool, successful gameplay substitute, postMonth
call-position change, report-body/stage-2 work, commit, push, closure or tick.

Version **0.67.0-alpha**, left **uncommitted** for parent acceptance.

## Monthly autoquit, links, quests and ranks (#89)

At the unchanged pin `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`, extract
complete original definitions, retaining source `#line` attribution:

| source | definitions and inclusive lines |
|---|---|
| day.c | `auto_drop` **1133–1167**, `link_decay` **1966–1997**, `quest_decay` **1737–1753**, `post_month` **2147–2203** |
| garr.c | `top_ruler` **101–114**, `nprovs_to_rank` **588–600**, `find_kings` **627–671**, `lower_noble_rank` **680–698**, `det_noble_rank_sup` **700–715**, `determine_noble_ranks` **718–761** |

The accepted ownership, pledge, strength, queue, mint, output, rank-display
and save-box closure is reused. `rankland` only constructs dense fixtures:
one castle with its first direct character owner, numbered provinces and
garrisons. It neither calculates ranks nor supplies expected behavior.
No new source dependency, pin, RNG adapter or upkeep tool is needed.

Nine added sections exercise:

- Asymmetric regular factions **4/5/6** turns since orders against a
  five-turn autoquit, plus missing player record and NPC subtype. The missing
  regular record is allocated even before its age test. Existing faction
  orders precede appended `quit`; repeats append another one. Names/email,
  special log **11**, output bits and actual saved order files are compared.
- Links **0/1/2/-1**, missing destinations/records, matching/nonmatching
  months and ship exclusion. Positive countdown happens before reopening to
  two; negative stays permanent. Turn **8 → 9** wraps zero-based month
  **7 → 0**, and turn 16 returns to month 7. `subloc_quest` is **quest_late
  itself**, not a separate quest ID (`oly.h:1082`); only positive counters
  age, covering **0/1/2/-1/-128/127** and repeated calls.
- Every rank threshold, both helper and actual owner-count paths:
  **0/1/5/6/12/13/24/25/37/38/50/51/63/64**, plus negative/max-int helper
  controls. Strength **9/10**, subtype versus first-character traversal,
  temporary resets, missing character allocation, rank reset, pledge chains,
  ties and a zero-rank overlord distinguish recursive minimum from a cap or
  iterative decrement. No cycle protection is added.
- Kingdoms **14/15**, weak garrisons, nonlocation children, an ungoverned
  province and mixed versus shared top rulers. Kingship deliberately does
  **not** test strength: fifteen weak garrisons still crown their ruler.
- Two complete original-C `post_month` calls on a two-province, two-garrison,
  castle/stack world: owned-province forwarding **21 then 15 gold**, castle
  taxes **30 then 40**, maintenance at both original positions, NP/unformed
  minting, aura/shroud, production, autoquit append and disabled-autoquit
  control, links/quests, real token callback, rank settlement and final flag.
  Actual character description and stack output, logs, orders and serialized
  state match. Every callback is real, not a phase trace; empty branches do
  not claim coverage of every possible world. No new gameplay dependency was
  reached. Full-world integration and reports remain outside this task.

Loop audit: `loop_pl_regular`, `loop_char`, `loop_loc` and `loop_garrison`
use the accepted next-prefetch chains (`loop.h:122,153–163,187`). Ownership
reads the next pledge before the body (`loop.h:359–367`), with no independent
guard. `find_kings` preserves the `loop_here` snapshot and caller assertion
**garr.c:642 `valid_box(reg)`**, not a loop.h source location. Existing
`garrison_strength` inventory assertion remains **garr.c:60**.

Normal and `--update` linker probes temporarily replaced `drv_auto_drop`
with a missing-dependency marker and abort. Both failed **2** at section
**723**, before ordinary status acceptance or recording replacement; the
recording hash was unchanged. The probe was removed and normal comparison
rerun. No newly verified upstream bug or blocker was found.

Two real C recordings and Go match **732 sections / 55172 lines / 948492
bytes**. Prior **723 sections / 53123 lines / 896467 output bytes** and
**16648 lines / 299004 input bytes** are exact byte prefixes. New input:
**16979 lines / 305393 bytes**. Output SHA-256:
`72f0d38d1b3fa9c04245003e08940318d648c3c52ec546db38948028354207b8`.
Raw C trailing whitespace is retained; `git diff --check` flags only those
recorded bytes, not source edits.

Sequential checks passed: command **732/55172**, move **135/7150**, helper
**37/1487**, RNG **321 results**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**,
default turn (status, raw stdout/stderr and one state file identical), then
bounded vet and fresh bounded default tests. All engine processes used shared
`run_limited`; all eight tools used `G3_REPO=$PWD/tmp/upstream.git`.
Source-backed clarification: [#89 comment](https://github.com/mdhender/g3g0/issues/89#issuecomment-6463214748).

Version **0.68.0-alpha**; changes left **uncommitted**. The final four phase
NYIs are implemented, so `turn_phases.go` is deleted. `postMonth` order,
tokens and final flag positions are unchanged. No full-world run, stage-exit
acceptance, unrelated command, report-body/stage-2 work, commit, push, issue
closure or checklist tick.

## Output ordering and player entity state (#121)

Pinned C remains `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`. Extract complete
original definitions with their original `#line` locations:

| source | definitions and inclusive lines |
|---|---|
| report.c | `output_order_comp` **10–20**, `sort_for_output` **23–27**, `determine_output_order` **30–70**, `inv_item_comp` **213–220**, `extra_item_info` **246–290**, `show_unclaimed` **746–780**, `player_ent_info` **783–798** |
| perm.c | `admit_comp` **141–148**, `print_admit_sup` **151–187**, `print_admit` **190–223**, `print_att_sup` **426–466**, `print_att` **469–488** |

The C oracle calls actual `determine_output_order` and the entire actual
`player_ent_info` rendering closure; no expected-state calculation is
rewritten. Existing code.c, sout.c, location traversal, inventory/admission
lists and save_box supply their real dependencies. The driver now calls
real `init_spaces`, as main's initialization does: original `print_att_sup`
uses it when a rendered attitude line wraps. An initial development run
exposed this missing harness initialization, not an upstream failure.
No new source dependency/pin or affected.sh change is needed.

Go adds `determineOutputOrder`/`sortForOutput` and `playerEntInfoState`.
These are stage-1 database/diagnostic projections, not stage-2 report
rendering or successful gameplay substitutes. They are not wired into main
in this task. Rendering's `plural_item_box` warnings on stderr remain real
and observable, using the existing Go formatter after the inventory sort.
Faction log/report artifacts alone are outside this comparison; stdout,
stderr and status are compared without normalization. No character-report
or summary function is called on either side.

Eight small added sections cover:

- Opposing region/entity IDs; region traversal followed by depth-first
  here-list order through provinces, ships, buildings and stacked characters.
  Established player units sort in place, including silent factions, while
  missing player records stay missing. Repeated calls and relocation change
  only the appropriate scratch positions and persisted unit ordering.
- Unlocated equal positions use entity IDs, also when mixed with located
  units; unvisited boxes retain stale positions rather than being reset.
- Descending inventories, admissions and hostile/neutral/defend lists.
  Raw observers show retained invalid references and nonpositive quantities
  that save_box filters out; the actual saved `un`, `il`, `am`, `an`, `ad`,
  `ah` and `kn` bytes are also compared. Admission members retain their
  original order, and equal target keys gain no sense/member secondary key.
- Silent players skip all player_ent_info state, including absent-record
  allocation. Switching one to regular enables those effects. Missing
  regular/system/NPC/savage records allocate; missing silent records do not.
  Missing attitude/character/magic records stay missing. Character inventory
  and attitudes are unchanged; existing knowledge is unchanged.
- Repeats, all-invalid attitude lists, actual rendered admission/attitude
  continuation lines and ordered missing-plural warnings. The latter are
  emitted once per eligible inventory entry on every call, not suppressed.

The original loop macros retain their caller locations through extraction.
The new Go callers walk valid player/region chains, so their inventory/here
valid-box assertions and print_admit's player-kind assertion cannot fail
independently on these paths; no replacement assertion location is invented.

Negative linker probes temporarily wrapped only `player_ent_info` with an
explicit missing-dependency marker and abort. Normal and `--update` modes
both rejected section **734** with status **2**, specifically at the marker
check before expected-status acceptance or recording replacement. The
recording SHA-256 was unchanged after each probe. Removed the probe sources
and compiler wrapper, then reran the ordinary command comparison successfully.

Two real C recordings and Go match **740 sections / 55920 lines / 963859
bytes**. Prior **732 sections / 55172 lines / 948492 output bytes** and
**16979 lines / 305393 input bytes** remain exact byte prefixes. New input:
**17227 lines / 310599 bytes**. Output SHA-256:
`108b291e86341fcf5acdd4c855e14ce7a872487f3865f72ed203e2ce79c71f8a`.
Unfiltered C trailing whitespace is retained: full `git diff --check` flags
only recorded bytes; the check excluding that recording passes.

Sequential checks passed: command **740/55920**, move **135/7150**, helper
**37/1487**, RNG **321 results**, code **43/969**, loadsave **61/3411**,
turn-driver **7/152**, phase **5/484**, calendar **2/119**, schedule **4/1801**,
default turn (status, raw stdout/stderr and one state file identical), then
bounded vet and fresh bounded default tests. All engine runs use shared
`run_limited`; all eight tools used `G3_REPO=$PWD/tmp/upstream.git`.
Before broad checks, confirmed the parent's roundtrip shell **2896569** had
exited; its log reported **196 pass / 0 fail / 0 error**, not a run restarted
or claimed as this task's own verification.

Version **0.69.0-alpha**, left **uncommitted**. No new bug/blocker was found.
No main integration, character/summary state (#124), lore (#122), stage-2
rendering, whole-world/stage-exit acceptance, commit, push, issue closure or
checklist tick is included.

## Character report state and remaining allocations (#124)

Pin remains `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`, read only from
`tmp/upstream.git` via archive. Complete original definitions, with original
`#line` locations, extend this existing oracle; real `save_box` is unchanged:

| source | added definitions and inclusive lines |
|---|---|
| report.c | `show_carry_capacity` 73–105, `show_item_skills_sup` 108–193, `show_item_skills` 196–210, `show_char_inventory` 293–344, `char_rep_stack_info` 410–433, `collect_pledge_backlinks` 438–463, `show_pledged` 466–496, `prisoner_health` 499–510, `char_rep_prisoners` 513–537, `char_rep_health` 540–563, `char_rep_combat` 565–606, `char_rep_misc` 609–639, `char_rep_magic` 642–678, `char_rep_sup` 681–710, `character_report` 713–743, `turn_end_loc_reports` 1589–1642, `player_banner` 1645–1722 |
| use.c | `exp_s` 977–994, `rep_skill_comp` 1199–1228, `flat_skill_comp` 1231–1238, active `list_skill_sup` 1262–1271, `list_skills` 1276–1335, `fractional_skill_qualifier` 1345–1361, `list_partial_skills` 1364–1406 |
| summary.c | `rank_comp` 15–22, `ranking` 33–41, `top_rank` 44–69, `out_rank_mine_top` 72–89, `out_ranking` 92–119, `summary_provinces` 348–389; original SHOW_TOP/ranks declarations |
| buy.c / lore.c / u.c | `list_pending_trades` 1123–1158 / `np_req_s` 120–134 / `ordinal` 1711–1724 |

The C character renderer and its full closure really execute: no copied
expected-state formulas or successful dependency stubs. The extractor selects
the active `#else` definition of `list_skill_sup`, not the disabled `#if 0`
alternative. The location-allocation probes use empty known-location sets;
unported `show_loc`/`show_loc_header` dependencies **abort with the rejected
dependency marker** if reached, never return success. These probes do not
claim complete location rendering. `player_banner` itself executes whole.
`summary_provinces` includes its real ranking/rendering dependencies.
**summary.c is a new dependency**, tracked in tools/README.md and affected.sh.

Go extends report.go, reusing #121's `reportAttitudeState`. Only persistent
state and observable diagnostics/assertions are projected; report log bodies,
output bits and scratch temp/output_order/ranking data are stage-2 artifacts
or scratch, not substituted saved state. In particular:

- `character_report` skips silent owners first. Prisoners only temporarily
  toggle their prisoner flag for output, so stay prisoners and do not sort.
  That initial output can still allocate a missing nonsilent owner player
  record via sout.c:306; it does not trigger pledge collection.
- Only the first eligible nonprisoner reaches `show_pledged`. Its two passes
  first `p_magic`/clear **every** character (including silent/prisoners), then
  append every nonzero pledge in character-chain order across factions.
  Missing character records gain CH; empty magic alone emits no CM. A
  non-null backlink pointer, even after clearing to zero entries, retains
  bare CM because the real writer zero-checks the entire structure.
- Duplicate stale backlink entries disappear on collection, but repeated
  reports with the ready flag set do not rebuild. Pledge edits alone leave
  stale backlinks until that flag is cleared. Allocated-empty lists stay
  allocated across clear/repeat; never-allocated lists remain nil.
- Eligible inventory and hostile/neutral/defend lists sort **in place**,
  retaining invalid references and nonpositive entries filtered by writers.
  Both skill displays sort copies and leave all stored skill rows in order.
  Plural-name stderr warnings and prisoner/partial-skill source assertions
  remain observable; report.c:504 and use.c:1358 failures record status 134.

Allocation caller audit against actual main.c:981–1001 order:

1. determine_output_order allocates no absent player record (#121).
2. turn_end_loc_reports allocates **every nonsilent** player, before order
   templates and player_ent_info. Templates' conditional p_player call
   (order.c:398), print_admit and output routing therefore add no earlier
   allocation. The new projection is intentionally not moved to summary.
3. character_report has the lazy CH/CM effects above. Its allocation-only
   owner route is also tested independently of the preceding player phase.
4. player_banner's p_player is nonsilent, after character_report and before
   summary. It is tested standalone as well as in the original composition.
5. summary_provinces calls p_player for **ALL subtypes**, including silent
   factions. Other summary mutations are scratch. Subsequent player_report
   (player_report_sup/unit_summary), GM knowledge statistics and include
   generation encounter records already established by this all-player pass;
   they introduce no additional absent-record allocation in the actual tail.
   Skill lore/knowledge remains #122, not this task.

Eleven appended small sections cover all-prisoner/all-silent controls before
first eligibility; missing CH/magic/player records; lazy cross-faction pledge
targets and sources (ordinary, silent, prisoner), duplicate/clear/repeat and
nil versus allocated-empty lists; all three original attitudes and descending
inventories including invalid entries; unchanged descending known/partial
skills; all player subtypes; original allocation-phase composition; and real
assertion failures with healthy/immortal prisoner controls. No new default
world test or private formatting-only test was added.

Negative linker probes wrapped only `character_report` with a dependency
marker, once returning status 0 and once aborting with 134. Normal and
`--update` runs both rejected section **740** with harness status **2** at the
marker check, before expected-status or recording acceptance. Recording hash
stayed unchanged. Removed both temporary probe files, then reran the real
ordinary command comparison successfully.

Two C recordings and Go match **751 sections / 56776 lines / 979131 bytes**.
The prior **740 sections / 55920 lines / 963859 output bytes** and **17227
lines / 310599 input bytes** remain exact byte prefixes. New input is **17487
lines / 316195 bytes**. Output SHA-256:
`2c9911876370b789a97e4a41ce550ce12a55ef5f1806ce9ecfa9ecfe11448d88`.

Eight tools passed **sequentially**, all with `G3_REPO=$PWD/tmp/upstream.git`:
command **751/56776**, move **135/7150**, helper **37/1487**, RNG **321**,
code **43/969**, loadsave **61/3411**, turn-driver **7/152**, phase **5/484**,
calendar **2/119**, schedule **4/1801**, then default turn (identical status,
raw stdout/stderr and one state file). Next, bounded `go vet ./...` and fresh
bounded `go test -count=1 ./...` passed. All engines use shared run_limited;
no whole-world, roundtrip or stage-exit run is claimed for this task.
Source-only `git diff --check` passes; literal C recording whitespace remains.

Version **0.70.0-alpha**; all changes left **uncommitted**. No new defect or
upstream blocker found. No main integration, skill lore, full Go report
renderer, issue closure/checklist tick, commit or push is included.

## Skill lore state and input diagnostics (#122)

The extraction now includes every complete function in lore.c:10–387:
`lore_function`, `deliver_lore_sheet`, `np_req_s`, `out_skill_line`,
`deliver_skill_lore`, `deliver_lore`, `queue_lore`, `lore_comp`,
`show_lore_sheets`, `gm_show_all_skills` and `scan_char_skill_lore`.
It keeps the accepted `scan_char_item_lore`, real `test_known`/`set_known`,
`skill_school`, output library and actual `save_box` writers; u.c's whole
`weeks` definition completes the small rendering closure. No compared C
function is altered, no new source dependency/pin or generic IO change.

Go extends **the existing lore.go** with state/input projections only:

- Skill scanning uses every copied skill row, not the known-only iterator.
  It includes silent/prisoner characters, shares faction knowledge, preserves
  stored skill order and delegates all assertions and anyway handling to the
  original `queueLore`. No new RNG calls or independent gameplay fixes.
- Show-lore sorts the original queue with the C int comparator. Adjacent
  duplicates suppress only delivery, not queue entries. Repeated calls read
  again and never clear the queue or change the saved known-bit order.
- GM lore retains initial player allocation and listing `skill_school`
  assertions before any sheet read, then schools followed by nonschools in
  skill-chain order, without granting knowledge or queuing lore.
- Each sheet is actually opened/read using existing `clib.ReadFile` and
  `getlin` semantics. Linux directory fopen succeeds; its subsequent read
  failure yields no missing-file diagnostic. ENOENT and ENOTDIR remain
  nonfatal exact perror diagnostics. Unknown `$` names (including empty or
  whitespace-suffixed names) print stderr then assert at **lore.c:43**.
  Only `capturable_animals` and `animal_fighters` are allowed. Comment, @ and
  legal expansion report text remains stage 2, not a file-existence shortcut.

Nineteen appended sections cover all skill states (including an unexpected
state that scanning does not inspect), shared/different factions, prisoners,
silent players, missing records, unowned/invalid skill assertions, known and
queue order, anyway duplicates, sorted retention/repeat reads, missing and
existing sheets, both legal functions with real animals, comments and @ runs,
item-to-skill/etc remapping, GM read ordering/allocation, directory/read-error
and ENOTDIR controls, modified-sheet rereads, unknown/empty/whitespace function
assertions, and skill-cycle assertions before reads. Real C executes all
rendering closure, but report artifact text is not a Go acceptance target.

The fixture-only `lore_root` action changes to each isolated section directory
and sets libdir to `.` on both engines, preserving **raw** diagnostic paths
without filtering. Go sections run as sequential subtests so cwd cleanup is
complete before reading the recording or starting the next section. The
`lore_dir`/`lore_file` setup actions create actual files; `|` encodes newline,
with no implied final newline. C runs twice, comparing raw streams/status and
actual saved state before any recording update. New lore sections have exact
expected statuses; the existing dependency-marker rejection runs first.

Validation at **0.71.0-alpha**, with all changes left **uncommitted**:

- Two actual pinned C runs and Go match **770 sections / 57313 lines /
  989824 output bytes**. The prior **751 sections / 56776 lines / 979131
  output bytes** and **17487 input lines / 316195 input bytes** remain exact
  byte prefixes. New input: **17728 lines / 321249 bytes**. Output SHA-256:
  `08da580438e0ae135eac7766106e2db7632cf20147c41cd3506cf3dc058899fc`.
- Negative linker probes wrapped only `gm_show_all_skills`, emitting the
  rejected dependency marker and either returning 0 or aborting 134. Normal
  and `--update` runs rejected section **754**, harness status **2**, before
  expected-status/update acceptance. Recording hash was unchanged. Both
  temporary probe source/wrapper files were deleted before the real rerun.
- Eight tools passed **sequentially**, `G3_REPO=$PWD/tmp/upstream.git`:
  command **770/57313**, move **135/7150**, helper **37/1487**, RNG **321**,
  code **43/969**, loadsave **61/3411**, turn-driver **7/152**, phase **5/484**,
  calendar **2/119**, schedule **4/1801**, then default turn (identical status,
  raw stdout/stderr and one state file).
- Bounded `go vet ./...` and fresh bounded `go test -count=1 ./...` passed;
  Olympia defaults took **1.635s**. Engine binaries used shared run_limited.
  Source-only `git diff --check` passes; raw C recording whitespace is kept.

No new defect/blocker found. No main integration, Go report text renderer,
whole-world/stage-exit acceptance, issue closure/tick, commit or push.

## Residual report-tail state (#125)

The oracle extracts unchanged complete `orders_template_sup`, `orders_other`,
`orders_template`, `list_order_templates`, `report_account_sup`/`report_account`,
all active `summary_report` helpers, the complete `player_report` closure,
`gm_report` and its gm.c helpers, and `gen_include_sup`/`gen_include_section`
with `rep_player`, `rep_char`, `rep_loc` and `inc`. Existing real weights,
inventory, ownership, skill, output and save writers are reused. Original
report declarations/tables and GM gold scalars also come from the pin.

Go extends report.go's **state projections**, not a report renderer:

- Templates filter system/silent players and allocate output recipients;
  autoquit's explicit player access remains separate. Existing commands and
  own/foreign order queues are unchanged, not reconstructed from text.
- Account reports stage, reset style and allocate only regular players.
  The account-summary command/file input is inside `#if 0`; this does not
  implement external `charge_account`.
- Summary reports open regular output before `summary_provinces` allocates
  every player. `loop_known` sorts the original known lists. Rankings and
  skill/spell/men/gold/land statistics are deferred output-only scratch.
- Player reports allocate nonsilent records through `player_report_sup` and
  `unit_summary`, including system players. Unit, location-stack, weight,
  storm and garrison text/scratch/render-only validation remain stage 2.
- GM reports stage and open GM output, allocate every player in the known
  gate/location passes (sorting original knowledge), then open skill-player
  output last. Pure statistics and item/notice bodies remain stage 2.
- Include generation validates **existing** output bits, allocates nonsilent
  player records, and allocates corpse misc records through `rep_char` after
  its initial output. Foreign players, ordinary items and other unsupported
  entity kinds retain report.c's original source-located assertions. Silent
  player inputs are not inspected. No deferred artifact memberships are made
  up. Include text and location headers remain stage 2.

Shared `openPlayerOutput` supplies ordinary sout.go output and report state
effects, preserving real raw open-failure diagnostics before output allocates
a player record. It emits no placeholder report lines. Transient report bits
are deferred as in #122. `show_loc_header`, `show_loc_posts` and `show_loc`
remain **aborting rejected-marker dependencies** on deliberately unexercised
paths, not successful renderer substitutes. No production main call is wired.

Nineteen appended cases compare actual C save bytes, raw list/record observers,
output-file presence, preserved order queues and raw stdout/stderr/status.
They cover all player subtypes/missing records, repeated calls, asymmetric
knowledge, nested stacks, prisoners, capacity/storm traversals, gate knowledge,
corpse allocation, recognized special/stale include inputs, exact include
assertions, six report-family output-open failures, account style through an
ordinary real output probe, and the original relative order of the new
projections on a small world. Pure report artifacts are not read/compared.
No new upstream source dependency or pin, renderer, admission handling,
file-lifecycle/main integration or full-world/stage-exit acceptance is implied.

Validation at **0.72.0-alpha**, left **uncommitted**:

- Pinned C twice and Go match **789 sections / 58175 lines / 1005735 output
  bytes**. The previous **770 sections / 57313 lines / 989824 output bytes**
  and **17728 input lines / 321249 input bytes** are exact unchanged prefixes.
  New input: **18083 lines / 328103 bytes**. Output SHA-256:
  `68f9453a5ebd6373f4d3932fd9965f3a4c24dba8f2386d2113da063742021713`.
- Negative linker probes wrapped only `report_account` at an already
  expected-abort case, emitted the missing-dependency marker, and either
  returned 0 or aborted 134. Both ordinary and `--update` runs rejected
  **section 783** with harness status **2**, before status/update acceptance.
  Recording hash remained unchanged in all four runs. Temporary wrapper
  source/compiler scripts were deleted before acceptance checks.
- Eight tools passed **sequentially**, `G3_REPO=$PWD/tmp/upstream.git`:
  command **789/58175**, move **135/7150**, helper **37/1487**, RNG **321**,
  code **43/969**, loadsave **61/3411**, turn-driver **7/152**, phase **5/484**,
  calendar **2/119**, schedule **4/1801**, then default turn (identical status,
  raw stdout/stderr and one state file).
- Bounded `go vet ./...` and fresh bounded `go test -count=1 ./...` passed;
  Olympia defaults took **1.709s**. Engine binaries used shared run_limited.

No new stage-blocking defect or additional report input/state closure found.
Main's explicit report-tail stop remains. No issue closure/tick, commit, push,
whole-world check or stage-exit claim.

## Cwd-relative admission probing (#127)

The oracle now extracts complete original add.c `fetch_inp`,
`make_new_players_sup`, `make_new_players`, `add_new_players`,
`new_player_list` and their actual admission/report-loop dependencies,
including `get_city_id`, `pick_starting_city` and `add_new_player`.
These definitions are unchanged and keep source `#line` locations. Static
wrappers expose the parser/prober; no successful admission or rendering stub
is introduced. The existing rejected-marker `show_loc` dependency remains
unexercised: complete admissions and their initial reports are stage 6.
`add.c` is newly tracked in tools/README.md and affected.sh.

Nineteen appended sections compare missing/empty/file `act`, dot-prefix
skips, ENOENT/ENOTDIR Join opens, successful directory opens followed by
read errors, every partial-field count, repeated changed inputs, cwd versus
contradictory libdir input, and complete invalid account assertions at
**add.c:325**. Independent fetch controls cover five distinct valid fields,
extra fields, blank skipping, tabs/control/high bytes, embedded NUL, CR and
final nonnewline input. Actual `save_box` bytes, absent player records,
new-player/vector lists and raw stdout/stderr/status are recorded. The
no-admission tail still calls `new_player_list`, including its real player
vector update; it neither admits anyone nor allocates absent player records.

Go stops with `notImplemented("new player admission")` after all five fields
and scode/range validation, immediately before C's **add.c:327 alloc_box**.
Twelve separate boundary controls cover direct and filesystem entry, numeric,
encoded and bracketed accounts, the highest valid ID and already occupied
IDs. They require unchanged saved entity bytes and new-player/vector lists,
and direct calls must leave the sixth field unread. They are explicit NYI
tests, NOT expected-success C admission comparisons.

Validation at **0.73.0-alpha**, left **uncommitted**:

- Real pinned C twice and Go match **808 sections / 58660 lines / 1015265
  output bytes**. Input: **18244 lines / 332405 bytes**. Prior **789 sections,
  18083 input lines / 328103 input bytes**, and **58175 output lines /
  1005735 output bytes** are exact unchanged prefixes, checked by SHA-256.
  New output SHA-256:
  `fecb1996344d054fb4ec41df5f0df7f58889bfd5b8083ff3f6fe4e35e363e57a`.
- Eight tools passed sequentially with `G3_REPO=$PWD/tmp/upstream.git`:
  command **808/58660**, move **135/7150**, helper **37/1487**, RNG **321**,
  code **43/969**, loadsave **61/3411**, turn-driver **7/152**, phase **5/484**,
  calendar **2/119**, schedule **4/1801**, then default turn (identical status,
  raw stdout/stderr and one state file). All engine invocations use shared
  run_limited; the missing-dependency rejection still precedes status/update.
- Bounded vet and fresh bounded defaults passed; Olympia defaults **1.641s**.

No new bug found. Post writers (#128), cached ownership/Times (#129),
stage-2 rendering and main integration remain untouched, including main's
explicit postMonth stop. No whole-world/stage-exit success is claimed.

## Post-turn writer opens and state (#128)

Pin remains `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`, read only through
the bare `tmp/upstream.git` archive. The existing extractor adds complete
unchanged main.c:50–290 definitions: `write_totimes`, `write_email`,
`list_a_player`, `write_player_list`, `write_forward_sup`, `write_forwards`,
`write_faction_sup` and `write_factions`. Original ownership/garrison helpers,
code.c, sout.c, IO and save_box writers supply actual dependencies; there is
no replacement expected-state algorithm or successful unported dependency.
main.c is newly tracked in tools/README.md and affected.sh.
The issue's 55–314 reference is approximate: complete definitions at this
pin occupy 50–290; 314 is inside unrelated send_rep, which is not extracted.

Go extends the existing olympia.go owner with explicit stage-1 projections.
Every writer actually opens with the existing clib fopen("w") and closes on
normal completion. Only artifact bodies are deferred, not file decisions,
write_player_list's stage or the exact nonfatal can't-write/perror streams.
The original order is players, email, totimes, forward, factions. Ownership
and garrison validation occur only after forward opens successfully. There
is no suppression option or main call-site change.

Allocation audit: list_a_player calls p_player only after its caller has
already required a nonnull player record and email. Thus these writers do
not allocate missing player records on valid input. Forward/faction helpers
use raw player_email accesses after player(); garrison traversal uses real
players_who_rule_here. Save bytes plus explicit record flags prove missing
records remain missing, including a regular player with no record, rather
than introducing an unconditional allocation while deferring artifact text.

Fourteen appended sections cover asymmetric identities, visible/CIS mail,
new/silent/system players, missing and empty records, repeated calls and
individual writers, garrison castle/pledge traversal, all-five ENOENT and
ENOTDIR failures, each individual EISDIR failure with the other four opens
successful, and paired failed/successful forward opens with ownership-cycle
and invalid-castle validation. Only successful opens reach swear.c:32 and
loc.c:395 assertions (status 134); failed opens remain nonfatal (status 0).
Actual file presence/type and raw stdout/stderr/status are compared along
with real saved records. Report artifact bytes are deliberately not compared.

Validation at **0.74.0-alpha**, left **uncommitted**:

- Two actual C recordings and Go match **822 sections / 59134 lines /
  1024033 output bytes**. All prior **808 sections / 58660 output lines /
  1015265 output bytes** and **18244 input lines / 332405 input bytes** are
  exact unchanged byte prefixes. New input: **18422 lines / 336093 bytes**.
  Output SHA-256:
  `a15056a752433ffcf5bc43a074bb5edd9ffb560c2aade1ea1581a3e3fb095de9`.
- Temporary negative linker probes wrapped only write_forwards at the
  expected-abort ownership-cycle section, emitting the rejected dependency
  marker and either returning 0 or aborting 134. Both ordinary and --update
  runs rejected **section 818** with harness status **2**, before expected
  status or recording replacement. The recording hash was unchanged in all
  four runs. Probe source/compiler wrapper were deleted before the real run.
- Eight tools passed **sequentially**, with `G3_REPO=$PWD/tmp/upstream.git`:
  command **822/59134**, move **135/7150**, helper **37/1487**, RNG **321**,
  code **43/969**, loadsave **61/3411**, turn-driver **7/152**, phase **5/484**,
  calendar **2/119**, schedule **4/1801**, then default turn (identical status,
  raw stdout/stderr and one state file). All engines use shared run_limited.
- Bounded vet and fresh bounded default tests passed; Olympia **1.752s**.
  Source-only diff whitespace check passes; the added literal C recording
  retains one trailing-space line from the unchanged save_box writer.

All open local issues were inspected, including unlabeled issues: none is
bug/follow-c work. No new verified bug/blocker was found. Existing intentional
source-linked C defects are untouched. Main's explicit postMonth stop, cached
output/Times lifecycle (#129), valid admissions/-A, stage-2 rendering and
whole-world/stage-exit acceptance remain outside this task. No commit, push,
issue closure or checklist tick.

## Cached output ownership and Times lifecycle (#129)

Pin remains `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`; C sources are read
only from the bare `tmp/upstream.git` archive under `tmp/command-compat/g3`.
The original whole sout.c provides open_fp/alloc_fp/grab_fp (131–214), both
log openers (217–238) and close_logfile (241–264). Original whole io.c
provides save_logdir (3031–3069) and real remove_dir. The extractor adds
complete unchanged c2.c open_times (338–377) and close_times (457–472),
with the original DASH_LINE and static FILE declarations (326–331).
The issue's broader IO range reaches into save_db, and its Times ranges
include the next function's declaration; only the complete owners above
are newly ported, not main's final save_db integration.
All called definitions and original save_box writers remain real; no
successful rendering/gameplay dependency substitutes or rewritten expected
logic. sout.c, io.c and c2.c were already tracked in tools/README.md and
affected.sh, so no new upstream source dependency or pin change is needed.

Go extends the existing openPlayerOutput owner with a most-recent-first
bounded cache (FOPEN_MAX-10: six on Linux, ten on Darwin), not a second cache
or output abstraction. Hits borrow the existing stream, full-cache misses
close/flush the tail before opening, and closeLogfile resets all ownership
unless immediate. Ordinary and report-projection callers no longer close
borrowed streams. Append mode and buffering reuse clib.File's registry and
FlushAll; -f selects unbuffered output at open time, not retroactively on hits.
Existing command/move fixtures retain their C drivers' flush_always=1;
explicit fflush observers still flush. closeTimes closes/resets press then
rumor. openTimes uses the same buffered FILE owner with the real explicit
header flush before the rumor open. Fatal diagnostics have no perror/newline.

saveLogdir retains ignored mkdir/remove failures, both ordered renames and
their independent nonfatal diagnostics, and final log mkdir. It uses the
rename syscall, not os.Rename's destination-directory EEXIST precheck: the
actual C empty-directory success and nonempty-directory ENOTEMPTY must be
preserved. The Windows flag changes only initial backslash paths and prints
its deletion notice; both renames and final mkdir still use forward slashes.
The controls test this branch on Linux, not native Windows syscall semantics.

Thirty-six appended sections cover renamed/blocked cached hits, first-miss
ENOENT/ENOTDIR/EISDIR and close/eviction reopen failures, capacity/reset and
MRU eviction, buffered versus -f opens, immediate close/open bypass, shared
ordinary/report ownership in both directions, nondestructive/destructive
openers, successful/failed mkdir/renames, existing turn directories, save/no-
save path lifecycles and Windows controls. Real C save_box bytes and raw
stdout/stderr/status remain compared. Artifact text is not compared; file
type/presence and empty/nonempty probes establish ownership/flush decisions.
Times repeated opens survive renamed paths, then fail only after close/reset.
An atexit observer runs before exit's stdio cleanup and proves press headers
have already flushed when rumor fopen fails; press failure never opens rumor.
Main's postMonth stop and report-body/admission/-A boundaries are untouched.

Validation at **0.75.0-alpha**, left **uncommitted**:

- Actual pinned C repeated twice deterministically and Go match **858
  sections / 59875 output lines / 1038904 output bytes**. All previous **822
  sections / 59134 output lines / 1024033 output bytes** and **18422 input
  lines / 336093 input bytes** remain exact unchanged byte prefixes.
  New inputs total **18862 lines / 344535 bytes**. Output SHA-256:
  `626d440d357c30bb97420c5ee294a5865ed34668aca2080c5f31f7186a33b4cc`.
- Temporary negative linker probes wrapped open_times at fatal section
  **848**, emitted the missing-dependency marker and exited either **0** or
  **1**. Ordinary and --update runs all rejected with harness status **2**
  before expected-status checks or recording replacement; the output hash
  stayed unchanged. Probe source/compiler wrapper were removed before the
  final real-C run. Logs: `tmp/issue129-negative-*.log`.
- Eight tools passed **sequentially** with `G3_REPO=$PWD/tmp/upstream.git`:
  command **858/59875**, move **135/7150**, helper **37/1487**, RNG **321**,
  code **43/969**, loadsave **61/3411**, turn-driver **7/152**, phase **5/484**,
  calendar **2/119**, schedule **4/1801**, then default turn (identical status,
  raw stdout/stderr and one state file). All engines use shared run_limited.
- Bounded vet and fresh bounded default tests passed; Olympia **1.705s**.
  Full sequential log: `tmp/issue129-validation.log`. Source-only diff
  whitespace check passes. The literal C recording retains one trailing-
  space raw stdout line from immediate routing, not normalized text. New
  buffered-flush probes inspect file state without printing artifact text.

All six open local issues were inspected, including unlabeled issues; none
is bug/follow-c work. No new verified port/C bug or upstream blocker was
found. Existing source-linked C bug reproductions are unchanged. No native
Windows or Darwin execution, main wiring, final save_db/-n composition,
valid admissions/-A, stage-2 renderer, whole-world or stage-exit acceptance
is claimed. No commit, push, issue closure or checklist tick.
