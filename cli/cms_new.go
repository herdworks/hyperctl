package cli

import (
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"strings"

	"github.com/herdworks/clic"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

func CMSNewCommand() clic.Command {
	return clic.Command{Meta: clic.CommandMeta{Name: "new", Description: "Create a new website or page.", Usage: []string{
		"hyperctl cms new [command] [options]",
	}}, Spec: clic.CommandSpec{
		Key: "new", Help: clic.DefaultHelp,
	}, Commands: clic.MustCommandsFrom(cmsNewWebsiteCommand, cmsNewPageCommand)}
}
func cmsNewWebsiteCommand() clic.Command {
	return clic.Command{Meta: clic.CommandMeta{Name: "website", Description: "Create a new website.", Usage: []string{
		"hyperctl cms new website [options] <path>",
	}}, Spec: clic.CommandSpec{
		Key: "website", Help: clic.DefaultHelp,
		Options: clic.Options{
			"baseurl":     {Meta: clic.OptionMeta{Name: "baseurl", ShortName: "b", Description: "Base URL for the website."}, Spec: clic.OptionSpec{Key: "baseurl", Kind: clic.StringOption, Required: true}},
			"title":       {Meta: clic.OptionMeta{Name: "title", ShortName: "t", Description: "Website title."}, Spec: clic.OptionSpec{Key: "title", Kind: clic.StringOption, Required: true}},
			"description": {Meta: clic.OptionMeta{Name: "description", ShortName: "d", Description: "Website description.", Default: "This website is powered by HyperTemplates."}, Spec: clic.OptionSpec{Key: "description", Kind: clic.StringOption}},
			"format":      {Meta: clic.OptionMeta{Name: "format", Description: "Output format for the website configuration.", Default: "yaml"}, Spec: clic.OptionSpec{Key: "format", Kind: clic.StringOption}},
		},
		Validator: cmsNewWebsiteValidator,
		Handler:   cmsNewWebsiteHandler,
	}}
}
func cmsNewWebsiteValidator(execution *clic.Execution) error {
	arguments := execution.Command.State.Arguments
	if len(arguments) != 1 {
		return fmt.Errorf("expected exactly 1 argument, got %v (%s)", len(arguments), strings.Join(arguments, ", "))
	}
	uri, err := url.Parse(arguments[0])
	if err != nil {
		return err
	}
	if !slices.Contains([]string{"file", ""}, uri.Scheme) {
		return fmt.Errorf("path scheme %q invalid", uri.Scheme)
	}
	baseURL, err := url.Parse(execution.Options["baseurl"].StringValue())
	if err != nil {
		return fmt.Errorf("baseurl %q invalid", execution.Options["baseurl"].StringValue())
	}
	if !slices.Contains([]string{"https", "http"}, baseURL.Scheme) {
		return fmt.Errorf("baseurl scheme %q invalid", baseURL.Scheme)
	}
	return nil
}
func cmsNewWebsiteHandler(execution *clic.Execution) error {
	directory, _ := url.Parse(execution.Command.State.Arguments[0])
	options := execution.Options
	website := htxt.Website{URI: &href.Href{URL: url.URL{Scheme: directory.Scheme, Path: filepath.Join(directory.Path, "site."+options["format"].StringValue())}}, BaseURL: options["baseurl"].StringValue(), Title: options["title"].StringValue(), Description: options["description"].StringValue(), Favicon: htxt.WebsiteFavicon{Href: "/favicon.png"}, Byline: htxt.Byline{Name: options["title"].StringValue(), Favicon: "/favicon.png", Href: "/"}}
	website.ResolveConfig(nil)
	return website.CreatePackage(nil)
}
func cmsNewPageCommand() clic.Command {
	return clic.Command{Meta: clic.CommandMeta{Name: "page", Description: "Create a new page.", Usage: []string{
		"hyperctl cms new page [options]",
	}}, Spec: clic.CommandSpec{
		Key: "page", Help: clic.DefaultHelp,
		Options: clic.Options{
			"config":       {Meta: clic.OptionMeta{Name: "config", ShortName: "c", EnvName: "HYPER_CONFIG", Description: "Path to the website configuration file.", Default: "site.yaml"}, Spec: clic.OptionSpec{Key: "config", Kind: clic.StringOption, Required: true}},
			"layout":       {Meta: clic.OptionMeta{Name: "layout", ShortName: "l", Description: "New page layout."}, Spec: clic.OptionSpec{Key: "layout", Kind: clic.StringOption}},
			"slug":         {Meta: clic.OptionMeta{Name: "slug", ShortName: "s", Description: "New page slug."}, Spec: clic.OptionSpec{Key: "slug", Kind: clic.StringOption}},
			"title":        {Meta: clic.OptionMeta{Name: "title", ShortName: "t", Description: "New page title."}, Spec: clic.OptionSpec{Key: "title", Kind: clic.StringOption}},
			"description":  {Meta: clic.OptionMeta{Name: "description", ShortName: "d", Description: "New page description."}, Spec: clic.OptionSpec{Key: "description", Kind: clic.StringOption}},
			"content-type": {Meta: clic.OptionMeta{Name: "content-type", Description: "New page content type."}, Spec: clic.OptionSpec{Key: "content-type", Kind: clic.StringOption}},
			"content":      {Meta: clic.OptionMeta{Name: "content", Description: "New page content."}, Spec: clic.OptionSpec{Key: "content", Kind: clic.StringOption}},
			"draft":        {Meta: clic.OptionMeta{Name: "draft", Description: "Create a draft page.", Default: false}, Spec: clic.OptionSpec{Key: "draft", Kind: clic.BoolOption}},
		},
		Handler: cmsNewPageHandler,
	}}
}
func cmsNewPageHandler(execution *clic.Execution) error {
	options := execution.Options
	website := htxt.Website{URI: &href.Href{URL: url.URL{Path: options["config"].StringValue()}}}
	if err := website.Load(false, nil, nil); err != nil {
		return err
	}
	page := htxt.Page{LayoutRef: options["layout"].StringValue(), Slug: options["slug"].StringValue(), Title: options["title"].StringValue(), Description: options["description"].StringValue(), Content: options["content"].StringValue()}
	if options["draft"].BoolValue() {
		page.Visibility = htxt.DraftVisibility
	}
	return htxt.CreatePage(&website, &page, website.GetContentTypeByName(options["content-type"].StringValue(), true, nil), nil)
}
