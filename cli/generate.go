package cli

import (
	"context"
	"net/url"
	"path/filepath"

	"github.com/google/uuid"
	hypertexting "github.com/herdworks/hypertexting-go/lib"
)

func GenerateCommand() (cmd *Command) {
	id := uuid.New()
	return &Command{
		Config: CommandConfig{
			Name:        "generate",
			Description: "Perform an incremental build of a single page and its associated contents.",
			Usage:       "hyperctl generate [options] [content]",
			Validator:   GenerateCommandValidator,
			Handler:     GenerateCommandHandler,
		},
		Options: CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Required:    true,
				Secret:      false,
			},
			{
				Name:        "page",
				ShortName:   "p",
				Description: "Path to the page index file.",
				Required:    true,
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
		},
	}
}

func GenerateCommandValidator(cmd *Command) (err error) {
	return nil
}

func GenerateCommandHandler(cmd *Command) (err error) {
	var website *hypertexting.Website
	var websiteURI *url.URL
	var page hypertexting.Page
	var pageURI *url.URL
	if websiteURI, err = url.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if pageURI, err = url.Parse(cmd.Options.Get("page")); err != nil {
		return err
	}
	page = hypertexting.Page{URI: pageURI, Path: hypertexting.PagePath(website, pageURI)}
	build := &hypertexting.HyperBuild{
		URI:         &url.URL{Path: filepath.Join(website.BuildsDir(), cmd.Options.Get("id"))},
		ID:          cmd.Options.Get("id"),
		Website:     website,
		Incremental: true,
		Pages:       &hypertexting.Pages{&page},
	}
	configOverrides := &hypertexting.WebsiteConfig{
		ContentDirectory: cmd.Options.Get("content-dir"),
		DataDirectory:    cmd.Options.Get("data-dir"),
		StaticDirectory:  cmd.Options.Get("static-dir"),
		ThemeSettings:    cmd.Options.Get("theme"),
	}
	return build.Execute(context.TODO(), nil, configOverrides)
}
