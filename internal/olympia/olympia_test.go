// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"os"
	"strings"
	"testing"
	"unsafe"
)

// The expected output in testdata comes from the pinned C olympia-g3:
// rng-selftest.txt is G3's tests/olympia/regress/rng-selftest/EXPECT.txt,
// usage-x.txt is the C engine's stderr for "olympia-g3 -x", and
// ilist-test.txt is its stdout for "olympia-g3 -t". The full
// comparison against the C build is tools/olympia-cli-compat/run.sh.

// header is what main() prints before it dispatches, on a 64-bit host.
const header = "The Olympia C code is not 64-bit clean.\n" +
	"at the least: it puts pointers into ilists\n" +
	"\tsizeof(struct box) = 144\n"

func run(tty bool, args ...string) (int, string, string) {
	var stdout, stderr strings.Builder
	status := Run(append([]string{"olympia-g3"}, args...), &stdout, &stderr, tty)
	return status, stdout.String(), stderr.String()
}

func readTestdata(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// requireFull skips t, a test that replays whole saved games, unless the
// opt-in G3G0_FULL=1 is set; -short always skips it. tools/README.md says
// when to run them.
func requireFull(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("a whole-game test; skipped under -short")
	}
	if os.Getenv("G3G0_FULL") != "1" {
		t.Skip("a whole-game test; set G3G0_FULL=1 to run it (see tools/README.md)")
	}
}

func TestSizeofBox(t *testing.T) {
	// The LP64 reference build prints 144.
	if unsafe.Sizeof(uintptr(0)) == 8 && sizeofBox != 144 {
		t.Errorf("sizeof(struct box) = %d, want 144", sizeofBox)
	}
}

func TestSelfTest(t *testing.T) {
	want := readTestdata(t, "rng-selftest.txt")
	status, stdout, stderr := run(false, "-R")
	if status != 0 || stdout != want || stderr != "" {
		t.Errorf("-R: status %d, stderr %q, stdout matches C: %v", status, stderr, stdout == want)
	}

	// On a terminal the first 100 draws are a 10x10 grid, row by row.
	status, tty, _ := run(true, "-R")
	if status != 0 {
		t.Fatalf("-R on a tty: status %d", status)
	}
	lines := strings.Split(tty, "\n")
	var flat []string
	for i, line := range lines {
		if i >= 3 && i < 13 {
			if len(line) != 48 {
				t.Errorf("grid row %q is not ten %%3d columns", line)
			}
			flat = append(flat, strings.Fields(line)...)
			continue
		}
		flat = append(flat, line)
	}
	if got := strings.Join(flat, "\n"); got != want {
		t.Errorf("-R on a tty: the grid is not the draws in draw order")
	}
}

func TestUsage(t *testing.T) {
	want := readTestdata(t, "usage-x.txt")
	for _, args := range [][]string{{"-x"}, {"-x", "-r"}, {"-r", "-x"}} {
		status, stdout, stderr := run(false, args...)
		if status != 1 || stdout != header || stderr != want {
			t.Errorf("%q: status %d, stdout %q, stderr %q", args, status, stdout, stderr)
		}
	}
	// No getopt() diagnostic: '?' is in the option string, and the long
	// options are checked before getopt() runs.
	usage := want[strings.IndexByte(want, '\n')+1:]
	for _, args := range [][]string{{"-?"}, {"--coverage"}, {"--state-export", ""}, {"--coverage", "", "-R"}} {
		status, stdout, stderr := run(false, args...)
		if status != 1 || stdout != header || stderr != usage {
			t.Errorf("%q: status %d, stdout %q, stderr %q", args, status, stdout, stderr)
		}
	}
	// -R and -t return from inside the option loop, before errflag is
	// looked at.
	if status, _, _ := run(false, "-x", "-R"); status != 0 {
		t.Errorf("-x -R: status %d, want 0", status)
	}
	if status, _, _ := run(false, "-x", "-t"); status != 0 {
		t.Errorf("-x -t: status %d, want 0", status)
	}
}

func TestIListTest(t *testing.T) {
	want := readTestdata(t, "ilist-test.txt")
	for _, args := range [][]string{{"-t"}, {"-t", "-R"}, {"-t", "-x"}} {
		status, stdout, stderr := run(false, args...)
		if status != 0 || stdout != want || stderr != "" {
			t.Errorf("%q: status %d, stdout %q, stderr %q", args, status, stdout, stderr)
		}
	}
}

func TestNotYetImplemented(t *testing.T) {
	for _, args := range [][]string{nil, {"-r"}, {"-S"}, {"--coverage", "-R"}, {"foo", "-R"}} {
		status, stdout, stderr := run(false, args...)
		if status == 0 || stdout != header || !strings.HasSuffix(stderr, ": not yet implemented\n") {
			t.Errorf("%q: status %d, stdout %q, stderr %q", args, status, stdout, stderr)
		}
	}
}
