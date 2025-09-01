package cmd

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	cli "github.com/herdworks/cli-go"
	hypergraph "github.com/herdworks/hypergraph-go"
)

var (
	supportedSchemes = []string{"http", "https"}
	supportedFormats = []string{"json", "JSON", "yaml", "YAML", "yml", "YML"}
)

func GraphCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "graph",
			Description: "Web scraping utilities.",
			Usage:       "hyperctl graph <subcommand> [options] <href>",
			Handler:     GraphCommandHandler,
		},
		SubCommands: cli.Commands{
			GraphDocumentCommand(),
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

func GraphDocumentCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "document",
			Description: "Get web page metadata.",
			Usage:       "hyperctl graph document [options] <href>",
			Validator:   GraphCommandValidator,
			Handler:     GraphDocumentCommandHandler,
		},
		Options: cli.CommandOptions{},
	}
}

func GraphDocumentCommandHandler(cmd *cli.Command) (err error) {
	var metadata map[string]any
	var metabytes []byte
	href := cmd.Args[0]
	metadata = hypergraph.GraphDocument(href, cmd.Options.Get("verbose") == "true")
	switch cmd.Options.Get("format") {
	case "yaml", "yml":
		metabytes, err = hypergraph.PrettyYAML(metadata)
	case "json":
		metabytes, err = hypergraph.PrettyJSON(metadata)
	}
	if err != nil {
		return fmt.Errorf("error serializing metadata to %s: %w", strings.ToUpper(cmd.Options.Get("format")), err)
	}
	fmt.Printf("%v\n", string(metabytes))
	return nil
}

func GraphEntityCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "entity",
			Description: "Get web author metadata.",
			Usage:       "hyperctl graph entity [options] <href>",
			Validator:   GraphCommandValidator,
			Handler:     GraphEntityCommandHandler,
		},
		Options: cli.CommandOptions{},
	}
}

func GraphEntityCommandHandler(cmd *cli.Command) (err error) {
	var metadata map[string]any
	var metabytes []byte
	href := cmd.Args[0]
	metadata = hypergraph.GraphEntity(href, cmd.Options.Get("verbose") == "true")
	switch cmd.Options.Get("format") {
	case "yaml", "yml":
		metabytes, err = hypergraph.PrettyYAML(metadata)
	case "json":
		metabytes, err = hypergraph.PrettyJSON(metadata)
	}
	if err != nil {
		return fmt.Errorf("error serializing metadata to %s: %w", strings.ToUpper(cmd.Options.Get("format")), err)
	}
	fmt.Printf("%v\n", string(metabytes))
	return nil
}

func GraphAttachmentCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "attachment",
			Description: "Get attachment metadata.",
			Usage:       "hyperctl graph attachment [options] <href>",
			Validator:   GraphCommandValidator,
			Handler:     GraphAttachmentCommandHandler,
		},
		Options: cli.CommandOptions{},
	}
}

func GraphAttachmentCommandHandler(cmd *cli.Command) (err error) {
	var metadata map[string]any
	var metabytes []byte
	href := cmd.Args[0]
	metadata = hypergraph.GraphAttachment(href, cmd.Options.Get("verbose") == "true")
	switch cmd.Options.Get("format") {
	case "yaml", "yml":
		metabytes, err = hypergraph.PrettyYAML(metadata)
	case "json":
		metabytes, err = hypergraph.PrettyJSON(metadata)
	}
	if err != nil {
		return fmt.Errorf("error serializing metadata to %s: %w", strings.ToUpper(cmd.Options.Get("format")), err)
	}
	fmt.Printf("%v\n", string(metabytes))
	return nil
}
