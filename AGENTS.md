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

<!-- See https://agents.md for more information -->
