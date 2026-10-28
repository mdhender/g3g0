// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

// wealthListComp and noblesListComp are gm.c:739–748 and 794–803.
func (w *world) wealthListComp(a, b int) int { return w.skillsKnownComp(a, b) }
func (w *world) noblesListComp(a, b int) int { return w.skillsKnownComp(a, b) }

// gmFactionWealth is complete gm.c:751–791. Only regular factions rank,
// but every character contributes to its owner's scratch total.
func (w *world) gmFactionWealth(pl int) {
	w.gmWealthList.Clear()
	w.clearTemps(tPlayer)
	for i := range w.loopKind(tChar) {
		w.bx[w.player(i)].temp += int32(w.hasItem(i, itemGold))
	}
	for i := range w.loopKind(tPlayer) {
		if w.subkind(i) == subPlRegular {
			w.gmWealthList.Append(i)
		}
	}
	safeQsort(w.gmWealthList, w.wealthListComp)
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "%4s %11s  %s", "rank", "gold", "faction")
	w.commandOut(pl, "%4s %11s  %s", "----", "----", "-------")
	for i, n := range w.gmWealthList {
		w.commandOut(pl, "%4d %11s  %s", i+1, commaNum(int(w.bx[n].temp)), w.boxName(n))
	}
	w.commandOut(pl, "")
	w.outPath, w.outAltWho = 0, 0
}

// gmNoblesList is complete gm.c:806–846, including all character subtypes.
func (w *world) gmNoblesList(pl int) {
	w.gmNoblesListScratch.Clear()
	w.clearTemps(tPlayer)
	for i := range w.loopKind(tChar) {
		w.bx[w.player(i)].temp++
	}
	for i := range w.loopKind(tPlayer) {
		if w.subkind(i) == subPlRegular {
			w.gmNoblesListScratch.Append(i)
		}
	}
	safeQsort(w.gmNoblesListScratch, w.noblesListComp)
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "%4s %11s  %s", "rank", "nobles", "faction")
	w.commandOut(pl, "%4s %11s  %s", "----", "------", "-------")
	for i, n := range w.gmNoblesListScratch {
		w.commandOut(pl, "%4d %11s  %s", i+1, commaNum(int(w.bx[n].temp)), w.boxName(n))
	}
	w.commandOut(pl, "")
	w.outPath, w.outAltWho = 0, 0
}

// gmPlayerDetails is complete gm.c:849–961. Separate unit traversals retain
// the source's ownership/read order; no player record is allocated for age.
func (w *world) gmPlayerDetails(pl int) {
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "%4s %3s %5s %5s %4s %6s %5s %6s", "who", "age", "gold", "units", "bld", "subloc", "ships", "skills")
	w.commandOut(pl, "%4s %3s %5s %5s %4s %6s %5s %6s", "---", "---", "----", "-----", "---", "------", "-----", "------")
	for i := range w.loopKind(tPlayer) {
		if w.subkind(i) == subPlSystem {
			continue
		}
		var gold, units, buildings, sublocs, ships, skills int32
		for j := range w.loopUnits(i) {
			gold += int32(w.hasItem(j, itemGold))
		}
		for range w.loopUnits(i) {
			units++
		}
		for j := range w.loopUnits(i) {
			where := w.subloc(j)
			if w.locDepth(where) == locBuild && w.buildingOwner(where) == j && !w.isShip(where) && !w.isShipNotdone(where) {
				buildings++
			}
		}
		for j := range w.loopUnits(i) {
			where := w.subloc(j)
			if w.locDepth(where) == locSubloc && w.firstCharacter(where) == j {
				sublocs++
			}
		}
		for j := range w.loopUnits(i) {
			where := w.subloc(j)
			if (w.isShip(where) || w.isShipNotdone(where)) && w.buildingOwner(where) == j {
				ships++
			}
		}
		w.clearTemps(tSkill)
		for j := range w.loopUnits(i) {
			if p := w.rpChar(j); p != nil {
				for _, e := range p.skills.Copy() {
					if e.know == skillKnow {
						w.bx[e.skill].temp++
					}
				}
			}
		}
		for j := range w.loopKind(tSkill) {
			if w.bx[j].temp != 0 {
				skills++
			}
		}
		age := "???"
		if p := w.rpPlayer(i); p != nil {
			age = sout("%d", int32(w.sysclock.turn)-p.firstTurn)
		}
		w.commandOut(pl, "%4s %3s %5s %5s %4s %6s %5s %6s  %s", boxCodeLess(i), age, knum(gold, false), knum(units, false), knum(buildings, true), knum(sublocs, true), knum(ships, true), knum(skills, true), w.justName(i))
	}
	w.commandOut(pl, "")
	w.outPath, w.outAltWho = 0, 0
}

// listAllNotices is complete gm.c:964–992; real posts validate before text.
func (w *world) listAllNotices(pl int) {
	w.clearTemps(tLoc)
	for i := range w.loopKind(tPost) {
		w.bx[w.subloc(i)].temp++
	}
	w.outPath, w.outAltWho = master, 2
	for i := range w.loopKind(tLoc) {
		if w.bx[i].temp != 0 {
			w.showLocPosts(pl, i, true)
		}
	}
	w.outPath, w.outAltWho = 0, 0
}

