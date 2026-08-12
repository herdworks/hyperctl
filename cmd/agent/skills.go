package agent

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	cli "github.com/herdworks/cli-go"
	htskills "github.com/herdworks/htskills"
)

const embeddedSkillsRoot = htskills.Root

var embeddedSkills = htskills.FS

type skillEntries []htskills.Skill

func (entries skillEntries) ColumnValues(field string) (values []string) {
	for _, entry := range entries {
		switch field {
		case "Name":
			values = append(values, entry.Name)
		case "Description":
			values = append(values, entry.Description)
		}
	}
	return values
}

func skillCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{Name: "skill", Description: "Read, install, or update bundled portable Agent Skills.", Usage: "hyperctl agent skill [command]", Handler: help}, SubCommands: cli.Commands{skillListCommand(), skillDescribeCommand(), skillInstallCommand(), skillUpdateCommand()}}
}
func skillListCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{Name: "ls", Description: "List bundled offline skills.", Usage: "hyperctl agent skill ls", Handler: func(*cli.Command) error {
		skills, err := bundledSkills()
		if err != nil {
			return fail(err)
		}
		fmt.Print(skillListOutput(skills))
		return nil
	}}}
}

func bundledSkills() (skillEntries, error) {
	skills, err := htskills.Skills()
	return skillEntries(skills), err
}

func skillListOutput(skills skillEntries) string {
	namecol := (&cli.Command{}).ColumnWidth(skills, "Name", 1)
	var output strings.Builder
	fmt.Fprintf(&output, "%-*s %s\n", namecol, "Name", "Description")
	fmt.Fprintf(&output, "%-*s %s\n", namecol, "----", "-----------")
	for _, skill := range skills {
		fmt.Fprintf(&output, "%-*s %s\n", namecol, skill.Name, skill.Description)
	}
	return output.String()
}
func skillDescribeCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{Name: "describe", Description: "Write the raw embedded SKILL.md bytes.", Usage: "hyperctl agent skill describe <name>", Handler: func(command *cli.Command) error {
		if err := oneArgument(command); err != nil {
			return fail(err)
		}
		content, err := skillContent(command.Args[0])
		if err != nil {
			return fail(err)
		}
		_, err = fmt.Print(string(content))
		return err
	}}}
}
func skillInstallCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{Name: "install", Description: "Copy a complete skill directory to a new destination.", Usage: "hyperctl agent skill install <name> <path>", Handler: func(command *cli.Command) error {
		if err := twoArguments(command); err != nil {
			return fail(err)
		}
		written, err := installSkill(command.Args[0], command.Args[1])
		if err != nil {
			return fail(err)
		}
		for _, name := range written {
			fmt.Println(name)
		}
		return nil
	}}}
}

func skillUpdateCommand() *cli.Command {
	return &cli.Command{Config: cli.CommandConfig{Name: "update", Description: "Update an installed skill when the bundled version is newer.", Usage: "hyperctl agent skill update <name> <path>", Handler: func(command *cli.Command) error {
		if err := twoArguments(command); err != nil {
			return fail(err)
		}
		written, err := updateSkill(command.Args[0], command.Args[1])
		if err != nil {
			return fail(err)
		}
		for _, name := range written {
			fmt.Println(name)
		}
		return nil
	}}}
}
func twoArguments(command *cli.Command) error {
	if len(command.Args) != 2 {
		return fmt.Errorf("expected exactly 2 arguments, got %d", len(command.Args))
	}
	return nil
}
func skillContent(name string) ([]byte, error) {
	content, err := htskills.Read(name)
	if err != nil {
		return nil, fmt.Errorf("%w (use hyperctl agent skill ls)", err)
	}
	return content, nil
}

func installSkill(name, destination string) ([]string, error) {
	if _, err := skillContent(name); err != nil {
		return nil, err
	}
	if _, err := os.Lstat(destination); err == nil {
		return nil, fmt.Errorf("destination already exists: %s", destination)
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if err := os.Mkdir(destination, 0755); err != nil {
		return nil, err
	}
	return copyEmbeddedSkill(name, destination)
}

func updateSkill(name, destination string) ([]string, error) {
	if _, err := skillContent(name); err != nil {
		return nil, err
	}
	info, err := os.Stat(destination)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("skill is not installed: %s", destination)
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("skill installation is not a directory: %s", destination)
	}
	installed, err := os.ReadFile(filepath.Join(destination, "SKILL.md"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("skill is not installed: %s", destination)
		}
		return nil, err
	}
	installedMetadata, err := htskills.Parse(installed)
	if err != nil {
		return nil, fmt.Errorf("installed skill %q: %w", destination, err)
	}
	if installedMetadata.Name != name {
		return nil, fmt.Errorf("installed skill name %q does not match requested skill %q", installedMetadata.Name, name)
	}
	embedded, err := skillContent(name)
	if err != nil {
		return nil, err
	}
	embeddedMetadata, err := htskills.Parse(embedded)
	if err != nil {
		return nil, fmt.Errorf("bundled skill %q: %w", name, err)
	}
	if embeddedMetadata.Version == "" {
		return nil, fmt.Errorf("bundled skill %q is unversioned", name)
	}
	if installedMetadata.Version == embeddedMetadata.Version {
		return nil, fmt.Errorf("skill %q is already at version %s", name, installedMetadata.Version)
	}
	if installedMetadata.Version > embeddedMetadata.Version {
		return nil, fmt.Errorf("skill %q has newer installed version %s (bundled version %s)", name, installedMetadata.Version, embeddedMetadata.Version)
	}
	return copyEmbeddedSkill(name, destination)
}

func copyEmbeddedSkill(name, destination string) ([]string, error) {
	root := embeddedSkillsRoot + "/" + name
	written := []string{}
	err := fs.WalkDir(embeddedSkills, root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative := strings.TrimPrefix(path, root)
		if relative == "" {
			return nil
		}
		output := filepath.Join(destination, filepath.FromSlash(relative))
		if entry.IsDir() {
			return os.MkdirAll(output, 0755)
		}
		content, err := embeddedSkills.ReadFile(path)
		if err != nil {
			return err
		}
		if err = os.WriteFile(output, content, 0644); err != nil {
			return err
		}
		written = append(written, output)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(written)
	return written, nil
}
