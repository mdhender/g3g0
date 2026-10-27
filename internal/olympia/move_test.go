// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Expected bytes come only from tools/move-compat's real pinned C helpers.
func TestMoveMatchesC(t *testing.T) {
	var sections [][]string
	for line := range strings.SplitSeq(readTestdata(t, "move-cases.txt"), "\n") {
		switch {
		case line == "" || strings.HasPrefix(line, "#"):
		case strings.HasPrefix(line, "section "):
			sections = append(sections, []string{line})
		default:
			sections[len(sections)-1] = append(sections[len(sections)-1], line)
		}
	}
	var out strings.Builder
	for _, section := range sections {
		out.WriteString(runMoveSection(t, section))
	}
	want := readTestdata(t, "move-c-output.txt")
	gotLines, wantLines := strings.Split(out.String(), "\n"), strings.Split(want, "\n")
	if len(gotLines) != len(wantLines) {
		t.Fatalf("line count = %d, want %d (including final empty lines)", len(gotLines), len(wantLines))
	}
	for i := range wantLines {
		if wantLines[i] != gotLines[i] {
			t.Fatalf("line %d:\n got %q\nwant %q", i+1, gotLines[i], wantLines[i])
		}
	}
}

func runMoveSection(t *testing.T, lines []string) (output string) {
	var stdout, stderr strings.Builder
	w := newWorld()
	w.stdout, w.stderr = &stdout, &stderr
	w.libdir = t.TempDir()
	if err := os.Mkdir(filepath.Join(w.libdir, "log"), 0o755); err != nil {
		t.Fatal(err)
	}
	status := 0
	defer func() {
		switch r := recover().(type) {
		case nil:
		case assertion:
			stderr.WriteString(r.Error())
			status = 134
		case runtime.Error:
			// Only the separately C-recorded, known unowned-character NULL
			// dereference may discard buffered output (G3 #353).
			if lines[0] != "section unowned hidden route G3 353" {
				panic(r)
			}
			stdout.Reset()
			stderr.Reset()
			status = 139
		default:
			panic(r)
		}
		output = fmt.Sprintf("stdout:\n%sstderr:\n%s\nexit %d\n", stdout.String(), stderr.String(), status)
	}()
	dump := func(l []*exitView) {
		fmt.Fprintf(&stdout, "exits %d\n", len(l))
		for _, v := range l {
			fmt.Fprintf(&stdout, "exit %d %d %d %d %d %d %d %d %d %d %d %d %d %d\n",
				v.direction, v.destination, v.orig, v.distance, v.impassable, v.destHidden,
				v.origHidden, v.hidden, v.inside, v.road, v.water, v.inTransit, v.magicBarrier, v.hadesCost)
		}
	}
	cmd, v := &command{}, &exitView{}
	fmt.Fprintln(&stdout, lines[0])
	for _, line := range lines[1:] {
		fmt.Fprintf(&stdout, "> %s\n", line)
		c := &caseLine{rest: line}
		op := c.word()
		if op == "argument" {
			cmd.parse = cstringsList{"move", strings.TrimSpace(c.rest)}
			continue
		}
		var a [4]int
		for i := range a {
			if strings.TrimSpace(c.rest) != "" {
				a[i] = c.num(t)
			}
		}
		n, b, sk, d := a[0], a[1], a[2], a[3]
		switch op {
		case "command":
			cmd.who, cmd.a, cmd.wait, cmd.cmd = int32(n), int32(b), int32(sk), int32(d)
		case "args":
			cmd.b, cmd.c, cmd.d, cmd.e = int32(n), int32(b), int32(sk), int32(d)
		case "tail":
			cmd.f, cmd.g, cmd.h, cmd.daysExecuting = int32(n), int32(b), int32(sk), int32(d)
		case "exit1":
			v.direction, v.destination, v.orig, v.distance = int32(n), int32(b), int32(sk), int32(d)
		case "exit2":
			v.impassable, v.destHidden, v.origHidden, v.hidden = int32(n), int32(b), int32(sk), int32(d)
		case "exit3":
			v.inside, v.road, v.water, v.inTransit = int32(n), int32(b), int32(sk), int32(d)
		case "exit4":
			v.magicBarrier, v.hadesCost = int32(n), int32(b)
		case "parse":
			var zeroArg *string
			if b != 0 {
				s := "move"
				zeroArg = &s
			}
			ret := w.parseExitDir(cmd, n, zeroArg)
			if ret == nil {
				fmt.Fprintln(&stdout, "no exit")
			} else {
				*v = *ret
				fmt.Fprintln(&stdout, "parsed")
			}
		case "check":
			fmt.Fprintf(&stdout, "check %d\n", b2i(w.landCheck(cmd, v, n != 0)))
		case "can":
			fmt.Fprintf(&stdout, "can %d\n", b2i(w.canMoveHere(n, cmd)))
		case "outer":
			fmt.Fprintf(&stdout, "outer %d\n", w.canMoveAtOuterLevel(n, cmd))
		case "save":
			saveVArray(cmd, v)
		case "restore":
			restoreVArray(cmd, v)
		case "dump_exit":
			dump([]*exitView{v})
		case "dump_command":
			fmt.Fprintf(&stdout, "command %d %d %d %d %d %d %d %d %d %d %d %d\n", cmd.who, cmd.a, cmd.b, cmd.c, cmd.d, cmd.e, cmd.f, cmd.g, cmd.h, cmd.wait, cmd.cmd, cmd.daysExecuting)
		case "see_all":
			w.immedSeeAll = int32(n)
		case "admission":
			w.pPlayer(n).admits = append(w.pPlayer(n).admits, &admit{targ: int32(b), sense: int32(sk)})
		case "entry":
			w.pPlayer(n).admits[b].l.Append(sk)
		case "admit":
			fmt.Fprintf(&stdout, "admit %d\n", b2i(w.willAdmit(n, b, sk)))
		case "default_garr":
			w.pMagic(n).defaultGarr = int8(b)
		case "guard":
			w.pChar(n).guard = int8(b)
		case "char_moving":
			w.pChar(n).moving = int32(b)
		case "swear_release":
			w.pMagic(n).swearOnRelease = int8(b)
		case "suspend":
			w.suspendStackActions(n)
		case "resume":
			w.restoreStackActions(n)
		case "clear_guard":
			w.clearGuardFlag(n)
		case "leave":
			w.leaveStack(n)
		case "unstack":
			w.unstack(n)
		case "movement_state":
			for n := range w.loopBoxes() {
				if p := w.rpChar(n); p != nil {
					fmt.Fprintf(&stdout, "action %d %d %d %d %d\n", n, p.moving, p.guard, p.loyKind, p.loyRate)
				}
				if p := w.rpMagic(n); p != nil {
					fmt.Fprintf(&stdout, "magic %d %d %d\n", n, p.defaultGarr, p.swearOnRelease)
				}
				if p := w.rpPlayer(n); p != nil {
					fmt.Fprintf(&stdout, "player %d admits %d\n", n, len(p.admits))
					for _, p := range p.admits {
						fmt.Fprintf(&stdout, "permission %d %d %d", p.targ, p.sense, p.flag)
						for _, j := range p.l {
							fmt.Fprintf(&stdout, " %d", j)
						}
						fmt.Fprintln(&stdout)
					}
				}
			}
		case "output":
			w.immediate = n != 0
		case "output_bit":
			fmt.Fprintf(&stdout, "output_bit %d\n", b2i(testBit(w.pPlayer(n).output, b)))
		case "report":
			fmt.Fprintf(&stdout, "report %d:\n", n)
			data, err := os.ReadFile(filepath.Join(w.libdir, "log", fmt.Sprint(n)))
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			stdout.Write(data)
		case "item":
			p := w.pItem(n)
			p.weight, p.landCap, p.rideCap = int16(b), int16(sk), int16(d)
		case "air":
			p := w.pItem(n)
			p.flyCap, p.animal, p.isManItem = int16(b), int8(sk), int8(d)
		case "fighter":
			p := w.pItem(n)
			p.attack, p.defense, p.missile = int16(b), int16(sk), int16(d)
		case "unit":
			p := w.pChar(n)
			p.unitItem, p.prisoner, p.loyKind = int32(b), int8(sk), int8(d)
		case "inv":
			var entry *itemEnt
			for _, e := range w.bx[n].items {
				if e.item == int32(b) {
					entry = e
					break
				}
			}
			if entry == nil {
				entry = &itemEnt{item: int32(b)}
				w.bx[n].items = append(w.bx[n].items, entry)
			}
			entry.qty = int32(sk)
		case "factor":
			w.armySlowFactor = int32(n)
		case "strength":
			w.pMisc(n).stormStr = int16(b)
		case "unit_weights", "stack_weights":
			var v weights
			if op == "unit_weights" {
				v = w.determineUnitWeights(n)
			} else {
				v = w.determineStackWeights(n)
			}
			fmt.Fprintf(&stdout, "weights %d %d %d %d %d %d %d %d\n", v.animals, v.totalWeight, v.landCap, v.landWeight, v.rideCap, v.rideWeight, v.flyCap, v.flyWeight)
		case "counts":
			fmt.Fprintf(&stdout, "counts %d %d\n", w.countStackMoveNobles(n), w.countStackFigures(n))
		case "land", "fly", "water":
			cmd, v := &command{who: int32(n)}, &exitView{destination: int32(b), distance: int32(sk)}
			var result int32
			switch op {
			case "land":
				result = w.moveExitLand(cmd, v, d != 0)
			case "fly":
				result = w.moveExitFly(cmd, v, d != 0)
			case "water":
				result = w.moveExitWater(cmd, v, b, d != 0)
			}
			fmt.Fprintf(&stdout, "%s %d\nunchanged %d %d %d %d\n", op, result, cmd.who, v.destination, v.distance, cmd.wait)
		case "state":
			fmt.Fprintf(&stdout, "state factor %d\n", w.armySlowFactor)
			for n := range w.loopBoxes() {
				fmt.Fprintf(&stdout, "box %d %d %d %d here", n, w.kind(n), w.subkind(n), w.loc(n))
				for _, i := range w.rpLocInfo(n).hereList {
					fmt.Fprintf(&stdout, " %d", i)
				}
				fmt.Fprintln(&stdout)
				if p := w.rpChar(n); p != nil {
					fmt.Fprintf(&stdout, "unit %d %d %d %d\n", p.unitItem, p.prisoner, p.loyKind, w.lord(n))
				}
				if p := w.rpItem(n); p != nil {
					fmt.Fprintf(&stdout, "item %d %d %d %d %d %d %d %d %d\n", p.weight, p.landCap, p.rideCap, p.flyCap, p.animal, p.isManItem, p.attack, p.defense, p.missile)
				}
				if p := w.rpMisc(n); p != nil {
					fmt.Fprintf(&stdout, "storm %d\n", p.stormStr)
				}
				for _, e := range w.bx[n].items {
					fmt.Fprintf(&stdout, "inv %d %d\n", e.item, e.qty)
				}
				if p := w.rpPlayer(n); p != nil {
					fmt.Fprint(&stdout, "output:")
					for i := 1; i < maxBoxes; i++ {
						if testBit(p.output, i) {
							fmt.Fprintf(&stdout, " %d", i)
						}
					}
					fmt.Fprintln(&stdout)
				}
			}
		case "box":
			w.allocBox(n, b, sk)
			w.pLocInfo(n).where = int32(d)
			if d != 0 {
				w.pLocInfo(d).hereList.Append(n)
			}
		case "dest":
			w.pLoc(n).provDest.Append(b)
		case "hidden":
			w.pLoc(n).hidden = int8(b)
		case "barrier":
			w.pLoc(n).barrier = int16(b)
		case "lane":
			w.pLoc(n).seaLane = int8(b)
		case "uldim":
			w.pSubloc(n).uldimFlag = int8(b)
		case "summer":
			w.pSubloc(n).summerFlag = int8(b)
		case "open":
			w.pSubloc(n).linkOpen = int8(b)
		case "from":
			w.pSubloc(n).linkFrom.Append(b)
		case "to":
			w.pSubloc(n).linkTo.Append(b)
		case "road":
			w.pGate(n).toLoc, w.pGate(n).roadHidden = int32(b), int8(sk)
		case "lord":
			w.pChar(n).unitLord = int32(b)
		case "known":
			setBit(&w.pPlayer(n).known, b)
		case "hades":
			w.hadesRegion = int32(n)
		case "clock":
			w.sysclock.turn, w.sysclock.daysSinceEpoch, w.evening, w.monthDone = int16(n), int32(b), int32(sk), int32(d)
		case "day":
			w.sysclock.day = int16(n)
		case "moving":
			w.pSubloc(n).moving = int32(b)
		case "master":
			seed := make([]byte, 16)
			for i := range seed {
				seed[i] = byte(n + i)
			}
			path := filepath.Join(t.TempDir(), "randseed")
			if err := os.WriteFile(path, seed, 0o600); err != nil {
				t.Fatal(err)
			}
			w.loadSeed(path)
		case "legacy":
			fmt.Fprintf(&stdout, "legacy %d\n", w.rnd.Rnd(n, b))
		case "direction":
			fmt.Fprintf(&stdout, "direction %d\n", w.locationDirection(n, b))
		case "distance":
			fmt.Fprintf(&stdout, "distance %d\n", w.exitDistance(n, b))
		case "port":
			fmt.Fprintf(&stdout, "port %d\n", b2i(w.isPortCity(n)))
		case "province_port":
			fmt.Fprintf(&stdout, "province_port %d\n", w.provinceHasPortCity(n))
		case "ocean":
			fmt.Fprintf(&stdout, "ocean %d\n", w.hasOceanAccess(n))
		case "exits":
			dump(w.exitsFromLoc(n, b))
		case "nsew":
			dump(w.exitsFromLocNSEW(n, b))
		case "select":
			dump(w.exitsFromLocNSEWSelect(n, b, sk, d))
		default:
			t.Fatalf("unknown move operation %q", op)
		}
	}
	return ""
}
