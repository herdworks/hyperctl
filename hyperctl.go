package main

import (
	"fmt"

	cli "github.com/herdworks/cli-go"
	cmd "github.com/herdworks/hyperctl/cmd"
	hypertexting "github.com/herdworks/hypertexting-go"
)

var HyperCTL = &cli.CLI{
	Name:        "hyperctl",
	Version:     fmt.Sprintf("%s (%s)", hypertexting.Version(), hypertexting.ReleaseDate()),
	Description: "hyperctl is a lightweight static site generator based on HyperTemplates.",
	Commands: []*cli.Command{
		cmd.AssetCommand(),
		cmd.BuildCommand(),
		cmd.ContentTypeCommand(),
		cmd.DataCommand(),
		cmd.DeployCommand(),
		cmd.GenerateCommand(),
		cmd.GraphCommand(),
		cmd.MimetypeCommand(),
		cmd.NewCommand(),
		cmd.PageCommand(),
		cmd.PublishCommand(),
		cmd.RenderCommand(),
		cmd.ServerCommand(),
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
