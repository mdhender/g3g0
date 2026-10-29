// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

const provOrCity = -1

// npc.c:407–466, including the original NULL sentinel.
type cookieMonsterEntry struct {
	cookie, kind, sk, ni, terrain, manKind, low, high int
	notHere, noCookies                                *string
}

var cookieMonster = [...]cookieMonsterEntry{
	{itemMobCookie, tChar, subNi, itemAngryPeasant, provOrCity, itemAngryPeasant, 12, 36,
		strSave("Mobs can only be raised in provinces and cities."), strSave("A mob has already been raised from this place.")},
	{itemUndeadCookie, tChar, subUndead, 0, subGraveyard, itemCorpse, 15, 25,
		strSave("Demon lords may only be summoned in graveyards."), strSave("A demon lord has already been summoned from this graveyard.")},
	{itemRainCookie, tStorm, subRain, 0, 0, 0, 0, 0,
		strSave("Rain may not be summoned here."), strSave("A storm has already been summoned from this province.")},
	{itemWindCookie, tStorm, subWind, 0, 0, 0, 0, 0,
		strSave("Wind may not be summoned here."), strSave("A storm has already been summoned from this province.")},
	{itemFogCookie, tStorm, subFog, 0, 0, 0, 0, 0,
		strSave("Fog may not be summoned here."), strSave("A storm has already been summoned from this province.")},
	{},
}

// npc.c:469–481.
func (w *world) findCookie(k int) *cookieMonsterEntry {
	assert(w.kind(k) == tItem, "olympia/npc.c", 474, "kind(k) == T_item")
	for i := range cookieMonster {
		if cookieMonster[i].cookie == 0 {
			break
		}
		if cookieMonster[i].cookie == k {
			return &cookieMonster[i]
		}
	}
	return nil
}

// npc.c:484–534. Terrain diagnostics precede depleted-cookie diagnostics;
// a zero summoner suppresses both messages, not validation.
func (w *world) mayCookieNPC(who, where, cookie int) bool {
	t := w.findCookie(cookie)
	assert(t != nil, "olympia/npc.c", 491, "t")
	sk := w.subkind(where)
	badPlace := t.terrain > 0 && sk != t.terrain
	if t.terrain == 0 {
		if cookie == itemWindCookie && sk != subPlain && sk != subMountain && sk != subDesert && sk != subOcean {
			badPlace = true
		}
		if cookie == itemRainCookie && sk != subForest && sk != subOcean {
			badPlace = true
		}
		if cookie == itemFogCookie && sk != subForest && sk != subSwamp && sk != subOcean {
			badPlace = true
		}
	}
	if t.terrain == provOrCity && sk != subCity && w.locDepth(where) != locProvince {
		badPlace = true
	}
	if badPlace {
		if who != 0 {
			w.commandOut(who, "%s", *t.notHere)
		}
		return false
	}
	if w.hasItem(where, cookie) == 0 {
		if who != 0 {
			w.commandOut(who, "%s", *t.noCookies)
		}
		return false
	}
	return true
}

// npc.c:537–590. Character placement uses newChar's independent ownership
// and stack semantics; storms instead use plain setWhere, with no troop draw.
func (w *world) doCookieNPC(who, where, cookie, place int) int {
	if !w.mayCookieNPC(who, where, cookie) {
		return 0
	}
	t := w.findCookie(cookie)
	assert(t != nil, "olympia/npc.c", 548, "t")
	var n int
	if t.kind == tChar {
		n = w.newChar(t.sk, t.ni, place, 100, int(w.indepPlayer), loyNpc, 0, "")
	} else {
		n = w.newEnt(t.kind, t.sk)
		if n > 0 {
			w.setWhere(n, place)
		}
	}
	if n <= 0 {
		return 0
	}
	if t.sk == subNi {
		w.pChar(n).health = -1
	}
	p := w.pMisc(n)
	p.npcHome, p.npcCookie = int32(where), int32(cookie)
	p.summonedBy, p.npcCreated = int32(who), int32(w.sysclock.turn)
	if t.manKind != 0 {
		w.streams.BeginNPC(where)
		w.genItem(n, t.manKind, w.streams.NPCQty(cookie, n, t.low, t.high))
	}
	w.consumeItem(where, cookie, 1)
	return n
}

