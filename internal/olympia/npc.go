// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Complete npc.c:87–145. Selection always shuffles the land exits before
// considering the previous direction, using the shared sequential exdr stream.
func getExitDir(l []*exitView, dir int) *exitView {
	for _, e := range l {
		if int(e.direction) == dir {
			return e
		}
	}
	return nil
}

func (w *world) chooseNPCDirection(who, where, dir int) *exitView {
	l := w.exitsFromLocNSEWSelect(who, where, exitLand, 1)
	if len(l) == 0 {
		return nil
	}
	w.streams.BeginNPC(where)
	if dir != 0 && w.streams.NPCBehavior(who, dir, 1, 10) < 10 {
		if e := getExitDir(l, dir); e != nil {
			return e
		}
	}
	return l[0]
}

func (w *world) npcMove(who int) {
	where := w.subloc(who)
	if w.locDepth(where) != locProvince {
		w.queue(who, "move out")
		return
	}
	dir := 0
	if p := w.rpMisc(who); p != nil {
		dir = int(p.npcDir)
	}
	if e := w.chooseNPCDirection(who, where, dir); e != nil {
		w.pMisc(who).npcDir = int8(e.direction)
		w.queue(who, "move %s", fullDirS[e.direction])
	}
}

// Complete npc.c:202–348 entry hooks. Autonomous hunting/programs remain
// separate work; these hooks create real units and queue real commands.
func (w *world) createHadesBandit(who, where int) int {
	var item int
	switch w.streams.HadsBanditKind(who, where) {
	case 1:
		item = itemSpirit
	case 2:
		item = itemCorpse
	case 3:
		item = itemSavage
	case 4:
		item = itemSkeleton
	case 5:
		item = itemGorgon
	default:
		assert(false, "olympia/npc.c", 232, "FALSE")
	}
	n := w.newChar(subNi, item, where, -1, int(w.indepPlayer), loyNpc, 0, "")
	// Preserve C's order: it dereferences new before testing for failure.
	w.pChar(n).breakPoint = 0
	w.rpChar(n).npcProg = progHadesBandit
	if n < 0 {
		return -1
	}
	w.genItem(n, item, w.streams.HadsBanditQty(n))
	w.commandOut(where, "%s appear.", w.boxName(n))
	return n
}

func (w *world) createFaeryBandit(who, where int) int {
	var item int
	switch w.streams.FaerBanditKind(who, where) {
	case 1, 2:
		item = itemFaery
	case 3:
		item = itemElf
	default:
		assert(false, "olympia/npc.c", 271, "FALSE")
	}
	n := w.newChar(subNi, item, where, -1, int(w.indepPlayer), loyNpc, 0, "")
	w.pChar(n).breakPoint = 0
	w.rpChar(n).npcProg = progFaeryBandit
	if n < 0 {
		return -1
	}
	w.genItem(n, item, w.streams.FaerBanditQty(n))
	w.commandOut(where, "%s appear.", w.boxName(n))
	w.genItem(n, itemGold, w.streams.FaerBanditGold(n))
	return n
}

func (w *world) hadesAttackCheck(who, where int) {
	// The draw precedes every predicate, including NPC/hidden/alone safety.
	if w.streams.HadsAmbush(who, where) > 6 {
		return
	}
	if w.isNPC(who) || w.kind(who) != tChar || w.charReallyHidden(who) ||
		(w.hasSkill(who, skTranscendDeath) != 0 && w.stackParent(who) == 0 && w.countStackAny(who) == 1) {
		return
	}
	n := w.createHadesBandit(who, where)
	if n < 0 {
		return
	}
	w.queue(n, "wait time 0")
	w.initLoadSup(n)
	if w.streams.HadsRetal(who, where) == 1 {
		w.queue(n, "attack %s", boxCodeLess(who))
	}
}

func (w *world) faeryAttackCheck(who, where int) {
	if w.streams.FaerAmbush(who, where) > 6 {
		return
	}
	if w.isNPC(who) || w.kind(who) != tChar || w.charReallyHidden(who) || w.stackHasUseKey(who, useFaeryStone) != 0 {
		return
	}
	n := w.createFaeryBandit(who, where)
	if n < 0 {
		return
	}
	w.queue(n, "wait time 0")
	w.initLoadSup(n)
	if w.streams.FaerRetal(who, where) == 1 {
		w.queue(n, "attack %s", boxCodeLess(who))
	}
}
