// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"syscall"

	"github.com/mdhender/g3g0/internal/clib"
	"github.com/mdhender/g3g0/internal/lists"
)

// Ported from olympia/io.c: the line helpers and the field scan and print
// helpers that the entity sections use.
//
// The data files hold one field per line, "xy data", where xy is a
// two-letter tag (see linehash). A list too long for one line continues on
// the next: the line ends in a backslash, and the scan helper reads the next
// line with readlin_ew(). The print helpers write to fp exactly what C's
// fprintf() and fputs() write; C ignores write errors there, and so do
// they.
//
// Where a list's last line ends in a backslash at the end of the file,
// readlin_ew() returns NULL and the scan helper stops there, keeping the
// entries it has read, as if that line had no backslash. (Before upstream
// fixed https://github.com/playbymail/olympia-g3/issues/313, C dereferenced
// the NULL and crashed.) A blank continuation line ends a list too (see
// lastIsBackslash).

// at is s[i] as C sees it: 0 (the terminating NUL) past the end of s.
func at(s string, i int) byte {
	if i < len(s) {
		return s[i]
	}
	return 0
}

// linehash is the linehash() macro: a line's two-letter tag crunched into
// an int, t[0] << 8 | t[1], or 0 if the line is shorter than two bytes. The
// bytes are plain chars, signed on the reference builds, so a byte from
// 0x80 up is negative, and t[1] sign-extends over t[0]. No tag that a case
// label names has such a byte, so those lines fall to the default case.
// Compare the result with 'x'<<8 | 'y'.
func linehash(t string) int {
	if len(t) < 2 {
		return 0
	}
	return int(int8(t[0]))<<8 | int(int8(t[1]))
}

// tString is the t_string() macro: the data after a line's tag and the
// one byte that follows it, or "" if the line is shorter than four bytes.
func tString(t string) string {
	if len(t) >= 4 {
		return t[3:]
	}
	return ""
}

// advance is the advance macro: it reads the next line into line.
func (w *world) advance() { w.line = w.readlin() }

// lastIsBackslash is the test *s && s[strlen(s)-1] == '\\' that the list
// scans use to see whether another entry follows. An empty s, t_string()'s
// "" or a blank continuation line, ends the list. (Before upstream fixed
// https://github.com/playbymail/olympia-g3/issues/318, which also fixed
// https://github.com/playbymail/olympia-g3/issues/312, C read the byte
// before an empty s, which was undefined.)
func lastIsBackslash(s string) bool {
	return len(s) > 0 && s[len(s)-1] == '\\'
}

// scanInts is io.c's scan_ints, with one conversion for each of ptrs,
// each an *int32, *int16 (for olytime) or *int. It stops at the
// first conversion that fails and leaves that target and the rest as they
// were. It returns the number of conversions made.
//
// Numbers have no token-length limit: strtoll saturates to signed 64 bits,
// then the int keeps its low 32 bits (olytime shorts keep the low 16).
// clib.Atoi has that same decimal saturation and narrowing contract.
func scanInts(s string, ptrs ...any) int {
	for i, p := range ptrs {
		start := 0
		for start < len(s) && clib.IsSpace(s[start]) {
			start++
		}
		end := start
		if end < len(s) && (s[end] == '+' || s[end] == '-') {
			end++
		}
		digits := end
		for end < len(s) && isdigit(s[end]) {
			end++
		}
		if end == digits {
			return i
		}
		v := clib.Atoi(s[start:end])
		s = s[end:]
		switch p := p.(type) {
		case *int32:
			*p = int32(v)
		case *int16:
			*p = int16(v)
		case *int:
			*p = v
		default:
			panic(fmt.Sprintf("olympia: scanInts: unsupported target %T", p))
		}
	}
	return len(ptrs)
}

// zeroCheck is zero_check(p, sizeof(*p)): whether the structure p points to
// is all zero bytes. Under the rules of entity.go that is
// reflect.Value.IsZero, padding aside: a non-NULL pointer, string or list,
// even an empty one, is not zero.
func zeroCheck[T any](p *T) bool {
	return reflect.ValueOf(p).Elem().IsZero()
}

// ifName is if_name(): " (name)" for box n, to make the data files easier
// to read, if pretty_data_files is set and the box exists and has a
// non-empty name; otherwise "". C formats it with sout(), which cuts it to
// BUF_LEN (LEN) bytes.
func (w *world) ifName(num int) string {
	if !w.prettyDataFiles {
		return ""
	}

	if !w.validBox(num) {
		return ""
	}

	s := w.name(num)

	if s != "" {
		s = fmt.Sprintf(" (%s)", s)
		return s[:min(len(s), lenMax)]
	}

	return ""
}

// boxScan is box_scan(): the box number at the start of t, or 0, with a
// diagnostic that shows the whole line, if it is not a valid box.
func (w *world) boxScan(t string) int {
	n := clib.Atoi(t)

	if w.validBox(n) {
		return n
	}

	fmt.Fprintf(w.stderr, "box_scan(%d): bad reference: %s\n", w.extBoxnum, clib.Str(w.line))

	return 0
}

// boxPrint is box_print(): header, n and its pretty name, if n is a valid
// box.
func (w *world) boxPrint(fp io.Writer, header string, n int) {
	if w.validBox(n) {
		fmt.Fprintf(fp, "%s%d%s\n", header, n, w.ifName(n))
	}
}

// scanNumbers is the loop that boxlist0_scan, boxlist_scan, admit_scan,
// ilist_scan and known_scan share: it calls f with each run of digits in s
// (atoi() of it), skipping white space, and reads a continuation line with
// readlin_ew() wherever s has a backslash. Anything else ends the list, as
// does the end of the file where a continuation line should be. A sign is
// not a digit, so "-5" ends it too.
func (w *world) scanNumbers(s string, f func(n int)) {
	for {
		i := 0
		for i < len(s) && iswhite(s[i]) {
			i++
		}
		s = s[i:]

		switch {
		case isdigit(at(s, 0)):
			f(clib.Atoi(s))

			i = 0
			for i < len(s) && isdigit(s[i]) {
				i++
			}
			s = s[i:]
		case at(s, 0) == '\\': // continuation line follows
			next := w.readlinEw()
			if next == nil { // file ends mid-list
				return
			}
			s = *next
		default:
			return
		}
	}
}

// boxlist0Scan is boxlist0_scan(): boxlistScan, but it keeps 0 as well as
// valid boxes.
func (w *world) boxlist0Scan(s string, boxNum int, l *lists.IList) {
	w.scanNumbers(s, func(n int) {
		if n == 0 || w.validBox(n) {
			l.Append(n)
		} else {
			fmt.Fprintf(w.stderr, "boxlist_scan(%d): bad box reference: %d\n", boxNum, n)
		}
	})
}

// boxlist0Print is boxlist0_print(): boxlistPrint, but it prints 0 as well
// as valid boxes.
func (w *world) boxlist0Print(fp io.Writer, header string, l lists.IList) {
	count := 0

	for i := 0; i < l.Len(); i++ {
		if l[i] == 0 || w.validBox(l[i]) {
			count++
			if count == 1 {
				io.WriteString(fp, header)
			}

			if count%11 == 10 { // continuation line
				io.WriteString(fp, "\\\n\t")
			}

			fmt.Fprintf(fp, "%d ", l[i])
		}
	}

	if count != 0 {
		fmt.Fprintf(fp, "\n")
	}
}

// boxlistScan is boxlist_scan(): it appends each valid box in s, and its
// continuation lines, to l, with a diagnostic for each number that is not
// a valid box.
func (w *world) boxlistScan(s string, boxNum int, l *lists.IList) {
	w.scanNumbers(s, func(n int) {
		if w.validBox(n) {
			l.Append(n)
		} else {
			fmt.Fprintf(w.stderr, "boxlist_scan(%d): bad box reference: %d\n", boxNum, n)
		}
	})
}

// boxlistPrint is boxlist_print(): header and the valid boxes of l, each
// followed by a space, or nothing if l has none. A continuation line starts
// before the 10th box, and every 11th after it.
func (w *world) boxlistPrint(fp io.Writer, header string, l lists.IList) {
	count := 0

	for i := 0; i < l.Len(); i++ {
		if w.validBox(l[i]) {
			count++
			if count == 1 {
				io.WriteString(fp, header)
			}

			if count%11 == 10 { // continuation line
				io.WriteString(fp, "\\\n\t")
			}

			fmt.Fprintf(fp, "%d ", l[i])
		}
	}

	if count != 0 {
		fmt.Fprintf(fp, "\n")
	}
}

// admitPrintSup is admit_print_sup(): one " am" line, if p's target is a
// valid box and p has a sense or a list. The count starts at 2, for the
// target and sense, so a continuation line starts before the 8th box of
// the list, and "\n" always follows, even if no box of the list is valid.
func (w *world) admitPrintSup(fp io.Writer, p *admit) {
	count := 2

	if !w.validBox(int(p.targ)) {
		return
	}

	if p.sense == 0 && p.l.Len() == 0 {
		return
	}

	fmt.Fprintf(fp, " am %d %d ", p.targ, p.sense)

	for i := 0; i < p.l.Len(); i++ {
		if w.validBox(p.l[i]) {
			count++

			if count%11 == 10 { // continuation line
				io.WriteString(fp, "\\\n\t")
			}

			fmt.Fprintf(fp, "%d ", p.l[i])
		}
	}

	if count != 0 {
		fmt.Fprintf(fp, "\n")
	}
}

// admitPrint is admit_print(): an " am" line for each admit of p.
func (w *world) admitPrint(fp io.Writer, p *entityPlayer) {
	for i := 0; i < p.admits.Len(); i++ {
		w.admitPrintSup(fp, p.admits[i])
	}
}

// admitScan is admit_scan(): an admit from s and its continuation lines,
// the target, the sense and then the boxes, appended to pp's admits if the
// target is a valid box. Each later number that is not a valid box, and an
// invalid target, gets a diagnostic.
func (w *world) admitScan(s string, boxNum int, pp *entityPlayer) {
	count := 0
	p := &admit{}

	w.scanNumbers(s, func(n int) {
		switch count {
		case 0:
			p.targ = int32(n)
		case 1:
			p.sense = int32(n)
		default:
			if w.validBox(n) {
				p.l.Append(n)
			} else {
				fmt.Fprintf(w.stderr, "admit_scan(%d): bad box reference: %d\n", boxNum, n)
			}
		}
		count++
	})

	if !w.validBox(int(p.targ)) {
		fmt.Fprintf(w.stderr, "admit_scan(%d): bad targ %d\n", boxNum, p.targ)
		return
	}

	pp.admits.Append(p)
}

// ilistPrint is ilist_print(): header and every element of l, each
// followed by a space, or nothing if l is empty. A continuation line starts
// before the 11th element, and every 11th after it.
func ilistPrint(fp io.Writer, header string, l lists.IList) {
	if l.Len() > 0 {
		io.WriteString(fp, header)

		for i := 0; i < l.Len(); i++ {
			if i%11 == 10 { // continuation line
				fmt.Fprintf(fp, "\\\n\t")
			}

			fmt.Fprintf(fp, "%d ", l[i])
		}

		fmt.Fprintf(fp, "\n")
	}
}

// ilistScan is ilist_scan(): it appends every number in s, and its
// continuation lines, to l.
func (w *world) ilistScan(s string, l *lists.IList) {
	w.scanNumbers(s, func(n int) {
		l.Append(n)
	})
}

// sortKnown is the safe_qsort(kn, ilist_len(kn), sizeof(int), int_comp)
// that loop_known() in olympia/loop.h does first: it sorts kn in place, in
// the caller's list.
func sortKnown(kn sparse) {
	safeQsort(kn, intComp)
}

