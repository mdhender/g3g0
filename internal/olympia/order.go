// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"os"

	"github.com/mdhender/g3g0/internal/clib"
)

// Ported from olympia/order.c: the parts of the order queues that
// get_command() uses, load_orders()/save_orders() and the complete order
// template renderers. Remaining command handlers belong to later stages.

// pOrderHead is p_order_head(): player pl's order queue for unit who,
// appended to the player's queues if it has none. pl gets an
// entity_player if it has none.
func (w *world) pOrderHead(pl, who int) *orderList {
	p := w.pPlayer(pl)

	for i := 0; i < p.orders.Len(); i++ {
		if int(p.orders[i].unit) == who {
			return p.orders[i]
		}
	}

	n := &orderList{unit: int32(who)}

	p.orders.Append(n)

	return n
}

// rpOrderHead is rp_order_head(): player pl's order queue for unit who, or
// nil if pl has no entity_player or no queue for who.
func (w *world) rpOrderHead(pl, who int) *orderList {
	p := w.rpPlayer(pl)

	if p == nil {
		return nil
	}

	for i := 0; i < p.orders.Len(); i++ {
		if int(p.orders[i].unit) == who {
			return p.orders[i]
		}
	}

	return nil
}

// topOrder is top_order(): the first order in player pl's queue for unit
// who, and true, or false if there is none.
func (w *world) topOrder(pl, who int) (string, bool) {
	p := w.rpOrderHead(pl, who)

	if p != nil && p.l.Len() > 0 {
		return p.l[0], true
	}

	return "", false
}

func isStopOrder(s string) bool {
	for len(s) > 0 && iswhite(s[0]) {
		s = s[1:]
	}
	i := 0
	for i < len(s) && !iswhite(s[i]) && i < lenMax {
		i++
	}
	// The pinned C fuzzy comparison uses the remainder, not the token.
	return iStrcmp(s[:i], "stop") == 0 || fuzzyStrcmp(s[i:], "stop")
}

func (w *world) stopOrder(pl, who int) bool {
	s, ok := w.topOrder(pl, who)
	return ok && isStopOrder(s)
}

func (w *world) flushUnitOrders(pl, who int) {
	if w.bx[who] == nil {
		return
	}
	for {
		if _, ok := w.topOrder(pl, who); !ok {
			break
		}
		w.popOrder(pl, who)
	}
	if w.player(who) == pl {
		if c := w.rpCommand(who); c != nil && c.state == stateLoad {
			w.commandDone(c)
		}
	}
}

// popOrder is pop_order(): it removes the first order in player pl's queue
// for unit who, which must have one.
func (w *world) popOrder(pl, who int) {
	p := w.rpOrderHead(pl, who)

	assert(p != nil, "olympia/order.c", 130, "p != NULL")
	assert(p.l.Len() > 0, "olympia/order.c", 131, "cstrings_len(p->l) > 0")

	p.l.Delete(0)
}

// queueOrder is queue_order(): it appends order s to player pl's queue for
// unit who, unless the queue already holds 250 orders.
func (w *world) queueOrder(pl, who int, s string) {
	p := w.pOrderHead(pl, who)

	if p.l.Len() >= 250 {
		return
	}

	p.l.Append(s)
}

// prependOrder is order.c:174–181. Unlike queue_order there is no 250-order
// limit: implicit outer-level movement must reinsert the original order.
func (w *world) prependOrder(pl, who int, s string) {
	w.pOrderHead(pl, who).l.Prepend(s)
}

// queue is order.c:199–208, the formatted interface to queue_order.
func (w *world) queue(who int, format string, args ...any) {
	w.queueOrder(w.player(who), who, fmt.Sprintf(format, args...))
}

