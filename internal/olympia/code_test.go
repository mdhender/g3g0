// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"hash"
	"hash/fnv"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/mdhender/g3g0/internal/lists"
)

// TestCodeMatchesC runs testdata/code-cases.txt the way the C driver in
// tools/code-compat/run.sh does and requires its output, line for line.
// testdata/code-c-output.txt is that driver's output from the pinned G3
// sources; regenerate it only with run.sh --update.
func TestCodeMatchesC(t *testing.T) {
	cases := readTestdata(t, "code-cases.txt")
	want := readTestdata(t, "code-c-output.txt")

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
		out.WriteString(runCodeSection(t, s))
	}

	wantLines := strings.Split(want, "\n")
	gotLines := strings.Split(out.String(), "\n")
	if len(gotLines) != len(wantLines) {
		t.Fatalf("line count = %d, want %d (including final empty lines)", len(gotLines), len(wantLines))
	}
	for i := range wantLines {
		if wantLines[i] != gotLines[i] {
			t.Fatalf("line %d:\n got %q\nwant %q", i+1, gotLines[i], wantLines[i])
		}
	}
}

// runCodeSection runs one section in a fresh world, as one run of the C
// driver, and returns what the driver prints, with the "exit STATUS" line
// run.sh adds.
func runCodeSection(t *testing.T, lines []string) (output string) {
	var out strings.Builder
	w := newWorld()
	w.stderr = &out // the driver sends stderr to stdout
	w.stdout = &out
	traceAbort := false

	out.WriteString(lines[0] + "\n")
	defer func() {
		switch r := recover().(type) {
		case nil:
			out.WriteString("exit 0\n")
		case assertion:
			out.WriteString(r.Error())
			if traceAbort {
				a, b := w.streams.MintAlloc(0, 1000000000), w.streams.MintAlloc(0, 1000000000)
				fmt.Fprintf(&out, "mint after assertion %d %d\n", a, b)
			}
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
		case "trace_abort":
			traceAbort = true
		case "new":
			fmt.Fprintf(&out, "new %d\n", w.newEnt(c.num(t), c.num(t)))
		case "scan":
			fmt.Fprintf(&out, "scan %d\n", w.rndAllocNum(c.num(t), c.num(t)))
		case "fill":
			low, high, except := c.num(t), c.num(t), c.num(t)
			for n := low; n <= high; n++ {
				if n != except {
					w.allocBox(n, tItem, 0)
				}
			}
		case "clock":
			w.sysclock.turn, w.sysclock.day = int16(c.num(t)), int16(c.num(t))
		case "master":
			n := c.num(t)
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
			fmt.Fprintf(&out, "legacy %d\n", w.rnd.Rnd(c.num(t), c.num(t)))
		case "city":
			fmt.Fprintf(&out, "city %d\n", w.streams.MintCity())
		case "password":
			fmt.Fprintf(&out, "password %d\n", w.streams.MintPassword(c.num(t), c.num(t), c.num(t)))
		case "shuffle":
			l := lists.IList{10, 11, 12, 13, 14, 15}
			w.streams.MintShuffle(l)
			out.WriteString("shuffle:")
			for _, n := range l {
				fmt.Fprintf(&out, " %d", n)
			}
			out.WriteString("\n")
		case "seedstate":
			m := w.rnd.MasterSeed()
			fmt.Fprintf(&out, "master %08x %08x %08x %08x\n", m[0], m[1], m[2], m[3])
			fmt.Fprintf(&out, "digest %x\n", w.rnd.Seed())
		case "alloc":
			n, k, sk := c.num(t), c.num(t), c.num(t)
			w.allocBox(n, k, sk)
		case "delete":
			w.deleteBox(c.num(t))
		case "kind":
			n := c.num(t)
			w.changeBoxKind(n, c.num(t))
		case "sub":
			n := c.num(t)
			w.changeBoxSubkind(n, c.num(t))
		case "chains":
			for k := range tMax {
				if w.kindFirst(k) != 0 {
					fmt.Fprintf(&out, "kind %d:", k)
					for i := w.kindFirst(k); i > 0; i = w.kindNext(i) {
						fmt.Fprintf(&out, " %d", i)
					}
					out.WriteString("\n")
				}
			}
			for k := range subMax {
				if w.subFirst(k) != 0 {
					fmt.Fprintf(&out, "sub %d:", k)
					for i := w.subFirst(k); i > 0; i = w.subNext(i) {
						fmt.Fprintf(&out, " %d", i)
					}
					out.WriteString("\n")
				}
			}
		case "loop":
			codeLoop(t, w, c, &out)
		case "usage":
			w.printBoxUsage()
		case "name":
			fmt.Fprintf(&out, "name [%s]\n", esc(w.name(c.num(t))))
		case "setname":
			n := c.num(t)
			w.setName(n, c.text())
		case "banner":
			n := c.num(t)
			w.setBanner(n, c.text())
			if d := w.rpMisc(n).display; d != nil {
				fmt.Fprintf(&out, "banner [%s]\n", esc(*d))
			} else {
				out.WriteString("banner NULL\n")
			}
		case "ifname":
			n := c.num(t)
			w.prettyDataFiles = c.num(t) != 0
			fmt.Fprintf(&out, "if_name [%s]\n", esc(w.ifName(n)))
		case "loc":
			fmt.Fprintf(&out, "loc %d\n", w.loc(c.num(t)))
		case "boxcode":
			n := c.num(t)
			fmt.Fprintf(&out, "box_code [%s] box_code_less [%s]\n", w.boxCode(n), boxCodeLess(n))
		case "code":
			s := intToCode(c.num(t))
			fmt.Fprintf(&out, "int_to_code [%s] code_to_int %d\n", s, codeToInt(s))
		case "int":
			fmt.Fprintf(&out, "code_to_int %d\n", codeToInt(c.text()))
		case "scode":
			fmt.Fprintf(&out, "scode %d\n", scode(c.text()))
		case "letter":
			out.WriteString("letter_val:")
			for _, b := range []byte(c.text()) {
				fmt.Fprintf(&out, " %d", letterVal(b, letters2))
			}
			out.WriteString("\n")
		case "codesum":
			codeSum(&out)
		default:
			t.Fatalf("unknown operation %q", line)
		}
	}
	return
}

