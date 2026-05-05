# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to [Semantic Versioning].

## 0.20.0 (2026-05-05)

Updated `hyperctl` for the current HyperTexting, HyperTemplates, and HyperProviders release set.
This release adds the `ht-apply` template value/plugin pipeline, support for multiple hosting providers, named deployment environments, git-based hosting providers, and support for computed template data (!).

> BREAKING CHANGES: hyperctl v0.20.0 drops support for the drafts_dir – drafts are now defined by setting `draft:true` in page data files. > This release also changes template data namespace names which are now derived from their file path.

* Added support for HyperTemplates `ht-apply` template values in generated sites, including `${ ... }` substitutions in element attributes and text nodes, `--ht-value(...)` CSS values, `get` lookups, string-literal fallbacks, and JavaScript plugin functions in value pipelines.
* Added theme and website JavaScript plugin discovery through the new `plugins_dir` settings in `theme.json` and `site.yaml`.
* Added named deployment environment support. `hyperctl deploy complete` and `hyperctl deploy incremental` can deploy through the selected `site.environments` entry.
* Added `--providers` / `HYPER_PROVIDERS` to deploy commands for comma-separated provider selection from `site.providers`.
* Added multi-provider deployment support through `hyperproviders-go` environments and env-backed secrets.
* Added git-based hosting provider support through the new `git` provider kind.
* Added `page.ugly_url` so pages can opt out of `<path>/index.html` index pages, and instead generate `<path>.html` pages.
* Added computed-data scripts. Drop `*.js` files into the data directory alongside YAML/JSON/OPML; each must define `function namespace()` and return the value to install under the file's keypath. Computed scripts see `ht`, `build`, `env`, `site`, `data`, and `theme` as ambient bindings — `site.pages` is populated because the pass runs after page loading.
* Added `--verbose` / `-v` to `hyperctl cms data ls` and `hyperctl cms data inspect`. Computed-script `console.log` output is silenced by default so stdout stays pipeable to `jq`; `--verbose` routes it to stderr.
* Changed deploy target selection so `--environment` and `--providers` are mutually exclusive, with no selector defaulting to the `default` environment.
* Changed build, deploy, dev server, CMS, and theme commands to use the new root-scoped HyperTexting path helpers.
* Changed `hyperctl cms page` draft filtering to use the page-level `draft` property instead of a separate drafts directory.
* Changed `hyperctl dev server` to render draft pages while keeping build and deploy commands draft-excluding by default.
* Changed `hyperctl cms asset ls` to use resolved asset entries from `hypertexting-go`, including theme, website, and page assets with their source, size, and final public path.
* Changed `hyperctl cms data ls` to enumerate the full build-context namespace set (`ht`, `build`, `env`, `site`, `data`, `theme`) instead of only website global data. Backed by the new `(*HyperBuild).LoadContext` / `Context()` API in `hypertexting-go` v0.26.0.
* Changed `hyperctl cms data inspect <namespace>` to `hyperctl cms data inspect <keypath>`. Accepts arbitrary dotted keypaths (e.g. `site.pages`, `data.social.twitter`) instead of only top-level namespaces.
* Changed data namespaces to derive from filenames (and subdirectory layout), not from a `namespace:` key or single-top-key inference inside the file. Migration: rename data files so `<filename>.<ext>` matches the desired top-level namespace; group related files into subdirectories for nested keypaths (`data/social/twitter.yaml` → `data.social.twitter.*`). Cross-format collisions (`foo.yaml` + `foo.json`) and directory-vs-file collisions (`foo.yaml` + `foo/`) now error.
* Changed plugin file format: each plugin file must define a top-level `function plugin(...)` (or bind the `plugin` identifier to a callable, e.g. `const plugin = (input) => ...`). The plugin's public name is its filename (sans `.js`); names must match `^[A-Za-z_][A-Za-z0-9_]*$` and cannot be `get`, `plugin`, `namespace`, or `data`. The previous `const plugin = name` selector inside the file is no longer recognized.
* Removed `drafts_dir` configuration and loading of pages from a drafts directory.
* Fixed build, deploy, and dev server commands to honor `site.config.builds_dir` when constructing build output paths.
* Updated `github.com/herdworks/hypertemplates-go` to v0.26.0.
* Updated `github.com/herdworks/hypertexting-go` to v0.26.0.
* Updated `github.com/herdworks/hyperproviders-go` to v0.19.0.
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
