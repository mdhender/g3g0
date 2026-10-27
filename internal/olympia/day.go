// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/day.c: location-touching state, including day 0.

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

func (w *world) stormOwnerTouchLoc() {
	for i := range w.loopKind(tStorm) {
		if p := w.rpMisc(i); p != nil && p.summonedBy != 0 {
			where := w.subloc(i)
			if pl := w.player(int(p.summonedBy)); pl != 0 {
				w.touchLocPl(pl, where)
			}
		}
	}
}

func (w *world) initLocsTouched() {
	for who := range w.loopKind(tChar) {
		if !w.isPrisoner(who) {
			w.touchLoc(who)
		}
	}
	w.stormOwnerTouchLoc()
}
