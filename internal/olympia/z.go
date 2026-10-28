// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/mdhender/g3g0/internal/clib"
	"github.com/mdhender/g3g0/internal/rnd"
	"github.com/mdhender/g3g0/internal/rng"
)

// testRandom is test_random() in olympia/z.c, the -R self-test. It draws
// from a local stream rooted at the master seed. -R runs before load_db(),
// so no seed has been loaded and the master seed is all zero.
//
// On a terminal the first 100 draws print as a 10x10 grid, each row drawn in
// order before it prints; otherwise they print one per line.
func testRandom(out *cStdout, tty bool) {
	var seed rnd.Source // never loaded: rng_master_seed() is all zero
	s := rng.Seed(seed.MasterSeed())

	if tty {
		for range 10 {
			var r [10]int
			for j := range r {
				r[j] = s.Draw(1, 10)
			}
			out.Printf("%3d  %3d  %3d  %3d  %3d  %3d  %3d  %3d  %3d  %3d\n",
				r[0], r[1], r[2], r[3], r[4], r[5], r[6], r[7], r[8], r[9])
		}
	} else {
		for range 100 {
			out.Printf("%d\n", s.Draw(1, 10))
		}
	}

	for i := -10; i >= -16; i-- {
		out.Printf("rng_draw(%d, %d) == %d\n", i, -3, s.Draw(i, -3))
	}

	for range 100 {
		out.Printf("%d\n", s.Draw(1000, 9999))
	}

	l := make([]int, 0, 10)
	for i := 1; i <= 10; i++ {
		l = append(l, i)
	}
	rng.Shuffle(l, &s)

	out.Printf("Scramble:\n")
	for _, v := range l {
		out.Printf("%d\n", v)
	}
}

// assertion is a failed assert(): asfail() in olympia/z.c prints it to
// stderr and calls abort(). The port panics with it, and whatever runs the
// engine recovers it, prints Error() to stderr and exits as abort() does:
// with status 134, after flushing every open stdio stream, as asfail()
// explicitly does before abort() (stdout and the files clib.FlushAll flushes).
//
// file stands in for the __FILE__ that C prints, which is whatever path the
// C build compiled the source as (an absolute path in the CMake build);
// the port uses the path in the G3 repository, "olympia/code.c".
type assertion struct {
	file string
	line int
	cond string
}

// Error is the line asfail() prints.
func (a assertion) Error() string {
	return fmt.Sprintf("assertion failure: %s (%d): %s\n", a.file, a.line, a.cond)
}

// assert is assert(cond): if ok is false it panics with an assertion for
// the C source line file:line and the condition's text.
func assert(ok bool, file string, line int, cond string) {
	if !ok {
		panic(assertion{file: file, line: line, cond: cond})
	}
}

// iswhite is the iswhite() macro of olympia/z.h: a space or a tab.
func iswhite(c byte) bool { return c == ' ' || c == '\t' }

// isdigit is the isdigit() macro of olympia/z.h, which replaces the C
// library's: '0' to '9'.
func isdigit(c byte) bool { return c >= '0' && c <= '9' }

// lenMax is LEN in olympia/z.h, the generic string max length. readlin()
// cuts each line to lenMax-1 bytes.
const lenMax = 2048

// getlinTrunc is where getlin() cuts the last line of a file when it has
// no newline (its "#if 1" block). A line that ends in a newline is not cut.
const getlinTrunc = 1023

// cString is the C string a buffer holds: b up to its first NUL byte.
func cString(b []byte) string {
	if i := bytes.IndexByte(b, 0); i >= 0 {
		b = b[:i]
	}
	return string(b)
}

// getlin is getlin() in olympia/z.c, the line reader "with no size limits":
// the next line of fp without its newline, or nil (NULL) at end of file.
// It reads with fgetc(), so a read error ends the input as end of file
// does. The last line of the file, if it has no newline, is cut to
// getlinTrunc bytes. Like every line reader here it returns the line as C
// sees it, up to its first NUL byte.
//
// C returns its static buffer, which the next call overwrites; callers that
// change the line in place (getlin_ew) change the buffer. Here each call
// returns a new string.
func getlin(fp io.ByteReader) *string {
	var buf []byte

	for {
		c, err := fp.ReadByte()
		if err != nil {
			break
		}

		if c == '\n' {
			s := cString(buf)
			return &s
		}

		buf = append(buf, c)
	}

	if len(buf) == 0 {
		return nil
	}

	if len(buf) >= getlinTrunc {
		buf = buf[:getlinTrunc]
	}

	s := cString(buf)
	return &s
}

