// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/mdhender/g3g0/internal/clib"
)

// TestBoxMatchesC runs testdata/box-cases.txt the way the C driver in
// tools/box-compat/run.sh does and requires its output, line for line.
// testdata/box-c-output.txt is that driver's output from the pinned G3
// sources; regenerate it only with run.sh --update.
func TestBoxMatchesC(t *testing.T) {
	cases := readTestdata(t, "box-cases.txt")
	want := readTestdata(t, "box-c-output.txt")

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
		out.WriteString(runBoxSection(t, s))
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

// runBoxSection runs one section in a fresh world, as one run of the C
// driver, and returns what the driver prints, with the "exit STATUS" line
// run.sh adds. A nil pointer dereference is C's NULL dereference, which
// crashes with SIGSEGV (status 139).
func runBoxSection(t *testing.T, lines []string) (output string) {
	var out strings.Builder
	w := newWorld()
	w.stderr = &out // the driver sends stderr to stdout
	var cmd command // the driver's static struct command

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
			fmt.Fprintf(&out, "%s [%s]\n", what, esc(*s))
		} else {
			fmt.Fprintf(&out, "%s NULL\n", what)
		}
	}
	showCmd := func(c *command) {
		fmt.Fprintf(&out, "cmd %d fuzzy %d conditional %d who %d use_skill %d\n",
			c.cmd, c.fuzzy, c.conditional, c.who, c.useSkill)
		showStr("line", c.line)
		showStr("parsed_line", c.parsedLine)
		out.WriteString("parse:")
		if c.parse == nil {
			out.WriteString(" NULL")
		}
		for _, s := range c.parse {
			fmt.Fprintf(&out, " [%s]", esc(s))
		}
		out.WriteString("\n")
	}

	for _, line := range lines[1:] {
		fmt.Fprintf(&out, "> %s\n", line)
		c := &caseLine{rest: line}
		op := c.word()
		switch op {
		case "alloc":
			n, k, sk := c.cnum(t), c.cnum(t), c.cnum(t)
			w.allocBox(n, k, sk)
		case "pretty":
			w.prettyDataFiles = c.cnum(t) != 0
		case "libdir":
			w.libdir = c.ctext()
		case "put", "putl":
			name := c.word()
			f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
			if err != nil {
				t.Fatal(err)
			}
			f.WriteString(c.text())
			if op == "putl" {
				f.WriteString("\n")
			}
			f.Close()
		case "readfile":
			fmt.Fprintf(&out, "readfile %d\n", b2i(w.readfile(c.word())))
		case "advance":
			w.advance()
		case "line":
			showStr("line", w.line)
		case "load":
			w.loadBox(c.cnum(t))
			showStr("line", w.line)
		case "loadall":
			loadAll(w, &out)
		case "save":
			n := c.cnum(t)
			w.saveBox(&out, n)
			if n > 0 && n < maxBoxes && w.bx[n] != nil {
				fmt.Fprintf(&out, "temp %d\n", w.bx[n].temp)
			}
		case "saveall":
			for n := range maxBoxes {
				if w.bx[n] != nil {
					w.saveBox(&out, n)
				}
			}
		case "cmd":
			showCmd(w.pCommand(c.cnum(t)))
		case "parse":
			s := c.ctext()
			k := w.olyParseCmd(&cmd, &s)
			fmt.Fprintf(&out, "oly_parse_cmd %d fuzzy_find %d\n", b2i(k), b2i(w.fuzzyFind))
			showStr("s", &s)
			showCmd(&cmd)
		case "find":
			k := w.findCommand(c.ctext())
			fmt.Fprintf(&out, "find_command %d fuzzy_find %d\n", k, b2i(w.fuzzyFind))
		case "p":
			what, n := c.word(), c.cnum(t)
			switch what {
			case "char":
				w.pChar(n)
			case "magic":
				w.pMagic(n)
			case "loc":
				w.pLoc(n)
			case "subloc":
				w.pSubloc(n)
			case "item":
				w.pItem(n)
			case "item_magic":
				w.pItemMagic(n)
			case "player":
				w.pPlayer(n)
			case "skill":
				w.pSkill(n)
			case "gate":
				w.pGate(n)
			case "misc":
				w.pMisc(n)
			case "disp":
				w.pDisp(n)
			case "command":
				w.pCommand(n)
			default:
				t.Fatalf("unknown p %q", line)
			}
		case "poke":
			what, n, k := c.word(), c.cnum(t), c.cnum(t)
			switch what {
			case "loc_next":
				w.pLoc(n).next = int32(k)
			case "prov_dest_empty":
				w.pLoc(n).provDest.Append(k)
				w.pLoc(n).provDest.Clear()
			case "mage_worked":
				w.pMagic(n).mageWorked = int8(k)
			case "one_turn_use":
				w.pItemMagic(n).oneTurnUse = int8(k)
			case "opium_double":
				w.pMisc(n).opiumDouble = int8(k)
			case "cmd":
				w.pCommand(n).cmd = int32(k)
			case "where":
				w.pLocInfo(n).where = int32(k)
			default:
				t.Fatalf("unknown poke %q", line)
			}
		default:
			t.Fatalf("unknown operation %q", line)
		}
	}
	return
}

