// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"github.com/mdhender/g3g0/internal/clib"
)

// Recordings come only from tools/command-compat's real pinned C input.c,
// table, handlers and output routines, never a model of the dispatcher.
func TestCommandMatchesC(t *testing.T) {
	var sections [][]string
	for line := range strings.SplitSeq(readTestdata(t, "command-cases.txt"), "\n") {
		if line == "" || line[0] == '#' {
			continue
		}
		if strings.HasPrefix(line, "section ") {
			sections = append(sections, []string{line})
		} else {
			sections[len(sections)-1] = append(sections[len(sections)-1], line)
		}
	}
	var out strings.Builder
	for _, section := range sections {
		t.Run(strings.TrimPrefix(section[0], "section "), func(t *testing.T) {
			out.WriteString(runCommandSection(t, section))
		})
	}
	want := strings.Split(readTestdata(t, "command-c-output.txt"), "\n")
	got := strings.Split(out.String(), "\n")
	if len(got) != len(want) {
		t.Fatalf("line count = %d, want %d (including final empty lines)", len(got), len(want))
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("line %d:\n got %q\nwant %q", i+1, got[i], want[i])
		}
	}
}

// Native directory order is not portable recording data. The tool supplies
// two fresh real-C runs over this same unchanged directory before Go reads it.
func TestCommandAdmissionNativeOrder(t *testing.T) {
	dir := os.Getenv("COMMAND_NATIVE_DIR")
	if dir == "" {
		t.Skip("live C expectation absent; run tools/command-compat/run.sh")
	}
	cases, err := os.ReadFile(filepath.Join(dir, "cases.txt"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join(dir, "out.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(want), "out-of-scope C dependency:") {
		t.Fatal("live C expectation reached a dependency marker")
	}
	got := runCommandSectionInDir(t, strings.Split(strings.TrimSuffix(string(cases), "\n"), "\n"), filepath.Join(dir, "input"))
	if got != string(want) {
		t.Fatalf("live raw C/Go mismatch:\n got %q\nwant %q", got, want)
	}
}

func runCommandSection(t *testing.T, lines []string) string {
	return runCommandSectionInDir(t, lines, "")
}

func runCommandSectionInDir(t *testing.T, lines []string, dir string) (output string) {
	var stdout, stderr strings.Builder
	w := newWorld()
	w.flushAlways = true // command C driver opens streams with flush_always=1.
	t.Cleanup(func() { w.immediate = false; w.closeLogfile(); w.closeTimes() })
	w.stdout, w.stderr = &stdout, &stderr
	w.libdir = dir
	if dir == "" {
		w.libdir = t.TempDir()
		if err := os.Mkdir(filepath.Join(w.libdir, "log"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(filepath.Join(w.libdir, "orders"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	w.allocBox(200, tPlayer, subPlRegular)
	w.allocBox(201, tPlayer, subPlRegular)
	w.pPlayer(200)
	w.pPlayer(201)
	w.allocBox(1000, tChar, 0)
	w.allocBox(1001, tChar, 0)
	w.allocBox(skHideSelf, tSkill, 0)
	w.pChar(1000).unitLord = 200
	w.pChar(1001).unitLord = 200
	c := w.pCommand(1000)
	c.who, c.state = 1000, stateDone
	created := 0
	savedVector := w.outVector
	var watched []string
	fileState := func(path string) {
		st, err := os.Stat(path)
		kind := "missing"
		if err == nil {
			switch {
			case st.IsDir():
				kind = "directory"
			case st.Size() != 0:
				kind = "nonempty"
			default:
				kind = "empty"
			}
		}
		fmt.Fprintf(&stdout, "path %s %s\n", path, kind)
	}
	waitDump := func() {
		fmt.Fprintf(&stdout, "wait args %d\n", c.waitParse.Len())
		for _, p := range c.waitParse {
			flag := "NULL"
			if p.flag != nil {
				flag = *p.flag
			}
			fmt.Fprintf(&stdout, "wait arg %d %d %d [%s]\n", p.tag, p.a1, p.a2, flag)
		}
	}
	fmt.Fprintln(&stdout, lines[0])
	defer func() {
		status := 0
		switch r := recover().(type) {
		case nil:
			if strings.HasPrefix(lines[0], "section combat boundary ") {
				t.Fatal("expected actual combat boundary, scan returned")
			}
		case notImplemented:
			if r != notImplemented("regular_combat") || !strings.HasPrefix(lines[0], "section combat boundary ") {
				panic(r)
			}
			// Observe the real production stop, never return fake combat success.
			fmt.Fprintln(&stdout, "combat log:")
			pl := w.player(int(c.who))
			data, err := os.ReadFile(filepath.Join(w.libdir, "log", fmt.Sprint(pl)))
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			stdout.Write(data)
			fmt.Fprintf(&stdout, "combat output %d\n", b2i(testBit(w.pPlayer(pl).output, int(c.who))))
			fmt.Fprintln(&stderr, "combat boundary: regular_combat")
			status = 134 // C harness abort, not a successful C combat recording.
		case assertion:
			fmt.Fprint(&stderr, r.Error())
			status = 134
		case exitStatus:
			status = int(r)
		case runtime.Error:
			if !strings.Contains(r.Error(), "nil pointer dereference") {
				panic(r)
			}
			// As in runEngine, C SIGSEGV does not flush buffered stdout.
			stdout.Reset()
			status = 139
		default:
			panic(r)
		}
		if status == 0 || status == 1 {
			for _, path := range watched {
				fileState(path) // C atexit: before stdio's exit flush.
			}
		}
		output = fmt.Sprintf("stdout:\n%sstderr:\n%s\nexit %d\n", stdout.String(), stderr.String(), status)
	}()
	for _, line := range lines[1:] {
		fmt.Fprintf(&stdout, "> %s\n", line)
		args := &caseLine{rest: line}
		switch op := args.word(); op {
		case "weather_map":
			count := args.num(t)
			w.allocBox(8000, tLoc, subRegion)
			for i := count - 1; i >= 0; i-- {
				n := 9000 + i
				w.allocBox(n, tLoc, 1+i%4)
				w.setWhere(n, 8000)
				setBit(&w.pPlayer(200).locs, n)
			}
		case "weather":
			switch action := args.word(); action {
			case "natural":
				w.naturalWeather()
			case "some":
				w.createSomeStorms(args.num(t), args.num(t))
			case "new":
				fmt.Fprintf(&stdout, "new storm %d\n", w.newStorm(args.num(t), args.num(t), args.num(t), args.num(t)))
			case "start":
				c.a = int32(args.num(t))
				fmt.Fprintf(&stdout, "bind start %d\n", w.vBindStorm(c))
			case "finish":
				c.a = int32(args.num(t))
				fmt.Fprintf(&stdout, "bind finish %d\n", w.dBindStorm(c))
			case "count":
				fmt.Fprintf(&stdout, "provinces %d\n", w.nprovinces())
			case "decay":
				w.stormDecay()
			case "monthly":
				w.stormMove()
			case "touch":
				w.stormOwnerTouchLoc()
			case "move":
				w.moveStorm(args.num(t), args.num(t))
			case "bound":
				w.moveBoundStorms(args.num(t), args.num(t))
			case "dissipate":
				w.dissipateStorm(args.num(t), args.num(t))
			case "pending":
				n := args.num(t)
				dir, dest := 0, 0
				if p := w.rpMisc(n); p != nil {
					dir, dest = int(p.npcDir), int(p.stormMove)
				}
				fmt.Fprintf(&stdout, "pending %d %d %d\n", n, dir, dest)
			case "state":
				for n := range w.loopKind(tStorm) {
					fmt.Fprintf(&stdout, "storm %d\n", n)
					w.saveBox(&stdout, n)
				}
				for n := range w.loopProvince() {
					fmt.Fprintf(&stdout, "province %d\n", n)
					helperList(&stdout, "here", w.pLocInfo(n).hereList)
				}
				fmt.Fprintf(&stdout, "show_garrison %d\n", b2i(w.showToGarrison))
				helperList(&stdout, "output", w.pPlayer(200).output)
			case "probe":
				n, a := args.num(t), args.num(t)
				l := make([]int, 11)
				for i := range l {
					l[i] = 300 + i*7
				}
				w.streams.WthrShuffle(l)
				helperList(&stdout, "weather shuffle", l)
				fmt.Fprintf(&stdout, "weather strength %d\n", w.streams.WthrStorm(n, a, 2, 3))
			default:
				t.Fatalf("unknown weather %s", action)
			}
		case "garrison":
			switch action := args.word(); action {
			case "ping":
				w.pingGarrisons()
			case "pay":
				w.defaultGarrisonPay()
			case "tax":
				w.garrisonGold()
			case "metric":
				w.goldProvinces = int32(args.num(t))
			case "subloc":
				n, sk := args.num(t), args.num(t)
				fmt.Fprintf(&stdout, "province subloc %d %d %d\n", n, sk, w.provinceSubloc(n, sk))
			case "taxstate":
				fmt.Fprintf(&stdout, "gold provinces %d\n", w.goldProvinces)
				for n := range w.loopSubkind(subGarrison) {
					p := w.rpMisc(n)
					tax, forward := int32(0), int32(0)
					if p != nil {
						tax, forward = p.garrTax, p.garrForward
					}
					fmt.Fprintf(&stdout, "tax garrison %d misc %d tax %d forward %d strength %d\n", n, b2i(p != nil), tax, forward, w.garrisonStrength(n))
				}
				for n := range w.loopKind(tLoc) {
					fmt.Fprintf(&stdout, "tax loc %d temp %d\n", n, w.bx[n].temp)
				}
			case "strength":
				fmt.Fprintf(&stdout, "garrison strength %d\n", w.garrisonStrength(args.num(t)))
			case "maint":
				fmt.Fprintf(&stdout, "maintenance %d\n", maintCost(args.num(t)))
			case "state":
				fmt.Fprintf(&stdout, "garrison pay %d gold %d show %d\n", w.garrisonPay, w.goldGarrison, b2i(w.showToGarrison))
				helperList(&stdout, "vector", w.outVector)
				for n := range w.loopKind(tChar) {
					guard := int8(0)
					if p := w.rpChar(n); p != nil {
						guard = p.guard
					}
					fmt.Fprintf(&stdout, "garrison char %d leader %d default %d prisoner %d npc %d guard %d gold %d\n",
						n, w.stackLeader(n), w.defaultGarrison(n), b2i(w.isPrisoner(n)), b2i(w.isNPC(n)), guard, w.hasItem(n, itemGold))
				}
				for n := range w.loopKind(tPlayer) {
					fmt.Fprintf(&stdout, "player %d\n", n)
					helperList(&stdout, "output", w.pPlayer(n).output)
				}
			default:
				t.Fatalf("unknown garrison %s", action)
			}
		case "npc":
			action, n := args.word(), args.num(t)
			switch action {
			case "unsworn":
				w.autoUnsworn(n)
			case "dispatch":
				w.queueNpcOrders()
			case "bandit":
				w.autoBandit(n, args.num(t))
			case "mob":
				w.autoMob(n)
			case "undead":
				w.autoUndead(n)
			case "begin":
				w.streams.BeginNPC(n)
			case "probe":
				fmt.Fprintf(&stdout, "behavior zero %d\n", w.streams.NPCBehavior(n, 0, 1, 2))
				fmt.Fprintf(&stdout, "behavior one %d\n", w.streams.NPCBehavior(n, 1, 1, 2))
				fmt.Fprintf(&stdout, "behavior wait %d\n", w.streams.NPCBehavior(n, 1, 10, 20))
				fmt.Fprintf(&stdout, "behavior wide %d\n", w.streams.NPCBehavior(n, 99, -100000, 100000))
			default:
				t.Fatalf("unknown npc %s", action)
			}
		case "savage":
			switch action := args.word(); action {
			case "init":
				w.initSavageAttacks()
				fmt.Fprintf(&stdout, "initial savages %d\n", w.numSavages)
			case "auto":
				w.autoSavage(args.num(t))
			case "hates":
				fmt.Fprintf(&stdout, "hates %d\n", b2i(w.savageHates(args.num(t))))
			case "here":
				fmt.Fprintf(&stdout, "hate here %d\n", w.savageHateHere(args.num(t)))
			case "humans":
				fmt.Fprintf(&stdout, "humans %d\n", b2i(w.controlledHumansHere(args.num(t))))
			default:
				t.Fatalf("unknown savage %s", action)
			}
		case "savage_population":
			for i, n := 0, args.num(t); i < n; i++ {
				w.allocBox(1100+i, tChar, subNi)
				w.pChar(1100 + i).unitItem = itemSavage
			}
		case "savage_rolls":
			n, a := args.num(t), args.num(t)
			w.streams.BeginNPC(n)
			fmt.Fprintf(&stdout, "savage spawn %d\n", w.streams.NPCSpawn(n, 0, 1, 100))
			fmt.Fprintf(&stdout, "savage qty %d\n", w.streams.NPCQty(n, a, 3, 25))
			fmt.Fprintf(&stdout, "savage wait %d\n", w.streams.NPCBehavior(a, 0, 35, 50))
		case "savage_state":
			for n := range w.loopUnits(int(w.indepPlayer)) {
				fmt.Fprintf(&stdout, "savage unit %d\n", n)
				w.saveBox(&stdout, n)
				fmt.Fprintf(&stdout, "savage output %d\n", b2i(testBit(w.pPlayer(w.player(n)).output, n)))
			}
			for i, q := range w.loadQ {
				helperList(&stdout, fmt.Sprintf("load_q %d", i), q)
			}
			helperList(&stdout, "run_q", w.runQ)
		case "hades":
			switch action := args.word(); action {
			case "auto":
				w.autoHades()
			case "unit":
				w.autoHadesSup(args.num(t))
			case "spawn":
				w.createHadesNasty(args.num(t))
			default:
				t.Fatalf("unknown hades %s", action)
			}
		case "hades_state":
			for n := range w.loopUnits(int(w.hadesPlayer)) {
				fmt.Fprintf(&stdout, "hades unit %d\n", n)
				w.saveBox(&stdout, n)
				fmt.Fprintf(&stdout, "hades output %d\n", b2i(testBit(w.pPlayer(w.player(n)).output, n)))
			}
			for i, q := range w.loadQ {
				helperList(&stdout, fmt.Sprintf("load_q %d", i), q)
			}
			helperList(&stdout, "run_q", w.runQ)
		case "nasty_rolls":
			n, a := args.num(t), args.num(t)
			fmt.Fprintf(&stdout, "nasty where %d\n", w.streams.HadsNasty(n, 0, 0, a-1))
			fmt.Fprintf(&stdout, "nasty kind %d\n", w.streams.HadsNasty(n, 1, 1, 4))
			fmt.Fprintf(&stdout, "nasty spirits %d\n", w.streams.HadsNasty(n, 2, 25, 75))
			fmt.Fprintf(&stdout, "nasty lesser %d\n", w.streams.HadsNasty(n, 2, 50, 150))
			fmt.Fprintf(&stdout, "nasty greater %d\n", w.streams.HadsNasty(n, 2, 100, 250))
		case "faery":
			switch action := args.word(); action {
			case "auto":
				w.autoFaery()
			case "spawn":
				w.createElvenHunt(args.num(t))
			case "unit":
				w.autoFaerySup(args.num(t))
			default:
				t.Fatalf("unknown faery %s", action)
			}
		case "hunt_rolls":
			n, a := args.num(t), args.num(t)
			for i := 0; i < 6; i++ {
				fmt.Fprintf(&stdout, "hunt loc %d %d\n", i, w.streams.FaerHuntLoc(n, i, 0, a-1))
			}
			fmt.Fprintf(&stdout, "hunt qty %d\n", w.streams.FaerHuntQty(n))
		case "faery_state":
			for n := range w.loopUnits(int(w.faeryPlayer)) {
				fmt.Fprintf(&stdout, "faery unit %d\n", n)
				w.saveBox(&stdout, n)
			}
			for i, q := range w.loadQ {
				helperList(&stdout, fmt.Sprintf("load_q %d", i), q)
			}
			helperList(&stdout, "run_q", w.runQ)
		case "npc_move":
			w.npcMove(args.num(t))
		case "npc_choose":
			e := w.chooseNPCDirection(args.num(t), args.num(t), args.num(t))
			if e == nil {
				fmt.Fprintln(&stdout, "chosen none")
			} else {
				fmt.Fprintf(&stdout, "chosen %d %d\n", e.direction, e.destination)
			}
		case "npc_rolls":
			n, a := args.num(t), args.num(t)
			w.streams.BeginNPC(a)
			for i := 1; i <= 4; i++ {
				fmt.Fprintf(&stdout, "npc direction %d roll %d\n", i, w.streams.NPCBehavior(n, i, 1, 10))
			}
		case "escape":
			action, n := args.word(), args.num(t)
			switch action {
			case "check":
				fmt.Fprintf(&stdout, "escaped %d\n", b2i(w.checkPrisonerEscape(n, args.num(t), args.num(t))))
			case "run":
				w.prisonerEscapes(n)
			case "move":
				w.prisonerMovementEscapeCheck(n)
			default:
				t.Fatalf("unknown escape %s", action)
			}
		case "land":
			fmt.Fprintf(&stdout, "land %d\n", w.findNearestLand(args.num(t)))
		case "routes":
			p := w.pLoc(args.num(t))
			p.provDest.Clear()
			for strings.TrimSpace(args.rest) != "" {
				p.provDest.Append(args.num(t))
			}
		case "road":
			p := w.pGate(args.num(t))
			p.toLoc, p.roadHidden = int32(args.num(t)), int8(args.num(t))
		case "delete_box":
			w.changeBoxKind(args.num(t), tDeleted)
		case "shuffle":
			who, where := args.num(t), args.num(t)
			for _, v := range w.exitsFromLocNSEWSelect(who, where, 1, 1) {
				fmt.Fprintf(&stdout, "exitdir %d %d\n", v.direction, v.destination)
			}
		case "move_state":
			helperList(&stdout, "vector", w.outVector)
			for n := range w.loopKind(tChar) {
				guard := int8(0)
				if p := w.rpChar(n); p != nil {
					guard = p.guard
				}
				fmt.Fprintf(&stdout, "move char %d where %d moving %d guard %d\n", n, w.loc(n), w.charMoving(n), guard)
				helperList(&stdout, "here", w.pLocInfo(n).hereList)
			}
			fmt.Fprintf(&stdout, "show_garrison %d\n", b2i(w.showToGarrison))
		case "entity_rolls":
			n, sub, where := args.num(t), args.num(t), args.num(t)
			fmt.Fprintf(&stdout, "prison roll %d\n", w.streams.EntPrisoner(n, sub))
			fmt.Fprintf(&stdout, "land dir99 %d\n", w.streams.EntLandDir(where, 99))
			fmt.Fprintf(&stdout, "land dir98 %d\n", w.streams.EntLandDir(where, 98))
			fmt.Fprintf(&stdout, "land pick %d\n", w.streams.EntLandPick(where, 3))
		case "escape_state":
			helperList(&stdout, "vector", w.outVector)
			for n := range w.loopKind(tChar) {
				p, m := w.rpChar(n), w.rpMagic(n)
				prisoner, moving, release := 0, 0, 0
				if p != nil {
					prisoner, moving = int(p.prisoner), int(p.moving)
				}
				if m != nil {
					release = int(m.swearOnRelease)
				}
				fmt.Fprintf(&stdout, "escape char %d where %d prisoner %d moving %d magic %d release %d\n", n, w.loc(n), prisoner, moving, b2i(m != nil), release)
				helperList(&stdout, "here", w.rpLocInfo(n).hereList)
			}
		case "arrival":
			action, n, a := args.word(), args.num(t), args.num(t)
			switch action {
			case "move":
				w.moveStack(n, a)
			case "known":
				w.markLocStackKnown(n, a)
			case "touch":
				w.touchLocAfterMove(n, a)
			case "weather":
				w.updateWeatherViewLocs(n, a)
			case "clear":
				w.clearContacts(n)
			case "match":
				w.matchTrades(n)
			case "market":
				fmt.Fprintf(&stdout, "market %d\n", w.marketHere(n))
			default:
				t.Fatalf("unknown arrival %s", action)
			}
		case "ocean_add":
			w.oceanChars.Append(args.num(t))
		case "ocean_init":
			w.initOceanChars()
		case "ocean_check":
			w.checkOceanChars()
		case "heal":
			w.healCharacters()
		case "weekly":
			w.weeklyPrisonerEscapeCheck()
		case "dogs":
			w.dogsBarkAtHiddenChars()
		case "bark":
			w.barkDogs(args.num(t))
		case "ring":
			w.incrementStoneRingAura()
		case "curse":
			w.noncreatorCurseErode()
		case "closure":
			switch action := args.word(); action {
			case "auto":
				w.autoDrop()
			case "links":
				w.linkDecay()
			case "quests":
				w.questDecay()
			case "post":
				w.postMonth()
			case "ranks":
				w.determineNobleRanks()
			case "kings":
				w.findKings()
			case "threshold":
				n := args.num(t)
				fmt.Fprintf(&stdout, "rank threshold %d %d\n", n, nprovsToRank(int32(n)))
			case "lower":
				n := args.num(t)
				fmt.Fprintf(&stdout, "rank lower %d %d\n", n, lowerNobleRank(n))
			case "sup":
				n := args.num(t)
				fmt.Fprintf(&stdout, "rank recursive %d %d\n", n, w.detNobleRankSup(n))
			case "ruler":
				n := args.num(t)
				fmt.Fprintf(&stdout, "top ruler %d %d\n", n, w.topRuler(n))
			case "describe":
				n := args.num(t)
				fmt.Fprintf(&stdout, "rank description %d [%s]\n", n, w.linerDesc(n))
				w.showCharsBelow(200, n)
			case "state":
				fmt.Fprintf(&stdout, "post flag %d month %d\n", w.postHasBeenRun, olyMonth(w.sysclock))
				for n := range w.loopKind(tChar) {
					p, rank := w.rpChar(n), int8(0)
					if p != nil {
						rank = p.rank
					}
					fmt.Fprintf(&stdout, "rank char %d allocated %d temp %d rank %d [%s]\n", n, b2i(p != nil), w.bx[n].temp, rank, w.rankS(n))
				}
			default:
				t.Fatalf("unknown closure %s", action)
			}
		case "rankland":
			n, count, strength, base, reg := args.num(t), args.num(t), args.num(t), args.num(t), args.num(t)
			w.allocBox(base, tLoc, subCastle)
			if old := w.loc(n); old != 0 {
				w.pLocInfo(old).hereList.RemValue(n)
			}
			w.pLocInfo(n).where = int32(base)
			w.pLocInfo(base).hereList.Append(n)
			for j := 1; j <= count; j++ {
				prov, garr := base+j, base+100+j
				w.allocBox(prov, tLoc, subPlain)
				w.allocBox(garr, tChar, subGarrison)
				w.pChar(garr).unitLord = 201
				w.pMisc(garr).garrCastle = int32(base)
				w.pLocInfo(prov).where = int32(reg)
				if reg != 0 {
					w.pLocInfo(reg).hereList.Append(prov)
				}
				w.pLocInfo(garr).where = int32(prov)
				w.pLocInfo(prov).hereList.Append(garr)
				w.genItem(garr, itemSoldier, strength)
			}
			w.pLocInfo(base).where = int32(base + 1)
			w.pLocInfo(base + 1).hereList.Append(base)
		case "decay":
			switch action := args.word(); action {
			case "mines":
				w.collapsedMineDecay()
			case "ghosts":
				w.ghostWarriorDecay()
			case "corpses":
				w.corpseDecay()
			case "bodies":
				w.deadBodyRot()
			case "roll":
				n := args.num(t)
				fmt.Fprintf(&stdout, "corpse roll %d %d\n", n, w.streams.UpCorpse(n, 0, 2))
			case "chains":
				fmt.Fprint(&stdout, "mines:")
				for n := range w.loopSubkind(subMineCollapsed) {
					fmt.Fprintf(&stdout, " %d", n)
				}
				fmt.Fprint(&stdout, "\nbodies:")
				for n := range w.loopSubkind(subDeadBody) {
					fmt.Fprintf(&stdout, " %d", n)
				}
				fmt.Fprintln(&stdout)
			default:
				t.Fatalf("unknown decay %s", action)
			}
		case "upkeep":
			switch action := args.word(); action {
			case "cost":
				fmt.Fprintf(&stdout, "unit maintenance %d\n", w.unitMaintCost(args.num(t)))
			case "charge":
				fmt.Fprintf(&stdout, "maintenance result %d\n", w.chargeMaintSup(args.num(t)))
			case "all":
				w.chargeMaintCosts()
			case "starve":
				w.menStarve(args.num(t), int32(args.num(t)))
			case "animals":
				w.animalDeaths()
			case "animal":
				n, item, idx := args.num(t), args.num(t), args.num(t)
				fmt.Fprintf(&stdout, "animal roll %d %d %d %d\n", n, item, idx, w.streams.UpAnimal(n, item, idx, 1, 1000))
			case "flavor":
				n, item := args.num(t), args.num(t)
				fmt.Fprintf(&stdout, "starve roll %d %d %d\n", n, item, w.streams.UpStarve(n, item, 1, 2))
			default:
				t.Fatalf("unknown upkeep %s", action)
			}
		case "income":
			switch action := args.word(); action {
			case "inn":
				w.innIncome()
			case "temple":
				w.templeIncome()
			case "roll":
				n, sub, low, high := args.num(t), args.num(t), args.num(t), args.num(t)
				fmt.Fprintf(&stdout, "income roll %d %d %d\n", n, sub, w.streams.UpIncome(n, sub, low, high))
			case "count":
				n, a, b := args.num(t), args.num(t), args.num(t)
				fmt.Fprintf(&stdout, "structures %d %d %d %d\n", n, a, b, w.countLocStructures(n, a, b))
			case "metrics":
				w.goldInn, w.goldTemple = int32(args.num(t)), int32(args.num(t))
			case "state":
				fmt.Fprintf(&stdout, "income gold inn %d temple %d\n", w.goldInn, w.goldTemple)
			default:
				t.Fatalf("unknown income %s", action)
			}
		case "monthly":
			switch action := args.word(); action {
			case "aura":
				w.incrementCurrentAura()
			case "ability":
				w.decrementAbilityShroud()
			case "shroud":
				w.decrementRegionShroud()
			case "hinder":
				w.decrementMeditationHinder()
			case "barrier":
				w.decrementLocBarrier()
			case "loyalty":
				w.loyaltyDecay()
			case "pillage":
				w.pillageDecay()
			case "relic":
				w.relicDecay()
			case "hide":
				w.hideMageDecay()
			case "roll":
				n, a := args.num(t), args.num(t)
				fmt.Fprintf(&stdout, "loyal roll %d %d %d\n", n, a, w.streams.UpLoyal(n, a, 1, 2))
			case "loyal":
				n := args.num(t)
				fmt.Fprintf(&stdout, "loyal %d %s\n", n, w.loyalS(n))
			case "state":
				for n := range w.loopKind(tChar) {
					p, m := w.rpChar(n), w.rpMagic(n)
					fresh, lord := 0, 0
					if p != nil {
						fresh, lord = int(p.freshHire), int(p.newLord)
					}
					fmt.Fprintf(&stdout, "monthly char %d char %d magic %d fresh %d new %d\n", n, b2i(p != nil), b2i(m != nil), fresh, lord)
					w.saveBox(&stdout, n)
				}
				for n := range w.loopKind(tLoc) {
					p, recent := w.rpSubloc(n), 0
					if p != nil {
						recent = int(p.recentLoot)
					}
					fmt.Fprintf(&stdout, "monthly loc %d loc %d subloc %d recent %d\n", n, b2i(w.rpLoc(n) != nil), b2i(p != nil), recent)
					w.saveBox(&stdout, n)
				}
			default:
				t.Fatalf("unknown monthly %s", action)
			}
		case "events":
			w.dailyEvents()
		case "event_days":
			fmt.Fprintf(&stdout, "curse day %d\n", w.streams.CalDay(0, 1, monthDays))
			fmt.Fprintf(&stdout, "faery day %d\n", w.streams.CalDay(1, monthDays/2, monthDays))
			fmt.Fprintf(&stdout, "dog day %d\n", w.streams.CalDay(2, 1, monthDays))
		case "bark_roll":
			fmt.Fprintf(&stdout, "bark roll %d\n", w.streams.EntBark(args.num(t), args.num(t)))
		case "prison_roll":
			fmt.Fprintf(&stdout, "weekly roll %d\n", w.streams.EntPrisoner(args.num(t), 0))
		case "heal_roll":
			fmt.Fprintf(&stdout, "heal roll %d\n", w.streams.UpHeal(args.num(t), args.num(t), args.num(t), args.num(t)))
		case "arrival_state":
			helperList(&stdout, "ocean", w.oceanChars)
			for n := range w.loopKind(tChar) {
				p := w.rpChar(n)
				time, allocated := 0, false
				if p != nil {
					time, allocated = int(p.timeFlying), p.contact != nil
				}
				fmt.Fprintf(&stdout, "arrival char %d allocated %d time %d contacts_alloc %d\n", n, b2i(p != nil), time, b2i(allocated))
				if p != nil {
					helperList(&stdout, "contacts", p.contact)
				}
			}
			for n := range w.loopKind(tPlayer) {
				p := w.rpPlayer(n)
				fmt.Fprintf(&stdout, "arrival player %d allocated %d\n", n, b2i(p != nil))
				if p != nil {
					helperList(&stdout, "known", p.known)
					helperList(&stdout, "touched", p.locs)
					helperList(&stdout, "weather", p.weatherSeen)
					helperList(&stdout, "output", p.output)
				}
			}
		case "accounts":
			switch action := args.word(); action {
			case "taxes":
				w.collectTaxes()
			case "claim":
				w.addClaimGold()
			case "np":
				w.addNoblePoints()
			case "unformed":
				w.addUnformed()
			case "sup":
				w.addUnformedSup(args.num(t))
			case "next":
				n := args.num(t)
				fmt.Fprintf(&stdout, "next np %d %d\n", n, w.nextNPTurn(n))
			case "clear":
				w.clearOrdersSent()
			case "open":
				w.specialLocsOpen()
			case "close":
				w.specialLocsClose()
			case "state":
				fmt.Fprintf(&stdout, "tax gold %d\n", w.goldTaxes)
				for n := range w.loopKind(tPlayer) {
					p := w.rpPlayer(n)
					fmt.Fprintf(&stdout, "account %d allocated %d\n", n, b2i(p != nil))
					if p != nil {
						fmt.Fprintf(&stdout, "np gained %d spent %d sent %d\n", p.npGained, p.npSpent, p.sentOrders)
						helperList(&stdout, "unformed", p.unformed)
						helperList(&stdout, "output", p.output)
					}
					w.saveBox(&stdout, n)
				}
				for n := range w.loopKind(tUnform) {
					w.saveBox(&stdout, n)
				}
			default:
				t.Fatalf("unknown accounts %s", action)
			}
		case "production":
			switch action := args.word(); action {
			case "civ":
				w.computeCivLevels()
			case "locations":
				w.locationProduction()
			case "mine":
				w.mineProduction(args.num(t))
			case "floor":
				w.replenish(args.num(t), args.num(t), args.num(t))
			case "taxes":
				w.seedTaxes()
			case "gold":
				w.moveCityGold()
			case "post":
				w.postProduction()
			case "state":
				for n := range w.loopKind(tLoc) {
					double := int8(0)
					if p := w.rpMisc(n); p != nil {
						double = p.opiumDouble
					}
					fmt.Fprintf(&stdout, "production loc %d temp %d double %d\n", n, w.bx[n].temp, double)
					w.saveBox(&stdout, n)
				}
				fmt.Fprintf(&stdout, "post %d\n", w.postHasBeenRun)
				helperList(&stdout, "pending trades", w.tradesToCheck)
			default:
				t.Fatalf("unknown production %s", action)
			}
		case "market":
			switch action := args.word(); action {
			case "add":
				w.addCityTrade(args.num(t), args.num(t), args.num(t), args.num(t), args.num(t), args.num(t))
			case "options":
				n, a, month, expire, cloak := args.num(t), args.num(t), args.num(t), args.num(t), args.num(t)
				p := w.bx[n].trades[a]
				p.monthProd, p.expire, p.cloak = int32(month), int32(expire), int32(cloak)
			case "opium":
				w.opiumMarketDelta(args.num(t))
			case "expire":
				w.expireTrades(args.num(t))
			case "restock":
				w.locTradeSup(args.num(t), args.num(t) != 0)
			case "rings":
				w.tradeSuffuseRing(args.num(t))
			case "all":
				w.locationTrades()
			case "rolls":
				n := args.num(t)
				w.streams.BeginEconomy(n)
				fmt.Fprintf(&stdout, "ring rolls %d %d %d\n", w.streams.EconStock(subSuffuseRing, n, 1, 3), w.streams.EconRing(subSuffuseRing, n, 1, 5), w.streams.EconCost(subSuffuseRing, n, 0, 12))
			case "new_ring":
				n := args.num(t)
				w.streams.BeginEconomy(n)
				created = w.newSuffuseRing(n)
				fmt.Fprintf(&stdout, "created %d\n", created)
			case "state":
				n := args.num(t)
				for i, p := range w.bx[n].trades {
					fmt.Fprintf(&stdout, "raw trade %d %d %d %d %d %d %d %d %d\n", i, p.who, p.kind, p.item, p.qty, p.cost, p.monthProd, p.expire, p.cloak)
				}
				w.saveBox(&stdout, n)
				for a := range w.loopSubkind(subSuffuseRing) {
					w.saveBox(&stdout, a)
				}
			case "buy_ring":
				n, a := args.num(t), args.num(t)
				for b := range w.loopSubkind(subSuffuseRing) {
					if w.itemUnique(b) == n {
						w.addCityTrade(a, tradeBuy, b, 1, 2000, 0)
						w.matchTrades(a)
						break
					}
				}
			default:
				t.Fatalf("unknown market %s", action)
			}
		case "trade":
			n := args.num(t)
			w.bx[n].trades.Append(&trade{kind: int32(args.num(t)), item: int32(args.num(t)), qty: int32(args.num(t)), cost: int32(args.num(t)), who: int32(args.num(t)), sort: int32(args.num(t))})
		case "economy":
			switch action := args.word(); action {
			case "all":
				w.matchAllTrades()
			case "check":
				w.checkValidatedTrades()
			case "queue":
				w.tradesToCheck.Append(args.num(t))
			case "options":
				n, a, b, d := args.num(t), args.num(t), args.num(t), args.num(t)
				w.bx[n].trades[a].haveLeft = int32(b)
				w.bx[n].trades[a].cloak = int32(d)
			case "reduce":
				n, a, b := args.num(t), args.num(t), args.num(t)
				fmt.Fprintf(&stdout, "reduced %d\n", w.reduceQty(w.bx[n].trades[a], int32(b)))
			case "attempt":
				n, a, b, d := args.num(t), args.num(t), args.num(t), args.num(t)
				w.attemptTrade(w.bx[n].trades[a], w.bx[b].trades[d])
			case "raw":
				n := args.num(t)
				for i, p := range w.bx[n].trades {
					fmt.Fprintf(&stdout, "trade raw %d %d %d %d %d %d %d %d %d %d %d\n", i, p.who, p.kind, p.item, p.qty, p.cost, p.cloak, p.haveLeft, p.monthProd, p.expire, p.sort)
				}
				w.saveBox(&stdout, n)
			case "state":
				fmt.Fprintf(&stdout, "trade gold %d pot basket %d opium %d\n", w.goldTrade, w.goldPotBasket, w.goldOpium)
				helperList(&stdout, "pending trades", w.tradesToCheck)
			default:
				t.Fatalf("unknown economy %s", action)
			}
		case "offers":
			seller, where, except := args.num(t), args.num(t), args.num(t)
			var l tradesList
			if seller != 0 {
				l = w.sellerList(where, except)
			} else {
				l = w.buyerList(where, except)
			}
			for _, p := range l {
				fmt.Fprintf(&stdout, "offer %d %d %d %d %d\n", p.who, p.item, p.qty, p.cost, p.sort)
			}
		case "master":
			n := args.num(t)
			seed := make([]byte, 16)
			for i := range seed {
				seed[i] = byte(n + i)
			}
			path := filepath.Join(w.libdir, "randseed")
			if err := os.WriteFile(path, seed, 0o600); err != nil {
				t.Fatal(err)
			}
			if err := w.rnd.Load(path); err != nil {
				t.Fatal(err)
			}
		case "ambush", "bandit":
			faery, who, where := args.num(t) != 0, args.num(t), args.num(t)
			if op == "bandit" {
				if faery {
					created = w.createFaeryBandit(who, where)
				} else {
					created = w.createHadesBandit(who, where)
				}
				fmt.Fprintf(&stdout, "created %d\n", created)
			} else if faery {
				w.faeryAttackCheck(who, where)
			} else {
				w.hadesAttackCheck(who, where)
			}
		case "region_rolls":
			n, a := args.num(t), args.num(t)
			fmt.Fprintf(&stdout, "hads ambush %d\n", w.streams.HadsAmbush(n, a))
			fmt.Fprintf(&stdout, "hads kind %d\n", w.streams.HadsBanditKind(n, a))
			fmt.Fprintf(&stdout, "hads retal %d\n", w.streams.HadsRetal(n, a))
			fmt.Fprintf(&stdout, "faer ambush %d\n", w.streams.FaerAmbush(n, a))
			fmt.Fprintf(&stdout, "faer kind %d\n", w.streams.FaerBanditKind(n, a))
			fmt.Fprintf(&stdout, "faer retal %d\n", w.streams.FaerRetal(n, a))
		case "bandits":
			for n := range w.loopKind(tChar) {
				if w.subkind(n) != subNi {
					continue
				}
				fmt.Fprintf(&stdout, "bandit %d\n", n)
				w.saveBox(&stdout, n)
				fmt.Fprintf(&stdout, "bandit output %d\n", b2i(testBit(w.pPlayer(w.player(n)).output, n)))
			}
			for i, q := range w.loadQ {
				helperList(&stdout, fmt.Sprintf("load_q %d", i), q)
			}
			helperList(&stdout, "run_q", w.runQ)
		case "select_bandit":
			for n := range w.loopKind(tChar) {
				if w.subkind(n) == subNi {
					c = w.rpCommand(n)
					break
				}
			}
		case "learn":
			n, a, b, d := args.num(t), args.num(t), args.num(t), args.num(t)
			w.pChar(n).skills.Append(&skillEnt{skill: int32(a), know: int8(b), experience: int16(d)})
		case "payment":
			kind, n, a := args.word(), args.num(t), args.num(t)
			var result int
			switch kind {
			case "can":
				result = b2i(w.canPay(n, a))
			case "charge":
				result = b2i(w.charge(n, a))
			case "auto":
				result = b2i(w.autocharge(n, a))
			case "stack":
				result = int(w.stackHasItem(n, a))
			case "sub":
				result = b2i(w.stackSubItem(n, a, args.num(t)))
			case "key":
				result = w.hasUseKey(n, a)
			case "stackkey":
				result = w.stackHasUseKey(n, a)
			default:
				t.Fatalf("unknown payment %q", kind)
			}
			fmt.Fprintf(&stdout, "payment %d\n", result)
		case "genitem":
			w.genItem(args.num(t), args.num(t), args.num(t))
		case "output":
			n, a := args.num(t), args.num(t)
			fmt.Fprintf(&stdout, "output %d\n", b2i(testBit(w.pPlayer(n).output, a)))
		case "orders_file":
			n := args.num(t)
			w.savePlayerOrders(n)
			fmt.Fprintf(&stdout, "orders file %d:\n", n)
			data, err := os.ReadFile(filepath.Join(w.libdir, "orders", fmt.Sprint(n)))
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			stdout.Write(data)
		case "emit":
			w.commandOut(args.num(t), "routing %d  ", args.num(t))
		case "sout":
			switch action := args.word(); action {
			case "route":
				w.outPath, w.outAltWho = args.num(t), args.num(t)
			case "frame":
				w.indent, w.secondIndent = int32(args.num(t)), int32(args.num(t))
				w.showDay, w.sysclock.day = args.num(t) != 0, int16(args.num(t))
			case "style":
				w.style(int32(args.num(t)))
			case "clear":
				w.vectorClear()
			case "add":
				w.vectorAdd(args.num(t))
			case "players":
				w.vectorPlayers()
			case "save":
				savedVector = w.saveOutputVector()
			case "restore":
				w.restoreOutputVector(savedVector)
				savedVector = nil
			case "state":
				fmt.Fprintf(&stdout, "route %d %d frame %d %d %d %d\n", w.outPath, w.outAltWho, w.indent, w.secondIndent, b2i(w.showDay), w.sysclock.day)
				helperList(&stdout, "vector", w.outVector)
				helperList(&stdout, "saved vector", savedVector)
			case "text":
				n, kind, ind, count := args.num(t), args.word(), args.num(t), args.num(t)
				s := args.ctext()
				if count < 0 || count*len(s) > 16384 {
					t.Fatal("oversized sout text")
				}
				s = strings.ReplaceAll(strings.Repeat(s, count), "|", "\n")
				switch kind {
				case "out":
					w.commandOut(n, "%s", s)
				case "wiout":
					w.wiout(n, int32(ind), "%s", s)
				case "html":
					w.html(n, "%s", s)
				case "lines":
					w.lines(n, s)
				case "match":
					w.matchLines(n, s)
				case "log":
					w.logWrite(n, "%s", s)
				default:
					t.Fatalf("unknown sout text %s", kind)
				}
			default:
				t.Fatalf("unknown sout %s", action)
			}
		case "garrison_view":
			w.showToGarrison = args.num(t) != 0
		case "lore_state":
			p := w.pPlayer(args.num(t))
			helperList(&stdout, "known", p.known)
			helperList(&stdout, "lore", p.deliverLore)
			helperList(&stdout, "touched", p.locs)
			helperList(&stdout, "weather", p.weatherSeen)
		case "queue_lore":
			w.queueLore(args.num(t), args.num(t), args.num(t) != 0)
		case "skill_lore":
			w.scanCharSkillLore()
		case "show_lore":
			w.showLoreSheetsState()
		case "gm_lore":
			w.gmShowAllSkillsState(args.num(t))
		case "lore_root":
			t.Chdir(w.libdir)
			w.libdir = "."
		case "lore_dir":
			if err := os.Mkdir(args.word(), 0o755); err != nil {
				t.Fatal(err)
			}
		case "lore_file":
			path := args.word()
			if err := os.WriteFile(path, []byte(strings.ReplaceAll(args.ctext(), "|", "\n")), 0o644); err != nil {
				t.Fatal(err)
			}
		case "life":
			switch action := args.word(); action {
			case "open":
				w.winFlag = args.num(t) != 0
				w.openLogfile(w.winFlag)
			case "nondestruct":
				w.openLogfileNondestruct()
			case "close":
				w.closeLogfile()
			case "flush":
				clib.FlushAll()
			case "save":
				w.winFlag = args.num(t) != 0
				w.saveLogdir(w.winFlag)
			case "times_open":
				w.openTimes()
			case "times_close":
				w.closeTimes()
			case "fill":
				for i := 201; i < 200+maxPlayerOutput(); i++ {
					if !w.validBox(i) {
						w.allocBox(i, tPlayer, subPlRegular)
					}
					w.commandOut(i, "cache fill")
				}
			case "extra":
				w.allocBox(220, tPlayer, subPlRegular)
				w.commandOut(220, "cache extra")
			case "path":
				fileState(args.word())
			case "watch":
				watched = append(watched, args.word())
			case "rename":
				if err := os.Rename(args.word(), args.word()); err != nil {
					t.Fatal(err)
				}
			case "unlink":
				if err := os.Remove(args.word()); err != nil {
					t.Fatal(err)
				}
			case "remove":
				removeDir(args.word())
			default:
				t.Fatalf("unknown lifecycle %s", action)
			}
		case "post":
			switch action := args.word(); action {
			case "players":
				w.writePlayerListState()
			case "email":
				w.writeEmailState()
			case "totimes":
				w.writeTotimesState()
			case "forward":
				w.writeForwardsState()
			case "factions":
				w.writeFactionsState()
			case "all":
				w.writePlayerListState()
				w.writeEmailState()
				w.writeTotimesState()
				w.writeForwardsState()
				w.writeFactionsState()
			case "file":
				path := args.word()
				st, err := os.Stat(path)
				if err != nil {
					fmt.Fprintf(&stdout, "file %s missing\n", path)
				} else if st.Mode().IsRegular() {
					fmt.Fprintf(&stdout, "file %s regular\n", path)
				} else {
					fmt.Fprintf(&stdout, "file %s directory\n", path)
				}
			case "empty_player":
				w.bx[args.num(t)].xPlayer = nil
			case "new":
				w.newPlayers.Append(args.num(t))
			case "mail":
				p := w.pPlayer(args.num(t))
				switch kind := args.word(); kind {
				case "null":
					p.email = nil
				case "empty":
					s := ""
					p.email = &s
				case "visible":
					s := "public@example.test"
					p.visEmail = &s
				case "cis":
					p.compuserve = 1
				default:
					t.Fatalf("unknown mail %s", kind)
				}
			default:
				t.Fatalf("unknown post %s", action)
			}
		case "admission":
			switch action := args.word(); action {
			case "run":
				w.addNewPlayers()
			case "probe":
				w.makeNewPlayers()
			case "parse":
				acct, path := args.word(), args.word()
				fp, err := os.Open(path)
				if err != nil {
					t.Fatal(err)
				}
				defer fp.Close()
				fmt.Fprintf(&stdout, "admission result %d\n", b2i(w.makeNewPlayersSup(acct, bufio.NewReader(fp))))
			case "fetch":
				path, count := args.word(), args.num(t)
				fp, err := os.Open(path)
				if err != nil {
					t.Fatal(err)
				}
				defer fp.Close()
				reader := bufio.NewReader(fp)
				for i := 0; i < count; i++ {
					if s := fetchInp(reader); s != nil {
						fmt.Fprintf(&stdout, "field %d [%s]\n", i, *s)
					} else {
						fmt.Fprintf(&stdout, "field %d NULL\n", i)
					}
				}
			case "hex":
				path := args.word()
				data, err := hex.DecodeString(args.word())
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, data, 0o644); err != nil {
					t.Fatal(err)
				}
			case "lib":
				w.libdir = args.word()
			case "state":
				helperList(&stdout, "new players", w.newPlayers)
				helperList(&stdout, "vector", w.outVector)
				for n := range w.loopKind(tPlayer) {
					fmt.Fprintf(&stdout, "admission player %d record %d\n", n, b2i(w.rpPlayer(n) != nil))
					w.saveBox(&stdout, n)
				}
			default:
				t.Fatalf("unknown admission %s", action)
			}
		case "set_known":
			w.setKnown(args.num(t), args.num(t))
		case "additem":
			w.addItem(args.num(t), args.num(t), args.num(t))
		case "damage":
			switch action := args.word(); action {
			case "coast":
				fmt.Fprintf(&stdout, "coast %d\n", w.nearRockyCoast(args.num(t)))
			case "ships":
				w.shipCoastalDamage()
			case "locations":
				w.randomLocDamage()
			case "mine":
				w.mineCalamity(args.num(t))
			case "inn":
				w.innCalamity(args.num(t))
			case "count":
				n := args.num(t)
				fmt.Fprintf(&stdout, "workers %d depth %d\n", w.countLocCharItem(n, itemWorker), w.mineDepth(n))
			case "roll":
				fmt.Fprintf(&stdout, "wreck %d\n", w.streams.WthrWreck(args.num(t), args.num(t), args.num(t), args.num(t)))
			default:
				t.Fatalf("unknown damage %s", action)
			}
		case "structure":
			switch action := args.word(); action {
			case "damage":
				fmt.Fprintf(&stdout, "destroyed %d\n", b2i(w.addStructureDamage(args.num(t), args.num(t), args.num(t) != 0)))
			case "sink":
				w.sinkShip(args.num(t))
			case "collapse":
				w.buildingCollapses(args.num(t))
			case "remove":
				w.getRidOfCollapsedMine(args.num(t))
			case "ocean":
				w.killStackOcean(args.num(t))
			case "state":
				n := args.num(t)
				damage, delay, castle, bind := -1, 0, 0, 0
				if p := w.rpSubloc(n); p != nil {
					damage = int(p.damage)
				}
				if m := w.rpMisc(n); m != nil {
					delay, castle, bind = int(m.mineDelay), int(m.garrCastle), int(m.bindStorm)
				}
				fmt.Fprintf(&stdout, "structure %d kind %d sub %d damage %d delay %d castle %d bind %d\n", n, w.kind(n), w.subkind(n), damage, delay, castle, bind)
				if p := w.rpSubloc(n); p != nil {
					helperList(&stdout, "bound", p.boundStorms)
				} else {
					helperList(&stdout, "bound", nil)
				}
				helperList(&stdout, "here", w.rpLocInfo(n).hereList)
			default:
				t.Fatalf("unknown structure %s", action)
			}
		case "death":
			switch action := args.word(); action {
			case "kill":
				w.killChar(args.num(t), args.num(t))
			case "reclaim":
				w.charReclaim(args.num(t))
			case "subkind":
				w.changeBoxSubkind(args.num(t), args.num(t))
			case "survive":
				fmt.Fprintf(&stdout, "survive %d\n", b2i(w.surviveFatal(args.num(t))))
			case "forget":
				fmt.Fprintf(&stdout, "forget %d\n", b2i(w.forgetSkill(args.num(t), args.num(t))))
			case "location":
				fmt.Fprintf(&stdout, "location [%s]\n", w.charRepLocation(args.num(t)))
			case "hades":
				fmt.Fprintf(&stdout, "hades location %d\n", w.randomHadesLoc(args.num(t)))
			case "state":
				for _, kind := range []int{tChar, tItem, tDeadchar} {
					fmt.Fprintf(&stdout, "kind chain %d:", kind)
					for n := range w.loopKind(kind) {
						fmt.Fprintf(&stdout, " %d", n)
					}
					fmt.Fprintln(&stdout)
				}
			default:
				t.Fatalf("unknown death %s", action)
			}
		case "inherit":
			switch action := args.word(); action {
			case "take":
				w.takeUnitItems(args.num(t), args.num(t), args.num(t))
			case "recipient":
				fmt.Fprintf(&stdout, "recipient %d\n", w.stackmateInheritor(args.num(t)))
			case "drop":
				w.dropStack(args.num(t), args.num(t))
			case "extract":
				w.extractStackedUnit(args.num(t))
			case "prisoner":
				fmt.Fprintf(&stdout, "move prisoner %d\n", w.movePrisoner(args.num(t), args.num(t), args.num(t)))
			case "desert":
				w.unitDeserts(args.num(t), args.num(t), args.num(t) != 0, args.num(t), args.num(t))
			case "cookie":
				w.putBackCookie(args.num(t))
			case "rolls":
				n, a, b := args.num(t), args.num(t), args.num(t)
				fmt.Fprintf(&stdout, "take gate %d\n", w.streams.EntTake(n, a, 0, 1, 2))
				fmt.Fprintf(&stdout, "take amount %d\n", w.streams.EntTake(n, a, 1, 0, b))
				fmt.Fprintf(&stdout, "drop gate %d\n", w.streams.EntDrop(n, a, 0, 1, 5))
				fmt.Fprintf(&stdout, "drop message %d\n", w.streams.EntDrop(n, a, 1, 1, 3))
			case "name":
				n, a := args.num(t), args.num(t)
				fmt.Fprintf(&stdout, "plural box [%s]\n", w.pluralItemBox(n, a))
				fmt.Fprintf(&stdout, "box qty [%s]\n", w.boxNameQty(n, a))
			case "state":
				fmt.Fprintf(&stdout, "combat gold %d independent %d indent %d show_day %d\n", w.goldCombat, w.goldCombatIndep, w.indent, b2i(w.showDay))
				for n := range w.loopKind(tPlayer) {
					fmt.Fprintf(&stdout, "inherit player %d\n", n)
					helperList(&stdout, "bits", w.pPlayer(n).output)
				}
				for i, q := range w.loadQ {
					helperList(&stdout, fmt.Sprintf("load_q %d", i), q)
				}
				helperList(&stdout, "run_q", w.runQ)
			default:
				t.Fatalf("unknown inherit %s", action)
			}
		case "moveitem":
			fmt.Fprintf(&stdout, "moveitem %d\n", b2i(w.moveItem(args.num(t), args.num(t), args.num(t), args.num(t))))
		case "dropitem", "consume":
			n, a, b := args.num(t), args.num(t), args.num(t)
			var result bool
			if op == "consume" {
				result = w.consumeItem(n, a, b)
			} else {
				result = w.dropItem(n, a, b)
			}
			fmt.Fprintf(&stdout, "inventory result %d\n", b2i(result))
		case "kindname":
			n := args.num(t)
			fmt.Fprintf(&stdout, "kind [%s] name [%s]\n", w.displayKind(n), w.boxNameKind(n))
		case "names":
			action, n, qty := args.word(), args.num(t), args.num(t)
			switch action {
			case "repeat":
				if qty < 0 || qty > 4096 {
					t.Fatal("oversized repeated name")
				}
				w.setName(n, strings.Repeat("a", qty))
			case "show":
				fmt.Fprintf(&stdout, "display [%s]\n", w.displayName(n))
				fmt.Fprintf(&stdout, "code [%s]\n", w.boxCode(n))
				fmt.Fprintf(&stdout, "code less [%s]\n", boxCodeLess(n))
				fmt.Fprintf(&stdout, "box [%s]\n", w.boxName(n))
				fmt.Fprintf(&stdout, "just [%s]\n", w.justName(n))
				fmt.Fprintf(&stdout, "plural [%s]\n", w.pluralItemName(n, qty))
				fmt.Fprintf(&stdout, "plural box [%s]\n", w.pluralItemBox(n, qty))
				fmt.Fprintf(&stdout, "just qty [%s]\n", w.justNameQty(n, qty))
				fmt.Fprintf(&stdout, "box qty [%s]\n", w.boxNameQty(n, qty))
			default:
				t.Fatalf("unknown names %s", action)
			}
		case "reconcile":
			w.checkTokenUnits()
		case "process":
			w.processOrders()
		case "token_owner":
			fmt.Fprintf(&stdout, "token player %d\n", w.tokenPlayer(args.num(t)))
		case "token_state":
			n := args.num(t)
			w.saveBox(&stdout, n)
			for a := range w.loopUnits(n) {
				fmt.Fprintf(&stdout, "token row %d\n", a)
				w.saveBox(&stdout, a)
			}
			helperList(&stdout, "pending trades", w.tradesToCheck)
			for i, q := range w.loadQ {
				helperList(&stdout, fmt.Sprintf("load_q %d", i), q)
			}
			helperList(&stdout, "run_q", w.runQ)
			for a := range w.loopKind(tPlayer) {
				fmt.Fprintf(&stdout, "token output %d\n", a)
				helperList(&stdout, "bits", w.pPlayer(a).output)
			}
		case "scan_lore":
			w.scanCharItemLore()
		case "init_touched":
			w.initLocsTouched()
		case "init_weather":
			w.initWeatherViews()
		case "newchar":
			sk, ni, where, health := args.num(t), args.num(t), args.num(t), args.num(t)
			pl, lk, ll := args.num(t), args.num(t), args.num(t)
			created = w.newChar(sk, ni, where, health, pl, lk, ll, strings.TrimLeft(args.rest, " \t"))
			fmt.Fprintf(&stdout, "created %d\n", created)
		case "newitem":
			created = w.createUniqueItem(args.num(t), args.num(t))
			fmt.Fprintf(&stdout, "created %d\n", created)
		case "save_new":
			w.saveBox(&stdout, created)
		case "delete_new":
			w.destroyUniqueItem(args.num(t), created)
		case "drop_new":
			fmt.Fprintf(&stdout, "drop new %d\n", b2i(w.dropItem(args.num(t), created, args.num(t))))
		case "created_state":
			fmt.Fprintf(&stdout, "created kind %d sub %d unique %d\n", w.kind(created), w.subkind(created), w.itemUnique(created))
		case "body_new":
			p := w.pMisc(created)
			p.oldLord = int32(args.num(t))
			q := w.pChar(created)
			q.prevLord, q.loyKind, q.loyRate = int32(args.num(t)), int8(args.num(t)), int32(args.num(t))
			name := strings.TrimLeft(args.rest, " \t")
			p.saveName = nil
			if name != "" {
				p.saveName = strSave(name)
			}
		case "destroy":
			w.destroyUniqueItem(args.num(t), args.num(t))
		case "subitem":
			fmt.Fprintf(&stdout, "subitem %d\n", b2i(w.subItem(args.num(t), args.num(t), args.num(t))))
		case "join":
			w.joinStack(args.num(t), args.num(t))
		case "leave":
			w.leaveStack(args.num(t))
		case "queue":
			w.pOrderHead(w.player(int(c.who)), int(c.who)).l.Append(strings.TrimPrefix(line, "queue "))
		case "orders":
			n, a := args.num(t), args.num(t)
			w.pOrderHead(n, a).l.Append(strings.TrimLeft(args.rest, " \t"))
		case "transfer":
			n, a, b, d := args.num(t), args.num(t), args.num(t), args.num(t)
			w.setLord(n, a, b, d)
		case "generic_interrupt":
			w.interruptOrder(args.num(t))
		case "stops":
			w.processInterruptedUnits()
		case "flush":
			n, a := args.num(t), args.num(t)
			w.flushUnitOrders(n, a)
		case "stop_check":
			n, a := args.num(t), args.num(t)
			fmt.Fprintf(&stdout, "stop %d\n", b2i(w.stopOrder(n, a)))
		case "np_total":
			fmt.Fprintf(&stdout, "np %d\n", w.charNPTotal(args.num(t)))
		case "save":
			w.saveBox(&stdout, args.num(t))
		case "report":
			switch action := args.word(); action {
			case "order":
				w.determineOutputOrder()
			case "inventory":
				n, owner := args.num(t), args.num(t)
				w.showCharInventory(n, owner)
			case "unclaimed":
				n, owner := args.num(t), args.num(t)
				w.showUnclaimed(n, owner)
			case "capacity":
				n, owner := args.num(t), args.num(t)
				w.showCarryCapacity(n, owner)
			case "item_skills":
				n, owner := args.num(t), args.num(t)
				w.showItemSkills(n, owner)
			case "skills":
				n, owner := args.num(t), args.num(t)
				w.listSkills(n, owner)
			case "partial_skills":
				n, owner := args.num(t), args.num(t)
				w.listPartialSkills(n, owner)
			case "stack_info", "pledged", "prisoners", "health", "combat", "detail_misc", "magic":
				n, owner := args.num(t), args.num(t)
				switch action {
				case "stack_info":
					w.charRepStackInfo(n, owner)
				case "pledged":
					w.showPledged(n, owner)
				case "prisoners":
					w.charRepPrisoners(n, owner)
				case "health":
					w.charRepHealth(n, owner)
				case "combat":
					w.charRepCombat(n, owner)
				case "detail_misc":
					w.charRepMisc(n, owner)
				case "magic":
					w.charRepMagic(n, owner)
				}
			case "prisoner_health":
				fmt.Fprintf(&stdout, "prisoner health [%s]\n", w.prisonerHealth(args.num(t)))
			case "experience":
				fmt.Fprintf(&stdout, "experience [%s]\n", expS(args.num(t)))
			case "qualifier":
				n, know, days := args.num(t), args.num(t), args.num(t)
				p := &skillEnt{skill: int32(n), know: int8(know), daysStudied: int32(days)}
				fmt.Fprintf(&stdout, "qualifier [%s]\n", w.fractionalSkillQualifier(p))
			case "skill_rows":
				n := args.num(t)
				p := w.rpChar(n)
				fmt.Fprintf(&stdout, "skill record %d %d\n", n, b2i(p != nil))
				if p != nil {
					for _, s := range p.skills {
						fmt.Fprintf(&stdout, "skill row %d %d %d %d\n", s.skill, s.know, s.daysStudied, s.experience)
					}
				}
				w.saveBox(&stdout, n)
			case "root":
				n, where := args.num(t), args.num(t)
				w.showLoc(n, where)
			case "no_header":
				w.showLocNoHeader = args.num(t) != 0
			case "root_state":
				p := w.rpPlayer(args.num(t))
				fmt.Fprintf(&stdout, "location globals %d %d %d %d player %d\n", w.outPath, w.outAltWho, b2i(w.showLocNoHeader), b2i(w.showDisplayString), b2i(p != nil))
				var touched, output sparse
				if p != nil {
					touched, output = p.locs, p.output
				}
				helperList(&stdout, "touched", touched)
				helperList(&stdout, "output", output)
			case "ancillary":
				mode, n, where := args.num(t), args.num(t), args.num(t)
				switch mode {
				case 0:
					w.showNearbyCities(n, where)
				case 1:
					w.showLocSkills(n, where)
				case 2:
					w.showWeather(n, where)
				case 3:
					w.showLocRuler(n, where)
				case 4:
					w.showLocStats(n, where)
				}
			case "posts":
				n, where, full := args.num(t), args.num(t), args.num(t)
				w.showLocPosts(n, where, full != 0)
			case "post_line":
				n := args.num(t)
				_, text, ok := strings.Cut(args.rest, "@")
				if !ok {
					t.Fatal("missing post text delimiter")
				}
				w.pMisc(n).postTxt = append(w.pMisc(n).postTxt, text)
			case "post_empty":
				w.pMisc(args.num(t))
			case "ancillary_state":
				n := args.num(t)
				s, p := w.rpSubloc(n), w.rpPlayer(n)
				fmt.Fprintf(&stdout, "ancillary records %d SL %d PL %d MI %d IM %d CH %d\n", n, b2i(s != nil), b2i(p != nil), b2i(w.rpMisc(n) != nil), b2i(w.rpItemMagic(n) != nil), b2i(w.rpChar(n) != nil))
				if s != nil {
					helperList(&stdout, "near", s.nearCities)
					helperList(&stdout, "teaches", s.teaches)
				} else {
					helperList(&stdout, "near", nil)
					helperList(&stdout, "teaches", nil)
				}
				if p != nil {
					helperList(&stdout, "weather", p.weatherSeen)
					helperList(&stdout, "output", p.output)
				} else {
					helperList(&stdout, "weather", nil)
					helperList(&stdout, "output", nil)
				}
			case "inner_locs", "sublocs_here", "ships_here":
				n, where := args.num(t), args.num(t)
				switch action {
				case "inner_locs":
					w.showInnerLocs(n, where)
				case "sublocs_here":
					w.showSublocsHere(n, where)
				case "ships_here":
					w.showShipsHere(n, where)
				}
			case "chars_here", "owner_stack", "chars_below":
				n, where := args.num(t), args.num(t)
				switch action {
				case "chars_here":
					w.showCharsHere(n, where)
				case "owner_stack":
					w.showOwnerStack(n, where)
				case "chars_below":
					w.showCharsBelow(n, where)
				}
			case "highlight", "below_highlight":
				n, who, depth := args.num(t), args.num(t), int32(args.num(t))
				if action == "highlight" {
					fmt.Fprintf(&stdout, "highlight [%s]\n", w.highlightUnits(n, who, depth))
				} else {
					w.showCharsBelowHighlight(n, who, depth)
				}
			case "market":
				w.marketReport(args.num(t), args.num(t))
			case "exits", "sailable":
				n, where := args.num(t), args.num(t)
				if action == "exits" {
					w.listExits(n, where)
				} else {
					w.listSailableRoutes(n, where)
				}
			case "route_view":
				n := args.num(t)
				v := &exitView{}
				for _, p := range []*int32{&v.road, &v.direction, &v.destination, &v.distance, &v.hidden, &v.destHidden, &v.inside, &v.impassable, &v.inTransit, &v.magicBarrier, &v.hadesCost} {
					*p = int32(args.num(t))
				}
				mode, oldIndent := args.num(t), w.indent
				first := "Direct routes:"
				if mode == 1 {
					first = ""
				}
				p := &first
				if mode == 0 {
					p = nil
				}
				if v.road != 0 {
					w.listRoadSup(n, 0, v, p)
				} else {
					w.listExitsSup(n, 0, v, p)
				}
				fmt.Fprintf(&stdout, "route first %d indent %d second %d\n", b2i(first != ""), w.indent, w.secondIndent)
				w.indent = oldIndent
			case "describe":
				n := args.num(t)
				fmt.Fprintf(&stdout, "description [%s]\n", w.linerDesc(n))
			case "describe_out":
				n, where := args.num(t), args.num(t)
				w.commandOut(n, "%s", w.linerDesc(where))
			case "owner":
				fmt.Fprintf(&stdout, "owner [%s]\n", w.displayOwner(args.num(t)))
			case "display":
				w.showDisplayString = args.num(t) != 0
			case "set_banner":
				n := args.num(t)
				w.setBanner(n, strings.TrimLeft(args.rest, " \t"))
			case "header":
				n, where := args.num(t), args.num(t)
				w.commandOut(n, "%s", w.showLocHeader(where))
			case "barrier":
				n, where := args.num(t), args.num(t)
				fmt.Fprintf(&stdout, "barrier %d\n", b2i(w.showLocBarrier(n, where)))
			case "header_parts":
				n := args.num(t)
				fmt.Fprintf(&stdout, "header records %d %d\n", b2i(w.rpLoc(n) != nil), b2i(w.rpSubloc(n) != nil))
				fmt.Fprintf(&stdout, "inside [%s]\n", w.locInsideString(n))
				fmt.Fprintf(&stdout, "civ [%s]\n", w.locCivS(n))
				fmt.Fprintf(&stdout, "header [%s]\n", w.showLocHeader(n))
				fmt.Fprintf(&stdout, "header records %d %d\n", b2i(w.rpLoc(n) != nil), b2i(w.rpSubloc(n) != nil))
			case "players":
				w.playerEntInfo()
			case "permissions":
				w.printAdmit(args.num(t))
			case "attitudes":
				n, target := args.num(t), args.num(t)
				w.printAtt(n, target)
			case "characters":
				w.characterReport()
			case "character":
				w.charRepSup(args.num(t), args.num(t))
			case "pending_trades":
				w.listPendingTrades(args.num(t), args.num(t))
			case "locations":
				w.turnEndLocReports()
			case "location":
				n := args.num(t)
				w.commandOut(n, "%s", w.charRepLocation(n))
			case "banner":
				w.playerBanner()
			case "banner_win":
				w.winFlag = args.num(t) != 0
			case "banner_title":
				w.gameTitle = args.ctext()
			case "rules_url":
				w.rulesURL = args.ctext()
			case "times_url":
				w.timesURL = args.ctext()
			case "provinces":
				w.summaryProvincesState()
			case "templates":
				w.listOrderTemplates()
			case "template":
				who, pl := args.num(t), args.num(t)
				w.ordersTemplate(who, pl)
			case "saved":
				w.readBoxes(args.word())
			case "account":
				w.reportAccount()
			case "account_sup":
				w.reportAccountSup(args.num(t))
			case "summary":
				w.summaryReportState()
			case "factions":
				w.playerReportState()
			case "gm":
				w.gmReportState(args.num(t))
			case "include":
				w.genIncludeSectionState()
			case "include_sup":
				w.genIncludeSupState(args.num(t))
			case "bit":
				n, bit := args.num(t), args.num(t)
				setBit(&w.pPlayer(n).output, bit)
			case "block":
				if err := os.Mkdir(filepath.Join(w.libdir, "log", fmt.Sprint(args.num(t))), 0o755); err != nil {
					t.Fatal(err)
				}
			case "files":
				for n := range w.loopKind(tPlayer) {
					_, err := os.Stat(filepath.Join(w.libdir, "log", fmt.Sprint(n)))
					fmt.Fprintf(&stdout, "report file %d %d\n", n, b2i(err == nil))
				}
			case "stored_orders":
				if p := w.rpPlayer(args.num(t)); p != nil {
					for _, q := range p.orders {
						fmt.Fprintf(&stdout, "stored orders %d:", q.unit)
						for _, s := range q.l {
							fmt.Fprintf(&stdout, " [%s]", s)
						}
						fmt.Fprintln(&stdout)
					}
				}
			case "known":
				if p := w.rpPlayer(args.num(t)); p != nil {
					helperList(&stdout, "raw known", p.known)
				}
			case "misc":
				n := args.num(t)
				fmt.Fprintf(&stdout, "report misc %d %d\n", n, b2i(w.rpMisc(n) != nil))
				w.saveBox(&stdout, n)
			case "style":
				w.style(int32(args.num(t)))
			case "style_probe":
				w.bottomOut(args.num(t), 0, "style probe")
			case "skill":
				n, a, b, d := args.num(t), args.num(t), args.num(t), args.num(t)
				w.pChar(n).skills.Append(&skillEnt{skill: int32(a), know: int8(b), daysStudied: int32(d)})
			case "backlinks":
				n := args.num(t)
				p := w.rpMagic(n)
				fmt.Fprintf(&stdout, "backlinks ready %d allocated %d\n", b2i(w.pledgeBacklinks), b2i(p != nil && p.pledgedToUs != nil))
				if p != nil {
					helperList(&stdout, "pledged", p.pledgedToUs)
				}
				w.saveBox(&stdout, n)
			case "backlink":
				n, a := args.num(t), args.num(t)
				w.pMagic(n).pledgedToUs.Append(a)
			case "admit":
				n, a, b, d := args.num(t), args.num(t), args.num(t), args.num(t)
				p := &admit{targ: int32(a), sense: int32(b)}
				if d != 0 {
					p.l.Append(d)
				}
				w.pPlayer(n).admits.Append(p)
			case "member":
				n, a := args.num(t), args.num(t)
				l := w.pPlayer(n).admits
				l[len(l)-1].l.Append(a)
			case "state":
				n := args.num(t)
				p, att := w.rpPlayer(n), w.rpDisp(n)
				fmt.Fprintf(&stdout, "report entity %d order %d player %d attitudes %d char %d magic %d\n",
					n, w.bx[n].outputOrder, b2i(p != nil), b2i(att != nil), b2i(w.rpChar(n) != nil), b2i(w.rpMagic(n) != nil))
				for _, e := range w.bx[n].items {
					fmt.Fprintf(&stdout, "raw item %d %d\n", e.item, e.qty)
				}
				if p != nil {
					helperList(&stdout, "raw units", p.units)
					for _, a := range p.admits {
						fmt.Fprintf(&stdout, "raw admit %d %d\n", a.targ, a.sense)
						helperList(&stdout, "members", a.l)
					}
				}
				if att != nil {
					helperList(&stdout, "raw hostile", att.hostile)
					helperList(&stdout, "raw neutral", att.neutral)
					helperList(&stdout, "raw defend", att.defend)
				}
				w.saveBox(&stdout, n)
			default:
				t.Fatalf("unknown report action %s", action)
			}
		case "pstate":
			n := args.num(t)
			p := w.pPlayer(n)
			fmt.Fprintf(&stdout, "player %d np %d gained %d spent %d backlinks %d\n", n, p.noblePoints, p.npGained, p.npSpent, b2i(w.pledgeBacklinks))
			helperList(&stdout, "units", p.units)
			fmt.Fprintf(&stdout, "output %d %d %d %d %d touched %d %d\n", b2i(testBit(p.output, 10)), b2i(testBit(p.output, 200)), b2i(testBit(p.output, 201)), b2i(testBit(p.output, 1000)), b2i(testBit(p.output, 1001)), b2i(testBit(p.locs, 9000)), b2i(testBit(p.locs, 9001)))
			for _, q := range p.orders {
				fmt.Fprintf(&stdout, "orders %d:", q.unit)
				for _, s := range q.l {
					fmt.Fprintf(&stdout, " [%s]", s)
				}
				fmt.Fprintln(&stdout)
			}
		case "readlog":
			n := args.num(t)
			clib.FlushAll() // Preserve the C observer's intentional fflush(NULL).
			fmt.Fprintf(&stdout, "log %d:\n", n)
			data, err := os.ReadFile(filepath.Join(w.libdir, "log", fmt.Sprint(n)))
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			stdout.Write(data)
		case "readfile":
			path := args.word()
			st, err := os.Stat(path)
			if err != nil {
				if !os.IsNotExist(err) && !errors.Is(err, syscall.ENOTDIR) {
					t.Fatal(err)
				}
				fmt.Fprintf(&stdout, "file bytes %s missing\n", path)
				break
			}
			if !st.Mode().IsRegular() || st.Size() > 65536 {
				t.Fatalf("readfile requires a regular file of at most 65536 bytes: %s", path)
			}
			// Borrow no production streams; observe disk without flushing caches.
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			data, err := io.ReadAll(io.LimitReader(f, 65537))
			closeErr := f.Close()
			if err != nil || closeErr != nil || len(data) > 65536 {
				t.Fatalf("readfile %s: read %v, close %v, size %d", path, err, closeErr, len(data))
			}
			fmt.Fprintf(&stdout, "file bytes %s %d:\n", path, len(data))
			stdout.Write(data)
			fmt.Fprint(&stdout, "\nend file bytes\n")
		case "load":
			fmt.Fprintf(&stdout, "load %d\n", b2i(w.loadCommand(c)))
		case "do":
			w.doCommand(c)
		case "finish":
			fmt.Fprintf(&stdout, "finish %d\n", w.finishCommand(c))
		case "interrupt":
			fmt.Fprintf(&stdout, "interrupt %d\n", cmdTbl[c.cmd].interrupt(w, c))
		case "init_wait":
			w.initWaitList()
		case "init_collect":
			w.initCollectList()
		case "poll_wait":
			w.checkAllWaits()
		case "clear_second":
			w.clearSecondWaits()
		case "wait_add":
			w.waitList.Append(args.num(t))
		case "lists":
			helperList(&stdout, "wait_list", w.waitList)
			helperList(&stdout, "collectors", w.collectors)
			helperList(&stdout, "second_wait_list", w.secondWaitList)
			fmt.Fprintf(&stdout, "second_wait %d\n", c.secondWait)
		case "done":
			w.commandDone(c)
		case "shift":
			w.cmdShift(c)
		case "wait", "wait_check":
			if op == "wait" {
				clearWaitParse(c)
				text := line
				if !w.olyParse(c, &text) {
					t.Fatal("wait parse failed")
				}
			}
			ret := w.checkWaitConditions(c)
			if ret == "" {
				ret = "NULL"
			}
			fmt.Fprintf(&stdout, "wait result [%s]\n", ret)
		case "wait_parse":
			ret := w.parseWaitArgs(c)
			if ret == "" {
				ret = "NULL"
			}
			fmt.Fprintf(&stdout, "wait parse [%s]\n", ret)
			waitDump()
		case "wait_clear":
			clearWaitParse(c)
			waitDump()
		case "wait_dump":
			waitDump()
		case "signal":
			n := args.num(t)
			w.flags = append(w.flags, &flagEnt{who: int32(n), flag: strings.TrimLeft(args.rest, " \t")})
		case "name", "plural":
			n := args.num(t)
			s := strings.TrimLeft(args.rest, " \t")
			if op == "name" {
				w.setName(n, s)
			} else {
				w.pItem(n).pluralName = &s
			}
		case "identity":
			n, name, email := args.num(t), args.word(), args.word()
			w.pPlayer(n).fullName, w.pPlayer(n).email = &name, &email
		case "entity":
			n, k, sk := args.num(t), args.num(t), args.num(t)
			w.allocBox(n, k, sk)
		case "link":
			n, a := args.num(t), args.num(t)
			if old := w.loc(a); old != 0 {
				w.pLocInfo(old).hereList.RemValue(a)
			}
			w.pLocInfo(a).where = int32(n)
			w.pLocInfo(n).hereList.Append(a)
		case "inventory":
			n, a, b := args.num(t), args.num(t), args.num(t)
			w.bx[n].items.Append(&itemEnt{item: int32(a), qty: int32(b)})
		case "field":
			n, field, a := args.num(t), args.word(), args.num(t)
			switch field {
			case "lord":
				w.pChar(n).unitLord = int32(a)
			case "format":
				w.pPlayer(n).format = int8(a)
			case "last_order":
				w.pPlayer(n).lastOrderTurn = int32(a)
			case "quest":
				w.pSubloc(n).questLate = int8(a)
			case "link_to":
				w.pSubloc(n).linkTo.Append(a)
			case "link_from":
				w.pSubloc(n).linkFrom.Append(a)
			case "link_open":
				w.pSubloc(n).linkOpen = int8(a)
			case "link_when":
				w.pSubloc(n).linkWhen = int8(a)
			case "rank":
				w.pChar(n).rank = int8(a)
			case "mine_delay":
				w.pMisc(n).mineDelay = int8(a)
			case "death_turn":
				w.pChar(n).deathTime.turn = int16(a)
			case "created_death":
				w.pChar(created).deathTime.turn = int16(a)
			case "melt":
				w.pChar(n).meltMe = int8(a)
			case "sick":
				w.pChar(n).sick = int8(a)
			case "char_attack":
				w.pChar(n).attack = int16(a)
			case "char_defense":
				w.pChar(n).defense = int16(a)
			case "char_missile":
				w.pChar(n).missile = int16(a)
			case "behind":
				w.pChar(n).behind = int8(a)
			case "break":
				w.pChar(n).breakPoint = int8(a)
			case "vision_protect":
				w.pMagic(n).visProtect = int8(a)
			case "project_cast":
				w.pMagic(n).projectCast = int32(a)
			case "quick_cast":
				w.pMagic(n).quickCast = int16(a)
			case "magician":
				w.pMagic(n).magician = int8(a)
			case "max_aura":
				w.pMagic(n).maxAura = int32(a)
			case "cur_aura":
				w.pMagic(n).curAura = int32(a)
			case "ability_shroud":
				w.pMagic(n).abilityShroud = int16(a)
			case "hinder":
				w.pMagic(n).hinderMeditation = int8(a)
			case "hide_mage":
				w.pMagic(n).hideMage = int8(a)
			case "shroud":
				w.pLoc(n).shroud = int16(a)
			case "fresh":
				w.pChar(n).freshHire = int8(a)
			case "recent":
				w.pSubloc(n).recentLoot = int8(a)
			case "relic_decay":
				w.pItemMagic(n).relicDecay = int16(a)
			case "aura_bonus":
				w.pItemMagic(n).auraBonus = int16(a)
			case "auraculum":
				w.pMagic(n).auraculum = int32(a)
			case "item_aura":
				w.pItemMagic(n).aura = int16(a)
			case "creator":
				w.pItemMagic(n).creator = int32(a)
			case "curse":
				w.pItemMagic(n).curseLoyalty = int8(a)
			case "token":
				w.pMagic(n).token = int32(a)
			case "may_study":
				w.pItemMagic(n).mayStudy.Append(a)
			case "may_use":
				w.pItemMagic(n).mayUse.Append(a)
			case "required":
				w.pSkill(n).requiredSkill = int32(a)
			case "noble_item":
				w.pChar(n).unitItem = int32(a)
			case "prev":
				w.pChar(n).prevLord = int32(a)
			case "rate":
				w.pChar(n).loyRate = int32(a)
			case "pledge":
				w.pMagic(n).pledge = int32(a)
			case "castle":
				w.pMisc(n).garrCastle = int32(a)
			case "unit":
				w.pPlayer(n).units.Append(a)
			case "np":
				w.pPlayer(n).noblePoints = int16(a)
			case "first_turn":
				w.pPlayer(n).firstTurn = int32(a)
			case "sent":
				w.pPlayer(n).sentOrders = int8(a)
			case "unformed":
				w.pPlayer(n).unformed.Append(a)
			case "summer":
				w.pSubloc(n).summerFlag = int8(a)
			case "ship_moving":
				w.pSubloc(n).moving = int32(a)
			case "uldim":
				w.pSubloc(n).uldimFlag = int8(a)
			case "skill_np":
				w.pSkill(n).npReq = int32(a)
			case "skill_time":
				w.pSkill(n).timeToLearn = int32(a)
			case "skill_no_exp":
				w.pSkill(n).noExp = int32(a)
			case "unique":
				w.pItem(n).whoHas = int32(a)
			case "token_num":
				w.pItemMagic(n).tokenNum = int8(a)
			case "token_ni":
				w.pItemMagic(n).tokenNi = int32(a)
			case "lore":
				w.pItemMagic(n).lore = int32(a)
			case "use_key":
				w.pItemMagic(n).useKey = int8(a)
			case "man":
				w.pItem(n).isManItem = int8(a)
			case "animal":
				w.pItem(n).animal = int8(a)
			case "attack":
				w.pItem(n).attack = int16(a)
			case "defense":
				w.pItem(n).defense = int16(a)
			case "missile":
				w.pItem(n).missile = int16(a)
			case "attack_bonus":
				w.pItemMagic(n).attackBonus = int8(a)
			case "defense_bonus":
				w.pItemMagic(n).defenseBonus = int8(a)
			case "missile_bonus":
				w.pItemMagic(n).missileBonus = int8(a)
			case "fly_cap":
				w.pItem(n).flyCap = int16(a)
			case "default_garrison":
				w.pMagic(n).defaultGarr = int8(a)
			case "weight":
				w.pItem(n).weight = int16(a)
			case "land_cap":
				w.pItem(n).landCap = int16(a)
			case "ride_cap":
				w.pItem(n).rideCap = int16(a)
			case "guard":
				w.pChar(n).guard = int8(a)
			case "queued_lore":
				w.pPlayer(n).deliverLore.Append(a)
			case "touched":
				setBit(&w.pPlayer(n).locs, a)
			case "weather_seen":
				setBit(&w.pPlayer(n).weatherSeen, a)
			case "weather_mage":
				w.pMagic(n).knowsWeather = int8(a)
			case "summoner":
				w.pMisc(n).summonedBy = int32(a)
			case "npc_home":
				w.pMisc(n).npcHome = int32(a)
			case "npc_cookie":
				w.pMisc(n).npcCookie = int32(a)
			case "npc_created":
				w.pMisc(n).npcCreated = int32(a)
			case "capturable":
				w.pItem(n).capturable = int8(a)
			case "release":
				w.pMagic(n).swearOnRelease = int8(a)
			case "health":
				w.pChar(n).health = int8(a)
			case "prisoner":
				w.pChar(n).prisoner = int8(a)
			case "npc":
				w.pChar(n).npcProg = int8(a)
			case "npc_dir":
				w.pMisc(n).npcDir = int8(a)
			case "storm_move":
				w.pMisc(n).stormMove = int32(a)
			case "bind_storm":
				w.pMisc(n).bindStorm = int32(a)
			case "bound_storm":
				w.pSubloc(n).boundStorms.Append(a)
			case "near_city":
				w.pSubloc(n).nearCities.Append(a)
			case "teaches":
				w.pSubloc(n).teaches.Append(a)
			case "damage":
				w.pSubloc(n).damage = uint8(a)
			case "capacity":
				w.pSubloc(n).capacity = int32(a)
			case "structure_defense":
				w.pSubloc(n).defense = int32(a)
			case "effort_given":
				w.pSubloc(n).effortGiven = int32(a)
			case "effort_required":
				w.pSubloc(n).effortRequired = int32(a)
			case "ram":
				w.pSubloc(n).galleyRam = int8(a)
			case "board_fee":
				w.pMagic(n).fee = int32(a)
			case "shaft_depth":
				w.pSubloc(n).shaftDepth = int16(a)
			case "loyal":
				w.pChar(n).loyKind = int8(a)
			case "moving":
				w.pChar(n).moving = int32(a)
			case "time_flying":
				w.pChar(n).timeFlying = int8(a)
			case "hidden":
				w.pMagic(n).hideSelf = int8(a)
			case "horn":
				w.pMagic(n).ferryFlag = int8(a)
			case "subkind":
				w.bx[n].skind = int8(a)
			case "vulnerable":
				w.pMisc(n).onlyVulnerable = int32(a)
			case "strength":
				w.pMisc(n).stormStr = int16(a)
			case "castle_level":
				w.pSubloc(n).castleLev = int8(a)
			case "major":
				w.pSubloc(n).major = int8(a)
			case "civ":
				w.pLoc(n).civ = int8(a)
			case "pillage":
				w.pSubloc(n).loot = int8(a)
			case "opium":
				w.pSubloc(n).opiumEcon = int32(a)
			case "double":
				w.pMisc(n).opiumDouble = int8(a)
			case "temp":
				w.bx[n].temp = int32(a)
			case "output_order":
				w.bx[n].outputOrder = int32(a)
			case "garr_tax":
				w.pMisc(n).garrTax = int32(a)
			case "garr_forward":
				w.pMisc(n).garrForward = int32(a)
			case "safe":
				w.pSubloc(n).safe = int8(a)
			case "loc_hidden":
				w.pLoc(n).hidden = int8(a)
			case "barrier":
				w.pLoc(n).barrier = int16(a)
			case "hostile":
				w.pDisp(n).hostile.Append(a)
			case "neutral":
				w.pDisp(n).neutral.Append(a)
			case "defend":
				w.pDisp(n).defend.Append(a)
			case "contact":
				w.pChar(n).contact.Append(a)
			case "garr_host":
				w.pMisc(n).garrHost.Append(a)
			case "known":
				setBit(&w.pPlayer(n).known, a)
			default:
				t.Fatalf("unknown field %s", field)
			}
		case "predicate":
			name := args.word()
			n := args.num(t)
			var a, b int
			if args.rest != "" {
				a = args.num(t)
			}
			if args.rest != "" {
				b = args.num(t)
			}
			var result int
			switch name {
			case "hostile":
				result = b2i(w.isHostile(n, a))
			case "contacted":
				result = b2i(w.contacted(n, a))
			case "char_where":
				result = b2i(w.charWhere(n, a, b))
			case "char_here":
				result = b2i(w.charHere(n, a))
			case "known":
				result = b2i(w.testKnown(n, a))
			case "loc_hidden":
				result = w.locHidden(n)
			case "has_item":
				result = int(w.hasItem(n, a))
			case "safe":
				result = b2i(w.inSafeNow(n))
			case "province":
				result = w.province(n)
			case "owner":
				result = w.buildingOwner(n)
			case "weather":
				result = int(w.weatherHere(n, a))
			case "really_hidden":
				result = b2i(w.charReallyHidden(n))
			case "npc":
				result = b2i(w.isNPC(n))
			case "vulnerable":
				result = int(w.onlyDefeatable(n))
			case "char_alloc":
				result = b2i(w.rpChar(n) != nil)
			case "fort_covers":
				result = w.fortCovers(n)
			case "depth":
				result = w.locDepth(n)
			default:
				t.Fatalf("unknown predicate %s", name)
			}
			fmt.Fprintf(&stdout, "predicate %d\n", result)
		case "initial":
			w.initialCommandLoad()
		case "start":
			w.startPhase()
		case "evening":
			w.eveningPhase()
		case "daily":
			w.dailyCommandLoop()
		case "players":
			w.processPlayerOrders()
		case "auto":
			w.checkAllAutoAttacks()
		case "empty_chars":
			w.changeBoxKind(1000, tDeleted)
			w.changeBoxKind(1001, tDeleted)
		case "min":
			fmt.Fprintf(&stdout, "min %d\n", w.minPriReady())
		case "sort_load":
			w.sortLoadQueue(w.loadQ[args.num(t)])
		case "sort_run":
			w.sortRunQueue(w.runQ)
		case "phase":
			fmt.Fprintf(&stdout, "phase evening %d auto %d second %d\n", w.evening, b2i(w.autoAttackFlag), w.secondWaitList.Len())
		case "second":
			c.secondWait = 1
			w.secondWaitList.Append(int(c.who))
		case "select":
			c = w.rpCommand(args.num(t))
		case "place":
			n := args.num(t)
			if !w.validBox(n) {
				w.allocBox(n, tLoc, subPlain)
			}
			w.pLocInfo(int(c.who)).where = int32(n)
			w.pLocInfo(n).hereList.Append(int(c.who))
		case "resolve_wait":
			w.cmdWait = int32(w.findCommand("wait"))
		case "skill":
			know, experience := args.num(t), args.num(t)
			w.pChar(int(c.who)).skills.Append(&skillEnt{skill: skHideSelf, know: int8(know), experience: int16(experience)})
		case "parse":
			text := strings.TrimPrefix(line, "parse ")
			fmt.Fprintf(&stdout, "parse %d\n", b2i(w.olyParse(c, &text)))
		case "allow":
			text := strings.TrimPrefix(line, "allow ")
			allow := &text
			if text == "NULL" {
				allow = nil
			} else if text == "EMPTY" {
				text = ""
			}
			fmt.Fprintf(&stdout, "allow %d\n", b2i(w.checkAllow(c, allow)))
		case "who":
			n := args.num(t)
			c = w.pCommand(n)
			c.who, c.state = int32(n), stateDone
		case "run":
			c.wait = int32(args.num(t))
			c.poll, c.status, c.inhibitFinish = int8(args.num(t)), int8(args.num(t)), int8(args.num(t))
			c.pri = int8(cmdTbl[c.cmd].pri)
			w.setState(c, stateRun, int(c.pri))
		case "set":
			field, n := args.word(), args.num(t)
			switch field {
			case "immediate":
				w.immediate = n != 0
			case "see_all":
				w.immedSeeAll = int32(n)
			case "month_done":
				w.monthDone = int32(n)
			case "evening":
				w.evening = int32(n)
			case "flush":
				w.flushAlways = n != 0
			case "auto_quit":
				w.autoQuitTurns = int32(n)
			case "post":
				w.postHasBeenRun = int32(n)
			case "indep":
				w.indepPlayer = int32(n)
			case "gm":
				w.gmPlayer = int32(n)
			case "skill_player":
				w.skillPlayer = int32(n)
			case "faery_region":
				w.faeryRegion = int32(n)
			case "hades_region":
				w.hadesRegion = int32(n)
			case "cloud_region":
				w.cloudRegion = int32(n)
			case "under_region":
				w.underRegion = int32(n)
			case "tunnel_region":
				w.tunnelRegion = int32(n)
			case "nowhere_region":
				w.nowhereRegion = int32(n)
			case "nowhere_loc":
				w.nowhereLoc = int32(n)
			case "hades_player":
				w.hadesPlayer = int32(n)
			case "faery_player":
				w.faeryPlayer = int32(n)
			case "backlinks":
				w.pledgeBacklinks = n != 0
			case "garrison_pay":
				w.garrisonPay = int32(n)
			case "gold_garrison":
				w.goldGarrison = int32(n)
			case "show_day":
				w.showDay = n != 0
			case "day":
				w.sysclock.day = int16(n)
			case "epoch":
				w.sysclock.daysSinceEpoch = int32(n)
			case "turn":
				w.sysclock.turn = int16(n)
			case "priority":
				w.curPri = int32(n)
			case "status":
				c.status = int8(n)
			case "wait":
				c.wait = int32(n)
			case "poll":
				c.poll = int8(n)
			case "debug":
				c.debug = int8(n)
			case "days":
				c.daysExecuting = int32(n)
			case "restrict":
				w.pMisc(int(c.who)).cmdAllow = int8(n)
			case "lord":
				w.pChar(int(c.who)).unitLord = int32(n)
			case "moving":
				w.pChar(int(c.who)).moving = int32(n)
			case "parent":
				w.pLocInfo(int(c.who)).where = int32(n)
			case "kind":
				w.bx[c.who].kind = int8(n)
			case "waitcmd":
				w.cmdWait = int32(n)
			case "count":
				w.pPlayer(w.player(int(c.who))).cmdCount = int16(n)
			case "prisoner":
				w.pChar(int(c.who)).prisoner = int8(n)
			case "health":
				w.pChar(int(c.who)).health = int8(n)
			case "second":
				c.secondWait = int8(n)
			case "cmdpri":
				c.pri = int8(n)
			case "state":
				c.state = int8(n)
			case "auto":
				w.autoAttackFlag = n != 0
			default:
				t.Fatalf("unknown field %q", field)
			}
		case "table":
			for i, row := range cmdTbl {
				fmt.Fprintf(&stdout, "handler %d %s %d %d %d\n", i, row.name, b2i(row.start != nil), b2i(row.finish != nil), b2i(row.interrupt != nil))
			}
		case "dump":
			pl := w.player(int(c.who))
			p, q := w.pPlayer(pl), w.pOrderHead(pl, int(c.who))
			fmt.Fprintf(&stdout, "cmd %d who %d state %d status %d wait %d poll %d pri %d cond %d fuzzy %d days %d debug %d inhibit %d cur %d\n",
				c.cmd, c.who, c.state, c.status, c.wait, c.poll, c.pri, c.conditional, c.fuzzy, c.daysExecuting, c.debug, c.inhibitFinish, w.curPri)
			fmt.Fprintf(&stdout, "args %d %d %d %d %d %d %d %d; parse", c.a, c.b, c.c, c.d, c.e, c.f, c.g, c.h)
			for _, s := range c.parse {
				fmt.Fprintf(&stdout, " [%s]", s)
			}
			var hidden int8
			if p := w.rpMagic(int(c.who)); p != nil {
				hidden = p.hideSelf
			}
			fmt.Fprintf(&stdout, "\nplayer %d count %d last %d format %d notab %d output %d hidden %d orders %d\n",
				pl, p.cmdCount, p.lastOrderTurn, p.format, p.notab, b2i(testBit(p.output, int(c.who))), hidden, q.l.Len())
			for i, q := range w.loadQ {
				helperList(&stdout, fmt.Sprintf("load_q %d", i), q)
			}
			helperList(&stdout, "run_q", w.runQ)
		case "log":
			clib.FlushAll() // actual C fixture calls fflush(NULL).
			fmt.Fprintln(&stdout, "log:")
			data, err := os.ReadFile(filepath.Join(w.libdir, "log", fmt.Sprint(w.player(int(c.who)))))
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			stdout.Write(data)
		case "null_do":
			w.doCommand(nil)
		case "null_finish":
			w.finishCommand(nil)
		case "null_load":
			w.loadCommand(nil)
		case "invalid_finish":
			w.finishCommand(&command{who: 9999})
		default:
			t.Fatalf("unknown operation %q", op)
		}
	}
	w.closeLogfile()
	w.closeTimes()
	return ""
}

// Unported production routines must stop, not be mistaken for C NULL or
// success. This is a port-boundary contract, not a fabricated C recording.
func TestCommandUnportedStops(t *testing.T) {
	w := newWorld()
	for _, name := range []string{"collect", "attack"} {
		i := w.findCommand(name)
		for phase, fn := range map[string]func(*world, *command) int{"start": cmdTbl[i].start, "finish": cmdTbl[i].finish, "interrupt": cmdTbl[i].interrupt} {
			if fn == nil {
				continue
			}
			t.Run(name+"/"+phase, func(t *testing.T) {
				defer func() {
					if r := recover(); r != notImplemented("command "+name+" "+phase) {
						t.Fatalf("unported handler returned or panicked unexpectedly: %v", r)
					}
				}()
				fn(w, &command{})
			})
		}
	}
}

// Real transfer/interrupt orchestration must not turn unported non-NULL
// callbacks into success. This tests the explicit boundary, not C results.
func TestCommandInterruptionBoundaries(t *testing.T) {
	for _, name := range []string{"sail", "collect"} {
		for _, transfer := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/transfer=%v", name, transfer), func(t *testing.T) {
				w := newWorld()
				w.allocBox(200, tPlayer, subPlRegular)
				w.allocBox(201, tPlayer, subPlRegular)
				w.allocBox(1000, tChar, 0)
				w.allocBox(9000, tLoc, subPlain)
				w.pChar(1000).unitLord = 200
				w.setWhere(1000, 9000)
				w.pPlayer(200).units.Append(1000)
				w.queueOrder(200, 1000, "format 9")
				c := w.pCommand(1000)
				c.who, c.cmd, c.state = 1000, int32(w.findCommand(name)), stateRun
				defer func() {
					if r := recover(); r != notImplemented("command "+name+" interrupt") {
						t.Fatalf("interrupt boundary = %v", r)
					}
					if w.player(1000) != 200 || c.state != stateRun {
						t.Fatal("unported interruption falsely completed transfer/command")
					}
					_, queued := w.topOrder(200, 1000)
					if queued == transfer {
						t.Fatal("transfer must flush before interruption; standalone interrupt must not flush")
					}
				}()
				if transfer {
					w.setLord(1000, 201, loyUnchanged, 0)
				} else {
					w.interruptOrder(1000)
				}
			})
		}
	}
}

// These check deliberate production boundaries, not fabricated C output.
// An actual hostile target must stop before priority-three commands, but only
// AFTER real lower-priority handlers have mutated state/output. Healthy leaders
// without targets complete normally in the C-recorded scheduler cases.
func TestCommandSchedulerBoundaries(t *testing.T) {
	for _, health := range []int8{100, -1} {
		t.Run(fmt.Sprint(health), func(t *testing.T) {
			w := newWorld()
			w.libdir = t.TempDir()
			t.Cleanup(w.closeLogfile)
			if err := os.Mkdir(filepath.Join(w.libdir, "log"), 0o755); err != nil {
				t.Fatal(err)
			}
			w.allocBox(200, tPlayer, subPlRegular)
			w.allocBox(1000, tChar, 0)
			w.allocBox(9000, tLoc, subPlain)
			w.pChar(1000).unitLord = 200
			w.pChar(1000).health = health
			w.pLocInfo(1000).where = 9000
			w.pLocInfo(9000).hereList.Append(1000)
			w.allocBox(201, tPlayer, subPlRegular)
			w.allocBox(1001, tChar, 0)
			w.pChar(1001).unitLord = 201
			w.pLocInfo(1001).where = 9000
			w.pLocInfo(9000).hereList.Append(1001)
			w.pDisp(1000).hostile.Append(1001)
			w.pOrderHead(200, 1000).l.Append("stop")
			w.pOrderHead(200, 1000).l.Append("format 17")
			w.pOrderHead(200, 1000).l.Append("hide 0")
			w.initialCommandLoad()
			defer func() {
				if r := recover(); r != notImplemented("regular_combat") {
					t.Fatalf("combat boundary = %v", r)
				}
				c := w.rpCommand(1000)
				if p := w.rpPlayer(200); p.format != 17 || p.cmdCount != 2 || c.cmd != int32(w.findCommand("hide")) || c.state != stateLoad || !w.autoAttackFlag || w.curPri != 3 {
					t.Fatalf("auto-attack ran before low priorities or after priority 3: player=%+v command=%+v auto=%v cur=%d", p, c, w.autoAttackFlag, w.curPri)
				}
				clib.FlushAll() // Same pre-abort fflush(NULL) as the real C fixture.
				log, err := os.ReadFile(filepath.Join(w.libdir, "log", "200"))
				if err != nil || !strings.Contains(string(log), "> format 17\n") || !strings.Contains(string(log), "> [auto-attack "+boxCodeLess(1001)+"]\n") || strings.Contains(string(log), "> hide") {
					t.Fatalf("wrong output before combat boundary: %q, %v", log, err)
				}
			}()
			w.dailyCommandLoop()
		})
	}
	t.Run("trades/evening", func(t *testing.T) {
		w := newWorld()
		w.allocBox(200, tPlayer, subPlRegular)
		w.allocBox(1000, tChar, 0)
		w.allocBox(9000, tLoc, subPlain)
		w.pChar(1000).unitLord = 200
		w.pLocInfo(1000).where = 9000
		w.pLocInfo(9000).hereList.Append(1000)
		c := w.pCommand(1000)
		c.who, c.state = 1000, stateDone
		w.tradesToCheck.Append(1000)
		c.cmd = int32(w.findCommand("xyzzy"))
		c.wait, c.status = 1, 1
		w.setState(c, stateRun, 0)
		w.eveningPhase()
		if c.state != stateDone || c.wait != 0 || w.runQ.Len() != 0 || w.tradesToCheck.Len() != 0 {
			t.Fatalf("command/validation not completed: %+v, queue %v, trades %v", c, w.runQ, w.tradesToCheck)
		}
	})
}

// A garrison without a castle has no output recipients, but still reaches
// the real regular_combat boundary after auto_attack's output routing.
func TestCommandGarrisonAutoAttackBoundary(t *testing.T) {
	w := newWorld()
	var stdout, stderr strings.Builder
	w.stdout, w.stderr = &stdout, &stderr
	w.allocBox(200, tPlayer, subPlRegular)
	w.allocBox(201, tPlayer, subPlRegular)
	w.allocBox(1000, tChar, subGarrison)
	w.allocBox(1001, tChar, 0)
	w.allocBox(9000, tLoc, subPlain)
	w.pChar(1000).unitLord, w.pChar(1000).health = 200, 100
	w.pChar(1001).unitLord = 201
	w.setWhere(1000, 9000)
	w.setWhere(1001, 9000)
	w.checkAllAutoAttacks()
	w.pMisc(1000).garrHost.Append(1001)
	defer func() {
		if r := recover(); r != notImplemented("regular_combat") {
			t.Fatalf("garrison combat boundary = %v", r)
		}
		if stdout.Len() != 0 || stderr.Len() != 0 || w.rpPlayer(200) != nil {
			t.Fatal("garrison boundary fabricated report output or player mutations")
		}
	}()
	w.checkAllAutoAttacks()
}
