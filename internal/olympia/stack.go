// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/stack.c: stack lookup, positioning, join/leave and
// their small dependencies, also used by new_char(), and movement escapes.
// Stack command handlers remain for later stages.

func (w *world) herePos(who int) int {
	p := w.rpLocInfo(w.loc(who))
	assert(p != nil, "olympia/stack.c", 16, "p != NULL")
	ret := p.hereList.Lookup(who)
	assert(ret >= 0, "olympia/stack.c", 20, "ret >= 0")
	return ret
}

// stackParent is stack_parent(): the unit who is stacked under, or 0.
func (w *world) stackParent(who int) int {
	n := w.loc(who)

	if w.kind(n) == tChar {
		return n
	}

	return 0
}

// stackLeader is stack_leader(): the unit at the top of who's stack, who
// itself if it is not stacked. who must be a unit, and a stack deeper than
// 1000 (a loop) fails an assertion.
func (w *world) stackLeader(who int) int {
	count := 0

	assert(w.kind(who) == tChar, "olympia/stack.c", 96, "kind(who) == T_char")

	n := who
	for w.kind(n) == tChar {
		who = n
		n = w.stackParent(n)

		assert(count < 1000, "olympia/stack.c", 104, "count++ < 1000") // infinite loop check
		count++
	}

	return who
}

func (w *world) stackedBeneath(a, b int) bool {
	assert(w.kind(a) == tChar, "olympia/stack.c", 119, "kind(a) == T_char")
	assert(w.kind(b) == tChar, "olympia/stack.c", 120, "kind(b) == T_char")
	if a == b {
		return false
	}
	for b > 0 {
		b = w.stackParent(b)
		if a == b {
			return true
		}
	}
	return false
}

func (w *world) promote(who, newPos int) {
	p := w.rpLocInfo(w.loc(who))
	assert(p != nil, "olympia/stack.c", 144, "p != NULL")
	whoPos := p.hereList.Lookup(who)
	assert(whoPos >= newPos, "olympia/stack.c", 148, "who_pos >= new_pos")
	for i := whoPos; i > newPos; i-- {
		p.hereList[i] = p.hereList[i-1]
	}
	p.hereList[newPos] = who
}

func (w *world) unstack(who int) {
	leader := w.stackLeader(who)
	assert(w.validBox(leader), "olympia/stack.c", 161, "valid_box(leader)")
	assert(w.subloc(leader) == w.loc(leader), "olympia/stack.c", 168, "subloc(leader) == loc(leader)")
	if p := w.rpMagic(who); p != nil && p.swearOnRelease != 0 {
		w.pMagic(who).swearOnRelease = 0
	}
	w.setWhere(who, w.subloc(leader))
	w.promote(who, w.herePos(leader)+1)
	w.restoreStackActions(who)
	if p := w.rpChar(who); p != nil && p.loyKind == loySummon {
		w.setLoyal(who, loyNpc, 0)
	}
}

func (w *world) leaveStack(who int) {
	leader := w.stackParent(who)
	if leader <= 0 {
		return
	}
	w.commandOut(leader, "%s unstacks from us.", w.boxName(who))
	w.commandOut(who, "%s unstacks from %s.", w.boxName(who), w.boxName(leader))
	w.vectorCharHere(who)
	w.commandOut(vect, "%s unstacks from %s.", w.boxName(who), w.boxName(leader))
	w.unstack(who)
}

func (w *world) stack(who, target int) {
	assert(w.stackParent(who) == 0, "olympia/stack.c", 215, "stack_parent(who) == 0")
	w.setWhere(who, target)
	w.pChar(who).moving = int32(w.charMoving(target))
	if !w.isPrisoner(who) {
		p := w.rpLocInfo(target)
		pos := -1
		for i, n := range p.hereList {
			if w.kind(n) == tChar && w.isPrisoner(n) {
				pos = i
				break
			}
		}
		if pos >= 0 {
			w.promote(who, pos)
		}
	}
}

func (w *world) joinStack(who, target int) {
	assert(!w.stackedBeneath(who, target), "olympia/stack.c", 238, "!stacked_beneath(who, target)")
	assert(!w.isPrisoner(target), "olympia/stack.c", 239, "!is_prisoner(target)")
	w.leaveStack(who)
	assert(w.subloc(target) == w.subloc(who), "olympia/stack.c", 243, "subloc(target) == subloc(who)")
	w.commandOut(who, "%s stacks beneath %s.", w.boxName(who), w.boxName(target))
	w.commandOut(target, "%s stacks beneath us.", w.boxName(who))
	w.vectorCharHere(who)
	w.commandOut(vect, "%s stacks beneath %s.", w.boxName(who), w.boxName(target))
	w.stack(who, target)
}

// checkPrisonerEscape is stack.c:255–286. The hounds belong to the
// prisoner's faction within the whole stack, not necessarily the captor's.
func (w *world) checkPrisonerEscape(who, chance, sub int) bool {
	leader := w.stackParent(who)
	hound := w.stackHasItem(who, itemHound)
	threshold := int32(chance) * 10
	if hound != 0 {
		threshold /= 2
	}
	n := w.streams.EntPrisoner(who, sub)
	if int32(n) > threshold {
		if hound != 0 {
			w.vectorStack(leader, true)
			if hound == 1 {
				w.commandOut(vect, "The hound is barking.")
			} else {
				w.commandOut(vect, "The hounds are barking.")
			}
		}
		return false
	}
	w.prisonerEscapes(who)
	return true
}

