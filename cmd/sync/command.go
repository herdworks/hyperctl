// Package sync exposes the hsync checkout workflow through hyperctl.
package sync

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	cli "github.com/herdworks/cli-go"
	hsync "github.com/herdworks/hsync"
)

// Command groups synchronization operations for a checked-out website.
func Command() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "sync",
			Description: "Synchronize checked-out websites.",
			Usage:       "hyperctl sync [command] [options]",
			Handler: func(cmd *cli.Command) error {
				cmd.Help(nil)
				return nil
			},
		},
		SubCommands: cli.Commands{
			downloadCommand(),
			initCommand(),
			pullCommand(),
			pushCommand(),
			resetCommand(),
			statusCommand(),
		},
	}
}

func downloadCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{
		Name:        "download",
		Description: "Download a website.",
		Usage:       "hyperctl sync download <url> <directory>",
		Validator: func(cmd *cli.Command) error {
			if len(cmd.Args) != 2 {
				return fmt.Errorf("download requires a URL and destination directory")
			}
			return nil
		},
		Handler: func(cmd *cli.Command) error {
			source, err := parseHTTPURL(cmd.Args[0])
			if err != nil {
				return fmt.Errorf("invalid URL %q", cmd.Args[0])
			}
			result, err := hsync.DownloadWithOptions(context.Background(), source, cmd.Args[1], hsync.DownloadOptions{Output: os.Stdout})
			if err != nil {
				return err
			}
			fmt.Println(result.ManifestID)
			return nil
		},
	}}
}

func initCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{
		Name:        "init",
		Description: "Initialize a website for synchronization.",
		Usage:       "hyperctl sync init <remote> [--config site.yaml]",
		Validator: func(cmd *cli.Command) error {
			if len(cmd.Args) != 1 {
				return fmt.Errorf("init requires one remote URL")
			}
			return nil
		},
		Handler: func(cmd *cli.Command) error {
			return hsync.InitializeConfig(cmd.Options.Get("config"), cmd.Args[0])
		},
	}, Options: configOption("Configuration file to create.")}
}

func pullCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{
		Name:        "pull",
		Description: "Pull remote changes into the workspace.",
		Usage:       "hyperctl sync pull [directory] [--config site.yaml]",
		Validator:   atMostOneDirectory,
		Handler: func(cmd *cli.Command) error {
			config, err := commandSyncConfig(commandDirectory(cmd), cmd.Options.Get("config"))
			if err != nil {
				return err
			}
			result, err := hsync.PullWithOptions(context.Background(), config.Directory(), *config.Sync, hsync.PullOptions{Output: os.Stdout})
			if err != nil {
				return err
			}
			fmt.Println(result.ManifestID)
			return nil
		},
	}, Options: configOption("Website YAML or JSON configuration.")}
}

func pushCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{
		Name:        "push",
		Description: "Push local changes to the remote.",
		Usage:       "hyperctl sync push [directory] [--config site.yaml] [options]",
		Validator: func(cmd *cli.Command) error {
			if err := atMostOneDirectory(cmd); err != nil {
				return err
			}
			if cmd.Options.Get("environment") != "" && cmd.Options.Get("providers") != "" {
				return fmt.Errorf("--providers and --environment are mutually exclusive")
			}
			return nil
		},
		Handler: func(cmd *cli.Command) error {
			config, err := commandSyncConfig(commandDirectory(cmd), cmd.Options.Get("config"))
			if err != nil {
				return err
			}
			return hsync.Push(context.Background(), hsync.PushOptions{
				ConfigPath:  config.Path,
				Environment: cmd.Options.Get("environment"),
				Providers:   hsync.SplitPatterns(cmd.Options.Get("providers")),
			})
		},
	}, Options: append(configOption("Website YAML or JSON configuration."),
		&cli.CommandOption{Name: "environment", ShortName: "e", Description: "Deployment environment from config.environments."},
		&cli.CommandOption{Name: "providers", Description: "Comma-separated provider names from config.providers."},
	)}
}

func resetCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{
		Name:        "reset",
		Description: "Reset modified files in the workspace.",
		Usage:       "hyperctl sync reset [directory] [--config site.yaml] [--interactive]",
		Validator:   atMostOneDirectory,
		Handler: func(cmd *cli.Command) error {
			config, err := commandSyncConfig(commandDirectory(cmd), cmd.Options.Get("config"))
			if err != nil {
				return err
			}
			repository, remote, err := hsync.OpenCheckout(config.Directory(), *config.Sync, nil)
			if err != nil {
				return err
			}
			result, err := repository.Reset(context.Background(), remote, hsync.ResetOptions{Output: os.Stdout, Interactive: cmd.Options.Get("interactive") == "true"})
			if err != nil {
				return err
			}
			fmt.Println(result.ManifestID)
			return nil
		},
	}, Options: append(configOption("Website YAML or JSON configuration."), &cli.CommandOption{Name: "interactive", ShortName: "i", Description: "Choose each tracked file to restore.", Boolean: true})}
}

func statusCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{
		Name:        "status",
		Description: "List modified files in the workspace.",
		Usage:       "hyperctl sync status [directory] [--config site.yaml]",
		Validator:   atMostOneDirectory,
		Handler: func(cmd *cli.Command) error {
			config, err := commandSyncConfig(commandDirectory(cmd), cmd.Options.Get("config"))
			if err != nil {
				return err
			}
			repository, err := hsync.OpenRepository(config.Directory())
			if err != nil {
				return err
			}
			repository.Config = config.Sync
			result, err := repository.Status()
			if err != nil {
				return err
			}
			for _, change := range result.Changes {
				fmt.Printf("%10.10s %s\n", change.Kind, change.Path)
			}
			return nil
		},
	}, Options: configOption("Website YAML or JSON configuration.")}
}

func configOption(description string) cli.CommandOptions {
	return cli.CommandOptions{{Name: "config", ShortName: "c", EnvName: "HYPER_CONFIG", Description: description, Default: "site.yaml"}}
}

func commandSyncConfig(directory, configPath string) (hsync.RepositoryConfig, error) {
	if configPath == "" {
		configPath = "site.yaml"
	}
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(directory, configPath)
	}
	config, _, err := hsync.ReadConfig(configPath)
	if err != nil {
		return hsync.RepositoryConfig{}, err
	}
	if config.Sync == nil {
		return hsync.RepositoryConfig{}, fmt.Errorf("configuration sync is empty")
	}
	if err := config.Sync.Validate(); err != nil {
		return hsync.RepositoryConfig{}, err
	}
	return config, nil
}

func atMostOneDirectory(cmd *cli.Command) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("%s accepts at most one directory", cmd.Config.Name)
	}
	return nil
}

func commandDirectory(cmd *cli.Command) string {
	if len(cmd.Args) == 1 {
		return cmd.Args[0]
	}
	return "."
}

func parseHTTPURL(value string) (*url.URL, error) {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid HTTP URL")
	}
	return parsed, nil
}
