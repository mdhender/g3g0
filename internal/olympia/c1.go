// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

// Complete c1.c:314–359, including the static permission helper. Item
// naming deliberately does not require possession except for an auraculum.
func (w *world) mayName(who, target int) bool {
	switch w.kind(target) {
	case tChar, tPlayer:
		return w.player(who) == w.player(target)
	case tLoc, tShip:
		if w.kind(target) == tLoc && w.safeHaven(target) != 0 {
			return false
		}
		if w.locDepth(target) == locBuild {
			return w.player(who) == w.player(w.buildingOwner(target))
		}
		return w.mayRuleHere(who, target)
	case tItem:
		if w.hasAuraculum(who) == target {
			return true
		}
		if p := w.rpItemMagic(target); p != nil {
			if int(p.creator) == who {
				return true
			}
			switch int(p.useKey) {
			case useDeathPotion, useHealPotion, useSlavePotion, useProjCast, useQuickCast:
				return true
			}
		} else if who == 0 { // item_creator's absent-record value is zero.
			return true
		}
	case tStorm:
		return w.npcSummoner(target) == who
	}
	return false
}

// Complete c1.c:362–419. cmdShift changes the saved parse and numeric
// fields; validation order and byte lengths are the original C contract.
func (w *world) vName(c *command) int {
	who, target := int(c.who), int(c.who)
	if c.parse.Len()-1 >= 2 && c.a > 0 {
		target = int(c.a)
		w.cmdShift(c)
	}
	newName := restName(c, 1)
	if newName == nil || *newName == "" {
		w.commandOut(who, "No new name given.")
		return 0
	}
	if !w.mayName(who, target) {
		w.commandOut(who, "Not allowed to change the name of %s.", w.boxCode(target))
		return 0
	}
	limit := 25
	if w.kind(target) == tChar {
		limit = 35
	}
	if len(*newName) > limit {
		w.commandOut(who, "Name is longer than %d characters.", limit)
		return 0
	}
	oldName := w.boxName(target)
	w.setName(target, *newName)
	w.commandOut(who, "%s will now be known as %s.", oldName, w.boxName(target))
	if target != who && (w.kind(target) == tChar || w.isLocOrShip(target)) {
		w.commandOut(target, "%s will now be known as %s.", oldName, w.boxName(target))
	}
	return 1
}

// Complete c1.c:422–441. Remap before narrowing to signed char; the
// receives-Times message tests the narrowed field, not the original int.
func (w *world) vTimes(c *command) int {
	if c.a == 1 {
		c.a = 2
	} else if c.a == 2 {
		c.a = 1
	}
	p := w.pPlayer(w.player(int(c.who)))
	p.compuserve = int8(c.a)
	if p.compuserve != 0 {
		w.commandOut(int(c.who), "Will not receive the Times.")
	} else {
		w.commandOut(int(c.who), "Will receive the Times.")
	}
	return 1
}

// Complete c1.c:444–470. Refusals do not replace an existing fullName,
// including a nonnil empty string, or allocate a player before output.
func (w *world) vFullname(c *command) int {
	newName := restName(c, 1)
	if newName == nil || *newName == "" {
		w.commandOut(int(c.who), "No new name given.")
		return 0
	}
	if len(*newName) > 60 {
		w.commandOut(int(c.who), "Name is longer than %d characters.", 60)
		return 0
	}
	w.pPlayer(w.player(int(c.who))).fullName = newName
	return 1
}

// Complete active c1.c:473–526. Only an explicit target takes the valid-box
// and may_name checks. setBanner owns allocation and NULL/empty clearing.
func (w *world) vBanner(c *command) int {
	who, target := int(c.who), int(c.who)
	if c.parse.Len()-1 >= 2 && c.a > 0 {
		target = int(c.a)
		w.cmdShift(c)
		if !w.validBox(target) {
			w.commandOut(who, "%s is not a valid entity.", w.boxCode(target))
			return 0
		}
		if !w.mayName(who, target) {
			w.commandOut(who, "You do not control %s.", w.boxCode(target))
			return 0
		}
	}
	newName := restName(c, 1)
	text := ""
	if newName != nil {
		text = *newName
	}
	if len(text) > 50 {
		w.commandOut(who, "Banner is longer than 50 characters.")
		return 0
	}
	w.setBanner(target, text)
	if text != "" {
		w.commandOut(who, "Banner set.")
	} else {
		w.commandOut(who, "Banner cleared.")
	}
	return 1
}

