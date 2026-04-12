package dev

import (
	cli "github.com/herdworks/cli-go"
)

func Command() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "dev",
			Description: "Developer tools",
			Usage:       "hyperctl dev [command] [options] <path>",
			Validator:   nil,
			Handler:     CommandHandler,
		},
		Options: cli.CommandOptions{},
		SubCommands: cli.Commands{
			GraphCommand(),
			MimetypeCommand(),
			RenderCommand(),
			ServerCommand(),
		},
	}
}

func CommandHandler(cmd *cli.Command) (err error) {
	cmd.Help(nil)
	return nil
}
