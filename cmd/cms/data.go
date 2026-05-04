package cms

import (
	"encoding/json"
	"fmt"
	"strings"

	cli "github.com/herdworks/cli-go"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
)

// hyperctl data
func DataCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "data",
			Description: "Manage global template data.",
			Usage:       "hyperctl cms data [command] [options]",
			Handler:     DataCommandHandler,
		},
		Options: cli.CommandOptions{},
		SubCommands: cli.Commands{
			DataListCommand(),
			DataInspectCommand(),
		},
	}
}
func DataCommandHandler(cmd *cli.Command) (err error) {
	cmd.Help(nil)
	return nil
}

// hyperctl data ls
func DataListCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "ls",
			Description: "List global data namespaces.",
			Usage:       "hyperctl data ls [options]",
			Handler:     DataListCommandHandler,
		},
		Options: cli.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Required:    true,
			},
		},
	}
}

func DataListCommandHandler(cmd *cli.Command) (err error) {
	var websiteURI *href.Href
	var website *hypertexting.Website
	var environment hypertexting.Environment = hypertexting.Environment{}
	if websiteURI, err = href.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if err = website.Load(false, nil); err != nil {
		return err
	}
	if err = environment.Load(); err != nil {
		return err
	}
	if website.GlobalData, err = hypertexting.NamespacedTemplateData(website, website.Config.Markdown, false); err != nil {
		return err
	}
	if environment.TemplateData != nil {
		website.GlobalData.Set("env", environment.TemplateData)
	} // add environment data to global data
	if len(website.GlobalData) == 0 {
		fmt.Println("No namespace data configured.")
		return nil
	}
	colwidth := cmd.ColumnWidth(hypertexting.NamespacedData(website.GlobalData), "Namespace", 1)
	fmt.Printf("%-*s %s\n", colwidth, "Namespace", "Size")
	fmt.Printf("%-*s %s\n", colwidth, "---------", "----")
	for _, namespace := range website.GlobalData.Keys("") {
		bytes, err := json.Marshal(website.GlobalData.Get(namespace))
		if err != nil {
			return fmt.Errorf("error inspecting namespace %q: %v", namespace, err)
		}
		fmt.Printf("%-*s %v\n", colwidth, namespace, len(bytes))
	}
	return nil
}

// hyperctl data inspect
func DataInspectCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "inspect",
			Description: "Inspect global data namespace.",
			Usage:       "hyperctl data inspect <namespace> [options]",
			Validator:   DataInspectValidator,
			Handler:     DataInspectCommandHandler,
		},
		Options: cli.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Required:    true,
				Secret:      false,
			},
		},
	}
}
func DataInspectValidator(cmd *cli.Command) (err error) {
	requiredArgs := 1
	if len(cmd.Args) != requiredArgs {
		return fmt.Errorf("expected %v argument(s), got %v (%s)", requiredArgs, len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}
func DataInspectCommandHandler(cmd *cli.Command) (err error) {
	var websiteURI *href.Href
	var website *hypertexting.Website
	var namespace string = cmd.Args[0]
	var namespaceBytes []byte
	if websiteURI, err = href.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if err = website.Load(false, nil); err != nil {
		return err
	}
	if website.GlobalData, err = hypertexting.NamespacedTemplateData(website, website.Config.Markdown, false); err != nil {
		return err
	}
	if namespaceBytes, err = json.Marshal(website.GlobalData.Get(namespace)); err != nil {
		return err
	}
	fmt.Println(string(namespaceBytes))
	return nil
}

func first(incoming ...any) string {
	var s string
	var ok bool
	for _, i := range incoming {
		if i == nil {
			continue
		} // skip if nil
		if s, ok = i.(string); !ok {
			continue
		} // skip if not a string
		if len(strings.TrimSpace(s)) > 0 {
			return s
		}
	}
	return ""
}
