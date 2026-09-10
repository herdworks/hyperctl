package cli

import (
	"github.com/herdworks/clic"
	hsynccli "github.com/herdworks/hsync/cli"
	"github.com/herdworks/htxt"
)

// SyncCommand groups synchronization operations for a checked-out website.
func SyncCommand() clic.Command {
	command := hsynccli.Command(htxt.SyncConfigManager{})
	command.Spec.Key = "sync"
	command.Meta.Description = "Website synchronization tools"
	command.Meta.EnvPrefix = ""
	command.Meta.Order = 3
	return command
}
