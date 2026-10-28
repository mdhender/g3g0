// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// pct is gm.c:31–37. Keep C's signed-int multiplication and truncation.
func pct(n, d int32) int32 {
	if d == 0 {
		return 0
	}
	return n * 100 / d
}

// gmShowInterestingAttributes is complete gm.c:183–280.
func (w *world) gmShowInterestingAttributes(pl int) {
	w.outPath, w.outAltWho = master, 2 // OUT_LORE
	ability, hinder, project, quick := 0, 0, 0, 0
	shroud, barrier, opium := 0, 0, 0
	format, players, rams, galleys := 0, 0, 0, 0
	for i := range w.loopKind(tChar) {
		if p := w.rpMagic(i); p != nil {
			if p.abilityShroud != 0 {
				ability++
			}
			if p.hinderMeditation != 0 {
				hinder++
			}
			if p.projectCast != 0 {
				project++
			}
			if p.quickCast != 0 {
				quick++
			}
		}
	}
	for i := range w.loopKind(tLoc) {
		if p := w.rpLoc(i); p != nil {
			if p.shroud != 0 {
				shroud++
			}
			if p.barrier != 0 {
				barrier++
			}
		}
		if p := w.rpSubloc(i); p != nil && p.opiumEcon != 0 {
			opium++
		}
	}
	for i := range w.loopKind(tShip) {
		if w.subkind(i) == subGalley {
			galleys++
			if p := w.rpSubloc(i); p != nil && p.galleyRam != 0 {
				rams++
			}
		}
	}
	for i := range w.loopKind(tPlayer) {
		if w.subkind(i) == subPlRegular {
			players++
			if p := w.rpPlayer(i); p != nil && p.format != 0 {
				format++
			}
		}
	}
	w.commandOut(pl, "Interesting attribute counts")
	w.commandOut(pl, "----------------------------")
	w.commandOut(pl, "")
	w.commandOut(pl, "char ability shroud:    %d", ability)
	w.commandOut(pl, "char hinder meditate:   %d", hinder)
	w.commandOut(pl, "char project cast:      %d", project)
	w.commandOut(pl, "char quick cast:        %d", quick)
	w.commandOut(pl, "loc shroud:             %d", shroud)
	w.commandOut(pl, "loc barrier:            %d", barrier)
	w.commandOut(pl, "loc opium:              %d", opium)
	w.commandOut(pl, "format one:             %d/%d", format, players)
	w.commandOut(pl, "")
	if galleys != 0 {
		w.commandOut(pl, "galleys with rams:      %d (%d%%)", rams, rams*100/galleys)
	}
	w.outPath, w.outAltWho = 0, 0
}

// gmListAnimateItems is complete gm.c:283–329. Capturability alone does
// not select an item: beast_capturable is the character predicate, as in C.
func (w *world) gmListAnimateItems(pl int) {
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "")
	w.commandOut(pl, "Animate items list")
	w.commandOut(pl, "------------------")
	w.commandOut(pl, "")
	w.commandOut(pl, "%25s %8s %8s %8s  %s", "name", "swamp", "man-like", "beast", "fighter")
	w.commandOut(pl, "%25s %8s %8s %8s  %s", "----", "-----", "--------", "-----", "-------")
	for i := range w.loopKind(tItem) {
		var attack, defense, missile int16
		var animal, man, capturable bool
		if p := w.rpItem(i); p != nil {
			attack, defense, missile = p.attack, p.defense, p.missile
			animal, man, capturable = p.animal != 0, p.isManItem != 0, p.capturable != 0
		}
		fighter := attack != 0 || defense != 0 || missile != 0 || i == itemGhostWarrior
		if !fighter && !man && !w.beastCapturable(i) && !animal {
			continue
		}
		buf := " -"
		if fighter {
			buf = sout("(%d,%d,%d)", attack, defense, missile)
		}
		flags := [3]string{"no  ", "no  ", "no  "}
		for j, yes := range []bool{animal, man, capturable} {
			if yes {
				flags[j] = "yes "
			}
		}
		w.commandOut(pl, "%25s %8s %8s %8s  %s", w.boxName(i), flags[0], flags[1], flags[2], buf)
	}
	w.outPath, w.outAltWho = 0, 0
}

