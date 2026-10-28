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
