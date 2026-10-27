// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Only the saved-command collector list is ported. Harvesting and COLLECT's
// callbacks remain explicit unported boundaries, not successful substitutes.
func (w *world) initCollectList() {
	cmdCollect := w.findCommand("collect")
	assert(cmdCollect > 0, "olympia/produce.c", 533, "cmd_collect > 0")
	for i := range w.loopKind(tChar) {
		c := w.rpCommand(i)
		if c != nil && c.state == stateRun && int(c.cmd) == cmdCollect {
			w.collectors.Append(i)
		}
	}
}
