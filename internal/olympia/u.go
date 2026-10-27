// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"time"
)

// Ported from olympia/u.c.

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
	assert(w.validBox(who), "olympia/u.c", 2332, "valid_box(who)")
	assert(w.validBox(i), "olympia/u.c", 2333, "valid_box(i)")
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
	assert(w.validBox(who), "olympia/u.c", 1820, "valid_box(who)")
	assert(w.validBox(item), "olympia/u.c", 1830, "valid_box(item)")
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

	case subTemple, subGalley, subRoundship, subCastle, subGalleyNotdone, subRoundshipNotdone, subGhostShip, subTempleNotdone, subInn, subInnNotdone, subCastleNotdone, subMine, subMineNotdone, subMineCollapsed, subTower, subTowerNotdone, subSewer:
		return locBuild

	default:
		fmt.Fprintf(w.stderr, "subkind is %d\n", w.subkind(n))
		assert(false, "olympia/u.c", 386, "FALSE")
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
// The lore branch (queue_lore() for a character that does not know the
// item's lore) is not yet ported; create_unique_item_alloc(), the one
// caller so far, calls it before the item has an item_magic, so its lore
// is always 0 there.
func (w *world) addItem(who, item, qty int) {
	assert(w.validBox(who), "olympia/u.c", 1847, "valid_box(who)")
	assert(w.validBox(item), "olympia/u.c", 1848, "valid_box(item)")
	assert(qty >= 0, "olympia/u.c", 1849, "qty >= 0")

	{
		var lore int
		if p := w.rpItemMagic(item); p != nil {
			lore = int(p.lore)
		}

		if lore != 0 && w.kind(who) == tChar {
			panic("olympia: add_item: test_known and queue_lore are not yet implemented")
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
