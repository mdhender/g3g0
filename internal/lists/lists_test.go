// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package lists_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/mdhender/g3g0/internal/lists"
	"github.com/mdhender/g3g0/internal/rnd"
	"github.com/mdhender/g3g0/internal/rng"
)

// TestNilAndEmpty checks the NULL-versus-empty distinctions the C lists
// make: clear keeps the block, copy of NULL is NULL, reclaim makes NULL.
func TestNilAndEmpty(t *testing.T) {
	var l lists.IList
	l.Clear()
	if l != nil || l.Copy() != nil || l.Len() != 0 || l.Lookup(0) != -1 {
		t.Fatalf("operations on a NULL list must leave it NULL")
	}
	l.Append(1)
	l.Clear()
	if l == nil || l.Len() != 0 {
		t.Fatalf("clear made the list NULL")
	}
	if c := l.Copy(); c == nil || c.Len() != 0 {
		t.Fatalf("copy of an empty list is %v, want empty and non-NULL", c)
	}
	l.Reclaim()
	if l != nil {
		t.Fatalf("reclaim did not make the list NULL")
	}
}

// TestCopyIsIndependent checks that a copy shares no storage.
func TestCopyIsIndependent(t *testing.T) {
	var l lists.IList
	for i := range 5 {
		l.Append(i)
	}
	c := l.Copy()
	c[0] = 99
	c.Append(5)
	l.Delete(4)
	if !slices.Equal(l, lists.IList{0, 1, 2, 3}) || !slices.Equal(c, lists.IList{99, 1, 2, 3, 4, 5}) {
		t.Fatalf("l %v, c %v", l, c)
	}
}

// TestDeleteAsserts checks that the deletes C's assert() rejects panic.
func TestDeleteAsserts(t *testing.T) {
	for _, tc := range []struct {
		l lists.IList
		i int
	}{
		{nil, 0},
		{lists.IList{}, 0},
		{lists.IList{1, 2}, -1},
		{lists.IList{1, 2}, 2},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("Delete(%d) on %v did not panic", tc.i, tc.l)
				}
			}()
			tc.l.Delete(tc.i)
		}()
	}
}

// TestPointerLists checks that the template lists compare elements by
// pointer, as C does.
func TestPointerLists(t *testing.T) {
	type tile struct{ n int }
	a, b := &tile{1}, &tile{1}
	var l lists.List[*tile]
	l.Append(a)
	l.Append(b)
	l.Prepend(b)
	if l.Lookup(a) != 1 || l.Lookup(b) != 0 || l.Lookup(&tile{1}) != -1 {
		t.Fatalf("lookup does not compare pointers")
	}
	l.RemValue(b)
	if len(l) != 1 || l[0] != a {
		t.Fatalf("rem_value left %v", l)
	}
}

// TestShuffleDraws checks that Shuffle and Scramble draw from rnd() once per
// index but the last, and ShuffleRng from the stream, as C does.
func TestShuffleDraws(t *testing.T) {
	seed := [16]byte{1, 2, 3}
	for n := range 6 {
		l := make(lists.IList, n)
		src, ref := rnd.New(seed), rnd.New(seed)
		l.Scramble(src)
		for i := 0; i < n-1; i++ {
			ref.Rnd(i, n-1)
		}
		if src.Seed() != ref.Seed() {
			t.Errorf("n %d: scramble drew a different number of times", n)
		}

		var s rng.Stream
		l.ShuffleRng(&s)
		if want := max(n-1, 0); s.Counter != uint32(want) {
			t.Errorf("n %d: shuffle_rng drew %d times, want %d", n, s.Counter, want)
		}
	}
}

// TestILTest checks the -t output; tools/olympia-cli-compat/run.sh compares
// it with the C engine byte for byte.
func TestILTest(t *testing.T) {
	var b strings.Builder
	lists.ILTest(func(format string, a ...any) { fmt.Fprintf(&b, format, a...) })
	want := "ilist_test: starting\n" +
		"ilist_test: len of empty list = 0\n" +
		"ilist_test: after append, len = 100\n" +
		"ilist_test: copy verified\n" +
		"ilist_test: after delete, len = 99\n" +
		"ilist_test: after prepend, len = 100\n" +
		"ilist_test: lookup verified\n" +
		"ilist_test: after clear, len = 0\n" +
		"ilist_test: all tests passed\n"
	if b.String() != want {
		t.Errorf("got %q", b.String())
	}
}
