// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "testing"

// TestAccessors checks the p_xxx and rp_xxx semantics TestEntityMatchesC
// compares with C: a p_xxx allocates a zeroed substructure once and then
// returns it, and an rp_xxx never allocates.
func TestAccessors(t *testing.T) {
	w := newWorld()
	w.bx[1001] = &box{kind: tChar}

	if w.rpChar(1001) != nil || w.rpMagic(1001) != nil {
		t.Fatal("new box has an entity_char")
	}
	if w.rpChar(1001) != nil {
		t.Fatal("rp_char allocated")
	}
	m := w.pMagic(1001)
	c := w.rpChar(1001)
	if m == nil || c == nil || c.xCharMagic != m {
		t.Fatal("p_magic did not allocate the entity_char and char_magic")
	}
	if w.pChar(1001) != c || w.pMagic(1001) != m || w.rpMagic(1001) != m {
		t.Fatal("p_xxx allocated again")
	}

	if w.rpItemMagic(1001) != nil {
		t.Fatal("rp_item_magic allocated")
	}
	if w.rpItem(1001) != nil {
		t.Fatal("rp_item_magic allocated the entity_item")
	}
	im := w.pItemMagic(1001)
	if w.rpItem(1001) == nil || w.rpItem(1001).xItemMagic != im {
		t.Fatal("p_item_magic did not allocate the entity_item")
	}

	if w.pLocInfo(1001) != &w.bx[1001].xLocInfo || w.rpLocInfo(1001) != w.pLocInfo(1001) {
		t.Fatal("p_loc_info is not the box's loc_info")
	}
}

// TestWorldsAreIndependent checks that two worlds share no state.
func TestWorldsAreIndependent(t *testing.T) {
	a, b := newWorld(), newWorld()
	a.allocBox(5, tItem, 0)
	a.indepPlayer = 101
	if b.kind(5) != tDeleted || b.boxHead[tItem] != 0 || b.subHead[0] != 0 || b.indepPlayer != 100 {
		t.Fatal("worlds share state")
	}
}

// TestCheckGlob is check_glob() in check.c: kind_s and subkind_s have T_MAX
// and SUB_MAX entries.
func TestCheckGlob(t *testing.T) {
	if len(kindS) != tMax {
		t.Errorf("len(kindS) = %d, want %d", len(kindS), tMax)
	}
	if len(subkindS) != subMax {
		t.Errorf("len(subkindS) = %d, want %d", len(subkindS), subMax)
	}
}
