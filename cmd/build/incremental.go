package build

import (
	"context"
	"net/url"
	"path/filepath"

	"github.com/google/uuid"

	cli "github.com/herdworks/cli-go"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
)

func GenerateCommand() (cmd *cli.Command) {
	id := uuid.New()
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "incremental",
			Description: "Perform an incremental build.",
			Usage:       "hyperctl build incremental [options] [content]",
			Validator:   GenerateCommandValidator,
			Handler:     GenerateCommandHandler,
		},
		Options: cli.CommandOptions{
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

func GenerateCommandValidator(cmd *cli.Command) (err error) {
	return nil
}

func GenerateCommandHandler(cmd *cli.Command) (err error) {
	var website *hypertexting.Website
	var websiteURI *href.Href
	var page hypertexting.Page
	var pageURI *href.Href
	if websiteURI, err = href.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if pageURI, err = href.Parse(cmd.Options.Get("page")); err != nil {
		return err
	}
	page = hypertexting.Page{URI: &pageURI.URL, Path: hypertexting.PagePath(website, &pageURI.URL)}
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
		ThemePath:        cmd.Options.Get("theme"),
	}
	return build.Execute(context.TODO(), nil, configOverrides)
}
