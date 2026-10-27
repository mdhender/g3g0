// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package island ports G3's island-g3, which adds one randomly generated
// island to a map read from stdin and prints the new map.
//
// The port follows island/island.c and the RNG in lib/rnd.c. It keeps the C
// program's draw order, scan order and integer semantics, so that the same
// map, flags and randseed give the same output byte for byte.
package island

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strings"

	"github.com/mdhender/g3g0/internal/clib"
	"github.com/mdhender/g3g0/internal/rnd"
)

const (
	maxMap      = 120
	lineMax     = 200
	distanceCap = 9
)

type location struct {
	x, y int
}

type terrain struct {
	symbol     byte
	min, max   int
	targetProb int
	prob       int
}

// The terrain glyphs are the map's shared vocabulary with mapgen; see the
// note above terrains[] in island/island.c.
var terrains = [...]terrain{
	{symbol: 'p', min: 12, max: 30, targetProb: 30},
	{symbol: 'f', min: 6, max: 14, targetProb: 30},
	{symbol: 'm', min: 6, max: 10, targetProb: 20},
	{symbol: 'd', min: 15, max: 30, targetProb: 10},
	{symbol: 's', min: 1, max: 3, targetProb: 10},
}

type generator struct {
	m        [maxMap][maxMap]byte // map[][] from main()
	working  [maxMap][maxMap]byte
	distance [maxMap][maxMap]int
	ids      [maxMap][maxMap]int

	xSize, ySize int
}

