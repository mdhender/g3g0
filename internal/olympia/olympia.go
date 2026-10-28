// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package olympia ports G3's olympia-g3, the game engine.
//
// So far it ports main() in olympia/main.c: the command line, what main()
// prints before it dispatches, the -R self-test, the -t ilist test, and a
// run that loads an existing game (load_db), runs a turn with -r, checks it,
// saves it with -S (save_db) and writes the Times masthead. Report paths
// preserve state and diagnostics; report artifacts remain stage 2. Other
// modes and unported gameplay stop with "not yet implemented", as does a
// load that has to build a new world.
//
// The entity model of olympia/oly.h and olympia/glob.c is in entity.go
// (the structures), oly.go (the constants), glob.go (the constant tables)
// and world.go (the global state, as a world, and the p_xxx and rp_xxx
// accessors). code.go ports olympia/code.c's box allocation, kind and
// subkind chains, entity codes and names, loop.go the entity loops of
// olympia/loop.h, and z.go the failed assert(). io.go ports
// olympia/io.c's entity reading and writing: the field helpers, the entity
// sections, load_box and save_box, the box-file passes and writers
// (fast_scan, scan_all_boxes, read_all_boxes, write_all_boxes and
// write_master), and load_db and save_db.
package olympia

import (
	"bufio"
	"fmt"
	"io"
	"runtime"
	"strings"
	"unsafe"

	"github.com/mdhender/g3g0/internal/clib"
	"github.com/mdhender/g3g0/internal/lists"
)

// usage is the text main() prints to stderr when errflag is set.
const usage = `usage: oly [options]
	-w		Windows mode
	-a		Add new players mode
	-e		Eat orders from libdir/spool
	-E		Eat orders from libdir/spool and terminate
	-f		Don't buffer files for debugging
	-i		Immediate mode
	-l dir		Specify libdir, default ./lib
	-p		Don't make data files pretty
	-r		Run a turn
	-R		Test the random number generator
	-s		Extract safe havens as start city list
	-S		Save the database at completion
	-t		Test ilist code
	-T		Print timing info
	-M		Mail reports and order acks
	-A		Charge player accounts
	-G seed		Set per-game trade-route seed (game creation only)
	--coverage path	Write the command coverage dump to path
	--use-coverage path	Write the use coverage dump to path
	--state-export path	Write the faction state export to path
`

// options holds the flags main() sets, as the globals and locals of the
// same names in olympia/main.c. runEngine copies those that load_db(),
// save_db() and the end of main() read into the world; the others belong to
// modes that are not ported yet.
type options struct {
	libdir string // -l; the C default is "lib" (libdir in glob.c)

	winFlag         bool // -w
	addFlag         bool // -a
	acctFlag        bool // -A
	flushAlways     bool // -f
	eatFlag         bool // -e, -E
	eatFast         bool // -E
	immediate       bool
	immedAfter      bool // -i
	prettyDataFiles bool // -p toggles it
	runFlag         bool // -r
	startlocFlag    bool // -s
	saveFlag        bool // -S
	timeSelf        bool // -T
	mailNow         bool // -M
	htmlFlag        bool // -h

	tradeRouteSeedPending bool   // -G
	tradeRouteSeedValue   uint64 // -G, as strtoull(optarg, NULL, 0)

	testUseConstReportDate bool   // test-use-const-report-date
	coveragePath           string // --coverage
	useCoveragePath        string // --use-coverage
	phaseCoveragePath      string // --phase-coverage
	stateExportPath        string // --state-export
}