// knownPrint is known_print(): header and the valid boxes of kn, each
// followed by a space, or nothing if kn has none. It sorts kn first, as
// loop_known() does. A continuation line starts before the 11th box, and
// every 11th after it.
func (w *world) knownPrint(fp io.Writer, header string, kn sparse) {
	count := 0
	first := true

	sortKnown(kn)
	for _, i := range kn {
		if !w.validBox(i) {
			continue
		}

		if first {
			io.WriteString(fp, header)
			first = false
		}

		if count%11 == 10 {
			fmt.Fprintf(fp, "\\\n\t")
		}
		count++

		fmt.Fprintf(fp, "%d ", i)
	}

	if !first {
		fmt.Fprintf(fp, "\n")
	}
}

// knownScan is known_scan(): it sets the bit in kn of each valid box in s,
// and its continuation lines, with a diagnostic for each number that is
// not a valid box.
func (w *world) knownScan(s string, kn *sparse, boxNum int) {
	w.scanNumbers(s, func(n int) {
		if w.validBox(n) {
			setBit(kn, n)
		} else {
			fmt.Fprintf(w.stderr, "known_scan(%d): bad box reference: %d\n", boxNum, n)
		}
	})
}

// skillListPrint is skill_list_print(): header and "skill know days
// experience 0" for each skill of l that is a valid box, one entry to a
// line, joined by " \\\n\t", or nothing if l has none.
func (w *world) skillListPrint(fp io.Writer, header string, l skillEntsList) {
	count := 0

	for i := 0; i < l.Len(); i++ {
		if w.validBox(int(l[i].skill)) {
			count++
			if count == 1 {
				io.WriteString(fp, header)
			}

			if count > 1 {
				io.WriteString(fp, " \\\n\t")
			}

			fmt.Fprintf(fp, "%d %d %d %d 0",
				l[i].skill,
				l[i].know,
				l[i].daysStudied,
				l[i].experience)
		}
	}

	if count != 0 {
		io.WriteString(fp, "\n")
	}
}

// skillListScan is skill_list_scan(): a skill entry from s, and from each
// continuation line while the last one ends in a backslash, appended to l
// if its skill is a valid box. The fifth number is read and ignored.
//
// A field that a short line leaves out is 0 for that entry. (Fixed
// upstream in https://github.com/playbymail/olympia-g3/issues/311; before,
// C read know and experience from uninitialized locals.)
func (w *world) skillListScan(s string, l *skillEntsList, boxNum int) {
	for {
		p := &skillEnt{}
		var dummy, know, experience int // a short line leaves them unset
		scanInts(s, &p.skill, &know, &p.daysStudied, &experience, &dummy)

		p.know = int8(know)
		p.experience = int16(experience)

		if w.validBox(int(p.skill)) {
			l.Append(p)
		} else {
			fmt.Fprintf(w.stderr, "skill_list_scan(%d): bad skill %d\n", boxNum, p.skill)
		}

		if lastIsBackslash(s) { // another entry follows
			next := w.readlinEw()
			if next == nil { // file ends mid-list
				break
			}
			s = *next
		} else {
			break
		}
	}
}

// itemListPrint is item_list_print(): header and "item qty" for each entry
// of l whose item is a valid box and whose qty is positive, joined by
// " \\\n\t", or nothing if l has none.
func (w *world) itemListPrint(fp io.Writer, header string, l itemEntsList) {
	count := 0

	for i := 0; i < l.Len(); i++ {
		if w.validBox(int(l[i].item)) && l[i].qty > 0 {
			count++
			if count == 1 {
				io.WriteString(fp, header)
			}

			if count > 1 {
				io.WriteString(fp, " \\\n\t")
			}

			fmt.Fprintf(fp, "%d %d",
				l[i].item,
				l[i].qty)
		}
	}

	if count != 0 {
		io.WriteString(fp, "\n")
	}
}

// itemListScan is item_list_scan(): an item entry from s, and from each
// continuation line while the last one ends in a backslash, appended to l
// if its item is a valid box.
func (w *world) itemListScan(s string, l *itemEntsList, boxNum int) {
	for {
		p := &itemEnt{}
		scanInts(s, &p.item, &p.qty)

		if w.validBox(int(p.item)) {
			l.Append(p)
		} else {
			fmt.Fprintf(w.stderr, "item_list_scan(%d): bad item %d\n", boxNum, p.item)
		}

		if lastIsBackslash(s) { // another entry follows
			next := w.readlinEw()
			if next == nil { // file ends mid-list
				break
			}
			s = *next
		} else {
			break
		}
	}
}

// tradeListPrint is trade_list_print(): header and the eight saved fields
// of each trade of l whose item is a valid box, joined by " \\\n\t", or
// nothing if l has none. BUY and SELL trades with no qty left are left out.
func (w *world) tradeListPrint(fp io.Writer, header string, l tradesList) {
	count := 0

	for i := 0; i < l.Len(); i++ {
		if w.validBox(int(l[i].item)) {
			// Weed out completed or cleared BUY and SELL trades, but don't
			// touch PRODUCE or CONSUME zero-qty trades.

			if (l[i].kind == tradeBuy || l[i].kind == tradeSell) &&
				l[i].qty <= 0 {
				continue
			}

			count++
			if count == 1 {
				io.WriteString(fp, header)
			}

			if count > 1 {
				io.WriteString(fp, " \\\n\t")
			}

			fmt.Fprintf(fp, "%d %d %d %d %d %d %d %d",
				l[i].kind,
				l[i].item,
				l[i].qty,
				l[i].cost,
				l[i].cloak,
				l[i].haveLeft,
				l[i].monthProd,
				l[i].expire)
		}
	}

	if count != 0 {
		io.WriteString(fp, "\n")
	}
}

// tradeListScan is trade_list_scan(): a trade from s, and from each
// continuation line while the last one ends in a backslash, appended to l
// if its item is valid and a positive quantity has a price. Each trade's
// who is boxNum.
func (w *world) tradeListScan(s string, l *tradesList, boxNum int) {
	for {
		p := &trade{}
		scanInts(s,
			&p.kind,
			&p.item,
			&p.qty,
			&p.cost,
			&p.cloak,
			&p.haveLeft,
			&p.monthProd,
			&p.expire)

		p.who = int32(boxNum)

		if !w.validBox(int(p.item)) {
			fmt.Fprintf(w.stderr, "trade_list_scan(%d): bad item %d\n", boxNum, p.item)
		} else if p.qty > 0 && p.cost < 1 {
			// BUY/SELL refuse these orders, and traders divide by the price.
			fmt.Fprintf(w.stderr, "trade_list_scan(%d): item %d qty %d has no price\n", boxNum, p.item, p.qty)
		} else {
			l.Append(p)
		}

		if lastIsBackslash(s) { // another entry follows
			next := w.readlinEw()
			if next == nil { // file ends mid-list
				break
			}
			s = *next
		} else {
			break
		}
	}
}

// reqListPrint is req_list_print(): header and "item qty consume" for each
// entry of l whose item is a valid box, joined by " \\\n\t", or nothing if
// l has none.
func (w *world) reqListPrint(fp io.Writer, header string, l reqEntsList) {
	count := 0

	for i := 0; i < l.Len(); i++ {
		if w.validBox(int(l[i].item)) {
			count++
			if count == 1 {
				io.WriteString(fp, header)
			}

			if count > 1 {
				io.WriteString(fp, " \\\n\t")
			}

			fmt.Fprintf(fp, "%d %d %d",
				l[i].item,
				l[i].qty,
				l[i].consume)
		}
	}

	if count != 0 {
		io.WriteString(fp, "\n")
	}
}

// reqListScan is req_list_scan(): an entry from s, and from each
// continuation line while the last one ends in a backslash, appended to l
// if its item is a valid box.
//
// A consume that a short line leaves out is 0 for that entry. (Fixed
// upstream in https://github.com/playbymail/olympia-g3/issues/311; before,
// C read it from an uninitialized local.)
func (w *world) reqListScan(s string, l *reqEntsList, boxNum int) {
	for {
		p := &reqEnt{}
		var consume int // a short line leaves it unset
		scanInts(s, &p.item, &p.qty, &consume)
		p.consume = int8(consume)

		if w.validBox(int(p.item)) {
			l.Append(p)
		} else {
			fmt.Fprintf(w.stderr, "req_list_scan(%d): bad item %d\n", boxNum, p.item)
		}

		if lastIsBackslash(s) { // another entry follows
			next := w.readlinEw()
			if next == nil { // file ends mid-list
				break
			}
			s = *next
		} else {
			break
		}
	}
}

// olytimeScan is olytime_scan(): "turn day days" from s into p. Fields
// that s does not give keep their values.
func olytimeScan(s string, p *olytime) {
	scanInts(s, &p.turn, &p.day, &p.daysSinceEpoch)
}

// olytimePrint is olytime_print(): header and "turn day days", unless all
// three are 0.
func olytimePrint(fp io.Writer, header string, p *olytime) {
	if p.turn != 0 || p.day != 0 || p.daysSinceEpoch != 0 {
		fmt.Fprintf(fp, "%s%d %d %d\n", header,
			p.turn, p.day, p.daysSinceEpoch)
	}
}

// hShift is the shift of linehash(): the case labels of the scanners
// write C's 'wh' as 'w'<<hShift | 'h'.
const hShift = 8

// scanSection is the loop that every scan_xxx() section scanner runs: for
// each line after the section's tag line that starts with white space, it
// moves line past that byte and calls field with the line's tag
// (linehash) and data (t_string). field returns false for a tag it does
// not know, and the line gets a diagnostic that names the scanner. The
// first line that does not start with white space, or nil at end of file,
// ends the section and is left in line for load_box().
//
// C tests *line == '#' first, to skip a comment, but *line is white space
// there, so it never is.
func (w *world) scanSection(scanner string, boxNum int, field func(c int, t string) bool) {
	w.advance()
	for w.line != nil && iswhite(at(*w.line, 0)) {
		s := (*w.line)[1:]
		w.line = &s

		c := linehash(s)
		t := tString(s)

		if !field(c, t) {
			fmt.Fprintf(w.stderr, "%s(%d):  bad line: %s\n", scanner, boxNum, s)
		}
		w.advance()
	}
}

// strSave is str_save(): a new copy of s.
func strSave(s string) *string { return &s }

// printLocInfo is print_loc_info(): the "LI" section, unless p is all
// zero.
func (w *world) printLocInfo(fp io.Writer, p *locInfo) {
	if zeroCheck(p) {
		return
	}

	fmt.Fprintf(fp, "LI\n")
	w.boxPrint(fp, " wh ", int(p.where))
	w.boxlistPrint(fp, " hl ", p.hereList)
}

// scanLocInfo is scan_loc_info().
func (w *world) scanLocInfo(p *locInfo, boxNum int) {
	w.scanSection("scan_loc_info", boxNum, func(c int, t string) bool {
		switch c {
		case 'w'<<hShift | 'h':
			p.where = int32(w.boxScan(t))
		case 'h'<<hShift | 'l':
			w.boxlistScan(t, boxNum, &p.hereList)
		default:
			return false
		}
		return true
	})
}