// Complete c1.c:1152–1179; share the actual WAIT signal list and lookup.
func (w *world) vFlag(c *command) int {
	if c.parse.Len()-1 < 1 {
		w.commandOut(int(c.who), "Must specify what message to signal.")
		return 0
	}
	flag := c.parse[1]
	if w.flagRaised(int(c.who), flag) >= 0 {
		w.commandOut(int(c.who), "%s has already given that signal this month.", w.boxName(int(c.who)))
		return 0
	}
	w.flags = append(w.flags, &flagEnt{who: c.who, flag: flag})
	return 1
}

// c1.c:1694–1736: SPLIT is unconditionally unsupported. The old report
// splitting implementation is a C comment, not active settings behavior.
func (w *world) vSplit(c *command) int {
	w.commandOut(int(c.who), "This command is no longer supported.")
	return 0
}

// Complete c1.c:529–651. ACCEPT deliberately appends unvalidated raw fields.
func (w *world) vAccept(c *command) int {
	p := w.pChar(int(c.who))
	p.accept.Append(&acceptEnt{fromWho: c.a, item: c.b, qty: c.c})
	return 1
}

func (w *world) willAcceptSup(who, item, from, qty int) bool {
	if p := w.rpChar(who); p != nil {
		for _, e := range p.accept {
			if (e.item == int32(item) || e.item == 0) &&
				(e.fromWho == int32(from) || e.fromWho == 0) &&
				(e.qty >= int32(qty) || e.qty == 0) {
				// Exactly consuming a finite allowance leaves C's unlimited zero
				// sentinel in place; do not remove or mark the row exhausted.
				if e.qty != 0 {
					e.qty -= int32(qty)
				}
				return true
			}
		}
	}
	return false
}

func (w *world) willAccept(who, item, from, qty int) bool {
	if item == itemGold || w.player(who) == w.player(from) {
		return true
	}
	if w.subkind(who) == subGarrison {
		if w.mayRuleHere(from, who) {
			return true
		}
		w.commandOut(from, "%s is not under your control.", w.boxName(who))
		return false
	}
	if w.willAcceptSup(who, item, from, qty) ||
		w.willAcceptSup(w.player(who), item, from, qty) ||
		w.willAcceptSup(who, item, w.player(from), qty) ||
		w.willAcceptSup(w.player(who), item, w.player(from), qty) {
		return true
	}
	w.commandOut(who, "Refused %s from %s.", w.justNameQty(item, qty), w.boxName(from))
	w.commandOut(from, "Refused by %s.", w.justName(who))
	return false
}

func (w *world) howMany(who, fromWho, item, qty, haveLeft int) int {
	numHas := w.hasItem(fromWho, item)
	if numHas <= 0 {
		w.commandOut(who, "%s has no %s.", w.justName(fromWho), w.justName(item))
		return 0
	}
	if numHas <= int32(haveLeft) {
		w.commandOut(who, "%s has only %s.", w.justName(fromWho), w.justNameQty(item, int(numHas)))
		return 0
	}
	q := int32(qty)
	if q == 0 {
		q = numHas
	}
	q = min(numHas-int32(haveLeft), q)
	assert(q > 0, "olympia/c1.c", 648, "qty > 0")
	return int(q)
}

