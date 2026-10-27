// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// turn-compat compares captured engine runs, not engine implementations.
package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// These are output artifacts in pinned main.c, sout.c, c2.c and
// io.c:save_logdir. Everything else is compared, including unknown paths.
// In particular, report routines' changes to fact/* or any other DB state
// are NOT excluded. Do not replace this with a database-file allowlist.
func reportArtifact(path string) bool {
	switch path {
	case "lib/players", "lib/email", "lib/forward", "lib/factions", "lib/totimes",
		"lib/times_0", "lib/times_press", "lib/times_rumor":
		return true
	}
	return strings.HasPrefix(path, "lib/log/") || strings.HasPrefix(path, "lib/save/")
}

func capture(dir string) (map[string][]byte, error) {
	files := make(map[string][]byte)
	for _, name := range []string{"stdout", "stderr", "status"} {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		files[name] = b
	}
	err := filepath.WalkDir(filepath.Join(dir, "cwd"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(filepath.Join(dir, "cwd"), path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if reportArtifact(rel) {
			return nil
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("not a regular database file: %s", path)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files["cwd/"+rel] = b
		return nil
	})
	return files, err
}

func differences(a, b map[string][]byte) []string {
	var diff []string
	for path, left := range a {
		right, ok := b[path]
		if !ok {
			diff = append(diff, "only reference: "+path)
		} else if !bytes.Equal(left, right) {
			diff = append(diff, "bytes differ: "+path)
		}
	}
	for path := range b {
		if _, ok := a[path]; !ok {
			diff = append(diff, "only candidate: "+path)
		}
	}
	sort.Strings(diff)
	return diff
}

func manifest(files map[string][]byte) string {
	var paths []string
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	var out strings.Builder
	for _, path := range paths {
		fmt.Fprintf(&out, "%x  %s\n", sha256.Sum256(files[path]), path)
	}
	return out.String()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: turn-compat REFERENCE CANDIDATE | --manifest RUN")
		os.Exit(2)
	}
	if os.Args[1] == "--manifest" {
		files, err := capture(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		fmt.Print(manifest(files))
		return
	}
	a, err := capture(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	b, err := capture(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	diff := differences(a, b)
	for _, line := range diff {
		fmt.Println(line)
	}
	if len(diff) != 0 {
		os.Exit(1)
	}
	fmt.Printf("PASS: status, raw stdout/stderr and %d state files identical\n", len(a)-3)
}
