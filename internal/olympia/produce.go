// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Original produce.c:10–88 terrain and depth tables. Production replenishes
// a stock floor; it does not add the table quantity each month.
const poppyOpium = 25

var terrProd = [...]struct{ terr, item, qty int }{
	{subForest, itemLumber, 30},
	{subSacredGrove, itemLumber, 5},
	{subTreeCircle, itemLumber, 5},
	{subMountain, itemStone, 50},
	{subRockyHill, itemStone, 50},
	{subDesert, itemStone, 10},
	{subCave, itemFarrenstone, 2},
	{subPlain, itemWildHorse, 5},
	{subPasture, itemWildHorse, 5},
	{subOcean, itemFish, 50},
	{subMallornGrove, itemAviniaLeaf, 2},
	{subMallornGrove, itemMallornWood, 2},
	{subBog, itemSpinyRoot, 4},
	{subPits, itemSpinyRoot, 4},
	{subSwamp, itemSpinyRoot, 1},
	{subYewGrove, itemYew, 5},
	{subGraveyard, itemCorpse, 15},
	{subTreeCircle, itemLanaBark, 3},
	{subSandPit, itemPretusBones, 1},
	{subSwamp, itemOpium, poppyOpium},
	{subPoppyField, itemOpium, poppyOpium},
	{subForest, itemPeasant, 10},
	{subMountain, itemPeasant, 10},
	{subPlain, itemPeasant, 10},
	{subCity, itemPeasant, 10},
	{0, 0, 0},
}

var mineProd = [...]struct{ iron, gold, mithril int }{
	{0, 0, 0}, {15, 25, 0}, {12, 100, 0}, {10, 200, 0},
	{8, 500, 1}, {5, 500, 0}, {3, 400, 0}, {0, 150, 0},
	{0, 50, 1}, {0, 0, 2}, {0, 0, 8}, {0, 0, 5},
	{0, 0, 0}, {10, 10, 0}, {10, 50, 2}, {0, 10, 1},
	{0, 0, 0}, {10, 0, 1}, {0, 500, 0}, {0, 0, 0}, {0, 0, 0},
}

// replenish, mine_production and location_production: produce.c:91–169.
func (w *world) replenish(where, item, qty int) {
	if n := w.hasItem(where, item); n < int32(qty) {
		w.genItem(where, item, int(int32(qty)-n))
	}
}

func (w *world) mineProduction(where int) {
	depth := w.mineDepth(where)
	assert(depth > 0, "olympia/produce.c", 109, "depth > 0")
	if depth > 20 {
		depth = 20
	}
	p := mineProd[depth]
	w.replenish(where, itemIron, p.iron)
	w.replenish(where, itemGold, p.gold)
	w.replenish(where, itemMithril, p.mithril)
}

func (w *world) locationProduction() {
	for where := range w.loopKind(tLoc) {
		terr := w.subkind(where)
		if terr == subMine {
			w.mineProduction(where)
		} else {
			for _, p := range terrProd {
				if p.terr == 0 {
					break
				}
				if p.terr == terr {
					w.replenish(where, p.item, p.qty)
				}
			}
		}
		if terr == subPoppyField {
			if n := w.hasItem(where, itemOpium); n > poppyOpium {
				w.consumeItem(where, itemOpium, int(n-poppyOpium))
			}
			if p := w.rpMisc(where); p != nil && p.opiumDouble != 0 {
				p.opiumDouble = 0
				w.genItem(where, itemOpium, int(w.hasItem(where, itemOpium)))
			}
		}
		if terr == subIsland || (w.locDepth(where) == locProvince && w.hasOceanAccess(where) != 0) {
			w.replenish(where, itemFlotsam, 30)
		}
	}
}

// Harvesting and COLLECT's callbacks remain explicit unported boundaries.
func (w *world) initCollectList() {
	cmdCollect := w.findCommand("collect")
	assert(cmdCollect > 0, "olympia/produce.c", 533, "cmd_collect > 0")
	for i := range w.loopKind(tChar) {
		c := w.rpCommand(i)
		if c != nil && c.state == stateRun && int(c.cmd) == cmdCollect {
			w.collectors.Append(i)
		}
	}
}