// Run is main() in olympia/main.c. args is argv, including argv[0], and
// stdoutIsTTY is isatty(1). It returns the process exit status olympia-g3
// would have.
func Run(args []string, stdout, stderr io.Writer, stdoutIsTTY bool) int {
	out := newCStdout(stdout, stdoutIsTTY)
	defer out.Flush()

	// sizeof(int) != sizeof(int *): true on every 64-bit build. The C code
	// prints the warning and carries on.
	if unsafe.Sizeof(int32(0)) != unsafe.Sizeof(uintptr(0)) {
		out.Printf("The Olympia C code is not 64-bit clean.\n")
		out.Printf("at the least: it puts pointers into ilists\n")
	}

	out.Printf("\tsizeof(struct box) = %d\n", sizeofBox)

	// call_init_routines() prints nothing. Its glob_init() is newWorld(),
	// which runEngine calls; nothing seen before then depends on it. The
	// rest of what it initializes belongs to later tasks.

	o := &options{libdir: "lib"}
	errflag := 0

	// Pull the long-form arguments out of argv before getopt() sees them,
	// compacting argv so the remaining options parse normally. An option
	// whose path is missing or empty is an error.
	argv := []string{args[0]}
	for ai := 1; ai < len(args); ai++ {
		var path *string
		switch args[ai] {
		case "test-use-const-report-date":
			o.testUseConstReportDate = true
			continue
		case "--coverage":
			path = &o.coveragePath
		case "--use-coverage":
			path = &o.useCoveragePath
		case "--phase-coverage":
			path = &o.phaseCoveragePath
		case "--state-export":
			path = &o.stateExportPath
		default:
			argv = append(argv, args[ai])
			continue
		}
		if ai+1 < len(args) && args[ai+1] != "" {
			ai++
			*path = args[ai]
		} else {
			errflag++
		}
	}

	opts := clib.NewGetopt(argv, "waeEfirl:pR?sStMTAhG:", stderr)
	for c := opts.Next(); c != -1; c = opts.Next() {
		switch c {
		case 'w':
			o.winFlag = true
		case 'a':
			o.addFlag = true
			o.immediate = false
		case 'A':
			o.acctFlag = true
		case 'f':
			o.flushAlways = true
			out.Unbuffer()
		case 'e':
			o.eatFlag = true
			o.eatFast = false
			o.immediate = false
		case 'E':
			o.eatFlag = true
			o.eatFast = true
			o.immediate = false
		case 'i':
			o.immedAfter = true
		case 'l':
			o.libdir = opts.Optarg
		case 'p':
			o.prettyDataFiles = !o.prettyDataFiles
		case 'r':
			o.immediate = false
			o.runFlag = true
		case 'R':
			testRandom(out, stdoutIsTTY)
			return 0
		case 's':
			o.startlocFlag = true
		case 'S':
			o.saveFlag = true
		case 't':
			// ilist_test() in lib/ilist.c starts with setbuf(stdout, NULL).
			out.Unbuffer()
			lists.ILTest(out.Printf)
			return 0
		case 'T':
			o.timeSelf = true
		case 'M':
			o.mailNow = true
		case 'h':
			o.htmlFlag = true
		case 'G':
			o.tradeRouteSeedPending = true
			o.tradeRouteSeedValue = clib.Strtoull0(opts.Optarg)
		default:
			// Includes '?', which is in the option string, so -? is an
			// error without a getopt() diagnostic.
			errflag++
		}
	}

	if errflag != 0 {
		fmt.Fprint(stderr, usage)
		return 1
	}

	// Every remaining path starts with load_db(). The modes that run after
	// it are not ported yet, so a run that asks for one stops before it.
	if o.unported() {
		return notYetImplemented(stderr, args[0], "load_db"+o.modes())
	}

	return runEngine(o, args[0], out, stderr)
}

// unported reports whether o asks for a mode that runs after load_db() and
// is not ported yet. Of the flags that are left, -w and -A change only
// what those modes do, and -f, -p, -T, -G and test-use-const-report-date
// are read by load_db(), save_db() and what main() runs after them.
func (o *options) unported() bool {
	return o.eatFlag || o.startlocFlag || o.addFlag ||
		o.immedAfter || o.mailNow || o.htmlFlag ||
		o.coveragePath != "" || o.useCoveragePath != "" ||
		o.phaseCoveragePath != "" || o.stateExportPath != ""
}

