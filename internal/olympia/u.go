// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"time"

	"github.com/mdhender/g3g0/internal/clib"
)

// Ported from olympia/u.c.

// findNearestLand is the complete u.c:2010–2092. In particular, retries do
// not reset where, and a missing exit rotates from the original location.
func (w *world) findNearestLand(where int) int {
	origWhere := where
	for tryTwo := 99; tryTwo >= 0; tryTwo-- {
		dir := w.streams.EntLandDir(origWhere, tryTwo)
		for tryOne := 999; tryOne >= 0; tryOne-- {
			if w.subkind(where) != subOcean {
				return where
			}
			assert(w.validBox(where), "olympia/u.c", 2030, "valid_box(where)")
			for i := range w.loopHere(where) {
				if w.subkind(i) == subIsland {
					assert(w.kind(i) == tLoc, "olympia/u.c", 2034, "kind(i) == T_loc")
					return i
				}
			}
			where = w.locationDirection(where, dir)
			check := 0
			for where == 0 {
				where = origWhere
				dir = dir%4 + 1
				check++
				assert(check <= 4, "olympia/u.c", 2053, "check <= 4")
				where = w.locationDirection(where, dir)
			}
		}
		if tryTwo == 99 {
			w.logWrite(10, "find_nearest_land: Plan B")
		}
	}
	w.logWrite(10, "find_nearest_land: Plan C")
	var l []int
	for i := range w.loopKind(tLoc) {
		if w.region(i) == w.region(origWhere) && w.locDepth(i) == locProvince && w.subkind(i) != subOcean {
			l = append(l, i)
		}
	}
	if len(l) < 1 {
		return 0
	}
	return l[w.streams.EntLandPick(origWhere, len(l))]
}

