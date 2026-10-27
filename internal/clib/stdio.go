// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package clib

import (
	"bufio"
	"bytes"
	"io"
	"math"
	"os"
	"slices"
	"strings"
	"sync"
	"syscall"
)

// This file holds the pieces of C stdio behavior: reading files and lines,
// printf conversions, and sscanf.

// ReadFile reads a file opened the way fopen(path, "r") would open it.
// A read error after a successful open (for example, a directory) ends the
// input, as fgets() returning NULL does in C.
func ReadFile(path string) ([]byte, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	data, _ := io.ReadAll(fp)
	return data, nil
}

// FgetsLines splits data into the buffers successive fgets(buf, size, fp)
// calls would return. Each buffer is cut at its first NUL byte, since the
// C code only ever treats it as a string.
func FgetsLines(data []byte, size int) []string {
	var lines []string
	for len(data) > 0 {
		n := min(size-1, len(data))
		if i := bytes.IndexByte(data[:n], '\n'); i >= 0 {
			n = i + 1
		}
		lines = append(lines, CString(data[:n]))
		data = data[n:]
	}
	return lines
}

// CString returns b up to its first NUL byte.
func CString(b []byte) string {
	b, _, _ = bytes.Cut(b, []byte{0})
	return string(b)
}

// Str formats a possibly NULL C string the way printf("%s") does on the
// supported C libraries.
func Str(s *string) string {
	if s == nil {
		return "(null)"
	}
	return *s
}

// PadRight and PadLeft implement %-Ns and %Ns. C counts bytes, while Go's
// fmt counts runes, so strings are padded here instead of by fmt.
func PadRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func PadLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}

// Ch formats a byte the way printf("%c") does. Go's string(byte) would
// encode bytes above 0x7f as UTF-8.
func Ch(c byte) string {
	return string([]byte{c})
}

// Scanner consumes a string the way sscanf() consumes its input.
type Scanner struct {
	s   string
	pos int
}

// NewScanner returns a Scanner over s.
func NewScanner(s string) *Scanner {
	return &Scanner{s: s}
}

// SkipSpace implements a whitespace directive.
func (sc *Scanner) SkipSpace() {
	for sc.pos < len(sc.s) && IsSpace(sc.s[sc.pos]) {
		sc.pos++
	}
}

// scanWidth is the most characters macOS scanf reads for one numeric
// conversion without a field width: its buffer holds 512 characters and a
// NUL.
const scanWidth = 512

// Int implements %d as the macOS (BSD) libc does. It skips white space and
// then reads an optional sign and decimal digits, at most scanWidth
// characters in all, the sign and leading zeros included; any digits past
// that stay in the input for the next directive. The value is converted
// with strtoimax(), which saturates to the 64-bit range, and stored in an
// int, which keeps its low 32 bits. Int reports false on a matching or
// input failure: no digit after the optional sign.
func (sc *Scanner) Int() (int, bool) {
	sc.SkipSpace()
	start := sc.pos
	i := start
	if i < len(sc.s) && (sc.s[i] == '+' || sc.s[i] == '-') {
		i++
	}
	digits := i
	for i < len(sc.s) && i-start < scanWidth && sc.s[i] >= '0' && sc.s[i] <= '9' {
		i++
	}
	if i == digits {
		return 0, false
	}
	sc.pos = i
	return Atoi(sc.s[start:i]), true
}

// Uint64 implements %llu as the macOS (BSD) libc does. It reads what Int
// reads: white space, then an optional sign and decimal digits, at most
// scanWidth characters in all. The value is converted with strtoumax(),
// which saturates to the unsigned 64-bit range, and negates it, as an
// unsigned value, after a minus sign. Uint64 reports false on a matching or
// input failure: no digit after the optional sign.
func (sc *Scanner) Uint64() (uint64, bool) {
	sc.SkipSpace()
	start := sc.pos
	i := start
	neg := false
	if i < len(sc.s) && (sc.s[i] == '+' || sc.s[i] == '-') {
		neg = sc.s[i] == '-'
		i++
	}
	digits := i
	var acc uint64
	overflow := false
	for i < len(sc.s) && i-start < scanWidth && sc.s[i] >= '0' && sc.s[i] <= '9' {
		d := uint64(sc.s[i] - '0')
		if overflow || acc > (math.MaxUint64-d)/10 {
			overflow = true
		} else {
			acc = acc*10 + d
		}
		i++
	}
	if i == digits {
		return 0, false
	}
	sc.pos = i
	switch {
	case overflow:
		return math.MaxUint64, true
	case neg:
		return -acc, true
	}
	return acc, true
}

