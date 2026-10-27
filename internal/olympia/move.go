// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

// moveExitLand, moveExitFly and moveExitWater are the complete calculation
// helpers from move.c. Admission, departure and arrivals are separate tasks.
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
