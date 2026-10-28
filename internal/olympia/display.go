// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "strings"

// Location reports from display.c.

// showLoc is complete display.c:1056–1101. The display flag is reset to
// false, not restored to its incoming value, just as in C.
func (w *world) showLoc(who, where int) {
	assert(w.validBox(where), "olympia/display.c", 1061, "valid_box(where)")
	w.showDisplayString = true
	if !w.showLocNoHeader {
		w.commandOut(who, "%s", w.showLocHeader(where))
		w.commandOut(who, "")
	}
	if w.showLocBarrier(who, where) {
		w.commandOut(who, "")
	}
	if p := w.rpSubloc(where); p != nil && p.loot != 0 {
		recent := ""
		if p.recentLoot != 0 {
			recent = " recent"
		}
		w.commandOut(who, "Recovery from%s pillaging will take %s month%s.", recent, niceNum(int(p.loot)), addS(int(p.loot)))
		w.commandOut(who, "")
	}
	w.showLocRuler(who, where)
	w.listSailableRoutes(who, where)
	w.listExits(who, where)
	w.showLocStats(who, where)
	w.showNearbyCities(who, where)
	w.showLocSkills(who, where)
	if w.subkind(where) == subCity {
		w.marketReport(who, where)
	}
	w.showSublocsHere(who, where)
	w.showLocPosts(who, where, false)
	w.showWeather(who, where)
	w.showCharsHere(who, where)
	w.showShipsHere(who, where)
	w.showDisplayString = false
}

// showLocStats is display.c:92–122. Keep signed C-int multiplication before
// division and use the existing shipWeight traversal without allocating records.
func (w *world) showLocStats(who, where int) {
	first := true
	if p := w.rpSubloc(where); p != nil && p.damage != 0 {
		w.commandOut(who, "")
		first = false
		w.commandOut(who, "Damage: %d%%", p.damage)
	}
	if sc := w.shipCap(where); sc != 0 {
		sw := w.shipWeight(where)
		if first {
			w.commandOut(who, "")
		}
		w.commandOut(who, "Ship capacity: %s/%s (%d%%)", commaNum(int(sw)), commaNum(int(sc)), sw*100/sc)
	}
}

// showNearbyCities is display.c:842–869; duplicates and stored order survive.
func (w *world) showNearbyCities(who, where int) {
	p := w.rpSubloc(where)
	if p == nil || len(p.nearCities) < 1 {
		return
	}
	w.commandOut(who, "")
	w.commandOut(who, "Cities rumored to be nearby:")
	w.indent += 3
	for _, city := range p.nearCities {
		suffix := w.safeHavenS(city)
		w.commandOut(who, "%s, in %s%s", w.boxName(city), w.boxName(w.province(city)), suffix)
	}
	w.indent -= 3
}

// showLocSkills is display.c:872–892, including loop_loc_teach's assertion
// and copied iteration rather than sorting the saved teaches list.
func (w *world) showLocSkills(who, where int) {
	assert(w.validBox(where), "olympia/display.c", 878, "valid_box(where)")
	s := ""
	if p := w.rpSubloc(where); p != nil {
		for _, i := range append([]int(nil), p.teaches...) {
			if s != "" {
				s += ", "
			}
			s += w.boxName(i)
		}
	}
	if s != "" {
		w.commandOut(who, "")
		w.commandOut(who, "Skills taught here:")
		w.indent += 3
		w.commandOut(who, "%s", s)
		w.indent -= 3
	}
}

// showLocPosts is complete display.c:895–963, not an allocation projection.
// Validate each post before printing its separator; text is transient, not saved.
func (w *world) showLocPosts(who, where int, showFullLoc bool) {
	flag := true
	assert(w.validBox(where), "olympia/display.c", 904, "valid_box(where)")
	for post := range w.loopHere(where) {
		if w.kind(post) != tPost {
			continue
		}
		p := w.rpMisc(post)
		assert(p != nil && len(p.postTxt) >= 1, "olympia/display.c", 912, "FALSE")
		l := p.postTxt
		if flag {
			w.commandOut(who, "")
			var name string
			if showFullLoc {
				name = w.boxName(where)
			} else {
				name = w.justName(where)
			}
			w.commandOut(who, "Posted in %s:", name)
			flag = false
			w.indent += 3
		} else {
			w.commandOut(who, "")
		}
		creator := 0
		if m := w.rpItemMagic(post); m != nil {
			creator = int(m.creator)
		}
		if creator != 0 {
			w.commandOut(who, "Posted by %s:", w.boxName(creator))
		} else {
			w.commandOut(who, "Posted:")
		}
		w.indent += 3
		first := true
		for i, line := range l {
			before, after := "", ""
			if first {
				before = "\""
			}
			if i+1 == len(l) {
				after = "\""
			}
			w.commandOut(who, "%s%s%s", before, line, after)
			if first {
				first = false
				w.indent++
			}
		}
		if !first {
			w.indent--
		}
		w.indent -= 3
	}
	if !flag {
		w.indent -= 3
	}
}