// Complete c1.c:658–760. Unlike CLAIM's local correction, GIVE swaps the
// saved numeric command fields, but leaves its original line and parse list.
func (w *world) vGive(c *command) int {
	who, target, item := int(c.who), int(c.a), int(c.b)
	qty, haveLeft := int(c.c), int(c.d)
	if c.parse.Len()-1 >= 2 && (w.kind(target) == tItem || w.hasPrisoner(who, target)) &&
		w.kind(item) == tChar && w.subloc(who) == w.subloc(item) && !w.hasPrisoner(who, item) {
		c.a, c.b = c.b, c.a
		target, item = int(c.a), int(c.b)
		switch c.parse.Len() - 1 {
		case 2:
			w.commandOut(who, "(assuming you meant 'give %d %d')", target, item)
		case 3:
			w.commandOut(who, "(assuming you meant 'give %d %d %d')", target, item, qty)
		default:
			w.commandOut(who, "(assuming you meant 'give %d %d %d %d')", target, item, qty, haveLeft)
		}
	}
	if !w.checkCharHere(who, target) || !w.checkCharGone(who, target) {
		return 0
	}
	if w.isPrisoner(target) {
		w.commandOut(who, "Prisoners may not be given anything.")
		return 0
	}
	if p := w.rpChar(target); p != nil && p.loyKind == loySummon {
		w.commandOut(who, "Summoned entities may not be given anything.")
		return 0
	}
	if w.hasPrisoner(who, item) {
		if !w.willAccept(target, item, who, 1) {
			return 0
		}
		if w.givePrisoner(who, target, item) {
			return 1
		}
		return 0
	}
	if w.kind(item) != tItem {
		w.commandOut(who, "%s is not an item or a prisoner.", w.boxCode(item))
		return 0
	}
	if w.hasItem(who, item) < 1 {
		w.commandOut(who, "%s does not have any %s.", w.boxName(who), w.boxCode(item))
		return 0
	}
	qty = w.howMany(who, who, item, qty, haveLeft)
	if qty <= 0 || !w.willAccept(target, item, who, qty) {
		return 0
	}
	ret := w.moveItem(who, target, item, qty)
	assert(ret, "olympia/c1.c", 751, "ret")
	w.commandOut(who, "Gave %s to %s.", w.justNameQty(item, qty), w.boxName(target))
	w.commandOut(target, "Received %s from %s.", w.boxNameQty(item, qty), w.boxName(who))
	return 1
}

// Complete c1.c:763–775. Use the real parser; doCommand owns accounting
// against its saved pre-reparse PAY index, not the resulting GIVE index.
func (w *world) vPay(c *command) int {
	line := sout("give %d 1 %d %d", c.a, c.b, c.c)
	ret := w.olyParse(c, &line)
	assert(ret, "olympia/c1.c", 772, "ret")
	return w.vGive(c)
}

// Complete static c1.c:778–806 prerequisite; garrison control takes
// precedence over the ordinary same-faction/direct-prisoner alternatives.
func (w *world) mayTake(who, target int) bool {
	if !w.checkCharHere(who, target) || !w.checkCharGone(who, target) {
		return false
	}
	if w.subkind(target) == subGarrison {
		if w.mayRuleHere(who, target) {
			return true
		}
		w.commandOut(who, "%s is not under your control.", w.boxName(target))
		return false
	}
	if !w.myPrisoner(who, target) && w.player(target) != w.player(who) {
		w.commandOut(who, "May only take items from other units in your faction.")
		return false
	}
	return true
}

// Complete c1.c:813–874, shared by GET and TAKE. Only this withdrawal path
// enforces the ten-man remainder and notifies a running captain of crew loss.
func (w *world) vGet(c *command) int {
	who, target, item := int(c.who), int(c.a), int(c.b)
	qty, haveLeft := int(c.c), int(c.d)
	if !w.mayTake(who, target) {
		return 0
	}
	if w.hasPrisoner(target, item) {
		if w.givePrisoner(target, who, item) {
			return 1
		}
		return 0
	}
	if w.kind(item) != tItem {
		w.commandOut(who, "%s is not an item or a prisoner.", w.boxCode(item))
		return 0
	}
	if w.hasItem(target, item) < 1 {
		w.commandOut(who, "%s does not have any %s.", w.boxName(target), w.boxCode(item))
		return 0
	}
	qty = w.howMany(who, target, item, qty, haveLeft)
	if qty <= 0 {
		return 0
	}
	if p := w.rpItem(item); w.subkind(target) == subGarrison && p != nil && p.isManItem != 0 {
		men := w.countManItems(target) - int32(qty)
		if men < 10 {
			w.commandOut(who, "Garrisons must be left with a minimum of ten men.")
			return 0
		}
	}
	ret := w.moveItem(target, who, item, qty)
	assert(ret, "olympia/c1.c", 862, "ret")
	w.commandOut(who, "Took %s from %s.", w.justNameQty(item, qty), w.boxName(target))
	w.commandOut(target, "%s took %s from us.", w.boxName(who), w.boxNameQty(item, qty))
	if item == itemSailor || item == itemPirate {
		w.checkCaptainLosesSailors(int32(qty), target, who)
	}
	return 1
}

// Complete c1.c:947–974. Ownership initializes the new noble's queued command
// before joinStack moves it beneath the forming noble; do not use newChar,
// which allocates a different entity and has different initialization.
func (w *world) formNewNoble(who int, name string, new int) {
	assert(w.kind(new) == tUnform, "olympia/c1.c", 953, "kind(new) == T_unform")
	w.changeBoxKind(new, tChar)
	p := w.pChar(new)
	op := w.pChar(who)
	p.behind = op.behind
	p.freshHire = 1
	p.health = 100
	p.attack = 80
	p.defense = 80
	p.breakPoint = 50
	w.setName(new, name)
	w.setWhere(new, w.subloc(who))
	w.setLord(new, w.player(who), loyContract, 500)
	w.joinStack(new, who)
}

