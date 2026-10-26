// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package mapgen

import (
	"bytes"

	"github.com/mdhender/g3g0/internal/clib"
)

// This file holds the helpers from mapgen/z.c and mapgen/z.h that mapgen's
// observable output depends on.

// lineReader reproduces getlin() from mapgen/z.c over an entire file.
type lineReader struct {
	data []byte
}

// getlin returns the next line without its newline, or false at end of input.
func (r *lineReader) getlin() (string, bool) {
	if len(r.data) == 0 {
		return "", false
	}
	if line, rest, found := bytes.Cut(r.data, []byte{'\n'}); found {
		r.data = rest
		return clib.CString(line), true
	}
	line := r.data
	r.data = nil
	// getlin() only truncates the final, unterminated line.
	if len(line) >= 1023 {
		line = line[:1023]
	}
	return clib.CString(line), true
}

// iswhite is the iswhite() macro from mapgen/z.h.
func iswhite(c byte) bool {
	return c == ' ' || c == '\t'
}
