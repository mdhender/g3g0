// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/swear.c: lord() and player(), which write_master()
// uses, and set_loyal(), set_lord() and is_unit(), which check_db() uses.

// swear.c:344–396. Start validation is deliberately stricter than completion.
func (w *world) vBribe(c *command) int {
	who, target, amount := int(c.who), int(c.a), int(c.b)
	if w.hasSkill(who, skBribeNoble) < 1 {
		w.commandOut(who, "BRIBE requires knowledge of %s.", capitalize(w.boxName(skBribeNoble)))
		return 0
	}
	if !w.checkCharHere(who, target) {
		return 0
	}
	if p := w.rpChar(target); p != nil && p.newLord != 0 {
		w.commandOut(who, "%s just switched employers this month, and is not looking for a new one so soon.", w.boxName(target))
		return 0
	}
	if w.isNPC(target) {
		w.commandOut(who, "NPC's cannot be bribed.")
		return 0
	}
	if w.player(target) == w.player(who) {
		w.commandOut(who, "%s already belongs to our faction.", w.boxName(target))
		return 0
	}
	if amount == 0 {
		w.commandOut(who, "Must specify an amount of gold to use as a bribe.")
		return 0
	}
	if !w.canPay(who, amount) {
		w.commandOut(who, "Don't have %s for a bribe.", goldS(amount))
		return 0
	}
	w.commandOut(who, "Attempt to bribe %s with a gift of %s.", w.boxName(target), goldS(amount))
	return 1
}

// swear.c:529–551. Gift flavor is keyed per turn, unlike the daily outcome.
func (w *world) thanksForGift(who, target int) {
	switch w.streams.SocGift(who, target) {
	case 1:
		w.commandOut(who, "%s graciously accepts our gift.", w.boxName(target))
	case 2:
		w.commandOut(who, "%s thanks us for the gift.", w.boxName(target))
	case 3:
		w.commandOut(who, "%s pockets the gold.", w.boxName(target))
	default:
		assert(false, "olympia/swear.c", 549, "FALSE")
	}
}

// swear.c:573–679. Do not repeat start's skill, NPC or same-faction checks.
// Bribes are intentionally a gold sink, including POCKET; only REPORT refunds:
// https://github.com/playbymail/olympia-g3/issues/159
func (w *world) dBribe(c *command) int {
	who, target, amount := int(c.who), int(c.a), int(c.b)
	if !w.checkStillHere(who, target) {
		return 0
	}
	if p := w.rpChar(target); p != nil && p.newLord != 0 {
		w.commandOut(who, "%s just switched employers this month, and is not looking for a new one so soon.", w.boxName(target))
		return 0
	}
	if !w.charge(who, amount) {
		w.commandOut(who, "Don't have %s for a bribe.", goldS(amount))
		return 0
	}
	var lk int8
	var rate int32
	if p := w.rpChar(target); p != nil {
		lk, rate = p.loyKind, p.loyRate
	}
	threshold := 0
	switch lk {
	case loyUnsworn, loyContract:
		threshold = max(int(rate), 250)
	case loyFear:
		threshold = 250
	case loyOath:
	default:
		assert(false, "olympia/swear.c", 614, "FALSE")
	}
	const (
		switchLord = iota + 1
		pocket
		report
		headForHills
	)
	outcome := pocket
	if threshold <= 0 || amount < threshold {
		if w.streams.SocBribe(who, target, 0, 1, 2) != 1 {
			outcome = report
		}
	} else {
		n := w.streams.SocBribe(who, target, 1, 1, 100)
		switch {
		case n <= 35:
			outcome = switchLord
		case n <= 65:
			outcome = pocket
		case n <= 90:
			outcome = report
		default:
			outcome = headForHills
		}
	}
	if outcome == switchLord && !w.enoughNPToAcquire(who, target) {
		outcome = pocket
	}
	switch outcome {
	case switchLord:
		w.commandOut(who, "%s accepts the gift, and has decided to join us.", w.boxName(target))
		w.unitDeserts(target, w.player(who), true, loyContract, 250)
		w.pChar(target).freshHire = 1
		if c.c != 0 {
			w.joinStack(target, who)
		}
	case headForHills:
		w.thanksForGift(who, target)
		w.commandOut(who, "%s left the service of %s, but didn't join us.", w.boxName(target), w.boxName(w.player(target)))
		w.unitDeserts(target, int(w.indepPlayer), true, loyUnsworn, 0)
	case pocket:
		w.thanksForGift(who, target)
	case report:
		w.thanksForGift(who, target)
		w.genItem(who, itemGold, amount)
		w.commandOut(target, "%s tried to bribe us with %s.", w.boxName(who), goldS(amount))
	default:
		assert(false, "olympia/swear.c", 675, "FALSE")
	}
	return 1
}

// swear.c:973–1082. RAISE checks availability at start, but allocates from
// the actor's current sublocation at completion, through the shared cookie owner.
func (w *world) vRaise(c *command) int {
	where := w.subloc(int(c.who))
	if !w.checkSkill(int(c.who), skRaiseMob) {
		return 0
	}
	if !w.mayCookieNPC(int(c.who), where, itemMobCookie) {
		return 0
	}
	return 1
}

func (w *world) dRaise(c *command) int {
	who := int(c.who)
	where := w.subloc(who)
	mob := w.doCookieNPC(who, where, itemMobCookie, where)
	if mob <= 0 {
		w.logWrite(10, "d_raise mob <= 0")
		w.commandOut(who, "Failed to raise peasant mob.")
		return 0
	}
	w.addSkillExperience(who, skRaiseMob)
	w.queue(mob, "guard 1")
	w.initLoadSup(mob)
	w.commandOut(who, "Raised %s.", w.boxName(mob))
	w.commandOut(where, "A speech by %s has raised %s.", w.boxName(who), w.linerDesc(mob))
	return 1
}

func (w *world) vRally(c *command) int {
	who, mob := int(c.who), int(c.a)
	if !w.checkSkill(who, skRallyMob) {
		return 0
	}
	if !w.checkCharHere(who, mob) {
		return 0
	}
	if w.nobleItem(mob) != itemPeasant && w.nobleItem(mob) != itemAngryPeasant {
		w.commandOut(who, "%s is not a peasant mob.", w.boxName(mob))
		return 0
	}
	return 1
}

func (w *world) dRally(c *command) int {
	who, mob := int(c.who), int(c.a)
	if !w.checkCharGone(who, mob) {
		return 0
	}
	if w.nobleItem(mob) != itemPeasant && w.nobleItem(mob) != itemAngryPeasant {
		w.commandOut(who, "%s is not a peasant mob.", w.boxName(mob))
		return 0
	}
	w.addSkillExperience(who, skRallyMob)
	if n := w.stackParent(mob); n != 0 {
		var rate int32
		if p := w.rpChar(mob); p != nil {
			rate = p.loyRate
		}
		w.setLoyal(mob, loySummon, min(int(rate+3), 5))
		w.commandOut(who, "Renewed enthusiasm of %s for %s.", w.boxName(mob), w.boxName(n))
		w.commandOut(who, "The peasants will stay spirited for %d months.", w.pChar(mob).loyRate)
		return 1
	}
	w.joinStack(mob, who)
	w.setLoyal(mob, loySummon, 3)
	w.flushUnitOrders(w.player(mob), mob)
	w.interruptOrder(mob)
	// C deliberately returns FALSE after joining, with all side effects retained.
	return 0
}

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
