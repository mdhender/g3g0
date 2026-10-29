// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

// incrementCurrentAura is day.c:722–774. Artifact recovery is one point
// per positive inventory entry with a nonzero bonus, not per quantity/bonus.
func (w *world) incrementCurrentAura() {
	w.phaseRan(phAura)
	for who := range w.loopKind(tChar) {
		if p := w.rpMagic(who); p == nil || p.magician == 0 {
			continue
		}
		ac, ma := w.hasAuraculum(who), w.maxEffAura(who)
		p := w.pMagic(who)
		before := p.curAura
		if p.curAura < ma {
			p.curAura++
		}
		if p.curAura < ma {
			p.curAura++
		}
		if ac != 0 {
			if p.curAura < ma {
				p.curAura++
			}
			if p.curAura < ma {
				p.curAura++
			}
		}
		for e := range w.loopInv(who) {
			if im := w.rpItemMagic(int(e.item)); im != nil && im.auraBonus != 0 && p.curAura < ma {
				p.curAura++
			}
		}
		if p.curAura > before {
			w.phaseActed(phAura)
		}
	}
}

// day.c:806–898 and 1945–1963: only positive shrouds/countdowns decay.
func (w *world) decrementAbilityShroud() {
	w.phaseRan(phAbilityShroud)
	for who := range w.loopKind(tChar) {
		if p := w.rpMagic(who); p != nil && p.abilityShroud > 0 {
			p.abilityShroud--
			w.phaseActed(phAbilityShroud)
		}
	}
}

func (w *world) decrementLocBarrier() {
	w.phaseRan(phBarrierDec)
	for where := range w.loopKind(tLoc) {
		if p := w.rpLoc(where); p != nil && p.barrier > 0 {
			p.barrier--
			w.phaseActed(phBarrierDec)
			if p.barrier == 0 {
				w.commandOut(where, "The barrier over %s has dissipated.", w.boxName(where))
				w.phaseActed(phBarrierGone)
			}
		}
	}
}

func (w *world) decrementRegionShroud() {
	w.phaseRan(phRegionShroud)
	for where := range w.loopKind(tLoc) {
		if p := w.rpLoc(where); p != nil && p.shroud > 0 {
			p.shroud--
			w.notifyLocShroud(where)
			w.phaseActed(phRegionShroud)
		}
	}
}

// notifyLocShroud is complete scry.c:333–359, including stacked recipients.
func (w *world) notifyLocShroud(where int) {
	p := w.rpLoc(where)
	if p == nil {
		return
	}
	for _, who := range w.allCharHere(where) {
		if w.hasSkill(who, skShroudRegion) != 0 {
			if p.shroud > 0 {
				w.commandOut(who, "The magical shroud over %s has diminished to %s aura.", w.boxName(where), niceNum(int(p.shroud)))
			} else {
				w.commandOut(who, "The magical shroud over %s has dissipated.", w.boxName(where))
			}
		}
	}
}

func (w *world) decrementMeditationHinder() {
	w.phaseRan(phMedHinder)
	for who := range w.loopKind(tChar) {
		if p := w.rpMagic(who); p != nil && p.hinderMeditation > 0 {
			p.hinderMeditation--
			w.phaseActed(phMedHinder)
		}
	}
}

// loyalS is u.c:1569–1605, also used in the real desertion log.
func (w *world) loyalS(who int) string {
	kind, rate := int8(0), int32(0)
	if p := w.rpChar(who); p != nil {
		kind, rate = p.loyKind, p.loyRate
	}
	var s string
	switch kind {
	case loyUnsworn:
		s = "unsworn"
	case loyContract:
		s = "contract"
	case loyOath:
		s = "oath"
	case loyFear:
		s = "fear"
	case loyNpc:
		s = "npc"
	case loySummon:
		s = "summon"
	default:
		assert(false, "olympia/u.c", 1601, "FALSE")
	}
	return fmt.Sprintf("%s-%d", s, rate)
}

// loyaltyDecay is day.c:1007–1075. Contract desertion tests the old rate;
// fear tests after decrement. Keep the live chain and full desertion closure.
func (w *world) loyaltyDecay() {
	w.phaseRan(phLoyDecay)
	for who := range w.loopKind(tChar) {
		p := w.rpChar(who)
		if p == nil || p.freshHire != 0 {
			continue
		}
		if p.loyKind == loyUnsworn || p.loyKind == loyOath || p.loyKind == loyNpc {
			continue
		}
		if p.loyRate <= 0 && p.loyKind == loySummon {
			w.leaveStack(who)
			w.setLoyal(who, loyNpc, 0)
			w.phaseActed(phLoyRelease)
			continue
		}
		if p.loyRate < 50 && p.loyKind == loyContract && w.streams.UpLoyal(who, 0, 1, 2) == 1 {
			w.logWrite(12, "%s deserts, %s", w.boxName(who), w.loyalS(who))
			w.unitDeserts(who, int(w.indepPlayer), true, loyUnsworn, 0)
			w.phaseActed(phLoyDesert)
			continue
		}
		switch p.loyKind {
		case loySummon:
			p.loyRate--
			w.phaseActed(phLoyDecay)
		case loyFear:
			p.loyRate--
			w.phaseActed(phLoyDecay)
			if p.loyRate <= 0 && w.streams.UpLoyal(who, 1, 1, 2) == 1 {
				w.logWrite(12, "%s deserts, %s", w.boxName(who), w.loyalS(who))
				w.unitDeserts(who, int(w.indepPlayer), true, loyUnsworn, 0)
				w.phaseActed(phLoyDesert)
			}
		case loyContract:
			p.loyRate -= max(int32(50), p.loyRate/10)
			w.phaseActed(phLoyDecay)
		default:
			assert(false, "olympia/day.c", 1068, "FALSE")
		}
		if p.loyRate < 0 {
			p.loyRate = 0
		}
	}
}

// relicDecay is day.c:1078–1108. Negative nonzero timers also decrement
// and transfer. Use real moveItem for unique ownership, inventory and lore.
func (w *world) relicDecay() {
	w.phaseRan(phRelicDecay)
	for item := range w.loopKind(tItem) {
		if w.subkind(item) != subRelic {
			continue
		}
		p := w.rpItemMagic(item)
		if p == nil || p.relicDecay == 0 {
			continue
		}
		p.relicDecay--
		w.phaseActed(phRelicDecay)
		if p.relicDecay > 0 {
			continue
		}
		owner := 0
		if ip := w.rpItem(item); ip != nil {
			owner = int(ip.whoHas)
		}
		if w.kind(owner) == tChar {
			w.commandOut(owner, "%s vanishes.", w.boxName(item))
		}
		w.moveItem(owner, int(w.nowhereLoc), item, 1)
		w.phaseActed(phRelicVanish)
	}
}

// pillageDecay is day.c:1111–1130; unlike shrouds, negative loot decays.
func (w *world) pillageDecay() {
	w.phaseRan(phPillage)
	for where := range w.loopKind(tLoc) {
		if p := w.rpSubloc(where); p != nil && p.loot != 0 && p.recentLoot == 0 {
			p.loot--
			w.phaseActed(phPillage)
		}
	}
}

func (w *world) hideMageDecay() {
	w.phaseRan(phHideMage)
	for who := range w.loopKind(tChar) {
		if p := w.rpMagic(who); p != nil && p.hideMage > 0 {
			p.hideMage--
			w.phaseActed(phHideMage)
		}
	}
}
