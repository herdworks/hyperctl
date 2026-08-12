package main

import (
	"fmt"

	"github.com/herdworks/clif"
	"github.com/herdworks/hyperctl/cmd/agent"
	"github.com/herdworks/hyperctl/cmd/build"
	"github.com/herdworks/hyperctl/cmd/cms"
	"github.com/herdworks/hyperctl/cmd/deploy"
	"github.com/herdworks/hyperctl/cmd/dev"
	"github.com/herdworks/hyperctl/cmd/sync"
	"github.com/herdworks/hyperctl/cmd/theme"
	"github.com/herdworks/hyperctl/cmd/version"
)

var HyperCTL = &clif.CLI{
	Name:        "hyperctl",
	Version:     fmt.Sprintf("%s (%s)", version.CLIVersion, version.CLIReleaseDate),
	Description: "The HyperTexting CLI.",
	Commands: []*clif.Command{
		agent.Command(),
		cms.Command(),
		theme.Command(),
		build.Command(),
		deploy.Command(),
		sync.Command(),
		dev.Command(),
	},
	Options: clif.CommandOptions{
		{
			Name:        "help",
			ShortName:   "h",
			Description: "Display help information.",
			Boolean:     true,
			Secret:      true,
			Global:      true,
		},
	},
}

func main() {
	HyperCTL.Execute()
}
