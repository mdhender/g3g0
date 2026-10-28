// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"strings"

	"github.com/mdhender/g3g0/internal/clib"
)

// Ported from olympia/code.c: the entity codes, the names, print_box_usage,
// and the allocation of boxes and the kind and subkind chains, and of the
// name helpers built on display_name, including display_kind and
// box_name_kind, plus new_ent and the shared mint allocator.
//
// G2 entity coding system:
//
//	range             extent  use
//	1-999                999  items
//	1000-8999           8000  chars
//	9000-9999           1000  skills
//
//	10,000-19,999     10,000  provinces       (CCNN: AA00-DV99)
//	20,000-49,999     30,000  more provinces  (CCNN: DW00-ZZ99)
//	50,000-56,759       6760  player entities (CCN)
//	56,760-58,759       2000  lucky locs      (CNN)
//	58,760-58,999        240  regions         (NNNNN)
//	59,000-78,999     20,000  sublocs, misc   (CNNN: A000-Z999)
//	79,000-100,000    21,000  storms          (NNNNN)
//
// Note: restricted alphabet, no vowels (except a) or l:
// "abcdfghjkmnpqrstvwxz"

const (
	letters  = "abcdefghijklmnopqrstuvwxyz"
	letters2 = "abcdfghjkmnpqrstvwxz"
)

// The C locale ctype tests code_to_int uses. On the macOS reference
// platform a negative (high-bit) char is neither a digit nor a letter.
func cIsDigit(c byte) bool { return '0' <= c && c <= '9' }
func cIsAlpha(c byte) bool { return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' }
func cToLower(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		return c + 'a' - 'A'
	}
	return c
}

// letterVal is letter_val(): the index of c in let, or 0 if it is not
// there.
func letterVal(c byte, let string) int {
	if i := strings.IndexByte(let, c); i >= 0 {
		return i
	}
	return 0 // error
}

// intToCode is int_to_code(), the code an entity number is written as.
func intToCode(l int) string {
	if l < 10000 {
		return fmt.Sprintf("%d", l)
	}

	if l < 50000 { // CCNN
		l -= 10000

		n := l % 100
		l /= 100

		a := l % 20
		b := l / 20

		return fmt.Sprintf("%c%c%02d", letters2[b], letters2[a], n)
	}

	if l < 56760 { // CCN
		l -= 50000

		n := l % 10
		l /= 10

		a := l % 26
		b := l / 26

		return fmt.Sprintf("%c%c%d", letters[b], letters[a], n)
	} else if l < 58760 { // CNN
		l -= 56760

		n := l % 10
		l /= 10

		a := l % 10
		b := l / 10

		return fmt.Sprintf("%c%d%d", letters2[b], a, n)
	}

	if l < 59000 {
		return fmt.Sprintf("%d", l)
	}

	if l < 79000 {
		l -= 59000

		a := l / 1000
		n := l % 1000

		return fmt.Sprintf("%c%03d", letters2[a], n)
	}

	return fmt.Sprintf("%d", l)
}

