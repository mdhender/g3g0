// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"github.com/mdhender/g3g0/internal/clib"
)

// TestCheckMatchesC runs testdata/check-cases.txt the way
// tools/check-compat/run.sh runs its C driver, and requires its output,
// line for line: what each section prints, how it ends, and the tree of
// files it leaves. testdata/check-c-output.txt is that output from the
// pinned G3 sources; regenerate it only with run.sh --update.
func TestCheckMatchesC(t *testing.T) {
	cases := readTestdata(t, "check-cases.txt")
	want := readTestdata(t, "check-c-output.txt")

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
		out.WriteString(runCheckSection(t, s))
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

// checkLoad is what load_db() runs of the box-file passes, and checkSave
// what save_db() runs of the writers.
func (w *world) checkLoad() {
	if !w.fastScan() { // pass 1: call alloc_box for each entity
		w.scanAllBoxes()
	}
	w.readAllBoxes() // pass 2: read the entity attributes
}

func (w *world) checkSave() {
	w.writeAllBoxes()
	w.writeMaster()
}

// checkExit recovers what ends a run of the C driver and returns its exit
// status, writing what C writes to stderr: a failed assertion aborts
// (status 134, after flushing the open files), and a nil pointer
// dereference is C's NULL dereference, which crashes with SIGSEGV (status
// 139).
func checkExit(r any, stderr *strings.Builder) int {
	if re, ok := r.(runtime.Error); ok && strings.Contains(re.Error(), "nil pointer dereference") {
		return 139
	}
	switch r := r.(type) {
	case nil:
		return 0
	case assertion:
		stderr.WriteString(r.Error())
		clib.FlushAll() // abort() flushes the open files
		return 134
	}
	panic(r)
}

// runCheckSection runs one section in a fresh world, as one run of the C
// driver, and returns what the driver prints and how the section ended, as
// run.sh writes them. put and putl write the database and are not echoed.
func runCheckSection(t *testing.T, lines []string) (output string) {
	var out strings.Builder
	w := newWorld()
	w.stderr = &out // the driver sends stderr to stdout

	out.WriteString(lines[0] + "\n")
	defer func() {
		fmt.Fprintf(&out, "exit %d\n", checkExit(recover(), &out))
		output = out.String()
	}()

	for _, line := range lines[1:] {
		c := &caseLine{rest: line}
		op := c.word()
		if op == "put" || op == "putl" {
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
			continue
		}

		fmt.Fprintf(&out, "> %s\n", line)
		switch op {
		case "mkdir":
			os.Mkdir(c.ctext(), 0o755)
		case "load":
			w.checkLoad()
		case "check":
			w.checkDb()
		case "save":
			w.checkSave()
		default:
			t.Fatalf("unknown operation %q", line)
		}
	}
	return
}

// TestCheckWorlds is the Go side of tools/check-compat/run.sh's check of
// real saved games. It runs only with G3G0_FULL=1 (see requireFull), and
// only when CHECK_COMPAT_GAMES names a list of games, one per line: a
// directory that holds a copy of the game's lib and an output prefix,
// separated by a tab. For each, from that directory with libdir "lib", it
// runs pass 1 (fastScan, or scanAllBoxes if it fails), pass 2
// (readAllBoxes), checkDb, writeAllBoxes and writeMaster, and writes stdout
// to PREFIX.out, followed by "exit STATUS", and stderr to PREFIX.err.
// run.sh compares them, and the lib trees, with the C driver's.
func TestCheckWorlds(t *testing.T) {
	requireFull(t)
	list := os.Getenv("CHECK_COMPAT_GAMES")
	if list == "" {
		t.Skip("CHECK_COMPAT_GAMES not set; run tools/check-compat/run.sh --full")
	}
	data, err := os.ReadFile(list)
	if err != nil {
		t.Fatal(err)
	}
	for line := range strings.SplitSeq(strings.TrimSuffix(string(data), "\n"), "\n") {
		f := strings.Split(line, "\t")
		if len(f) != 2 {
			t.Fatalf("bad line in %s: %q", list, line)
		}
		t.Chdir(f[0])
		var stderr strings.Builder // nothing writes to stdout
		w := newWorld()
		w.stderr = &stderr

		status := func() (status int) {
			defer func() { status = checkExit(recover(), &stderr) }()
			w.checkLoad()
			w.checkDb()
			w.checkSave()
			return
		}()

		if err := os.WriteFile(f[1]+".out", fmt.Appendf(nil, "exit %d\n", status), 0o666); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(f[1]+".err", []byte(stderr.String()), 0o666); err != nil {
			t.Fatal(err)
		}
	}
}
