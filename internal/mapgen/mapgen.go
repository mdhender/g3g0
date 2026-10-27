// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package mapgen ports the Olympia G3 map generator (mapgen/mapgen.c) from C.
//
// The port is a compatibility port: given the same Map, Regions, Land,
// Cities and randseed files it writes byte-identical loc, gate, road and
// randseed files, and the same diagnostics, as G3's mapgen-g3. Function and
// variable names follow the C source so the two can be read side by side.
// Defects in the C code are reproduced on purpose; see the comments that link
// to the upstream issues.
package mapgen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mdhender/g3g0/internal/clib"
	"github.com/mdhender/g3g0/internal/rnd"
)

const (
	lineLen = 256

	maxRow = 100
	maxCol = 100

	regionOff = 58760 // where to start numbering regions
)

const (
	terrLand        = 1
	terrOcean       = 2
	terrForest      = 3
	terrSwamp       = 4
	terrMountain    = 5
	terrPlain       = 6
	terrDesert      = 7
	terrWater       = 8
	terrIsland      = 9
	terrStoneCir    = 10 // circle of stones
	terrGrove       = 11 // mallorn grove
	terrBog         = 12
	terrCave        = 13
	terrCity        = 14
	terrGuild       = 15
	terrGrave       = 16
	terrRuins       = 17
	terrBattlefield = 18
	terrEnchFor     = 19 // enchanted forest
	terrRockyHill   = 20
	terrTreeCir     = 21
	terrPits        = 22
	terrPasture     = 23
	terrOasis       = 24
	terrYewGrove    = 25
	terrSandPit     = 26
	terrSacGrove    = 27 // sacred grove
	terrPopField    = 28 // poppy field
	terrTemple      = 29
	terrLair        = 30 // dragon lair
)

const (
	dirN  = 1
	dirE  = 2
	dirS  = 3
	dirW  = 4
	dirNE = 5
	dirSE = 6
	dirSW = 7
	dirNW = 8

	maxDir = 9
)

const maxBox = 100000

// G2 entity coding system; see the table in mapgen.c.
const (
	sublocLow  = 59000
	sublocHigh = 78999

	cityLow  = 56760
	cityHigh = 58759
)

const (
	maxInside = 200 // max continents/regions
	maxSubloc = 20000
)

var terrS = []string{
	"<null>",
	"land",
	"ocean",
	"forest",
	"swamp",
	"mountain",
	"plain",
	"desert",
	"water",
	"island",
	"ring of stones",
	"mallorn grove",
	"bog",
	"cave",
	"city",
	"guild",
	"graveyard",
	"ruins",
	"battlefield",
	"enchanted forest",
	"rocky hill",
	"circle of trees",
	"pits",
	"pasture",
	"oasis",
	"yew grove",
	"sand pit",
	"sacred grove",
	"poppy field",
	"temple",
	"lair",
}

// Names are *string because the C code compares them by pointer as well as
// by value, and NULL is a distinct state. These are the string literals.
var (
	nameUldimPass    = strp("Uldim pass")
	nameSummerbridge = strp("Summerbridge")
	nameMtOlympus    = strp("Mt. Olympus")
)

// strp is str_save(): a fresh copy with its own identity.
func strp(s string) *string {
	return &s
}

type road struct {
	entNum int
	name   *string
	toLoc  int
	hidden int
}

type tile struct {
	saveChar byte
	region   int
	name     *string
	terrain  int
	hidden   int
	city     int
	mark     int
	inside   int
	color    int // map coloring for
	row, col int // map tile we're inside
	depth    int

	safeHaven        int
	seaLane          int
	uldimFlag        int
	summerbridgeFlag int
	regionBoundary   int
	majorCity        int

	subs []int

	gatesDest []int // gates from here
	gatesNum  []int // gates from here
	gatesKey  []int

	roads []*road
}

// generator holds what mapgen.c keeps in globals and function statics.
type generator struct {
	dir    string
	stdout io.Writer
	stderr io.Writer
	rng    *rnd.Source

	allocFlag [maxBox]bool

	dirVector [maxDir]int

	insideNames     [maxInside]*string
	insideList      [maxInside][]*tile // list of provinces in each region
	insideGatesTo   [maxInside]int     // for info gathering only
	insideGatesFrom [maxInside]int     // for info gathering only
	insideNumCities [maxInside]int     // for info gathering only
	insideTop       int

	maxCol int
	maxRow int

	waterCount int // count of water provinces
	landCount  int // count of land provinces
	numIslands int

	m [maxRow][maxCol]*tile

	subloc    [maxSubloc]*tile
	topSubloc int

	files  []*os.File
	locFp  *bufio.Writer
	gateFp *bufio.Writer
	roadFp *bufio.Writer

	// statics
	holeRoadNameCnt   int
	cornerRoadNameCnt int
	cities            *lineReader
}

// exitStatus unwinds the generator for exit(3).
type exitStatus int

// assertion unwinds the generator for asfail(), which aborts.
type assertion struct {
	line int
	cond string
}

// Run generates a map from the input files in dir, writing the output files
// there. It returns the process exit status mapgen-g3 would have.
func Run(dir string, stdout, stderr io.Writer) (status int) {
	g := &generator{dir: dir, stdout: stdout, stderr: stderr, rng: rnd.New([16]byte{})}

	defer func() {
		switch r := recover().(type) {
		case nil:
		case exitStatus:
			status = int(r)
		case assertion:
			// asfail() prints the compiler's __FILE__, which is the absolute
			// path of the upstream source; the relative path stands in for it.
			fmt.Fprintf(g.stderr, "assertion failure: %s (%d): %s\n", "mapgen/mapgen.c", r.line, r.cond)
			status = 134 // abort()
		case runtime.Error:
			// The C code dereferences NULL on several inputs and dies with
			// SIGSEGV, printing nothing. The port reaches the same points as nil
			// dereferences; report them the same way rather than as a Go panic.
			if !strings.Contains(r.Error(), "nil pointer dereference") {
				panic(r)
			}
			status = 139 // 128 + SIGSEGV
		default:
			panic(r)
		}
	}()

	g.main()
	return 0
}

func (g *generator) path(name string) string {
	return filepath.Join(g.dir, name)
}

func (g *generator) errf(format string, args ...any) {
	fmt.Fprintf(g.stderr, format, args...)
}

func (g *generator) exit(status int) {
	panic(exitStatus(status))
}

func (g *generator) assert(ok bool, line int, cond string) {
	if !ok {
		panic(assertion{line: line, cond: cond})
	}
}

func (g *generator) rnd(low, high int) int {
	return g.rng.Rnd(low, high)
}

func (g *generator) main() {
	g.clearAllocFlag()
	g.dirAssert()
	g.openFps()
	g.loadSeed("randseed")
	g.mapInit()
	g.readMap()
	g.fixTerrainLand()
	g.setRegions()
	g.setProvinceClumps()
	g.unnamedProvinceClumps()
	g.makeIslands()
	g.makeGraveyards()
	g.placeSublocations()
	g.makeGates()
	g.makeRoads()
	g.printMap()
	g.printSublocs()
	g.dumpContinents()
	g.countCities()
	g.countContinents()
	g.countSublocs()
	g.countSublocCoverage()
	g.dumpRoads()
	g.dumpGates()

	g.closeFps()

	g.countTiles()

	g.errf("\n")
	g.errf("highest province = %d\n", g.rcToRegion(g.maxRow, g.maxCol))
	g.errf("\n")

	g.saveSeed("randseed")
}

// openFps creates the output files. They are buffered like stdio streams and
// written when closed.
func (g *generator) openFps() {
	open := func(name string) *bufio.Writer {
		fp, err := os.Create(g.path(name))
		if err != nil {
			clib.Perror(g.stderr, "can't write "+name, err)
			g.exit(1)
		}
		g.files = append(g.files, fp)
		return bufio.NewWriter(fp)
	}
	g.locFp = open("loc")
	g.gateFp = open("gate")
	g.roadFp = open("road")
}

func (g *generator) closeFps() {
	for _, w := range []*bufio.Writer{g.locFp, g.gateFp, g.roadFp} {
		_ = w.Flush()
	}
	for _, fp := range g.files {
		_ = fp.Close()
	}
	g.files = nil
}

