// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Complete hades.c:404–526. World creation and the autonomous NPC dispatcher
// are separate work; these routines preserve the real faction population count
// (not a count of monsters, leaders, living units or units currently in Hades).
func (w *world) createHadesNasty(slot int) {
	p := w.rpLocInfo(int(w.hadesRegion))
	assert(p != nil, "olympia/hades.c", 412, "p")
	where := p.hereList[w.streams.HadsNasty(slot, 0, 0, len(p.hereList)-1)]
	var n int
	switch w.streams.HadsNasty(slot, 1, 1, 4) {
	case 1:
		n = w.newChar(subNi, itemSpirit, where, 100, int(w.hadesPlayer), loyNpc, 0, "Tortured spirits")
		if n < 0 {
			return
		}
		w.genItem(n, itemSpirit, w.streams.HadsNasty(slot, 2, 25, 75))
	case 2:
		n = w.newChar(0, 0, where, 100, int(w.hadesPlayer), loyNpc, 0, "Ghostly presence")
		if n < 0 {
			return
		}
		w.pChar(n).attack = 100
		w.rpChar(n).defense = 100
	case 3:
		n = w.newChar(0, 0, where, 100, int(w.hadesPlayer), loyNpc, 0, "Lesser Demon")
		if n < 0 {
			return
		}
		w.pChar(n).attack = 250
		w.rpChar(n).defense = 250
		w.genItem(n, itemSpirit, w.streams.HadsNasty(slot, 2, 50, 150))
	case 4:
		n = w.newChar(0, 0, where, 100, int(w.hadesPlayer), loyNpc, 0, "Greater Demon")
		if n < 0 {
			return
		}
		w.pChar(n).attack = 500
		w.rpChar(n).defense = 500
		w.genItem(n, itemSpirit, w.streams.HadsNasty(slot, 2, 100, 250))
	default:
		assert(false, "olympia/hades.c", 465, "FALSE")
	}
	w.queue(n, "wait time 0")
	w.initLoadSup(n)
}

func (w *world) autoHadesSup(who int) {
	where := w.subloc(who)
	queued := false
	for i := range w.loopHere(where) {
		if w.kind(i) != tChar || w.subkind(w.player(i)) != subPlRegular {
			continue
		}
		if w.hasSkill(i, skTranscendDeath) != 0 && w.stackParent(i) == 0 && w.countStackAny(i) == 1 {
			continue
		}
		queued = true
		w.queue(who, "attack %s", boxCodeLess(i))
	}
	if !queued {
		w.npcMove(who)
	}
}

func (w *world) autoHades() {
	n := 0
	for range w.loopUnits(int(w.hadesPlayer)) {
		n++
	}
	for n < 25 {
		w.createHadesNasty(n)
		n++
	}
	for i := range w.loopUnits(int(w.hadesPlayer)) {
		w.autoHadesSup(i)
	}
}

// randomHadesLoc is hades.c:529–555, using the persistent keyed Hades stream.
func (w *world) randomHadesLoc(who int) int {
	var candidates []int
	for i := range w.loopKind(tLoc) {
		if w.region(i) == int(w.hadesRegion) && w.subkind(i) == subUnder {
			candidates = append(candidates, i)
		}
	}
	if len(candidates) == 0 {
		return 0
	}
	return candidates[w.streams.HadsRandomLoc(who, len(candidates))]
}
