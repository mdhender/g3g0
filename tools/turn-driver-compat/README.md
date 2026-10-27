# Clock and turn-driver helpers

Pinned reference: [playbymail/olympia-g3 at
6a25f2d252ce45bd5adbe3baaccb3be3eed56279](https://github.com/playbymail/olympia-g3/tree/6a25f2d252ce45bd5adbe3baaccb3be3eed56279).

```sh
tools/turn-driver-compat/run.sh
tools/turn-driver-compat/run.sh --update  # record only real C output
```

Requires Bash, git, Python 3, a C11 compiler (`CC`, default `cc`), Go,
Perl (the shared timeout supervisor), cmp and diff. `G3_REPO` defaults to
`~/Software/playbymail/olympia-g3`; a bare cache containing the pin can also
be supplied. `G3_REV` overrides the pin only deliberately. `WORK` defaults
to `tmp/turn-driver-compat` and is emptied first. No upstream working tree or
build is used or modified.

The phase driver requires compiler section flags `-ffunction-sections` and
`-fdata-sections`, plus linker dead-code elimination: `-Wl,-dead_strip` on
Darwin, or `-Wl,--gc-sections` elsewhere, matching `command-compat`.

The script archives `olympia/` and `lib/` into `tmp/`, extracts the real
`olytime_increment`, `olytime_turn_change`, `print_dot` (including its static
initializer) and `stage` definitions from `u.c`, and retains their original
`#line` locations. The driver includes the exported headers, so `oly_month`
and `oly_year` are the actual `oly.h` macros, not rewritten C expectations.
The clock driver has no engine stubs. Both C runs must be deterministic and match the
recording byte for byte, then the Go test must match that recording. stdout,
stderr and exit status are captured separately and concatenated with labels;
no bytes are filtered, including progress whitespace and final newlines.
Each section is a fresh C process / Go world. All C runs and Go test binaries
run under `tools/lib/timeout.sh` via `run_limited` / `run-limited.sh`.

`TestTurnDriverMatchesC` uses the existing `caseLine` and `readTestdata`
helpers and runs in default Go tests without an upstream checkout. It reuses
the already-ported `stage` / `stageSup`. `timesMasthead` in `c2.go` now uses
the shared month/year helpers while retaining C's local negative-month clamp
and assertion; `olympia-roundtrip` covers its output. `-r` now enters the
actual driver under `runEngine` recovery, stopping at `queue_npc_orders`.
Whole-world turns remain opt-in in `turn-compat` and are still failures.

## Isolated orchestration callbacks (not gameplay acceptance)

The same script also extracts **unchanged real C bodies** of
`process_orders`, `daily_command_loop`, `daily_events`, `post_month` and
`clear_orders_sent`. `phase-driver.c` supplies tracing callbacks for their
unfinished world phases. The real clock, progress output, list append and
length, sorting, conditionals, flags and clear-orders mutation execute.
Calendar/shuffle callbacks deliberately collide on day 15 and supply weather
prefix `15,7,1,3`, followed by the remaining days; these are **instrumentation
inputs, not implementations or RNG expectations**. Coverage includes day-zero
load order, day-one trade matching before dots/commands, all thirty days,
weekly/timed event collisions, the four-weather-day cap, persistent statics,
seasonal openings/closures and positive/zero/negative auto-quit settings.

`prepare-phases.py` generates a **test-only Go build overlay** in `tmp/`:
only leaf calls in the four corresponding production orchestration bodies
become trace callbacks. The bodies themselves and real `clearOrdersSent`
are exercised, not hand-rewritten Go orchestration. Production has no phase
callback interface and never uses this overlay. `TestTurnPhaseMatchesC`
runs only through this tool. Scratch validation of the current real NYI
boundaries confirmed day zero changed turn 13/day 27/epoch 913 to turn 14/day
0/epoch 913, resolved WAIT, left completion flags unset, and cleared sent-orders
before the first unfinished settlement phase. That temporary test was removed:
the private NYI names are not lasting contracts as world phases are ported.
There are five callback sections / 484 raw recorded lines, in addition
to the unchanged seven clock sections / 152 lines. Both recordings come only
from two deterministic pinned C executions; stdout/stderr/status are unfiltered.

A successful callback comparison proves **call order only**, not NPC, calendar,
weather, economy, decay or gameplay completion. `turn-compat --full` runs the
unmodified actual engines and must continue to fail until the unfinished world
phases and main report-path mutations are implemented. At this task's boundary
Go prints `process_orders()` and stops with status 2 at `queue_npc_orders`,
before any save, instead of pretending the bare map completed a turn.

## Discriminating cases

- Nonzero turn **13**, day **27**, epoch **913**: neither field can be
  reconstructed from the others. Increment changes only day/epoch;
  turn change resets day to **0**, increments turn, and leaves epoch alone.
- Day **29 → 30 → 31**: increment does not silently roll a month/turn.
  A complete 30-day increment after day zero advances epoch by exactly 30.
- Turn **0 → 1**, **7 → 8 → 9**, **16 → 17**, **24 → 25**: distinguish
  `(turn-1)` from `turn`, eight-month years from other calendars, and
  zero-based macro results from the masthead's one-based display year.
- Negative turns and **-32768**: distinguish C's truncation/remainder from
  floor division and promotion before subtraction from short arithmetic.
  Day and turn **32767 → -32768** preserve short widths; epoch reaches
  **2147483647** without invoking C signed-int overflow.
- Character markers on dots **1, 58, 59, 60, 118, 119, 120** distinguish
  indentation and wrapping *before* character 59/119 from conventional
  60-column wrapping. Named, empty and NULL stages interleave without
  resetting the count. Character **321** checks fputc's low-byte conversion.
  A fresh section checks that a new world resets progress state.

The stage cases cover `time_self=0` and the first timed named/NULL call:
the former omits the newline and the latter prints `0 seconds`. Untimed calls
must not set the timer's old/first statics. Each timing section enables timing
for only one call, so these bytes do not depend on actual elapsed time; later
elapsed-time output is not recorded. The existing stage implementation is
unchanged. These routines consume no RNG and do not iterate entities.
Seven sections / 152 lines were
recorded on Linux x86-64 and pass the byte comparison.
