// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package clib

import (
	"math"
	"reflect"
	"strings"
	"testing"
)

// The atoi and strtoull expectations were observed from the C library on
// macOS, the reference platform.

func TestAtoi(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want int
	}{
		{"10", 10},
		{" 10", 10},
		{"+7", 7},
		{"-3", -3},
		{"20x", 20},
		{"abc", 0},
		{"", 0},
		{"4294967306", 10},           // (int)4294967306
		{"-4294967286", 10},          // (int)-4294967286
		{"99999999999999999999", -1}, // strtol saturates to LONG_MAX
		{"-99999999999999999999", 0}, // LONG_MIN truncates to 0
	} {
		if got := Atoi(tc.in); got != tc.want {
			t.Errorf("Atoi(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestStrtoull0(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want uint64
	}{
		{"0", 0},
		{"12345", 12345},
		{"0x10", 16},
		{"0X10", 16},
		{"010", 8},
		{"08", 0},
		{"0x", 0},
		{"0xg", 0},
		{"  +12", 12},
		{"42abc", 42},
		{"nope", 0},
		{"-1", math.MaxUint64},
		{"18446744073709551615", math.MaxUint64},
		{"99999999999999999999999", math.MaxUint64},
		{"-99999999999999999999999", math.MaxUint64},
	} {
		if got := Strtoull0(tc.in); got != tc.want {
			t.Errorf("Strtoull0(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestFgetsLines(t *testing.T) {
	long := strings.Repeat("x", 300)
	got := FgetsLines([]byte("ab\ncd\x00ef\n\n"+long+"\nlast"), 256)
	want := []string{"ab\n", "cd", "\n", long[:255], long[255:] + "\n", "last"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FgetsLines = %q, want %q", got, want)
	}
}

func TestScanner(t *testing.T) {
	tests := []struct {
		in       string
		row, col int
		typ      byte
		n        int // fields assigned, as sscanf would return
	}{
		{"41,6\tp\tEndless Prairie", 41, 6, 'p', 3},
		{"  -3,+7   f", -3, 7, 'f', 3},
		{"4, 28 x", 4, 28, 'x', 3},
		{"4 ,28 x", 4, 0, 0, 1},
		{"12,", 12, 0, 0, 1},
		{"9,9", 9, 9, 0, 2},
		{"", 0, 0, 0, 0},
		{"Great Sea", 0, 0, 0, 0},
	}
	for _, tc := range tests {
		sc := NewScanner(tc.in)
		var row, col, n int
		var typ byte
		if v, ok := sc.Int(); ok {
			row, n = v, 1
			if sc.Literal(',') {
				if v, ok := sc.Int(); ok {
					col, n = v, 2
					sc.SkipSpace()
					if c, ok := sc.Char(); ok {
						typ, n = c, 3
					}
				}
			}
		}
		if row != tc.row || col != tc.col || typ != tc.typ || n != tc.n {
			t.Errorf("scan %q = %d,%d %q (%d), want %d,%d %q (%d)",
				tc.in, row, col, typ, n, tc.row, tc.col, tc.typ, tc.n)
		}
	}
}

func TestPad(t *testing.T) {
	// C pads by bytes; "é" is two bytes.
	if got := PadRight("é", 4); got != "é  " {
		t.Errorf("PadRight = %q", got)
	}
	if got := PadLeft("ab", 4); got != "  ab" {
		t.Errorf("PadLeft = %q", got)
	}
}

func TestGetopt(t *testing.T) {
	type opt struct {
		c   int
		arg string
	}
	for _, tc := range []struct {
		args []string
		want []opt
	}{
		{[]string{"p", "-s", "5", "-b3"}, []opt{{'s', "5"}, {'b', "3"}}},
		{[]string{"p", "foo", "-s", "5"}, nil}, // no permuting
		{[]string{"p", "-s", "5", "foo", "-b", "1"}, []opt{{'s', "5"}}},
		{[]string{"p", "--", "-s", "5"}, nil},
		{[]string{"p", "-", "-s", "5"}, nil},
		{[]string{"p", "--s", "5"}, []opt{{'?', ""}, {'s', "5"}}}, // "-" is unknown; the cluster goes on
		{[]string{"p", "-x"}, []opt{{'?', ""}}},
		{[]string{"p", "-:"}, []opt{{'?', ""}}},
		{[]string{"p", "-s"}, []opt{{'?', ""}}},
		{[]string{"p", "-b1c1"}, []opt{{'b', "1c1"}}},
		{[]string{"p", "-s", "-b"}, []opt{{'s', "-b"}}},
	} {
		var stderr strings.Builder
		g := NewGetopt(tc.args, "b:c:s:", &stderr)
		g.Opterr = false
		var got []opt
		for c := g.Next(); c != -1; c = g.Next() {
			got = append(got, opt{c, g.Optarg})
		}
		if !reflect.DeepEqual(got, tc.want) || stderr.Len() != 0 {
			t.Errorf("getopt %q = %v (stderr %q), want %v", tc.args, got, stderr.String(), tc.want)
		}
	}
}
