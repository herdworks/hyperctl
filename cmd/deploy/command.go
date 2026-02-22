package deploy

import (
	cli "github.com/herdworks/cli-go"
)

func Command() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "deploy",
			Description: "Deployment management tools",
			Usage:       "hyperctl deploy [command] [options] <path>",
			Validator:   nil,
			Handler:     CommandHandler,
		},
		Options: cli.CommandOptions{},
		SubCommands: cli.Commands{
			DeployCommand(),
			PublishCommand(),
		},
	}
}

func CommandHandler(cmd *cli.Command) (err error) {
	cmd.Help(nil)
	return nil
}
