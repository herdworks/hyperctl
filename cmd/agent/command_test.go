package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	cli "github.com/herdworks/cli-go"
	htskills "github.com/herdworks/htskills"
)

func testSkills(t *testing.T) skillEntries {
	t.Helper()
	skills, err := bundledSkills()
	if err != nil {
		t.Fatal(err)
	}
	return skills
}

func TestAgentCommandRegistersRequiredCommands(t *testing.T) {
	command := Command()
	for _, name := range []string{"readme", "schema", "skill"} {
		if _, err := command.SubCommands.Lookup(name); err != nil {
			t.Fatalf("agent command does not register %q: %v", name, err)
		}
	}
	for _, name := range []string{"ls", "describe"} {
		if _, err := schemaCommand().SubCommands.Lookup(name); err != nil {
			t.Fatalf("schema command does not register %q: %v", name, err)
		}
	}
	for _, name := range []string{"ls", "describe", "install", "update"} {
		if _, err := skillCommand().SubCommands.Lookup(name); err != nil {
			t.Fatalf("skill command does not register %q: %v", name, err)
		}
	}
}

func TestAgentReadmeIntroducesCompleteWorkflow(t *testing.T) {
	readme := readmeCommand()
	if readme.Config.Description != "Show onboarding context for a new agent." {
		t.Fatalf("readme description = %q", readme.Config.Description)
	}
	for _, phrase := range []string{
		"HyperTemplates", "HyperTexting CMS", "agent schema", "agent skill ls", "agent skill describe",
		"hyperctl cms new website", "hyperctl cms new page", "hyperctl dev server", "hyperctl build", "hyperctl build --id <build-id>", "hyperctl deploy", "established commands directly for mutations", "ht-apply",
	} {
		if !strings.Contains(readmeContext, phrase) {
			t.Fatalf("agent readme context is missing %q", phrase)
		}
	}
	if err := noArguments(&cli.Command{}); err != nil {
		t.Fatalf("start command rejects no arguments: %v", err)
	}
}

func TestSkillListUsesAlignedColumns(t *testing.T) {
	skills := testSkills(t)
	lines := strings.Split(strings.TrimSuffix(skillListOutput(skills), "\n"), "\n")
	if len(lines) != len(skills)+2 {
		t.Fatalf("skill list has %d lines, want %d", len(lines), len(skills)+2)
	}
	descriptionColumn := strings.Index(lines[0], "Description")
	if descriptionColumn <= len("Name") || !strings.HasPrefix(lines[1], "----") {
		t.Fatalf("skill list header is not tabular: %q / %q", lines[0], lines[1])
	}
	for i, skill := range skills {
		if got := strings.Index(lines[i+2], skill.Description); got != descriptionColumn {
			t.Fatalf("%s description starts in column %d, want %d", skill.Name, got, descriptionColumn)
		}
	}
}

func TestSchemaListUsesAlignedColumns(t *testing.T) {
	entries := schemaEntries(schemas.Entries())
	lines := strings.Split(strings.TrimSuffix(schemaListOutput(entries), "\n"), "\n")
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
	want := []string{"attachment", "byline", "feed", "html-link", "html-meta", "markdown", "mention-config", "page", "provider", "provider-filter-rules", "provider-secret", "sync", "sync-rules", "theme", "website", "website-config", "website-favicon"}
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
	provider := document["$defs"].(map[string]any)["github.com/herdworks/hyperproviders-go.HostingProvider"].(map[string]any)
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
	filter := document["$defs"].(map[string]any)["github.com/herdworks/hyperproviders-go.FilterRules"].(map[string]any)
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

func TestEmbeddedSkillsAreCompleteAndExact(t *testing.T) {
	skills := testSkills(t)
	for _, skill := range skills {
		content, err := skillContent(skill.Name)
		if err != nil {
			t.Fatalf("%s: %v", skill.Name, err)
		}
		if !strings.HasPrefix(string(content), "---\nname: "+skill.Name+"\n") {
			t.Fatalf("%s has invalid skill frontmatter", skill.Name)
		}
		metadata, err := htskills.Parse(content)
		if err != nil {
			t.Fatalf("%s metadata: %v", skill.Name, err)
		}
		if metadata.Version == nil {
			t.Fatalf("%s has no skill version", skill.Name)
		}
		for _, heading := range []string{"## Overview", "## Workflow", "## Instructions", "## Tools", "## References"} {
			if !strings.Contains(string(content), heading) {
				t.Fatalf("%s is missing %q", skill.Name, heading)
			}
		}
	}
}

func TestDataAndLayoutSkillsCoverCoreAuthoringRules(t *testing.T) {
	data, err := skillContent("hypertexting-data")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"custom properties", "content/index.yaml", "data/homepage.yaml"} {
		if !strings.Contains(string(data), phrase) {
			t.Fatalf("data skill is missing %q", phrase)
		}
	}
	layout, err := skillContent("hypertexting-layouts")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"ht-with", "ht-debug", "source design dictate the authoring workflow", "full-page mockups", "component-level designs", "both paths must converge", "required outcome: reusable fragments", "Do not split a layout merely to anticipate reuse", "clear data contract or reuse opportunity", "rootless component", "fixture data", "top-level section", "component a fragment", "Compare extracted section fragments", "smaller shared fragments", "Use `hypertexting-html` to choose semantic elements", "Colocate a fragment’s component styles"} {
		if !strings.Contains(string(layout), phrase) {
			t.Fatalf("layout skill is missing %q", phrase)
		}
	}
	templateData, err := embeddedSkills.ReadFile(embeddedSkillsRoot + "/hypertexting-layouts/references/template-data.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"direct text-node children", "implicit `ht-apply`", "root-shallow", "do not reach descendant elements", "plain object renders one fragment", "each.index", "site", "page", "theme"} {
		if !strings.Contains(string(templateData), phrase) {
			t.Fatalf("layout template-data reference is missing %q", phrase)
		}
	}
	css, err := skillContent("hypertexting-css")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"Colocate component-specific CSS", "Include these two style elements", "style#components` before `style#layout", "source order precedence", "fixture JSON", "Style custom elements and their custom attributes directly", "few CSS classes", "fragment-owned component styles", "generic static `style.css`", "Do not create a fragment variant solely for layout-specific styling", "markup, data contract, or reusable behavior differs"} {
		if !strings.Contains(string(css), phrase) {
			t.Fatalf("CSS skill is missing %q", phrase)
		}
	}
}

