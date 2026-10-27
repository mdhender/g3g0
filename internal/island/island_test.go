// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package island

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The expected values below were observed from the pinned C island-g3 on
// macOS. The full byte-for-byte comparison against the C build is
// tools/island-compat/run.sh.

// ocean40 is tests/island/fixtures/ocean-40.map.
var ocean40 = strings.Repeat(strings.Repeat(".", 40)+"\n", 40)

// run runs island in a fresh directory holding seed as randseed (none when
// seed is nil) and returns the directory, exit status, stdout and stderr.
func run(t *testing.T, seed []byte, stdin string, args ...string) (string, int, string, string) {
	t.Helper()
	dir := t.TempDir()
	if seed != nil {
		if err := os.WriteFile(filepath.Join(dir, "randseed"), seed, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var stdout, stderr bytes.Buffer
	status := Run(dir, append([]string{"island-g3"}, args...), strings.NewReader(stdin), &stdout, &stderr)
	return dir, status, stdout.String(), stderr.String()
}

func sha(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func TestMatchesC(t *testing.T) {
	for _, tc := range []struct {
		name     string
		seed     string
		args     []string
		stdout   string // sha256
		randseed string // hex
		stderr   string
	}{
		{"default-a", "ISLAND000000000A", nil,
			"61f7719fe19cb40293751d22862b1fe6051081b96a0c38cc57e2735c6b7f7b84",
			"c9c8d073e7356b6e72c36cffb2a0d1bb", "Added island of 14 provinces.\n"},
		{"size-50", "ISLAND000000000A", []string{"-s", "50"},
			"2c6ab337aa5f84f655a9d60b4c7f9c7a785e37c5e7618bda928ca73bbdc31a4a",
			"cd1eae7c3f27233a4114697aeec27152", "Added island of 50 provinces.\n"},
		{"default-zero", string(make([]byte, 16)), nil,
			"1b861919c52d4a93be3b78a7903747c555772b14cfba64cfea4f56a68c30c9d5",
			"1593cdfd8006e0d991b038966142bc44", "Added island of 322 provinces.\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir, status, stdout, stderr := run(t, []byte(tc.seed), ocean40, tc.args...)
			if status != 0 || stderr != tc.stderr {
				t.Fatalf("status %d, stderr %q; want 0, %q", status, stderr, tc.stderr)
			}
			if got := sha(stdout); got != tc.stdout {
				t.Errorf("stdout sha256 %s, want %s", got, tc.stdout)
			}
			seed, err := os.ReadFile(filepath.Join(dir, "randseed"))
			if err != nil {
				t.Fatal(err)
			}
			if got := hex.EncodeToString(seed); got != tc.randseed {
				t.Errorf("randseed %s, want %s", got, tc.randseed)
			}
		})
	}
}

func TestMissingSeedIsZero(t *testing.T) {
	// With no randseed, island runs from the zero seed and says so on
	// stderr; stdout is the map alone (G3 #298).
	_, status, stdout, stderr := run(t, nil, ocean40)
	if status != 0 || stderr != "randseed could not be opened.\nAdded island of 322 provinces.\n" {
		t.Fatalf("status %d, stderr %q", status, stderr)
	}
	if got := sha(stdout); got != "1b861919c52d4a93be3b78a7903747c555772b14cfba64cfea4f56a68c30c9d5" {
		t.Errorf("stdout sha256 %s", got)
	}
}

func TestErrors(t *testing.T) {
	const usage = "Usage: island-g3 -b border_width -c continental_shelf_width -s size\n"
	seed := []byte("ISLAND000000000A")
	for _, tc := range []struct {
		name   string
		stdin  string
		args   []string
		stderr string
	}{
		{"unknown option", ocean40, []string{"-x"}, usage},
		{"missing argument", ocean40, []string{"-s"}, usage},
		{"long option", ocean40, []string{"--s", "5"}, usage},
		{"empty", "", nil, "No map on stdin!\n"},
		{"blank lines", "\n\r\n", nil, "No map on stdin!\n"},
		{"wide", strings.Repeat(".", 120) + "\n", nil, "Map row longer than 119 characters!\n"},
		{"tall", strings.Repeat(".\n", 121), nil, "Map taller than 120 rows!\n"},
		{"allocation", ocean40, []string{"-s", "2147483647"}, "Unable to allocate memory for island data!\n"},
		{"all land", strings.Repeat("pppp\n", 4), nil, "No room for an island!\n"},
		{"border", ocean40, []string{"-b", "20"}, "No room for an island!\n"},
		{"wide border", ocean40, []string{"-b", "1000"}, "No room for an island!\n"},
		{"nul lines", "\x00\n\x00x\n", nil, "No map on stdin!\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir, status, stdout, stderr := run(t, seed, tc.stdin, tc.args...)
			if status != 1 || stdout != "" || stderr != tc.stderr {
				t.Errorf("status %d, stdout %q, stderr %q; want 1, \"\", %q", status, stdout, stderr, tc.stderr)
			}
			got, err := os.ReadFile(filepath.Join(dir, "randseed"))
			if err != nil || !bytes.Equal(got, seed) {
				t.Errorf("randseed changed: %x, %v", got, err)
			}
		})
	}
}

func TestSeedNotWritable(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "randseed"), 0o755); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	status := Run(dir, []string{"island-g3"}, strings.NewReader(ocean40), &stdout, &stderr)
	if status != 1 || stdout.Len() != 0 || stderr.String() != "couldn't write randseed:Is a directory\n" {
		t.Errorf("status %d, stdout %q, stderr %q", status, stdout.String(), stderr.String())
	}
}

func TestOperandsStopOptions(t *testing.T) {
	// getopt(3) on macOS does not permute, so -s after an operand is ignored.
	seed := []byte("ISLAND000000000A")
	_, _, want, _ := run(t, seed, ocean40)
	_, status, got, stderr := run(t, seed, ocean40, "foo", "-s", "5")
	if status != 0 || got != want || stderr != "Added island of 14 provinces.\n" {
		t.Errorf("status %d, stderr %q, same map %v", status, stderr, got == want)
	}
}

func TestNULLinesAreSkipped(t *testing.T) {
	// A line starting with NUL is not a row (G3 #303).
	seed := []byte("ISLAND000000000A")
	_, _, want, _ := run(t, seed, ocean40)
	half := strings.Repeat(strings.Repeat(".", 40)+"\n", 20)
	for _, stdin := range []string{"\x00x\n" + ocean40, half + "\x00x\n" + half, ocean40 + "\x00"} {
		_, status, got, _ := run(t, seed, stdin)
		if status != 0 || got != want {
			t.Errorf("status %d, same map %v for %q...", status, got == want, stdin[:4])
		}
	}
}
