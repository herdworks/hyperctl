package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/herdworks/clic"
	htskills "github.com/herdworks/htskills"
)

func SkillCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "skill",
			Description: "Read, install, or update bundled portable Agent Skills.",
			Usage: []string{
				"hyperctl agent skill [command]",
			},
		},
		Spec: clic.CommandSpec{
			Key:  "skill",
			Help: clic.DefaultHelp,
		},
		Commands: clic.MustCommandsFrom(
			SkillListCommand,
			SkillDescribeCommand,
			SkillInstallCommand,
			SkillUpdateCommand,
		),
	}
}
func SkillListCommand() clic.Command {
	var skills htskills.Skills
	var err error
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "ls",
			Description: "List bundled offline skills.",
			Usage: []string{
				"hyperctl agent skill ls",
			},
		},
		Spec: clic.CommandSpec{
			Key:  "ls",
			Help: clic.DefaultHelp,
			Handler: func(execution *clic.Execution) error {
				if skills, err = htskills.Catalog(); err != nil {
					return err
				}
				return clic.WriteList(execution.Stdout, clic.ColumnarFormat, skills,
					clic.Field{Key: "Name", Header: "Name"},
					clic.Field{Key: "Description", Header: "Description"},
				)
			},
		},
	}
}

func SkillDescribeCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "describe",
			Description: "Write the raw embedded SKILL.md bytes.",
			Usage: []string{
				"hyperctl agent skill describe <name>",
			},
		},
		Spec: clic.CommandSpec{
			Key:       "describe",
			Help:      clic.DefaultHelp,
			Validator: expectedPositionalArguments(1),
			Handler: func(execution *clic.Execution) error {
				content, err := skillContent(execution.Command.State.Arguments[0])
				if err != nil {
					return err
				}
				_, err = fmt.Fprint(execution.Stdout, string(content))
				return err
			},
		},
	}
}

func SkillInstallCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "install",
			Description: "Copy a complete skill directory to a new destination.",
			Usage: []string{
				"hyperctl agent skill install <name> <path>",
			},
		},
		Spec: clic.CommandSpec{
			Key:       "install",
			Help:      clic.DefaultHelp,
			Validator: expectedPositionalArguments(2),
			Handler: func(execution *clic.Execution) error {
				written, err := installSkill(execution.Command.State.Arguments[0], execution.Command.State.Arguments[1])
				if err != nil {
					return err
				}
				for _, name := range written {
					if _, err := fmt.Fprintln(execution.Stdout, name); err != nil {
						return err
					}
				}
				return nil
			},
		},
	}
}

func SkillUpdateCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "update",
			Description: "Update an installed skill when the bundled version is newer.",
			Usage: []string{
				"hyperctl agent skill update <name> <path>",
			},
		},
		Spec: clic.CommandSpec{
			Key:       "update",
			Help:      clic.DefaultHelp,
			Validator: expectedPositionalArguments(2),
			Handler: func(execution *clic.Execution) error {
				written, err := updateSkill(execution.Command.State.Arguments[0], execution.Command.State.Arguments[1])
				if err != nil {
					return err
				}
				for _, name := range written {
					if _, err := fmt.Fprintln(execution.Stdout, name); err != nil {
						return err
					}
				}
				return nil
			},
		},
	}
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
	root := htskills.Root + "/" + name
	written := []string{}
	err := fs.WalkDir(htskills.FS, root, func(path string, entry fs.DirEntry, err error) error {
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
		content, err := htskills.FS.ReadFile(path)
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
