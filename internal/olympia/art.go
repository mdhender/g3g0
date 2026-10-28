// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// hasAuraculum is art.c:7–18 has_auraculum.
func (w *world) hasAuraculum(who int) int {
	var ac int
	if p := w.rpMagic(who); p != nil {
		ac = int(p.auraculum)
	}
	if ac != 0 && w.validBox(ac) && w.hasItem(who, ac) > 0 {
		return ac
	}
	return 0
}

// maxEffAura is art.c:25–50 max_eff_aura. Bonuses are per positive
// inventory entry, not per item quantity; only innate aura is clamped.
func (w *world) maxEffAura(who int) int32 {
	var a int32
	if p := w.rpMagic(who); p != nil {
		a = p.maxAura
	}
	if a < 0 {
		a = 0
	}
	if ac := w.hasAuraculum(who); ac != 0 {
		if p := w.rpItemMagic(ac); p != nil {
			a += int32(p.aura)
		}
	}
	assert(w.validBox(who), "olympia/art.c", 41, "valid_box(who)")
	for e := range w.loopInv(who) {
		if p := w.rpItemMagic(int(e.item)); p != nil {
			a += int32(p.auraBonus)
		}
	}
	return a
}
