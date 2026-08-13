package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/herdworks/clic"
	"github.com/herdworks/hq"
	"gopkg.in/yaml.v2"
)

var supportedSchemes = []string{"http", "https"}
var supportedFormats = []string{"json", "JSON", "yaml", "YAML", "yml", "YML"}
var supportedNamespaces = []string{"entity", "hypergraph", "opengraph", "twitter", "oembed", "all", "metadata"}

func DevGraphCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "graph",
			Description: "Web scraper utilities.",
			Usage: []string{
				"hyperctl dev graph <subcommand> [options] <href>",
			},
		},
		Spec: clic.CommandSpec{
			Key: "graph", Help: clic.DefaultHelp,
			Options: clic.Options{
				"format":  {Meta: clic.OptionMeta{Name: "format", ShortName: "f", Description: "Output format for website metadata.", Default: "json", Global: true}, Spec: clic.OptionSpec{Key: "format", Kind: clic.StringOption}},
				"verbose": {Meta: clic.OptionMeta{Name: "verbose", Description: "Enable verbose mode.", Default: false, Global: true}, Spec: clic.OptionSpec{Key: "verbose", Kind: clic.BoolOption}},
			},
		},
		Commands: clic.MustCommandsFrom(
			devGraphAllCommand,
			devGraphEntityCommand,
			devGraphAttachmentCommand,
		),
	}
}
func devGraphAllCommand() clic.Command { return graphCommand("all", "Get all webpage metadata.") }
func devGraphEntityCommand() clic.Command {
	return graphCommand("entity", "Get website entity metadata.")
}
func devGraphAttachmentCommand() clic.Command {
	return graphCommand("attachment", "Get attachment metadata.")
}
func graphCommand(name, description string) clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        name,
			Description: description,
			Usage: []string{
				"hyperctl dev graph " + name + " [options] <href>",
			},
		},
		Spec: clic.CommandSpec{
			Key: name, Help: clic.DefaultHelp,
			Validator: graphValidator,
			Handler:   graphHandler,
		},
	}
}
func graphValidator(execution *clic.Execution) error {
	arguments := execution.Command.State.Arguments
	if len(arguments) != 1 {
		return fmt.Errorf("expected exactly 1 arguments, got %v (%s)", len(arguments), strings.Join(arguments, ", "))
	}
	uri, err := url.Parse(arguments[0])
	if err != nil {
		return fmt.Errorf("URL %q invalid", arguments[0])
	}
	if !slices.Contains(supportedSchemes, uri.Scheme) {
		return fmt.Errorf("scheme %q invalid", uri.Scheme)
	}
	if !slices.Contains(supportedFormats, execution.Options["format"].StringValue()) {
		return fmt.Errorf("format %q invalid", execution.Options["format"].StringValue())
	}
	return nil
}
func graphHandler(execution *clic.Execution) error {
	document := &hq.Document{Href: execution.Command.State.Arguments[0]}
	if err := document.Download(execution.Context, execution.Options["verbose"].BoolValue(), nil, 10, nil); err != nil {
		return errors.Join(fmt.Errorf("downloading document"), err)
	}
	if err := document.Scrape(execution.Context, execution.Options["verbose"].BoolValue(), nil, 10, nil); err != nil {
		return errors.Join(fmt.Errorf("scraping document"), err)
	}
	var data []byte
	var err error
	switch execution.Command.Meta.Name {
	case "entity":
		data, err = json.MarshalIndent(document.Metadata.Entity, "", "  ")
	case "attachment":
		data, err = json.MarshalIndent(document.Metadata.HyperGraph, "", "  ")
	case "all":
		data, err = json.MarshalIndent(document.Metadata, "", "  ")
	default:
		data = []byte(fmt.Sprintf("{\"error\": \"namespace \\\"%s\\\" not found (supported namespaces: %s)\"}", execution.Command.Meta.Name, strings.Join(supportedNamespaces, ", ")))
	}
	if err != nil {
		return fmt.Errorf("error serializing metadata: %w", err)
	}
	if format := execution.Options["format"].StringValue(); format == "yaml" || format == "yml" {
		metadata := map[string]any{}
		if err = json.Unmarshal(data, &metadata); err != nil {
			return fmt.Errorf("format=yaml error: %w", err)
		}
		data, err = yaml.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("format=yaml error: %w", err)
		}
		data = append([]byte("---\n"), data...)
	}
	_, err = fmt.Fprintln(execution.Stdout, string(data))
	return err
}
