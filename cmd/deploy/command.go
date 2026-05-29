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

func Command() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "deploy",
			Description: "Deployment management tools",
			Usage:       "hyperctl deploy [options]",
			Validator:   DeployCommandValidator,
			Handler:     DeployCommandHandler,
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
				Name:        "environment",
				ShortName:   "e",
				EnvName:     "HYPER_ENVIRONMENT",
				Description: "Deployment environment to use from site.environments.",
				Secret:      false,
			},
			{
				Name:        "providers",
				EnvName:     "HYPER_PROVIDERS",
				Description: "Comma-separated provider names to deploy from site.providers.",
				Secret:      false,
			},
			{
				Name:        "incremental",
				Description: "Switch from full-site deployment mode to incremental deployment mode. Requires --page.",
				Boolean:     true,
			},
			{
				Name:        "page",
				ShortName:   "p",
				Description: "Path to the page index file for incremental deployments.",
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

func DeployCommandValidator(cmd *cli.Command) (err error) {
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

func DeployCommandHandler(cmd *cli.Command) (err error) {
	var websiteURI *href.Href
	var website *hypertexting.Website
	var buildID string = uuid.New().String()
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
		URI:         &url.URL{Path: filepath.Join(website.BuildsDirPath(), buildID)},
		ID:          buildID,
		Website:     website,
		Config:      config,
		Incremental: incremental,
		Pages:       pages,
	}
	if err = build.Execute(context.TODO()); err != nil {
		return err
	}
	environment, environmentName, err := deploymentEnvironment(website, cmd.Options.Get("environment"), cmd.Options.Get("providers"))
	if err != nil {
		return err
	}
	if err = deployEnvironment(context.TODO(), website, environment, environmentName, build.URI.Path, deploymentMode(incremental)); err != nil {
		return err
	}
	fmt.Printf("\nBuild \"%s\" deployed\n", buildID)
	return nil
}
