// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/seed.c: int_comp(). The rest of seed.c, which builds
// a new world, is left for stage 5.

// intComp is int_comp(), the qsort() comparison of ints: a - b, as a C int,
// which wraps if the two are more than INT_MAX apart. The lists it sorts
// hold box numbers, which never are.
func intComp(a, b int) int {
	return int(int32(a) - int32(b))
}
