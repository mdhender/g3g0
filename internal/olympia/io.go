// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

// Ported from olympia/io.c.

// ifName is if_name(): " (name)" for box n, to make the data files easier
// to read, if pretty_data_files is set and the box exists and has a
// non-empty name; otherwise "".
func (w *world) ifName(num int) string {
	if !w.prettyDataFiles {
		return ""
	}

	if !w.validBox(num) {
		return ""
	}

	s := w.name(num)

	if s != "" {
		return fmt.Sprintf(" (%s)", s)
	}

	return ""
}
