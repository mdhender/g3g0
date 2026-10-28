// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"strings"
)

// discoverRoad is move.c:140–162. allCharHere is C's snapshot walk,
// including prisoners; each reverse road marks both ends for the stack.
func (w *world) discoverRoad(who, where int, v *exitView) {
	for _, reverse := range w.exitsFromLoc(who, int(v.destination)) {
		if reverse.road != 0 && int(reverse.destination) == where {
			w.setKnown(who, int(reverse.road))
			w.setKnown(who, int(v.road))
			for _, j := range w.allCharHere(who) {
				w.setKnown(j, int(reverse.road))
				w.setKnown(j, int(v.road))
			}
		}
	}
}

// vMove is the complete move.c:584–689 start callback. In particular payment
// precedes leaving the old stack, while escapes precede departure messages.
func (w *world) vMove(c *command) int {
	who := int(c.who)
	where, checkOuter := w.subloc(who), true
	var v *exitView
	if c.parse.Len()-1 < 1 {
		w.commandOut(who, "Specify direction or destination to MOVE.")
		return 0
	}
	for c.parse.Len()-1 > 0 {
		move := "move"
		v = w.parseExitDir(c, where, &move)
		if v != nil {
			checkOuter = false
			if w.landCheck(c, v, true) {
				break
			}
			v = nil
		}
		if c.parse.Len()-1 == 1 {
			break
		}
		w.cmdShift(c)
	}
	if v == nil && checkOuter && w.canMoveAtOuterLevel(where, c) != 0 {
		c.a = int32(w.subloc(where))
		v = w.parseExitDir(c, where, nil)
		assert(v != nil, "olympia/move.c", 625, "v")
		if v != nil {
			assert(w.moveExitLand(c, v, false) >= 0, "olympia/move.c", 629, "move_exit_land(c, v, FALSE) >= 0")
			w.commandOut(who, "(assuming 'move out' first)")
			w.prependOrder(w.player(who), who, *c.line)
		}
	}
	if v == nil {
		return 0
	}
	delay := w.moveExitLand(c, v, true)
	if delay < 0 {
		return 0
	}
	if v.hadesCost != 0 {
		cost := v.hadesCost * w.countStackFigures(who)
		w.logWrite(11, "%s (%s) tries to enter Hades", w.boxName(w.player(who)), w.boxName(who))
		if !w.autocharge(who, int(cost)) {
			w.commandOut(who, "Can't afford %s to enter Hades.", goldS(int(cost)))
			return 0
		}
		w.commandOut(who, "The Gatekeeper Spirit of Hades took %s from us.", goldS(int(cost)))
	}
	v.distance, c.wait = delay, delay
	saveVArray(c, v)
	w.leaveStack(who)
	if delay > 0 {
		w.vectorStack(who, true)
		w.commandOut(vect, "Travel to %s will take %s day%s.", w.boxName(int(v.destination)), niceNum(int(delay)), addS(int(delay)))
	}
	w.suspendStackActions(who)
	w.clearGuardFlag(who)
	if delay > 1 {
		w.prisonerMovementEscapeCheck(who)
	}
	w.departureMessage(who, v)
	return 1
}

// dMove is move.c:776–832. The disabled #if 0 move_skip block is not active
// gameplay. Missing destinations fail before restoring moving flags, as C does.
func (w *world) dMove(c *command) int {
	var v exitView
	restoreVArray(c, &v)
	who, dest := int(c.who), int(v.destination)
	if !w.validBox(dest) {
		w.commandOut(who, "Your destination no longer exists!")
		return 0
	}
	if v.road != 0 {
		w.discoverRoad(who, w.subloc(who), &v)
	}
	w.vectorStack(who, true)
	w.commandOut(vect, "Arrival at %s.", w.boxName(dest))
	if w.locDepth(dest) == locProvince && w.viewloc(w.subloc(who)) != w.viewloc(dest) && w.weatherHere(dest, subFog) != 0 {
		w.commandOut(vect, "The province is blanketed in fog.")
	}
	w.restoreStackActions(who)
	w.moveStack(who, dest)
	if w.viewloc(int(v.orig)) != w.viewloc(dest) {
		w.arrivalMessage(who, &v)
	}
	return 1
}

