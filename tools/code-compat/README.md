# entity allocation and codes differential test

`run.sh` checks that the Go port's box allocation, kind and subkind chains, entity codes and names (`internal/olympia`: `code.go`, `loop.go`, `io.go`'s `ifName`) match G3's `olympia/code.c`, the entity loops of `olympia/loop.h` and `if_name()` in `olympia/io.c`.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `6a25f2d252ce45bd5adbe3baaccb3be3eed56279`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.
This advances the former retained pin (unchanged code-tool inputs) to the
stage pin before exercising the new allocator functions. `tools/affected.sh`
for that move does not list code-compat; its previous 20 sections / 652 lines
remain byte-identical.

The port follows, from `olympia/code.c`: `letter_val`, `int_to_code`, `code_to_int`, `print_box_usage`, `scode`, `name`, `set_name`, `set_banner`, `box_code`, `box_code_less`, the chains (`add_next_chain`, `remove_next_chain`, `add_sub_chain`, `remove_sub_chain` and their cache), `delete_box`, `change_box_kind`, `change_box_subkind` and `alloc_box`.
From `olympia/oly.h`: `kind_first`, `kind_next`, `sub_first`, `sub_next` and `loc`.
From `olympia/loop.h`: `loop_boxes`, `loop_kind`, `loop_subkind` and `loop_loc_or_ship`.
Also ported: `begin_mint`, `rnd_alloc_num` and `new_ent` (code.c:706–903).
The existing `internal/rng.Streams` implements the mint tag, turn guard and
shared sequential/keyed APIs; each world owns one persistent instance,
constructed before load with live clock/master callbacks. Loading the seed
does not reset an already-seeded turn. `display_name` and the name helpers built on it that `check_db` uses (`box_name`, `just_name`, `plural_item_name`) are checked by [`helper-compat`](../helper-compat). `display_kind`, `box_name_kind` and the quantity helpers also run in [`command-compat`](../command-compat#bounded-namescodes-audit-137); #137 audits their remaining fallback/composition paths without adding another naming implementation.

For #137, all upstream inputs read by this tool (including headers and linked
lib files) are unchanged between this retained pin and stage-2's
`e733acb63bacb21167a1f2589d8dcdc023364896`; code.c has identical blob
`1a9996f1c7ea3da603a35a3c534a00549c4375a7`. The independent code-tool pin
and its 43-section / 969-line recording deliberately remain unchanged.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another.
- A C11 compiler (`CC`, default `cc`), Go, Bash, Python 3, `git`, `cmp` and `diff`.

## Running

```sh
tools/code-compat/run.sh            # compare
tools/code-compat/run.sh --update   # rewrite internal/olympia/testdata/code-c-output.txt from C
```

The script:

1. Exports `lib/` and `olympia/` at the pinned revision with `git archive` into `tmp/code-compat/g3`.
2. Generates **aborting** stubs for the unused command/display dependencies,
   the existing formatting-only `sout()` ring, and copies `if_name()` (static
   in `io.c`) out of the export. Reaching an aborting stub fails even `--update`.
3. Writes a C driver to `tmp/code-compat/driver.c` and links it with the exported `olympia/code.c`, `olympia/glob.c`, `olympia/z.c`, `lib/checked_alloc.c`, `lib/ilist.c`, `lib/rng.c` and `lib/rnd.c`. A generated translation unit includes the **unchanged actual code.c** and exposes its static scan/mint through trivial wrappers; it copies no allocator/RNG logic. It compiles from the export's root, so `asfail()` prints `olympia/code.c` as the file, as the Go port does.
4. Runs each section of [`internal/olympia/testdata/code-cases.txt`](../../internal/olympia/testdata/code-cases.txt) in a fresh process, since a section may end in a failed assertion, which aborts. stderr goes to stdout, in order, and `exit STATUS` follows each section. It does this twice and requires the same output both times.
5. Requires that output to equal [`internal/olympia/testdata/code-c-output.txt`](../../internal/olympia/testdata/code-c-output.txt), or with `--update`, writes it there.
6. Runs `go test -run TestCodeMatchesC ./internal/olympia`, which runs the same cases on the Go port and requires the recorded C output line for line.

