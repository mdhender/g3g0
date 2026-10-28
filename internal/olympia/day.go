// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "slices"

// Ported from olympia/day.c: location-touching state, including day 0.

// healCharSup is day.c:519–578, using the world's guarded, per-day upkeep
// leaves. Recovery precedes the independent amount draw and healing branch.
func (w *world) healCharSup(who int) {
	h, sick := 0, false
	if p := w.rpChar(who); p != nil {
		h, sick = int(p.health), p.sick != 0
	}
	if h >= 100 {
		return
	}
	chance := 5
	if w.subkind(w.subloc(who)) == subInn {
		chance = 10
	}
	if sick && w.streams.UpHeal(who, 0, 1, 100) <= chance {
		w.commandOut(who, "%s defeated illness and is now recovering.", w.boxName(who))
		w.pChar(who).sick = 0
		sick = false
	}
	amount := w.streams.UpHeal(who, 1, 3, 15)
	if sick {
		if h-amount < 0 {
			amount = h
		}
		w.commandOut(who, "%s lost %s health.", w.boxName(who), niceNum(amount))
		w.pChar(who).health -= int8(amount)
		if w.pChar(who).health <= 0 {
			w.killChar(who, mates)
		}
	} else {
		if h+amount > 100 {
			amount = 100 - h
		}
		w.commandOut(who, "%s gained %s health.", w.boxName(who), niceNum(amount))
		w.pChar(who).health += int8(amount)
		if h+amount >= 100 {
			w.commandOut(who, "%s is fully healed.", w.boxName(who))
		}
	}
}

// healCharacters is day.c:616–632, retaining the deletion-safe kind chain.
func (w *world) healCharacters() {
	for who := range w.loopKind(tChar) {
		n := 0
		if p := w.rpChar(who); p != nil {
			n = int(p.health)
		}
		if n >= 0 && n < 100 {
			w.healCharSup(who)
		}
	}
}

// defaultGarrisonPay is day.c:581–613, including payments to non-immediate
// descendants, zero-pay output and gen_item's negative-pay assertion.
func (w *world) defaultGarrisonPay() {
	for who := range w.loopKind(tChar) {
		if w.isPrisoner(who) {
			continue
		}
		parent := w.stackLeader(who)
		if parent == 0 || w.defaultGarrison(parent) == 0 {
			continue
		}
		if w.defaultGarrison(who) != 0 {
			continue
		}
		w.genItem(who, itemGold, int(w.garrisonPay))
		if !w.isNPC(who) {
			w.goldGarrison += w.garrisonPay
		}
		w.commandOut(who, "Received %s from %s.", goldS(int(w.garrisonPay)), w.boxName(parent))
	}
}

// maintCost is the complete day.c:1171–1199 item upkeep tariff. Charging
// maintenance and starvation remain separate unfinished month-end phases.
func maintCost(item int) int {
	switch item {
	case itemPeasant:
		return 1
	case itemWorker, itemSoldier, itemSailor, itemAngryPeasant, itemCrossbowman:
		return 2
	case itemBlessedSoldier, itemPikeman, itemSwordsman, itemPirate, itemArcher:
		return 3
	case itemKnight, itemEliteArch:
		return 4
	case itemEliteGuard:
		return 5
	}
	return 0
}

// touchLocPl is touch_loc_pl(): it marks where's viewing location, and the
// location around that if it is inside a province, as touched by player pl.
func (w *world) touchLocPl(pl, where int) {
	p := w.pPlayer(pl)
	where = w.viewloc(where)

	if where == 0 {
		return
	}

	setBit(&p.locs, where)

	if w.locDepth(where) > locProvince {
		inside := w.loc(where)
		setBit(&p.locs, inside)
	}
}

// touchLoc is touch_loc(): touchLocPl for who's player and location.
func (w *world) touchLoc(who int) {
	where := w.subloc(who)

	w.touchLocPl(w.player(who), where)
}

// stormDecay is day.c:1655–1678, including signed-short decrement and
// live kind-chain traversal while dissipation deletes adjacent storms.
func (w *world) stormDecay() {
	for i := range w.loopKind(tStorm) {
		p := w.pMisc(i)
		p.stormStr--
		if p.stormStr > 0 {
			continue
		}
		p.stormStr = 0
		w.dissipateStorm(i, 1)
	}
}

// stormOwnerTouchLoc is day.c:1681–1702, called by the shared day-zero
// initializer below, not merely by a weather-specific harness.
func (w *world) stormOwnerTouchLoc() {
	for i := range w.loopKind(tStorm) {
		if p := w.rpMisc(i); p != nil && p.summonedBy != 0 {
			where := w.subloc(i)
			if pl := w.player(int(p.summonedBy)); pl != 0 {
				w.touchLocPl(pl, where)
			}
		}
	}
}