// codeToInt is code_to_int(), the entity number a code stands for, or 0 if
// it is not a code. A code that starts with a digit is read with atoi(), so
// "12ab" is 12. A letter outside the restricted alphabet counts as 'a'
// where C uses letter_val().
func codeToInt(s string) int {
	if i := strings.IndexByte(s, 0); i >= 0 {
		s = s[:i] // a C string ends at the first NUL
	}
	at := func(i int) byte { // s[i], or C's terminating NUL
		if i < len(s) {
			return s[i]
		}
		return 0
	}

	if cIsDigit(at(0)) {
		return clib.Atoi(s)
	}

	if !cIsAlpha(at(0)) {
		return 0
	}

	switch len(s) {
	case 3:
		if cIsAlpha(s[1]) && cIsDigit(s[2]) { // CCN
			a := int(cToLower(s[0]) - 'a')
			b := int(cToLower(s[1]) - 'a')
			c := int(s[2] - '0')

			return a*260 + b*10 + c + 50000
		}

		if cIsDigit(s[1]) && cIsDigit(s[2]) { // CNN
			a := letterVal(cToLower(s[0]), letters2)
			b := int(s[1] - '0')
			c := int(s[2] - '0')

			return a*100 + b*10 + c + 56760
		}

		return 0

	case 4: // CCNN
		if cIsAlpha(s[1]) && cIsDigit(s[2]) && cIsDigit(s[3]) {
			a := letterVal(cToLower(s[0]), letters2)
			b := letterVal(cToLower(s[1]), letters2)
			c := int(s[2] - '0')
			d := int(s[3] - '0')

			return a*2000 + b*100 + c*10 + d + 10000
		}

		if cIsDigit(s[1]) && cIsDigit(s[2]) && cIsDigit(s[3]) {
			a := letterVal(cToLower(s[0]), letters2)
			b := int(s[1] - '0')
			c := int(s[2] - '0')
			d := int(s[3] - '0')

			return a*1000 + b*100 + c*10 + d + 59000
		}
		return 0

	default:
		return 0
	}
}

// printBoxUsageSup is print_box_usage_sup(), counting an inclusive range.
func (w *world) printBoxUsageSup(low, high int, s string) {
	used, unused := 0, 0

	for i := low; i <= high; i++ {
		if w.bx[i] != nil {
			used++
		} else {
			unused++
		}
	}

	fmt.Fprintf(w.stderr, "\t%5d - %5d  %4s  %d/%d used (%d%%)\n",
		low, high,
		s,
		used, high-low+1, used*100/(high-low+1))
}

// printBoxUsage is print_box_usage(), which load_db() writes to stderr.
func (w *world) printBoxUsage() {
	fmt.Fprintf(w.stderr, "entity space usage:\n")

	w.printBoxUsageSup(1000, 9999, "char")
	w.printBoxUsageSup(10000, 49999, "locs")
	w.printBoxUsageSup(50000, 56759, "play")
	w.printBoxUsageSup(56760, 58759, "lloc")
	w.printBoxUsageSup(59000, 78999, "subl")
	w.printBoxUsageSup(79000, maxBoxes-1, "rest")
}

// scode is scode(): code_to_int() of s, less one leading '[' or '('. The
// closing bracket is not removed, so "[ab12]" is 0 but "[1234]" is 1234.
func scode(s string) int {
	if len(s) > 0 && (s[0] == '[' || s[0] == '(') {
		s = s[1:]
	}

	return codeToInt(s)
}

// name is name(): the name of box n, or "" if it has none.
func (w *world) name(n int) string {
	assert(w.validBox(n), "olympia/code.c", 228, "valid_box(n)")

	if w.bx[n].name != nil {
		return *w.bx[n].name
	}

	return ""
}

// setName is set_name(): it names box n s, with each '[' changed to '{'
// and each ']' to '}'. The name is set even when s is "".
func (w *world) setName(n int, s string) {
	assert(w.validBox(n), "olympia/code.c", 241, "valid_box(n)")

	t := strings.NewReplacer("[", "{", "]", "}").Replace(s)
	w.bx[n].name = &t
}

// setBanner is set_banner(): it sets box n's display banner to s, cut to 50
// bytes, or clears it if s is "" (C's NULL or ""). Either way it allocates
// the entity_misc. C cuts the caller's string in place; Go strings cannot
// be changed, and the one caller (v_banner) rejects a banner longer than 50
// bytes first.
func (w *world) setBanner(n int, s string) {
	p := w.pMisc(n)

	p.display = nil

	if s != "" {
		if len(s) > 50 {
			s = s[:50]
		}
		p.display = &s
	}
}

// boxCodeLess is box_code_less(): box_code() without the brackets.
func boxCodeLess(n int) string {
	return intToCode(n)
}