// printMagic is print_magic(): the "CM" section, unless p is all zero.
func (w *world) printMagic(fp io.Writer, p *charMagic) {
	if zeroCheck(p) {
		return
	}

	fmt.Fprintf(fp, "CM\n")

	if p.magician != 0 {
		fmt.Fprintf(fp, " im %d\n", p.magician)
	}

	if p.maxAura != 0 {
		fmt.Fprintf(fp, " ma %d\n", p.maxAura)
	}

	if p.curAura != 0 {
		fmt.Fprintf(fp, " ca %d\n", p.curAura)
	}

	if p.abilityShroud != 0 {
		fmt.Fprintf(fp, " as %d\n", p.abilityShroud)
	}

	if p.hinderMeditation != 0 {
		fmt.Fprintf(fp, " hm %d\n", p.hinderMeditation)
	}

	if p.quickCast != 0 {
		fmt.Fprintf(fp, " qc %d\n", p.quickCast)
	}

	if p.auraReflect != 0 {
		fmt.Fprintf(fp, " rb %d\n", p.auraReflect)
	}

	if p.hideSelf != 0 {
		fmt.Fprintf(fp, " hs %d\n", p.hideSelf)
	}

	if p.hideMage != 0 {
		fmt.Fprintf(fp, " cm %d\n", p.hideMage)
	}

	if p.pray != 0 {
		fmt.Fprintf(fp, " pr %d\n", p.pray)
	}

	if p.knowsWeather != 0 {
		fmt.Fprintf(fp, " kw %d\n", p.knowsWeather)
	}

	if p.defaultGarr != 0 {
		fmt.Fprintf(fp, " dg %d\n", p.defaultGarr)
	}

	if p.swearOnRelease != 0 {
		fmt.Fprintf(fp, " sr %d\n", p.swearOnRelease)
	}

	if p.fee != 0 {
		fmt.Fprintf(fp, " bf %d\n", p.fee)
	}

	if p.visProtect != 0 {
		fmt.Fprintf(fp, " vp %d\n", p.visProtect)
	}

	w.boxPrint(fp, " pl ", int(p.pledge))
	w.boxPrint(fp, " pc ", int(p.projectCast))
	w.boxPrint(fp, " ar ", int(p.auraculum))
	w.boxPrint(fp, " ot ", int(p.token)) // our token artifact
	w.knownPrint(fp, " vi ", p.visions)
}

// scanMagic is scan_magic().
func (w *world) scanMagic(p *charMagic, boxNum int) {
	w.scanSection("scan_magic", boxNum, func(c int, t string) bool {
		switch c {
		case 'i'<<hShift | 'm':
			p.magician = int8(clib.Atoi(t))
		case 'm'<<hShift | 'a':
			p.maxAura = int32(clib.Atoi(t))
		case 'c'<<hShift | 'a':
			p.curAura = int32(clib.Atoi(t))
		case 'a'<<hShift | 's':
			p.abilityShroud = int16(clib.Atoi(t))
		case 'h'<<hShift | 'm':
			p.hinderMeditation = int8(clib.Atoi(t))
		case 'p'<<hShift | 'c':
			p.projectCast = int32(w.boxScan(t))
		case 'q'<<hShift | 'c':
			p.quickCast = int16(clib.Atoi(t))
		case 'o'<<hShift | 't':
			p.token = int32(w.boxScan(t))
		case 'p'<<hShift | 'l':
			p.pledge = int32(w.boxScan(t))
		case 'a'<<hShift | 'r':
			p.auraculum = int32(w.boxScan(t))
		case 'r'<<hShift | 'b':
			p.auraReflect = int8(clib.Atoi(t))
		case 'h'<<hShift | 's':
			p.hideSelf = int8(clib.Atoi(t))
		case 'c'<<hShift | 'm':
			p.hideMage = int8(clib.Atoi(t))
		case 'p'<<hShift | 'r':
			p.pray = int8(clib.Atoi(t))
		case 's'<<hShift | 'r':
			p.swearOnRelease = int8(clib.Atoi(t))
		case 'k'<<hShift | 'w':
			p.knowsWeather = int8(clib.Atoi(t))
		case 'v'<<hShift | 'p':
			p.visProtect = int8(clib.Atoi(t))
		case 'd'<<hShift | 'g':
			p.defaultGarr = int8(clib.Atoi(t))
		case 'b'<<hShift | 'f':
			p.fee = int32(clib.Atoi(t))
		case 'v'<<hShift | 'i':
			w.knownScan(t, &p.visions, boxNum)
		default:
			return false
		}
		return true
	})
}

// printChar is print_char(): the "CH" section, always.
func (w *world) printChar(fp io.Writer, p *entityChar) {
	fmt.Fprintf(fp, "CH\n")

	w.boxPrint(fp, " ni ", int(p.unitItem))
	w.boxPrint(fp, " lo ", int(p.unitLord))
	w.boxPrint(fp, " pl ", int(p.prevLord))

	if p.health != 0 {
		fmt.Fprintf(fp, " he %d\n", p.health)
	}

	if p.sick != 0 {
		fmt.Fprintf(fp, " si %d\n", p.sick)
	}

	if p.loyKind != 0 {
		fmt.Fprintf(fp, " lk %d\n", p.loyKind)
	}
	if p.loyRate != 0 {
		fmt.Fprintf(fp, " lr %d\n", p.loyRate)
	}

	w.skillListPrint(fp, " sl\t", p.skills)

	if p.prisoner != 0 {
		fmt.Fprintf(fp, " pr %d\n", p.prisoner)
	}

	if p.moving != 0 {
		fmt.Fprintf(fp, " mo %d\n", p.moving)
	}

	if p.behind != 0 {
		fmt.Fprintf(fp, " bh %d\n", p.behind)
	}

	if p.guard != 0 {
		fmt.Fprintf(fp, " gu %d\n", p.guard)
	}

	if p.timeFlying != 0 {
		fmt.Fprintf(fp, " tf %d\n", p.timeFlying)
	}

	if p.breakPoint != 0 {
		fmt.Fprintf(fp, " bp %d\n", p.breakPoint)
	}

	if p.rank != 0 {
		fmt.Fprintf(fp, " ra %d\n", p.rank)
	}

	if p.attack != 0 {
		fmt.Fprintf(fp, " at %d\n", p.attack)
	}

	if p.defense != 0 {
		fmt.Fprintf(fp, " df %d\n", p.defense)
	}

	if p.missile != 0 {
		fmt.Fprintf(fp, " mi %d\n", p.missile)
	}

	if p.npcProg != 0 {
		fmt.Fprintf(fp, " po %d\n", p.npcProg)
	}

	w.boxlistPrint(fp, " ct ", p.contact)

	olytimePrint(fp, " dt ", &p.deathTime)
}

// scanChar is scan_char().
func (w *world) scanChar(p *entityChar, boxNum int) {
	w.scanSection("scan_char", boxNum, func(c int, t string) bool {
		switch c {
		case 'n'<<hShift | 'i':
			p.unitItem = int32(w.boxScan(t))
		case 'l'<<hShift | 'o':
			p.unitLord = int32(w.boxScan(t))
		case 'p'<<hShift | 'l':
			p.prevLord = int32(w.boxScan(t))
		case 'h'<<hShift | 'e':
			p.health = int8(clib.Atoi(t))
		case 's'<<hShift | 'i':
			p.sick = int8(clib.Atoi(t))
		case 'p'<<hShift | 'r':
			p.prisoner = int8(clib.Atoi(t))
		case 'm'<<hShift | 'o':
			p.moving = int32(clib.Atoi(t))
		case 'b'<<hShift | 'h':
			p.behind = int8(clib.Atoi(t))
		case 'l'<<hShift | 'k':
			p.loyKind = int8(clib.Atoi(t))
		case 'l'<<hShift | 'r':
			p.loyRate = int32(clib.Atoi(t))
		case 'g'<<hShift | 'u':
			p.guard = int8(clib.Atoi(t))
		case 't'<<hShift | 'f':
			p.timeFlying = int8(clib.Atoi(t))
		case 'b'<<hShift | 'p':
			p.breakPoint = int8(clib.Atoi(t))
		case 'r'<<hShift | 'a':
			p.rank = int8(clib.Atoi(t))
		case 'a'<<hShift | 't':
			p.attack = int16(clib.Atoi(t))
		case 'd'<<hShift | 'f':
			p.defense = int16(clib.Atoi(t))
		case 'm'<<hShift | 'i':
			p.missile = int16(clib.Atoi(t))
		case 'p'<<hShift | 'o':
			p.npcProg = int8(clib.Atoi(t))
		case 'c'<<hShift | 't':
			w.boxlistScan(t, boxNum, &p.contact)
		case 's'<<hShift | 'l':
			w.skillListScan(t, &p.skills, boxNum)
		case 'd'<<hShift | 't':
			olytimeScan(t, &p.deathTime)
		default:
			return false
		}
		return true
	})
}

// printLoc is print_loc(): the "LO" section, unless p is all zero.
func (w *world) printLoc(fp io.Writer, p *entityLoc) {
	if zeroCheck(p) {
		return
	}

	fmt.Fprintf(fp, "LO\n")

	w.boxlist0Print(fp, " pd ", p.provDest)

	if p.hidden != 0 {
		fmt.Fprintf(fp, " hi %d\n", p.hidden)
	}

	if p.shroud != 0 {
		fmt.Fprintf(fp, " sh %d\n", p.shroud)
	}

	if p.barrier != 0 {
		fmt.Fprintf(fp, " ba %d\n", p.barrier)
	}

	if p.distFromGate != 0 {
		fmt.Fprintf(fp, " dg %d\n", p.distFromGate)
	}

	if p.seaLane != 0 {
		fmt.Fprintf(fp, " sl %d\n", p.seaLane)
	}

	if p.civ != 0 {
		fmt.Fprintf(fp, " lc %d\n", p.civ)
	}
}

// scanLoc is scan_loc().
func (w *world) scanLoc(p *entityLoc, boxNum int) {
	w.scanSection("scan_loc", boxNum, func(c int, t string) bool {
		switch c {
		case 'h'<<hShift | 'i':
			p.hidden = int8(clib.Atoi(t))
		case 's'<<hShift | 'h':
			p.shroud = int16(clib.Atoi(t))
		case 'b'<<hShift | 'a':
			p.barrier = int16(clib.Atoi(t))
		case 'd'<<hShift | 'g':
			p.distFromGate = int8(clib.Atoi(t))
		case 'l'<<hShift | 'c':
			p.civ = int8(clib.Atoi(t))
		case 's'<<hShift | 'l':
			p.seaLane = int8(clib.Atoi(t))
		case 'p'<<hShift | 'd':
			w.boxlist0Scan(t, boxNum, &p.provDest)
		default:
			return false
		}
		return true
	})
}

// printSubloc is print_subloc(): the "SL" section, always.
func (w *world) printSubloc(fp io.Writer, p *entitySubloc) {
	fmt.Fprintf(fp, "SL\n")

	w.boxlistPrint(fp, " te ", p.teaches)

	if p.damage != 0 {
		fmt.Fprintf(fp, " da %d\n", p.damage)
	}

	if p.defense != 0 {
		fmt.Fprintf(fp, " de %d\n", p.defense)
	}

	if p.capacity != 0 {
		fmt.Fprintf(fp, " ca %d\n", p.capacity)
	}

	if p.buildMaterials != 0 {
		fmt.Fprintf(fp, " bm %d\n", p.buildMaterials)
	}

	if p.effortRequired != 0 {
		fmt.Fprintf(fp, " er %d\n", p.effortRequired)
	}

	if p.effortGiven != 0 {
		fmt.Fprintf(fp, " eg %d\n", p.effortGiven)
	}

	if p.moving != 0 {
		fmt.Fprintf(fp, " mo %d\n", p.moving)
	}

	if p.galleyRam != 0 {
		fmt.Fprintf(fp, " gr %d\n", p.galleyRam)
	}

	if p.shaftDepth != 0 {
		fmt.Fprintf(fp, " sd %d\n", p.shaftDepth)
	}

	if p.castleLev != 0 {
		fmt.Fprintf(fp, " cl %d\n", p.castleLev)
	}

	if p.safe != 0 {
		fmt.Fprintf(fp, " sh %d\n", p.safe)
	}

	if p.major != 0 {
		fmt.Fprintf(fp, " mc %d\n", p.major)
	}

	if p.opiumEcon != 0 {
		fmt.Fprintf(fp, " op %d\n", p.opiumEcon)
	}

	if p.loot != 0 {
		fmt.Fprintf(fp, " lo %d\n", p.loot)
	}

	if p.prominence != 0 {
		fmt.Fprintf(fp, " cp %d\n", p.prominence)
	}

	if p.uldimFlag != 0 {
		fmt.Fprintf(fp, " uf %d\n", p.uldimFlag)
	}

	if p.summerFlag != 0 {
		fmt.Fprintf(fp, " sf %d\n", p.summerFlag)
	}

	if p.questLate != 0 {
		fmt.Fprintf(fp, " ql %d\n", p.questLate)
	}

	if p.tunnelLevel != 0 {
		fmt.Fprintf(fp, " td %d\n", p.tunnelLevel)
	}

	w.boxlistPrint(fp, " nc ", p.nearCities)

	w.boxlistPrint(fp, " lt ", p.linkTo)
	w.boxlistPrint(fp, " lf ", p.linkFrom)
	w.boxlistPrint(fp, " bs ", p.boundStorms)

	if p.linkWhen != 0 {
		fmt.Fprintf(fp, " lw %d\n", p.linkWhen)
	}

	if p.linkOpen != 0 {
		fmt.Fprintf(fp, " lp %d\n", p.linkOpen)
	}
}