// stormMove is day.c:1705–1734. npc_dir gates a previously selected
// destination; monthly relocation deliberately bypasses moveStorm output.
func (w *world) stormMove() {
	for i := range w.loopKind(tStorm) {
		p := w.pMisc(i)
		if p.npcDir == 0 {
			continue
		}
		assert(w.locDepth(int(p.stormMove)) == locProvince, "olympia/day.c", 1721, "loc_depth(p->storm_move) == LOC_province")
		w.setWhere(i, int(p.stormMove))
		if owner := w.npcSummoner(i); owner != 0 && w.validBox(owner) {
			w.setKnown(owner, int(p.stormMove))
		}
		p.npcDir = 0
		p.stormMove = 0
	}
}

func (w *world) initLocsTouched() {
	for who := range w.loopKind(tChar) {
		if !w.isPrisoner(who) {
			w.touchLoc(who)
		}
	}
	w.stormOwnerTouchLoc()
}

func (w *world) clearOrdersSent() {
	for pl := range w.loopKind(tPlayer) {
		if p := w.rpPlayer(pl); p != nil {
			p.sentOrders = 0
		}
	}
}

// postMonth is day.c's settlement sequence. Even report-producing phases
// may mutate the world; none may be skipped on the way to save_db.
func (w *world) postMonth() {
	w.stage("post_month()")
	w.clearOrdersSent()
	if olyMonth(w.sysclock) == 2 {
		w.specialLocsOpen()
	}
	if olyMonth(w.sysclock) == 6 {
		w.specialLocsClose()
	}
	w.moveCityGold()
	w.garrisonGold()
	w.collectTaxes()
	w.addClaimGold()
	w.addNoblePoints()
	w.addUnformed()
	w.incrementCurrentAura()
	w.decrementAbilityShroud()
	w.decrementRegionShroud()
	w.decrementMeditationHinder()
	w.decrementLocBarrier()
	w.loyaltyDecay()
	w.pillageDecay()
	w.relicDecay()
	w.hideMageDecay()
	w.innIncome()
	w.templeIncome()
	w.chargeMaintCosts()
	w.animalDeaths()
	w.ghostWarriorDecay()
	w.corpseDecay()
	w.deadBodyRot()
	w.stormDecay()
	w.stormMove()
	w.collapsedMineDecay()
	w.postProduction()
	if w.autoQuitTurns > 0 {
		w.autoDrop()
	}
	w.linkDecay()
	w.questDecay()
	w.checkTokenUnits()
	w.determineNobleRanks()
	w.postHasBeenRun = 1
}

// dailyEvents retains the C statics on the world, including the weather
// cursor: turn changes do not reset this schedule. The named streams' turn
// guards are independent. Only the first four shuffled days are sorted (G3 #229).
func (w *world) dailyEvents() {
	if w.curseErodeDay == 0 {
		w.curseErodeDay = w.streams.CalDay(0, 1, monthDays)
	}
	if w.faeryDay == 0 {
		w.faeryDay = w.streams.CalDay(1, monthDays/2, monthDays)
	}
	if w.dogBarkDay == 0 {
		w.dogBarkDay = w.streams.CalDay(2, 1, monthDays)
	}
	if len(w.weatherDays) == 0 {
		for i := 1; i <= monthDays; i++ {
			w.weatherDays = append(w.weatherDays, i)
		}
		assert(4 <= monthDays, "olympia/day.c", 2231, "4 <= MONTH_DAYS")
		w.streams.WthrShuffle(w.weatherDays)
		slices.Sort(w.weatherDays[:4])
	}
	w.defaultGarrisonPay()
	w.shipCoastalDamage()
	w.randomLocDamage()
	w.checkOceanChars()
	if w.sysclock.day%7 == 0 {
		w.healCharacters()
		w.weeklyPrisonerEscapeCheck()
	}
	if int(w.sysclock.day) == w.curseErodeDay {
		w.noncreatorCurseErode()
	}
	if w.sysclock.day == 15 {
		w.incrementStoneRingAura()
	}
	if int(w.sysclock.day) == w.dogBarkDay {
		w.dogsBarkAtHiddenChars()
	}
	if w.wdayIndex < 4 && int(w.sysclock.day) == w.weatherDays[w.wdayIndex] {
		w.wdayIndex++
		w.naturalWeather()
	}
	if int(w.sysclock.day) == w.faeryDay {
		w.autoFaery()
	}
}