// Literal matches a single ordinary character in the format.
func (sc *Scanner) Literal(c byte) bool {
	if sc.pos < len(sc.s) && sc.s[sc.pos] == c {
		sc.pos++
		return true
	}
	return false
}

// Char implements %c, which does not skip whitespace.
func (sc *Scanner) Char() (byte, bool) {
	if sc.pos < len(sc.s) {
		c := sc.s[sc.pos]
		sc.pos++
		return c, true
	}
	return 0, false
}

// File is a FILE * that fopen(path, "w") opened: writes go through a
// buffer of the size the macOS libc gives a file (its st_blksize, or BUFSIZ
// if that is not set), which Close flushes. Write errors are ignored, as
// fprintf() and fclose() callers that do not check them ignore them.
//
// abort() on macOS flushes every open stream before it kills the process,
// so a failed assert() leaves each file with all that was written to it;
// whatever recovers the port's assertion panic calls FlushAll to do the
// same. A crash (SIGSEGV) flushes nothing, and a file keeps only what its
// buffer had already written.
type File struct {
	*bufio.Writer
	f *os.File
}

// files is every File that is open, in the order Fopen opened them, for
// FlushAll.
var (
	filesMu sync.Mutex
	files   []*File
)

// FlushAll flushes every open File, as abort() does on macOS (and as
// fflush(NULL) does). Each stays open.
func FlushAll() {
	filesMu.Lock()
	defer filesMu.Unlock()

	for _, fp := range files {
		_ = fp.Flush()
	}
}

// bufsiz is BUFSIZ in the macOS <stdio.h>.
const bufsiz = 1024

// Fopen is fopen(path, "w"): it creates path (mode 0666, less the umask)
// or truncates it, and opens it for writing only.
func Fopen(path string) (*File, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return nil, err
	}

	size := bufsiz
	var st syscall.Stat_t
	if syscall.Fstat(int(f.Fd()), &st) == nil && st.Blksize > 0 {
		size = int(st.Blksize)
	}

	fp := &File{Writer: bufio.NewWriterSize(f, size), f: f}

	filesMu.Lock()
	files = append(files, fp)
	filesMu.Unlock()

	return fp, nil
}

// Close is fclose(): it flushes the buffer and closes the file.
func (fp *File) Close() {
	filesMu.Lock()
	if i := slices.Index(files, fp); i >= 0 {
		files = slices.Delete(files, i, i+1)
	}
	filesMu.Unlock()

	_ = fp.Flush()
	_ = fp.f.Close()
}

// Dir is a directory stream, as opendir() returns.
type Dir struct {
	f *os.File
}

// Opendir is opendir(): it fails on a path that is not a directory
// (ENOTDIR) or cannot be opened.
func Opendir(path string) (*Dir, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_DIRECTORY, 0)
	if err != nil {
		return nil, err
	}
	return &Dir{f: f}, nil
}

// Readdir is readdir(): the next name in the file system's order, less "."
// and "..", or false at the end or on a read error, where readdir()
// returns NULL. Names are read one at a time (on macOS, with libc's
// readdir()), so a caller that removes entries between calls sees what C
// would.
func (d *Dir) Readdir() (string, bool) {
	names, err := d.f.Readdirnames(1)
	if err != nil || len(names) == 0 {
		return "", false
	}
	return names[0], true
}

// Closedir is closedir().
func (d *Dir) Closedir() {
	_ = d.f.Close()
}

// ReadDir is the names opendir(path) and readdir() return, in readdir()'s
// order (the file system's, not sorted), less "." and "..". opendir()
// fails on a path that is not a directory (ENOTDIR) or cannot be opened. A
// read error after a successful open ends the list, as readdir() returning
// NULL does.
func ReadDir(path string) ([]string, error) {
	d, err := os.OpenFile(path, os.O_RDONLY|syscall.O_DIRECTORY, 0)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	names, _ := d.Readdirnames(-1)
	return names, nil
}
