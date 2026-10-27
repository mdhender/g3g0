// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/swear.c: lord() and player(), which write_master()
// uses.

// lord is lord(): who box n's owner is: n itself if it is a player,
// otherwise its unit_lord, or 0 if it has no entity_char.
func (w *world) lord(n int) int {
	if w.kind(n) == tPlayer {
		return n
	}

	if p := w.rpChar(n); p != nil {
		return int(p.unitLord)
	}
	return 0
}

// player is player(): the player box n belongs to, following its lords up
// to a player, or 0. A chain longer than 5000 lords (a loop) fails an
// assertion.
func (w *world) player(n int) int {
	loopCheck := 5000

	for n != 0 && w.kind(n) != tPlayer {
		n = w.lord(n)
		assert(loopCheck > 0, "olympia/swear.c", 32, "loop_check-- > 0")
		loopCheck--
	}

	return n
}
