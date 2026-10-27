// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package rnd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mdhender/g3g0/internal/rnd"
)

// The expected values were produced by lib/rnd.c at G3 revision
// 03b5cea880c44fb42b40173dacf20d32896c0aa4, seeded with the mapgen fixture
// randseed, calling rnd(low, high) for each row in order.
func TestRndMatchesC(t *testing.T) {
	var seed [16]byte
	copy(seed[:], "slartibartfast?\n")
	s := rnd.New(seed)

	draws := []struct{ low, high, want int }{
		{1, 10, 8},
		{0, 1, 0},
		{111, 333, 302},
		{56760, 58759, 58701},
		{59000, 78999, 73008},
		{5, 5, 5},
		{-16, -3, -9},
		{0, 2147483647, 1640603723},
		{7, 7, 7},
		{1, 10, 9},
		{0, 1, 1},
		{111, 333, 145},
		{56760, 58759, 57298},
		{59000, 78999, 77212},
		{5, 5, 5},
		{-16, -3, -3},
		{0, 2147483647, 888269455},
		{7, 7, 7},
		{1, 10, 3},
		{0, 1, 1},
		{111, 333, 190},
		{56760, 58759, 58044},
		{59000, 78999, 73121},
		{5, 5, 5},
		{-16, -3, -16},
		{0, 2147483647, 1913720182},
		{7, 7, 7},
	}
	for i, d := range draws {
		if got := s.Rnd(d.low, d.high); got != d.want {
			t.Fatalf("draw %d: Rnd(%d, %d) = %d, want %d", i, d.low, d.high, got, d.want)
		}
	}

	want := [16]byte{
		0x4d, 0x8f, 0x0f, 0x13, 0x18, 0x0b, 0x27, 0x3e,
		0xc4, 0x1e, 0x2e, 0x38, 0x58, 0xc2, 0x5c, 0x3d,
	}
	if got := s.Seed(); got != want {
		t.Errorf("final seed = % x, want % x", got, want)
	}
}

func TestLoadSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "randseed")

	// A short file replaces only the start of the state.
	if err := os.WriteFile(path, []byte{1, 2, 3}, 0o644); err != nil {
		t.Fatal(err)
	}
	s := rnd.New([16]byte{})
	if err := s.Load(path); err != nil {
		t.Fatal(err)
	}
	if got, want := s.Seed(), [16]byte{1, 2, 3}; got != want {
		t.Errorf("Load short seed = %x, want %x", got, want)
	}
	if err := s.Save(path); err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(path); err != nil || len(got) != 16 || got[0] != 1 {
		t.Errorf("Save wrote %x, %v", got, err)
	}

	// A missing file is an error and leaves the state alone; a directory
	// opens, reads nothing, and cannot be saved.
	if err := s.Load(filepath.Join(dir, "missing")); err == nil {
		t.Error("Load of a missing file succeeded")
	}
	if err := s.Load(dir); err != nil {
		t.Errorf("Load of a directory: %v", err)
	}
	if got, want := s.Seed(), [16]byte{1, 2, 3}; got != want {
		t.Errorf("state after failed loads = %x, want %x", got, want)
	}
	if err := s.Save(dir); err == nil {
		t.Error("Save to a directory succeeded")
	}
}
