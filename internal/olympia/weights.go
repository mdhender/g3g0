// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// weights is u.h's struct weights. All arithmetic is C int arithmetic,
// including item products and intermediate travel ratios.
type weights struct {
	animals, totalWeight, landCap, landWeight int32
	rideCap, rideWeight, flyCap, flyWeight    int32
}

func (w *world) addItemWeight(item int, qty int32, v *weights) {
	var wt, lc, rc, fc int32
	if p := w.rpItem(item); p != nil {
		wt = int32(p.weight) * qty
		lc, rc, fc = int32(p.landCap), int32(p.rideCap), int32(p.flyCap)
		if p.animal != 0 {
			v.animals += qty
		}
	}
	if lc != 0 {
		v.landCap += max(lc, 0) * qty
	} else {
		v.landWeight += wt
	}
	if rc != 0 {
		v.rideCap += max(rc, 0) * qty
	} else {
		v.rideWeight += wt
	}
	if fc != 0 {
		v.flyCap += max(fc, 0) * qty
	} else {
		v.flyWeight += wt
	}
	v.totalWeight += wt
}

func (w *world) determineUnitWeights(who int) weights {
	assert(w.kind(who) == tChar, "olympia/u.c", 1480, "kind(who) == T_char")
	var v weights
	base := w.nobleItem(who)
	if base == 0 {
		base = itemPeasant
	}
	w.addItemWeight(base, 1, &v)
	assert(w.validBox(who), "olympia/u.c", 1490, "valid_box(who)")
	for e := range w.loopInv(who) {
		w.addItemWeight(int(e.item), e.qty, &v)
	}
	return v
}

func (w *world) determineStackWeights(who int) weights {
	v := w.determineUnitWeights(who)
	for _, i := range w.allHere(who) {
		u := w.determineUnitWeights(i)
		v.totalWeight += u.totalWeight
		v.landCap += u.landCap
		v.rideCap += u.rideCap
		v.flyCap += u.flyCap
		v.landWeight += u.landWeight
		v.rideWeight += u.rideWeight
		v.flyWeight += u.flyWeight
		v.animals += u.animals
	}
	return v
}

func (w *world) countStackMoveNobles(who int) int32 {
	sum := int32(1)
	pl := w.player(who)
	for _, i := range w.allCharHere(who) {
		if w.player(i) == pl || (!w.isNPC(i) && !w.isPrisoner(i)) {
			sum++
		}
	}
	return sum
}

func (w *world) countStackFigures(who int) int32 {
	var sum int32
	for i := range w.loopStack(who) {
		sum += w.countAny(i)
	}
	return sum
}
