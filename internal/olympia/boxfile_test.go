// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"syscall"
	"testing"

	"github.com/mdhender/g3g0/internal/clib"
)

// TestBoxfileMatchesC runs testdata/boxfile-cases.txt the way
// tools/boxfile-compat/run.sh runs its C driver, and requires its output,
// line for line: what each section prints, how it ends, and the tree of
// files it leaves. testdata/boxfile-c-output.txt is that output from the
// pinned G3 sources; regenerate it only with run.sh --update.
func TestBoxfileMatchesC(t *testing.T) {
	cases := readTestdata(t, "boxfile-cases.txt")
	want := readTestdata(t, "boxfile-c-output.txt")

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

	// run.sh runs the C driver with umask 022
	old := syscall.Umask(0o022)
	t.Cleanup(func() { syscall.Umask(old) })

	var out strings.Builder
	for _, s := range sections {
		dir := t.TempDir() // each section runs in a fresh directory
		t.Chdir(dir)
		out.WriteString(runBoxfileSection(t, s))
		out.WriteString(boxfileTree(dir))
		// so that the directory can be removed
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if d != nil && d.IsDir() {
				os.Chmod(path, 0o755)
			}
			return nil
		})
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

// runBoxfileSection runs one section in a fresh world, as one run of the C
// driver, and returns what the driver prints and how the section ended, as
// run.sh writes them. A nil pointer dereference is C's NULL dereference,
// which crashes with SIGSEGV (status 139).
func runBoxfileSection(t *testing.T, lines []string) (output string) {
	var out strings.Builder
	w := newWorld()
	w.stderr = &out // the driver sends stderr to stdout

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
				clib.FlushAll() // abort() flushes the open files
				out.WriteString("exit 134\n")
			default:
				panic(r)
			}
		}
		output = out.String()
	}()

	for _, line := range lines[1:] {
		fmt.Fprintf(&out, "> %s\n", line)
		c := &caseLine{rest: line}
		op := c.word()
		switch op {
		case "alloc":
			n, k, sk := c.cnum(t), c.cnum(t), c.cnum(t)
			w.allocBox(n, k, sk)
		case "delete":
			w.deleteBox(c.cnum(t))
		case "name":
			n := c.cnum(t)
			w.setName(n, c.ctext())
		case "unit":
			n, k := c.cnum(t), c.cnum(t)
			w.pPlayer(n).units.Append(k)
		case "lord":
			n, k := c.cnum(t), c.cnum(t)
			w.pChar(n).unitLord = int32(k)
		case "pretty":
			w.prettyDataFiles = c.cnum(t) != 0
		case "win":
			// win_flag: nothing these cases run reads it since G3 fixed
			// #320, and the cases check that it changes nothing.
			c.cnum(t)
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
		case "mkdir":
			os.Mkdir(c.ctext(), 0o755)
		case "chmod":
			mode := c.cnum(t)
			os.Chmod(c.ctext(), fs.FileMode(mode))
		case "symlink":
			name := c.word()
			os.Symlink(c.ctext(), name)
		case "line":
			if w.line != nil {
				fmt.Fprintf(&out, "line [%s]\n", esc(*w.line))
			} else {
				fmt.Fprintf(&out, "line NULL\n")
			}
		case "boxes":
			for n := range maxBoxes {
				if b := w.bx[n]; b != nil {
					fmt.Fprintf(&out, "box %d kind %d skind %d temp %d\n", n, b.kind, b.skind, b.temp)
				}
			}
			for k := range tMax {
				if w.boxHead[k] != 0 {
					fmt.Fprintf(&out, "kind %d:", k)
					for n := range w.loopKind(k) {
						fmt.Fprintf(&out, " %d", n)
					}
					out.WriteString("\n")
				}
			}
		case "fast":
			fmt.Fprintf(&out, "fast_scan %d\n", b2i(w.fastScan()))
		case "scan":
			w.scanBoxes(c.word())
		case "scanchars":
			w.scanChars()
		case "scanall":
			w.scanAllBoxes()
		case "read":
			w.readBoxes(c.word())
		case "readchars":
			w.readChars()
		case "readall":
			w.readAllBoxes()
		case "writekind":
			k := c.cnum(t)
			w.writeKind(k, c.word())
		case "writeplayer":
			w.writePlayer(c.cnum(t))
		case "writechars":
			w.writeChars()
		case "leftovers":
			w.writeLeftovers()
		case "writeall":
			w.writeAllBoxes()
		case "master":
			w.writeMaster()
		default:
			t.Fatalf("unknown operation %q", line)
		}
	}
	return
}

