package theme

import "github.com/herdworks/clif"

func Command() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "theme",
			Description: "Theme management tools",
			Usage:       "hyperctl theme [command] [options]",
			Handler:     ThemeCommandHandler,
		},
		Options: clif.CommandOptions{},
		SubCommands: clif.Commands{
			ThemeListCommand(),
			ThemeInstallCommand(),
			ThemePackageCommand(),
			ThemeMigrateCommand(),
		},
	}
}

func ThemeCommandHandler(cmd *clif.Command) (err error) {
	cmd.Help(nil)
	return nil
}
