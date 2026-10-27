// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// checkAllAutoAttacks is combat.c's leader/health gate and target scan.
// Only an actual selected attack reaches the unported combat boundary.
func (w *world) checkAllAutoAttacks() {
	for i := range w.loopKind(tChar) {
		if w.stackParent(i) != 0 {
			continue
		}
		var health int8
		if p := w.rpChar(i); p != nil {
			health = p.health
		}
		if health != 100 && health != -1 {
			continue
		}
		_ = w.rpCommand(i)
		w.checkAutoAttackSup(i)
	}
}

func (w *world) checkAutoAttackSup(who int) {
	where := w.subloc(who)
	okay := true
	if w.inSafeNow(who) || w.charMoving(who) != 0 { // oly.h char_gone
		return
	}
	for i := range w.loopHere(where) {
		if !w.isHostile(who, i) {
			continue
		}
		var target int
		if w.kind(i) == tChar {
			if w.isPrisoner(i) || !w.charHere(who, i) {
				continue
			}
			target = w.stackLeader(i)
			if w.stackLeader(who) == target {
				continue
			}
			if n := w.onlyDefeatable(target); n != 0 && w.hasItem(who, int(n)) == 0 {
				continue
			}
		} else if w.kind(i) == tLoc || w.kind(i) == tShip {
			if w.kind(i) == tLoc && w.locHidden(i) != 0 && !w.testKnown(who, i) {
				continue
			}
			target = i
			var targWho int
			if w.locDepth(i) == locBuild {
				targWho = w.buildingOwner(target)
			} else {
				targWho = w.firstCharacter(target)
			}
			if targWho == 0 || w.player(who) == w.player(targWho) {
				continue
			}
		} else {
			continue
		}
		w.autoAttack(who, target)
		okay = false
		break
	}
	if okay && w.isShip(where) {
		outer := w.subloc(where)
		for i := range w.loopHere(outer) {
			if i == where || !w.isShip(i) {
				continue
			}
			if !w.isHostile(who, i) {
				continue
			}
			if w.buildingOwner(i) == 0 {
				continue
			}
			w.autoAttack(who, i)
			okay = false
			break
		}
	}
}

// isShip is the oly.h subkind macro, not a kind(T_ship) test.
func (w *world) isShip(n int) bool {
	return w.subkind(n) == subGalley || w.subkind(n) == subRoundship || w.subkind(n) == subRaft
}

func (w *world) autoAttack(who, target int) {
	// Garrison report fan-out is still stage 2. Stop at the actual
	// auto_attack entry, before its out(), rather than fake that output or
	// its player output-bit mutations. Ordinary units can reach regular_combat.
	if w.subkind(who) == subGarrison && !w.isPrisoner(who) {
		panic(notImplemented("auto_attack"))
	}
	w.commandOut(who, "> [auto-attack %s]", boxCodeLess(target))
	panic(notImplemented("regular_combat"))
}

func (w *world) clearSecondWaits() {
	for _, i := range w.secondWaitList {
		w.rpCommand(i).secondWait = 0
	}
	w.secondWaitList.Clear()
}
