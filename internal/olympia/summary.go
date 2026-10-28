// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Active summary.c at e733acb. SHOW_TOP is FALSE:
// top_rank is deliberately not evaluated, even for an empty ranks list.
func (w *world) rankComp(a, b int) int {
	return int(w.bx[b].temp - w.bx[a].temp)
}

func (w *world) ranking(n int) int {
	for n > 0 && w.bx[w.summaryRanks[n]].temp == w.bx[w.summaryRanks[n-1]].temp {
		n--
	}
	return n + 1
}

func (w *world) outRankMineTop(who, num int, title string, total int32) {
	val := w.bx[w.summaryRanks[num]].temp
	s := commaNum(int(val))
	if total != 0 {
		s = sout("%s/%s", s, commaNum(int(total)))
	}
	w.wiout(who, 20+1+11+2+5+2, "%-20s %16s  %-5s  %s", title, s, ordinal(w.ranking(num)), "")
}

func (w *world) outRanking(title string, total int32) {
	w.summaryRanks.Clear()
	for i := range w.loopKind(tPlayer) {
		if w.subkind(i) != subPlRegular {
			continue
		}
		w.summaryRanks.Append(i)
	}
	safeQsort(w.summaryRanks, w.rankComp)
	for i, pl := range w.summaryRanks {
		w.outRankMineTop(pl, i, title, total)
	}
}

func (w *world) collectGameTotals() {
	w.summaryNPlayers, w.summaryNControlled, w.summaryNOther = 0, 0, 0
	w.clearTemps(tChar)
	for i := range w.loopKind(tPlayer) {
		if w.subkind(i) != subPlRegular {
			continue
		}
		w.summaryNPlayers++
		for j := range w.loopUnits(i) {
			w.summaryNControlled++
			w.bx[j].temp = 1
		}
	}
	for i := range w.loopKind(tChar) {
		if w.bx[i].temp == 0 {
			w.summaryNOther++
		}
	}
}

func (w *world) summaryUnits() {
	w.clearTemps(tPlayer)
	for i := range w.loopKind(tChar) {
		w.bx[w.player(i)].temp++
	}
	w.outRanking("Characters:", 0)
}

func (w *world) summaryMen() {
	w.clearTemps(tPlayer)
	for i := range w.loopKind(tChar) {
		for e := range w.loopInv(i) {
			if p := w.rpItem(int(e.item)); p != nil && p.isManItem != 0 {
				w.bx[w.player(i)].temp += e.qty
			}
		}
	}
	w.outRanking("Men:", 0)
}

func (w *world) summaryGold() {
	w.clearTemps(tPlayer)
	for i := range w.loopKind(tChar) {
		w.bx[w.player(i)].temp += w.hasItem(i, itemGold)
	}
	w.outRanking("Gold:", 0)
}

func (w *world) summaryLandOwned() {
	w.clearTemps(tPlayer)
	for i := range w.loopSubkind(subGarrison) {
		if w.garrisonStrength(i) < 10 {
			continue
		}
		for _, pl := range w.playersWhoRuleHere(i) {
			if pl != 0 {
				w.bx[pl].temp++
			}
		}
	}
	w.outRanking("Land controlled:", 0)
}

// summarySkills and summarySpells are complete summary.c:179–280.
func (w *world) summarySkills() {
	var nskills int32
	for i := range w.loopKind(tSkill) {
		if w.subkind(w.skillSchool(i)) == subMagic {
			continue
		}
		nskills++
	}
	w.clearTemps(tPlayer)
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) != subPlRegular {
			continue
		}
		w.clearTemps(tSkill)
		for who := range w.loopUnits(pl) {
			assert(w.validBox(who), "olympia/summary.c", 207, "valid_box(who)")
			if p := w.rpChar(who); p != nil {
				for _, e := range p.skills.Copy() {
					if e.know == skillKnow {
						w.bx[e.skill].temp++
					}
				}
			}
		}
		for i := range w.loopKind(tSkill) {
			if w.subkind(w.skillSchool(i)) == subMagic {
				continue
			}
			if w.bx[i].temp != 0 {
				w.bx[pl].temp++
			}
		}
	}
	w.outRanking("Skills known:", nskills)
}

func (w *world) summarySpells() {
	var nskills int32
	for i := range w.loopKind(tSkill) {
		if w.subkind(w.skillSchool(i)) != subMagic {
			continue
		}
		nskills++
	}
	w.clearTemps(tPlayer)
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) != subPlRegular {
			continue
		}
		w.clearTemps(tSkill)
		for who := range w.loopUnits(pl) {
			assert(w.validBox(who), "olympia/summary.c", 259, "valid_box(who)")
			if p := w.rpChar(who); p != nil {
				for _, e := range p.skills.Copy() {
					if e.know == skillKnow {
						w.bx[e.skill].temp++
					}
				}
			}
		}
		for i := range w.loopKind(tSkill) {
			if w.subkind(w.skillSchool(i)) != subMagic {
				continue
			}
			if w.bx[i].temp != 0 {
				w.bx[pl].temp++
			}
		}
	}
	w.outRanking("Spells known:", nskills)
}

// summaryProvinces is complete summary.c:348–389. loop_known sorts the
// original stored list, and p_player allocates for ALL player subtypes.
func (w *world) summaryProvinces() {
	var nlocs int32
	for i := range w.loopKind(tLoc) {
		if w.locDepth(i) == locProvince {
			nlocs++
		}
	}
	w.clearTemps(tPlayer)
	for pl := range w.loopKind(tPlayer) {
		w.clearTemps(tLoc)
		known := w.pPlayer(pl).known
		sortKnown(known)
		for _, i := range known {
			w.bx[i].temp++
		}
		for i := range w.loopKind(tLoc) {
			if w.locDepth(i) != locProvince {
				continue
			}
			if w.bx[i].temp != 0 {
				w.bx[pl].temp++
			}
		}
	}
	w.outRanking("Provinces visited:", nlocs)
}

// summaryReport is complete summary.c:438–502. summary_sublocs is #if 0.
func (w *world) summaryReport() {
	w.stage("summary_report()")
	w.collectGameTotals()
	w.outPath, w.outAltWho = master, 0 // OUT_SUMMARY
	w.indent += 3
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) != subPlRegular {
			continue
		}
		w.commandOut(pl, "")
		w.commandOut(pl, "Game totals:")
		w.indent += 3
		w.commandOut(pl, "%-20s  %5s", "Players:", commaNum(int(w.summaryNPlayers)))
		w.commandOut(pl, "%-20s  %5s", "Controlled units:", commaNum(int(w.summaryNControlled)))
		w.commandOut(pl, "%-20s  %5s", "Other units:", commaNum(int(w.summaryNOther)))
		w.indent -= 3
		w.commandOut(pl, "")
		s := sout("Faction %s", w.boxCode(pl))
		w.commandOut(pl, "{<b>}%-20s %16s  %-5s  %s{</b>}", s, "", "rank", "")
		line := make([]byte, len(s))
		for i := range line {
			line[i] = '-'
		}
		w.commandOut(pl, "%-20s %16s  %-5s  %s", string(line), "", "----", "")
	}
	w.summaryUnits()
	w.summaryMen()
	w.summaryGold()
	w.summaryLandOwned()
	w.summarySkills()
	w.summarySpells()
	w.summaryProvinces()
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) != subPlRegular {
			continue
		}
		w.commandOut(pl, "")
		w.commandOut(pl, "")
	}
	w.indent -= 3
	w.outPath, w.outAltWho = 0, 0
}
