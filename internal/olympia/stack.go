// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/stack.c: stack lookup, positioning, join/leave and
// their small dependencies, also used by new_char(), and movement escapes.
// STACK, UNSTACK and PROMOTE share these same topology owners.

func (w *world) herePos(who int) int {
	p := w.rpLocInfo(w.loc(who))
	assert(p != nil, "olympia/stack.c", 16, "p != NULL")
	ret := p.hereList.Lookup(who)
	assert(ret >= 0, "olympia/stack.c", 20, "ret >= 0")
	return ret
}

// herePrecedes is stack.c:30–51. The a test intentionally comes first:
// equal operands return true, and finding either operand is sufficient.
func (w *world) herePrecedes(a, b int) bool {
	if w.loc(a) != w.loc(b) {
		return false
	}
	p := w.rpLocInfo(w.loc(a))
	assert(p != nil, "olympia/stack.c", 40, "p != NULL")
	for _, i := range p.hereList {
		if i == a {
			return true
		} else if i == b {
			return false
		}
	}
	assert(false, "olympia/stack.c", 48, "FALSE")
	return false
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

// vStack is the complete stack.c:703–741 validation and join command.
func (w *world) vStack(c *command) int {
	who, target := int(c.who), int(c.a)
	if !w.checkCharGone(who, target) {
		return 0
	}
	if target == who {
		w.commandOut(who, "Can't stack beneath oneself.")
		return 0
	}
	if w.stackedBeneath(who, target) {
		w.commandOut(who, "Cannot stack beneath %s since %s is stacked under you.", w.boxName(target), w.justName(target))
		return 0
	}
	if w.isPrisoner(target) {
		w.commandOut(who, "Can't stack beneath prisoners.")
		return 0
	}
	if !w.willAdmit(target, who, target) {
		w.commandOut(who, "%s refuses to let us stack.", w.boxName(target))
		w.commandOut(target, "Refused to let %s stack with us.", w.boxName(who))
		return 0
	}
	w.joinStack(who, target)
	return 1
}

// vUnstack is the complete stack.c:744–787. Explicit self extraction leaves
// the rest of the stack behind; an argumentless leave takes its children.
func (w *world) vUnstack(c *command) int {
	who, target := int(c.who), int(c.a)
	if c.parse.Len()-1 < 1 {
		if w.stackParent(who) <= 0 {
			w.commandOut(who, "Not stacked under anyone.")
			return 0
		}
		w.leaveStack(who)
		return 1
	}
	if who == target {
		w.extractStackedUnit(who)
		return 1
	}
	if !w.validBox(target) || w.stackParent(target) != who {
		w.commandOut(who, "%s is not stacked beneath us.", c.parse[1])
		return 0
	}
	w.dropStack(who, target)
	return 1
}

// promoteAfter is stack.c:828–861: depth-first location order, not just
// siblings. Snapshot traversal and a-before-b tests are the original order.
func (w *world) promoteAfter(a, b int) bool {
	where := w.subloc(a)
	ret := 0
	assert(w.subloc(b) == where, "olympia/stack.c", 835, "subloc(b) == where")
	for _, i := range w.allCharHere(where) {
		if i == a {
			ret = a
			break
		} else if i == b {
			ret = b
			break
		}
	}
	if ret == a {
		return true
	}
	if ret == b {
		return false
	}
	assert(false, "olympia/stack.c", 858, "FALSE")
	return false
}

// vPromote is the complete stack.c:864–925. Foreign targets must be direct
// children or later siblings; own targets only need later location order.
func (w *world) vPromote(c *command) int {
	who, target := int(c.who), int(c.a)
	if c.parse.Len()-1 < 1 {
		w.commandOut(who, "Must specify which character to promote.")
		return 0
	}
	if w.kind(target) != tChar {
		w.commandOut(who, "%s is not a character.", c.parse[1])
		return 0
	}
	if w.isPrisoner(target) {
		w.commandOut(who, "Can't promote prisoners.")
		return 0
	}
	targPar := w.stackParent(target)
	if !w.checkCharHere(who, target) {
		return 0
	}
	if target == who {
		w.commandOut(who, "Can't promote oneself.")
		return 0
	}
	if w.player(who) == w.player(target) {
		if !w.promoteAfter(who, target) {
			w.commandOut(who, "%s already comes before us in location order.", w.boxName(target))
			return 0
		}
	} else if targPar != who && !w.herePrecedes(who, target) {
		w.commandOut(who, "Characters to be promoted must be stacked immediately beneath the promoter, or be listed after the promoter at the same level.")
		return 0
	}
	w.promoteStack(target, who)
	return 1
}

// promoteStack is the complete stack.c:527–558. setWhere appends the target
// before promote moves it, even when both units already have the same parent.
// Unlike unstack this does not restore moving flags or summon loyalty.
func (w *world) promoteStack(lower, higher int) {
	assert(!w.stackedBeneath(lower, higher), "olympia/stack.c", 533, "!stacked_beneath(lower, higher)")
	w.setWhere(lower, w.loc(higher))
	p := w.rpLocInfo(w.loc(higher))
	assert(p != nil, "olympia/stack.c", 544, "p != NULL")
	assert(p.hereList[len(p.hereList)-1] == lower, "olympia/stack.c", 546, "p->here_list[ilist_len(p->here_list)-1] == lower")
	pos := p.hereList.Lookup(higher)
	w.promote(lower, pos)
	w.commandOut(higher, "Promoted %s.", w.boxName(lower))
	w.commandOut(lower, "%s promoted us.", w.boxName(higher))
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

// takePrisoner is stack.c:561–612. Announce before suppressing the target's
// output; inheritance moves its prisoners before extraction reparents children.
func (w *world) takePrisoner(who, target int) {
	assert(who != target, "olympia/stack.c", 566, "who != target")
	assert(w.kind(who) == tChar, "olympia/stack.c", 567, "kind(who) == T_char")
	assert(w.kind(target) == tChar, "olympia/stack.c", 568, "kind(target) == T_char")
	ni := w.subkind(target) == subNi && w.beastCapturable(target)
	w.vectorStack(w.stackLeader(target), true)
	w.vectorAdd(who)
	if ni {
		w.commandOut(vect, "%s disbands.", w.boxName(target))
	} else {
		w.commandOut(vect, "%s is %s by %s.", w.boxName(target), "taken prisoner", w.boxName(who))
	}
	w.pChar(target).prisoner = 1
	howMany := takeAll
	if ni {
		howMany = takeNI
	}
	w.takeUnitItems(target, who, howMany)
	w.extractStackedUnit(target)
	w.interruptOrder(target)
	if ni {
		w.unitDeserts(target, 0, true, loyUnsworn, 0)
		w.putBackCookie(target)
		w.setWhere(target, 0)
		w.changeBoxKind(target, tDeadchar)
	} else {
		w.stack(target, who)
	}
}

// hasPrisoner is stack.c:615–632: direct here-list membership, not a recursive
// stack search, and unlike myPrisoner it does not explicitly test kind.
func (w *world) hasPrisoner(who, pris int) bool {
	for i := range w.loopHere(who) {
		if i == pris && w.isPrisoner(i) {
			return true
		}
	}
	return false
}

// prisonerGivesToday is stack.c:662–681. Append even when the ensuing escape
// check refuses transfer. Clear retains the allocated list across either clock
// change; each prisoner counts only its own earlier entries.
func (w *world) prisonerGivesToday(pris int) int {
	var n int32
	if w.givesTurn != int32(w.sysclock.turn) || w.givesDay != int32(w.sysclock.day) {
		w.givesToday.Clear()
		w.givesTurn = int32(w.sysclock.turn)
		w.givesDay = int32(w.sysclock.day)
	}
	for _, i := range w.givesToday {
		if i == pris {
			n++
		}
	}
	w.givesToday.Append(pris)
	return int(n)
}

// givePrisoner is stack.c:684–700. ESC_GIVE+n selects the keyed leaf for
// this prisoner's previous-give count, not the global length of givesToday.
func (w *world) givePrisoner(who, target, pris int) bool {
	if w.checkPrisonerEscape(pris, 2, int(int32(2)+int32(w.prisonerGivesToday(pris)))) {
		return false
	}
	w.movePrisoner(who, target, pris)
	w.commandOut(who, "Transferred prisoner %s to %s.", w.boxName(pris), w.boxName(target))
	w.commandOut(target, "%s transferred the prisoner %s to us.", w.boxName(who), w.boxName(pris))
	return true
}
