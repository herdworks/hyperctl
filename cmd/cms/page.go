package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/herdworks/clif"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

// hyperctl cms page
func PageCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "page",
			Description: "Manage pages and drafts.",
			Usage:       "hyperctl cms page [command] [options]",
			Handler:     PageCommandHandler,
		},
		Options: clif.CommandOptions{},
		SubCommands: clif.Commands{
			PageListCommand(),
			PageInspectCommand(),
		},
	}
}
func PageCommandHandler(cmd *clif.Command) (err error) {
	cmd.Help(nil)
	return nil
}

// hyperctl cms page ls
func PageListCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "ls",
			Description: "List website pages.",
			Usage:       "hyperctl cms page ls [options] [<filter>]",
			Validator:   PageListCommandValidator,
			Handler:     PageListCommandHandler,
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
		},
	}
}

func PageListCommandValidator(cmd *clif.Command) (err error) {
	filter := "all"
	filters := []string{"all", "pages", "drafts"}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("expected no more than 1 argument(s), got %v", len(cmd.Args))
	}
	if len(cmd.Args) == 1 {
		filter = cmd.Args[0]
	}
	if !slices.Contains(filters, filter) {
		return fmt.Errorf("filter \"%s\" unknown (expected one of: %s)", cmd.Args[0], strings.Join(filters, ", "))
	}
	return nil
}

func PageListCommandHandler(cmd *clif.Command) (err error) {
	var filter string = "all"
	var pages htxt.Pages
	var website *htxt.Website
	var websiteURI *href.Href
	if websiteURI, err = href.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &htxt.Website{URI: websiteURI}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	if len(cmd.Args) > 0 {
		filter = cmd.Args[0]
	}
	if pages, err = website.ListPages(); err != nil {
		return err
	}
	if err = htxt.LoadPages(context.TODO(), website, pages, nil, false, nil); err != nil {
		return err
	}
	pages = FilterPages(pages, filter)
	pages.SortBy("slug", true)
	PageListCommandOutput(cmd, pages)
	return nil
}

func FilterPages(pages htxt.Pages, filter string) htxt.Pages {
	if filter != "pages" && filter != "drafts" {
		return pages
	}

	filtered := make(htxt.Pages, 0, len(pages))
	for _, page := range pages {
		if page == nil || (filter == "pages" && page.Visibility.IsDraft()) || (filter == "drafts" && !page.Visibility.IsDraft()) {
			continue
		}
		filtered = append(filtered, page)
	}
	return filtered
}

func PageListCommandOutput(cmd *clif.Command, pages htxt.Pages) {
	pathcol := cmd.ColumnWidth(pages, "Path", 1)
	lastmodcol := cmd.ColumnWidth(pages, "LastModified", 1)
	titlecol := cmd.ColumnWidth(pages, "Title", 1)
	draftcol := 6
	fmt.Printf("%-*s %-*s %-*s %-*s\n", lastmodcol, "Last Modified", titlecol, "Title", draftcol, "Draft", pathcol, "Path")
	fmt.Printf("%-*s %-*s %-*s %-*s\n", lastmodcol, "-------------", titlecol, "-----", draftcol, "-----", pathcol, "----")
	for i := range pages {
		fmt.Printf(
			"%-*s %-*v %-*v %-*s\n",
			lastmodcol, pages[i].LastModified().Format(time.RFC3339),
			titlecol, pages[i].Title,
			draftcol, pages[i].Visibility.IsDraft(),
			pathcol, pages[i].Path,
		)
	}
}

// hyperctl cms page inspect
func PageInspectCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "inspect",
			Description: "Inspect website pages.",
			Usage:       "hyperctl cms page inspect [options] <path>",
			Validator:   PageInspectCommandValidator,
			Handler:     PageInspectCommandHandler,
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
		},
	}
}

func PageInspectCommandValidator(cmd *clif.Command) (err error) {
	requiredArgs := 1
	if len(cmd.Args) != requiredArgs {
		return fmt.Errorf("expected %v argument(s), got %v (%s)", requiredArgs, len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}

func PageInspectCommandHandler(cmd *clif.Command) (err error) {
	var website *htxt.Website
	var page *htxt.Page
	var pageBytes []byte
	var websiteURI *href.Href
	if websiteURI, err = href.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &htxt.Website{URI: websiteURI}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	build := &htxt.Build{Website: website}
	if err = build.LoadContext(context.TODO(), nil); err != nil {
		return err
	}
	page = build.Pages.PageByPath(cmd.Args[0])
	if page == nil {
		return fmt.Errorf("page \"%s\" not found", cmd.Args[0])
	}
	pageBytes, err = json.Marshal(page.TemplateData)
	if err != nil {
		return err
	}
	fmt.Println(string(pageBytes))
	return nil
}
