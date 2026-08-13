package cli

import (
	"github.com/herdworks/clic"
	hsynccli "github.com/herdworks/hsync/cli"
)

// SyncCommand groups synchronization operations for a checked-out website.
func SyncCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "sync",
			Description: "Synchronize checked-out websites.",
			Usage: []string{
				"hyperctl sync [command] [options]",
			},
		},
		Spec: clic.CommandSpec{
			Key:  "sync",
			Help: clic.DefaultHelp,
		},
		Commands: clic.MustCommandsFrom(
			DownloadCommand,
			InitSyncCommand,
			PullCommand,
			PushCommand,
			ResetCommand,
			StatusCommand,
		),
	}
}

func DownloadCommand() clic.Command {
	command := hsynccli.DownloadCommand()
	command.Meta.Usage = []string{"hyperctl sync download <url> <directory>"}
	return command
}

func InitSyncCommand() clic.Command {
	command := hsynccli.InitCommand()
	command.Meta.Usage = []string{"hyperctl sync init <remote> [--config site.yaml]"}
	command.Spec.Options["config"].Meta.Description = "Configuration file to create."
	command.Spec.Options["config"].Meta.EnvName = "HYPER_CONFIG"
	command.Spec.Options["config"].Meta.Default = "site.yaml"
	return command
}

func PullCommand() clic.Command {
	command := hsynccli.PullCommand()
	command.Meta.Usage = []string{"hyperctl sync pull [directory] [--config path]"}
	command.Spec.Options["config"].Meta.Description = "Website YAML or JSON configuration."
	command.Spec.Options["config"].Meta.EnvName = "HYPER_CONFIG"
	command.Spec.Options["config"].Meta.Default = nil
	return command
}

func PushCommand() clic.Command {
	command := hsynccli.PushCommand()
	command.Meta.Usage = []string{"hyperctl sync push", "hyperctl sync push [--config path]", "hyperctl sync push [--config path] [--environment environment]", "hyperctl sync push [--config path] [--providers provider,provider]"}
	command.Spec.Options["config"].Meta.Description = "Website YAML or JSON configuration."
	command.Spec.Options["config"].Meta.EnvName = "HYPER_CONFIG"
	command.Spec.Options["config"].Meta.Default = nil
	command.Spec.Options["environment"].Meta.Description = "Deployment environment from config.environments."
	command.Spec.Options["providers"].Meta.Description = "Comma-separated provider names from config.providers."
	return command
}

func ResetCommand() clic.Command {
	command := hsynccli.ResetCommand()
	command.Meta.Usage = []string{"hyperctl sync reset [directory] [--config path] [--interactive]"}
	command.Spec.Options["config"].Meta.Description = "Website YAML or JSON configuration."
	command.Spec.Options["config"].Meta.EnvName = "HYPER_CONFIG"
	command.Spec.Options["config"].Meta.Default = nil
	return command
}

func StatusCommand() clic.Command {
	command := hsynccli.StatusCommand()
	command.Meta.Usage = []string{"hyperctl sync status [directory] [--config path]"}
	command.Spec.Options["config"].Meta.Description = "Website YAML or JSON configuration."
	command.Spec.Options["config"].Meta.EnvName = "HYPER_CONFIG"
	command.Spec.Options["config"].Meta.Default = nil
	return command
}
