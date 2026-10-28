// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"strings"

	"github.com/mdhender/g3g0/internal/lists"
)

// Complete report.c:1062–1185. temp is the original reverse-built linked
// list, shared with the ancestor catch-up markers; it is not saved state.
func (w *world) locIndS(where int) string {
	ld := w.locDepth(where) - 1
	if ld <= 0 {
		return w.justName(where)
	}
	return strings.Repeat(" ", ld*2) + w.boxName(where)
}

func (w *world) locStackCatchup(pl, where int) {
	if where == 0 || w.bx[where].temp != 0 {
		return
	}
	w.locStackCatchup(pl, w.loc(where))
	w.commandOut(pl, "%s", w.locIndS(where))
	w.bx[where].temp = -1
}

func (w *world) locStackRepSup(pl, where, who int) {
	whereS, star, ind := "", "", ""
	if where != 0 {
		w.locStackCatchup(pl, w.loc(where))
		whereS = w.locIndS(where)
	}
	if w.player(who) != pl {
		star = " *"
		w.locStackExplain = true
	}
	if w.kind(w.loc(who)) == tChar {
		ind = "  "
	}
	w.commandOut(pl, "%-34s %s%s%s", whereS, ind, w.boxName(who), star)
}

func (w *world) locStackReport(pl int) {
	var l, locs lists.IList
	w.clearTemps(tLoc)
	w.clearTemps(tShip)
	w.clearTemps(tChar)
	w.locStackExplain = false
	for i := range w.loopUnits(pl) {
		if w.isPrisoner(i) {
			continue
		}
		l.Append(i)
	}
	if len(l) < 1 {
		return
	}
	w.commandOut(pl, "")
	w.commandOut(pl, "{<b>}%-34s %s{</b>}", "Location", "Stack")
	w.commandOut(pl, "%-34s %s", "--------", "-----")
	w.sortForOutput(l)
	for i := len(l) - 1; i >= 0; i-- {
		where := w.subloc(l[i])
		if w.bx[where].temp == 0 {
			locs.Append(where)
		}
		w.bx[l[i]].temp = w.bx[where].temp
		w.bx[where].temp = int32(l[i])
	}
	w.sortForOutput(locs)
	for _, where := range locs {
		j := int(w.bx[where].temp)
		assert(w.validBox(j), "olympia/report.c", 1168, "valid_box(j)")
		w.locStackRepSup(pl, where, j)
		for j = int(w.bx[j].temp); j != 0; j = int(w.bx[j].temp) {
			assert(w.validBox(j), "olympia/report.c", 1174, "valid_box(j)")
			w.locStackRepSup(pl, 0, j)
		}
	}
	if w.locStackExplain {
		w.commandOut(pl, "")
		w.commandOut(pl, "%-34s    %s", "", "* -- unit belongs to another faction")
	}
}

// Complete report.c:1215–1320, using the actual recursive weight owner.
// Ratios multiply as C int before division, including overflow behavior.
func (w *world) stackCapacityReport(pl int) {
	first := true
	for who := range w.loopUnits(pl) {
		if first {
			w.commandOut(pl, "")
			w.commandOut(pl, "{<b>}%*s  %10s %15s %15s %15s{</b>}", charField, "stack", "total wt", "   walk   ", "   ride   ", "   fly    ")
			w.commandOut(pl, "%*s  %10s %15s %15s %15s", charField, "-----", "--------", "-----------", "-----------", "-----------")
			first = false
		}
		v := w.determineStackWeights(who)
		walkS, rideS, flyS := "", "", ""
		if v.landCap > 0 {
			n := v.landWeight * 100 / v.landCap
			s := " -- "
			if n <= 999 {
				s = sout("%3d%%", n)
			}
			if v.landWeight > v.landCap {
				walkS = sout("(%s) %s", commaNum(int(v.landWeight-v.landCap)), s)
			} else {
				walkS = sout("%s %s", commaNum(int(v.landCap-v.landWeight)), s)
			}
		}
		if v.rideCap > 0 {
			n := v.rideWeight * 100 / v.rideCap
			s := " -- "
			if n <= 999 {
				s = sout("%3d%%", n)
			}
			if v.rideWeight > v.rideCap {
				rideS = sout("(%s) %s", commaNum(int(v.rideWeight-v.rideCap)), s)
			} else {
				rideS = sout("%s %s", commaNum(int(v.rideCap-v.rideWeight)), s)
			}
		}
		if v.flyCap > 0 {
			n := v.flyWeight * 100 / v.flyCap
			s := " -- "
			if n <= 999 {
				s = sout("%3d%%", n)
			}
			if v.flyWeight > v.flyCap {
				flyS = sout("(%s) %s", commaNum(int(v.flyWeight-v.flyCap)), s)
			} else {
				flyS = sout("%s %s", commaNum(int(v.flyCap-v.flyWeight)), s)
			}
		}
		w.commandOut(pl, "%*s  %10s %15s %15s %15s", charField, boxCodeLess(who), commaNum(int(v.totalWeight)), walkS, rideS, flyS)
	}
}