// boxfileTree is run.sh's listing of a section's directory: every path
// under dir, sorted, "d MODE PATH" for a directory ("unreadable" after it
// if it cannot be listed), and "f MODE SIZE PATH" for a file, followed by
// its lines as "| LINE", escaped, and "\ no newline at end" if the last
// line has no newline.
func boxfileTree(dir string) string {
	var out strings.Builder
	var walk func(rel string)
	walk = func(rel string) {
		path := filepath.Join(dir, rel)
		ents, err := os.ReadDir(path)
		if err != nil {
			st, _ := os.Lstat(path)
			fmt.Fprintf(&out, "d %03o %s unreadable\n", st.Mode().Perm(), rel)
			return
		}
		if rel != "" {
			st, _ := os.Lstat(path)
			fmt.Fprintf(&out, "d %03o %s\n", st.Mode().Perm(), rel)
		}
		names := make([]string, 0, len(ents))
		for _, e := range ents {
			names = append(names, e.Name())
		}
		slices.Sort(names)
		for _, name := range names {
			r := name
			if rel != "" {
				r = rel + "/" + name
			}
			p := filepath.Join(dir, r)
			st, _ := os.Lstat(p)
			switch {
			case st.IsDir():
				walk(r)
			case st.Mode().IsRegular():
				fmt.Fprintf(&out, "f %03o %d %s\n", st.Mode().Perm(), st.Size(), r)
				data, err := os.ReadFile(p)
				if err != nil {
					out.WriteString("unreadable\n")
					continue
				}
				lines := bytes.Split(data, []byte("\n"))
				last := lines[len(lines)-1]
				for _, l := range lines[:len(lines)-1] {
					fmt.Fprintf(&out, "| %s\n", esc(string(l)))
				}
				if len(last) > 0 {
					fmt.Fprintf(&out, "| %s\n", esc(string(last)))
					out.WriteString("\\ no newline at end\n")
				}
			default:
				fmt.Fprintf(&out, "? %s\n", r)
			}
		}
	}
	walk("")
	return out.String()
}

// TestBoxfileWorlds is the Go side of tools/boxfile-compat/run.sh's check
// of real saved games. It runs only when BOXFILE_COMPAT_GAMES names a list
// of games, one per line: a directory that holds a copy of the game's lib,
// pretty_data_files (0 or 1) and an output prefix, separated by tabs. For
// each, from that directory with libdir "lib", as load_db() and save_db()
// do, it runs pass 1 (fastScan, or scanAllBoxes if it fails), pass 2
// (readAllBoxes), writeAllBoxes and writeMaster, and writes stdout to
// PREFIX.out and stderr to PREFIX.err. run.sh compares them, and the lib
// trees, with the C driver's.
func TestBoxfileWorlds(t *testing.T) {
	list := os.Getenv("BOXFILE_COMPAT_GAMES")
	if list == "" {
		t.Skip("BOXFILE_COMPAT_GAMES not set; run tools/boxfile-compat/run.sh")
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
		t.Chdir(f[0])
		var stdout, stderr bytes.Buffer // nothing writes to stdout
		w := newWorld()
		w.stderr = &stderr
		w.prettyDataFiles = f[1] == "1"

		if !w.fastScan() { // pass 1: call alloc_box for each entity
			w.scanAllBoxes()
		}
		w.readAllBoxes() // pass 2: read the entity attributes
		w.writeAllBoxes()
		w.writeMaster()

		if err := os.WriteFile(f[2]+".out", stdout.Bytes(), 0o666); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(f[2]+".err", stderr.Bytes(), 0o666); err != nil {
			t.Fatal(err)
		}
	}
}
