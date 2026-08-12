package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/herdworks/clif"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

func InitCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "init",
			Description: "Init a HyperTemplates website project.",
			Usage:       "hyperctl init [options] [content]",
			Validator:   InitCommandValidator,
			Handler:     InitCommandHandler,
		},
		Options: clif.CommandOptions{
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

func InitCommandValidator(cmd *clif.Command) (err error) {
	project := cmd.Options.Get("project") // TODO: the --project flag has been removed
	config := filepath.Join(project, "site.json")
	if _, err := os.Stat(config); err == nil {
		return fmt.Errorf("website configuration already exists at %s", config)
	}
	return nil
}

func InitCommandHandler(cmd *clif.Command) (err error) {
	project := cmd.Options.Get("project") // TODO: the --project flag has been removed
	config := filepath.Join(project, "site.json")
	uri, _ := href.Parse(config)
	website := htxt.Website{
		URI:         uri,
		BaseURL:     cmd.Options.Get("base_url"),
		Title:       cmd.Options.Get("title"),
		Description: cmd.Options.Get("description"),
		Links:       []htxt.HTMLLink{},
		Metadata:    []htxt.HTMLMeta{},
		Config:      htxt.WebsiteConfig{},
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
