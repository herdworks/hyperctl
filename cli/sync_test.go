package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/herdworks/clic"
	"github.com/herdworks/htxt"
)

func TestSyncCommandRegistersCheckoutWorkflow(t *testing.T) {
	command := SyncCommand()
	if command.Spec.Key != "sync" {
		t.Fatalf("sync key = %q", command.Spec.Key)
	}
	if command.Meta.EnvPrefix != "" {
		t.Fatalf("sync environment prefix = %q, want inherited", command.Meta.EnvPrefix)
	}
	want := []string{"download", "init", "pull", "push", "remote", "reset", "status"}
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
	for _, key := range []string{"init", "pull", "push", "remote", "reset", "status"} {
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

func TestSyncInitDefaultsToSiteConfiguration(t *testing.T) {
	config := SyncCommand().Commands["init"].Spec.Options["config"]
	if config == nil || config.Meta.Default != htxt.DefaultSyncConfigPath {
		t.Fatalf("init config default = %#v", config)
	}
}

func TestSyncStatusUsesWebsiteConfiguration(t *testing.T) {
	directory := t.TempDir()
	site := `base_url: https://example.com
title: Example
description: Example website
byline:
  name: Example
sync:
  remote: https://example.com
`
	if err := os.WriteFile(filepath.Join(directory, "site.yaml"), []byte(site), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "theme.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	result := CLI.Run(clic.Execution{
		Context:          context.Background(),
		WorkingDirectory: directory,
		Stdout:           &output,
		LookupEnv:        func(string) (string, bool) { return "", false },
	}, clic.Arguments{"sync", "status"})
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if !bytes.Contains(output.Bytes(), []byte("added site.yaml")) {
		t.Fatalf("status output = %q", output.String())
	}
}
