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

type entry struct {
	mode fs.FileMode
	data []byte
}

func regularBytes(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("unsupported nonregular entry: %s (%s)", path, info.Mode().Type())
	}
	return os.ReadFile(path)
}

func capture(dir string) (map[string]entry, error) {
	files := make(map[string]entry)
	for _, name := range []string{"stdout", "stderr", "status"} {
		b, err := regularBytes(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		// Stream/status modes are capture plumbing, not engine output.
		files[name] = entry{data: b}
	}
	err := filepath.WalkDir(filepath.Join(dir, "cwd"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(filepath.Join(dir, "cwd"), path)
		if err != nil {
			return err
		}
		key := "cwd"
		if rel != "." {
			key += "/" + filepath.ToSlash(rel)
		}
		mode := info.Mode()
		if key == "cwd" && !mode.IsDir() {
			return fmt.Errorf("capture cwd is not a directory: %s", path)
		}
		var b []byte
		if mode.IsRegular() {
			b, err = regularBytes(path)
			if err != nil {
				return err
			}
		} else if !mode.IsDir() {
			return fmt.Errorf("unsupported nonregular entry: %s (%s)", path, mode.Type())
		}
		files[key] = entry{mode: mode, data: b}
		return nil
	})
	return files, err
}

func differences(a, b map[string]entry) []string {
	var diff []string
	for path, left := range a {
		right, ok := b[path]
		if !ok {
			diff = append(diff, "only reference: "+path)
		} else if left.mode.Type() != right.mode.Type() {
			diff = append(diff, "type differs: "+path)
		} else {
			if left.mode != right.mode {
				diff = append(diff, "mode differs: "+path)
			}
			if !bytes.Equal(left.data, right.data) {
				diff = append(diff, "bytes differ: "+path)
			}
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

// Like tree-modes.py, capture includes the cwd root and permission/special
// bits on ALL entries, not only regular files. There are no path exclusions.
func manifest(files map[string]entry) string {
	var paths []string
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	var out strings.Builder
	for _, path := range paths {
		e := files[path]
		if e.mode.IsDir() {
			fmt.Fprintf(&out, "%s  -  %s\n", e.mode, path)
		} else {
			fmt.Fprintf(&out, "%s  %x  %s\n", e.mode, sha256.Sum256(e.data), path)
		}
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
	fmt.Printf("PASS: status, raw stdout/stderr and %d cwd entries (types, modes, bytes) identical\n", len(a)-3)
}
