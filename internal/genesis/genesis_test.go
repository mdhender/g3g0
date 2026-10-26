// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package genesis

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The expected values below were observed from the pinned C genesis-g3 on
// macOS. The full byte-for-byte comparison against the C build is
// tools/genesis-compat/run.sh.

func TestAtoi(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want int
	}{
		{"10", 10},
		{" 10", 10},
		{"+7", 7},
		{"-3", -3},
		{"20x", 20},
		{"abc", 0},
		{"", 0},
		{"4294967306", 10},           // (int)4294967306
		{"-4294967286", 10},          // (int)-4294967286
		{"99999999999999999999", -1}, // strtol saturates to LONG_MAX
		{"-99999999999999999999", 0}, // LONG_MIN truncates to 0
	} {
		if got := atoi(tc.in); got != tc.want {
			t.Errorf("atoi(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestStrtoull0(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want uint64
	}{
		{"0", 0},
		{"12345", 12345},
		{"0x10", 16},
		{"0X10", 16},
		{"010", 8},
		{"08", 0},
		{"0x", 0},
		{"0xg", 0},
		{"  +12", 12},
		{"42abc", 42},
		{"nope", 0},
		{"-1", math.MaxUint64},
		{"18446744073709551615", math.MaxUint64},
		{"99999999999999999999999", math.MaxUint64},
		{"-99999999999999999999999", math.MaxUint64},
	} {
		if got := strtoull0(tc.in); got != tc.want {
			t.Errorf("strtoull0(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

// run runs genesis in a fresh directory and returns it with the exit status
// and stderr.
func run(t *testing.T, args ...string) (string, int, string) {
	t.Helper()
	dir := t.TempDir()
	var stderr bytes.Buffer
	status := Run(dir, append([]string{"genesis-g3"}, args...), &stderr)
	return dir, status, stderr.String()
}

func TestArgumentErrors(t *testing.T) {
	for _, tc := range []struct {
		args  []string
		first string // the first line of stderr
	}{
		{[]string{"-x"}, "genesis-g3: invalid option -- x"},
		{[]string{"-:"}, "genesis-g3: invalid option -- :"},
		{[]string{"-s"}, "genesis-g3: option requires an argument -- s"},
		{[]string{"--bogus"}, "genesis-g3: unrecognized option `--bogus'"},
		{[]string{"--bogus=3"}, "genesis-g3: unrecognized option `--bogus=3'"},
		{[]string{"--s", "5"}, "genesis-g3: option `--s' is ambiguous"},
		{[]string{"--i", "3"}, "genesis-g3: option `--i' is ambiguous"},
		{[]string{"--=x"}, "genesis-g3: option `--' is ambiguous"},
		{[]string{"---size", "5"}, "genesis-g3: unrecognized option `---size'"},
		{[]string{"--seed"}, "genesis-g3: option `--seed' requires an argument"},
		{[]string{"--defer-oceans=1"}, "genesis-g3: option `--defer-oceans' doesn't allow an argument"},
		{[]string{"-D-x"}, "genesis-g3: unrecognized option `--x'"},
		{[]string{"-D-"}, "Usage: genesis-g3 [--seed N] [--size 10..99] [--number-of-oceans N]"},
		{[]string{"x", "-x"}, "genesis-g3: invalid option -- x"},
		{[]string{"--size=5"}, "genesis: --size must be 10..99 (got 5)"},
		{[]string{"--size", "99999999999999999999"}, "genesis: --size must be 10..99 (got -1)"},
		{[]string{"--size", "9", "-n", "0"}, "genesis: --size must be 10..99 (got 9)"},
		{[]string{"-n", "0"}, "genesis: --number-of-oceans must be >= 1 (got 0)"},
		{[]string{"-i", "11"}, "genesis: --island-iterations must be 2..10 (got 11)"},
		{[]string{"--ocean-growth=BFS"}, "genesis: --ocean-growth must be frontier or bfs (got BFS)"},
		{[]string{"--ocean-growth=bfs"}, "genesis: --ocean-growth bfs requires --defer-oceans"},
		{[]string{"--start-cities", "17"}, "genesis: --start-cities must be 0..16 (got 17)"},
		{[]string{"--image", "x.bmp", "--pixels", "0"}, "genesis: --pixels must be 1..64 (got 0)"},
		{[]string{"--image", "x.bmp", "--grid", "5"}, "genesis: --grid must be 0..4 (got 5)"},
	} {
		dir, status, stderr := run(t, tc.args...)
		first, _, _ := strings.Cut(stderr, "\n")
		if status != 1 || first != tc.first {
			t.Errorf("%q: status %d, stderr %q; want status 1, %q", tc.args, status, first, tc.first)
		}
		if entries, _ := os.ReadDir(dir); len(entries) != 0 {
			t.Errorf("%q: wrote %d files, want none", tc.args, len(entries))
		}
	}
}

func TestPosixlyCorrect(t *testing.T) {
	t.Setenv("POSIXLY_CORRECT", "1")
	if _, status, stderr := run(t, "-x"); status != 1 || !strings.HasPrefix(stderr, "genesis-g3: illegal option -- x\n") {
		t.Errorf("-x: status %d, stderr %q", status, stderr)
	}
	// Option parsing stops at the first operand, so --bogus is never seen.
	if _, status, stderr := run(t, "--size", "20", "stray", "--bogus"); status != 0 ||
		!strings.HasPrefix(stderr, "genesis: seed=0 size=20 ") {
		t.Errorf("stray operand: status %d, stderr %q", status, stderr)
	}
}

func TestHelp(t *testing.T) {
	for _, args := range [][]string{{"-h"}, {"--help"}, {"--he"}, {"-hx"}, {"--size", "5", "--help"}} {
		dir, status, stderr := run(t, args...)
		if status != 0 || !strings.HasPrefix(stderr, "Usage: genesis-g3 [--seed N]") {
			t.Errorf("%q: status %d, stderr %q", args, status, stderr)
		}
		if entries, _ := os.ReadDir(dir); len(entries) != 0 {
			t.Errorf("%q: wrote %d files, want none", args, len(entries))
		}
	}
}

// world reads the files a run wrote.
func world(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	files := map[string][]byte{}
	for _, name := range []string{"Map", "Regions", "Land", "Cities", "randseed"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		files[name] = data
	}
	return files
}

func TestDeterministicAndSeeded(t *testing.T) {
	for _, args := range [][]string{
		{"--size", "30"},
		{"--size", "30", "-D"},
		{"--size", "30", "-D", "--ocean-growth", "bfs"},
	} {
		dir1, status1, stderr1 := run(t, append([]string{"--seed", "7"}, args...)...)
		dir2, status2, stderr2 := run(t, append([]string{"--seed", "7"}, args...)...)
		dir3, status3, _ := run(t, append([]string{"--seed", "8"}, args...)...)
		if status1 != 0 || status2 != 0 || status3 != 0 {
			t.Fatalf("%q: statuses %d %d %d", args, status1, status2, status3)
		}
		if stderr1 != stderr2 {
			t.Errorf("%q: stderr differs between runs", args)
		}
		w1, w2, w3 := world(t, dir1), world(t, dir2), world(t, dir3)
		for name := range w1 {
			if !bytes.Equal(w1[name], w2[name]) {
				t.Errorf("%q: %s differs between runs with the same seed", args, name)
			}
		}
		if bytes.Equal(w1["Map"], w3["Map"]) {
			t.Errorf("%q: seeds 7 and 8 gave the same Map", args)
		}
		if len(w1["randseed"]) != 16 {
			t.Errorf("%q: randseed is %d bytes, want 16", args, len(w1["randseed"]))
		}
	}
}

func TestWriteFailure(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "Regions"), 0o755); err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	status := Run(dir, []string{"genesis-g3", "--size", "20"}, &stderr)
	if status != 1 || !strings.HasSuffix(stderr.String(), "genesis: can't write Regions: Is a directory\n") {
		t.Errorf("status %d, stderr %q", status, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "Map")); err != nil {
		t.Errorf("Map was not written before the failure: %v", err)
	}
}
