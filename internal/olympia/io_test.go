// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/mdhender/g3g0/internal/lists"
)

// TestIOMatchesC runs testdata/io-cases.txt the way the C driver in
// tools/io-compat/run.sh does and requires its output, line for line.
// testdata/io-c-output.txt is that driver's output from the pinned G3
// sources; regenerate it only with run.sh --update.
func TestIOMatchesC(t *testing.T) {
	cases := readTestdata(t, "io-cases.txt")
	want := readTestdata(t, "io-c-output.txt")

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
		t.Chdir(t.TempDir()) // each section runs in a fresh directory
		out.WriteString(runIOSection(t, s))
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

// ioState is what the C driver holds for the helpers to scan into and
// print from.
type ioState struct {
	l      lists.IList
	k      sparse
	skills skillEntsList
	items  itemEntsList
	trades tradesList
	reqs   reqEntsList
	player entityPlayer
	ot     olytime
	gfp    *bufio.Reader
}

// runIOSection runs one section in a fresh world, as one run of the C
// driver, and returns what the driver prints, with the "exit STATUS" line
// run.sh adds. A nil pointer dereference is C's NULL dereference, which
// crashes with SIGSEGV (status 139).
func runIOSection(t *testing.T, lines []string) (output string) {
	var out strings.Builder
	w := newWorld()
	w.stderr = &out // the driver sends stderr to stdout
	var st ioState

	out.WriteString(lines[0] + "\n")
	defer func() {
		r := recover()
		if re, ok := r.(runtime.Error); ok && strings.Contains(re.Error(), "nil pointer dereference") {
			out.WriteString("exit 139\n")
		} else {
			switch r := r.(type) {
			case nil:
				out.WriteString("exit 0\n")
			case assertion:
				out.WriteString(r.Error())
				out.WriteString("exit 134\n")
			default:
				panic(r)
			}
		}
		output = out.String()
	}()

	showStr := func(what string, s *string) {
		if s != nil {
			fmt.Fprintf(&out, "%s [%s] len %d\n", what, esc(*s), len(*s))
		} else {
			fmt.Fprintf(&out, "%s NULL\n", what)
		}
	}
	field := func() string { return tString(*w.line) }

	for _, line := range lines[1:] {
		fmt.Fprintf(&out, "> %s\n", line)
		c := &caseLine{rest: line}
		op := c.word()
		switch op {
		case "alloc":
			n, k, sk := c.cnum(t), c.cnum(t), c.cnum(t)
			w.allocBox(n, k, sk)
		case "setname":
			n := c.cnum(t)
			w.setName(n, c.ctext())
		case "namen":
			n, k := c.cnum(t), c.cnum(t)
			w.setName(n, strings.Repeat(c.ctext(), k))
		case "pretty":
			w.prettyDataFiles = c.cnum(t) != 0
		case "extbox":
			w.extBoxnum = int32(c.cnum(t))
		case "put", "putn":
			name := c.word()
			k := 1
			if op == "putn" {
				k = c.cnum(t)
			}
			f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
			if err != nil {
				t.Fatal(err)
			}
			f.WriteString(strings.Repeat(c.text(), k))
			f.Close()
		case "readfile":
			fmt.Fprintf(&out, "readfile %d\n", b2i(w.readfile(c.word())))
		case "closefile":
			w.closefile(c.word())
		case "readlin":
			showStr("readlin", w.readlin())
		case "readlin_ew":
			showStr("readlin_ew", w.readlinEw())
		case "skip":
			for k := c.cnum(t); k > 0; k-- {
				w.readlin()
			}
		case "drain":
			k, last := 0, ""
			for s := w.readlin(); s != nil; s = w.readlin() {
				k, last = k+1, *s
			}
			fmt.Fprintf(&out, "drain %d lines, last [%s]\n", k, esc(last))
		case "fopen":
			if f, err := os.Open(c.word()); err != nil {
				st.gfp = nil
				out.WriteString("fopen NULL\n")
			} else {
				t.Cleanup(func() { f.Close() })
				st.gfp = bufio.NewReader(f)
				out.WriteString("fopen ok\n")
			}
		case "getlin":
			showStr("getlin", getlin(st.gfp))
		case "getlin_ew":
			showStr("getlin_ew", getlinEw(st.gfp))
		case "ew":
			s := eatLeadingTrailingWhitespace(c.ctext())
			showStr("ew", &s)
		case "hash":
			s := c.ctext()
			fmt.Fprintf(&out, "linehash %d t_string [%s]\n", linehash(s), esc(tString(s)))
		case "next":
			w.advance()
			if w.line != nil && iswhite(at(*w.line, 0)) {
				s := (*w.line)[1:]
				w.line = &s
			}
			if w.line != nil {
				fmt.Fprintf(&out, "line [%s] hash %d t [%s]\n", esc(*w.line), linehash(*w.line), esc(tString(*w.line)))
			} else {
				out.WriteString("line NULL\n")
			}
		case "box_scan":
			fmt.Fprintf(&out, "box_scan %d\n", w.boxScan(field()))
		case "boxlist_scan":
			w.boxlistScan(field(), c.cnum(t), &st.l)
			showIList(&out, "L:", st.l)
		case "boxlist0_scan":
			w.boxlist0Scan(field(), c.cnum(t), &st.l)
			showIList(&out, "L:", st.l)
		case "ilist_scan":
			w.ilistScan(field(), &st.l)
			showIList(&out, "L:", st.l)
		case "known_scan":
			w.knownScan(field(), &st.k, c.cnum(t))
			showIList(&out, "K:", st.k)
		case "admit_scan":
			w.admitScan(field(), c.cnum(t), &st.player)
			st.showAdmits(&out)
		case "skill_scan":
			w.skillListScan(field(), &st.skills, c.cnum(t))
			st.showSkills(&out)
		case "item_scan":
			w.itemListScan(field(), &st.items, c.cnum(t))
			st.showItems(&out)
		case "trade_scan":
			w.tradeListScan(field(), &st.trades, c.cnum(t))
			st.showTrades(&out)
		case "req_scan":
			w.reqListScan(field(), &st.reqs, c.cnum(t))
			st.showReqs(&out)
		case "olytime_scan":
			olytimeScan(field(), &st.ot)
			st.showOlytime(&out)
		case "box_print":
			n := c.cnum(t)
			w.boxPrint(&out, c.ctext(), n)
		case "boxlist_print":
			w.boxlistPrint(&out, c.ctext(), st.l)
		case "boxlist0_print":
			w.boxlist0Print(&out, c.ctext(), st.l)
		case "ilist_print":
			ilistPrint(&out, c.ctext(), st.l)
		case "known_print":
			w.knownPrint(&out, c.ctext(), st.k)
			showIList(&out, "K:", st.k)
		case "admit_print":
			w.admitPrint(&out, &st.player)
		case "skill_print":
			w.skillListPrint(&out, c.ctext(), st.skills)
		case "item_print":
			w.itemListPrint(&out, c.ctext(), st.items)
		case "trade_print":
			w.tradeListPrint(&out, c.ctext(), st.trades)
		case "req_print":
			w.reqListPrint(&out, c.ctext(), st.reqs)
		case "olytime_print":
			olytimePrint(&out, c.ctext(), &st.ot)
		case "set":
			st.l.Clear()
			c.nums(t, &st.l)
		case "setk":
			st.k.Clear()
			c.nums(t, &st.k)
		case "skill":
			st.skills.Append(&skillEnt{
				skill:       int32(c.cnum(t)),
				know:        int8(c.cnum(t)),
				daysStudied: int32(c.cnum(t)),
				experience:  int16(c.cnum(t)),
			})
		case "item":
			st.items.Append(&itemEnt{item: int32(c.cnum(t)), qty: int32(c.cnum(t))})
		case "trade":
			st.trades.Append(&trade{
				kind:      int32(c.cnum(t)),
				item:      int32(c.cnum(t)),
				qty:       int32(c.cnum(t)),
				cost:      int32(c.cnum(t)),
				cloak:     int32(c.cnum(t)),
				haveLeft:  int32(c.cnum(t)),
				monthProd: int32(c.cnum(t)),
				expire:    int32(c.cnum(t)),
			})
		case "req":
			st.reqs.Append(&reqEnt{item: int32(c.cnum(t)), qty: int32(c.cnum(t)), consume: int8(c.cnum(t))})
		case "admit":
			e := &admit{targ: int32(c.cnum(t)), sense: int32(c.cnum(t))}
			c.nums(t, &e.l)
			st.player.admits.Append(e)
		case "olytime":
			st.ot = olytime{turn: int16(c.cnum(t)), day: int16(c.cnum(t)), daysSinceEpoch: int32(c.cnum(t))}
		case "zero":
			li := locInfo{where: int32(c.cnum(t))}
			if c.cnum(t) != 0 {
				li.hereList.Append(1)
				li.hereList.Clear()
			}
			fmt.Fprintf(&out, "zero_check %d\n", b2i(zeroCheck(&li)))
		default:
			t.Fatalf("unknown operation %q", line)
		}
	}
	return
}

// cnum is num() stored in a C int, as the driver stores it.
func (c *caseLine) cnum(t *testing.T) int {
	return int(int32(c.num(t)))
}

// ctext is text() as the driver passes it to C: up to its first NUL byte.
func (c *caseLine) ctext() string {
	return cString([]byte(c.text()))
}

// nums appends the rest of the line's numbers to l, each stored in a C
// int.
func (c *caseLine) nums(t *testing.T, l *lists.IList) {
	for strings.TrimLeft(c.rest, " \t") != "" {
		l.Append(c.cnum(t))
	}
}

func showIList(out io.Writer, what string, l lists.IList) {
	io.WriteString(out, what)
	if l == nil {
		io.WriteString(out, " NULL")
	}
	for _, n := range l {
		fmt.Fprintf(out, " %d", n)
	}
	io.WriteString(out, "\n")
}

func (st *ioState) showSkills(out io.Writer) {
	io.WriteString(out, "skills:")
	for _, e := range st.skills {
		fmt.Fprintf(out, " [%d %d %d %d]", e.skill, e.know, e.daysStudied, e.experience)
	}
	io.WriteString(out, "\n")
}

func (st *ioState) showItems(out io.Writer) {
	io.WriteString(out, "items:")
	for _, e := range st.items {
		fmt.Fprintf(out, " [%d %d]", e.item, e.qty)
	}
	io.WriteString(out, "\n")
}

func (st *ioState) showTrades(out io.Writer) {
	io.WriteString(out, "trades:")
	for _, e := range st.trades {
		fmt.Fprintf(out, " [%d %d %d %d %d %d %d %d who %d]", e.kind, e.item, e.qty, e.cost,
			e.cloak, e.haveLeft, e.monthProd, e.expire, e.who)
	}
	io.WriteString(out, "\n")
}

func (st *ioState) showReqs(out io.Writer) {
	io.WriteString(out, "reqs:")
	for _, e := range st.reqs {
		fmt.Fprintf(out, " [%d %d %d]", e.item, e.qty, e.consume)
	}
	io.WriteString(out, "\n")
}

func (st *ioState) showAdmits(out io.Writer) {
	io.WriteString(out, "admits:")
	for _, e := range st.player.admits {
		fmt.Fprintf(out, " [%d %d", e.targ, e.sense)
		for _, n := range e.l {
			fmt.Fprintf(out, " %d", n)
		}
		io.WriteString(out, "]")
	}
	io.WriteString(out, "\n")
}

func (st *ioState) showOlytime(out io.Writer) {
	fmt.Fprintf(out, "olytime: turn %d day %d days %d\n", st.ot.turn, st.ot.day, st.ot.daysSinceEpoch)
}

// TestGetlinReadError checks that a read error ends getlin()'s input as
// end of file does: fgetc() returns EOF for both.
func TestGetlinReadError(t *testing.T) {
	r := bufio.NewReader(io.MultiReader(strings.NewReader("ab\ncd"), errReader{}))
	for _, want := range []string{"ab", "cd"} {
		if got := getlin(r); got == nil || *got != want {
			t.Fatalf("getlin: got %v, want %q", got, want)
		}
	}
	if got := getlin(r); got != nil {
		t.Fatalf("getlin after error: got %q, want nil", *got)
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("read error") }

// TestSscanD checks the conversions sscanD stops at.
func TestSscanD(t *testing.T) {
	for _, tc := range []struct {
		s    string
		n    int
		want [3]int32
	}{
		{"", 0, [3]int32{7, 7, 7}},
		{"1 2 3", 3, [3]int32{1, 2, 3}},
		{" \t\n1\v2\f3", 3, [3]int32{1, 2, 3}},
		{"1 x 3", 1, [3]int32{1, 7, 7}},
		{"-1 +2 -", 2, [3]int32{-1, 2, 7}},
		{"4294967297", 1, [3]int32{1, 7, 7}},
		{"99999999999999999999 -99999999999999999999", 2, [3]int32{-1, 0, 7}},
	} {
		got := [3]int32{7, 7, 7}
		if n := sscanD(tc.s, &got[0], &got[1], &got[2]); n != tc.n || got != tc.want {
			t.Errorf("sscanD(%q) = %d %v, want %d %v", tc.s, n, got, tc.n, tc.want)
		}
	}
}