// loadSeed is load_seed() from lib/rnd.c.
func (g *generator) loadSeed(fnam string) {
	if err := g.rng.Load(g.path(fnam)); err != nil {
		fmt.Fprintf(g.stdout, "%s could not be opened.\n", fnam)
	}
}

// saveSeed is save_seed() from lib/rnd.c.
func (g *generator) saveSeed(fnam string) {
	if err := g.rng.Save(g.path(fnam)); err != nil {
		g.errf("couldn't write %s:", fnam)
		clib.Perror(g.stderr, "", err)
		g.exit(1)
	}
}

func (g *generator) clearAllocFlag() {
	for i := range g.allocFlag {
		g.allocFlag[i] = false
	}
}

func (g *generator) rndAllocNum(low, high int) int {
	n := g.rnd(low, high)

	for i := n; i <= high; i++ {
		if !g.allocFlag[i] {
			g.allocFlag[i] = true
			return i
		}
	}

	for i := low; i < n; i++ {
		if !g.allocFlag[i] {
			g.allocFlag[i] = true
			return i
		}
	}

	g.errf("rnd_alloc_num(%d,%d) failed\n", low, high)

	return -1
}

func (g *generator) mapInit() {
	for i := range g.m {
		for j := range g.m[i] {
			g.m[i][j] = nil
		}
	}
}

func (g *generator) readMap() {
	data, err := clib.ReadFile(g.path("Map"))
	if err != nil {
		clib.Perror(g.stderr, "can't open Map", err)
		g.exit(1)
	}

	var terrain int
	row := 0

	for _, buf := range clib.FgetsLines(data, lineLen) {
		for col := 0; col < len(buf) && buf[col] != '\n'; col++ {
			if buf[col] == '#' { // hole in map
				continue
			}

			if row > g.maxRow {
				g.maxRow = row
			}

			if col > g.maxCol {
				g.maxCol = col
			}

			t := &tile{}
			g.m[row][col] = t
			t.row = row
			t.col = col
			t.region = g.rcToRegion(row, col)
			t.depth = 2

			color := 0

			switch buf[col] {
			case ';', ',':
				if buf[col] == ';' {
					t.seaLane = 1
				}
				terrain = terrOcean
				color = 1

			case ':', '.':
				if buf[col] == ':' {
					t.seaLane = 1
				}
				terrain = terrOcean
				color = 2

			case '~', ' ':
				if buf[col] == '~' {
					t.seaLane = 1
				}
				terrain = terrOcean
				color = 3

			case '"', '\'':
				if buf[col] == '"' {
					t.seaLane = 1
				}
				terrain = terrOcean
				color = 4

			case 'p':
				color = 5
				terrain = terrPlain

			case 'P':
				color = 6
				terrain = terrPlain

			case 'd':
				color = 7
				terrain = terrDesert

			case 'D':
				color = 8
				terrain = terrDesert

			case 'm':
				color = 9
				terrain = terrMountain

			case 'M':
				color = 10
				terrain = terrMountain

			case 's':
				color = 11
				terrain = terrSwamp

			case 'S':
				color = 12
				terrain = terrSwamp

			case 'f':
				color = 13
				terrain = terrForest

			case 'F':
				color = 14
				terrain = terrForest

			case 'o':
				switch g.rnd(1, 10) {
				case 1, 2, 3:
					terrain = terrForest
				case 4, 5, 6:
					terrain = terrPlain
				case 7, 8:
					terrain = terrMountain
				case 9:
					terrain = terrSwamp
				case 10:
					terrain = terrDesert
				}
				color = -1

			// Special stuff

			case '^':
				color = 9 // was 15, unique
				terrain = terrMountain
				t.uldimFlag = 1
				t.regionBoundary = 1

			case 'v':
				color = 9 // was 15, unique
				terrain = terrMountain
				t.uldimFlag = 2
				t.regionBoundary = 1

			case '{':
				color = 16
				terrain = terrMountain
				t.uldimFlag = 3
				t.name = nameUldimPass
				t.regionBoundary = 1

			case '}':
				color = 16
				terrain = terrMountain
				t.uldimFlag = 4
				t.name = nameUldimPass
				t.regionBoundary = 1

			case ']':
				terrain = terrSwamp
				t.summerbridgeFlag = 1
				t.name = nameSummerbridge
				t.regionBoundary = 1

			case '[':
				terrain = terrSwamp
				t.summerbridgeFlag = 2
				t.name = nameSummerbridge
				t.regionBoundary = 1

			case 'O':
				terrain = terrMountain
				color = -1
				t.name = nameMtOlympus

			case '1', '2', '3', '4', '5', '6', '7', '9', '0':
				terrain = terrForest
				color = 19
				t.safeHaven = 1

				var name *string
				label := "Start city"
				switch buf[col] {
				case '1':
					name = strp("Drassa")
				case '2':
					name = strp("Rimmon")
				case '3':
					name = strp("Harn")
				case '4':
					name = strp("Imperial City")
					label = "Imperical City"
				case '5':
					name = strp("Port Aurnos")
				case '6':
					name = strp("Greyfell")
				case '7':
					name = strp("Yellowleaf")
				}

				n := g.createACity(row, col, name, 1)
				g.subloc[n].safeHaven = 1

				g.errf("%s #%s %s at (%d,%d)\n", label, clib.Ch(buf[col]), clib.Str(g.subloc[n].name), row, col)

			case '8':
				terrain = terrForest
				color = 19

				n := g.createACity(row, col, strp("Golden City"), 1)

				g.errf("Golden City #%s %s at (%d,%d)\n", clib.Ch(buf[col]), clib.Str(g.subloc[n].name), row, col)

			case '*':
				terrain = terrLand
				g.createACity(row, col, nil, 1)

			case '%':
				terrain = terrLand
				g.createACity(row, col, nil, 0)

			default:
				// Every digit has a case above, so this is never a digit.
				g.errf("unknown terrain %s\n", clib.Ch(buf[col]))
				g.assert(false, 680, "FALSE")
			}

			t.saveChar = buf[col]
			t.terrain = terrain
			t.color = color

			if terrain == terrWater || terrain == terrOcean {
				g.waterCount++
			} else {
				g.landCount++
			}
		}

		row++
	}
}

func (g *generator) addRoad(from *tile, toLoc int, hidden int, name *string) {
	r := &road{}
	r.entNum = g.rndAllocNum(sublocLow, sublocHigh)
	r.toLoc = toLoc
	r.hidden = hidden
	r.name = name

	from.roads = append(from.roads, r)
}

func (g *generator) linkRoads(from, to *tile, hidden int, name *string) {
	// If there is a sublocation at an endpoint of the secret road,
	// move the road to come from the sublocation instead of the province.
	//
	// Since only 1/3 of the locations have sublocs, this doesn't happen
	// all the time.  A very few locations will have the route between
	// two sublocs

	n := 0
	for i := 1; i <= g.topSubloc; i++ {
		if g.subloc[i].inside == from.region && g.subloc[i].terrain != terrCity {
			n = i
			break
		}
	}

	if n != 0 {
		from = g.subloc[n]
	}

	n = 0
	for i := 1; i <= g.topSubloc; i++ {
		if g.subloc[i].inside == to.region && g.subloc[i].terrain != terrCity {
			n = i
			break
		}
	}

	if n != 0 {
		to = g.subloc[n]
	}

	g.addRoad(from, to.region, hidden, name)
	g.addRoad(to, from.region, hidden, name)
}

func (g *generator) allocInside() int {
	g.insideTop++
	g.assert(g.insideTop < maxInside, 768, "inside_top < MAX_INSIDE")

	return g.insideTop
}

func (g *generator) dumpContinents() {
	for i := 1; i <= g.insideTop; i++ {
		fmt.Fprintf(g.locFp, "%d loc region\n", regionOff+i)
		if g.insideNames[i] != nil {
			fmt.Fprintf(g.locFp, "na %s\n", *g.insideNames[i])
		}
		g.printInsideLocs(i)
		fmt.Fprintf(g.locFp, "\n")
	}
}

func (g *generator) printInsideLocs(n int) {
	count := 0

	for _, t := range g.insideList[n] {
		count++

		if count == 1 {
			fmt.Fprintf(g.locFp, "LI\n")
			fmt.Fprintf(g.locFp, " hl ")
		}

		if count%11 == 10 { // continuation line
			fmt.Fprintf(g.locFp, "\\\n\t")
		}

		fmt.Fprintf(g.locFp, "%d ", t.region)
	}

	if count != 0 {
		fmt.Fprintf(g.locFp, "\n")
	}
}

