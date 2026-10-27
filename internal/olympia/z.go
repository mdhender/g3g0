// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"

	"github.com/mdhender/g3g0/internal/rnd"
	"github.com/mdhender/g3g0/internal/rng"
)

// testRandom is test_random() in olympia/z.c, the -R self-test. It draws
// from a local stream rooted at the master seed. -R runs before load_db(),
// so no seed has been loaded and the master seed is all zero.
//
// On a terminal the first 100 draws print as a 10x10 grid, each row drawn in
// order before it prints; otherwise they print one per line.
func testRandom(out *cStdout, tty bool) {
	var seed rnd.Source // never loaded: rng_master_seed() is all zero
	s := rng.Seed(seed.MasterSeed())

	if tty {
		for range 10 {
			var r [10]int
			for j := range r {
				r[j] = s.Draw(1, 10)
			}
			out.Printf("%3d  %3d  %3d  %3d  %3d  %3d  %3d  %3d  %3d  %3d\n",
				r[0], r[1], r[2], r[3], r[4], r[5], r[6], r[7], r[8], r[9])
		}
	} else {
		for range 100 {
			out.Printf("%d\n", s.Draw(1, 10))
		}
	}

	for i := -10; i >= -16; i-- {
		out.Printf("rng_draw(%d, %d) == %d\n", i, -3, s.Draw(i, -3))
	}

	for range 100 {
		out.Printf("%d\n", s.Draw(1000, 9999))
	}

	l := make([]int, 0, 10)
	for i := 1; i <= 10; i++ {
		l = append(l, i)
	}
	rng.Shuffle(l, &s)

	out.Printf("Scramble:\n")
	for _, v := range l {
		out.Printf("%d\n", v)
	}
}

// assertion is a failed assert(): asfail() in olympia/z.c prints it to
// stderr and calls abort(). The port panics with it, and whatever runs the
// engine recovers it, prints Error() to stderr and exits as abort() does
// (status 134, without flushing stdout, which abort() discards).
//
// file stands in for the __FILE__ that C prints, which is whatever path the
// C build compiled the source as (an absolute path in the CMake build);
// the port uses the path in the G3 repository, "olympia/code.c".
type assertion struct {
	file string
	line int
	cond string
}

// Error is the line asfail() prints.
func (a assertion) Error() string {
	return fmt.Sprintf("assertion failure: %s (%d): %s\n", a.file, a.line, a.cond)
}

// assert is assert(cond): if ok is false it panics with an assertion for
// the C source line file:line and the condition's text.
func assert(ok bool, file string, line int, cond string) {
	if !ok {
		panic(assertion{file: file, line: line, cond: cond})
	}
}
