package cli

import (
	"fmt"

	"github.com/herdworks/clic"
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

func SchemaCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "schema",
			Description: "Describe Hyperctl data and CLI contracts as JSON.",
			Usage: []string{
				"hyperctl agent schema [command]",
			},
		},
		Spec: clic.CommandSpec{
			Key:  "schema",
			Help: clic.DefaultHelp,
		},
		Commands: clic.MustCommandsFrom(
			SchemaListCommand,
			SchemaDescribeCommand,
		),
	}
}

func SchemaListCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "ls",
			Description: "List stable schema kinds.",
			Usage: []string{
				"hyperctl agent schema ls",
			},
		},
		Spec: clic.CommandSpec{
			Key:  "ls",
			Help: clic.DefaultHelp,
			Handler: func(execution *clic.Execution) error {
				if err := schemas.Validate(); err != nil {
					return err
				}
				return clic.WriteList(execution.Stdout, clic.ColumnarFormat, schemas.Entries(),
					clic.Field{Key: "Name", Header: "Name"},
					clic.Field{Key: "Description", Header: "Description"},
				)
			},
		},
	}
}

func SchemaDescribeCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "describe",
			Description: "Write one raw JSON Schema document.",
			Usage: []string{
				"hyperctl agent schema describe <kind>",
			},
		},
		Spec: clic.CommandSpec{
			Key:       "describe",
			Help:      clic.DefaultHelp,
			Validator: expectedPositionalArguments(1),
			Handler: func(execution *clic.Execution) error {
				document, err := schemas.Describe(execution.Command.State.Arguments[0])
				if err != nil {
					return err
				}
				_, err = fmt.Fprint(execution.Stdout, string(document))
				return err
			},
		},
	}
}

func schemaNames() []string { return schemas.Names() }
