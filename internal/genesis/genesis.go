// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package genesis ports the Olympia G3 world generator (genesis/genesis.c)
// from C.
//
// The port is a compatibility port: given the same arguments it writes
// byte-identical Map, Regions, Land, Cities and randseed files, the same
// --image BMP, and the same diagnostics, as G3's genesis-g3. Function and
// variable names follow the C source so the two can be read side by side.
package genesis

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/mdhender/g3g0/internal/clib"
)

const (
	maxDim      = 100 // mapgen MAX_ROW/MAX_COL; size is 10..99
	distanceCap = 9
	numTerrains = 5

	// Continental-shelf and border widths (island-g3 defaults).
	shelfWidth  = 3
	borderWidth = 2

	// Per-step PCG stream tags.
	stepPartition   = 1
	stepNames       = 2
	stepCities      = 3
	stepRandseed    = 4
	stepOrigins     = 5   // --defer-oceans origin selection and frontier
	stepStartCities = 6   // start-city placement
	stepIslandBase  = 100 // island streams: stepIslandBase + iteration

	// --defer-oceans tuning.
	alphaSq       = 0.4225 // (0.65)^2: minimum origin spacing
	lakeMax       = 12     // a landlocked water body smaller than this is a lake
	minNamedWater = 4      // landlocked water bodies smaller than this become land
)

type location struct {
	x, y int
}

type terrain struct {
	symbol           byte
	min, max         int
	targetProb, prob int
}

// islandSizeTab maps an iteration (1..10) to its island size; iteration N
// grows N islands. Index 0 is unused.
var islandSizeTab = [11]int{0, 521, 257, 257, 127, 61, 61, 31, 17, 11, 7}

// terrainsBase holds the original terrain weights. growOneIsland normalizes
// a fresh copy for each island.
var terrainsBase = [numTerrains]terrain{
	{'p', 12, 30, 30, 0},
	{'f', 6, 14, 30, 0},
	{'m', 6, 10, 20, 0},
	{'d', 15, 30, 10, 0},
	{'s', 1, 3, 10, 0},
}

// The four plain ocean glyphs, indexed by 4-coloring result 0..3.
var oceanGlyph = [4]byte{',', '.', ' ', '\''}

// mapgen's 8-neighbor offsets. The cardinal directions are at even indexes.
var (
	dr8 = [8]int{-1, -1, 0, 1, 1, 1, 0, -1}
	dc8 = [8]int{0, 1, 1, 1, 0, -1, -1, -1}
)

// exitStatus unwinds the generator for exit(3).
type exitStatus int

// assertion unwinds the generator for a failed assert().
type assertion struct {
	function string
	line     int
	cond     string
}

// generator holds what genesis.c keeps in static variables.
type generator struct {
	dir    string
	stderr io.Writer

	gSize    int
	mapp     [maxDim][maxDim]byte
	oceanID  [maxDim][maxDim]int
	working  [maxDim][maxDim]byte
	distance [maxDim][maxDim]int
	ids      [maxDim][maxDim]int

	islandCells [maxDim * maxDim]location

	label      [maxDim][maxDim]int
	floodStack [maxDim * maxDim]int

	gDeferOceans    bool
	gOceanBFS       bool
	deferN          int
	oceanSize       [maxDim]int
	oceanOriginR    [maxDim]int
	oceanOriginC    [maxDim]int
	oceanLandlocked [maxDim]bool

	regionOrder  [len(regionNames)]int
	regionPoolN  int
	regionCursor int
}

func (g *generator) exit(status int) {
	panic(exitStatus(status))
}

func (g *generator) assert(ok bool, function string, line int, cond string) {
	if !ok {
		panic(assertion{function: function, line: line, cond: cond})
	}
}

// path resolves a file name against the working directory.
func (g *generator) path(name string) string {
	if filepath.IsAbs(name) {
		return name
	}
	return filepath.Join(g.dir, name)
}

// ---------------------------------------------------------------------------
// Glyph helpers (must agree with mapgen.c's reader).
// ---------------------------------------------------------------------------

func isOceanGlyph(c byte) bool {
	switch c {
	case ',', '.', ' ', '\'', ';', ':', '~', '"':
		return true
	}
	return false
}

// oceanColor is mapgen's ocean color (1..4), or 0 for land.
func oceanColor(c byte) int {
	switch c {
	case ',', ';':
		return 1
	case '.', ':':
		return 2
	case ' ', '~':
		return 3
	case '\'', '"':
		return 4
	}
	return 0
}

// neighbor8 applies mapgen's 8-neighbor rule: columns wrap, rows stop at
// the edge. It reports whether the neighbor is on the map.
func (g *generator) neighbor8(r, c, d int) (int, int, bool) {
	rr := r + dr8[d]
	cc := c + dc8[d]
	if cc < 0 {
		cc = g.gSize - 1
	} else if cc > g.gSize-1 {
		cc = 0
	}
	if rr < 0 || rr > g.gSize-1 {
		return 0, 0, false
	}
	return rr, cc, true
}

// ---------------------------------------------------------------------------
// Partition the canvas into oceans: a randomized frontier flood over the
// cardinal neighbors, without wrapping.
// ---------------------------------------------------------------------------

func (g *generator) partitionOceans(nOceans int, rng *pcg32) {
	var frontier [maxDim * maxDim]int
	fn := 0
	total := g.gSize * g.gSize
	claimed := 0

	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			g.oceanID[r][c] = -1
		}
	}

	// Seed each ocean at a distinct random cell.
	for o := 0; o < nOceans; o++ {
		var idx int
		for {
			idx = rng.rangeInt(0, total-1)
			if g.oceanID[idx/g.gSize][idx%g.gSize] == -1 {
				break
			}
		}
		g.oceanID[idx/g.gSize][idx%g.gSize] = o
		claimed++
		frontier[fn] = idx
		fn++
	}

	// Grow from the frontier until every cell is claimed.
	for claimed < total {
		var nb [4]int

		g.assert(fn > 0, "partition_oceans", 214, "fn > 0")
		fi := rng.rangeInt(0, fn-1)
		cell := frontier[fi]
		cr := cell / g.gSize
		cc := cell % g.gSize

		k := 0
		if cr > 0 && g.oceanID[cr-1][cc] == -1 {
			nb[k] = (cr-1)*g.gSize + cc
			k++
		}
		if cr < g.gSize-1 && g.oceanID[cr+1][cc] == -1 {
			nb[k] = (cr+1)*g.gSize + cc
			k++
		}
		if cc > 0 && g.oceanID[cr][cc-1] == -1 {
			nb[k] = cr*g.gSize + (cc - 1)
			k++
		}
		if cc < g.gSize-1 && g.oceanID[cr][cc+1] == -1 {
			nb[k] = cr*g.gSize + (cc + 1)
			k++
		}

		if k == 0 {
			fn--
			frontier[fi] = frontier[fn] // drop the boxed-in cell
			continue
		}

		i := nb[rng.rangeInt(0, k-1)]
		g.oceanID[i/g.gSize][i%g.gSize] = g.oceanID[cr][cc]
		claimed++
		frontier[fn] = i
		fn++
	}
}