// showWeather is display.c:966–1012. Weather aggregates at the province,
// but privileged storm descriptions traverse the original direct location.
func (w *world) showWeather(who, where int) {
	rain := w.weatherHere(where, subRain)
	wind := w.weatherHere(where, subWind)
	fog := w.weatherHere(where, subFog)
	if rain == 0 && wind == 0 && fog == 0 {
		return
	}
	w.commandOut(who, "")
	if rain != 0 {
		w.commandOut(who, "It is raining.")
	}
	if wind != 0 {
		w.commandOut(who, "It is windy.")
	}
	if fog != 0 {
		w.commandOut(who, "The province is blanketed in fog.")
	}
	if w.canSeeWeatherHere(who, where) {
		first := true
		assert(w.validBox(where), "olympia/display.c", 994, "valid_box(where)")
		for i := range w.loopHere(where) {
			if w.kind(i) != tStorm {
				continue
			}
			if first {
				w.indent += 3
				first = false
			}
			w.commandOut(who, "%s", w.linerDesc(i))
		}
		if !first {
			w.indent -= 3
		}
	}
}

// showLocRuler is display.c:1015–1046; topRuler runs even without a castle.
func (w *world) showLocRuler(who, where int) {
	if w.locDepth(where) != locProvince || w.subkind(where) == subOcean {
		return
	}
	garr := w.garrisonHere(where)
	if garr == 0 {
		return
	}
	castle := 0
	if p := w.rpMisc(garr); p != nil {
		castle = int(p.garrCastle)
	}
	ruler := w.topRuler(garr)
	if castle != 0 {
		prov := w.province(castle)
		w.commandOut(who, "Province controlled by %s, in %s", w.boxNameKind(castle), w.boxName(prov))
		if ruler != 0 {
			w.commandOut(who, "Ruled by %s%s", w.boxName(ruler), w.rankS(ruler))
		}
		w.commandOut(who, "")
	}
}

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

// displayOwner is display.c:208–217. first_char_here reads the direct here
// list, not the recursive character snapshot used for ship weights.
func (w *world) displayOwner(who int) string {
	assert(w.validBox(who), "olympia/u.c", 2395, "valid_box(where)")
	for i := range w.loopHere(who) {
		if w.kind(i) == tChar {
			return ", owner:"
		}
	}
	return ""
}

func (w *world) incompleteString(n int) string {
	p := w.rpSubloc(n)
	if p == nil || p.effortRequired == 0 {
		return ""
	}
	return sout(", %d%% completed", p.effortGiven*100/p.effortRequired)
}

func (w *world) shipCapS(n int) string {
	if sc := w.shipCap(n); sc != 0 {
		sw := w.shipWeight(n)
		return sout(", %d%% loaded", sw*100/sc)
	}
	return ""
}

// linerDescShip is display.c:237–272, in the original suffix/read order.
func (w *world) linerDescShip(n int) string {
	buf := w.boxNameKind(n) + w.incompleteString(n)
	buf += w.shipCapS(n)
	if p := w.rpSubloc(n); p != nil && p.defense != 0 {
		buf += sout(", defense~%d", p.defense)
	}
	if p := w.rpSubloc(n); p != nil && p.damage != 0 {
		buf += sout(", %d%%~damaged", p.damage)
	}
	if w.showDisplayString {
		if p := w.rpMisc(n); p != nil && p.display != nil && *p.display != "" {
			buf += sout(", \"{<i>}%s{</i>}\"", *p.display)
		}
	}
	if w.locHidden(n) != 0 {
		buf += ", hidden"
	}
	if p := w.rpSubloc(n); p != nil && p.galleyRam != 0 {
		buf += ", with ram"
	}
	if p := w.rpMagic(n); p != nil && p.fee != 0 {
		buf += sout(", %s per 100 wt. to board", goldS(int(p.fee)))
	}
	return sout("%s", buf)
}

