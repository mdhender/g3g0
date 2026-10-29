// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Complete hinder-meditation callbacks from olympia/basic.c:331–420.
func (w *world) vHinderMed(c *command) int {
	target := int(c.a)
	who := int(c.who)
	if c.b < 1 {
		c.b = 1
	}
	if c.b > 3 {
		c.b = 3
	}
	if !w.checkAura(who, int(c.b)) {
		return 0
	}
	where := w.resetCastWhere(who)
	c.d = int32(where)
	if !w.checkCharWhere(where, who, target) {
		return 0
	}
	w.commandOut(who, "Attempt to hinder attempts at meditation by %s.", w.boxCode(target))
	return 1
}

func (w *world) hinderMedOmen(who, other int) {
	switch w.streams.MagcOmen(who, other) {
	case 1:
		w.commandOut(who, "A disturbing image of %s appeared last night in a dream.", w.boxName(other))
	case 2:
		w.commandOut(who, "As a cloud drifts across the moon, it seems for an instant that it takes the shape of a ghoulish face, looking straight at you.")
	case 3:
		w.commandOut(who, "You are shocked out of your slumber in the middle of the night by cold fingers touching your neck, but when you glance about, there is no one to be seen.")
	case 4:
	default:
		assert(false, "olympia/basic.c", 388, "FALSE")
	}
}

func (w *world) dHinderMed(c *command) int {
	target, aura, where := int(c.a), int(c.b), int(c.d)
	who := int(c.who)
	if !w.chargeAura(who, aura) {
		return 0
	}
	if !w.checkCharWhere(where, who, target) {
		return 0
	}
	w.commandOut(who, "Successfully cast %s on %s.", w.boxName(skHinderMed), w.boxName(target))
	p := w.pMagic(target)
	// C stores the signed-char addition before comparing against the cap.
	p.hinderMeditation += int8(aura)
	if p.hinderMeditation > 5 {
		p.hinderMeditation = 5
	}
	w.hinderMedOmen(target, who)
	return 1
}
