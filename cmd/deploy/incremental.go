package deploy

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/google/uuid"
	cli "github.com/herdworks/cli-go"
	hyperproviders "github.com/herdworks/hyperproviders-go"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
)

func PublishCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "incremental",
			Description: "Deploy an incremental build.",
			Usage:       "hyperctl deploy incremental [options] <content>",
			Validator:   PublishCommandValidator,
			Handler:     PublishCommandHandler,
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

func PublishCommandValidator(cmd *cli.Command) (err error) {
	return nil
}

func PublishCommandHandler(cmd *cli.Command) (err error) {
	var website *hypertexting.Website
	var websiteURI *href.Href
	var page hypertexting.Page
	var pageURI *href.Href
	var provider *hyperproviders.HyperProvider
	var buildID string = uuid.New().String()
	if websiteURI, err = href.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if pageURI, err = href.Parse(cmd.Options.Get("page")); err != nil {
		return err
	}
	page = hypertexting.Page{URI: &pageURI.URL, Path: hypertexting.PagePath(website, &pageURI.URL)}
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
		ThemePath:        cmd.Options.Get("theme"),
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
