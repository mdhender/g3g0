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

## CI

Do not create GitHub Actions workflows or other GitHub CI actions.
Run the checks above locally instead.
