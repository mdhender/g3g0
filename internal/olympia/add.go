// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"syscall"

	"github.com/mdhender/g3g0/internal/clib"
)

// Ported from olympia/add.c: rename_act_join_files(), which save_db() runs
// last. The rest of add.c, which adds the new players, is left for later
// stages.

// renameActJoinFiles is rename_act_join_files(): for each player added
// this turn, in the order of newPlayers, it renames act/CODE/Join-g3 to
// act/CODE/Join-g3-, where CODE is the player's box code. The act
// directory is relative to the current directory, not to libdir. A rename
// that fails gets a diagnostic. It is rename(2) itself (syscall.Rename),
// since os.Rename refuses some renames that rename(2) makes.
func (w *world) renameActJoinFiles() {
	acctDir := "act"

	for i := 0; i < w.newPlayers.Len(); i++ {
		pl := w.newPlayers[i]
		acct := boxCodeLess(pl)

		oldName := sout("%s/%s/Join-g3", acctDir, acct)
		newName := sout("%s/%s/Join-g3-", acctDir, acct)

		if err := syscall.Rename(oldName, newName); err != nil {
			fmt.Fprintf(w.stderr, "rename(%s, %s) failed:",
				oldName, newName)
			clib.Perror(w.stderr, "", err)
		}
	}
}
