// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/mdhender/g3g0/internal/lists"
)

// TestHelperMatchesC runs testdata/helper-cases.txt the way the C driver in
// tools/helper-compat/run.sh does and requires its output, line for line.
// testdata/helper-c-output.txt is that driver's output from the pinned G3
// sources; regenerate it only with run.sh --update.
func TestHelperMatchesC(t *testing.T) {
	cases := readTestdata(t, "helper-cases.txt")
	want := readTestdata(t, "helper-c-output.txt")

	var sections [][]string
	for line := range strings.SplitSeq(strings.TrimSuffix(cases, "\n"), "\n") {
		switch {
		case line == "" || line[0] == '#':
		case strings.HasPrefix(line, "section "):
			sections = append(sections, []string{line})
		case len(sections) > 0:
			sections[len(sections)-1] = append(sections[len(sections)-1], line)
		}
	}

	var out strings.Builder
	for _, s := range sections {
		out.WriteString(runHelperSection(t, s))
	}

	wantLines := strings.Split(want, "\n")
	gotLines := strings.Split(out.String(), "\n")
	for i := 0; i < len(wantLines) || i < len(gotLines); i++ {
		var w, g string
		if i < len(wantLines) {
			w = wantLines[i]
		}
		if i < len(gotLines) {
			g = gotLines[i]
		}
		if w != g {
			t.Fatalf("line %d:\n got %q\nwant %q", i+1, g, w)
		}
	}
}

