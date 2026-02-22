package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"

	cli "github.com/herdworks/cli-go"
	hypertexting "github.com/herdworks/hypertexting-go"
)

// hyperctl page
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

// hyperctl page ls
func PageListCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "ls",
			Description: "List website pages.",
			Usage:       "hyperctl page ls [options] [<filter>]",
			Validator:   PageListCommandValidator,
			Handler:     PageListCommandHandler,
		},
		Options: cli.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
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
	var websiteURI *url.URL
	if websiteURI, err = url.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if err = website.Load(false, nil); err != nil {
		return err
	}
	if len(cmd.Args) > 0 {
		filter = cmd.Args[0]
	}
	if pages, err = (website.ListPages(filter)); err != nil {
		return err
	}
	if err = hypertexting.LoadPages(context.TODO(), website, pages, nil, false); err != nil {
		return err
	}
	pages.SortBy("slug", true)
	PageListCommandOutput(cmd, website, pages)
	return nil
}

func PageListCommandOutput(cmd *cli.Command, website *hypertexting.Website, pages hypertexting.Pages) {
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
			draftcol, hypertexting.PageIsDraft(website, pages[i]),
			pathcol, pages[i].Path,
		)
	}
}

// hyperctl page inspect
func PageInspectCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "inspect",
			Description: "Inspect website pages.",
			Usage:       "hyperctl page inspect [options] <path>",
			Validator:   PageInspectCommandValidator,
			Handler:     PageInspectCommandHandler,
		},
		Options: cli.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
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
	var filter string = "all"
	var website *hypertexting.Website
	var pages hypertexting.Pages
	var page *hypertexting.Page
	var pageBytes []byte
	var websiteURI *url.URL
	if websiteURI, err = url.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if err = website.Load(false, nil); err != nil {
		return err
	}
	if len(cmd.Args) > 0 {
		filter = cmd.Args[0]
	}
	if pages, err = (website.ListPages(filter)); err != nil {
		return err
	}
	if err = hypertexting.LoadPages(context.TODO(), website, pages, nil, false); err != nil {
		return err
	}
	if err = hypertexting.LoadPageFeeds(context.TODO(), website, pages); err != nil {
		return err
	} // populate page feeds
	pages.ResolveTemplateData(website)
	page = pages.GetPath(cmd.Args[0])
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
