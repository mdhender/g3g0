// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

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
// helpers from move.c. Departure and arrivals are separate tasks (#67).
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
// has already removed qty sailors/pirates. SAIL callbacks remain NYI (#67).
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
