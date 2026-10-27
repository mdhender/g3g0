// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package clib

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// TestFopen checks fopen(path, "w")'s behavior: it truncates, it opens a
// write-only file, and what is written reaches the file only when the
// buffer is flushed, by Close or FlushAll.
func TestFopen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f")
	if err := os.WriteFile(path, []byte("old contents\n"), 0o200); err != nil {
		t.Fatal(err)
	}

	fp, err := Fopen(path)
	if err != nil {
		t.Fatalf("Fopen of a write-only file: %v", err)
	}
	fp.WriteString("new\n")
	if data, _ := os.ReadFile(path); len(data) != 0 {
		t.Errorf("before a flush the file holds %q, want it truncated", data)
	}
	FlushAll()
	os.Chmod(path, 0o600)
	if data, _ := os.ReadFile(path); string(data) != "new\n" {
		t.Errorf("after FlushAll the file holds %q", data)
	}
	fp.WriteString("more\n")
	fp.Close()
	if data, _ := os.ReadFile(path); string(data) != "new\nmore\n" {
		t.Errorf("after Close the file holds %q", data)
	}

	if _, err := Fopen(filepath.Join(dir, "missing", "f")); Strerror(err) != "No such file or directory" {
		t.Errorf("Fopen in a missing directory: %v", err)
	}
	if _, err := Fopen(dir); Strerror(err) != "Is a directory" {
		t.Errorf("Fopen of a directory: %v", err)
	}
}

// TestReadDir checks opendir()'s errors and that "." and ".." are left out.
func TestReadDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a"), nil, 0o644)
	os.Mkdir(filepath.Join(dir, "b"), 0o755)

	names, err := ReadDir(dir)
	slices.Sort(names)
	if err != nil || !slices.Equal(names, []string{"a", "b"}) {
		t.Errorf("ReadDir = %q, %v", names, err)
	}
	if _, err := ReadDir(filepath.Join(dir, "a")); Strerror(err) != "Not a directory" {
		t.Errorf("ReadDir of a file: %v", err)
	}
	if _, err := ReadDir(filepath.Join(dir, "c")); Strerror(err) != "No such file or directory" {
		t.Errorf("ReadDir of a missing directory: %v", err)
	}
}

// TestOpendir checks opendir()'s errors, and that readdir() leaves out "."
// and ".." and returns each name once, even when entries are removed
// between calls.
func TestOpendir(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a", "b", "c", "d"} {
		os.WriteFile(filepath.Join(dir, name), nil, 0o644)
	}

	d, err := Opendir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for {
		name, ok := d.Readdir()
		if !ok {
			break
		}
		names = append(names, name)
		os.Remove(filepath.Join(dir, name))
	}
	d.Closedir()
	slices.Sort(names)
	if !slices.Equal(names, []string{"a", "b", "c", "d"}) {
		t.Errorf("Readdir = %q", names)
	}

	if _, err := Opendir(filepath.Join(dir, "a")); Strerror(err) != "No such file or directory" {
		t.Errorf("Opendir of a removed file: %v", err)
	}
	os.WriteFile(filepath.Join(dir, "f"), nil, 0o644)
	if _, err := Opendir(filepath.Join(dir, "f")); Strerror(err) != "Not a directory" {
		t.Errorf("Opendir of a file: %v", err)
	}
}
