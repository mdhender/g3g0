// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/day.c: touch_loc() and touch_loc_pl(), which
// set_lord() calls. The rest of day.c is left for later stages.

// touchLocPl is touch_loc_pl(): it marks where's viewing location, and the
// location around that if it is inside a province, as touched by player pl.
func (w *world) touchLocPl(pl, where int) {
	p := w.pPlayer(pl)
	where = w.viewloc(where)

	if where == 0 {
		return
	}

	setBit(&p.locs, where)

	if w.locDepth(where) > locProvince {
		inside := w.loc(where)
		setBit(&p.locs, inside)
	}
}

// touchLoc is touch_loc(): touchLocPl for who's player and location.
func (w *world) touchLoc(who int) {
	where := w.subloc(who)

	w.touchLocPl(w.player(who), where)
}
