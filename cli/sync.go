package cli

import (
	"github.com/herdworks/clic"
	hsynccli "github.com/herdworks/hsync/cli"
)

// SyncCommand groups synchronization operations for a checked-out website.
func SyncCommand() clic.Command {
	command := hsynccli.Command()
	command.Meta.Description = "Website synchronization tools"
	command.Meta.EnvPrefix = ""
	command.Meta.Order = 5
	command.Spec.Key = "sync"
	return command
}