func (w *world) vForm(c *command) int {
	if w.subkind(w.subloc(int(c.who))) != subCity {
		w.commandOut(int(c.who), "Nobles may only be formed in cities.")
		return 0
	}
	// noble_cost(pl) is the constant 1 in c1.c:878, not a scaling tariff.
	if p := w.rpPlayer(w.player(int(c.who))); p == nil || p.noblePoints < 1 {
		w.commandOut(int(c.who), "To form another noble requires 1 Noble Point.")
		return 0
	}
	return 1
}

func (w *world) dForm(c *command) int {
	pl := w.player(int(c.who))
	if p := w.rpPlayer(pl); p == nil || p.noblePoints < 1 {
		w.commandOut(int(c.who), "To form another noble requires 1 Noble Point.")
		return 0
	}
	p := w.pPlayer(w.player(int(c.who)))
	new := int(c.a)
	if new != 0 && (w.kind(new) != tUnform || p.unformed.Lookup(new) < 0) {
		w.commandOut(int(c.who), "%s is not a valid unformed noble entity.", w.boxCode(new))
		new = 0
	}
	if new == 0 && len(p.unformed) > 0 {
		new = p.unformed[0]
	}
	if new == 0 {
		w.commandOut(int(c.who), "No further nobles may be formed this turn.")
		return 0
	}
	name := "New noble"
	if c.parse.Len() >= 3 && c.parse[2] != "" {
		name = c.parse[2]
	}
	w.formNewNoble(int(c.who), name, new)
	p.unformed.RemValue(new)
	w.deductNP(pl, 1)
	return 1
}

// Complete c1.c:904–944; unformed displays at most five IDs, but names the
// full queue length. nextNPTurn owns its original allocation and schedule.
func (w *world) printHiringStatus(pl int) {
	assert(w.kind(pl) == tPlayer, "olympia/c1.c", 909, "kind(pl) == T_player")
	if w.subkind(pl) != subPlRegular {
		return
	}
	n := w.nextNPTurn(pl)
	if n == 0 {
		n += numMonths
	}
	w.commandOut(pl, "The next NP will be received at the end of turn %d.", int(w.sysclock.turn)+n)
}

func (w *world) printUnformed(pl int) {
	p := w.rpPlayer(pl)
	if p == nil || len(p.unformed) < 1 {
		return
	}
	buf := ""
	for i, n := range p.unformed {
		if i >= 5 {
			break
		}
		buf += " " + boxCodeLess(n)
	}
	w.commandOut(pl, "")
	w.commandOut(pl, "The next %s nobles formed will be: %s", niceNum(len(p.unformed)), buf)
	w.commandOut(pl, "")
}

// nextNPTurn is c1.c:881–901. The original first_turn schedule is disabled:
// every regular faction receives NP on turns 8, 16, ... . Even a query
// allocates a missing player record, as p_player does in C.
func (w *world) nextNPTurn(pl int) int {
	w.pPlayer(pl)
	ct := 7 - (int(w.sysclock.turn)+1)%numMonths
	ft := 2
	return (ft + ct) % numMonths
}

// Empty results represent C's NULL; every actual diagnostic/completion string
// is nonempty. The signal command remains unported.
var waitTags = []string{
	"time", "day", "unit", "gold", "item", "flag", "loc", "stack", "top",
	"ferry", "ship", "rain", "fog", "wind", "not", "owner", "raining",
	"foggy", "windy", "clear", "shiploc", "month", "turn",
}

type flagEnt struct {
	who  int32
	flag string
}

func (w *world) flagRaised(who int, flag string) int {
	for i, p := range w.flags {
		if who != 0 && w.player(int(p.who)) != who && int(p.who) != who {
			continue
		}
		if iStrcmp(p.flag, flag) == 0 {
			return i
		}
	}
	return -1
}

func clearWaitParse(c *command) {
	for i := range c.waitParse {
		c.waitParse[i] = nil
	}
	c.waitParse.Clear()
}

