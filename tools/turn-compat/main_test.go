// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRecordedCProvenance(t *testing.T) {
	metadata, err := os.ReadFile("testdata/c-determinism.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(metadata), "pin 2567e6a5dd0ba0d95516d8bed9e1bafb87a51985\n") {
		t.Fatal("recording is not from the stage pin")
	}
	for _, name := range []string{"world-sample", "rejected-option"} {
		files, err := capture(filepath.Join("testdata", name))
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(strings.TrimSpace(manifest(files)), "\n") {
			if !strings.Contains(string(metadata), line+"\n") {
				t.Errorf("%s: recorded bytes absent from C determinism manifest: %s", name, line)
			}
		}
	}
}

// Mutations below are comparator negative controls, NOT engine expectations.
// The unmodified bytes all come from two actual pinned C turns.
func TestComparisonContract(t *testing.T) {
	reference := "testdata/world-sample"
	a, err := capture(reference)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, path, action, want string
	}{
		{"identical", "", "", ""},
		{"database mutation", "cwd/lib/system", "append", "bytes differ: cwd/lib/system"},
		{"missing database", "cwd/lib/randseed", "remove", "only reference: cwd/lib/randseed"},
		{"missing empty database", "cwd/lib/ship", "remove", "only reference: cwd/lib/ship"},
		{"extra unknown database", "cwd/lib/new-db", "append", "only candidate: cwd/lib/new-db"},
		{"extra fact", "cwd/lib/fact/123", "append", "only candidate: cwd/lib/fact/123"},
		{"extra orders", "cwd/lib/orders/123", "append", "only candidate: cwd/lib/orders/123"},
		{"trade routes", "cwd/lib/trade_routes", "append", "only candidate: cwd/lib/trade_routes"},
		{"report path state", "cwd/lib/fact/players", "append", "only candidate: cwd/lib/fact/players"},
		{"unknown times", "cwd/lib/times_state", "append", "only candidate: cwd/lib/times_state"},
		{"outside lib", "cwd/state", "append", "only candidate: cwd/state"},
		{"stdout whitespace", "stdout", "append", "bytes differ: stdout"},
		{"stderr whitespace", "stderr", "append", "bytes differ: stderr"},
		{"exit status", "status", "append", "bytes differ: status"},
	}
	for _, path := range []string{"log/1", "save/1/1", "players", "email", "forward", "factions", "totimes", "times_0", "times_press", "times_rumor"} {
		cases = append(cases, struct{ name, path, action, want string }{"deferred " + path, "cwd/lib/" + path, "append", ""})
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.CopyFS(dir, os.DirFS(reference)); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, tc.path)
			switch tc.action {
			case "remove":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			case "append":
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatal(err)
				}
				f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					t.Fatal(err)
				}
				_, err = f.WriteString(" \n")
				closeErr := f.Close()
				if err != nil || closeErr != nil {
					t.Fatalf("write: %v; close: %v", err, closeErr)
				}
			}
			b, err := capture(dir)
			if err != nil {
				t.Fatal(err)
			}
			if got := strings.Join(differences(a, b), "\n"); got != tc.want {
				t.Fatalf("diff = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestIncompleteCaptureFails(t *testing.T) {
	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS("testdata/world-sample")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dir, "stderr")); err != nil {
		t.Fatal(err)
	}
	if _, err := capture(dir); err == nil {
		t.Fatal("missing stderr accepted as empty stderr")
	}
}
