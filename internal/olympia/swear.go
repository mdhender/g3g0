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
func (w *world) setLord(who, newLord, k, lev int) {
	var newPl int
	flag := false

	oldPl := w.player(who)

	if newLord != 0 {
		newPl = w.player(newLord)
	} else {
		newPl = 0
	}

	if oldPl != 0 && oldPl != newPl {
		w.flushUnitOrders(oldPl, who)
		w.interruptOrder(who)
		w.clearAllAtt(who)
		w.reswearAllSworn(who)
		if p := w.rpMagic(who); p != nil && p.pledge != 0 {
			w.pMagic(who).pledge = 0
			w.pledgeBacklinks = false
		}
		w.pPlayer(oldPl).units.RemValue(who)
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
		if w.kind(oldPl) == tPlayer {
			nps := w.charNPTotal(who)
			w.addNP(oldPl, nps)
			w.commandOut(oldPl, "Received %d NP%s for %s.", nps, addS(nps), w.boxName(who))
		}
		flag = true
	}

	if k != loyUnchanged {
		w.setLoyal(who, k, lev)
	}
	if flag {
		nps := w.charNPTotal(who)
		if !w.deductNP(newPl, nps) {
			assert(false, "olympia/swear.c", 180, "FALSE")
		}
		w.commandOut(newPl, "Paid %d NP%s for %s.", nps, addS(nps), w.boxName(who))
	}
}

// Dependents stay with the old direct lord, before who's own lord changes.
func (w *world) reswearAllSworn(who int) {
	pl, newLord := w.player(who), w.lord(who)
	for i := range w.loopUnits(pl) {
		if w.lord(i) == who {
			w.setLord(i, newLord, loyUnchanged, 0)
		}
	}
}

// isUnit is is_unit(): whether v is in the units list of player pl, which
// must be a player.
func (w *world) isUnit(pl, v int) bool {
	assert(w.kind(pl) == tPlayer, "olympia/swear.c", 294, "kind(pl) == T_player")

	return w.pPlayer(pl).units.Lookup(v) >= 0
}

// unitDeserts is swear.c:300–340. loyCheck is unused in the original too.
func (w *world) unitDeserts(who, toWho int, loyCheck bool, k, lev int) {
	sp := w.player(who)
	if toWho != 0 && sp != 0 {
		w.commandOut(sp, "%s renounces loyalty to us.", w.boxName(who))
		w.commandOut(who, "%s renounces loyalty.", w.boxName(who))
		w.showDay = false
		w.commandOut(who, "")
		w.showDay = true
	}
	if toWho != 0 && w.isPrisoner(who) && w.player(toWho) == w.player(w.stackParent(who)) {
		w.pChar(who).prisoner = 0
	} else if !w.isPrisoner(who) {
		w.extractStackedUnit(who)
	}
	w.setLord(who, toWho, k, lev)
	if toWho != 0 {
		w.commandOut(who, "%s pledges fealty to us.", w.boxName(who))
		w.commandOut(toWho, "%s pledges fealty to us.", w.boxName(who))
		w.pChar(who).newLord = 1
	}
}
