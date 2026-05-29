package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	cli "github.com/herdworks/cli-go"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
)

// hyperctl cms page
func PageCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "page",
			Description: "Manage pages and drafts.",
			Usage:       "hyperctl cms page [command] [options]",
			Handler:     PageCommandHandler,
		},
		Options: cli.CommandOptions{},
		SubCommands: cli.Commands{
			PageListCommand(),
			PageInspectCommand(),
		},
	}
}
func PageCommandHandler(cmd *cli.Command) (err error) {
	cmd.Help(nil)
	return nil
}

// hyperctl cms page ls
func PageListCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "ls",
			Description: "List website pages.",
			Usage:       "hyperctl cms page ls [options] [<filter>]",
			Validator:   PageListCommandValidator,
			Handler:     PageListCommandHandler,
		},
		Options: cli.CommandOptions{
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

func PageListCommandValidator(cmd *cli.Command) (err error) {
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

func PageListCommandHandler(cmd *cli.Command) (err error) {
	var filter string = "all"
	var pages hypertexting.Pages
	var website *hypertexting.Website
	var websiteURI *href.Href
	if websiteURI, err = href.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	if len(cmd.Args) > 0 {
		filter = cmd.Args[0]
	}
	if pages, err = website.ListPages(); err != nil {
		return err
	}
	if err = hypertexting.LoadPages(context.TODO(), website, pages, nil, false, nil); err != nil {
		return err
	}
	pages = FilterPages(pages, filter)
	pages.SortBy("slug", true)
	PageListCommandOutput(cmd, pages)
	return nil
}

func FilterPages(pages hypertexting.Pages, filter string) hypertexting.Pages {
	switch filter {
	case "pages":
		return pages.FilterDraft(false)
	case "drafts":
		return pages.FilterDraft(true)
	default:
		return pages
	}
}

func PageListCommandOutput(cmd *cli.Command, pages hypertexting.Pages) {
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
			draftcol, pages[i].Draft,
			pathcol, pages[i].Path,
		)
	}
}

// hyperctl cms page inspect
func PageInspectCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "inspect",
			Description: "Inspect website pages.",
			Usage:       "hyperctl cms page inspect [options] <path>",
			Validator:   PageInspectCommandValidator,
			Handler:     PageInspectCommandHandler,
		},
		Options: cli.CommandOptions{
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

func PageInspectCommandValidator(cmd *cli.Command) (err error) {
	requiredArgs := 1
	if len(cmd.Args) != requiredArgs {
		return fmt.Errorf("expected %v argument(s), got %v (%s)", requiredArgs, len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}

func PageInspectCommandHandler(cmd *cli.Command) (err error) {
	var website *hypertexting.Website
	var pages hypertexting.Pages
	var page *hypertexting.Page
	var pageBytes []byte
	var websiteURI *href.Href
	if websiteURI, err = href.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	if pages, err = website.ListPages(); err != nil {
		return err
	}
	if err = hypertexting.LoadPages(context.TODO(), website, pages, nil, false, nil); err != nil {
		return err
	}
	if err = hypertexting.LoadPageFeeds(context.TODO(), website, pages, nil); err != nil {
		return err
	} // populate page feeds
	pages.ResolveTemplateData(website)
	page = pages.PageByPath(cmd.Args[0])
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