// gmReport is complete gm.c:1025–1047. Item output has its own recipient.
func (w *world) gmReport(pl int) {
	w.stage("gm_report()")
	w.gmShowInterestingAttributes(pl)
	w.gmShowGold(pl)
	w.gmShowControlArts(pl)
	w.gmListAnimateItems(pl)
	w.gmShowSkillUseCounts(pl)
	w.gmCountStuff(pl)
	w.gmLandStats(pl)
	w.gmShowGateStats(pl)
	w.gmShowLocsVisited(pl)
	w.gmLoyaltyStats(pl)
	w.gmShowSkillsKnown(pl)
	w.gmFactionWealth(pl)
	w.gmNoblesList(pl)
	w.gmPlayerDetails(pl)
	w.listAllNotices(pl)
	w.listAllItems(int(w.skillPlayer))
}

// skillUseComp is gm.c:40–54, including the pinned ID tie-break.
func (w *world) skillUseComp(a, b int) int {
	pa, pb := w.rpSkill(a), w.rpSkill(b)
	if pa.useCount != pb.useCount {
		return int(pb.useCount - pa.useCount)
	}
	return a - b
}

// gmShowSkillUseCounts is complete gm.c:57–104.
func (w *world) gmShowSkillUseCounts(pl int) {
	w.gmSkillUseList.Clear()
	for sk := range w.loopKind(tSkill) {
		if w.skillSchool(sk) == sk {
			continue
		}
		if p := w.rpSkill(sk); p != nil && p.useCount != 0 {
			w.gmSkillUseList.Append(sk)
		}
	}
	safeQsort(w.gmSkillUseList, w.skillUseComp)
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "")
	w.commandOut(pl, "Skill use counts:")
	w.commandOut(pl, "-----------------")
	w.commandOut(pl, "")
	w.commandOut(pl, "%5s  %4s  %5s  %s", "count", "who", "skill", "name")
	w.commandOut(pl, "%5s  %4s  %5s  %s", "-----", "---", "-----", "----")
	for _, sk := range w.gmSkillUseList {
		p := w.rpSkill(sk)
		w.commandOut(pl, "%4d   %4s  %5s  %s", p.useCount, boxCodeLess(w.player(int(p.lastUseWho))), boxCodeLess(sk), w.justName(sk))
	}
	w.commandOut(pl, "")
	w.outPath, w.outAltWho = 0, 0
}

// skillsKnownComp is gm.c:107–116. Region occupancy uses the same total order.
func (w *world) skillsKnownComp(a, b int) int {
	if w.bx[a].temp != w.bx[b].temp {
		return int(w.bx[b].temp - w.bx[a].temp)
	}
	return a - b
}

// gmShowSkillsKnown is complete gm.c:118–179, including zero-count skills.
func (w *world) gmShowSkillsKnown(pl int) {
	w.gmSkillsKnownList.Clear()
	w.clearTemps(tSkill)
	for i := range w.loopKind(tChar) {
		if p := w.rpChar(i); p != nil {
			for _, e := range p.skills.Copy() {
				if e.know == skillKnow {
					w.bx[e.skill].temp++
				}
			}
		}
	}
	for sk := range w.loopKind(tSkill) {
		if w.skillSchool(sk) != sk {
			w.gmSkillsKnownList.Append(sk)
		}
	}
	safeQsort(w.gmSkillsKnownList, w.skillsKnownComp)
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "")
	w.commandOut(pl, "Skills known by players:")
	w.commandOut(pl, "------------------------")
	w.commandOut(pl, "")
	w.commandOut(pl, "%5s  %5s  %4s  %s", "count", "skill", "use", "name")
	w.commandOut(pl, "%5s  %5s  %4s  %s", "-----", "-----", "---", "----")
	for _, sk := range w.gmSkillsKnownList {
		p := w.rpSkill(sk)
		use := ""
		if p.useCount != 0 {
			use = commaNum(int(p.useCount))
		}
		w.commandOut(pl, "%4d   %5s  %4s  %s", w.bx[sk].temp, boxCodeLess(sk), use, w.justName(sk))
	}
	w.commandOut(pl, "")
	w.outPath, w.outAltWho = 0, 0
}

