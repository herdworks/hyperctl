package deploy

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/herdworks/clif"
	hsync "github.com/herdworks/hsync"
	"github.com/herdworks/htup"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

func Command() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "deploy",
			Description: "Deployment management tools",
			Usage:       "hyperctl deploy [options]",
			Validator:   DeployCommandValidator,
			Handler:     DeployCommandHandler,
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

func DeployCommandValidator(cmd *clif.Command) (err error) {
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

func DeployCommandHandler(cmd *clif.Command) (err error) {
	var websiteURI *href.Href
	var website *htxt.Website
	var buildID string = uuid.New().String()
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
		URI:         &url.URL{Path: filepath.Join(website.BuildsDirPath(), buildID)},
		ID:          buildID,
		Website:     website,
		Config:      config,
		Incremental: incremental,
		Page:        page,
	}
	if err = build.Execute(context.TODO()); err != nil {
		return err
	}
	environment, environmentName, err := deploymentEnvironment(build.Website, cmd.Options.Get("environment"), cmd.Options.Get("providers"))
	if err != nil {
		return err
	}
	var providers htup.HostingProviders
	if environmentName == "" {
		providers, err = environment.Select(build.Website.Providers)
		if err != nil {
			return err
		}
	}
	ctx := htup.WithDeploymentRoot(context.TODO(), build.URI.Path)
	var writer io.Writer
	var finish func(error)
	if environmentName != "" {
		ctx, writer, finish = htup.BeginDeployment(ctx, environmentName, nil)
		defer func() { finish(err) }()
	}
	deployment := hsync.Deployment{
		ConfigPath:  build.Website.URI.Path,
		PayloadPath: build.URI.Path,
		Providers:   providers,
		Environment: environmentName,
		Kind:        deploymentKind(incremental),
	}
	if err = deployment.Execute(ctx, DeploymentSecrets(build.Website), writer); err != nil {
		return err
	}
	fmt.Printf("\nBuild \"%s\" deployed\n", buildID)
	return nil
}
