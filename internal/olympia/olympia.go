// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package olympia ports G3's olympia-g3, the game engine.
//
// So far it ports main() in olympia/main.c: the command line, what main()
// prints before it dispatches, the -R self-test and the -t ilist test. Every
// other mode stops with a "not yet implemented" message until a later task
// ports it.
//
// The entity model of olympia/oly.h and olympia/glob.c is in entity.go
// (the structures), oly.go (the constants), glob.go (the constant tables)
// and world.go (the global state, as a world, and the p_xxx and rp_xxx
// accessors). code.go ports olympia/code.c's box allocation, kind and
// subkind chains, entity codes and names, loop.go the entity loops of
// olympia/loop.h, and z.go the failed assert(). io.go ports
// olympia/io.c's entity reading and writing: the field helpers, the entity
// sections, load_box and save_box, and the box-file passes and writers
// (fast_scan, scan_all_boxes, read_all_boxes, write_all_boxes and
// write_master). Nothing in main() uses them yet.
package olympia

import (
	"bufio"
	"fmt"
	"io"
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
// same names in olympia/main.c. Nothing reads them yet; the tasks that port
// load_db and the modes will.
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
	// which the task that ports load_db will call; the rest of what it
	// initializes belongs to later tasks.

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

	// Every remaining path starts with load_db().
	return notYetImplemented(stderr, args[0], "load_db"+o.modes())
}

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
	w        *bufio.Writer
	lineBuf  bool
	unbuffer bool
}

func newCStdout(w io.Writer, tty bool) *cStdout {
	return &cStdout{w: bufio.NewWriterSize(w, 1<<16), lineBuf: tty}
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

// Flush writes out anything still buffered, as exit() does.
func (c *cStdout) Flush() {
	c.w.Flush()
}
