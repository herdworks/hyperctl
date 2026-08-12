package theme

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/herdworks/clif"
	httheme "github.com/herdworks/httt/theme"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

func ThemePackageCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "package",
			Description: "Package a theme into a zip archive.",
			Usage:       "hyperctl theme package [options]",
			Handler:     ThemePackageCommandHandler,
		},
		Options: clif.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Default:     "site.yaml",
				Required:    true,
			},
			{
				Name:        "theme",
				ShortName:   "t",
				EnvName:     "HYPER_THEME",
				Description: "Path to the theme configuration file.",
				Required:    true,
				Default:     "theme.json",
			},
			{
				Name:        "file",
				ShortName:   "f",
				Description: "Output file path for the theme package (zip archive).",
				Required:    false,
			},
		},
	}
}

func ThemePackageCommandHandler(cmd *clif.Command) (err error) {
	var website *htxt.Website
	var theme httheme.Theme = httheme.Theme{URI: &url.URL{Path: cmd.Options.Get("theme")}}
	var archive string = cmd.Options.Get("file")

	// Load the website configuration
	if uri, err := href.Parse(cmd.Options.Get("config")); err == nil && uri != nil {
		website = &htxt.Website{URI: uri}
	} else {
		return fmt.Errorf("invalid website config URI: %s", cmd.Options.Get("config"))
	}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}

	// Load the theme configuration
	if err = theme.Parse(nil); err == nil {
		theme.ResolveConfig()
	} else {
		return err
	}

	// Ensure output file does not already exist
	if archive == "" {
		archive = fmt.Sprintf("%s-v%s.zip", theme.Name, theme.Version)
	}
	if _, statErr := os.Stat(archive); statErr == nil {
		return fmt.Errorf("error: %s already exists. Please remove the existing file before packaging", archive)
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return statErr
	}

	// Create the theme package archive
	if err = theme.Archive(archive, nil); err != nil {
		return err
	}
	fmt.Printf("Theme packaged successfully: %s\n", archive)
	return nil
}