// gmShowGateStats is complete gm.c:476–536.
func (w *world) gmShowGateStats(pl int) {
	w.clearTemps(tGate)
	for i := range w.loopKind(tPlayer) {
		known := w.pPlayer(i).known
		sortKnown(known)
		for _, j := range known {
			if w.kind(j) == tGate {
				w.bx[j].temp++
			}
		}
	}
	var gates, found, sealed, jump, unseal int32
	for i := range w.loopKind(tGate) {
		gates++
		if w.bx[i].temp != 0 {
			found++
		}
		if p := w.rpGate(i); p != nil {
			if p.sealKey != 0 {
				sealed++
			}
			if p.notifyJumps != 0 {
				jump++
			}
			if p.notifyUnseal != 0 {
				unseal++
			}
		}
	}
	w.outPath, w.outAltWho = master, 2
	// The pinned zero-gate return deliberately leaves both routing globals set.
	if gates == 0 {
		return
	}
	w.commandOut(pl, "%d/%d gates found (%d%%)", found, gates, found*100/gates)
	w.commandOut(pl, "    %d sealed (%d%%), %d notify jump (%d%%), %d notify unseal (%d%%)", sealed, sealed*100/gates, jump, jump*100/gates, unseal, unseal*100/gates)
	w.commandOut(pl, "")
	w.outPath, w.outAltWho = 0, 0
}

// gmShowLocsVisited is complete gm.c:539–634.
func (w *world) gmShowLocsVisited(pl int) {
	w.clearTemps(tLoc)
	for i := range w.loopKind(tPlayer) {
		known := w.pPlayer(i).known
		sortKnown(known)
		for _, j := range known {
			if w.kind(j) == tLoc {
				w.bx[j].temp++
			}
		}
	}
	var prov, pv, sub, sv, hid, vis int32
	for i := range w.loopKind(tLoc) {
		switch w.locDepth(i) {
		case locProvince:
			prov++
			if w.bx[i].temp != 0 {
				pv++
			}
		case locSubloc:
			sub++
			if w.bx[i].temp != 0 {
				sv++
				if w.locHidden(i) != 0 {
					hid++
				} else {
					vis++
				}
			}
		}
	}
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "%d/%d provinces visited (%d%%)", pv, prov, pct(pv, prov))
	w.commandOut(pl, "%d/%d sublocs visited (%d%%)", sv, sub, pct(sv, sub))
	w.commandOut(pl, "    %d%% visible, %d%% hidden", pct(vis, sv), pct(hid, sv))
	var nf, nt, fv, fh int32
	for i := range w.loopKind(tLoc) {
		if w.locDepth(i) != locProvince || w.bx[i].temp == 0 {
			continue
		}
		for j := range w.loopHere(i) {
			if w.kind(j) != tLoc || w.locDepth(j) != locSubloc {
				continue
			}
			nt++
			if w.bx[j].temp != 0 {
				nf++
				if w.locHidden(j) != 0 {
					fh++
				} else {
					fv++
				}
			}
		}
	}
	w.commandOut(pl, "    %d%% of visisted province's sublocs found", pct(nf, nt))
	w.commandOut(pl, "    %d%% visible, %d%% hidden", pct(fv, nf), pct(fh, nf))
	w.commandOut(pl, "")
	w.outPath, w.outAltWho = 0, 0
}

// gmLoyaltyStats is complete gm.c:637–678; NPC loyalty is in the denominator.
func (w *world) gmLoyaltyStats(pl int) {
	var tot, oath, fear, cont, unsw int32
	for i := range w.loopKind(tChar) {
		if w.subkind(i) != 0 {
			continue
		}
		tot++
		var loy int8
		if p := w.rpChar(i); p != nil {
			loy = p.loyKind
		}
		switch loy {
		case loyOath:
			oath++
		case loyFear:
			fear++
		case loyContract:
			cont++
		case loyUnsworn:
			unsw++
		case loyNpc:
		default:
			fmt.Fprintf(w.stderr, "loy_kind = %d\n", loy)
			assert(false, "olympia/gm.c", 659, "FALSE")
		}
	}
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "%d chars: %d oath (%d%%), %d fear (%d%%), %d contract (%d%%), %d unsworn (%d%%)", tot, oath, pct(oath, tot), fear, pct(fear, tot), cont, pct(cont, tot), unsw, pct(unsw, tot))
	w.commandOut(pl, "")
	w.outPath, w.outAltWho = 0, 0
}

// regionOccupyComp and gmLandStats are complete gm.c:681–736.
func (w *world) regionOccupyComp(a, b int) int { return w.skillsKnownComp(a, b) }

func (w *world) gmLandStats(pl int) {
	w.gmRegionList.Clear()
	w.clearTemps(tLoc)
	var n int32
	for i := range w.loopKind(tChar) {
		w.bx[w.region(i)].temp++
		n++
	}
	for i := range w.loopKind(tLoc) {
		if w.bx[i].temp != 0 {
			w.gmRegionList.Append(i)
		}
	}
	safeQsort(w.gmRegionList, w.regionOccupyComp)
	w.outPath, w.outAltWho = master, 2
	w.commandOut(pl, "%10s  %s", "nobles", "region")
	w.commandOut(pl, "%10s  %s", "------", "------")
	for _, i := range w.gmRegionList {
		w.commandOut(pl, "%10s  %s", commaNum(int(w.bx[i].temp)), w.justName(i))
	}
	w.commandOut(pl, "%10s  %s", "======", "")
	w.commandOut(pl, "%10s  %s", commaNum(int(n)), "")
	w.commandOut(pl, "")
	w.outPath, w.outAltWho = 0, 0
}

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
