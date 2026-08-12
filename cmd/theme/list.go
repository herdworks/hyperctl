package theme

import (
	"fmt"

	"github.com/herdworks/clif"
	httheme "github.com/herdworks/httt/theme"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

func ThemeListCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "ls",
			Description: "List installed themes.",
			Usage:       "hyperctl theme ls [options]",
			Validator:   nil,
			Handler:     ThemeListCommandHandler,
		},
		Options: clif.CommandOptions{
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

func ThemeListCommandHandler(cmd *clif.Command) (err error) {
	var themes []httheme.Theme
	var website *htxt.Website
	if uri, err := href.Parse(cmd.Options.Get("config")); err == nil && uri != nil {
		website = &htxt.Website{URI: uri}
	} else {
		return fmt.Errorf("invalid website config URI: %s", cmd.Options.Get("config"))
	}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	if themes, err = website.ListThemes(nil); err != nil {
		return err
	}
	ThemeListCommandOutput(cmd, website, themes)
	return
}

func ThemeListCommandOutput(cmd *clif.Command, website *htxt.Website, themes Themes) {
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
