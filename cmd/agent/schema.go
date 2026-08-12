package agent

import (
	"fmt"
	"strings"

	"github.com/herdworks/clif"
	hsync "github.com/herdworks/hsync"
	"github.com/herdworks/httt/hypermark"
	"github.com/herdworks/httt/markdown"
	"github.com/herdworks/httt/theme"
	"github.com/herdworks/htup"
	"github.com/herdworks/htxt"
	"github.com/herdworks/htxt/jsonschema"
)

var schemas = jsonschema.Registry{
	htxt.Attachment{},
	htxt.Byline{},
	htxt.FeedConfig{},
	htxt.FeedData{},
	htxt.FeedSelector{},
	htxt.HTMLLink{},
	htxt.HTMLMeta{},
	markdown.HyperMarkConfig{},
	hypermark.MentionConfig{},
	htxt.Page{},
	htup.HostingProvider{},
	htup.FilterRules{},
	htup.Secret{},
	hsync.SyncConfig{},
	hsync.Rules{},
	theme.Theme{},
	htxt.Website{},
	htxt.WebsiteConfig{},
	htxt.WebsiteFavicon{},
}

func schemaCommand() *clif.Command {
	return &clif.Command{Config: clif.CommandConfig{Name: "schema", Description: "Describe Hyperctl data and CLI contracts as JSON.", Usage: "hyperctl agent schema [command]", Handler: help}, SubCommands: clif.Commands{schemaListCommand(), schemaDescribeCommand()}}
}

func schemaListCommand() *clif.Command {
	return &clif.Command{Config: clif.CommandConfig{Name: "ls", Description: "List stable schema kinds.", Usage: "hyperctl agent schema ls", Handler: func(*clif.Command) error {
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
	namecol := (&clif.Command{}).ColumnWidth(entries, "Name", 1)
	var output strings.Builder
	fmt.Fprintf(&output, "%-*s %s\n", namecol, "Name", "Description")
	fmt.Fprintf(&output, "%-*s %s\n", namecol, "----", "-----------")
	for _, entry := range entries {
		fmt.Fprintf(&output, "%-*s %s\n", namecol, entry.Name, entry.Summary)
	}
	return output.String()
}

func schemaDescribeCommand() *clif.Command {
	return &clif.Command{Config: clif.CommandConfig{Name: "describe", Description: "Write one raw JSON Schema document.", Usage: "hyperctl agent schema describe <kind>", Handler: func(command *clif.Command) error {
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

func oneArgument(command *clif.Command) error {
	if len(command.Args) != 1 {
		return fmt.Errorf("expected exactly 1 argument, got %d", len(command.Args))
	}
	return nil
}

func noArguments(command *clif.Command) error {
	if len(command.Args) != 0 {
		return fmt.Errorf("expected no arguments, got %d", len(command.Args))
	}
	return nil
}

func schemaNames() []string { return schemas.Names() }