// scanSubloc is scan_subloc().
func (w *world) scanSubloc(p *entitySubloc, boxNum int) {
	w.scanSection("scan_subloc", boxNum, func(c int, t string) bool {
		switch c {
		case 'd'<<hShift | 'a':
			p.damage = uint8(clib.Atoi(t))
		case 'd'<<hShift | 'e':
			p.defense = int32(clib.Atoi(t))
		case 'c'<<hShift | 'a':
			p.capacity = int32(clib.Atoi(t))
		case 'e'<<hShift | 'r':
			p.effortRequired = int32(clib.Atoi(t))
		case 'e'<<hShift | 'g':
			p.effortGiven = int32(clib.Atoi(t))
		case 'b'<<hShift | 'm':
			p.buildMaterials = int8(clib.Atoi(t))
		case 'm'<<hShift | 'o':
			p.moving = int32(clib.Atoi(t))
		case 'g'<<hShift | 'r':
			p.galleyRam = int8(clib.Atoi(t))
		case 's'<<hShift | 'd':
			p.shaftDepth = int16(clib.Atoi(t))
		case 's'<<hShift | 'h':
			p.safe = int8(clib.Atoi(t))
		case 'm'<<hShift | 'c':
			p.major = int8(clib.Atoi(t))
		case 'o'<<hShift | 'p':
			p.opiumEcon = int32(clib.Atoi(t))
		case 'l'<<hShift | 'o':
			p.loot = int8(clib.Atoi(t))
		case 'c'<<hShift | 'p':
			p.prominence = int8(clib.Atoi(t))
		case 'l'<<hShift | 'w':
			p.linkWhen = int8(clib.Atoi(t))
		case 'l'<<hShift | 'p':
			p.linkOpen = int8(clib.Atoi(t))
		case 'u'<<hShift | 'f':
			p.uldimFlag = int8(clib.Atoi(t))
		case 's'<<hShift | 'f':
			p.summerFlag = int8(clib.Atoi(t))
		case 'q'<<hShift | 'l':
			p.questLate = int8(clib.Atoi(t))
		case 't'<<hShift | 'd':
			p.tunnelLevel = int8(clib.Atoi(t))
		case 'c'<<hShift | 'l':
			p.castleLev = int8(clib.Atoi(t))
		case 'l'<<hShift | 't':
			w.boxlistScan(t, boxNum, &p.linkTo)
		case 'l'<<hShift | 'f':
			w.boxlistScan(t, boxNum, &p.linkFrom)
		case 't'<<hShift | 'e':
			w.boxlistScan(t, boxNum, &p.teaches)
		case 'n'<<hShift | 'c':
			w.boxlistScan(t, boxNum, &p.nearCities)
		case 'b'<<hShift | 's':
			w.boxlistScan(t, boxNum, &p.boundStorms)
		default:
			return false
		}
		return true
	})
}

// printItem is print_item(): the "IT" section, unless p is all zero. The
// " fc" (fly_cap) line is printed when fly_cap is set, as G3 does since it
// fixed https://github.com/playbymail/olympia-g3/issues/314 (before, it
// was printed when ride_cap was set).
func (w *world) printItem(fp io.Writer, p *entityItem) {
	if zeroCheck(p) {
		return
	}

	fmt.Fprintf(fp, "IT\n")

	if p.pluralName != nil && *p.pluralName != "" {
		fmt.Fprintf(fp, " pl %s\n", *p.pluralName)
	}

	if p.weight != 0 {
		fmt.Fprintf(fp, " wt %d\n", p.weight)
	}

	if p.landCap != 0 {
		fmt.Fprintf(fp, " lc %d\n", p.landCap)
	}

	if p.rideCap != 0 {
		fmt.Fprintf(fp, " rc %d\n", p.rideCap)
	}

	if p.flyCap != 0 {
		fmt.Fprintf(fp, " fc %d\n", p.flyCap)
	}

	if p.isManItem != 0 {
		fmt.Fprintf(fp, " mu %d\n", p.isManItem)
	}

	if p.prominent != 0 {
		fmt.Fprintf(fp, " pr %d\n", p.prominent)
	}

	if p.animal != 0 {
		fmt.Fprintf(fp, " an %d\n", p.animal)
	}

	if p.attack != 0 {
		fmt.Fprintf(fp, " at %d\n", p.attack)
	}

	if p.defense != 0 {
		fmt.Fprintf(fp, " df %d\n", p.defense)
	}

	if p.missile != 0 {
		fmt.Fprintf(fp, " mi %d\n", p.missile)
	}

	if p.basePrice != 0 {
		fmt.Fprintf(fp, " bp %d\n", p.basePrice)
	}

	if p.capturable != 0 {
		fmt.Fprintf(fp, " ca %d\n", p.capturable)
	}

	w.boxPrint(fp, " un ", int(p.whoHas))
}

// scanItem is scan_item().
func (w *world) scanItem(p *entityItem, boxNum int) {
	w.scanSection("scan_item", boxNum, func(c int, t string) bool {
		switch c {
		case 'p'<<hShift | 'l':
			p.pluralName = strSave(t)
		case 'w'<<hShift | 't':
			p.weight = int16(clib.Atoi(t))
		case 'l'<<hShift | 'c':
			p.landCap = int16(clib.Atoi(t))
		case 'r'<<hShift | 'c':
			p.rideCap = int16(clib.Atoi(t))
		case 'f'<<hShift | 'c':
			p.flyCap = int16(clib.Atoi(t))
		case 'm'<<hShift | 'u':
			p.isManItem = int8(clib.Atoi(t))
		case 'p'<<hShift | 'r':
			p.prominent = int8(clib.Atoi(t))
		case 'a'<<hShift | 'n':
			p.animal = int8(clib.Atoi(t))
		case 'u'<<hShift | 'n':
			p.whoHas = int32(w.boxScan(t))
		case 'a'<<hShift | 't':
			p.attack = int16(clib.Atoi(t))
		case 'd'<<hShift | 'f':
			p.defense = int16(clib.Atoi(t))
		case 'm'<<hShift | 'i':
			p.missile = int16(clib.Atoi(t))
		case 'b'<<hShift | 'p':
			p.basePrice = int32(clib.Atoi(t))
		case 'c'<<hShift | 'a':
			p.capturable = int8(clib.Atoi(t))
		default:
			return false
		}
		return true
	})
}

// printItemMagic is print_item_magic(): the "IM" section, unless p is all
// zero.
func (w *world) printItemMagic(fp io.Writer, p *itemMagic) {
	if zeroCheck(p) {
		return
	}

	fmt.Fprintf(fp, "IM\n")

	if p.aura != 0 {
		fmt.Fprintf(fp, " au %d\n", p.aura)
	}

	if p.curseLoyalty != 0 {
		fmt.Fprintf(fp, " cl %d\n", p.curseLoyalty)
	}

	if p.cloakRegion != 0 {
		fmt.Fprintf(fp, " cr %d\n", p.cloakRegion)
	}

	if p.cloakCreator != 0 {
		fmt.Fprintf(fp, " cc %d\n", p.cloakCreator)
	}

	if p.useKey != 0 {
		fmt.Fprintf(fp, " uk %d\n", p.useKey)
	}

	if p.quickCast != 0 {
		fmt.Fprintf(fp, " qc %d\n", p.quickCast)
	}

	if p.attackBonus != 0 {
		fmt.Fprintf(fp, " ab %d\n", p.attackBonus)
	}

	if p.defenseBonus != 0 {
		fmt.Fprintf(fp, " db %d\n", p.defenseBonus)
	}

	if p.missileBonus != 0 {
		fmt.Fprintf(fp, " mb %d\n", p.missileBonus)
	}

	if p.auraBonus != 0 {
		fmt.Fprintf(fp, " ba %d\n", p.auraBonus)
	}

	if p.relicDecay != 0 {
		fmt.Fprintf(fp, " rd %d\n", p.relicDecay)
	}

	if p.tokenNum != 0 {
		fmt.Fprintf(fp, " tn %d\n", p.tokenNum)
	}

	if p.orbUseCount != 0 {
		fmt.Fprintf(fp, " oc %d\n", p.orbUseCount)
	}

	w.boxPrint(fp, " ti ", int(p.tokenNi))

	w.boxPrint(fp, " rc ", int(p.regionCreated))
	w.boxPrint(fp, " pc ", int(p.projectCast))
	w.boxPrint(fp, " ct ", int(p.creator))
	w.boxPrint(fp, " lo ", int(p.lore))

	w.boxlistPrint(fp, " mu ", p.mayUse)
	w.boxlistPrint(fp, " ms ", p.mayStudy)
}

// scanItemMagic is scan_item_magic(). Unlike the other box fields, "ti"
// (token_ni) is read with atoi(), not box_scan(), so it is not checked.
func (w *world) scanItemMagic(p *itemMagic, boxNum int) {
	w.scanSection("scan_item_magic", boxNum, func(c int, t string) bool {
		switch c {
		case 'a'<<hShift | 'u':
			p.aura = int16(clib.Atoi(t))
		case 'c'<<hShift | 'l':
			p.curseLoyalty = int8(clib.Atoi(t))
		case 'c'<<hShift | 'r':
			p.cloakRegion = int8(clib.Atoi(t))
		case 'c'<<hShift | 'c':
			p.cloakCreator = int8(clib.Atoi(t))
		case 'u'<<hShift | 'k':
			p.useKey = int8(clib.Atoi(t))
		case 'r'<<hShift | 'c':
			p.regionCreated = int32(w.boxScan(t))
		case 'p'<<hShift | 'c':
			p.projectCast = int32(w.boxScan(t))
		case 'c'<<hShift | 't':
			p.creator = int32(w.boxScan(t))
		case 'l'<<hShift | 'o':
			p.lore = int32(w.boxScan(t))
		case 'q'<<hShift | 'c':
			p.quickCast = int16(clib.Atoi(t))
		case 'a'<<hShift | 'b':
			p.attackBonus = int8(clib.Atoi(t))
		case 'd'<<hShift | 'b':
			p.defenseBonus = int8(clib.Atoi(t))
		case 'm'<<hShift | 'b':
			p.missileBonus = int8(clib.Atoi(t))
		case 'b'<<hShift | 'a':
			p.auraBonus = int16(clib.Atoi(t))
		case 'r'<<hShift | 'd':
			p.relicDecay = int16(clib.Atoi(t))
		case 't'<<hShift | 'n':
			p.tokenNum = int8(clib.Atoi(t))
		case 't'<<hShift | 'i':
			p.tokenNi = int32(clib.Atoi(t))
		case 'o'<<hShift | 'c':
			p.orbUseCount = int8(clib.Atoi(t))
		case 'm'<<hShift | 'u':
			w.boxlistScan(t, boxNum, &p.mayUse)
		case 'm'<<hShift | 's':
			w.boxlistScan(t, boxNum, &p.mayStudy)
		default:
			return false
		}
		return true
	})
}

