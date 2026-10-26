// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package clib

import (
	"bytes"
	"io"
	"os"
	"strings"
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

// Int implements %d. It reports false on a matching or input failure.
// Overflow is undefined in C; the value simply wraps here.
func (sc *Scanner) Int() (int, bool) {
	sc.SkipSpace()
	i := sc.pos
	neg := false
	if i < len(sc.s) && (sc.s[i] == '+' || sc.s[i] == '-') {
		neg = sc.s[i] == '-'
		i++
	}
	start := i
	var n int32
	for i < len(sc.s) && sc.s[i] >= '0' && sc.s[i] <= '9' {
		n = n*10 + int32(sc.s[i]-'0')
		i++
	}
	if i == start {
		return 0, false
	}
	sc.pos = i
	if neg {
		n = -n
	}
	return int(n), true
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
