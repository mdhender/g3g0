// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package rnd ports the MD5-based random number generator in G3's lib/rnd.c.
//
// The generator state is a 16-byte digest. Every draw replaces the digest with
// the MD5 of itself and takes the low 32 bits of the new digest, read as a
// little-endian word, as the candidate value. This matches the C code on the
// little-endian hosts G3 supports (x86-64 and arm64), where the digest is
// stored as four native uint32 words.
package rnd

import (
	"crypto/md5"
	"encoding/binary"
)

// Source is the mutable generator state.
type Source struct {
	digest [16]byte
}

// New returns a source whose state is seed.
func New(seed [16]byte) *Source {
	return &Source{digest: seed}
}

// Seed returns the current state; it is what save_seed() writes.
func (s *Source) Seed() [16]byte {
	return s.digest
}

// SetSeed replaces the state. Like load_seed(), a short seed leaves the
// remaining bytes of the previous state in place, so callers that want the C
// behavior for a truncated randseed file should start from a zero Source.
func (s *Source) SetSeed(b []byte) {
	copy(s.digest[:], b)
}

// Rnd returns a value in [low, high], matching rnd() in lib/rnd.c, including
// its 32-bit integer semantics. As in C, high < low never terminates: the
// mask loop shifts a negative int forever.
func (s *Source) Rnd(low, high int) int {
	rng := uint32(int32(high) - int32(low))
	var mask uint32
	for r := int32(rng); r != 0; r >>= 1 {
		mask |= uint32(r)
	}

	var num uint32
	for {
		s.digest = md5.Sum(s.digest[:])
		num = binary.LittleEndian.Uint32(s.digest[:4]) & mask
		if num <= rng {
			break
		}
	}

	return int(int32(num + uint32(int32(low))))
}
