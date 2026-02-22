package theme

import (
	"fmt"
	"net/url"

	cli "github.com/herdworks/cli-go"
	hypertemplates "github.com/herdworks/hypertemplates-go"
	hypertexting "github.com/herdworks/hypertexting-go"
)

func ThemeListCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "ls",
			Description: "List installed themes.",
			Usage:       "hyperctl theme ls [options]",
			Validator:   nil,
			Handler:     ThemeListCommandHandler,
		},
		Options: cli.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Required:    true,
				Secret:      false,
			},
		},
	}
}

func ThemeInstallCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "install",
			Description: "Install a theme from a package archive.",
			Usage:       "hyperctl theme install [options]",
			Handler:     ThemeInstallCommandHandler,
		},
		Options: cli.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Required:    true,
			},
			{
				Name:        "file",
				ShortName:   "f",
				Description: "Path to theme package file (zip archive).",
				Required:    true,
			},
		},
	}
}

func ThemeListCommandHandler(cmd *cli.Command) (err error) {
	var themes []hypertemplates.Theme
	var website *hypertexting.Website
	if uri, err := url.Parse(cmd.Options.Get("config")); err == nil && uri != nil {
		website = &hypertexting.Website{URI: uri}
	} else {
		return fmt.Errorf("invalid website config URI: %s", cmd.Options.Get("config"))
	}
	if err = website.Load(false, nil); err != nil {
		return err
	}
	if themes, err = website.ListThemes(); err != nil {
		return err
	}
	ThemeListCommandOutput(cmd, website, themes)
	return
}

func ThemeListCommandOutput(cmd *cli.Command, website *hypertexting.Website, themes Themes) {
	pathcol := cmd.ColumnWidth(themes, "Path", 1)
	namecol := cmd.ColumnWidth(themes, "Name", 1)
	versioncol := cmd.ColumnWidth(themes, "Version", 1)
	fmt.Printf("%-*s %-*s %-*s\n", pathcol, "Path", namecol, "Name", versioncol, "Version")
	fmt.Printf("%-*s %-*s %-*s\n", pathcol, "----", namecol, "----", versioncol, "-------")
	for i := range themes {
		fmt.Printf(
			"%-*s %-*s %-*s\n",
			pathcol, themes[i].URI.Path,
			namecol, themes[i].Name,
			versioncol, themes[i].Version,
		)
	}
}