// These six synonyms retain the original reparsing (move.c:1096–1172):
// olyParse changes c.cmd to MOVE, whose finish callback then handles arrival.
func (w *world) vExit(c *command) int {
	s := "move out"
	assert(w.olyParse(c, &s), "olympia/move.c", 1102, "ret")
	return w.vMove(c)
}

func (w *world) vEnter(c *command) int {
	if c.parse.Len()-1 < 1 {
		s := "move in"
		assert(w.olyParse(c, &s), "olympia/move.c", 1116, "ret")
		return w.vMove(c)
	}
	s := sout("move %s", c.parse[1])
	assert(w.olyParse(c, &s), "olympia/move.c", 1121, "ret")
	return w.vMove(c)
}

func (w *world) vNorth(c *command) int {
	s := "move north"
	assert(w.olyParse(c, &s), "olympia/move.c", 1133, "ret")
	return w.vMove(c)
}

func (w *world) vSouth(c *command) int {
	s := "move south"
	assert(w.olyParse(c, &s), "olympia/move.c", 1145, "ret")
	return w.vMove(c)
}

func (w *world) vEast(c *command) int {
	s := "move east"
	assert(w.olyParse(c, &s), "olympia/move.c", 1157, "ret")
	return w.vMove(c)
}

func (w *world) vWest(c *command) int {
	s := "move west"
	assert(w.olyParse(c, &s), "olympia/move.c", 1169, "ret")
	return w.vMove(c)
}

// touchLocAfterMove is move.c:692–707. where is deliberately unused by C:
// touching observes each character's actual sublocation after set_where.
func (w *world) touchLocAfterMove(who, where int) {
	if w.kind(who) == tChar {
		w.touchLoc(who)
	}
	for _, i := range w.allCharHere(who) {
		if !w.isPrisoner(i) {
			w.touchLoc(i)
		}
	}
}

// clearContacts is stealth.c:285–299: unlike the other arrival walks this
// includes prisoners, but does nothing for a container root.
func (w *world) clearContacts(stack int) {
	if w.kind(stack) == tChar {
		for _, i := range w.allStack(stack) {
			w.pChar(i).contact.Clear()
		}
	}
}

// moveStack is the complete move.c:710–773 arrival state transition. MOVE
// callbacks and their messages remain separate; do not reorder these hooks.
func (w *world) moveStack(who, where int) {
	assert(w.kind(who) == tChar, "olympia/move.c", 714, "kind(who) == T_char")
	if w.region(w.subloc(who)) != int(w.faeryRegion) && w.region(where) == int(w.faeryRegion) {
		w.logWrite(11, "%s enters Faery at %s.", w.boxName(who), w.boxName(where))
	}
	w.setWhere(who, where)
	w.markLocStackKnown(who, where)
	w.touchLocAfterMove(who, where)
	w.updateWeatherViewLocs(who, where)
	w.clearContacts(who)
	if w.subkind(where) == subCity {
		for _, i := range w.allStack(who) {
			w.matchTrades(i)
		}
	}
	if w.subkind(where) == subOcean {
		p := w.pChar(who)
		if p.timeFlying == 0 {
			p.timeFlying++
			w.oceanChars.Append(who)
		}
	}
	if w.subkind(where) != subOcean {
		if p := w.rpChar(who); p != nil && p.timeFlying != 0 {
			p.timeFlying = 0
			w.oceanChars.RemValue(who)
		}
	}
	if w.locDepth(where) == locProvince && w.subkind(where) != subOcean && w.region(where) == int(w.faeryRegion) {
		w.faeryAttackCheck(who, where)
	}
	if w.locDepth(where) == locProvince && w.subkind(where) != subOcean && w.region(where) == int(w.hadesRegion) {
		w.hadesAttackCheck(who, where)
	}
}

