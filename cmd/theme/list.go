package theme

import (
	"fmt"

	cli "github.com/herdworks/cli-go"
	hypertemplates "github.com/herdworks/hypertemplates-go"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
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
				Default:     "site.yaml",
				Required:    true,
				Secret:      false,
			},
		},
	}
}

func ThemeListCommandHandler(cmd *cli.Command) (err error) {
	var themes []hypertemplates.Theme
	var website *hypertexting.Website
	if uri, err := href.Parse(cmd.Options.Get("config")); err == nil && uri != nil {
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
