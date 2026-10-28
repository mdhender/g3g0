// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "github.com/mdhender/g3g0/internal/lists"

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

// summaryProvincesState is the persistent projection of summary.c:348–389.
// Unlike other report loops it allocates records for ALL player subtypes,
// including silent players. Province counts and ranking lists are scratch.
func (w *world) summaryProvincesState() {
	for pl := range w.loopKind(tPlayer) {
		sortKnown(w.pPlayer(pl).known) // loop_known sorts the original list.
	}
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

// summaryReportState projects summary.c:140–389,438–502. Totals, rankings,
// skill/spell/land counters and their render-only traversal are scratch and
// deferred with report text. Regular output opens BEFORE summary_provinces
// allocates every player and sorts their original knowledge lists.
func (w *world) summaryReportState() {
	w.stage("summary_report()")
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) == subPlRegular {
			w.reportOutputState(pl)
		}
	}
	w.summaryProvincesState()
}

// playerReportState projects report.c:803–1353. player_report_sup allocates
// nonsystem records before output; unit_summary allocates even system records
// after its initial output. Unit/location-stack/capacity/storm/garrison text,
// scratch links and render-only validation remain stage 2. None changes saved
// units, inventory, skills, commands, ownership, storm or garrison records.
func (w *world) playerReportState() {
	w.stage("player_report()")
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) == subPlSilent {
			continue
		}
		if w.subkind(pl) != subPlSystem {
			w.pPlayer(pl)
		}
		w.reportOutputState(pl)
		w.pPlayer(pl)
	}
}

// gmReportState projects gm.c:1–1048, not gm_show_all_skills (lore.go).
// Statistics, sorted scratch copies, item/notice text and render-only entity
// traversals remain stage 2. gate_stats/locs_visited allocate ALL players and
// loop_known sorts the stored lists; neither pass grants knowledge. No active
// input files are read by this closure. list_all_items outputs to skill_player
// last, independently of the GM recipient and its subtype.
func (w *world) gmReportState(pl int) {
	w.stage("gm_report()")
	w.reportOutputState(pl)
	for i := range w.loopKind(tPlayer) {
		sortKnown(w.pPlayer(i).known)
	}
	w.reportOutputState(int(w.skillPlayer))
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
