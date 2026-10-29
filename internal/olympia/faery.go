// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Complete faery.c:529–634. Hunts use slot/retry keyed faer leaves and
// the shared mint allocator; NPC routes use the existing exdr stream.
func (w *world) createElvenHunt(slot int) {
	p := w.rpLocInfo(int(w.faeryRegion))
	assert(p != nil, "olympia/faery.c", 538, "p")
	var where int
	for retry := 0; ; retry++ {
		where = p.hereList[w.streams.FaerHuntLoc(slot, retry, 0, len(p.hereList)-1)]
		if w.subkind(where) != subOcean {
			break
		}
	}
	n := w.newChar(subNi, itemElf, where, 100, int(w.faeryPlayer), loyNpc, 0, "Faery Hunt")
	if n < 0 {
		return
	}
	w.genItem(n, itemElf, w.streams.FaerHuntQty(slot))
	w.queue(n, "wait time 0")
	w.initLoadSup(n)
}

func (w *world) warnHuman(who, target int) {
	w.queue(who, "message 1 %s", boxCodeLess(target))
	w.queue(who, "You are not welcome in Faery.  Leave, or you will be killed.")
	w.logWrite(11, "Faery hunt warned %s.", w.boxName(target)) // LOG_SPECIAL
}

func (w *world) autoFaerySup(who int) {
	where := w.subloc(who)
	p := w.pMisc(w.player(who))
	assert(w.validBox(where), "olympia/faery.c", 581, "valid_box(where)")
	queued := false
	for i := range w.loopHere(where) {
		if w.kind(i) != tChar || w.subkind(w.player(i)) != subPlRegular {
			continue
		}
		if w.stackHasUseKey(i, useFaeryStone) != 0 {
			continue
		}
		queued = true
		if !testBit(p.npcMemory, i) {
			w.warnHuman(who, i)
			setBit(&p.npcMemory, i)
			continue
		}
		w.queue(who, "attack %s", boxCodeLess(i))
	}
	if !queued {
		w.npcMove(who)
	}
}

func (w *world) autoFaery() {
	n := 0
	w.phaseRan(phFaerySpawn)
	// C counts every faction unit, regardless of kind or current location.
	for range w.loopUnits(int(w.faeryPlayer)) {
		n++
	}
	for n < 15 {
		w.createElvenHunt(n)
		w.phaseActed(phFaerySpawn)
		n++
	}
	for i := range w.loopUnits(int(w.faeryPlayer)) {
		w.autoFaerySup(i)
		w.phaseActed(phFaeryUnit)
	}
}