// Complete report.c:801–1059. These are the original byte-oriented name
// operations, not a word-splitting abbreviation algorithm.
const truncName = 15

var stupidWords = [...]string{"a", "the", "of", "de", "des", "la", "and", "du", "aux", "et", "ses", "avec", "un", "van", "von", "-", "--"}

func stripLeadingStupidWord(s string) string {
	for _, word := range stupidWords {
		if len(s) > len(word) && iStrcmp(s[:len(word)], word) == 0 && s[len(word)] == ' ' {
			t := len(word)
			for t < len(s) && s[t] == ' ' {
				t++
			}
			if t < len(s) {
				return s[t:]
			}
			return s
		}
	}
	return s
}

func stupidWord(s string) bool {
	for _, word := range stupidWords {
		if iStrcmp(s, word) == 0 {
			return true
		}
	}
	return false
}

// -1 represents C's NULL. Start at t itself, including a previously written
// NUL, and walk backwards over bytes; only a literal space is a separator.
// Preserve native C's empty backward traversal when t is NULL; its relational
// pointer comparison is nonportable: https://github.com/playbymail/olympia-g3/issues/359.
func prevWord(s []byte, t int) int {
	for t > 0 && s[t] != ' ' {
		t--
	}
	if t > 0 {
		return t
	}
	return -1
}

func (w *world) summaryTruncName(who int) string {
	s := sout("%s", w.justName(who))
	if len(s) <= truncName {
		return s
	}
	s = stripLeadingStupidWord(s)
	if len(s) <= truncName {
		return s
	}
	b := append([]byte(s), 0)
	t := prevWord(b, truncName)
	if t >= 0 {
		b[t] = 0
	}
	for {
		t = prevWord(b, t)
		if t < 0 || b[t] != ' ' {
			break
		}
		end := t + 1
		for b[end] != 0 {
			end++
		}
		if !stupidWord(string(b[t+1 : end])) {
			break
		}
		b[t] = 0
	}
	b[truncName] = 0
	end := 0
	for b[end] != 0 {
		end++
	}
	return string(b[:end])
}

func (w *world) sumFighters(who int) int32 {
	var sum int32
	for e := range w.loopInv(who) {
		p := w.rpItem(int(e.item))
		if p != nil && p.isManItem != 0 && (p.attack != 0 || p.defense != 0 || p.missile != 0 || int(e.item) == itemGhostWarrior) {
			if max(p.attack, max(p.defense, p.missile)) > 1 {
				sum += e.qty
			}
		}
	}
	return sum
}

func (w *world) unitSummarySup(pl, who int) {
	pr := w.isPrisoner(who)
	nam := w.summaryTruncName(who)
	var health, sick, behind, loyalty int8
	var rate int32
	if p := w.rpChar(who); p != nil {
		health, sick, behind, loyalty, rate = p.health, p.sick, p.behind, p.loyKind, p.loyRate
	}
	healthS := ""
	switch {
	case health == 100:
		healthS = "100 "
	case health == -1:
		healthS = "n/a "
	case sick != 0:
		healthS = sout("%d-", health)
	default:
		healthS = sout("%d+", health)
	}
	under := ""
	if pr {
		under = " ?? "
	} else if n := w.stackParent(who); n != 0 {
		under = boxCodeLess(n)
	}
	aura := ""
	if p := w.rpMagic(who); p != nil && p.magician != 0 {
		aura = sout("%d", p.curAura)
	}
	gold, peas := w.hasItem(who, itemGold), w.hasItem(who, itemPeasant)
	work, sail, fight := w.hasItem(who, itemWorker), w.hasItem(who, itemSailor), w.sumFighters(who)
	loyal := sout("%c%s", "ucofns"[loyalty], knum(rate, false))
	where := " ?? "
	if !pr {
		where = boxCodeLess(w.subloc(who))
	}
	w.commandOut(pl, "%-*s %-*s %-5s %4s%2d%3s %4s %4s %4s %4s %4s %-*s %s",
		charField, boxCodeLess(who), charField, where, loyal, healthS, behind, aura,
		knum(gold, true), knum(peas, true), knum(work, true), knum(sail, true), knum(fight, true), charField, under, nam)
	for i, n := range [...]int32{gold, peas, work, sail, fight} {
		w.unitSummaryTotals[i] += n
	}
}

