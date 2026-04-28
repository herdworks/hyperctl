# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to [Semantic Versioning].

## Unreleased

## 0.20.0 (2026-04-27)

Added a new `ht-apply` templating directive and all-new variable subtitution system with `${ ... }` style template variables. 

_NOTE: variable substitution was designed to look and function similar to variable substitution in [Javascript template literals], but the HyperTemplates implementation is **not** based on Javascript._

* Added new template variable declaration, with `${ ... }` delimiter. 
  Template variables can be used in element attributes, _and_ text nodes. 

  ```html
  <title ht-apply>${site.title} | ${page.title}</title>
  <meta name='description' content='${page.description}' ht-apply>
  ```

* Added new `ht-apply` template directive that enables variable substitution on template element attributes and text nodes.

  _NOTE: text node variable substitution is scoped to **direct descendents** of a template element. 
  In the following example, the `${page.greeting.planet}` variable will be replaced, but the `${page.greeting.galaxy}` variable will not. The former is a a direct descendent text node of the parent `<a>` element, but the latter is a direct descendant text node of the nested `<span>` element (a grandchild of the parent `<a>` element)._

  ```html
  <a href='${page.path}' ht-apply>
    Hello, ${page.greeting.planet} world! Welcome to the <span>${page.greeting.galaxy}</span> galaxy.
  </a>
  ```

* Added a CSS template variable declaration, with `--ht-value( ... )` delimiter. 
  
  > NOTE: variable subtitution in `<style>` elements and `style=''` attributes must use the CSS syntax. 
  > 
  > Example:
  > 
  > ```html
  > <style ht-apply>
  >     :root article {
  >         color: --ht-value("page.colors.text, page.colors.body", rgba(0,0,0,0.5));
  >     }
  > </style>
  > ```
  >
  > HyperTemplates CSS variables support two arguments: 
  > 
  > * a comma-separated list of template data keys (see `"page.colors.text, page.colors.body"` above)
  > * a fallback value (see `rgba(0,0,0,0.5)` above), which can be any valid CSS value
  >
  > NOTE: The variable substitution syntax for CSS conforms to the [CSS custom function] specification, 
  > which means they are valid CSS. 
  > Syntax highlighters rejoice!

* Templating variables support functions. The only function that is currently supported is: `default`, for providing fallback values.

  > **Example**
  >
  > ```html
  > <meta name='description' content='${page.description | default "Lorem ipsum, placeholder page description."}' ht-apply>
  > ```

* Added support for `site.environments`.
* Added a new `git` provider kind for deploying to git endpoints
* Fixed a bug that caused build and deploy commands to ignore `site.config.builds_dir`
* Updated `github.com/herdworks/hyperproviders-go` to v0.18.1.
* Changed minimum Go toolchain to `1.25.0`.

<!-- Links -->
[Javascript template literals]: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Template_literals
[CSS custom function]: https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Custom_functions_and_mixins/Using_custom_functions#function_basics

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