// linerDescLoc is display.c:288–333. Location banners deliberately have no
// italic markup, unlike ship and character banners.
func (w *world) linerDescLoc(n int) string {
	buf := w.boxNameKind(n) + w.safeHavenS(n) + w.incompleteString(n)
	if w.locDepth(n) == locBuild {
		if p := w.rpSubloc(n); p != nil && p.defense != 0 {
			buf += sout(", defense~%d", p.defense)
		}
		if d := w.mineDepth(n); d != 0 {
			buf += sout(", depth~%d", d)
		}
		if p := w.rpSubloc(n); p != nil && p.damage != 0 {
			buf += sout(", %d%%~damaged", p.damage)
		}
		if p := w.rpSubloc(n); p != nil && p.castleLev != 0 {
			buf += sout(", level~%d", p.castleLev)
		}
	}
	if w.showDisplayString {
		if p := w.rpMisc(n); p != nil && p.display != nil && *p.display != "" {
			buf += sout(", \"%s\"", *p.display)
		}
	}
	if w.locHidden(n) != 0 {
		buf += ", hidden"
	}
	if w.locDepth(n) == locSubloc {
		if w.subkind(n) == subHadesPit {
			buf += ", 28 days"
		} else {
			buf += ", 1 day"
		}
	}
	return sout("%s", buf)
}

func (w *world) linerDescRoad(n int) string {
	dest := w.roadDest(n)
	hid := ""
	if w.roadHidden(n) != 0 {
		hid = ", hidden"
	}
	dist := w.exitDistance(w.loc(n), dest)
	return sout("%s, to %s%s, %d~day%s", w.boxName(n), w.boxName(dest), hid, dist, addS(int(dist)))
}

func (w *world) linerDescStorm(n int) string {
	buf := w.boxNameKind(n)
	p := w.rpMisc(n)
	if owner := w.npcSummoner(n); owner != 0 {
		buf += sout(", owner~%s", boxCodeLess(owner))
	}
	buf += sout(", strength~%s", commaNum(int(w.stormStrength(n))))
	if p != nil && p.npcDir != 0 {
		buf += sout(", heading %s", fullDirS[p.npcDir])
	}
	return sout("%s", buf)
}

func (w *world) linerDesc(n int) string {
	switch w.kind(n) {
	case tChar:
		return w.linerDescChar(n)
	case tShip:
		return w.linerDescShip(n)
	case tLoc:
		return w.linerDescLoc(n)
	case tRoad:
		return w.linerDescRoad(n)
	case tStorm:
		return w.linerDescStorm(n)
	default:
		assert(false, "olympia/display.c", 541, "FALSE")
	}
	return ""
}

