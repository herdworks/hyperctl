# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to [Semantic Versioning].

## Unreleased

* **Added:** `hyperctl agent skill up [skills-directory]` synchronizes bundled
  HyperTexting skills into every discovered project-local skill root and
  rejects global Claude and Codex skill directories.
* **Changed:** `hyperctl agent skill` now remounts the reusable `htskills/cli`
  command tree, including `readme`, `ls`, `describe`, `install`, `update`, and
  `up`; `hyperctl agent readme` remains a compatibility link to `skill readme`.
* **Changed:** Narrowed `hyperctl agent` skill tests to command mounting and
  Hyperctl namespace integration; `htskills` owns bundled-skill content and
  command-behavior tests.
* **Removed:** The unregistered `hyperctl dev mimetype` helper.
* **Added:** `hyperctl agent` onboarding and discovery commands:
  `readme`; `schema ls` and `schema describe`; and `skill ls`, `describe`, `install`, and `update` for bundled portable Agent Skills.
* **Added:** `hyperctl sync` commands for hsync checkout workflows:
  `download`, `init`, `pull`, `push`, `remote`, `reset`, and `status`.
  Configuration-consuming commands resolve `site.yaml` through HyperTexting's
  website-aware sync manager, defaulting omitted include, ignore, and maximum
  file-size policy while honoring `HYPER_CONFIG`.
* **Added:** `hyperctl download`, `pull`, and `push` shortcuts for their
  corresponding `hyperctl sync` commands.
* **Added:** `hyperctl version` prints the CLI version;
  use `hyperctl version --verbose` to list the individual HyperTexting package versions.
* **Added:** `hyperctl server` is a shortcut for `hyperctl website server`.
* **Added:** `hyperctl htxt` is an alias for `hyperctl website`.
* **Added:** `hyperctl httt` is an alias for `hyperctl theme`.
* **Added:** `hyperctl tmpl` is an alias for `hyperctl theme`.
* **Changed:** `hyperctl cms` is now `hyperctl website`.
* **Changed:** `hyperctl dev server` is now `hyperctl website server`.
* **Changed:** `hyperctl dev render` is now `hyperctl theme render`.
* **Changed:** `hyperctl theme` now remixes the standalone `httt` theme workflow.
* **Changed:** `hyperctl build` is a shortcut for `hyperctl website build`.
* **Changed:** `hyperctl cms` is now an alias for `hyperctl website`.
* **Changed:** `hyperctl dev` is now an alias for `hyperctl website`.
* **Changed:** `hyperctl website build` uses HyperTexting's private website
  page catalog and selected build output sets.
* **Changed:** `hyperctl deploy` now uses HyperTexting's private website page
  catalog and selected build output sets, and its hsync-enabled deployment
  lifecycle. Sites with `site.sync` prepare synchronization artifacts before
  provider deployment and finalize their local synchronization state only after
  every provider succeeds.
* **Changed:** `hyperctl website deploy` now delegates its build and hsync
  deployment lifecycle to `htxt.Deployment`; `hyperctl deploy` remains a
  hidden compatibility link to that command.
* **Changed:** `hyperctl website server` uses HyperTexting's private website page catalog
  and selected build output sets.
* **Changed:** `hyperctl website` inspection commands use HyperTexting's
  private website page catalog and selected build output sets.
* **Changed:** `hyperctl agent schema` replaces the retired `feed` schema with
  Feed 2.0's `feed-config`, `feed-data`, and `feed-selector` schemas.
* **Changed:** `hyperctl theme ls` supports `--format columns|json` output.
* **Changed:** `hyperctl website page ls` replaces positional filters with
  `--visibility`.
* **Changed:** moved the CLI command tree to `hyperctl/cli`; the executable remains at `cmd/hyperctl`.
* **Updated:** package dependencies, including:
  `cli-go` v0.1.0 to `clic` v0.4.3
  `hyperfeed-go` v0.26.5 to `hyperfeed` v0.28.0
  `hyperproviders-go` v0.23.0 to `htup` v0.27.3
  `hypertemplates-go` v0.35.0 to `httt` v0.39.1
  `hypertexting-go` v0.38.0 to `htxt` v0.45.4
  `hypergraph-go` v0.22.3 to `hq` v0.24.0
  `hsync` to v0.7.1 and `htskills` to v0.8.0;
  and added _new_ `hypermedia` v0.2.0 and `htxtd` v0.2.0 (indirectly).


## v0.23.0 (2026-07-30)

Updates `hyperctl` to use the next-generation templating directives from the `hypertemplates-go` v0.35.0 and `hypertexting-go` v0.38.0 releases, and adds `hyperctl theme migrate` for migrating themes with legacy directives.