// departureMessage is the complete move.c:11–73 departure_message. Early
// returns preserve showToGarrison; a completed message resets it to false.
func (w *world) departureMessage(who int, v *exitView) {
	assert(w.validBox(who), "olympia/move.c", 19, "valid_box(who)")
	if w.charReallyHidden(who) {
		return
	}
	orig, dest := int(v.orig), int(v.destination)
	if w.locDepth(orig) == locProvince && w.weatherHere(orig, subFog) != 0 {
		return
	}
	if v.destHidden != 0 {
		return
	}
	desc := w.linerDesc(who)
	var to string
	if w.subloc(dest) == orig {
		to = sout(" entered %s", w.boxName(dest))
	} else if w.subloc(orig) == dest {
		to = sout(" exited %s", w.boxName(orig))
	} else if w.viewloc(orig) != w.viewloc(dest) {
		if v.direction >= 1 && v.direction <= 4 {
			to = sout(" went %s", fullDirS[v.direction])
		} else {
			to = sout(" left for %s", w.boxName(dest))
		}
	} else {
		return
	}
	with := w.displayWith(who)
	comma := ""
	if strings.Contains(desc, ",") {
		comma = ","
	}
	if with == "" {
		with = "."
	}
	if w.viewloc(orig) != w.viewloc(dest) {
		if garr := w.garrisonHere(orig); garr != 0 && w.garrisonNotices(garr, who) {
			w.showToGarrison = true
		}
	}
	w.commandOut(orig, "%s%s%s%s", desc, comma, to, with)
	w.showCharsBelow(orig, who)
	w.showToGarrison = false
}

// arrivalMessage is the complete move.c:76–133 arrival_message. Unlike
// departure, hidden origins suppress only the source phrase, not the message.
func (w *world) arrivalMessage(who int, v *exitView) {
	if w.charReallyHidden(who) {
		return
	}
	orig, dest := int(v.orig), int(v.destination)
	if w.locDepth(dest) == locProvince && w.weatherHere(dest, subFog) != 0 {
		return
	}
	desc := w.linerDesc(who)
	from := ""
	if v.origHidden == 0 {
		if v.direction >= 1 && v.direction <= 4 {
			from = sout(" from the %s", fullDirS[exitOpposite[v.direction]])
		} else {
			from = sout(" from %s", w.boxName(orig))
		}
	}
	with := w.displayWith(who)
	comma := ""
	if strings.Contains(desc, ",") {
		comma = ","
	}
	if with == "" {
		with = "."
	}
	if w.viewloc(orig) != w.viewloc(dest) {
		if garr := w.garrisonHere(dest); garr != 0 {
			if w.garrisonNotices(garr, who) {
				w.showToGarrison = true
			}
			if w.garrisonSpotCheck(garr, who) {
				w.indent += 3
				w.commandOut(garr, "%s%s", desc, with)
				w.showCharsBelow(garr, who)
				w.indent -= 3
			}
		}
	}
	w.commandOut(dest, "%s%s arrived%s%s", desc, comma, from, with)
	w.showCharsBelow(dest, who)
	w.showToGarrison = false
}

// parseExitDir is the complete parse_exit_dir from pinned move.c. A nil
// zeroArg suppresses diagnostics, not exit construction or its side effects.
func (w *world) parseExitDir(c *command, where int, zeroArg *string) *exitView {
	l := w.exitsFromLoc(int(c.who), where)
	if w.validBox(int(c.a)) {
		if where == int(c.a) {
			if zeroArg != nil {
				w.commandOut(int(c.who), "Already in %s.", w.boxName(where))
			}
			return nil
		}
		var ret, impassRet *exitView
		for _, v := range l {
			if v.destination == c.a && (v.hidden == 0 || w.immedSeeAll != 0) {
				if v.impassable != 0 {
					impassRet = v
				} else {
					ret = v
				}
			}
		}
		if ret != nil {
			return ret
		}
		if impassRet != nil {
			return impassRet
		}
		if zeroArg != nil {
			w.commandOut(int(c.who), "No visible route from %s to %s.", w.boxName(where), c.parse[1])
		}
		return nil
	}
	dir := lookup(fullDirS, c.parse[1])
	if dir < 0 {
		dir = lookup(shortDirS, c.parse[1])
	}
	if dir < 0 {
		if zeroArg != nil {
			w.commandOut(int(c.who), "Unknown direction or destination '%s'.", c.parse[1])
		}
		return nil
	}
	for _, v := range l {
		if int(v.direction) == dir && (v.hidden == 0 || w.immedSeeAll != 0) {
			if dir == dirIn && zeroArg != nil {
				w.commandOut(int(c.who), "(assuming '%s %s')", *zeroArg, boxCodeLess(int(v.destination)))
			}
			return v
		}
	}
	if zeroArg != nil {
		w.commandOut(int(c.who), "No visible %s route from %s.", fullDirS[dir], w.boxName(where))
	}
	return nil
}

