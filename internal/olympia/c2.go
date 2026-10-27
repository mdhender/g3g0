// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"time"

	"github.com/mdhender/g3g0/internal/clib"
)

// Ported from olympia/c2.c: times_masthead(), which main() runs at the
// end of every run that loads the game.
//
// open_times() and close_times() are not ported: open_times() opens
// libdir/times_press and libdir/times_rumor only for a turn (-r), so the
// close_times() that main() runs after times_masthead() has nothing to
// close.

// months is months[] in c2.c, the calendar months of the masthead date.
var months = []string{
	"January", "February", "March", "April", "May", "June", "July",
	"August", "September", "October", "November", "December",
}

// timesMasthead is times_masthead(): the masthead of the Olympia Times,
// written to libdir/times_0. The date is today's, in local time, or
// January 1, 2000 with test-use-const-report-date. A file that cannot be
// written gets a diagnostic, with no newline, and the engine exits with
// status 1.
func (w *world) timesMasthead() {
	fnam := sout("%s/times_0", w.libdir)
	fp, err := clib.Fopen(fnam)
	if err != nil {
		fmt.Fprintf(w.stderr, "can't write %s", fnam)
		panic(exitStatus(1))
	}

	var date string
	if w.testUseConstReportDate {
		// fixed date for deterministic golden output (see main.c)
		date = "January 1, 2000"
	} else {
		tm := time.Now()
		date = fmt.Sprintf("%s %d, %d", months[tm.Month()-1], tm.Day(), tm.Year())
	}

	nplayers := 0
	for i := range w.loopKind(tPlayer) {
		if w.subkind(i) == subPlRegular {
			nplayers++
		}
	}

	plural := "s"
	if nplayers == 1 {
		plural = ""
	}
	turnS := fmt.Sprintf("turn %d, %d player%s", w.sysclock.turn, nplayers, plural)

	// oly_month() is ((turn-1) % NUM_MONTHS), which C clamps to 0 when it
	// is negative (turn 0 and below).
	month := (int(w.sysclock.turn) - 1) % numMonths
	if month < 0 {
		month = 0
	}
	assert(month >= 0 && month < numMonths, "olympia/c2.c", 434, "month >= 0 && month < NUM_MONTHS")
	issueS := fmt.Sprintf("%s year %d", monthNames[month], (int(w.sysclock.turn)-1)/numMonths+1)

	fmt.Fprintf(fp, "From: %s\n", w.fromHost)
	fmt.Fprintf(fp, "Subject: %s - The Olympia Times - Issue %d\n", w.gameTitle, w.sysclock.turn)
	fmt.Fprintf(fp, "\n")

	fmt.Fprintf(fp, "   +----------------------------------------------------------------------+\n")
	fmt.Fprintf(fp, "   | The Olympia Times %s |\n", clib.PadLeft(issueS, 50))
	fmt.Fprintf(fp, "   | %s |\n", clib.PadRight(date, 68))
	fmt.Fprintf(fp, "   |                                                                      |\n")
	fmt.Fprintf(fp, "   | %s %s |\n", turnS, padStar(w.fromHost, 67-len(turnS)))
	fmt.Fprintf(fp, "   +----------------------------------------------------------------------+\n\n")
	fmt.Fprintf(fp, "     Questions, comments, to play: %s\n\n", w.gameURL)
	fmt.Fprintf(fp, "                             Olympia PBEM\n\n")
	fmt.Fprintf(fp, "                                *  *  *\n\n")
	fmt.Fprintf(fp, "                                *  *  *\n\n")

	fp.Close()
}

// padStar is printf("%*s", width, s): s padded with spaces on the left to
// width bytes, or, for a negative width, on the right to -width bytes.
func padStar(s string, width int) string {
	if width < 0 {
		return clib.PadRight(s, -width)
	}
	return clib.PadLeft(s, width)
}
