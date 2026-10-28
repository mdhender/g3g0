// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "github.com/mdhender/g3g0/internal/lists"

// provinceSubloc is build.c:160–177. Direct province children take
// precedence over children of the first city; it does not recurse further.
func (w *world) provinceSubloc(where, sk int) int {
	prov := w.province(where)
	city := w.sublocHere(prov, subCity)
	if n := w.sublocHere(prov, sk); n != 0 {
		return n
	}
	if city != 0 {
		return w.sublocHere(city, sk)
	}
	return 0
}

// countLocStructures is the complete loc.c:352–366 definition. Only direct
// location children count, including unowned structures and either subtype.
func (w *world) countLocStructures(where, a, b int) int32 {
	var sum int32
	assert(w.validBox(where), "olympia/loc.c", 358, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.kind(i) == tLoc && (w.subkind(i) == a || w.subkind(i) == b) {
			sum++
		}
	}
	return sum
}

// Complete loc.c:332–351, also used by the city_here macro (loc.h:22).
func (w *world) sublocHere(where, sk int) int {
	assert(w.validBox(where), "olympia/loc.c", 338, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.kind(i) == tLoc && w.subkind(i) == sk {
			return i
		}
	}
	return 0
}

// allHere is all_here(): all descendants, including noncharacters. Unlike
// all_stack this does not prune branches at noncharacter boxes.
func (w *world) allHere(who int) lists.IList {
	assert(w.validBox(who), "olympia/loc.c", 53, "valid_box(who)")
	var l lists.IList
	for _, i := range w.rpLocInfo(who).hereList {
		w.addHere(i, &l)
	}
	return l
}

func (w *world) addHere(who int, l *lists.IList) {
	assert(w.validBox(who), "olympia/loc.c", 35, "valid_box(who)")
	l.Append(who)
	for _, i := range w.rpLocInfo(who).hereList {
		w.addHere(i, l)
	}
}

func (w *world) allCharHere(who int) lists.IList {
	assert(w.validBox(who), "olympia/loc.c", 92, "valid_box(who)")
	var l lists.IList
	for _, i := range w.rpLocInfo(who).hereList {
		if w.kind(i) == tChar {
			w.addCharHere(i, &l)
		}
	}
	return l
}

// Ported from olympia/loc.c: the here lists, all_stack(), set_where() and
// the location walks that check_db() reaches (subloc() and viewloc(),
// through set_lord() and touch_loc()) and determine_map_edges() reaches
// (region()). The rest of loc.c is left for later stages.
//
// rp_loc_info() is the address of a member of the box, never NULL, so the
// "p == NULL" tests of all_stack() and in_here_list(), and the assert(p !=
// NULL) of add_char_here(), cannot fire; they are left out. A box number
// with no box is a NULL dereference in C and a panic here.

// addCharHere is add_char_here(): it appends who to l, then, depth first,
// each unit in who's here list.
func (w *world) addCharHere(who int, l *lists.IList) {
	assert(w.validBox(who), "olympia/loc.c", 73, "valid_box(who)")

	l.Append(who)

	p := w.rpLocInfo(who)

	for i := 0; i < p.hereList.Len(); i++ {
		if w.kind(p.hereList[i]) == tChar {
			w.addCharHere(p.hereList[i], l)
		}
	}
}

// allStack is all_stack(): who, then every unit stacked under it, depth
// first in here-list order.
func (w *world) allStack(who int) lists.IList {
	assert(w.validBox(who), "olympia/loc.c", 113, "valid_box(who)")

	l := lists.IList{who}

	p := w.rpLocInfo(who)

	for i := 0; i < p.hereList.Len(); i++ {
		if w.kind(p.hereList[i]) == tChar {
			w.addCharHere(p.hereList[i], &l)
		}
	}

	return l
}

// region is region(): who, or the first location out from it, that is a
// region, or 0.
func (w *world) region(who int) int {
	for who > 0 &&
		(w.kind(who) != tLoc || w.locDepth(who) != locRegion) {
		who = w.loc(who)
	}

	return who
}

