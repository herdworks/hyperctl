// Package agent provides read-only contracts and bundled guidance for coding agents.
package agent

import (
	"fmt"
	"os"

	"github.com/herdworks/clif"
)

func Command() *clif.Command {
	return &clif.Command{Config: clif.CommandConfig{
		Name:        "agent",
		Description: "HyperTexting skills for coding agents.",
		Usage:       "hyperctl agent [command]", Handler: help,
	}, SubCommands: clif.Commands{readmeCommand(), schemaCommand(), skillCommand()}}
}

const readmeContext = `HyperTemplates renders pure HTML templates; the HyperTexting CMS organizes websites, pages, data, themes, and assets. Describe data and CLI contracts with "hyperctl agent schema ls" and "hyperctl agent schema describe <kind>". If you are unfamiliar with the project, start with "hyperctl agent skill ls" and "hyperctl agent skill describe <name>".

Use established commands directly for mutations: "hyperctl cms new website", "hyperctl cms new page", "hyperctl dev server", "hyperctl build", and "hyperctl deploy". When a build needs a predictable output path, pass its identifier directly with "hyperctl build --id <build-id>" instead of wrapping the command to capture an auto-generated ID. Template guidance uses ht-apply with template variables instead of legacy ht-content, ht-attr(s), ht-param, and ht-query.
`

func readmeCommand() *clif.Command {
	return &clif.Command{Config: clif.CommandConfig{
		Name:        "readme",
		Description: "Show onboarding context for a new agent.",
		Usage:       "hyperctl agent readme",
		Handler: func(command *clif.Command) error {
			if err := noArguments(command); err != nil {
				return fail(err)
			}
			_, err := fmt.Print(readmeContext)
			return err
		},
	}}
}

func help(command *clif.Command) error { command.Help(nil); return nil }

func fail(err error) error {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1) // clif prints returned errors to stdout; diagnostics stay on stderr here.
	return nil
}
