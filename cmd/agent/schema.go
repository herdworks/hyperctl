package agent

import (
	"fmt"
	"strings"

	cli "github.com/herdworks/cli-go"
	hsync "github.com/herdworks/hsync"
	hyperproviders "github.com/herdworks/hyperproviders-go"
	hypermark "github.com/herdworks/hypertemplates-go/hypermark"
	markdown "github.com/herdworks/hypertemplates-go/markdown"
	theme "github.com/herdworks/hypertemplates-go/theme"
	hypertexting "github.com/herdworks/hypertexting-go"
	jsonschema "github.com/herdworks/hypertexting-go/jsonschema"
)

var schemas = jsonschema.Registry{
	hypertexting.Attachment{},
	hypertexting.Byline{},
	hypertexting.HyperFeed{},
	hypertexting.HTMLLink{},
	hypertexting.HTMLMeta{},
	markdown.HyperMarkConfig{},
	hypermark.MentionConfig{},
	hypertexting.Page{},
	hyperproviders.HostingProvider{},
	hyperproviders.FilterRules{},
	hyperproviders.Secret{},
	hsync.SyncConfig{},
	hsync.Rules{},
	theme.Theme{},
	hypertexting.Website{},
	hypertexting.WebsiteConfig{},
	hypertexting.WebsiteFavicon{},
}

func schemaCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{Name: "schema", Description: "Describe Hyperctl data and CLI contracts as JSON.", Usage: "hyperctl agent schema [command]", Handler: help}, SubCommands: cli.Commands{schemaListCommand(), schemaDescribeCommand()}}
}

func schemaListCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{Name: "ls", Description: "List stable schema kinds.", Usage: "hyperctl agent schema ls", Handler: func(*cli.Command) error {
		if err := schemas.Validate(); err != nil {
			return fail(err)
		}
		fmt.Print(schemaListOutput(schemaEntries(schemas.Entries())))
		return nil
	}}}
}

type schemaEntries []jsonschema.Schema

func (entries schemaEntries) ColumnValues(field string) (values []string) {
	for _, entry := range entries {
		switch field {
		case "Name":
			values = append(values, entry.Name)
		case "Description":
			values = append(values, entry.Summary)
		}
	}
	return values
}

func schemaListOutput(entries schemaEntries) string {
	namecol := (&cli.Command{}).ColumnWidth(entries, "Name", 1)
	var output strings.Builder
	fmt.Fprintf(&output, "%-*s %s\n", namecol, "Name", "Description")
	fmt.Fprintf(&output, "%-*s %s\n", namecol, "----", "-----------")
	for _, entry := range entries {
		fmt.Fprintf(&output, "%-*s %s\n", namecol, entry.Name, entry.Summary)
	}
	return output.String()
}

func schemaDescribeCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{Name: "describe", Description: "Write one raw JSON Schema document.", Usage: "hyperctl agent schema describe <kind>", Handler: func(command *cli.Command) error {
		if err := oneArgument(command); err != nil {
			return fail(err)
		}
		document, err := schemas.Describe(command.Args[0])
		if err != nil {
			return fail(err)
		}
		fmt.Print(string(document))
		return nil
	}}}
}

func oneArgument(command *cli.Command) error {
	if len(command.Args) != 1 {
		return fmt.Errorf("expected exactly 1 argument, got %d", len(command.Args))
	}
	return nil
}

func noArguments(command *cli.Command) error {
	if len(command.Args) != 0 {
		return fmt.Errorf("expected no arguments, got %d", len(command.Args))
	}
	return nil
}

func schemaNames() []string { return schemas.Names() }
