// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// isHostile is perm.c's unit-target predicate. In particular, a hostile
// disposition naming only the target's faction does not match here.
func (w *world) isHostile(who, targ int) bool {
	if w.player(who) == w.player(targ) {
		return false
	}
	if w.npcProgram(who) == progSublocMonster && !w.isNPC(targ) && w.onlyDefeatable(who) == 0 {
		return true
	}
	if w.subkind(who) == subGarrison {
		if p := w.rpMisc(who); p != nil && p.garrHost.Lookup(targ) >= 0 {
			return true
		}
	}
	if p := w.rpDisp(who); p != nil && p.hostile.Lookup(targ) >= 0 {
		return true
	}
	if p := w.rpDisp(w.player(who)); p != nil && p.hostile.Lookup(targ) >= 0 {
		return true
	}
	return false
}

func (w *world) clearAllAtt(who int) {
	if p := w.rpDisp(who); p != nil {
		p.neutral.Clear()
		p.hostile.Clear()
		p.defend.Clear()
	}
}