func TestDataHtmlPluginsAndThemeSkillsDeclareOwnershipAndThemePatterns(t *testing.T) {
	data, err := skillContent("hypertexting-data")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Use `hypertexting-plugins` to implement the computed-data plugin") {
		t.Fatal("data skill does not hand computed-data implementation to the plugins skill")
	}
	namespaces, err := embeddedSkills.ReadFile(embeddedSkillsRoot + "/hypertexting-data/references/namespaces.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"social/twitter.yaml", "social.twitter", "archive.js", "static file and one computed `.js` file"} {
		if !strings.Contains(string(namespaces), phrase) {
			t.Fatalf("data namespaces reference is missing %q", phrase)
		}
	}
	html, err := skillContent("hypertexting-html")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"Author the HTML first", "semantic HTML", "Avoid div soup", "hyphenated custom element", "few or none"} {
		if !strings.Contains(string(html), phrase) {
			t.Fatalf("HTML skill is missing %q", phrase)
		}
	}
	plugins, err := skillContent("hypertexting-plugins")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"`hypertexting-data` has defined a reusable derived namespace", "social/twitter.js` becomes `social.twitter`", "template-value plugin", "text-transform: uppercase", "alternate HTML authoring system"} {
		if !strings.Contains(string(plugins), phrase) {
			t.Fatalf("plugins skill is missing %q", phrase)
		}
	}
	pluginReference, err := embeddedSkills.ReadFile(embeddedSkillsRoot + "/hypertexting-plugins/references/choosing-plugins.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"duration_iso", "export default function durationISO", "export default function durationLabel", "Computed-data plugin", "Web Component"} {
		if !strings.Contains(string(pluginReference), phrase) {
			t.Fatalf("plugin reference is missing %q", phrase)
		}
	}
	themes, err := skillContent("hypertexting-themes")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"`theme.json`", "themes/<theme>/theme.json", "themes/<theme>/<version>/theme.json", "site.config.theme", "embedded builtins only as a fallback"} {
		if !strings.Contains(string(themes), phrase) {
			t.Fatalf("themes skill is missing %q", phrase)
		}
	}
	config, err := embeddedSkills.ReadFile(embeddedSkillsRoot + "/hypertexting-themes/references/config.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"JSON", "layouts_dir", "fragments_dir", "static_dir", "data"} {
		if !strings.Contains(string(config), phrase) {
			t.Fatalf("theme config reference is missing %q", phrase)
		}
	}
}

func TestHypertextingSkillRoutesAndManagesFocusedSkills(t *testing.T) {
	content, err := skillContent("hypertexting")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"hypertexting-data",
		"hypertexting-themes",
		"hypertexting-layouts",
		"hypertexting-debugging",
		"hypertexting-assets",
		"hypertexting-html",
		"hypertexting-plugins",
		"hypertexting-css",
		"hypertexting-javascript",
	} {
		for _, command := range []string{
			"hyperctl agent skill install " + name + " <path>",
			"hyperctl agent skill update " + name + " <path>",
			"hyperctl agent skill describe " + name,
		} {
			if !strings.Contains(string(content), command) {
				t.Fatalf("orchestration skill is missing %q", command)
			}
		}
	}
	for _, phrase := range []string{"hypertexting-html` to author semantic elements", "hypertexting-layouts` to compose HTML", "hypertexting-themes` to organize reusable layouts", "hyperctl build --id <build-id>", "Do not wrap `hyperctl build`"} {
		if !strings.Contains(string(content), phrase) {
			t.Fatalf("orchestration skill is missing %q", phrase)
		}
	}
	architecture, err := embeddedSkills.ReadFile(embeddedSkillsRoot + "/hypertexting/references/architecture.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"pre-renders HTML and assets at build time", "islands-style enhancement", "serverless functions"} {
		if !strings.Contains(string(architecture), phrase) {
			t.Fatalf("architecture reference is missing %q", phrase)
		}
	}
}