func (w *world) unitSummary(pl int) {
	w.clearTemps(tChar)
	w.unitSummaryTotals = [5]int32{}
	w.commandOut(pl, "")
	count := len(w.pPlayer(pl).units)
	if count <= 0 {
		return
	}
	w.commandOut(pl, "{<b>}%-*s %-*s loyal heal B CA gold peas work sail figh %-*s name{</b>}", charField, "unit", charField, "where", charField, "under")
	w.commandOut(pl, "%-*s %-*s ----- ---- - -- ---- ---- ---- ---- ---- %-*s ----", charField, "----", charField, "-----", charField, "-----")
	for who := range w.loopUnits(pl) {
		w.unitSummarySup(pl, who)
	}
	if count > 1 {
		w.commandOut(pl, "%*s %-*s                 ==== ==== ==== ==== ====", charField, "", charField, "")
		s := w.unitSummaryTotals
		w.commandOut(pl, "%*s %-*s                 %4s %4s %4s %4s %4s", charField, "", charField, "", knum(s[0], true), knum(s[1], true), knum(s[2], true), knum(s[3], true), knum(s[4], true))
	}
}

// Character detail sections are complete report.c:410–678.
func (w *world) charRepStackInfo(who, num int) {
	if n := w.stackParent(num); n != 0 {
		w.wiout(who, 16, "Stacked under:  %s", w.boxName(n))
	}
	first := true
	assert(w.validBox(num), "olympia/report.c", 419, "valid_box(num)")
	for n := range w.loopHere(num) {
		if w.kind(n) == tChar && !w.isPrisoner(n) {
			if first {
				w.commandOut(who, "Stacked over:   %s", w.boxName(n))
				first = false
			} else {
				w.commandOut(who, "                %s", w.boxName(n))
			}
		}
	}
}

func (w *world) showPledged(who, num int) {
	if !w.pledgeBacklinks {
		w.collectPledgeBacklinks()
	}
	if p := w.rpMagic(num); p != nil {
		if p.pledge != 0 {
			w.wiout(who, 16, "Pledged to:     %s", w.boxName(int(p.pledge)))
		}
		for i, n := range p.pledgedToUs {
			if i == 0 {
				w.commandOut(who, "Pledged to us:  %s", w.boxName(n))
			} else {
				w.commandOut(who, "                %s", w.boxName(n))
			}
		}
	}
}

func (w *world) prisonerHealth(who int) string {
	var health int8
	if p := w.rpChar(who); p != nil {
		health = p.health
	}
	assert(health != 0, "olympia/report.c", 504, "health != 0")
	if health < 0 {
		return ""
	}
	return sout(", health %d", health)
}

func (w *world) charRepPrisoners(who, num int) {
	first := true
	assert(w.validBox(num), "olympia/report.c", 519, "valid_box(num)")
	for n := range w.loopHere(num) {
		if w.kind(n) == tChar && w.isPrisoner(n) {
			if first {
				w.commandOut(who, "Prisoners:      %s%s", w.boxName(n), w.prisonerHealth(n))
				first = false
			} else {
				w.commandOut(who, "                %s%s", w.boxName(n), w.prisonerHealth(n))
			}
		}
	}
}

func (w *world) charRepHealth(who, num int) {
	var health, sick int8
	if p := w.rpChar(num); p != nil {
		health, sick = p.health, p.sick
	}
	if health == -1 {
		w.commandOut(who, "Health:         n/a")
		return
	}
	s := ""
	if health > 0 && health < 100 {
		if sick != 0 {
			s = " (getting worse)"
		} else {
			s = " (getting better)"
		}
	}
	w.commandOut(who, "Health:         %d%%%s", health, s)
}