// highlightUnits is display.c:547–558. The literal spaces are part of the
// line, independently of wiout's secondary indentation metadata.
func (w *world) highlightUnits(who, n int, depth int32) string {
	assert(depth >= 3, "olympia/display.c", 551, "depth >= 3")
	assert(w.indent == 0, "olympia/display.c", 552, "indent == 0")
	if w.kind(who) == tPlayer && w.player(n) == who {
		return sout(" *%s", strings.Repeat(" ", int(depth-2)))
	}
	return strings.Repeat(" ", int(depth))
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

// showCharsBelowHighlight is display.c:579–593: all descendants have the
// same depth, not one extra indentation level per nested character.
func (w *world) showCharsBelowHighlight(who, n int, depth int32) {
	depth += 3
	for _, i := range w.allCharHere(n) {
		w.wiout(who, depth, "%s%s", w.highlightUnits(who, i, depth), w.linerDesc(i))
	}
}

func (w *world) showOwnerStack(who, n int) {
	depth := w.indent + 3
	w.indent = 0
	first := true
	assert(w.validBox(n), "olympia/display.c", 606, "valid_box(n)")
	for i := range w.loopHere(n) {
		if w.kind(i) == tChar {
			if !first && w.charReallyHidden(i) {
				continue
			}
			w.wiout(who, depth, "%s%s%s", w.highlightUnits(who, i, depth), w.linerDesc(i), w.displayWith(i))
			w.showCharsBelowHighlight(who, i, depth)
			first = false
		}
	}
	w.indent = depth - 3
}

func (w *world) showCharsHere(who, where int) {
	if w.locDepth(where) == locProvince && w.weatherHere(where, subFog) != 0 {
		w.commandOut(who, "")
		w.commandOut(who, "No one can be seen through the fog.")
		return
	}
	depth := w.indent
	w.indent = 0
	flying := ""
	if w.subkind(where) == subOcean {
		flying = ", flying"
	}
	first := true
	assert(w.validBox(where), "olympia/display.c", 650, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.kind(i) == tChar {
			if w.charReallyHidden(i) {
				continue
			}
			if first {
				first = false
				w.commandOut(who, "")
				w.commandOut(who, "Seen here:")
				depth += 3
			}
			w.wiout(who, depth, "%s%s%s%s", w.highlightUnits(who, i, depth), w.linerDesc(i), flying, w.displayWith(i))
			w.showCharsBelowHighlight(who, i, depth)
		}
	}
	if !first {
		depth -= 3
	}
	w.indent = depth
}

// showInnerLocs is display.c:685–717. Each visible container is rendered
// before its owners, barrier and recursively nested locations and ships.
func (w *world) showInnerLocs(who, where int) {
	first := true
	assert(w.validBox(where), "olympia/display.c", 691, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.kind(i) == tLoc || w.kind(i) == tShip {
			if w.locHidden(i) != 0 && !w.testKnown(who, i) && w.immedSeeAll == 0 {
				continue
			}
			if first {
				first = false
				w.indent += 3
			}
			w.commandOut(who, "%s%s", w.linerDesc(i), w.displayOwner(i))
			w.showOwnerStack(who, i)
			w.showLocBarrier(who, i)
			w.showInnerLocs(who, i)
		}
	}
	if !first {
		w.indent -= 3
	}
}

// showSublocsHere is display.c:720–797. Direct cities omit owners and
// recursion; open link_from rows only describe the target, retaining duplicates.
func (w *world) showSublocsHere(who, where int) {
	first := true
	assert(w.validBox(where), "olympia/display.c", 727, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.kind(i) == tLoc {
			if w.locHidden(i) != 0 && !w.testKnown(who, i) && w.immedSeeAll == 0 {
				continue
			}
			if first {
				first = false
				w.commandOut(who, "")
				w.commandOut(who, "Inner locations:")
				w.indent += 3
			}
			if w.subkind(i) == subCity {
				w.commandOut(who, "%s", w.linerDesc(i))
				w.showLocBarrier(who, i)
			} else {
				w.commandOut(who, "%s%s", w.linerDesc(i), w.displayOwner(i))
				w.showOwnerStack(who, i)
				w.showLocBarrier(who, i)
				w.showInnerLocs(who, i)
			}
		}
	}
	if p := w.rpSubloc(where); p != nil {
		for _, i := range p.linkFrom {
			if w.locHidden(i) != 0 && !w.testKnown(who, i) && w.immedSeeAll == 0 {
				continue
			}
			if w.locLinkOpen(i) != 0 {
				if first {
					first = false
					w.commandOut(who, "")
					w.commandOut(who, "Inner locations:")
					w.indent += 3
				}
				w.commandOut(who, "%s", w.linerDesc(i))
			}
		}
	}
	if !first {
		w.indent -= 3
	}
}

func (w *world) showShipsHere(who, where int) {
	first := true
	assert(w.validBox(where), "olympia/display.c", 806, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.kind(i) == tShip {
			if w.locHidden(i) != 0 && !w.testKnown(who, i) && w.immedSeeAll == 0 {
				continue
			}
			if first {
				first = false
				w.commandOut(who, "")
				if w.subkind(where) == subOcean {
					w.commandOut(who, "Ships sighted:")
				} else {
					w.commandOut(who, "Ships docked at port:")
				}
				w.indent += 3
			}
			w.wiout(who, 3, "%s%s", w.linerDesc(i), w.displayOwner(i))
			w.showOwnerStack(who, i)
			w.showLocBarrier(who, i)
		}
	}
	if !first {
		w.indent -= 3
	}
}
