// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package rng ports G3's addressable RNG layer, lib/rng.c, and the named
// streams the engine derives from it (the begin_*() helpers listed in G3's
// doc/claude/rng-streams.md).
//
// A stream is a 128-bit seed plus a counter. Children are derived by hashing
// the parent seed with a key and a purpose tag, so a stream depends only on
// its address, not on how many draws or siblings came before it. Every hash
// is MD5 over native 32-bit words, read and written little-endian, which is
// what the C code does on the little-endian hosts G3 supports. The range
// reduction is the legacy rnd() mask-and-reject loop, so draws reuse
// internal/rnd.
package rng

import (
	"crypto/md5"
	"encoding/binary"

	"github.com/mdhender/g3g0/internal/rnd"
)

// Stream is rng_stream. Copy it by value; there is no global state.
type Stream struct {
	Seed    [4]uint32
	Counter uint32
}

// Seed is rng_seed(): the root stream for a master seed.
func Seed(master [4]uint32) Stream {
	return Stream{Seed: master}
}

// Of is rng_stream_of(): the child of s for key and tag. The key is stored
// as a 32-bit word, so it wraps the way C's (uint32_t) cast of an int does.
func (s *Stream) Of(key int, tag uint32) Stream {
	return Stream{Seed: words(md5Words(s.Seed[0], s.Seed[1], s.Seed[2], s.Seed[3], uint32(key), tag))}
}

// Draw is rng_draw(): a sequential draw in [low, high] that advances the
// counter. As in C, high < low never terminates.
func (s *Stream) Draw(low, high int) int {
	h := md5Words(s.Seed[0], s.Seed[1], s.Seed[2], s.Seed[3], s.Counter)
	s.Counter++
	return reduceRange(h, low, high)
}

// Keyed is rng_keyed(): a draw in [low, high] that depends only on the
// stream seed, k1, k2 and tag. It does not touch the counter. As in C,
// high < low never terminates.
func (s *Stream) Keyed(k1, k2 int, tag uint32, low, high int) int {
	h := md5Words(s.Seed[0], s.Seed[1], s.Seed[2], s.Seed[3], uint32(k1), uint32(k2), tag)
	return reduceRange(h, low, high)
}

// Shuffle is ilist_shuffle_rng() in lib/ilist.c, and exit_views_shuffle_rng()
// in lib/exit_views.c, which is the same loop over a different element type:
// for each index but the last, draw a partner from the rest of the list and
// swap. Lists of fewer than two elements draw nothing.
func Shuffle[T any](l []T, s *Stream) {
	n := len(l) - 1
	for i := 0; i < n; i++ {
		if r := s.Draw(i, n); r != i {
			l[i], l[r] = l[r], l[i]
		}
	}
}

// reduceRange is reduce_range() in lib/rng.c. It is rnd() run on a private
// copy of the digest h, so it reuses internal/rnd.
func reduceRange(h [16]byte, low, high int) int {
	return rnd.New(h).Rnd(low, high)
}

// md5Words is C's MD5(out, in, sizeof in) over an array of uint32.
func md5Words(in ...uint32) [16]byte {
	buf := make([]byte, 4*len(in))
	for i, w := range in {
		binary.LittleEndian.PutUint32(buf[4*i:], w)
	}
	return md5.Sum(buf)
}

// words reads a digest as the four uint32 words C stores it in.
func words(h [16]byte) [4]uint32 {
	var w [4]uint32
	for i := range w {
		w[i] = binary.LittleEndian.Uint32(h[4*i:])
	}
	return w
}