// ---------------------------------------------------------------------------
// 4-color the oceans by deterministic backtracking over the mapgen-rule
// adjacency graph.
// ---------------------------------------------------------------------------

func colorSolve(v, n int, adj *[maxDim][maxDim]bool, colors *[maxDim]int) bool {
	if v == n {
		return true
	}
	for k := 0; k < 4; k++ {
		ok := true
		for u := 0; u < n; u++ {
			if adj[v][u] && colors[u] == k {
				ok = false
				break
			}
		}
		if ok {
			colors[v] = k
			if colorSolve(v+1, n, adj, colors) {
				return true
			}
			colors[v] = -1
		}
	}
	return false
}

// colorOceans paints the ocean glyphs and reports true, or reports false if
// the oceans are not 4-colorable.
func (g *generator) colorOceans(nOceans int) bool {
	var adj [maxDim][maxDim]bool
	var colors [maxDim]int

	for i := 0; i < nOceans; i++ {
		colors[i] = -1
	}

	// Build the ocean-adjacency graph with mapgen's 8-neighbor rule.
	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			for d := 0; d < 8; d++ {
				nr, nc, ok := g.neighbor8(r, c, d)
				if !ok {
					continue
				}
				a := g.oceanID[r][c]
				b := g.oceanID[nr][nc]
				if a != b {
					adj[a][b] = true
					adj[b][a] = true
				}
			}
		}
	}

	if !colorSolve(0, nOceans, &adj, &colors) {
		return false
	}

	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			g.mapp[r][c] = oceanGlyph[colors[g.oceanID[r][c]]]
		}
	}
	return true
}

// oceanCap is min(ceil(size / 3), 12).
func (g *generator) oceanCap() int {
	c := (g.gSize + 2) / 3
	if c > 12 {
		return 12
	}
	return c
}

// addOceans partitions and colors the oceans, clamping the request to the
// cap and retrying with one fewer ocean until the oceans color. It returns
// the number of oceans used.
func (g *generator) addOceans(requested int, seed uint64) int {
	var rng pcg32
	oceanCap := g.oceanCap()
	n := requested

	if n > oceanCap {
		fmt.Fprintf(g.stderr, "genesis: --number-of-oceans %d exceeds the cap for size %d; using %d\n",
			requested, g.gSize, oceanCap)
		n = oceanCap
	}

	for ; n >= 1; n-- {
		// Reseed the partition stream each attempt.
		rng.seedStep(seed, stepPartition, 0)
		g.partitionOceans(n, &rng)
		if g.colorOceans(n) {
			return n
		}
		fmt.Fprintf(g.stderr, "genesis: %d oceans not 4-colorable on this map; retrying with %d\n", n, n-1)
	}

	fmt.Fprintf(g.stderr, "genesis: internal error: could not color even one ocean\n")
	g.exit(1)
	return 0
}

// ---------------------------------------------------------------------------
// --defer-oceans: generate the oceans after the islands and Mt. Olympus.
// ---------------------------------------------------------------------------

// fillOceanCanvas resets the whole canvas to open water.
func (g *generator) fillOceanCanvas() {
	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			g.mapp[r][c] = ' '
		}
	}
}

// dist2 is the squared Euclidean distance, without column wrap.
func dist2(r1, c1, r2, c2 int) int64 {
	dr := int64(r1 - r2)
	dc := int64(c1 - c2)
	return dr*dr + dc*dc
}

// dist2Wrap is the squared distance with mapgen's column wrap.
func (g *generator) dist2Wrap(r1, c1, r2, c2 int) int64 {
	dr := int64(r1 - r2)
	var dc int64
	if c1 > c2 {
		dc = int64(c1 - c2)
	} else {
		dc = int64(c2 - c1)
	}
	if dc > int64(g.gSize)-dc {
		dc = int64(g.gSize) - dc
	}
	return dr*dr + dc*dc
}

// bestCandidates is Mitchell's best-candidate sampling over the wn
// candidate cells (wr, wc). It picks n points spaced at least sqrt(r2)
// apart into orow and ocol, or reports false if a point can't clear r2.
func (g *generator) bestCandidates(n int, r2 int64, wrap bool, rng *pcg32,
	wr, wc []int, wn int, orow, ocol []int) bool {
	if wn < n {
		return false
	}

	for i := 0; i < n; i++ {
		best := -1
		bestD2 := int64(-1)

		if i == 0 {
			best = rng.rangeInt(0, wn-1)
		} else {
			budget := n * n
			for a := 0; a < budget; a++ {
				cand := rng.rangeInt(0, wn-1)
				d2 := int64(math.MaxInt64)
				for j := 0; j < i; j++ {
					var dd int64
					if wrap {
						dd = g.dist2Wrap(wr[cand], wc[cand], orow[j], ocol[j])
					} else {
						dd = dist2(wr[cand], wc[cand], orow[j], ocol[j])
					}
					if dd < d2 {
						d2 = dd
					}
				}
				if d2 > bestD2 {
					bestD2 = d2
					best = cand
				}
				if d2 >= r2 {
					break // clears the spacing floor
				}
			}
			if bestD2 < r2 {
				return false
			}
		}
		orow[i] = wr[best]
		ocol[i] = wc[best]
	}
	return true
}

// placeOrigins picks n ocean-origin cells among the water (' ') cells.
func (g *generator) placeOrigins(n int, r2 int64, rng *pcg32, orow, ocol []int) bool {
	wr := make([]int, maxDim*maxDim)
	wc := make([]int, maxDim*maxDim)
	wn := 0

	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			if g.mapp[r][c] == ' ' {
				wr[wn] = r
				wc[wn] = c
				wn++
			}
		}
	}
	return g.bestCandidates(n, r2, false, rng, wr, wc, wn, orow, ocol)
}

// growOceansBFS is --ocean-growth bfs: a multi-source breadth-first search
// over the water, one shuffled ring at a time. ring holds the rn seeded
// origins on entry.
func (g *generator) growOceansBFS(ring []int, rn int, rng *pcg32) {
	next := make([]int, maxDim*maxDim)

	for rn > 0 {
		for i := rn - 1; i > 0; i-- {
			j := rng.rangeInt(0, i)
			ring[i], ring[j] = ring[j], ring[i]
		}
		nn := 0
		for i := 0; i < rn; i++ {
			cr := ring[i] / g.gSize
			cc := ring[i] % g.gSize
			for d := 0; d < 8; d += 2 { // cardinal only
				nr, nc, ok := g.neighbor8(cr, cc, d)
				if !ok {
					continue
				}
				if g.mapp[nr][nc] == ' ' && g.oceanID[nr][nc] == -1 {
					g.oceanID[nr][nc] = g.oceanID[cr][cc]
					next[nn] = nr*g.gSize + nc
					nn++
				}
			}
		}
		copy(ring, next[:nn])
		rn = nn
	}
}

