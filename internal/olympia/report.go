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

// characterReportState preserves character_report/char_rep_sup's persistent
// effects (report.c:681–743), not their stage-2 report text. show_pledged is
// lazy: an all-silent/all-prisoner world must not collect backlinks. Prisoner
// output temporarily toggles prisoner off and on, leaving saved state alone.
func (w *world) characterReportState() {
	w.stage("character_report()")
	for who := range w.loopKind(tChar) {
		pl := w.player(who)
		if w.subkind(pl) == subPlSilent {
			continue
		}
		// The initial out() (also with the prisoner's temporary toggle)
		// reaches sout.c:306's p_player before char_rep_sup does anything.
		if !w.immediate {
			if w.subkind(who) == subGarrison {
				for _, ruler := range w.playersWhoRuleHere(who) {
					if ruler != 0 && w.subkind(ruler) != subPlSilent {
						w.pPlayer(ruler)
					}
				}
			} else if pl != 0 {
				w.pPlayer(pl)
			}
		}
		if w.isPrisoner(who) {
			continue
		}
		w.charRepLocation(who)
		w.loyalS(who)
		if !w.pledgeBacklinks {
			w.collectPledgeBacklinks()
		}
		for _, n := range w.rpLocInfo(who).hereList {
			if w.kind(n) == tChar && w.isPrisoner(n) {
				assert(w.rpChar(n).health != 0, "olympia/report.c", 504, "health != 0")
			}
		}
		w.reportAttitudeState(who)
		// list_skills and list_partial_skills sort COPIES. Preserve their
		// assertion/diagnostic path, but never reorder the stored skills.
		if p := w.rpChar(who); p != nil {
			l := p.skills.Copy()
			safeQsort(l, func(a, b *skillEnt) int {
				if a.know != skillKnow && b.know == skillKnow {
					return -1
				}
				if b.know != skillKnow && a.know == skillKnow {
					return 1
				}
				pa, pb := w.skillSchool(int(a.skill)), w.skillSchool(int(b.skill))
				if pa != pb {
					return intComp(pa, pb)
				}
				return int(a.skill - b.skill)
			})
			for _, e := range l {
				if e.know == skillKnow {
					if s := w.rpSkill(int(e.skill)); s == nil || s.noExp == 0 {
						w.skillSchool(int(e.skill))
					}
				}
			}
			l = p.skills.Copy()
			safeQsort(l, func(a, b *skillEnt) int { return int(a.skill - b.skill) })
			for _, e := range l {
				if e.know != skillKnow && e.know != skillDont {
					assert(e.know == skillLearning, "olympia/use.c", 1358, "p->know == SKILL_learning")
				}
			}
		}
		safeQsort(w.bx[who].items, func(a, b *itemEnt) int { return int(a.item - b.item) })
		for e := range w.loopInv(who) {
			w.pluralItemBox(int(e.item), int(e.qty))
		}
		w.determineUnitWeights(who)
		for e := range w.loopInv(who) {
			if p := w.rpItemMagic(int(e.item)); p != nil {
				for _, sk := range p.mayStudy {
					assert(w.validBox(sk), "olympia/report.c", 126, "valid_box(sk)")
					w.skillSchool(sk)
				}
				for _, sk := range p.mayUse {
					assert(w.validBox(sk), "olympia/report.c", 162, "valid_box(sk)")
					w.skillSchool(sk)
				}
			}
		}
	}
}

// turnEndLocReportsState preserves report.c:1589–1642's nonsilent p_player
// allocations BEFORE player_ent_info. Location report artifacts stay stage 2.
func (w *world) turnEndLocReportsState() {
	w.stage("turn_end_loc_reports()")
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) != subPlSilent {
			w.pPlayer(pl)
		}
	}
}

// playerBannerState is report.c:1645–1722's allocation projection. It occurs
// after character_report, before summary_provinces allocates silent records.
func (w *world) playerBannerState() {
	w.stage("player_banner()")
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) != subPlSilent {
			w.pPlayer(pl)
		}
	}
}

// summaryProvincesState is the persistent projection of summary.c:348–389.
// Unlike other report loops it allocates records for ALL player subtypes,
// including silent players. Province counts and ranking lists are scratch.
func (w *world) summaryProvincesState() {
	for pl := range w.loopKind(tPlayer) {
		w.pPlayer(pl)
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
