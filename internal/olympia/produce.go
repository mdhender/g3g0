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

// COLLECT dispatch remains unsupported until the daily/reparse closure is
// ported. These complete start and summary owners are shared by that closure.
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

// produce.c:172–182. Do not infer this from current stock: terrain controls
// the distinction between exhausted stock and a resource never found here.
func itemGenHere(terr, item int) bool {
	for _, p := range terrProd {
		if p.terr == 0 {
			break
		}
		if p.terr == terr && p.item == item {
			return true
		}
	}
	return false
}

// Original struct and table, produce.c:315–506. The sentinel's message
// pointers really are NULL. visItem/mult are consumed by the daily owner.
type harvest struct {
	item, visItem, mult, skill, worker, chance int
	gotEm, noneNow, noneEver, taskDesc         *string
	public                                     bool
}

var harvTbl = [...]harvest{
	{itemPeasant, 0, 0, 0, 0, 0, new("recruited"), new("There are no more peasants here to recruit."), new("Peasants must be recruited in provinces."), new("recruit peasants"), true},
	{itemCorpse, 0, 0, skRaiseCorpses, 0, 0, new("raised"), new("There are no more corpses here to raise."), new("Corpses are found in graveyards."), new("raise corpses"), false},
	{itemMallornWood, 0, 0, skHarvestMallorn, 0, 20, new("cut"), new("All mallorn wood ready this month has been cut here."), new("Mallorn wood is found only in mallorn groves."), new("cut mallorn wood"), true},
	{itemOpium, 0, 0, skHarvestOpium, 0, 0, new("harvested"), new("All opium  ready this month has been harvested."), new("Opium is harvested only in poppy fields."), new("harvest opium"), true},
	{itemStone, 0, 0, skQuarryStone, itemWorker, 0, new("quarried"), new("No further stone may be quarried here this month."), new("Stone must be quarried in mountain provinces."), new("quarry stone"), true},
	{itemFish, 0, 0, skFishing, itemSailor, 50, new("caught"), new("No further fish may be caught here this month."), new("Fish must be caught in ocean provinces."), new("catch fish"), true},
	{itemLumber, 0, 0, skHarvestLumber, itemWorker, 0, new("cut"), new("All ready timber has already been cut this month."), new("Wood must be cut in forest provinces."), new("cut timber"), true},
	{itemYew, 0, 0, skHarvestYew, itemWorker, 0, new("cut"), new("All yew available this month has already been cut."), new("Yew must be cut in yew groves"), new("cut yew"), true},
	{itemWildHorse, 0, 0, skCatchHorse, 0, 50, new("caught"), new("No wild horses can be found roaming here now."), new("Wild horses are found on the plains and in pastures."), new("catch horses"), true},
	{itemAviniaLeaf, 0, 0, skCollectFoliage, 0, 20, new("collected"), new("All of the avinia leaves here have been collected."), new("Avinia leaves are found in mallorn groves."), new("collect avinia leaves"), true},
	{itemSpinyRoot, 0, 0, skCollectFoliage, 0, 25, new("collected"), new("All of the spiny roots here have been collected."), new("Avinia leaves are found in swamps, pits and bogs."), new("collect spiny roots"), true},
	{itemLanaBark, 0, 0, skCollectFoliage, 0, 50, new("collected"), new("All of the lana bark here has been collected."), new("Lana bark is found in circles of trees."), new("collect lana bark"), true},
	{itemFarrenstone, 0, 0, skCollectElem, 0, 100, new("collected"), new("This cave's supply of farrenstone for this month has been exhausted."), new("Farrenstone is found in caves."), new("collect farrenstone"), true},
	{itemPretusBones, 0, 0, skCollectElem, 0, 100, new("collected"), new("No pretus bones can be found."), new("Pretus bones are found in sand pits."), new("collect pretus bones"), true},
	{itemMageMenial, itemGold, 10, skMageMenial, 0, 100, new("earned"), new("No work at common magic can be found here."), new("No work at common magic can be found here."), new("work at common magic"), true},
	{},
}