func (w *world) charRepCombat(who, num int) {
	var att, def, mis int16
	var behind, point int8
	if p := w.rpChar(num); p != nil {
		att, def, mis = p.attack, p.defense, p.missile
		behind, point = p.behind, p.breakPoint
	}
	if mk := w.nobleItem(num); mk != 0 {
		att, def, mis = 0, 0, 0
		if p := w.rpItem(mk); p != nil {
			att, def, mis = p.attack, p.defense, p.missile
		}
	}
	w.commandOut(who, "Combat:         attack %d, defense %d, missile %d", att, def, mis)
	s := " (front line in combat)"
	if behind != 0 {
		s = " (stay behind in combat)"
	}
	w.commandOut(who, "                behind %d %s", behind, s)
	s = ""
	if point == 0 {
		s = " (fight to the death)"
	} else if point == 100 {
		s = " (break almost immediately)"
	}
	w.commandOut(who, "Break point:    %d%%%s", point, s)
}

func (w *world) charRepMisc(who, num int) {
	p := w.rpMagic(num)
	if p != nil && p.abilityShroud != 0 {
		w.commandOut(who, "Ability shroud: %d aura", p.abilityShroud)
	}
	if w.hasSkill(num, skHideSelf) != 0 {
		var hidden int8
		if p != nil {
			hidden = p.hideSelf
		}
		s := "not concealing self"
		if hidden != 0 {
			if w.charReallyHidden(num) {
				s = "concealing self"
			} else {
				s = "concealing self, but not alone"
			}
		}
		w.commandOut(who, "use %4d %d      (%s)", skHideSelf, hidden, s)
	}
	if p != nil && p.visProtect > 0 {
		w.commandOut(who, "Receive Vision: %d protection", p.visProtect)
	}
}

func (w *world) charRepMagic(who, num int) {
	var ca, ma int32
	p := w.rpMagic(num)
	if p != nil {
		ca, ma = p.curAura, p.maxAura
	}
	mea := w.maxEffAura(num)
	if ma < 0 {
		ma = 0
	}
	w.commandOut(who, "")
	w.commandOut(who, "Current aura:   %d", ca)
	if ma < mea {
		w.commandOut(who, "Maximum aura:   %d (%d+%d)", mea, ma, mea-ma)
	} else {
		w.commandOut(who, "Maximum aura:   %d", ma)
	}
	if p != nil {
		if p.abilityShroud > 0 {
			w.commandOut(who, "Ability shroud: %d aura", p.abilityShroud)
		}
		if w.isLocOrShip(int(p.projectCast)) {
			w.commandOut(who, "Project cast:   %s", w.boxName(int(p.projectCast)))
		}
		if p.quickCast != 0 {
			w.commandOut(who, "Quicken cast:   %d", p.quickCast)
		}
		if p.hideMage > 0 {
			w.commandOut(who, "Appear common:  turn %d", int32(w.sysclock.turn)+int32(p.hideMage))
		}
	}
}

// showCarryCapacity is complete report.c:73–105. Ratios use the separate
// travel weights, including C int multiplication before division.
func (w *world) showCarryCapacity(who, num int) {
	w.commandOut(who, "")
	v := w.determineUnitWeights(num)
	walkPercent := ""
	if v.landCap > 0 {
		walkPercent = sout(" (%d%%)", v.landWeight*100/v.landCap)
	}
	buf := sout("%s/%s land%s", commaNum(int(v.landWeight)), commaNum(int(v.landCap)), walkPercent)
	if v.rideCap > 0 {
		buf += sout(", %s/%s ride (%d%%)", commaNum(int(v.rideWeight)), commaNum(int(v.rideCap)), v.rideWeight*100/v.rideCap)
	}
	if v.flyCap > 0 {
		buf += sout(", %s/%s fly (%d%%)", commaNum(int(v.flyWeight)), commaNum(int(v.flyCap)), v.flyWeight*100/v.flyCap)
	}
	w.wiout(who, int32(len("Capacity:  ")), "Capacity:  %s", buf)
}

