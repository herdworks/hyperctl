// Package cli composes Hyperctl's importable command-line interface.
package cli

import (
	"github.com/herdworks/clic"
)

var CLI = clic.CLI{
	Command: clic.Command{
		Meta: clic.CommandMeta{
			Name:        "hyperctl",
			Description: "The HyperTexting CLI.",
		},
		Spec: clic.CommandSpec{
			Key:     "hyperctl",
			Help:    clic.DefaultHelp,
			Options: clic.Options{},
		},
		Commands: clic.MustCommandsFrom(
			AgentCommand,
			BuildCommand,
			CMSCommand,
			DeployCommand,
			DevCommand,
			SyncCommand,
			ThemeCommand,
			VersionCommand,
		),
	},
}

// global options
func GlobalVersionOption() clic.Option {
	return clic.Option{
		Meta: clic.OptionMeta{
			Name:        "version",
			Description: "Print version information and exit.",
		},
		Spec: clic.OptionSpec{
			Key:  "version",
			Kind: clic.BoolOption,
		},
	}
}
func GlobalWebsiteConfigOption() clic.Option {
	return clic.Option{
		Meta: clic.OptionMeta{
			Name:        "config",
			ShortName:   "c",
			EnvName:     "HYPER_CONFIG",
			Description: "Path to the website configuration file.",
			Default:     "site.yaml",
		},
		Spec: clic.OptionSpec{
			Key:      "config",
			Kind:     clic.StringOption,
			Required: true,
		},
	}
}

func GlobalThemeConfigOption() clic.Option {
	return clic.Option{
		Meta: clic.OptionMeta{
			Name:        "theme",
			ShortName:   "t",
			EnvName:     "HYPER_THEME",
			Description: "Path to the theme configuration file.",
			Default:     "theme/theme.json",
		},
		Spec: clic.OptionSpec{
			Key:      "theme",
			Kind:     clic.StringOption,
			Required: true,
		},
	}
}

func GlobalVerboseOption() clic.Option {
	return clic.Option{
		Meta: clic.OptionMeta{
			Name:        "verbose",
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
