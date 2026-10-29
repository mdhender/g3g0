// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

// These are Go port boundaries, not invented C gameplay outcomes. The oracle
// checks all C NULL bits; every other absent implementation must remain NYI.
func TestCommandUseCallbackBoundaries(t *testing.T) {
	implemented := map[int][2]func(*world, *command) int{
		skHideSelf:      {(*world).vHide, (*world).dHide},
		skBindStorm:     {(*world).vBindStorm, (*world).dBindStorm},
		skHinderMed:     {(*world).vHinderMed, (*world).dHinderMed},
		skQuarryStone:   {(*world).vQuarry, nil},
		skCatchHorse:    {(*world).vCatch, nil},
		skHarvestLumber: {(*world).vWood, nil},
		skHarvestYew:    {(*world).vYew, nil},
		skRaiseCorpses:  {(*world).vRaiseCorpses, nil},
		skMageMenial:    {(*world).vMageMenial, nil},
		skFishing:       {(*world).vFish, nil},
		skTrainWild:     {(*world).vUseTrainRiding, nil},
		skTrainWarmount: {(*world).vUseTrainWar, nil},
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
	// Shared dispatch must be wired to its real owners, not the NYI callbacks.
	row := cmdTbl[newWorld().findCommand("use")]
	for phase, callback := range []func(*world, *command) int{row.start, row.finish, row.interrupt} {
		real := []func(*world, *command) int{(*world).vUse, (*world).dUse, (*world).iUse}[phase]
		if reflect.ValueOf(callback).Pointer() != reflect.ValueOf(real).Pointer() {
			t.Fatalf("USE phase %d not wired", phase)
		}
	}
}

// These tests assert only explicit Go port boundaries. They are deliberately
// not recordings of C successes with missing gameplay replaced by a stub.
func TestCommandUseSharedNYIBoundaries(t *testing.T) {
	for _, row := range useTbl[1 : len(useTbl)-1] {
		if row.start == nil || row.skill == skHideSelf || row.skill == skBindStorm || row.skill == skHinderMed {
			continue
		}
		switch row.skill {
		case skQuarryStone, skCatchHorse, skHarvestLumber, skHarvestYew, skRaiseCorpses, skMageMenial, skFishing,
			skTrainWild, skTrainWarmount:
			continue // Real reparse/dispatch outcomes are covered by the C oracle.
		}
		t.Run(fmt.Sprint(row.skill), func(t *testing.T) {
			w := newWorld()
			w.stdout, w.stderr = io.Discard, io.Discard
			w.libdir = t.TempDir()
			w.openLogfileNondestruct()
			t.Cleanup(w.closeLogfile)
			w.allocBox(200, tPlayer, subPlRegular)
			w.allocBox(1000, tChar, 0)
			w.allocBox(50000, tSkill, 0)
			w.allocBox(row.skill, tSkill, 0)
			w.pSkill(row.skill).requiredSkill = 50000
			w.pChar(1000).unitLord = 200
			w.pChar(1000).skills.Append(&skillEnt{skill: 50000, know: skillKnow})
			w.pChar(1000).skills.Append(&skillEnt{skill: int32(row.skill), know: skillKnow})
			c := w.pCommand(1000)
			c.who = 1000
			line := fmt.Sprintf("use %d", row.skill)
			if !w.olyParse(c, &line) {
				t.Fatal("parse failed")
			}
			path := ""
			w.coveragePath, w.useCoveragePath = &path, &path
			defer func() {
				if _, ok := recover().(notImplemented); !ok {
					t.Fatal("missing callback did not stop NYI")
				}
				if w.covCounts != nil || w.ucovSkills != nil {
					t.Fatal("NYI fabricated a start outcome")
				}
				if w.pPlayer(200).cmdCount != 1 || int(c.useEnt) != findUseEntry(row.skill) || c.h != int32(row.skill) {
					t.Fatal("shared pre-callback mutation timing differs")
				}
			}()
			w.doCommand(c)
		})
	}
}

func TestCommandUseItemNYIBoundaries(t *testing.T) {
	for key := 1; key <= useKeyMax; key++ {
		if key == useDrum {
			continue // The real drum callback is covered by the C command oracle.
		}
		t.Run(fmt.Sprint(key), func(t *testing.T) {
			w := newWorld()
			w.allocBox(1000, tChar, 0)
			w.allocBox(60000, tItem, 0)
			w.pItemMagic(60000).useKey = int8(key)
			w.genItem(1000, 60000, 1)
			switch key {
			case useDeathPotion, useHealPotion, useSlavePotion, useFaeryStone, useBtaSkull:
				// These keys are not magic-safe-haven rejects in C. Missing
				// gameplay must still stop NYI, not invent a shared-check fail.
				w.allocBox(9000, tLoc, subPlain)
				w.pLocInfo(1000).where = 9000
				w.pSubloc(9000).safe = 1
			}
			path := ""
			w.useCoveragePath = &path
			c := &command{who: 1000, wait: 99, poll: 1}
			line := "use 60000 7"
			if !w.olyParse(c, &line) {
				t.Fatal("parse failed")
			}
			defer func() {
				if _, ok := recover().(notImplemented); !ok {
					t.Fatal("item gameplay did not stop NYI")
				}
				if w.ucovSkills != nil || w.ucovItems != [useKeyMax + 1][ucovMax]int32{} {
					t.Fatal("NYI fabricated a use outcome")
				}
				if c.wait != 0 || c.poll != 0 || c.inhibitFinish != 0 {
					t.Fatal("item mutation crossed the callback boundary")
				}
				if c.a != 60000 {
					t.Fatal("unported item argument changed")
				}
			}()
			w.vUse(c)
		})
	}
}
