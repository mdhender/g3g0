// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"strings"
	"testing"
)

// Expected bytes are recorded only by tools/turn-driver-compat/run.sh --update,
// using real pinned C. Default tests need neither C nor a whole saved world.
func TestTurnDriverMatchesC(t *testing.T) {
	var out, stdout, stderr strings.Builder
	var w *world
	finish := func() {
		if w != nil {
			fmt.Fprintf(&out, "stdout:\n%sstderr:\n%s\nexit 0\n", stdout.String(), stderr.String())
		}
	}
	for line := range strings.SplitSeq(readTestdata(t, "turn-driver-cases.txt"), "\n") {
		if line == "" || line[0] == '#' {
			continue
		}
		if strings.HasPrefix(line, "section ") {
			finish()
			stdout.Reset()
			stderr.Reset()
			w = newWorld()
			w.stdout, w.stderr = &stdout, &stderr
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
		default:
			t.Fatalf("unknown operation %q", op)
		}
	}
	finish()
	want := readTestdata(t, "turn-driver-c-output.txt")
	if got := out.String(); got != want {
		// Quote the streams so whitespace and missing final newlines show.
		t.Fatalf("C recording mismatch:\n got %q\nwant %q", got, want)
	}
}