* **BREAKING:** Theme rendering now uses HyperTemplates' typed `${ ... }` variable expressions and current directive model. Replace removed legacy directives with their migrated equivalents, and update CSS interpolation from `--ht-value(...)` to `--ht-var(...)`.
* Added `hyperctl theme migrate` with idempotent migration of legacy `ht-not`, `ht-template`, `ht-attr`, `ht-param`, `ht-content`, and `ht-query` declarations. Push-based `ht-pipe` declarations are reported for manual conversion to destination-side `ht-pipe='from "<selector>" as <type>'` directives. It supports explicit theme selection, optional complete-theme output directories, and interactive review; unsupported constructs and optional `ht-debug` cleanup are reported with actionable locations.
* Updated HyperTexting to v0.38.0 and HyperTemplates to v0.35.0 for the current theme, file, and directive APIs.

## 0.22.1 (2026-07-20)

* Updated `github.com/herdworks/hypertemplates-go` to `v0.32.0`, adding `TemplateData` mutation helpers and stricter mutation-path validation for add, replace, and append operations.
* Updated `github.com/herdworks/hypertexting-go` to `v0.36.3`, adding `.csv` discovery for site and theme static data directories, content-type overlay and preview APIs, page-data MIME-type helpers, and feed-reader integration updates.
* Updated `github.com/herdworks/hyperfeed-go` to `v0.26.5` and `github.com/herdworks/hypergraph-go` to `v0.22.3`, adding following/recommendation discovery improvements, following URL metadata persistence, and feed import/discovery fixes.

## 0.22.0 (2026-06-03)

Added support for build logs, and simplified the `hyperctl build` and `hyperctl deploy` commands. 

* Added support for converting `hypertexting.Page` link attachments to Atom entry "related" links
* Added support for converting `hypertexting.Page` image, audio, video, document, and file attachments to Atom entry "enclosure" links
* Added `--data-file` and `--log-file` options to `hyperctl build`, `hyperctl deploy`, and `hyperctl dev server`.
* Combined `hyperctl build complete` and `hyperctl build incremental` into a single `hyperctl build` command, with new `hyperctl build --incremental --page <path>` arguments.
* Combined `hyperctl deploy complete` and `hyperctl deploy incremental` into a single `hyperctl deploy` command, with new `hyperctl deploy --incremental --page <path>` arguments.
* Changed deployments to pass explicit hyperproviders deployment modes: full deploys use `replacement`, while incremental deploys use `overlay`.
* Updated command call sites for the new `hypertexting-go` `Build` / `BuildConfig` APIs and explicit writer logging APIs.
* Updated `github.com/herdworks/hypertexting-go` to `v0.34.0` with renamed `hypertexting.HyperBuild` → `hypertexting.Build` and new/simplified `hypertexting.BuildConfig` options that replace the growing number of `build.Execute()` arguments. 
* Updated `github.com/herdworks/hyperproviders-go` to `v0.23.0` with support for new `.deploymentignore` files
* Updated `github.com/herdworks/hypertemplates-go` to `v0.30.0` with io.Writer support for improved logging.
* Updated `github.com/herdworks/hyperfeed-go` to `v0.20.0` with improved Atom, RSS, and JSON Feed parser handling for related/external links versus enclosure attachments, deduplicates item attachments by normalized href, and enclosure link classification fixes.
* Updated `github.com/herdworks/hypergraph-go` to `v0.19.1` with updated feed discovery that skips `.well-known` URIs.

## 0.21.0 (2026-05-27)

Added support for a "files" secret provider and removed the git provider `ssh_key_path` secret.

* Added support for file-based secrets. 
  The files secret provider looks up secrets using the secret key as a file path, with support for absolute paths, and `./`, `../`, `~/`, `$HOME/`, or `${HOME}/` resolution.
* Removed support for the `git` provider `ssh_key_path` secret.
  To use an SSH private key with the git provider, configure the `ssh_private_key` secret with a file path key; e.g.:

  ```yaml
  providers:
    github:
      kind: git
      endpoint: https://github.com/example/website.git
      config:
        branch: gh-pages
      secrets:
        - name: ssh_private_key
          key: $HOME/.ssh/example
  ```

  _NOTE: relative path secret keys (e.g. `key: deploy_key.pem`) are resolved from the website root directory.

## 0.20.1 (2026-05-27)

Updated `github.com/herdworks/hypertexting-go` to 0.32.0 for generated atom feed canonical url fixes.

## 0.20.0 (2026-05-12)

