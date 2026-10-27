// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"slices"
	"strings"

	"github.com/mdhender/g3g0/internal/clib"
)

// Ported from olympia/input.c: the command line parser that loading a
// command (scan_command() in io.c) uses, and the loading of a unit's next
// order into its command (init_load_sup() and what it calls), which
// set_lord() uses, plus individual command dispatch, completion and wait polling.
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
// An order that does not parse is reported immediately; completion loads
// the next order, so this returns whether that order is ready to load.
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
		w.doCommand(c)
		return c.state == stateLoad
	}

	return true
}

// cmdShift deletes argument 1, shifts the parsed numeric arguments and
// parses the newly exposed eighth argument (not necessarily zero).
func (w *world) cmdShift(c *command) {
	if c.parse.Len() > 1 {
		c.parse = append(c.parse[:1], c.parse[2:]...)
	}
	c.a, c.b, c.c, c.d, c.e, c.f, c.g = c.b, c.c, c.d, c.e, c.f, c.g, c.h
	if c.parse.Len()-1 >= 8 {
		c.h = int32(w.parseArg(int(c.who), c.parse[8]))
	} else {
		c.h = 0
	}
}

// checkAllow is check_allow. nil is C NULL (unrestricted); a pointer to
// the empty string is not NULL. The actual table contains empty strings.
func (w *world) checkAllow(c *command, allow *string) bool {
	if allow == nil || (w.immediate && strings.ContainsRune(*allow, 'i')) {
		return true
	}
	var t byte
	switch w.kind(int(c.who)) {
	case tPlayer:
		t = 'p'
	case tChar:
		if p := w.rpMisc(int(c.who)); p != nil {
			t = byte(p.cmdAllow)
		}
		if t == 0 {
			t = 'c'
		}
	default:
		fmt.Fprintf(w.stderr, "check_allow: bad kind: %s\n", w.boxName(int(c.who)))
		assert(false, "olympia/input.c", 352, "FALSE")
	}
	if strings.ContainsRune(*allow, 'm') && w.player(int(c.who)) == int(w.gmPlayer) {
		return true
	}
	if !strings.ContainsRune(*allow, rune(t)) {
		w.commandOut(int(c.who), "%s may not issue that order.", w.boxName(int(c.who)))
		return false
	}
	return true
}

func (w *world) commandDone(c *command) {
	if w.immediate {
		w.setState(c, stateDone, 0)
		return
	}
	if w.loadCommand(c) && int32(c.pri) < w.curPri {
		w.curPri = int32(c.pri)
	}
}

func (w *world) finishCommand(c *command) int {
	assert(c != nil, "olympia/input.c", 594, "c != NULL")
	assert(w.validBox(int(c.who)), "olympia/input.c", 595, "valid_box(c->who)")
	if w.kind(int(c.who)) == tDeadchar {
		w.commandDone(c)
		return 0
	}
	if w.charMoving(int(c.who)) != 0 && w.stackLeader(int(c.who)) != int(c.who) && c.cmd != w.cmdWait {
		return 1
	}
	if c.wait > 0 {
		c.wait--
	}
	if c.wait > 0 && int16(c.debug) == w.sysclock.day {
		w.commandOut(int(c.who), "finish_command called twice, wait=%d", c.wait)
	}
	c.debug = int8(w.sysclock.day)
	if c.wait <= 0 || c.poll != 0 {
		if finish := cmdTbl[c.cmd].finish; finish != nil && c.inhibitFinish == 0 {
			c.status = int8(finish(w, c))
		}
	}
	if c.state == stateRun && (c.status == 0 || c.wait == 0) {
		w.commandDone(c)
	}
	return int(c.status)
}

