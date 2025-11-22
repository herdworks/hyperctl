package cmd

import (
	"context"
	"net/url"
	"path/filepath"

	"github.com/google/uuid"

	cli "github.com/herdworks/cli-go"
	hypertexting "github.com/herdworks/hypertexting-go"
)

func BuildCommand() *cli.Command {
	id := uuid.New()
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "build",
			Description: "Build a static site from HyperTemplates.",
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

func BuildCommandValidator(cmd *cli.Command) (err error) {
	return nil
}

func BuildCommandHandler(cmd *cli.Command) (err error) {
	var websiteURI *url.URL
	var website *hypertexting.Website
	if websiteURI, err = url.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	configOverrides := &hypertexting.WebsiteConfig{
		ContentDirectory: cmd.Options.Get("content-dir"),
		DataDirectory:    cmd.Options.Get("data-dir"),
		StaticDirectory:  cmd.Options.Get("static-dir"),
		ThemeSettings:    cmd.Options.Get("theme"),
	}
	build := &hypertexting.HyperBuild{
		URI:     &url.URL{Path: filepath.Join(website.BuildsDir(), cmd.Options.Get("id"))},
		ID:      cmd.Options.Get("id"),
		Website: website,
		Pages:   &website.Pages, // all pages
	}
	return build.Execute(context.TODO(), nil, configOverrides)
}