// showItemSkillsSup is complete active report.c:108–193. Study and use are
// independent stored lists; validate each entry before its first heading.
func (w *world) showItemSkillsSup(who, item int, p *itemMagic) {
	first := true
	for _, sk := range p.mayStudy {
		assert(w.validBox(sk), "olympia/report.c", 126, "valid_box(sk)")
		parent := w.skillSchool(sk)
		if first {
			w.commandOut(who, "")
			w.commandOut(who, "%s permits study of the following skills:", w.boxName(item))
			w.indent += 3
			first = false
		}
		req := ""
		if sk != parent {
			req = sout(" (requires %s)", w.justName(parent))
		}
		if sk != parent && w.hasSkill(who, parent) < 1 {
			w.wiout(who, 3, "???%s", req)
		} else {
			w.wiout(who, 3, "%s%s", w.boxName(sk), req)
		}
	}
	if !first {
		w.indent -= 3
	}
	first = true
	for _, sk := range p.mayUse {
		assert(w.validBox(sk), "olympia/report.c", 162, "valid_box(sk)")
		parent := w.skillSchool(sk)
		if first {
			w.commandOut(who, "")
			w.commandOut(who, "%s permits use of the following skills:", w.boxName(item))
			w.indent += 3
			first = false
		}
		req := ""
		if sk != parent {
			req = sout(" (requires %s)", w.justName(parent))
		}
		if sk != parent && w.hasSkill(who, parent) < 1 {
			w.wiout(who, 3, "???%s", req)
		} else {
			w.wiout(who, 3, "%s%s", w.boxName(sk), req)
		}
	}
	if !first {
		w.indent -= 3
	}
}

// showItemSkills is report.c:196–210. Do not sort the inventory here or
// allocate missing magic records. Later roots retain their projections.
func (w *world) showItemSkills(who, num int) {
	for e := range w.loopInv(num) {
		if p := w.rpItemMagic(int(e.item)); p != nil {
			w.showItemSkillsSup(who, int(e.item), p)
		}
	}
}

// invItemComp is report.c:213–220. Sort the original list, including entries
// which loopInv will subsequently filter; subtraction has C int width.
func invItemComp(a, b *itemEnt) int { return int(a.item - b.item) }

// extraItemInfo is complete report.c:246–290. Capacity display chooses fly,
// ride, then land (not their sum). Aura visibility depends on the recipient,
// not the inventory owner, and is always suppressed for unclaimed items.
func (w *world) extraItemInfo(who, item, qty int) string {
	s := ""
	if p := w.rpItem(item); p != nil {
		if p.flyCap > 0 {
			s = sout("fly %s", niceNum(int(int32(p.flyCap)*int32(qty))))
		} else if p.rideCap > 0 {
			s = sout("ride %s", niceNum(int(int32(p.rideCap)*int32(qty))))
		} else if p.landCap > 0 {
			s = sout("cap %s", niceNum(int(int32(p.landCap)*int32(qty))))
		}
		if p.attack != 0 || p.defense != 0 || p.missile != 0 || item == itemGhostWarrior {
			s += sout(" (%d,%d,%d)", p.attack, p.defense, p.missile)
		}
	} else if item == itemGhostWarrior {
		s += " (0,0,0)"
	}
	if p := w.rpItemMagic(item); p != nil {
		if p.attackBonus != 0 {
			s += sout("+%d attack", p.attackBonus)
		}
		if p.defenseBonus != 0 {
			s += sout("+%d defense", p.defenseBonus)
		}
		if p.missileBonus != 0 {
			s += sout("+%d missile", p.missileBonus)
		}
		if p.auraBonus != 0 && who != 0 {
			if m := w.rpMagic(who); m != nil && m.magician != 0 {
				s += sout("+%d aura", p.auraBonus)
			}
		}
	}
	return s
}

// showCharInventory is complete report.c:294–344. Later report roots retain
// their projections until their own tasks; callers must use one path only.
func (w *world) showCharInventory(who, num int) {
	safeQsort(w.bx[num].items, invItemComp)
	first := true
	var total int32
	for e := range w.loopInv(num) {
		var weight int32
		if p := w.rpItem(int(e.item)); p != nil {
			weight = int32(p.weight) * e.qty
		}
		if first {
			w.commandOut(who, "")
			w.commandOut(who, "Inventory:")
			w.commandOut(who, "{<b>}%9s  %-30s %9s{</b>}", "qty", "name", "weight")
			w.commandOut(who, "%9s  %-30s %9s", "---", "----", "------")
			first = false
		}
		w.commandOut(who, "%9s  %-30s %9s  %s", commaNum(int(e.qty)), w.pluralItemBox(int(e.item), int(e.qty)), commaNum(int(weight)), w.extraItemInfo(who, int(e.item), int(e.qty)))
		total += weight
	}
	if !first {
		w.commandOut(who, "%9s  %-30s %9s", "", "", "======")
		w.commandOut(who, "%9s  %-30s %9s", "", "", commaNum(int(total)))
	} else {
		w.commandOut(who, "")
		w.commandOut(who, "%s has no possessions.", w.boxName(num))
	}
}