// partitionOceansDeferred grows each of the n origins into an ocean through
// the water, then sweeps for water no origin reached and makes each such
// body its own ocean. It sets deferN.
func (g *generator) partitionOceansDeferred(n int, orow, ocol []int, rng *pcg32) {
	frontier := make([]int, maxDim*maxDim)
	fn := 0

	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			g.oceanID[r][c] = -1
		}
	}

	for i := 0; i < n; i++ {
		g.oceanID[orow[i]][ocol[i]] = i
		frontier[fn] = orow[i]*g.gSize + ocol[i]
		fn++
	}

	if g.gOceanBFS {
		g.growOceansBFS(frontier, fn, rng)
	} else {
		for fn > 0 {
			fi := rng.rangeInt(0, fn-1)
			cell := frontier[fi]
			cr := cell / g.gSize
			cc := cell % g.gSize
			var nb [8]int
			k := 0

			// Cardinal neighbors only; diagonal growth would interleave
			// oceans into a checkerboard.
			for d := 0; d < 8; d += 2 {
				nr, nc, ok := g.neighbor8(cr, cc, d)
				if !ok {
					continue
				}
				if g.mapp[nr][nc] == ' ' && g.oceanID[nr][nc] == -1 {
					nb[k] = nr*g.gSize + nc
					k++
				}
			}
			if k == 0 {
				fn--
				frontier[fi] = frontier[fn] // boxed in
				continue
			}
			i := nb[rng.rangeInt(0, k-1)]
			g.oceanID[i/g.gSize][i%g.gSize] = g.oceanID[cr][cc]
			frontier[fn] = i
			fn++
		}
	}

	// Discover origin-less water bodies.
	g.deferN = n
	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			if g.mapp[r][c] != ' ' || g.oceanID[r][c] != -1 {
				continue
			}
			top := 0
			g.oceanID[r][c] = g.deferN
			g.floodStack[top] = r*g.gSize + c
			top++
			for top > 0 {
				top--
				pc := g.floodStack[top]
				pr := pc / g.gSize
				qc := pc % g.gSize
				for d := 0; d < 8; d++ {
					nr, nc, ok := g.neighbor8(pr, qc, d)
					if !ok {
						continue
					}
					if g.mapp[nr][nc] == ' ' && g.oceanID[nr][nc] == -1 {
						g.oceanID[nr][nc] = g.deferN
						g.floodStack[top] = nr*g.gSize + nc
						top++
					}
				}
			}
			g.deferN++
		}
	}
}

// colorOceansDeferred 4-colors the deferred oceans, painting only water.
func (g *generator) colorOceansDeferred() bool {
	var adj [maxDim][maxDim]bool
	var colors [maxDim]int

	g.assert(g.deferN <= maxDim, "color_oceans_deferred", 632, "defer_n <= MAX_DIM")

	for i := 0; i < g.deferN; i++ {
		colors[i] = -1
	}

	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			if g.oceanID[r][c] < 0 {
				continue
			}
			for d := 0; d < 8; d++ {
				nr, nc, ok := g.neighbor8(r, c, d)
				if !ok {
					continue
				}
				if g.oceanID[nr][nc] < 0 {
					continue
				}
				a := g.oceanID[r][c]
				b := g.oceanID[nr][nc]
				if a != b {
					adj[a][b] = true
					adj[b][a] = true
				}
			}
		}
	}

	if !colorSolve(0, g.deferN, &adj, &colors) {
		return false
	}

	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			if g.oceanID[r][c] >= 0 {
				g.mapp[r][c] = oceanGlyph[colors[g.oceanID[r][c]]]
			}
		}
	}
	return true
}

// computeOceanMetadata records each ocean's size, whether it is landlocked,
// and its anchor cell: the origin for seeded oceans, otherwise the first
// cell in row-major order.
func (g *generator) computeOceanMetadata(seededN int, orow, ocol []int) {
	for i := 0; i < g.deferN; i++ {
		g.oceanSize[i] = 0
		g.oceanLandlocked[i] = true
		g.oceanOriginR[i] = -1
		g.oceanOriginC[i] = -1
	}
	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			id := g.oceanID[r][c]
			if id < 0 {
				continue
			}
			g.oceanSize[id]++
			if r == 0 || r == g.gSize-1 {
				g.oceanLandlocked[id] = false
			}
			if g.oceanOriginR[id] < 0 {
				g.oceanOriginR[id] = r
				g.oceanOriginC[id] = c
			}
		}
	}
	for i := 0; i < seededN; i++ {
		g.oceanOriginR[i] = orow[i]
		g.oceanOriginC[i] = ocol[i]
	}
}

// fillTinyWater turns landlocked water bodies smaller than minNamedWater
// into 'o' land and returns how many it filled.
func (g *generator) fillTinyWater() int {
	var size [maxDim]int
	var lock [maxDim]bool
	filled := 0

	for i := 0; i < g.deferN; i++ {
		size[i] = 0
		lock[i] = true
	}
	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			id := g.oceanID[r][c]
			if id < 0 {
				continue
			}
			size[id]++
			if r == 0 || r == g.gSize-1 {
				lock[id] = false
			}
		}
	}
	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			id := g.oceanID[r][c]
			if id < 0 {
				continue
			}
			if size[id] < minNamedWater && lock[id] {
				g.mapp[r][c] = 'o' // random terrain; merges into the enclosing land
				g.oceanID[r][c] = -1
			}
		}
	}
	for i := 0; i < g.deferN; i++ {
		if size[i] < minNamedWater && lock[i] {
			filled++
		}
	}
	return filled
}