// gmShowGold is complete gm.c:332–377. Player combat transfers are shown
// separately and are deliberately not part of generated gold.
func (w *world) gmShowGold(pl int) {
	sum := w.goldCommonMagic + w.goldLeadToGold + w.goldPotBasket + w.goldTrade + w.goldOpium + w.goldInn + w.goldTaxes + w.goldProvinces + w.goldGarrison + w.goldTimes + w.goldCombatIndep + w.goldPettyThief + w.goldTemple + w.goldPillage + w.goldFerry
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "")
	w.commandOut(pl, "Gold report")
	w.commandOut(pl, "-----------")
	w.commandOut(pl, "")
	if sum > 0 {
		for _, row := range []struct {
			label string
			n     int32
		}{
			{"Common magic:", w.goldCommonMagic}, {"Lead to gold:", w.goldLeadToGold},
			{"Pots and baskets:", w.goldPotBasket}, {"Opium:", w.goldOpium},
			{"Trade to cities:", w.goldTrade}, {"Inn income:", w.goldInn},
			{"Taxes:", w.goldTaxes}, {"Provinces:", w.goldProvinces},
			{"Times press:", w.goldTimes}, {"Combat with indeps:", w.goldCombatIndep},
			{"Petty thievery:", w.goldPettyThief}, {"Temple income:", w.goldTemple},
			{"Pillaging:", w.goldPillage}, {"Ferry boarding:", w.goldFerry}, {"Garrison Pay:", w.goldGarrison},
		} {
			w.commandOut(pl, "%-22s%10s %3d%%", row.label, commaNum(int(row.n)), pct(row.n, sum))
		}
		w.commandOut(pl, "                      %10s %4s", "", "----")
		w.commandOut(pl, "Total:                %10s", commaNum(int(sum)))
		w.commandOut(pl, "")
		w.commandOut(pl, "Player combat:        %10s", commaNum(int(w.goldCombat)))
		w.commandOut(pl, "")
	}
	w.outPath, w.outAltWho = 0, 0
}

// gmShowControlArts is complete gm.c:380–406.
func (w *world) gmShowControlArts(pl int) {
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "")
	w.commandOut(pl, "Captured control artifacts")
	w.commandOut(pl, "--------------------------")
	w.commandOut(pl, "")
	for item := range w.loopSubkind(subNpcToken) {
		n := w.tokenPlayer(w.itemUnique(item))
		if n != int(w.indepPlayer) {
			w.commandOut(pl, "%-33s  %s", w.boxName(item), w.boxName(n))
		}
	}
	w.commandOut(pl, "")
	w.outPath, w.outAltWho = 0, 0
}

// gmCountStuff is complete gm.c:409–473.
func (w *world) gmCountStuff(pl int) {
	counts := make(map[int]int)
	for i := range w.loopLocOrShip() {
		counts[w.subkind(i)]++
	}
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "%10s  %9s  %s", "", "finished", "unfinished")
	w.commandOut(pl, "%10s +----------------------", "")
	for _, row := range []struct {
		name             string
		done, unfinished int
	}{
		{"galley", subGalley, subGalleyNotdone}, {"roundship", subRoundship, subRoundshipNotdone},
		{"tower", subTower, subTowerNotdone}, {"castle", subCastle, subCastleNotdone},
		{"mine", subMine, subMineNotdone}, {"inn", subInn, subInnNotdone}, {"temple", subTemple, subTempleNotdone},
	} {
		w.commandOut(pl, "%10s |%9s  %8s", row.name, commaNum(counts[row.done]), commaNum(counts[row.unfinished]))
	}
	w.commandOut(pl, "")
	w.outPath, w.outAltWho = 0, 0
}

// listAllItems is complete gm.c:995–1022; missing item records read as zero.
func (w *world) listAllItems(pl int) {
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "%4s %-24s %6s %4s %4s %4s", "item", "name", "weight", "land", "ride", "fly")
	w.commandOut(pl, "%4s %-24s %6s %4s %4s %4s", "----", "----", "------", "----", "----", "---")
	for i := range w.loopKind(tItem) {
		var weight, land, ride, fly int16
		if p := w.rpItem(i); p != nil {
			weight, land, ride, fly = p.weight, p.landCap, p.rideCap, p.flyCap
		}
		w.commandOut(pl, "%4s %-24s %6s %4s %4s %4s", boxCodeLess(i), w.justName(i), commaNum(int(weight)), commaNum(int(land)), commaNum(int(ride)), commaNum(int(fly)))
	}
	w.commandOut(pl, "")
	w.outPath, w.outAltWho = 0, 0
}
