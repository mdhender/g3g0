// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

const whereShip = -1

// make.c:34–173. No original recipe uses worker or public one-day
// production, but the generic helpers retain those contracts.
type makeRecipe struct {
	item, inp1, inp2, reqSkill, worker int
	gotEm                              *string
	public                             bool
	where, aura, factor, days          int
}

var makeTbl = [...]makeRecipe{
	{itemBlankScroll, itemLanaBark, 0, skAlchemy, 0, strSave("made"), false, 0, 0, 1, 1},
	{itemEliteArch, itemArcher, 0, skArchery, 0, strSave("trained"), false, subCastle, 0, 1, 1},
	{itemAngryPeasant, itemPeasant, 0, skTrainAngry, 0, strSave("trained"), false, 0, 0, 1, 1},
	{itemPeasant, itemAngryPeasant, 0, skTrainAngry, 0, strSave("trained"), false, 0, 0, 1, 1},
	{itemArcher, itemSoldier, itemLongbow, skArchery, 0, strSave("trained"), false, 0, 0, 1, 1},
	{itemEliteGuard, itemKnight, itemPlate, skSwordplay, 0, strSave("trained"), false, subCastle, 0, 1, 1},
	{itemKnight, itemSwordsman, itemWarmount, skSwordplay, 0, strSave("trained"), false, 0, 0, 1, 1},
	{itemBlessedSoldier, itemSoldier, 0, skReligion, 0, strSave("trained"), false, subTemple, 0, 1, 1},
	{itemGhostWarrior, 0, 0, skSummonGhost, 0, strSave("summoned"), false, 0, 1, 2, 1},
	{itemSwordsman, itemSoldier, itemLongsword, skSwordplay, 0, strSave("trained"), false, 0, 0, 1, 1},
	{itemPirate, itemSailor, itemLongsword, skSwordplay, 0, strSave("trained"), false, whereShip, 0, 1, 1},
	{itemPikeman, itemSoldier, itemPike, skCombat, 0, strSave("trained"), false, 0, 0, 1, 1},
	{itemSoldier, itemPeasant, 0, skCombat, 0, strSave("trained"), false, 0, 0, 1, 1},
	{itemCrossbowman, itemPeasant, itemCrossbow, skCombat, 0, strSave("trained"), false, 0, 0, 1, 1},
	{itemSailor, itemPeasant, 0, skPilotShip, 0, strSave("trained"), false, 0, 0, 1, 1},
	{itemWorker, itemPeasant, 0, 0, 0, strSave("trained"), false, 0, 0, 1, 1},
	{itemBasket, 0, 0, 0, 0, strSave("made"), false, 0, 0, 1, 1},
	{itemPot, 0, 0, 0, 0, strSave("made"), false, 0, 0, 1, 1},
	{itemCrossbow, itemLumber, 0, skWeaponsmith, 0, strSave("made"), false, 0, 0, 1, 1},
	{itemPike, itemLumber, 0, skWeaponsmith, 0, strSave("made"), false, 0, 0, 1, 1},
	{itemLongsword, itemIron, 0, skWeaponsmith, 0, strSave("made"), false, 0, 0, 1, 1},
	{itemPlate, itemIron, 0, skWeaponsmith, 0, strSave("made"), false, 0, 0, 1, 1},
	{itemLongbow, itemYew, 0, skWeaponsmith, 0, strSave("made"), false, 0, 0, 1, 1},
	{itemDrum, itemMallornWood, 0, skSummonSavage, 0, strSave("made"), false, 0, 0, 1, 1},
	{itemHide, itemOx, 0, 0, 0, strSave("made"), false, 0, 0, 1, 1},
	{itemRidingHorse, itemWildHorse, 0, skTrainWild, 0, strSave("trained"), true, 0, 0, 1, 3},
	{itemWarmount, itemWildHorse, 0, skTrainWarmount, 0, strSave("trained"), true, 0, 0, 1, 7},
	{0, 0, 0, 0, 0, nil, false, 0, 0, 1, 1},
}

// make.c:176–186.
func findMake(item int) *makeRecipe {
	for i := range makeTbl {
		if makeTbl[i].item == 0 {
			break
		}
		if makeTbl[i].item == item {
			return &makeTbl[i]
		}
	}
	return nil
}

