// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// collapsedMineDecay is day.c:1525–1553. Assignment to the narrow field
// precedes the clamp; the subtype iterator saves next before deletion.
func (w *world) collapsedMineDecay() {
	for i := range w.loopSubkind(subMineCollapsed) {
		p := w.pMisc(i)
		p.mineDelay--
		if p.mineDelay < 0 {
			p.mineDelay = 0
		}
		if p.mineDelay == 0 {
			w.getRidOfCollapsedMine(i)
		}
	}
}

// ghostWarriorDecay is day.c:1556–1588: exactly one, non-NPCs only.
func (w *world) ghostWarriorDecay() {
	for i := range w.loopKind(tChar) {
		if w.isNPC(i) || w.hasItem(i, itemGhostWarrior) <= 0 {
			continue
		}
		w.commandOut(i, "%s evaporated.", capitalize(w.boxNameQty(itemGhostWarrior, 1)))
		w.consumeItem(i, itemGhostWarrior, 1)
	}
}

// corpseDecay is day.c:1591–1623; the persistent upkeep leaf is keyed by
// character, not a shared sequential draw or a freshly initialized stream.
func (w *world) corpseDecay() {
	for i := range w.loopKind(tChar) {
		if w.isNPC(i) {
			continue
		}
		has := int(w.hasItem(i, itemCorpse))
		if has <= 0 {
			continue
		}
		has = min(has, w.streams.UpCorpse(i, 0, 2))
		if has != 0 {
			w.commandOut(i, "%s decomposed.", capitalize(w.boxNameQty(itemCorpse, has)))
			w.consumeItem(i, itemCorpse, has)
		}
	}
}

// deadBodyRot is day.c:1626–1652. p_char intentionally allocates even
// for a body with no character record, before testing the exact 12-turn age.
func (w *world) deadBodyRot() {
	for i := range w.loopSubkind(subDeadBody) {
		owner := w.itemUnique(i)
		assert(owner != 0, "olympia/day.c", 1638, "owner")
		t := int(w.sysclock.turn) - int(w.pChar(i).deathTime.turn)
		if t < 12 {
			continue
		}
		if w.kind(owner) == tChar {
			w.commandOut(owner, "%s decomposed.", w.boxName(i))
		}
		w.destroyUniqueItem(owner, i)
	}
}
