package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	cli "github.com/herdworks/cli-go"
	"github.com/herdworks/hypertemplates-go/scripting"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
)

// hyperctl data
func DataCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "data",
			Description: "Manage global template data.",
			Usage:       "hyperctl cms data [command] [options]",
			Handler:     DataCommandHandler,
		},
		Options: cli.CommandOptions{},
		SubCommands: cli.Commands{
			DataListCommand(),
			DataInspectCommand(),
		},
	}
}
func DataCommandHandler(cmd *cli.Command) (err error) {
	cmd.Help(nil)
	return nil
}

// hyperctl data ls
func DataListCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "ls",
			Description: "List build context namespaces (ht, build, env, site, data, theme).",
			Usage:       "hyperctl data ls [options]",
			Handler:     DataListCommandHandler,
		},
		Options: cli.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Default:     "site.yaml",
				Required:    true,
			},
			{
				Name:        "verbose",
				ShortName:   "v",
				Description: "Route computed-data console.log output to stderr (default: silenced).",
				Boolean:     true,
				Default:     "false",
			},
		},
	}
}

func DataListCommandHandler(cmd *cli.Command) (err error) {
	build, err := loadBuildContext(cmd)
	if err != nil {
		return err
	}
	ctx := build.Context()
	if len(ctx) == 0 {
		fmt.Println("No build context configured.")
		return nil
	}
	header := "Namespace"
	colwidth := max(cmd.ColumnWidth(hypertexting.NamespacedData(ctx), header, 1), len(header)+1)
	fmt.Printf("%-*s %s\n", colwidth, header, "Size")
	fmt.Printf("%-*s %s\n", colwidth, strings.Repeat("-", len(header)), "----")
	for _, namespace := range ctx.Keys("") {
		bytes, err := json.Marshal(ctx.Get(namespace))
		if err != nil {
			return fmt.Errorf("error inspecting namespace %q: %v", namespace, err)
		}
		fmt.Printf("%-*s %v\n", colwidth, namespace, len(bytes))
	}
	return nil
}

// hyperctl data inspect
func DataInspectCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "inspect",
			Description: "Inspect a value in the build context by keypath (e.g. site.pages, data.social.twitter).",
			Usage:       "hyperctl data inspect <keypath> [options]",
			Validator:   DataInspectValidator,
			Handler:     DataInspectCommandHandler,
		},
		Options: cli.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Default:     "site.yaml",
				Required:    true,
				Secret:      false,
			},
			{
				Name:        "verbose",
				ShortName:   "v",
				Description: "Route computed-data console.log output to stderr (default: silenced).",
				Boolean:     true,
				Default:     "false",
			},
		},
	}
}
func DataInspectValidator(cmd *cli.Command) (err error) {
	requiredArgs := 1
	if len(cmd.Args) != requiredArgs {
		return fmt.Errorf("expected %v argument(s), got %v (%s)", requiredArgs, len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}
func DataInspectCommandHandler(cmd *cli.Command) (err error) {
	build, err := loadBuildContext(cmd)
	if err != nil {
		return err
	}
	keypath := cmd.Args[0]
	value := build.Context().Get(keypath)
	if value == nil {
		return fmt.Errorf("no value at keypath %q", keypath)
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}

// loadBuildContext constructs a Build for the website at the
// --config path and runs LoadContext so build.Context() is fully
// populated for inspection. The build is never written; URI is set to
// an empty URL placeholder since the write phase is skipped.
//
// Computed-data scripts may call console.log. By default those writes
// are redirected to io.Discard so callers can pipe `cms data inspect`
// output to tools like `jq` without noise on stdout. When the --verbose
// flag is set the output is routed to stderr instead, keeping stdout
// clean while still surfacing diagnostics.
func loadBuildContext(cmd *cli.Command) (*hypertexting.Build, error) {
	websiteURI, err := href.Parse(cmd.Options.Get("config"))
	if err != nil {
		return nil, err
	}
	build := &hypertexting.Build{
		URI:     &url.URL{},
		ID:      "inspect",
		Website: &hypertexting.Website{URI: websiteURI},
		Pages:   &hypertexting.Pages{},
	}
	logWriter := io.Writer(io.Discard)
	if cmd.Options.Get("verbose") == "true" {
		scripting.SetConsoleOutput(os.Stderr)
		logWriter = os.Stderr
	} else {
		scripting.SetConsoleOutput(io.Discard)
	}
	if err := build.LoadContext(context.Background(), logWriter); err != nil {
		return nil, err
	}
	return build, nil
}

func first(incoming ...any) string {
	var s string
	var ok bool
	for _, i := range incoming {
		if i == nil {
			continue
		} // skip if nil
		if s, ok = i.(string); !ok {
			continue
		} // skip if not a string
		if len(strings.TrimSpace(s)) > 0 {
			return s
		}
	}
	return ""
}
