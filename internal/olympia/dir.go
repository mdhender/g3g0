// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/dir.c: determine_map_edges(), which load_db() runs to
// set up the map routines. The rest of dir.c is left for later stages.
//
// A province's box number gives its place on the map: 10101 is row 1,
// column 1, and [10199] row 1, column 99.

// regionRow is region_row(): the map row of location where.
func regionRow(where int) int {
	return (where / 100) % 100
}

// regionCol is region_col(): the map column of location where.
func regionCol(where int) int {
	return where % 100
}

// determineMapEdges is determine_map_edges(): it raises maxMapRow and
// maxMapCol to the largest row and column of every province, in chain
// order, that is not in Faery, Hades, the Cloudlands, the tunnels or
// Nowhere, and sets maxMapInit. A special region not yet created is ignored.
func (w *world) determineMapEdges() {
	for i := range w.loopProvince() {
		if (w.faeryRegion != 0 && w.region(i) == int(w.faeryRegion)) ||
			(w.hadesRegion != 0 && w.region(i) == int(w.hadesRegion)) ||
			(w.cloudRegion != 0 && w.region(i) == int(w.cloudRegion)) ||
			(w.tunnelRegion != 0 && w.region(i) == int(w.tunnelRegion)) ||
			(w.nowhereRegion != 0 && w.region(i) == int(w.nowhereRegion)) {
			continue
		}

		row := regionRow(i)
		col := regionCol(i)

		if row > int(w.maxMapRow) {
			w.maxMapRow = int32(row)
		}

		if col > int(w.maxMapCol) {
			w.maxMapCol = int32(col)
		}
	}

	w.maxMapInit = true
}
