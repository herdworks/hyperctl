package dev

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"

	cli "github.com/herdworks/cli-go"
	hypergraph "github.com/herdworks/hypergraph-go/types"
	"gopkg.in/yaml.v2"
)

var (
	supportedSchemes = []string{"http", "https"}
	supportedFormats = []string{"json", "JSON", "yaml", "YAML", "yml", "YML"}
	supportedNamespaces = []string{"entity", "hypergraph", "opengraph", "twitter", "oembed", "all", "metadata"}
)

func GraphCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "graph",
			Description: "Web scraper utilities.",
			Usage:       "hyperctl dev graph <subcommand> [options] <href>",
			Handler:     GraphCommandHandler,
		},
		SubCommands: cli.Commands{
			GraphAllCommand(),
			GraphEntityCommand(),
			GraphAttachmentCommand(),
		},
		Options: cli.CommandOptions{
			{
				Name:        "format",
				ShortName:   "f",
				Description: "Output format for website metadata.",
				Default:     "json",
				Global:      true,
			},
			{
				Name:        "verbose",
				Description: "Enable verbose mode.",
				Boolean:     true,
				Default:     "false",
				Global:      true,
			},
		},
	}
}

func GraphCommandValidator(cmd *cli.Command) (err error) {
	var uri *url.URL
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expected exactly 1 arguments, got %v (%s)", len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	if uri, err = url.Parse(cmd.Args[0]); err != nil {
		return fmt.Errorf("URL \"%s\" invalid", cmd.Args[1])
	}
	if !slices.Contains(supportedSchemes, uri.Scheme) {
		return fmt.Errorf("scheme \"%s\" invalid", uri.Scheme)
	}
	if !slices.Contains(supportedFormats, cmd.Options.Get("format")) {
		return fmt.Errorf("format \"%s\" invalid", cmd.Options.Get("format"))
	}
	return nil
}

func GraphCommandHandler(cmd *cli.Command) (err error) {
	cmd.Help(nil) // display help if no subcommand is provided
	return
}

func GraphAllCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "all",
			Description: "Get all webpage metadata.",
			Usage:       "hyperctl graph all [options] <href>",
			Validator:   GraphCommandValidator,
			Handler:     GraphMetadataCommandHandler,
		},
		Options: cli.CommandOptions{},
	}
}

func GraphEntityCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "entity",
			Description: "Get website entity metadata.",
			Usage:       "hyperctl graph entity [options] <href>",
			Validator:   GraphCommandValidator,
			Handler:     GraphMetadataCommandHandler,
		},
		Options: cli.CommandOptions{},
	}
}

func GraphAttachmentCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "attachment",
			Description: "Get attachment metadata.",
			Usage:       "hyperctl graph attachment [options] <href>",
			Validator:   GraphCommandValidator,
			Handler:     GraphMetadataCommandHandler,
		},
		Options: cli.CommandOptions{},
	}
}

func GraphMetadataCommandHandler(cmd *cli.Command) (err error) {
	var ctx context.Context = context.Background()
	var verbose bool = cmd.Options.Get("verbose") == "true"
	var jsonBytes []byte
	var yamlBytes []byte
	var metabytes []byte
	href := cmd.Args[0]
	var document *hypergraph.Document = &hypergraph.Document{Href: href}
	if err = document.Download(ctx, verbose, nil, 10, nil); err != nil {
		return errors.Join(fmt.Errorf("downloading document"), err)
	}
	if err = document.Scrape(ctx, verbose, nil, 10, nil); err != nil {
		return errors.Join(fmt.Errorf("scraping document"), err)
	}

	switch cmd.Config.Name {
	case "entity":
		jsonBytes, err = json.MarshalIndent(document.Metadata.Entity, "", "  ")
	case "hypergraph", "attachment":
		jsonBytes, err = json.MarshalIndent(document.Metadata.HyperGraph, "", "  ")
	case "opengraph":
		jsonBytes, err = json.MarshalIndent(document.Metadata.OpenGraph, "", "  ")
	case "twitter":
		jsonBytes, err = json.MarshalIndent(document.Metadata.Twitter, "", "  ")
	case "oembed":
		jsonBytes, err = json.MarshalIndent(document.Metadata.OEmbed, "", "  ")
	case "all", "metadata":
		jsonBytes, err = json.MarshalIndent(document.Metadata, "", "  ")
	default:
		jsonBytes = []byte(fmt.Sprintf("{\"error\": \"namespace \\\"%s\\\" not found (supported namespaces: %s)\"}", cmd.Config.Name, strings.Join(supportedNamespaces, ", ")))
	}
	switch cmd.Options.Get("format") {
	case "yaml", "yml":
		metadata := map[string]any{}
		if err = json.Unmarshal(jsonBytes, &metadata); err != nil {
			return fmt.Errorf("format=yaml error: %w", err)
		}
		if yamlBytes, err = yaml.Marshal(metadata); err != nil {
			return fmt.Errorf("format=yaml error: %w", err)
		}
		metabytes = fmt.Appendf(metabytes, "---\n%s", yamlBytes)
	case "json":
		metabytes = jsonBytes
	}
	if err != nil {
		return fmt.Errorf("error serializing metadata to %s: %w", strings.ToUpper(cmd.Options.Get("format")), err)
	}
	fmt.Printf("%v\n", string(metabytes))
	return nil
}

