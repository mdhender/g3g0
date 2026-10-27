// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "unsafe"

// cBox mirrors the layout of struct box in olympia/oly.h, field for field,
// so that unsafe.Sizeof gives sizeof(struct box) under the host's C ABI.
// main() prints that size, and the Go engine has to print the same bytes as
// the C build it is compared with: a C build for the same host. The types
// stand in for the C ones by size and alignment only (schar is int8, every
// pointer and ilist is a pointer, int is int32), which Go lays out with the
// same padding rules as the C compilers G3 supports. On the 64-bit (LP64)
// reference hosts this is 144; on a 32-bit host it would be 80.
//
// It is a layout description only. The entity model is ported separately;
// keep this in step with struct box when the pinned revision changes it.
type cBox struct {
	kind  int8     // schar kind
	skind int8     // schar skind
	name  *byte    // char *name
	where int32    // struct loc_info x_loc_info: int where
	here  *int32   // ... ilist here_list
	ents  [9]*byte // x_player, x_char, x_loc, x_subloc, x_item, x_skill, x_gate, x_misc, x_disp
	cmd   *byte    // struct command *cmd
	items *byte    // item_ents_list items
	trade *byte    // trades_list trades

	temp        int32 // int temp
	outputOrder int32 // int output_order
	nextKind    int32 // int x_next_kind
	nextSub     int32 // int x_next_sub
}

// sizeofBox is sizeof(struct box) as main() prints it.
const sizeofBox = unsafe.Sizeof(cBox{})
