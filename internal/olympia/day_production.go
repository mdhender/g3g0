// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// compute_civ_levels is day.c:1792–1929. Each building type contributes only
// once per province, in all_here's depth-first order, not once per sublocation.
func (w *world) computeCivLevels() {
	crownLoc := -1
	if w.bx[relicCrown] != nil {
		crownLoc = w.province(w.itemUnique(relicCrown))
	}
	w.stage("compute_civ_levels()")
	w.clearTemps(tLoc)
	if w.validBox(crownLoc) {
		w.bx[crownLoc].temp += 8
	}
	prevLoc := -1
	var castle, tower, temple, inn, mine, city bool
	for where := range w.loopProvince() {
		if where != prevLoc {
			prevLoc = where
			castle, tower, temple, inn, mine, city = true, true, true, true, true, true
		}
		for _, i := range w.allHere(where) {
			var amount int32
			switch w.subkind(i) {
			case subCastle:
				if castle {
					amount = 6
					if p := w.rpSubloc(i); p != nil {
						amount += int32(p.castleLev)
					}
					castle = false
				}
			case subTower:
				if tower {
					amount, tower = 4, false
				}
			case subTemple:
				if temple {
					amount, temple = 4, false
				}
			case subInn:
				if inn {
					amount, inn = 4, false
				}
			case subMine:
				if mine {
					amount, mine = 4, false
				}
			case subCity:
				if city {
					amount, city = 4, false
					if p := w.rpSubloc(i); p != nil && p.major != 0 {
						amount = 8
					}
				}
			}
			if amount != 0 {
				w.bx[w.province(i)].temp += amount
			}
		}
	}
	for {
		changed := false
		for where := range w.loopProvince() {
			if w.subkind(where) == subOcean {
				continue
			}
			// loop_prov_dest reads the live list, and skips an absent entity_loc.
			p := w.rpLoc(where)
			if p == nil {
				continue
			}
			for _, i := range p.provDest {
				if i == 0 {
					continue
				}
				destCiv := w.bx[i].temp / 2
				if destCiv > w.bx[where].temp {
					w.bx[where].temp = destCiv
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}
	for where := range w.loopProvince() {
		w.pLoc(where).civ = int8(w.bx[where].temp / 4)
	}
}

// post_production is day.c:1932–1942. Keep the real order: #83 must implement
// location_trades before the later helpers can execute in this composition.
// loadDb's new-world initialization remains an explicit stage-5 boundary.
func (w *world) postProduction() {
	w.computeCivLevels()
	w.locationTrades()
	w.locationProduction()
	w.seedTaxes()
	w.postHasBeenRun = 1
}

func (w *world) locationTrades() { panic(notImplemented("location_trades")) }

// move_city_gold is day.c:2005–2024; use real move_item, including its zero
// quantity behavior and the city subtype chain (not a province traversal).
func (w *world) moveCityGold() {
	for i := range w.loopSubkind(subCity) {
		prov := w.province(i)
		has := w.hasItem(i, itemTaxCookie)
		w.moveItem(i, prov, itemTaxCookie, int(has))
	}
}
