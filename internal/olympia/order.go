// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"os"

	"github.com/mdhender/g3g0/internal/clib"
)

// Ported from olympia/order.c: the parts of the order queues that
// get_command() uses (rp_order_head(), top_order() and pop_order()), and
// load_orders() and save_orders(), which keep the queues in libdir/orders,
// one file per player, with what they call. The rest of order.c is left
// for later stages.

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
