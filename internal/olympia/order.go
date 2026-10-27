// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/order.c: the parts of the order queues that
// get_command() uses (rp_order_head(), top_order() and pop_order()). The
// rest of order.c, with load_orders() and save_orders(), is left for the
// load_db and save_db task.

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