// showUnclaimed is complete report.c:746–780: no empty message or total row.
func (w *world) showUnclaimed(who, num int) {
	safeQsort(w.bx[num].items, invItemComp)
	first := true
	for e := range w.loopInv(num) {
		var weight int32
		if p := w.rpItem(int(e.item)); p != nil {
			weight = int32(p.weight) * e.qty
		}
		if first {
			w.commandOut(who, "")
			w.commandOut(who, "Unclaimed items:")
			w.commandOut(who, "")
			w.commandOut(who, "%9s  %-30s %9s", "qty", "name", "weight")
			w.commandOut(who, "%9s  %-30s %9s", "---", "----", "------")
			first = false
		}
		w.commandOut(who, "%9s  %-30s %9s  %s", commaNum(int(e.qty)), w.pluralItemBox(int(e.item), int(e.qty)), commaNum(int(weight)), w.extraItemInfo(0, int(e.item), int(e.qty)))
	}
}

// sortForOutput is report.c:11–27. Equal scratch positions use entity IDs;
// subtraction retains C int arithmetic rather than comparing as Go ints.
func (w *world) sortForOutput(l lists.IList) {
	safeQsort(l, func(a, b int) int {
		if w.bx[a].outputOrder != w.bx[b].outputOrder {
			return int(w.bx[a].outputOrder - w.bx[b].outputOrder)
		}
		return intComp(a, b)
	})
}

// determineOutputOrder is report.c:30–70. Only descendants of regions get
// new scratch positions; unlocated boxes retain their previous positions.
// Faction units, unlike outputOrder itself, are persistent database state.
func (w *world) determineOutputOrder() {
	var count int32
	for reg := range w.loopKind(tLoc) {
		if w.locDepth(reg) != locRegion {
			continue
		}
		w.bx[reg].outputOrder = count
		count++
		for _, i := range w.allHere(reg) {
			w.bx[i].outputOrder = count
			count++
		}
	}
	for pl := range w.loopKind(tPlayer) {
		if p := w.rpPlayer(pl); p != nil {
			w.sortForOutput(p.units)
		}
	}
}

// playerEntInfo is complete report.c:783–798. Each renderer owns its original
// allocations, in-place sorts and output; silent players skip all three.
func (w *world) playerEntInfo() {
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) == subPlSilent {
			continue
		}
		w.printAdmit(pl)
		w.printAtt(pl, pl)
		w.showUnclaimed(pl, pl)
	}
}

// collectPledgeBacklinks is report.c:438–463. The two passes are deliberate:
// p_magic allocates CH records even for silent/prisoner characters. Clearing
// a previously allocated list retains its pointer, so print_magic can emit a
// bare CM section even though pledged_to_us itself is not serialized.
func (w *world) collectPledgeBacklinks() {
	w.pledgeBacklinks = true
	for i := range w.loopKind(tChar) {
		w.pMagic(i).pledgedToUs.Clear()
	}
	for i := range w.loopKind(tChar) {
		if p := w.rpMagic(i); p != nil && p.pledge != 0 {
			w.pMagic(int(p.pledge)).pledgedToUs.Append(i)
		}
	}
}

// charRepSup is complete report.c:681–710. Each accepted subordinate owns
// its output and side effects; in particular showPledged collects lazily.
func (w *world) charRepSup(who, num int) {
	w.wiout(who, 16, "Location:       %s", w.charRepLocation(num))
	w.commandOut(who, "Loyalty:        %s", capitalize(w.loyalS(num)))
	w.charRepStackInfo(who, num)
	w.charRepHealth(who, num)
	w.charRepCombat(who, num)
	w.charRepMisc(who, num)
	if p := w.rpMisc(num); p != nil && p.display != nil {
		w.commandOut(who, "Banner:         %s", *p.display)
	}
	w.showPledged(who, num)
	if p := w.rpMagic(num); p != nil && p.magician != 0 {
		w.charRepMagic(who, num)
	}
	w.charRepPrisoners(who, num)
	w.printAtt(who, num)
	w.listSkills(who, num)
	w.listPartialSkills(who, num)
	w.showCharInventory(who, num)
	w.showCarryCapacity(who, num)
	w.showItemSkills(who, num)
	w.listPendingTrades(who, num)
	w.commandOut(who, "")
}

