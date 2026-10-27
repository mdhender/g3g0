// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/input.c: the command line parser that loading a
// command (scan_command() in io.c) uses, and the loading of a unit's next
// order into its command (init_load_sup() and what it calls), which
// set_lord() uses. Running commands is left for later stages; so is
// reporting an order that does not parse (do_command()).
//
// C works on its strings in place, writing NUL bytes into them. The ports
// work on a []byte copy where C writes into the string, and give each
// string as C then sees it: up to its first NUL byte (cString).

// bat is s[i] as C sees it: 0 (the terminating NUL) past the end of s.
func bat(s []byte, i int) byte {
	if i < len(s) {
		return s[i]
	}
	return 0
}

// parseLine is parse_line(): it clears l and appends to it each argument of
// s, which it cuts up in place, as C does, with a NUL after each argument.
// An argument is a run of bytes that are not white space, or the text
// between a pair of double or single quotes (to the end of s if the quote
// is not closed), without its leading and trailing white space.
//
// C's list holds pointers into s; here each element is the argument as
// that pointer sees it once parse_line() returns.
func parseLine(l cstringsList, s []byte) cstringsList {
	l.Clear()

	i := 0
	for bat(s, i) != 0 {
		for bat(s, i) != 0 && iswhite(bat(s, i)) { // eat whitespace
			i++
		}

		if bat(s, i) == 0 {
			break
		}

		var prev int
		if bat(s, i) == '"' { // handle "double quoted" text
			i++
			prev = i
			for bat(s, i) != 0 && bat(s, i) != '"' {
				i++
			}
		} else if bat(s, i) == '\'' { // handle 'single quoted' text
			i++
			prev = i
			for bat(s, i) != 0 && bat(s, i) != '\'' {
				i++
			}
		} else { // unquoted argument
			prev = i
			for bat(s, i) != 0 && !iswhite(bat(s, i)) {
				i++
			}
		}

		if bat(s, i) != 0 {
			s[i] = 0
			i++
		}

		// Strip leading and trailing whitespace from argument

		for bat(s, prev) != 0 && iswhite(bat(s, prev)) {
			prev++
		}

		p := prev
		for bat(s, p) != 0 {
			p++
		}
		p--

		for p >= prev && iswhite(s[p]) {
			s[p] = 0
			p--
		}

		l.Append(cString(s[prev:]))
	}

	return l
}

// removeComment is remove_comment(): it ends s at its first '#' that is
// not inside quoted text, in place.
//
// As in parse_line(), a double or single quote starts quoted text only at
// the start of a token: at the start of s, after white space, or right
// after a closing quote. A '#' inside the quotes is kept, an unterminated
// quote runs to the end of s, and a quote inside a word (Joe's) is
// ordinary text. (Fixed upstream in
// https://github.com/playbymail/olympia-g3/issues/315; before, C's quote
// loops stopped at once on the opening quote, so a '#' between quotes
// ended the line too.)
func removeComment(s []byte) {
	tok := true // s[i] is at the start of a token
	i := 0
	for bat(s, i) != 0 {
		if c := s[i]; tok && (c == '"' || c == '\'') {
			q := c
			i++
			for bat(s, i) != 0 && bat(s, i) != q {
				i++
			}
			if bat(s, i) == 0 { // unterminated: runs to the end
				break
			}
			i++ // parse_line() splits after the quote
			continue
		}

		if s[i] == '#' {
			s[i] = 0
			break
		}

		tok = iswhite(s[i])
		i++
	}
}

// removeCtrlChars is remove_ctrl_chars(): in place, each byte of s from
// 0x80 up loses its high bit, and then each byte below 32 becomes a space,
// so a byte of 0x80 becomes a space, not the end of s.
func removeCtrlChars(s []byte) {
	for i := 0; bat(s, i) != 0; i++ {
		if s[i]&0x80 != 0 {
			s[i] &= 0x7F
		}
		if s[i] < 32 {
			s[i] = ' '
		}
	}
}

// findCommand is find_command(): the index in cmdTbl of the command s
// names, matched first exactly and then fuzzily (see fuzzyStrcmp), both
// without regard to case, or -1. It sets fuzzyFind if the match is fuzzy.
func (w *world) findCommand(s string) int {
	w.fuzzyFind = false

	if s == "" {
		return -1
	}

	for i := 1; i < len(cmdTbl); i++ {
		if iStrcmp(cmdTbl[i].name, s) == 0 {
			return i
		}
	}

	for i := 1; i < len(cmdTbl); i++ {
		if fuzzyStrcmp(cmdTbl[i].name, s) {
			w.fuzzyFind = true
			return i
		}
	}

	return -1
}

// olyParseCmd is oly_parse_cmd(): it parses the command line s into c. It
// drops s's leading white space and comment, makes its control bytes
// spaces and saves it as c's line; takes a leading '&' or '?' as the
// command's condition; and cuts the rest into c's parse list, whose first
// word names the command. It returns false, leaving c's cmd 0, if that
// word names no command. A line with no words is not an error. It clears
// c's cmd, fuzzy and use_skill first.
//
// C drops the comment and control bytes in place, in the caller's string,
// after its leading white space; here they are dropped from *s.
func (w *world) olyParseCmd(c *command, s *string) bool {
	c.cmd = 0
	c.fuzzy = 0
	c.useSkill = 0

	i := 0
	for iswhite(at(*s, i)) {
		i++
	}

	b := []byte((*s)[i:])
	removeComment(b)
	b = []byte(cString(b))
	removeCtrlChars(b)

	line := string(b)
	c.line = &line
	*s = (*s)[:i] + line

	switch at(line, 0) {
	case '&':
		c.conditional = 1
		b = b[1:]
	case '?':
		c.conditional = 2
		b = b[1:]
	default:
		c.conditional = 0
	}

	parsed := []byte(string(b))
	c.parse = parseLine(c.parse, parsed)
	pl := cString(parsed)
	c.parsedLine = &pl

	if c.parse.Len() > 0 {
		i := w.findCommand(c.parse[0])

		if i > 0 {
			c.cmd = int32(i)
			if w.fuzzyFind {
				c.fuzzy = 1
			}
		} else {
			return false
		}
	}

	return true
}

