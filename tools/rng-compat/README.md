# rng differential test

`run.sh` checks that the Go RNG streams (`internal/rng`, with the master seed and `md5_int()` in `internal/rnd`) match G3's C code.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `9a03d126b329ed466c4af48425553c025eaca26f`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.
The line ranges in `run.sh` (`SNIPPETS`) are for that revision; check them when the revision changes.

For #291, the USE guard moves the complete unchanged skills RNG block from
use.c **1623–1740 to 1630–1747**. Every extracted RNG block is byte-identical
to the independent old e733acb pin; no assertion remapping is introduced.
Two real C runs and the fresh Go default retain all **321 operations** and
the entire recording unchanged. See [merged-pin proof](../README.md#merged-pin-proof-291)
for the isolated source, bounded commands and acceptance limits.

For #131/G3 #357, swear.c's social prefix moves from 428–528 to 425–525.
The complete original definitions (begin_social through soc_breed) are
byte-identical; all 321 deterministic C results remain unchanged.

The port follows:

- `lib/rng.c`: `rng_seed`, `rng_stream_of`, `rng_draw`, `rng_keyed`.
- `lib/rnd.c` and `olympia/rnd.c`: the master seed that `load_seed()` captures, `rng_master_seed()`, and `md5_int()`.
- `lib/ilist.c` `ilist_shuffle_rng()` and `lib/exit_views.c` `exit_views_shuffle_rng()`, as one generic `rng.Shuffle`.
- The named streams in `olympia/*.c`: every `begin_*()` seeder and the helpers that draw from its stream, as methods on `rng.Streams`.
  The two draws that C writes inline (the `auto_bandit` victim pick in `npc.c` and the `random_hades_loc` pick in `hades.c`) are `NPCVictim` and `HadsRandomLoc`.
  The `-R` self-test in `z.c` (`test_random`) is left to the `cmd/olympia` task.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another.
- A C11 compiler (`CC`, default `cc`), Go, Bash, `git`, `cmp` and `diff`.
- A little-endian host. The C code hashes native 32-bit words, and the port reads and writes them little-endian.

## Running

```sh
tools/rng-compat/run.sh            # compare
tools/rng-compat/run.sh --update   # rewrite internal/rng/testdata/c-output.txt from C
```

The script:

1. Exports `lib/` and `olympia/` at the pinned revision with `git archive` into `tmp/rng-compat/g3`.
2. Cuts the named-stream code out of the exported `olympia/*.c` files at the line ranges in `SNIPPETS`: the purpose tags, the static streams and their turn guards, and the seeders and helpers, verbatim.
3. Writes a C driver to `tmp/rng-compat/driver.c` that includes those snippets and links the exported `lib/rng.c`, `lib/rnd.c` and `lib/ilist.c`.
   The driver defines `sysclock` with the engine's `olytime` layout (short `day` and `turn`).
4. Runs the driver twice over [`internal/rng/testdata/cases.txt`](../../internal/rng/testdata/cases.txt) and requires the same output both times.
5. Requires that output to equal [`internal/rng/testdata/c-output.txt`](../../internal/rng/testdata/c-output.txt), or with `--update`, writes it there.
6. Runs `go test -run TestMatchesC ./internal/rng`, which interprets the same cases in Go and requires the recorded C output line for line.

`go test ./...` runs `TestMatchesC` too, so the port is held to the recorded C output without a G3 clone.
`TestRngCheck` ports the assertions in G3's `tests/rng/rng_check.c`.

## Case matrix

`cases.txt` runs in one process, so state carries from line to line as it does in the engine. It covers:

| area | cases |
|---|---|
| primitives | the `rng_check.c` master and its draws; keys at 0, -1, `INT_MAX`, `INT_MIN` and 2^32-1; tags at 0, 1, `0x80000000` and `0xffffffff`; one-value ranges at `INT_MIN` and `INT_MAX`, `[0, INT_MAX]`, negative ranges, and small ranges that reject often |
| shuffles | lists of 0 and 1 (no draws), 2, 3, 10 and 52 |
| master seed | seeds of 16, 2 and 0 bytes (a short seed keeps the rest of the state), a missing file (the master becomes the state as advanced by `rnd()`), `rnd()` between loads, and `md5_int()` |
| named streams | every seeder and helper; reseeding a per-battle, per-pillage, per-quest, per-market and per-build stream; sequential streams whose counters carry between calls (weather, mint, exit directions) |
| zero state | every kind of helper before anything seeds it: unseeded streams are all zero, and a turn guard that already matches (turn -1) never seeds |
| clock | a second day of the same turn (the per-day leaves change, the per-turn leaves repeat), a new turn, going back to an earlier turn, and turns and days that overflow a `short` |
| master changes mid-turn | streams reseeded on every call pick up the new master; turn-guarded streams keep the stream from earlier in the turn until the turn changes |
| key arithmetic | `target*2+sub`, `where*16+sub`, `item*4+sub` and `prisoner*2+sub` overflowing a C `int`, and a negative item folded into `up_animal`'s tag |

C's `reduce_range()`, like `rnd()`, never returns when `high < low` or when `high - low` overflows an `int`; the port does the same, so the matrix has no such case.

## Results

Last run: 2026-10-26, macOS 26.6.2 arm64, Apple clang 21.0.0, Go 1.26.4.

- 321 operations. The C driver is deterministic, its output matches `c-output.txt`, and `TestMatchesC` passes.

## Upstream defects reproduced on purpose

None are known at the pinned revision.
