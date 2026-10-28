# Movement helper compatibility

Pin: `6a25f2d252ce45bd5adbe3baaccb3be3eed56279` in
[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3/tree/6a25f2d252ce45bd5adbe3baaccb3be3eed56279).

```sh
G3_REPO=/home/user/workspace/repo/tmp/upstream.git tools/move-compat/run.sh
# Only actual pinned C may replace the expected bytes:
G3_REPO=/home/user/workspace/repo/tmp/upstream.git tools/move-compat/run.sh --update
```

`G3_REPO` defaults to the workspace's bare `tmp/upstream.git`. `G3_REV`, `WORK`,
`CC` and `RUN_TIMEOUT` may be overridden. Source is read only via `git archive`
into `tmp/move-compat/g3`; nothing builds in an upstream working tree.

The harness extracts **complete original definitions**, under original `#line`
locations, for dir.c's topology/distance/selection closure; loc.c's region,
province and stack walks; u.c's depth/hidden/knowledge, weight/count and number
helpers; swear.c's ownership walks and loyalty setter; move.c's complete
land/fly/water travel calculations, parsing, land admission/outer feasibility,
exit save/restore, stack action/guard bookkeeping and `check_captain_loses_sailors`;
input.c's `find_command` using the original glob.c command table; u.c's real
`sub_item` for mid-voyage crew removal; perm.c's `rp_admit` and
`will_admit`; stack.c's actual leave/unstack closure; and storm.c's actual
`weather_here`. Garrison routing uses complete garr.c `garrison_here`
(41–52), `province_admin` (73–98), `players_who_rule_here` (168–201),
and day.c `touch_loc_pl` (2274–2294). It preserves dir.c's original
includes, tag constants, stream declaration, turn guard and `begin_exitdir`,
and u.c's original number table. It links original code.c/glob.c/z.c/**sout.c**
and the actual lib checked allocator, ilist, **exit_views.c**, rng.c, rnd.c and
item_ents.c, admits.c and cstrings.c. Initialization calls real `glob_init` and
`init_lower`, as main.c does. Original headers provide structs, macros and
integer field widths.
There are **no successful gameplay stubs**: unreachable command callbacks
(including SAIL) abort if reached, and the tool rejects their out-of-scope marker
even in assertion cases or `--update` mode. The driver only constructs fixtures,
invokes helpers and dumps observed state; output comes through real sout.

Each section runs in a fresh process, twice, under `run_limited`. Successful
cases must exit 0, assertion cases 134, and the separately identified known
G3 #353 case 139. Both runs must match raw stdout, stderr and status, then match
`internal/olympia/testdata/move-c-output.txt`. The Go recorded-case test uses
the absolute run-limited adapter. Default checks need no whole world.

Coverage:

- Asymmetric terrain distances, depth swaps, buildings, Hades pit, ocean links
  and signed sea-lane flags; missing and zero direction slots.
- Province/sublocation/building/ship/sibling-ship exits; ordering of directions,
  inner locations, links and roads; recursive port insertion and first-city choice.
- Hidden known/unknown destinations and roads, OUT exceptions, barriers,
  collapsed mines, regional/Hades metadata, and the distinct road rules.
- Vertical cloud routes, permanent Uldim blocks, seasonal Uldim/Summerbridge
  boundaries, and month-done behavior.
- Ship departure day/evening, differing sibling timestamps, negative and wrapped
  C-int timestamp arithmetic; completed and unfinished ships and dry-dock links.
- LAND/WATER/both/empty filters; sequential repeated shuffles, turn reseeding,
  day changes, empty/singleton no-draw controls and unrelated legacy-RNG draws.