// prisonerEscapes is the complete stack.c:289–328, including actual swimming
// and recursive arrival effects through moveStack, never a modeled landing.
func (w *world) prisonerEscapes(who int) {
	leader := w.stackParent(who)
	w.commandOut(leader, "Prisoner %s escaped!", w.boxName(who))
	w.pChar(who).prisoner = 0
	w.pMagic(who).swearOnRelease = 0
	w.unstack(who)
	w.touchLoc(who)
	w.commandOut(who, "We escaped!")
	where := w.subloc(who)
	if w.locDepth(where) <= locProvince {
		return
	}
	outOne := w.loc(where)
	if w.isShip(where) && w.subkind(outOne) == subOcean {
		outOne = w.findNearestLand(outOne)
		w.commandOut(who, "After jumping over the side of the boat and enduring a long, grueling, swim, we finally washed ashore at %s.", w.boxName(outOne))
		w.commandOut(leader, "%s jumped overboard and presumably drowned.", w.justName(who))
		w.logWrite(11, "!! Someone swam ashore, who=%s", boxCodeLess(who))
	}
	w.moveStack(who, outOne)
}

// prisonerMovementEscapeCheck is stack.c:331–342. loop_char_here builds a
// depth-first snapshot before any escapes; nested prisoners still get checked
// after their parent's escape has changed the stack and restored moving flags.
func (w *world) prisonerMovementEscapeCheck(who int) {
	for _, i := range w.allCharHere(who) {
		if w.isPrisoner(i) {
			w.checkPrisonerEscape(i, 2, 1) // stack.h ESC_MOVE
		}
	}
}

// dropStack is stack.c:377–442, including the natural disappearance ->
// desertion -> extraction closure. It does not invoke death/reclamation.
func (w *world) dropStack(who, toDrop int) {
	releaseSwear := false
	assert(w.stackParent(toDrop) == who, "olympia/stack.c", 382, "stack_parent(to_drop) == who")
	if w.isPrisoner(toDrop) {
		w.pChar(toDrop).prisoner = 0
		w.touchLoc(toDrop)
		w.commandOut(who, "Freed prisoner %s.", w.boxName(toDrop))
		w.commandOut(toDrop, "%s set us free.", w.boxName(who))
		if p := w.rpMagic(toDrop); p != nil && p.swearOnRelease != 0 {
			releaseSwear = true
		}
	} else {
		w.commandOut(who, "Dropped %s from stack.", w.boxName(toDrop))
		w.commandOut(toDrop, "%s dropped us from the stack.", w.boxName(who))
		w.vectorCharHere(toDrop)
		w.commandOut(vect, "%s dropped %s from the stack.", w.boxName(who), w.boxName(toDrop))
	}
	w.unstack(toDrop)
	if releaseSwear {
		w.logWrite(11, "%s frees a swear_on_release prisoner", w.boxName(who))
		if w.streams.EntDrop(who, toDrop, 0, 1, 5) < 5 {
			w.commandOut(who, "%s is grateful for your gallantry.", w.boxName(toDrop))
			w.commandOut(who, "%s pledges fealty to us.", w.boxName(toDrop))
			w.setLord(toDrop, w.player(who), loyOath, 1)
		} else {
			switch w.streams.EntDrop(who, toDrop, 1, 1, 3) {
			case 1:
				w.commandOut(who, "%s spits on you, and vanishes in a cloud of orange smoke.", w.boxName(toDrop))
			case 2:
				w.commandOut(who, "%s cackles wildly and vanishes.", w.boxName(toDrop))
			case 3:
				w.commandOut(who, "%s smiles briefly at you, then vanishes.", w.boxName(toDrop))
			}
			w.unitDeserts(toDrop, 0, true, loyUnsworn, 0)
			w.putBackCookie(toDrop)
			w.setWhere(toDrop, 0)
			w.changeBoxKind(toDrop, tDeadchar)
		}
	}
}

func (w *world) extractStackedUnit(who int) {
	first := 0
	for i := range w.loopHere(who) {
		if w.kind(i) == tChar && !w.isPrisoner(i) {
			first = i
			break
		}
	}
	if first != 0 {
		for i := range w.loopHere(who) {
			if i != first && w.kind(i) == tChar {
				w.setWhere(i, first)
			}
		}
		w.setWhere(first, w.loc(who))
		w.promote(first, w.herePos(who)+1)
	}
	for i := range w.loopHere(who) {
		if w.kind(i) == tChar && w.isPrisoner(i) {
			w.prisonerEscapes(i)
		}
	}
	w.leaveStack(who)
}

func (w *world) movePrisoner(who, target, pris int) int {
	rs := false
	if p := w.rpMagic(pris); p != nil {
		rs = p.swearOnRelease != 0
	}
	w.unstack(pris)
	w.stack(pris, target)
	if rs {
		w.pMagic(pris).swearOnRelease = 1
	}
	assert(w.rpChar(pris).prisoner != 0, "olympia/stack.c", 646, "rp_char(pris)->prisoner")
	return 0
}
