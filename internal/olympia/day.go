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
		w.phaseActed(phHealRecover)
	}
	amount := w.streams.UpHeal(who, 1, 3, 15)
	if sick {
		if h-amount < 0 {
			amount = h
		}
		w.commandOut(who, "%s lost %s health.", w.boxName(who), niceNum(amount))
		w.pChar(who).health -= int8(amount)
		w.phaseActed(phHealSicken)
		if w.pChar(who).health <= 0 {
			w.killChar(who, mates)
		}
	} else {
		if h+amount > 100 {
			amount = 100 - h
		}
		w.commandOut(who, "%s gained %s health.", w.boxName(who), niceNum(amount))
		w.pChar(who).health += int8(amount)
		w.phaseActed(phHealHeal)
		if h+amount >= 100 {
			w.commandOut(who, "%s is fully healed.", w.boxName(who))
		}
	}
}

// healCharacters is day.c:616–632, retaining the deletion-safe kind chain.
func (w *world) healCharacters() {
	w.phaseRan(phHealHeal)
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
	w.phaseRan(phGarrPay)
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
		w.phaseActed(phGarrPay)
	}
}

// maintCost is the complete day.c:1171–1199 item upkeep tariff.
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

// menStarve is day.c:1200–1301. Payment and casualties both cycle through
// inventory order, not tariff order. Unlike unitMaintCost, this includes the
// noble item. All arithmetic that is C int remains signed 32-bit.
func (w *world) menStarve(who int, have int32) {
	var items []int
	var qty, cost, starve []int32
	var nmen, npaid int32
	assert(w.validBox(who), "olympia/day.c", 1223, "valid_box(who)")
	for e := range w.loopInv(who) {
		if n := maintCost(int(e.item)); n != 0 {
			items = append(items, int(e.item))
			qty = append(qty, e.qty)
			cost = append(cost, int32(n))
			starve = append(starve, 0)
			nmen += e.qty
		}
	}
	gold := have
	for {
		hitOne := false
		for i := range items {
			if qty[i] > 0 && have >= cost[i] {
				have -= cost[i]
				qty[i]--
				npaid++
				hitOne = true
			}
		}
		if !hitOne || have <= 0 {
			break
		}
	}
	gold -= have
	nstarve := nmen - npaid
	nstarve = (nstarve + 2) / 3
	assert(nstarve > 0, "olympia/day.c", 1259, "nstarve > 0")
	var failcheck int32
	for i := 0; nstarve > 0; {
		assert(failcheck < 10000, "olympia/day.c", 1264, "failcheck++ < 10000")
		failcheck++
		if qty[i] != 0 {
			nstarve--
			starve[i]++
			qty[i]--
		}
		i++
		if i >= len(items) {
			i = 0
		}
	}
	w.autocharge(who, int(gold))
	for i, item := range items {
		if starve[i] == 0 {
			continue
		}
		s := "starved"
		if item != itemPeasant {
			if w.streams.UpStarve(who, item, 1, 2) == 1 {
				s = "left service"
			} else {
				s = "deserted"
			}
		}
		w.commandOut(who, "%s %s.", capitalize(w.justNameQty(item, int(starve[i]))), s)
		w.consumeItem(who, item, int(starve[i]))
		if item == itemSailor || item == itemPirate {
			w.checkCaptainLosesSailors(starve[i], who, 0)
		}
	}
}

// unitMaintCost is day.c:1304–1318, including noble-item exclusion.
func (w *world) unitMaintCost(who int) int32 {
	var cost int32
	assert(w.validBox(who), "olympia/day.c", 1310, "valid_box(who)")
	for e := range w.loopInv(who) {
		if int(e.item) != w.nobleItem(who) {
			cost += int32(maintCost(int(e.item))) * e.qty
		}
	}
	return cost
}

// chargeMaintSup is day.c:1325–1349: 0 owes nothing, 1 paid, 2 starved.
func (w *world) chargeMaintSup(who int) int {
	cost := w.unitMaintCost(who)
	if cost < 1 {
		return 0
	}
	if w.autocharge(who, int(cost)) {
		w.commandOut(who, "Paid maintenance of %s.", goldS(int(cost)))
		return 1
	}
	have := w.stackHasItem(who, itemGold)
	w.commandOut(who, "Maintenance costs are %s, can afford %s.", goldS(int(cost)), goldS(int(have)))
	w.menStarve(who, have)
	return 2
}

