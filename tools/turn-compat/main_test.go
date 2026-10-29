// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestRecordedCProvenance(t *testing.T) {
	metadata, err := os.ReadFile("testdata/c-determinism.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(metadata), "pin f73e48118805ffe8ff5e763bfa4f20ae00c3a705\n") {
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
	reference := t.TempDir()
	if err := os.CopyFS(reference, os.DirFS("testdata/world-sample")); err != nil {
		t.Fatal(err)
	}
	// Synthetic paths are comparator controls, not engine expectations.
	for _, dir := range []string{"empty", "lib/fact", "lib/orders", "lib/log", "lib/save/1"} {
		if err := os.MkdirAll(filepath.Join(reference, "cwd", dir), 0755); err != nil {
			t.Fatal(err)
		}
	}
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
		{"outside lib whitespace name", "cwd/unknown \t\n", "append", "only candidate: cwd/unknown \t\n"},
		{"empty directory", "cwd/empty", "remove", "only reference: cwd/empty"},
		{"file mode only", "cwd/lib/system", "chmod", "mode differs: cwd/lib/system"},
		{"directory mode", "cwd/empty", "chmod", "mode differs: cwd/empty"},
		{"cwd mode", "cwd", "chmod", "mode differs: cwd"},
		{"special mode", "cwd/lib/system", "setuid", "mode differs: cwd/lib/system"},
		{"file to directory", "cwd/lib/system", "directory", "type differs: cwd/lib/system"},
		{"directory to file", "cwd/empty", "file", "type differs: cwd/empty"},
		{"stdout whitespace", "stdout", "append", "bytes differ: stdout"},
		{"stderr whitespace", "stderr", "append", "bytes differ: stderr"},
		{"exit status", "status", "append", "bytes differ: status"},
	}
	for _, path := range []string{"log/1", "save/1/1", "players", "email", "forward", "factions", "totimes", "times_0", "times_press", "times_rumor"} {
		want := "only candidate: cwd/lib/" + path
		if path == "times_0" {
			want = "bytes differ: cwd/lib/" + path
		}
		cases = append(cases, struct{ name, path, action, want string }{"report " + path, "cwd/lib/" + path, "append", want})
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.CopyFS(dir, os.DirFS(reference)); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, tc.path)
			switch tc.action {
			case "chmod", "setuid":
				mode := os.FileMode(0700)
				if tc.action == "setuid" {
					mode = 0644 | os.ModeSetuid
				}
				if err := os.Chmod(path, mode); err != nil {
					t.Fatal(err)
				}
			case "directory", "file":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if tc.action == "directory" {
					err = os.Mkdir(path, 0755)
				} else {
					err = os.WriteFile(path, nil, 0644)
				}
				if err != nil {
					t.Fatal(err)
				}
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
	for _, name := range []string{"stdout", "stderr", "status", "cwd"} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.CopyFS(dir, os.DirFS("testdata/world-sample")); err != nil {
				t.Fatal(err)
			}
			if err := os.RemoveAll(filepath.Join(dir, name)); err != nil {
				t.Fatal(err)
			}
			if _, err := capture(dir); err == nil {
				t.Fatalf("missing %s accepted", name)
			}
		})
	}
	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS("testdata/rejected-option")); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(dir, "cwd")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "cwd"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := capture(dir); err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Fatalf("non-directory cwd: %v", err)
	}
}

