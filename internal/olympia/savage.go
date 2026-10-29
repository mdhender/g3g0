// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

// Complete savage.c:13–146,221–345,353–409 at the stage pin. The summon
// command remains separate work.
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
	w.numSavages++
	return true
}

// vUseDrum is complete savage.c:70–146. Raw speeds other than 0/1
// still emit sound and shuffle exits, but cannot call savages.
func (w *world) vUseDrum(c *command) int {
	where, speed := w.subloc(int(c.who)), int(c.a)
	speedS := ""
	switch speed {
	case 1:
		speedS = "slow "
	case 2:
		speedS = "fast "
	}
	w.commandOut(int(c.who), "%s sounds a %sdrumbeat.", w.boxName(int(c.who)), speedS)
	w.commandOut(where, "%s sounds a %sdrumbeat.", w.boxName(int(c.who)), speedS)
	if w.locDepth(where) != locProvince {
		s := sout("%sbeating drums may be heard coming from %s.", speedS, w.boxName(where))
		w.commandOut(w.subloc(where), "%s", capitalize(s))
		return 1
	}
	l := w.exitsFromLocNSEWSelect(int(c.who), w.province(where), exitLand, 1)
	n := false
	for _, v := range l {
		dir := exitOpposite[v.direction]
		s := sout("%sbeating drums may be heard to the %s.", speedS, fullDirS[dir])
		w.commandOut(int(v.destination), "%s", capitalize(s))
		if w.numSavages < 200 && w.subkind(int(v.destination)) != subOcean &&
			(speed == 0 || speed == 1) && !n {
			n = w.callSavage(int(v.destination), where, int(c.who), speed)
		}
	}
	s := "call"
	switch speed {
	case 0:
		s = "battle challenge"
	case 1:
		s = "call to arms"
	}
	if !n && (speed == 0 || speed == 1) {
		w.commandOut(int(c.who), "No savages are responding to the %s.", s)
	} else if n && (speed == 0 || speed == 1) {
		w.commandOut(int(c.who), "Savages will surely respond to the %s.", s)
	}
	return 1
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
	assert(w.validBox(where), "olympia/savage.c", 252, "valid_box(where)")
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
	assert(w.validBox(who), "olympia/savage.c", 331, "valid_box(who)")
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
	w.phaseRan(phSavageSpawn)
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
		if w.numSavages >= 200 {
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
				w.phaseActed(phSavageSpawn)
				break
			}
		}
	}
}
