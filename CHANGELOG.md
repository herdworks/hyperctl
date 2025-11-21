# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to [Semantic Versioning].

## Unreleased

* Changed `github.com/herdworks/hypergraph-go` dependency to `0.18.0-dev` branch
* Changed `hyperctl new website` commands to use new `hypergraph.Username()` function
* Changed `hyperctl new website` commands to skip entity validation
* Changed `hyperctl graph document` subcommand to "all"
* Changed `hyperctl graph` subcommands to use shared handler function
* Removed `hyperctl data transform` command

## 0.17.3

Initial release of `hyperctl` from the standalone `herdworks/hyperctl` repo.

* Moved `cli.go` to standalone `herdworks/cli-go` package
* Changed `go.mod` to use `herdworks/cli-go` package
* Changed `go.mod` to use new versioned `herdworks/hyper*` packages
* Changed/moved CLI command source files to subdirectory of entrypoint (`hyperctl.go`)

## 0.17.2

See `herdworks/hypertexting-monorepo` for version 0.17.2 and prior changes.
