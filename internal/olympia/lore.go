// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// setKnown is u.c set_known. Invalid entities are silently ignored by C.
func (w *world) setKnown(who, i int) {
	if !w.validBox(who) || !w.validBox(i) {
		return
	}
	pl := w.player(who)
	if !w.validBox(pl) {
		return
	}
	setBit(&w.pPlayer(pl).known, i)
}

// queueLore preserves insertion order and anyway duplicates. Rendering and
// duplicate removal at delivery are stage 2, not part of queue mutation.
func (w *world) queueLore(who, num int, anyway bool) {
	pl := w.player(who)
	assert(w.kind(pl) == tPlayer, "olympia/lore.c", 236, "kind(pl) == T_player")
	if !anyway && w.testKnown(pl, num) {
		return
	}
	w.pPlayer(pl).deliverLore.Append(num)
	w.setKnown(pl, num)
}

func (w *world) scanCharItemLore() {
	for who := range w.loopKind(tChar) {
		for e := range w.loopInv(who) {
			p := w.rpItemMagic(int(e.item))
			if p != nil && p.lore != 0 && !w.testKnown(who, int(e.item)) {
				w.queueLore(who, int(e.item), false)
			}
		}
	}
}