// deferOceans picks origins, partitions and colors, retrying with one fewer
// seeded ocean on a placement or coloring failure. The spacing floor is
// computed once from the requested count.
func (g *generator) deferOceans(requested int, seed uint64) {
	orow := make([]int, maxDim)
	ocol := make([]int, maxDim)
	oceanCap := g.oceanCap()
	n0 := requested
	wn := 0

	if n0 > oceanCap {
		fmt.Fprintf(g.stderr, "genesis: --number-of-oceans %d exceeds the cap for size %d; using %d\n",
			requested, g.gSize, oceanCap)
		n0 = oceanCap
	}

	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			if g.mapp[r][c] == ' ' {
				wn++
			}
		}
	}
	if wn < 1 {
		fmt.Fprintf(g.stderr, "genesis: no open water left for oceans\n")
		g.exit(1)
	}

	r2 := int64(alphaSq * float64(wn) / float64(n0))
	if r2 < 1 {
		r2 = 1 // keep origins on distinct cells
	}

	for n := n0; n >= 1; n-- {
		var rng pcg32
		rng.seedStep(seed, stepOrigins, uint64(n))
		if !g.placeOrigins(n, r2, &rng, orow, ocol) {
			fmt.Fprintf(g.stderr, "genesis: couldn't place %d ocean origins; retrying with %d\n", n, n-1)
			continue
		}
		g.partitionOceansDeferred(n, orow, ocol, &rng)
		nfilled := g.fillTinyWater()
		if g.colorOceansDeferred() {
			g.computeOceanMetadata(n, orow, ocol)
			fmt.Fprintf(g.stderr, "genesis: deferred oceans: %d seeded, %d total, %d tiny filled\n",
				n, g.deferN, nfilled)
			return
		}
		fmt.Fprintf(g.stderr, "genesis: %d deferred oceans not 4-colorable; retrying with %d\n", n, n-1)
	}

	fmt.Fprintf(g.stderr, "genesis: internal error: could not color even one deferred ocean\n")
	g.exit(1)
}

// ---------------------------------------------------------------------------
// Grow one island onto the current map (the island-g3 algorithm on PCG).
// ---------------------------------------------------------------------------

func (g *generator) makeShelf(y, x, shelf int) {
	if g.working[y][x] == '~' {
		g.working[y][x] = '_'
	}
	if shelf < 1 {
		return
	}
	if y > 0 {
		g.makeShelf(y-1, x, shelf-1)
	}
	if y < g.gSize-1 {
		g.makeShelf(y+1, x, shelf-1)
	}
	if x > 0 {
		g.makeShelf(y, x-1, shelf-1)
	}
	if x < g.gSize-1 {
		g.makeShelf(y, x+1, shelf-1)
	}
}

