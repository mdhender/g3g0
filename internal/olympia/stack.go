// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/stack.c: stack lookup, positioning, join/leave and
// their small dependencies, also used by new_char(). Prisoner handling and
// stack command handlers are left for later stages.

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