func (w *world) doCommand(c *command) {
	assert(c != nil, "olympia/input.c", 644, "c != NULL")
	if !w.immediate {
		w.commandOut(int(c.who), "> %s", clib.Str(c.line))
		if c.fuzzy != 0 {
			w.commandOut(int(c.who), "(assuming you meant '%s')", cmdTbl[c.cmd].name)
		}
	}
	if c.state == stateError {
		w.commandOut(int(c.who), "Unrecognized command.")
		c.status = 0
	} else if !w.checkAllow(c, &cmdTbl[c.cmd].allow) {
		c.status = 0
	} else if cmdTbl[c.cmd].start == nil {
		w.commandOut(int(c.who), "Unimplemented command.")
		c.status = 0
	} else {
		w.pPlayer(w.player(int(c.who))).cmdCount++
		w.setState(c, stateRun, 0)
		assert(c.daysExecuting == 0, "olympia/input.c", 680, "c->days_executing == 0")
		c.debug = 0
		c.inhibitFinish = 0
		cmd := c.cmd // start may re-parse c
		c.status = int8(cmdTbl[cmd].start(w, c))
	}
	// As in the pinned C, conditional is parsed and saved but not consulted
	// here: neither '&' nor '?' suppresses execution based on old status.
	// coverage_count is the stage-3 instrumentation boundary, not game state.
	if c.status == 0 {
		w.commandDone(c)
	} else if c.wait == 0 && c.state == stateRun {
		c.status = int8(w.finishCommand(c))
	}
	if w.tradesToCheck.Len() > 0 {
		panic(notImplemented("check_validated_trades"))
	}
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

func (w *world) execPrecedence(who int) int32 {
	if w.kind(who) != tChar {
		return 0
	}
	var depth int32
	for n := who; n != 0; n = w.stackParent(n) {
		depth++
	}
	pos := w.herePos(who)
	assert(depth < 100, "olympia/input.c", 429, "stack_depth < 100")
	assert(pos < 10000, "olympia/input.c", 430, "pos < 10000")
	return depth*10000 + int32(pos)
}

func (w *world) execComp(a, b int) int {
	if w.bx[a].temp != w.bx[b].temp {
		return int(w.bx[a].temp - w.bx[b].temp)
	}
	return a - b
}

func (w *world) sortLoadQueue(l []int) {
	for _, i := range l {
		w.bx[i].temp = w.execPrecedence(i)
	}
	slices.SortFunc(l, w.execComp)
}

func (w *world) sortRunQueue(l []int) {
	for _, i := range l {
		c := w.rpCommand(i)
		assert(c != nil, "olympia/input.c", 476, "c")
		w.bx[i].temp = int32(c.pri)*1000000 + w.execPrecedence(i)
	}
	slices.SortFunc(l, w.execComp)
}

func (w *world) initialCommandLoad() {
	for i := range w.loopKind(tChar) {
		w.initLoadSup(i)
	}
	for i := range w.loopKind(tPlayer) {
		w.initLoadSup(i)
	}
}

func (w *world) minPriReady() int {
	for pri := 0; pri < maxPri; pri++ {
		for _, i := range w.loadQ[pri] {
			c := w.rpCommand(i)
			assert(c != nil, "olympia/input.c", 776, "c")
			assert(c.state == stateLoad, "olympia/input.c", 777, "c->state == STATE_LOAD")
			assert(int(c.pri) == pri, "olympia/input.c", 778, "c->pri == pri")
			if !w.isPrisoner(int(c.who)) && w.charMoving(int(c.who)) == 0 && c.secondWait == 0 {
				return pri
			}
		}
	}
	return 99
}

func (w *world) initWaitList() {
	for i := range w.loopKind(tChar) {
		c := w.rpCommand(i)
		if c != nil && c.state == stateRun && c.cmd == w.cmdWait {
			w.waitList.Append(i)
		}
	}
}

func (w *world) checkAllWaits() {
	// Index the live list: d_wait may remove the current entry, shifting the
	// next one into its slot. C advances past that shifted entry.
	for i := 0; i < w.waitList.Len(); i++ {
		c := w.rpCommand(w.waitList[i])
		if c != nil && c.state == stateRun && c.cmd == w.cmdWait {
			assert(c.wait == -1, "olympia/input.c", 821, "c->wait == -1")
			w.finishCommand(c)
		}
	}
}

func (w *world) startPhase() {
	for {
		pri := w.minPriReady()
		w.curPri = int32(pri)
		if w.autoAttackFlag && pri >= 3 {
			w.checkAllAutoAttacks()
			w.autoAttackFlag = false
		}
		if pri == 99 {
			return
		}
		w.sortLoadQueue(w.loadQ[pri])
		for _, i := range w.loadQ[pri].Copy() {
			c := w.rpCommand(i)
			assert(c != nil, "olympia/input.c", 868, "c != NULL")
			if c.state == stateLoad && int(c.pri) == pri && !w.isPrisoner(i) && w.charMoving(i) == 0 && c.secondWait == 0 {
				w.doCommand(c)
				w.checkAllWaits()
				if int32(pri) != w.curPri {
					break
				}
			}
		}
	}
}

func (w *world) eveningPhase() {
	w.evening = 1
	w.sortRunQueue(w.runQ)
	for _, i := range w.runQ.Copy() {
		c := w.rpCommand(i)
		assert(c != nil, "olympia/input.c", 905, "c")
		if c.state != stateRun || c.secondWait != 0 {
			continue
		}
		assert(c.state == stateRun, "olympia/input.c", 913, "c->state == STATE_RUN")
		c.daysExecuting++
		w.finishCommand(c)
		if w.tradesToCheck.Len() != 0 {
			panic(notImplemented("check_validated_trades"))
		}
	}
	w.checkAllWaits()
	w.evening = 0
}

func (w *world) dailyCommandLoop() {
	w.autoAttackFlag = true
	w.startPhase()
	w.eveningPhase()
	w.clearSecondWaits()
}

func (w *world) processPlayerOrders() {
	for pl := range w.loopKind(tPlayer) {
		c := w.rpCommand(pl)
		if c == nil {
			continue
		}
		for w.kind(pl) == tPlayer && c.state == stateLoad {
			w.doCommand(c)
		}
	}
}

func (w *world) interruptOrder(who int) {
	if w.stackLeader(who) == who {
		w.restoreStackActions(who)
	}
	c := w.rpCommand(who)
	if c == nil {
		return
	}
	if c.state == stateRun {
		if fn := cmdTbl[c.cmd].interrupt; fn != nil {
			c.status = int8(fn(w, c))
		}
		w.commandDone(c)
		assert(c.state != stateRun, "olympia/input.c", 1016, "c->state != STATE_RUN")
	}
}

func (w *world) processInterruptedUnits() {
	for who := range w.loopKind(tChar) {
		if !w.stopOrder(w.player(who), who) {
			continue
		}
		w.popOrder(w.player(who), who)
		w.commandOut(who, "> stop")
		c := w.rpCommand(who)
		if c != nil && c.cmd != 0 && c.wait != 0 {
			w.commandOut(who, "Interrupt current order.")
			w.interruptOrder(who)
		} else {
			w.commandOut(who, "No order is currently executing.")
		}
	}
}
