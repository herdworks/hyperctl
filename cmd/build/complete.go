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

func BuildCommand() *cli.Command {
	id := uuid.New()
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "complete",
			Description: "Perform a complete build.",
			Usage:       "hyperctl build complete [options]",
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
	var websiteURI *href.Href
	var website *hypertexting.Website
	configOverrides := &hypertexting.WebsiteConfig{
		ContentDirectory: cmd.Options.Get("content-dir"),
		DataDirectory:    cmd.Options.Get("data-dir"),
		StaticDirectory:  cmd.Options.Get("static-dir"),
		ThemePath:        cmd.Options.Get("theme"),
	}
	if websiteURI, err = href.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if err = website.Load(false, configOverrides); err != nil {
		return err
	}
	build := &hypertexting.HyperBuild{
		URI:     &url.URL{Path: filepath.Join(website.BuildsDirPath(), cmd.Options.Get("id"))},
		ID:      cmd.Options.Get("id"),
		Website: website,
		Pages:   &website.Pages, // all pages
	}
	return build.Execute(context.TODO(), nil, nil)
}
