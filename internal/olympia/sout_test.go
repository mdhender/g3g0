// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"strings"
	"testing"
)

type outputFormatProbe struct{ calls int }

func (p *outputFormatProbe) String() string {
	p.calls++
	return "formatted"
}

// Go's variadic boundary must have the same single formatting operation as
// C vout's buffer, even when recipients duplicate, skip, or do not exist.
// The exact resulting bytes/styles/routes are covered by the real-C oracle.
func TestOutputFormatsOnceBeforeRouting(t *testing.T) {
	for _, recipients := range [][]int{nil, {200, 200, 201}, {1000}} {
		w := newWorld()
		w.immediate = true
		w.outPath = master
		var out strings.Builder
		w.stdout = &out
		w.allocBox(200, tPlayer, subPlRegular)
		w.allocBox(201, tPlayer, subPlRegular)
		w.allocBox(1000, tChar, 0)
		w.pChar(1000).prisoner = 1
		for _, n := range recipients {
			w.vectorAdd(n)
		}
		p := &outputFormatProbe{}
		w.commandOut(vect, "%s", p)
		if p.calls != 1 {
			t.Fatalf("recipients %v: formatted %d times, want once", recipients, p.calls)
		}
	}
}
