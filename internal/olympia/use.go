// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/use.c: skill_school(), which check_db() uses. The
// rest of use.c is left for later stages.

// skillSchool is skill_school(): the school of skill sk, found by
// following its required skills up to one that requires none. A chain
// longer than 1000 (a loop) fails an assertion.
func (w *world) skillSchool(sk int) int {
	count := 0

	for {
		assert(count < 1000, "olympia/use.c", 1182, "count++ < 1000")
		count++

		n := w.reqSkill(sk)
		if n == 0 {
			return sk
		}
		sk = n
	}
}
