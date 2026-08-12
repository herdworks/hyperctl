package dev

import (
	"fmt"

	"github.com/herdworks/clif"
	"github.com/herdworks/httt"
)

func RenderCommand() *clif.Command {
	return &clif.Command{
		Config: clif.CommandConfig{
			Name:        "render",
			Description: "HyperTemplate rendering utilities.",
			Usage:       "hyperctl dev render [options]",
			Validator:   RenderCommandValidator,
			Handler:     RenderCommandHandler,
		},
		Options: clif.CommandOptions{
			{
				Name:        "data",
				ShortName:   "d",
				EnvName:     "HYPER_DATA",
				Description: "Path to a HyperTemplates data file (json, yaml, or markdown).",
				Required:    true,
				Secret:      false,
			},
			{
				Name:        "layout",
				ShortName:   "l",
				Description: "Path to a HyperTemplates layout file (html).",
				Required:    false,
				Secret:      false,
			},
			{
				Name:        "theme",
				ShortName:   "t",
				Description: "Name of HyperTemplates theme.",
				Required:    false,
				Secret:      false,
			},
			{
				Name:        "untidy",
				ShortName:   "u",
				Description: "Disable tidy mode.",
				Boolean:     true,
			},
		},
	}
}

func RenderCommandValidator(cmd *clif.Command) (err error) {
	return nil
}

func RenderCommandHandler(cmd *clif.Command) (err error) {
	html, err := httt.Render(
		cmd.Options.Get("data"),
		cmd.Options.Get("layout"),
		cmd.Options.Get("theme"),
		cmd.Options.Get("untidy") == "true",
		&httt.DefaultHyperMarkConfig,
		nil,
	)
	if err != nil {
		return err
	}
	fmt.Println(string(html))
	return nil
}
