# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to [Semantic Versioning].

## Unreleased

* Changed `hyperctl new` → `hyperctl cms new`
* Changed `hyperctl asset` → `hyperctl cms asset`
* Changed `hyperctl content-type` → `hyperctl cms content-type`
* Changed `hyperctl data` → `hyperctl cms data`
* Changed `hyperctl page` → `hyperctl cms page`
* Changed `hyperctl build` → `hyperctl build complete`
* Changed `hyperctl generate` → `hyperctl build incremental`
* Changed `hyperctl deploy` → `hyperctl deploy complete`
* Changed `hyperctl publish` → `hyperctl deploy incremental` 
* Changed `hyperctl graph` → `hyperctl dev graph`
* Changed `hyperctl mimetype` → `hyperctl dev mimetype`
* Changed `hyperctl render` → `hyperctl dev render`
* Changed `hyperctl server` → `hyperctl dev server`
* Fixed several bugs in the live-reloading development server

## 0.17.3

Initial release of `hyperctl` from the standalone `herdworks/hyperctl` repo.

* Moved `cli.go` to standalone `herdworks/cli-go` package
* Changed `go.mod` to use `herdworks/cli-go` package
* Changed `go.mod` to use new versioned `herdworks/hyper*` packages
* Changed/moved CLI command source files to subdirectory of entrypoint (`hyperctl.go`)

## 0.17.2

See `herdworks/hypertexting-monorepo` for version 0.17.2 and prior changes.

<!-- Links -->
[Keep a Changelog]: https://keepachangelog.com/
[Semantic Versioning]: https://semver.org
