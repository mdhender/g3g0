// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"time"

	"github.com/mdhender/g3g0/internal/clib"
)

// Complete c2.c:563–583. publicTurn is unsaved monthly state; the first
// call alone allocates the record and generates CLAIM gold for the faction.
func (w *world) vPublic(c *command) int {
	pl := w.player(int(c.who))
	if p := w.rpPlayer(pl); p != nil && p.publicTurn != 0 {
		w.commandOut(int(c.who), "Already a public turn.")
		return 0
	}
	w.pPlayer(pl).publicTurn = 1
	w.genItem(pl, itemGold, 100)
	w.commandOut(int(c.who), "Received 100 CLAIM gold.")
	return 1
}

// Complete c2.c:11–43; validation precedes quantity assertions and dropItem
// owns consumption, unique ownership, ocean routing and token side effects.
func (w *world) vDiscard(c *command) int {
	who, item := int(c.who), int(c.a)
	if w.kind(item) != tItem {
		w.commandOut(who, "%s is not an item.", w.boxCode(item))
		return 0
	}
	if w.hasItem(who, item) < 1 {
		w.commandOut(who, "%s does not have any %s.", w.boxName(who), w.boxCode(item))
		return 0
	}
	qty := w.howMany(who, who, item, int(c.b), int(c.c))
	if qty <= 0 {
		return 0
	}
	ret := w.dropItem(who, item, qty)
	assert(ret, "olympia/c2.c", 38, "ret")
	w.commandOut(who, "Dropped.")
	return 1
}

// Complete c2.c:761–819. The common mistake correction changes local values,
// not the saved command or parse fields. Zero realm IDs still compare equal.
func (w *world) vClaim(c *command) int {
	who, item, qty := int(c.who), int(c.a), int(c.b)
	pl := w.player(who)
	if w.region(who) == int(w.cloudRegion) || w.region(who) == int(w.hadesRegion) || w.region(who) == int(w.faeryRegion) {
		w.commandOut(who, "CLAIM may not be used in the Cloudlands, Hades or Faery.")
		return 0
	}
	if c.parse.Len()-1 < 2 && (w.kind(item) != tItem || w.hasItem(pl, item) < 1) && qty == 0 {
		w.logWrite(10, "correcting CLAIM for %s:  %s", boxCodeLess(w.player(who)), clib.Str(c.line))
		w.commandOut(who, "(assuming you meant CLAIM %d %d)", itemGold, item)
		qty, item = item, itemGold
	}
	if w.kind(item) != tItem {
		w.commandOut(who, "%s is not an item.", w.boxCode(item))
		return 0
	}
	if w.hasItem(pl, item) < 1 {
		w.commandOut(who, "No %s for you to claim.", w.boxCode(item))
		return 0
	}
	qty = w.howMany(who, pl, item, qty, int(c.c))
	if qty <= 0 {
		return 0
	}
	ret := w.moveItem(pl, who, item, qty)
	assert(ret, "olympia/c2.c", 815, "ret")
	w.commandOut(who, "Claimed %s.", w.justNameQty(item, qty))
	return 1
}

// Minimal command handlers, ported whole from c2.c. In particular these
// report settings persist on the player and control report formatting.
func (w *world) vFormat(c *command) int {
	pl := w.player(int(c.who))
	w.pPlayer(pl).format = int8(c.a)
	w.commandOut(int(c.who), "Formatting for %s set to %d.", w.boxName(pl), c.a)
	return 1
}

func (w *world) vNotab(c *command) int {
	pl := w.player(int(c.who))
	w.pPlayer(pl).notab = int8(c.a)
	if c.a != 0 {
		w.commandOut(int(c.who), "No TAB characters will appear in turn reports.")
	} else {
		w.commandOut(int(c.who), "TAB characters may appear in turn reports.")
	}
	return 1
}

// openTimes preserves main.c's file-open failures before process_orders.
// Headers are explicitly flushed before the next open, including on failure.
func (w *world) openTimes() {
	const dashLine = "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=\n"
	if w.pressFP == nil {
		fnam := sout("%s/times_press", w.libdir)
		var err error
		w.pressFP, err = clib.Fopen(fnam)
		if err != nil {
			fmt.Fprintf(w.stderr, "can't write %s", fnam)
			panic(exitStatus(1))
		}
	}
	fmt.Fprint(w.pressFP, "\n"+dashLine+"\n                        Player-contributed press\n\n"+dashLine+"\n")
	_ = w.pressFP.Flush()
	if w.rumorFP == nil {
		fnam := sout("%s/times_rumor", w.libdir)
		var err error
		w.rumorFP, err = clib.Fopen(fnam)
		if err != nil {
			fmt.Fprintf(w.stderr, "can't write %s", fnam)
			panic(exitStatus(1))
		}
	}
	fmt.Fprint(w.rumorFP, "                                 Rumors\n\n"+dashLine+"\n")
	_ = w.rumorFP.Flush()
}

// closeTimes is the complete c2.c:457–472 close/reset, press before rumor.
func (w *world) closeTimes() {
	if w.pressFP != nil {
		w.pressFP.Close()
		w.pressFP = nil
	}
	if w.rumorFP != nil {
		w.rumorFP.Close()
		w.rumorFP = nil
	}
}

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
	month := olyMonth(w.sysclock)
	if month < 0 {
		month = 0
	}
	assert(month >= 0 && month < numMonths, "olympia/c2.c", 434, "month >= 0 && month < NUM_MONTHS")
	issueS := fmt.Sprintf("%s year %d", monthNames[month], olyYear(w.sysclock)+1)

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