// produce.c:509–519, including its pointer return-type line.
func findHarv(item int) *harvest {
	for i := range harvTbl {
		if harvTbl[i].item == 0 {
			break
		}
		if harvTbl[i].item == item {
			return &harvTbl[i]
		}
	}
	return nil
}

// produce.c:579–623. Skill refusal precedes all saved-argument mutation;
// empty-stock completion precedes even the worker requirement.
func (w *world) vGenericHarvest(c *command, number, days int, t *harvest) int {
	who := int(c.who)
	where := w.subloc(who)
	if t.item == itemFish && w.isShip(where) {
		where = w.loc(where)
	}
	if t.skill != 0 && w.hasSkill(who, t.skill) < 1 {
		w.commandOut(who, "Requires %s.", w.boxName(t.skill))
		return 0
	}
	if days < 1 {
		days = -1
	}
	c.c, c.d = int32(number), 0
	if w.hasItem(where, t.item) <= 0 {
		return w.iGenericHarvest(c, t)
	}
	if t.worker != 0 && w.hasItem(who, t.worker) < 1 {
		w.commandOut(who, "Need at least one %s to %s.", w.boxName(t.worker), *t.taskDesc)
		return 0
	}
	w.collectors.Append(who)
	c.wait = int32(days)
	return 1
}

// produce.c:701–722. EntMenial is keyed by actual entity and turn, not day.
func (w *world) mageMenialHow(who int) string {
	switch w.streams.EntMenial(who) {
	case 1:
		return " curing runny noses"
	case 2:
		return " dowsing for water"
	case 3:
		return " selling love potions"
	case 4:
		return " selling good luck charms"
	case 5:
		return " predicting the future"
	case 6:
		return " reading palms"
	case 7, 8, 9:
		return ""
	default:
		assert(false, "olympia/produce.c", 718, "FALSE")
	}
	return ""
}

// produce.c:725–780. This summary does not clear c.d or manufacture inventory.
// The caller owns exactly-once completion; experience has its real month guard.
func (w *world) iGenericHarvest(c *command, t *harvest) int {
	who := int(c.who)
	where := w.subloc(who)
	if t.item == itemFish && w.isShip(where) {
		where = w.loc(where)
	}
	if c.d == 0 {
		text := t.noneEver
		if itemGenHere(w.subkind(where), t.item) {
			text = t.noneNow
		}
		w.commandOut(who, "%s", *text)
	} else {
		item := t.item
		if t.visItem != 0 {
			item = t.visItem
		}
		if t.item == itemMageMenial {
			w.commandOut(who, "Earned %s%s.", goldS(int(c.d)), w.mageMenialHow(who))
		} else {
			w.commandOut(who, "%s %s.", capitalize(*t.gotEm), w.justNameQty(item, int(c.d)))
		}
		if t.public {
			w.showToGarrison = true
			if t.item == itemMageMenial {
				w.commandOut(where, "%s earned %s working at common magic.", w.boxName(who), goldS(int(c.d)))
			} else {
				w.commandOut(where, "%s %s %s.", w.boxName(who), *t.gotEm, w.justNameQty(item, int(c.d)))
			}
			w.showToGarrison = false
		}
		if t.skill != 0 {
			w.addSkillExperience(who, t.skill)
		}
	}
	w.collectors.RemValue(who)
	c.wait = 0
	if c.d > 0 && c.d >= c.c {
		return 1
	}
	return 0
}

// produce.c:783–801. Deliberately not installed in cmdTbl yet.
func (w *world) vCollect(c *command) int {
	t := findHarv(int(c.a))
	if t == nil {
		w.commandOut(int(c.who), "Don't know how to collect %s.", w.boxCode(int(c.a)))
		return 0
	}
	return w.vGenericHarvest(c, int(c.b), int(c.c), t)
}
