package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/herdworks/clic"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

func PageCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "page",
			Description: "Manage pages and drafts.",
			Usage: []string{
				"hyperctl cms page [command] [options]",
			},
		},
		Spec: clic.CommandSpec{
			Key:  "page",
			Help: clic.DefaultHelp,
		},
		Commands: clic.MustCommandsFrom(
			PageListCommand,
			PageInspectCommand,
		),
	}
}
func PageListCommand() clic.Command {
	var config clic.Option = GlobalWebsiteConfigOption()
	var visibility clic.Option = clic.Option{
		Meta: clic.OptionMeta{Name: "visibility", Description: "Page visibility: all, public, unlisted, status, draft, or none.", Default: "all"},
		Spec: clic.OptionSpec{Key: "visibility", Kind: clic.StringOption},
	}
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "ls",
			Description: "List website pages.",
			Usage: []string{
				"hyperctl cms page ls [options]",
			}},
		Spec: clic.CommandSpec{
			Key:  "ls",
			Help: clic.DefaultHelp,
			Options: clic.Options{
				config.Spec.Key:     &config,
				visibility.Spec.Key: &visibility,
			},
			Validator: pageListValidator,
			Handler:   pageListHandler,
		},
	}
}
func PageInspectCommand() clic.Command {
	var config clic.Option = GlobalWebsiteConfigOption()
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "inspect",
			Description: "Inspect website pages.",
			Usage: []string{
				"hyperctl cms page inspect [options] <path>",
			},
		},
		Spec: clic.CommandSpec{
			Key:  "inspect",
			Help: clic.DefaultHelp,
			Options: clic.Options{
				config.Spec.Key: &config,
			},
			Validator: expectedPositionalArguments(1),
			Handler:   pageInspectHandler,
		},
	}
}
func pageListValidator(execution *clic.Execution) error {
	if err := expectedPositionalArguments(0)(execution); err != nil {
		return err
	}
	visibility := execution.Options["visibility"].StringValue()
	switch visibility {
	case "all", "public", "unlisted", "status", "draft", "none":
		return nil
	default:
		return fmt.Errorf("visibility %q unknown (expected one of: all, public, unlisted, status, draft, none)", visibility)
	}
}
func pageListHandler(execution *clic.Execution) error {
	uri, err := href.Parse(execution.Options["config"].StringValue())
	if err != nil {
		return err
	}
	website := &htxt.Website{URI: uri}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	pages, err := website.ListPages()
	if err != nil {
		return err
	}
	if err = htxt.LoadPages(execution.Context, website, pages, nil, false, nil); err != nil {
		return err
	}
	pages = filterPages(pages, execution.Options["visibility"].StringValue())
	pages.SortBy("slug", true)
	return clic.WriteList(execution.Stdout, clic.ColumnarFormat, pages,
		clic.Field{Key: "lastmod", Header: "Last Modified"},
		clic.Field{Key: "title", Header: "Title"},
		clic.Field{Key: "visibility", Header: "Visibility"},
		clic.Field{Key: "path", Header: "Path"},
	)
}
func filterPages(pages htxt.Pages, visibility string) htxt.Pages {
	result := make(htxt.Pages, 0, len(pages))
	for _, page := range pages {
		if page == nil || !matchesPageVisibility(page.Visibility, visibility) {
			continue
		}
		result = append(result, page)
	}
	return result
}

func matchesPageVisibility(visibility htxt.PageVisibilityKind, filter string) bool {
	switch filter {
	case "all":
		return true
	case "public":
		return visibility.IsPublic()
	case "unlisted":
		return visibility.IsUnlisted()
	case "status":
		return visibility.IsStatus()
	case "draft":
		return visibility.IsDraft()
	case "none":
		return visibility.IsNone()
	default:
		return false
	}
}
func pageInspectHandler(execution *clic.Execution) error {
	uri, err := href.Parse(execution.Options["config"].StringValue())
	if err != nil {
		return err
	}
	website := &htxt.Website{URI: uri}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	build := &htxt.Build{Website: website}
	if err = build.LoadContext(context.Background(), nil); err != nil {
		return err
	}
	page := build.Pages.PageByPath(execution.Command.State.Arguments[0])
	if page == nil {
		return fmt.Errorf("page %q not found", execution.Command.State.Arguments[0])
	}
	data, err := json.Marshal(page.TemplateData)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(execution.Stdout, string(data))
	return err
}
