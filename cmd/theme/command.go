package theme

import cli "github.com/herdworks/cli-go"

func Command() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "theme",
			Description: "Theme management tools",
			Usage:       "hyperctl theme [command] [options]",
			Handler:     ThemeCommandHandler,
		},
		Options: cli.CommandOptions{},
		SubCommands: cli.Commands{
			ThemeListCommand(),
			ThemeInstallCommand(),
		},
	}
}

func ThemeCommandHandler(cmd *cli.Command) (err error) {
	cmd.Help(nil)
	return nil
}