- Original diagnostic/assertion bytes and assertion flushing. The known
  [G3 #353](https://github.com/playbymail/olympia-g3/issues/353) unowned-character
  NULL dereference retains SIGSEGV's unflushed output and status 139.
- All eight asymmetric weight totals, noble-item fallback, negative/zero/positive
  capacities, signed char/short fields, C-int products, inventory filtering,
  nested stacks/prisoners, fighter/man counts, ownership/NPC noble counting,
  and the `loop_all_here` noncharacter assertion.
- Land zero bypass, riding/equality and odd-distance truncation, all three
  difficult terrains, animals, load-capacity and double-capacity boundaries,
  percentage truncation, army denominator/threshold/cap and singular/plural
  messages, including actual report output and output-bit mutation when show=FALSE.
- Fly zero/one/three/four/seven/eight/nine distances, equality/overload and
  underground refusal. Calculations leave the command and exit untouched.
- Completed **uncrewed raft at distance seven**, wind off/on, retains base
  distance and emits no crew/wind message (G3 #350/#351). Roundships at distance
  seven and two with no/partial/full crew and wind off/on discriminate the
  roundship-only bonus. Galley no/partial/full crew, ignored slaves, overcrew,
  negative wind, building suppression and other-province/weather controls.
- State dumps retain every configured travel input, inventory and here-list
  order, and player output bits. Ordinary output is checked both immediately
  and through real sout's log-file route, without normalizing bytes.
- Parsing honors valid destination before literal direction, full/short and
  mixed-case directions, last-passable/last-impassable duplicate destination
  versus first direction match, hidden/known/see-all visibility, first IN and
  port-city choices, quiet mode, and exact already-there/unknown/absent diagnostics.
- Admission checks default garrisons before owner lookup, same owner, absent
  permissions and player allocation, allow/deny sense with unit/player entries,
  first matching list, sewer/OUT/no-owner exceptions, both refusal recipients,
  and water-before-impassable-before-transit-before-owner precedence.
- Outer feasibility uses real nested exits and real travel costs, excludes IN,
  preserves zero-distance overload bypass and same-depth zero results, and
  keeps army slowdown's real output/output-bit mutation even in quiet checks.
- Asymmetric command and exit save/restore preserves untouched command fields
  and zeroes every omitted exit field. Nested/prisoner suspend/restore checks
  epoch extrema/zero; guard clearing handles character and container roots and
  prunes noncharacter branches. Real leave/unstack checks output order, here-list
  promotion, summoned loyalty, release-magic clearing and descendant resume.
- Actual running SAIL command state on a completed raft, roundship and galley.
  Removing real sailors/pirates mid-voyage leaves raft wait/base water cost
  unchanged with no crew-loss output or log; ships incur actual singular/plural,
  incremental/capped penalties. Enough, partial and absent crew; NULL command,
  DONE/LOAD/ERROR, non-SAIL, nonship and unfinished ship controls; raft early
  return before penalty/wait assertions and the real ship assertion bytes.
  Dumps preserve command arguments/state/wait, ship/character moving fields,
  inventory and output vector. Real sout/log_write check immediate output,
  day-prefixed faction and GM logs, output bits, inform zero/self/other-player/
  same-player and prisoner suppression. No output text is synthesized by the driver.

## Location/garrison output routing (#109)

The original linked sout.c supplies complete `bottom_out` (288–362),
`can_view_loc` (376–404), `out_location` (407–429), `out_garrison` (432–445),
`out_sup` (448–482) and `vout` (485–492), including original format strings.
No recipient decision or output text is modeled in the C driver.

Thirteen new cases cover direct/outer touched versus untapped/known-only,
hidden known/unknown, city versus province view boundaries, nested buildings
and ships, real touching and its outer-bit mutation, and allocation of player
data even for factions denied output. Garrison controls check first-character
selection, absent/invalid/empty castles, ownership and pledge order with
duplicates/zero owners, silent recipients, visibility off/on, hidden sublocation
refusal even when known, hidden provinces, and the touched-hidden early return.
Prisoners bypass even malformed garrison routes; mixed vectors retain duplicates
and order. All eight day/secondary-indent/unit framing combinations use signed,
asymmetric values. Real faction logs and output bits distinguish location IDs
from OUT_GARR (6), and preserve two literal trailing spaces in routing messages.
Sparse touched/known/output lists are dumped in actual insertion order, with
nonascending touched and output controls rather than sorted membership only.
Assertions retain the actual file/line/condition, including eager ruler lookup
before a touched-location shortcut and the 5000-iteration pledge guard.

The recording has **83 sections / 4067 lines / 81524 bytes**, from two
deterministic bounded C runs. The previous **70 sections / 3380 lines / 69134
bytes** are an unchanged prefix (including 59/2868, 40/1703 and 25/814).
Command compat adds serialized-state and real garrison auto-attack boundary
controls: **207 sections / 7193 lines / 137909 bytes**; its previous
**204 sections / 7002 lines / 135225 bytes** remain an unchanged prefix.
Move callbacks remain explicit NYI for #67; departure/arrival and prisoner
escape dependencies belong to the separately precreated #110–#113 tasks.
This harness accepts movement helpers, not whole movement execution or
arrival callbacks. Report rendering and unrelated LOS-distance helpers are not
extracted; report-path output mutations are checked, not skipped.

## Character-description prerequisites (#114)

Complete pinned definitions are extracted for garr.c `rank_s` (603–624),
art.c `has_auraculum` (7–18) and `max_eff_aura` (25–50), combat.c
`find_wield` (489–544) and `wield_s` (547–588). The original combat.c
`fight` and `wield` structures (110–146) are exported once into a generated
header shared by the driver and extracted functions, not re-created by the
driver. Original oly.h/garr.h macros, signed item-magic bonus fields,
loop.h inventory filtering and code.c box names remain authoritative.
Go's new art.go owns the original art.c aura responsibilities; no formatter,
monthly rank computation, movement messages or combat execution is ported here.

Sixteen additional sections cover every rank string, absent character data,
signed-byte rank wrapping and both invalid-rank assertions; absent, invalid,
unheld, zero/negative-quantity and held auracula; negative innate aura clamped
**before** signed auraculum and per-entry bonuses; quantity independence,
signed-short wrapping and inventory filtering. Equipment cases use independent
asymmetric maxima, nonascending inventory IDs, unequal quantities, equal ties,
negative/zero/wrapped signed-byte bonuses, all four NULL/non-NULL selection/stat
combinations and all single-category strings. All 19 fight fields are dumped:
only attack/defense/missile ratings change, not selected-item fields or other
sentinels. Attack/defense, attack/missile and all-three duplicates are suppressed;
the original defense/missile-only duplicate is intentionally retained.
Before/after state dumps include inventories in insertion order, character,
magic and item-magic allocation presence and fields, ownership, here lists and
output bits. Invalid inventory owners check actual art.c/combat.c assertion bytes;
the aura assertion uses an allocated deleted box so the earlier rp_magic macro
can run (an unallocated owner violates that macro's precondition).

The new recording is **99 sections / 4789 lines / 96680 bytes**, generated
only from two fresh bounded deterministic C executions, then compared with Go.
The entire old move **83-section input** and **4067-line / 81524-byte output**
remain exact prefixes. Command input and its **207-section / 7193-line /
137909-byte output** are unchanged. Recorded trailing whitespace is preserved.
No new successful dependency stubs or expected-success NYI cases are introduced.

## Complete character and stack descriptions (#115)

The same harness now extracts complete display.c definitions at their original
source locations: `with_inventory_string` (169–194), `display_with` (197–205),
`mage_s` (336–358), `priest_s` (361–369), `emperor` (372–399),
`liner_desc_char` (406–477), `liner_desc` (528–544) and `show_chars_below`
(561–576), plus the original display globals (402–403). Linked sout.c supplies
the actual `wiout` (506–519), including its reset of `second_indent` to **zero**,
not to the caller's previous value. The full original `liner_desc` switch is
retained: its four noncharacter helpers abort with the rejected unfinished
dependency marker, never invented successful strings. Go retains explicit NYI
for those branches. A temporary expected-abort noncharacter case was rejected
in both normal and `--update` mode without changing the recording; it was
removed after the check. Go tests retain checks of all four NYI boundaries.

Real priest lookup uses complete use.c `exp_level` (957–974), `rp_skill_ent`
(997–1013) and `has_skill` (1145–1156), with original lib/skill_ents.c. Existing
Go `hasSkill`, `rankS`, `maxEffAura`, `wieldS`, name/quantity helpers and
depth-first character walks are reused, not reimplemented in the driver.
Serialization uses the linked original io.c `save_box` (2073–2123) and its
actual print closure, lib/trades.c and req_ents.c, and seed.c `int_comp`
(386–393). No serializer output is modeled. Output-bit insertion order is
dumped separately because these bits are deliberately not saved to the DB.

Seventeen appended sections discriminate complete suffix order (rank, aura,
religion, Emperor, three independent equipment choices, guard/combat flags,
banner, inventory, prisoner); ordinary and NPC names/counts including signed
C-int wrap; noble-item exclusion, nonascending inventory order, prominence,
quantity and invalid-entry filtering, and real missing-plural warnings. Every
mage threshold has upper-endpoint and next-value controls, with absent/hidden/
nonmagician controls and actual effective-aura bonuses. Priest controls retain
known versus learning/unknown/unrelated skills and signed field wrapping.
Emperor checks absent/zero throne owner, unowned/NPC/regular factions,
wrong province/building/unfinished castle and changed unique owner. Leader-only
guard, zero/nonzero signed behind and empty/nonempty ally suffixes differ.
Banner controls cover absent, empty, truncation and the display toggle.

Stack cases distinguish raw nonempty here lists from filtered character walks:
an item-only list is accompanied but emits no descendant, and its character
child is pruned. Nonascending nested characters are shown depth first,
including hidden units and prisoners, without repeating the root. Actual
immediate output, faction log bytes, location recipient order, vector duplicates,
prisoner suppression, day/unit/primary/secondary framing, output bits and
secondary reset are recorded. Empty walks retain the previous secondary indent;
even a suppressed `wiout` resets it. Before/after complete state and real
serialized character/item/player bytes check that descriptions do not mutate
the database. Assertions retain the actual display.c/loc.c/u.c file, line,
condition and status 134. Prior SIGSEGV/no-flush status 139 remains unchanged.

The recording is now **116 sections / 5944 lines / 119914 bytes**, generated
only by two fresh bounded deterministic C runs and compared exactly with Go.
All prior **99-section input bytes** and **4789 lines / 96680 output bytes**
are unchanged prefixes. Command's input and **207 sections / 7193 lines /
137909 output bytes** are unchanged. Literal trailing spaces, including the
empty-banner fixture and actual log messages, are intentional.

This accepts only character formatting and stack display. Movement messages
remain #110; report assembly, monthly rank computation, noncharacter
descriptions and full MOVE/turn acceptance remain explicitly unfinished.
