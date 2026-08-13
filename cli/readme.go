package cli

import (
	"fmt"
	"strings"

	"github.com/herdworks/clic"
)

func ReadmeCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "readme",
			Description: "Show onboarding context for a new agent.",
			Usage: []string{
				"hyperctl agent readme",
			},
		},
		Spec: clic.CommandSpec{
			Key:       "readme",
			Help:      clic.DefaultHelp,
			Validator: expectedPositionalArguments(0),
			Handler: func(execution *clic.Execution) error {
				_, err := fmt.Fprint(execution.Stdout, readmeContent())
				return err
			},
		},
	}
}

func readmeContent() string {
	var content string = `
HyperTemplates renders pure HTML templates; the HyperTexting CMS organizes websites, pages, data, themes, and assets. Describe data and CLI contracts with "hyperctl agent schema ls" and "hyperctl agent schema describe <kind>". If you are unfamiliar with the project, start with "hyperctl agent skill ls" and "hyperctl agent skill describe <name>".

Use established commands directly for mutations: "hyperctl cms new website", "hyperctl cms new page", "hyperctl dev server", "hyperctl build", and "hyperctl deploy". When a build needs a predictable output path, pass its identifier directly with "hyperctl build --id <build-id>" instead of wrapping the command to capture an auto-generated ID. Template guidance uses ht-apply with template variables instead of legacy ht-content, ht-attr(s), ht-param, and ht-query.
	`
	return fmt.Sprintf("%s\n", strings.TrimSpace(content))
}
