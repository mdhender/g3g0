// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/mdhender/g3g0/internal/clib"
)

// Ported from olympia/add.c. Admission input and unsuccessful/no-admission
// traversal belong to stage 1; creating new players remains stage 6.

// fetchInp is add.c:34–46. Blank cleaned lines do not consume a field.
func fetchInp(fp io.ByteReader) *string {
	for s := getlinEw(fp); s != nil; s = getlinEw(fp) {
		if *s != "" {
			return s // getlinEw already returns an independent saved string.
		}
	}
	return nil
}

// makeNewPlayersSup preserves add.c:302–325 before the stage-6 allocation.
// Even an invalid account is only inspected AFTER all five input fields.
func (w *world) makeNewPlayersSup(acct string, fp io.ByteReader) bool {
	_ = fetchInp(fp) // faction
	_ = fetchInp(fp) // character
	_ = fetchInp(fp) // start_city
	_ = fetchInp(fp) // full_name
	email := fetchInp(fp)
	if email == nil {
		fmt.Fprintf(w.stderr, "error: partial read for '%s'\n", acct)
		return false
	}
	pl := scode(acct)
	assert(pl > 0 && pl < maxBoxes, "olympia/add.c", 325, "pl > 0 && pl < MAX_BOXES")
	// Immediately before C's alloc_box, not after a fake faction or noble.
	panic(notImplemented("new player admission"))
}

// makeNewPlayers is add.c:336–378. act is cwd-relative, NOT libdir-relative.
// ReadDir preserves native iteration; a successful fopen followed by a read
// error (including a directory) is partial input, not a failed open.
func (w *world) makeNewPlayers() {
	names, err := clib.ReadDir("act")
	if err != nil {
		fmt.Fprint(w.stderr, "make_new_players: can't open act: ")
		clib.Perror(w.stderr, "", err)
		return
	}
	for _, acct := range names {
		if acct[0] == '.' {
			continue
		}
		fp, err := os.Open(sout("act/%s/Join-g3", acct))
		if err != nil {
			continue
		}
		func() {
			defer fp.Close()
			w.makeNewPlayersSup(acct, bufio.NewReader(fp))
		}()
	}
}

// addNewPlayers is add.c:722–734's unsuccessful/no-admission traversal.
// The intervening new-character/banner/report/template loops are empty:
// complete admissions stop above, and must never fabricate report success.
func (w *world) addNewPlayers() {
	w.stage("add_new_players()")
	w.makeNewPlayers()
	if w.newPlayers.Len() != 0 {
		panic(notImplemented("new player reports"))
	}
	w.stage("new_player_list()")
	// new_player_list calls sout.c:568–582 vector_players even with no
	// admissions. Silent and system players participate; only these two skip.
	w.outVector.Clear()
	for pl := range w.loopKind(tPlayer) {
		if pl != int(w.eatPl) && pl != int(w.skillPlayer) {
			w.outVector.Append(pl)
		}
	}
}

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
