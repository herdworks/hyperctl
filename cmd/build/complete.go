package build

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/google/uuid"
	cli "github.com/herdworks/cli-go"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
)

func BuildCommand() *cli.Command {
	id := uuid.New()
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "build",
			Description: "Build management tools",
			Usage:       "hyperctl build [options]",
			Validator:   BuildCommandValidator,
			Handler:     BuildCommandHandler,
		},
		Options: cli.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Default:     "site.yaml",
				Required:    true,
				Secret:      false,
			},
			{
				Name:        "id",
				ShortName:   "i",
				Description: "Identifier to use for the build.",
				Required:    false,
				Secret:      false,
				Default:     id.String(),
			},
			{
				Name:        "incremental",
				Description: "Switch from full-site build mode to incremental build mode. Requires --page.",
				Boolean:     true,
			},
			{
				Name:        "page",
				ShortName:   "p",
				Description: "Path to the page index file for incremental builds.",
			},
			{
				Name:        "content-dir",
				EnvName:     "HYPER_CONTENT_DIR",
				Description: "Relative path to the content directory to use for the build (overwrites site.config.content_dir).",
			},
			{
				Name:        "data-dir",
				EnvName:     "HYPER_DATA_DIR",
				Description: "Relative path to the data directory to use for the build (overwrites site.config.data_dir).",
			},
			{
				Name:        "static-dir",
				EnvName:     "HYPER_STATIC_DIR",
				Description: "Relative path to the static directory to use for the build (overwrites site.config.static_dir).",
			},
			{
				Name:        "theme",
				EnvName:     "HYPER_THEME",
				Description: "Relative path to theme config to use for the build (overwrites site.config.theme).",
			},
			{
				Name:        "data-file",
				Description: "Build-root-relative or absolute path for build data JSON, for example build.json.",
			},
			{
				Name:        "log-file",
				Description: "Build-root-relative or absolute path for build logs, for example build.log.",
			},
		},
	}
}

func BuildCommandValidator(cmd *cli.Command) (err error) {
	incremental := cmd.Options.Get("incremental") == "true"
	page := cmd.Options.Get("page")
	if incremental && page == "" {
		return fmt.Errorf("--page is required when --incremental is set")
	}
	if !incremental && page != "" {
		return fmt.Errorf("--page can only be used when --incremental is set")
	}
	return nil
}

func BuildCommandHandler(cmd *cli.Command) (err error) {
	var websiteURI *href.Href
	var website *hypertexting.Website
	incremental := cmd.Options.Get("incremental") == "true"
	config := hypertexting.BuildConfig{
		WebsiteConfig: hypertexting.WebsiteConfig{
			ContentDirectory: cmd.Options.Get("content-dir"),
			DataDirectory:    cmd.Options.Get("data-dir"),
			StaticDirectory:  cmd.Options.Get("static-dir"),
			ThemePath:        cmd.Options.Get("theme"),
		},
		DataFile: cmd.Options.Get("data-file"),
		LogFile:  cmd.Options.Get("log-file"),
		Verbose:  true,
	}
	if websiteURI, err = href.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if err = website.Load(false, &config.WebsiteConfig, nil); err != nil {
		return err
	}
	var pages *hypertexting.Pages = &website.Pages
	if incremental {
		pageURI, err := href.Parse(cmd.Options.Get("page"))
		if err != nil {
			return err
		}
		page := &hypertexting.Page{URI: &pageURI.URL, Path: hypertexting.PagePath(website, &pageURI.URL)}
		pages = &hypertexting.Pages{page}
	}
	build := &hypertexting.Build{
		URI:         &url.URL{Path: filepath.Join(website.BuildsDirPath(), cmd.Options.Get("id"))},
		ID:          cmd.Options.Get("id"),
		Website:     website,
		Config:      config,
		Incremental: incremental,
		Pages:       pages,
	}
	return build.Execute(context.TODO())
}
