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
	"io"
	"os"
)

// Source is the mutable generator state.
type Source struct {
	digest [16]byte

	// master is the state as Load left it, before any draw advanced it:
	// master_seed in lib/rnd.c and olympia/rnd.c. It roots the addressable
	// streams in internal/rng.
	master [16]byte
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

// Load is load_seed() without its message: it replaces the start of the
// state with up to 16 bytes read from path. A read error after a successful
// open (for example, a directory) reads nothing, as fread() does.
//
// Either way, like load_seed(), it then captures the state as the master
// seed. When path cannot be opened, the master seed is the current state,
// including any draws made since the last load.
func (s *Source) Load(path string) error {
	defer func() { s.master = s.digest }()
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	var buf [16]byte
	n, _ := io.ReadFull(fp, buf[:])
	s.SetSeed(buf[:n])
	return nil
}

// MasterSeed is rng_master_seed(): the state captured by the last Load, as
// the four native 32-bit words C copies out of it (little-endian here, as
// for the draws). It is zero until the first Load.
func (s *Source) MasterSeed() [4]uint32 {
	var out [4]uint32
	for i := range out {
		out[i] = binary.LittleEndian.Uint32(s.master[4*i:])
	}
	return out
}

// Save is save_seed() without its message: it writes the state to path,
// opened the way fopen(path, "wb+") opens it.
func (s *Source) Save(path string) error {
	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}
	_, _ = fp.Write(s.digest[:])
	_ = fp.Close()
	return nil
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

// MD5Int is md5_int() in lib/rnd.c: the low word of the MD5 of the four
// arguments, each stored as a 32-bit word, read back as a C int.
func MD5Int(a, b, c, d int) int {
	var buf [16]byte
	for i, v := range [4]int{a, b, c, d} {
		binary.LittleEndian.PutUint32(buf[4*i:], uint32(v))
	}
	sum := md5.Sum(buf[:])
	return int(int32(binary.LittleEndian.Uint32(sum[:4])))
}
