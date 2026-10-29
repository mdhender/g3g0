// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// These are Go port boundaries, not invented C gameplay outcomes. The oracle
// checks all C NULL bits; every other absent implementation must remain NYI.
func TestCommandUseCallbackBoundaries(t *testing.T) {
	implemented := map[int][2]func(*world, *command) int{
		skHideSelf:  {(*world).vHide, (*world).dHide},
		skBindStorm: {(*world).vBindStorm, (*world).dBindStorm},
	}
	for i, row := range useTbl {
		for phase, callback := range []func(*world, *command) int{row.start, row.finish, row.interrupt} {
			if callback == nil {
				continue
			}
			t.Run(fmt.Sprintf("%d/%d", i, phase), func(t *testing.T) {
				if real, ok := implemented[row.skill]; ok && phase < 2 {
					if reflect.ValueOf(callback).Pointer() != reflect.ValueOf(real[phase]).Pointer() {
						t.Fatal("did not reuse existing gameplay callback")
					}
					return
				}
				defer func() {
					r, ok := recover().(notImplemented)
					if !ok || !strings.HasPrefix(string(r), "command USE ") {
						t.Fatalf("missing gameplay callback = %v, want NYI", r)
					}
				}()
				callback(nil, nil)
			})
		}
	}
	// Neither shared USE dispatch nor its item dispatcher is enabled by metadata.
	row := cmdTbl[newWorld().findCommand("use")]
	for phase, callback := range []func(*world, *command) int{row.start, row.finish, row.interrupt} {
		t.Run(fmt.Sprintf("command/%d", phase), func(t *testing.T) {
			defer func() {
				if _, ok := recover().(notImplemented); !ok {
					t.Fatal("USE command prematurely enabled")
				}
			}()
			callback(nil, nil)
		})
	}
}
