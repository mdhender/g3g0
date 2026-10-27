// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// Expected bytes are recorded only by tools/turn-driver-compat/run.sh --update,
// using real pinned C. Default tests need neither C nor a whole saved world.
func TestTurnDriverMatchesC(t *testing.T) {
	runTurnDriverCases(t, "turn-driver")
}

// Only the compat tool's build overlay may run instrumented phase callbacks.
func TestTurnPhaseMatchesC(t *testing.T) {
	if os.Getenv("G3G0_PHASE_TRACE") != "1" {
		t.Skip("isolated callbacks require tools/turn-driver-compat/run.sh overlay")
	}
	runTurnDriverCases(t, "turn-phase")
}

func runTurnDriverCases(t *testing.T, prefix string) {
	t.Helper()
	var out, stdout, stderr strings.Builder
	var w *world
	finish := func() {
		if w != nil {
			fmt.Fprintf(&out, "stdout:\n%sstderr:\n%s\nexit 0\n", stdout.String(), stderr.String())
		}
	}
	for line := range strings.SplitSeq(readTestdata(t, prefix+"-cases.txt"), "\n") {
		if line == "" || line[0] == '#' {
			continue
		}
		if strings.HasPrefix(line, "section ") {
			finish()
			stdout.Reset()
			stderr.Reset()
			w = newWorld()
			w.stdout, w.stderr = &stdout, &stderr
			if prefix == "turn-phase" {
				for n := 200; n <= 202; n++ {
					w.allocBox(n, tPlayer, subPlRegular)
				}
				w.pPlayer(200).sentOrders = 7
				w.pPlayer(202).sentOrders = -5
			}
			fmt.Fprintln(&stdout, line)
			continue
		}
		fmt.Fprintf(&stdout, "> %s\n", line)
		c := &caseLine{rest: line}
		switch op := c.word(); op {
		case "clock":
			w.sysclock = olytime{turn: int16(c.num(t)), day: int16(c.num(t)), daysSinceEpoch: int32(c.num(t))}
		case "increment":
			for range c.num(t) {
				olytimeIncrement(&w.sysclock)
			}
		case "turn_change":
			olytimeTurnChange(&w.sysclock)
		case "dump":
			fmt.Fprintf(&stdout, "clock: %d %d %d; month %d year %d\n", w.sysclock.turn,
				w.sysclock.day, w.sysclock.daysSinceEpoch, olyMonth(w.sysclock), olyYear(w.sysclock))
		case "dot":
			n, ch := c.num(t), c.num(t)
			for range n {
				w.printDot(ch)
			}
		case "time_self":
			w.timeSelf = c.num(t) != 0
		case "stage":
			w.stage(strings.TrimPrefix(c.rest, " "))
		case "stage_null":
			w.stageSup(nil)
		case "stage_empty":
			w.stage("")
		case "process_orders":
			w.processOrders()
		case "daily_events":
			w.dailyEvents()
		case "post_month":
			w.postMonth()
		case "auto_quit":
			w.autoQuitTurns = int32(c.num(t))
		case "state":
			fmt.Fprintf(&stdout, "state: %d %d %d wait=%d auto=%d month=%d post=%d\n",
				w.sysclock.turn, w.sysclock.day, w.sysclock.daysSinceEpoch,
				w.cmdWait, b2i(w.autoAttackFlag), w.monthDone, w.postHasBeenRun)
			fmt.Fprintf(&stdout, "sent: %d %d\n", w.rpPlayer(200).sentOrders, w.rpPlayer(202).sentOrders)
		default:
			t.Fatalf("unknown operation %q", op)
		}
	}
	finish()
	want := readTestdata(t, prefix+"-c-output.txt")
	if got := out.String(); got != want {
		// Quote the streams so whitespace and missing final newlines show.
		t.Fatalf("C recording mismatch:\n got %q\nwant %q", got, want)
	}
}

func (w *world) phaseTrace(name string) {
	fmt.Fprintf(w.stdout, "phase %s: %d %d %d auto=%d month=%d post=%d\n", name,
		w.sysclock.turn, w.sysclock.day, w.sysclock.daysSinceEpoch,
		b2i(w.autoAttackFlag), w.monthDone, w.postHasBeenRun)
}

func (w *world) phaseFindCommand(name string) int {
	w.phaseTrace(name)
	return 42
}

func (w *world) phaseCalDay(pick, low, high int) int {
	fmt.Fprintf(w.stdout, "cal_day %d %d %d\n", pick, low, high)
	if pick == 0 {
		return 7
	}
	return 15
}

func (w *world) phaseWthrShuffle(days []int) {
	w.phaseTrace("wthr_shuffle")
	copy(days, []int{15, 7, 1, 3})
	j := 4
	for i := 1; i <= monthDays; i++ {
		if i != 15 && i != 7 && i != 1 && i != 3 {
			days[j] = i
			j++
		}
	}
}
