// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/move.c: restore_stack_actions(), which check_db()
// uses. The rest of move.c is left for later stages.

// restoreStackActions is restore_stack_actions(): it clears the moving day
// of who and every unit stacked under it.
func (w *world) restoreStackActions(who int) {
	for i := range w.loopStack(who) {
		w.pChar(i).moving = 0
	}
}