// boxCode is box_code(): "[code]", or "Garrison" for garrison_magic.
func (w *world) boxCode(n int) string {
	if n == int(w.garrisonMagic) {
		return "Garrison"
	}

	return "[" + intToCode(n) + "]"
}

// displayName is display_name(): box n's name, or if it has none, a name
// from what it is: "Player", "Gate" or "Sign" for those kinds, the item a
// noble is made of, or its subkind, with the first letter capitalized. It
// is "" for a box that is not valid.
func (w *world) displayName(n int) string {
	if !w.validBox(n) {
		return ""
	}

	s := w.name(n)

	if s != "" {
		return s
	}

	switch w.kind(n) {
	case tPlayer:
		return "Player"
	case tGate:
		return "Gate"
	case tPost:
		return "Sign"
	}

	if i := w.nobleItem(n); i != 0 {
		return capitalize(w.pluralItemName(i, 1))
	}

	return capitalize(subkindS[w.subkind(n)])
}

func (w *world) displayKind(n int) string {
	switch w.subkind(n) {
	case subCity:
		if w.isPortCity(n) {
			return "port city"
		}
		return "city"
	case subFog, subRain, subWind:
		return "storm"
	default:
		return subkindS[w.subkind(n)]
	}
}

func (w *world) boxNameKind(n int) string {
	return sout("%s, %s", w.boxName(n), w.displayKind(n))
}

// boxName is box_name(): box n's display name and code, as "Name~[code]",
// or only the code if it has no display name or is not valid, or
// "Garrison" for garrison_magic.
func (w *world) boxName(n int) string {
	if n == int(w.garrisonMagic) {
		return "Garrison"
	}

	if w.validBox(n) {
		s := w.displayName(n)
		if s != "" {
			return sout("%s~%s", s, w.boxCode(n))
		}
	}

	return w.boxCode(n)
}

// justName is just_name(): box n's display name, or its code if it has
// none or is not valid, or "Garrison" for garrison_magic.
func (w *world) justName(n int) string {
	if n == int(w.garrisonMagic) {
		return "Garrison"
	}

	if w.validBox(n) {
		s := w.displayName(n)
		if s != "" {
			return s
		}
	}

	return w.boxCode(n)
}

// pluralItemName is plural_item_name(): the display name of item for one,
// or its plural name for any other qty. An item with no plural name is
// reported on stderr and gets its display name. The item's box must exist.
func (w *world) pluralItemName(item, qty int) string {
	if qty == 1 {
		return w.displayName(item)
	}

	s := ""
	if p := w.rpItem(item); p != nil && p.pluralName != nil {
		s = *p.pluralName
	}

	if s == "" {
		fmt.Fprintf(w.stderr, "warning: plural name not set for item %s\n", w.boxCode(item))
		s = w.displayName(item)
	}

	return s
}

// justNameQty is code.c's just_name_qty, used by wait item predicates.
func (w *world) justNameQty(item, qty int) string {
	return sout("%s~%s", niceNum(qty), w.pluralItemName(item, qty))
}

func (w *world) pluralItemBox(item, qty int) string {
	if qty == 1 {
		return w.boxName(item)
	}
	return sout("%s~%s", w.pluralItemName(item, qty), w.boxCode(item))
}

func (w *world) boxNameQty(item, qty int) string {
	return sout("%s~%s", niceNum(qty), w.pluralItemBox(item, qty))
}

// Allocating entities and threading like entities together (kind_first,
// kind_next).
//
// addNextChain and addSubChain cache the last box they added, so that boxes
// added in increasing order are appended without a walk. removeNextChain
// and removeSubChain drop the cache when they remove that box: a box added
// after it must not be linked after a box that is no longer on the chain
// (G3 #246). The cache is the static state of code.c, kept in world
// (nextCacheLast and the rest).

