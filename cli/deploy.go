package cli

import (
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/herdworks/clic"
	hsync "github.com/herdworks/hsync"
	"github.com/herdworks/htup"
	envsecrets "github.com/herdworks/htup/providers/env"
	filesecrets "github.com/herdworks/htup/providers/files"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

func DeployCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "deploy",
			Description: "Deployment management tools",
			Usage: []string{
				"hyperctl deploy [options]",
			},
		},
		Spec: clic.CommandSpec{
			Key: "deploy", Help: clic.DefaultHelp,
			Options:   deployOptions(),
			Validator: buildValidator,
			Handler:   deployHandler,
		},
	}
}
func deployOptions() clic.Options {
	return clic.Options{
		"config":      {Meta: clic.OptionMeta{Name: "config", ShortName: "c", EnvName: "HYPER_CONFIG", Description: "Path to the website configuration file.", Default: "site.yaml"}, Spec: clic.OptionSpec{Key: "config", Kind: clic.StringOption, Required: true}},
		"id":          {Meta: clic.OptionMeta{Name: "id", ShortName: "i", Description: "Identifier to use for the build.", Default: uuid.New().String()}, Spec: clic.OptionSpec{Key: "id", Kind: clic.StringOption}},
		"environment": {Meta: clic.OptionMeta{Name: "environment", ShortName: "e", EnvName: "HYPER_ENVIRONMENT", Description: "Deployment environment to use from site.environments."}, Spec: clic.OptionSpec{Key: "environment", Kind: clic.StringOption}},
		"providers":   {Meta: clic.OptionMeta{Name: "providers", EnvName: "HYPER_PROVIDERS", Description: "Comma-separated provider names to deploy from site.providers."}, Spec: clic.OptionSpec{Key: "providers", Kind: clic.StringOption}},
		"incremental": {Meta: clic.OptionMeta{Name: "incremental", Description: "Switch from full-site deployment mode to incremental deployment mode. Requires --page.", Default: false}, Spec: clic.OptionSpec{Key: "incremental", Kind: clic.BoolOption}},
		"page":        {Meta: clic.OptionMeta{Name: "page", ShortName: "p", Description: "Path to the page index file for incremental deployments."}, Spec: clic.OptionSpec{Key: "page", Kind: clic.StringOption}},
		"content-dir": {Meta: clic.OptionMeta{Name: "content-dir", EnvName: "HYPER_CONTENT_DIR", Description: "Relative path to the content directory to use for the build (overwrites site.config.content_dir)."}, Spec: clic.OptionSpec{Key: "content-dir", Kind: clic.StringOption}},
		"data-dir":    {Meta: clic.OptionMeta{Name: "data-dir", EnvName: "HYPER_DATA_DIR", Description: "Relative path to the data directory to use for the build (overwrites site.config.data_dir)."}, Spec: clic.OptionSpec{Key: "data-dir", Kind: clic.StringOption}},
		"static-dir":  {Meta: clic.OptionMeta{Name: "static-dir", EnvName: "HYPER_STATIC_DIR", Description: "Relative path to the static directory to use for the build (overwrites site.config.static_dir)."}, Spec: clic.OptionSpec{Key: "static-dir", Kind: clic.StringOption}},
		"theme":       {Meta: clic.OptionMeta{Name: "theme", EnvName: "HYPER_THEME", Description: "Relative path to theme config to use for the build (overwrites site.config.theme)."}, Spec: clic.OptionSpec{Key: "theme", Kind: clic.StringOption}},
		"data-file":   {Meta: clic.OptionMeta{Name: "data-file", Description: "Relative path to write an optional build data JSON file (e.g. build.json)."}, Spec: clic.OptionSpec{Key: "data-file", Kind: clic.StringOption}},
		"log-file":    {Meta: clic.OptionMeta{Name: "log-file", Description: "Relative path to write an optional build log file (e.g. build.log)."}, Spec: clic.OptionSpec{Key: "log-file", Kind: clic.StringOption}},
	}
}
func deployHandler(execution *clic.Execution) (err error) {
	options := execution.Options
	incremental := options["incremental"].BoolValue()
	config := htxt.BuildConfig{WebsiteConfig: htxt.WebsiteConfig{ContentDirectory: options["content-dir"].StringValue(), DataDirectory: options["data-dir"].StringValue(), StaticDirectory: options["static-dir"].StringValue(), ThemePath: options["theme"].StringValue()}, DataFile: options["data-file"].StringValue(), LogFile: options["log-file"].StringValue(), Verbose: true}
	uri, err := href.Parse(options["config"].StringValue())
	if err != nil {
		return err
	}
	website := &htxt.Website{URI: uri}
	if err = website.Load(false, &config.WebsiteConfig, nil); err != nil {
		return err
	}
	var page *htxt.Page
	if incremental {
		pageURI, err := href.Parse(options["page"].StringValue())
		if err != nil {
			return err
		}
		page = &htxt.Page{URI: &pageURI.URL, Path: htxt.PagePath(website, &pageURI.URL)}
	}
	id := options["id"].StringValue()
	if id == "" {
		id = "deploy"
	}
	build := &htxt.Build{URI: &url.URL{Path: filepath.Join(website.BuildsDirPath(), id)}, ID: id, Website: website, Config: config, Incremental: incremental, Page: page}
	if err = build.Execute(execution.Context); err != nil {
		return err
	}
	environment, name, err := deployEnvironment(build.Website, options["environment"].StringValue(), options["providers"].StringValue(), execution)
	if err != nil {
		return err
	}
	var providers htup.HostingProviders
	if name == "" {
		providers, err = environment.Select(build.Website.Providers)
		if err != nil {
			return err
		}
	}
	context := htup.WithDeploymentRoot(execution.Context, build.URI.Path)
	var writer io.Writer
	if name != "" {
		var finish func(error)
		context, writer, finish = htup.BeginDeployment(context, name, nil)
		defer func() { finish(err) }()
	}
	deployment := hsync.Deployment{ConfigPath: build.Website.URI.Path, PayloadPath: build.URI.Path, Providers: providers, Environment: name, Kind: deployKind(incremental)}
	if err = deployment.Execute(context, deploymentSecrets(build.Website), writer); err != nil {
		return err
	}
	_, err = fmt.Fprintf(execution.Stdout, "\nBuild %q deployed\n", id)
	return err
}
func deployEnvironment(website *htxt.Website, name, providerNames string, execution *clic.Execution) (*htup.Environment, string, error) {
	if name != "" && providerNames != "" {
		return nil, "", fmt.Errorf("--providers and --environment are mutually exclusive")
	}
	if providerNames != "" {
		names := []string{}
		for _, part := range strings.Split(providerNames, ",") {
			if part = strings.TrimSpace(part); part != "" {
				names = append(names, part)
			}
		}
		environment, err := website.DeploymentProviders(names)
		return &environment, "", err
	}
	if name == "" {
		name = "default"
		if _, err := fmt.Fprintf(execution.Stdout, "No deployment target specified; using default environment %q\n", name); err != nil {
			return nil, "", err
		}
	}
	environment, err := website.DeploymentEnvironment(name)
	return environment, name, err
}
func deployKind(incremental bool) htup.DeploymentKind {
	if incremental {
		return htup.IncrementalDeployment
	}
	return htup.CompleteDeployment
}
func deploymentSecrets(website *htxt.Website) htup.SecretsProviders {
	files := filesecrets.New()
	if website != nil && website.URI != nil && website.URI.Path != "" {
		files = filesecrets.NewWithBaseDirectory(filepath.Dir(website.URI.Path))
	}
	return htup.SecretsProviders{envsecrets.New(), files}
}