// loadAll is the driver's load_all(): load_box() for each box line from
// line on, as read_boxes() does, but a line that is not a box number ends
// the file, where read_boxes() would skip it (tools/boxfile-compat tests
// readBoxes itself).
func loadAll(w *world, out io.Writer) {
	for w.line != nil {
		n := int(int32(clib.Atoi(*w.line)))
		if n > 0 {
			w.loadBox(n)
		} else {
			fmt.Fprintf(out, "stop [%s]\n", esc(*w.line))
			break
		}
	}
}

// TestBoxWorlds is the Go side of tools/box-compat/run.sh's check of real
// saved games. It runs only with G3G0_FULL=1 (see requireFull), and only
// when BOX_COMPAT_GAMES names a list of games, one per line: a lib
// directory, pretty_data_files (0 or 1) and an output prefix, separated by
// tabs. For each, as the C driver's -w mode does, it
// allocates every box from the box files' "n kind subkind" lines, loads
// every box file with loadBox and saves every box with saveBox, and writes
// the diagnostics to PREFIX.err and the saved boxes to PREFIX.out.
func TestBoxWorlds(t *testing.T) {
	requireFull(t)
	list := os.Getenv("BOX_COMPAT_GAMES")
	if list == "" {
		t.Skip("BOX_COMPAT_GAMES not set; run tools/box-compat/run.sh --full")
	}
	data, err := os.ReadFile(list)
	if err != nil {
		t.Fatal(err)
	}
	for line := range strings.SplitSeq(strings.TrimSuffix(string(data), "\n"), "\n") {
		f := strings.Split(line, "\t")
		if len(f) != 3 {
			t.Fatalf("bad line in %s: %q", list, line)
		}
		var stdout, stderr bytes.Buffer
		boxWorld(t, f[0], f[1] == "1", &stdout, &stderr)
		if err := os.WriteFile(f[2]+".out", stdout.Bytes(), 0o666); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(f[2]+".err", stderr.Bytes(), 0o666); err != nil {
			t.Fatal(err)
		}
	}
}

// boxWorld is the C driver's world().
func boxWorld(t *testing.T, dir string, pretty bool, stdout, stderr io.Writer) {
	w := newWorld()
	w.stderr = stderr
	w.prettyDataFiles = pretty
	w.libdir = dir

	files := []string{"loc", "item", "skill", "gate", "road", "ship", "unform", "misc"}
	if ents, err := os.ReadDir(filepath.Join(dir, "fact")); err == nil {
		var fact []string
		for _, e := range ents {
			if name := e.Name(); name[0] >= '0' && name[0] <= '9' {
				fact = append(fact, "fact/"+name)
			}
		}
		slices.Sort(fact)
		files = append(files, fact...)
	}

	for _, f := range files {
		allocFile(t, w, dir+"/"+f)
	}
	for _, f := range files {
		path := dir + "/" + f
		if !w.readfile(path) {
			continue
		}
		w.advance()
		loadAll(w, stdout)
		w.closefile(path)
	}

	bw := bufio.NewWriter(stdout)
	for n := range maxBoxes {
		if w.bx[n] != nil {
			w.saveBox(bw, n)
		}
	}
	bw.Flush()
}

// allocFile is the C driver's alloc_file(): allocBox for every "n kind
// subkind" line of a box file, the first line and each after a blank one.
func allocFile(t *testing.T, w *world, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	header := true
	for line := range strings.SplitSeq(string(data), "\n") {
		if line == "" {
			header = true
			continue
		}
		if !header || line[0] == '#' {
			continue
		}
		header = false
		n := int(int32(clib.Atoi(line)))
		_, rest, ok := strings.Cut(line, " ")
		if n <= 0 || !ok {
			continue
		}
		ks, sks, ok := strings.Cut(rest, " ")
		if !ok {
			continue
		}
		k := slices.Index(kindS, ks)
		sk := 0
		if sks != "0" {
			sk = slices.Index(subkindS, sks)
		}
		if k < 0 || sk < 0 {
			continue
		}
		w.allocBox(n, k, sk)
	}
}
