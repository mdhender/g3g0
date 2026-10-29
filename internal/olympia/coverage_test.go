// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestCommandCoverageLaterModesRemainGated(t *testing.T) {
	t.Chdir(t.TempDir())
	for _, flag := range []string{"--coverage", "--use-coverage", "--phase-coverage"} {
		for _, later := range [][]string{
			{"--state-export", "s"},
			{"-e"}, {"-E"}, {"-s"}, {"-a"}, {"-i"}, {"-M"}, {"-h"},
		} {
			args := append([]string{flag, "c"}, later...)
			status, stdout, stderr := run(false, args...)
			if status != 2 || stdout != header || !strings.Contains(stderr, "load_db (") || !strings.HasSuffix(stderr, ": not yet implemented\n") {
				t.Fatalf("%q: status %d stdout %q stderr %q", args, status, stdout, stderr)
			}
		}
	}
	entries, err := os.ReadDir(".")
	if err != nil || len(entries) != 0 {
		t.Fatalf("gated modes wrote files: %v, %v", entries, err)
	}
}

// This is a Go port boundary, not an invented C outcome: a real but unported
// C handler must still stop, without counting NULL-start, failure or success.
func TestCommandCoverageNYIBoundary(t *testing.T) {
	w := newWorld()
	w.stdout, w.stderr = io.Discard, io.Discard
	w.libdir = t.TempDir()
	w.openLogfileNondestruct()
	t.Cleanup(w.closeLogfile)
	w.allocBox(200, tPlayer, subPlRegular)
	w.allocBox(1000, tChar, 0)
	w.pChar(1000).unitLord = 200
	c := w.pCommand(1000)
	c.who = 1000
	line := "breed"
	if !w.olyParse(c, &line) {
		t.Fatal("breed did not parse")
	}
	path := ""
	w.coveragePath = &path
	defer func() {
		if r := recover(); r != notImplemented("command breed start") {
			t.Fatalf("panic = %v, want breed NYI", r)
		}
		if w.covCounts != nil {
			t.Fatal("NYI fabricated a C start outcome")
		}
	}()
	w.doCommand(c)
}
