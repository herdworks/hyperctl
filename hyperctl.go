package main

import (
	"fmt"

	cli "github.com/herdworks/cli-go"
	agent "github.com/herdworks/hyperctl/cmd/agent"
	build "github.com/herdworks/hyperctl/cmd/build"
	cms "github.com/herdworks/hyperctl/cmd/cms"
	deploy "github.com/herdworks/hyperctl/cmd/deploy"
	dev "github.com/herdworks/hyperctl/cmd/dev"
	sync "github.com/herdworks/hyperctl/cmd/sync"
	theme "github.com/herdworks/hyperctl/cmd/theme"
	version "github.com/herdworks/hyperctl/cmd/version"
)

var HyperCTL = &cli.CLI{
	Name:        "hyperctl",
	Version:     fmt.Sprintf("%s (%s)", version.CLIVersion, version.CLIReleaseDate),
	Description: "The HyperTexting CLI.",
	Commands: []*cli.Command{
		agent.Command(),
		cms.Command(),
		theme.Command(),
		build.Command(),
		deploy.Command(),
		sync.Command(),
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
