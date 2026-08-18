// Package cli composes Hyperctl's importable command-line interface.
package cli

import (
	"github.com/herdworks/clic"
	htskillscli "github.com/herdworks/htskills/cli"
)

// AgentCommand exposes coding-agent guidance and contracts.
func AgentCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "agent",
			Description: "Coding agent tools",
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
			SchemaCommand,
			func() clic.Command {
				command := htskillscli.SkillCommand()
				command.Meta.Name = "skill"
				command.Meta.Usage = []string{"hyperctl agent skill [command]"}
				command.Spec.Key = "skill"
				return command
			},
		),
		Links: clic.Links{
			{From: "readme", To: []string{"skill", "readme"}},
		},
	}
}
