// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package lists ports G3's reallocing arrays: lib/ilist.c and the list
// modules generated from the same template (lib/accept_ents.c, admits.c,
// cstrings.c, exit_views.c, fights.c, flag_ents.c, item_ents.c, orders.c,
// req_ents.c, roads.c, skill_ents.c, tiles.c, trades.c and wait_args.c),
// declared in lib/lists.h.
//
// The fourteen template modules are the same code over a different element
// pointer type, and ilist.c is that code over int plus ilist_add,
// ilist_rem_value_uniq and ilist_shuffle_rng. One generic type, List, ports
// them all. In C a list is a pointer to the first element of a block whose
// two hidden leading slots hold the length and capacity. Here it is a slice:
//
//   - A NULL list is a nil List. Clear keeps a non-nil list non-nil (C keeps
//     the block), Copy of a nil list is nil, and Reclaim makes a list nil,
//     so code that tests a list against NULL can test it against nil.
//   - The capacity and its growth policy (ILIST_ALLOC slots, doubling) are
//     not visible to the engine, so they are not reproduced.
//   - C keeps the length in the block, so every copy of the pointer sees a
//     change of length (until a realloc moves the block). A Go slice header
//     carries its own length. Ports must keep one owner for each list and
//     pass a pointer to it, as the C code passes &list, never mutate a copy
//     of the header and expect the owner to see it.
//   - Lookup and RemValue compare elements with ==. For the pointer lists
//     that is pointer identity, as in C. For cstrings, C compares char
//     pointers, not their contents; a port that stores Go strings must not
//     rely on Lookup or RemValue to find a string. (The engine at the pinned
//     revision never calls cstrings_lookup or cstrings_rem_value.)
//   - Out-of-range deletes fail C's assert() and abort; here they panic.
//
// lib/checked_alloc.c has nothing to port: my_malloc, my_realloc and my_free
// are thin wrappers over calloc, realloc and free at the pinned revision (the
// canary allocator is commented out), and they print nothing. The lists call
// calloc and realloc directly; their "out of memory" exits cannot be reached
// from Go, where allocation failure is fatal to the runtime.
package lists

import (
	"fmt"

	"github.com/mdhender/g3g0/internal/rnd"
	"github.com/mdhender/g3g0/internal/rng"
)

// List is a reallocing array: ilist or one of the <name>_list types.
type List[T comparable] []T

// IList is ilist, the list of ints (entity numbers, mostly).
type IList = List[int]

// Add is ilist_add(): Append, unless n is already in the list.
func (l *List[T]) Add(n T) {
	if l.Lookup(n) == -1 {
		l.Append(n)
	}
}

// Append is ilist_append() and <name>_append(): add n at the end.
func (l *List[T]) Append(n T) {
	*l = append(*l, n)
}

// Clear is ilist_clear() and <name>_clear(): empty the list, keeping it
// non-nil if it was.
func (l *List[T]) Clear() {
	if *l != nil {
		*l = (*l)[:0]
	}
}

// Copy is ilist_copy() and <name>_copy(): a new list with the same elements,
// or nil for a nil list.
func (l List[T]) Copy() List[T] {
	if l == nil {
		return nil
	}
	return append(make(List[T], 0, cap(l)), l...)
}

// Delete is ilist_delete() and <name>_delete(): remove the element at index
// i, shifting the rest down. C asserts that the list is not NULL and that i
// is in range.
func (l *List[T]) Delete(i int) {
	if *l == nil {
		panic("lists: Delete: assertion *l != NULL failed")
	}
	if !(i >= 0 && i < len(*l)) {
		panic(fmt.Sprintf("lists: Delete: assertion i >= 0 && i < len failed (i %d, len %d)", i, len(*l)))
	}
	*l = append((*l)[:i], (*l)[i+1:]...)
}

// Len is ilist_len() and <name>_len().
func (l List[T]) Len() int {
	return len(l)
}

// Lookup is ilist_lookup() and <name>_lookup(): the index of the first
// element equal to n, or -1.
func (l List[T]) Lookup(n T) int {
	for i, v := range l {
		if v == n {
			return i
		}
	}
	return -1
}

// Prepend is ilist_prepend() and <name>_prepend(): add n at the front.
func (l *List[T]) Prepend(n T) {
	var zero T
	*l = append(*l, zero)
	copy((*l)[1:], (*l)[:len(*l)-1])
	(*l)[0] = n
}

// Reclaim is ilist_reclaim() and <name>_reclaim(): free the list, leaving
// it nil.
func (l *List[T]) Reclaim() {
	*l = nil
}

// RemValue is ilist_rem_value() and <name>_rem_value(): delete every element
// equal to n. (ilist.c walks the list backwards and the template forwards;
// the result is the same.)
func (l *List[T]) RemValue(n T) {
	for i := len(*l) - 1; i >= 0; i-- {
		if (*l)[i] == n {
			l.Delete(i)
		}
	}
}

// RemValueUniq is ilist_rem_value_uniq(): delete the last element equal to
// n, if any.
func (l *List[T]) RemValueUniq(n T) {
	for i := len(*l) - 1; i >= 0; i-- {
		if (*l)[i] == n {
			l.Delete(i)
			break
		}
	}
}

// Scramble is ilist_scramble() and <name>_scramble(), which only call the
// shuffle.
func (l List[T]) Scramble(src *rnd.Source) {
	l.Shuffle(src)
}

// Shuffle is ilist_shuffle() and <name>_shuffle(): Knuth's shuffle, drawing
// each partner with the global rnd(). Lists of fewer than two elements draw
// nothing.
func (l List[T]) Shuffle(src *rnd.Source) {
	n := len(l) - 1
	for i := 0; i < n; i++ {
		if r := src.Rnd(i, n); r != i {
			l[i], l[r] = l[r], l[i]
		}
	}
}

// ShuffleRng is ilist_shuffle_rng() and exit_views_shuffle_rng(): the same
// shuffle, drawing from the stream s. It is rng.Shuffle.
func (l List[T]) ShuffleRng(s *rng.Stream) {
	rng.Shuffle(l, s)
}
