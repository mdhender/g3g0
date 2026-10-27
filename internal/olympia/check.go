// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

// Ported from olympia/check.c: check_db() and its checks, which check
// database integrity and make minor repairs. Every repair is noted on
// stderr. The checks run in C's order, and each walks the boxes in its C
// loop's order, so the messages come out in the same order too.

// checkHere is check_here():
//
//  1. Go through every box. If the box claims it's in a location, but it
//     doesn't show up in the here list of the location, then add it to the
//     here list.
//  2. Go through every location. If the loc has a here list, see if each
//     box in the list claims that it's in that location. If not, remove it
//     from the here list.
//
// This gives precedence to the location the unit claims to be in over the
// location's here-list of units.
//
// The second loop walks a copy, so it checks that a stale duplicate has
// not already been removed before reporting and removing it.
//
// loop_here asserts valid_box(i) at check.c line 49, but loop_boxes yields
// only valid boxes and nothing here deletes one, so the assertion cannot
// fail and is left out.
func (w *world) checkHere() {
	for i := range w.loopBoxes() {
		where := w.loc(i)
		if where > 0 && !w.inHereList(where, i) {
			fmt.Fprintf(w.stderr, "\tcheck_here: adding [%d] to here list of [%d]\n", i, where)

			w.addToHereList(where, i)
		}
	}

	for i := range w.loopBoxes() {
		for j := range w.loopHere(i) {
			where := w.loc(j)
			if where != i && w.inHereList(i, j) {
				fmt.Fprintf(w.stderr, "\tcheck_here: removing [%d] from here list of [%d]\n", j, i)

				w.removeFromHereList(i, j)
			}
		}
	}
}

// checkSwear is check_swear():
//
//  1. For every unit, if it belongs to a player, make sure it appears in
//     the player's unit list.
//  2. For every player, go through the unit list seeing if the units
//     actually belong to the player. If not, remove them from the list.
//
// This gives precedence to the lord chain over the player's unit list.
func (w *world) checkSwear() {
	for i := range w.loopKind(tChar) {
		over := w.player(i)
		if over > 0 && !w.isUnit(over, i) {
			fmt.Fprintf(w.stderr, "\tcheck_swear: adding [%d] to player [%d]\n", i, over)
			w.pPlayer(over).units.Append(i)

			safeQsort(w.pPlayer(over).units, intComp)
		}
	}

	for i := range w.loopKind(tPlayer) {
		for j := range w.loopUnits(i) {
			over := w.player(j)
			if over != i {
				fmt.Fprintf(w.stderr, "\tcheck_swear: removing [%d] from player list of [%d]\n", j, i)

				w.pPlayer(i).units.RemValue(j)
			}
		}
	}
}

// checkSpecialPlayer is the part check_indep(), check_gm(),
// check_skill_player(), check_eat_player(), check_npc_player() and
// check_garr_player() share: if player box n does not exist, it reports
// that it creates it, as "\tFUNC: creating WHAT player [N]", and allocates
// it with subkind sk; it asserts, at check.c line assertLine, that n is a
// player (the condition's text is cond); and it names it name if it has no
// name.
func (w *world) checkSpecialPlayer(n, sk int, fn, what string, assertLine int, cond, name string) {
	if w.bx[n] == nil {
		fmt.Fprintf(w.stderr, "\t%s: creating %s player [%d]\n", fn, what, n)

		w.allocBox(n, tPlayer, sk)
	}

	assert(w.kind(n) == tPlayer, "olympia/check.c", assertLine, cond)

	if w.name(n) == "" {
		w.setName(n, name)
	}
}

// checkIndep is check_indep(): it makes sure the independent player
// exists, and swears every unit that belongs to no player to it.
func (w *world) checkIndep() {
	w.checkSpecialPlayer(int(w.indepPlayer), subPlNpc, "check_indep", "independent", 137,
		"kind(indep_player) == T_player", "Independent player")

	for i := range w.loopKind(tChar) {
		if w.player(i) == 0 {
			fmt.Fprintf(w.stderr, "\tcheck_indep: swearing unit [%d] to %s\n", i, w.boxName(int(w.indepPlayer)))

			w.setLord(i, int(w.indepPlayer), loyUnsworn, 0)
		}
	}
}

// checkGm is check_gm(): it makes sure the gm player exists.
func (w *world) checkGm() {
	w.checkSpecialPlayer(int(w.gmPlayer), subPlSystem, "check_gm", "gm", 168,
		"kind(gm_player) == T_player", "Gamemaster")
}

