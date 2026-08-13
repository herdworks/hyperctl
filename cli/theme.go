package cli

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/herdworks/clic"
	"github.com/herdworks/httt/migrate"
	httheme "github.com/herdworks/httt/theme"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

// ThemeCommand manages website themes.
func ThemeCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "theme",
			Description: "Theme management tools",
			Usage: []string{
				"hyperctl theme [command]",
			},
		},
		Spec: clic.CommandSpec{
			Key:  "theme",
			Help: clic.DefaultHelp,
		},
		Commands: clic.MustCommandsFrom(
			ThemeListCommand,
			ThemeInstallCommand,
			ThemePackageCommand,
			ThemeMigrateCommand,
		),
	}
}

func ThemeListCommand() clic.Command {
	var config clic.Option = GlobalWebsiteConfigOption()
	var format clic.Option = clic.Option{
		Meta: clic.OptionMeta{Name: "format", ShortName: "f", Description: "Output format: columns or json.", Default: string(clic.ColumnarFormat)},
		Spec: clic.OptionSpec{Key: "format", Kind: clic.StringOption},
	}
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "ls",
			Description: "List installed themes.",
			Usage: []string{
				"hyperctl theme ls [options]",
			},
		},
		Spec: clic.CommandSpec{
			Key: "ls", Help: clic.DefaultHelp,
			Options: clic.Options{
				config.Spec.Key: &config,
				format.Spec.Key: &format,
			},
			Handler: themeListHandler,
		},
	}
}

func ThemeInstallCommand() clic.Command {
	var config clic.Option = GlobalWebsiteConfigOption()
	var file clic.Option = ThemeFileOption()
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "install",
			Description: "Install a theme from a package archive.",
			Usage: []string{
				"hyperctl theme install [options]",
			},
		},
		Spec: clic.CommandSpec{
			Key: "install", Help: clic.DefaultHelp,
			Options: clic.Options{
				config.Spec.Key: &config,
				file.Spec.Key:   &file,
			},
			Handler: themeInstallHandler,
		},
	}
}

func ThemePackageCommand() clic.Command {
	var config clic.Option = GlobalWebsiteConfigOption()
	var theme clic.Option = GlobalThemeConfigOption()
	var file clic.Option = ThemeFileOption()
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "package",
			Description: "Package a theme into a zip archive.",
			Usage: []string{
				"hyperctl theme package [options]",
			},
		},
		Spec: clic.CommandSpec{
			Key: "package", Help: clic.DefaultHelp,
			Options: clic.Options{
				config.Spec.Key: &config,
				theme.Spec.Key:  &theme,
				file.Spec.Key:   &file,
			},
			Handler: themePackageHandler,
		},
	}
}

func ThemeMigrateCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "migrate",
			Description: "Migrate a theme to current directive syntax.",
			Usage: []string{
				"hyperctl theme migrate --theme theme.json [options]",
			},
		},
		Spec: clic.CommandSpec{
			Key: "migrate", Help: clic.DefaultHelp,
			Options: clic.Options{
				"theme":       {Meta: clic.OptionMeta{Name: "theme", ShortName: "t", Description: "Path to the theme configuration file."}, Spec: clic.OptionSpec{Key: "theme", Kind: clic.StringOption, Required: true}},
				"output":      {Meta: clic.OptionMeta{Name: "output", Description: "Write a complete migrated theme to an empty output directory."}, Spec: clic.OptionSpec{Key: "output", Kind: clic.StringOption}},
				"interactive": {Meta: clic.OptionMeta{Name: "interactive", Description: "Review each migration decision interactively.", Default: false}, Spec: clic.OptionSpec{Key: "interactive", Kind: clic.BoolOption}},
			},
			Handler: themeMigrateHandler,
		},
	}
}

func ThemeFileOption() clic.Option {
	return clic.Option{
		Meta: clic.OptionMeta{
			Name:        "file",
			ShortName:   "f",
			Description: "Path to theme package file (zip archive).",
		},
		Spec: clic.OptionSpec{
			Key:      "file",
			Kind:     clic.StringOption,
			Required: false,
		},
	}
}

func themeInstallHandler(execution *clic.Execution) error {
	options := execution.Options
	uri, err := href.Parse(options["config"].StringValue())
	if err != nil || uri == nil {
		return fmt.Errorf("invalid website config URI: %s", options["config"].StringValue())
	}
	website := &htxt.Website{URI: uri}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	filePath := options["file"].StringValue()
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return fmt.Errorf("failed to open package %q: %w", filePath, err)
	}
	defer reader.Close()
	theme, err := readThemeFromZip(reader)
	if err != nil {
		return err
	}
	installPath := filepath.Join(website.ThemesDirPath(), theme.Name, theme.Version)
	if _, err = os.Stat(installPath); err == nil {
		return fmt.Errorf("theme %q version %q is already installed at %s", theme.Name, theme.Version, installPath)
	}
	if _, err = fmt.Fprintf(execution.Stdout, "installing theme %q version %q to %s\n", theme.Name, theme.Version, installPath); err != nil {
		return err
	}
	if err = extractZip(reader, installPath); err != nil {
		return fmt.Errorf("failed to extract package: %w", err)
	}
	themeFile := filepath.Join(installPath, "theme.json")
	if _, err = os.Stat(themeFile); err != nil {
		return fmt.Errorf("installation verification failed: %s not found", themeFile)
	}
	_, err = fmt.Fprintf(execution.Stdout, "theme %q version %q installed successfully\n", theme.Name, theme.Version)
	return err
}