func (w *world) vWait(c *command) int {
	if c.parse.Len()-1 < 1 {
		w.commandOut(int(c.who), "Must say what condition to wait for.")
		return 0
	}
	clearWaitParse(c)
	if s := w.checkWaitConditions(c); s != "" {
		w.commandOut(int(c.who), "Wait finished: %s", s)
		c.wait = 0
		c.inhibitFinish = 1
		return 1
	}
	w.waitList.Append(int(c.who))
	return 1
}

func (w *world) dWait(c *command) int {
	if s := w.checkWaitConditions(c); s != "" {
		w.commandOut(int(c.who), "Wait finished: %s", s)
		w.waitList.RemValue(int(c.who))
		c.wait = 0
		c.inhibitFinish = 1
	}
	return 1
}

func (w *world) iWait(c *command) int {
	w.waitList.RemValue(int(c.who))
	return 1
}

func (w *world) parseWaitArgs(c *command) string {
	assert(c.waitParse.Len() == 0, "olympia/c1.c", 1236, "wait_args_len(c->wait_parse) == 0")
	i := 1
	for i < c.parse.Len() {
		tagS := c.parse[i]
		tag := lookup(waitTags, tagS)
		switch tag {
		case 16, 17, 18:
			tag -= 5
		}
		if tag < 0 {
			return sout("Unknown condition '%s'.", tagS)
		}
		i++
		p := &waitArg{tag: int32(tag)}
		c.waitParse.Append(p)
		switch tag {
		case 0, 1, 21, 22, 2, 3, 6, 20, 7, 9, 10:
			if i >= c.parse.Len() {
				return sout("Argument missing for '%s'.", tagS)
			}
			p.a1 = int32(w.parseArg(int(c.who), c.parse[i]))
			i++
		case 4:
			if i >= c.parse.Len() {
				return sout("Argument missing for '%s'.", tagS)
			}
			p.a1 = int32(w.parseArg(int(c.who), c.parse[i]))
			i++
			if i < c.parse.Len() {
				p.a2 = int32(w.parseArg(int(c.who), c.parse[i]))
				i++
			} else {
				p.a2 = 1
			}
		case 5:
			if i >= c.parse.Len() {
				return "Flag missing."
			}
			p.flag = &c.parse[i]
			i++
			p.a1 = int32(w.player(int(c.who)))
			if i < c.parse.Len() && ((len(c.parse[i]) > 0 && cIsDigit(c.parse[i][0])) || w.parseArg(int(c.who), c.parse[i]) != 0) {
				p.a1 = int32(w.parseArg(int(c.who), c.parse[i]))
				i++
			}
		case 8, 11, 12, 13, 14, 15, 19:
		default:
			fmt.Fprintf(w.stderr, "tag is %d\n", tag)
			assert(false, "olympia/c1.c", 1319, "FALSE")
		}
	}
	return ""
}

func (w *world) waitLeader(who int) int {
	if w.kind(who) != tChar {
		return who
	}
	return w.stackLeader(who)
}

// isShipNotdone is the oly.h subkind macro, including unfinished rafts.
func (w *world) isShipNotdone(n int) bool {
	return w.subkind(n) == subGalleyNotdone || w.subkind(n) == subRoundshipNotdone || w.subkind(n) == subRaftNotdone
}