func (g *generator) printMap() {
	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			t := g.m[row][col]
			if t == nil {
				continue
			}

			flag := true
			sl := false

			fmt.Fprintf(g.locFp, "%d loc %s\n", t.region, terrS[t.terrain])

			if t.name != nil && *t.name != "Unnamed" {
				fmt.Fprintf(g.locFp, "na %s\n", *t.name)
			}

			if t.uldimFlag != 0 {
				if !sl {
					fmt.Fprintf(g.locFp, "SL\n")
					sl = true
				}
				fmt.Fprintf(g.locFp, " uf %d\n", t.uldimFlag)
			}

			if t.summerbridgeFlag != 0 {
				if !sl {
					fmt.Fprintf(g.locFp, "SL\n")
					sl = true
				}
				fmt.Fprintf(g.locFp, " sf %d\n", t.summerbridgeFlag)
			}

			if t.safeHaven != 0 {
				if !sl {
					fmt.Fprintf(g.locFp, "SL\n")
				}
				fmt.Fprintf(g.locFp, " sh 1\n")
			}

			if t.inside != 0 {
				fmt.Fprintf(g.locFp, "LI\n")
				flag = false

				fmt.Fprintf(g.locFp, " wh %d\n", t.inside+regionOff)
			}

			g.printInsideSublocs(flag, row, col)

			fmt.Fprintf(g.locFp, "LO\n")
			fmt.Fprintf(g.locFp, " pd %d %d %d %d\n",
				g.provDest(t, dirN),
				g.provDest(t, dirE),
				g.provDest(t, dirS),
				g.provDest(t, dirW))

			if t.hidden != 0 {
				fmt.Fprintf(g.locFp, " hi %d\n", t.hidden)
			}
			if t.seaLane != 0 {
				fmt.Fprintf(g.locFp, " sl 1\n")
			}

			fmt.Fprintf(g.locFp, "\n")
		}
	}
}

func (g *generator) printSublocs() {
	for i := 1; i <= g.topSubloc; i++ {
		s := g.subloc[i]
		sl := false

		fmt.Fprintf(g.locFp, "%d loc %s\n", s.region, terrS[s.terrain])

		if s.name != nil && *s.name != "Unnamed" {
			fmt.Fprintf(g.locFp, "na %s\n", *s.name)
		}

		g.assert(s.inside != 0, 921, "subloc[i]->inside")
		fmt.Fprintf(g.locFp, "LI\n")
		fmt.Fprintf(g.locFp, " wh %d\n", s.inside)
		g.printSublocGates(i)

		fmt.Fprintf(g.locFp, "LO\n")

		if s.hidden != 0 {
			fmt.Fprintf(g.locFp, " hi %d\n", s.hidden)
		}

		if s.majorCity != 0 {
			if !sl {
				fmt.Fprintf(g.locFp, "SL\n")
				sl = true
			}
			fmt.Fprintf(g.locFp, " mc %d\n", s.majorCity)
		}

		if s.safeHaven != 0 {
			if !sl {
				fmt.Fprintf(g.locFp, "SL\n")
			}
			fmt.Fprintf(g.locFp, " sh 1\n")
		}

		fmt.Fprintf(g.locFp, "\n")
	}
}

// printHereList writes the " hl" list shared by print_subloc_gates() and
// print_inside_sublocs(). An "LI" header precedes it when li is set.
func (g *generator) printHereList(li bool, t *tile) {
	count := 0
	item := func(n int) {
		count++
		if count == 1 {
			if li {
				fmt.Fprintf(g.locFp, "LI\n")
			}
			fmt.Fprintf(g.locFp, " hl ")
		}

		if count%11 == 10 { // continuation line
			fmt.Fprintf(g.locFp, "\\\n\t")
		}

		fmt.Fprintf(g.locFp, "%d ", n)
	}

	for _, r := range t.roads {
		item(r.entNum)
	}
	for _, n := range t.gatesNum {
		item(n)
	}
	for _, n := range t.subs {
		item(n)
	}

	if count != 0 {
		fmt.Fprintf(g.locFp, "\n")
	}
}

func (g *generator) printSublocGates(n int) { // and inside buildings...
	g.printHereList(false, g.subloc[n])
}

func (g *generator) printInsideSublocs(flag bool, row, col int) {
	g.printHereList(flag, g.m[row][col])
}

// The entity number of a region determines where it is on the map.
// Entity 10101 corresponds to (1,1).  [10199] corresponds to (1,99).

func regionRow(where int) int {
	return (where / 100) % 100
}

func regionCol(where int) int {
	return where % 100
}

func (g *generator) rcToRegion(row, col int) int {
	g.assert(row >= 0 && row < 100, 1055, "row >= 0 && row < 100")
	g.assert(col >= 0 && col < 100, 1056, "col >= 0 && col < 100")

	return 10000 + (row * 100) + col
}

func (g *generator) dirAssert() {
	reg := g.rcToRegion(1, 1)
	g.assert(reg == 10101, 1071, "reg == 10101")
	g.assert(regionRow(reg) == 1, 1072, "region_row(reg) == 1")
	g.assert(regionCol(reg) == 1, 1073, "region_col(reg) == 1")

	reg = g.rcToRegion(99, 99)
	g.assert(reg == 19999, 1078, "reg == 19999")
	g.assert(regionRow(reg) == 99, 1079, "region_row(reg) == 99")
	g.assert(regionCol(reg) == 99, 1080, "region_col(reg) == 99")

	reg = g.rcToRegion(57, 63)
	g.assert(reg == 15763, 1085, "reg == 15763")
	g.assert(regionRow(reg) == 57, 1086, "region_row(reg) == 57")
	g.assert(regionCol(reg) == 63, 1087, "region_col(reg) == 63")
}

func (g *generator) adjacentTileWater(row, col int) *tile {
	var p *tile

	g.randomizeDirVector()

	for i := 1; !(p != nil && p.terrain == terrOcean) && i < maxDir; i++ {
		p = g.adjacentTileSup(row, col, g.dirVector[i])
	}

	if p != nil && p.terrain == terrOcean {
		return p
	}
	return nil
}

func (g *generator) adjacentTileTerr(row, col int) *tile {
	var p *tile

	g.randomizeDirVector()

	isTerr := func(p *tile) bool {
		return p != nil && p.terrain != terrLand && p.terrain != terrOcean
	}

	for i := 1; !isTerr(p) && i < maxDir; i++ {
		p = g.adjacentTileSup(row, col, g.dirVector[i])
	}

	if isTerr(p) {
		return p
	}
	return nil
}

// provDest returns the region immediately adjacent to t in direction dir,
// or 0 if there is no adjacent location in that direction.
func (g *generator) provDest(t *tile, dir int) int {
	row := t.row
	col := t.col

	switch dir {
	case dirN:
		row--
	case dirE:
		col++
	case dirS:
		row++
	case dirW:
		col--
	default:
		g.errf("location_direction: bad dir %d\n", dir)
		g.assert(false, 1157, "FALSE")
	}

	if row < 0 || row > 99 {
		return 0
	}

	if col < 0 {
		col = g.maxCol
	}
	if col > g.maxCol {
		col = 0
	}

	if g.m[row][col] == nil {
		return 0
	}

	return g.m[row][col].region
}

// adjacentTileSup returns the tile immediately adjacent to (row, col) in
// direction dir, or nil if there is none.
func (g *generator) adjacentTileSup(row, col, dir int) *tile {
	switch dir {
	case dirN:
		row--
	case dirNE:
		row--
		col++
	case dirE:
		col++
	case dirSE:
		row++
		col++
	case dirS:
		row++
	case dirSW:
		row++
		col--
	case dirW:
		col--
	case dirNW:
		row--
		col--
	default:
		g.errf("location_direction: bad dir %d\n", dir)
		g.assert(false, 1227, "FALSE")
	}

	if col < 0 {
		col = g.maxCol
	}
	if col > g.maxCol {
		col = 0
	}

	if row < 0 || row > 99 || col < 0 || col > 99 {
		return nil // off the map
	}

	return g.m[row][col]
}