func (g *generator) extendDistance(y, x int) {
	if y > 0 && g.distance[y-1][x] > g.distance[y][x]+1 {
		g.distance[y-1][x] = g.distance[y][x] + 1
		g.extendDistance(y-1, x)
	}
	if y < g.gSize-1 && g.distance[y+1][x] > g.distance[y][x]+1 {
		g.distance[y+1][x] = g.distance[y][x] + 1
		g.extendDistance(y+1, x)
	}
	if x > 0 && g.distance[y][x-1] > g.distance[y][x]+1 {
		g.distance[y][x-1] = g.distance[y][x] + 1
		g.extendDistance(y, x-1)
	}
	if x < g.gSize-1 && g.distance[y][x+1] > g.distance[y][x]+1 {
		g.distance[y][x+1] = g.distance[y][x] + 1
		g.extendDistance(y, x+1)
	}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

// growOneIsland returns the number of provinces actually grown, which may
// be less than target.
func (g *generator) growOneIsland(target int, rng *pcg32) int {
	xSize, ySize := g.gSize, g.gSize
	cells := &g.islandCells
	working := &g.working

	if target < 1 {
		return 0
	}

	terrains := terrainsBase // a fresh copy of the weights for each island

	// 1. Classify: ocean is available ('~'), everything else is land ('p').
	for y := 0; y < ySize; y++ {
		for x := 0; x < xSize; x++ {
			if isOceanGlyph(g.mapp[y][x]) {
				working[y][x] = '~'
			} else {
				working[y][x] = 'p'
			}
			g.ids[y][x] = 0
		}
	}

	// 2. Continental shelf: exclude water within shelfWidth of land.
	for y := 0; y < ySize; y++ {
		for x := 0; x < xSize; x++ {
			if working[y][x] == 'p' {
				g.makeShelf(y, x, shelfWidth)
			}
		}
	}

	// 3. Border: exclude the outer borderWidth rows and columns.
	for y := 0; y < ySize; y++ {
		for x := 0; x < borderWidth; x++ {
			if working[y][x] == '~' {
				working[y][x] = '_'
			}
			if working[y][xSize-x-1] == '~' {
				working[y][xSize-x-1] = '_'
			}
		}
	}
	for x := 0; x < xSize; x++ {
		for y := 0; y < borderWidth; y++ {
			if working[y][x] == '~' {
				working[y][x] = '_'
			}
			if working[ySize-y-1][x] == '~' {
				working[ySize-y-1][x] = '_'
			}
		}
	}

	// 4. Distance field: distance to the nearest non-water, capped.
	for y := 0; y < ySize; y++ {
		for x := 0; x < xSize; x++ {
			if working[y][x] == '~' {
				g.distance[y][x] = distanceCap
			} else {
				g.distance[y][x] = 0
			}
		}
	}
	for y := 0; y < ySize; y++ {
		for x := 0; x < xSize; x++ {
			g.extendDistance(y, x)
		}
	}

	// 5a. Seed: pick a random cell at the maximum distance.
	maxd := 0
	count := 0
	islandSize := 0
	for y := 0; y < ySize; y++ {
		for x := 0; x < xSize; x++ {
			if g.distance[y][x] > maxd {
				maxd = g.distance[y][x]
				count = 1
			} else if g.distance[y][x] == maxd {
				count++
			}
		}
	}
	if maxd < 1 || count < 1 {
		return 0 // no open water at all
	}

	d := rng.rangeInt(1, count)
	for y := 0; d > 0 && y < ySize; y++ {
		for x := 0; d > 0 && x < xSize; x++ {
			if g.distance[y][x] == maxd {
				d--
				if d == 0 {
					cells[0] = location{x: x, y: y}
					islandSize = 1
					working[y][x] = 'o'
					g.mapp[y][x] = 'o'
					g.ids[y][x] = 0
				}
			}
		}
	}

	// 5b. Grow: claim a random available cardinal neighbor until target.
	// Each neighbor is counted once per adjacent island cell, and the C
	// code's post-decrement test (!d--) picks the d'th of them.
	for islandSize < target {
		count = 0
		for i := 0; i < islandSize; i++ {
			c := cells[i]
			if c.x > 0 && working[c.y][c.x-1] == '~' {
				count++
			}
			if c.x < xSize-1 && working[c.y][c.x+1] == '~' {
				count++
			}
			if c.y > 0 && working[c.y-1][c.x] == '~' {
				count++
			}
			if c.y < ySize-1 && working[c.y+1][c.x] == '~' {
				count++
			}
		}
		if count < 1 {
			break // boxed in: stop early
		}

		d = rng.rangeInt(0, count-1)
		for i := 0; i < islandSize; i++ {
			c := cells[i]
			if c.x > 0 && working[c.y][c.x-1] == '~' {
				if d == 0 {
					cells[islandSize] = location{x: c.x - 1, y: c.y}
				}
				d--
			}
			if c.x < xSize-1 && working[c.y][c.x+1] == '~' {
				if d == 0 {
					cells[islandSize] = location{x: c.x + 1, y: c.y}
				}
				d--
			}
			if c.y > 0 && working[c.y-1][c.x] == '~' {
				if d == 0 {
					cells[islandSize] = location{x: c.x, y: c.y - 1}
				}
				d--
			}
			if c.y < ySize-1 && working[c.y+1][c.x] == '~' {
				if d == 0 {
					cells[islandSize] = location{x: c.x, y: c.y + 1}
				}
				d--
			}
		}
		nc := cells[islandSize]
		working[nc.y][nc.x] = 'o'
		g.mapp[nc.y][nc.x] = 'o'
		g.ids[nc.y][nc.x] = islandSize
		islandSize++
	}

	// 5c. Terrain: integer weights by GCD and LCM, then weighted flood-fill
	// clusters.
	for i := 0; i < numTerrains; i++ {
		terrains[i].prob = terrains[i].min + terrains[i].max
		div := gcd(terrains[i].targetProb, terrains[i].prob)
		terrains[i].targetProb /= div
		terrains[i].prob /= div
	}
	mult := 1
	for i := 0; i < numTerrains; i++ {
		mult = lcm(mult, terrains[i].prob)
	}
	count = 0
	for i := 0; i < numTerrains; i++ {
		terrains[i].targetProb *= mult / terrains[i].prob
		count += terrains[i].targetProb
	}

	o := islandSize
	size := 0
	terr := 0
	for o > 0 {
		clusterEnd := o
		o--
		d = rng.rangeInt(0, o)
		cells[d], cells[o] = cells[o], cells[d]
		g.ids[cells[o].y][cells[o].x] = o
		g.ids[cells[d].y][cells[d].x] = d
		if size < 1 {
			d = rng.rangeInt(1, count)
			for terr = 0; d > terrains[terr].targetProb; terr++ {
				d -= terrains[terr].targetProb
			}
			size = rng.rangeInt(terrains[terr].min, terrains[terr].max)
		}
		g.mapp[cells[o].y][cells[o].x] = terrains[terr].symbol
		working[cells[o].y][cells[o].x] = terrains[terr].symbol
		size--
		if size > o {
			size = o
		}
		for size > 0 {
			nopt := 0
			for i := o; i < clusterEnd; i++ {
				c := cells[i]
				if c.y > 0 && working[c.y-1][c.x] == 'o' {
					nopt++
				}
				if c.y < ySize-1 && working[c.y+1][c.x] == 'o' {
					nopt++
				}
				if c.x > 0 && working[c.y][c.x-1] == 'o' {
					nopt++
				}
				if c.x < xSize-1 && working[c.y][c.x+1] == 'o' {
					nopt++
				}
			}
			if nopt < 1 {
				break
			}
			d = rng.rangeInt(0, nopt-1)
			id := 0
		pick:
			for i := o; i < clusterEnd; i++ {
				c := cells[i]
				if c.y > 0 && working[c.y-1][c.x] == 'o' {
					if d == 0 {
						id = g.ids[c.y-1][c.x]
						break pick
					}
					d--
				}
				if c.y < ySize-1 && working[c.y+1][c.x] == 'o' {
					if d == 0 {
						id = g.ids[c.y+1][c.x]
						break pick
					}
					d--
				}
				if c.x > 0 && working[c.y][c.x-1] == 'o' {
					if d == 0 {
						id = g.ids[c.y][c.x-1]
						break pick
					}
					d--
				}
				if c.x < xSize-1 && working[c.y][c.x+1] == 'o' {
					if d == 0 {
						id = g.ids[c.y][c.x+1]
						break pick
					}
					d--
				}
			}
			o--
			size--
			cells[id], cells[o] = cells[o], cells[id]
			g.mapp[cells[o].y][cells[o].x] = terrains[terr].symbol
			working[cells[o].y][cells[o].x] = terrains[terr].symbol
			g.ids[cells[o].y][cells[o].x] = o
			g.ids[cells[id].y][cells[id].x] = id
		}
	}

	return islandSize
}

func (g *generator) addIslands(iterations int, seed uint64) {
	var total int64

	for it := 1; it <= iterations; it++ {
		target := islandSizeTab[it]
		for j := 0; j < it; j++ {
			var rng pcg32
			rng.seedStep(seed, uint64(stepIslandBase+it), uint64(j))
			total += int64(g.growOneIsland(target, &rng))
		}
	}
	fmt.Fprintf(g.stderr, "genesis: islands grown, %d land provinces total\n", total)
}

// addOlympus writes Mt. Olympus: a center 'O' ringed by eight 'M'.
func (g *generator) addOlympus() {
	cr := g.gSize / 2
	cc := g.gSize / 2

	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			r := cr + dr
			c := cc + dc
			if r < 0 || r >= g.gSize || c < 0 || c >= g.gSize {
				continue
			}
			g.mapp[r][c] = 'M'
		}
	}
	g.mapp[cr][cc] = 'O'
}

// ---------------------------------------------------------------------------
// Region labeling (iterative floods that mirror mapgen's flood rules).
// ---------------------------------------------------------------------------

// floodWater labels the water region of the same color reachable from
// (sr, sc).
func (g *generator) floodWater(sr, sc, lab int) {
	top := 0
	col0 := oceanColor(g.mapp[sr][sc])

	g.floodStack[top] = sr*g.gSize + sc
	top++
	g.label[sr][sc] = lab

	for top > 0 {
		top--
		cell := g.floodStack[top]
		r := cell / g.gSize
		c := cell % g.gSize
		for d := 0; d < 8; d++ {
			nr, nc, ok := g.neighbor8(r, c, d)
			if !ok {
				continue
			}
			if g.label[nr][nc] != -1 {
				continue
			}
			if !isOceanGlyph(g.mapp[nr][nc]) || oceanColor(g.mapp[nr][nc]) != col0 {
				continue
			}
			g.label[nr][nc] = lab
			g.floodStack[top] = nr*g.gSize + nc
			top++
		}
	}
}

// floodLand labels the whole land mass reachable from (sr, sc).
func (g *generator) floodLand(sr, sc, lab int) {
	top := 0

	g.floodStack[top] = sr*g.gSize + sc
	top++
	g.label[sr][sc] = lab

	for top > 0 {
		top--
		cell := g.floodStack[top]
		r := cell / g.gSize
		c := cell % g.gSize
		for d := 0; d < 8; d++ {
			nr, nc, ok := g.neighbor8(r, c, d)
			if !ok {
				continue
			}
			if g.label[nr][nc] != -1 {
				continue
			}
			if isOceanGlyph(g.mapp[nr][nc]) {
				continue
			}
			g.label[nr][nc] = lab
			g.floodStack[top] = nr*g.gSize + nc
			top++
		}
	}
}

