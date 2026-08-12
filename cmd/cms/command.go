package cms

import "github.com/herdworks/clif"

func Command() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "cms",
			Description: "Content management tools",
			Usage:       "hyperctl cms [command] [options]",
			Handler:     NewCommandHandler,
		},
		Options: clif.CommandOptions{},
		SubCommands: clif.Commands{
			NewCommand(),
			AssetCommand(),
			ContentTypeCommand(),
			DataCommand(),
			PageCommand(),
		},
	}
}