// addNextChain is add_next_chain(): it links box n into the chain of its
// kind, in increasing order. A box of kind 0 is not linked.
func (w *world) addNextChain(n int) {
	assert(w.bx[n] != nil, "olympia/code.c", 496, "bx[n] != NULL")
	kind := int(w.bx[n].kind)

	if kind == 0 {
		return
	}

	// optim!

	if int(w.nextCacheKind) == kind &&
		n > int(w.nextCacheLast) &&
		w.bx[w.nextCacheLast].xNextKind == 0 {
		w.bx[w.nextCacheLast].xNextKind = int32(n)
		w.bx[n].xNextKind = 0
		w.nextCacheLast = int32(n)
		return
	}

	w.nextCacheLast = int32(n)
	w.nextCacheKind = int32(kind)

	if w.boxHead[kind] == 0 {
		w.boxHead[kind] = int32(n)
		w.bx[n].xNextKind = 0
		return
	}

	if n < int(w.boxHead[kind]) {
		w.bx[n].xNextKind = w.boxHead[kind]
		w.boxHead[kind] = int32(n)
		return
	}

	i := w.boxHead[kind]
	for w.bx[i].xNextKind > 0 && int(w.bx[i].xNextKind) < n {
		i = w.bx[i].xNextKind
	}

	w.bx[n].xNextKind = w.bx[i].xNextKind
	w.bx[i].xNextKind = int32(n)
}

// removeNextChain is remove_next_chain(): it unlinks box n from the chain
// of its kind. A box not on that chain fails the assertion, as C does.
func (w *world) removeNextChain(n int) {
	assert(w.bx[n] != nil, "olympia/code.c", 545, "bx[n] != NULL")

	i := w.boxHead[w.bx[n].kind]

	if int(i) == n {
		w.boxHead[w.bx[n].kind] = w.bx[n].xNextKind
	} else {
		for i > 0 && int(w.bx[i].xNextKind) != n {
			i = w.bx[i].xNextKind
		}

		assert(i > 0, "olympia/code.c", 558, "i > 0")

		w.bx[i].xNextKind = w.bx[n].xNextKind
	}

	w.bx[n].xNextKind = 0

	if n == int(w.nextCacheLast) {
		w.nextCacheLast = 0
		w.nextCacheKind = 0
	}
}

// addSubChain is add_sub_chain(): it links box n into the chain of its
// subkind, in increasing order. Unlike the kind chain, subkind 0 has a
// chain.
func (w *world) addSubChain(n int) {
	assert(w.bx[n] != nil, "olympia/code.c", 579, "bx[n] != NULL")
	kind := int(w.bx[n].skind)

	// optim!

	if int(w.subCacheKind) == kind &&
		n > int(w.subCacheLast) &&
		w.bx[w.subCacheLast].xNextSub == 0 {
		w.bx[w.subCacheLast].xNextSub = int32(n)
		w.bx[n].xNextSub = 0
		w.subCacheLast = int32(n)
		return
	}

	w.subCacheLast = int32(n)
	w.subCacheKind = int32(kind)

	if w.subHead[kind] == 0 {
		w.subHead[kind] = int32(n)
		w.bx[n].xNextSub = 0
		return
	}

	if n < int(w.subHead[kind]) {
		w.bx[n].xNextSub = w.subHead[kind]
		w.subHead[kind] = int32(n)
		return
	}

	i := w.subHead[kind]
	for w.bx[i].xNextSub > 0 && int(w.bx[i].xNextSub) < n {
		i = w.bx[i].xNextSub
	}

	w.bx[n].xNextSub = w.bx[i].xNextSub
	w.bx[i].xNextSub = int32(n)
}