// ---------------------------------------------------------------------------
// Start cities.
// ---------------------------------------------------------------------------

// startCityGlyph is the digit for the i'th placed city: the named safe
// havens 1..7, then 9 and 0 alternately. Never 8, which is not a safe haven.
func startCityGlyph(i int) byte {
	if i < 7 {
		return byte('1' + i)
	}
	if (i-7)%2 == 0 {
		return '9'
	}
	return '0'
}

func isLandTerrain(c byte) bool {
	switch c {
	case 'p', 'f', 'm', 'd', 's', 'o':
		return true
	}
	return false
}

// placeStartCities writes up to n0 start-city digits onto land and returns
// the number placed.
func (g *generator) placeStartCities(n0 int, seed uint64) int {
	var excl [maxDim][maxDim]bool
	wr := make([]int, maxDim*maxDim)
	wc := make([]int, maxDim*maxDim)
	orow := make([]int, maxDim*maxDim)
	ocol := make([]int, maxDim*maxDim)
	wn := 0

	if n0 < 1 {
		return 0
	}

	// Mt. Olympus and every cell touching it.
	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			if g.mapp[r][c] != 'O' && g.mapp[r][c] != 'M' {
				continue
			}
			excl[r][c] = true
			for d := 0; d < 8; d++ {
				if nr, nc, ok := g.neighbor8(r, c, d); ok {
					excl[nr][nc] = true
				}
			}
		}
	}

	// Anchors: write_regions' (the first land cell) and write_land's (which
	// skips 'O').
	for pass := 0; pass < 2; pass++ {
		for r := 0; r < g.gSize; r++ {
			for c := 0; c < g.gSize; c++ {
				g.label[r][c] = -1
			}
		}
		for r := 0; r < g.gSize; r++ {
			for c := 0; c < g.gSize; c++ {
				if g.label[r][c] == -1 && !isOceanGlyph(g.mapp[r][c]) &&
					(pass == 0 || g.mapp[r][c] != 'O') {
					excl[r][c] = true
					g.floodLand(r, c, 0)
				}
			}
		}
	}

	// Prefer plains and forest; fall back to any land terrain.
	for pass := 0; pass < 2 && wn == 0; pass++ {
		for r := 0; r < g.gSize; r++ {
			for c := 0; c < g.gSize; c++ {
				gl := g.mapp[r][c]
				if excl[r][c] {
					continue
				}
				if pass == 0 {
					if gl != 'p' && gl != 'f' {
						continue
					}
				} else if !isLandTerrain(gl) {
					continue
				}
				wr[wn] = r
				wc[wn] = c
				wn++
			}
		}
	}
	if wn < 1 {
		fmt.Fprintf(g.stderr, "genesis: no land for start cities; placed none\n")
		return 0
	}

	r2 := int64(alphaSq * float64(wn) / float64(n0))
	if r2 < 4 {
		r2 = 4 // at least 2 cells apart: start cities never touch
	}

	n := n0
	for ; n >= 1; n-- {
		var rng pcg32
		rng.seedStep(seed, stepStartCities, uint64(n))
		if g.bestCandidates(n, r2, true, &rng, wr, wc, wn, orow, ocol) {
			break
		}
		fmt.Fprintf(g.stderr, "genesis: couldn't space %d start cities; retrying with %d\n", n, n-1)
	}
	g.assert(n >= 1, "place_start_cities", 1343, "n >= 1")

	for i := 0; i < n; i++ {
		g.mapp[orow[i]][ocol[i]] = startCityGlyph(i)
	}
	fmt.Fprintf(g.stderr, "genesis: placed %d start cities (spacing^2 >= %d)\n", n, r2)
	return n
}

// ---------------------------------------------------------------------------
// The region name pool, drawn without replacement.
// ---------------------------------------------------------------------------

func (g *generator) initRegionNames(seed uint64) {
	var rng pcg32

	g.regionPoolN = len(regionNames)
	for i := 0; i < g.regionPoolN; i++ {
		g.regionOrder[i] = i
	}

	// Fisher-Yates shuffle.
	rng.seedStep(seed, stepNames, 0)
	for i := g.regionPoolN - 1; i > 0; i-- {
		j := rng.rangeInt(0, i)
		g.regionOrder[i], g.regionOrder[j] = g.regionOrder[j], g.regionOrder[i]
	}
	g.regionCursor = 0
}

// nextRegionName returns the next pool name, with a numeric suffix once the
// pool wraps.
func (g *generator) nextRegionName() string {
	idx := g.regionOrder[g.regionCursor%g.regionPoolN]
	wrap := g.regionCursor / g.regionPoolN

	g.regionCursor++
	if wrap == 0 {
		return regionNames[idx]
	}
	return fmt.Sprintf("%s %d", regionNames[idx], wrap+1)
}

// ---------------------------------------------------------------------------
// Output writers.
// ---------------------------------------------------------------------------

