// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Complete day.c:237–486 at the stage pin.
func (w *world) nearRockyCoast(where int) int {
	if w.subkind(where) != subOcean {
		return 0
	}
	ret := 3
	for _, v := range w.exitsFromLocNSEW(0, where) {
		if w.subkind(int(v.destination)) != subOcean {
			if w.subkind(int(v.destination)) == subMountain {
				return 2
			}
			ret = 1
		}
	}
	return ret
}

func (w *world) shipCoastalDamage() {
	w.phaseRan(phCoastDamage)
	for ship := range w.loopKind(tShip) {
		if !w.isShip(ship) {
			continue
		}
		n := w.nearRockyCoast(w.subloc(ship))
		sub := int(w.sysclock.day) * 16
		switch n {
		case 0:
		case 1:
			if w.streams.WthrWreck(ship, sub, 1, 75) == 1 {
				w.commandOut(ship, "%s struck a coastal reef.  There is minor damage to the ship.", w.boxName(ship))
				w.addStructureDamage(ship, w.streams.WthrWreck(ship, sub+1, 3, 5), true)
				w.phaseActed(phCoastDamage)
			}
		case 2:
			if w.streams.WthrWreck(ship, sub+2, 1, 50) == 1 {
				w.commandOut(ship, "%s struck some submerged rocks.  There is minor damage to the ship.", w.boxName(ship))
				w.addStructureDamage(ship, w.streams.WthrWreck(ship, sub+3, 6, 9), true)
				w.phaseActed(phCoastDamage)
			}
		case 3:
			switch w.streams.WthrWreck(ship, sub+4, 1, 200) {
			case 1, 2:
				w.commandOut(ship, "Hungry looking birds circle overhead.")
				w.phaseActed(phCoastOmen)
			case 3:
				w.commandOut(ship, "Sharks circle in the water a short distance from the ship.")
				w.phaseActed(phCoastOmen)
			}
		default:
			assert(false, "olympia/day.c", 323, "FALSE")
		}
	}
}

func (w *world) mineCalamity(mine int) {
	sub := int(w.sysclock.day) * 16
	what := w.streams.WthrWreck(mine, sub, 1, 4)
	w.vectorCharHere(mine)
	w.outVector.Append(mine)
	switch what {
	case 1:
		w.commandOut(vect, "A tunnel in %s has caved in!", w.boxName(mine))
	case 2:
		w.commandOut(vect, "A roof has collapsed in a shaft in %s!", w.boxName(mine))
	case 3:
		w.commandOut(vect, "A new fissure in a shaft wall has allowed poisonous gas to seep into %s.", w.boxName(mine))
	case 4:
		w.commandOut(vect, "Rock dust has caused an explosion in %s!", w.boxName(mine))
	default:
		assert(false, "olympia/day.c", 367, "FALSE")
	}
	toKill := min(w.streams.WthrWreck(mine, sub+1, 1, 5), w.countLocCharItem(mine, itemWorker))
	if toKill > 0 {
		verb := "were"
		if toKill == 1 {
			verb = "was"
		}
		w.commandOut(vect, "%s miner%s %s killed in the accident.", capitalize(niceNum(toKill)), addS(toKill), verb)
		for _, i := range w.allCharHere(mine) {
			if m := w.hasItem(i, itemWorker); m != 0 {
				n := min(int(m), toKill)
				w.commandOut(i, "%s lost %s worker%s.", w.boxName(i), niceNum(n), addS(n))
				w.consumeItem(i, itemWorker, n)
				toKill -= n
				if toKill <= 0 {
					break
				}
			}
		}
	}
	w.addStructureDamage(mine, w.streams.WthrWreck(mine, sub+2, 1, 15), true)
}

func (w *world) innCalamity(where int) {
	sub := int(w.sysclock.day) * 16
	own := w.buildingOwner(where)
	dam := w.streams.WthrWreck(where, sub, 5, 15)
	subject := [...]string{"Some customers ", "Some patrons ", "An irate customer ", "Two large, angry men ", "A surly local ", "A party of traveling entertainers "}
	action := [...]string{"got drunk, ", "started a fight, ", "got drunk and started a fight, ", "insulted the chef, ", "refused to pay, "}
	consequence := [...]string{"and broke some furniture", "and damaged a wall", "and kicked in the door", "and knocked over a keg of beer", "and set a fire in the closet", "and knocked over the smokehouse", "and broke some chairs"}
	buf := subject[w.streams.WthrWreck(where, sub+1, 1, 6)-1]
	buf += action[w.streams.WthrWreck(where, sub+2, 1, 5)-1]
	buf += consequence[w.streams.WthrWreck(where, sub+3, 1, 7)-1]
	if own != 0 {
		w.commandOut(own, "%s:  %s, causing %d points of damage.", w.boxName(where), buf, dam)
	}
	w.addStructureDamage(where, dam, true)
}

func (w *world) randomLocDamage() {
	w.phaseRan(phLocMine)
	for where := range w.loopKind(tLoc) {
		switch w.subkind(where) {
		case subMine:
			depth := w.mineDepth(where)
			if w.streams.WthrWreck(where, int(w.sysclock.day)*16+8, 1, 90) <= depth {
				w.mineCalamity(where)
				w.phaseActed(phLocMine)
			}
		case subInn:
			if w.streams.WthrWreck(where, int(w.sysclock.day)*16+9, 1, 100) == 1 {
				w.innCalamity(where)
				w.phaseActed(phLocInn)
			}
		}
	}
}

// mine_depth is oly.h:1102, with signed-short shaft depth and C truncation.
func (w *world) mineDepth(where int) int {
	if p := w.rpSubloc(where); p != nil {
		return int(p.shaftDepth) / 3
	}
	return 0
}

// count_loc_char_item is u.c:1336–1349. all_char_here takes a depth-first
// snapshot, descending through characters only, just as loop_char_here does.
func (w *world) countLocCharItem(where, item int) int {
	var sum int32
	for _, i := range w.allCharHere(where) {
		sum += w.hasItem(i, item)
	}
	return int(sum)
}
