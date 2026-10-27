// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

// Ported from olympia/sout.c.

// sout is sout(): the formatted string, cut to BUF_LEN (LEN) bytes, the
// size of the ring buffers C prints into. C's ring of 50 buffers is not
// needed: each call returns a new string.
func sout(format string, a ...any) string {
	s := fmt.Sprintf(format, a...)
	return s[:min(len(s), lenMax)]
}
