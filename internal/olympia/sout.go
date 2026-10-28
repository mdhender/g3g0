// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/mdhender/g3g0/internal/clib"
	"github.com/mdhender/g3g0/internal/lists"
)

// Ported from olympia/sout.c.

const (
	vect   = -1 // sout.h VECT
	master = -2
)

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
	s := outputText(format, a...)
	w.secondIndent = ind
	w.outSup(who, s)
	w.secondIndent = 0
}

// vout formats exactly once, before routing (including an empty vector).
func outputText(format string, a ...any) string {
	s := fmt.Sprintf(format, a...)
	s = s[:min(len(s), lenMax-1)]
	if i := strings.IndexByte(s, 0); i >= 0 {
		s = s[:i]
	}
	return s
}

// commandOut is out/vout; wout is only an alias, not a wrapping engine.
func (w *world) commandOut(who int, format string, a ...any) {
	w.outSup(who, outputText(format, a...))
}

func (w *world) outSup(who int, s string) {
	if who == vect {
		for _, n := range w.outVector {
			assert(n != vect, "olympia/sout.c", 457, "out_vector[i] != VECT")
			w.outSup(n, s)
		}
		return
	}
	if w.isPrisoner(who) {
		return
	}
	if w.outPath == master {
		w.bottomOut(who, w.outAltWho, s)
	} else if w.subkind(who) == subGarrison {
		w.outGarrison(who, s)
	} else if w.isLocOrShip(who) {
		w.outLocation(who, s)
	} else {
		w.bottomOut(w.player(who), who, s)
	}
}

func (w *world) style(n int32) {
	if n == stylePrev { // One previous state, not a stack or swap.
		w.styleHTML = w.prevStyle
	} else {
		w.prevStyle, w.styleHTML = w.styleHTML, n
	}
}

func (w *world) html(who int, format string, a ...any) {
	s := outputText(format, a...)
	w.style(styleHTML)
	w.outSup(who, s)
	w.style(stylePrev)
}

func (w *world) saveOutputVector() lists.IList {
	t := w.outVector
	w.outVector = nil
	return t
}

func (w *world) restoreOutputVector(t lists.IList) { w.outVector = t }
func (w *world) vectorClear()                      { w.outVector.Clear() }
func (w *world) vectorAdd(who int)                 { w.outVector.Append(who) }

func (w *world) vectorPlayers() {
	w.vectorClear()
	for pl := range w.loopKind(tPlayer) {
		if pl != int(w.eatPl) && pl != int(w.skillPlayer) {
			w.vectorAdd(pl)
		}
	}
}

func (w *world) lines(who int, s string) {
	w.style(styleHTML)
	before := sout("%s", s)
	after := ""
	if i := strings.IndexByte(before, ','); i >= 0 {
		after, before = sout("%s", before[i:]), before[:i]
	}
	w.commandOut(who, "")
	w.commandOut(who, "<name=\"%s\"><font size=+1><b>%s</b>%s</font>", before, before, after)
	w.commandOut(who, "<hr>")
	w.style(styleText)
	w.commandOut(who, "%s", s)
	w.commandOut(who, "%s", strings.Repeat("-", 72))
	w.style(0)
}

func (w *world) matchLines(who int, s string) {
	w.commandOut(who, "%s", s)
	w.commandOut(who, "%s", strings.Repeat("-", len(outputText("%s", s))))
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

func (w *world) outLocation(where int, s string) {
	outer := w.viewloc(where)
	for pl := range w.loopKind(tPlayer) {
		if w.canViewLoc(pl, w.pPlayer(pl), where, outer) {
			w.bottomOut(pl, outer, s)
		}
	}
}

func (w *world) outGarrison(garr int, s string) {
	for _, pl := range w.playersWhoRuleHere(garr) {
		if pl != 0 && w.subkind(pl) != subPlSilent {
			w.bottomOutUnit(pl, 6, garr, s) // sout.h OUT_GARR
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

// logWrite borrows MASTER routing and restores both ordinary route fields.
func (w *world) logWrite(k int, format string, a ...any) {
	path, alt := w.outPath, w.outAltWho
	assert(k >= 10 && k <= 20, "olympia/sout.c", 708, "k >= 10 && k <= 20")
	w.outPath, w.outAltWho = master, k
	w.commandOut(int(w.gmPlayer), format, a...)
	w.outPath, w.outAltWho = path, alt
}

func (w *world) bottomOut(pl, who int, s string) {
	w.bottomOutUnit(pl, who, 0, s)
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
// Ordinary output and report renderers borrow, never close, its FILE.
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

func (w *world) bottomOutUnit(pl, who, unit int, s string) {
	var fp io.Writer = w.stdout
	if !w.immediate {
		if pl == 0 {
			return
		}
		f := w.openPlayerOutput(pl)
		fp = f
		setBit(&w.pPlayer(pl).output, who)
	}
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
