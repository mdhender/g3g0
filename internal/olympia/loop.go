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
//   - loopLocOrShip and loopProvince read the next box after the body
//     runs, as next_loc_or_ship and next_province do, so the body must not
//     unlink the current box.
//   - loopBoxes tests each box number in turn as the loop reaches it.
//
// break works as it does in C. The loops that copy a list first (loopUnits,
// loopHere, loopStack and loopInv) take the copy when the range starts, so
// the body may change the list. loop_here and loop_inv assert
// valid_box(where) or valid_box(who) at the call site, so the text and line
// of the assertion are the caller's; the port's callers assert before the
// loop, where it is reachable, and the loops do not.

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

// loopProvince is loop_province: the locations, in chain order, whose
// loc_depth() is LOC_province. Like loop_province, it reads the next
// location after the body runs, so the body must not unlink the current
// one.
func (w *world) loopProvince() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := w.kindFirst(tLoc); i > 0; i = w.kindNext(i) {
			if w.locDepth(i) == locProvince {
				if !yield(i) {
					return
				}
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

// loopHere is loop_here: the boxes in the here list of where, from a copy
// taken when the loop starts.
func (w *world) loopHere(where int) iter.Seq[int] {
	return func(yield func(int) bool) {
		l := w.rpLocInfo(where).hereList.Copy()
		for _, i := range l {
			if !yield(i) {
				return
			}
		}
	}
}

// loopStack is loop_stack: who and every unit stacked under it, from
// all_stack() when the loop starts.
func (w *world) loopStack(who int) iter.Seq[int] {
	return func(yield func(int) bool) {
		l := w.allStack(who)
		for _, i := range l {
			if !yield(i) {
				return
			}
		}
	}
}

// loopInv is loop_inv: who's inventory entries whose item is a valid box
// and whose qty is positive. The list of entries is copied when the loop
// starts; each entry is tested when the loop reaches it, and the body gets
// a copy of it, so changing *e changes nothing.
func (w *world) loopInv(who int) iter.Seq[*itemEnt] {
	return func(yield func(*itemEnt) bool) {
		l := w.bx[who].items.Copy()
		for _, p := range l {
			if w.validBox(int(p.item)) && p.qty > 0 {
				e := *p
				if !yield(&e) {
					return
				}
			}
		}
	}
}

// loopTrade is loop_trade: who's trades whose item is a valid box and whose
// qty is positive, from a copy of the list taken when the loop starts. The
// copy shares the trades, as trades_copy() does, so the body may change
// them. loop_trade asserts valid_box(who) at the call site; the one caller
// so far, investigatePossibleTrade, is reached only after add_item() has
// asserted it.
func (w *world) loopTrade(who int) iter.Seq[*trade] {
	return func(yield func(*trade) bool) {
		l := w.bx[who].trades.Copy()
		for _, t := range l {
			if w.validBox(int(t.item)) && t.qty > 0 {
				if !yield(t) {
					return
				}
			}
		}
	}
}