// printPlayer is print_player(): the "PL" section, always.
func (w *world) printPlayer(fp io.Writer, p *entityPlayer) {
	fmt.Fprintf(fp, "PL\n")

	if p.fullName != nil && *p.fullName != "" {
		fmt.Fprintf(fp, " fn %s\n", *p.fullName)
	}

	if p.email != nil && *p.email != "" {
		fmt.Fprintf(fp, " em %s\n", *p.email)
	}

	if p.visEmail != nil && *p.visEmail != "" {
		fmt.Fprintf(fp, " ve %s\n", *p.visEmail)
	}

	if p.lastEmail != nil && *p.lastEmail != "" {
		fmt.Fprintf(fp, " le %s\n", *p.lastEmail)
	}

	if p.password != nil && *p.password != "" {
		fmt.Fprintf(fp, " pw %s\n", *p.password)
	}

	if p.noblePoints != 0 {
		fmt.Fprintf(fp, " np %d\n", p.noblePoints)
	}

	if p.fastStudy != 0 {
		fmt.Fprintf(fp, " fs %d\n", p.fastStudy)
	}

	if p.firstTurn != 0 {
		fmt.Fprintf(fp, " ft %d\n", p.firstTurn)
	}

	if p.format != 0 {
		fmt.Fprintf(fp, " fo %d\n", p.format)
	}

	if p.notab != 0 {
		fmt.Fprintf(fp, " nt %d\n", p.notab)
	}

	if p.firstTower != 0 {
		fmt.Fprintf(fp, " tf %d\n", p.firstTower)
	}

	if p.splitLines != 0 {
		fmt.Fprintf(fp, " sl %d\n", p.splitLines)
	}

	if p.splitBytes != 0 {
		fmt.Fprintf(fp, " sb %d\n", p.splitBytes)
	}

	if p.sentOrders != 0 {
		fmt.Fprintf(fp, " so %d\n", p.sentOrders)
	}

	if p.dontRemind != 0 {
		fmt.Fprintf(fp, " dr %d\n", p.dontRemind)
	}

	if p.compuserve != 0 {
		fmt.Fprintf(fp, " ci %d\n", p.compuserve)
	}

	if p.brokenMailer != 0 {
		fmt.Fprintf(fp, " bm %d\n", p.brokenMailer)
	}

	if p.lastOrderTurn != 0 {
		fmt.Fprintf(fp, " lt %d\n", p.lastOrderTurn)
	}

	w.knownPrint(fp, " kn ", p.known)
	w.boxlistPrint(fp, " un ", p.units)
	w.boxlistPrint(fp, " uf ", p.unformed)
	w.admitPrint(fp, p)
}

// scanPlayer is scan_player().
func (w *world) scanPlayer(p *entityPlayer, boxNum int) {
	w.scanSection("scan_player", boxNum, func(c int, t string) bool {
		switch c {
		case 'f'<<hShift | 'n':
			p.fullName = strSave(t)
		case 'e'<<hShift | 'm':
			p.email = strSave(t)
		case 'v'<<hShift | 'e':
			p.visEmail = strSave(t)
		case 'l'<<hShift | 'e':
			p.lastEmail = strSave(t)
		case 'p'<<hShift | 'w':
			p.password = strSave(t)
		case 'n'<<hShift | 'p':
			p.noblePoints = int16(clib.Atoi(t))
		case 'f'<<hShift | 's':
			p.fastStudy = int16(clib.Atoi(t))
		case 'f'<<hShift | 't':
			p.firstTurn = int32(clib.Atoi(t))
		case 'f'<<hShift | 'o':
			p.format = int8(clib.Atoi(t))
		case 'n'<<hShift | 't':
			p.notab = int8(clib.Atoi(t))
		case 't'<<hShift | 'f':
			p.firstTower = int8(clib.Atoi(t))
		case 's'<<hShift | 'o':
			p.sentOrders = int8(clib.Atoi(t))
		case 'l'<<hShift | 't':
			p.lastOrderTurn = int32(clib.Atoi(t))
		case 's'<<hShift | 'l':
			p.splitLines = int32(clib.Atoi(t))
		case 's'<<hShift | 'b':
			p.splitBytes = int32(clib.Atoi(t))
		case 'c'<<hShift | 'i':
			p.compuserve = int8(clib.Atoi(t))
		case 'b'<<hShift | 'm':
			p.brokenMailer = int8(clib.Atoi(t))
		case 'd'<<hShift | 'r':
			p.dontRemind = int8(clib.Atoi(t))
		case 'k'<<hShift | 'n':
			w.knownScan(t, &p.known, boxNum)
		case 'u'<<hShift | 'n':
			w.boxlistScan(t, boxNum, &p.units)
		case 'u'<<hShift | 'f':
			w.boxlistScan(t, boxNum, &p.unformed)
		case 'a'<<hShift | 'm':
			w.admitScan(t, boxNum, p)
		default:
			return false
		}
		return true
	})
}

// printSkill is print_skill(): the "SK" section, always.
func (w *world) printSkill(fp io.Writer, p *entitySkill) {
	fmt.Fprintf(fp, "SK\n")
	fmt.Fprintf(fp, " tl %d\n", p.timeToLearn)
	w.boxPrint(fp, " rs ", int(p.requiredSkill))

	w.boxlistPrint(fp, " of ", p.offered)
	w.boxlistPrint(fp, " re ", p.research)
	w.reqListPrint(fp, " rq\t", p.req)
	w.boxPrint(fp, " pr ", int(p.produced))

	if p.npReq != 0 {
		fmt.Fprintf(fp, " np %d\n", p.npReq)
	}
	if p.noExp != 0 {
		fmt.Fprintf(fp, " ne %d\n", p.noExp)
	}
}

// scanSkill is scan_skill().
func (w *world) scanSkill(p *entitySkill, boxNum int) {
	w.scanSection("scan_skill", boxNum, func(c int, t string) bool {
		switch c {
		case 't'<<hShift | 'l':
			p.timeToLearn = int32(clib.Atoi(t))
		case 'n'<<hShift | 'e':
			p.noExp = int32(clib.Atoi(t))
		case 'n'<<hShift | 'p':
			p.npReq = int32(clib.Atoi(t))
		case 'r'<<hShift | 's':
			p.requiredSkill = int32(w.boxScan(t))
		case 'p'<<hShift | 'r':
			p.produced = int32(w.boxScan(t))
		case 'o'<<hShift | 'f':
			w.boxlistScan(t, boxNum, &p.offered)
		case 'r'<<hShift | 'e':
			w.boxlistScan(t, boxNum, &p.research)
		case 'r'<<hShift | 'q':
			w.reqListScan(t, &p.req, boxNum)
		default:
			return false
		}
		return true
	})
}

// printCommand is print_command(): the "CO" section, unless p has no
// command.
func (w *world) printCommand(fp io.Writer, p *command) {
	if p.cmd == 0 {
		return
	}

	fmt.Fprintf(fp, "CO\n")
	fmt.Fprintf(fp, " li %s\n", clib.Str(p.line))

	fmt.Fprintf(fp, " ar %d %d %d %d %d %d %d %d\n",
		p.a, p.b, p.c, p.d, p.e, p.f, p.g, p.h)

	if p.state != 0 {
		fmt.Fprintf(fp, " cs %d\n", p.state)
	}

	if p.wait != 0 {
		fmt.Fprintf(fp, " wa %d\n", p.wait)
	}

	if p.status != 0 {
		fmt.Fprintf(fp, " st %d\n", p.status)
	}

	if p.useSkill != 0 {
		w.boxPrint(fp, " us ", int(p.useSkill))
	}

	if p.useExp != 0 {
		fmt.Fprintf(fp, " ue %d\n", p.useExp)
	}

	if p.daysExecuting != 0 {
		fmt.Fprintf(fp, " de %d\n", p.daysExecuting)
	}

	if p.poll != 0 {
		fmt.Fprintf(fp, " po %d\n", p.poll)
	}

	if p.pri != 0 {
		fmt.Fprintf(fp, " pr %d\n", p.pri)
	}

	if p.inhibitFinish != 0 {
		fmt.Fprintf(fp, " if %d\n", p.inhibitFinish)
	}
}

// scanCommand is scan_command(). It sets p's who to boxNum first. The "li"
// line goes through olyParseCmd, which clears use_skill, so a "us" line
// before it is lost; its diagnostic shows the line as olyParseCmd left it,
// without its comment.
func (w *world) scanCommand(p *command, boxNum int) {
	p.who = int32(boxNum)

	w.scanSection("scan_command", boxNum, func(c int, t string) bool {
		switch c {
		case 'l'<<hShift | 'i':
			if !w.olyParseCmd(p, &t) {
				fmt.Fprintf(w.stderr, "scan_command(%d): "+
					"bad cmd %s\n", boxNum, t)
			}
		case 'c'<<hShift | 's':
			p.state = int8(clib.Atoi(t))
		case 'w'<<hShift | 'a':
			p.wait = int32(clib.Atoi(t))
		case 's'<<hShift | 't':
			p.status = int8(clib.Atoi(t))
		case 'd'<<hShift | 'e':
			p.daysExecuting = int32(clib.Atoi(t))
		case 'p'<<hShift | 'o':
			p.poll = int8(clib.Atoi(t))
		case 'p'<<hShift | 'r':
			p.pri = int8(clib.Atoi(t))
		case 'i'<<hShift | 'f':
			p.inhibitFinish = int8(clib.Atoi(t))
		case 'u'<<hShift | 's':
			p.useSkill = int32(w.boxScan(t))
		case 'u'<<hShift | 'e':
			p.useExp = int32(clib.Atoi(t))
		case 'a'<<hShift | 'r':
			scanInts(t,
				&p.a, &p.b, &p.c, &p.d,
				&p.e, &p.f, &p.g, &p.h)
		default:
			return false
		}
		return true
	})
}

// printGate is print_gate(): the "GA" section, always.
func (w *world) printGate(fp io.Writer, p *entityGate) {
	fmt.Fprintf(fp, "GA\n")

	w.boxPrint(fp, " tl ", int(p.toLoc))

	if p.notifyJumps != 0 {
		w.boxPrint(fp, " nj ", int(p.notifyJumps))
	}

	if p.notifyUnseal != 0 {
		w.boxPrint(fp, " nu ", int(p.notifyUnseal))
	}

	if p.sealKey != 0 {
		fmt.Fprintf(fp, " sk %d\n", p.sealKey)
	}

	if p.roadHidden != 0 {
		fmt.Fprintf(fp, " rh %d\n", p.roadHidden)
	}
}

// scanGate is scan_gate().
func (w *world) scanGate(p *entityGate, boxNum int) {
	w.scanSection("scan_gate", boxNum, func(c int, t string) bool {
		switch c {
		case 't'<<hShift | 'l':
			p.toLoc = int32(w.boxScan(t))
		case 'n'<<hShift | 'j':
			p.notifyJumps = int32(w.boxScan(t))
		case 'n'<<hShift | 'u':
			p.notifyUnseal = int32(w.boxScan(t))
		case 's'<<hShift | 'k':
			p.sealKey = int16(clib.Atoi(t))
		case 'r'<<hShift | 'h':
			p.roadHidden = int8(clib.Atoi(t))
		default:
			return false
		}
		return true
	})
}

// printMisc is print_misc(): the "MI" section, unless p is all zero.
// cmd_allow prints with %c, as the one byte it holds.
func (w *world) printMisc(fp io.Writer, p *entityMisc) {
	if zeroCheck(p) {
		return
	}

	fmt.Fprintf(fp, "MI\n")

	w.boxPrint(fp, " sb ", int(p.summonedBy))

	if p.npcDir != 0 {
		fmt.Fprintf(fp, " di %d\n", p.npcDir)
	}

	if p.npcCreated != 0 {
		fmt.Fprintf(fp, " mc %d\n", p.npcCreated)
	}

	if p.mineDelay != 0 {
		fmt.Fprintf(fp, " md %d\n", p.mineDelay)
	}

	if p.stormStr != 0 {
		fmt.Fprintf(fp, " ss %d\n", p.stormStr)
	}

	if p.cmdAllow != 0 {
		fmt.Fprintf(fp, " ca %s\n", []byte{byte(p.cmdAllow)})
	}

	w.boxPrint(fp, " gc ", int(p.garrCastle))
	w.boxPrint(fp, " mh ", int(p.npcHome))
	w.boxPrint(fp, " co ", int(p.npcCookie))
	w.boxPrint(fp, " ov ", int(p.onlyVulnerable))
	w.boxPrint(fp, " ol ", int(p.oldLord))
	w.boxPrint(fp, " bs ", int(p.bindStorm))

	if p.saveName != nil && *p.saveName != "" {
		fmt.Fprintf(fp, " sn %s\n", *p.saveName)
	}

	if p.display != nil && *p.display != "" {
		fmt.Fprintf(fp, " ds %s\n", *p.display)
	}

	w.knownPrint(fp, " nm ", p.npcMemory)
}

