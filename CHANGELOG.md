# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to [Semantic Versioning].

## Unreleased

* Added `--providers` / `HYPER_PROVIDERS` to deploy commands for comma-separated provider selection from `site.providers`.
* Changed deploy commands so `--environment` and `--providers` are mutually exclusive, with no selector defaulting to the `default` environment.
* Updated `github.com/herdworks/hyperproviders-go` to version 0.19.0
* Updated `github.com/herdworks/hypertexting-go` to version 0.23.0

## 0.20.0 (2026-04-30)

Updated `hyperctl` for the current HyperTexting, HyperTemplates, and HyperProviders release set.
This release adds the `ht-apply` template value/plugin pipeline, named deployment environments with multi-provider deploys, git-based hosting providers, and the redesigned website resource model from `hypertexting-go`.

* Added support for HyperTemplates `ht-apply` template values in generated sites, including `${ ... }` substitutions in element attributes and text nodes, `--ht-value(...)` CSS values, `get` lookups, string-literal fallbacks, and JavaScript plugin functions in value pipelines.
* Added theme and website JavaScript plugin discovery through the new `plugins_dir` settings in `theme.json` and `site.yaml`.
* Added named deployment environment support. `hyperctl deploy complete` and `hyperctl deploy incremental` can deploy through the selected `site.environments` entry.
* Added multi-provider deployment support through `hyperproviders-go` environments and env-backed secrets.
* Added git-based hosting provider support through the new `git` provider kind.
* Changed build, deploy, dev server, CMS, and theme commands to use the new root-scoped HyperTexting path helpers.
* Changed `hyperctl cms page` draft filtering to use the page-level `draft` property instead of a separate drafts directory.
* Changed `hyperctl dev server` to render draft pages while keeping build and deploy commands draft-excluding by default.
* Changed `hyperctl cms asset ls` to use resolved asset entries from `hypertexting-go`, including theme, website, and page assets with their source, size, and final public path.
* Removed `drafts_dir` configuration and loading of pages from a drafts directory.
* Fixed build, deploy, and dev server commands to honor `site.config.builds_dir` when constructing build output paths.
* Updated `github.com/herdworks/hypertemplates-go` to v0.23.0.
* Updated `github.com/herdworks/hypertexting-go` to v0.22.1.
* Updated `github.com/herdworks/hyperproviders-go` to v0.18.2.
* Updated minimum Go toolchain to `1.25.0`.

## 0.19.0 (2026-04-21)

Added `site.pages`, `site.drafts`, and `site.assets` template data. 
Made several improvements to `ht-block` elements, including access to `page.*` template data.

* Added `page.*` template data to `ht-block`
* Added `site.pages`, `site.drafts`, and `site.assets` to template data.
* Updated `github.com/herdworks/hypertemplates-go` to version v0.20.0, including: 
  * added support for template data array accessors (e.g. `page.attachments.1.kind`)
  * added support for`ht-block` user data
  * improved markdown detection of inline vs block-level `ht-block` elements
  * improved markdown detection of `<!-- more -->` comments (now supports whitespace)
* Updated `github.com/herdworks/hypertexting-go` to version v0.20.0, including: 
  * support for `page.*` template data in `ht-block` elements


## 0.18.2 (2026-04-12)

* Added `byline.username` field for backwards compatibility with very early versions of `hyperctl`
* Fixed a race condition in `hyperctl dev server` that would cause the server to crash intermittently
* Fixed a bug that generated broken URLs in Atom feed `<link rel='self'>` elements;
  we're not properly _joining_ "atom.xml" instead of _concatenating_ "atom.xml" to page paths

## 0.18.1 (2026-03-01)

* Added `github.com/herdworks/cmd/version` package for local versioning
* Fixed a bug in `hyperctl theme package` that failed to include the `theme.json` config file in the theme archive.

## 0.18.0 (2026-02-26)

* Added `hyperctl theme` command group with `ls`, `install`, and `package` subcommands for theme management.
* Added theme package installation validation and extraction for zip archives.
* Changed the CLI to use namespaced command groups (`cms`, `theme`, `build`, `deploy`, and `dev`).
* Changed the top-level command order to group the most common workflows (`cms`, `theme`, `build`, `deploy`, and `dev`) together.
* Changed `hyperctl new` → `hyperctl cms new`.
* Changed `hyperctl asset` → `hyperctl cms asset`.
* Changed `hyperctl content-type` → `hyperctl cms content-type`.
* Changed `hyperctl data` → `hyperctl cms data`.
* Changed `hyperctl page` → `hyperctl cms page`.
* Changed `hyperctl build` → `hyperctl build complete`.
* Changed `hyperctl generate` → `hyperctl build incremental`.
* Changed `hyperctl deploy` → `hyperctl deploy complete`.
* Changed `hyperctl publish` → `hyperctl deploy incremental`.
* Changed `hyperctl graph` → `hyperctl dev graph`.
* Changed `hyperctl graph document` → `hyperctl dev graph all`.
* Changed `hyperctl mimetype` → `hyperctl dev mimetype`.
* Changed `hyperctl render` → `hyperctl dev render`.
* Changed `hyperctl server` → `hyperctl dev server`.
* Changed `hyperctl dev graph` subcommands to use a shared metadata scraping handler based on `hypergraph.Document`.
* Changed `hyperctl cms new website` to derive author usernames with `hypergraph.Username()` and skip entity validation.
* Changed template data support via dependency updates to expose root `ht.*`, `env.*`, and `theme.*` namespaces.
* Changed block template data to use the `block.*` namespace via updated `hypertemplates-go` support.
* Changed GitHub Actions release workflows to build on pull requests and only publish tag releases when the tag commit is on `main`.
* Removed `hyperctl data transform`.
* Fixed the live-reloading development server to reload `site.yaml` changes after file updates.
* Fixed `hyperctl dev server` to respect `site.config.builds_dir` when selecting the build output directory.
* Fixed erroneous debug logging in `hyperctl build`.

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
