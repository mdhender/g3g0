// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/swear.c: lord() and player(), which write_master()
// uses, and set_loyal(), set_lord() and is_unit(), which check_db() uses.

// vSwear is the active unconditional refusal in swear.c:216–287. The
// fealty code below its #if 1 return is unreachable in the pinned engine.
func (w *world) vSwear(c *command) int {
	w.commandOut(int(c.who), "The SWEAR order has been removed from the game.")
	return 0
}

// vHonor is the complete swear.c:682–718 start (also HONOUR). In
// particular, negative raw amounts reach charge's actual subItem assertion;
// only zero and oath loyalty are rejected before that boundary.
func (w *world) vHonor(c *command) int {
	who, amount := int(c.who), int(c.a)
	if amount == 0 {
		w.commandOut(who, "Must specify an amount of gold to use as a gift.")
		return 0
	}
	if p := w.rpChar(who); p != nil && p.loyKind == loyOath {
		w.commandOut(who, "%s graciously declines the offer.", w.boxName(who))
		return 0
	}
	if !w.charge(who, amount) {
		w.commandOut(who, "Do not have %s.", goldS(amount))
		return 0
	}
	if p := w.rpChar(who); p == nil || p.loyKind != loyContract {
		w.pChar(who).loyKind = loyContract
		w.pChar(who).loyRate = 0
	}
	w.pChar(who).loyRate += int32(amount)
	w.commandOut(who, "%s now bound with %s.", w.boxName(who), w.loyalS(who))
	return 1
}

// vOath is complete swear.c:721–791, with the original active #if 1
// NP rule. Clamp only the local flag, not c.a; oath-1 -> oath-2 costs one
// extra NP. NP reads do not allocate player records, but refusal output can
// allocate the owner through the shared output machinery.
func (w *world) vOath(c *command) int {
	who, flag := int(c.who), int(c.a)
	pl := w.player(who)
	var lk int8
	var lr int32
	if p := w.rpChar(who); p != nil {
		lk, lr = p.loyKind, p.loyRate
	}
	if flag > 2 {
		flag = 2
	}
	if flag < 1 {
		flag = 1
	}
	if lk == loyOath && lr >= 2 {
		w.commandOut(who, "%s already is at %s, the maximum loyalty.", w.boxName(who), w.loyalS(who))
		return 0
	}
	if flag == 2 && lk == loyOath && lr == 1 {
		flag = 1
	}
	npCost := flag
	assert(npCost > 0, "olympia/swear.c", 762, "np_cost > 0")
	np := 0
	if p := w.rpPlayer(pl); p != nil {
		np = int(p.noblePoints)
	}
	if np < 1 {
		w.commandOut(who, "Player %s has no Noble Points.", w.boxCode(pl))
		return 0
	}
	if np < npCost {
		w.commandOut(who, "Player %s only has %d Noble Points.", w.boxCode(pl), np)
		npCost = np
	}
	if lk != loyOath {
		w.pChar(who).loyKind = loyOath
		w.pChar(who).loyRate = 0
	}
	w.pChar(who).loyRate += int32(npCost)
	w.deductNP(pl, npCost)
	w.commandOut(who, "%s now bound with %s.", w.boxName(who), w.loyalS(who))
	return 1
}

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

// npToAcquire is swear.c:187–196. Only the independent unit's previous
// owner is exempt; ordinary acquisitions still use the complete NP total.
func (w *world) npToAcquire(who, target int) int {
	if w.player(target) == int(w.indepPlayer) && int(w.pChar(target).prevLord) == w.player(who) {
		return 0
	}
	return w.charNPTotal(target)
}

func (w *world) enoughNPToAcquire(who, target int) bool {
	nps := w.npToAcquire(who, target)
	np := 0
	if p := w.rpPlayer(w.player(who)); p != nil {
		np = int(p.noblePoints)
	}
	if np < nps {
		w.commandOut(who, "Don't have %d NP%s to take control of %s.", nps, addS(nps), w.boxName(target))
		return false
	}
	return true
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
