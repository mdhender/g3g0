// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"time"
)

// Ported from olympia/u.c.

// testBit is test_bit(): whether i is in the sparse set kr.
func testBit(kr sparse, i int) bool {
	return kr.Lookup(i) != -1
}

// setBit is set_bit(): it adds i to the sparse set kr, at the end, unless
// it is there already.
func setBit(kr *sparse, i int) {
	if kr.Lookup(i) == -1 {
		kr.Append(i)
	}
}

// lookup is lookup(): the index of the first entry of table that is s,
// without regard to case, or -1. C's tables end with a NULL entry; the
// slices here have no terminator (see glob.go).
func lookup(table []string, s string) int {
	for i, t := range table {
		if iStrcmp(s, t) == 0 {
			return i
		}
	}

	return -1
}

// stage is stage(s), with s not NULL: it prints s and a newline to
// stderr, or, if time_self (-T) is set, the seconds since the last stage()
// on a line of their own and then s without a newline.
func (w *world) stage(s string) {
	w.stageSup(&s)
}

// stageSup is stage(), where s may be NULL: stage(NULL) prints nothing
// without time_self, and the seconds since the first stage() with it.
func (w *world) stageSup(s *string) {
	if !w.timeSelf {
		if s != nil {
			fmt.Fprintf(w.stderr, "%s\n", *s)
		}
		return
	}

	t := time.Now().Unix()

	if w.stageOld != 0 {
		fmt.Fprintf(w.stderr, "\t%d sec\n", int32(t-w.stageOld))
	} else {
		w.stageFirst = t
	}
	w.stageOld = t

	if s != nil {
		fmt.Fprintf(w.stderr, "%s", *s)
	} else {
		fmt.Fprintf(w.stderr, "%d seconds\n", int32(t-w.stageFirst))
	}
}
