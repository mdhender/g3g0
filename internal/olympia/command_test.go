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
		out.WriteString(runCommandSection(t, section))
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

func runCommandSection(t *testing.T, lines []string) (output string) {
	var stdout, stderr strings.Builder
	w := newWorld()
	w.stdout, w.stderr = &stdout, &stderr
	w.libdir = t.TempDir()
	if err := os.Mkdir(filepath.Join(w.libdir, "log"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(w.libdir, "orders"), 0o755); err != nil {
		t.Fatal(err)
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
		case "trade":
			n := args.num(t)
			w.bx[n].trades.Append(&trade{kind: int32(args.num(t)), item: int32(args.num(t)), qty: int32(args.num(t)), cost: int32(args.num(t)), who: int32(args.num(t)), sort: int32(args.num(t))})
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
		case "set_known":
			w.setKnown(args.num(t), args.num(t))
		case "additem":
			w.addItem(args.num(t), args.num(t), args.num(t))
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
			fmt.Fprintf(&stdout, "log %d:\n", n)
			data, err := os.ReadFile(filepath.Join(w.libdir, "log", fmt.Sprint(n)))
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			stdout.Write(data)
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
			case "melt":
				w.pChar(n).meltMe = int8(a)
			case "sick":
				w.pChar(n).sick = int8(a)
			case "magician":
				w.pMagic(n).magician = int8(a)
			case "max_aura":
				w.pMagic(n).maxAura = int32(a)
			case "cur_aura":
				w.pMagic(n).curAura = int32(a)
			case "aura_bonus":
				w.pItemMagic(n).auraBonus = int16(a)
			case "token":
				w.pMagic(n).token = int32(a)
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
			case "skill_np":
				w.pSkill(n).npReq = int32(a)
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
			case "attack":
				w.pItem(n).attack = int16(a)
			case "defense":
				w.pItem(n).defense = int16(a)
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
			case "indep":
				w.indepPlayer = int32(n)
			case "gm":
				w.gmPlayer = int32(n)
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
				log, err := os.ReadFile(filepath.Join(w.libdir, "log", "200"))
				if err != nil || !strings.Contains(string(log), "> format 17\n") || !strings.Contains(string(log), "> [auto-attack "+boxCodeLess(1001)+"]\n") || strings.Contains(string(log), "> hide") {
					t.Fatalf("wrong output before combat boundary: %q, %v", log, err)
				}
			}()
			w.dailyCommandLoop()
		})
	}
	for _, boundary := range []string{"trades/evening"} {
		t.Run(boundary, func(t *testing.T) {
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
			want := notImplemented("check_validated_trades")
			defer func() {
				if r := recover(); r != want {
					t.Fatalf("boundary = %v, want %v", r, want)
				}
				if c.state != stateDone || c.wait != 0 || w.runQ.Len() != 0 {
					t.Fatalf("command not completed before boundary: %+v, queue %v", c, w.runQ)
				}
			}()
			c.cmd = int32(w.findCommand("xyzzy"))
			c.wait, c.status = 1, 1
			w.setState(c, stateRun, 0)
			w.eveningPhase()
		})
	}
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
