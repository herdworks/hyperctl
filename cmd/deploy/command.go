package deploy

import (
	cli "github.com/herdworks/cli-go"
)

func Command() *cli.Command {
	return DeployCommand()
}