// make.c:193–256. Set desired/obtained before any refusal, but preserve wait
// until validation succeeds. All arithmetic on saved command values is C32.
func (w *world) vGenericMake(c *command, number int, t *makeRecipe) int {
	who := int(c.who)
	where := w.subloc(who)
	days := int32(-1)
	if days < 1 && int32(number) == 0 && t.inp1 == 0 && t.inp2 == 0 {
		days = monthDays + 1 - int32(w.sysclock.day)
	}
	c.c, c.d = int32(number), 0
	if t.reqSkill != 0 && w.hasSkill(who, t.reqSkill) < 1 {
		w.commandOut(who, "Requires %s.", w.boxName(t.reqSkill))
		return 0
	}
	if t.worker != 0 && w.hasItem(who, t.worker) < 1 {
		w.commandOut(who, "Need at least one %s.", w.boxName(t.worker))
		return 0
	}
	if t.inp1 != 0 && w.hasItem(who, t.inp1) < 1 {
		w.commandOut(who, "Don't have any %s.", w.pluralItemBox(t.inp1, 2))
		return 0
	}
	if t.inp2 != 0 && w.hasItem(who, t.inp2) < 1 {
		w.commandOut(who, "Don't have any %s.", w.pluralItemBox(t.inp2, 2))
		return 0
	}
	if t.where == whereShip && !w.isShip(where) && !w.isShipNotdone(where) {
		w.commandOut(who, "Must be on a ship.")
		return 0
	}
	if t.where > 0 && w.subkind(where) != t.where {
		w.commandOut(who, "Must be in a %s.", subkindS[t.where])
		return 0
	}
	if t.aura != 0 && w.charCurAura(who) < t.aura {
		w.commandOut(who, "Need at least %d aura.", t.aura)
		return 0
	}
	c.wait = days
	return 1
}

// make.c:259–326, including the original active #if 1 continuation branch.
func (w *world) dGenericMake(c *command, t *makeRecipe) int {
	who, number := int(c.who), c.c
	qty := int32(1)
	if t.worker != 0 {
		qty = w.hasItem(who, t.worker)
	}
	if t.inp1 != 0 {
		qty = min(qty, w.hasItem(who, t.inp1))
	}
	if t.inp2 != 0 {
		qty = min(qty, w.hasItem(who, t.inp2))
	}
	if t.aura != 0 {
		qty = min(qty, int32(w.charCurAura(who)))
	}
	if qty > 0 {
		if number > 0 && c.d+qty > number {
			qty = number - c.d
		}
		assert(qty >= 0, "olympia/make.c", 284, "qty >= 0")
		if t.inp1 != 0 {
			w.consumeItem(who, t.inp1, int(qty))
		}
		if t.inp2 != 0 {
			w.consumeItem(who, t.inp2, int(qty))
		}
		if t.aura != 0 {
			w.deductAura(who, t.aura)
		}
		w.genItem(who, t.item, int(qty*int32(t.factor)))
		c.d += qty
		if t.reqSkill != 0 {
			w.addSkillExperience(who, t.reqSkill)
		}
		if (t.inp1 == 0 || w.hasItem(who, t.inp1) > 0) &&
			(t.inp2 == 0 || w.hasItem(who, t.inp2) > 0) &&
			c.wait != 0 && !(number > 0 && c.d >= number) {
			return 1
		}
	}
	return w.iGenericMake(c, t)
}

// make.c:329–350. Public summary deliberately does not multiply by factor.
func (w *world) iGenericMake(c *command, t *makeRecipe) int {
	who := int(c.who)
	where := w.subloc(who)
	w.commandOut(who, "%s %s.", capitalize(*t.gotEm), w.justNameQty(t.item, int(c.d*int32(t.factor))))
	if t.public {
		w.commandOut(where, "%s %s %s.", w.boxName(who), *t.gotEm, w.justNameQty(t.item, int(c.d)))
	}
	c.wait = 0
	if c.d > 0 && c.d >= c.c {
		return 1
	}
	return 0
}

