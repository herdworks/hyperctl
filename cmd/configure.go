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

// TODO: refactor to support reading/writing from/to yaml or json
func ConfigureCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "configure",
			Description: "Manage HyperTemplates website configuration.",
			Usage:       "hyperctl configure [options]",
			Validator:   ConfigureCommandValidator,
			Handler:     ConfigureCommandHandler,
		},
		Options: cli.CommandOptions{
			{
				Name:        "base_url",
				ShortName:   "u",
				EnvName:     "HYPER_BASEURL",
				Description: "Base URL to use for the website.",
				Secret:      false,
			},
			{
				Name:        "title",
				ShortName:   "t",
				EnvName:     "HYPER_TITLE",
				Description: "Title of the website.",
				Secret:      false,
			},
			{
				Name:        "description",
				ShortName:   "d",
				EnvName:     "HYPER_DESCRIPTION",
				Description: "Description of the website.",
				Secret:      false,
			},
			{
				Name:        "theme",
				EnvName:     "HYPER_THEME",
				Description: "Name of the HyperTemplates theme.",
			},
			{
				Name:        "metadata",
				ShortName:   "m",
				Description: "Metadata key=value pairs (comma-separated).",
			},
		},
	}
}

func ConfigureCommandValidator(cmd *cli.Command) (err error) {
	if len(os.Args[2:]) == 0 {
		return fmt.Errorf("no configuration parameters provided, nothing to do")
	}
	project := cmd.Options.Get("project") // TODO: the --project flag has been removed
	config := filepath.Join(project, "site.json")
	if _, err := os.Stat(config); err != nil {
		return fmt.Errorf("error loading website configuration at %s; %s", config, err)
	}
	return nil
}

func ConfigureCommandHandler(cmd *cli.Command) (err error) {
	project := cmd.Options.Get("project") // TODO: the --project flag has been removed
	configfile := filepath.Join(project, "site.json")
	uri, _ := url.Parse(configfile)
	website := &hypertexting.Website{URI: uri}
	var jsonBytes []byte
	if jsonBytes, err = os.ReadFile(configfile); err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, website)
	if err != nil {
		return err
	}

	if cmd.Options.Get("base_url") != "" {
		fmt.Printf("updating site.base_url to \"%s\"\n", cmd.Options.Get("base_url"))
		website.BaseURL = cmd.Options.Get("base_url")
	}
	if cmd.Options.Get("title") != "" {
		fmt.Printf("updating site.title to \"%s\"\n", cmd.Options.Get("title"))
		website.Title = cmd.Options.Get("title")
	}
	if cmd.Options.Get("description") != "" {
		fmt.Printf("updating site.description to \"%s\"\n", cmd.Options.Get("description"))
		website.Description = cmd.Options.Get("description")
	}
	// if cmd.Options.Get("metadata") != "" {
	// 	properties := strings.Split(cmd.Options.Get("metadata"), ",")
	// 	for _, prop := range properties {
	// 		kv := strings.Split(prop, "=")
	// 		if len(kv) == 2 {
	// 			fmt.Printf("updating site.metadata.%s to \"%s\"\n", kv[0], kv[1])
	// 			website.Metadata[kv[0]] = kv[1]
	// 		}
	// 	}
	// }

	// write updated configuration
	jsonBytes, err = json.MarshalIndent(website, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configfile, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