`go test ./...` runs `TestCodeMatchesC` too, so the port is held to the recorded C output without a G3 clone.

## What the cases cover

| section | covers |
|---|---|
| codes | `int_to_code` at the edges of every range; `code_to_int` and `scode` on valid, malformed and high-bit input; `letter_val`; and checksums of `int_to_code` and back over every box number from -10 to `MAX_BOXES`+10, of `code_to_int` and `scode` over every string of up to four bytes from a 19-byte alphabet, and of `code_to_int` over every three- and four-byte code shape in both cases |
| allocation, deletion, changes | chain order after allocation in increasing, decreasing and mixed order across kinds and subkinds (which exercises the chain caches), after `delete_box` of the head, middle and tail and allocation after it, and after `change_box_kind` and `change_box_subkind`, including no-op changes and kind 0 |
| loops | `loop_kind` and `loop_subkind` while the body deletes the current box or allocates one before or after it; `loop_boxes`; `loop_loc_or_ship` with and without ships |
| names | `name`, `set_name` (brackets become braces), `set_banner` (cut to 50 bytes, cleared by an empty string), `if_name` with and without `pretty_data_files`, `box_code` and `box_code_less` |
| usage | `print_box_usage` on an empty world and with boxes at the ends of each range |
| assertions | `alloc_box` of a box in use (deleted or not) or out of range, `delete_box` twice or of a missing box, changing the kind or subkind of a deleted box, and `name`/`set_name` of a missing or deleted box: the stderr text and exit status 134 |
| mint ranges and scans | all `new_ent` kind/subkind branches; inclusive upper scan, wrap to low, singleton, sparse, full and deleted-occupied ranges; every fallback level, total exhaustion and forbidden T_skill; sorted kind/subkind chains after randomized minting |
| mint RNG isolation | INIT and signed-short turn boundaries, same-turn/day progression, turn reset/revisit, keyed passwords vs sequential city/shuffle/allocator draws, load-master changes inside/across turns, and controls with/without legacy rnd advancement; raw master/digest bytes |

Exhaustion cases enable `trace_abort`: only **after real asfail** a
single-threaded SIGABRT observer draws two probes from the actual mint stream,
then returns to abort (still status 134). Go's test observer does the same
after recovering the actual assertion. These are harness diagnostics, not
successful allocation or a replacement assertion: they discriminate the
otherwise invisible repeated fallback draws (four for exhausted char/NI,
two for storm, zero for T_skill). No production hook or RNG model is added.

Discriminating recordings: the unconstrained `[1200,1209]` start is 1202,
so one-hole cases prove upper scan to 1209 and wrap to 1200. Legacy/control
draws both give 3377, 4707, 3543 then turn-1 7748 despite differing saved
digest bytes. Reloading master in the same turn keeps second draw 4707;
on turn 1 the changed master gives 5708 instead of unchanged-master 7748.
Revisiting turn 1 mints 60589 instead of its previously occupied 60588.
Character/item state and actual `save_box` serialization are checked in
[`command-compat`](../command-compat).

## Results

Last run: 2026-10-27, Linux amd64, GCC 12.2.0, Go 1.26.4.

- 43 sections, 969 lines. Two actual C runs are byte-identical, their output matches `code-c-output.txt`, and `TestCodeMatchesC` passes. Defaults use only these small recordings; no whole-world turn.

## Corrected historical diagnostic note

The former final section alleged outstanding `print_box_usage` range and
denominator defects. They were fixed in [G3 #310 / #340,
73243150ef4ac0e3cbe7d9427263726a9534bcb3](https://github.com/playbymail/olympia-g3/commit/73243150ef4ac0e3cbe7d9427263726a9534bcb3).
Pinned code.c:178–209, the Go implementation and the retained usage recording
already agree: locs 10000–49999, rest 79000–MAX_BOXES-1, denominator
`high-low+1`. No new bug or engine behavior change is called for.