// make.c:357–393. Only the initial delay receives the experience speedup.
func (w *world) vSecondMake(c *command, number int, t *makeRecipe) int {
	who := int(c.who)
	c.c, c.d = int32(number), 0
	if t.reqSkill != 0 && w.hasSkill(who, t.reqSkill) < 1 {
		w.commandOut(who, "Requires %s.", w.boxName(t.reqSkill))
		return 0
	}
	if t.inp1 != 0 && w.hasItem(who, t.inp1) < 1 {
		w.commandOut(who, "Don't have any %s.", w.pluralItemBox(t.inp1, 2))
		return 0
	}
	if t.inp2 != 0 && w.hasItem(who, t.inp2) < 1 {
		w.commandOut(who, "Don't have any %s.", w.pluralItemBox(t.inp2, 2))
		return 0
	}
	c.wait, c.poll = int32(t.days), 0
	if t.reqSkill != 0 {
		c.useExp = int32(w.hasSkill(who, t.reqSkill))
		experienceUseSpeedup(c)
	}
	return 1
}

// make.c:396–442. Revalidate inputs only, consuming and producing one each.
func (w *world) dSecondMake(c *command, t *makeRecipe) int {
	who := int(c.who)
	if t.inp1 != 0 && w.hasItem(who, t.inp1) < 1 {
		w.commandOut(who, "Don't have %s.", w.boxNameQty(t.inp1, 2))
		return 0
	}
	if t.inp2 != 0 && w.hasItem(who, t.inp2) < 1 {
		w.commandOut(who, "Don't have %s.", w.boxNameQty(t.inp2, 2))
		return 0
	}
	if t.inp1 != 0 {
		w.consumeItem(who, t.inp1, 1)
	}
	if t.inp2 != 0 {
		w.consumeItem(who, t.inp2, 1)
	}
	w.genItem(who, t.item, 1)
	w.commandOut(who, "%s %s.", capitalize(*t.gotEm), w.justNameQty(t.item, 1))
	if t.public {
		w.commandOut(w.subloc(who), "%s %s %s.", w.boxName(who), *t.gotEm, w.justNameQty(t.item, 1))
	}
	c.d++
	if (t.inp1 == 0 || w.hasItem(who, t.inp1) > 0) &&
		(t.inp2 == 0 || w.hasItem(who, t.inp2) > 0) &&
		!(c.c > 0 && c.d >= c.c) {
		c.wait = int32(t.days)
	}
	if t.reqSkill != 0 {
		w.addSkillExperience(who, t.reqSkill)
	}
	return 1
}

// make.c:445–509. Both MAKE and TRAIN use these same callbacks.
func (w *world) vMake(c *command) int {
	t := findMake(int(c.a))
	if t == nil {
		w.commandOut(int(c.who), "Don't know how to make %s.", w.boxCode(int(c.a)))
		return 0
	}
	if t.days == 1 {
		return w.vGenericMake(c, int(c.b), t)
	}
	return w.vSecondMake(c, int(c.b), t)
}

func (w *world) dMake(c *command) int {
	t := findMake(int(c.a))
	if t == nil {
		w.commandOut(int(c.who), "Internal error.")
		w.logWrite(10, "d_make: t is NULL, who=%d", c.who)
		return 0
	}
	if t.days == 1 {
		return w.dGenericMake(c, t)
	}
	return w.dSecondMake(c, t)
}

func (w *world) iMake(c *command) int {
	t := findMake(int(c.a))
	if t == nil {
		w.commandOut(int(c.who), "Internal error.")
		w.logWrite(10, "i_make: t is NULL, who=%d", c.who)
		return 0
	}
	if t.days == 1 {
		return w.iGenericMake(c, t)
	}
	return 1
}

// make.c:512–535. Reparse replaces USE with MAKE; keep USE table metadata.
func (w *world) vUseTrainRiding(c *command) int {
	s := fmt.Sprintf("make %s %d", boxCodeLess(itemRidingHorse), c.a)
	ret := w.olyParse(c, &s)
	assert(ret, "olympia/make.c", 519, "ret")
	return w.vMake(c)
}

func (w *world) vUseTrainWar(c *command) int {
	s := fmt.Sprintf("make %s %d", boxCodeLess(itemWarmount), c.a)
	ret := w.olyParse(c, &s)
	assert(ret, "olympia/make.c", 532, "ret")
	return w.vMake(c)
}