// characterReport is complete report.c:713–743. Silent/prisoner-only worlds
// never reach showPledged. The temporary prisoner toggle enables only the
// two original lines, using normal output routing and allocation timing.
func (w *world) characterReport() {
	w.stage("character_report()")
	w.indent += 3
	for who := range w.loopKind(tChar) {
		if w.subkind(w.player(who)) == subPlSilent {
			continue
		}
		if w.isPrisoner(who) {
			w.pChar(who).prisoner = 0
			w.commandOut(who, "%s is being held prisoner.", w.boxName(who))
			w.commandOut(who, "")
			w.pChar(who).prisoner = 1
		} else {
			w.commandOut(who, "")
			w.charRepSup(who, who)
		}
	}
	w.indent -= 3
}

// turnEndLocReports is complete report.c:1589–1642, including the original
// global resets and allocation before sorted touched-location iteration.
func (w *world) turnEndLocReports() {
	w.stage("turn_end_loc_reports()")
	w.outPath = master
	w.showLocNoHeader = true
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) == subPlSilent {
			continue
		}
		separate := false
		if p := w.rpPlayer(pl); p != nil {
			separate = p.format != 0
		}
		p := w.pPlayer(pl)
		sortKnown(p.locs)
		for _, i := range p.locs {
			if !w.validBox(i) {
				continue
			}
			if separate {
				w.outAltWho = 4 // sout.h OUT_LOC
				w.lines(pl, w.showLocHeader(i))
				w.showLoc(pl, i)
				w.commandOut(pl, "")
			} else {
				w.outAltWho = i
				w.commandOut(pl, "")
				w.showLoc(pl, i)
			}
		}
	}
	w.outPath = 0
	w.outAltWho = 0
	w.showLocNoHeader = false
}

// playerBanner is complete report.c:1645–1722. The initial HTML output
// precedes p_player, preserving the cached output owner's allocation timing.
func (w *world) playerBanner() {
	w.stage("player_banner()")
	w.outPath = master
	w.outAltWho = 0 // sout.h OUT_BANNER
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) == subPlSilent {
			continue
		}
		w.html(pl, "<pre>")
		p := w.pPlayer(pl)
		if p.compuserve != 0 {
			w.indent += 3
			w.commandOut(pl, "Note:  Olympia Times mailing is turned off for you.To begin receiving the Times via email again, issue the order TIMES 0")
			w.indent -= 3
			w.commandOut(pl, "")
		}
		w.html(pl, "<center>")
		if w.winFlag {
			w.html(pl, "<img src=\"head.gif\" align=middle width=100 height=100 alt=\"\">")
		} else {
			w.html(pl, "<img src=\"../head.gif\" align=middle width=100 height=100 alt=\"\">")
		}
		w.html(pl, "<h1>")
		w.commandOut(pl, "%s turn %d", w.gameTitle, w.sysclock.turn)
		w.commandOut(pl, "Report for %s.", w.boxName(pl))
		w.html(pl, "</h1>")
		month, year := olyMonth(w.sysclock), olyYear(w.sysclock)
		if month < 0 {
			month = 0
		}
		assert(month >= 0 && month < numMonths, "olympia/report.c", 1706, "month >= 0 && month < NUM_MONTHS")
		w.commandOut(pl, "{<i>}Season \"%s\", month %d, in the year %d.{</i>}", monthNames[month], month+1, year+1)
		w.html(pl, "</center>")
		w.commandOut(pl, "")
	}
	w.outPath = 0
	w.outAltWho = 0
}

// reportOutputState preserves a MASTER out()'s open-before-allocation path
// (sout.c:295–306,465–467), not its stage-2 body or transient output bits.
// No blank/successful rendering substitute is emitted. Like ordinary output,
// this uses openPlayerOutput, so an inaccessible log is still a real failure.
func (w *world) reportOutputState(pl int) {
	if w.immediate || pl == 0 || w.isPrisoner(pl) {
		return
	}
	w.openPlayerOutput(pl)
	w.pPlayer(pl)
}

// reportAccountSup is complete active report.c:1725–1791. The external
// account command and file read are #if 0; only the footer is active.
func (w *world) reportAccountSup(pl int) {
	w.style(0)
	w.html(pl, "<hr><center><font size=+1>")
	w.html(pl, "<a href=\"%s\">Rules</a>   <a href=\"%s\"><i>The Olympia Times</i></a>", w.rulesURL, w.timesURL)
	w.html(pl, "</font></center>")
	w.commandOut(pl, "")
}