func saveVArray(c *command, v *exitView) {
	c.b, c.c, c.d, c.e = v.direction, v.destination, v.road, v.destHidden
	c.f, c.g, c.h = v.distance, v.orig, v.origHidden
}

func restoreVArray(c *command, v *exitView) {
	*v = exitView{}
	v.direction, v.destination, v.road, v.destHidden = c.b, c.c, c.d, c.e
	v.distance, v.orig, v.origHidden = c.f, c.g, c.h
}

func (w *world) suspendStackActions(who int) {
	for i := range w.loopStack(who) {
		w.pChar(i).moving = w.sysclock.daysSinceEpoch
	}
}

func (w *world) clearGuardFlag(who int) {
	if w.kind(who) == tChar {
		w.pChar(who).guard = 0
	}
	for _, i := range w.allCharHere(who) {
		w.pChar(i).guard = 0
	}
}

func (w *world) landCheck(c *command, v *exitView, show bool) bool {
	who := int(c.who)
	if v.water != 0 {
		if show {
			w.commandOut(who, "A sea-worthy ship is required for travel across water.")
		}
		return false
	}
	if v.impassable != 0 {
		if show {
			w.commandOut(who, "That route is impassable.")
		}
		return false
	}
	if v.inTransit != 0 {
		if show {
			w.commandOut(who, "%s is underway.  Boarding is not possible.", w.boxName(int(v.destination)))
		}
		return false
	}
	if w.locDepth(int(v.destination)) == locBuild && w.subkind(int(v.destination)) != subSewer {
		owner := w.buildingOwner(int(v.destination))
		if owner != 0 && !w.willAdmit(owner, who, int(v.destination)) && v.direction != dirOut {
			if show {
				w.commandOut(who, "%s refused to let us enter.", w.boxName(owner))
				w.commandOut(owner, "Refused to let %s enter.", w.boxName(who))
			}
			return false
		}
	}
	return true
}

func (w *world) canMoveHere(where int, c *command) bool {
	v := w.parseExitDir(c, where, nil)
	return v != nil && v.direction != dirIn && w.landCheck(c, v, false) && w.moveExitLand(c, v, false) >= 0
}

func (w *world) canMoveAtOuterLevel(where int, c *command) int {
	for outer := w.subloc(where); w.locDepth(outer) > locRegion; outer = w.subloc(outer) {
		if w.canMoveHere(outer, c) {
			return w.locDepth(outer) - w.locDepth(where)
		}
	}
	return 0
}

