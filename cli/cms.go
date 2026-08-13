package cli

import (
	"github.com/herdworks/clic"
)

// CMSCommand manages website content.
func CMSCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "cms",
			Description: "Content management tools",
			Usage: []string{
				"hyperctl cms [command]",
			},
		},
		Spec: clic.CommandSpec{
			Key:  "cms",
			Help: clic.DefaultHelp,
		},
		Commands: clic.MustCommandsFrom(
			AssetCommand,
			ContentTypeCommand,
			DataCommand,
			CMSNewCommand,
			PageCommand,
		),
	}
}
