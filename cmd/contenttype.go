package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	cli "github.com/herdworks/cli-go"
	hypertexting "github.com/herdworks/hypertexting-go"
)

// hyperctl content-type
func ContentTypeCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "content-type",
			Description: "Inspect HyperTemplate theme content types.",
			Usage:       "hyperctl content-type [command] [options]",
			Handler:     ContentTypeCommandHandler,
		},
		Options: cli.CommandOptions{},
		SubCommands: cli.Commands{
			ContentTypeListCommand(),
			ContentTypeInspectCommand(),
		},
	}
}
func ContentTypeCommandHandler(cmd *cli.Command) (err error) {
	cmd.Help(nil)
	return nil
}

// hyperctl content-type ls
func ContentTypeListCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "ls",
			Description: "List content types.",
			Usage:       "hyperctl content-type ls [options]",
			Handler:     ContentTypeListCommandHandler,
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

func ContentTypeListCommandHandler(cmd *cli.Command) (err error) {
	website := hypertexting.Website{URI: &url.URL{Path: cmd.Options.Get("config")}}
	if err = website.Load(false, nil); err != nil {
		return err
	}

	supportedtypes, err := website.ContentTypes()
	if err != nil {
		return err
	}
	if len(supportedtypes) == 0 {
		fmt.Println("No content types configured.")
		return nil
	}
	colwidth := cmd.ColumnWidth(supportedtypes, "Name", 1)
	fmt.Printf("%-*s %s\n", colwidth, "Name", "Description")
	fmt.Printf("%-*s %s\n", colwidth, "----", "-----------")
	for _, contenttype := range supportedtypes {
		fmt.Printf("%-*s %s\n", colwidth, contenttype.Metadata.Name, contenttype.Metadata.Description)
	}
	return nil
}

// hyperctl content-type inspect
func ContentTypeInspectCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "inspect",
			Description: "Inspect content types.",
			Usage:       "hyperctl content-type inspect [options] <type>",
			Validator:   ContentTypeInspectValidator,
			Handler:     ContentTypeInspectCommandHandler,
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
func ContentTypeInspectValidator(cmd *cli.Command) (err error) {
	requiredArgs := 1
	if len(cmd.Args) != requiredArgs {
		return fmt.Errorf("expected %v argument(s), got %v (%s)", requiredArgs, len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}
func ContentTypeInspectCommandHandler(cmd *cli.Command) (err error) {
	var jsonBytes []byte
	var typeRef string = cmd.Args[0]
	website := hypertexting.Website{URI: &url.URL{Path: cmd.Options.Get("config")}}
	if err = website.Load(false, nil); err != nil {
		return err
	}
	contentType := website.GetContentTypeByName(typeRef, false)
	if contentType == nil {
		return fmt.Errorf("content-type \"%s\" not found", typeRef)
	}
	if jsonBytes, err = json.Marshal(contentType.Spec.TemplateData); err != nil {
		return err
	}
	fmt.Println(string(jsonBytes))
	return nil
}