// create opens an output file the way fopen(name, "w") does, exiting with
// perror's message on failure. Write errors are ignored, as the C code
// ignores them.
func (g *generator) create(name string) (*os.File, *bufio.Writer) {
	fp, err := os.OpenFile(g.path(name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		clib.Perror(g.stderr, "genesis: can't write "+name, err)
		g.exit(1)
	}
	return fp, bufio.NewWriter(fp)
}

func closeOutput(fp *os.File, w *bufio.Writer) {
	w.Flush()
	fp.Close()
}

func (g *generator) writeMap() {
	fp, w := g.create("Map")
	for r := 0; r < g.gSize; r++ {
		w.Write(g.mapp[r][:g.gSize])
		w.WriteByte('\n')
	}
	closeOutput(fp, w)
}

// lakeName recasts a sea name as a lake name: the first "Sea" or "Ocean"
// becomes "Lake", or else "Lake " is prepended.
func lakeName(base string) string {
	if i := strings.Index(base, "Sea"); i >= 0 {
		return base[:i] + "Lake" + base[i+3:]
	}
	if i := strings.Index(base, "Ocean"); i >= 0 {
		return base[:i] + "Lake" + base[i+5:]
	}
	return "Lake " + base
}

// writeWaterRegionsDeferred names the --defer-oceans oceans by size: the
// largest is the Great Sea, and small landlocked bodies become lakes.
func (g *generator) writeWaterRegionsDeferred(w io.Writer) {
	great, maxsize := -1, -1
	id00 := g.oceanID[0][0]

	for i := 0; i < g.deferN; i++ {
		if g.oceanSize[i] > maxsize {
			maxsize = g.oceanSize[i]
		}
	}

	if id00 >= 0 && g.oceanSize[id00] == maxsize {
		great = id00
	} else {
		best := int64(-1)
		for i := 0; i < g.deferN; i++ {
			if g.oceanSize[i] != maxsize {
				continue
			}
			d := dist2(g.oceanOriginR[i], g.oceanOriginC[i], 0, 0)
			if best < 0 || d < best { // ties go to the lowest id
				best = d
				great = i
			}
		}
	}

	fmt.Fprintf(w, "%d,%d\tGreat Sea\n", g.oceanOriginR[great], g.oceanOriginC[great])
	for i := 0; i < g.deferN; i++ {
		if i == great {
			continue
		}
		if g.oceanSize[i] == 0 {
			continue // filled in by fillTinyWater
		}
		nm := g.nextRegionName()
		if g.oceanLandlocked[i] && g.oceanSize[i] < lakeMax {
			nm = lakeName(nm)
		}
		fmt.Fprintf(w, "%d,%d\t%s\n", g.oceanOriginR[i], g.oceanOriginC[i], nm)
	}
}

// writeRegions names every sea and land mass.
func (g *generator) writeRegions() {
	fp, w := g.create("Regions")
	next := 0

	if g.gDeferOceans {
		g.writeWaterRegionsDeferred(w)
	} else {
		// Label all water regions.
		for r := 0; r < g.gSize; r++ {
			for c := 0; c < g.gSize; c++ {
				g.label[r][c] = -1
			}
		}
		for r := 0; r < g.gSize; r++ {
			for c := 0; c < g.gSize; c++ {
				if g.label[r][c] == -1 && isOceanGlyph(g.mapp[r][c]) {
					g.floodWater(r, c, next)
					next++
				}
			}
		}

		greatSeaLabel := g.label[0][0]
		fmt.Fprintf(w, "0,0\tGreat Sea\n")

		// Name every other water region at its first cell in row-major order.
		for lab := 0; lab < next; lab++ {
			fr, fc := -1, -1
			if lab == greatSeaLabel {
				continue
			}
		find:
			for r := 0; r < g.gSize; r++ {
				for c := 0; c < g.gSize; c++ {
					if g.label[r][c] == lab {
						fr, fc = r, c
						break find
					}
				}
			}
			if fr >= 0 {
				fmt.Fprintf(w, "%d,%d\t%s\n", fr, fc, g.nextRegionName())
			}
		}
	}

	// Label and name each land mass at its first cell.
	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			g.label[r][c] = -1
		}
	}
	next = 0
	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			if g.label[r][c] == -1 && !isOceanGlyph(g.mapp[r][c]) {
				fmt.Fprintf(w, "%d,%d\t%s\n", r, c, g.nextRegionName())
				g.floodLand(r, c, next)
				next++
			}
		}
	}

	closeOutput(fp, w)
}

// writeLand names one same-terrain province clump per land mass.
func (g *generator) writeLand() {
	fp, w := g.create("Land")
	next := 0

	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			g.label[r][c] = -1
		}
	}
	for r := 0; r < g.gSize; r++ {
		for c := 0; c < g.gSize; c++ {
			if g.label[r][c] == -1 && !isOceanGlyph(g.mapp[r][c]) && g.mapp[r][c] != 'O' {
				fmt.Fprintf(w, "%d,%d\t%c\t%s\n", r, c, g.mapp[r][c], g.nextRegionName())
				g.floodLand(r, c, next)
				next++
			}
		}
	}

	closeOutput(fp, w)
}

// writeCities writes the city name pool in a shuffled order.
func (g *generator) writeCities(seed uint64) {
	fp, w := g.create("Cities")
	n := len(cityNames)
	order := make([]int, n)
	var rng pcg32

	for i := 0; i < n; i++ {
		order[i] = i
	}
	rng.seedStep(seed, stepCities, 0)
	for i := n - 1; i > 0; i-- {
		j := rng.rangeInt(0, i)
		order[i], order[j] = order[j], order[i]
	}
	for i := 0; i < n; i++ {
		fmt.Fprintf(w, "%s\n", cityNames[order[i]])
	}

	closeOutput(fp, w)
}

// writeRandseed writes 16 bytes for mapgen's MD5 RNG. The C code writes two
// uint64_t values with fwrite, so the bytes are in the host's byte order.
func (g *generator) writeRandseed(seed uint64) {
	fp, w := g.create("randseed")
	a := splitmix64(seed ^ (uint64(stepRandseed) << 32))
	b := splitmix64(a)
	var buf [16]byte
	binary.NativeEndian.PutUint64(buf[0:], a)
	binary.NativeEndian.PutUint64(buf[8:], b)
	w.Write(buf[:])
	closeOutput(fp, w)
}

// ---------------------------------------------------------------------------
// The command.
// ---------------------------------------------------------------------------

func usage(w io.Writer, prog string) {
	fmt.Fprintf(w,
		"Usage: %s [--seed N] [--size 10..99] [--number-of-oceans N]\n"+
			"          [--island-iterations 2..10] [--defer-oceans]\n"+
			"          [--ocean-growth frontier|bfs]\n"+
			"          [--start-cities 0..16]\n"+
			"          [--image FILE.bmp] [--pixels 1..64] [--grid 0..4]\n"+
			"  --number-of-oceans is capped to min(ceil(size/3), 12) and clamped to\n"+
			"  fit; larger requests are reduced with a warning.\n"+
			"  --defer-oceans generates the oceans AFTER the islands and Mt. Olympus\n"+
			"  (origin-seeded, size-named; largest sea = Great Sea).  See issue #80.\n"+
			"  --ocean-growth picks how --defer-oceans grows each ocean from its\n"+
			"  origin: frontier (default; random frontier flood, ragged seas) or bfs\n"+
			"  (ring-by-ring, nearest origin through the water).  See issue #87.\n"+
			"  --start-cities places up to N safe-haven start cities on land (default\n"+
			"  6; fewer if they can't be spaced; 0 = none).  Without at least one,\n"+
			"  olympia-g3 can't add players.  See issue #85.\n"+
			"  --image dumps the generated Map as a 24-bit BMP for eyeballing the\n"+
			"  world (debug aid, issue #81; emitted only when given, never affects\n"+
			"  the Map/Regions/Land/Cities output).  --pixels sets N x N pixels per\n"+
			"  province (default 8); --grid sets black grid-line thickness in pixels\n"+
			"  between and around provinces (default 0 = no grid).\n"+
			"Generates the mapgen input set (Map, Regions, Land, Cities, randseed)\n"+
			"in the current directory.  See doc/genesis.md.\n", prog)
}

const (
	optImage = 256 + iota
	optPixels
	optGrid
	optStartCities
	optOceanGrowth
)

