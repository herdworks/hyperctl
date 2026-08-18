package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/herdworks/clic"
)

func TestAgentCommandMountsSkillTreeAndLegacyReadmeLink(t *testing.T) {
	command := AgentCommand()
	for _, name := range []string{"schema", "skill"} {
		if _, ok := command.SubCommand(name); !ok {
			t.Fatalf("agent command does not register %q", name)
		}
	}
	if _, ok := command.Link("readme"); !ok {
		t.Fatal("agent command does not link the legacy readme command")
	}
	for _, name := range []string{"ls", "describe"} {
		if _, ok := SchemaCommand().SubCommand(name); !ok {
			t.Fatalf("schema command does not register %q", name)
		}
	}
	skill, ok := command.SubCommand("skill")
	if !ok {
		t.Fatal("agent command does not register skill")
	}
	if got, want := skill.Meta.Usage, []string{"hyperctl agent skill [command]"}; !slices.Equal(got, want) {
		t.Fatalf("skill usage = %v, want %v", got, want)
	}
}

func TestAgentReadmeUsesHyperctlSkillNamespace(t *testing.T) {
	var output bytes.Buffer
	result := CLI.Run(clic.Execution{
		Context: context.Background(),
		Stdout:  &output,
	}, clic.Arguments{"agent", "readme"})
	if result.Error != nil {
		t.Fatalf("hyperctl agent readme = %#v", result)
	}
	if !strings.Contains(output.String(), `"hyperctl agent skill up"`) {
		t.Fatalf("agent readme does not use the Hyperctl namespace: %q", output.String())
	}
	if strings.Contains(output.String(), "htskills") {
		t.Fatalf("agent readme leaks the imported command name: %q", output.String())
	}
}

func TestAgentSkillUpUsesHyperctlNamespace(t *testing.T) {
	result := CLI.Run(clic.Execution{Context: context.Background(), WorkingDirectory: t.TempDir(), LookupEnv: func(string) (string, bool) { return t.TempDir(), true }}, clic.Arguments{"agent", "skill", "up"})
	if result.ExitCode != 1 || result.Error == nil {
		t.Fatalf("hyperctl agent skill up = %#v, want validation error", result)
	}
	if !strings.Contains(result.Error.Error(), "hyperctl agent skill up") || strings.Contains(result.Error.Error(), "htskills up") {
		t.Fatalf("hyperctl agent skill up error does not use the Hyperctl namespace: %q", result.Error)
	}
}

func TestSchemaListUsesAlignedColumns(t *testing.T) {
	entries := schemas.Entries()
	command := SchemaListCommand()
	var output strings.Builder
	if err := command.Spec.Handler(&clic.Execution{Stdout: &output}); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSuffix(output.String(), "\n"), "\n")
	if len(lines) != len(entries)+2 {
		t.Fatalf("schema list has %d lines, want %d", len(lines), len(entries)+2)
	}
	descriptionColumn := strings.Index(lines[0], "Description")
	if descriptionColumn <= len("Name") || !strings.HasPrefix(lines[1], "----") {
		t.Fatalf("schema list header is not tabular: %q / %q", lines[0], lines[1])
	}
	for i, entry := range entries {
		if got := strings.Index(lines[i+2], entry.Summary); got != descriptionColumn {
			t.Fatalf("%s description starts in column %d, want %d", entry.Name, got, descriptionColumn)
		}
	}
}

func TestSchemasAreJSONAndHaveStableNames(t *testing.T) {
	want := []string{"attachment", "byline", "feed-config", "feed-data", "feed-selector", "html-link", "html-meta", "markdown", "mention-config", "page", "provider", "provider-filter-rules", "provider-secret", "sync", "sync-rules", "theme", "website", "website-config", "website-favicon"}
	if got := schemaNames(); !slices.Equal(got, want) {
		t.Fatalf("schema names = %v, want %v", got, want)
	}
	for _, entry := range schemas.Entries() {
		bytes, err := schemas.Describe(entry.Name)
		if err != nil {
			t.Fatalf("%s: %v", entry.Name, err)
		}
		var document map[string]any
		if err := json.Unmarshal(bytes, &document); err != nil {
			t.Fatalf("%s is not JSON: %v", entry.Name, err)
		}
	}
}

func schemaDocument(t *testing.T, name string) []byte {
	t.Helper()
	document, err := schemas.Describe(name)
	if err != nil {
		t.Fatal(err)
	}
	return document
}

func TestWebsiteAndPageDatesUseDateTime(t *testing.T) {
	for _, name := range []string{"website", "page"} {
		bytes := schemaDocument(t, name)
		if !strings.Contains(string(bytes), `"format": "date-time"`) {
			t.Fatalf("%s schema does not include date-time fields", name)
		}
	}
}

func TestDataSchemasExplicitlyAllowCustomProperties(t *testing.T) {
	for _, name := range []string{"website", "theme", "page", "provider"} {
		bytes := schemaDocument(t, name)
		if !strings.Contains(string(bytes), `"additionalProperties": true`) {
			t.Fatalf("%s schema does not explicitly allow custom properties", name)
		}
	}
}

func TestThemeSchemaMatchesThemeJSONContract(t *testing.T) {
	provider, ok := schemas.Lookup("theme")
	if !ok {
		t.Fatal("theme schema is not registered")
	}
	if provider.JSONSchemaDescription() != "A theme source configuration file stored as theme.json." {
		t.Fatalf("theme schema description = %q", provider.JSONSchemaDescription())
	}
	content := schemaDocument(t, "theme")
	if !strings.Contains(string(content), `"data"`) || strings.Contains(string(content), `"template"`) {
		t.Fatalf("theme schema has an invalid data contract: %s", content)
	}
}

func TestProviderSchemaMatchesHostingProviderContract(t *testing.T) {
	content := schemaDocument(t, "provider")
	var document map[string]any
	if err := json.Unmarshal(content, &document); err != nil {
		t.Fatal(err)
	}
	provider := document["$defs"].(map[string]any)["github.com/herdworks/htup.HostingProvider"].(map[string]any)
	if got := provider["required"]; !slices.Equal(got.([]any), []any{"kind", "endpoint"}) {
		t.Fatalf("provider required = %#v", got)
	}
	properties := provider["properties"].(map[string]any)
	if got := properties["kind"].(map[string]any)["description"]; !strings.Contains(got.(string), "'git'") || !strings.Contains(got.(string), "'s3'") {
		t.Fatalf("provider kind description = %#v, want documented built-in kinds", got)
	}
	if properties["secrets"].(map[string]any)["type"] != "array" {
		t.Fatalf("provider secrets = %#v, want array", properties["secrets"])
	}
	filter := document["$defs"].(map[string]any)["github.com/herdworks/htup.FilterRules"].(map[string]any)
	for _, name := range []string{"include", "exclude"} {
		property := properties[name].(map[string]any)
		if property["$ref"] == nil {
			t.Fatalf("provider %s = %#v, want a filter-rules reference", name, property)
		}
		if filter["type"] != "object" {
			t.Fatalf("provider %s = %#v, want object", name, property)
		}
		if _, ok := filter["properties"].(map[string]any)["paths"]; !ok {
			t.Fatalf("provider %s has no paths", name)
		}
	}
}

func TestUnknownSchemaFails(t *testing.T) {
	if _, ok := schemas.Lookup("unknown"); ok {
		t.Fatal("unknown schema found")
	}
}