// scanMisc is scan_misc(). "ca" (cmd_allow) is the first byte of the
// line's data, or 0 if it has none.
func (w *world) scanMisc(p *entityMisc, boxNum int) {
	w.scanSection("scan_misc", boxNum, func(c int, t string) bool {
		switch c {
		case 'd'<<hShift | 'i':
			p.npcDir = int8(clib.Atoi(t))
		case 'm'<<hShift | 'c':
			p.npcCreated = int32(clib.Atoi(t))
		case 'm'<<hShift | 'd':
			p.mineDelay = int8(clib.Atoi(t))
		case 's'<<hShift | 's':
			p.stormStr = int16(clib.Atoi(t))
		case 'm'<<hShift | 'h':
			p.npcHome = int32(w.boxScan(t))
		case 'g'<<hShift | 'c':
			p.garrCastle = int32(w.boxScan(t))
		case 's'<<hShift | 'b':
			p.summonedBy = int32(w.boxScan(t))
		case 'c'<<hShift | 'o':
			p.npcCookie = int32(w.boxScan(t))
		case 'o'<<hShift | 'v':
			p.onlyVulnerable = int32(w.boxScan(t))
		case 'b'<<hShift | 's':
			p.bindStorm = int32(w.boxScan(t))
		case 'o'<<hShift | 'l':
			p.oldLord = int32(w.boxScan(t))
		case 's'<<hShift | 'n':
			p.saveName = strSave(t)
		case 'd'<<hShift | 's':
			p.display = strSave(t)
		case 'c'<<hShift | 'a':
			p.cmdAllow = int8(at(t, 0))
		case 'n'<<hShift | 'm':
			w.knownScan(t, &p.npcMemory, boxNum)
		default:
			return false
		}
		return true
	})
}

// loadBox is load_box(): it reads box n's fields and sections from line
// on, for read_boxes(), which has read the box's "n kind subkind" line
// into line and left it there. Comment lines ('#') are skipped, and a line
// with a tag load_box() does not know gets a diagnostic, and it and the
// lines after it that start with white space are skipped. A blank line or
// the end of the file ends the box; line moves past the blank line.
//
// Box n must have been allocated by the first pass; if not, C prints a
// diagnostic and fails an assertion.
func (w *world) loadBox(n int) {
	// The fast scan of libdir/master should have allocated all of the
	// boxes. If one was later added manually, it won't have been
	// pre-allocated by the master fast scan. Remove master file and let
	// io.c do the slow scan to allocate all the boxes. Save the database to
	// recreate master.

	if !w.validBox(n) {
		fmt.Fprintf(w.stderr, "Unforseen box %d found in load phase.\n", n)
		fmt.Fprintf(w.stderr, "Remove %s/master and retry.\n", w.libdir)
		assert(false, "olympia/io.c", 2002, "FALSE")
	}

	w.extBoxnum = int32(n)

	w.advance()
	for w.line != nil && *w.line != "" {
		if (*w.line)[0] == '#' {
			w.advance()
			continue
		}

		c := linehash(*w.line)
		t := tString(*w.line)

		switch c {
		case 'n'<<hShift | 'a':
			w.setName(n, t)
			w.advance()

		case 'i'<<hShift | 'l':
			w.itemListScan(t, &w.bx[n].items, n)
			w.advance()

		case 't'<<hShift | 'l':
			w.tradeListScan(t, &w.bx[n].trades, n)
			w.advance()

		case 'a'<<hShift | 'n':
			w.boxlistScan(t, n, &w.pDisp(n).neutral)
			w.advance()

		case 'a'<<hShift | 'd':
			w.boxlistScan(t, n, &w.pDisp(n).defend)
			w.advance()

		case 'a'<<hShift | 'h':
			w.boxlistScan(t, n, &w.pDisp(n).hostile)
			w.advance()

		case 'C'<<hShift | 'H':
			w.scanChar(w.pChar(n), n)
		case 'C'<<hShift | 'M':
			w.scanMagic(w.pMagic(n), n)
		case 'L'<<hShift | 'I':
			w.scanLocInfo(w.pLocInfo(n), n)
		case 'L'<<hShift | 'O':
			w.scanLoc(w.pLoc(n), n)
		case 'S'<<hShift | 'L':
			w.scanSubloc(w.pSubloc(n), n)
		case 'I'<<hShift | 'T':
			w.scanItem(w.pItem(n), n)
		case 'P'<<hShift | 'L':
			w.scanPlayer(w.pPlayer(n), n)
		case 'S'<<hShift | 'K':
			w.scanSkill(w.pSkill(n), n)
		case 'G'<<hShift | 'A':
			w.scanGate(w.pGate(n), n)
		case 'M'<<hShift | 'I':
			w.scanMisc(w.pMisc(n), n)
		case 'I'<<hShift | 'M':
			w.scanItemMagic(w.pItemMagic(n), n)
		case 'C'<<hShift | 'O':
			w.scanCommand(w.pCommand(n), n)

		default:
			fmt.Fprintf(w.stderr, "load_box(%d):  bad line: %s\n", n, *w.line)

			w.advance()
			for w.line != nil && iswhite(at(*w.line, 0)) {
				w.advance()
			}
		}
	}

	if w.line != nil {
		assert(*w.line == "", "olympia/io.c", 2080, "*line == '\\0'")
		w.advance() // advance over blank separating line
	}
}

// saveBox is save_box(): box n to fp, unless it is deleted: its
// "n kind subkind" line, its fields and its sections, in C's order, and a
// blank line. It marks the box saved (temp = 1) for write_leftovers().
func (w *world) saveBox(fp io.Writer, n int) {
	if w.kind(n) == tDeleted {
		return
	}

	assert(w.validBox(n), "olympia/io.c", 2096, "valid_box(n)")

	p := w.bx[n]

	if p.skind != 0 {
		fmt.Fprintf(fp, "%d %s %s\n",
			n, kindS[p.kind], subkindS[p.skind])
	} else {
		fmt.Fprintf(fp, "%d %s 0\n",
			n, kindS[p.kind])
	}

	if p.name != nil && *p.name != "" {
		fmt.Fprintf(fp, "na %s\n", *p.name)
	}

	w.itemListPrint(fp, "il\t", p.items)
	w.tradeListPrint(fp, "tl\t", p.trades)

	if pd := w.rpDisp(n); pd != nil {
		w.boxlistPrint(fp, "an ", pd.neutral)
		w.boxlistPrint(fp, "ad ", pd.defend)
		w.boxlistPrint(fp, "ah ", pd.hostile)
	}

	w.printLocInfo(fp, w.rpLocInfo(n))
	if vp := w.rpChar(n); vp != nil {
		w.printChar(fp, vp)
	}
	if vp := w.rpMagic(n); vp != nil {
		w.printMagic(fp, vp)
	}
	if vp := w.rpLoc(n); vp != nil {
		w.printLoc(fp, vp)
	}
	if vp := w.rpSubloc(n); vp != nil {
		w.printSubloc(fp, vp)
	}
	if vp := w.rpItem(n); vp != nil {
		w.printItem(fp, vp)
	}
	if vp := w.rpItemMagic(n); vp != nil {
		w.printItemMagic(fp, vp)
	}
	if vp := w.rpPlayer(n); vp != nil {
		w.printPlayer(fp, vp)
	}
	if vp := w.rpSkill(n); vp != nil {
		w.printSkill(fp, vp)
	}
	if vp := w.rpGate(n); vp != nil {
		w.printGate(fp, vp)
	}
	if vp := w.rpMisc(n); vp != nil {
		w.printMisc(fp, vp)
	}
	if vp := w.rpCommand(n); vp != nil {
		w.printCommand(fp, vp)
	}

	fmt.Fprintf(fp, "\n")

	w.bx[n].temp = 1 // mark for write_leftovers()
}

// The box files. Pass 1 (fastScan, or scanAllBoxes if libdir/master
// cannot be read) allocates every box; pass 2 (readAllBoxes) loads each
// box with loadBox. writeAllBoxes and writeMaster write them back. The
// files are, under libdir: loc, item, skill, gate, road, ship, unform and
// misc, one per kind (misc holding every box the others do not), fact/N for
// player N and its units, and master, the index pass 1 reads first.
//
// Files are written with fopen(path, "w") (clib.Fopen): created with mode
// 0666 less the umask, or truncated. C ignores write errors, and so do
// these.

// openWriteFp is open_write_fp(): libdir/fnam opened for writing, or nil,
// with a diagnostic, if it cannot be.
func (w *world) openWriteFp(fnam string) *clib.File {
	path := sout("%s/%s", w.libdir, fnam)

	fp, err := clib.Fopen(path)

	if err != nil {
		fmt.Fprintf(w.stderr, "open_write_fp: can't open %s: ", path)
		clib.Perror(w.stderr, "", err)
		return nil
	}

	return fp
}

// writeKind is write_kind(): every box of kind boxKind, in chain order, to
// libdir/fnam.
func (w *world) writeKind(boxKind int, fnam string) {
	fp := w.openWriteFp(fnam)

	if fp != nil {
		for i := range w.loopKind(boxKind) {
			w.saveBox(fp, i)
		}

		fp.Close()
	}
}

// writePlayer is write_player(): player pl and then its units, in the
// order of its units list, to libdir/fact/pl. A unit that is neither a
// character nor deleted fails an assertion, leaving the file with what
// stdio had written of it.
func (w *world) writePlayer(pl int) {
	fp := w.openWriteFp(sout("fact/%d", pl))

	if fp == nil {
		return
	}

	w.saveBox(fp, pl)

	for who := range w.loopUnits(pl) {
		assert(w.kind(who) == tChar || w.kind(who) == tDeleted,
			"olympia/io.c", 2195, "kind(who) == T_char || kind(who) == T_deleted")

		w.saveBox(fp, who)
	}

	fp.Close()
}

// writeChars is write_chars(): writePlayer for every player.
func (w *world) writeChars() {
	for i := range w.loopKind(tPlayer) {
		w.writePlayer(i)
	}
}

// writeLeftovers is write_leftovers(): every box that is not deleted and
// that no other file has saved (temp is 0), in box number order, to
// libdir/misc.
func (w *world) writeLeftovers() {
	fp := w.openWriteFp("misc")

	if fp != nil {
		for i := range maxBoxes {
			if w.bx[i] != nil && w.bx[i].temp == 0 {
				if w.kind(i) != tDeleted {
					w.saveBox(fp, i)
				}
			}
		}

		fp.Close()
	}
}

// closeBox is close_box(): closefile() of libdir/fnam.
func (w *world) closeBox(fnam string) {
	path := sout("%s/%s", w.libdir, fnam)
	w.closefile(path)
}

// readBoxes is read_boxes(), pass 2 for one file: loadBox for each box of
// libdir/fnam. A file that cannot be opened gets readfile()'s diagnostic
// and is skipped.
//
// A comment line is skipped. Any other line that is not a box number
// (atoi() of it is not positive, such as a second blank line between
// boxes) gets an "unexpected line" diagnostic and is skipped. (Fixed
// upstream in https://github.com/playbymail/olympia-g3/issues/316; before,
// C never moved past such a line and looped forever.)
func (w *world) readBoxes(fnam string) {
	path := sout("%s/%s", w.libdir, fnam)
	if !w.readfile(path) {
		return
	}

	w.advance()

	for w.line != nil {
		if at(*w.line, 0) == '#' {
			w.advance() // skip comment lines
			continue
		}

		boxNum := clib.Atoi(*w.line)
		if boxNum > 0 {
			w.loadBox(boxNum)
		} else {
			fmt.Fprintf(w.stderr, "read_boxes: unexpected line %s\n", *w.line)
			w.advance()
		}
	}

	w.closeBox(fnam)
}

