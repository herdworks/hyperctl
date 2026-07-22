package theme

import (
	"testing"
)

func TestThemeCommandRegistersMigrate(t *testing.T) {
	command := Command()
	for _, subcommand := range command.SubCommands {
		if subcommand.Config.Name == "migrate" {
			return
		}
	}
	t.Fatal("theme migrate command is not registered")
}

func TestThemeMigrateCommandOptions(t *testing.T) {
	command := ThemeMigrateCommand()
	theme := command.Options.Lookup("theme")
	if theme == nil || theme.ShortName != "t" || !theme.Required || theme.Default != "" || theme.EnvName != "" {
		t.Fatalf("theme option = %#v", theme)
	}
	output := command.Options.Lookup("output")
	if output == nil || output.Boolean || output.Required {
		t.Fatalf("output option = %#v", output)
	}
	interactive := command.Options.Lookup("interactive")
	if interactive == nil || !interactive.Boolean {
		t.Fatalf("interactive option = %#v", interactive)
	}
	if command.Options.Lookup("dry-run") != nil {
		t.Fatal("dry-run option is still registered")
	}
	if err := command.Validate(); err == nil {
		t.Fatal("command without --theme unexpectedly validated")
	}
	theme.RawValue = "theme.json"
	if err := command.Validate(); err != nil {
		t.Fatalf("command with --theme did not validate: %v", err)
	}
}

func TestThemeMigrateCommandParsesOutputAndInteractiveMode(t *testing.T) {
	command := ThemeMigrateCommand()
	if err := command.ParseArguments([]string{"-t", "theme.json", "--output", "migrated", "--interactive"}); err != nil {
		t.Fatal(err)
	}
	if got := command.Options.Get("theme"); got != "theme.json" {
		t.Fatalf("theme = %q, want theme.json", got)
	}
	if got := command.Options.Get("output"); got != "migrated" {
		t.Fatalf("output = %q, want migrated", got)
	}
	if command.Options.Get("interactive") != "true" {
		t.Fatalf("interactive = %q, want true", command.Options.Get("interactive"))
	}
}
