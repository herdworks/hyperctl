package main

import (
	"fmt"

	cli "github.com/herdworks/cli-go"
	build "github.com/herdworks/hyperctl/cmd/build"
	cms "github.com/herdworks/hyperctl/cmd/cms"
	deploy "github.com/herdworks/hyperctl/cmd/deploy"
	dev "github.com/herdworks/hyperctl/cmd/dev"
	theme "github.com/herdworks/hyperctl/cmd/theme"

	hypertexting "github.com/herdworks/hypertexting-go"
)

var HyperCTL = &cli.CLI{
	Name:        "hyperctl",
	Version:     fmt.Sprintf("%s (%s)", hypertexting.Version(), hypertexting.ReleaseDate()),
	Description: "hyperctl is a simple static site generator.",
	Commands: []*cli.Command{
		cms.Command(),
		theme.Command(),
		build.Command(),
		deploy.Command(),
		dev.Command(),
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
