// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package lists

// ILTest is ilist_test() in lib/ilist.c, which olympia-g3 -t runs. It prints
// its progress with printf. The caller does the setbuf(stdout, NULL) it
// starts with. Its assertions cannot fail; if one did, C would abort, and
// this panics.
func ILTest(printf func(format string, a ...any)) {
	var l, c IList

	printf("ilist_test: starting\n")

	printf("ilist_test: len of empty list = %d\n", l.Len())
	assert(l.Len() == 0)

	const N = 100
	var items [N]int
	for i := range N {
		items[i] = i + 1
	}

	for i := range N {
		l.Append(items[i])
	}

	printf("ilist_test: after append, len = %d\n", l.Len())
	assert(l.Len() == N)

	for i := range N {
		assert(l[i] == items[i])
	}

	c = l.Copy()
	assert(c.Len() == N)
	for i := range N {
		assert(&c[i] != &l[i])
		assert(c[i] == l[i])
	}
	printf("ilist_test: copy verified\n")

	l.Delete(50)
	assert(l.Len() == N-1)
	printf("ilist_test: after delete, len = %d\n", l.Len())

	l.Prepend(items[50])
	assert(l.Len() == N)
	assert(l[0] == items[50])
	printf("ilist_test: after prepend, len = %d\n", l.Len())

	assert(l.Lookup(items[50]) == 0)
	printf("ilist_test: lookup verified\n")

	l.Clear()
	assert(l.Len() == 0)
	printf("ilist_test: after clear, len = %d\n", l.Len())

	l.Reclaim()
	c.Reclaim()
	assert(l == nil)
	assert(c == nil)

	printf("ilist_test: all tests passed\n")
}

// assert stands in for C's assert(), which aborts.
func assert(ok bool) {
	if !ok {
		panic("lists: ilist_test assertion failed")
	}
}
