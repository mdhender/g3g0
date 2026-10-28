# Compat tools

Each directory here holds a `run.sh` that checks part of the Go port against G3's C code at a pinned revision, and a README that says what it checks.
This page says which checks to run when, and which upstream files each tool reads, so that moving the pin re-records only the tools whose sources changed.

Task notes below retain their original acceptance scope and temporary boundaries;
the current integration evidence and stage-exit status are under **Stage exit**.

## Complete order templates (#152)

Existing order.go implements complete pinned order.c orders_template_sup
**370–454**, orders_other **457–505**, orders_template **508–547** and
list_order_templates **550–570**. Main/oracle replace listOrderTemplatesState
with the renderer; the copied report.go projection is removed. Real loops,
saved commands, order queues, names/loyalty and sout.go's cached FILE owner
are reused, not modeled or duplicated. A trivial wrapper exposes unchanged
io.c's static read_boxes to the existing report fixture adapter; existing
Go readBoxes and lore_file supply saved PL/CO inputs. No command handler,
new upstream dependency/pin or later report root.

At **0.92.0-alpha**, eleven appended sections match complete unchanged real
C twice and Go: **1072 sections / 72518 lines / 1388090 output bytes**.
Inherited **1061 sections / 71543 lines / 1368616 output bytes** and **432470
input bytes** remain exact prefixes. Actual raw logs and records cover absent,
empty and present passwords, broken mailer, autoquit half thresholds,
contract **49/50**, fear **1/2**, saved RUN **-1/0/1/12001**, LOAD/DONE,
foreign empty/queued orders, dead/invalid/prisoner filtering, all subtypes,
immediate allocation, route/frame resets and renamed cached inode repeats.
Numeric constants and reached branches were verified against the pin and
recording. Four temporary marker probes (0/134, normal/update) reject before
expected-failure acceptance/update, retaining the recording hash; removed.
See [complete closure, hashes and evidence](command-compat#complete-order-templates-152).

Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded defaults pass with explicit isolated G3_REPO
and run_limited supervision. Logs: `tmp/issue152-*.log`. Source-only whitespace
check passes; literal C recording whitespace is retained. Independent parent
review and the same tool/vet/fresh-test tier pass; logs `tmp/parent152-*.log`.
No stage-exit claim. Fresh post-#152 strict evidence remains deterministic
C repeats, C/Go **0**, equal streams/paths/types/modes, eight report-byte
differences, exit **1** before lifecycle; log `tmp/parent152-strict.log`.
G3 #358 stays nonblocking and
unchanged; division/Go panic with malformed **C136/Go2** is NOT passing parity.

## Complete character-report root (#151)

Existing report.go replaces characterReportState with complete pinned
report.c char_rep_sup **681–710** / character_report **713–743**, wired once
in main; copied projection logic and reportAttitudeState are removed.
Existing buy.go adds already-extracted buy.c list_pending_trades
**1123–1158**. All accepted detail/inventory/capacity/permissions/skills and
cached output owners are reused in actual C order, not callback traces or
parallel state/render paths. No new dependency/pin, later root or command.

At **0.91.0-alpha**, eleven appended sections match complete unchanged real
C twice and Go: **1061 sections / 71543 lines / 1368616 output bytes**.
Accepted **1050 sections / 70245 lines / 1339950 output bytes** and **426753
input bytes** remain exact prefixes. Full asymmetric faction report bytes,
silent/prisoner-only lazy controls, temporary prisoner toggle, missing owner
records/CH, retained CM, repeat backlinks, garrison routing, stored trades
and renamed cached inode are actually observed. Numeric fixtures and reached
branches were checked against pinned C. #168 native traversal still matches
over unchanged input. Four temporary marker probes (0/134, ordinary/update)
reject before expected-failure acceptance or recording update; probes removed.
See [closure, hashes and evidence](command-compat#complete-character-report-root-151).

Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded defaults pass with explicit isolated G3_REPO
and run_limited engine supervision; logs `tmp/issue151-*.log`. Independent
parent acceptance passes the same tier. The fresh post-#151 strict run has identical C
repeats, C/Go 0, equal streams and paths/types/modes, eight report-byte
differences, exit 1 before lifecycle. No full-world/stage-exit acceptance;
logs `tmp/parent151-strict.log`. G3 #358 remains nonblocking; original
division/Go panic unchanged, malformed C136/Go2 NOT passing parity.

## Character detail sections (#150)

Existing report.go implements complete report.c stack **410–433**, pledged
**466–496**, prisoner health/list **499–537** and health/combat/misc/magic
**540–678**, reusing location **357–407**, the existing two-pass backlink
collector **438–463** and lazy world flag, names, topology, hidden/skill,
inventory/aura and accepted output/cache owners. New direct renderers leave
the character root's projection and original collection timing intact until
#151; no duplicate effect pass, new dependency/pin, command or later root.

At **0.90.0-alpha**, fifteen appended sections match complete actual C twice
and Go: **1050 sections / 70245 lines / 1339950 output bytes**. The accepted
**1035-section / 69581-line / 1324865-output-byte** and **420081-input-byte**
prefixes remain exact. Cases actually reach health/recovery, noble versus
character stats, exact behind/break spacing, descending stack/prisoner order,
lazy all-character allocations and retained bare CM, every concealment phrase,
shroud/vision, aura bonuses/auraculum, project target kind filters and signed
quick/hide-mage boundaries. #168 live native traversal still matches unchanged
shared input. Four temporary missing-marker controls (0/134, ordinary/update)
reject before expected-failure acceptance or recording replacement; probes
removed. See [closure, hashes and evidence](command-compat#character-detail-sections-150).

Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded defaults pass with explicit isolated `G3_REPO`
and engine run_limited supervision. No new source dependency or pin.

The latest full strict evidence remains the parent's post-#145 diagnostic:
deterministic C repeats, equal raw streams and complete paths/types/modes,
eight saved-report byte differences, exit 1 before lifecycle. This task does
not refresh or accept that gate; no full-world/stage-exit claim. G3 #358 stays
nonblocking and unchanged; malformed C136/Go2 is not accepted as parity.

## Known and partial skills (#149)

Existing use.go adds complete pinned use.c exp_s **977–994**, comparators
**1199–1238**, active list_skill_sup **1262–1271**, list_skills **1276–1335**
and fractional qualifier/partial list **1345–1406**; lore.go adds complete
np_req_s **120–134**. Existing skill/magic/name/output/cache owners supply
the closure. Both sorts copy pointer lists, preserving raw/saved skill order.
The later character root remains a projection until #151, not a second pass.
No new dependency/pin, command or later root.

At **0.89.0-alpha**, eighteen appended cases match complete actual C twice
and Go: **1035 sections / 69581 lines / 1324865 output bytes**. The inherited
**1017-section / 68918-line / 1310660-output-byte** and **414715-input-byte**
prefixes remain exact. Numeric constants and reached branches include all
know states, asymmetric school grouping, nested prerequisites, no-exp short
circuiting, all experience bands/boundaries, NP text and malformed assertion
order. #168 live native traversal still matches. Four temporary marker
controls (exit 0/134, ordinary/update) reject before expected-failure acceptance
or recording replacement, preserving the recording hash; probes removed.
Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded defaults pass with explicit isolated `G3_REPO`.
See [complete closure, hashes, branches and evidence](command-compat#known-and-partial-skills-149)
and `tmp/issue149-*.log`.

No full-world/full-corpus/strict-lifecycle or stage-exit acceptance. Latest
strict post-#145 evidence remains identical C repeats, C/Go 0, matching
streams and paths/types/modes, eight saved-report byte differences, exit 1
before lifecycle. G3 #358 remains nonblocking: original division/Go panic
untouched; malformed C136/Go2 is NOT passing parity.

## Admission permissions, attitudes and player entity reports (#148)

Existing perm.go adds complete pinned perm.c admit comparator/renderers
**141–223** and attitude renderers **426–488** (complete closing definition);
report.go/main replace playerEntInfoState with real player_ent_info
**783–798**, reusing accepted showUnclaimed and original cached output.
The still-deferred character root retains only reportAttitudeState until
#151. Original in-place admission/inventory/three attitude sorts, retained
invalid references, unsorted admission members, allocation/output order and
silent skip are unchanged. No new dependency/pin, commands, valid new-player
admissions or later roots.

At **0.88.0-alpha**, eleven appended sections match unchanged C twice and Go:
**1017 sections / 68918 lines / 1310660 output bytes**. The inherited
**1006-section / 68035-line / 1290206-output-byte** and **408111-input-byte**
prefixes remain exact. Numeric constants and actual branch reach include
equal targets, duplicates, all player subtypes and 11/12/13 continuations;
#168 live native traversal still matches. Four temporary marker controls
(exit 0/134, normal/update) reject before expected-failure acceptance or
recording replacement, preserving the recording hash. Probe files removed.
Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded defaults pass with explicit isolated
`G3_REPO` throughout. See [closure, hashes, evidence and caveats](command-compat#admission-permissions-attitudes-and-player-entity-reports-148)
and `tmp/issue148-*.log`. Independent parent review and the same acceptance
tier pass; parent logs are `tmp/parent148-*.log`.

No new strict/full-world/stage-exit acceptance: latest strict post-#145 C
repeats match, C/Go 0, streams and entry paths/types/modes match, eight report
byte differences exit 1 before lifecycle. G3 #358 remains nonblocking;
original division/Go panic untouched, malformed C136/Go2 NOT passing parity.

## Capacity and item-granted skills (#147)

Existing report.go adds complete pinned report.c show_carry_capacity
**73–105**, show_item_skills_sup **108–193** and show_item_skills **196–210**.
The existing oracle already extracts their entire real C closure, reusing
weights/skills, names, inventory and cached output. No new dependency/pin,
later character/player root or projection-plus-render double pass. Preserve
separate travel weights, signed-int ratios, independent stored study/use
lists, recipient prerequisites, original assertions and indentation order.

At **0.87.0-alpha**, fourteen appended sections match real C/C/Go:
**1006 sections / 68035 lines / 1290206 output bytes**, preserving the exact
inherited **992-section / 67520-line / 1278600-output-byte** and
**403193-input-byte** prefixes. Numeric constants and reached branches were
checked against the immutable pin; #168 live native traversal still matches.
Four temporary missing-dependency controls (exit 0/134, ordinary/update)
reject an expected-abort section before status acceptance or recording
replacement, leaving its hash unchanged. Probe sources are removed.
Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded default tests pass, with explicit isolated
`G3_REPO` throughout. Logs: `tmp/issue147-*.log`; see
[exact closure, fixture branches, hashes and limits](command-compat#capacity-and-item-granted-skills-147).

No full-world/full-corpus/strict-lifecycle or stage-exit acceptance. Latest
strict evidence remains the parent's post-#145 gate: identical C repeats,
C/Go 0, matching streams and paths/types/modes, eight saved-report byte
differences, exit 1 before lifecycle; not refreshed by #146 or this task.
G3 #358 remains nonblocking, original division/Go panic unchanged; malformed
C136/Go2 is not passing parity. Independent parent acceptance passes the
routine tier; the strict stage gate remains pending.

## Inventories and unclaimed items (#146)

Existing report.go adds complete pinned report.c inv_item_comp **213–220**,
extra_item_info **246–290**, show_char_inventory **294–344** and
show_unclaimed **746–780**. Existing inventory, magic, naming/number and
cached-output owners supply the real closure; the C oracle already extracts
these complete unchanged definitions. Direct section calls do not implement
later character/player roots or run their projections alongside rendering.
Original in-place sort timing, invalid/nonpositive rows, weight arithmetic,
capacity precedence, fighter/equipment/aura visibility, totals/empty output
and literal plural warning bytes are preserved.

At **0.86.0-alpha**, ten appended cases pass complete actual C/C/Go:
**992 sections / 67520 lines / 1278600 output bytes**, with the inherited
**982-section / 66938-line / 1262604-output-byte** and **399178-input-byte**
prefixes unchanged. #168 live native traversal still matches. Sequential
command/move/helper/RNG/code/loadsave/turn-driver/default-turn, bounded vet
and fresh bounded default tests pass. Logs are `tmp/issue146-*.log`; exact
closure, hashes, fixture constants and reached branches are documented in
[command-compat](command-compat#inventories-and-unclaimed-items-146).
No new source dependency/pin, later root/command or full-world/stage-exit
acceptance. The parent's latest post-#145 strict full run still owns the
eight saved-report byte differences and exit 1 before lifecycle controls;
this task does not refresh that gate. G3 #358 remains nonblocking and
unchanged; malformed C136/Go2 is not passing parity.

## Complete location report root (#145)

Existing display.go adds complete pinned display.c show_loc **1056–1101**;
report.go replaces turnEndLocReportsState with complete report.c
turn_end_loc_reports **1589–1642**, and main calls the renderer once. The
existing command oracle extracts unchanged real C and the original header
global, removing the final rejected show_loc dependency now that all accepted
#138–#144 sections are real. It reuses output/cache ownership and preserves
section order, display/header/routing flag semantics, allocations, touched-list
sorting/invalid skips, all nonsilent subtypes and both player formats.
Nine appended small-world sections compare complete-root report bytes and
raw/saved state, repeated/global resets, empty/silent controls, cache inode
reuse, shared RNG aftermath and source assertions. Numeric fixture constants
and actual reached branches were verified against the immutable pin.

At **0.85.0-alpha**, C/C/Go matches **982 sections / 66938 lines / 1262604
output bytes**, preserving the full **973-section / 66015-line / 1241870-byte
output** and **393356-byte input** prefixes. #168's separate live native-order
check passes unchanged. Expected-abort marker injection is rejected before
status acceptance/recording update, with the recording hash unchanged.
Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded default tests pass. See
[exact closure, hashes and evidence](command-compat#complete-location-report-root-145).
Independent parent acceptance passes the same tier. A fresh strict full-world
run has identical C repeats, C/Go status 0, matching raw streams and entry
paths/types/modes, but still eight saved-report byte differences at
lib/save/1/{100,200,202,203,204,205,210,players}. It exits 1 before strict
lifecycle controls; this is not stage acceptance. No new dependency/pin,
later report roots or commands. G3 #358 remains
nonblocking; division/Go panic are unchanged, malformed C136/Go2 is not parity.

## Ancillary location sections and posts (#144)

Existing display.go implements complete stage-pin display.c nearby cities
842–869, taught skills 872–892, posts 895–963, weather 966–1012, ruler
1015–1046 and location stats 92–122. The existing command oracle extracts
all six unchanged original definitions and storm.c can_see_weather_here
1661–1676, replacing the rejected posts dependency. Existing topology,
weather views, garrison ownership, names, ship weights and cached output
remain the owners; no POST/weather command or #145 location root is added.
Nineteen appended sections cover asymmetric valid/malformed posts, exact
quote/author/name order, empty text, stored city/skill order and safe havens,
all weather combinations and direct-versus-province storm visibility,
castle/ruler branches, ship stats, original assertions and raw state/cache.
At **0.84.0-alpha**, real C/C/Go compares **973 sections / 66015 lines /
1241870 bytes**, retaining exact **954 sections / 65118 lines / 1224413
output-byte** and **385519 input-byte** prefixes. See
[closure and task evidence](command-compat#ancillary-location-sections-and-posts-144).
Recording refresh/default command and #168's live native C/C/Go comparison
pass, as do sequential move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded default tests. Expected-abort marker injection
is rejected before recording update; temporary controls are removed.
No full-world or stage-exit claim. G3 #358 remains nonblocking; original
division/Go panic are untouched and malformed C136/Go2 is not passing parity.

## Inner locations and ship sightings (#143)

Existing display.go adds complete stage-pin display.c show_inner_locs 685–717,
show_sublocs_here 720–797 and show_ships_here 800–839, reusing descriptions,
owners/owner stacks, knowledge, barrier rendering and cached output. Fourteen
appended command-oracle sections compare complete unchanged real C, raw state
and report bytes: known/unknown hidden containers and ships, see_all 0/1/2,
direct-city versus building/recursive-city branches, deep nesting, open/closed
crosslinks with duplicates/order, actual ocean/port headings and owner/barrier
ordering. Original assertions and expected-abort dependency-marker rejection
are checked; numeric fixture constants are verified against pinned oly.h.
No root/posts, commands, new dependency/pin or output owner.
At **0.83.0-alpha**, C/C/Go matches **954 sections / 65118 lines / 1224413
bytes**, with the entire inherited **940/64201/1206911 output** and **378933
input-byte** prefixes unchanged, plus #168's live native-order comparison.
Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded defaults pass in independent parent acceptance.
See [exact closure, evidence and limits](command-compat#inner-locations-and-ship-sightings-143).
No full-world/stage-exit claim. G3 #358 remains nonblocking: existing division
and Go panic are untouched, malformed C136/Go2 is still not passing parity.

## Visible characters and owner stacks (#142)

Existing display.go implements complete stage-pin display.c highlight_units
547–558, show_chars_below_highlight 579–593, show_owner_stack 596–626 and
show_chars_here 629–682, leaving show_chars_below 561–576 and character
descriptions unchanged. Fifteen appended command-oracle sections execute the
complete unchanged C closure and compare hidden-first-owner versus later
hidden owners, nested preorder/prisoners/noncharacter pruning, own-unit
highlight spaces, province fog short-circuit and controls, ocean flying
suffix, original assertions, indentation, saved state and raw cached log bytes.
No new root, gameplay command, source dependency, pin or output owner.
At **0.82.0-alpha**, real C/C/Go matches **940 sections / 64201 lines /
1206911 bytes**; all inherited **925 sections / 63547 lines / 1195280 output
bytes** and **373845 input bytes** are exact unchanged prefixes. #168's live
raw native traversal check still passes. See [closure and bounded evidence](command-compat#visible-characters-and-owner-stacks-142).
Sequential command/move/helper/RNG/code/loadsave/turn-driver/default-turn,
bounded vet and fresh bounded default tests pass in independent parent
acceptance after correcting the ocean fixture. This is not full-world or
stage-exit acceptance. G3 #358's existing division
and Go panic are unchanged; the nonblocking malformed C136/Go2 failure
mismatch is not counted as passing parity.

## City market rendering (#141)

Existing buy.go implements complete stage-pin buy.c list_market_items
1027–1066 and market_report 1069–1120 using actual buyerList/sellerList/
reduceQty and the accepted names/output/cache owners. The command oracle
extracts both unchanged real C definitions. Nine appended sections compare
report bytes, BUY-before-SELL and native/stable seller order, actual sort-field
mutation, saved trade/inventory equality, reserves/balances, cloak 0/1/2,
zero availability, original production notices, record nonallocation,
widths/trailing spaces and MASTER/style/renamed-cache routing. No settlement,
BUY/SELL/movement handler, location root, new source dependency or pin change.
At **0.81.0-alpha**, C/C/Go matches **925 sections / 63547 lines / 1195280
bytes**, with all inherited **916 sections / 62819 lines** and input/recording
bytes exact unchanged prefixes, plus #168's live native C/C/Go check.
Sequential routine command/move/helper/RNG/code/loadsave/turn-driver/default
turn, bounded vet and fresh bounded default tests pass; this is not full-world
or stage-exit acceptance. See [exact closure, evidence and limits](command-compat#city-market-rendering-141).

**G3 #358 remains open and owner-declared nonblocking.** Actual C loading a
positive-quantity zero-price character BUY and rendering its market report
fails with SIGFPE on this x86-64 host; real cancellation controls do not.
The original division is retained and source-linked, with no local game fix
or new engine failure handling. A separate real renderer probe establishes
C **136** versus Go divide-by-zero panic **2**: raw/status failure parity
does **not** pass and is not counted among the recorded passing cases.
Normal validated-command generation of that malformed saved trade is not
established. See the command-tool evidence and
[G3 #358](https://github.com/playbymail/olympia-g3/issues/358).

## Portable admission oracle (#168)

The inherited multi-account admission recording encoded filesystem-specific
readdir order, not a port defect. Only its one section is replaced by five
single-candidate traversal probes for zero through four nonblank fields;
all other **902 saved sections** and input sections remain byte-identical.
Native traversal stays unsorted and is checked separately by a small live
real-C/C/Go raw comparison over one unchanged shared directory. It uses the
existing actual pinned `make_new_players` closure and rejects dependency
markers, mutations and unexpected status; no modeled expected output or
diagnostic normalization is used, and any native filesystem order is valid.
Plain Go tests explicitly skip only this live case without the tool's C
expectation; portable recordings always run. See
[exact byte and negative-control evidence](command-compat#portable-admission-recording-and-live-native-order-168).
Command `--update` and default comparison pass **903 sections / 62047 lines /
1160604 bytes**, plus a **28-line / 426-byte** live comparison; sorted and
reversed traversal overlays both fail. Bounded vet and fresh bounded tests
pass at **0.79.1-alpha**. Pin/dependencies and engine/report behavior are
unchanged. This local harness prerequisite does not include #140, full-world
or stage-exit acceptance.

## Route rendering (#140)

Existing dir.go implements complete stage-pin dir.c list_exit_extras,
list_exits_sup, list_road_sup, list_exits and list_sailable_routes (839–1072),
using actual exitsFromLoc and accepted names/knowledge/seasonal/output owners.
No movement handler, new source dependency, pin change or location root.
Thirteen appended command sections exercise byte/state, hidden/see_all/IN,
road/transit/barrier/Hades, ocean and seasonal branches, real assertions and
paired shuffle/RNG aftermath. See [exact closure and evidence](command-compat#route-rendering-140).
The inherited native-directory recording defect was fixed separately in #168;
its accepted input/recording prefix is preserved exactly. Independent C/C/Go
route acceptance passes **916 sections / 62819 lines / 1177391 bytes**, plus
live native traversal, at **0.80.0-alpha**. Sequential command/move/helper/RNG/
code/loadsave/turn-driver/default-turn, bounded vet and fresh bounded default
tests pass. No full-world or stage-exit acceptance is included.

## Noncharacter descriptions (#139)

display.go extends its accepted headers/character owner with complete unchanged
stage-pin ship, location, road and storm liner descriptions, capacity and
completion suffixes and display_owner. Character rendering is unchanged.
Existing names, topology, mine depth, weather and output/cache owners are
reused; weights.go adds the previously absent real ship-capacity/weight
helpers on top of accepted unit weights. The command oracle removes the four
rejected noncharacter dependencies and extracts actual unchanged C, including
first_char_here; show_loc and show_loc_posts remain rejected. No source
dependency or pin change, movement handler or full location root is included.
Thirteen asymmetric appended sections cover suffix/read order, signed widths,
recursive load, road distances, ownership, markup and source assertions, with
actual save_box state and raw MASTER/style log bytes. See
[exact closure and byte evidence](command-compat#noncharacter-descriptions-139).
C/C/Go is **903 sections / 62031 lines / 1160307 bytes**, preserving the entire
**890 input sections / 61457 recording lines** prefix. Sequential command/move/
helper/code/loadsave/turn-driver/default-turn/RNG, bounded vet and fresh bounded
defaults pass at **0.79.0-alpha**. Full strict stage-2 report-tree parity remains
deferred; no full corpus/full turn or stage-exit acceptance is claimed.

## Location headers (#138)

The existing display.go owner ports complete unchanged stage-pin display.c
8–49, 52–63, 66–74 and 125–166, reusing accepted topology, names, mine depth,
safe-haven, hidden and output/cache helpers. Command-compat extracts the real
definitions and replaces only its rejected show_loc_header dependency;
show_loc/show_loc_posts and noncharacter descriptions remain explicit stops.
Seven appended cases compare containment/suffix branches, raw record presence,
saved state, assertions and real header/barrier log bytes with banner/include
MASTER sections. See [the exact closure and byte evidence](command-compat#location-headers-138).
C/C/Go is **890 sections / 61457 lines**; the entire **883/61012** prefix is
unchanged. Sequential command/move/helper/code/loadsave/turn-driver/default
turn/RNG checks, bounded vet and fresh bounded defaults pass. No source
dependency or pin changes, full corpus/full turn run or strict stage-2
report-tree acceptance is implied. Full strict parity remains a later task;
its gate is unchanged.

## Bounded names/codes audit (#137)

Complete unchanged stage-pin code.c 47–174, 212–343 and 350–466 already
match the existing code.go owner; no engine changes or replacement naming
layer were needed. Existing code/helper/move coverage is retained. Six new
command cases directly observe garrison/raw-code distinctions, named/empty
port cities, noble fallback, absent/null/empty plural warnings, raw brackets,
quantity separators and nested 2048-byte truncation. See
[command-compat's source audit and byte evidence](command-compat#bounded-namescodes-audit-137).
After the separately accepted #167 parser fix, C/C/Go is **883/61012**;
the accepted **877/60622** prefix remains byte-exact. Sequential code/helper/
move/command/default-turn, bounded vet and fresh bounded tests pass. All
code-tool upstream dependencies are unchanged by G3 #357, so its own pin
stays `6a25f2d`. No full-world or stage-exit acceptance is claimed.

## Stage-2 pin move and follow-C (#131)

Current affected-tool pin: `e733acb63bacb21167a1f2589d8dcdc023364896`
(G3 #357), moved from `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`.
`G3_REPO="$PWD/tmp/upstream.git" tools/affected.sh OLD NEW` lists rng,
helper, io, box, boxfile, check, loadsave, olympia-cli, olympia-roundtrip,
olympia-spec, turn, turn-driver, command and move. Their script defaults and
current README references move together; entity/code/lists and peripheral
pins do not move. Historical task evidence, hashes and old bug descriptions
below describe their original pins and are not rewritten as new acceptance.

Go now follows the NP failure assertion after the existing transfer/refund
mutations, the explicit knowledge-owner assertion, the savage tally and
pre-draw cap, and the regular-player token disappearance notice. Command,
move and helper recordings change only from two deterministic real-C runs:
NP aborts replace the old alert/Paid continuation; unowned knowledge now
asserts rather than segfaulting; savage cap/state/mint aftermath and token
notice/output bits change; source assertion offsets follow the new files.
One appended command case distinguishes skipping capped forts **before**
reseeding the shared NPC stream. Original case/input prefixes are retained.
RNG's complete unchanged social helpers move from swear.c:428–528 to
425–525; the RNG recording stays identical.

USE 886 was unported (the USE callbacks in glob.go remain explicit NYI).
Stage 4 must preserve v_curse_noncreat's c->b minimum of 1 and existing
possession/destroyable-item refusals, then print exactly
`The noncreator possession curse has no effect yet.` and return FALSE.
It must not reach d_curse_noncreat or charge aura on that path. This task
does not implement the command. Existing turn erosion still logs NOTYET,
without loyalty mutation, exactly as the unchanged day.c branch does.

The turn C-only recordings are refreshed with `--full --record`; only the
pin/binary provenance changes, not captured world/diagnostic bytes or tree
hashes. The comparator still has stage-1 report exclusions until #134.
At #131, roundtrip case construction was deliberately unchanged; #133 below
resolves its expected-abort prerequisite without changing engine behavior.

#131 validation (Linux amd64, GCC 12.2.0, Go 1.26.4): all affected small
tools plus unchanged-pin code pass sequentially. Command is 859 sections /
59841 lines, move 135 / 7158, helper 37 / 1487; all C repeats deterministic.
Ten targeted deterministic roundtrips pass (world, release-swear-np without
turn, savage-cap/token-vanish/use-art/use-art-refused with and without turn).
Three sampled loadsave worlds compare 118097 boxes with zero check_db
messages and identical trees/raw streams/status. Refreshed full turn matches
all 190 state files/raw streams/status and all 13 lifecycle controls after
identical full C/C repeats. Spec remains exactly 19 pass / 104 fail / 0 error,
with no lost passes; only its generated provenance line changes. Vet and
fresh bounded default tests pass (Olympia 1.508s). Logs: `tmp/issue131-*.log`.
Full 202-world checks are the parent's accepted inherited baseline, not
rerun here. One changed raw C save line (`hl 110158 `) deliberately retains
its trailing space; all non-recording diff whitespace checks pass.

## Expected-abort corpus contract (#133)

At the same immutable `e733acb` pin, roundtrip has **201 meaningful saved-game
cases**, not 202 successful snapshots: all 196 inherited cases, `savage-cap`
and `token-vanish` before/after a turn, and ordinary `release-swear-np`.
The missing `release-swear-np+turn` is explicitly unavailable. Its unchanged
upstream check requires C status 134 and a `swear.c` FALSE assertion (G3 #357).
The harness validates the precise pinned line 180 twice under `run_limited`,
rejects dependency markers, sanitizer reports, wrong status/assertion and
timeouts, and requires identical raw streams/status/partial trees. It reports
one separate `EXPECTED_ABORT`, not a load/save pass or a manufactured game.
Other C turn failures remain errors. Other `G3_REV` overrides retain ordinary
`+turn` enumeration/status-0 preparation; no valid older snapshots are dropped.

`--list` lists available cases without building/running engines; an exclusively
unavailable selection errors explicitly. Successful cached inputs/C results
and format 1 remain unchanged, reusing `tmp/stage2-roundtrip`; controls always
rerun. Matching Go captures are removed per case. Full-game tools still see
202 worlds (their separate `world.tgz` plus these 201 cached inputs), exactly
the input corpus already accepted by the parent in `tmp/stage2-baseline-*.log`.
The four full-world tools were not rerun for this harness-only accounting
change. This does not accept Go turn parity for later unported commands,
full report artifacts, or stage-2 exit.

#133 validation: full `--determinism-all` gives 201 PASS / 0 FAIL / 0 ERROR
and one separate deterministic expected abort (153s); all 196 inherited names
are retained and PASS. The 201 existing C/C checks remain valid, every Go
load/save reruns freshly, and both C abort runs are fresh. Negative controls
reject wrong statuses/assertions, sanitizer/dependency markers and a real
unexpected C status-134 turn on ordinary snapshot preparation. Fresh `--list`
builds/runs no engines; the old `6a25f2d` override retains exactly 196 names.
Vet, fresh bounded default tests, default turn and command comparisons pass
(command 859 sections / 59841 lines, deterministic). Logs and enumeration:
`tmp/issue133-*.log`, full TSV `tmp/issue133-roundtrip-results.tsv`.

## Strict whole-tree turn gate (#134)

`turn-compat` now excludes no artifacts: C/C and C/Go compare every cwd
entry, including empty directories, types, all-entry permission/special bits
(the unchanged `lib/tree-modes.py` policy), regular bytes and raw streams/status.
Nonregular entries and failed listings/captures are errors, not omissions or
ordinary mismatches. Defaults remain tiny recorded cases with no world turn.
Two fresh unchanged pinned C turns match all 233 cwd entries. The strict
C/Go gate honestly exits 1 on eight `lib/save/1` report-byte differences;
both engines exit 0, raw streams match and no paths/types/modes differ.
Lifecycle controls now use the same strict comparator but are not reached
past this ordinary-turn failure. See `turn-compat/README.md` and
`tmp/issue134-turn-full{,-record}.log`; prior stage-1 artifact-excluding
acceptance below remains historical, not strict stage-2 acceptance.

## Output primitives (#136)

MASTER routing, styles/HTML, vector ownership and section primitives extend
the existing sout.go owner and command oracle, not a second report layer.
Complete unchanged sout.c is already linked; there is no pin or dependency
change. Twelve new exact-C cases cover routing/frame/state and raw log bytes,
including buffered renamed-inode output, single previous style, LEN truncation,
unwrapped newlines and 72-dash sections. Old input/output byte prefixes remain
unchanged. See command-compat's #136 evidence and negative controls.
For shared output changes run command/move/helper/code/loadsave/turn-driver/
default-turn sequentially, then bounded vet and fresh bounded defaults, always
with explicit `G3_REPO="$PWD/tmp/upstream.git"`. Full report renderers and -h
generation remain separate; the strict eight saved-report gap is not accepted.

## Check tiers

**Routine task acceptance.**
Seconds to a few minutes; writes little to `tmp/`.

```sh
go vet ./...
go test ./...                        # recorded cases only; the whole-game tests skip
tools/<tool>/run.sh                  # each tool whose area the task touched, default mode
```

For load and save work, add `tools/olympia-roundtrip/run.sh` (about 25 seconds once its C cache is built).
`tools/box-compat/run.sh --games 10`, `tools/boxfile-compat/run.sh --games 10`, `tools/check-compat/run.sh --games 10` and `tools/loadsave-compat/run.sh --games 10` compare a fixed sample of saved games when a change could show only on real data.
For turn work, add `tools/turn-compat/run.sh`: small C recordings, comparator
contracts and a byte-identical diagnostic run, never a whole-world turn by default.
For clock/calendar or progress output work, also run `tools/turn-driver-compat/run.sh`.
Its calendar/schedule suites use real pinned caln/wthr streams and preserve
process-lifetime picks/list/cursor separately from the streams' turn guards.
Only test overlays trace daily gameplay callbacks: scheduling acceptance does
not accept the unfinished weather/damage/healing/faery/event implementations.
Run rng/command/move/helper/loadsave/default turn comparisons alongside it;
ordinary tests use the small direct-stream recording, not a full world.
Natural weather creation and BIND (#75) extend command-compat rather than
duplicating its real mint, placement, garrison output and save-box lifecycle
in a weather-only tool. The new cases call actual storm.c definitions, not
callback traces. Run command/move/helper/rng/code/loadsave/turn-driver/default
turn comparisons sequentially. All four scheduled weather days together
remain opt-in integration; unrelated weather magic remains explicitly
unported. Additional source dependencies cloud.c, tunnel.c and
quest.c supply only original greater-region globals (tracked by affected.sh).
Storm movement/decay (#76) extends the same command comparison with complete
day.c storm_decay/storm_move/storm_owner_touch_loc and storm.c move_storm,
move_bound_storms and dissipate_storm. Monthly relocation uses set_where
directly; ship-bound motion emits weather and touches owners. Shared day-zero
location initialization and month-end call positions are preserved. Run the
same sequential comparisons above. No new source dependency/pin or SAIL
handler, weather-magic command or full-world turn success is implied.
Faery hunts (#77) also extend command-compat, extracting complete faery.c
create_elven_hunt/warn_human/auto_faery_sup/auto_faery definitions alongside
the existing original faer prefix. This reuses #111 entry hooks, actual
creation/mint, NPC/exdr routing, queue/dispatch, logs and save-box/order writers,
not a separate modeled NPC harness. Cases cover the 15-unit cap (all faction
units), slot/retry ocean rejection, quantities/order, shared warning memory,
stone/stack filtering, and real spawned WAIT → eight-day MOVE → WAIT.
Daily call position is unchanged. Run command/move/helper/rng/code/loadsave/
turn-driver/default turn sequentially. MESSAGE/ATTACK and unrelated commands
remain explicitly unported; no new source dependency/pin or full-world success.
Structural destruction (#120), the bounded prerequisite split from #78,
extends command-compat with complete u.c add_structure_damage, sink_ship,
building_collapses, get_rid_of_collapsed_mine and kill_stack_ocean. It reuses
real death/survival/inheritance/extraction, arrival/nearest-land routing,
output/VECT and bound-storm state, not successful gameplay substitutes.
Run the same eight comparisons sequentially before vet/default tests. No new
source dependency or pin: u.c was already tracked. Daily coastal/mine/inn
events remain #78, character damage/healing/ocean polling #79, and monthly
mine decay #88. Full worlds remain opt-in; no stage-exit success is implied.
Daily coastal and location damage (#78) extends that same command harness
with complete day.c near_rocky_coast, ship_coastal_damage, mine_calamity,
inn_calamity and random_loc_damage, plus u.c count_loc_char_item and the
original mine_depth macro. It uses the real structural destruction above,
worker consumption, VECT/logging, topology/ownership and shared wthr stream.
Actual-hit cases, independent C gate/range endpoint controls and paired
shuffle/mint aftermath distinguish the per-day leaf keys and turn guards.
Run command/move/helper/rng/code/loadsave/turn-driver/default turn sequentially
before vet/default tests. No source dependency or pin change; no duplicate
weather harness, character damage helper, unrelated command or full turn.
Ocean polling and weekly healing (#79) extend command-compat, not a new
upkeep-compat tool: complete move.c init_ocean_chars/check_ocean_chars and
day.c heal_char_sup/heal_characters use the accepted real stack-death,
survival/shore-arrival, MATES inheritance, command queues, output and save
closure. Existing load-time initialization and dailyEvents positions/weekly
gates are unchanged. The world's existing UpHeal stream supplies both keyed
subdraws; the C driver extracts the original guarded upkeep prefix. Actual
gameplay cases include real init_ocean_chars → checks, snapshot mutation,
plain/inn recovery thresholds, health clamps/deaths and per-day/turn replay.
Weekly boundaries remain in the real-C schedule suite, which proves schedule
order only, not healing gameplay. Run the same eight tools sequentially
before bounded vet/default tests. No new source dependency/pin or RNG adapter.
The older broad boundary notes do not add add_char_damage to #79: that
unrelated helper is not required by these functions and remains unported.
Daily prisoner, ring-aura, curse and dog events (#80) extend command-compat,
not a duplicate upkeep-compat: complete stack.c weekly_prisoner_escape_check,
day.c dogs_bark_at_hidden_chars/increment_stone_ring_aura/noncreator_curse_erode
and u.c bark_dogs reuse the accepted escape/movement/death, VECT/logging,
inventory, knowledge and effective-aura closure. The active C curse branch
only logs NOTYET; its public USE 886 limitation is tracked as G3 #356 and is
intentionally reproduced without loyalty mutation. The harness also extracts
real daily_events and its original calendar prefix for small-map composition
through actual gameplay callbacks, distinct from schedule-only traces.
Run command/move/helper/rng/code/loadsave/turn-driver/default turn sequentially
before bounded vet/default tests. No new source dependency/pin or RNG adapter;
production dailyEvents order/gating stay intact.
Day-1 trade matching (#81) extends command-compat rather than creating a
duplicate economy-compat: complete buy.c reduce_qty, attempt_trade,
match_all_trades and check_validated_trades reuse the real market lists,
inventory/lore/token, command/arrival and output/save lifecycle. The validator
iterates its live growing queue and clears after the loop; real completion
guards and the day-1 call position remain intact. Small cases cover reserves,
partial balances/quantities, seller ties/city-last order, city generation and
gold metrics, unique/token transfers, cloaks, a live settlement cascade, real
command/evening/arrival triggers and actual process_orders on a tiny map.
Run the same eight tools sequentially before bounded vet/default tests.
No new source dependency/pin or RNG adapter. BUY/SELL/FIND commands and
month-end restock (#83) remain unported; no full-world/stage-exit success or
new observed full-world boundary is inferred from the small-map check.
City gold and existing-world production (#82) extend command-compat with
complete day.c compute_civ_levels, move_city_gold and post_production,
produce.c's original terrain/mine tables plus replenish, mine_production and
location_production, seed.c seed_taxes and u.c clear_temps. The shared real
inventory, traversal, topology, trade queue and save-box closure is reused.
Small recorded cases check asymmetric civ diffusion, first-of-type buildings,
crown bonuses, every resource/depth row, stock floors, poppy capping/doubling,
coastal flotsam and city-tax transfer/opium/pillage rounding. postProduction
now executes real location_trades (#83), location production and tax reseeding
before setting the post flag; command-compat records successful composition.
New-world load initialization remains explicitly stage 5. Run the same eight
comparisons sequentially before vet/fresh bounded defaults; no new source
dependency/pin, RNG adapter, command handler or full-world success is implied.
Monthly market restock (#83) extends command-compat with complete buy.c
find_trade/new_trade/add_city_trade, opium_market_delta, expire_trades,
loc_trade_sup, trade_suffuse_ring, location_trades and is_tradegood, plus
art.c new_suffuse_ring. The original economy prefix/macros and opium_data
are extracted from the pin. It reuses persistent w.streams.BeginEconomy/Econ*
and real mint, settlement, lore, inventory and serialization, not an adapter.
Positive-quantity filtered loop_trade remains a snapshot sharing pointers;
raw opium/expiration loops and the live validated settlement queue are preserved.
Small recordings cover sold-out/partial/no/absent opium orders and clamps,
annual seasons and override, surplus floors/metadata, held tradegood expiration,
asymmetric markets and all five ring kinds, gates/costs, sold-out replacement,
real purchases/lore and next-turn creation. Successful post_production now
owns the former temporary NYI test's ordering contract. seed_city_trade is
INIT-only, and existing route initialization already works: neither FIND/BUY/
SELL commands nor new-world initialization is added to monthly acceptance.
See the source-verified scope clarification on #83 and command-compat README.
Run command/move/helper/rng/code/loadsave/turn-driver/default-turn sequentially
before bounded vet/fresh defaults. No full-world/stage-exit acceptance implied.
For command parsing, dispatch, handlers or completion, also run
`tools/command-compat/run.sh` (small real-C dispatch/scheduler recordings, no world turn).
Monthly faction/tax/seasonal settlement (#84) extends that same command tool:
complete day.c collect_taxes, add_claim_gold, add_noble_points,
add_unformed_sup/add_unformed, special_locs_open/close and clear_orders_sent,
plus c1.c next_np_turn. Small asymmetric accounts reuse real mint, inventory,
NP/lord, topology, output and save closure. Cases cover player subtype/missing
records, exact 0/4/5/6 unformed floors without pruning, shared-mint aftermath,
NP turn/width boundaries, castle ownership/order/tax truncation and all seasonal
branches/precedence with real logs. The existing turn-driver phase comparison
owns post_month's month-2/month-6 gates; it accepts scheduling, not skipped
gameplay. garrison_gold remains an explicit stop at its unchanged post_month
position (#116 after maintenance #87), so #84 does not accept full post_month
integration. Run command/move/helper/rng/code/loadsave/turn-driver/default-turn
sequentially before bounded vet/fresh defaults. No duplicate upkeep/garrison
harness, new source dependency/pin or full-world success is implied.
Monthly aura/shroud/loyalty decay (#85) extends command-compat with complete
day.c increment_current_aura, decrement_ability_shroud/region_shroud,
decrement_meditation_hinder/loc_barrier, loyalty_decay, pillage_decay,
relic_decay and hide_mage_decay, plus scry.c notify_loc_shroud and u.c loyal_s.
Small recordings distinguish per-entry artifact recovery, signed countdown
quirks, skilled shroud recipients, all loyalty kinds and pre/post-decrement
gates, keyed hit/miss and turn guards, real desertion/WAIT/stack/dependents,
independent NP exclusion, and relic ownership/inventory/lore. The real shared
upkeep stream and lifecycle/output/save closure are reused; post_month's
order is unchanged and garrison_gold remains the #116 stop. scry.c is newly
tracked by affected.sh; there is no new pin or upkeep tool. Run the same
eight comparisons sequentially before bounded vet/fresh default tests.
No full-world, maintenance/income, report-body or stage-2 scope is implied.
Monthly inn/temple income (#86) extends command-compat with complete day.c
inn_income/temple_income and loc.c count_loc_structures. It extracts the
original guarded up_income prefix and gm.c gold_inn/gold_temple scalars,
reusing persistent w.streams.UpIncome and real ownership, inventory, logs and
save-box writers. Inn sharing counts both matching direct location subtypes,
including unowned structures, at subloc(inn), not its province. Pillage divides
before sharing; traveller bonuses are outside sharing and gated by no pillage.
Temples have no sharing or occupant/capacity scaling: known religion on the
first direct character gates a flat 100 offering before local pillage.
Small cases cover unequal incomes, nested locations, exact subkeys/ranges,
all pillage flavors, absent records, same-turn master/day guards, next-turn
reseeding, signed-int accounting and paired mint/weather/exit-shuffle controls.
The raw invalid count-helper case preserves loop_here's source assertion at
loc.c:358 and requires status 134; ordinary income cases require status 0.
Post-month order and the #116 garrison stop are unchanged. Run command/move/
helper/rng/code/loadsave/turn-driver/default-turn sequentially, then bounded
vet and fresh default tests. No new source dependency/pin, duplicate upkeep
tool, full-world success, report-body or stage-2 acceptance is implied.
Monthly maintenance/starvation/animal deaths (#87) extends command-compat,
not a duplicate upkeep tool. Complete day.c men_starve, unit_maint_cost,
charge_maint_sup, charge_maint_costs and animal_deaths reuse the accepted
maint_cost tariff, real stack funds/autocharge, inventory consumption,
formatters and persistent UpStarve/UpAnimal leaves. The harness extracts
complete move.c check_captain_loses_sailors for actual moving-captain crew
loss and wait penalties, not a successful SAIL substitute. Cases distinguish
per-item round-robin payment/casualties from cheaper-first allocation, stack
funds and prisoners, noble-item exclusion only in the cost calculation,
ceiling-one-third losses, both starvation flavors, one-based animal indices
and exact 9/10 threshold, faction filtering and same/next-turn guards.
Macro assertions keep their source file/line. Run the same eight comparisons
sequentially before bounded vet/fresh defaults. Post-month call positions and
the preceding #116 garrison stop remain unchanged; no new source dependency,
pin, report-body, full-world or stage-exit acceptance is implied.
MOVE/GO and all six movement synonyms use command-compat's actual C callbacks
through dispatch/evenings, including queued multi-day NPC MOVE → WAIT. Also run
move-compat for the shared travel, message and raft crew-loss contracts; full
SAIL/FLY and the eventual whole-engine raft fixture (#90) are separate work.
Stack-arrival knowledge, contacts, weather views, ocean bookkeeping and city
matching also use command-compat, reusing its actual Hades/Faery hooks, keyed
streams, entity creation, order queues and real database writer. Candidate
settlement now executes real `attempt_trade`; empty and unmatched city paths
also execute real `match_trades`, not a skipped callback.
For exit topology, distances, randomized exit selection, weights or travel
calculations, movement parsing/admission or movement bookkeeping, and character
descriptions/stack display (including rank/aura/equipment), departure/arrival
messages or garrison notice/spotting, also run
`tools/move-compat/run.sh` (small real-C exit-field/order/RNG/weight/travel and
parsing/admission/stack/message recordings, no world turn). Its message helpers
reuse the already listed `u.c`, `garr.c`, `move.c`, `glob.c`, `display.c` and
`sout.c` dependencies; output includes real faction logs and serialized state.
For mint allocation or character/unique-item lifecycle changes, run
`tools/code-compat/run.sh` and `tools/command-compat/run.sh`: the former checks
the allocator and shared mint draws, the latter real creation/deletion,
stack placement, ownership and serialized state (no full combat or death).
Movement prisoner escapes and nearest-land routing use command-compat's real
stack release, entity RNG, arrival, swimming, logging and serialization closure;
also run move-compat and rng-compat for their shared movement/stream contracts.
Hades population spawning and NPC route queueing also extend command-compat:
actual seeded zero/existing/capped populations, direct-target filtering,
asymmetric topology, direction persistence, shared exit shuffles and generated
MOVE → WAIT execution. Run move/helper/rng/code comparisons alongside it for
the shared topology, creation and stream closure. Autonomous NPC dispatch
also runs in command-compat; it does not substitute successful combat.

Savage fort spawning and savage autonomous order generation likewise extend
command-compat, including initial population caps, per-fort keyed checks,
recursive human filtering, ordered targets and real spawned MOVE execution.
Run command/move/helper/rng/code/loadsave/turn-driver/default turn comparisons.
Queued USE/ATTACK/RAZE/UNSTACK/DIE are not successful command execution; their
unported handlers retain explicit stops. The command tool additionally reads
`olympia/savage.c` (also tracked by `affected.sh`).

Unsworn/mob/undead autonomous programs also extend command-compat, including
the 4/5-turn mob boundary, summoner clearing, location-keyed RNG, saved-order
append and generated MOVE → WAIT. Run the same command/move/helper/rng/code/
loadsave/turn-driver/default-turn checks. ATTACK/PILLAGE/RECLAIM execution
remains explicitly unported. The command
tool additionally reads `olympia/necro.c` (also tracked by `affected.sh`).

Bandit programs and day-0 NPC dispatch extend the same command comparison:
keyed victim identity in saved target order, special-program safety filters,
all child UNSTACK orders, fallback directions and actual MOVE → WAIT;
summon/prisoner/command/queued-order skips and every program branch; real
savage-before-Hades creation and shared RNG/queue/serialization state. The
existing process_orders callback now invokes this dispatcher. Garrison guard
setup, daily default-garrison pay and token reconciliation are real; calendar
and later daily and monthly phases remain explicit stops. After #73 the
opt-in full-world diagnostic reaches `match_all_trades` and still fails;
it is not stage acceptance.

Garrison guard setup and default-garrison payroll extend command-compat, using
the real inventory/strength/maintenance tariff, subtype iteration, descriptions,
ownership and location/garrison output routes. Cases cover the 9/10 boundary,
every guard reset, nested stack leaders, captive/default/NPC skips, unequal
balances/pay, zero-pay generation/output and invalid negative-pay assertions.
Month-end garrison tax forwarding remains #116 after real maintenance #87.
Run command/move/helper/code/loadsave/turn-driver/default-turn comparisons.

Month-end garrison tax forwarding (#116) extends command-compat with complete
garr.c `garrison_gold` **446–549**, build.c `province_subloc` **160–177** and
gm.c's original `gold_provinces` scalar. The supplied 446–558 range includes
the start of unrelated `count_castle_garrisons`; the complete settlement ends
at 549 in the pin. It reuses real #87 maintenance/starvation, stack funds,
strength, ownership, inventory, output and signed 32-bit accounting.
Small cases distinguish consumed versus restored local tax, own-gold spending,
9/10 strength before/after starvation, odd halves, multiple castles/garrisons,
owner/default/missing-misc skips and persistent tax/forward fields across calls.
build.c is a new command-compat dependency, tracked by affected.sh. The
postMonth call position is unchanged; the next explicit boundary is #88's
ghost_warrior_decay, not a successful full postMonth. Run command/move/helper/
rng/code/loadsave/turn-driver/default-turn sequentially before bounded vet and
fresh bounded defaults. No new pin, substitute gameplay, full-world acceptance
or report-body/stage-2 scope is implied.

Monthly body, ghost and mine decay (#88) extends command-compat, superseding
the preceding #116 boundary note. Complete day.c collapsed_mine_decay
1525–1553, ghost_warrior_decay 1556–1588, corpse_decay 1591–1623 and
dead_body_rot 1626–1652 reuse real collapsed-mine evacuation/deletion,
unique-body destruction/NP refunds, inventory, output and save writers.
The original upkeep prefix supplies up_corpse; Go uses persistent UpCorpse.
Small recordings distinguish signed mine-delay wrapping before clamp,
full is_npc filtering, exact 12-turn body age, missing p_char allocation,
old/previous-lord refunds, keyed 0/1/2 corpse draws and repeated/next-turn
calls. Nested mines and stacks exercise snapshot relocation and deletion-safe
subtype chains. Macro assertions retain their caller source file and line.
See the source-backed scope clarification on #88 and command-compat README.
Run command/move/helper/rng/code/loadsave/turn-driver/default-turn sequentially
before bounded vet/fresh defaults. No new tool, source dependency, pin, RNG
adapter or postMonth call-position change. No full-world acceptance, unrelated
command, report-body or stage-2 scope is implied.

Monthly autoquit, links, quests and ranks (#89) extend command-compat with
complete day.c auto_drop/link_decay/quest_decay and garr.c top_ruler,
nprovs_to_rank/find_kings/lower_noble_rank/det_noble_rank_sup/
determine_noble_ranks. Accepted ownership/pledge, strength, queue, output,
rank display, mint and save writers are reused. Cases cover asymmetric
autoquit 4/5/6-turn boundaries and missing records, appended faction orders,
signed countdowns and month-8 wrap, every rank threshold through actual
owner counting, strength 9/10, recursive pledge ties, kingdom 14/15 and mixed
top rulers. subloc_quest is quest_late itself (oly.h:1082), not a separate
quest ID; kingship intentionally does not require strength 10.
The tool now extracts complete real post_month 2147–2203 and runs two
successful small-world compositions through all actual callbacks, with
maintenance, owned-province forwarding, taxes, production, decay, queue and
rank state/output. Phase traces still establish scheduling only. The final
four turn_phases.go NYIs are gone and the empty boundary file is deleted;
postMonth order, token callback and final flag assignment are unchanged.
See the source-backed clarification on #89 and command-compat README.
Run command/move/helper/rng/code/loadsave/turn-driver/default-turn sequentially
before bounded vet/fresh defaults. No new tool, source dependency, pin or RNG
adapter; no full-world/stage-exit acceptance, unrelated command or stage-2
report-body work is implied.

Output ordering and player entity state (#121) extend command-compat with
complete report.c determine_output_order/sort_for_output and player_ent_info,
including original show_unclaimed/extra_item_info and perm.c print_admit/
print_att rendering dependencies. C executes those definitions unchanged;
Go preserves their persistent state and observable plural-name warnings,
not deferred report artifacts. Region traversal and depth-first here order
sort established faction unit lists (including silent players) without
allocating absent player records. player_ent_info instead allocates every
nonsilent player record, sorts inventory, admissions and all three attitude
lists in place, and retains invalid entries and unsorted admission members.
Raw observers supplement actual save_box bytes where the writer filters
invalid references. Cases include opposing IDs/nested stacks/ships, stale
unlocated order and ID ties, repeats/relocation, silent and missing-record
controls, all player subtypes, equal admission targets, continuation lines,
and unchanged character state/knowledge. The C driver initializes spaces
as main does, needed by original attitude continuation padding. report.c,
perm.c and sout.c were already tracked: no new source dependency or pin.
Run the same eight comparisons sequentially before bounded vet/fresh defaults.
Main integration, character/summary closure (#124), lore (#122) and stage-2
rendering remain separate; no full-world/stage-exit acceptance is implied.

Character report state and remaining allocations (#124) extend the same
command comparison with complete original character_report/char_rep_sup and
their rendering closure. Go extends report.go's state projection: lazy
two-pass pledge backlinks allocate CH across ALL characters only at the first
eligible nonsilent/nonprisoner report; allocated-empty temporary lists can
retain bare CM sections. Only eligible character inventories and original
hostile/neutral/defend lists sort; skill reporting sorts copies, not saved
skills. Initial prisoner output can allocate a missing owner player record
without collecting backlinks. Real turn_end_loc_reports/player_banner and
summary_provinces allocation callers are compared, including all player
subtypes and the original phase order. summary.c is a new source dependency,
now tracked below and in affected.sh. Cases compare actual save_box bytes,
raw list/record observers, streams/status and source-located assertions.
Run command/move/helper/rng/code/loadsave/turn-driver/default-turn sequentially
before bounded vet/fresh defaults. No main integration, lore or full Go
report rendering, new tool/pin or full-world/stage-exit claim is implied.

Skill-lore state and input diagnostics (#122) extend command-compat with
complete original lore.c:10–387 functions and u.c weeks, reusing the real
queue_lore/test_known/set_known, skill_school and save writers. Go extends
the existing lore.go, not a duplicate knowledge implementation. ALL character
skill rows queue, regardless of learning state, prisoner or silent owner.
Delivery sorts the original queue, skips adjacent duplicate reads and retains
the queue, including duplicates. GM lore reads schools before nonschools;
listing assertions still precede reads. Sheets are actually read with getlin
semantics; missing opens are nonfatal, read errors after successful open are
not open failures, and only the two original $ functions are legal. Stage-2
text is deferred, not knowledge, queue mutation, allocation or raw diagnostics.
Small recorded cases include real legal expansion/substitution execution in
C, modified files, directories, ENOENT/ENOTDIR and source-located assertions.
No new upstream source dependency or pin: lore.c, u.c and use.c were tracked.
Run the same eight comparisons sequentially before bounded vet/fresh defaults.
No main integration, full Go renderer or whole-world/stage-exit claim.

Residual report-tail state (#125) extends report.go and the existing
command-compat oracle with complete original templates, account, summary,
player-report, GM-report and include functions. Go preserves player/corpse
record allocations, in-place known-list ordering, account style reset and
existing-output include validation, not report bodies or output-only scratch
statistics/validation. Shared sout.go log opening preserves real open-failure
diagnostics before output allocations; deferred artifacts do not fabricate
output bits. The account-summary input/command block is #if 0 in this pin.
Original C executes its actual report closure; unexercised location-header
and notice-rendering dependencies still abort with rejected markers, never
successful rendering stubs. No new source dependency or pin: gm.c/order.c/
report.c/summary.c and all called source files were already tracked below
and in affected.sh. Run the same eight comparisons sequentially, then bounded
vet/fresh defaults. Admissions, file lifecycle and main wiring remain #126/
#123, including main's explicit stop; full-world/stage exit remains #90.

Cwd-relative admission probing (#127, child of #126) extends add.go and the
same command oracle with complete original add.c fetch/parser/prober and
no-admission traversal definitions. Input opens, native directory iteration,
dot skips, partial/directory reads, account validation and raw diagnostics
are real stage-1 work. Go stops explicitly before alloc_box on complete
valid input; new-player gameplay remains stage 6. Empty report loops do not
skip new_player_list's stage or vector_players side effects. Saved bytes and
raw streams/status are compared for no-admission paths; independent fetch
controls and Go boundary tests do not claim successful admission parity.
add.c is a new command dependency, tracked below and by affected.sh. Run
command/move/helper/rng/code/loadsave/turn-driver/default-turn sequentially,
then bounded vet/fresh defaults. Post writers (#128), cached output/Times
(#129), main wiring (#123) and stage-2 artifacts remain separate. Main's
postMonth stop is unchanged; there is no full-world/stage-exit claim.

Post-turn writers (#128) extend the existing olympia.go owner and command
oracle with complete unchanged main.c:50–290 definitions and real existing
dependencies/save writers. Players/email/totimes/forward/factions open and
close actual files; nonfatal perror diagnostics and successful-open ownership
and garrison validation are stage 1, while artifact bodies alone remain
stage 2. Cases compare raw streams/status, saved bytes, absent-record flags
and file presence, with asymmetric success, ENOENT/ENOTDIR and each EISDIR
failure in original order. main.c is newly tracked below and in affected.sh.
Run command/move/helper/rng/code/loadsave/turn-driver/default-turn sequentially,
then bounded vet/fresh defaults. No suppression option, main wiring, cached
output/Times lifecycle (#129), valid admission or stage-exit acceptance.

Cached output ownership and Times lifecycle (#129) extend that same command
oracle, not a second cache or rendering layer. The existing openPlayerOutput
owner now retains sout.c's bounded move-to-front FILE cache: hits reuse the
open inode, misses append/open, full caches close the least-recently-used
stream first, and close_logfile closes/resets unless immediate. Both ordinary
output and report state projections borrow that owner. Existing clib stdio
supplies append buffering, fflush and the open-time flush_always setting.
Complete unchanged sout.c/io.c remain linked, including real save_logdir;
complete open_times/close_times and their original c2.c declarations are now
extracted. These source dependencies were already tracked below and in
affected.sh; the pin and dependency set are unchanged.
Cases distinguish cache hits from reopen failures, MRU/eviction flushes,
close/reset, ignored mkdir failures, both rename decisions and diagnostics,
save/no-save path controls, Windows flag path branches, Times open/close/reset
and fatal press-before-rumor ordering. An atexit filesystem observer runs
before exit's stdio flush, so rumor failure checks the actual press flush.
Raw streams/status and original save_box bytes are compared; artifact bodies
are not. Windows controls exercise the original win_flag branch on the host,
not a native Windows build. Run command/move/helper/rng/code/loadsave/
turn-driver/default-turn sequentially, then bounded vet/fresh defaults.
Main's explicit postMonth stop is unchanged: main wiring, final save_db/-n
integration, admissions/-A, stage-2 rendering and stage exit remain separate.

Main tail integration (#123) wires the accepted projections in actual main.c
order after post_month, including lore reads and cwd-relative admissions,
gen_include_section last, close_logfile and the post writers. check_db follows
reports; save_db precedes save_logdir, which runs only with both -S and -r.
Times masthead/close, final stage and stdout follow. Valid admissions and -A's
external charging remain explicit unsupported branches; only report artifacts
are deferred. Routine command/move/helper/rng/code/loadsave/turn-driver/default
turn comparisons, bounded vet and fresh defaults pass. Actual turn --full
matches raw streams/status and all 190 state files, with entire-capture C/C
determinism. Thirteen opt-in main lifecycle controls add no-save, act probing
and fatal Times versus nonfatal post-writer failures, repeating actual C and
comparing unfiltered streams/status and state. Matching captures are removed
per case; defaults stay fast. See turn-compat/README.md for exact evidence.
No recordings, pins, dependencies or comparator exclusions change. #90 still
owns final all-tools/all-world stage acceptance; stage 2 is not implemented.

Inventory transfer and token acquisition (#117) extend command-compat with
the complete real `consume_item`, mutually recursive `move_item`/`drop_item`,
and compiled `move_token` acquisition closure. Real lord interruption/NP,
creation/shared mint, inventory/lore/trade, naming and output dependencies
are reused. Run command/move/helper/rng/code/loadsave/turn-driver/default-turn
comparisons. No new upstream source dependency or pin change is needed:
`code.c` is already linked whole and `art.c`/`u.c` are already tracked.
Independent melting and `check_token_units` remain #73 after #118/#119;
the acquisition path contains no successful death/reclamation substitute.
GIVE/DROP/USE handlers and `create_npc_token` are outside this helper task.

Inheritance, prisoner release and desertion (#118) extend that same command
comparison with complete `stackmate_inheritor`, `take_unit_items`,
`put_back_cookie`, `drop_stack`, `extract_stacked_unit`, `move_prisoner`,
`unit_deserts`, `plural_item_box` and `box_name_qty`. The natural disappearance
path runs real desertion/extraction, cookie return and dead-character
conversion, without calling unfinished death/reclamation helpers. Existing
`code.c` is linked whole; `u.c`, `stack.c`, `swear.c` and `gm.c` are already
tracked dependencies (including original combat gold scalars). No pin or
dependency-table/affected.sh change is needed. Run command/move/helper/rng/
code/loadsave/turn-driver/default-turn comparisons. Full shared death remains
#119 and token reconciliation remains #73; no UNSTACK/DIE/RECLAIM handler or
whole-turn/stage-exit acceptance is implied.

Shared character death/reclamation (#119) extends command-compat with complete
`survive_fatal`, `char_reclaim`, `dead_char_body`, `kill_char`,
`hack_unique_item`, `forget_skill`, `random_hades_loc` and `char_rep_location`.
The last helper is a **new report.c dependency**: its location strings feed
the real death log, not deferred stage-2 report assembly. Original use.c
`skill_school` is also extracted, reusing its existing Go implementation.
Run command/move/helper/rng/code/loadsave/turn-driver/default-turn comparisons.
Actual survival, corpses, transcendence/arrival, token removal/destruction,
lord/NP, inheritance, extraction and interruption paths are exercised; no
RECLAIM/FORGET/DIE command, sibling damage/combat/resurrection helper or token
reconciliation is implied by #119 alone (reconciliation follows in #73).
G3 #355's missing token-vanishing notification is deliberately
reproduced and source-linked, not independently fixed.

Token melting/reconciliation (#73) extends command-compat with complete
art.c `melt_token_units` 1054–1080 and `check_token_units` 1197–1259, reusing
the full accepted death/reclamation, creation, lord/NP, inventory, output,
orders and shared RNG closure. The real process_orders definition and existing
clock/progress helpers are also extracted for a day-zero owner assertion
after actual NPC/guard setup. Only phase instrumentation is omitted; later
unported turn dependencies fail explicitly and cannot be recorded as success.
Existing processOrders/postMonth callback positions are unchanged. Run command/
move/helper/rng/code/loadsave/turn-driver/default-turn comparisons. No new
source dependency or pin change, unrelated command handlers or artifact
generation. Full-world turn success and stage exit remain unfinished.

**Pin move.**
Run `tools/affected.sh OLD NEW` and, for each tool it lists, move the pin in its `run.sh` and README, then rerun it (with `--update` where the tool records C output, and review the diff of the recorded file).
Then run the routine checks.
`olympia-roundtrip` rebuilds its C cache for the new revision by itself, which takes about 3 minutes.
`turn-compat` uses `--full --record` to refresh its small C recordings and
determinism metadata; review those bytes. The full Go turn succeeds after
#123; no NYI failure is accepted as a successful turn baseline.

**Stage exit.**
Everything, with the whole-game comparisons:

```sh
go vet ./...
G3G0_FULL=1 go test ./...
tools/<tool>/run.sh                  # every tool pinned to the olympia revision
tools/olympia-roundtrip/run.sh --determinism-all
tools/box-compat/run.sh --full       # after olympia-roundtrip, whose games it reads
tools/boxfile-compat/run.sh --full
tools/check-compat/run.sh --full
tools/loadsave-compat/run.sh --full
tools/turn-compat/run.sh --full      # strict stage-2 whole-tree gate; reports included
```

Run these one at a time: `box-compat`, `boxfile-compat`, `check-compat` and `loadsave-compat` read the games `olympia-roundtrip` caches, so do not run them while `olympia-roundtrip --fresh` rebuilds that cache.

**Current acceptance evidence (#90, 2026-10-28).** Linux amd64, GCC 12.2.0,
Go 1.26.4, stage pin `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`.
Roundtrip `--determinism-all` passes 196/196 cases with C/C determinism,
identical raw streams/status and every file; full box comparison passes
197 games / 7,947,575 boxes, plain and pretty. Logs are under
`tmp/issue90-evidence`; the bare source is `tmp/upstream.git`, the reusable C
cache is `tmp/issue104-roundtrip`, and full tools use
`GAMES="$PWD/tmp/issue104-roundtrip/cases"`.

Default engine compat checks and vet pass. The actual full turn repeats C
deterministically, then matches Go status 0, raw stdout/stderr and all 190
state files; all 13 lifecycle controls pass. Fresh bounded
`G3G0_FULL=1 go test -count=1 -v -exec "bash $PWD/tools/lib/run-limited.sh" ./...`
passes 11 packages, with nine explicitly opt-in/child-entry skips; the four
whole-game tests actually run through the full compat tools, not that standalone
suite. Spec progress is 19 pass / 101 fail / 0 error, no recorded passes lost;
golden initialization still stops at unsupported `-i`. See the turn/spec READMEs.

The initial run exposed false mode failures from BSD `stat` arguments on Linux;
local #130 resolved that harness defect and was independently accepted/pushed
at `686717f70ec3684b8b6b2f6c35fa8dcb3c18d615` (0.76.1-alpha), without changing
the engine, C oracle or recordings. Parent reruns now pass all required full
gates: boxfile, 197 games in each of fast/fast-pretty/slow, 7,947,575 boxes;
check, 197 games / 7,947,575 boxes; loadsave, 197 games / 7,947,567 saved boxes.
Check/loadsave report zero check_db messages; C/Go trees, raw streams and
statuses match. Evidence: `tmp/parent130-{boxfile,check,loadsave}-full.log`
and #130's acceptance comment. The earlier failure logs remain historical
discovery evidence, not current blockers. Final #90 acceptance passed the
parent's independent rerun and documentation review: default engine tools,
196 roundtrips, full box/turn comparisons, vet, full-enabled tests and unchanged
spec passes. Final logs are `tmp/parent90-*.log`; the three full mode-aware
comparisons are in the `parent130` logs above.

Separate disclosed platform/peripheral results: scanf's 78 deterministic
glibc probes differ from the macOS recording on 11 numeric-width cases
(exit 1; expectations unchanged). Island has 114/119 matching cases, five
case differences / seven failed checks: two huge-allocation cases and three
GNU-versus-BSD operand-parsing cases; all 17 upstream island gates pass for
both engines. These known exceptions are not stage-1 engine blockers and are
not described as passing tools. Genesis (131 cases, 13 gates) and mapgen
(126 comparisons, 11 gates) pass at their separate pins.

Open bugs rechecked after #130 closed: no local bug/follow-C issue is open.
G3 #309 (invalid RNG ranges, won't fix), #352 (unpaid NP transfer),
#353 (unowned-character knowledge crash), #354 (savage cap overshoot), #355
(missing token-vanishing notice) and #356 (log-only curse erosion) remain open.
Their existing C behavior is preserved; none is a known prerequisite fix for
stage 2. No unresolved stage-1 blocker is currently known.

## The opt-in switches

| switch | effect |
|---|---|
| `G3G0_FULL=1` | lets the whole-game Go tests (`TestBoxWorlds`, `TestBoxfileWorlds`, `TestCheckWorlds`, `TestLoadSaveWorlds`) run; `-short` always skips them. They compare nothing on their own: they need the game list that `box-compat`, `boxfile-compat`, `check-compat` or `loadsave-compat` hands them, and skip without it. |
| `--full` | `box-compat`, `boxfile-compat`, `check-compat`, `loadsave-compat`: also compare every saved game (`world.tgz` and each game cached by `olympia-roundtrip`). |
| `--full` | `turn-compat`: repeat pinned C world turns, then compare Go's complete cwd (paths/types/modes/empty directories/regular bytes, reports included) and raw stdout/stderr/status. Currently red until pending report tasks land. |
| `--games N` | `box-compat`, `boxfile-compat`, `check-compat`, `loadsave-compat`: compare `N` saved games, always the same ones: `world.tgz` and cached games spread evenly over the sorted list. |
| `--determinism-all` | `olympia-roundtrip`: check that C gives the same result twice on every selected case, not just the fixed sample. |

By default `box-compat`, `boxfile-compat`, `check-compat` and `loadsave-compat` run only their recorded cases, and `olympia-roundtrip` checks C's determinism only on its fixed sample of cases when it builds its cache.
The game comparisons delete each game's output as soon as it matches, and keep only failing games' output under the tool's `tmp/` directory.

## Time limits

The affected engine tools (`rng`, `entity`, `code`, `helper`, `io`, `box`, `boxfile`,
`check`, `loadsave`, `mapgen` and `olympia-cli` compat, plus `olympia-roundtrip`
and `olympia-spec`, `turn-compat`, `turn-driver-compat` and `command-compat`) run their C/Go processes under `run_limited` from
[`lib/timeout.sh`](lib/timeout.sh), with a limit of `RUN_TIMEOUT` seconds.
`move-compat` uses the same supervision for both C recordings and its Go test.
[`lib/run-limited.sh`](lib/run-limited.sh) adapts it for `go test -exec` and
unchanged upstream fixtures. `genesis-compat` and `island-compat` also use it.
At the limit it sends SIGTERM, then SIGKILL 2 seconds later, so a hung Go binary is stopped as surely as a C one ([#47](https://github.com/mdhender/g3g0/issues/47): the Go runtime ignores the SIGALRM the scripts used to rely on).
A run that times out has status 124, for both engines; any other status is the program's own, with 128 + N for death by signal N.
The scripts report a status of 124 as a timeout.

## Upstream files each tool reads

Paths are in the G3 repository.
A tool needs rerunning after a pin move if any file it reads changed; [`affected.sh`](affected.sh) applies this table to `git diff --name-only OLD NEW`.

| tool | records | reads |
|---|---|---|
| [`rng-compat`](rng-compat) | `internal/rng/testdata/c-output.txt` | `lib/rng.[ch]`, `lib/rnd.[ch]`, `lib/ilist.c`, `lib/lists.h`; line ranges of `olympia/` `buy.c`, `combat.c`, `npc.c`, `quest.c`, `cloud.c`, `faery.c`, `hades.c`, `tunnel.c`, `basic.c`, `c1.c`, `code.c`, `dir.c`, `day.c`, `seed.c`, `storm.c`, `swear.c`, `use.c`, `u.c` (any change can move a range) |
| [`lists-compat`](lists-compat) | `internal/lists/testdata/c-output.txt` | `lib/*.[ch]` |
| [`entity-compat`](entity-compat) | `internal/olympia/testdata/entity-c-output.txt` | `olympia/*.h`, `lib/*.h`, `olympia/glob.c`, `lib/checked_alloc.c` |
| [`code-compat`](code-compat) | `internal/olympia/testdata/code-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `code.c`, `glob.c`, `z.c`, `io.c`, `sout.c`; `lib/` `checked_alloc.c`, `ilist.c`, `rng.c`, `rnd.c` |
| [`helper-compat`](helper-compat) | `internal/olympia/testdata/helper-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `code.c`, `glob.c`, `z.c`, `sout.c`, `u.c`, `seed.c`, `loc.c`, `stack.c`, `swear.c`, `use.c`, `move.c`, `day.c`, `garr.c`, `order.c`, `input.c`; `lib/` `checked_alloc.c`, `ilist.c`, `rng.c`, `rnd.c`, `cstrings.c`, `item_ents.c`, `orders.c` |
| [`io-compat`](io-compat) | `internal/olympia/testdata/io-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `code.c`, `glob.c`, `z.c`, `io.c`, `sout.c`, `u.c`, `seed.c`; `lib/` `checked_alloc.c`, `ilist.c`, `rng.c`, `rnd.c`, `admits.c`, `item_ents.c`, `req_ents.c`, `skill_ents.c`, `trades.c` |
| [`box-compat`](box-compat) | `internal/olympia/testdata/box-c-output.txt` | what `io-compat` reads, and `olympia/input.c`, `lib/cstrings.c`, `tests/olympia/fixtures/world.tgz`; with `--full` or `--games`, the games `olympia-roundtrip` caches |
| [`boxfile-compat`](boxfile-compat) | `internal/olympia/testdata/boxfile-c-output.txt` | what `box-compat` reads, and `olympia/swear.c` |
| [`check-compat`](check-compat) | `internal/olympia/testdata/check-c-output.txt` | what `helper-compat` reads, and `olympia/check.c`, `olympia/io.c`, `lib/admits.c`, `lib/req_ents.c`, `lib/skill_ents.c`, `lib/trades.c`, `tests/olympia/fixtures/world.tgz`; with `--full` or `--games`, the games `olympia-roundtrip` caches |
| [`loadsave-compat`](loadsave-compat) | `internal/olympia/testdata/loadsave-c-output.txt` | what `check-compat` reads, and `olympia/buy.c`, `olympia/dir.c`, `olympia/add.c`, `olympia/quest.c`; with `--full` or `--games`, the games `olympia-roundtrip` caches |
| [`scanf-compat`](scanf-compat) | `internal/clib/testdata/scanf-c-output.txt` | nothing upstream: it tests the host's C library |
| [`olympia-cli-compat`](olympia-cli-compat) | `results.tsv` | the whole tree: it builds `olympia-g3` |
| [`olympia-roundtrip`](olympia-roundtrip) | nothing | the whole tree: it builds `olympia-g3` and runs the regress fixtures' saved games |
| [`olympia-spec`](olympia-spec) | `passing.txt` | the whole tree: it runs G3's spec |
| [`turn-compat`](turn-compat) | `testdata/` small real C snapshots and determinism manifest | the whole tree: it builds `olympia-g3` and turns `tests/olympia/fixtures/world.tgz` |
| [`turn-driver-compat`](turn-driver-compat) | `internal/olympia/testdata/turn-{driver,phase}-c-output.txt`, `{calendar,schedule}-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `u.c`, `input.c`, `day.c`, `z.c`, `storm.c`, `seed.c`, `glob.c`, `rnd.c`; `lib/` `ilist.c`, `rng.c` |
| [`move-compat`](move-compat) | `internal/olympia/testdata/move-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `dir.c`, `loc.c`, `u.c`, `swear.c`, `move.c`, `input.c`, `perm.c`, `stack.c`, `storm.c`, `garr.c`, `art.c`, `combat.c`, `day.c`, `code.c`, `glob.c`, `z.c`, `sout.c`, `display.c`, `use.c`, `io.c`, `seed.c`; `lib/` `checked_alloc.c`, `ilist.c`, `exit_views.c`, `rng.c`, `rnd.c`, `item_ents.c`, `admits.c`, `cstrings.c`, `skill_ents.c`, `trades.c`, `req_ents.c` |
| [`command-compat`](command-compat) | `internal/olympia/testdata/command-c-output.txt` | `olympia/*.h`, `lib/*.h`; `olympia/` `code.c`, `glob.c`, `z.c`, `sout.c`, `input.c`, `combat.c`, `art.c`, `display.c`, `order.c`, `swear.c`, `stack.c`, `loc.c`, `garr.c`, `build.c`, `u.c`, `dir.c`, `perm.c`, `storm.c`, `use.c`, `c1.c`, `produce.c`, `c2.c`, `immed.c`, `stealth.c`, `move.c`, `day.c`, `gm.c`, `seed.c`, `io.c`, `buy.c`, `lore.c`, `npc.c`, `hades.c`, `faery.c`, `savage.c`, `necro.c`, `report.c`, `summary.c`, `cloud.c`, `tunnel.c`, `quest.c`, `scry.c`, `add.c`, `main.c`; `lib/` `checked_alloc.c`, `ilist.c`, `exit_views.c`, `rng.c`, `rnd.c`, `cstrings.c`, `item_ents.c`, `skill_ents.c`, `orders.c`, `wait_args.c`, `flag_ents.c`, `admits.c`, `trades.c`, `req_ents.c` |

`genesis-compat`, `island-compat` and `mapgen-compat` have pins of their own and are not covered here.