func (g *generator) isPortCity(row, col int) bool {
	n := g.adjacentTileSup(row, col, dirN)
	s := g.adjacentTileSup(row, col, dirS)
	e := g.adjacentTileSup(row, col, dirE)
	w := g.adjacentTileSup(row, col, dirW)

	return n != nil && n.terrain == terrOcean ||
		s != nil && s.terrain == terrOcean ||
		e != nil && e.terrain == terrOcean ||
		w != nil && w.terrain == terrOcean
}

func (g *generator) randomizeDirVector() {
	g.dirVector[0] = 0
	for i := 1; i < maxDir; i++ {
		g.dirVector[i] = i
	}

	for i := 1; i < maxDir; i++ {
		swap := g.rnd(i, maxDir-1)

		if i != swap {
			g.dirVector[i], g.dirVector[swap] = g.dirVector[swap], g.dirVector[i]
		}
	}
}

// roadName cycles through the secret road names.
func roadName(cnt *int) *string {
	var name string
	switch *cnt % 3 {
	case 0:
		name = "Secret pass"
	case 1:
		name = "Secret route"
	case 2:
		name = "Old road"
	}
	*cnt++
	return strp(name)
}

// isMarked reports whether a neighbor exists and is marked. Off-map squares
// and holes are nil.
func isMarked(t *tile) bool {
	return t != nil && t.mark != 0
}

// The horror, the horror: bump each neighbor's mark by a coin flip.
func (g *generator) horror(tiles ...*tile) {
	for _, t := range tiles {
		if t != nil {
			t.mark += g.rnd(0, 1)
		}
	}
}

func (g *generator) bridgeMapHoleSup(row, col int) int {
	var l []int

	// Find all squares neighboring the hole

	n := g.adjacentTileSup(row, col, dirN)
	s := g.adjacentTileSup(row, col, dirS)
	e := g.adjacentTileSup(row, col, dirE)
	w := g.adjacentTileSup(row, col, dirW)
	nw := g.adjacentTileSup(row, col, dirNW)
	sw := g.adjacentTileSup(row, col, dirSW)
	ne := g.adjacentTileSup(row, col, dirNE)
	se := g.adjacentTileSup(row, col, dirSE)

	if isMarked(n) || isMarked(s) || isMarked(e) || isMarked(w) ||
		isMarked(nw) || isMarked(sw) || isMarked(ne) || isMarked(se) {
		return 0
	}

	// For every path across the hole, add it to the list of possibilities
	// if it's land-land and we haven't already used one of the destination
	// tiles for another hole-crossing.

	if n != nil && s != nil &&
		n.terrain != terrOcean && s.terrain != terrOcean &&
		g.m[n.row][n.col].mark+g.m[s.row][s.col].mark == 0 {
		l = append(l, 1)
	}

	if e != nil && w != nil &&
		e.terrain != terrOcean && w.terrain != terrOcean &&
		g.m[e.row][e.col].mark+g.m[w.row][w.col].mark == 0 {
		l = append(l, 2)
	}

	if ne != nil && sw != nil &&
		ne.terrain != terrOcean && sw.terrain != terrOcean &&
		g.m[ne.row][ne.col].mark+g.m[sw.row][sw.col].mark == 0 {
		l = append(l, 3)
	}

	if se != nil && nw != nil &&
		se.terrain != terrOcean && nw.terrain != terrOcean &&
		g.m[se.row][se.col].mark+g.m[nw.row][nw.col].mark == 0 {
		l = append(l, 4)
	}

	i := len(l)

	if i <= 0 {
		return 0
	}

	name := roadName(&g.holeRoadNameCnt)

	g.horror(n, e, w, s, nw, ne, sw, se)

	i = g.rnd(0, i-1)

	switch l[i] {
	case 1:
		g.linkRoads(n, s, 1, name)
		n.mark = 1
		s.mark = 1
	case 2:
		g.linkRoads(e, w, 1, name)
		e.mark = 1
		w.mark = 1
	case 3:
		g.linkRoads(ne, sw, 1, name)
		ne.mark = 1
		sw.mark = 1
	case 4:
		g.linkRoads(se, nw, 1, name)
		se.mark = 1
		nw.mark = 1
	}

	return l[i]
}

var bridgeDirS = []string{
	"-invalid-",
	"  n-s",
	"  e-w",
	"ne-sw",
	"nw-se",
}

// bridgeMapHoles bridges each # map hole with a secret road.
// It won't put two roads in the same square.
func (g *generator) bridgeMapHoles() {
	for row := 0; row <= g.maxRow; row++ {
		for col := 0; col <= g.maxCol; col++ {
			if g.m[row][col] == nil {
				if n := g.bridgeMapHoleSup(row, col); n != 0 {
					g.errf("%s map hole bridge at (%d,%d)\n", bridgeDirS[n], row, col)
				}
			}
		}
	}

	g.errf("\n")
}

func (g *generator) bridgeCornerSup(row, col int) int {
	var l []int

	// Find all squares neighboring the hole

	n := g.adjacentTileSup(row, col, dirN)
	s := g.adjacentTileSup(row, col, dirS)
	e := g.adjacentTileSup(row, col, dirE)
	w := g.adjacentTileSup(row, col, dirW)
	nw := g.adjacentTileSup(row, col, dirNW)
	sw := g.adjacentTileSup(row, col, dirSW)
	ne := g.adjacentTileSup(row, col, dirNE)
	se := g.adjacentTileSup(row, col, dirSE)

	if isMarked(n) || isMarked(s) || isMarked(e) || isMarked(w) ||
		isMarked(nw) || isMarked(sw) || isMarked(ne) || isMarked(se) {
		return 0
	}

	// The name counter advances even when no road is laid.
	name := roadName(&g.cornerRoadNameCnt)

	if nw != nil && nw.terrain != terrOcean {
		l = append(l, 1)
	}
	if ne != nil && ne.terrain != terrOcean {
		l = append(l, 2)
	}
	if se != nil && se.terrain != terrOcean {
		l = append(l, 3)
	}
	if sw != nil && sw.terrain != terrOcean {
		l = append(l, 4)
	}

	i := len(l)

	if i <= 0 {
		return 0
	}

	g.horror(n, e, w, s, nw, ne, sw, se)

	i = g.rnd(0, i-1)

	here := g.m[row][col]
	var to *tile
	switch l[i] {
	case 1:
		to = nw
	case 2:
		to = ne
	case 3:
		to = se
	case 4:
		to = sw
	}
	g.linkRoads(here, to, 1, name)
	here.mark = 1
	to.mark = 1

	return l[i]
}

func (g *generator) bridgeCaddyCorners() {
	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			if g.m[row][col] != nil && g.m[row][col].terrain != terrOcean &&
				g.rnd(1, 35) == 35 {
				g.bridgeCornerSup(row, col)
			}
		}
	}
}

func (g *generator) bridgeMountainSup(row, col int) {
	from := g.m[row][col]
	to := g.adjacentTileWater(row, col)

	g.assert(to.terrain == terrOcean, 1579, "to->terrain == terr_ocean")

	var name *string
	switch g.rnd(1, 3) {
	case 1:
		name = strp("Narrow channel")
	case 2:
		name = strp("Rocky channel")
	case 3:
		name = strp("Secret sea route")
	}

	g.addRoad(from, to.region, 1, name)
	g.addRoad(to, from.region, 1, name)

	g.errf("secret sea route at (%d,%d)\n", from.row, from.col)
}

func (g *generator) bridgeMountainPorts() {
	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			if g.m[row][col] != nil &&
				g.m[row][col].terrain == terrMountain &&
				g.isPortCity(row, col) &&
				g.rnd(1, 7) == 7 {
				g.bridgeMountainSup(row, col)
			}
		}
	}
}

func (g *generator) makeRoads() {
	g.clearProvinceMarks()
	g.bridgeMapHoles()
	g.bridgeCaddyCorners()
	g.bridgeMountainPorts()
}

const continentFormat = "%s  %s  %s  %s  %s\n"

func (g *generator) printContinentRow(name, coord, nprovs, ncities, gates string) {
	g.errf(continentFormat,
		clib.PadRight(name, 25),
		clib.PadLeft(coord, 8),
		clib.PadLeft(nprovs, 6),
		clib.PadLeft(ncities, 7),
		gates)
}

