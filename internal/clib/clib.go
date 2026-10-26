// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package clib reproduces the C library behavior that the ported programs'
// observable output depends on. It follows the BSD libc that macOS ships,
// which is the reference platform.
package clib

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"syscall"
)

// IsSpace is the C locale isspace(), which scanf uses to skip input.
func IsSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\v' || c == '\f' || c == '\r'
}

// digitValue returns the value of c as a digit in bases up to 36, or 36 if
// c is not a digit.
func digitValue(c byte) uint64 {
	switch {
	case c >= '0' && c <= '9':
		return uint64(c - '0')
	case c >= 'A' && c <= 'Z':
		return uint64(c-'A') + 10
	case c >= 'a' && c <= 'z':
		return uint64(c-'a') + 10
	}
	return 36
}

// Atoi is atoi(3), which is (int)strtol(s, NULL, 10): out-of-range values
// saturate to the 64-bit long range and are then truncated to 32 bits.
func Atoi(s string) int {
	i := 0
	for i < len(s) && IsSpace(s[i]) {
		i++
	}
	neg := false
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		neg = s[i] == '-'
		i++
	}
	var acc uint64
	limit := uint64(math.MaxInt64)
	if neg {
		limit++
	}
	overflow := false
	for ; i < len(s); i++ {
		d := digitValue(s[i])
		if d >= 10 {
			break
		}
		if overflow || acc > (limit-d)/10 {
			overflow = true
			continue
		}
		acc = acc*10 + d
	}
	var n int64
	switch {
	case overflow && neg:
		n = math.MinInt64
	case overflow:
		n = math.MaxInt64
	case neg:
		n = -int64(acc)
	default:
		n = int64(acc)
	}
	return int(int32(n))
}

// Strtoull0 is strtoull(s, NULL, 0): a 0x prefix selects hexadecimal and a
// leading 0 octal. A minus sign negates the result as an unsigned value, and
// out-of-range values saturate to ULLONG_MAX.
func Strtoull0(s string) uint64 {
	i := 0
	for i < len(s) && IsSpace(s[i]) {
		i++
	}
	neg := false
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		neg = s[i] == '-'
		i++
	}
	base := uint64(10)
	if i < len(s) && s[i] == '0' {
		base = 8
		if i+2 < len(s) && (s[i+1] == 'x' || s[i+1] == 'X') && digitValue(s[i+2]) < 16 {
			base = 16
			i += 2
		}
	}
	var acc uint64
	overflow := false
	for ; i < len(s); i++ {
		d := digitValue(s[i])
		if d >= base {
			break
		}
		if overflow || acc > (math.MaxUint64-d)/base {
			overflow = true
			continue
		}
		acc = acc*base + d
	}
	switch {
	case overflow:
		return math.MaxUint64
	case neg:
		return -acc
	}
	return acc
}

// Strerror returns the C library message for err, as perror() prints it.
func Strerror(err error) string {
	if errno, ok := errors.AsType[syscall.Errno](err); ok {
		msg := errno.Error()
		if msg != "" {
			// Go's messages are the C messages with a lowercase first letter.
			return strings.ToUpper(msg[:1]) + msg[1:]
		}
	}
	return err.Error()
}

// Perror writes s and the message for err the way perror(3) does. An empty
// s prints the message alone.
func Perror(w io.Writer, s string, err error) {
	if s != "" {
		fmt.Fprintf(w, "%s: %s\n", s, Strerror(err))
		return
	}
	fmt.Fprintf(w, "%s\n", Strerror(err))
}
