// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Active summary.c:18–176,283–345 at e733acb. SHOW_TOP is FALSE:
// top_rank is deliberately not evaluated, even for an empty ranks list.
// The summary root and its persistent projection remain owned by #158.
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