func (g *generator) printContinent(i int) {
	p := g.insideList[i][0]

	var name string
	if g.insideNames[i] == nil {
		name = fmt.Sprintf("?? (%d,%d)", p.row, p.col)
	} else {
		name = *g.insideNames[i]
	}

	coord := fmt.Sprintf("(%d,%d)", p.row, p.col)
	gates := fmt.Sprintf("%d/%d", g.insideGatesFrom[i], g.insideGatesTo[i])
	nprovs := fmt.Sprintf("%d", len(g.insideList[i]))
	ncities := fmt.Sprintf("%d", g.insideNumCities[i])

	g.printContinentRow(name, coord, nprovs, ncities, gates)
}

func (g *generator) countContinents() {
	g.errf("\nLand regions:\n\n")

	g.printContinentRow("name", "coord", "nprovs", "ncities", "gates (out/in)")
	g.printContinentRow("----", "-----", "------", "-------", "--------------")

	for i := 1; i <= g.insideTop; i++ {
		if g.insideList[i][0].terrain != terrOcean {
			g.printContinent(i)
		}
	}

	g.errf("\n\nOceans:\n\n")

	g.printContinentRow("name", "coord", "nprovs", "ncities", "gates (out/in)")
	g.printContinentRow("----", "-----", "------", "-------", "--------------")

	for i := 1; i <= g.insideTop; i++ {
		if g.insideList[i][0].terrain == terrOcean {
			g.printContinent(i)
		}
	}

	g.errf("\n\n%d continents\n", g.insideTop)
	g.errf("%d land locs\n", g.landCount)
	g.errf("%d water locs\n", g.waterCount)
}

func (g *generator) fixTerrainLand() {
	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			t := g.m[row][col]
			if t == nil || t.terrain != terrLand {
				continue
			}
			p := g.adjacentTileTerr(row, col)
			if p != nil && p.terrain != terrLand && p.terrain != terrOcean {
				t.terrain = p.terrain
				t.color = p.color
			} else {
				g.errf("fix_terrain: could not infer type of (%d,%d)\n", row, col)
				g.errf("    assuming 'forest'\n")
				t.terrain = terrForest
			}
		}
	}
}

// province returns the tile at (row, col), or nil when there is none: the
// square is a hole, past the end of its line, or off the map.
func (g *generator) province(row, col int) *tile {
	if row < 0 || row >= maxRow || col < 0 || col >= maxCol {
		return nil
	}
	return g.m[row][col]
}

// restOfLine returns the text after the first whitespace-delimited field,
// as set_regions() finds the region name.
func restOfLine(buf string) string {
	i := 0
	for i < len(buf) && !iswhite(buf[i]) {
		i++
	}
	for i < len(buf) && iswhite(buf[i]) {
		i++
	}
	return buf[i:]
}

// chomp truncates buf at its first newline.
func chomp(buf string) string {
	for i := 0; i < len(buf); i++ {
		if buf[i] == '\n' {
			return buf[:i]
		}
	}
	return buf
}

func (g *generator) setRegions() {
	// row and col are uninitialized in C; a Regions file whose first line
	// does not start with "row,col" reads indeterminate values there.
	var row, col int
	landCount := 0
	waterCount := 0

	data, err := clib.ReadFile(g.path("Regions"))
	if err != nil {
		clib.Perror(g.stderr, "can't read Regions", err)
		return
	}

	for _, buf := range clib.FgetsLines(data, lineLen) {
		buf = chomp(buf)

		if buf == "" {
			continue
		}

		sc := clib.NewScanner(buf)
		if v, ok := sc.Int(); ok {
			row = v
			if sc.Literal(',') {
				if v, ok := sc.Int(); ok {
					col = v
				}
			}
		}

		p := restOfLine(buf)

		if g.province(row, col) == nil {
			g.errf("Regions '%s': no province at (%d,%d)\n", p, row, col)
			continue
		}

		if ins := g.m[row][col].inside; ins != 0 {
			g.errf("collision between %s and %s at (%d,%d)\n", p, clib.Str(g.insideNames[ins]), row, col)
			continue
		}

		ins := g.allocInside()
		g.insideNames[ins] = strp(p)

		if g.m[row][col].terrain == terrOcean {
			waterCount++
			g.floodWaterInside(row, col, ins)
		} else {
			landCount++
			g.floodLandInside(row, col, ins)
		}
	}

	g.errf("set_regions: named %d land regions, %d water regions\n", landCount, waterCount)

	// Now locate unnamed regions

	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			if g.m[row][col] != nil && g.m[row][col].inside == 0 {
				ins := g.allocInside()

				if g.m[row][col].terrain == terrOcean {
					n := g.floodWaterInside(row, col, ins)
					g.errf("\tunnamed sea at  %d,%d (%d locs)\n", row, col, n)
				} else {
					n := g.floodLandInside(row, col, ins)
					g.errf("\tunnamed land at %d,%d (%d locs)\n", row, col, n)
				}
			}
		}
	}

	// Now collect the list of provinces in each region

	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			if t := g.m[row][col]; t != nil && t.inside != 0 {
				g.insideList[t.inside] = append(g.insideList[t.inside], t)
			}
		}
	}
}

func (g *generator) floodLandInside(row, col, ins int) int {
	count := 1

	g.m[row][col].inside = ins

	if g.m[row][col].regionBoundary != 0 {
		return count
	}

	for dir := 1; dir < maxDir; dir++ {
		p := g.adjacentTileSup(row, col, dir)

		if p == nil || p.terrain == terrOcean {
			continue
		}

		if p.inside == ins {
			continue // already been here
		}

		if p.inside != 0 && p.regionBoundary != 0 {
			continue // another region's boundary
		}

		if p.inside != 0 {
			g.errf("flood_land_inside(%d,%d,%s) error\n", row, col, clib.Str(g.insideNames[ins]))
			g.assert(false, 1858, "FALSE")
		}

		count += g.floodLandInside(p.row, p.col, ins)
	}

	return count
}

func (g *generator) floodWaterInside(row, col, ins int) int {
	count := 1

	g.m[row][col].inside = ins

	for dir := 1; dir < maxDir; dir++ {
		p := g.adjacentTileSup(row, col, dir)

		if p == nil || p.color == -1 || p.color != g.m[row][col].color {
			continue
		}

		if p.inside == ins {
			continue // already been here
		}

		if p.inside != 0 {
			g.errf("flood_water_inside(%d,%d,%s) error\n", row, col, clib.Str(g.insideNames[ins]))
			g.assert(false, 1890, "FALSE")
		}

		count += g.floodWaterInside(p.row, p.col, ins)
	}

	return count
}

// setProvinceClumps names groups of provinces.
func (g *generator) setProvinceClumps() {
	// As in setRegions, these keep their values from the previous line when
	// sscanf() cannot assign them.
	var row, col int
	var typ byte
	count := 0

	data, err := clib.ReadFile(g.path("Land"))
	if err != nil {
		clib.Perror(g.stderr, "can't read Land", err)
		return
	}

	for _, buf := range clib.FgetsLines(data, lineLen) {
		buf = chomp(buf)

		sc := clib.NewScanner(buf)
		if v, ok := sc.Int(); ok {
			row = v
			if sc.Literal(',') {
				if v, ok := sc.Int(); ok {
					col = v
					sc.SkipSpace()
					if c, ok := sc.Char(); ok {
						typ = c
					}
				}
			}
		}

		// Skip two whitespace-delimited fields to reach the name.
		i := 0
		for i < len(buf) && !iswhite(buf[i]) {
			i++
		}
		for i < len(buf) && iswhite(buf[i]) {
			i++
		}
		for i < len(buf) && !iswhite(buf[i]) {
			i++
		}
		for i < len(buf) && iswhite(buf[i]) {
			i++
		}
		p := buf[i:]

		if g.province(row, col) == nil {
			g.errf("Land '%s': no province at (%d,%d)\n", p, row, col)
			continue
		}

		if g.m[row][col].saveChar != typ {
			g.errf("Land '%s' expects '%s' at (%d,%d), got '%s'\n",
				p, clib.Ch(typ), row, col, clib.Ch(g.m[row][col].saveChar))
		}

		if g.m[row][col].name != nil {
			g.errf("clump collision between %s and %s at (%d,%d)\n", p, *g.m[row][col].name, row, col)
		}

		g.floodLandClumps(row, col, strp(p))
		count++
	}

	g.errf("set_province_clumps: named %d areas\n", count)
}

