// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

// Ported from olympia/quest.c: create_relics(), which load_db() runs on
// every load. The rest of quest.c (create_nowhere(), the quest stream and
// the QUEST command) is left for later stages.

// createARelic is create_a_relic(): unless box n exists, it creates n as a
// relic in nowhere_loc's inventory and reports it on stderr. No new_ent()
// is involved: the box number is fixed, so no RNG stream is drawn from.
func (w *world) createARelic(n int, name string, use, weight int) {
	if w.bx[n] != nil {
		return
	}

	w.createUniqueItemAlloc(n, int(w.nowhereLoc), subRelic)

	w.setName(n, name)
	w.pItemMagic(n).useKey = int8(use)
	w.pItemMagic(n).lore = int32(n)
	w.pItem(n).weight = int16(weight)

	fmt.Fprintf(w.stderr, "Created %s\n", w.boxName(n))
}

// createRelics is create_relics(): the three relics every game has.
func (w *world) createRelics() {
	w.createARelic(relicThrone, "Imperial Throne", 0, 500)
	w.createARelic(relicCrown, "Crown of Prosperity", 0, 10)
	w.createARelic(relicBtaSkull, "Skull of Bastrestric", useBtaSkull, 10)
}
