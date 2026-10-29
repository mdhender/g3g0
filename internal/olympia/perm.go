// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"strings"

	"github.com/mdhender/g3g0/internal/lists"
)

// printAdmit is perm.c:141–223. Sorting mutates the saved admission list,
// but neither invalid references nor the original member order are changed.
func (w *world) printAdmit(pl int) {
	assert(w.kind(pl) == tPlayer, "olympia/perm.c", 197, "kind(pl) == T_player")
	p := w.pPlayer(pl)
	safeQsort(p.admits, func(a, b *admit) int { return int(a.targ - b.targ) })
	first := true
	for _, a := range p.admits {
		if !w.validBox(int(a.targ)) {
			continue
		}
		if first {
			w.commandOut(pl, "")
			w.commandOut(pl, "Admit permissions:")
			w.commandOut(pl, "")
			w.indent += 3
			first = false
		}
		w.printAdmitSup(pl, a)
	}
	if !first {
		w.indent -= 3
	}
}

func (w *world) printAdmitSup(pl int, p *admit) {
	buf := sout("admit %4s", boxCodeLess(int(p.targ)))
	count := 0
	if p.sense != 0 {
		buf += "  all"
		count++
	}
	for _, n := range p.l {
		if !w.validBox(n) {
			continue
		}
		count++
		if count >= 12 {
			w.commandOut(pl, "%s", buf)
			buf = "          "
			count = 1
		}
		buf += sout(" %4s", boxCodeLess(n))
	}
	if count != 0 {
		w.commandOut(pl, "%s", buf)
	}
}

// printAtt is complete perm.c:426–488. Each original list is sorted only
// when its renderer is reached; no disposition record is allocated.
func (w *world) printAtt(who, n int) {
	p := w.rpDisp(n)
	if p == nil {
		return
	}
	first := true
	w.printAttSup(who, p.hostile, "hostile", &first)
	w.printAttSup(who, p.neutral, "neutral", &first)
	w.printAttSup(who, p.defend, "defend ", &first)
	if !first {
		w.indent -= 3
	}
}

func (w *world) printAttSup(who int, l lists.IList, header string, first *bool) {
	if len(l) == 0 {
		return
	}
	buf, count := header, 0
	safeQsort(l, intComp)
	for _, n := range l {
		if !w.validBox(n) {
			continue
		}
		if *first {
			w.commandOut(who, "")
			w.commandOut(who, "Declared attitudes:")
			w.indent += 3
			*first = false
		}
		count++
		if count >= 12 {
			w.commandOut(who, "%s", buf)
			buf = strings.Repeat(" ", len(header))
			count = 1
		}
		buf += sout(" %4s", boxCodeLess(n))
	}
	if count != 0 {
		w.commandOut(who, "%s", buf)
	}
}

// rpAdmit deliberately uses pPlayer, as does C: a lookup may allocate it.
func (w *world) rpAdmit(pl, targ int) *admit {
	assert(w.kind(pl) == tPlayer, "olympia/perm.c", 16, "kind(pl) == T_player")
	for _, p := range w.pPlayer(pl).admits {
		if int(p.targ) == targ {
			return p
		}
	}
	return nil
}

func (w *world) pAdmit(pl, targ int) *admit {
	assert(w.kind(pl) == tPlayer, "olympia/perm.c", 34, "kind(pl) == T_player")
	p := w.pPlayer(pl)
	for _, a := range p.admits {
		if int(a.targ) == targ {
			return a
		}
	}
	a := &admit{targ: int32(targ)}
	p.admits.Append(a)
	return a
}

func (w *world) vAdmit(c *command) int {
	targ, pl := int(c.a), w.player(int(c.who))
	if !w.validBox(targ) {
		w.commandOut(int(c.who), "Must specify an entity for admit.")
		return 0
	}
	w.cmdShift(c)
	p := w.pAdmit(pl, targ)
	if p.flag == 0 {
		p.sense = 0
		p.l.Clear()
		p.flag = 1
	}
	for c.parse.Len() > 1 {
		if iStrcmp(c.parse[1], "all") == 0 {
			p.sense = 1
		} else if k := w.kind(int(c.a)); k == tChar || k == tPlayer || k == tUnform {
			p.l.Append(int(c.a))
		} else {
			w.commandOut(int(c.who), "%s isn't a valid entity to admit.", c.parse[1])
		}
		w.cmdShift(c)
	}
	return 1
}

