package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/herdworks/clic"
)

func TestThemeCommandRemixesHTTTCommands(t *testing.T) {
	command := ThemeCommand()
	for _, name := range []string{"init", "ls", "install", "package", "migrate", "render"} {
		child, ok := command.Commands[name]
		if !ok {
			t.Fatalf("theme command missing %q", name)
		}
		if len(child.Meta.Usage) == 0 || !strings.Contains(child.Meta.Usage[0], "{{ .Command }}") {
			t.Fatalf("%s usage = %#v", name, child.Meta.Usage)
		}
		if config := child.Spec.Options["config"]; config != nil && config.Meta.EnvName != "{{ .EnvPrefix }}CONFIG" {
			t.Fatalf("%s config env = %q", name, config.Meta.EnvName)
		}
	}
}

func TestThemeCommandRendersMountedUsage(t *testing.T) {
	var output bytes.Buffer
	result := CLI.Run(
		clic.Execution{
			Context: context.Background(),
			Stdout:  &output,
		},
		clic.Arguments{
			"theme",
			"render",
			"--help",
		},
	)
	if result.Error != nil {
		t.Fatalf("CLI.Run(theme render --help) = %#v", result)
	}
	if !strings.Contains(output.String(), "hyperctl theme render --theme theme.json") {
		t.Fatalf("render help = %q", output.String())
	}
}

func TestHyperctlKeepsDevAsAWebsiteAlias(t *testing.T) {
	command, ok := CLI.SubCommand("dev")
	if !ok || command.Name() != "website" {
		t.Fatalf("dev alias = %#v, %t", command, ok)
	}
}
