// Package cli composes Hyperctl's importable command-line interface.
package cli

import (
	"github.com/herdworks/clic"
	htskillscli "github.com/herdworks/htskills/cli"
)

// SkillCommand mounts the portable HyperTexting Agent Skills workflow.
func SkillCommand() clic.Command {
	command := htskillscli.SkillCommand()
	command.Meta.Name = "skill"
	command.Meta.Usage = []string{"hyperctl skill [command]"}
	command.Meta.Order = 4
	command.Spec.Key = "skill"
	return command
}
