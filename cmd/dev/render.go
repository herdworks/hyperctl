package dev

import (
	"fmt"

	cli "github.com/herdworks/cli-go"
	hypertemplates "github.com/herdworks/hypertemplates-go"
)

func RenderCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "render",
			Description: "HyperTemplate rendering utilities.",
			Usage:       "hyperctl dev render [options]",
			Validator:   RenderCommandValidator,
			Handler:     RenderCommandHandler,
		},
		Options: cli.CommandOptions{
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

func RenderCommandValidator(cmd *cli.Command) (err error) {
	return nil
}

func RenderCommandHandler(cmd *cli.Command) (err error) {
	html, err := hypertemplates.Render(
		cmd.Options.Get("data"),
		cmd.Options.Get("layout"),
		cmd.Options.Get("theme"),
		cmd.Options.Get("untidy") == "true",
		&hypertemplates.DefaultHyperMarkConfig,
		nil,
	)
	if err != nil {
		return err
	}
	fmt.Println(string(html))
	return nil
}