// savePlayerOrders is save_player_orders(): player pl's queued orders, as
// "unit:order" lines, queue by queue, to libdir/orders/pl. A queue whose
// unit is not a valid box, or is a dead character, is skipped. The file is
// created only when there is a line to write; one that cannot be gets a
// diagnostic and ends the player's orders.
func (w *world) savePlayerOrders(pl int) {
	var fp *clib.File

	assert(w.validBox(pl), "olympia/order.c", 219, "valid_box(pl)")

	p := w.rpPlayer(pl)

	if p == nil {
		return
	}

	for i := 0; i < p.orders.Len(); i++ {
		if !w.validBox(int(p.orders[i].unit)) ||
			w.kind(int(p.orders[i].unit)) == tDeadchar {
			continue
		}

		for j := 0; j < p.orders[i].l.Len(); j++ {
			if fp == nil {
				fnam := sout("%s/orders/%d", w.libdir, pl)
				var err error
				fp, err = clib.Fopen(fnam)

				if err != nil {
					fmt.Fprintf(w.stderr, "can't write %s: ",
						fnam)
					clib.Perror(w.stderr, "", err)
					return
				}
			}

			fmt.Fprintf(fp, "%d:%s\n",
				p.orders[i].unit,
				p.orders[i].l[j])
		}
	}

	if fp != nil {
		fp.Close()
	}
}

// loadPlayerOrders is load_player_orders(): it queues each line of
// libdir/orders/pl, read with readlin(), for the unit atoi() of the line
// gives, as the text after the line's first ':' (empty, if it has none).
// pl must have an entity_player with no queues yet: a box with none is a
// NULL dereference in C, and a panic here, and a second file for the same
// player fails the assertion (upstream defect:
// https://github.com/playbymail/olympia-g3/issues/332). A file that cannot
// be opened gets readfile()'s diagnostic.
func (w *world) loadPlayerOrders(pl int) {
	assert(w.validBox(pl), "olympia/order.c", 265, "valid_box(pl)")
	assert(w.rpPlayer(pl).orders == nil, "olympia/order.c", 266, "rp_player(pl)->orders == NULL")

	fnam := sout("%s/orders/%d", w.libdir, pl)

	if !w.readfile(fnam) {
		return
	}

	for line := w.readlin(); line != nil; line = w.readlin() {
		unit := clib.Atoi(*line)

		p := 0
		for p < len(*line) && (*line)[p] != ':' {
			p++
		}
		if p < len(*line) { // *p == ':'
			p++
		}

		w.queueOrder(pl, unit, (*line)[p:])
	}

	w.closefile(fnam)
}

// loadOrders is load_orders(): loadPlayerOrders for every file in
// libdir/orders whose name starts with a digit, in readdir()'s order, for
// the player atoi() of the name gives. A name that gives no valid box gets
// a diagnostic and is skipped. A directory that cannot be opened gets a
// diagnostic.
func (w *world) loadOrders() {
	fnam := sout("%s/orders", w.libdir)
	d, err := clib.Opendir(fnam)

	if err != nil {
		fmt.Fprintf(w.stderr, "can't open %s: ", fnam)
		clib.Perror(w.stderr, "", err)
		return
	}

	for {
		name, ok := d.Readdir()
		if !ok {
			break
		}

		if isdigit(name[0]) {
			fact := clib.Atoi(name)

			if !w.validBox(fact) {
				fmt.Fprintf(w.stderr,
					"ERROR: orders/%d but no box [%d]\n",
					fact, fact)
				continue
			}

			if w.rpPlayer(fact) == nil {
				fmt.Fprintf(w.stderr, "ERROR: orders/%s but [%d] is not a player\n", name, fact)
				continue
			}
			if w.rpPlayer(fact).orders != nil {
				fmt.Fprintf(w.stderr, "ERROR: orders/%s but orders for [%d] are already loaded\n", name, fact)
				continue
			}

			w.loadPlayerOrders(fact)
		}
	}

	d.Closedir()
}

// saveOrders is save_orders(): it removes libdir/orders with removeDir,
// makes it again (mode 0755 less the umask), ignoring failure, and saves
// every player's orders, in chain order.
func (w *world) saveOrders() {
	removeDir(sout("%s/orders", w.libdir))
	_ = os.Mkdir(sout("%s/orders", w.libdir), 0o755)

	for i := range w.loopKind(tPlayer) {
		w.savePlayerOrders(i)
	}
}

