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
`v_hide`, `d_hide`, `v_wait`, `d_wait`, `i_wait`. Real dependencies include order-head/top/pop operations,
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
unused display formatting (`is_port_city`) and validated-trade settlement. Reaching any
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
complete a persisted MOVE without calling unported movement start/finish code.
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
