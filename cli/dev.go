package cli

import (
	"github.com/herdworks/clic"
)

// DevCommand provides development tools.
func DevCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "dev",
			Description: "Developer tools",
			Usage: []string{
				"hyperctl dev [command]",
			},
		},
		Spec: clic.CommandSpec{
			Key:  "dev",
			Help: clic.DefaultHelp,
		},
		Commands: clic.MustCommandsFrom(DevGraphCommand, DevMimetypeCommand, DevRenderCommand, DevServerCommand),
	}
}
