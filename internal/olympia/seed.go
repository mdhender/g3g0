// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Existing-world seed_taxes and int_comp are ported. New-world creation
// remains an explicit stage-5 boundary.

// intComp is int_comp(), the qsort() comparison of ints: a - b, as a C int,
// which wraps if the two are more than INT_MAX apart. The lists it sorts
// hold box numbers, which never are.
func intComp(a, b int) int {
	return int(int32(a) - int32(b))
}

// seed_taxes is seed.c:641–704. Arithmetic is C int, including division
// before opium multiplication and the later pillage division.
func (w *world) seedTaxes() {
	for where := range w.loopKind(tLoc) {
		if w.locDepth(where) != locProvince && w.subkind(where) != subCity {
			continue
		}
		if w.subkind(where) == subOcean {
			continue
		}
		if w.subkind(where) == subCity {
			w.consumeItem(where, itemPettyThief, int(w.hasItem(where, itemPettyThief)))
			w.genItem(where, itemPettyThief, 1)
		}
		w.consumeItem(where, itemMageMenial, int(w.hasItem(where, itemMageMenial)))
		w.consumeItem(where, itemTaxCookie, int(w.hasItem(where, itemTaxCookie)))
		var civ, pil, opium int32
		if p := w.rpLoc(w.province(where)); p != nil {
			civ = int32(p.civ)
		}
		if p := w.rpSubloc(where); p != nil {
			pil, opium = int32(p.loot), p.opiumEcon
		}
		base := civ * 5
		if pil != 0 {
			base /= pil + 1
		}
		w.genItem(where, itemMageMenial, int(base))
		assert(w.hasItem(where, itemTaxCookie) == 0, "olympia/seed.c", 679, "has_item(where, item_tax_cookie) == 0")
		if w.subkind(where) == subCity {
			base = 100
		} else {
			base = 50 + civ*50
		}
		base -= base / 10 * opium
		assert(base > 0, "olympia/seed.c", 696, "base > 0")
		if pil != 0 {
			base /= pil + 1
		}
		w.genItem(where, itemTaxCookie, int(base))
	}
}
