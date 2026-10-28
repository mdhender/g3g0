// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

// Complete savage.c:13–65,219–343,351–403 at the stage pin. The drum
// command handlers and autonomous NPC dispatcher are separate work.
func (w *world) createSavage(where int) int {
	n := w.newChar(subNi, itemSavage, where, 100, int(w.indepPlayer), loyNpc, 0, "")
	if n < 0 {
		return -1
	}
	w.streams.BeginNPC(where)
	w.genItem(n, itemDrum, 1)
	w.genItem(n, itemSavage, w.streams.NPCQty(where, n, 3, 25))
	return n
}

func (w *world) callSavage(where, toWhere, who, why int) bool {
	if w.controlledHumansHere(where) {
		return false
	}
	n := w.createSavage(where)
	w.queue(n, "move %s", boxCodeLess(toWhere))
	switch why {
	case 0:
		w.queue(n, "attack %s", boxCodeLess(who))
	case 1:
		w.setLoyal(n, loySummon, 3)
		w.queue(n, "stack %s", boxCodeLess(who))
	case 2:
		w.queue(n, "use 98 1")
		w.queue(n, "wait time %d", w.streams.NPCBehavior(n, 0, 35, 50))
		w.queue(n, "attack %s", boxCodeLess(who))
	}
	w.initLoadSup(n)
	return true
}

func (w *world) savageHates(where int) bool {
	switch w.subkind(where) {
	case subInn, subInnNotdone, subCastle, subCastleNotdone,
		subTower, subTowerNotdone, subGalley, subGalleyNotdone,
		subRoundship, subRoundshipNotdone, subTemple, subTempleNotdone:
		return true
	}
	return false
}

func (w *world) savageHateHere(where int) int {
	assert(w.validBox(where), "olympia/savage.c", 250, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.isLocOrShip(i) && w.locDepth(i) == locBuild && w.savageHates(i) {
			return i
		}
	}
	return 0
}

func (w *world) autoSavage(who int) {
	where := w.subloc(who)
	if w.stackParent(who) != 0 {
		return
	}
	if w.savageHates(where) {
		w.queue(who, "raze")
		return
	}
	if target := w.savageHateHere(where); target != 0 {
		if w.buildingOwner(target) == 0 {
			w.queue(who, "enter %s", boxCodeLess(target))
			w.queue(who, "raze")
			return
		}
		if w.controlledHumansHere(target) {
			if w.hasItem(who, itemDrum) < 1 {
				w.genItem(who, itemDrum, 1)
			}
			w.streams.BeginNPC(where)
			w.queue(who, "use %d 1", itemDrum)
			w.queue(who, "wait time %d", w.streams.NPCBehavior(who, 0, 35, 50))
			w.queue(who, "attack %s", boxCodeLess(target))
			return
		}
	}
	assert(w.validBox(who), "olympia/savage.c", 329, "valid_box(who)")
	for i := range w.loopHere(who) {
		if w.kind(i) == tChar {
			w.queue(who, "unstack %d", i)
		}
	}
	if !w.controlledHumansHere(where) {
		w.queue(who, "die")
		return
	}
	w.npcMove(who)
}

func (w *world) initSavageAttacks() {
	// Reproduce G3 #354: the initial tally is never incremented after spawning,
	// so multiple eligible forts can exceed the 200-unit limit within a turn.
	// https://github.com/playbymail/olympia-g3/issues/354
	w.numSavages = 0
	for i := range w.loopKind(tChar) {
		if w.nobleItem(i) == itemSavage {
			w.numSavages++
		}
	}
	if w.numSavages >= 200 {
		return
	}
	for fort := range w.loopKind(tLoc) {
		if w.locDepth(fort) != locBuild {
			continue
		}
		w.streams.BeginNPC(fort)
		if w.streams.NPCSpawn(fort, 0, 1, 100) != 1 {
			continue
		}
		where := w.subloc(fort)
		if w.locDepth(where) != locProvince {
			continue
		}
		l := w.exitsFromLocNSEWSelect(0, where, exitLand, 1)
		if len(l) == 0 {
			fmt.Fprintln(w.stderr, "init_savage_attacks: no exits?")
			continue
		}
		for _, e := range l {
			if w.callSavage(int(e.destination), where, fort, 2) {
				break
			}
		}
	}
}
