// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/mdhender/g3g0/internal/clib"
)

// TestLoadSaveMatchesC runs testdata/loadsave-cases.txt the way
// tools/loadsave-compat/run.sh runs its C driver, and requires its output,
// line for line: what each section prints, how it ends, and the tree of
// files it leaves. testdata/loadsave-c-output.txt is that output from the
// pinned G3 sources; regenerate it only with run.sh --update.
func TestLoadSaveMatchesC(t *testing.T) {
	cases := readTestdata(t, "loadsave-cases.txt")
	want := readTestdata(t, "loadsave-c-output.txt")

	var sections [][]string
	for line := range strings.SplitSeq(strings.TrimSuffix(cases, "\n"), "\n") {
		switch {
		case line == "" || line[0] == '#':
		case strings.HasPrefix(line, "section "):
			sections = append(sections, []string{line})
		case len(sections) > 0:
			sections[len(sections)-1] = append(sections[len(sections)-1], line)
		}
	}

	// run.sh runs the C driver with umask 022
	old := syscall.Umask(0o022)
	t.Cleanup(func() { syscall.Umask(old) })

	var out strings.Builder
	for _, s := range sections {
		dir := t.TempDir() // each section runs in a fresh directory
		t.Chdir(dir)
		out.WriteString(runLoadSaveSection(t, s))
		out.WriteString(boxfileTree(dir))
		// so that the directory can be removed
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if d != nil && d.IsDir() {
				os.Chmod(path, 0o755)
			}
			return nil
		})
	}

	wantLines := strings.Split(want, "\n")
	gotLines := strings.Split(out.String(), "\n")
	for i := 0; i < len(wantLines) || i < len(gotLines); i++ {
		var w, g string
		if i < len(wantLines) {
			w = wantLines[i]
		}
		if i < len(gotLines) {
			g = gotLines[i]
		}
		if w != g {
			t.Fatalf("line %d:\n got %q\nwant %q", i+1, g, w)
		}
	}
}

// loadSaveLoad is what load_db() runs of the routines tools/loadsave-compat
// checks, in its order, and loadSaveSave what save_db() runs of them.
func (w *world) loadSaveLoad() {
	w.initTradeRoutes()
	w.loadSystem()
	w.checkLoad() // the box-file passes
	w.loadOrders()
	w.checkDb()
	w.determineMapEdges()
	w.createRelics()
	w.initOceanChars()
	w.deleteDeadchars()
}

func (w *world) loadSaveSave() {
	w.cleanupPosts()
	w.saveTradeRoutes()
	w.saveSystem()
	w.checkSave() // write_all_boxes() and write_master()
	w.saveOrders()
	w.renameActJoinFiles()
}

// loadSavePrint is the driver's print(): the globals the routines set, and
// every player's order queues.
func (w *world) loadSavePrint(out io.Writer) {
	fmt.Fprintf(out, "sysclock %d %d %d\n", w.sysclock.turn, w.sysclock.day,
		w.sysclock.daysSinceEpoch)
	fmt.Fprintf(out, "from_host=%s\n", w.fromHost)
	fmt.Fprintf(out, "reply_host=%s\n", w.replyHost)
	fmt.Fprintf(out, "gm_address=%s\n", w.gmAddress)
	fmt.Fprintf(out, "game_title=%s\n", w.gameTitle)
	fmt.Fprintf(out, "game_url=%s\n", w.gameURL)
	fmt.Fprintf(out, "rules_url=%s\n", w.rulesURL)
	fmt.Fprintf(out, "times_url=%s\n", w.timesURL)
	fmt.Fprintf(out, "htpasswd_loc=%s\n", w.htpasswdLoc)
	fmt.Fprintf(out, "garrison_pay=%d auto_quit=%d army_slow_factor=%d\n",
		w.garrisonPay, w.autoQuitTurns, w.armySlowFactor)
	fmt.Fprintf(out, "indep=%d gm=%d skill=%d post=%d init=%d\n",
		w.indepPlayer, w.gmPlayer, w.skillPlayer,
		w.postHasBeenRun, w.seedHasBeenRun)
	fmt.Fprintf(out, "fr=%d tr=%d ur=%d fp=%d hr=%d hp=%d hl=%d\n",
		w.faeryRegion, w.tunnelRegion, w.underRegion,
		w.faeryPlayer, w.hadesRegion, w.hadesPit, w.hadesPlayer)
	fmt.Fprintf(out, "nr=%d nl=%d np=%d cr=%d cp=%d mo=%d\n",
		w.nowhereRegion, w.nowhereLoc, w.npcPl, w.cloudRegion,
		w.combatPl, w.mountOlympus)
	fmt.Fprintf(out, "max_map %d %d %d\n", w.maxMapRow, w.maxMapCol, b2i(w.maxMapInit))
	fmt.Fprintf(out, "trade_route %d %d %d\n", b2i(w.tradeRouteHaveSeed),
		w.tradeRouteSeed, w.tradeRouteSecret)
	fmt.Fprintf(out, "ocean_chars:")
	for _, i := range w.oceanChars {
		fmt.Fprintf(out, " %d", i)
	}
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "trades_to_check:")
	for _, i := range w.tradesToCheck {
		fmt.Fprintf(out, " %d", i)
	}
	fmt.Fprintf(out, "\n")

	for pl := range w.loopKind(tPlayer) {
		p := w.rpPlayer(pl)
		if p == nil {
			continue
		}
		for _, o := range p.orders {
			fmt.Fprintf(out, "orders %d %d:", pl, o.unit)
			for _, l := range o.l {
				fmt.Fprintf(out, " [%s]", l)
			}
			fmt.Fprintf(out, "\n")
		}
	}
}

