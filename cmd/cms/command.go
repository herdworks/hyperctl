package cms

import cli "github.com/herdworks/cli-go"

func Command() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "cms",
			Description: "Content management tools",
			Usage:       "hyperctl cms [command] [options]",
			Handler:     NewCommandHandler,
		},
		Options: cli.CommandOptions{},
		SubCommands: cli.Commands{
			NewCommand(),
			AssetCommand(),
			ContentTypeCommand(),
			DataCommand(),
			PageCommand(),
		},
	}
}
