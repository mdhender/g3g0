// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"iter"

	"github.com/mdhender/g3g0/internal/lists"
)

// The entity loops of olympia/loop.h, as iterators. Each visits boxes in
// the order its C macro does, and has the macro's behavior when the body
// changes the chain it walks:
//
//   - loopKind and loopSubkind read the next box before the body runs, so
//     the body may delete the current box (the "delete problem" loop.h
//     guards against). A box the body links in after the current one is
//     visited; one linked in before it is not.
//   - loopLocOrShip reads the next box after the body runs, as
//     next_loc_or_ship does, so the body must not unlink the current box.
//   - loopBoxes tests each box number in turn as the loop reaches it.
//
// break works as it does in C. Of the loops that copy a list first,
// loopUnits is here; loop_here and the rest belong with the code that
// builds those lists.

// loopBoxes is loop_boxes: every box that is not T_deleted, in box number
// order.
func (w *world) loopBoxes() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 1; i < maxBoxes; i++ {
			if w.kind(i) != tDeleted {
				if !yield(i) {
					return
				}
			}
		}
	}
}

// loopKind is loop_kind: the chain of kind k (loop_char is loopKind(tChar),
// and so on).
func (w *world) loopKind(k int) iter.Seq[int] {
	return func(yield func(int) bool) {
		next := w.kindFirst(k)
		for i := next; i > 0; i = next {
			next = w.kindNext(i)
			if !yield(i) {
				return
			}
		}
	}
}

// loopSubkind is loop_subkind: the chain of subkind sk (loop_city is
// loopSubkind(subCity), and so on).
func (w *world) loopSubkind(sk int) iter.Seq[int] {
	return func(yield func(int) bool) {
		next := w.subFirst(sk)
		for i := next; i > 0; i = next {
			next = w.subNext(i)
			if !yield(i) {
				return
			}
		}
	}
}

// loopLocOrShip is loop_loc_or_ship: the ships, then the locations.
func (w *world) loopLocOrShip() iter.Seq[int] {
	return func(yield func(int) bool) {
		state := true
		i := w.kindFirst(tShip)
		if i <= 0 {
			i = w.kindFirst(tLoc)
			state = false
		}
		for i > 0 {
			if !yield(i) {
				return
			}
			i = w.kindNext(i)
			if i <= 0 && state {
				i = w.kindFirst(tLoc)
				state = false
			}
		}
	}
}

// loopUnits is loop_units: the units of player pl, from a copy of its
// units list taken when the loop starts, so the body may change the list.
// A player with no entity_player has none.
func (w *world) loopUnits(pl int) iter.Seq[int] {
	return func(yield func(int) bool) {
		var l lists.IList
		if p := w.rpPlayer(pl); p != nil {
			l = p.units.Copy()
		}
		for _, i := range l {
			if !yield(i) {
				return
			}
		}
	}
}
