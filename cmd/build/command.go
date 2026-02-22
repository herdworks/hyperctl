package build

import (
	cli "github.com/herdworks/cli-go"
)

func Command() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "build",
			Description: "Build management tools",
			Usage:       "hyperctl build [command] [options] <path>",
			Validator:   nil,
			Handler:     CommandHandler,
		},
		Options: cli.CommandOptions{},
		SubCommands: cli.Commands{
			BuildCommand(),
			GenerateCommand(),
		},
	}
}

func CommandHandler(cmd *cli.Command) (err error) {
	cmd.Help(nil)
	return nil
}
