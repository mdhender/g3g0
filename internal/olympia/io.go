// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"io"
	"reflect"
	"slices"

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
// Where C reads a continuation line at the end of the file, readlin_ew()
// returns NULL and C dereferences it, which crashes the engine. The ports
// dereference the nil *string that readlinEw returns, which panics.

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

// lastIsBackslash is the test s[strlen(s)-1] == '\\' that the list scans
// use to see whether another entry follows. On an empty s, C reads the byte
// before it, which is undefined: when s is t_string()'s "" it is whatever
// precedes that literal, and when s is an empty continuation line it is the
// white space readlin_ew() skipped, or the byte before the line buffer.
// Here an empty s ends the list. Upstream defect:
// https://github.com/playbymail/olympia-g3/issues/312
func lastIsBackslash(s string) bool {
	return len(s) > 0 && s[len(s)-1] == '\\'
}

// sscanD is sscanf(s, "%d %d ...") with one conversion for each of ptrs,
// each an *int32, *int16 (for %hd) or *int. Like sscanf, it stops at the
// first conversion that fails and leaves that target and the rest as they
// were. It returns the number of conversions made.
//
// Each conversion is clib.Scanner.Int, which follows the macOS libc
// (FreeBSD's vfscanf), truncated to the target's width as %hd stores the
// low bits of the converted value.
func sscanD(s string, ptrs ...any) int {
	sc := clib.NewScanner(s)
	for i, p := range ptrs {
		v, ok := sc.Int()
		if !ok {
			return i
		}
		switch p := p.(type) {
		case *int32:
			*p = int32(v)
		case *int16:
			*p = int16(v)
		case *int:
			*p = v
		default:
			panic(fmt.Sprintf("olympia: sscanD: unsupported target %T", p))
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
// readlin_ew() wherever s has a backslash. Anything else ends the list. A
// sign is not a digit, so "-5" ends it too.
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
			s = *w.readlinEw()
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
// the caller's list. int_comp() (olympia/seed.c) returns *a - *b, which
// overflows if two elements are more than INT_MAX apart; a known list holds
// box numbers, which never are, so the order is ascending.
func sortKnown(kn sparse) {
	if kn.Len() > 1 {
		slices.SortFunc(kn, func(a, b int) int { return int(int32(a) - int32(b)) })
	}
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
// C reads know, experience and the fifth number into uninitialized locals,
// so a line with fewer than four numbers gives know or experience whatever
// was on the stack, which is undefined. Here they start at 0 and keep their
// values from one entry to the next, as the locals of an unoptimized build
// do. Upstream defect:
// https://github.com/playbymail/olympia-g3/issues/311
func (w *world) skillListScan(s string, l *skillEntsList, boxNum int) {
	var dummy, know, experience int

	for {
		p := &skillEnt{}
		sscanD(s, &p.skill, &know, &p.daysStudied, &experience, &dummy)

		p.know = int8(know)
		p.experience = int16(experience)

		if w.validBox(int(p.skill)) {
			l.Append(p)
		} else {
			fmt.Fprintf(w.stderr, "skill_list_scan(%d): bad skill %d\n", boxNum, p.skill)
		}

		if lastIsBackslash(s) { // another entry follows
			s = *w.readlinEw()
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
		sscanD(s, &p.item, &p.qty)

		if w.validBox(int(p.item)) {
			l.Append(p)
		} else {
			fmt.Fprintf(w.stderr, "item_list_scan(%d): bad item %d\n", boxNum, p.item)
		}

		if lastIsBackslash(s) { // another entry follows
			s = *w.readlinEw()
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
// if its item is a valid box. Each trade's who is boxNum.
func (w *world) tradeListScan(s string, l *tradesList, boxNum int) {
	for {
		p := &trade{}
		sscanD(s,
			&p.kind,
			&p.item,
			&p.qty,
			&p.cost,
			&p.cloak,
			&p.haveLeft,
			&p.monthProd,
			&p.expire)

		p.who = int32(boxNum)

		if w.validBox(int(p.item)) {
			l.Append(p)
		} else {
			fmt.Fprintf(w.stderr, "trade_list_scan(%d): bad item %d\n", boxNum, p.item)
		}

		if lastIsBackslash(s) { // another entry follows
			s = *w.readlinEw()
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
// As in skillListScan, consume is an uninitialized local in C, undefined
// on a line with fewer than three numbers; here it starts at 0 and keeps
// its value from one entry to the next. Upstream defect:
// https://github.com/playbymail/olympia-g3/issues/311
func (w *world) reqListScan(s string, l *reqEntsList, boxNum int) {
	var consume int

	for {
		p := &reqEnt{}
		sscanD(s, &p.item, &p.qty, &consume)
		p.consume = int8(consume)

		if w.validBox(int(p.item)) {
			l.Append(p)
		} else {
			fmt.Fprintf(w.stderr, "req_list_scan(%d): bad item %d\n", boxNum, p.item)
		}

		if lastIsBackslash(s) { // another entry follows
			s = *w.readlinEw()
		} else {
			break
		}
	}
}

// olytimeScan is olytime_scan(): "turn day days" from s into p. Fields
// that s does not give keep their values.
func olytimeScan(s string, p *olytime) {
	sscanD(s, &p.turn, &p.day, &p.daysSinceEpoch)
}

// olytimePrint is olytime_print(): header and "turn day days", unless all
// three are 0.
func olytimePrint(fp io.Writer, header string, p *olytime) {
	if p.turn != 0 || p.day != 0 || p.daysSinceEpoch != 0 {
		fmt.Fprintf(fp, "%s%d %d %d\n", header,
			p.turn, p.day, p.daysSinceEpoch)
	}
}
