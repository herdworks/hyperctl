package theme

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	cli "github.com/herdworks/cli-go"
	hypertemplates "github.com/herdworks/hypertemplates-go"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
)

func ThemePackageCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "package",
			Description: "Package a theme into a zip archive.",
			Usage:       "hyperctl theme package [options]",
			Handler:     ThemePackageCommandHandler,
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

func ThemePackageCommandHandler(cmd *cli.Command) (err error) {
	var website *hypertexting.Website
	var theme hypertemplates.Theme = hypertemplates.Theme{URI: &url.URL{Path: cmd.Options.Get("theme")}}
	var archive string = cmd.Options.Get("file")

	// Load the website configuration
	if uri, err := href.Parse(cmd.Options.Get("config")); err == nil && uri != nil {
		website = &hypertexting.Website{URI: uri}
	} else {
		return fmt.Errorf("invalid website config URI: %s", cmd.Options.Get("config"))
	}
	if err = website.Load(false, nil); err != nil {
		return err
	}

	// Load the theme configuration
	if err = theme.Parse(); err == nil {
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
	if err = theme.Archive(archive); err != nil {
		return err
	}
	fmt.Printf("Theme packaged successfully: %s\n", archive)
	return nil
}
