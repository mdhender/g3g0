// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"io"
	"os"

	"github.com/mdhender/g3g0/internal/clib"
)

// Ported from olympia/sout.c.

const vect = -1 // sout.h VECT

// sout is sout(): the formatted string, cut to BUF_LEN (LEN) bytes, the
// size of the ring buffers C prints into. C's ring of 50 buffers is not
// needed: each call returns a new string.
func sout(format string, a ...any) string {
	s := fmt.Sprintf(format, a...)
	return s[:min(len(s), lenMax)]
}

// commandOut is the default character/player route through out/vout,
// out_sup and bottom_out. wout is an alias for out in C. Keep raw output
// framing even in immediate mode, and keep the player's output bit mutation.
// Other routing/formatting and report assembly are stage 2, not no-op stubs.
func (w *world) commandOut(who int, format string, a ...any) {
	if who == vect {
		for _, n := range w.outVector {
			assert(n != vect, "olympia/sout.c", 457, "out_vector[i] != VECT")
			w.commandOut(n, format, a...)
		}
		return
	}
	if w.isPrisoner(who) {
		return
	}
	if w.subkind(who) == subGarrison || (w.kind(who) != tChar && w.kind(who) != tPlayer) {
		panic(notImplemented("command output routing"))
	}
	w.bottomOut(w.player(who), who, format, a...)
}

// vector_char_here uses all_char_here: descendants only, depth first,
// excluding non-character branches. The vector persists until replaced.
func (w *world) vectorCharHere(where int) {
	assert(w.validBox(where), "olympia/loc.c", 92, "valid_box(who)")
	w.outVector.Clear()
	// all_char_here's NULL loc_info check cannot fire after valid_box:
	// oly.h embeds x_loc_info in the box, and rp_loc_info takes its address.
	// An unlocated character still has this structure (where == 0).
	for _, n := range w.rpLocInfo(where).hereList {
		if w.kind(n) == tChar {
			w.addCharHere(n, &w.outVector)
		}
	}
}

// logWrite is the small MASTER route used for code alerts. It bypasses
// ordinary character routing but shares bottom_out framing and output bits.
func (w *world) logWrite(k int, format string, a ...any) {
	assert(k >= 10 && k <= 20, "olympia/sout.c", 708, "k >= 10 && k <= 20")
	if !w.isPrisoner(int(w.gmPlayer)) {
		w.bottomOut(int(w.gmPlayer), k, format, a...)
	}
}

func (w *world) bottomOut(pl, who int, format string, a ...any) {
	var fp io.Writer = w.stdout
	if !w.immediate {
		if pl == 0 {
			return
		}
		name := sout("%s/log/%d", w.libdir, pl)
		f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o666)
		if err != nil {
			fmt.Fprintf(w.stderr, "alloc_fp: can't open %s: ", name)
			clib.Perror(w.stderr, "", err)
			assert(false, "olympia/sout.c", 158, "FALSE")
		}
		defer f.Close()
		fp = f
		setBit(&w.pPlayer(pl).output, who)
	}
	// vout's buffer has LEN bytes, unlike sout's LEN+1.
	s := fmt.Sprintf(format, a...)
	s = s[:min(len(s), lenMax-1)]
	day := ""
	if w.showDay {
		day = fmt.Sprintf("%d", w.sysclock.day)
	}
	fmt.Fprintf(fp, "%d:0::0:%s:%s\n", who, day, s)
}
