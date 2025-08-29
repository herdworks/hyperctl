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

func PublishCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "publish",
			Description: "Publish a single page build to a HyperProvider.",
			Usage:       "hyperctl publish [options] <content>",
			Validator:   PublishCommandValidator,
			Handler:     PublishCommandHandler,
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
				Description: "Path to the content index file.",
				Required:    true,
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

func PublishCommandValidator(cmd *Command) (err error) {
	return nil
}

func PublishCommandHandler(cmd *Command) (err error) {
	var website *hypertexting.Website
	var websiteURI *url.URL
	var page hypertexting.Page
	var pageURI *url.URL
	var provider *hyperproviders.HyperProvider
	var buildID string = uuid.New().String()
	if websiteURI, err = url.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if pageURI, err = url.Parse(cmd.Options.Get("page")); err != nil {
		return err
	}
	page = hypertexting.Page{URI: pageURI, Path: hypertexting.PagePath(website, pageURI)}
	build := &hypertexting.HyperBuild{
		URI:         &url.URL{Path: filepath.Join(website.BuildsDir(), buildID)},
		ID:          buildID,
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
