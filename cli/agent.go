// Package cli composes Hyperctl's importable command-line interface.
package cli

import (
	"github.com/herdworks/clic"
)

// AgentCommand exposes coding-agent guidance and contracts.
func AgentCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "agent",
			Description: "HyperTexting skills for coding agents.",
			Usage: []string{
				"hyperctl agent [command]",
			},
			Order: 1,
		},
		Spec: clic.CommandSpec{
			Key:  "agent",
			Help: clic.DefaultHelp,
		},
		Commands: clic.MustCommandsFrom(
			ReadmeCommand,
			SchemaCommand,
			SkillCommand,
		),
	}
}
