// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package genesis

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// This file reproduces getopt_long(3) from the BSD libc on macOS, the
// reference platform, for the one option table genesis uses. Its error
// messages differ from glibc's. Like the C library it permutes: options
// after a non-option argument are still parsed, unless POSIXLY_CORRECT is
// set. genesis ignores non-option arguments, so they are only skipped here.

const (
	noArgument       = 0
	requiredArgument = 1

	badCh  = '?' // unknown option
	badArg = '?' // missing or unexpected argument; ':' only with a leading ':' in optstring
)

// longOption is struct option without the flag pointer, which genesis
// never sets.
type longOption struct {
	name   string
	hasArg int
	val    int
}

// getopt holds the state that getopt_long keeps in globals.
type getopt struct {
	args      []string
	optstring string
	longopts  []longOption
	stderr    io.Writer

	optind int
	optarg string
	place  string // rest of the current short-option cluster

	posixlyCorrect bool
}

func newGetopt(args []string, optstring string, longopts []longOption, stderr io.Writer) *getopt {
	_, posixly := os.LookupEnv("POSIXLY_CORRECT")
	return &getopt{
		args:           args,
		optstring:      optstring,
		longopts:       longopts,
		stderr:         stderr,
		optind:         1,
		posixlyCorrect: posixly,
	}
}

// warnx prints a diagnostic prefixed with the program name, which is the
// last element of argv[0].
func (g *getopt) warnx(format string, a ...any) {
	prog := ""
	if len(g.args) > 0 {
		prog = g.args[0]
		if i := strings.LastIndexByte(prog, '/'); i >= 0 {
			prog = prog[i+1:]
		}
	}
	fmt.Fprintf(g.stderr, "%s: %s\n", prog, fmt.Sprintf(format, a...))
}

// next returns the next option character or long option value, or -1 when
// the options are exhausted. The option's argument is left in g.optarg.
func (g *getopt) next() int {
	g.optarg = ""
	if g.place == "" {
		for {
			if g.optind >= len(g.args) {
				return -1
			}
			arg := g.args[g.optind]
			if arg == "" || arg[0] != '-' || arg == "-" {
				// A non-option argument.
				if g.posixlyCorrect {
					return -1
				}
				g.optind++
				continue
			}
			if arg == "--" {
				g.optind++
				return -1
			}
			g.place = arg[1:]
			break
		}
	}

	// A '-' after the leading dash, at the start of an argument or inside a
	// cluster such as "-D-size", introduces a long option.
	if g.place[0] == '-' {
		current := g.place[1:]
		g.place = ""
		if current == "" {
			// "-D-": malformed, and reported without a message.
			return badArg
		}
		return g.parseLong(current)
	}

	optchar := g.place[0]
	g.place = g.place[1:]
	oli := strings.IndexByte(g.optstring, optchar)
	if optchar == ':' || oli < 0 {
		if g.place == "" {
			g.optind++
		}
		if g.posixlyCorrect {
			g.warnx("illegal option -- %c", optchar)
		} else {
			g.warnx("invalid option -- %c", optchar)
		}
		return badCh
	}
	if oli+1 >= len(g.optstring) || g.optstring[oli+1] != ':' {
		// Takes no argument.
		if g.place == "" {
			g.optind++
		}
		return int(optchar)
	}
	if g.place != "" {
		g.optarg = g.place
	} else {
		g.optind++
		if g.optind >= len(g.args) {
			g.warnx("option requires an argument -- %c", optchar)
			return badArg
		}
		g.optarg = g.args[g.optind]
	}
	g.place = ""
	g.optind++
	return int(optchar)
}

// parseLong handles a long option; current is its text after "--".
func (g *getopt) parseLong(current string) int {
	g.optind++

	name, value, hasEqual := strings.Cut(current, "=")
	match := -1
	exact := false
	ambiguous := false
	for i, o := range g.longopts {
		if !strings.HasPrefix(o.name, name) {
			continue
		}
		if len(o.name) == len(name) {
			match = i
			exact = true
			break
		}
		if match == -1 {
			match = i
		} else if o.hasArg != g.longopts[match].hasArg || o.val != g.longopts[match].val {
			ambiguous = true
		}
	}
	if !exact && ambiguous {
		g.warnx("option `--%s' is ambiguous", name)
		return badCh
	}
	if match == -1 {
		g.warnx("unrecognized option `--%s'", current)
		return badCh
	}

	o := g.longopts[match]
	if o.hasArg == noArgument && hasEqual {
		g.warnx("option `--%s' doesn't allow an argument", name)
		return badArg
	}
	if o.hasArg == requiredArgument {
		switch {
		case hasEqual:
			g.optarg = value
		case g.optind < len(g.args):
			g.optarg = g.args[g.optind]
			g.optind++
		default:
			g.warnx("option `--%s' requires an argument", current)
			return badArg
		}
	}
	return o.val
}
