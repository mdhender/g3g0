# lists differential test

`run.sh` checks that the Go reallocing arrays (`internal/lists`) match G3's C list modules.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `47535543d83fbc5649d54b63b49e578f2b11ab1c`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.

The port follows:

- `lib/ilist.c`: every operation `lib/lists.h` declares for `ilist` and defines (`ilist_insert` is declared but never defined), and `ilist_test()`, which `olympia-g3 -t` runs.
- The fourteen modules generated from one template: `accept_ents`, `admits`, `cstrings`, `exit_views`, `fights`, `flag_ents`, `item_ents`, `orders`, `req_ents`, `roads`, `skill_ents`, `tiles`, `trades` and `wait_args`, plus `exit_views_shuffle_rng()`.
  They differ from each other only in names and element type, and from `ilist.c` only in the element type and the operations `ilist.c` adds, so one generic type, `lists.List[T]`, ports all fifteen.
  Their `<name>_test()` functions are never called by the engine and are not ported.
- `lib/checked_alloc.c` has nothing observable to port: at the pinned revision `my_malloc`, `my_realloc` and `my_free` are thin wrappers over `calloc`, `realloc` and `free` that print nothing.

The capacity and its growth policy (`ILIST_ALLOC`, doubling) are not visible to the engine and are not reproduced.
The package comment in `internal/lists/lists.go` lists the other differences between a C list and a Go slice that ports must keep in mind.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another.
- A C11 compiler (`CC`, default `cc`), Go, Bash, `git`, `cmp` and `diff`.

## Running

```sh
tools/lists-compat/run.sh            # compare
tools/lists-compat/run.sh --update   # rewrite internal/lists/testdata/c-output.txt from C
```

The script:

1. Exports `lib/` at the pinned revision with `git archive` into `tmp/lists-compat/g3`.
2. Writes a C driver to `tmp/lists-compat/driver.c` with one handler per list module, and links the exported `lib/ilist.c`, the fourteen template modules, `lib/rnd.c` and `lib/rng.c`.
3. Runs the driver twice over [`internal/lists/testdata/cases.txt`](../../internal/lists/testdata/cases.txt) and requires the same output both times.
4. Requires that output to equal [`internal/lists/testdata/c-output.txt`](../../internal/lists/testdata/c-output.txt), or with `--update`, writes it there.
5. Runs `go test -run TestMatchesC ./internal/lists`, which interprets the same cases in Go and requires the recorded C output line for line.

`go test ./...` runs `TestMatchesC` too, so the port is held to the recorded C output without a G3 clone.

## Case matrix

`cases.txt` runs in one process, so lists and the RNG state carry from line to line. The same block runs for each of the fifteen modules:

| area | cases |
|---|---|
| NULL lists | length, lookup, clear, copy, shuffle, scramble, rem_value and reclaim on a NULL list, which stays NULL and draws nothing |
| empty lists | clear after prepend keeps the list non-NULL; a copy of it is empty and non-NULL and does not share its length |
| growth | 40 appends through several reallocs, a copy, and an append that the copy does not see |
| delete | at the front, the back and the middle, then prepends |
| duplicates | lookup finds the first; rem_value of an absent value, then of values present several times, down to empty |
| shuffles | `<name>_shuffle` and `<name>_scramble` on lists of 0, 1, 2 and 40 elements, with an `rnd()` after each group to show the draws they made |
| `ilist` only | `ilist_add` of new and present values; `ilist_rem_value_uniq` removing the last match, and of absent values |
| stream shuffles | `ilist_shuffle_rng` and `exit_views_shuffle_rng` on NULL, 1-element and 40-element lists, printing the stream counter |

A delete out of range, or of a NULL list, fails C's `assert()` and aborts, so the matrix has no such case; `TestDeleteAsserts` checks that the port panics.

## Results

Last run: 2026-10-26, macOS 26.6.2 arm64, Apple clang 21.0.0, Go 1.26.4.

- 899 operations. The C driver is deterministic, its output matches `c-output.txt`, and `TestMatchesC` passes.

## Upstream defects reproduced on purpose

None are known at the pinned revision.