// codeLoop is the driver's loop operation.
func codeLoop(t *testing.T, w *world, c *caseLine, out *strings.Builder) {
	what := c.word()
	k, del, an, ak, ask := 0, false, 0, 0, 0
	if what == "kind" || what == "sub" {
		k = c.num(t)
		switch c.word() {
		case "del":
			del = true
		case "add":
			an, ak, ask = c.num(t), c.num(t), c.num(t)
		}
	}

	out.WriteString("visit:")

	first := true
	body := func(i int) {
		fmt.Fprintf(out, " %d", i)
		if del {
			w.deleteBox(i)
		}
		if an != 0 && first {
			w.allocBox(an, ak, ask)
		}
		first = false
	}
	switch what {
	case "kind":
		for i := range w.loopKind(k) {
			body(i)
		}
	case "sub":
		for i := range w.loopSubkind(k) {
			body(i)
		}
	case "boxes":
		for i := range w.loopBoxes() {
			fmt.Fprintf(out, " %d", i)
		}
	case "locship":
		for i := range w.loopLocOrShip() {
			fmt.Fprintf(out, " %d", i)
		}
	}
	out.WriteString("\n")
}

// codeSum is the driver's codesum operation: FNV-1a checksums over the
// whole code space.
func codeSum(out *strings.Builder) {
	var sum hash.Hash64

	// int_to_code over the whole box space and past it, and back.
	sum = fnv.New64a()
	n := 0
	for l := -10; l <= maxBoxes+10; l++ {
		code := intToCode(l)
		back := codeToInt(code)
		fmt.Fprintf(sum, "%d %s %d\n", l, code, back)
		if back != l {
			fmt.Fprintf(out, "int_to_code(%d) = [%s] code_to_int = %d\n", l, code, back)
			n++
		}
	}
	fmt.Fprintf(out, "int_to_code %d..%d: %d do not convert back, sum %016x\n", -10, maxBoxes+10, n, sum.Sum64())

	// code_to_int and scode over every string of up to four bytes of alpha.
	const alpha = "abelzAEZ059[( -+\x80\xff"
	sum = fnv.New64a()
	n = 0
	for length := 0; length <= 4; length++ {
		total := 1
		for range length {
			total *= len(alpha)
		}
		for l := range total {
			s := make([]byte, length)
			v := l
			for i := range s {
				s[i] = alpha[v%len(alpha)]
				v /= len(alpha)
			}
			fmt.Fprintf(sum, "%d %d\n", codeToInt(string(s)), scode(string(s)))
			n++
		}
	}
	fmt.Fprintf(out, "strings: %d, sum %016x\n", n, sum.Sum64())

	// code_to_int over every letter-letter/digit-digit[-digit] shape.
	const shape1 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const shape2 = shape1 + "0123456789"
	sum = fnv.New64a()
	n = 0
	for a := range len(shape1) {
		for b := range len(shape2) {
			for c := range 10 {
				for d := -1; d < 10; d++ {
					s := []byte{shape1[a], shape2[b], byte('0' + c)}
					if d >= 0 {
						s = append(s, byte('0'+d))
					}
					fmt.Fprintf(sum, "%d\n", codeToInt(string(s)))
					n++
				}
			}
		}
	}
	fmt.Fprintf(out, "shapes: %d, sum %016x\n", n, sum.Sum64())
}

// caseLine reads the arguments of one operation the way the C driver
// does.
type caseLine struct {
	rest string
}

// word is sscanf("%s"): the next word, skipping white space.
func (c *caseLine) word() string {
	c.rest = strings.TrimLeft(c.rest, " \t")
	i := strings.IndexAny(c.rest, " \t")
	if i < 0 {
		i = len(c.rest)
	}
	w := c.rest[:i]
	c.rest = c.rest[i:]
	return w
}

// num is strtol(pos, &pos, 0) on the next word.
func (c *caseLine) num(t *testing.T) int {
	w := c.word()
	n, err := strconv.ParseInt(w, 0, 64)
	if err != nil {
		t.Fatalf("bad number %q", w)
	}
	return int(n)
}

// text is the rest of the line after one space, with \xHH and \\ decoded.
func (c *caseLine) text() string {
	s := strings.TrimPrefix(c.rest, " ")
	var b []byte
	for i := 0; i < len(s); i++ {
		switch {
		case s[i] == '\\' && i+3 < len(s) && s[i+1] == 'x':
			v, _ := strconv.ParseUint(s[i+2:i+4], 16, 8)
			b = append(b, byte(v))
			i += 3
		case s[i] == '\\' && i+1 < len(s) && s[i+1] == '\\':
			b = append(b, '\\')
			i++
		default:
			b = append(b, s[i])
		}
	}
	return string(b)
}

// esc is the driver's esc(): bytes outside ' '..'~', and '\', as \xHH.
func esc(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if c := s[i]; c < ' ' || c > '~' || c == '\\' {
			fmt.Fprintf(&b, "\\x%02x", c)
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}