// runEngine is main() from load_db() on. A turn runs until the first real
// unfinished phase; -S must not save a partially executed turn.
//
// It recovers the ways the engine ends a C process early. A failed
// assert() (an assertion) prints its message and calls abort(), which on
// macOS flushes every stdio stream (stdout and the files clib.FlushAll
// flushes) and exits 134. exit() (an exitStatus) flushes them too. A NULL
// dereference, which Go reaches as a nil pointer dereference, kills C with
// SIGSEGV: nothing is flushed, and the status is 139. A step that is not
// ported (a notImplemented) stops with the "not yet implemented" message.
// Other panics are not recovered; normally they indicate a bug in the port.
//
// C's integer division by zero would die with SIGFPE (136) on x86-64 and
// give 0 on arm64. A malformed loaded zero-price BUY can reach that division
// in reduceQty during market reporting (G3 #358). Go's divide-by-zero panic
// remains unrecovered; host-dependent C failure parity is not implemented.
func runEngine(o *options, prog string, out *cStdout, stderr io.Writer) (status int) {
	// call_init_routines(): glob_init()
	w := newWorld()
	w.stdout = out
	w.stderr = stderr
	w.libdir = o.libdir
	w.prettyDataFiles = o.prettyDataFiles
	w.flushAlways = o.flushAlways
	w.winFlag = o.winFlag
	w.timeSelf = o.timeSelf
	w.tradeRouteSeedPending = o.tradeRouteSeedPending
	w.tradeRouteSeedValue = o.tradeRouteSeedValue
	w.testUseConstReportDate = o.testUseConstReportDate

	defer func() {
		// Process termination closes open stdio files, including at an NYI.
		if w.pressFP != nil {
			defer w.pressFP.Close()
		}
		if w.rumorFP != nil {
			defer w.rumorFP.Close()
		}
		switch r := recover().(type) {
		case nil:
		case assertion:
			fmt.Fprint(stderr, r.Error())
			clib.FlushAll()
			status = 134 // 128 + SIGABRT
		case exitStatus:
			clib.FlushAll()
			status = int(r)
		case notImplemented:
			clib.FlushAll()
			status = notYetImplemented(stderr, prog, string(r))
		case runtime.Error:
			if !strings.Contains(r.Error(), "nil pointer dereference") {
				panic(r)
			}
			out.Discard()
			status = 139 // 128 + SIGSEGV
		default:
			panic(r)
		}
	}()

	w.loadDb()

	// -e, -E, -s, -a and -i, and the coverage dumps, stop before
	// load_db() (unported). immediate is never set without them.
	if o.runFlag {
		w.openLogfile(o.winFlag)
		w.openTimes()
		w.showDay = true
		w.processOrders()
		w.postMonth()
		w.showDay = false
		// main.c:981–1008, in C order. Remaining projections defer report
		// artifacts only, not persistent state, input reads or diagnostics.
		w.determineOutputOrder()
		w.turnEndLocReports()
		w.listOrderTemplates()
		w.playerEntInfo()
		w.characterReport()
		w.playerBanner()
		if o.acctFlag {
			panic(notImplemented("charge_account"))
		}
		w.reportAccount()
		w.summaryReport()
		w.playerReport()
		w.scanCharSkillLore()
		w.showLoreSheets()
		w.gmReportState(int(w.gmPlayer))
		w.gmShowAllSkills(int(w.skillPlayer))
		w.addNewPlayers()
		w.genIncludeSectionState() // must be last, before closing output
		w.closeLogfile()
		w.writePlayerListState()
		w.writeEmailState()
		w.writeTotimesState()
		w.writeForwardsState()
		w.writeFactionsState()
	}

	w.checkDb() // check database integrity

	if o.saveFlag {
		w.saveDb()
	}
	if o.saveFlag && o.runFlag {
		w.saveLogdir(o.winFlag)
	}

	// write_state_export() and -M and -h stop before load_db() (unported).

	w.timesMasthead()
	w.closeTimes()
	w.stageSup(nil)

	// malloc_size and realloc_size in z.c are never counted.
	out.Printf("\tmalloc, realloc = %d, %d\n", 0, 0)

	// return from main(): exit() flushes every open stream
	clib.FlushAll()
	return 0
}

// The main.c:50–290 post writers retain real fopen("w")/fclose decisions
// and nonfatal diagnostics. Artifact bodies alone belong to stage 2. In
// particular, forward ownership validation must happen AFTER a successful
// open, not as a projection run before (or instead of) the real IO.
func (w *world) openPostWriter(name string) *clib.File {
	fnam := sout("%s/%s", w.libdir, name)
	fp, err := clib.Fopen(fnam)
	if err != nil {
		fmt.Fprintf(w.stderr, "can't write %s", fnam)
		clib.Perror(w.stderr, "", err)
		return nil
	}
	return fp
}

// writeTotimesState and writeEmailState have no persistent body effects:
// their rp_player/email/compuserve reads only select deferred artifact text.
func (w *world) writeTotimesState() {
	if fp := w.openPostWriter("totimes"); fp != nil {
		fp.Close()
	}
}

func (w *world) writeEmailState() {
	if fp := w.openPostWriter("email"); fp != nil {
		fp.Close()
	}
}