// moveExitLand, moveExitFly and moveExitWater are the complete calculation
// helpers from move.c, shared by movement callbacks and direct comparisons.
func (w *world) moveExitLand(c *command, v *exitView, show bool) int32 {
	who := int(c.who)
	delay := int32(v.distance)
	if delay == 0 {
		return 0
	}
	terr := w.subkind(int(v.destination))
	swamp := terr == subSwamp || terr == subBog || terr == subPits
	u := w.determineStackWeights(who)
	if delay > 1 && u.rideCap >= u.rideWeight && !swamp {
		delay -= delay / 2
	} else {
		if u.landWeight > u.landCap*2 {
			if show {
				w.commandOut(who, "%s is too overloaded to travel.", w.boxName(who))
			}
			return -1
		}
		if swamp && u.animals != 0 {
			if show {
				w.commandOut(who, "Difficult terrain slows the animals.  Travel will take an extra day.")
			}
			delay++
		}
		if u.landWeight > u.landCap {
			ratio := (u.landWeight - u.landCap) * 100 / u.landCap
			additional := delay * ratio / 100
			if show {
				if additional == 1 {
					w.commandOut(who, "Excess inventory slows movement.  Travel will take an extra day.")
				} else if additional > 1 {
					w.commandOut(who, "Excess inventory slows movement.  Travel will take an extra %s days.", niceNum(int(additional)))
				}
			}
			delay += additional
		}
	}
	nobles := w.countStackMoveNobles(who)
	men := w.countStackFigures(who) - nobles
	assert(nobles > 0, "olympia/move.c", 332, "nobles > 0")
	extra := men / (nobles * w.armySlowFactor)
	if extra > int32(v.distance)*2 {
		extra = int32(v.distance) * 2
	}
	if extra == 1 {
		w.commandOut(who, "%s noble%s, %s %s: travel will take an extra day.", capitalize(niceNum(int(nobles))), addS(int(nobles)), niceNum(int(men)), manOrMen(men))
	} else if extra > 1 {
		w.commandOut(who, "%s noble%s, %s %s: travel will take an extra %d days.", capitalize(niceNum(int(nobles))), addS(int(nobles)), niceNum(int(men)), manOrMen(men), extra)
	}
	return delay + extra
}

func manOrMen(n int32) string {
	if n == 1 {
		return "man"
	}
	return "men"
}

func (w *world) moveExitFly(c *command, v *exitView, show bool) int32 {
	who := int(c.who)
	delay := int32(v.distance)
	if w.subkind(int(v.destination)) == subUnder || w.subkind(int(v.destination)) == subTunnel {
		if show {
			w.commandOut(who, "Cannot fly underground.")
		}
		return -1
	}
	if delay < 8 {
		if delay > 3 {
			delay = 3
		}
	} else {
		delay = 4
	}
	u := w.determineStackWeights(who)
	if u.flyCap < u.flyWeight {
		if show {
			w.commandOut(who, "%s is too overloaded to fly.", w.boxName(who))
		}
		return -1
	}
	nobles := w.countStackMoveNobles(who)
	men := w.countStackFigures(who) - nobles
	assert(nobles > 0, "olympia/move.c", 401, "nobles > 0")
	extra := men / (nobles * w.armySlowFactor)
	if extra > int32(v.distance)*2 {
		extra = int32(v.distance) * 2
	}
	if extra == 1 {
		w.commandOut(who, "%s noble%s, %s %s: travel will take an extra day.", capitalize(niceNum(int(nobles))), addS(int(nobles)), niceNum(int(men)), manOrMen(men))
	} else if extra > 1 {
		w.commandOut(who, "%s noble%s, %s %s: travel will take an extra %d days.", capitalize(niceNum(int(nobles))), addS(int(nobles)), niceNum(int(men)), manOrMen(men), extra)
	}
	return delay + extra
}

// checkCaptainLosesSailors is the complete pinned move.c helper. The caller
// has already removed qty sailors/pirates. Full SAIL callbacks remain NYI.
func (w *world) checkCaptainLosesSailors(qty int32, target, inform int) {
	where := w.subloc(target)
	if w.cmdSail < 0 {
		w.cmdSail = int32(w.findCommand("sail"))
		assert(w.cmdSail > 0, "olympia/move.c", 1189, "cmd_sail > 0")
	}
	if !w.isShip(w.subloc(target)) {
		return
	}
	c := w.rpCommand(target)
	if c == nil || c.state != stateRun || c.cmd != w.cmdSail {
		return
	}
	now := int32(w.hasItem(target, itemSailor)) + int32(w.hasItem(target, itemPirate))
	before := now + qty
	var shouldHave int32
	switch w.subkind(where) {
	case subRaft:
		// G3 #350/#351: a completed raft needs no crew, even mid-voyage.
		return
	case subGalley:
		shouldHave = 14
	case subRoundship:
		shouldHave = 8
	default:
		fmt.Fprintf(w.stderr, "kind is %d\n", w.subkind(where))
		assert(false, "olympia/move.c", 1218, "FALSE")
	}
	if now >= shouldHave {
		return
	}
	if before > shouldHave {
		before = shouldHave
	}
	penalty := before - now
	assert(penalty > 0, "olympia/move.c", 1229, "penalty > 0")
	w.outVector.Clear()
	w.outVector.Append(target)
	if inform != 0 && target != inform {
		w.outVector.Append(inform)
	}
	if penalty == 1 {
		w.commandOut(vect, "Loss of crew will cause travel to take an extra day.")
	} else {
		w.commandOut(vect, "Loss of crew will cause travel to take %s extra days.", niceNum(int(penalty)))
	}
	assert(c.wait > 0, "olympia/move.c", 1242, "c->wait > 0")
	c.wait += penalty
	w.logWrite(11, "Loss of sailors incurs penalty for %s.", w.boxCode(w.player(target))) // LOG_SPECIAL
}