// runLoadSaveSection runs one section in a fresh world, as one run of the
// C driver, and returns what the driver prints and how the section ended,
// as run.sh writes them. put and putl write the files and are not echoed.
func runLoadSaveSection(t *testing.T, lines []string) (output string) {
	var out strings.Builder
	w := newWorld()
	w.stderr = &out // the driver sends stderr to stdout

	out.WriteString(lines[0] + "\n")
	defer func() {
		fmt.Fprintf(&out, "exit %d\n", checkExit(recover(), &out))
		output = out.String()
	}()

	for _, line := range lines[1:] {
		c := &caseLine{rest: line}
		op := c.word()
		if op == "put" || op == "putl" {
			name := c.word()
			f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
			if err != nil {
				t.Fatal(err)
			}
			f.WriteString(c.text())
			if op == "putl" {
				f.WriteString("\n")
			}
			f.Close()
			continue
		}

		fmt.Fprintf(&out, "> %s\n", line)
		switch op {
		case "mkdir":
			os.Mkdir(c.ctext(), 0o755)
		case "rm":
			syscall.Unlink(c.ctext())
		case "chmod":
			mode, err := strconv.ParseUint(c.word(), 8, 32)
			if err != nil {
				t.Fatalf("bad mode in %q", line)
			}
			os.Chmod(c.ctext(), os.FileMode(mode))
		case "libdir":
			w.libdir = c.ctext()
		case "seed":
			w.tradeRouteSeedPending = true
			w.tradeRouteSeedValue = clib.Strtoull0(c.ctext())
		case "newplayer":
			w.newPlayers.Append(clib.Atoi(c.ctext()))
		case "trade":
			w.initTradeRoutes()
		case "system":
			w.loadSystem()
		case "boxes":
			w.checkLoad()
		case "orders":
			w.loadOrders()
		case "check":
			w.checkDb()
		case "edges":
			w.determineMapEdges()
		case "relics":
			w.createRelics()
		case "ocean":
			w.initOceanChars()
		case "deadchars":
			w.deleteDeadchars()
		case "load":
			w.loadSaveLoad()
		case "posts":
			w.cleanupPosts()
		case "savetrade":
			w.saveTradeRoutes()
		case "savesystem":
			w.saveSystem()
		case "saveboxes":
			w.checkSave()
		case "saveorders":
			w.saveOrders()
		case "rename":
			w.renameActJoinFiles()
		case "save":
			w.loadSaveSave()
		case "print":
			w.loadSavePrint(&out)
		default:
			t.Fatalf("unknown operation %q", line)
		}
	}
	return
}

// TestLoadSaveWorlds is the Go side of tools/loadsave-compat/run.sh's
// check of real saved games. It runs only with G3G0_FULL=1 (see
// requireFull), and only when LOADSAVE_COMPAT_GAMES names a list of games,
// one per line: a directory that holds a copy of the game's lib and an
// output prefix, separated by a tab. For each, from that directory with
// libdir "lib", it runs loadSaveLoad, loadSavePrint and loadSaveSave, and
// writes stdout to PREFIX.out, followed by "exit STATUS", and stderr to
// PREFIX.err. run.sh compares them, and the lib trees, with the C
// driver's.
func TestLoadSaveWorlds(t *testing.T) {
	requireFull(t)
	list := os.Getenv("LOADSAVE_COMPAT_GAMES")
	if list == "" {
		t.Skip("LOADSAVE_COMPAT_GAMES not set; run tools/loadsave-compat/run.sh --full")
	}
	data, err := os.ReadFile(list)
	if err != nil {
		t.Fatal(err)
	}
	for line := range strings.SplitSeq(strings.TrimSuffix(string(data), "\n"), "\n") {
		f := strings.Split(line, "\t")
		if len(f) != 2 {
			t.Fatalf("bad line in %s: %q", list, line)
		}
		t.Chdir(f[0])
		var stdout, stderr strings.Builder
		w := newWorld()
		w.stderr = &stderr

		status := func() (status int) {
			defer func() { status = checkExit(recover(), &stderr) }()
			w.loadSaveLoad()
			w.loadSavePrint(&stdout)
			w.loadSaveSave()
			return
		}()

		fmt.Fprintf(&stdout, "exit %d\n", status)
		if err := os.WriteFile(f[1]+".out", []byte(stdout.String()), 0o666); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(f[1]+".err", []byte(stderr.String()), 0o666); err != nil {
			t.Fatal(err)
		}
	}
}