// runHelperSection runs one section in a fresh world, as one run of the C
// driver, and returns what the driver prints, with the "exit STATUS" line
// run.sh adds.
func runHelperSection(t *testing.T, lines []string) (output string) {
	var out strings.Builder
	w := newWorld()
	w.stderr = &out // the driver sends stderr to stdout

	out.WriteString(lines[0] + "\n")
	defer func() {
		switch r := recover().(type) {
		case nil:
			out.WriteString("exit 0\n")
		case assertion:
			out.WriteString(r.Error())
			out.WriteString("exit 134\n")
		default:
			panic(r)
		}
		output = out.String()
	}()

	for _, line := range lines[1:] {
		fmt.Fprintf(&out, "> %s\n", line)
		c := &caseLine{rest: line}
		op := c.word()
		switch op {
		// setup
		case "alloc":
			n, k, sk := c.num(t), c.num(t), c.num(t)
			w.allocBox(n, k, sk)
		case "delete":
			w.deleteBox(c.num(t))
		case "setname":
			n := c.num(t)
			w.setName(n, c.text())
		case "where":
			n := c.num(t)
			w.pLocInfo(n).where = int32(c.num(t))
		case "here":
			n := c.num(t)
			w.pLocInfo(n).hereList.Append(c.num(t))
		case "lord":
			n := c.num(t)
			w.pChar(n).unitLord = int32(c.num(t))
		case "prevlord":
			n := c.num(t)
			w.pChar(n).prevLord = int32(c.num(t))
		case "units":
			n := c.num(t)
			w.pPlayer(n).units.Append(c.num(t))
		case "unititem":
			n := c.num(t)
			w.pChar(n).unitItem = int32(c.num(t))
		case "plural":
			n := c.num(t)
			w.pItem(n).pluralName = strSave(c.text())
		case "whohas":
			n := c.num(t)
			w.pItem(n).whoHas = int32(c.num(t))
		case "req":
			n := c.num(t)
			w.pSkill(n).requiredSkill = int32(c.num(t))
		case "learn":
			n := c.num(t)
			w.pSkill(n).timeToLearn = int32(c.num(t))
		case "moving":
			n := c.num(t)
			w.pChar(n).moving = int32(c.num(t))
		case "prisoner":
			n := c.num(t)
			w.pChar(n).prisoner = int8(c.num(t))
		case "cmd":
			n := c.num(t)
			cmd := w.pCommand(n)
			cmd.who = int32(c.num(t))
			cmd.state = int8(c.num(t))
			cmd.pri = int8(c.num(t))
		case "order":
			pl, who := c.num(t), c.num(t)
			helperOrderHead(w, pl, who).l.Append(c.text())
		case "turn":
			w.sysclock.turn = int16(c.num(t))
		case "temp":
			n := c.num(t)
			w.bx[n].temp = int32(c.num(t))
		case "inv":
			n := c.num(t)
			item, qty := c.num(t), c.num(t)
			w.bx[n].items.Append(&itemEnt{item: int32(item), qty: int32(qty)})

		// queries and operations under test
		case "names":
			n := c.num(t)
			fmt.Fprintf(&out, "display_name [%s]\n", esc(w.displayName(n)))
			fmt.Fprintf(&out, "box_name [%s]\n", esc(w.boxName(n)))
			fmt.Fprintf(&out, "just_name [%s]\n", esc(w.justName(n)))
		case "pluralname":
			n := c.num(t)
			fmt.Fprintf(&out, "plural_item_name [%s]\n", esc(w.pluralItemName(n, c.num(t))))
		case "cap":
			fmt.Fprintf(&out, "cap [%s]\n", esc(capitalize(c.text())))
		case "depth":
			fmt.Fprintf(&out, "loc_depth %d\n", w.locDepth(c.num(t)))
		case "inhere":
			n := c.num(t)
			fmt.Fprintf(&out, "in_here_list %d\n", b2i(w.inHereList(n, c.num(t))))
		case "addhere":
			n := c.num(t)
			w.addToHereList(n, c.num(t))
		case "remhere":
			n := c.num(t)
			w.removeFromHereList(n, c.num(t))
		case "parent":
			fmt.Fprintf(&out, "stack_parent %d\n", w.stackParent(c.num(t)))
		case "leader":
			fmt.Fprintf(&out, "stack_leader %d\n", w.stackLeader(c.num(t)))
		case "school":
			fmt.Fprintf(&out, "skill_school %d\n", w.skillSchool(c.num(t)))
		case "isunit":
			n := c.num(t)
			fmt.Fprintf(&out, "is_unit %d\n", b2i(w.isUnit(n, c.num(t))))
		case "player":
			fmt.Fprintf(&out, "player %d\n", w.player(c.num(t)))
		case "subloc":
			fmt.Fprintf(&out, "subloc %d\n", w.subloc(c.num(t)))
		case "viewloc":
			fmt.Fprintf(&out, "viewloc %d\n", w.viewloc(c.num(t)))
		case "firstchar":
			fmt.Fprintf(&out, "first_character %d\n", w.firstCharacter(c.num(t)))
		case "garrison":
			fmt.Fprintf(&out, "garrison_here %d\n", w.garrisonHere(c.num(t)))
		case "parsearg":
			n := c.num(t)
			fmt.Fprintf(&out, "parse_arg %d\n", w.parseArg(n, c.text()))
		case "setlord":
			n, l, k, lev := c.num(t), c.num(t), c.num(t), c.num(t)
			w.setLord(n, l, k, lev)
		case "restore":
			w.restoreStackActions(c.num(t))
		case "touch":
			w.touchLoc(c.num(t))
		case "initload":
			w.initLoadSup(c.num(t))
		case "clear":
			w.clearTemps(c.num(t))
		case "macros":
			n := c.num(t)
			fmt.Fprintf(&out, "noble_item %d char_moving %d is_prisoner %d item_unique %d learn_time %d req_skill %d\n",
				w.nobleItem(n), w.charMoving(n), helperPrisoner(w, n),
				w.itemUnique(n), w.learnTime(n), w.reqSkill(n))
		case "sort":
			var v lists.IList
			for {
				s := c.word()
				if s == "" {
					break
				}
				n, err := strconv.ParseInt(s, 0, 64)
				if err != nil {
					t.Fatalf("bad number %q", s)
				}
				v.Append(int(n))
			}
			safeQsort(v, intComp)
			helperList(&out, "sorted", v)
		case "loop":
			helperLoop(t, w, c, &out)
		case "dump":
			helperDump(w, c.num(t), &out)
		case "queues":
			for k := range maxPri {
				helperList(&out, fmt.Sprintf("load_q %d", k), w.loadQ[k])
			}
			helperList(&out, "run_q", w.runQ)
		default:
			t.Fatalf("unknown operation %q", line)
		}
	}
	return
}