var longopts = []clib.LongOption{
	{Name: "size", HasArg: clib.RequiredArgument, Val: 's'},
	{Name: "number-of-oceans", HasArg: clib.RequiredArgument, Val: 'n'},
	{Name: "island-iterations", HasArg: clib.RequiredArgument, Val: 'i'},
	{Name: "seed", HasArg: clib.RequiredArgument, Val: 'd'},
	{Name: "defer-oceans", HasArg: clib.NoArgument, Val: 'D'},
	{Name: "start-cities", HasArg: clib.RequiredArgument, Val: optStartCities},
	{Name: "ocean-growth", HasArg: clib.RequiredArgument, Val: optOceanGrowth},
	{Name: "image", HasArg: clib.RequiredArgument, Val: optImage},
	{Name: "pixels", HasArg: clib.RequiredArgument, Val: optPixels},
	{Name: "grid", HasArg: clib.RequiredArgument, Val: optGrid},
	{Name: "help", HasArg: clib.NoArgument, Val: 'h'},
}

// Run generates a world from the command-line arguments args (args[0] is
// the program name), writing the output files into dir. It returns the
// process exit status genesis-g3 would have.
func Run(dir string, args []string, stderr io.Writer) (status int) {
	g := &generator{dir: dir, stderr: stderr}

	defer func() {
		switch r := recover().(type) {
		case nil:
		case exitStatus:
			status = int(r)
		case assertion:
			// assert() prints the compiler's __FILE__, which is the absolute
			// path of the upstream source; the relative path stands in for it.
			fmt.Fprintf(stderr, "Assertion failed: (%s), function %s, file %s, line %d.\n",
				r.cond, r.function, "genesis/genesis.c", r.line)
			// abort() exits 134 and, on macOS, flushes every open stdio
			// stream. genesis.c never writes stdout and opens no output
			// file before its last assert(), so there is nothing to flush.
			status = 134
		default:
			panic(r)
		}
	}()

	return g.main(args)
}

func (g *generator) main(args []string) int {
	size := 99
	oceans := 9
	iterations := 10
	startCities := 6
	var seed uint64
	imagePath, haveImage := "", false
	imagePixels := 8
	imageGrid := 0

	prog := ""
	if len(args) > 0 {
		prog = args[0]
	}

	opts := clib.NewGetoptLong(args, "s:n:i:d:Dh", longopts, g.stderr)
	// gen_warn uses a fixed prefix; usage still prints the actual argv[0].
	opts.ProgramName = "genesis-g3"
	for {
		opt := opts.Next()
		if opt == -1 {
			break
		}
		optarg := opts.Optarg
		switch opt {
		case 's':
			size = clib.Atoi(optarg)
		case 'n':
			oceans = clib.Atoi(optarg)
		case 'i':
			iterations = clib.Atoi(optarg)
		case 'd':
			seed = clib.Strtoull0(optarg)
		case 'D':
			g.gDeferOceans = true
		case optImage:
			imagePath, haveImage = optarg, true
		case optPixels:
			imagePixels = clib.Atoi(optarg)
		case optGrid:
			imageGrid = clib.Atoi(optarg)
		case optStartCities:
			startCities = clib.Atoi(optarg)
		case optOceanGrowth:
			switch optarg {
			case "bfs":
				g.gOceanBFS = true
			case "frontier":
				g.gOceanBFS = false
			default:
				fmt.Fprintf(g.stderr, "genesis: --ocean-growth must be frontier or bfs (got %s)\n", optarg)
				return 1
			}
		case 'h':
			usage(g.stderr, prog)
			return 0
		default:
			usage(g.stderr, prog)
			return 1
		}
	}

	if size < 10 || size > 99 {
		fmt.Fprintf(g.stderr, "genesis: --size must be 10..99 (got %d)\n", size)
		return 1
	}
	if oceans < 1 {
		fmt.Fprintf(g.stderr, "genesis: --number-of-oceans must be >= 1 (got %d)\n", oceans)
		return 1
	}
	if iterations < 2 || iterations > 10 {
		fmt.Fprintf(g.stderr, "genesis: --island-iterations must be 2..10 (got %d)\n", iterations)
		return 1
	}
	if g.gOceanBFS && !g.gDeferOceans {
		fmt.Fprintf(g.stderr, "genesis: --ocean-growth bfs requires --defer-oceans\n")
		return 1
	}
	if startCities < 0 || startCities > 16 {
		fmt.Fprintf(g.stderr, "genesis: --start-cities must be 0..16 (got %d)\n", startCities)
		return 1
	}
	if haveImage && (imagePixels < 1 || imagePixels > 64) {
		fmt.Fprintf(g.stderr, "genesis: --pixels must be 1..64 (got %d)\n", imagePixels)
		return 1
	}
	if haveImage && (imageGrid < 0 || imageGrid > 4) {
		fmt.Fprintf(g.stderr, "genesis: --grid must be 0..4 (got %d)\n", imageGrid)
		return 1
	}

	g.gSize = size
	mode := ""
	if g.gDeferOceans {
		mode = " (defer-oceans)"
		if g.gOceanBFS {
			mode = " (defer-oceans, bfs)"
		}
	}
	fmt.Fprintf(g.stderr, "genesis: seed=%d size=%d oceans=%d iterations=%d%s\n",
		seed, size, oceans, iterations, mode)

	// Each step reseeds its own stream from the user seed.
	if g.gDeferOceans {
		// Islands need water to grow into, so fill the canvas first, grow
		// the world, then carve the oceans out of the remaining water.
		g.fillOceanCanvas()
		g.addIslands(iterations, seed)
		g.addOlympus()
		g.deferOceans(oceans, seed)
	} else {
		g.addOceans(oceans, seed)
		g.addIslands(iterations, seed)
		g.addOlympus()
	}
	// The last terrain step: the land and ocean layout is final.
	g.placeStartCities(startCities, seed)

	g.initRegionNames(seed)
	g.writeMap()
	g.writeRegions()
	g.writeLand()
	g.writeCities(seed)
	g.writeRandseed(seed)

	fmt.Fprintf(g.stderr, "genesis: wrote Map, Regions, Land, Cities, randseed\n")

	// --image is a debugging aid that runs after every output step.
	if haveImage {
		grid := make([]byte, 0, g.gSize*g.gSize)
		for r := 0; r < g.gSize; r++ {
			grid = append(grid, g.mapp[r][:g.gSize]...)
		}
		// fopen("") fails, but the empty path would resolve to dir.
		if imagePath == "" || !writeBMP(g.path(imagePath), grid, g.gSize, imagePixels, imageGrid) {
			fmt.Fprintf(g.stderr, "genesis: failed to write image '%s'\n", imagePath)
			return 1
		}
		px := g.gSize*imagePixels + (g.gSize+1)*imageGrid
		fmt.Fprintf(g.stderr, "genesis: wrote image %s (%dx%d px)\n", imagePath, px, px)
	}

	return 0
}
