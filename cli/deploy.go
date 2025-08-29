package cli

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/google/uuid"
	hyperproviders "github.com/herdworks/hyperproviders-go/lib"
	hypertexting "github.com/herdworks/hypertexting-go/lib"
)

func DeployCommand() (cmd *Command) {
	return &Command{
		Config: CommandConfig{
			Name:        "deploy",
			Description: "Deploy a static site build to the configured provider.",
			Usage:       "hyperctl deploy [options]",
			Validator:   DeployCommandValidator,
			Handler:     DeployCommandHandler,
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

func DeployCommandValidator(cmd *Command) (err error) {
	return nil
}

func DeployCommandHandler(cmd *Command) (err error) {
	var websiteURI *url.URL
	var website *hypertexting.Website
	var provider *hyperproviders.HyperProvider
	var buildID string = uuid.New().String()
	if websiteURI, err = url.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	build := &hypertexting.HyperBuild{
		URI:     &url.URL{Path: filepath.Join(website.BuildsDir(), buildID)},
		ID:      buildID,
		Website: website,
		Pages:   &website.Pages, // all pages
	}
	configOverrides := &hypertexting.WebsiteConfig{
		ContentDirectory: cmd.Options.Get("content-dir"),
		DataDirectory:    cmd.Options.Get("data-dir"),
		StaticDirectory:  cmd.Options.Get("static-dir"),
		ThemeSettings:    cmd.Options.Get("theme"),
	}
	if err = build.Execute(context.TODO(), nil, configOverrides); err != nil {
		return err
	}
	if err = website.Provider.Validate(); err != nil {
		return err
	}
	if provider, err = hyperproviders.New(&website.Provider); err != nil {
		return err
	}
	if err = provider.Deploy(context.TODO(), build.URI.Path, website.BaseURL); err != nil {
		return err
	}
	fmt.Printf("\nBuild \"%s\" deployed\n", buildID)
	return nil
}