// Run adds an island to the map read from stdin and prints the result to
// stdout, reading and advancing randseed in dir. args is argv, including
// argv[0]. It returns the process exit status island-g3 would have.
func Run(dir string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	g := &generator{}
	terrains := terrains // main() changes the probabilities in place

	targetSize := 0
	border := 2
	shelf := 3

	opts := clib.NewGetopt(args, "b:c:s:", stderr)
	opts.Opterr = false
	for {
		opt := opts.Next()
		if opt == -1 {
			break
		}
		switch opt {
		case 'b':
			border = clib.Atoi(opts.Optarg)
		case 'c':
			shelf = clib.Atoi(opts.Optarg)
		case 's':
			targetSize = clib.Atoi(opts.Optarg)
		default:
			fmt.Fprintf(stderr, "Usage: island-g3 -b border_width -c continental_shelf_width -s size\n")
			return 1
		}
	}

	if border < 0 {
		border = 0
	}

	if shelf < 0 {
		shelf = 0
	}

	// Read the starting map from stdin. A read error ends the input, as
	// fgets() returning NULL does.
	data, _ := io.ReadAll(stdin)
	for _, buffer := range clib.FgetsLines(data, lineMax) {
		n := strings.IndexAny(buffer, "\r\n") // strcspn(buffer, "\r\n")
		if n < 0 {
			n = len(buffer)
		}
		if n > maxMap-1 {
			fmt.Fprintf(stderr, "Map row longer than %d characters!\n", maxMap-1)
			return 1
		}
		if g.ySize == maxMap && n > 0 {
			fmt.Fprintf(stderr, "Map taller than %d rows!\n", maxMap)
			return 1
		}
		// sscanf(buffer, "%[^\r\n]", map[y_size]) stores the row and
		// returns 1 unless the line starts with '\r', '\n' or a NUL. Only a
		// match counts (G3 #303).
		// A row of a different length is skipped.
		if n > 0 {
			copy(g.m[g.ySize][:], buffer[:n])
			if g.xSize == 0 {
				g.xSize = n
			}
			if n == g.xSize {
				g.ySize++
			}
		}
	}

	// No rows: nothing to place an island in (G3 #295)
	if g.ySize == 0 {
		fmt.Fprintf(stderr, "No map on stdin!\n")
		return 1
	}

	// island checks that randseed opens before calling load_seed(), so a
	// missing seed is reported on stderr, not in the map (G3 #298).
	rng := rnd.New([16]byte{})
	seedPath := filepath.Join(dir, "randseed")
	if err := rng.Load(seedPath); err != nil {
		fmt.Fprintf(stderr, "randseed could not be opened.\n")
	}

	xSize, ySize := g.xSize, g.ySize
	m, working, distance, ids := &g.m, &g.working, &g.distance, &g.ids

	// Divide the map into available (water) and unavailable (land) space.
	// Water is any glyph mapgen treats as ocean.
	for y := range ySize {
		for x := range xSize {
			switch m[y][x] {
			case ';', ',', ':', '.', '~', ' ', '"', '\'':
				working[y][x] = '~'
			default:
				working[y][x] = 'p'
			}
			ids[y][x] = 0
		}
	}
	// exclude any square within shelf of a land square
	g.makeShelf(shelf)
	// exclude border squares around each edge
	// A border wider than the map stops at its edge (G3 #302)
	for y := range ySize {
		for x := 0; x < border && x < xSize; x++ {
			if working[y][x] == '~' {
				working[y][x] = '_'
			}
			if working[y][xSize-x-1] == '~' {
				working[y][xSize-x-1] = '_'
			}
		}
	}
	for x := range xSize {
		for y := 0; y < border && y < ySize; y++ {
			if working[y][x] == '~' {
				working[y][x] = '_'
			}
			if working[ySize-y-1][x] == '~' {
				working[ySize-y-1][x] = '_'
			}
		}
	}
	// Calculate distance from land
	for y := range ySize {
		for x := range xSize {
			if working[y][x] == '~' {
				distance[y][x] = distanceCap
			} else {
				distance[y][x] = 0
			}
		}
	}
	for y := range ySize {
		for x := range xSize {
			g.extendDistance(y, x)
		}
	}

	// No water left outside the shelf and border: every square has distance
	// 0, and the start-square draw below would pick from all of them, land
	// included. Stop before any draw, so randseed is not rewritten (G3 #305).
	maxDist := 0
	for y := range ySize {
		for x := range xSize {
			maxDist = max(maxDist, distance[y][x])
		}
	}
	if maxDist == 0 {
		fmt.Fprintf(stderr, "No room for an island!\n")
		return 1
	}

	if targetSize < 1 {
		// Decide how big to make the island
		// Count from 0: any size below 1 means random (G3 #294)
		targetSize = 0
		for range 2 {
			for rng.Rnd(0, 99) != 0 {
				targetSize++
			}
		}
	}
	// malloc((size_t)(target_size + 1) * sizeof(*island)) fails only when
	// target_size + 1 overflows; the macOS allocator grants anything smaller
	// (up to 16 GiB) without touching it. The island never has more squares
	// than the map, so that is all the port allocates.
	if targetSize == math.MaxInt32 {
		fmt.Fprintf(stderr, "Unable to allocate memory for island data!\n")
		return 1
	}
	island := make([]location, min(targetSize, xSize*ySize)+1)

	// select a square to start the new island
	maxDist = 0
	count := 0
	islandSize := 0
	for y := range ySize {
		for x := range xSize {
			if distance[y][x] > maxDist {
				maxDist = distance[y][x]
				count = 1
			} else if distance[y][x] == maxDist {
				count++
			}
		}
	}
	d := rng.Rnd(1, count)
	for y := 0; d > 0 && y < ySize; y++ {
		for x := 0; d > 0 && x < xSize; x++ {
			if distance[y][x] == maxDist {
				d--
				if d == 0 {
					island[0].x = x
					island[0].y = y
					islandSize = 1
					working[y][x] = 'o'
					m[y][x] = 'o'
					ids[y][x] = 0
				}
			}
		}
	}

	// next is the C idiom !d--: true when d reaches zero, and d counts
	// down either way.
	next := func() bool {
		hit := d == 0
		d--
		return hit
	}

	for islandSize < targetSize {
		// Find all the possible expansions
		// Count multiple access paths multiple times to
		// make the island more likely to fill interior space
		count = 0
		for i := range islandSize {
			p := island[i]
			if p.x > 0 && working[p.y][p.x-1] == '~' {
				count++
			}
			if p.x < xSize-1 && working[p.y][p.x+1] == '~' {
				count++
			}
			if p.y > 0 && working[p.y-1][p.x] == '~' {
				count++
			}
			if p.y < ySize-1 && working[p.y+1][p.x] == '~' {
				count++
			}
		}
		if count < 1 {
			fmt.Fprintf(stderr, "Not enough room to expand island!\n")
			break
		}
		d = rng.Rnd(0, count-1)
		for i := range islandSize {
			p := island[i]
			if p.x > 0 && working[p.y][p.x-1] == '~' {
				if next() {
					island[islandSize] = location{x: p.x - 1, y: p.y}
				}
			}
			if p.x < xSize-1 && working[p.y][p.x+1] == '~' {
				if next() {
					island[islandSize] = location{x: p.x + 1, y: p.y}
				}
			}
			if p.y > 0 && working[p.y-1][p.x] == '~' {
				if next() {
					island[islandSize] = location{x: p.x, y: p.y - 1}
				}
			}
			if p.y < ySize-1 && working[p.y+1][p.x] == '~' {
				if next() {
					island[islandSize] = location{x: p.x, y: p.y + 1}
				}
			}
		}
		p := island[islandSize]
		working[p.y][p.x] = 'o'
		m[p.y][p.x] = 'o'
		ids[p.y][p.x] = islandSize
		islandSize++
	}

	// Calculate terrain ratios
	// Use GCD and LCM to render these down to integer ratios
	for i := range terrains {
		terrains[i].prob = terrains[i].min + terrains[i].max
		div := gcd(terrains[i].targetProb, terrains[i].prob)
		terrains[i].targetProb /= div
		terrains[i].prob /= div
	}
	mult := 1
	for i := range terrains {
		mult = lcm(mult, terrains[i].prob)
	}
	count = 0
	for i := range terrains {
		terrains[i].targetProb *= mult / terrains[i].prob
		count += terrains[i].targetProb
	}

	// swap exchanges island[a] and island[b] and renumbers both squares.
	swap := func(a, b int) {
		island[a], island[b] = island[b], island[a]
		ids[island[b].y][island[b].x] = b
		ids[island[a].y][island[a].x] = a
	}

	o := islandSize
	size := 0
	terr := 0
	id := 0
	for o > 0 {
		// pick a random square on the island
		clusterEnd := o
		o--
		d = rng.Rnd(0, o)
		// Move that square to the end of the list
		swap(d, o)
		if size < 1 {
			// Pick a random terrain type
			d = rng.Rnd(1, count)
			for terr = 0; d > terrains[terr].targetProb; terr++ {
				d -= terrains[terr].targetProb
			}
			// Determine how big this cluster should be
			size = rng.Rnd(terrains[terr].min, terrains[terr].max)
		}
		// set the first square to this terrain type
		p := island[o]
		m[p.y][p.x] = terrains[terr].symbol
		working[p.y][p.x] = terrains[terr].symbol
		size--
		if size > o {
			size = o
		}
		for size > 0 {
			opt := 0
			for i := o; i < clusterEnd; i++ {
				p := island[i]
				if p.y > 0 && working[p.y-1][p.x] == 'o' {
					opt++
				}
				if p.y < ySize-1 && working[p.y+1][p.x] == 'o' {
					opt++
				}
				if p.x > 0 && working[p.y][p.x-1] == 'o' {
					opt++
				}
				if p.x < xSize-1 && working[p.y][p.x+1] == 'o' {
					opt++
				}
			}
			if opt < 1 {
				break
			}
			d = rng.Rnd(0, opt-1)
			for i := o; i < clusterEnd; i++ {
				p := island[i]
				if p.y > 0 && working[p.y-1][p.x] == 'o' {
					if next() {
						id = ids[p.y-1][p.x]
						break
					}
				}
				if p.y < ySize-1 && working[p.y+1][p.x] == 'o' {
					if next() {
						id = ids[p.y+1][p.x]
						break
					}
				}
				if p.x > 0 && working[p.y][p.x-1] == 'o' {
					if next() {
						id = ids[p.y][p.x-1]
						break
					}
				}
				if p.x < xSize-1 && working[p.y][p.x+1] == 'o' {
					if next() {
						id = ids[p.y][p.x+1]
						break
					}
				}
			}
			o--
			size--
			swap(id, o)
			p := island[o]
			m[p.y][p.x] = terrains[terr].symbol
			working[p.y][p.x] = terrains[terr].symbol
		}
	}

	if err := rng.Save(seedPath); err != nil {
		fmt.Fprintf(stderr, "couldn't write %s:", "randseed")
		clib.Perror(stderr, "", err)
		return 1
	}

	fmt.Fprintf(stderr, "Added island of %d provinces.\n", islandSize)

	// Print the map with the added island
	w := bufio.NewWriter(stdout)
	for y := range ySize {
		w.Write(m[y][:xSize])
		w.WriteByte('\n')
	}
	_ = w.Flush()

	return 0
}