func (w *world) moveExitWater(c *command, v *exitView, ship int, show bool) int32 {
	who := int(c.who)
	delay, handsShort, windBonus := int32(v.distance), int32(0), int32(0)
	where := w.subloc(ship)
	switch w.subkind(ship) {
	case subRaft:
		// G3 #350/#351: no crew penalty, and wind remains roundship-only.
	case subRoundship, subGalley:
		n := int32(w.hasItem(who, itemSailor)) + int32(w.hasItem(who, itemPirate))
		need, crew := int32(8), "roundship is eight sailors"
		if w.subkind(ship) == subGalley {
			need, crew = 14, "galley is fourteen slaves or sailors"
		}
		if n < need {
			handsShort = need - n
			s := "day"
			if handsShort != 1 {
				s = sout("%s days", niceNum(int(handsShort)))
			}
			if show {
				have := "none"
				if n != 0 {
					have = niceNum(int(n))
				}
				w.commandOut(who, "The crew of a %s, but you have %s.  Travel will take an extra %s.", crew, have, s)
			}
		}
	default:
		fmt.Fprintf(w.stderr, "subkind = %d\n", w.subkind(ship))
		assert(false, "olympia/move.c", 1309, "FALSE")
	}
	if w.subkind(ship) == subRoundship && w.weatherHere(where, subWind) != 0 && delay > 2 {
		windBonus = 1
		if show {
			w.commandOut(who, "Favorable winds speed our progress.")
		}
	}
	return delay + handsShort - windBonus
}

// Ported from olympia/move.c: restore_stack_actions(), which check_db()
// uses, and init_ocean_chars(), which load_db() runs. The rest of move.c is
// left for later stages.

// restoreStackActions is restore_stack_actions(): it clears the moving day
// of who and every unit stacked under it.
func (w *world) restoreStackActions(who int) {
	for i := range w.loopStack(who) {
		w.pChar(i).moving = 0
	}
}

// initOceanChars is init_ocean_chars(): every character whose subloc is an
// ocean, in chain order, is appended to ocean_chars, for
// check_ocean_chars() to watch.
func (w *world) initOceanChars() {
	for i := range w.loopKind(tChar) {
		where := w.subloc(i)

		if w.subkind(where) == subOcean {
			w.oceanChars.Append(i)
		}
	}
}

// checkOceanChars is move.c:852–898. The snapshot survives deaths and
// removals from oceanChars; even dead entries allocate char state before cleanup.
func (w *world) checkOceanChars() {
	for _, who := range w.oceanChars.Copy() {
		where := w.subloc(who)
		p := w.pChar(who)
		if w.kind(who) != tChar || w.subkind(where) != subOcean {
			p.timeFlying = 0
			w.oceanChars.RemValue(who)
			continue
		}
		p.timeFlying++
		if p.timeFlying <= 15 || w.stackParent(who) != 0 {
			continue
		}
		w.vectorStack(who, true)
		w.commandOut(vect, "Flight can no longer be maintained.  %s plunges into the sea.", w.boxName(who))
		w.killStackOcean(who)
	}
}
