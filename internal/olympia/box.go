// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"reflect"
	"unsafe"
)

// sizeofBox is sizeof(struct box) as main() prints it. The Go engine has to
// print the same bytes as the C build it is compared with, a C build for the
// same host, so it is the size of struct box under the host's C ABI, derived
// from the box type by cSizeof. On the 64-bit (LP64) reference hosts it is
// 144; on a 32-bit host it would be 80.
var sizeofBox = cSizeof(reflect.TypeFor[box]())

// cPtrSize is the size and alignment of a C pointer on the host.
const cPtrSize = unsafe.Sizeof(uintptr(0))

// cLayout gives the size and alignment that a field of Go type t has in the
// C struct it ports, under the rules of entity.go: int8, uint8, int16 and
// int32 are the C integer types of the same width, and every pointer, list,
// string, slice and func stands for a C pointer. A struct is laid out field
// by field with C's padding rules, which all the C compilers G3 supports
// follow. Any other type is a mistake in a ported struct and panics.
func cLayout(t reflect.Type) (size, align uintptr) {
	switch t.Kind() {
	case reflect.Int8, reflect.Uint8:
		return 1, 1
	case reflect.Int16:
		return 2, 2
	case reflect.Int32:
		return 4, 4
	case reflect.Pointer, reflect.Slice, reflect.String, reflect.Func:
		return cPtrSize, cPtrSize
	case reflect.Struct:
		align = 1
		for f := range t.Fields() {
			fs, fa := cLayout(f.Type)
			size = (size + fa - 1) &^ (fa - 1)
			size += fs
			align = max(align, fa)
		}
		return (size + align - 1) &^ (align - 1), align
	}
	panic(fmt.Sprintf("olympia: cLayout: no C type for %s", t))
}

// cSizeof is sizeof() for the C struct that t ports.
func cSizeof(t reflect.Type) uintptr {
	size, _ := cLayout(t)
	return size
}
