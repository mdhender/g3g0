// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// rpAdmit deliberately uses pPlayer, as does C: a lookup may allocate it.
func (w *world) rpAdmit(pl, targ int) *admit {
	assert(w.kind(pl) == tPlayer, "olympia/perm.c", 16, "kind(pl) == T_player")
	for _, p := range w.pPlayer(pl).admits {
		if int(p.targ) == targ {
			return p
		}
	}
	return nil
}

func (w *world) defaultGarrison(n int) int8 {
	if p := w.rpMagic(n); p != nil {
		return p.defaultGarr
	}
	return 0
}

func (w *world) willAdmit(pl, who, targ int) bool {
	if w.defaultGarrison(targ) != 0 {
		return true
	}
	pl = w.player(pl)
	if w.player(who) == pl {
		return true
	}
	p := w.rpAdmit(pl, targ)
	if p == nil {
		return false
	}
	found := p.l.Lookup(who) >= 0 || p.l.Lookup(w.player(who)) >= 0
	if p.sense != 0 {
		return !found
	}
	return found
}

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
