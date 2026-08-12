package cms

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/herdworks/clif"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

// hyperctl content-type
func ContentTypeCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "content-type",
			Description: "Manage content-types.",
			Usage:       "hyperctl cms content-type [command] [options]",
			Handler:     ContentTypeCommandHandler,
		},
		Options: clif.CommandOptions{},
		SubCommands: clif.Commands{
			ContentTypeListCommand(),
			ContentTypeInspectCommand(),
		},
	}
}
func ContentTypeCommandHandler(cmd *clif.Command) (err error) {
	cmd.Help(nil)
	return nil
}

// hyperctl content-type ls
func ContentTypeListCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "ls",
			Description: "List content types.",
			Usage:       "hyperctl content-type ls [options]",
			Handler:     ContentTypeListCommandHandler,
		},
		Options: clif.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Default:     "site.yaml",
				Required:    true,
			},
		},
	}
}

func ContentTypeListCommandHandler(cmd *clif.Command) (err error) {
	website := htxt.Website{URI: &href.Href{URL: url.URL{Path: cmd.Options.Get("config")}}}
	if err = website.Load(false, nil, nil); err != nil {
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
func ContentTypeInspectCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "inspect",
			Description: "Inspect content types.",
			Usage:       "hyperctl content-type inspect [options] <type>",
			Validator:   ContentTypeInspectValidator,
			Handler:     ContentTypeInspectCommandHandler,
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
		},
	}
}
func ContentTypeInspectValidator(cmd *clif.Command) (err error) {
	requiredArgs := 1
	if len(cmd.Args) != requiredArgs {
		return fmt.Errorf("expected %v argument(s), got %v (%s)", requiredArgs, len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}
func ContentTypeInspectCommandHandler(cmd *clif.Command) (err error) {
	var jsonBytes []byte
	var typeRef string = cmd.Args[0]
	website := htxt.Website{URI: &href.Href{URL: url.URL{Path: cmd.Options.Get("config")}}}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	contentType := website.GetContentTypeByName(typeRef, false, nil)
	if contentType == nil {
		return fmt.Errorf("content-type \"%s\" not found", typeRef)
	}
	if jsonBytes, err = json.Marshal(contentType.Spec.TemplateData); err != nil {
		return err
	}
	fmt.Println(string(jsonBytes))
	return nil
}