// chargeMaintCosts is day.c:1352–1376. Prisoners and unsworn regular
// faction characters are not excluded by the active C source.
func (w *world) chargeMaintCosts() {
	w.phaseRan(phMaintPaid)
	for who := range w.loopKind(tChar) {
		if w.subkind(w.player(who)) == subPlRegular {
			switch w.chargeMaintSup(who) {
			case 1:
				w.phaseActed(phMaintPaid)
			case 2:
				w.phaseActed(phMaintStarve)
			}
		}
	}
}

// animalDeaths is day.c:1379–1418, using the original per-individual,
// one-based leaf keys and strict <10 threshold (nine chances per thousand).
func (w *world) animalDeaths() {
	w.stage("animal_deaths()")
	w.phaseRan(phAnimal)
	for who := range w.loopKind(tChar) {
		if w.subkind(w.player(who)) != subPlRegular {
			continue
		}
		assert(w.validBox(who), "olympia/day.c", 1394, "valid_box(who)")
		for e := range w.loopInv(who) {
			p := w.rpItem(int(e.item))
			if e.qty <= 0 || p == nil || p.animal == 0 {
				continue
			}
			var dead int32
			for i := int32(1); i <= e.qty; i++ {
				if w.streams.UpAnimal(who, int(e.item), int(i), 1, 1000) < 10 {
					dead++
				}
			}
			if dead > 0 {
				w.consumeItem(who, int(e.item), int(dead))
				w.phaseActed(phAnimal)
				w.commandOut(who, "%s %s died.", capitalize(niceNum(int(dead))), w.pluralItemName(int(e.item), int(dead)))
			}
		}
	}
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
	w.phaseRan(phStormWeaken)
	for i := range w.loopKind(tStorm) {
		p := w.pMisc(i)
		p.stormStr--
		w.phaseActed(phStormWeaken)
		if p.stormStr > 0 {
			continue
		}
		p.stormStr = 0
		w.dissipateStorm(i, 1)
		w.phaseActed(phStormGone)
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
	w.phaseRan(phStormMove)
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
		w.phaseActed(phStormMove)
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

// Faction settlement, day.c:635–719. Claim/NP grants filter by subtype;
// unformed replenishment instead visits every existing player record. Keep
// list entries verbatim, including stale entries and lists longer than five.
func (w *world) addClaimGold() {
	w.phaseRan(phClaim)
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) == subPlRegular {
			w.genItem(pl, itemGold, 25)
			w.phaseActed(phClaim)
		}
	}
}

func (w *world) addNoblePoints() {
	w.phaseRan(phNp)
	for pl := range w.loopKind(tPlayer) {
		if w.subkind(pl) == subPlRegular && w.nextNPTurn(pl) == 0 {
			w.addNP(pl, 1)
			w.phaseActed(phNp)
		}
	}
}

func (w *world) addUnformedSup(pl int) {
	p := w.rpPlayer(pl)
	if p == nil {
		return
	}
	for len(p.unformed) < 5 {
		n := w.newEnt(tUnform, 0)
		if n <= 0 {
			break
		}
		p.unformed.Append(n)
	}
}

func (w *world) addUnformed() {
	w.phaseRan(phUnformed)
	for pl := range w.loopKind(tPlayer) {
		p := w.rpPlayer(pl)
		n := 0
		if p != nil {
			n = len(p.unformed)
		}
		w.addUnformedSup(pl)
		if p != nil && len(p.unformed) > n {
			w.phaseActed(phUnformed)
		}
	}
}

// collectTaxes is day.c:1756–1789: each castle consumes the province's
// remaining cookies before truncating, including zero-tax output to its owner.
func (w *world) collectTaxes() {
	w.phaseRan(phTaxes)
	for fort := range w.loopSubkind(subCastle) {
		prov, owner := w.province(fort), w.buildingOwner(fort)
		if owner == 0 {
			continue
		}
		amount := w.hasItem(prov, itemTaxCookie)
		w.consumeItem(prov, itemTaxCookie, int(amount))
		amount /= 2
		w.goldTaxes += amount
		w.genItem(owner, itemGold, int(amount))
		w.commandOut(owner, "Collected %s in taxes.", goldS(int(amount)))
		w.phaseActed(phTaxes)
	}
}

