// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package mapgen

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetlin(t *testing.T) {
	r := &lineReader{data: []byte("one\n\nthree\x00x\n" + strings.Repeat("y", 1100))}
	for _, want := range []string{"one", "", "three", strings.Repeat("y", 1023)} {
		got, ok := r.getlin()
		if !ok || got != want {
			t.Fatalf("getlin = %q, %v, want %q", got, ok, want)
		}
	}
	if _, ok := r.getlin(); ok {
		t.Errorf("getlin at end of input reported a line")
	}
}

// g3Root returns the G3 checkout named by $G3_ROOT, skipping the test when it
// is not set. The checkout should be at the revision in tools/mapgen-compat.
func g3Root(t *testing.T) string {
	t.Helper()
	root := os.Getenv("G3_ROOT")
	if root == "" {
		t.Skip("set G3_ROOT to a G3 checkout to run the upstream mapgen regressions")
	}
	return root
}

// stage copies the mapgen fixture inputs into a fresh directory.
func stage(t *testing.T, root string) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range []string{"Map", "Regions", "Land", "Cities"} {
		data, err := os.ReadFile(filepath.Join(root, "tests", "mapgen", "fixtures", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// TestSecretSeaRoute mirrors tests/mapgen/regress/secret-sea-route/check.sh,
// using its committed seed and EXPECT.sha256.
func TestSecretSeaRoute(t *testing.T) {
	root := g3Root(t)
	regress := filepath.Join(root, "tests", "mapgen", "regress", "secret-sea-route")
	dir := stage(t, root)
	seed, err := os.ReadFile(filepath.Join(regress, "randseed"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "randseed"), seed, 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if status := Run(dir, &stdout, &stderr); status != 0 {
		t.Fatalf("Run exited %d, want 0\n%s", status, stderr.String())
	}

	var manifest strings.Builder
	for _, name := range []string{"gate", "loc", "road"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(&manifest, "%x  %s\n", sha256.Sum256(data), name)
	}
	want, err := os.ReadFile(filepath.Join(regress, "EXPECT.sha256"))
	if err != nil {
		t.Fatal(err)
	}
	if manifest.String() != string(want) {
		t.Errorf("manifest:\n%s\nwant:\n%s", manifest.String(), want)
	}
}

// TestSaveSeedUnwritable mirrors
// tests/mapgen/regress/save-seed-unwritable/check.sh.
func TestSaveSeedUnwritable(t *testing.T) {
	root := g3Root(t)
	dir := stage(t, root)
	if err := os.Mkdir(filepath.Join(dir, "randseed"), 0o755); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if status := Run(dir, &stdout, &stderr); status != 1 {
		t.Fatalf("Run exited %d, want 1", status)
	}
	found := false
	for _, line := range strings.Split(stderr.String(), "\n") {
		if strings.HasPrefix(line, "couldn't write randseed:") {
			found = true
		}
	}
	if !found {
		t.Errorf(`no "couldn't write randseed:" on stderr`)
	}
}
