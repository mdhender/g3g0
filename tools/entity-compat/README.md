# entity model differential test

`run.sh` checks that the Go entity model (`internal/olympia`: `entity.go`, `oly.go`, `glob.go`, `world.go`) matches G3's `olympia/oly.h` and `olympia/glob.c`.
G3 is the reference.
Expected results always come from running a C driver built from the pinned sources; nothing here is re-baselined from Go output.

## Upstream revision

[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3) at `47535543d83fbc5649d54b63b49e578f2b11ab1c`.
The revision is pinned in `run.sh` (`G3_REV`). Change it only on purpose, and rerun the test when you do.

The port follows:

- `olympia/oly.h`: every struct (and the `olytime` typedef), every numeric `#define`, and the `p_xxx`/`rp_xxx` accessors, `kind()`, `subkind()` and `valid_box()`.
  The phase coverage enum (`PH_xxx`) and the other helper macros (`char_health()` and so on) are left for the tasks that need them.
- `olympia/glob.c`: the globals that loading, saving and `check_db()` read or write (as fields of `world`), `glob_init()` (as `newWorld()`), and the tables `kind_s`, `subkind_s` and `cmd_tbl`.
  `cmd_tbl`'s command routines are nil until the commands are ported.
  `show_day`, `short_dir_s`, `full_dir_s`, `exit_opposite`, `loc_depth_s` and `month_names` are not ported: loading and saving do not use them. `libdir` is ported by `options.libdir`.
- The globals saved in `libdir/system` but defined outside `glob.c` (`seed_has_been_run`, `faery_region`, `tunnel_region` and the rest) are fields of `world` too. They start at zero; the driver cannot link the files that define them, so it does not print them.

## Prerequisites

- A G3 clone that contains the pinned revision. The default is `~/Software/playbymail/olympia-g3`; set `G3_REPO` to use another.
- A C11 compiler (`CC`, default `cc`), Go, Bash, Python 3, `git`, `cmp` and `diff`.

## Running

```sh
tools/entity-compat/run.sh            # compare
tools/entity-compat/run.sh --update   # rewrite internal/olympia/testdata/entity-c-output.txt from C
```

The script:

1. Exports `lib/` and `olympia/` at the pinned revision with `git archive` into `tmp/entity-compat/g3`.
2. Reads `oly.h` for its numeric `#define`s and the fields of its structs, and `glob.c` for the command routines `cmd_tbl` names, and generates `consts.h`, `structs.h` and stubs for the routines.
3. Writes a C driver to `tmp/entity-compat/driver.c` and links it with the exported `olympia/glob.c`, the stubs and `lib/checked_alloc.c`.
4. Runs the driver twice and requires the same output both times.
5. Requires that output to equal [`internal/olympia/testdata/entity-c-output.txt`](../../internal/olympia/testdata/entity-c-output.txt), or with `--update`, writes it there.
6. Runs `go test -run TestEntityMatchesC ./internal/olympia`, which prints the same lines from the Go model and requires the recorded C output line for line.

`go test ./...` runs `TestEntityMatchesC` too, so the port is held to the recorded C output without a G3 clone.
The struct layouts are those of the recording host (LP64); on a host with another pointer size the test skips.

## What the driver prints

| section | lines |
|---|---|
| constants | `const NAME VALUE` for each numeric `#define` in `oly.h`; every one must be ported with the same value, and nothing else |
| structs | `sizeof`, `_Alignof`, and each field's `offsetof` and size, for every struct; the Go side derives them from the Go types with `cLayout`, so this checks field names, order and widths, and `sizeof(struct box)`, which `main()` prints |
| tables | `kind_s`, `subkind_s` and `cmd_tbl` (allow, name, time, poll, pri), up to the NULL terminator |
| globals | the `glob.c` globals after `glob_init()`, and that `bx[]`, `box_head[]` and `sub_head[]` start empty |
| accessors | `kind()`, `valid_box()` and `subkind()` at and beyond the ends of the table; each `p_xxx` allocates a zeroed substructure once, each `rp_xxx` never allocates, and `p_magic()` and `p_item_magic()` allocate the `entity_char` or `entity_item` they hang from. Pointers print as labels in order of first appearance |

## Results

Last run: 2026-10-26, macOS 26.6.2 arm64, Apple clang 21.0.0, Go 1.26.4.

- 1090 lines: 430 constants, 24 structs, the three tables, the globals and 84 accessor results. The C driver is deterministic, its output matches `entity-c-output.txt`, and `TestEntityMatchesC` passes.

## Upstream defects reproduced on purpose

None. Two macros in `oly.h` name fields that no struct has, `char_persuaded()` (`persuaded`) and `skill_aura_req()` (`aura_req`); C compiles only because nothing uses them, and they are not ported.
