// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/input.c: the command line parser that loading a
// command (scan_command() in io.c) uses.
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
