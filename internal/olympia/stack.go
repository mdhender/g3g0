// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/stack.c: stack_parent(), stack_leader() and here_pos(),
// used by check_db() and command scheduling. The rest is left for later stages.

// A loader-accepted nowhere character reaches a NULL dereference before
// C's assertion here. Preserve that upstream defect rather than repairing it:
// https://github.com/playbymail/olympia-g3/issues/347
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