func (g *generator) unnamedProvinceClumps() {
	g.errf("Unnamed areas\n\n")

	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			t := g.m[row][col]
			if t == nil || t.name != nil {
				continue
			}
			if t.terrain == terrOcean {
				continue
			}

			n := g.floodLandClumps(row, col, strp("Unnamed"))

			if t.saveChar != 'o' {
				g.errf("%d,%d\t%s\t%d unnamed\n", row, col, clib.Ch(t.saveChar), n)
			}
		}
	}

	g.errf("\n")
}

func (g *generator) floodLandClumps(row, col int, name *string) int {
	count := 1

	g.m[row][col].name = name

	for dir := 1; dir < maxDir; dir++ {
		p := g.adjacentTileSup(row, col, dir)

		if p == nil || p.terrain == terrOcean || p.color == -1 ||
			p.color != g.m[row][col].color {
			continue
		}

		if p.name == name {
			continue // already been here
		}

		if p.name != nil {
			g.errf("flood_land_clumps(%d,%d,%s) error\n", row, col, *name)
			g.assert(false, 2010, "FALSE")
		}

		count += g.floodLandClumps(p.row, p.col, name)
	}

	return count
}

func (g *generator) makeIslands() {
	g.numIslands = g.waterCount / 3

	for i := 1; i <= g.numIslands; i++ {
		row := g.rnd(0, g.maxRow)
		col := g.rnd(0, g.maxCol)

		if g.m[row][col] != nil && g.m[row][col].terrain == terrOcean &&
			g.islandAllowed(row, col) {
			g.createASubloc(row, col, g.rnd(0, 1), terrIsland)
		} else {
			i--
		}
	}
}

func (g *generator) islandAllowed(row, col int) bool {
	inside := g.m[row][col].inside

	if inside == 0 {
		return true
	}

	// An unnamed region has no name, so it is never "Deep".
	name := ""
	if g.insideNames[inside] != nil {
		name = *g.insideNames[inside]
	}
	for p := 0; p < len(name); p++ {
		if len(name)-p >= 4 && name[p:p+4] == "Deep" {
			return false
		}
	}

	return true
}

func (g *generator) createASubloc(row, col, hidden, kind int) int {
	g.topSubloc++
	g.assert(g.topSubloc < maxSubloc, 2122, "top_subloc < MAX_SUBLOC")

	s := &tile{}
	g.subloc[g.topSubloc] = s
	if kind == terrCity {
		s.region = g.rndAllocNum(cityLow, cityHigh)
	} else {
		s.region = g.rndAllocNum(sublocLow, sublocHigh)
	}
	s.inside = g.m[row][col].region
	s.row = row
	s.col = col
	s.hidden = hidden
	s.terrain = kind
	s.depth = 3

	if kind == terrCity {
		g.m[row][col].city = 2
	}

	g.m[row][col].subs = append(g.m[row][col].subs, s.region)

	return g.topSubloc
}

// islandCountFormat and its sibling below print one histogram line.
func (g *generator) printHistogram(count []int, total int, what string) {
	for i := len(count) - 1; i >= 0; i-- {
		if count[i] != 0 {
			break
		}
		count[i] = -1
	}

	for i := range count {
		if count[i] == -1 {
			break
		}

		has := "locs have"
		if count[i] == 1 {
			has = "loc has"
		}
		plural := "s"
		if i == 1 {
			plural = " "
		}
		g.errf("%6d %s %d %s%s (%d%%)\n", count[i], has, i, what, plural, count[i]*100/total)
	}
}

func (g *generator) countSublocs() {
	var count [100]int

	g.errf("\nsubloc counts:\n")

	g.clearProvinceMarks()

	for i := 1; i <= g.topSubloc; i++ {
		if g.subloc[i].terrain == terrIsland {
			g.m[g.subloc[i].row][g.subloc[i].col].mark++
		}
	}

	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			if t := g.m[row][col]; t != nil && t.terrain == terrOcean {
				count[t.mark]++
			}
		}
	}

	g.printHistogram(count[:], g.waterCount, "island")
}

// dumpRoadList writes the road entities that leave t.
func (g *generator) dumpRoadList(t *tile) {
	for _, r := range t.roads {
		fmt.Fprintf(g.roadFp, "%d road 0\n", r.entNum)
		if r.name != nil {
			fmt.Fprintf(g.roadFp, "na %s\n", *r.name)
		}
		fmt.Fprintf(g.roadFp, "LI\n")
		fmt.Fprintf(g.roadFp, " wh %d\n", t.region)
		fmt.Fprintf(g.roadFp, "GA\n")
		fmt.Fprintf(g.roadFp, " tl %d\n", r.toLoc)
		if r.hidden != 0 {
			fmt.Fprintf(g.roadFp, " rh %d\n", r.hidden)
		}
		fmt.Fprintf(g.roadFp, "\n")
	}
}

func (g *generator) dumpRoads() {
	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			if g.m[row][col] != nil {
				g.dumpRoadList(g.m[row][col])
			}
		}
	}

	for i := 1; i <= g.topSubloc; i++ {
		g.dumpRoadList(g.subloc[i])
	}
}

// dumpGateList writes the gate entities that leave t. Provinces iterate over
// gates_dest and sublocations over gates_num; the lists are the same length.
func (g *generator) dumpGateList(t *tile) {
	for j := range t.gatesNum {
		fmt.Fprintf(g.gateFp, "%d gate 0\n", t.gatesNum[j])
		fmt.Fprintf(g.gateFp, "LI\n")
		fmt.Fprintf(g.gateFp, " wh %d\n", t.region)
		fmt.Fprintf(g.gateFp, "GA\n")
		fmt.Fprintf(g.gateFp, " tl %d\n", t.gatesDest[j])
		if t.gatesKey[j] != 0 {
			fmt.Fprintf(g.gateFp, " sk %d\n", t.gatesKey[j])
		}
		fmt.Fprintf(g.gateFp, "\n")
	}
}

func (g *generator) dumpGates() {
	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			if g.m[row][col] != nil {
				g.dumpGateList(g.m[row][col])
			}
		}
	}

	for i := 1; i <= g.topSubloc; i++ {
		g.dumpGateList(g.subloc[i])
	}
}

func (g *generator) randomTileFromEachRegion() []*tile {
	var l []*tile

	for i := 1; i <= g.insideTop; i++ {
		if g.insideList[i][0].terrain == terrOcean {
			continue
		}

		if g.isImpassable(i) {
			continue
		}

		var j int
		for {
			j = g.rnd(0, len(g.insideList[i])-1)
			if g.insideList[i][j].safeHaven == 0 {
				break
			}
		}

		l = append(l, g.insideList[i][j])
	}

	g.tilesShuffle(l)

	return l
}

func (g *generator) shiftTourEndpoints(l []*tile) []*tile {
	var other []*tile

	for i := range l {
		p := g.adjacentTileTerr(l[i].row, l[i].col)

		if p == nil {
			p = l[i]
		}

		q := g.adjacentTileTerr(p.row, p.col)
		if q == l[i] { // doubled back, retry
			q = g.adjacentTileTerr(p.row, p.col)
		}

		if q == nil || q.terrain == terrOcean || q.safeHaven != 0 {
			g.errf("couldn't shift tour (%d,%d)\n", l[i].row, l[i].col)
			other = append(other, l[i])
		} else {
			other = append(other, q)
		}
	}

	return other
}

// Gate laying plan
//
//	province->island->province	(gate_province_islands)
//	province->province		(random_province_gates)
//	continental tour		(gate_continental_tour)
//
//	each region gets a stone ring	(gate_stone_circles)
//	each with links to 2 other
//	rings, and 5 random provinces
//
//	7 rings of 5-12 provinces	(gate_land_ring)
//	12 rings of 6-12 islands	(gate_link_islands)

func (g *generator) makeGates() {
	g.gateProvinceIslands((g.landCount + 199) / 200)
	g.randomProvinceGates((g.landCount + 199) / 200)
	g.gateContinentalTour()
	g.gateStoneCircles()

	g.gateLandRing((g.landCount + 749) / 750)
	g.gateLinkIslands((g.numIslands + 149) / 150) // disjoint
	g.gateLinkIslands(g.numIslands / 450)         // might cross
}

