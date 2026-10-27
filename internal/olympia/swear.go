// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/swear.c: lord() and player(), which write_master()
// uses, and set_loyal(), set_lord() and is_unit(), which check_db() uses.

// lord is lord(): who box n's owner is: n itself if it is a player,
// otherwise its unit_lord, or 0 if it has no entity_char.
func (w *world) lord(n int) int {
	if w.kind(n) == tPlayer {
		return n
	}

	if p := w.rpChar(n); p != nil {
		return int(p.unitLord)
	}
	return 0
}

// player is player(): the player box n belongs to, following its lords up
// to a player, or 0. A chain longer than 5000 lords (a loop) fails an
// assertion.
func (w *world) player(n int) int {
	loopCheck := 5000

	for n != 0 && w.kind(n) != tPlayer {
		n = w.lord(n)
		assert(loopCheck > 0, "olympia/swear.c", 32, "loop_check-- > 0")
		loopCheck--
	}

	return n
}

// setLoyal is set_loyal(): it sets who's loyalty kind (stored as a schar)
// and rate.
func (w *world) setLoyal(who, k, lev int) {
	p := w.pChar(who)

	p.loyKind = int8(k)
	p.loyRate = int32(lev)
}

// setLord is set_lord(): it swears who to newLord (0 for none), moving it
// to newLord's player's unit list, and sets its loyalty to k and lev unless
// k is LOY_UNCHANGED.
//
// Only the case check_db() reaches is ported: who belongs to no player.
// Moving a unit from one player to another flushes its orders, interrupts
// its command, clears its attitudes, reswears the units sworn to it, and
// moves noble points, all through code not ported yet, so a unit that
// belongs to a player panics. With no old player, the noble point exchange
// at the end cannot happen either.
func (w *world) setLord(who, newLord, k, lev int) {
	var newPl int

	oldPl := w.player(who)

	if newLord != 0 {
		newPl = w.player(newLord)
	} else {
		newPl = 0
	}

	if oldPl != 0 && oldPl != newPl {
		// flush_unit_orders, interrupt_order, clear_all_att,
		// reswear_all_sworn, the pledge, and the removal from the old
		// player's units.
		panic("olympia: set_lord: changing a unit's player is not yet implemented")
	}

	prevLord := int(w.pChar(who).prevLord)
	if prevLord == 0 {
		prevLord = oldPl
	}

	if newPl != 0 && oldPl != newPl {
		w.pChar(who).prevLord = int32(oldPl)
	}

	w.pChar(who).unitLord = int32(newLord)

	if newLord != 0 && newPl != oldPl {
		p := w.pPlayer(newPl)
		p.units.Append(who)

		safeQsort(p.units, intComp)

		w.initLoadSup(who) // load command from new owner
		w.touchLoc(who)
	}

	// player -> player	NP cost and return
	// player -> indep	no return, no cost
	// indep  -> player	charge cost, return to NP prev owner

	if oldPl == int(w.indepPlayer) {
		oldPl = prevLord
	}

	if oldPl != 0 && newPl != 0 && oldPl != newPl && newPl != int(w.indepPlayer) {
		// Not reached: oldPl is 0, or the panic above was. With it, the
		// noble points (char_np_total, add_np, deduct_np) and wout.
		panic("olympia: set_lord: noble point exchange is not yet implemented")
	}

	if k != loyUnchanged {
		w.setLoyal(who, k, lev)
	}
}

// isUnit is is_unit(): whether v is in the units list of player pl, which
// must be a player.
func (w *world) isUnit(pl, v int) bool {
	assert(w.kind(pl) == tPlayer, "olympia/swear.c", 297, "kind(pl) == T_player")

	return w.pPlayer(pl).units.Lookup(v) >= 0
}