// removeSubChain is remove_sub_chain(): it unlinks box n from the chain of
// its subkind.
func (w *world) removeSubChain(n int) {
	assert(w.bx[n] != nil, "olympia/code.c", 625, "bx[n] != NULL")

	i := w.subHead[w.bx[n].skind]

	if int(i) == n {
		w.subHead[w.bx[n].skind] = w.bx[n].xNextSub
	} else {
		for i > 0 && int(w.bx[i].xNextSub) != n {
			i = w.bx[i].xNextSub
		}

		assert(i > 0, "olympia/code.c", 638, "i > 0")

		w.bx[i].xNextSub = w.bx[n].xNextSub
	}

	w.bx[n].xNextSub = 0

	if n == int(w.subCacheLast) {
		w.subCacheLast = 0
		w.subCacheKind = -1
	}
}

// deleteBox is delete_box(): it unlinks box n from its chains and makes it
// T_deleted. The box itself stays in bx[], keeping its subkind and
// substructures, so it cannot be allocated again; deleting it twice fails
// an assertion.
func (w *world) deleteBox(n int) {
	w.removeNextChain(n)
	w.removeSubChain(n)
	w.bx[n].kind = tDeleted
}

// changeBoxKind is change_box_kind(). The kind is stored as a schar.
func (w *world) changeBoxKind(n, kind int) {
	w.removeNextChain(n)
	w.bx[n].kind = int8(kind)
	w.addNextChain(n)
}

// changeBoxSubkind is change_box_subkind(). The subkind is stored as a
// schar.
func (w *world) changeBoxSubkind(n, sk int) {
	if w.subkind(n) == sk {
		return
	}

	w.removeSubChain(n)
	w.bx[n].skind = int8(sk)
	w.addSubChain(n)
}

// allocBox is alloc_box(): it creates box n with the given kind and
// subkind (each stored as a schar) and links it into their chains. A box
// number out of range, or one already in use, fails an assertion; a box in
// use is reported on stderr first.
func (w *world) allocBox(n, kind, sk int) {
	assert(n > 0 && n < maxBoxes, "olympia/code.c", 690, "n > 0 && n < MAX_BOXES")

	if w.bx[n] != nil {
		fmt.Fprintf(w.stderr, "alloc_box: DUP box %d\n", n)
		assert(false, "olympia/code.c", 695, "FALSE")
	}

	w.bx[n] = &box{kind: int8(kind), skind: int8(sk)}
	w.addNextChain(n)
	w.addSubChain(n)
}

// rndAllocNum is code.c's one-draw, inclusive upper scan then wrap. Deleted
// boxes still occupy their slots; valid_box is deliberately not used here.
func (w *world) rndAllocNum(low, high int) int {
	n := w.streams.MintAlloc(low, high)
	for i := n; i <= high; i++ {
		if w.bx[i] == nil {
			return i
		}
	}
	for i := low; i < n; i++ {
		if w.bx[i] == nil {
			return i
		}
	}
	return -1
}

func (w *world) newEnt(kind, sk int) int {
	n := -1
	switch kind {
	case tPlayer:
		n = w.rndAllocNum(50000, 56759)
	case tChar, tUnform:
		if sk == subNi {
			n = w.rndAllocNum(79000, maxBoxes-1)
		} else {
			n = w.rndAllocNum(1000, 9999)
		}
		if n < 0 {
			n = w.rndAllocNum(59000, 78999)
		}
		if n < 0 {
			n = w.rndAllocNum(79000, maxBoxes-1)
		}
	case tSkill:
		assert(false, "olympia/code.c", 864, "FALSE")
	case tLoc:
		switch sk {
		case subCity:
			n = w.rndAllocNum(56760, 58759)
		case subRegion:
			n = w.rndAllocNum(58760, 58999)
		case subUnder, subForest, subOcean, subCloud, subTunnel:
			n = w.rndAllocNum(20000, 49999)
		default:
			n = w.rndAllocNum(59000, 78999)
		}
	case tStorm:
		n = w.rndAllocNum(79000, maxBoxes-1)
	default:
		n = w.rndAllocNum(59000, 78999)
	}
	if n < 0 {
		n = w.rndAllocNum(79000, maxBoxes-1)
	}
	w.allocBox(n, kind, sk)
	return n
}
