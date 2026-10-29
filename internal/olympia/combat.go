// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

// vGuard is the complete combat.c:2947–2964 command start.
func (w *world) vGuard(c *command) int {
	where := w.subloc(int(c.who))
	if c.a != 0 {
		w.pChar(int(c.who)).guard = 1
		w.commandOut(int(c.who), "Will guard %s.", w.boxName(where))
		return 1
	}
	w.pChar(int(c.who)).guard = 0
	w.commandOut(int(c.who), "Will not guard %s.", w.boxName(where))
	return 1
}

// fight and wield retain combat.c:110–146's original fields and C int widths.
// findWield mutates only the three ratings, not the selected-item fields.
type fight struct {
	unit, kind, savNum, num, newHealth   int32
	attack, defense                      int32
	behind, missile                      int32
	protects, nprot, ally                int32
	prisoner, inside                     int32
	seizeSlot, survivedFatal             int32
	attackItem, defenseItem, missileItem int32
}

type wield struct {
	attack, defense, missile int32
}

// findWield is combat.c:489–544 find_wield, including the optional fight
// mutation and NULL selection destination. Strict comparison makes ties first-wins.
func (w *world) findWield(selected *wield, who int, f *fight) bool {
	var local wield
	if selected == nil {
		selected = &local
	}
	*selected = wield{}
	attackMax, defenseMax, missileMax := int32(-1), int32(-1), int32(-1)
	assert(w.validBox(who), "olympia/combat.c", 506, "valid_box(who)")
	for e := range w.loopInv(who) {
		p := w.rpItemMagic(int(e.item))
		if p == nil {
			continue
		}
		if n := int32(p.attackBonus); n != 0 && n > attackMax {
			attackMax, selected.attack = n, e.item
		}
		if n := int32(p.defenseBonus); n != 0 && n > defenseMax {
			defenseMax, selected.defense = n, e.item
		}
		if n := int32(p.missileBonus); n != 0 && n > missileMax {
			missileMax, selected.missile = n, e.item
		}
	}
	if f != nil {
		if selected.attack != 0 {
			f.attack += attackMax
		}
		if selected.defense != 0 {
			f.defense += defenseMax
		}
		if selected.missile != 0 {
			f.missile += missileMax
		}
	}
	return selected.attack != 0 || selected.defense != 0 || selected.missile != 0
}

// wieldS is combat.c:547–588 wield_s. Only attack duplicates are removed:
// a shared defense/missile item is deliberately still named twice.
func (w *world) wieldS(who int) string {
	var selected wield
	if !w.findWield(&selected, who, nil) {
		return ""
	}
	if selected.attack == selected.missile {
		selected.missile = 0
	}
	if selected.attack == selected.defense {
		selected.defense = 0
	}
	var s string
	if selected.attack != 0 || selected.missile != 0 {
		switch {
		case selected.attack == 0:
			s = sout(", wielding %s", w.boxName(int(selected.missile)))
		case selected.missile == 0:
			s = sout(", wielding %s", w.boxName(int(selected.attack)))
		default:
			s = sout(", wielding %s and %s", w.boxName(int(selected.attack)), w.boxName(int(selected.missile)))
		}
	}
	if selected.defense != 0 {
		s += sout(", wearing %s", w.boxName(int(selected.defense)))
	}
	return s
}

// fortCovers is combat.c's fort_covers(): the number a fort protects.
// Rafts, completed or not, provide no fort protection.
func (w *world) fortCovers(n int) int {
	switch w.subkind(n) {
	case subCastle, subCastleNotdone:
		return 500
	case subTower, subTowerNotdone:
		return 100
	case subGalley, subGalleyNotdone, subRoundship, subRoundshipNotdone,
		subTemple, subTempleNotdone, subInn, subInnNotdone,
		subMine, subMineNotdone, subSewer:
		return 50
	case subRaft, subRaftNotdone:
		return 0
	default:
		fmt.Fprintf(w.stderr, "subkind is %s\n", subkindS[w.subkind(n)])
		assert(false, "olympia/combat.c", 296, "FALSE")
	}
	return 0
}

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
	w.commandOut(who, "> [auto-attack %s]", boxCodeLess(target))
	panic(notImplemented("regular_combat"))
}

func (w *world) clearSecondWaits() {
	for _, i := range w.secondWaitList {
		w.rpCommand(i).secondWait = 0
	}
	w.secondWaitList.Clear()
}
