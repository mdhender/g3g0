# Movement helper compatibility

Pin: `6a25f2d252ce45bd5adbe3baaccb3be3eed56279` in
[playbymail/olympia-g3](https://github.com/playbymail/olympia-g3/tree/6a25f2d252ce45bd5adbe3baaccb3be3eed56279).

```sh
G3_REPO=/home/user/workspace/repo/tmp/upstream.git tools/move-compat/run.sh
# Only actual pinned C may replace the expected bytes:
G3_REPO=/home/user/workspace/repo/tmp/upstream.git tools/move-compat/run.sh --update
```

`G3_REPO` defaults to the workspace's bare `tmp/upstream.git`. `G3_REV`, `WORK`,
`CC` and `RUN_TIMEOUT` may be overridden. Source is read only via `git archive`
into `tmp/move-compat/g3`; nothing builds in an upstream working tree.

The harness extracts **complete original definitions**, under original `#line`
locations, for dir.c's topology/distance/selection closure, loc.c's region and
province walks, u.c's depth/hidden/knowledge helpers, and swear.c's ownership
walks. It preserves dir.c's original includes, tag constants, stream declaration,
turn guard and `begin_exitdir`. It links original code.c/glob.c/z.c and the actual
lib checked allocator, ilist, **exit_views.c**, rng.c and rnd.c. Original headers
provide structs, macros and integer field widths. There are **no gameplay stubs**.
The driver only constructs fixtures, invokes helpers and dumps every exit field.

Each section runs in a fresh process, twice, under `run_limited`. Successful
cases must exit 0, assertion cases 134, and the separately identified known
G3 #353 case 139. Both runs must match raw stdout, stderr and status, then match
`internal/olympia/testdata/move-c-output.txt`. The Go recorded-case test uses
the absolute run-limited adapter. Default checks need no whole world.

Coverage:

- Asymmetric terrain distances, depth swaps, buildings, Hades pit, ocean links
  and signed sea-lane flags; missing and zero direction slots.
- Province/sublocation/building/ship/sibling-ship exits; ordering of directions,
  inner locations, links and roads; recursive port insertion and first-city choice.
- Hidden known/unknown destinations and roads, OUT exceptions, barriers,
  collapsed mines, regional/Hades metadata, and the distinct road rules.
- Vertical cloud routes, permanent Uldim blocks, seasonal Uldim/Summerbridge
  boundaries, and month-done behavior.
- Ship departure day/evening, differing sibling timestamps, negative and wrapped
  C-int timestamp arithmetic; completed and unfinished ships and dry-dock links.
- LAND/WATER/both/empty filters; sequential repeated shuffles, turn reseeding,
  day changes, empty/singleton no-draw controls and unrelated legacy-RNG draws.
- Original diagnostic/assertion bytes and assertion flushing. The known
  [G3 #353](https://github.com/playbymail/olympia-g3/issues/353) unowned-character
  NULL dereference retains SIGSEGV's unflushed output and status 139.

Weights/travel calculations belong to #105, parsing/admission/bookkeeping to
#106, and move callbacks/arrivals to #67. This harness does not claim gameplay
execution, reports, or raft crew/travel acceptance; those obligations remain in
#105/#67. Report rendering and unrelated LOS-distance helpers are not extracted.
