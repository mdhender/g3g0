// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "github.com/mdhender/g3g0/internal/lists"

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

// playerEntInfoState preserves the database effects of player_ent_info
// (report.c:783–799), not its deferred stage-2 report artifacts. print_admit
// allocates the player record and sorts the original admissions; print_att
// sorts all three original attitude lists; show_unclaimed sorts inventory.
// Invalid references remain in these lists, even when rendering/save writers
// filter them. Admission members are not sorted. Silent players do none of it.
func (w *world) playerEntInfoState() {
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) == subPlSilent {
			continue
		}
		p := w.pPlayer(pl)
		safeQsort(p.admits, func(a, b *admit) int { return int(a.targ - b.targ) })
		w.reportAttitudeState(pl)
		safeQsort(w.bx[pl].items, func(a, b *itemEnt) int { return int(a.item - b.item) })
		// show_unclaimed's plural formatter can warn on stderr, which is
		// observable even though the inventory report itself is deferred.
		for e := range w.loopInv(pl) {
			w.pluralItemBox(int(e.item), int(e.qty))
		}
	}
}

// reportAttitudeState is the persistent part of perm.c:427–498. A missing
// attitude record is not allocated, and invalid entries are not removed.
func (w *world) reportAttitudeState(n int) {
	if p := w.rpDisp(n); p != nil {
		safeQsort(p.hostile, intComp)
		safeQsort(p.neutral, intComp)
		safeQsort(p.defend, intComp)
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
