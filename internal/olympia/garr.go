// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"

	"github.com/mdhender/g3g0/internal/lists"
)

// Ported from olympia/garr.c: the garrison ownership/output closure.

// garrisonGold is the complete garr.c:446–549 month-end settlement. In
// particular, spent measures this unit's gold, not its stack's deductions;
// garrForward is written only by the two forwarding branches, not reset.
func (w *world) garrisonGold() {
	w.clearTemps(tLoc)
	for garr := range w.loopSubkind(subGarrison) {
		if w.defaultGarrison(garr) != 0 {
			continue
		}
		where := w.subloc(garr)
		assert(w.locDepth(where) == locProvince, "olympia/garr.c", 471, "loc_depth(where) == LOC_province")
		base := w.hasItem(where, itemTaxCookie)
		w.consumeItem(where, itemTaxCookie, int(base))
		w.genItem(garr, itemGold, int(base))
		w.pMisc(garr).garrTax = base
		hasBefore := w.hasItem(garr, itemGold)
		w.chargeMaintSup(garr)
		hasNow := w.hasItem(garr, itemGold)
		spent := hasBefore - hasNow
		if spent >= base {
			continue
		}
		remain := base - spent
		w.consumeItem(garr, itemGold, int(remain))
		if w.garrisonStrength(garr) < 10 {
			continue
		}
		if w.provinceSubloc(where, subCastle) != 0 {
			w.genItem(where, itemTaxCookie, int(remain))
			w.pMisc(garr).garrForward = -1
			continue
		}
		remain /= 2
		castle := 0
		if p := w.rpMisc(garr); p != nil {
			castle = int(p.garrCastle)
		}
		if castle != 0 && remain > 0 {
			w.bx[castle].temp += remain
			w.pMisc(garr).garrForward = remain
		}
	}
	for i := range w.loopKind(tLoc) {
		if w.subkind(i) != subCastle {
			continue
		}
		owner := w.buildingOwner(i)
		if owner == 0 || w.bx[i].temp == 0 {
			continue
		}
		w.commandOut(owner, "Collected %s from owned provinces.", goldS(int(w.bx[i].temp)))
		w.genItem(owner, itemGold, int(w.bx[i].temp))
		w.goldProvinces += w.bx[i].temp
	}
}

// garrisonStrength is garr.c:54–71. Inventory iteration excludes invalid
// boxes and nonpositive quantities before applying all four fighter tests.
func (w *world) garrisonStrength(garr int) int32 {
	var strength int32
	assert(w.validBox(garr), "olympia/garr.c", 60, "valid_box(garr)")
	for e := range w.loopInv(garr) {
		p := w.rpItem(int(e.item))
		if p == nil || p.isManItem == 0 || maintCost(int(e.item)) <= 0 {
			continue
		}
		soldierAttack, soldierDefense := int16(0), int16(0)
		if soldier := w.rpItem(itemSoldier); soldier != nil {
			soldierAttack, soldierDefense = soldier.attack, soldier.defense
		}
		if p.attack >= soldierAttack && p.defense >= soldierDefense {
			strength += e.qty
		}
	}
	return strength
}

// pingGarrisons is garr.c:1064–1093. Visit the subtype chain, including
// garrisons that are not the first character in their province.
func (w *world) pingGarrisons() {
	w.showToGarrison = true
	for garr := range w.loopSubkind(subGarrison) {
		where := w.subloc(garr)
		p := w.rpChar(garr)
		assert(p != nil, "olympia/garr.c", 1077, "rp_char(garr) != NULL")
		p.guard = 0
		if w.garrisonStrength(garr) < 10 {
			continue
		}
		w.commandOut(where, "%s guards %s.", w.linerDesc(garr), w.boxName(where))
		p.guard = 1
	}
	w.showToGarrison = false
}

// garrisonNotices is garr.c:764–780. Only the target itself is checked in
// the watch list here; spotting below separately scans its entire stack.
func (w *world) garrisonNotices(garr, target int) bool {
	if w.isNPC(target) || w.countStackUnits(target) >= 5 || w.countStackFigures(target) >= 20 {
		return true
	}
	p := w.rpMisc(garr)
	return p != nil && p.garrWatch.Lookup(target) >= 0
}

// garrisonSpotCheck is garr.c:783–810, including its output before the
// caller prints the spotted stack and the ordinary arrival message.
func (w *world) garrisonSpotCheck(garr, target int) bool {
	assert(w.validBox(garr), "olympia/garr.c", 790, "valid_box(garr)")
	p := w.rpMisc(garr)
	if p == nil {
		return false
	}
	for i := range w.loopStack(target) {
		if p.garrWatch.Lookup(i) >= 0 {
			w.commandOut(garr, "Spotted in %s:", w.boxName(w.province(garr)))
			return true
		}
	}
	return false
}

// rankS is garr.c:603–624 rank_s; ranks are the garr.h constants.
func (w *world) rankS(who int) string {
	var rank int8
	if p := w.rpChar(who); p != nil {
		rank = p.rank
	}
	switch rank {
	case 0:
		return ""
	case 10:
		return ", lord"
	case 20:
		return ", knight"
	case 30:
		return ", baron"
	case 40:
		return ", count"
	case 50:
		return ", earl"
	case 60:
		return ", marquess"
	case 70:
		return ", duke"
	case 80:
		return ", king"
	default:
		assert(false, "olympia/garr.c", 621, "FALSE")
	}
	return ""
}

// garrisonHere is garrison_here(): the first unit in where's here list if
// it is a garrison, or 0.
func (w *world) garrisonHere(where int) int {
	n := w.firstCharacter(where)

	if n != 0 && w.subkind(n) == subGarrison {
		return n
	}

	return 0
}

func (w *world) provinceAdmin(n int) int {
	var garr int
	if w.kind(n) == tLoc {
		assert(w.locDepth(n) == locProvince, "olympia/garr.c", 81, "loc_depth(n) == LOC_province")
		garr = w.garrisonHere(n)
		if garr == 0 {
			return 0
		}
	} else {
		assert(w.subkind(n) == subGarrison, "olympia/garr.c", 88, "subkind(n) == sub_garrison")
		garr = n
	}
	castle := 0
	if p := w.rpMisc(garr); p != nil {
		castle = int(p.garrCastle)
	}
	if !w.validBox(castle) {
		return 0
	}
	return w.buildingOwner(castle)
}

func (w *world) playersWhoRuleHere(where int) lists.IList {
	var l lists.IList
	loopCheck := 5000
	for i := w.provinceAdmin(where); i > 0; {
		// loop_loc_owner reads the next pledge before entering its body.
		next := 0
		if p := w.rpMagic(i); p != nil {
			next = int(p.pledge)
		}
		pl := w.player(i)
		if pl != 0 && l.Lookup(pl) < 0 {
			l.Append(pl)
		}
		if loopCheck <= 0 {
			for j, pl := range l {
				fmt.Fprintf(w.stderr, "l[%d] = %d\n", j, pl)
			}
			fmt.Fprintf(w.stderr, "where = %d, i = %d\n", where, i)
		}
		assert(loopCheck > 0, "olympia/garr.c", 195, "loop_check-- > 0")
		loopCheck--
		i = next
	}
	return l
}
