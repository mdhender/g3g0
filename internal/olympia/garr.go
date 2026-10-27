// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/garr.c: garrison_here(), which parse_arg() uses. The
// rest of garr.c is left for later stages.

// garrisonHere is garrison_here(): the first unit in where's here list if
// it is a garrison, or 0.
func (w *world) garrisonHere(where int) int {
	n := w.firstCharacter(where)

	if n != 0 && w.subkind(n) == subGarrison {
		return n
	}

	return 0
}
