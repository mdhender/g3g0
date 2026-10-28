// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Complete necro.c:363–387. Attack/pillage are queued, not executed here.
func (w *world) autoUndead(who int) {
	where := w.subloc(who)
	master := 0
	if p := w.rpMisc(who); p != nil {
		master = int(p.summonedBy)
	}
	if master != 0 && w.subloc(who) == w.subloc(master) {
		w.queue(who, "attack %s", boxCodeLess(master))
		w.pMisc(who).summonedBy = 0
		return
	}
	w.streams.BeginNPC(where)
	if w.locDepth(where) != locProvince || w.streams.NPCBehavior(who, 0, 1, 2) == 1 {
		w.npcMove(who)
		return
	}
	w.queue(who, "pillage 1")
}
