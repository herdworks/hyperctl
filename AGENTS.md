# AGENTS.md

## Project-specific Instructions

<!-- Project-specific instructions take priority over general instructions -->

* Build: `sh scripts/build.sh`
* Test: `go test ./... -v`
* Always run `go fmt ./...` before committing
* This is a CLI tool and static site generator; it is a sibling project to the HyperTexting app (NOT an app dependency)
* Both hyperctl and the HyperTexting app consume the same Go libraries, demonstrating that the lower-level libraries are designed for use in multiple contexts (iOS app via libhypertexting, and CLI)
* This project depends on:
  * `../hypergraph-go` — Web scraper (OpenGraph/OEmbed/Twitter Cards)
  * `../hypertexting-go` — Core CMS data model and logic
  * `../hypertemplates-go` — HTML templating and static site generation
  * `../hyperproviders-go` — Multi-provider deployment
  * `../hyperfeed-go` — Feed reader (transitive dependency)
* To test uncommitted changes across repos, uncomment the `replace` directives in `go.mod` to point at local directories;
  re-comment them before committing
* Only bump the tool version from the `main` branch. On other branches, record breaking changes (and other notable changes) under the `Unreleased` section of `CHANGELOG.md`; the version bump happens at release time on `main`

### Release Instructions

* Release commits and tags must only be cut from the `main` branch
* Before creating a release commit, verify that `main` is cleanly formatted, passes tests, and builds by running `go fmt ./...`, `go test ./... -v`, and `go build ./...`
* Move the relevant `CHANGELOG.md` entries from `Unreleased` into `## vX.Y.Z (YYYY-MM-DD)`
* Update `version.go` with `LibVersion = "X.Y.Z"`
* Set `LibReleaseDate` to the current local time rounded down to the most recent hour; for example, `2026-05-11T10:36:31-07:00` becomes `2026-05-11T10:00:00-07:00`
* After the version/changelog edits, run `go fmt ./...`, `go test ./... -v`, and `go build ./...` again
* Commit as `updated changelog and bumped version for X.Y.Z`
* Create an annotated tag `vX.Y.Z` with message `release version X.Y.Z`
* Tags should not skip version numbers (e.g. don't jump from v0.1.0 to v0.3.0)

<!-- See https://agents.md for more information -->
