package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	cli "github.com/herdworks/cli-go"
	hypertexting "github.com/herdworks/hypertexting-go"
)

func InitCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "init",
			Description: "Init a HyperTemplates website project.",
			Usage:       "hyperctl init [options] [content]",
			Validator:   InitCommandValidator,
			Handler:     InitCommandHandler,
		},
		Options: cli.CommandOptions{
			{
				Name:        "base_url",
				ShortName:   "u",
				EnvName:     "HYPER_BASEURL",
				Description: "Base URL to use for the website.",
				Required:    true,
				Secret:      false,
			},
			{
				Name:        "title",
				ShortName:   "T",
				EnvName:     "HYPER_TITLE",
				Description: "Title of the website.",
				Required:    true,
				Secret:      false,
			},
			{
				Name:        "description",
				ShortName:   "d",
				EnvName:     "HYPER_DESCRIPTION",
				Description: "Description of the website.",
				Required:    true,
				Secret:      false,
			},
		},
	}
}

func InitCommandValidator(cmd *cli.Command) (err error) {
	project := cmd.Options.Get("project") // TODO: the --project flag has been removed
	config := filepath.Join(project, "site.json")
	if _, err := os.Stat(config); err == nil {
		return fmt.Errorf("website configuration already exists at %s", config)
	}
	return nil
}

func InitCommandHandler(cmd *cli.Command) (err error) {
	project := cmd.Options.Get("project") // TODO: the --project flag has been removed
	config := filepath.Join(project, "site.json")
	uri, _ := url.Parse(config)
	website := hypertexting.Website{
		URI:         uri,
		BaseURL:     cmd.Options.Get("base_url"),
		Title:       cmd.Options.Get("title"),
		Description: cmd.Options.Get("description"),
		Links:       []hypertexting.HTMLLink{},
		Metadata:    []hypertexting.HTMLMeta{},
		Config:      hypertexting.WebsiteConfig{},
	}
	jsonBytes, err := json.Marshal(website)
	if err != nil {
		return err
	}
	err = os.WriteFile(config, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
