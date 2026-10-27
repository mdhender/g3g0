// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package rng_test

import (
	"testing"

	"github.com/mdhender/g3g0/internal/rng"
)

// Purpose tags from G3's tests/rng/rng_check.c.
const (
	tagCombat = 0x636f6d62 // "comb"
	tagMarket = 0x6d726b74 // "mrkt"
	tagHit    = 0x68697420 // "hit "
	tagDmg    = 0x646d6720 // "dmg "
	tagMorale = 0x6d726c20 // "mrl "
)

// TestRngCheck ports the assertions in G3's tests/rng/rng_check.c.
func TestRngCheck(t *testing.T) {
	root := rng.Seed([4]uint32{0x01234567, 0x89abcdef, 0xfedcba98, 0x76543210})

	// 1. Determinism: identical inputs give identical output.
	if a, b := root.Keyed(501, 17, tagHit, 1, 100), root.Keyed(501, 17, tagHit, 1, 100); a != b {
		t.Errorf("determinism: same key gave %d and %d", a, b)
	}

	// 2. Keyed independence: inserting a morale draw between hit and damage
	// moves neither.
	battle := root.Of(501, tagCombat)
	hit1 := battle.Keyed(1, 0, tagHit, 1, 20)
	dmg1 := battle.Keyed(1, 0, tagDmg, 1, 6)
	hit2 := battle.Keyed(1, 0, tagHit, 1, 20)
	_ = battle.Keyed(1, 0, tagMorale, 1, 10)
	dmg2 := battle.Keyed(1, 0, tagDmg, 1, 6)
	if hit1 != hit2 {
		t.Errorf("keyed independence: hit moved from %d to %d", hit1, hit2)
	}
	if dmg1 != dmg2 {
		t.Errorf("keyed independence: damage moved from %d to %d", dmg1, dmg2)
	}

	// 3. Order-independent derivation: deriving combat before or after
	// market gives the same seeds.
	combatA := root.Of(0, tagCombat)
	marketA := root.Of(0, tagMarket)
	marketB := root.Of(0, tagMarket)
	combatB := root.Of(0, tagCombat)
	if combatA.Seed != combatB.Seed {
		t.Error("order-independence: combat stream moved")
	}
	if marketA.Seed != marketB.Seed {
		t.Error("order-independence: market stream moved")
	}
	if combatA.Seed == marketA.Seed {
		t.Error("distinct tags gave the same stream")
	}

	// The contrast rng_check.c prints for information: sequential draws are
	// sequence-dependent. Here it is asserted with the C build's values.
	s1 := root.Of(7, tagCombat)
	s2 := s1
	a1, a2 := s1.Draw(1, 1000000), s1.Draw(1, 1000000)
	_ = s2.Draw(1, 1000000)
	b1, b2 := s2.Draw(1, 1000000), s2.Draw(1, 1000000)
	if a1 != 86357 || a2 != 488217 || b1 != 488217 || b2 != 799350 {
		t.Errorf("sequential draws = %d,%d vs %d,%d, want 86357,488217 vs 488217,799350", a1, a2, b1, b2)
	}
}

// TestKeyedLeavesCounter checks that keyed draws and derivation leave the
// counter alone, while sequential draws advance it.
func TestKeyedLeavesCounter(t *testing.T) {
	s := rng.Seed([4]uint32{1, 2, 3, 4})
	_ = s.Keyed(1, 2, 3, 0, 9)
	_ = s.Of(1, 2)
	if s.Counter != 0 {
		t.Fatalf("counter = %d after keyed draw and derivation, want 0", s.Counter)
	}
	_ = s.Draw(0, 9)
	if s.Counter != 1 {
		t.Fatalf("counter = %d after one draw, want 1", s.Counter)
	}
	if c := s.Of(1, 2); c.Counter != 0 {
		t.Fatalf("child counter = %d, want 0", c.Counter)
	}
}
