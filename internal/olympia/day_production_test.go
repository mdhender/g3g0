// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"strings"
	"testing"
)

// This is a boundary/order check, not a successful composition recording.
// Real C helper results are independently checked by command-compat.
func TestPostProductionStopsAtLocationTrades(t *testing.T) {
	w := newWorld()
	var stderr strings.Builder
	w.stderr = &stderr
	w.allocBox(9000, tLoc, subForest)
	w.allocBox(9100, tLoc, subCity)
	w.allocBox(itemLumber, tItem, 0)
	w.allocBox(itemTaxCookie, tItem, 0)
	w.setWhere(9100, 9000)
	w.pLoc(9000).civ = 9
	w.genItem(9000, itemTaxCookie, 7)
	defer func() {
		if r := recover(); r != notImplemented("location_trades") {
			t.Fatalf("postProduction boundary = %v", r)
		}
		if stderr.String() != "compute_civ_levels()\n" || w.pLoc(9000).civ == 9 {
			t.Fatal("postProduction did not run real civilization calculation first")
		}
		if w.postHasBeenRun != 0 || w.hasItem(9000, itemLumber) != 0 || w.hasItem(9000, itemTaxCookie) != 7 {
			t.Fatal("postProduction reordered past unfinished location_trades")
		}
	}()
	w.postProduction()
}