func (w *world) checkWaitConditions(c *command) string {
	whereShip := w.subloc(int(c.who))
	if w.isShip(whereShip) || w.isShipNotdone(whereShip) {
		whereShip = w.subloc(whereShip)
	}
	if c.waitParse.Len() < 1 {
		if ret := w.parseWaitArgs(c); ret != "" {
			return ret
		}
		assert(c.waitParse.Len() > 0, "olympia/c1.c", 1363, "wait_args_len(c->wait_parse) > 0")
	}
	// C uses `not - cond`, not conjunction: the first true condition wins.
	// NOT lasts one condition; repeated NOT does not toggle it back off.
	not, setnot := false, false
	for _, p := range c.waitParse {
		if setnot {
			setnot = false
		} else if not {
			not = false
		}
		neg := ""
		if not {
			neg = " not"
		}
		a1, a2, who := int(p.a1), int(p.a2), int(c.who)
		switch p.tag {
		case 0:
			if (c.daysExecuting >= p.a1) != not {
				suffix := "s have"
				if p.a1 == 1 {
					suffix = " has"
				}
				return sout("%s day%s%s passed.", niceNum(a1), suffix, neg)
			}
		case 1:
			if (int32(w.sysclock.day) >= p.a1) != not {
				if not {
					return sout("today is not day %d.", a1)
				}
				return sout("today is day %d.", w.sysclock.day)
			}
		case 21, 22:
			if (int32(w.sysclock.turn) >= p.a1) != not {
				if not {
					return sout("it is not turn %d.", a1)
				}
				return sout("it is turn %d.", w.sysclock.turn)
			}
		case 2:
			if !w.validBox(a1) {
				return sout("%s does not exist.", w.boxCode(a1))
			}
			if w.charHere(who, a1) != not {
				return sout("%s is%s here.", w.boxCode(a1), neg)
			}
		case 3:
			if (w.hasItem(who, itemGold) >= p.a1) != not {
				if not {
					return sout("%s doesn't have %s.", w.justName(who), goldS(a1))
				}
				return sout("%s has %s.", w.justName(who), goldS(int(w.hasItem(who, itemGold))))
			}
		case 4:
			if (w.kind(a1) == tItem && w.hasItem(who, a1) >= p.a2) != not {
				if not {
					return sout("%s doesn't have %s.", w.justName(who), w.justNameQty(a1, a2))
				}
				return sout("%s has %s.", w.justName(who), w.justNameQty(a1, int(w.hasItem(who, a1))))
			}
		case 5:
			if a1 != 0 && !w.validBox(a1) {
				return sout("%s does not exist.", w.boxCode(a1))
			}
			j := w.flagRaised(a1, *p.flag)
			if not {
				if j < 0 {
					return "received no signal"
				}
			} else if j >= 0 {
				return sout("%s signaled '%s'", w.boxName(int(w.flags[j].who)), w.flags[j].flag)
			}
		case 6:
			if w.kind(a1) != tLoc && w.kind(a1) != tShip {
				return sout("%s is not a location or ship.", w.boxCode(a1))
			}
			if (w.subloc(who) == a1) != not {
				prefix := ""
				if not {
					prefix = "not "
				}
				return sout("%sat %s.", prefix, w.boxName(a1))
			}
		case 20:
			ship := w.subloc(who)
			if !w.isShip(ship) && !w.isShipNotdone(ship) {
				return sout("%s is not on a ship.", w.boxName(who))
			}
			if w.kind(a1) != tLoc && w.kind(a1) != tShip {
				return sout("%s is not a location or ship.", w.boxCode(a1))
			}
			if (w.subloc(ship) == a1) != not {
				prefix := ""
				if not {
					prefix = "not "
				}
				return sout("%sat %s.", prefix, w.boxName(a1))
			}
		case 7:
			if w.kind(a1) != tChar {
				break
			}
			if (w.waitLeader(who) == w.stackLeader(a1)) != not {
				return sout("%s is%s stacked with us.", w.boxName(a1), neg)
			}
		case 8:
			if (w.waitLeader(who) == who) != not {
				return sout("we are%s the stack leader", neg)
			}
		case 9:
			if !w.validBox(a1) || !w.isShip(a1) {
				return sout("%s is not a ship", w.boxCode(a1))
			}
			var horn int8
			if m := w.rpMagic(a1); m != nil {
				horn = m.ferryFlag
			}
			if (w.subloc(a1) == w.subloc(who) && horn != 0) != not {
				return sout("the ferry has%s signaled.", neg)
			}
		case 10:
			if w.kind(a1) != tShip {
				return sout("%s is not a ship.", w.boxCode(a1))
			}
			if (whereShip == w.subloc(a1)) != not {
				return sout("%s is%s here.", w.boxCode(a1), neg)
			}
		case 11:
			if (w.weatherHere(w.province(who), subRain) > 0) != not {
				return sout("it is%s raining.", neg)
			}
		case 12:
			if (w.weatherHere(w.province(who), subFog) > 0) != not {
				return sout("it is%s foggy.", neg)
			}
		case 13:
			if (w.weatherHere(w.province(who), subWind) > 0) != not {
				return sout("it is%s windy.", neg)
			}
		case 14:
			not, setnot = true, true
		case 15:
			if (w.firstCharacter(w.subloc(who)) == who) != not {
				return sout("we are%s the first character here", neg)
			}
		case 19:
			where := w.subloc(who)
			if (w.weatherHere(where, subFog) == 0 && w.weatherHere(where, subRain) == 0 && w.weatherHere(where, subWind) == 0) != not {
				return sout("it is%s clear.", neg)
			}
		default:
			fmt.Fprintf(w.stderr, "tag is %d\n", p.tag)
			assert(false, "olympia/c1.c", 1628, "FALSE")
		}
	}
	return ""
}
