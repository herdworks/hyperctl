package cms

import (
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"strings"

	"github.com/herdworks/clif"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

func NewCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "new",
			Description: "Create a new website or page.",
			Usage:       "hyperctl cms new [command] [options]",
			Handler:     NewCommandHandler,
		},
		Options: clif.CommandOptions{},
		SubCommands: clif.Commands{
			NewWebsiteCommand(),
			NewPageCommand(),
		},
	}
}

func NewCommandHandler(cmd *clif.Command) (err error) {
	cmd.Help(nil)
	return nil
}

func NewWebsiteCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "website",
			Description: "Create a new website.",
			Usage:       "hyperctl new website [options] <path>",
			Validator:   NewWebsiteCommandValidator,
			Handler:     NewWebsiteCommandHandler,
		},
		Options: clif.CommandOptions{
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

func NewWebsiteCommandValidator(cmd *clif.Command) (err error) {
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
func NewWebsiteCommandHandler(cmd *clif.Command) (err error) {
	websiteDirectory, _ := url.Parse(cmd.Args[0])
	websiteConfigFileName := "site." + cmd.Options.Get("format")
	websiteConfigURI := &href.Href{URL: url.URL{Scheme: websiteDirectory.Scheme, Path: filepath.Join(websiteDirectory.Path, websiteConfigFileName)}}
	website := htxt.Website{
		URI:         websiteConfigURI,
		BaseURL:     cmd.Options.Get("baseurl"),
		Title:       cmd.Options.Get("title"),
		Description: cmd.Options.Get("description"),
		Favicon: htxt.WebsiteFavicon{
			Href: "/favicon.png",
		},
	}
	website.Byline = htxt.Byline{
		Name:    cmd.Options.Get("title"),
		Favicon: "/favicon.png",
		Href:    "/",
	}
	website.ResolveConfig(nil)
	if err = website.CreatePackage(nil); err != nil {
		return err
	}
	return nil
}

func NewPageCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "page",
			Description: "Create a new page.",
			Usage:       "hyperctl new page [options]",
			Handler:     NewPageCommandHandler,
		},
		Options: clif.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Default:     "site.yaml",
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

func NewPageCommandHandler(cmd *clif.Command) (err error) {
	website := htxt.Website{URI: &href.Href{URL: url.URL{Path: cmd.Options.Get("config")}}}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	page := htxt.Page{
		LayoutRef:   cmd.Options.Get("layout"),
		Slug:        cmd.Options.Get("slug"),
		Title:       cmd.Options.Get("title"),
		Description: cmd.Options.Get("description"),
		Content:     cmd.Options.Get("content"),
	}
	if cmd.Options.Get("draft") == "true" {
		page.Visibility = htxt.DraftVisibility
	}
	contenttype := website.GetContentTypeByName(cmd.Options.Get("content-type"), true, nil)
	return htxt.CreatePage(&website, &page, contenttype, nil)
}
