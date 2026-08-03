package sync

import (
	"testing"

	cli "github.com/herdworks/cli-go"
)

func TestCommandRegistersCheckoutWorkflow(t *testing.T) {
	command := Command()
	want := []string{"download", "init", "pull", "push", "reset", "status"}
	if len(command.SubCommands) != len(want) {
		t.Fatalf("subcommands = %d, want %d", len(command.SubCommands), len(want))
	}
	for index, name := range want {
		if command.SubCommands[index].Config.Name != name {
			t.Errorf("subcommand %d = %q, want %q", index, command.SubCommands[index].Config.Name, name)
		}
	}
}

func TestConfigurationCommandsDefaultToSiteYAML(t *testing.T) {
	for _, command := range []*cli.Command{initCommand(), pullCommand(), pushCommand(), resetCommand(), statusCommand()} {
		option := command.Options.Lookup("config")
		if option == nil || option.Default != "site.yaml" || option.EnvName != "HYPER_CONFIG" {
			t.Errorf("%s config default = %#v, want site.yaml", command.Config.Name, option)
		}
	}
}

func TestConfigurationCommandsHonorHyperConfig(t *testing.T) {
	t.Setenv("HYPER_CONFIG", "alternate.yaml")
	for _, command := range []*cli.Command{initCommand(), pullCommand(), pushCommand(), resetCommand(), statusCommand()} {
		if got := command.Options.Get("config"); got != "alternate.yaml" {
			t.Errorf("%s config = %q, want alternate.yaml", command.Config.Name, got)
		}
	}
}
