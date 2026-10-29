// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// weeklyPrisonerEscapeCheck is the complete stack.c:345–375. The outer
// character chain is live; each direct here-list is copied before escapes.
func (w *world) weeklyPrisonerEscapeCheck() {
	w.phaseRan(phPrisonerEscape)
	for who := range w.loopKind(tChar) {
		if w.isPrisoner(who) || w.subkind(w.subloc(who)) == subOcean {
			continue
		}
		assert(w.validBox(who), "olympia/stack.c", 361, "valid_box(who)")
		for i := range w.loopHere(who) {
			if w.kind(i) == tChar && w.isPrisoner(i) {
				if p := w.rpMagic(i); p != nil && p.swearOnRelease != 0 {
					continue
				}
				chance := 2
				if w.locDepth(w.subloc(who)) >= locBuild {
					chance = 1
				}
				if w.checkPrisonerEscape(i, chance, 0) { // stack.h ESC_WEEKLY
					w.phaseActed(phPrisonerEscape)
				}
			}
		}
	}
}

// dogsBarkAtHiddenChars is day.c:489–516. Only directly present, really
// hidden characters trigger the location; nested characters alone do not.
func (w *world) dogsBarkAtHiddenChars() {
	w.phaseRan(phDogsBark)
	for where := range w.loopKind(tLoc) {
		flag := false
		assert(w.validBox(where), "olympia/day.c", 502, "valid_box(where)")
		for i := range w.loopHere(where) {
			if w.kind(i) == tChar && w.charReallyHidden(i) {
				flag = true
			}
		}
		if flag {
			w.barkDogs(where)
			w.phaseActed(phDogsBark)
		}
	}
}

// barkDogs is u.c:2634–2666, retaining VECT recipients, stack-owner item
// sums, direct-character keys and the per-hound enty leaf draw order.
func (w *world) barkDogs(where int) {
	w.vectorCharHere(where)
	w.outVector.Append(where) // vector_add appends unconditionally.
	var sum, bark int32
	assert(w.validBox(where), "olympia/u.c", 2645, "valid_box(where)")
	for who := range w.loopHere(where) {
		if w.kind(who) != tChar {
			continue
		}
		n := int32(w.stackHasItem(who, itemHound))
		sum += n
		for i := int32(1); i <= n; i++ {
			if w.streams.EntBark(who, int(i)) == 1 {
				bark++
			}
		}
	}
	if bark == 1 && sum == 1 {
		w.commandOut(vect, "The hound is barking.")
	} else if bark == 1 {
		w.commandOut(vect, "A hound is barking.")
	} else if bark > 1 {
		w.commandOut(vect, "The hounds are barking.")
	}
}

// incrementStoneRingAura is day.c:777–803, using the existing effective
// maximum contract (auraculum and per-entry artifact bonuses included).
func (w *world) incrementStoneRingAura() {
	w.phaseRan(phRingAura)
	for who := range w.loopKind(tChar) {
		if p := w.rpMagic(who); p == nil || p.magician == 0 {
			continue
		}
		if w.subkind(w.subloc(who)) != subStoneCir {
			continue
		}
		ma := w.maxEffAura(who)
		p := w.pMagic(who)
		if p.curAura < ma {
			p.curAura++
			w.commandOut(who, "Current aura is now %d.", p.curAura)
			w.phaseActed(phRingAura)
		}
	}
}

// noncreatorCurseErode is the active day.c:901–963 path, NOT its disabled
// loyalty implementation. Erosion still only logs NOTYET after G3 #357;
// the unported USE 886 validator now refuses the curse (stage 4).
func (w *world) noncreatorCurseErode() {
	w.phaseRan(phCurseErode)
	for who := range w.loopKind(tChar) {
		assert(w.validBox(who), "olympia/day.c", 914, "valid_box(who)")
		for e := range w.loopInv(who) {
			im := w.rpItemMagic(int(e.item))
			if im == nil || im.curseLoyalty == 0 || int(im.creator) == who {
				continue
			}
			if p := w.rpChar(who); p != nil && p.loyKind == loyOath {
				continue
			}
			w.logWrite(10, "noncreator_curse_erode: NOTYET!")
			w.phaseActed(phCurseErode)
		}
	}
}
