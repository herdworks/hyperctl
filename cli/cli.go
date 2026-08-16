// Package cli composes Hyperctl's importable command-line interface.
package cli

import "github.com/herdworks/clic"

var CLI = clic.CLI{
	Command: clic.Command{
		Meta: clic.CommandMeta{
			Description: "The HyperTexting CLI.",
			EnvPrefix:   "HYPER_",
		},
		Spec: clic.CommandSpec{
			Key:     "hyperctl",
			Help:    clic.DefaultHelp,
			Options: clic.Options{},
		},
		Commands: clic.MustCommandsFrom(
			AgentCommand,
			WebsiteCommand,
			SyncCommand,
			ThemeCommand,
			VersionCommand,
		),
		Links: clic.Links{
			{From: "build", To: []string{"website", "build"}},
			{From: "deploy", To: []string{"website", "deploy"}},
			{From: "download", To: []string{"sync", "download"}},
			{From: "pull", To: []string{"sync", "pull"}},
			{From: "push", To: []string{"sync", "push"}},
			{From: "server", To: []string{"website", "server"}},
		},
	},
}

// global options

func GlobalVerboseOption() clic.Option {
	return clic.Option{
		Meta: clic.OptionMeta{
			ShortName:   "v",
			Description: "Enable verbose mode",
		},
		Spec: clic.OptionSpec{
			Key:      "verbose",
			Kind:     clic.BoolOption,
			Required: false,
		},
	}
}

// global validators

// expectedPositionalArguments closure
func expectedPositionalArguments(count int) func(*clic.Execution) error {
	return func(execution *clic.Execution) error {
		if execution.Command.PositionalArguments() == count {
			return nil
		}
		return positionalArgumentsError(execution, count)
	}
}
