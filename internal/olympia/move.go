// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/move.c: restore_stack_actions(), which check_db()
// uses, and init_ocean_chars(), which load_db() runs. The rest of move.c is
// left for later stages.

// restoreStackActions is restore_stack_actions(): it clears the moving day
// of who and every unit stacked under it.
func (w *world) restoreStackActions(who int) {
	for i := range w.loopStack(who) {
		w.pChar(i).moving = 0
	}
}

// initOceanChars is init_ocean_chars(): every character whose subloc is an
// ocean, in chain order, is appended to ocean_chars, for
// check_ocean_chars() to watch.
func (w *world) initOceanChars() {
	for i := range w.loopKind(tChar) {
		where := w.subloc(i)

		if w.subkind(where) == subOcean {
			w.oceanChars.Append(i)
		}
	}
}
