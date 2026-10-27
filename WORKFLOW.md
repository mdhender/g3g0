# Project Workflow

Guidance for agents when working in this repository.

## Priorities

Fix bugs in our port before adding features, subject to the upstream policy below.
Before offering to implement a change or work on documentation, check for open bugs (`gh issue list --label bug --state open`).
If any are open, do not make the offer; propose working on the bugs instead.

## Build and test

Prefer `go run` over `go build` for your work.
Do not create build artifacts in the root of the repository - put them in tmp/.

```sh
go vet ./...
go test ./...
```

Both must pass before committing.
`go test ./...` skips the whole-game tests; [tools/README.md](tools/README.md) says which compat checks to run for a task, a pin move and a stage exit.

## Versioning

- Semantic versioning.
  The current version lives in `version.go` as a `semver.Version` literal (`Major`, `Minor`, `Patch`, `PreRelease`).
- **Every commit that changes code must bump the version in `version.go`.**
  Documentation-only commits are encouraged to bump it but are not required to.
- Choose the bump by the change: patch for fixes, minor for new features or any breaking change while below v1.0.0.
  Keep `PreRelease: "alpha"` until told otherwise.
- Do not create git tags.
  Versions are tagged on GitHub by the user.

## Git workflow

- Assign all issues and PRs to `@me`.
- Commit directly to `main`.
  No feature branches or pull requests; history is the recovery mechanism.
  This may change after v1.0.0.
- Agents are authorized to commit and push each task as it is completed, provided all checks pass and code changes include a version bump.
- Agents are authorized to create GitHub labels as needed.
- Agents are authorized to create issues in the upstream github.com/playbymail/olympia-g3 repository.

## Upstream behavior, bugs, and harness tests

- [Olympia-G3](https://github.com/playbymail/olympia-g3) drives the port.
  Implement the game that its code implements, including known bugs; do not introduce independent fixes or behavior changes.
- File bugs discovered in upstream code in the [G3 repository](https://github.com/playbymail/olympia-g3/issues), not as local deviations.
  Do not modify their code or implement fixes for them.
- If an upstream bug blocks the port, stop the affected work and wait for upstream to implement a fix.
  Do not work around it locally.
- If an upstream bug is not blocking, reproduce the buggy code and behavior in the Go port.
  Document the intentional duplication in a code comment with a link to the upstream issue.
- Bug fixes, additional harness tests, and intentional behavior changes must go through the G3 repository first.
  Open an upstream issue when additional harness coverage is needed; do not modify their harness for them or change our expected behavior independently.
- When upstream implements a fix, they should open a bug in this repository so we can update the port.
  If no local issue is opened, failing harness tests should reveal the mismatch.
  Follow the upstream fix and update the port and its issue references to match; do not close an upstream issue merely because our port changes.

## Port stages

The engine port is tracked in #19, one stage at a time.

1. Review #19 for the stage's scope and create a sub-issue of #19 for the stage, with a detailed plan pinned to an upstream revision.
2. Create a sub-issue of the stage for every task, all before work starts.
   Size each task so one agent can finish it without exhausting its context; split anything larger.
3. Ask the user before starting the stage's work, since they may want to start it in a fresh session.
4. Work tasks in order, one at a time.
   Before each task, work any open `follow-c` issues first.
   A sub-agent implements the task and leaves its changes uncommitted.
5. Accept the task only when `go vet ./...` and `go test ./...` pass, the task's byte-for-byte check against the C engine matches, the diff makes no independent fixes or behavior changes, and the version is bumped.
   Then commit, push, and tick the task in #19.
6. File every bug the task uncovers, in this repository or upstream, following the policy above.
   We match the C engine; upstream bugs are reproduced, not fixed.
   Bugs in this port, including those G3 files here after fixing their side, are fixed as tasks in the current stage.
7. A blocker is an upstream bug that prevents further progress.
   Stop and ask the user.
8. When the stage is done, summarize the open bugs and call out any that must be fixed before the next stage.
   If none must be, continue with the next stage.

The local upstream clone is `~/Software/playbymail/olympia-g3`.

## CI

Do not create GitHub Actions workflows or other GitHub CI actions.
Run the checks above locally instead.
