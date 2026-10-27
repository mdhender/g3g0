// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package clib

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestScannerMatchesC runs testdata/scanf-cases.txt through Scanner and
// requires the output that tools/scanf-compat/run.sh recorded from
// sscanf() on macOS, the reference platform, in testdata/scanf-c-output.txt.
func TestScannerMatchesC(t *testing.T) {
	cases, err := os.ReadFile(filepath.Join("testdata", "scanf-cases.txt"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "scanf-c-output.txt"))
	if err != nil {
		t.Fatal(err)
	}
	got := runScanfCases(t, cases)

	wantLines := strings.Split(string(want), "\n")
	gotLines := strings.Split(string(got), "\n")
	for i := 0; i < len(wantLines) || i < len(gotLines); i++ {
		var w, g string
		if i < len(wantLines) {
			w = wantLines[i]
		}
		if i < len(gotLines) {
			g = gotLines[i]
		}
		if w != g {
			t.Fatalf("line %d:\n got %q\nwant %q", i+1, g, w)
		}
	}
}

// runScanfCases is the Go side of the C driver in tools/scanf-compat: the
// same cases, decoded the same way, with the same output.
func runScanfCases(t *testing.T, cases []byte) []byte {
	t.Helper()
	var out bytes.Buffer
	sc := bufio.NewScanner(bytes.NewReader(cases))
	for sc.Scan() {
		line := sc.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		format, raw, ok := strings.Cut(line, "\t")
		if !ok {
			t.Fatalf("no tab: %q", line)
		}
		input := decodeScanfInput(t, raw)

		// The targets start as the C driver's do; n is %n, last in every
		// format, so it is set only when every directive matched.
		a, b, c, n := 7, 7, 7, -1
		s := NewScanner(input)
		r := 0
		switch format {
		case "d", "dd":
			if v, ok := s.Int(); ok {
				a, r = v, 1
				if format == "d" {
					n = s.pos
				} else if v, ok := s.Int(); ok {
					b, r, n = v, 2, s.pos
				}
			}
		case "d,d", "d,d c":
			if v, ok := s.Int(); ok {
				a, r = v, 1
				if s.Literal(',') {
					if v, ok := s.Int(); ok {
						b, r = v, 2
						if format == "d,d" {
							n = s.pos
						} else {
							s.SkipSpace()
							if ch, ok := s.Char(); ok {
								c, r, n = int(int8(ch)), 3, s.pos
							}
						}
					}
				}
			}
		default:
			t.Fatalf("unknown format %q", format)
		}
		// sscanf returns EOF for an input failure before the first
		// conversion: nothing but white space was left for it.
		if r == 0 && s.pos == len(input) {
			r = -1
		}
		fmt.Fprintf(&out, "%s\t%s\t-> r=%d a=%d b=%d c=%d n=%d\n", format, raw, r, a, b, c, n)
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}

// decodeScanfInput decodes the escapes of scanf-cases.txt: \xHH, \\ and
// {N:c}.
func decodeScanfInput(t *testing.T, s string) string {
	t.Helper()
	decode1 := func() byte {
		if strings.HasPrefix(s, `\x`) && len(s) >= 4 {
			v, err := strconv.ParseUint(s[2:4], 16, 8)
			if err != nil {
				t.Fatalf("bad escape in %q", s)
			}
			s = s[4:]
			return byte(v)
		}
		if strings.HasPrefix(s, `\\`) {
			s = s[2:]
			return '\\'
		}
		c := s[0]
		s = s[1:]
		return c
	}
	var b strings.Builder
	for s != "" {
		if s[0] != '{' {
			b.WriteByte(decode1())
			continue
		}
		count, rest, ok := strings.Cut(s[1:], ":")
		if !ok {
			t.Fatalf("bad repeat in %q", s)
		}
		k, err := strconv.Atoi(count)
		if err != nil || rest == "" {
			t.Fatalf("bad repeat in %q", s)
		}
		s = rest
		c := decode1()
		if !strings.HasPrefix(s, "}") {
			t.Fatalf("bad repeat in %q", s)
		}
		s = s[1:]
		b.WriteString(strings.Repeat(string([]byte{c}), k))
	}
	return b.String()
}