// makeShelf marks water within shelf squares of a land square as '_', the
// continental-shelf exclusion marker. '_' is only ever written to working,
// never to the printed map.
//
// "Within shelf squares" counts north/south/east/west steps through any
// square on the map. One breadth-first pass from every land square marks the
// same squares the old per-square recursion did, which took 4^shelf calls per
// land square and shelf stack frames (G3 #307).
func (g *generator) makeShelf(shelf int) {
	var steps [maxMap][maxMap]int
	queue := make([]location, 0, maxMap*maxMap)
	dy := [4]int{-1, 1, 0, 0}
	dx := [4]int{0, 0, -1, 1}

	for y := range g.ySize {
		for x := range g.xSize {
			if g.working[y][x] == 'p' {
				steps[y][x] = 0
				queue = append(queue, location{x: x, y: y})
			} else {
				steps[y][x] = -1
			}
		}
	}

	for head := 0; head < len(queue); head++ {
		p := queue[head]
		if steps[p.y][p.x] >= shelf {
			continue
		}
		for k := range 4 {
			ny, nx := p.y+dy[k], p.x+dx[k]
			if ny < 0 || ny >= g.ySize || nx < 0 || nx >= g.xSize || steps[ny][nx] >= 0 {
				continue
			}
			steps[ny][nx] = steps[p.y][p.x] + 1
			if g.working[ny][nx] == '~' {
				g.working[ny][nx] = '_'
			}
			queue = append(queue, location{x: nx, y: ny})
		}
	}
}

func (g *generator) extendDistance(y, x int) {
	distance := &g.distance
	if y > 0 && distance[y-1][x] > distance[y][x]+1 {
		distance[y-1][x] = distance[y][x] + 1
		g.extendDistance(y-1, x)
	}
	if y < g.ySize-1 && distance[y+1][x] > distance[y][x]+1 {
		distance[y+1][x] = distance[y][x] + 1
		g.extendDistance(y+1, x)
	}
	if x > 0 && distance[y][x-1] > distance[y][x]+1 {
		distance[y][x-1] = distance[y][x] + 1
		g.extendDistance(y, x-1)
	}
	if x < g.xSize-1 && distance[y][x+1] > distance[y][x]+1 {
		distance[y][x+1] = distance[y][x] + 1
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