func (w *world) writePlayerListState() {
	w.stage("write_player_list()")
	fp := w.openPostWriter("players")
	if fp == nil {
		return
	}
	for pl := range w.loopKind(tPlayer) {
		if p := w.rpPlayer(pl); p != nil && p.email != nil && w.subkind(pl) == subPlRegular {
			// list_a_player's allocation and name access, not its text.
			w.pPlayer(pl)
			w.justName(pl)
		}
	}
	fp.Close()
}

// writeForwardSupState covers both write_forward_sup and write_faction_sup:
// player() follows/validates ownership even if no email will be rendered.
func (w *world) writeForwardSupState(whoFor int) {
	pl := w.player(whoFor)
	w.rpPlayer(pl) // player_email: preserve its raw box access, including 0.
}

func (w *world) writeForwardsState() {
	fp := w.openPostWriter("forward")
	if fp == nil {
		return
	}
	for i := range w.loopKind(tPlayer) {
		w.writeForwardSupState(i)
	}
	for i := range w.loopKind(tChar) {
		w.writeForwardSupState(i)
	}
	for i := range w.loopSubkind(subGarrison) {
		for _, pl := range w.playersWhoRuleHere(i) {
			if pl != 0 {
				w.writeForwardSupState(pl)
			}
		}
	}
	fp.Close()
}

func (w *world) writeFactionsState() {
	fp := w.openPostWriter("factions")
	if fp == nil {
		return
	}
	for i := range w.loopKind(tPlayer) {
		w.writeForwardSupState(i)
	}
	fp.Close()
}

// notImplemented stops the engine at a step of a mode it runs that is not
// ported yet; runEngine recovers it.
type notImplemented string

// exitStatus is exit(status) called from inside the engine; runEngine
// recovers it.
type exitStatus int

// modes lists the requested modes, for the "not yet implemented" message.
func (o *options) modes() string {
	var m []string
	for _, f := range []struct {
		set  bool
		name string
	}{
		{o.eatFlag && !o.eatFast, "-e"},
		{o.eatFlag && o.eatFast, "-E"},
		{o.startlocFlag, "-s"},
		{o.runFlag, "-r"},
		{o.addFlag, "-a"},
		{o.immedAfter, "-i"},
		{o.saveFlag, "-S"},
		{o.mailNow, "-M"},
		{o.acctFlag, "-A"},
		{o.htmlFlag, "-h"},
		{o.timeSelf, "-T"},
		{o.tradeRouteSeedPending, "-G"},
		{o.coveragePath != "", "--coverage"},
		{o.useCoveragePath != "", "--use-coverage"},
		{o.phaseCoveragePath != "", "--phase-coverage"},
		{o.stateExportPath != "", "--state-export"},
	} {
		if f.set {
			m = append(m, f.name)
		}
	}
	if len(m) == 0 {
		return ""
	}
	return " (" + strings.Join(m, " ") + ")"
}

// notYetImplemented stops a mode that has not been ported.
func notYetImplemented(stderr io.Writer, prog, what string) int {
	fmt.Fprintf(stderr, "%s: %s: not yet implemented\n", prog, what)
	return 2
}

// cStdout buffers stdout the way C stdio does, so that stdout and stderr
// interleave as they do in C when both go to the same file: stdout is line
// buffered on a terminal and otherwise held until exit, unless -f turns
// buffering off. (C flushes a full buffer early, at a size that depends on
// the file; nothing ported so far writes enough to stdout for that to show
// before something is written to stderr.)
type cStdout struct {
	dst      io.Writer
	w        *bufio.Writer
	lineBuf  bool
	unbuffer bool
}

func newCStdout(w io.Writer, tty bool) *cStdout {
	return &cStdout{dst: w, w: bufio.NewWriterSize(w, 1<<16), lineBuf: tty}
}

// Printf is printf().
func (c *cStdout) Printf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	c.w.WriteString(s)
	if c.unbuffer || (c.lineBuf && strings.Contains(s, "\n")) {
		c.w.Flush()
	}
}

// Unbuffer is setbuf(stdout, NULL), which flushes what is buffered first.
func (c *cStdout) Unbuffer() {
	c.w.Flush()
	c.unbuffer = true
}

// Write is fwrite() to stdout, buffered as Printf buffers it.
func (c *cStdout) Write(p []byte) (int, error) {
	c.Printf("%s", p)
	return len(p), nil
}

// Discard drops what is buffered, as a process killed by a signal loses
// it.
func (c *cStdout) Discard() {
	c.w.Reset(c.dst)
}

// Flush writes out anything still buffered, as exit() does.
func (c *cStdout) Flush() {
	c.w.Flush()
}
