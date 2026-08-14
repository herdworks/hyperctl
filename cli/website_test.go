package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/herdworks/clic"
)

func TestWebsiteCommandMountsHTXTWorkflow(t *testing.T) {
	command := WebsiteCommand()
	if command.Meta.Name != "website" || command.Spec.Key != "htxt" {
		t.Fatalf("website command = %#v", command)
	}
	if command.Meta.EnvPrefix != "" {
		t.Fatalf("website environment prefix = %q", command.Meta.EnvPrefix)
	}
	for _, key := range []string{"init", "build", "server", "data", "page", "asset", "content-type"} {
		if _, ok := command.SubCommand(key); !ok {
			t.Errorf("website command %q is not registered", key)
		}
	}
	var output bytes.Buffer
	result := CLI.Run(
		clic.Execution{
			Context: context.Background(),
			Stdout:  &output,
		},
		clic.Arguments{
			"website",
			"page",
			"new",
			"--help",
		},
	)
	if result.Error != nil {
		t.Fatalf("CLI.Run(website page new --help) = %#v", result)
	}
	if strings.Contains(output.String(), "htxt") || !strings.Contains(output.String(), "hyperctl website page new") {
		t.Errorf("mounted page-new help = %q", output.String())
	}
	if !strings.Contains(output.String(), "HYPER_CONFIG") {
		t.Errorf("mounted page-new environment = %q", output.String())
	}
}

func TestHyperctlKeepsCMSAsAnAliasAndBuildAsALink(t *testing.T) {
	command, ok := CLI.SubCommand("cms")
	if !ok || command.Name() != "website" {
		t.Fatalf("cms alias = %#v, %t", command, ok)
	}
	if _, ok := CLI.SubCommand("build"); ok {
		t.Error("top-level build command is mounted instead of linked")
	}
	link, ok := CLI.Link("build")
	if !ok {
		t.Fatal("top-level build link is not registered")
	}
	if got, want := link.To, []string{"website", "build"}; !sameStrings(got, want) {
		t.Fatalf("build link destination = %q, want %q", got, want)
	}
	var output bytes.Buffer
	result := CLI.Run(
		clic.Execution{
			Context: context.Background(),
			Stdout:  &output,
		},
		clic.Arguments{
			"build",
			"--help",
		},
	)
	if result.Error != nil {
		t.Fatalf("CLI.Run(build --help) = %#v", result)
	}
	if !strings.Contains(output.String(), "hyperctl website build") {
		t.Fatalf("build help = %q", output.String())
	}
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
