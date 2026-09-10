package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/herdworks/clic"
)

func TestSkillCommandMountsSkillTree(t *testing.T) {
	command := SkillCommand()
	for _, name := range []string{"readme", "ls", "describe", "install", "update", "up"} {
		if _, ok := command.SubCommand(name); !ok {
			t.Fatalf("skill command does not register %q", name)
		}
	}
	if got, want := command.Meta.Usage, []string{"hyperctl skill [command]"}; !sameStrings(got, want) {
		t.Fatalf("skill usage = %v, want %v", got, want)
	}
	if _, ok := CLI.SubCommand("agent"); ok {
		t.Fatal("hyperctl still registers agent")
	}
}

func TestSkillReadmeUsesHyperctlNamespace(t *testing.T) {
	var output bytes.Buffer
	result := CLI.Run(clic.Execution{
		Context: context.Background(),
		Stdout:  &output,
	}, clic.Arguments{"skill", "readme"})
	if result.Error != nil {
		t.Fatalf("hyperctl skill readme = %#v", result)
	}
	if !strings.Contains(output.String(), `"hyperctl skill up"`) {
		t.Fatalf("skill readme does not use the Hyperctl namespace: %q", output.String())
	}
	if strings.Contains(output.String(), "htskills") {
		t.Fatalf("skill readme leaks the imported command name: %q", output.String())
	}
	output.Reset()
	result = CLI.Run(clic.Execution{Context: context.Background(), Stdout: &output}, clic.Arguments{"readme"})
	if result.Error != nil {
		t.Fatalf("hyperctl readme = %#v", result)
	}
	if !strings.Contains(output.String(), `"hyperctl skill up"`) {
		t.Fatalf("hyperctl readme does not target skill: %q", output.String())
	}
}

func TestSkillUpUsesHyperctlNamespace(t *testing.T) {
	result := CLI.Run(clic.Execution{Context: context.Background(), WorkingDirectory: t.TempDir(), LookupEnv: func(string) (string, bool) { return t.TempDir(), true }}, clic.Arguments{"skill", "up"})
	if result.ExitCode != 1 || result.Error == nil {
		t.Fatalf("hyperctl skill up = %#v, want validation error", result)
	}
	if !strings.Contains(result.Error.Error(), "hyperctl skill up") || strings.Contains(result.Error.Error(), "htskills up") {
		t.Fatalf("hyperctl skill up error does not use the Hyperctl namespace: %q", result.Error)
	}
}
