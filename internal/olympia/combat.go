// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// checkAllAutoAttacks is the actual combat.c scan, not a successful
// substitute for combat. Eligible leaders reach an explicit boundary.
func (w *world) checkAllAutoAttacks() {
	for i := range w.loopKind(tChar) {
		if w.stackParent(i) != 0 {
			continue
		}
		var health int8
		if p := w.rpChar(i); p != nil {
			health = p.health
		}
		if health != 100 && health != -1 {
			continue
		}
		_ = w.rpCommand(i)
		panic(notImplemented("check_auto_attack_sup"))
	}
}

func (w *world) clearSecondWaits() {
	for _, i := range w.secondWaitList {
		w.rpCommand(i).secondWait = 0
	}
	w.secondWaitList.Clear()
}
