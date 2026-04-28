package deploy

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

func DeployCommand() (cmd *cli.Command) {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "complete",
			Description: "Deploy a complete build.",
			Usage:       "hyperctl deploy complete [options]",
			Validator:   DeployCommandValidator,
			Handler:     DeployCommandHandler,
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
				Name:        "environment",
				ShortName:   "e",
				EnvName:     "HYPER_ENVIRONMENT",
				Description: "Deployment environment to use from site.environments.",
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

func DeployCommandValidator(cmd *cli.Command) (err error) {
	return nil
}

func DeployCommandHandler(cmd *cli.Command) (err error) {
	var websiteURI *href.Href
	var website *hypertexting.Website
	var buildID string = uuid.New().String()
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
		URI:     &url.URL{Path: filepath.Join(website.BuildsDirPath(), buildID)},
		ID:      buildID,
		Website: website,
		Pages:   &website.Pages, // all pages
	}
	if err = build.Execute(context.TODO(), nil, nil); err != nil {
		return err
	}
	environment, err := deploymentEnvironment(website, cmd.Options.Get("environment"))
	if err != nil {
		return err
	}
	if err = environment.Deploy(context.TODO(), build.URI.Path, deploymentSession()); err != nil {
		return err
	}
	fmt.Printf("\nBuild \"%s\" deployed\n", buildID)
	return nil
}
