package cli

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	hypergraph "github.com/herdworks/hypergraph-go/lib"
)

var (
	supportedSchemes = []string{"http", "https"}
	supportedFormats = []string{"json", "JSON", "yaml", "YAML", "yml", "YML"}
)

func GraphCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "graph",
			Description: "Web scraping utilities.",
			Usage:       "hyperctl graph <subcommand> [options] <href>",
			Handler:     GraphCommandHandler,
		},
		SubCommands: Commands{
			GraphDocumentCommand(),
			GraphEntityCommand(),
			GraphAttachmentCommand(),
		},
		Options: CommandOptions{
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

func GraphCommandValidator(cmd *Command) (err error) {
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

func GraphCommandHandler(cmd *Command) (err error) {
	cmd.Help(nil) // display help if no subcommand is provided
	return
}

func GraphDocumentCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "document",
			Description: "Get web page metadata.",
			Usage:       "hyperctl graph document [options] <href>",
			Validator:   GraphCommandValidator,
			Handler:     GraphDocumentCommandHandler,
		},
		Options: CommandOptions{},
	}
}

func GraphDocumentCommandHandler(cmd *Command) (err error) {
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

func GraphEntityCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "entity",
			Description: "Get web author metadata.",
			Usage:       "hyperctl graph entity [options] <href>",
			Validator:   GraphCommandValidator,
			Handler:     GraphEntityCommandHandler,
		},
		Options: CommandOptions{},
	}
}

func GraphEntityCommandHandler(cmd *Command) (err error) {
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

func GraphAttachmentCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "attachment",
			Description: "Get attachment metadata.",
			Usage:       "hyperctl graph attachment [options] <href>",
			Validator:   GraphCommandValidator,
			Handler:     GraphAttachmentCommandHandler,
		},
		Options: CommandOptions{},
	}
}

func GraphAttachmentCommandHandler(cmd *Command) (err error) {
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
