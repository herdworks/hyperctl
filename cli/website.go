package cli

import (
	"github.com/herdworks/clic"
	htxtcli "github.com/herdworks/htxt/cli"
)

// WebsiteCommand mounts htxt's website workflow under Hyperctl.
func WebsiteCommand() clic.Command {
	command := htxtcli.Command()
	command.Meta.Aliases = []string{"htxt", "cms", "dev"} // new htxt + legacy cms alias
	command.Meta.Name = "website"
	command.Meta.Description = "Website management tools"
	command.Meta.Order = 1
	command.Meta.EnvPrefix = ""
	return command
}
