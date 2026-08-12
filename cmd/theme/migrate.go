package theme

import (
	"errors"
	"os"

	"github.com/herdworks/clif"
	"github.com/herdworks/httt/migrate"
)

var errMigrationIncomplete = errors.New("migration incomplete")

func ThemeMigrateCommand() *clif.Command {
	command := &clif.Command{
		Config: clif.CommandConfig{
			Name:        "migrate",
			Description: "Migrate a theme to current directive syntax.",
			Usage:       "hyperctl theme migrate --theme theme.json [options]",
		},
		Options: clif.CommandOptions{
			{
				Name:        "theme",
				ShortName:   "t",
				Description: "Path to the theme configuration file.",
				Required:    true,
			},
			{
				Name:        "output",
				Description: "Write a complete migrated theme to an empty output directory.",
			},
			{
				Name:        "interactive",
				Description: "Review each migration decision interactively.",
				Boolean:     true,
			},
		},
	}
	command.Config.Handler = func(command *clif.Command) error {
		status, err := (migrate.Migration{
			ThemePath:   command.Options.Get("theme"),
			OutputDir:   command.Options.Get("output"),
			Interactive: command.Options.Get("interactive") == "true",
			Input:       os.Stdin,
			Output:      os.Stdout,
		}).Run()
		if err == nil || errors.Is(err, migrate.ErrInterrupted) {
			if writeErr := status.Print(os.Stdout); writeErr != nil {
				return writeErr
			}
		}
		if err != nil {
			return err
		}
		if !status.Complete() {
			return errMigrationIncomplete
		}
		return nil
	}
	return command
}
