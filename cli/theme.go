package cli

import (
	"github.com/herdworks/clic"
	htttcli "github.com/herdworks/httt/cli"
)

// ThemeCommand embeds the httt CLI at `hyperctl theme`
func ThemeCommand() clic.Command {
	theme := htttcli.Command()
	theme.Spec.Key = "theme"
	theme.Meta.Aliases = []string{"httt", "tmpl"} // new httt + tmpl aliases
	theme.Meta.Description = "Theme management tools"
	theme.Meta.Order = 2
	theme.Meta.EnvPrefix = "" // inherit the root hyperctl prefix
	return theme
}
