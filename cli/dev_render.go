package cli

import (
	"fmt"

	"github.com/herdworks/clic"
	"github.com/herdworks/httt"
)

func DevRenderCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "render",
			Description: "HyperTemplate rendering utilities.",
			Usage: []string{
				"hyperctl dev render [options]",
			},
		},
		Spec: clic.CommandSpec{
			Key: "render", Help: clic.DefaultHelp,
			Options: clic.Options{
				"data":   {Meta: clic.OptionMeta{Name: "data", ShortName: "d", EnvName: "HYPER_DATA", Description: "Path to a HyperTemplates data file (json, yaml, or markdown)."}, Spec: clic.OptionSpec{Key: "data", Kind: clic.StringOption, Required: true}},
				"layout": {Meta: clic.OptionMeta{Name: "layout", ShortName: "l", Description: "Path to a HyperTemplates layout file (html)."}, Spec: clic.OptionSpec{Key: "layout", Kind: clic.StringOption}},
				"theme":  {Meta: clic.OptionMeta{Name: "theme", ShortName: "t", Description: "Name of HyperTemplates theme."}, Spec: clic.OptionSpec{Key: "theme", Kind: clic.StringOption}},
				"untidy": {Meta: clic.OptionMeta{Name: "untidy", ShortName: "u", Description: "Disable tidy mode.", Default: false}, Spec: clic.OptionSpec{Key: "untidy", Kind: clic.BoolOption}},
			},
			Handler: renderHandler,
		},
	}
}

func renderHandler(execution *clic.Execution) error {
	options := execution.Options
	html, err := httt.Render(options["data"].StringValue(), options["layout"].StringValue(), options["theme"].StringValue(), options["untidy"].BoolValue(), &httt.DefaultHyperMarkConfig, nil)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(execution.Stdout, string(html))
	return err
}