// ordersTemplateSup is the complete order.c:370–454 template body. Commands
// are only read here; rendering never loads, executes or consumes an order.
func (w *world) ordersTemplateSup(who, num, pl int) {
	nam := ""
	if pl == w.player(num) {
		if w.isPrisoner(num) || w.kind(num) != tChar {
			nam = sout("  # %s", w.justName(num))
		} else {
			nam = sout("  # %s in %s", w.justName(num), w.boxName(w.subloc(num)))
		}
	}
	w.commandOut(who, "unit %s%s", boxCodeLess(num), nam)
	w.indent += 3
	if pl == w.player(num) {
		c := w.rpCommand(num)
		if w.kind(num) == tPlayer && w.autoQuitTurns > 0 {
			p := w.pPlayer(num)
			turns := int(int32(w.sysclock.turn) - p.lastOrderTurn)
			if turns >= int(w.autoQuitTurns/2) {
				w.commandOut(who, "#")
				w.commandOut(who, "# You have not sent in orders for %d turn%s.", turns, addS(turns))
				w.commandOut(who, "# If you do not submit orders for %d turn%s,", w.autoQuitTurns, addS(int(w.autoQuitTurns)))
				w.commandOut(who, "# your faction will be removed from the game.")
				w.commandOut(who, "#")
			}
		}
		if p := w.rpChar(num); p != nil && ((p.loyKind == loyContract && p.loyRate < 50) ||
			(p.loyKind == loyFear && p.loyRate < 2)) {
			w.commandOut(who, "#")
			w.commandOut(who, "# %s has loyalty %s and will renounce", boxCodeLess(num), w.loyalS(num))
			w.commandOut(who, "# loyalty at the end of this turn.")
			w.commandOut(who, "#")
		}
		if c != nil && (c.state == stateRun || c.state == stateLoad) {
			if c.state == stateRun {
				timeLeft := " (still~executing)"
				if c.wait >= 0 {
					timeLeft = sout(" (executing for %s more day%s)", niceNum(int(c.wait)), addS(int(c.wait)))
				}
				w.commandOut(who, "# > %s%s", clib.Str(c.line), timeLeft)
			} else {
				w.commandOut(who, "%s", clib.Str(c.line))
			}
		}
	}
	if l := w.rpOrderHead(pl, num); l != nil {
		for _, s := range l.l {
			w.commandOut(who, "%s", s)
		}
	}
	w.indent -= 3
	w.commandOut(who, "")
}

// ordersOther is order.c:457–505, including active foreign commands with
// empty queues (which cause a template, but are not themselves printed).
func (w *world) ordersOther(who, pl int) {
	assert(w.kind(pl) == tPlayer, "olympia/order.c", 464, "kind(pl) == T_player")
	p := w.rpPlayer(pl)
	if p == nil {
		return
	}
	first := true
	for _, q := range p.orders {
		num := int(q.unit)
		if pl == num || !w.validBox(num) || w.kind(num) == tDeadchar || p.units.Lookup(num) >= 0 {
			continue
		}
		c, l := w.rpCommand(num), w.rpOrderHead(pl, num)
		if (c == nil || c.state == stateDone) && (l == nil || l.l.Len() == 0) {
			continue
		}
		if first {
			w.commandOut(who, "#")
			w.commandOut(who, "# Orders for units you do not control as of now")
			w.commandOut(who, "#")
			first = false
		}
		w.ordersTemplateSup(who, num, pl)
	}
}

// ordersTemplate is order.c:508–547. Keep password lookup before the first
// output, since output may allocate the absent player record.
func (w *world) ordersTemplate(who, pl int) {
	pass := ""
	if p := w.rpPlayer(pl); p != nil && p.password != nil && *p.password != "" {
		pass = sout(" \"%s\"", *p.password)
	}
	w.commandOut(who, "")
	prefix := ""
	if p := w.rpPlayer(pl); p != nil && p.brokenMailer != 0 {
		prefix = " "
	}
	w.commandOut(who, "%sbegin %s%s  # %s", prefix, boxCodeLess(pl), pass, w.justName(pl))
	if w.player(who) == pl && w.kind(who) == tChar {
		if p := w.rpChar(who); p != nil && p.loyKind == loyContract && p.loyRate < 50 {
			w.commandOut(who, "")
			w.commandOut(who, "# WARNING -- loyalty is %s", w.loyalS(who))
		}
	}
	w.commandOut(who, "")
	w.ordersTemplateSup(who, pl, pl)
	for i := range w.loopUnits(pl) {
		w.ordersTemplateSup(who, i, pl)
	}
	w.ordersOther(who, pl)
	w.commandOut(who, "end")
}

func (w *world) listOrderTemplates() {
	w.outPath, w.outAltWho = master, 5 // OUT_TEMPLATE
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) != subPlSystem && w.subkind(pl) != subPlSilent {
			w.ordersTemplate(pl, pl)
		}
	}
	w.outPath, w.outAltWho = 0, 0
}