// reportAccount is complete report.c:1793–1815, not external -A charging.
func (w *world) reportAccount() {
	w.stage("report_account()")
	w.outPath = master
	w.outAltWho = 0 // sout.h OUT_BANNER
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) == subPlRegular {
			w.reportAccountSup(pl)
		}
	}
	w.outPath = 0
	w.outAltWho = 0
}

// playerReportSup is complete report.c:1188–1212.
func (w *world) playerReportSup(pl int) {
	if w.subkind(pl) == subPlSystem {
		return
	}
	p := w.pPlayer(pl)
	w.commandOut(pl, "Noble points:  %d     (%d gained, %d spent)", p.noblePoints, p.npGained, p.npSpent)
	w.printHiringStatus(pl)
	w.printUnformed(pl)
	if p.fastStudy > 0 {
		w.commandOut(pl, "%d fast study day%s are left.", p.fastStudy, addS(int(p.fastStudy)))
		w.commandOut(pl, "")
	}
}

// playerReport is complete report.c:1323–1350. Each accepted table runs
// once, in original order, borrowing the ordinary MASTER output owner.
func (w *world) playerReport() {
	w.stage("player_report()")
	w.outPath = master
	w.outAltWho = 0
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) == subPlSilent {
			continue
		}
		w.playerReportSup(pl)
		w.unitSummary(pl)
		w.locStackReport(pl)
		w.stackCapacityReport(pl)
		w.stormReport(pl)
		w.garrisonSummary(pl)
		w.commandOut(pl, "")
	}
	w.outPath = 0
	w.outAltWho = 0
}

// genIncludeSupState preserves report.c:1423–1566's existing-output input
// validation. It must not synthesize memberships for omitted report bodies.
// rep_char (1371–1389) sorts a COPY, then p_misc allocates corpse records;
// that happens after the initial MASTER output, never before validation.
// Include body text, lines(), and rep_loc/show_loc_header are stage 2.
func (w *world) genIncludeSupState(pl int) {
	p := w.pPlayer(pl)
	sortKnown(p.output)
	var chars lists.IList
	for _, n := range p.output {
		switch n {
		case 0, 1, 2, 3, 4, 5, 6, 7, 10, 11, 12, 13, 14, 20, 21, 22, 23, 24, 25:
			// OUT_BANNER/INCLUDE/LORE/NEW/LOC/TEMPLATE/GARR/SHOW_POSTS,
			// LOG_CODE/SPECIAL/DEATH/MISC/DROP and EAT_* (sout.h).
			continue
		}
		if !w.validBox(n) {
			continue
		}
		switch w.kind(n) {
		case tChar, tDeadchar:
			chars.Append(n)
		case tLoc, tShip:
			// Only deferred location-header text.
		case tPlayer:
			assert(n == pl, "olympia/report.c", 1495, "n == pl")
		case tItem:
			if w.subkind(n) == subDeadBody {
				chars.Append(n)
			} else {
				assert(false, "olympia/report.c", 1503, "FALSE")
			}
		default:
			assert(false, "olympia/report.c", 1507, "FALSE")
		}
	}
	w.reportOutputState(pl)
	w.sortForOutput(chars)
	for _, n := range chars {
		if w.subkind(n) == subDeadBody {
			w.pMisc(n)
		}
	}
}

// genIncludeSectionState is report.c:1568–1586's nonsilent loop. The helper
// remains a state/input projection, not a successful include renderer.
func (w *world) genIncludeSectionState() {
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) != subPlSilent {
			w.genIncludeSupState(pl)
		}
	}
}

// charRepLocation is report.c:357–407. Death logging needs this helper;
// report-body assembly remains separate stage-2 work.
func (w *world) charRepLocation(who int) string {
	where := w.subloc(who)
	if where == 0 {
		return "nowhere"
	}
	s := ""
	for w.locDepth(where) > locProvince {
		if s != "" {
			s = sout("%s, in %s", s, w.boxName(where))
		} else {
			s = w.boxName(where)
		}
		where = w.loc(where)
	}
	if w.subkind(w.province(where)) == subOcean {
		if s != "" {
			s = sout("%s, in %s", s, w.boxName(where))
		} else {
			s = w.boxName(where)
		}
	} else {
		if s != "" {
			s = sout("%s, in province %s", s, w.boxName(where))
		} else {
			s = w.boxName(where)
		}
	}
	if regName := w.name(w.region(where)); regName != "" {
		s = sout("%s, in %s", s, regName)
	}
	return s
}