// readChars is read_chars(): readBoxes for every file in libdir/fact whose
// name starts with a digit, in the directory's order.
func (w *world) readChars() {
	fnam := sout("%s/fact", w.libdir)
	names, err := clib.ReadDir(fnam)

	if err != nil {
		fmt.Fprintf(w.stderr, "read_chars: can't open %s: ", fnam)
		clib.Perror(w.stderr, "", err)
		return
	}

	for _, e := range names {
		if e[0] >= '0' && e[0] <= '9' {
			lastE := sout("fact/%s", e)
			w.readBoxes(lastE)
		}
	}
}

// masterLine is master_line(): parse "num<white>kind.subkind...".
func masterLine(s string) (num, kind, sk int) {
	num = clib.Atoi(s)
	p := 0
	for p < len(s) && isdigit(s[p]) {
		p++
	}
	for p < len(s) && iswhite(s[p]) {
		p++
	}
	q := p
	for q < len(s) && s[q] != '.' {
		q++
	}
	kind = clib.Atoi(s[p:q])
	if q < len(s) {
		sk = clib.Atoi(s[q+1:])
	}
	return
}

// fastScanCheck validates every master entry before any box is allocated.
func (w *world) fastScanCheck(path string) bool {
	if !w.readfile(path) {
		return false
	}
	ok := true
	for sp := w.readlin(); ok && sp != nil; sp = w.readlin() {
		num, kind, sk := masterLine(*sp)
		if num <= 0 || num >= maxBoxes || kind <= tDeleted || kind >= tMax || sk < 0 || sk >= subMax {
			fmt.Fprintf(w.stderr, "fast_scan: bad master entry: %d %d.%d\n", num, kind, sk)
			ok = false
			break
		}
	}
	w.closefile(path)
	return ok
}

// fastScan is fast_scan(): validate master, then reopen it and allocate.
// A missing or invalid master returns false for scan_all_boxes fallback.
func (w *world) fastScan() bool {
	path := sout("%s/master", w.libdir)
	if !w.fastScanCheck(path) {
		return false
	}
	if !w.readfile(path) {
		return false
	}

	for sp := w.readlin(); sp != nil; sp = w.readlin() {
		num, kind, sk := masterLine(*sp)
		w.allocBox(num, kind, sk)
	}

	w.closefile(path)
	return true
}

// scanBoxes is scan_boxes(), the slow pass 1 for one file: allocBox for
// each box of libdir/fnam, from its "n kind subkind" line. Kind and
// subkind are names, matched without regard to case; a subkind that starts
// with '0' is 0. An unknown kind or subkind gets a diagnostic and becomes
// 0. The lines after the box's first, up to one
// that is blank or all white space, are skipped, and so are comment lines.
// A file that cannot be opened gets readfile()'s diagnostic.
func (w *world) scanBoxes(fnam string) {
	path := sout("%s/%s", w.libdir, fnam)
	if !w.readfile(path) {
		return
	}

	for sp := w.readlin(); sp != nil; sp = w.readlin() {
		s := *sp
		if at(s, 0) == '#' {
			continue // skip comment lines
		}

		// Parse something of the form: box-number kind subkind
		// example:  10 item artifact

		boxNum := clib.Atoi(s)

		i := 0
		for i < len(s) && s[i] != ' ' { // skip over space to kind
			i++
		}
		if at(s, i) == ' ' {
			i++
		}
		s = s[i:]

		j := 0
		for j < len(s) && s[j] != ' ' { // to subkind
			j++
		}
		t := ""
		if at(s, j) == ' ' {
			t = s[j+1:]
		}
		s = s[:j]

		kind := lookup(kindS, s)

		if kind < 0 {
			fmt.Fprintf(w.stderr, "scan_boxes(%d): bad kind: %s\n", boxNum, s)
			kind = 0
		}

		var sk int
		if at(t, 0) == '0' {
			sk = 0
		} else {
			sk = lookup(subkindS, t)
		}

		if sk < 0 {
			fmt.Fprintf(w.stderr, "scan_boxes(%d): bad subkind: %s\n", boxNum, t)
			sk = 0
		}

		w.allocBox(boxNum, kind, sk)

		for sp = w.readlin(); sp != nil; sp = w.readlin() {
			s := *sp
			if at(s, 0) == '#' { // skip comments
				continue
			}

			i := 0
			for i < len(s) && iswhite(s[i]) {
				i++
			}

			if i == len(s) { // blank line, end of entry
				break
			}
		}
	}
	w.closefile(path)
}

// scanChars is scan_chars(): scanBoxes for every file in
// libdir/fact whose name starts with a digit, in the directory's order.
func (w *world) scanChars() {
	fnam := sout("%s/fact", w.libdir)
	names, err := clib.ReadDir(fnam)

	if err != nil {
		fmt.Fprintf(w.stderr, "scan_chars: can't open %s: ", fnam)
		clib.Perror(w.stderr, "", err)
		return
	}

	for _, e := range names {
		if e[0] >= '0' && e[0] <= '9' {
			w.scanBoxes(sout("fact/%s", e))
		}
	}
}

// scanAllBoxes is scan_all_boxes(), the slow pass 1, which load_db() runs
// when fast_scan() fails: scanBoxes for each box file, then scanChars.
//
// Scanning through all of the entity data files, calling alloc_box for
// each entity once its number and kind are known, makes it possible to
// perform type and sanity checking when the contents of the boxes are read
// in the second pass (read_boxes).
func (w *world) scanAllBoxes() {
	w.stage("fast_scan failed, scan_all_boxes()")

	w.scanBoxes("loc")
	w.scanBoxes("item")
	w.scanBoxes("skill")
	w.scanBoxes("gate")
	w.scanBoxes("road")
	w.scanBoxes("ship")
	w.scanBoxes("unform")
	w.scanBoxes("misc")

	w.scanChars()
}

// readAllBoxes is read_all_boxes(), pass 2: readBoxes for each box file,
// then readChars.
func (w *world) readAllBoxes() {
	w.readBoxes("loc")
	w.readBoxes("item")
	w.readBoxes("skill")
	w.readBoxes("gate")
	w.readBoxes("road")
	w.readBoxes("ship")
	w.readBoxes("unform")
	w.readBoxes("misc")

	w.readChars()
}

// removeDir is remove_dir(): it removes dir and everything under it,
// without going through the shell, as G3 does since it fixed
// https://github.com/playbymail/olympia-g3/issues/320. Each entry, in
// readdir()'s order, is removed with remove_dir() if lstat() says it is a
// directory (a symbolic link is not followed), and unlink()ed otherwise;
// then dir is rmdir()ed. A dir that cannot be opened (missing, not a
// directory, unreadable) is left alone. Every error is ignored and nothing
// is printed.
func removeDir(dir string) {
	d, err := clib.Opendir(dir)

	if err == nil {
		for {
			name, ok := d.Readdir()
			if !ok {
				break
			}

			fnam := sout("%s/%s", dir, name)

			if st, err := os.Lstat(fnam); err == nil && st.IsDir() {
				removeDir(fnam)
			} else {
				_ = syscall.Unlink(fnam)
			}
		}

		d.Closedir()
		_ = syscall.Rmdir(dir)
	}
}

// saveLogdir is io.c:3044–3082. win_flag changes only the initial removal
// and parent mkdir: the two renames and final mkdir deliberately use '/'.
// mkdir/remove failures are ignored; each failed rename reports and continues.
func (w *world) saveLogdir(win bool) {
	if win {
		fmt.Fprintln(w.stdout, "Deleting any existing save dirs (this might fail):")
		removeDir(sout("%s\\save\\%d", w.libdir, w.sysclock.turn))
		_ = os.Mkdir(sout("%s\\save", w.libdir), 0o755)
	} else {
		removeDir(sout("%s/save/%d", w.libdir, w.sysclock.turn))
		_ = os.Mkdir(sout("%s/save", w.libdir), 0o755)
	}
	s, target := sout("%s/log", w.libdir), sout("%s/save/%d", w.libdir, w.sysclock.turn)
	// os.Rename rejects even an empty destination directory with EEXIST;
	// C rename must reach the syscall, including its ENOTEMPTY diagnostic.
	if err := syscall.Rename(s, target); err != nil {
		fmt.Fprintf(w.stderr, "couldn't rename %s to %s:", s, target)
		clib.Perror(w.stderr, "", err)
	}
	s, target = sout("%s/players", w.libdir), sout("%s/save/%d/players", w.libdir, w.sysclock.turn)
	if err := syscall.Rename(s, target); err != nil {
		fmt.Fprintf(w.stderr, "couldn't rename %s to %s:", s, target)
		clib.Perror(w.stderr, "", err)
	}
	_ = os.Mkdir(sout("%s/log", w.libdir), 0o755)
}

// writeAllBoxes is write_all_boxes(): it clears every box's temp mark,
// removes libdir/fact with removeDir and makes it again (mode 0755 less the
// umask), and writes the box files: loc, item, skill, gate, road, ship and
// unform by kind, fact/N by player, and misc for every box those did not
// save. Neither failing is reported; the writes that follow report what
// they cannot open. win_flag (-w) does not change any of this.
func (w *world) writeAllBoxes() {
	for i := range maxBoxes {
		if w.bx[i] != nil {
			w.bx[i].temp = 0
		}
	}

	removeDir(sout("%s/fact", w.libdir))
	_ = os.Mkdir(sout("%s/fact", w.libdir), 0o755)

	w.writeKind(tLoc, "loc")
	w.writeKind(tItem, "item")
	w.writeKind(tSkill, "skill")
	w.writeKind(tGate, "gate")
	w.writeKind(tRoad, "road")
	w.writeKind(tShip, "ship")
	w.writeKind(tUnform, "unform")

	w.writeChars()
	w.writeLeftovers()
}

// writeMaster is write_master(): libdir/master, the index fast_scan()
// reads, one line per box that is not deleted, in box number order: first
// the boxes of the kinds that have their own file, with that file's name
// (players and characters name their player's fact/N), and then every
// other box, naming misc. It clears every box's temp mark,
// and leaves it set on the boxes of the first lines.
func (w *world) writeMaster() {
	fnam := sout("%s/master", w.libdir)
	fp, err := clib.Fopen(fnam)

	if err != nil {
		fmt.Fprintf(w.stderr, "can't write %s: ", fnam)
		clib.Perror(w.stderr, "", err)
		return
	}

	for i := range maxBoxes {
		if w.bx[i] != nil {
			w.bx[i].temp = 0
		}
	}

	for i := range maxBoxes {
		if w.kind(i) != tDeleted {
			s := w.name(i)
			b := w.bx[i]

			// C's switch has a case 0 too, which is T_deleted and so
			// never reached.
			switch w.kind(i) {
			case tLoc, tItem, tSkill, tGate, tRoad, tShip, tUnform:
				fmt.Fprintf(fp, "%d\t%d.%d\t%s\t\t%s\n",
					i,
					b.kind, b.skind,
					kindS[b.kind], s)
				b.temp = 1

			case tPlayer:
				fmt.Fprintf(fp, "%d\t%d.%d\tfact/%d\t%s\n",
					i, b.kind, b.skind, i, s)
				b.temp = 1

			case tChar:
				fmt.Fprintf(fp, "%d\t%d.%d\tfact/%d\t%s\n",
					i,
					b.kind, b.skind,
					w.player(i), s)
				b.temp = 1
			}
		}
	}

	for i := range maxBoxes {
		if w.kind(i) != tDeleted && w.bx[i].temp == 0 {
			s := w.name(i)
			fmt.Fprintf(fp, "%d\t%d.%d\tmisc\t%s\n",
				i, w.bx[i].kind, w.bx[i].skind, s)
		}
	}

	fp.Close()
}

