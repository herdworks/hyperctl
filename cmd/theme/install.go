package theme

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	cli "github.com/herdworks/cli-go"
	hypertemplates "github.com/herdworks/hypertemplates-go"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
)

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
				Default:     "site.yaml",
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

func ThemeInstallCommandHandler(cmd *cli.Command) (err error) {
	var website *hypertexting.Website
	var theme hypertemplates.Theme
	var reader *zip.ReadCloser

	// Load the website configuration.
	if uri, err := href.Parse(cmd.Options.Get("config")); err == nil && uri != nil {
		website = &hypertexting.Website{URI: uri}
	} else {
		return fmt.Errorf("invalid website config URI: %s", cmd.Options.Get("config"))
	}
	if err = website.Load(false, nil); err != nil {
		return err
	}

	// Open the theme package.
	filePath := cmd.Options.Get("file")
	if reader, err = zip.OpenReader(filePath); err != nil {
		return fmt.Errorf("failed to open package %q: %w", filePath, err)
	}
	defer reader.Close()

	// Read and validate theme metadata from the package.
	if theme, err = readThemeFromZip(reader); err != nil {
		return err
	}

	// Construct the install path and check for existing installation.
	installPath := filepath.Join(website.ThemesDirPath(), theme.Name, theme.Version)
	if _, err = os.Stat(installPath); err == nil {
		return fmt.Errorf("theme %q version %q is already installed at %s", theme.Name, theme.Version, installPath)
	}

	// Extract the package contents into the install path.
	fmt.Printf("installing theme %q version %q to %s\n", theme.Name, theme.Version, installPath)
	if err = extractZip(reader, installPath); err != nil {
		return fmt.Errorf("failed to extract package: %w", err)
	}

	// Verify installation by confirming theme.json exists at the install path.
	themeFile := filepath.Join(installPath, "theme.json")
	if _, err = os.Stat(themeFile); err != nil {
		return fmt.Errorf("installation verification failed: %s not found", themeFile)
	}

	fmt.Printf("theme %q version %q installed successfully\n", theme.Name, theme.Version)
	return nil
}

// readThemeFromZip finds and unmarshals theme.json from the root of a zip archive.
func readThemeFromZip(reader *zip.ReadCloser) (theme hypertemplates.Theme, err error) {
	for _, f := range reader.File {
		if f.Name != "theme.json" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return theme, fmt.Errorf("failed to read theme.json from package: %w", err)
		}
		defer rc.Close()

		data, err := io.ReadAll(rc)
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

// extractZip extracts all entries from a zip archive into destDir.
func extractZip(reader *zip.ReadCloser, destDir string) error {
	for _, f := range reader.File {
		target := filepath.Join(destDir, f.Name)

		// Zip slip protection: ensure the target path stays within destDir.
		rel, err := filepath.Rel(destDir, target)
		if err != nil || strings.HasPrefix(rel, "..") {
			return fmt.Errorf("invalid file path in archive: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}

		// Ensure the parent directory exists.
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		out, err := os.Create(target)
		if err != nil {
			rc.Close()
			return err
		}

		if _, err = io.Copy(out, rc); err != nil {
			out.Close()
			rc.Close()
			return err
		}

		out.Close()
		rc.Close()
	}
	return nil
}
