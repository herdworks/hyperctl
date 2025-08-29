package main

import (
	"fmt"

	"github.com/herdworks/hyperctl/cli"
	hypertexting "github.com/herdworks/hypertexting-go/lib"
)

var HyperCTL = &cli.CLI{
	Name:        "hyperctl",
	Version:     fmt.Sprintf("%s (%s)", hypertexting.Version(), hypertexting.ReleaseDate()),
	Description: "hyperctl is a lightweight static site generator based on HyperTemplates.",
	Commands: []*cli.Command{
		cli.AssetCommand(),
		cli.BuildCommand(),
		cli.ContentTypeCommand(),
		cli.DataCommand(),
		cli.DeployCommand(),
		cli.GenerateCommand(),
		cli.GraphCommand(),
		cli.MimetypeCommand(),
		cli.NewCommand(),
		cli.PageCommand(),
		cli.PublishCommand(),
		cli.RenderCommand(),
		cli.ServerCommand(),
	},
	Options: cli.CommandOptions{
		{
			Name:        "help",
			ShortName:   "h",
			Description: "Display help information.",
			Boolean:     true,
			Secret:      true, // don't display a default value for this option
			Global:      true, // pass this option to all subcommands
		},
	},
}

func main() {
	HyperCTL.Execute()
}