// helperOrderHead is p_order_head() in order.c, which the driver uses to
// queue orders: player pl's order queue for unit who, created at the end
// if it has none.
func helperOrderHead(w *world, pl, who int) *orderList {
	p := w.pPlayer(pl)

	for _, o := range p.orders {
		if int(o.unit) == who {
			return o
		}
	}

	o := &orderList{unit: int32(who)}
	p.orders.Append(o)
	return o
}

// helperPrisoner is is_prisoner() as an int: the prisoner flag itself.
func helperPrisoner(w *world, n int) int {
	if p := w.rpChar(n); p != nil {
		return int(p.prisoner)
	}
	return 0
}

// helperList is the driver's list().
func helperList(out *strings.Builder, label string, l lists.IList) {
	fmt.Fprintf(out, "%s:", label)
	for _, i := range l {
		fmt.Fprintf(out, " %d", i)
	}
	out.WriteString("\n")
}

// helperDump is the driver's dump().
func helperDump(w *world, n int, out *strings.Builder) {
	b := w.bx[n]
	if b == nil {
		fmt.Fprintf(out, "box %d: none\n", n)
		return
	}
	fmt.Fprintf(out, "box %d: kind %d sub %d where %d temp %d\n", n, b.kind, b.skind, b.xLocInfo.where, b.temp)
	helperList(out, "  here", b.xLocInfo.hereList)
	if p := b.xChar; p != nil {
		fmt.Fprintf(out, "  char: lord %d prev %d loy %d/%d moving %d prisoner %d item %d\n",
			p.unitLord, p.prevLord, p.loyKind, p.loyRate, p.moving, p.prisoner, p.unitItem)
	}
	if p := b.xPlayer; p != nil {
		fmt.Fprintf(out, "  player: last_order_turn %d\n", p.lastOrderTurn)
		helperList(out, "  units", p.units)
		helperList(out, "  locs", p.locs)
		for _, o := range p.orders {
			fmt.Fprintf(out, "  orders %d:", o.unit)
			for _, s := range o.l {
				fmt.Fprintf(out, " [%s]", esc(s))
			}
			out.WriteString("\n")
		}
	}
	if c := b.cmd; c != nil {
		fmt.Fprintf(out, "  cmd: who %d state %d pri %d wait %d poll %d days %d cmd %d cond %d fuzzy %d\n",
			c.who, c.state, c.pri, c.wait, c.poll, c.daysExecuting, c.cmd, c.conditional, c.fuzzy)
		line := "(null)"
		if c.line != nil {
			line = esc(*c.line)
		}
		fmt.Fprintf(out, "  args: %d %d %d %d %d %d %d %d line [%s]\n",
			c.a, c.b, c.c, c.d, c.e, c.f, c.g, c.h, line)
	}
}

// helperLoop is the driver's loop operation.
func helperLoop(t *testing.T, w *world, c *caseLine, out *strings.Builder) {
	what := c.word()
	n := c.num(t)

	out.WriteString("visit:")
	switch what {
	case "here":
		for i := range w.loopHere(n) {
			fmt.Fprintf(out, " %d", i)
		}
	case "stack":
		for i := range w.loopStack(n) {
			fmt.Fprintf(out, " %d", i)
		}
	case "inv":
		for e := range w.loopInv(n) {
			fmt.Fprintf(out, " %d/%d", e.item, e.qty)
		}
	case "kind":
		for i := range w.loopKind(n) {
			fmt.Fprintf(out, " %d", i)
		}
	case "units":
		for i := range w.loopUnits(n) {
			fmt.Fprintf(out, " %d", i)
		}
	}
	out.WriteString("\n")
}
