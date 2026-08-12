package build

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/herdworks/clif"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

func Command() *clif.Command {
	id := uuid.New()
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "build",
			Description: "Build management tools",
			Usage:       "hyperctl build [options]",
			Validator:   BuildCommandValidator,
			Handler:     BuildCommandHandler,
		},
		Options: clif.CommandOptions{
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
				Description: "Relative path to write an optional build data JSON file (e.g. build.json).",
			},
			{
				Name:        "log-file",
				Description: "Relative path to write an optional build log file (e.g. build.log).",
			},
		},
	}
}

func BuildCommandValidator(cmd *clif.Command) (err error) {
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

func BuildCommandHandler(cmd *clif.Command) (err error) {
	var websiteURI *href.Href
	var website *htxt.Website
	incremental := cmd.Options.Get("incremental") == "true"
	config := htxt.BuildConfig{
		WebsiteConfig: htxt.WebsiteConfig{
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
	website = &htxt.Website{URI: websiteURI}
	if err = website.Load(false, &config.WebsiteConfig, nil); err != nil {
		return err
	}
	var page *htxt.Page
	if incremental {
		pageURI, err := href.Parse(cmd.Options.Get("page"))
		if err != nil {
			return err
		}
		page = &htxt.Page{URI: &pageURI.URL, Path: htxt.PagePath(website, &pageURI.URL)}
	}
	build := &htxt.Build{
		URI:         &url.URL{Path: filepath.Join(website.BuildsDirPath(), cmd.Options.Get("id"))},
		ID:          cmd.Options.Get("id"),
		Website:     website,
		Config:      config,
		Incremental: incremental,
		Page:        page,
	}
	return build.Execute(context.TODO())
}
