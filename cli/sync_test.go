package cli

import "testing"

func TestSyncCommandRegistersCheckoutWorkflow(t *testing.T) {
	command := SyncCommand()
	if command.Spec.Key != "sync" {
		t.Fatalf("sync key = %q", command.Spec.Key)
	}
	if command.Meta.EnvPrefix != "" {
		t.Fatalf("sync environment prefix = %q, want inherited", command.Meta.EnvPrefix)
	}
	want := []string{"download", "init", "pull", "push", "reset", "status"}
	if len(command.Commands) != len(want) {
		t.Fatalf("sync commands = %d, want %d", len(command.Commands), len(want))
	}
	for _, key := range want {
		if _, ok := command.Commands[key]; !ok {
			t.Errorf("sync command %q is not registered", key)
		}
	}
}

func TestSyncConfigurationOptionsInheritHyperPrefix(t *testing.T) {
	command := SyncCommand()
	for _, key := range []string{"init", "pull", "push", "reset", "status"} {
		config := command.Commands[key].Spec.Options["config"]
		if config == nil {
			t.Errorf("%s config option is not registered", key)
			continue
		}
		if config.Meta.EnvName != "{{ .EnvPrefix }}CONFIG" {
			t.Errorf("%s config environment template = %q", key, config.Meta.EnvName)
		}
	}
}
