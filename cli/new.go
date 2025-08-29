package cli

import (
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"strings"

	hypergraph "github.com/herdworks/hypergraph-go/lib"
	hypertexting "github.com/herdworks/hypertexting-go/lib"
)

func NewCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "new",
			Description: "Create a new website or page.",
			Usage:       "hyperctl new [command] [options]",
			Handler:     NewCommandHandler,
		},
		Options: CommandOptions{},
		SubCommands: Commands{
			NewWebsiteCommand(),
			NewPageCommand(),
		},
	}
}

func NewCommandHandler(cmd *Command) (err error) {
	cmd.Help(nil)
	return nil
}

func NewWebsiteCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "website",
			Description: "Create a new website.",
			Usage:       "hyperctl new website [options] <path>",
			Validator:   NewWebsiteCommandValidator,
			Handler:     NewWebsiteCommandHandler,
		},
		Options: CommandOptions{
			{
				Name:        "baseurl",
				ShortName:   "b",
				Description: "Base URL for the website.",
				Required:    true,
			},
			{
				Name:        "title",
				ShortName:   "t",
				Description: "Website title.",
				Required:    true,
			},
			{
				Name:        "description",
				ShortName:   "d",
				Description: "Website description.",
				Required:    false,
				Default:     "This website is powered by HyperTemplates.",
			},
			{
				Name:        "format",
				Description: "Output format for the website configuration.",
				Default:     "yaml",
			},
		},
	}
}

func NewWebsiteCommandValidator(cmd *Command) (err error) {
	var uri *url.URL
	var baseurl *url.URL
	var supportedSchemes = []string{"https", "http"}
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expected exactly 1 argument, got %v (%s)", len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	if uri, err = url.Parse(cmd.Args[0]); err != nil {
		return err
	}
	if !slices.Contains([]string{"file", ""}, uri.Scheme) {
		return fmt.Errorf("path scheme \"%s\" invalid", uri.Scheme)
	}
	if baseurl, err = url.Parse(cmd.Options.Get("baseurl")); err != nil {
		return fmt.Errorf("baseurl \"%s\" invalid", cmd.Options.Get("baseurl"))
	}
	if !slices.Contains(supportedSchemes, baseurl.Scheme) {
		return fmt.Errorf("baseurl scheme \"%s\" invalid", baseurl.Scheme)
	}
	return nil
}
func NewWebsiteCommandHandler(cmd *Command) (err error) {
	websiteDirectory, _ := url.Parse(cmd.Args[0])
	websiteConfigFileName := "site." + cmd.Options.Get("format")
	websiteConfigURI := &url.URL{Scheme: websiteDirectory.Scheme, Path: filepath.Join(websiteDirectory.Path, websiteConfigFileName)}
	website := hypertexting.Website{
		URI:         websiteConfigURI,
		BaseURL:     cmd.Options.Get("baseurl"),
		Title:       cmd.Options.Get("title"),
		Description: cmd.Options.Get("description"),
		Favicon: hypertexting.WebsiteFavicon{
			Href: "/favicon.png",
		},
	}
	author := &hypergraph.Entity{
		Href:    "/",
		Name:    cmd.Options.Get("title"),
		Favicon: "/favicon.png",
	}
	author.ResolveUsername(cmd.Options.Get("baseurl"))
	if err = author.Validate("author"); err == nil {
		website.Author = author
	}
	website.ResolveConfig(nil)
	if err = website.CreatePackage(); err != nil {
		return err
	}
	return nil
}

func NewPageCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "page",
			Description: "Create a new page.",
			Usage:       "hyperctl new page [options]",
			Handler:     NewPageCommandHandler,
		},
		Options: CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Required:    true,
			},
			{
				Name:        "layout",
				ShortName:   "l",
				Description: "New page layout.",
				Default:     "",
			},
			{
				Name:        "slug",
				ShortName:   "s",
				Description: "New page slug.",
				Default:     "",
			},
			{
				Name:        "title",
				ShortName:   "t",
				Description: "New page title.",
			},
			{
				Name:        "description",
				ShortName:   "d",
				Description: "New page description.",
			},
			{
				Name:        "content-type",
				Description: "New page content type.",
				Required:    false,
			},
			{
				Name:        "content",
				Description: "New page content.",
			},
			{
				Name:        "draft",
				Description: "Create a draft page.",
				Boolean:     true,
				Default:     "false",
			},
		},
	}
}

func NewPageCommandHandler(cmd *Command) (err error) {
	website := hypertexting.Website{URI: &url.URL{Path: cmd.Options.Get("config")}}
	if err = website.Load(false, nil); err != nil {
		return err
	}
	page := hypertexting.Page{
		LayoutRef:   cmd.Options.Get("layout"),
		Slug:        cmd.Options.Get("slug"),
		Title:       cmd.Options.Get("title"),
		Description: cmd.Options.Get("description"),
		Content:     cmd.Options.Get("content"),
		Draft:       cmd.Options.Get("draft") == "true",
	}
	contenttype := website.GetContentTypeByName(cmd.Options.Get("content-type"), true)
	return hypertexting.CreatePage(&website, &page, contenttype)
}
