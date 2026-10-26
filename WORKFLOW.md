# Project Workflow

Guidance for agents when working in this repository.

## Priorities

Fix bugs before adding features.
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

## C changes and upstream deviations

- When fixing a bug found in the upstream C code or changing its meaning or behavior, label the issue `bug` for a defect fix or `enhancement` for an intentional semantic or behavioral change.
- Reference the issue in the implementing commit with `Fixes #<issue-number>` or `Closes #<issue-number>`.
  Close it when that commit lands on `main` so the issue and commit history document our updates and deviations from upstream.

## CI

Do not create GitHub Actions workflows or other GitHub CI actions.
Run the checks above locally instead.
