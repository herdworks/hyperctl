package cli

import (
	"fmt"

	"github.com/herdworks/clic"
	htskills "github.com/herdworks/htskills"
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
				content, err := htskills.Readme()
				if err != nil {
					return err
				}
				_, err = fmt.Fprint(execution.Stdout, string(content))
				return err
			},
		},
	}
}
