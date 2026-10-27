// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"os"
	"path/filepath"
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
	fmt.Fprintln(&stdout, lines[0])
	defer func() {
		status := 0
		switch r := recover().(type) {
		case nil:
		case assertion:
			fmt.Fprint(&stderr, r.Error())
			status = 134
		default:
			panic(r)
		}
		output = fmt.Sprintf("stdout:\n%sstderr:\n%s\nexit %d\n", stdout.String(), stderr.String(), status)
	}()
	for _, line := range lines[1:] {
		fmt.Fprintf(&stdout, "> %s\n", line)
		args := &caseLine{rest: line}
		switch op := args.word(); op {
		case "queue":
			w.pOrderHead(w.player(int(c.who)), int(c.who)).l.Append(strings.TrimPrefix(line, "queue "))
		case "load":
			fmt.Fprintf(&stdout, "load %d\n", b2i(w.loadCommand(c)))
		case "do":
			w.doCommand(c)
		case "finish":
			fmt.Fprintf(&stdout, "finish %d\n", w.finishCommand(c))
		case "done":
			w.commandDone(c)
		case "shift":
			w.cmdShift(c)
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
			case "show_day":
				w.showDay = n != 0
			case "day":
				w.sysclock.day = int16(n)
			case "turn":
				w.sysclock.turn = int16(n)
			case "priority":
				w.curPri = int32(n)
			case "status":
				c.status = int8(n)
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
	for _, name := range []string{"move", "wait", "collect", "attack"} {
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
