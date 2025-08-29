package cli

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"

	hypergraph "github.com/herdworks/hypergraph-go/lib"
	hypertemplates "github.com/herdworks/hypertemplates-go/lib"
	hypertexting "github.com/herdworks/hypertexting-go/lib"

	opml "github.com/herdworks/hypertemplates-go/lib/opml"
)

// hyperctl data
func DataCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "data",
			Description: "Manage global data.",
			Usage:       "hyperctl data [command] [options]",
			Handler:     DataCommandHandler,
		},
		Options: CommandOptions{},
		SubCommands: Commands{
			DataListCommand(),
			DataInspectCommand(),
			DataTransformCommand(),
		},
	}
}
func DataCommandHandler(cmd *Command) (err error) {
	cmd.Help(nil)
	return nil
}

// hyperctl data ls
func DataListCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "ls",
			Description: "List global data namespaces.",
			Usage:       "hyperctl data ls [options]",
			Handler:     DataListCommandHandler,
		},
		Options: CommandOptions{
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

func DataListCommandHandler(cmd *Command) (err error) {
	var websiteURI *url.URL
	var website *hypertexting.Website
	var environment hypertexting.Environment = hypertexting.Environment{}
	if websiteURI, err = url.Parse(cmd.Options.Get("config")); err != nil {
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
	if website.GlobalData.IsEmpty() {
		fmt.Println("No namespace data configured.")
		return nil
	}
	colwidth := cmd.ColumnWidth(website.GlobalData, "Namespace", 1)
	fmt.Printf("%-*s %s\n", colwidth, "Namespace", "Size")
	fmt.Printf("%-*s %s\n", colwidth, "---------", "----")
	for _, namespace := range website.GlobalData.Keys() {
		if bytes, err := website.GlobalData.Inspect([]string{namespace}, "json"); err == nil {
			fmt.Printf("%-*s %v\n", colwidth, namespace, len(bytes))
		} else {
			return fmt.Errorf("error inspecting namespace \"%s\": %v", namespace, err)
		}
	}
	return nil
}

// hyperctl data inspect
func DataInspectCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "inspect",
			Description: "Inspect global data namespace.",
			Usage:       "hyperctl data inspect <namespace> [options]",
			Validator:   DataInspectValidator,
			Handler:     DataInspectCommandHandler,
		},
		Options: CommandOptions{
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
func DataInspectValidator(cmd *Command) (err error) {
	requiredArgs := 1
	if len(cmd.Args) != requiredArgs {
		return fmt.Errorf("expected %v argument(s), got %v (%s)", requiredArgs, len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}
func DataInspectCommandHandler(cmd *Command) (err error) {
	var websiteURI *url.URL
	var website *hypertexting.Website
	var namespace string = cmd.Args[0]
	var namespaceBytes []byte
	if websiteURI, err = url.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if err = website.Load(false, nil); err != nil {
		return err
	}
	if website.GlobalData, err = hypertexting.NamespacedTemplateData(website, website.Config.Markdown, false); err != nil {
		return err
	}
	if namespaceBytes, err = website.GlobalData.Inspect([]string{namespace}, "json"); err != nil {
		return err
	}
	fmt.Println(string(namespaceBytes))
	return nil
}

// hyperctl data transform
func DataTransformCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "transform",
			Description: "Transform global data namespaces.",
			Usage:       "hyperctl data transform <namespace> [options]",
			Validator:   DataTransformValidator,
			Handler:     DataTransformCommandHandler,
		},
		Options: CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Required:    true,
				Secret:      false,
			},
			{
				Name:        "opml-enrichment",
				ShortName:   "o",
				Description: "Enrich OPML data with additional properties.",
				Boolean:     true,
			},
		},
	}
}
func DataTransformValidator(cmd *Command) (err error) {
	requiredArgs := 1
	if len(cmd.Args) != requiredArgs {
		return fmt.Errorf("expected %v argument(s), got %v (%s)", requiredArgs, len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}
func DataTransformCommandHandler(cmd *Command) (err error) {
	var websiteURI *url.URL
	var website *hypertexting.Website
	var namespace string = cmd.Args[0]
	if websiteURI, err = url.Parse(cmd.Options.Get("config")); err != nil {
		return err
	}
	website = &hypertexting.Website{URI: websiteURI}
	if err = website.Load(false, nil); err != nil {
		return err
	}
	if website.GlobalData, err = hypertexting.NamespacedTemplateData(website, website.Config.Markdown, false); err != nil {
		return err
	}
	if cmd.Options.Get("opml-enrichment") == "true" {
		if err = DataTransformOPMLEnrichmentHandler(website, namespace); err != nil {
			return err
		}
	}
	return nil
}
func DataTransformOPMLEnrichmentHandler(website *hypertexting.Website, namespace string) (err error) {
	var bytes []byte
	var opml opml.OPML
	var opmlBytes []byte
	if bytes, err = website.GlobalData.Inspect([]string{namespace}, "json"); err != nil {
		return err
	}
	if err = json.Unmarshal(bytes, &opml); err != nil {
		return err
	}
	for index, _ := range opml.Body.Outlines {
		TransformOutline(opml.Body.Outlines[index])
	}
	if opmlBytes, err = xml.MarshalIndent(opml, "", "    "); err != nil {
		return err
	}
	fmt.Println(string(opmlBytes))
	return nil
}

func TransformOutline(outline *opml.Outline) {
	document := hypergraph.Document{Href: outline.HTMLURL}
	m := hypertemplates.TemplateData(document.Metadata(true))
	outline.Title = first(outline.Title, m.Get("author.name"), m.Get("hypergraph.title")) // use existing title if available
	outline.Text = first(m.Get("author.username"), outline.Text)                          // always set text to username (if discovered)
	outline.FaviconURL = first(m.Get("hypergraph.favicon"), outline.FaviconURL)           // always update faviconURL (if discovered)
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
