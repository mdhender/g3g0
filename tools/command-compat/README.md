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
unreachable garrison report routes (`players_who_rule_here`), unused display formatting
(`is_port_city`), and validated-trade settlement. Reaching any
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

The minimal `commandOut` route implements default character/player output,
prisoner suppression, real raw C line framing and LEN-1 truncation, append-only
log writes and the player's output-bit mutation. The small `logWrite` MASTER
route shares that same bottom writer; code alerts preserve GM output bits
and raw log bytes, not just diagnostics. Immediate output goes to
stdout with C framing, but without command echo; it does **not** bypass all
allow checks. Other routing/formatting fails explicitly. Report assembly and
the broader sout interface remain stage 2. No report-side mutations on these
paths are omitted. `-r` remains disabled; these scheduler methods are not full
turn orchestration. Wait callbacks, polling, generic order interruption and
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
For a nonprisoner garrison, stage-2 report fan-out remains unported: its selected
attack stops at the actual `auto_attack` entry, before `out`, with no fabricated
report output or player-bit changes. A separate production-boundary test covers
its no-target return versus actual attack stop; no C garrison-output stub is
recorded or blessed. Prisoners retain C's output suppression.
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
Whole-game checks remain opt-in in `turn-compat`. No cases consume RNG.
