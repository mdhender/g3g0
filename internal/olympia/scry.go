// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// resetCastWhere is olympia/scry.c:37–51. Invalid projections are retained.
func (w *world) resetCastWhere(who int) int {
	where := 0
	if p := w.rpMagic(who); p != nil {
		where = int(p.projectCast)
	}
	if w.isLocOrShip(where) {
		w.pMagic(who).projectCast = 0
		return where
	}
	return w.subloc(who)
}