Updated `hyperctl` for the current HyperTexting, HyperTemplates, and HyperProviders release set.
This release adds support for a new `ht-apply` directive, template variables, templating plugins (!), computed template data, git-based hosting providers, multi-provider hosting environments, and more.

> BREAKING CHANGES: hyperctl v0.20.0 drops support for the drafts_dir – drafts are now defined by setting `draft:true` in page data files. 
> This release also changes template data namespace names which are now derived from their file path.

* Added a new `ht-apply` templating directive which performs variable substitution. 
* Added template variables, formatted `${ ... }` strings that can be used to insert template data into element attributes and text nodes.
* Added support for `<style>` element and `style` attribute variable substitution using the `--ht-value(...)` variable syntax.
* Added variable substitution _functions_ including a built-in `get` function for looking up template data with string-literal fallbacks. 
* Added template variable plugins, user-defined variable substitution functions written in Javascript and authored as ES modules with an `export default` function. 
* Added `site.config.plugins_dir` and `theme.config.plugins_dir` settings to configure template variable plugin discovery.
* Added new "builtins" – automatic template features powered by layout fragments that can be disabled or overriden by adding layout fragments to themes. 
* Added named deployment environments, configured via `site.environments`. Use `hyperctl deploy --environment` for full-site deployments, or `hyperctl deploy --incremental --page <path> --environment` for incremental deployments.
* Added support for multiple named providers, configured via `site.providers`. Use `hyperctl deploy --providers` for full-site deployments, or `hyperctl deploy --incremental --page <path> --providers` for incremental deployments.
* Added git-based hosting provider support through the new `git` provider kind.
* Added `page.ugly_url` so pages can opt out of `<path>/index.html` index pages, and instead generate `<path>.html` pages.
* Added computed-data scripts. Drop `*.js` files into the data directory alongside YAML/JSON/OPML; each must `export default` a function that returns the value to install under the file's keypath. Computed scripts see `ht`, `build`, `env`, `site`, `data`, and `theme` as ambient bindings — `site.pages` is populated because the pass runs after page loading.
* Added `--verbose` / `-v` to `hyperctl cms data ls` and `hyperctl cms data inspect`. Computed-script `console.log` output is silenced by default so stdout stays pipeable to `jq`; `--verbose` routes it to stderr.
* Added support for CSV/TSV template data files. Delimited data must contain a header row, and header rows must have unique non-empty column names. 
* Changed deploy target selection so `--environment` and `--providers` are mutually exclusive, with no selector defaulting to the `default` environment.
* Changed `hyperctl cms asset lookup` to enumerate every theme, website, and page candidate for the queried path in precedence order and identify the winning source. Previously only printed the winning source path.
* Changed `hyperctl cms page` draft filtering to use the page-level `draft` property instead of a separate drafts directory.
* Changed `hyperctl dev server` to render draft pages while keeping build and deploy commands draft-excluding by default.
* Changed `hyperctl cms data ls` to enumerate the full build-context namespace set (`ht`, `build`, `env`, `site`, `data`, `theme`) instead of only website global data. Backed by the shared build context methods in `hypertexting-go` v0.26.0 and newer.
* Changed `hyperctl cms data inspect <namespace>` to `hyperctl cms data inspect <keypath>`. Accepts template data keypaths (e.g. `site.pages`, `data.social.twitter`) instead of only top-level namespaces.
* Changed data namespaces to derive from filenames (and subdirectory layout), not from a `namespace:` key or single-top-key inference inside the file. Migration: rename data files so `<filename>.<ext>` matches the desired top-level namespace; group related files into subdirectories for nested keypaths (`data/social/twitter.yaml` → `data.social.twitter.*`). Cross-format collisions (`foo.yaml` + `foo.json`) and directory-vs-file collisions (`foo.yaml` + `foo/`) now error.
* Changed `-c` / `--config` flag to default to `site.yaml` across all commands that accept it.
* Removed `drafts_dir` configuration and loading of pages from a drafts directory.
* Fixed build, deploy, and dev server commands to honor `site.config.builds_dir` when constructing build output paths.
* Fixed an infinite loop that could occur on websites with multiple page.feed fields configured
* Updated `github.com/herdworks/hypertemplates-go` to v0.29.0. 
* Updated `github.com/herdworks/hypertexting-go` to v0.29.2. 
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
* Changed `hyperctl build` into the full-site build mode for the namespaced `build` command group used by that release.
* Changed `hyperctl generate` → `hyperctl build incremental`.
* Changed `hyperctl deploy` into the full-site deployment mode for the namespaced `deploy` command group used by that release.
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
