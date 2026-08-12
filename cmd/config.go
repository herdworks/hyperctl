package cmd

import (
	"fmt"
	"net/url"

	"github.com/herdworks/clif"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

func ConfigCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "config",
			Description: "View HyperTemplates configuration.",
			Usage:       "hyperctl config [options]",
			Validator:   ConfigCommandValidator,
			Handler:     ConfigCommandHandler,
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

func ConfigCommandValidator(cmd *clif.Command) (err error) {
	return nil
}

func ConfigCommandHandler(cmd *clif.Command) (err error) {
	website := htxt.Website{URI: &href.Href{URL: url.URL{Path: cmd.Options.Get("config")}}}
	if err = website.Load(true, nil, nil); err != nil {
		return err
	}
	fmt.Printf("%v\n", website)
	return nil
}
