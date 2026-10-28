# Command dispatch, scheduler, waits, interruptions and ownership transfers

Pinned reference: [Olympia G3 at
6a25f2d252ce45bd5adbe3baaccb3be3eed56279](https://github.com/playbymail/olympia-g3/tree/6a25f2d252ce45bd5adbe3baaccb3be3eed56279).

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
