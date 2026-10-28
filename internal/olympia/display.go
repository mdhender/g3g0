// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Location headers and character descriptions from display.c. Noncharacter
// descriptions and complete location-report assembly remain explicit boundaries.

// locInsideString is display.c:8–49. Nested buildings still name the province,
// not their immediate container; region names are raw names without box codes.
func (w *world) locInsideString(where int) string {
	if w.locDepth(where) == locBuild {
		if w.subkind(w.loc(where)) == subOcean {
			return sout(", in %s", w.boxName(w.province(where)))
		}
		if w.locDepth(w.loc(where)) == locProvince {
			return sout(", in province %s", w.boxName(w.province(where)))
		}
		return sout(", in %s", w.boxName(w.province(where)))
	} else if w.locDepth(where) == locSubloc {
		if w.subkind(w.province(where)) == subOcean {
			return sout(", in %s", w.boxName(w.province(where)))
		}
		return sout(", in province %s", w.boxName(w.province(where)))
	} else if w.subkind(where) == subOcean {
		if regName := w.name(w.region(where)); regName != "" {
			return sout(", in %s", regName)
		}
		return ""
	}
	if regName := w.name(w.region(where)); regName != "" {
		return sout(", in %s", regName)
	}
	return ""
}

func (w *world) showLocBarrier(who, where int) bool {
	if w.locBarrier(where) != 0 {
		w.commandOut(who, "A magical barrier surrounds %s.", w.boxName(where))
		return true
	}
	return false
}

func (w *world) safeHavenS(n int) string {
	if w.safeHaven(n) != 0 {
		return ", safe haven"
	}
	return ""
}

// locCivS is display.c:125–143, including the unguarded region comparisons
// in oly.h's in_faery/in_hades macros (zero special-region IDs also compare).
func (w *world) locCivS(where int) string {
	if w.locDepth(where) != locProvince || w.subkind(where) == subOcean ||
		w.region(where) == int(w.faeryRegion) || w.region(where) == int(w.hadesRegion) {
		return ""
	}
	var n int8
	if p := w.rpLoc(where); p != nil {
		n = p.civ
	}
	if n == 0 {
		return ", wilderness"
	}
	return sout(", civ-%d", n)
}

// showLocHeader is the complete display.c:147–166 suffix sequence.
func (w *world) showLocHeader(where int) string {
	buf := w.boxNameKind(where)
	buf += w.locInsideString(where)
	if w.mineDepth(where) != 0 {
		buf += sout(", depth~%d", w.mineDepth(where))
	}
	buf += w.safeHavenS(where)
	if w.locHidden(where) != 0 {
		buf += ", hidden"
	}
	if w.subkind(where) != subTunnel {
		buf += w.locCivS(where)
	}
	return sout("%s", buf)
}

// withInventoryString is display.c:169–194. loop_inv preserves insertion order
// and excludes invalid items and nonpositive quantities, not hidden items.
func (w *world) withInventoryString(who int) string {
	mk := w.nobleItem(who)
	with := ""
	assert(w.validBox(who), "olympia/display.c", 179, "valid_box(who)")
	for e := range w.loopInv(who) {
		if int(e.item) == mk {
			continue
		}
		if p := w.rpItem(int(e.item)); p == nil || p.prominent == 0 {
			continue
		}
		if with != "" {
			with += ", "
		} else {
			with = ", with "
		}
		with += w.justNameQty(int(e.item), int(e.qty))
	}
	return sout("%s", with)
}

// displayWith is display.c:197–205: the raw here list, not all_char_here.
func (w *world) displayWith(who int) string {
	if w.rpLocInfo(who).hereList.Len() > 0 {
		return ", accompanied~by:"
	}
	return ""
}

func (w *world) mageS(n int) string {
	p := w.rpMagic(n)
	if p == nil || p.magician == 0 || p.hideMage != 0 {
		return ""
	}
	a := w.maxEffAura(n)
	for i, ceiling := range [...]int32{5, 10, 15, 20, 30, 40, 50, 60, 70, 80} {
		if a <= ceiling {
			return [...]string{"", ", conjurer", ", mage", ", wizard", ", sorcerer", ", 6th black circle", ", 5th black circle", ", 4th black circle", ", 3rd black circle", ", 2nd black circle"}[i]
		}
	}
	return ", master of the black arts"
}

func (w *world) priestS(n int) string {
	if w.hasSkill(n, skReligion) == 0 {
		return ""
	}
	return ", priest"
}

// emperor is display.c:372–399. The throne's unique owner, not inventory,
// determines the claimant; only a regular player in Olympus's castle qualifies.
func (w *world) emperor() int {
	who := 0
	if w.validBox(401) { // oly.h RELIC_THRONE
		who = w.itemUnique(401)
	}
	if who == 0 || w.player(who) == 0 || w.subkind(w.player(who)) != subPlRegular {
		return 0
	}
	where := w.subloc(who)
	if w.subkind(where) != subCastle || w.province(where) != int(w.mountOlympus) {
		return 0
	}
	return who
}

// linerDescChar is the complete display.c:406–477 formatter, in C suffix order.
func (w *world) linerDescChar(n int) string {
	buf := w.boxName(n)
	sk := w.subkind(n)
	if sk == subNi {
		mk := w.nobleItem(n)
		num := int32(w.hasItem(n, mk)) + 1
		buf += sout(", %s", w.pluralItemName(mk, int(num)))
		if num != 1 {
			buf += sout(", number:~%s", commaNum(int(num)))
		}
	} else if sk != 0 {
		buf += sout(", %s", subkindS[sk])
	}
	buf += w.rankS(n) + w.mageS(n) + w.priestS(n)
	if w.emperor() == n {
		buf += ", Emperor"
	}
	buf += w.wieldS(n)
	if w.showCombatFlag {
		if p := w.rpChar(n); p != nil && p.behind != 0 {
			buf += sout(", behind~%d%s", p.behind, w.combatAlly)
		} else {
			buf += w.combatAlly
		}
	} else if p := w.rpChar(n); p != nil && p.guard != 0 && w.stackLeader(n) == n {
		buf += ", on guard"
	}
	if w.showDisplayString {
		if p := w.rpMisc(n); p != nil && p.display != nil && *p.display != "" {
			buf += sout(", \"{<i>}%s{</i>}\"", *p.display)
		}
	}
	buf += w.withInventoryString(n)
	if w.isPrisoner(n) {
		buf += ", prisoner"
	}
	return sout("%s", buf)
}

func (w *world) linerDesc(n int) string {
	switch w.kind(n) {
	case tChar:
		return w.linerDescChar(n)
	case tShip, tLoc, tRoad, tStorm:
		panic(notImplemented("liner_desc noncharacter"))
	default:
		assert(false, "olympia/display.c", 541, "FALSE")
	}
	return ""
}

// showCharsBelow is display.c:561–576. all_char_here takes a depth-first
// snapshot, prunes noncharacter branches, and deliberately does not hide units.
func (w *world) showCharsBelow(who, n int) {
	assert(w.validBox(who), "olympia/display.c", 566, "valid_box(who)")
	w.indent += 3
	for _, i := range w.allCharHere(n) {
		assert(w.validBox(who), "olympia/display.c", 571, "valid_box(who)")
		w.wiout(who, 3, "%s", w.linerDesc(i))
	}
	w.indent -= 3
}