func TestDebuggingSkillUsesDeclarativeDebuggingAndKnownIssues(t *testing.T) {
	debugging, err := skillContent("hypertexting-debugging")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"`ht-debug`", "Do not write a custom script", "known issues"} {
		if !strings.Contains(string(debugging), phrase) {
			t.Fatalf("debugging skill is missing %q", phrase)
		}
	}
	debugReference, err := embeddedSkills.ReadFile(embeddedSkillsRoot + "/hypertexting-debugging/references/ht-debug.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"ht-debug=\"element\"", "${ page.hero }", "exactly one", "Repeated `ht-debug` attributes"} {
		if !strings.Contains(string(debugReference), phrase) {
			t.Fatalf("ht-debug reference is missing %q", phrase)
		}
	}
}

func TestJavaScriptSkillSeparatesBrowserEnhancementsFromBuildPlugins(t *testing.T) {
	javascript, err := skillContent("hypertexting-javascript")
	if err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"pre-rendered HTML component", "build-time JavaScript plugins", "hypertexting-plugins", "document.addEventListener(\"<event>\", this)", "destination-side `ht-pipe='from \"script.component\" as js'`"} {
		if !strings.Contains(string(javascript), phrase) {
			t.Fatalf("JavaScript skill is missing %q", phrase)
		}
	}
	for _, path := range []string{
		embeddedSkillsRoot + "/hypertexting-javascript/references/runtime-boundaries.md",
		embeddedSkillsRoot + "/hypertexting-javascript/references/web-components.md",
	} {
		if _, err = embeddedSkills.ReadFile(path); err != nil {
			t.Fatalf("JavaScript skill reference %q: %v", path, err)
		}
	}
}

func TestInstallSkillCopiesCompleteDirectoryAndRejectsExistingDestination(t *testing.T) {
	skill := testSkills(t)[0]
	destination := filepath.Join(t.TempDir(), "installed-skill")
	written, err := installSkill(skill.Name, destination)
	if err != nil {
		t.Fatal(err)
	}
	if len(written) == 0 {
		t.Fatal("install wrote no files")
	}
	installed, err := os.ReadFile(filepath.Join(destination, "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	embedded, err := skillContent(skill.Name)
	if err != nil {
		t.Fatal(err)
	}
	if string(installed) != string(embedded) {
		t.Fatal("installed SKILL.md differs from embedded bytes")
	}
	if info, err := os.Stat(filepath.Join(destination, "references")); err != nil || !info.IsDir() {
		t.Fatalf("installed skill does not contain references directory: %v", err)
	}
	if _, err := installSkill(skill.Name, destination); err == nil {
		t.Fatal("install succeeded for an existing destination")
	}
}

func TestUpdateSkillReplacesOlderAndUnversionedSkillsOnly(t *testing.T) {
	skill := testSkills(t)[0]
	destination := filepath.Join(t.TempDir(), "installed-skill")
	if _, err := installSkill(skill.Name, destination); err != nil {
		t.Fatal(err)
	}
	installedPath := filepath.Join(destination, "SKILL.md")
	embedded, err := skillContent(skill.Name)
	if err != nil {
		t.Fatal(err)
	}
	metadata, err := htskills.Parse(embedded)
	if err != nil {
		t.Fatal(err)
	}
	if metadata.Version == nil {
		t.Fatal("embedded skill has no version")
	}
	versionLine := fmt.Sprintf("version: %d\n", *metadata.Version)
	if err := os.WriteFile(installedPath, []byte(strings.Replace(string(embedded), versionLine, "", 1)), 0644); err != nil {
		t.Fatal(err)
	}
	written, err := updateSkill(skill.Name, destination)
	if err != nil {
		t.Fatalf("update unversioned skill: %v", err)
	}
	if len(written) == 0 {
		t.Fatal("update wrote no files")
	}
	updated, err := os.ReadFile(installedPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(updated) != string(embedded) {
		t.Fatal("updated SKILL.md differs from embedded bytes")
	}
	if _, err = updateSkill(skill.Name, destination); err == nil {
		t.Fatal("update succeeded for an equal skill version")
	}
	newerVersionLine := fmt.Sprintf("version: %d\n", *metadata.Version+1)
	if err = os.WriteFile(installedPath, []byte(strings.Replace(string(embedded), versionLine, newerVersionLine, 1)), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err = updateSkill(skill.Name, destination); err == nil {
		t.Fatal("update succeeded for a newer installed skill")
	}
	if _, err = updateSkill(skill.Name, filepath.Join(t.TempDir(), "missing-skill")); err == nil {
		t.Fatal("update succeeded for a missing installation")
	}
}

func TestUnknownNamesFail(t *testing.T) {
	if _, ok := schemas.Lookup("unknown"); ok {
		t.Fatal("unknown schema found")
	}
	if _, err := skillContent("unknown"); err == nil {
		t.Fatal("unknown skill did not fail")
	}
}
