// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// vHide and dHide are the complete hide handlers from stealth.c. They
// provide real failure, timed completion and zero-time inhibited completion
// paths for dispatch without implementing any part of wait or scheduling.
func (w *world) vHide(c *command) int {
	if !w.checkSkill(int(c.who), skHideSelf) {
		return 0
	}
	if c.a != 0 && (w.stackParent(int(c.who)) != 0 || w.countStackAny(int(c.who)) != 1) {
		w.commandOut(int(c.who), "Must be alone to hide.")
		return 0
	}
	if c.a == 0 {
		w.pMagic(int(c.who)).hideSelf = 0
		w.commandOut(int(c.who), "No longer hidden.")
		c.wait = 0
		c.inhibitFinish = 1
	}
	return 1
}

func (w *world) dHide(c *command) int {
	if w.stackParent(int(c.who)) != 0 || w.countStackAny(int(c.who)) != 1 {
		w.commandOut(int(c.who), "Must be alone to hide.")
		return 0
	}
	w.pMagic(int(c.who)).hideSelf = 1
	w.commandOut(int(c.who), "Now hidden.")
	return 1
}