// Seasonal announcements, day.c:2046–2127. These do not change route flags;
// postMonth owns the calendar gates. The else-if precedence is significant.
func (w *world) specialLocsOpen() {
	w.phaseRan(phSpecialOpen)
	for i := range w.loopProvince() {
		if w.summerbridge(i) == 1 {
			w.logWrite(10, "%s open to the north.", w.boxName(i))
			w.commandOut(i, "The swamps of Summerbridge have dried enough to permit passage north.")
			w.phaseActed(phSpecialOpen)
		} else if w.summerbridge(i) == 2 {
			w.logWrite(10, "%s open to the south.", w.boxName(i))
			w.commandOut(i, "The swamps of Summerbridge have dried enough to permit passage south.")
			w.phaseActed(phSpecialOpen)
		} else if w.uldim(i) == 3 {
			w.logWrite(10, "%s open to the south.", w.boxName(i))
			w.commandOut(i, "The snows blocking Uldim pass to the south have melted.")
			w.phaseActed(phSpecialOpen)
		} else if w.uldim(i) == 4 {
			w.logWrite(10, "%s open to the north.", w.boxName(i))
			w.commandOut(i, "The snows blocking Uldim pass to the north have melted.")
			w.phaseActed(phSpecialOpen)
		}
	}
}

func (w *world) specialLocsClose() {
	w.phaseRan(phSpecialClose)
	for i := range w.loopProvince() {
		if w.summerbridge(i) == 1 {
			w.logWrite(10, "%s closed to the north.", w.boxName(i))
			w.commandOut(i, "Seasonal rains have made Summerbridge an impassable bog to the north.")
			w.phaseActed(phSpecialClose)
		} else if w.summerbridge(i) == 2 {
			w.logWrite(10, "%s closed to the south.", w.boxName(i))
			w.commandOut(i, "Seasonal rains have made Summerbridge an impassable bog to the south.")
			w.phaseActed(phSpecialClose)
		} else if w.uldim(i) == 3 {
			w.logWrite(10, "%s closed to the south.", w.boxName(i))
			w.commandOut(i, "Falling snow blocks Uldim pass to the south for the winter.")
			w.phaseActed(phSpecialClose)
		} else if w.uldim(i) == 4 {
			w.logWrite(10, "%s closed to the north.", w.boxName(i))
			w.commandOut(i, "Falling snow blocks Uldim pass to the north for the winter.")
			w.phaseActed(phSpecialClose)
		}
	}
}

// autoDrop is the complete day.c:1133–1167 autoquit queue/log pass.
func (w *world) autoDrop() {
	w.phaseRan(phAutoDrop)
	for pl := range w.loopSubkind(subPlRegular) {
		p := w.pPlayer(pl)
		if int32(w.sysclock.turn)-p.lastOrderTurn >= w.autoQuitTurns {
			name, email := "", ""
			if p.fullName != nil {
				name = *p.fullName
			}
			if p.email != nil {
				email = *p.email
			}
			w.queue(pl, "quit")
			w.phaseActed(phAutoDrop)
			w.logWrite(11, "Queued drop for %s", w.boxName(pl))
			w.logWrite(11, "    %s <%s>", name, email)
		}
	}
}

// linkDecay is day.c:1966–1997: positive countdown precedes month reopening;
// negative openings remain permanent, and absent destinations do not decay.
func (w *world) linkDecay() {
	w.phaseRan(phLinkDec)
	for i := range w.loopKind(tLoc) {
		p := w.rpSubloc(i)
		if p == nil || len(p.linkTo) < 1 {
			continue
		}
		if p.linkOpen > 0 {
			p.linkOpen--
			w.phaseActed(phLinkDec)
		}
		if int(p.linkWhen) == olyMonth(w.sysclock) && p.linkOpen < 2 && p.linkOpen >= 0 {
			p.linkOpen = 2
			w.phaseActed(phLinkOpen)
		}
	}
}

// questDecay is day.c:1737–1753. Despite its name, subloc_quest is the
// quest_late counter itself (oly.h:1082), not a separate quest ID.
func (w *world) questDecay() {
	w.phaseRan(phQuest)
	for where := range w.loopKind(tLoc) {
		if p := w.rpSubloc(where); p != nil && p.questLate > 0 {
			w.pSubloc(where).questLate--
			w.phaseActed(phQuest)
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
	w.phaseRan(phProduction) // not the initialization call from loadDb
	w.phaseActed(phProduction)
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
