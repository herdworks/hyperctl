package cmd

import (
	"fmt"
	"net/url"

	cli "github.com/herdworks/cli-go"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
)

func ConfigCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "config",
			Description: "View HyperTemplates configuration.",
			Usage:       "hyperctl config [options]",
			Validator:   ConfigCommandValidator,
			Handler:     ConfigCommandHandler,
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
			{
				Name:        "id",
				ShortName:   "i",
				Description: "Identifier to use for the build.",
				Required:    false,
				Secret:      false,
			},
		},
	}
}

func ConfigCommandValidator(cmd *cli.Command) (err error) {
	return nil
}

func ConfigCommandHandler(cmd *cli.Command) (err error) {
	website := hypertexting.Website{URI: &href.Href{URL: url.URL{Path: cmd.Options.Get("config")}}}
	if err = website.Load(true, nil, nil); err != nil {
		return err
	}
	fmt.Printf("%v\n", website)
	return nil
}