func (g *generator) newGate(from, to *tile, key int) {
	gateNum := g.rndAllocNum(sublocLow, sublocHigh)

	from.gatesNum = append(from.gatesNum, gateNum)
	from.gatesDest = append(from.gatesDest, to.region)
	from.gatesKey = append(from.gatesKey, key)

	// Gather statistics:

	g.insideGatesFrom[g.m[from.row][from.col].inside]++
	g.insideGatesTo[g.m[to.row][to.col].inside]++
}

func (g *generator) randomProvinceGates(n int) {
	g.clearProvinceMarks()
	g.markBadLocs()

	for i := 0; i < n; i++ {
		r1, c1 := g.randomProvince(0)
		r2, c2 := g.randomProvince(0)

		g.newGate(g.m[r1][c1], g.m[r2][c2], 0)
	}
}

func (g *generator) gateProvinceIslands(times int) {
	g.clearProvinceMarks()
	g.markBadLocs()
	g.clearSublocMarks()

	for j := 1; j <= times; j++ {
		r1, c1 := g.randomProvince(0)
		isle := g.randomIsland()
		r2, c2 := g.randomProvince(0)

		g.newGate(g.m[r1][c1], g.subloc[isle], 0)
		g.newGate(g.subloc[isle], g.m[r2][c2], 0)
	}
}

func (g *generator) gateContinentalTour() {
	l := g.randomTileFromEachRegion()
	m := g.shiftTourEndpoints(l)

	g.assert(len(l) == len(m), 2493, "tiles_len(l) == tiles_len(m)")

	g.errf("\nContinental gate tour:\n")

	i := 0
	for ; i < len(l)-1; i++ {
		g.errf("\t(%2d,%2d) -> (%2d,%2d)\n", l[i].row, l[i].col, m[i+1].row, m[i+1].col)

		g.newGate(l[i], m[i+1], 0)
	}

	g.errf("\t(%2d,%2d) -> (%2d,%2d)\n\n", l[i].row, l[i].col, m[0].row, m[0].col)

	g.newGate(l[i], m[0], g.rnd(111, 333))
}

func (g *generator) gateLinkIslands(rings int) {
	g.clearSublocMarks()

	for j := 1; j <= rings; j++ {
		num := g.rnd(5, 10)

		first := g.randomIsland()
		n := first

		for i := 1; i < num; i++ {
			next := g.randomIsland()
			g.newGate(g.subloc[n], g.subloc[next], 0)
			n = next
		}

		g.newGate(g.subloc[n], g.subloc[first], 0)
	}
}

func (g *generator) gateLandRing(rings int) {
	g.clearProvinceMarks()
	g.markBadLocs()

	for j := 1; j <= rings; j++ {
		num := g.rnd(5, 10)
		rFirst, cFirst := g.randomProvince(0)

		rN, cN := rFirst, cFirst

		for i := 1; i < num; i++ {
			rNext, cNext := g.randomProvince(0)
			g.newGate(g.m[rN][cN], g.m[rNext][cNext], 0)

			rN, cN = rNext, cNext
		}

		g.newGate(g.m[rN][cN], g.m[rFirst][cFirst], 0)
	}
}

func (g *generator) chooseRandomStoneCircle(l []*tile, avoid1, avoid2 *tile) *tile {
	for {
		i := g.rnd(0, len(l)-1)
		if l[i] != avoid1 && l[i] != avoid2 {
			return l[i]
		}
	}
}

// gateStoneCircles gives every region a hidden stone circle. Each stone
// circle has two gates to other stone circles, chosen at random, and five
// gates to random provinces.
func (g *generator) gateStoneCircles() {
	var circs []*tile

	l := g.randomTileFromEachRegion()

	g.errf("\nRing of stones:\n")

	for i := range l {
		n := g.createASubloc(l[i].row, l[i].col, 1, terrStoneCir)
		circs = append(circs, g.subloc[n])

		name := "??"
		if g.insideNames[l[i].inside] != nil {
			name = *g.insideNames[l[i].inside]
		}
		g.errf("\t(%2d,%2d) in %s\n", l[i].row, l[i].col, name)
	}

	// Cross-circle links need at least three circles (G3 issue #77).
	if len(circs) >= 3 {
		for i := range circs {
			first := g.chooseRandomStoneCircle(circs, circs[i], nil)
			second := g.chooseRandomStoneCircle(circs, circs[i], first)

			g.newGate(circs[i], first, g.rnd(111, 333))
			g.newGate(circs[i], second, g.rnd(111, 333))
		}
	}

	g.clearProvinceMarks()
	g.markBadLocs()

	for i := range circs {
		for j := 1; j <= 5; j++ {
			row, col := g.randomProvince(0)
			key := 0
			if g.rnd(0, 1) == 0 {
				key = g.rnd(111, 333)
			}
			g.newGate(circs[i], g.m[row][col], key)
		}
	}
}

// isImpassable reports whether region i is the Impassable Mountains.
// An unnamed region never is.
func (g *generator) isImpassable(i int) bool {
	return g.insideNames[i] != nil && *g.insideNames[i] == "Impassable Mountains"
}

func (g *generator) clearProvinceMarks() {
	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			if g.m[row][col] != nil {
				g.m[row][col].mark = 0
			}
		}
	}
}

func (g *generator) clearSublocMarks() {
	for i := 1; i <= g.topSubloc; i++ {
		g.subloc[i].mark = 0
	}
}

func (g *generator) markBadLocs() {
	for i := 1; i <= g.insideTop; i++ {
		if g.isImpassable(i) {
			for _, t := range g.insideList[i] {
				t.mark = 1
			}
		}
	}

	// don't put gates in locations where magic can't be used
	for r := 0; r <= g.maxRow; r++ {
		for c := 0; c <= g.maxCol; c++ {
			if g.m[r][c] != nil && g.m[r][c].safeHaven != 0 {
				g.m[r][c].mark = 1
			}
		}
	}
}

// randomProvince picks an unmarked province of terrain terr, or of any land
// terrain when terr is 0, and marks it.
func (g *generator) randomProvince(terr int) (int, int) {
	ok := func(t *tile) bool {
		if t == nil || t.mark != 0 {
			return false
		}
		if terr == 0 {
			return t.terrain != terrOcean
		}
		return t.terrain == terr
	}

	sum := 0
	for r := 0; r <= g.maxRow; r++ {
		for c := 0; c <= g.maxCol; c++ {
			if ok(g.m[r][c]) {
				sum++
			}
		}
	}

	n := g.rnd(1, sum)

	for r := 0; r <= g.maxRow; r++ {
		for c := 0; c <= g.maxCol; c++ {
			if ok(g.m[r][c]) {
				n--
				if n <= 0 {
					g.m[r][c].mark = 1
					return r, c
				}
			}
		}
	}

	g.assert(false, 2794, "FALSE")
	return 0, 0
}

func (g *generator) randomIsland() int {
	var i int
	for {
		n := g.rnd(1, g.numIslands)
		i = 1

		for i <= g.topSubloc {
			if g.subloc[i].terrain == terrIsland {
				n--
				if n <= 0 {
					break
				}
			}
			i++
		}

		g.assert(n <= g.topSubloc, 2824, "n <= top_subloc")

		if g.subloc[i].mark == 0 {
			break
		}
	}

	g.subloc[i].mark = 1

	return i
}

type guildName struct {
	skill  int
	weight int
	name   string // "" is NULL
}

var guildNames = []guildName{
	{terrStoneCir, 1, ""},
	{terrGrove, 1, ""},
	{terrBog, 1, ""},
	{terrCave, 1, ""},

	{terrGrave, 20, ""},
	{terrGrave, 1, "Barrows"},
	{terrGrave, 1, "Barrow Downs"},
	{terrGrave, 1, "Barrow Hills"},
	{terrGrave, 1, "Cairn Hills"},
	{terrGrave, 1, "Catacombs"},
	{terrGrave, 1, "Grave Mounds"},
	{terrGrave, 1, "Place of the Dead"},
	{terrGrave, 1, "Cemetary Hill"},
	{terrGrave, 1, "Fields of Death"},

	{terrRuins, 1, ""},

	{terrBattlefield, 3, "Old battlefield"},
	{terrBattlefield, 1, "Ancient battlefield"},
	{terrBattlefield, 1, ""},

	{terrEnchFor, 1, ""},
	{terrRockyHill, 1, ""},
	{terrTreeCir, 1, ""},
	{terrPits, 1, "Cursed Pits"},
	{terrPasture, 1, ""},
	{terrPasture, 1, "Grassy field"},
	{terrSacGrove, 1, ""},
	{terrOasis, 1, ""},
	{terrPopField, 1, ""},
	{terrSandPit, 1, ""},
	{terrYewGrove, 1, ""},

	{terrTemple, 1, ""},
	{terrLair, 1, ""},
}

