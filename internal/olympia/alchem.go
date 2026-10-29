// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// newPotion is alchem.c:8–42 new_potion. Minting, ownership and the keyed
// magic draw use the existing shared owners; brewing/item callbacks are separate.
func (w *world) newPotion(who int) int {
	n := w.createUniqueItem(who, 0)
	if n < 0 {
		return -1
	}
	var name string
	switch w.streams.MagcPotion(who, n, 0, 1, 2) {
	case 1:
		name = "Magic potion"
	case 2:
		name = "Strange potion"
	default:
		assert(false, "olympia/alchem.c", 30, "FALSE")
	}
	w.setName(n, name)
	p := w.pItemMagic(n)
	p.creator = int32(who)
	p.regionCreated = int32(w.province(who))
	w.pItem(n).weight = 1
	w.commandOut(who, "Produced one %s", w.boxName(n))
	return n
}
