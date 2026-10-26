// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package genesis

// This file ports genesis/pcg.c: the PCG32 generator and splitmix64 seed
// mixing that genesis uses instead of the engine's MD5 RNG.

// pcgMult is the classic PCG32 multiplier.
const pcgMult = 6364136223846793005

// pcg32 is pcg32_t.
type pcg32 struct {
	state uint64
	inc   uint64
}

// splitmix64 is the 64-bit avalanche mixer.
func splitmix64(x uint64) uint64 {
	x += 0x9e3779b97f4a7c15
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	x = x ^ (x >> 31)
	return x
}

// next returns the next raw 32-bit draw.
func (r *pcg32) next() uint32 {
	old := r.state
	r.state = old*pcgMult + r.inc
	xorshifted := uint32(((old >> 18) ^ old) >> 27)
	rot := uint32(old >> 59)
	return (xorshifted >> rot) | (xorshifted << (-rot & 31))
}

// init seeds the generator from an explicit (state, sequence) pair.
func (r *pcg32) init(initstate, initseq uint64) {
	r.state = 0
	r.inc = (initseq << 1) | 1
	r.next()
	r.state += initstate
	r.next()
}

// seedStep seeds the generator for one generation step, so each step's
// stream depends only on the user seed.
func (r *pcg32) seedStep(seed, step, island uint64) {
	k := splitmix64(seed + step*0x9e3779b97f4a7c15 + island*0xd1b54a32d192ed03)
	s0 := splitmix64(k)
	s1 := splitmix64(s0)
	r.init(s0, s1)
}

// rangeInt is pcg32_range: an unbiased draw in [lo, hi] by bitmask and
// rejection. The arithmetic follows the C int and uint32_t conversions.
func (r *pcg32) rangeInt(lo, hi int) int {
	span := uint32(int32(hi - lo))
	var mask uint32
	for x := span; x != 0; x >>= 1 {
		mask |= x
	}
	var n uint32
	for {
		n = r.next() & mask
		if n <= span {
			break
		}
	}
	return int(int32(n + uint32(int32(lo))))
}