// nameGuild returns a weighted random name for a subloc of kind skill.
// Callers treat NULL and "" alike, so "" stands in for NULL.
func (g *generator) nameGuild(skill int) string {
	sum := 0
	for _, gn := range guildNames {
		if gn.skill == skill {
			sum += gn.weight
		}
	}

	g.assert(sum > 0, 2892, "sum > 0")

	n := g.rnd(1, sum)

	for _, gn := range guildNames {
		if gn.skill == skill {
			n -= gn.weight

			if n <= 0 {
				return gn.name
			}
		}
	}

	g.assert(false, 2905, "FALSE")
	return ""
}

func (g *generator) countCities() {
	for i := 1; i <= g.topSubloc; i++ {
		if g.subloc[i].terrain == terrCity {
			ins := g.m[g.subloc[i].row][g.subloc[i].col].inside
			g.insideNumCities[ins]++
		}
	}
}

func (g *generator) countSublocCoverage() {
	var count [100]int

	g.clearProvinceMarks()

	for i := 1; i <= g.topSubloc; i++ {
		if g.subloc[i].depth == 3 {
			t := g.m[g.subloc[i].row][g.subloc[i].col]
			t.mark++

			if t.mark >= 5 {
				g.errf("(%d,%d) has %d sublocs (region %d)\n", g.subloc[i].row, g.subloc[i].col, t.mark, t.region)
			}
		}
	}

	g.errf("\nsubloc coverage:\n")

	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			if t := g.m[row][col]; t != nil && t.terrain != terrOcean {
				count[t.mark]++
			}
		}
	}

	g.printHistogram(count[:], g.landCount, "subloc")
}

type locEntry struct {
	terr   int // terrain appropriate
	kind   int // what to make there
	weight int // weight given to selection
	hidden int // 0=no, 1=yes, 2=rnd(0,1)
}

var locTable = []locEntry{
	{terrDesert, terrCave, 10, 1},
	{terrDesert, terrOasis, 10, 2},
	{terrDesert, terrSandPit, 10, 2},

	{terrMountain, terrRuins, 10, 1},
	{terrMountain, terrCave, 10, 1},
	{terrMountain, terrYewGrove, 8, 2},
	{terrMountain, terrTreeCir, 5, 1},
	{terrMountain, terrPits, 5, 1},
	{terrMountain, terrLair, 10, 2},
	{terrMountain, terrBattlefield, 6, 2},

	{terrSwamp, terrBog, 10, 2},
	{terrSwamp, terrPits, 10, 2},
	{terrSwamp, terrBattlefield, 6, 2},
	{terrSwamp, terrLair, 5, 2},

	{terrForest, terrRuins, 10, 1},
	{terrForest, terrRockyHill, 10, 0},
	{terrForest, terrTreeCir, 10, 1},
	{terrForest, terrEnchFor, 8, 1},
	{terrForest, terrPasture, 5, 0},
	{terrForest, terrYewGrove, 10, 2},
	{terrForest, terrCave, 10, 1},
	{terrForest, terrGrove, 9, 1},
	{terrForest, terrBattlefield, 6, 2},
	{terrForest, terrLair, 3, 1},

	{terrPlain, terrRuins, 10, 1},
	{terrPlain, terrPasture, 10, 0},
	{terrPlain, terrRockyHill, 10, 0},
	{terrPlain, terrSacGrove, 6, 2},
	{terrPlain, terrTreeCir, 6, 1},
	{terrPlain, terrPopField, 10, 0},
	{terrPlain, terrCave, 10, 1},
	{terrPlain, terrBattlefield, 6, 2},
}

func (g *generator) makeAppropriateSubloc(row, col int) {
	terr := g.m[row][col].terrain

	sum := 0
	for _, e := range locTable {
		if e.terr == terr {
			sum += e.weight
		}
	}

	if sum <= 0 {
		g.errf("no subloc appropriate for (%d,%d)\n", row, col)
		return
	}

	n := g.rnd(1, sum)

	for _, e := range locTable {
		if e.terr != terr {
			continue
		}

		n -= e.weight

		if n <= 0 {
			var hid int
			if e.hidden == 2 {
				hid = g.rnd(0, 1)
			} else {
				hid = e.hidden
			}

			sl := g.createASubloc(row, col, hid, e.kind)
			if s := g.nameGuild(e.kind); s != "" {
				g.subloc[sl].name = strp(s)
			}
			break
		}
	}
}

func (g *generator) createAGraveyard(row, col int) {
	n := g.createASubloc(row, col, g.rnd(0, 1), terrGrave)

	if s := g.nameGuild(terrGrave); s != "" {
		g.subloc[n].name = strp(s)
	}
}

func (g *generator) makeGraveyards() {
	for i := 1; i <= g.insideTop; i++ {
		if g.insideList[i][0].terrain == terrOcean {
			continue
		}

		n := len(g.insideList[i])

		if n < 10 {
			continue
		}

		l := append([]*tile(nil), g.insideList[i]...)
		g.tilesShuffle(l)

		for j := 0; j < n/10; j++ {
			g.createAGraveyard(l[j].row, l[j].col)
		}
	}
}

func (g *generator) createACity(row, col int, name *string, major int) int {
	// The C code keeps the Cities stream in a static; a failed open is
	// retried, and reported, on every call.
	if g.cities == nil {
		data, err := clib.ReadFile(g.path("Cities"))
		if err != nil {
			clib.Perror(g.stderr, "can't open Cities", err)
		} else {
			g.cities = &lineReader{data: data}
		}
	}

	if name == nil && g.cities != nil {
		if line, ok := g.cities.getlin(); ok {
			name = strp(line)
		}
	}

	n := g.createASubloc(row, col, 0, terrCity)

	g.subloc[n].name = name
	g.subloc[n].majorCity = major

	return n
}

func (g *generator) placeSublocations() {
	var l []int

	for row := 0; row < maxRow; row++ {
		for col := 0; col < maxCol; col++ {
			if g.m[row][col] != nil && g.m[row][col].terrain != terrOcean {
				l = append(l, row*1000+col)
			}
		}
	}

	g.ilistShuffle(l)

	for _, rc := range l {
		row := rc / 1000
		col := rc % 1000

		// Put a city everywhere there is a * or every 1 in 15 locs,
		// randomly.  Don't put one where there already is a city (city != 2).
		//
		// Nothing sets city to 1, so only the random test applies, and it
		// always draws.
		if g.m[row][col].city == 1 || (g.rnd(1, 15) == 1 && g.m[row][col].city != 2) {
			g.createACity(row, col, nil, 0)
		}

		if g.rnd(1, 100) <= 35 {
			g.makeAppropriateSubloc(row, col)
		}
		if g.rnd(1, 100) <= 35 {
			g.makeAppropriateSubloc(row, col)
		}
		if g.rnd(1, 100) <= 35 {
			g.makeAppropriateSubloc(row, col)
		}
	}
}

func (g *generator) countTiles() {
	var count [1000]int

	for r := 0; r < maxRow; r++ {
		for c := 0; c < maxCol; c++ {
			if g.m[r][c] != nil {
				count[g.m[r][c].terrain]++
			}
		}
	}

	for i := 1; i <= g.topSubloc; i++ {
		count[g.subloc[i].terrain]++
	}

	for i := 1; i < len(terrS); i++ {
		g.errf("%s %d\n", clib.PadRight(terrS[i], 30), count[i])
	}
}

// tilesShuffle and ilistShuffle are the Knuth shuffles in lib/tiles.c and
// lib/ilist.c.
func (g *generator) tilesShuffle(l []*tile) {
	n := len(l) - 1
	for i := 0; i < n; i++ {
		if r := g.rnd(i, n); r != i {
			l[i], l[r] = l[r], l[i]
		}
	}
}

func (g *generator) ilistShuffle(l []int) {
	n := len(l) - 1
	for i := 0; i < n; i++ {
		if r := g.rnd(i, n); r != i {
			l[i], l[r] = l[r], l[i]
		}
	}
}