// checkSkillPlayer is check_skill_player(): it makes sure the skill player
// exists.
func (w *world) checkSkillPlayer() {
	w.checkSpecialPlayer(int(w.skillPlayer), subPlSystem, "check_skill_player", "skill", 187,
		"kind(skill_player) == T_player", "Skill list")
}

// checkEatPlayer is check_eat_player(): it makes sure the order eater
// player exists.
func (w *world) checkEatPlayer() {
	w.checkSpecialPlayer(int(w.eatPl), subPlSystem, "check_eat_player", "eat", 206,
		"kind(eat_pl) == T_player", "Order eater")
}

// checkNpcPlayer is check_npc_player(): it makes sure the npc player
// exists.
func (w *world) checkNpcPlayer() {
	w.checkSpecialPlayer(int(w.npcPl), subPlSilent, "check_npc_player", "npc", 225,
		"kind(npc_pl) == T_player", "NPC control")
}

// checkGarrPlayer is check_garr_player(): it makes sure the garrison
// player exists.
func (w *world) checkGarrPlayer() {
	w.checkSpecialPlayer(int(w.garrPl), subPlSilent, "check_garr_player", "garrison", 244,
		"kind(garr_pl) == T_player", "Garrison units")
}

// checkGlob is check_glob(): it checks that kind_s has T_MAX entries and
// subkind_s SUB_MAX, before their NULLs (see glob.go).
func (w *world) checkGlob() {
	assert(len(kindS) == tMax, "olympia/check.c", 263, "i == T_MAX")
	assert(len(subkindS) == subMax, "olympia/check.c", 267, "i == SUB_MAX")
}

// checkNowhere is check_nowhere(): it warns of every unit, and of every
// location or ship below a region, that is nowhere. Units are moved to
// nowhere_loc when it is valid. loc_depth() fails an
// assertion for a location or ship whose subkind is not a location's.
func (w *world) checkNowhere() {
	for i := range w.loopKind(tChar) {
		if w.loc(i) == 0 {
			fmt.Fprintf(w.stderr, "\twarning: unit %s is nowhere\n", w.boxCode(i))
			if w.validBox(int(w.nowhereLoc)) {
				w.setWhere(i, int(w.nowhereLoc))
				fmt.Fprintf(w.stderr, "\t\tmoved to %s\n", w.boxCode(int(w.nowhereLoc)))
			}
		}
	}

	for i := range w.loopLocOrShip() {
		if w.locDepth(i) > locRegion && w.loc(i) == 0 {
			fmt.Fprintf(w.stderr, "\twarning: loc %s is nowhere\n", w.boxCode(i))
		}
	}
}

// checkSkills is check_skills(): it warns of orphaned subskills, skills
// with no learn time, skills two schools offer, skills offered by a school
// they are not in, and skills that no school offers. It leaves in each
// skill's temp the last school that offers it.
func (w *world) checkSkills() {
	for sk := range w.loopKind(tSkill) {
		if sk >= 9000 && w.skillSchool(sk) == sk {
			fmt.Fprintf(w.stderr, "\twarning: orphaned subskill %s\n", w.boxCode(sk))
		}
		w.bx[sk].temp = 0
	}

	// offered and research are checked alike
	check := func(sk int, l []int) {
		for _, n := range l {
			if w.bx[n].temp != 0 {
				fmt.Fprintf(w.stderr, "\twarning: both %s and %s offer skill %d\n",
					w.boxName(sk), w.boxName(int(w.bx[n].temp)), n)
			} else {
				w.bx[n].temp = int32(sk)
			}

			if w.skillSchool(n) != sk {
				fmt.Fprintf(w.stderr, "\twarning: %s offers %d, but %d is in school %d\n",
					w.boxName(sk), n, n, w.skillSchool(n))
			}
		}
	}

	for sk := range w.loopKind(tSkill) {
		p := w.rpSkill(sk)

		if w.learnTime(sk) == 0 {
			fmt.Fprintf(w.stderr, "\twarning: learn time of %s is 0\n", w.boxName(sk))
		}

		if p == nil {
			continue
		}

		check(sk, p.offered)
		check(sk, p.research)
	}

	for sk := range w.loopKind(tSkill) {
		if w.skillSchool(sk) == sk {
			continue
		}

		if w.bx[sk].temp == 0 {
			fmt.Fprintf(w.stderr, "\twarning: non-offered skill %s\n", w.boxName(sk))
		}
	}
}

