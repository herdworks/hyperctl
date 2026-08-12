package dev

import (
	"github.com/herdworks/clif"
)

func Command() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "dev",
			Description: "Developer tools",
			Usage:       "hyperctl dev [command] [options] <path>",
			Validator:   nil,
			Handler:     CommandHandler,
		},
		Options: clif.CommandOptions{},
		SubCommands: clif.Commands{
			GraphCommand(),
			MimetypeCommand(),
			RenderCommand(),
			ServerCommand(),
		},
	}
}

func CommandHandler(cmd *clif.Command) (err error) {
	cmd.Help(nil)
	return nil
}