// stripEw is the clean-up getlin_ew() and readlin_ew() apply to a line:
// leading spaces and tabs are skipped, then every byte below 32 becomes a
// space, and then trailing white space is removed.
//
// The test is *p < 32 on a plain char, which is signed on the reference
// builds (macOS, arm64 and x86_64), so every byte from 0x80 up is also
// below 32 and becomes a space. Leading white space is skipped before the
// bytes are changed, so a leading control byte becomes a space that stays.
func stripEw(s string) string {
	i := 0
	for i < len(s) && iswhite(s[i]) {
		i++ // eat leading whitespace
	}

	b := []byte(s[i:])
	for j, c := range b {
		if int8(c) < 32 || c == '\t' { // remove ctrl chars
			b[j] = ' '
		}
	}

	j := len(b)
	for j > 0 && iswhite(b[j-1]) {
		j-- // eat trailing whitespace
	}

	return string(b[:j])
}

// getlinEw is getlin_ew(): getlin() with leading and trailing white space
// removed and control characters made spaces (see stripEw).
func getlinEw(fp io.ByteReader) *string {
	line := getlin(fp)

	if line != nil {
		s := stripEw(*line)
		line = &s
	}

	return line
}

// eatLeadingTrailingWhitespace is eat_leading_trailing_whitespace(): s
// without its leading and trailing spaces and tabs. C trims the trailing
// white space in place, in the caller's buffer, and returns a pointer past
// the leading white space; callers here use the result instead.
func eatLeadingTrailingWhitespace(s string) string {
	return strings.TrimRight(strings.TrimLeft(s, " \t"), " \t")
}

// maxBuf is MAX_BUF in olympia/z.c: readfile() and readlin() read the file
// in read() calls of this size.
const maxBuf = 8192

// lineReader is the state of the readfile()/readlin() line reader of
// olympia/z.c: the file that readfile() opened (linebuf, nread, line_fd
// and point in C). It holds the file's whole contents, read at once; C
// reads them maxBuf bytes at a time, which gives the same lines.
type lineReader struct {
	data  []byte // what read() will return, from the start of the file
	point int    // offset in data of C's point
}

// loaded is the offset in data just past the read() buffer that C holds:
// the end of the maxBuf chunk that point is in (or just finished), and of
// the first chunk before any line is read.
func (r *lineReader) loaded() int {
	return maxBuf * max(1, (r.point+maxBuf-1)/maxBuf)
}

// readfile is readfile() in olympia/z.c: it opens path for readlin(), or
// prints "can't open PATH: " and the error to stderr and returns false.
//
// C opens the file with open() and reads the first chunk. It never closes
// the file it had open before. A failed open sets line_fd to -1 but leaves
// the buffer as it was, so readlin() still returns what was left in the
// chunk C holds, and then NULL, when read(-1) fails. A directory opens but
// cannot be read, which C sees as an empty file.
func (w *world) readfile(path string) bool {
	data, err := clib.ReadFile(path)
	if err != nil {
		fmt.Fprintf(w.stderr, "can't open %s: ", path)
		clib.Perror(w.stderr, "", err)
		w.closefile(path)
		return false
	}

	w.lr = lineReader{data: data}

	return true
}

// closefile is closefile() in olympia/z.c: it closes the file readfile()
// opened. C resets line_fd but not the buffer, so as after a failed
// readfile(), readlin() returns what is left in the chunk C holds and then
// NULL.
func (w *world) closefile(path string) {
	r := &w.lr
	r.data = r.data[:min(len(r.data), r.loaded())]
}

// readlin is readlin() in olympia/z.c: the next line of the file readfile()
// opened, without its newline, or nil (NULL) at end of file. Every line,
// whether or not it ends in a newline, is cut to lenMax-1 bytes, and then,
// as C sees it, at its first NUL byte.
//
// C returns its static buffer, which the next call overwrites. io.c keeps
// that pointer in line and readlin_ew() changes the buffer in place, but
// io.c never looks at line again after a continuation line is read without
// first reading the next line. Here each call returns a new string.
func (w *world) readlin() *string {
	r := &w.lr

	if r.point >= len(r.data) {
		return nil
	}

	b := r.data[r.point:]
	if i := bytes.IndexByte(b, '\n'); i >= 0 {
		b = b[:i]
		r.point += i + 1
	} else {
		r.point = len(r.data)
	}

	if len(b) >= lenMax {
		b = b[:lenMax-1]
	}

	s := cString(b)
	return &s
}

