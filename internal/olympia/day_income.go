// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// innIncome is the complete day.c:1421–1489 definition, excluding phase
// instrumentation. Sharing and pillage use the immediate sublocation, not
// the province; the traveller bonus is added after both integer divisions.
func (w *world) innIncome() {
	for i := range w.loopSubkind(subInn) {
		owner := w.buildingOwner(i)
		if owner == 0 {
			continue
		}
		where := w.subloc(i)
		n := w.countLocStructures(where, subInn, 0)
		base := int32(w.streams.UpIncome(i, 0, 50, 75))
		var pil int32
		if p := w.rpSubloc(where); p != nil {
			pil = int32(p.loot)
		}
		if pil != 0 {
			base /= pil + 1
		}
		base /= n
		if pil == 0 && w.streams.UpIncome(i, 1, 1, 8) == 1 {
			bonus := int32(w.streams.UpIncome(i, 2, 5, 13)) * 10
			w.commandOut(owner, "A rich traveller stayed in %s this month, spending %s.", w.boxName(i), goldS(int(bonus)))
			base += bonus
		}
		w.genItem(owner, itemGold, int(base))
		w.goldInn += base
		w.commandOut(owner, "%s yielded %s in income.", w.boxName(i), goldS(int(base)))
		if pil != 0 {
			switch w.streams.UpIncome(i, 3, 1, 3) {
			case 1:
				w.commandOut(owner, "Patrons were scared away by recent looting in the province.")
			case 2:
				w.commandOut(owner, "Profits were hurt by pillaging in the area.")
			case 3:
				w.commandOut(owner, "Recent pillaging in the area lowered profits.")
			}
		}
	}
}

// templeIncome is the complete day.c:1492–1522 definition. There is no
// sharing, capacity or occupant-count cap: only the building owner's religion
// skill gates the flat offering (before immediate-sublocation pillage).
func (w *world) templeIncome() {
	for i := range w.loopSubkind(subTemple) {
		owner := w.buildingOwner(i)
		if owner == 0 || w.hasSkill(owner, skReligion) == 0 {
			continue
		}
		amount := int32(100)
		if p := w.rpSubloc(w.subloc(i)); p != nil && p.loot != 0 {
			amount /= int32(p.loot) + 1
		}
		w.genItem(owner, itemGold, int(amount))
		w.goldTemple += amount
		w.commandOut(owner, "%s collected offerings of %s.", w.boxName(i), goldS(int(amount)))
	}
}