func addS(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func (w *world) deductNP(pl, num int) bool {
	assert(w.kind(pl) == tPlayer, "olympia/u.c", 1746, "kind(pl) == T_player")
	p := w.pPlayer(pl)
	if int(p.noblePoints) < num {
		return false
	}
	p.noblePoints = int16(int32(p.noblePoints) - int32(num))
	p.npSpent = int16(int32(p.npSpent) + int32(num))
	return true
}

func (w *world) addNP(pl, num int) {
	assert(w.kind(pl) == tPlayer, "olympia/u.c", 1764, "kind(pl) == T_player")
	p := w.pPlayer(pl)
	p.noblePoints = int16(int32(p.noblePoints) + int32(num))
	p.npGained = int16(int32(p.npGained) + int32(num))
}

func commaNum(n int) string {
	further := n / 1000000000
	n %= 1000000000
	millions := n / 1000000
	n %= 1000000
	thousands, ones := n/1000, n%1000
	if further == 0 && millions == 0 && thousands == 0 {
		return sout("%d", ones)
	} else if further == 0 && millions == 0 {
		return sout("%d,%.3d", thousands, ones)
	} else if further == 0 {
		return sout("%d,%.3d,%.3d", millions, thousands, ones)
	}
	return sout("%d,%.3d,%.3d,%.3d", further, millions, thousands, ones)
}

var numS = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

func niceNum(n int) string {
	if n > 10 || n < 0 {
		return sout("%s", commaNum(n))
	}
	return numS[n]
}

func goldS(n int) string { return sout("%s~gold", commaNum(n)) }

func (w *world) contacted(a, b int) bool {
	p := w.pChar(a)
	if p.contact.Lookup(b) >= 0 {
		return true
	}
	return p.contact.Lookup(w.player(b)) >= 0
}

func (w *world) charWhere(where, who, target int) bool {
	if where != w.subloc(target) {
		return false
	}
	if w.charReallyHidden(target) || (w.locDepth(where) == locProvince && w.weatherHere(where, subFog) != 0) {
		if w.player(who) == w.player(target) {
			return true
		}
		if target == w.garrisonHere(where) {
			return true
		}
		if w.contacted(target, who) {
			return true
		}
		return w.stackLeader(target) == w.stackLeader(who)
	}
	return true
}

func (w *world) charHere(who, target int) bool {
	return w.charWhere(w.subloc(who), who, target)
}

func (w *world) testKnown(who, i int) bool {
	if who == 0 {
		return false
	}
	assert(w.validBox(who), "olympia/u.c", 2334, "valid_box(who)")
	assert(w.validBox(i), "olympia/u.c", 2335, "valid_box(i)")
	// Preserve C's NULL bx[0] dereference for valid but unowned characters:
	// https://github.com/playbymail/olympia-g3/issues/353 (nonblocking).
	ep := w.rpPlayer(w.player(who))
	return ep != nil && testBit(ep.known, i)
}

func (w *world) locHidden(n int) int {
	if p := w.rpLoc(n); p != nil {
		return int(p.hidden)
	}
	return 0
}

func (w *world) hasItem(who, item int) int32 {
	assert(w.validBox(who), "olympia/u.c", 1822, "valid_box(who)")
	assert(w.validBox(item), "olympia/u.c", 1832, "valid_box(item)")
	for _, e := range w.bx[who].items {
		if int(e.item) == item {
			return e.qty
		}
	}
	return 0
}

// Actual oly.h macros used by the target predicates, with C field widths.
func (w *world) npcProgram(n int) int {
	if p := w.rpChar(n); p != nil {
		return int(p.npcProg)
	}
	return 0
}

func (w *world) isNPC(n int) bool {
	if w.subkind(n) != 0 {
		return true
	}
	p := w.rpChar(n)
	return p != nil && (p.loyKind == loyNpc || p.loyKind == loySummon)
}

func (w *world) onlyDefeatable(n int) int32 {
	if p := w.rpMisc(n); p != nil {
		return p.onlyVulnerable
	}
	return 0
}

func (w *world) safeHaven(n int) int {
	if p := w.rpSubloc(n); p != nil {
		return int(p.safe)
	}
	return 0
}

func (w *world) stormStrength(n int) int32 {
	if p := w.rpMisc(n); p != nil {
		return int32(p.stormStr)
	}
	return 0
}

func (w *world) charReallyHidden(n int) bool {
	p := w.rpMagic(n)
	return p != nil && p.hideSelf != 0 && w.stackParent(n) == 0 && w.countStackAny(n) == 1
}

func (w *world) checkSkill(who, skill int) bool {
	if w.hasSkill(who, skill) < 1 {
		w.commandOut(who, "Requires %s.", w.boxName(skill))
		return false
	}
	return true
}

// countAny and countStackAny are count_any and count_stack_any, used by
// char_alone in d_hide. Accumulate in C int width, in inventory/stack order.
func (w *world) countAny(who int) int32 {
	sum := int32(1)
	if w.subkind(who) == subGarrison {
		sum = 0
	}
	for e := range w.loopInv(who) {
		p := w.rpItem(int(e.item))
		if int(e.item) == itemGhostWarrior || (p != nil && (p.isManItem != 0 || p.attack != 0 || p.defense != 0 || p.missile != 0)) {
			sum += e.qty
		}
	}
	return sum
}

func (w *world) countStackAny(who int) int32 {
	var sum int32
	for i := range w.loopStack(who) {
		sum += w.countAny(i)
	}
	return sum
}

// countStackUnits is u.c:1200–1213, including prisoners and NPC descendants.
func (w *world) countStackUnits(who int) int32 {
	sum := int32(1)
	for range w.allCharHere(who) {
		sum++
	}
	return sum
}

// olytimeIncrement is olytime_increment(): advance the epoch and day,
// without rolling over the day or changing the turn.
func olytimeIncrement(p *olytime) {
	p.daysSinceEpoch++
	p.day++
}

// olytimeTurnChange is olytime_turn_change(): ready the next turn at day
// zero. The daily loop follows this with olytimeIncrement.
func olytimeTurnChange(p *olytime) {
	p.day = 0
	p.turn++
}

// printDot is print_dot(): progress on stderr, indented once and wrapped
// before every 59th, 119th, ... character. The counter survives stages.
func (w *world) printDot(c int) {
	if w.dotCount == 0 {
		fmt.Fprint(w.stderr, "   ")
		w.dotCount++
	}
	w.dotCount++
	if w.dotCount%60 == 0 {
		fmt.Fprint(w.stderr, "\n   ")
	}
	w.stderr.Write([]byte{byte(c)})
}

// testBit is test_bit(): whether i is in the sparse set kr.
func testBit(kr sparse, i int) bool {
	return kr.Lookup(i) != -1
}

// setBit is set_bit(): it adds i to the sparse set kr, at the end, unless
// it is there already.
func setBit(kr *sparse, i int) {
	if kr.Lookup(i) == -1 {
		kr.Append(i)
	}
}

// lookup is lookup(): the index of the first entry of table that is s,
// without regard to case, or -1. C's tables end with a NULL entry; the
// slices here have no terminator (see glob.go).
func lookup(table []string, s string) int {
	for i, t := range table {
		if iStrcmp(s, t) == 0 {
			return i
		}
	}

	return -1
}

// stage is stage(s), with s not NULL: it prints s and a newline to
// stderr, or, if time_self (-T) is set, the seconds since the last stage()
// on a line of their own and then s without a newline.
func (w *world) stage(s string) {
	w.stageSup(&s)
}

// stageSup is stage(), where s may be NULL: stage(NULL) prints nothing
// without time_self, and the seconds since the first stage() with it.
func (w *world) stageSup(s *string) {
	if !w.timeSelf {
		if s != nil {
			fmt.Fprintf(w.stderr, "%s\n", *s)
		}
		return
	}

	t := time.Now().Unix()

	if w.stageOld != 0 {
		fmt.Fprintf(w.stderr, "\t%d sec\n", int32(t-w.stageOld))
	} else {
		w.stageFirst = t
	}
	w.stageOld = t

	if s != nil {
		fmt.Fprintf(w.stderr, "%s", *s)
	} else {
		fmt.Fprintf(w.stderr, "%d seconds\n", int32(t-w.stageFirst))
	}
}

// locDepth is loc_depth(): the LOC_xxx depth of location n, by its
// subkind. Any other subkind is reported on stderr and fails an assertion.
func (w *world) locDepth(n int) int {
	switch w.subkind(n) {
	case subRegion:
		return locRegion

	case subOcean, subForest, subPlain, subMountain, subDesert, subSwamp, subUnder, subCloud, subTunnel, subChamber:
		return locProvince

	case subIsland, subStoneCir, subMallornGrove, subBog, subCave, subCity, subLair, subGraveyard, subRuins, subBattlefield, subEnchForest, subRockyHill, subTreeCircle, subPits, subPasture, subOasis, subYewGrove, subSandPit, subSacredGrove, subPoppyField, subFaeryHill, subHadesPit:
		return locSubloc

	case subTemple, subGalley, subRoundship, subCastle, subGalleyNotdone, subRoundshipNotdone, subRaft, subRaftNotdone, subGhostShip, subTempleNotdone, subInn, subInnNotdone, subCastleNotdone, subMine, subMineNotdone, subMineCollapsed, subTower, subTowerNotdone, subSewer:
		return locBuild

	default:
		fmt.Fprintf(w.stderr, "subkind is %d\n", w.subkind(n))
		assert(false, "olympia/u.c", 388, "FALSE")
	}

	return 0
}

// clearTemps is clear_temps(): it zeroes the temp of every box of kind k.
func (w *world) clearTemps(k int) {
	for i := range w.loopKind(k) {
		w.bx[i].temp = 0
	}
}

// capitalize is cap(): s with its first byte made upper case, by C's
// toupper() in the "C" locale (the engine never calls setlocale()), which
// changes only 'a' to 'z'. "" is returned as is. C copies s into a sout()
// buffer first, cutting it to LEN bytes.
func capitalize(s string) string {
	if s == "" {
		return s
	}

	t := []byte(sout("%s", s))
	if 'a' <= t[0] && t[0] <= 'z' {
		t[0] -= 'a' - 'A'
	}

	return string(t)
}

// addItem is add_item(), static in u.c: it gives who qty more of item,
// adding an inventory entry at the end if who has none, and lets
// investigate_possible_trade() see whether that completes a pending SELL.
// Item lore is queued before even a zero-quantity inventory change.
func (w *world) addItem(who, item, qty int) {
	assert(w.validBox(who), "olympia/u.c", 1849, "valid_box(who)")
	assert(w.validBox(item), "olympia/u.c", 1850, "valid_box(item)")
	assert(qty >= 0, "olympia/u.c", 1851, "qty >= 0")

	{
		var lore int
		if p := w.rpItemMagic(item); p != nil {
			lore = int(p.lore)
		}

		if lore != 0 && w.kind(who) == tChar && !w.testKnown(who, item) {
			w.queueLore(who, item, false)
		}
	}

	for _, e := range w.bx[who].items {
		if int(e.item) == item {
			old := int(e.qty)

			e.qty += int32(qty)

			w.investigatePossibleTrade(who, item, old)
			return
		}
	}

	w.bx[who].items.Append(&itemEnt{item: int32(item), qty: int32(qty)})

	w.investigatePossibleTrade(who, item, 0)
}

// createUniqueItemAlloc is create_unique_item_alloc(): it creates box new
// as a unique item of subkind sk and gives it to who.
func (w *world) createUniqueItemAlloc(new, who, sk int) int {
	w.allocBox(new, tItem, sk)

	w.pItem(new).whoHas = int32(who)
	w.addItem(who, new, 1)

	return new
}

func (w *world) newChar(sk, ni, where, health, pl, loyKind, loyLev int, name string) int {
	n := w.newEnt(tChar, sk)
	if n < 0 {
		return -1
	}
	if name != "" {
		w.setName(n, name)
	}
	p := w.pChar(n)
	p.health, p.unitItem, p.breakPoint = int8(health), int32(ni), 50
	p.attack, p.defense = 60, 60
	if w.kind(where) == tLoc || w.kind(where) == tShip {
		w.setWhere(n, where)
	} else {
		w.setWhere(n, w.subloc(where))
	}
	w.setLord(n, pl, loyKind, loyLev)
	if w.kind(where) == tChar {
		w.joinStack(n, where)
	}
	if w.beastCapturable(n) || w.isNPC(n) {
		p.breakPoint = 0
	}
	return n
}

func (w *world) beastCapturable(who int) bool {
	if w.subkind(who) != subNi {
		return false
	}
	p := w.rpItem(w.nobleItem(who))
	return p != nil && p.capturable != 0
}

func (w *world) createUniqueItem(who, sk int) int {
	n := w.newEnt(tItem, sk)
	if n < 0 {
		return -1
	}
	w.pItem(n).whoHas = int32(who)
	w.addItem(who, n, 1)
	return n
}

func (w *world) subItem(who, item, qty int) bool {
	assert(w.validBox(who), "olympia/u.c", 643, "valid_box(who)")
	assert(w.validBox(item), "olympia/u.c", 644, "valid_box(item)")
	assert(qty >= 0, "olympia/u.c", 645, "qty >= 0")
	for _, e := range w.bx[who].items {
		if int(e.item) == item {
			if int(e.qty) < qty {
				return false
			}
			e.qty -= int32(qty)
			return true
		}
	}
	return false
}

// genItem is u.c:1885–1891. addItem retains lore and trade bookkeeping.
func (w *world) genItem(who, item, qty int) {
	assert(w.itemUnique(item) == 0, "olympia/u.c", 1889, "!item_unique(item)")
	w.addItem(who, item, qty)
}

// Inheritance selectors and quantities are the original u.h constants.
const (
	mates       = -1
	matesSilent = -2
	takeAll     = 1
	takeSome    = 2
	takeNI      = 3
)

func (w *world) stackmateInheritor(who int) int {
	for i := range w.loopHere(who) {
		if w.kind(i) == tChar && !w.isPrisoner(i) {
			return i
		}
	}
	return w.stackParent(who)
}

// takeUnitItems is u.c:422–538. loopInv supplies copied quantities: the
// remainder is computed from the snapshot, not the mutated live inventory.
func (w *world) takeUnitItems(from, inherit, howMany int) {
	to, silent := inherit, false
	switch inherit {
	case 0:
		silent = true
	case mates:
		to = w.stackmateInheritor(from)
	case matesSilent:
		to, silent = w.stackmateInheritor(from), true
	}
	first := true
	if howMany == takeNI {
		w.genItem(from, w.nobleItem(from), 1)
	}
	for e := range w.loopInv(from) {
		item, qty := int(e.item), int(e.qty)
		if e.qty != 0 && howMany == takeSome && w.streams.EntTake(from, item, 0, 1, 2) == 1 {
			qty = w.streams.EntTake(from, item, 1, 0, int(e.qty))
		}
		if qty > 0 && !silent && w.validBox(to) {
			if first {
				first = false
				w.commandOut(to, "Taken from %s:", w.boxName(from))
				w.indent += 3
			}
			w.commandOut(to, "%s", w.boxNameQty(item, qty))
		}
		w.moveItem(from, to, item, qty)
		if item == itemGold && w.player(from) != w.player(to) {
			if w.player(from) == int(w.indepPlayer) {
				w.goldCombatIndep += int32(qty)
			} else {
				w.goldCombat += int32(qty)
			}
		}
		if qty != int(e.qty) {
			w.moveItem(from, 0, item, int(e.qty)-qty)
		}
	}
	for i := range w.loopHere(from) {
		if w.kind(i) != tChar || !w.isPrisoner(i) {
			continue
		}
		if to > 0 {
			if first && !silent {
				w.commandOut(to, "Taken from %s:", w.boxName(from))
				w.indent += 3
				first = false
			}
			w.movePrisoner(from, to, i)
			if !silent {
				w.commandOut(to, "%s", w.linerDesc(i))
			}
			if w.player(i) == w.player(to) {
				w.pChar(i).prisoner = 0
			}
		} else {
			w.pMagic(i).swearOnRelease = 0
			w.dropStack(from, i)
		}
	}
	if !first && !silent {
		w.indent -= 3
	}
}

func (w *world) putBackCookie(who int) {
	p := w.rpMisc(who)
	if p == nil || p.npcHome == 0 {
		return
	}
	w.genItem(int(p.npcHome), int(p.npcCookie), 1)
}

func (w *world) consumeItem(who, item, qty int) bool {
	assert(w.itemUnique(item) == 0, "olympia/u.c", 1898, "!item_unique(item)")
	return w.subItem(who, item, qty)
}

func (w *world) moveItem(from, to, item, qty int) bool {
	if qty <= 0 {
		return true
	}
	if to == 0 {
		return w.dropItem(from, item, qty)
	}
	if !w.subItem(from, item, qty) {
		return false
	}
	w.addItem(to, item, qty)
	if w.itemUnique(item) != 0 {
		assert(qty == 1, "olympia/u.c", 1924, "qty == 1")
		w.pItem(item).whoHas = int32(to)
		if w.subkind(item) == subNpcToken {
			w.moveToken(item, from, to)
		}
	}
	return true
}

func (w *world) dropItem(who, item, qty int) bool {
	if w.itemUnique(item) == 0 {
		return w.consumeItem(who, item, qty)
	}
	whoGets := w.province(who)
	if w.subkind(item) == subDeadBody && whoGets == 0 {
		w.destroyUniqueItem(who, item)
		return true
	}
	if w.subkind(whoGets) == subOcean {
		whoGets = w.findNearestLand(whoGets)
	}
	if whoGets == 0 {
		whoGets = w.province(who)
	}
	w.logWrite(10, "drop_item: %s from %s to %s", w.boxName(item), w.boxName(w.subloc(who)), w.boxName(whoGets))
	return w.moveItem(who, whoGets, item, qty)
}

func (w *world) canPay(who, amount int) bool {
	return w.hasItem(who, itemGold) >= int32(amount)
}

func (w *world) charge(who, amount int) bool {
	return w.subItem(who, itemGold, amount)
}

func (w *world) stackHasItem(who, item int) int32 {
	var sum int32
	for i := range w.loopStack(w.stackLeader(who)) {
		if w.player(i) != w.player(who) {
			continue
		}
		sum += w.hasItem(i, item)
	}
	return sum
}

func (w *world) hasUseKey(who, key int) int {
	for e := range w.loopInv(who) {
		if p := w.rpItemMagic(int(e.item)); p != nil && int(p.useKey) == key {
			return int(e.item)
		}
	}
	return 0
}

func (w *world) stackHasUseKey(who, key int) int {
	for i := range w.loopStack(w.stackLeader(who)) {
		if w.player(i) != w.player(who) {
			continue
		}
		if item := w.hasUseKey(i, key); item != 0 {
			return item
		}
	}
	return 0
}

// stackSubItem preserves C's actor-first, then friendly stack-order borrowing.
func (w *world) stackSubItem(who, item, qty int) bool {
	if w.stackHasItem(who, item) < int32(qty) {
		return false
	}
	n := min(w.hasItem(who, item), int32(qty))
	if n > 0 {
		qty = int(int32(qty) - n)
		w.subItem(who, item, int(n))
	}
	assert(qty >= 0, "olympia/u.c", 2241, "qty >= 0")
	if qty == 0 {
		return true
	}
	for i := range w.loopStack(w.stackLeader(who)) {
		if qty <= 0 {
			break
		}
		if w.player(i) != w.player(who) {
			continue
		}
		n = min(w.hasItem(i, item), int32(qty))
		if n > 0 {
			qty = int(int32(qty) - n)
			w.subItem(i, item, int(n))
		}
	}
	assert(qty == 0, "olympia/u.c", 2283, "qty == 0")
	return true
}

func (w *world) autocharge(who, amount int) bool {
	return w.stackSubItem(who, itemGold, amount)
}

func (w *world) destroyUniqueItem(who, item int) {
	assert(w.kind(item) == tItem, "olympia/u.c", 1982, "kind(item) == T_item")
	assert(w.itemUnique(item) != 0, "olympia/u.c", 1983, "item_unique(item)")
	if w.subkind(item) == subDeadBody {
		pl := 0
		if p := w.rpMisc(item); p != nil {
			pl = int(p.oldLord)
		}
		if pl == int(w.indepPlayer) {
			pl = int(w.pChar(item).prevLord)
		}
		if w.kind(pl) == tPlayer {
			nps := w.charNPTotal(item)
			w.commandOut(pl, "%s~%s has passed on.  Gained %d NP%s.", clib.Str(w.rpMisc(item).saveName), w.boxCode(item), nps, addS(nps))
			w.addNP(pl, nps)
		}
	}
	ret := w.subItem(who, item, 1)
	assert(ret, "olympia/u.c", 2004, "ret")
	w.deleteBox(item)
}