// readlinEw is readlin_ew(): readlin() with leading and trailing white
// space removed and control characters made spaces (see stripEw). It is
// nil at end of file, and io.c's continuation scans then end the list, so
// a backslash on a file's last line is ignored. (Fixed upstream in
// https://github.com/playbymail/olympia-g3/issues/313; before, the scans
// dereferenced the NULL and crashed.)
func (w *world) readlinEw() *string {
	line := w.readlin()

	if line != nil {
		s := stripEw(*line)
		line = &s
	}

	return line
}

// tolower is the tolower() macro of olympia/z.h, lower_array[c] after
// init_lower(): 'A' to 'Z' become 'a' to 'z', and every other byte is
// itself. C indexes the array with a plain char, so a byte from 0x80 up
// reads before it, which is undefined; no caller here passes one (see
// olyParseCmd), and such a byte is left as it is.
func tolower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c - 'A' + 'a'
	}
	return c
}

// iStrcmp is i_strcmp(): s and t compared without regard to case. Only
// whether it is 0 matters to callers; otherwise it is the difference of
// the first bytes that differ, as plain (signed) chars.
func iStrcmp(s, t string) int {
	for i := 0; ; i++ {
		a := tolower(at(s, i))
		b := tolower(at(t, i))
		if a != b {
			return int(int8(a)) - int(int8(b))
		}
		if a == 0 {
			return 0
		}
	}
}

// fuzzyTranspose is fuzzy_transpose(): whether two, with two adjacent
// bytes swapped, is one, without regard to case.
func fuzzyTranspose(one, two string, l1, l2 int) bool {
	if l1 != l2 {
		return false
	}

	buf := []byte(two)

	for i := 0; i < l2-1; i++ {
		buf[i], buf[i+1] = buf[i+1], buf[i]

		if iStrcmp(one, string(buf)) == 0 {
			return true
		}

		buf[i], buf[i+1] = buf[i+1], buf[i]
	}

	return false
}

// fuzzyOneLess is fuzzy_one_less(): whether two is one with one byte left
// out, without regard to case.
func fuzzyOneLess(one, two string, l1, l2 int) bool {
	count := 0

	if l1 != l2+1 {
		return false
	}

	for i, j := 0, 0; j < l2; i, j = i+1, j+1 {
		if tolower(at(one, i)) != tolower(at(two, j)) {
			if count != 0 {
				return false
			}
			count++
			j--
		}
	}

	return true
}

// fuzzyOneExtra is fuzzy_one_extra(): whether two is one with one byte
// added, without regard to case.
func fuzzyOneExtra(one, two string, l1, l2 int) bool {
	count := 0

	if l1 != l2-1 {
		return false
	}

	for i, j := 0, 0; i < l1; i, j = i+1, j+1 {
		if tolower(at(one, i)) != tolower(at(two, j)) {
			if count != 0 {
				return false
			}
			count++
			i--
		}
	}

	return true
}

// fuzzyOneBad is fuzzy_one_bad(): whether two is one with at most one byte
// changed, without regard to case.
func fuzzyOneBad(one, two string, l1, l2 int) bool {
	count := 0

	if l1 != l2 {
		return false
	}

	for i := 0; i < l2; i++ {
		if tolower(one[i]) != tolower(two[i]) {
			if count != 0 {
				return false
			}
			count++
		}
	}

	return true
}

// fuzzyStrcmp is fuzzy_strcmp(): whether two is a near miss for one: two
// adjacent bytes swapped (if two has at least 4 bytes), or one byte left
// out, added or changed (if it has at least 5).
func fuzzyStrcmp(one, two string) bool {
	l1 := len(one)
	l2 := len(two)

	if l2 >= 4 && fuzzyTranspose(one, two, l1, l2) {
		return true
	}

	if l2 >= 5 && fuzzyOneLess(one, two, l1, l2) {
		return true
	}

	if l2 >= 5 && fuzzyOneExtra(one, two, l1, l2) {
		return true
	}

	if l2 >= 5 && fuzzyOneBad(one, two, l1, l2) {
		return true
	}

	return false
}

// safeQsort is safe_qsort(base, nmemb, size, comp): qsort() of l by comp,
// skipped for fewer than two elements. qsort() is not stable, and neither
// is slices.SortFunc, so the two can order elements that comp finds equal
// differently. Summary ranks deliberately have no ID tie-break: equal
// values receive the same ordinal regardless of their relative order.
func safeQsort[T any](l []T, comp func(a, b T) int) {
	if len(l) > 1 {
		slices.SortFunc(l, comp)
	}
}