// subloc is subloc(): the location or ship who is in, irrespective of who
// it is stacked with, or 0 if there is none.
func (w *world) subloc(who int) int {
	for {
		who = w.loc(who)
		if !(who > 0 && w.kind(who) != tLoc && w.kind(who) != tShip) {
			return who
		}
	}
}

// viewloc is viewloc(): who, or the first location out from it, that is a
// province, a city, a graveyard, a sewer or a faery hill, or 0.
func (w *world) viewloc(who int) int {
	for who > 0 &&
		w.locDepth(who) != locProvince &&
		w.subkind(who) != subCity &&
		w.subkind(who) != subGraveyard &&
		w.subkind(who) != subSewer &&
		w.subkind(who) != subFaeryHill {
		who = w.loc(who)
	}

	return who
}

// addToHereList is add_to_here_list(): it appends who to the here list of
// loc, which must not hold it.
func (w *world) addToHereList(loc, who int) {
	assert(!w.inHereList(loc, who), "olympia/loc.c", 231, "!in_here_list(loc, who)")
	w.pLocInfo(loc).hereList.Append(who)
	assert(w.inHereList(loc, who), "olympia/loc.c", 233, "in_here_list(loc, who)")
}

// removeFromHereList is remove_from_here_list(): it removes who, every
// time it appears, from the here list of loc, which must hold it.
func (w *world) removeFromHereList(loc, who int) {
	assert(w.inHereList(loc, who), "olympia/loc.c", 241, "in_here_list(loc, who)")
	w.rpLocInfo(loc).hereList.RemValue(who)
	assert(!w.inHereList(loc, who), "olympia/loc.c", 243, "!in_here_list(loc, who)")
}

// setWhere is set_where(): it moves who from the here list of the
// location it is in, if any, to that of newLoc, if not 0, and makes newLoc
// its location.
func (w *world) setWhere(who, newLoc int) {
	// This check could be expanded to make sure that new_loc is not
	// anywhere inside of who, by walking up new_loc to the top and making
	// sure we don't go through who.
	assert(who != newLoc, "olympia/loc.c", 258, "who != new_loc")

	oldLoc := w.loc(who)

	if oldLoc > 0 {
		w.removeFromHereList(oldLoc, who)
	}

	if newLoc > 0 {
		w.addToHereList(newLoc, who)
	}

	w.pLocInfo(who).where = int32(newLoc)
}

// markLocStackKnown is loc.c:281–295. The character root is included even
// when a prisoner; captive descendants are excluded, not their children.
func (w *world) markLocStackKnown(stack, where int) {
	if w.kind(stack) == tChar {
		w.setKnown(stack, where)
	}
	for _, i := range w.allCharHere(stack) {
		if !w.isPrisoner(i) {
			w.setKnown(i, where)
		}
	}
}

// inHereList is in_here_list(): whether who is in the here list of loc.
func (w *world) inHereList(loc, who int) bool {
	return w.rpLocInfo(loc).hereList.Lookup(who) != -1
}

// firstCharacter is first_character(): the first unit in the here list of
// where, or 0.
func (w *world) firstCharacter(where int) int {
	assert(w.validBox(where), "olympia/loc.c", 318, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.kind(i) == tChar {
			return i
		}
	}

	return 0
}

// inSafeNow checks who itself, then every ancestor, even while moving.
func (w *world) inSafeNow(who int) bool {
	for {
		if w.safeHaven(who) != 0 {
			return true
		}
		who = w.loc(who)
		if who <= 0 {
			return false
		}
	}
}

func (w *world) province(who int) int {
	for who > 0 && (w.kind(who) != tLoc || w.locDepth(who) != locProvince) {
		who = w.loc(who)
	}
	return who
}

func (w *world) buildingOwner(where int) int {
	assert(w.locDepth(where) == locBuild, "olympia/loc.c", 395, "loc_depth(where) == LOC_build")
	return w.firstCharacter(where)
}