// Complete npc.c:351–403. All direct children get queued UNSTACK orders;
// victim order is the saved here order, with a keyed index, not a shuffle.
func (w *world) autoBandit(who, prog int) {
	where := w.subloc(who)
	for i := range w.loopHere(who) {
		w.queue(who, "unstack %s", boxCodeLess(i))
	}
	var targets []int
	for i := range w.loopHere(where) {
		if w.kind(i) != tChar || w.isNPC(i) || w.charReallyHidden(i) {
			continue
		}
		if prog == progFaeryBandit && w.stackHasUseKey(i, useFaeryStone) != 0 {
			continue
		}
		if prog == progHadesBandit && w.hasSkill(i, skTranscendDeath) != 0 && w.stackParent(i) == 0 && w.countStackAny(i) == 1 {
			continue
		}
		targets = append(targets, i)
	}
	if len(targets) > 0 {
		w.streams.BeginNPC(where)
		victim := targets[w.streams.NPCVictim(who, len(targets))]
		w.queue(who, "attack %s", boxCodeLess(victim))
	} else {
		w.npcMove(who)
	}
}

// Complete npc.c:613–691. Savage initialization precedes Hades, then a
// snapshot of independent units dispatches without executing queued orders.
func (w *world) queueNpcOrders() {
	w.stage("queue_npc_orders()")
	w.phaseRan(phNpcUnsworn)
	w.initSavageAttacks()
	w.autoHades()
	for who := range w.loopUnits(int(w.indepPlayer)) {
		p := w.rpChar(who)
		if p != nil && p.loyKind == loySummon {
			continue
		}
		if w.isPrisoner(who) {
			continue
		}
		if c := w.rpCommand(who); c != nil && c.state != stateDone {
			continue
		}
		if _, ok := w.topOrder(int(w.indepPlayer), who); ok {
			continue
		}
		prog := 0
		if p != nil {
			prog = int(p.npcProg)
		}
		switch prog {
		case 0:
			switch w.subkind(who) {
			case 0:
				w.autoUnsworn(who)
				w.phaseActed(phNpcUnsworn)
			case subUndead:
				w.autoUndead(who)
				w.phaseActed(phNpcUndead)
			case subNi:
				switch w.nobleItem(who) {
				case itemSavage:
					w.autoSavage(who)
					w.phaseActed(phNpcSavage)
				case itemPeasant, itemAngryPeasant:
					w.autoMob(who)
					w.phaseActed(phNpcMob)
				}
			}
		case progBandit, progHadesBandit, progFaeryBandit:
			w.autoBandit(who, prog)
			w.phaseActed(phNpcBandit)
		case progSublocMonster:
			// Intentional C no-op, not an unfinished program.
		case progNpcToken:
			w.npcMove(who)
			w.phaseActed(phNpcToken)
		default:
			assert(false, "olympia/npc.c", 687, "FALSE")
		}
	}
}

// Complete npc.c:149–165. Building occupants do not even seed the stream.
func (w *world) autoUnsworn(who int) {
	where := w.subloc(who)
	if w.locDepth(where) == locBuild {
		return
	}
	w.streams.BeginNPC(where)
	if w.streams.NPCBehavior(who, 0, 1, 2) == 1 {
		if n := w.sublocHere(where, subCity); n != 0 && w.streams.NPCBehavior(who, 1, 1, 2) == 1 {
			w.queue(who, "move %s", boxCodeLess(n))
		} else {
			w.npcMove(who)
		}
	}
}

// Complete npc.c:168–199. Queued reclamation is not execution.
func (w *world) autoMob(who int) {
	p := w.rpMisc(who)
	if p == nil {
		fmt.Fprintf(w.stderr, "warning: mob's rp_misc is NULL, who=%d\n", who)
		return
	}
	w.streams.BeginNPC(w.subloc(who))
	if w.subloc(who) != int(p.npcHome) ||
		(int32(w.sysclock.turn)-p.npcCreated >= 5 && w.streams.NPCBehavior(who, 0, 1, 2) == 1) {
		w.queue(who, "wait time %d", w.streams.NPCBehavior(who, 1, 10, 20))
		w.queue(who, "reclaim \"disperses.\"")
		return
	}
}

// Complete npc.c:65–84. This is recursive all_here, not direct occupants;
// regular subkind and non-unsworn loyalty are the only human predicates.
func (w *world) controlledHumansHere(where int) bool {
	for _, i := range w.allHere(where) {
		if w.kind(i) == tChar && w.subkind(i) == 0 {
			if p := w.rpChar(i); p != nil && p.loyKind != loyUnsworn {
				return true
			}
		}
	}
	return false
}

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