func TestUnsupportedEntries(t *testing.T) {
	for _, kind := range []string{"symlink", "fifo", "socket", "stream symlink"} {
		t.Run(kind, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.CopyFS(dir, os.DirFS("testdata/world-sample")); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, "cwd/lib/times_press")
			var err error
			switch kind {
			case "symlink":
				err = os.Symlink("system", path)
			case "fifo":
				err = exec.Command("mkfifo", path).Run()
			case "socket":
				// Keep the Unix-domain socket path below sockaddr_un's limit.
				path = filepath.Join(t.TempDir(), "s")
				listener, listenErr := net.Listen("unix", path)
				if listenErr != nil {
					t.Fatal(listenErr)
				}
				defer listener.Close()
				err = os.Rename(path, filepath.Join(dir, "cwd/socket"))
			case "stream symlink":
				path = filepath.Join(dir, "stdout")
				if err = os.Remove(path); err == nil {
					err = os.Symlink("stderr", path)
				}
			}
			if err != nil {
				t.Fatal(err)
			}
			if _, err := capture(dir); err == nil || !strings.Contains(err.Error(), "unsupported nonregular entry:") {
				t.Fatalf("%s: expected clear nonregular-entry rejection, got %v", kind, err)
			}
		})
	}
}

func shellFunction(t *testing.T, name string) string {
	t.Helper()
	script, err := os.ReadFile("run.sh")
	if err != nil {
		t.Fatal(err)
	}
	header := name + "() {"
	_, function, ok := strings.Cut(string(script), header)
	if !ok {
		t.Fatalf("missing %s", name)
	}
	function, _, ok = strings.Cut(function, "\n}\n")
	if !ok {
		t.Fatalf("missing %s end", name)
	}
	return header + function + "\n}\n"
}

func TestShellComparisonStatuses(t *testing.T) {
	// Exercise the actual shell function with controlled command statuses,
	// including a failed mode listing after an otherwise equal comparison.
	function := shellFunction(t, "compare_captures")
	for _, tc := range []struct{ compare, modes, want int }{
		{0, 0, 0}, {1, 0, 1}, {0, 1, 1}, {1, 1, 1},
		{2, 0, 2}, {127, 0, 2}, {0, 2, 2}, {1, 2, 2}, {0, 127, 2},
	} {
		cmd := exec.Command("bash", "-c", "set -euo pipefail\n"+
			function+
			"fake_compare() { return "+strconv.Itoa(tc.compare)+"; }\n"+
			"python3() { return "+strconv.Itoa(tc.modes)+"; }\n"+
			"COMPARE=fake_compare HERE=.\ncompare_captures left right")
		output, err := cmd.CombinedOutput()
		got := cmd.ProcessState.ExitCode()
		if got != tc.want {
			t.Fatalf("compare %d modes %d: exit %d, want %d (%v): %s", tc.compare, tc.modes, got, tc.want, err, output)
		}
	}
}

func TestCReferenceMarkers(t *testing.T) {
	function := shellFunction(t, "run_c_engine")
	for _, stream := range []string{"stdout", "stderr"} {
		for _, marker := range []string{"ordinary rejected-option diagnostic", "out-of-scope C dependency: helper", "not yet implemented", "missing C dependency", "C stub", "AddressSanitizer", "UndefinedBehaviorSanitizer", "runtime error: signed integer overflow", "missing stream"} {
			t.Run(stream+"/"+marker, func(t *testing.T) {
				dir := t.TempDir()
				for _, name := range []string{"stdout", "stderr"} {
					if err := os.WriteFile(filepath.Join(dir, name), nil, 0644); err != nil {
						t.Fatal(err)
					}
				}
				path := filepath.Join(dir, stream)
				var err error
				if marker == "missing stream" {
					err = os.Remove(path)
				} else {
					err = os.WriteFile(path, []byte(marker), 0644)
				}
				if err != nil {
					t.Fatal(err)
				}
				cmd := exec.Command("bash", "-c", "set -euo pipefail\n"+function+
					"run_engine() { :; }\ndie() { echo \"$*\" >&2; exit 2; }\nrun_c_engine ignored ignored \"$1\"", "bash", dir)
				output, err := cmd.CombinedOutput()
				want := 2
				if strings.HasPrefix(marker, "ordinary") {
					want = 0
				}
				if got := cmd.ProcessState.ExitCode(); got != want {
					t.Fatalf("exit %d, want %d (%v): %s", got, want, err, output)
				}
			})
		}
	}
}
