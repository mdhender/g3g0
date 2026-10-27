// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

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