// maxPri is MAX_PRI in input.c: the number of command priorities, and of
// load queues.
const maxPri = 5

// parseArg is parse_arg(): the entity or number s names (see scode), or 0
// if it is negative. If it is 0 and s is a near miss for "garrison" (see
// fuzzyStrcmp), it is the garrison in who's location, or garrison_magic if
// there is none.
func (w *world) parseArg(who int, s string) int {
	n := scode(s)

	if n < 0 {
		n = 0
	}

	if n == 0 && who != 0 && w.subloc(who) != 0 && fuzzyStrcmp(s, "garrison") {
		n = w.garrisonHere(w.subloc(who))
		if n == 0 {
			n = int(w.garrisonMagic)
		}
	}

	return n
}

// olyParse is oly_parse(): olyParseCmd, then the first eight arguments
// after the command name, each through parseArg, into c's a to h, which it
// clears first.
func (w *world) olyParse(c *command, s *string) bool {
	c.a = 0
	c.b = 0
	c.c = 0
	c.d = 0
	c.e = 0
	c.f = 0
	c.g = 0
	c.h = 0

	if !w.olyParseCmd(c, s) {
		return false
	}

	args := []*int32{&c.a, &c.b, &c.c, &c.d, &c.e, &c.f, &c.g, &c.h}
	for i := min(c.parse.Len(), 9) - 1; i >= 1; i-- {
		*args[i-1] = int32(w.parseArg(int(c.who), c.parse[i]))
	}

	return true
}

// getCommand is get_command(): it takes the next order queued for unit
// who by the player who controls it now, and returns it and true, or false
// if there is none or who is not a valid, living unit. Taking an order
// sets the player's last_order_turn to this turn.
func (w *world) getCommand(who int) (string, bool) {
	if !w.validBox(who) || w.kind(who) == tDeadchar {
		return "", false
	}

	fact := w.player(who) // who controls us now?

	if fact == 0 {
		return "", false
	}

	order, ok := w.topOrder(fact, who)

	// If we don't have any orders, then fail.

	if !ok {
		return "", false
	}

	// Update the player's last turn field so we know if he's been playing
	// or not.

	w.pPlayer(fact).lastOrderTurn = int32(w.sysclock.turn)

	assert(len(order) < lenMax, "olympia/input.c", 403, "strlen(order) < LEN") // maybe not a valid assert

	w.popOrder(fact, who)

	return order, true
}

// setState is set_state(): it moves c's unit from the queue of its state
// (the run queue, or the load queue of its priority) to the queue of the
// new state (the load queue of newPri for STATE_LOAD). STATE_DONE and
// STATE_ERROR have no queue. The state is stored as a schar.
func (w *world) setState(c *command, state, newPri int) {
	switch c.state {
	case stateRun:
		w.runQ.RemValueUniq(int(c.who))

	case stateLoad:
		w.loadQ[c.pri].RemValueUniq(int(c.who))
	}

	c.state = int8(state)

	switch c.state {
	case stateRun:
		w.runQ.Append(int(c.who))

	case stateLoad:
		w.loadQ[newPri].Append(int(c.who))
	}
}

// loadCommand is load_command(): it loads c with its unit's next order
// and queues it to run (STATE_LOAD), or, with no order left, sets it to
// STATE_DONE. It returns whether there was an order.
//
// An order that does not parse is set to STATE_ERROR and passed to
// do_command(), which reports it and loads the next order; that is not
// ported yet, and panics.
func (w *world) loadCommand(c *command) bool {
	assert(c != nil, "olympia/input.c", 537, "c != NULL")

	buf, ok := w.getCommand(int(c.who))
	if !ok {
		w.setState(c, stateDone, 0)
		return false
	}

	if w.olyParse(c, &buf) {
		pri := int(cmdTbl[c.cmd].pri)

		w.setState(c, stateLoad, pri)

		c.pri = int8(pri)
		c.wait = cmdTbl[c.cmd].time
		c.poll = int8(cmdTbl[c.cmd].poll)
		c.daysExecuting = 0
	} else {
		// No queue holds STATE_ERROR, so report the bad order now.
		// do_command fails it, and command_done loads the next order.
		w.setState(c, stateError, 0)
		panic("olympia: load_command: do_command is not yet implemented")
	}

	return true
}

// initLoadSup is init_load_sup(): it puts who's command on the queue for
// its state, or, if it is done, loads the next order into it. A unit with
// no command gets one, done. A command must belong to its unit, and be
// loading, running or done.
func (w *world) initLoadSup(who int) {
	c := w.rpCommand(who)

	// All characters should have a command structure. Create one if they
	// don't.

	if c == nil {
		c = w.pCommand(who)
		c.who = int32(who)
		c.state = stateDone
	}

	assert(who == int(c.who), "olympia/input.c", 722, "who == c->who")

	switch c.state {
	case stateLoad:
		w.loadQ[c.pri].Append(int(c.who))

	case stateRun:
		w.runQ.Append(int(c.who))

	case stateDone:
		w.loadCommand(c)

	default:
		assert(false, "olympia/input.c", 739, "FALSE")
	}
}