// checkItemCounts is check_item_counts(): it reports inventory entries
// that are not items, and checks that every unique item is held, once, by
// the box its who_has names, which it repairs to the box that holds it. It
// leaves in each item's temp the total quantity held.
//
// loop_inv asserts valid_box(i) at check.c line 401, but loop_boxes yields
// only valid boxes and nothing here deletes one, so the assertion cannot
// fail and is left out.
func (w *world) checkItemCounts() {
	w.clearTemps(tItem)

	for i := range w.loopBoxes() {
		for e := range w.loopInv(i) {
			item := int(e.item)

			if w.kind(item) != tItem {
				fmt.Fprintf(w.stderr, "\t%s has non-item %s\n", w.boxName(i), w.boxName(item))
				continue
			}

			if w.itemUnique(item) == 0 {
				continue
			}

			if w.itemUnique(item) != i {
				fmt.Fprintf(w.stderr, "\tunique item %s: whohas=%s, actual=%s\n",
					w.boxName(item), w.boxName(w.itemUnique(item)), w.boxName(i))
				w.pItem(item).whoHas = int32(i)
			}

			if e.qty != 1 {
				fmt.Fprintf(w.stderr, "\t%s has qty %d of unique item %s\n",
					w.boxName(i), e.qty, w.boxName(item))
			}

			w.bx[item].temp += e.qty
		}
	}

	for i := range w.loopKind(tItem) {
		if w.itemUnique(i) != 0 {
			if w.bx[i].temp != 1 {
				fmt.Fprintf(w.stderr, "\tunique item %s count %d\n", w.boxName(i), w.bx[i].temp)
			}
		}
	}
}

// checkLocNameLengths is check_loc_name_lengths(): it warns of every
// location whose name is longer than 25 bytes.
func (w *world) checkLocNameLengths() {
	for i := range w.loopKind(tLoc) {
		if len(w.justName(i)) > 25 {
			fmt.Fprintf(w.stderr, "\twarning: %s name too long\n", w.boxName(i))
		}
	}
}

// checkMoving is check_moving(): it stops every moving stack whose leader
// is not running a command, and makes every unit's moving day its stack
// leader's.
func (w *world) checkMoving() {
	for i := range w.loopKind(tChar) {
		if w.stackLeader(i) != i || w.charMoving(i) == 0 {
			continue
		}

		c := w.rpCommand(i)

		if c == nil || c.state != stateRun {
			fmt.Fprintf(w.stderr, "\t%s moving but no command\n", w.boxName(i))

			w.restoreStackActions(i)
		}
	}

	for i := range w.loopKind(tChar) {
		leader := w.stackLeader(i)

		if leader == i || w.charMoving(i) == w.charMoving(leader) {
			continue
		}

		fmt.Fprintf(w.stderr, "\t%s moving disagrees with leader\n", w.boxName(i))
		w.pChar(i).moving = int32(w.charMoving(leader))
	}
}

// checkPrisoner is check_prisoner(): it frees every prisoner that is not
// stacked under anyone.
func (w *world) checkPrisoner() {
	for who := range w.loopKind(tChar) {
		if !w.isPrisoner(who) {
			continue
		}

		if w.stackParent(who) == 0 {
			fmt.Fprintf(w.stderr, "\t%s prisoner but unstacked\n", w.boxName(who))
			w.pChar(who).prisoner = 0 // FALSE
		}
	}
}

// checkDb is check_db(): it checks database integrity and fixes minor
// problems in backlinks and lists, noting every correction on stderr.
func (w *world) checkDb() {
	w.stage("check_db()")

	w.checkGlob()
	w.checkHere()
	w.checkSwear()
	w.checkIndep()
	w.checkGm()
	w.checkSkillPlayer()
	w.checkEatPlayer()
	w.checkNpcPlayer()
	w.checkGarrPlayer()
	w.checkNowhere()
	w.checkSkills()
	w.checkItemCounts()
	w.checkLocNameLengths()
	w.checkMoving()
	w.checkPrisoner()

	if w.bx[w.garrisonMagic] != nil {
		fmt.Fprintf(w.stderr, "\twarning: %s should not be allocated,\n", w.boxName(int(w.garrisonMagic)))
		fmt.Fprintf(w.stderr, "\t\treserved for garrison_magic\n")
	}
}