// loadSystem is load_system(): the game's settings from libdir/system,
// read with getlin_ew(). Each line sets the global its prefix names:
// "sysclock:" is scanned as olytime_scan() scans it, the host, address,
// title, URL and htpasswd lines keep the rest of the line, and the others
// are atoi() of it. Blank lines and comments are skipped, and any other
// line gets a diagnostic. A file that cannot be opened gets a diagnostic
// and leaves every global as it was.
func (w *world) loadSystem() {
	fname := sout("%s/system", w.libdir)

	data, err := clib.ReadFile(fname)

	if err != nil {
		fmt.Fprintf(w.stderr, "load_system: can't read %s: ", fname)
		clib.Perror(w.stderr, "", err)
		return
	}

	fp := bytes.NewReader(data)

	for sp := getlinEw(fp); sp != nil; sp = getlinEw(fp) {
		s := *sp
		if s == "" || s[0] == '#' { // comments and blank lines
			continue
		}

		atoi := func(n int) int32 { return int32(clib.Atoi(s[n:])) }

		switch {
		case strings.HasPrefix(s, "sysclock:"):
			olytimeScan(s[9:], &w.sysclock)
		case strings.HasPrefix(s, "from_host="):
			w.fromHost = s[10:]
		case strings.HasPrefix(s, "reply_host="):
			w.replyHost = s[11:]
		case strings.HasPrefix(s, "gm_address="):
			w.gmAddress = s[11:]
		case strings.HasPrefix(s, "game_title="):
			w.gameTitle = s[11:]
		case strings.HasPrefix(s, "game_url="):
			w.gameURL = s[9:]
		case strings.HasPrefix(s, "rules_url="):
			w.rulesURL = s[10:]
		case strings.HasPrefix(s, "times_url="):
			w.timesURL = s[10:]
		case strings.HasPrefix(s, "htpasswd_loc="):
			w.htpasswdLoc = s[13:]
		case strings.HasPrefix(s, "garrison_pay="):
			w.garrisonPay = atoi(13)
		case strings.HasPrefix(s, "auto_quit="):
			w.autoQuitTurns = atoi(10)
		case strings.HasPrefix(s, "army_slow_factor="):
			w.armySlowFactor = atoi(17)
		case strings.HasPrefix(s, "indep_player="):
			w.indepPlayer = atoi(13)
		case strings.HasPrefix(s, "gm_player="):
			w.gmPlayer = atoi(10)
		case strings.HasPrefix(s, "skill_player="):
			w.skillPlayer = atoi(13)
		case strings.HasPrefix(s, "post="):
			w.postHasBeenRun = atoi(5)
		case strings.HasPrefix(s, "init="):
			w.seedHasBeenRun = atoi(5)
		case strings.HasPrefix(s, "fr="):
			w.faeryRegion = atoi(3)
		case strings.HasPrefix(s, "tr="):
			w.tunnelRegion = atoi(3)
		case strings.HasPrefix(s, "ur="):
			w.underRegion = atoi(3)
		case strings.HasPrefix(s, "fp="):
			w.faeryPlayer = atoi(3)
		case strings.HasPrefix(s, "hr="):
			w.hadesRegion = atoi(3)
		case strings.HasPrefix(s, "hp="):
			w.hadesPit = atoi(3)
		case strings.HasPrefix(s, "hl="):
			w.hadesPlayer = atoi(3)
		case strings.HasPrefix(s, "nr="):
			w.nowhereRegion = atoi(3)
		case strings.HasPrefix(s, "nl="):
			w.nowhereLoc = atoi(3)
		case strings.HasPrefix(s, "np="):
			w.npcPl = atoi(3)
		case strings.HasPrefix(s, "cr="):
			w.cloudRegion = atoi(3)
		case strings.HasPrefix(s, "cp="):
			w.combatPl = atoi(3)
		case strings.HasPrefix(s, "mo="):
			w.mountOlympus = atoi(3)
		default:
			fmt.Fprintf(w.stderr, "load_system: unrecognized "+
				"line: %s\n", s)
		}
	}
}

// saveSystem is save_system(): the game's settings to libdir/system, in
// the order and form loadSystem reads them. auto_quit and mo are written
// only when positive, and sysclock only when not all zero. A file that
// cannot be opened gets a diagnostic.
func (w *world) saveSystem() {
	fname := sout("%s/system", w.libdir)

	fp, err := clib.Fopen(fname)

	if err != nil {
		fmt.Fprintf(w.stderr, "save_system: can't write %s: ", fname)
		clib.Perror(w.stderr, "", err)
		return
	}

	olytimePrint(fp, "sysclock: ", &w.sysclock)
	fmt.Fprintf(fp, "indep_player=%d\n", w.indepPlayer)
	fmt.Fprintf(fp, "gm_player=%d\n", w.gmPlayer)
	fmt.Fprintf(fp, "skill_player=%d\n", w.skillPlayer)
	fmt.Fprintf(fp, "from_host=%s\n", w.fromHost)
	fmt.Fprintf(fp, "reply_host=%s\n", w.replyHost)
	fmt.Fprintf(fp, "gm_address=%s\n", w.gmAddress)
	fmt.Fprintf(fp, "game_title=%s\n", w.gameTitle)
	fmt.Fprintf(fp, "game_url=%s\n", w.gameURL)
	fmt.Fprintf(fp, "rules_url=%s\n", w.rulesURL)
	fmt.Fprintf(fp, "times_url=%s\n", w.timesURL)
	fmt.Fprintf(fp, "htpasswd_loc=%s\n", w.htpasswdLoc)
	fmt.Fprintf(fp, "garrison_pay=%d\n", w.garrisonPay)
	if w.autoQuitTurns > 0 {
		fmt.Fprintf(fp, "auto_quit=%d\n", w.autoQuitTurns)
	}
	fmt.Fprintf(fp, "army_slow_factor=%d\n", w.armySlowFactor)
	fmt.Fprintf(fp, "post=%d\n", w.postHasBeenRun)
	fmt.Fprintf(fp, "init=%d\n", w.seedHasBeenRun)
	fmt.Fprintf(fp, "fr=%d\n", w.faeryRegion)
	fmt.Fprintf(fp, "tr=%d\n", w.tunnelRegion)
	fmt.Fprintf(fp, "ur=%d\n", w.underRegion)
	fmt.Fprintf(fp, "fp=%d\n", w.faeryPlayer)
	fmt.Fprintf(fp, "hr=%d\n", w.hadesRegion)
	fmt.Fprintf(fp, "hp=%d\n", w.hadesPit)
	fmt.Fprintf(fp, "hl=%d\n", w.hadesPlayer)
	fmt.Fprintf(fp, "nr=%d\n", w.nowhereRegion)
	fmt.Fprintf(fp, "nl=%d\n", w.nowhereLoc)
	fmt.Fprintf(fp, "np=%d\n", w.npcPl)
	fmt.Fprintf(fp, "cr=%d\n", w.cloudRegion)
	fmt.Fprintf(fp, "cp=%d\n", w.combatPl)
	if w.mountOlympus > 0 {
		fmt.Fprintf(fp, "mo=%d\n", w.mountOlympus)
	}

	fp.Close()
}

// deleteDeadchars is delete_deadchars(), the last step of load_db() before
// print_box_usage(): every dead character, in chain order, is taken out of
// its location if it is still in one, and deleted.
func (w *world) deleteDeadchars() {
	for i := range w.loopKind(tDeadchar) {
		// Loc should have been zeroed already in kill_char

		if w.loc(i) != 0 {
			w.setWhere(i, 0)
		}
		w.deleteBox(i)
	}
}

// cleanupPosts is cleanup_posts(), which save_db() runs first: every post,
// in chain order, is taken out of its location and deleted.
func (w *world) cleanupPosts() {
	for i := range w.loopKind(tPost) {
		w.setWhere(i, 0)
		w.deleteBox(i)
	}
}

// loadDb is load_db(): it reads the game in libdir, in C's order, checks
// it and readies it for a turn.
//
// For a world that does not exist yet, load_db() runs the steps that
// build one (post_production(), seed_initial_locations(), the create_*
// region builders, seed_phase_two() and creating the combat player, with
// read_pw()). They are not ported yet: the first one a load reaches stops
// the engine with a "not yet implemented" message, after any stage() line
// C prints before it.
func (w *world) loadDb() {
	w.stage("load_db()")

	// Assertions to verify the sanity of the linehash macro
	// Switch the byte ordering if this fails

	assert(linehash("ab f") == 'a'<<8|'b', "olympia/io.c", 2891, "linehash(\"ab f\") == 'ab'")
	assert(linehash("") == 0, "olympia/io.c", 2892, "linehash(\"\") == 0")
	assert(linehash("na") == 'n'<<8|'a', "olympia/io.c", 2893, "linehash(\"na\") == 'na'")
	assert(linehash("ab ") == 'a'<<8|'b', "olympia/io.c", 2894, "linehash(\"ab \") == 'ab'")

	w.loadSeed(sout("%s/randseed", w.libdir))
	w.initTradeRoutes() // issue #46: per-game trade-route buyer secret
	w.loadSystem()

	if !w.fastScan() { // pass 1: call alloc_box for each entity
		w.scanAllBoxes()
	}

	w.readAllBoxes() // pass 2: read the entity attributes
	w.loadOrders()

	w.checkDb()           // check database integrity
	w.determineMapEdges() // initialization for map routines

	if w.postHasBeenRun == 0 {
		w.stage("INIT: post_production()")
		panic(notImplemented("post_production"))
	}

	if w.seedHasBeenRun == 0 {
		w.stage("INIT: seed_initial_locations()")
		panic(notImplemented("seed_initial_locations"))
	}

	if w.faeryRegion == 0 {
		panic(notImplemented("create_faery"))
	}

	if w.hadesRegion == 0 {
		panic(notImplemented("create_hades"))
	}

	if w.nowhereRegion == 0 {
		panic(notImplemented("create_nowhere"))
	}

	if w.cloudRegion == 0 {
		panic(notImplemented("create_cloudlands"))
	}

	if w.tunnelRegion == 0 {
		panic(notImplemented("create_tunnels"))
	}

	// seed_phase_two() runs only when seed_has_been_run is FALSE, which
	// has stopped the load above.

	if w.combatPl == 0 {
		panic(notImplemented("creating the combat player"))
	}

	w.createRelics()
	w.initOceanChars()
	w.deleteDeadchars()

	w.printBoxUsage()
}

// saveDb is save_db(): it writes the game back to libdir.
func (w *world) saveDb() {
	w.stage("save_db()")
	w.cleanupPosts()
	w.saveSeed(sout("%s/randseed", w.libdir))
	w.saveTradeRoutes() // issue #46: persist per-game seed if customized
	w.saveSystem()
	w.writeAllBoxes()
	w.writeMaster()
	w.saveOrders()
	w.renameActJoinFiles()
}

// loadSeed is load_seed() in olympia/rnd.c: the generator's state from
// fnam, which also becomes the master seed of the addressable streams. A
// file that cannot be opened leaves the state as it was, with a message on
// stdout.
func (w *world) loadSeed(fnam string) {
	if err := w.rnd.Load(fnam); err != nil {
		fmt.Fprintf(w.stdout, "%s could not be opened.\n", fnam)
	}
}

// saveSeed is save_seed() in olympia/rnd.c: the generator's state to
// fnam. A file that cannot be written gets a diagnostic, and the engine
// exits with status 1.
func (w *world) saveSeed(fnam string) {
	if err := w.rnd.Save(fnam); err != nil {
		fmt.Fprintf(w.stderr, "couldn't write %s:", fnam)
		clib.Perror(w.stderr, "", err)
		panic(exitStatus(1))
	}
}
