// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/mdhender/g3g0/internal/clib"
)

// Ported from olympia/sout.c.

const vect = -1 // sout.h VECT

func (w *world) openLogfileNondestruct() {
	_ = os.Mkdir(sout("%s/log", w.libdir), 0o755)
}

func (w *world) openLogfile(win bool) {
	if w.immediate {
		return
	}
	name := sout("%s/log", w.libdir)
	if win {
		name = sout("%s\\log", w.libdir)
	}
	removeDir(name)
	_ = os.Mkdir(name, 0o755)
}

// sout is sout(): the formatted string, cut to BUF_LEN (LEN) bytes, the
// size of the ring buffers C prints into. C's ring of 50 buffers is not
// needed: each call returns a new string.
func sout(format string, a ...any) string {
	s := fmt.Sprintf(format, a...)
	return s[:min(len(s), lenMax)]
}

// wiout is sout.c:506–519. It resets secondary indentation to zero, even
// when the caller had a nonzero secondary indent before this call.
func (w *world) wiout(who int, ind int32, format string, a ...any) {
	s := sout(format, a...)
	w.secondIndent = ind
	w.commandOut(who, "%s", s)
	w.secondIndent = 0
}

// commandOut is the default character/player route through out/vout,
// out_sup and bottom_out. wout is an alias for out in C. Keep raw output
// framing even in immediate mode, and keep the player's output bit mutation.
// Report assembly remains stage 2; location and garrison routing is real.
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
	if w.subkind(who) == subGarrison {
		w.outGarrison(who, format, a...)
	} else if w.isLocOrShip(who) {
		w.outLocation(who, format, a...)
	} else {
		w.bottomOut(w.player(who), who, format, a...)
	}
}

func (w *world) canViewLoc(pl int, p *entityPlayer, where, outer int) bool {
	var rulers []int
	if w.showToGarrison && w.locDepth(outer) == locProvince {
		rulers = w.playersWhoRuleHere(outer)
	}
	if testBit(p.locs, where) {
		return true
	}
	if testBit(p.locs, outer) {
		if w.locHidden(where) != 0 && !w.testKnown(pl, where) {
			return false
		}
		return true
	}
	if w.showToGarrison {
		for _, n := range rulers {
			if n == pl {
				return where == outer || w.locHidden(where) == 0
			}
		}
	}
	return false
}

func (w *world) outLocation(where int, format string, a ...any) {
	outer := w.viewloc(where)
	for pl := range w.loopKind(tPlayer) {
		if w.canViewLoc(pl, w.pPlayer(pl), where, outer) {
			w.bottomOut(pl, outer, format, a...)
		}
	}
}

func (w *world) outGarrison(garr int, format string, a ...any) {
	for _, pl := range w.playersWhoRuleHere(garr) {
		if pl != 0 && w.subkind(pl) != subPlSilent {
			w.bottomOutUnit(pl, 6, garr, format, a...) // sout.h OUT_GARR
		}
	}
}

// vectorStack is sout.c:585–602. Unlike vectorCharHere it includes the root;
// clear=false preserves prior recipients and appends only missing entries.
func (w *world) vectorStack(who int, clear bool) {
	if clear {
		w.outVector.Clear()
	}
	if clear || w.outVector.Lookup(who) < 0 {
		w.outVector.Append(who)
	}
	for _, i := range w.allCharHere(who) {
		if clear || w.outVector.Lookup(i) < 0 {
			w.outVector.Append(i)
		}
	}
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
	w.bottomOutUnit(pl, who, 0, format, a...)
}

type playerOutput struct {
	player int
	fp     *clib.File
}

// sout.c MAX_FP = FOPEN_MAX - 10 (stdio.h: Darwin 20, Linux 16).
func maxPlayerOutput() int {
	if runtime.GOOS == "darwin" {
		return 10
	}
	return 6
}

// openPlayerOutput owns sout.c's bounded move-to-front cache (131–214).
// Both ordinary output and report projections borrow, never close, its FILE.
// A hit must not reopen a renamed/unlinked path or change buffering mode.
func (w *world) openPlayerOutput(pl int) *clib.File {
	for i, p := range w.playerOutput {
		if p.player == pl {
			copy(w.playerOutput[1:i+1], w.playerOutput[:i])
			w.playerOutput[0] = p
			return p.fp
		}
	}
	if len(w.playerOutput) == maxPlayerOutput() {
		w.playerOutput[len(w.playerOutput)-1].fp.Close()
		w.playerOutput = w.playerOutput[:len(w.playerOutput)-1]
	}
	name := sout("%s/log/%d", w.libdir, pl)
	f, err := clib.FopenAppend(name)
	if err != nil {
		fmt.Fprintf(w.stderr, "alloc_fp: can't open %s: ", name)
		clib.Perror(w.stderr, "", err)
		assert(false, "olympia/sout.c", 158, "FALSE")
	}
	if w.flushAlways {
		f.SetUnbuffered()
	}
	w.playerOutput = append(w.playerOutput, playerOutput{})
	copy(w.playerOutput[1:], w.playerOutput[:len(w.playerOutput)-1])
	w.playerOutput[0] = playerOutput{pl, f}
	return f
}

// closeLogfile is sout.c:241–264. Immediate mode leaves the cache intact.
func (w *world) closeLogfile() {
	if w.immediate {
		return
	}
	for _, p := range w.playerOutput {
		p.fp.Close()
	}
	w.playerOutput = nil
}

func (w *world) bottomOutUnit(pl, who, unit int, format string, a ...any) {
	var fp io.Writer = w.stdout
	if !w.immediate {
		if pl == 0 {
			return
		}
		f := w.openPlayerOutput(pl)
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
	ind := fmt.Sprint(w.indent)
	if w.secondIndent != 0 {
		ind += fmt.Sprintf("/%d", w.secondIndent)
	}
	code := ""
	if unit != 0 {
		code = boxCodeLess(unit)
	}
	fmt.Fprintf(fp, "%d:%d:%s:%s:%s:%s\n", who, w.styleHTML, code, ind, day, s)
}
