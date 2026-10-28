// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// These test the explicit unsupported boundary, NOT successful C admissions.
// C-recorded fetch/parser/failure controls live in command-cases.txt. Complete
// valid input must stop before alloc_box, even for an already occupied ID.
func TestCommandAdmissionBoundaries(t *testing.T) {
	for _, acct := range []string{"aa0", "[50001]", "(50002)", "149999", "200", "1"} {
		for _, directory := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/directory=%v", acct, directory), func(t *testing.T) {
				w := newWorld()
				var stdout, stderr strings.Builder
				w.stdout, w.stderr = &stdout, &stderr
				w.allocBox(200, tPlayer, subPlRegular)
				w.pPlayer(200).noblePoints = 17
				w.allocBox(1, tItem, 0)
				w.allocBox(1000, tChar, 0)
				w.pChar(1000).unitLord = 200
				w.queueOrder(200, 1000, "wait 9")
				snapshot := func() string {
					var s strings.Builder
					for n, b := range w.bx {
						if b != nil {
							w.saveBox(&s, n)
						}
					}
					helperList(&s, "new players", w.newPlayers)
					helperList(&s, "vector", w.outVector)
					return s.String()
				}
				before := snapshot()
				input := "\n Faction \n\n Character \n City \n Full Name \n email@example.test \ntrailing field\n"
				reader := strings.NewReader(input)
				if directory {
					root := t.TempDir()
					t.Chdir(root)
					path := filepath.Join("act", acct)
					if err := os.MkdirAll(path, 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(path, "Join-g3"), []byte(input), 0o644); err != nil {
						t.Fatal(err)
					}
					w.libdir = filepath.Join(root, "not-the-account-directory")
				}
				defer func() {
					if r := recover(); r != notImplemented("new player admission") {
						t.Fatalf("valid admission boundary = %v", r)
					}
					if got := snapshot(); got != before {
						t.Fatalf("unsupported admission changed saved state:\n%s", got)
					}
					if !directory {
						if s := getlin(reader); s == nil || *s != "trailing field" {
							t.Fatalf("parser did not consume exactly five nonblank fields: %v", s)
						}
					}
					if strings.Contains(stderr.String(), "added player") || strings.Contains(stderr.String(), "new_player_list") {
						t.Fatalf("unsupported admission reported success: %q", stderr.String())
					}
				}()
				if directory {
					w.addNewPlayers()
				} else {
					w.makeNewPlayersSup(acct, reader)
				}
			})
		}
	}
}