func readThemeFromZip(reader *zip.ReadCloser) (theme httheme.Theme, err error) {
	for _, file := range reader.File {
		if file.Name != "theme.json" {
			continue
		}
		content, err := file.Open()
		if err != nil {
			return theme, fmt.Errorf("failed to read theme.json from package: %w", err)
		}
		defer content.Close()
		data, err := io.ReadAll(content)
		if err != nil {
			return theme, fmt.Errorf("failed to read theme.json from package: %w", err)
		}
		if err = json.Unmarshal(data, &theme); err != nil {
			return theme, fmt.Errorf("failed to parse theme.json: %w", err)
		}
		if theme.Name == "" {
			return theme, fmt.Errorf("theme.json is missing required field \"name\"")
		}
		if theme.Version == "" {
			return theme, fmt.Errorf("theme.json is missing required field \"version\"")
		}
		return theme, nil
	}
	return theme, fmt.Errorf("package does not contain a theme.json file")
}

func extractZip(reader *zip.ReadCloser, destination string) error {
	for _, file := range reader.File {
		target := filepath.Join(destination, file.Name)
		relative, err := filepath.Rel(destination, target)
		if err != nil || strings.HasPrefix(relative, "..") {
			return fmt.Errorf("invalid file path in archive: %s", file.Name)
		}
		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}
		if err = os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		input, err := file.Open()
		if err != nil {
			return err
		}
		output, err := os.Create(target)
		if err != nil {
			input.Close()
			return err
		}
		if _, err = io.Copy(output, input); err != nil {
			output.Close()
			input.Close()
			return err
		}
		output.Close()
		input.Close()
	}
	return nil
}

func themeListHandler(execution *clic.Execution) error {
	config := execution.Options["config"].StringValue()
	uri, err := href.Parse(config)
	if err != nil || uri == nil {
		return fmt.Errorf("invalid website config URI: %s", config)
	}
	website := &htxt.Website{URI: uri}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	themes, err := website.ListThemes(nil)
	if err != nil {
		return err
	}
	return clic.WriteList(execution.Stdout, clic.Format(execution.Options["format"].StringValue()), Themes(themes),
		clic.Field{Key: "path", Header: "Path"},
		clic.Field{Key: "name", Header: "Name"},
		clic.Field{Key: "version", Header: "Version"},
	)
}

func themePackageHandler(execution *clic.Execution) error {
	options := execution.Options
	uri, err := href.Parse(options["config"].StringValue())
	if err != nil || uri == nil {
		return fmt.Errorf("invalid website config URI: %s", options["config"].StringValue())
	}
	website := &htxt.Website{URI: uri}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	theme := httheme.Theme{URI: &url.URL{Path: options["theme"].StringValue()}}
	if err = theme.Parse(nil); err != nil {
		return err
	}
	theme.ResolveConfig()
	archive := options["file"].StringValue()
	if archive == "" {
		archive = fmt.Sprintf("%s-v%s.zip", theme.Name, theme.Version)
	}
	if _, statErr := os.Stat(archive); statErr == nil {
		return fmt.Errorf("error: %s already exists. Please remove the existing file before packaging", archive)
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return statErr
	}
	if err = theme.Archive(archive, nil); err != nil {
		return err
	}
	_, err = fmt.Fprintf(execution.Stdout, "Theme packaged successfully: %s\n", archive)
	return err
}

func themeMigrateHandler(execution *clic.Execution) error {
	options := execution.Options
	status, err := (migrate.Migration{ThemePath: options["theme"].StringValue(), OutputDir: options["output"].StringValue(), Interactive: options["interactive"].BoolValue(), Input: execution.Stdin, Output: execution.Stdout}).Run()
	if err == nil || errors.Is(err, migrate.ErrInterrupted) {
		if writeErr := status.Print(execution.Stdout); writeErr != nil {
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

// TODO: this belongs in httt!!
type Themes []httheme.Theme

func (themes Themes) ColumnValues(field string) []string {
	values := []string{}
	for _, theme := range themes {
		switch field {
		case "Path", "path":
			values = append(values, theme.URI.Path)
		case "Name", "name":
			values = append(values, theme.Name)
		case "Version", "version":
			values = append(values, theme.Version)
		}
	}
	return values
}

var errMigrationIncomplete = errors.New("migration incomplete")