// setAtt uses ilist_rem_value's complete duplicate removal and sorts only
// the selected destination, leaving the other two lists in their old order.
func (w *world) setAtt(who, targ, disp int) {
	p := w.pDisp(who)
	p.neutral.RemValue(targ)
	p.hostile.RemValue(targ)
	p.defend.RemValue(targ)
	var l *lists.IList
	switch disp {
	case attNeutral:
		l = &p.neutral
	case attHostile:
		l = &p.hostile
	case attDefend:
		l = &p.defend
	case attNone:
		return
	default:
		assert(false, "olympia/perm.c", 273, "FALSE")
	}
	l.Append(targ)
	safeQsort(*l, intComp)
}

func (w *world) isDefend(who, targ int) bool {
	if w.isHostile(who, targ) {
		return false
	}
	if w.defaultGarrison(who) != 0 {
		return true
	}
	pl := w.player(who)
	for _, n := range []int{who, pl} {
		if p := w.rpDisp(n); p != nil {
			if p.defend.Lookup(targ) >= 0 {
				return true
			}
			if p.neutral.Lookup(targ) >= 0 {
				return false
			}
			if p.defend.Lookup(w.player(targ)) >= 0 {
				return true
			}
			if p.neutral.Lookup(w.player(targ)) >= 0 {
				return false
			}
		}
	}
	return pl == w.player(targ) && pl != int(w.indepPlayer) && w.cloakLord(who) == 0
}

func (w *world) vSetAtt(c *command, k int) int {
	for c.parse.Len() > 1 {
		if !w.validBox(int(c.a)) {
			w.commandOut(int(c.who), "%s is not a valid entity.", c.parse[1])
		} else if k == attHostile && w.player(int(c.who)) == w.player(int(c.a)) && w.player(int(c.who)) != int(w.indepPlayer) {
			w.commandOut(int(c.who), "Can't be hostile to a unit in the same faction.")
		} else {
			w.setAtt(int(c.who), int(c.a), k)
		}
		w.cmdShift(c)
	}
	return 1
}

func (w *world) vHostile(c *command) int  { return w.vSetAtt(c, attHostile) }
func (w *world) vDefend(c *command) int   { return w.vSetAtt(c, attDefend) }
func (w *world) vNeutral(c *command) int  { return w.vSetAtt(c, attNeutral) }
func (w *world) vAttClear(c *command) int { return w.vSetAtt(c, attNone) }

func (w *world) defaultGarrison(n int) int8 {
	if p := w.rpMagic(n); p != nil {
		return p.defaultGarr
	}
	return 0
}

func (w *world) willAdmit(pl, who, targ int) bool {
	if w.defaultGarrison(targ) != 0 {
		return true
	}
	pl = w.player(pl)
	if w.player(who) == pl {
		return true
	}
	p := w.rpAdmit(pl, targ)
	if p == nil {
		return false
	}
	found := p.l.Lookup(who) >= 0 || p.l.Lookup(w.player(who)) >= 0
	if p.sense != 0 {
		return !found
	}
	return found
}

// isHostile is perm.c's unit-target predicate. In particular, a hostile
// disposition naming only the target's faction does not match here.
func (w *world) isHostile(who, targ int) bool {
	if w.player(who) == w.player(targ) {
		return false
	}
	if w.npcProgram(who) == progSublocMonster && !w.isNPC(targ) && w.onlyDefeatable(who) == 0 {
		return true
	}
	if w.subkind(who) == subGarrison {
		if p := w.rpMisc(who); p != nil && p.garrHost.Lookup(targ) >= 0 {
			return true
		}
	}
	if p := w.rpDisp(who); p != nil && p.hostile.Lookup(targ) >= 0 {
		return true
	}
	if p := w.rpDisp(w.player(who)); p != nil && p.hostile.Lookup(targ) >= 0 {
		return true
	}
	return false
}

func (w *world) clearAllAtt(who int) {
	if p := w.rpDisp(who); p != nil {
		p.neutral.Clear()
		p.hostile.Clear()
		p.defend.Clear()
	}
}
