package cli

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	hypertemplates "github.com/herdworks/hypertemplates-go/lib"
	hypertexting "github.com/herdworks/hypertexting-go/lib"
)

func AssetCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "asset",
			Description: "Manage website assets.",
			Usage:       "hyperctl asset [command] [options]",
			Handler:     NewCommandHandler,
		},
		Options: CommandOptions{},
		SubCommands: Commands{
			AssetAddCommand(),
			AssetListCommand(),
			AssetLookupCommand(),
		},
	}
}

func AssetAddCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "add",
			Description: "Add a static asset to a website.",
			Usage:       "hyperctl asset add [options] <asset> <path>",
			Validator:   AssetAddCommandValidator,
			Handler:     AssetAddCommandHandler,
		},
		Options: CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Required:    true,
				Secret:      false,
			},
		},
	}
}

func AssetAddCommandValidator(cmd *Command) (err error) {
	requiredArgs := 2
	if len(cmd.Args) != requiredArgs {
		return fmt.Errorf("expected %v argument(s), got %v (%s)", requiredArgs, len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}

func AssetAddCommandHandler(cmd *Command) (err error) {
	sourcePath := cmd.Args[0]
	assetPath := cmd.Args[1]
	website := hypertexting.Website{URI: &url.URL{Path: cmd.Options.Get("config")}}
	if err = website.Load(false, nil); err != nil {
		return err
	}
	if err = website.CopyAsset(sourcePath, assetPath); err != nil {
		return err
	}
	return nil
}

func AssetListCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "ls",
			Description: "List all static assets.",
			Usage:       "hyperctl asset ls [options]",
			Validator:   AssetListCommandValidator,
			Handler:     AssetListCommandHandler,
		},
		Options: CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Required:    true,
				Secret:      false,
			},
		},
	}
}

func AssetListCommandValidator(cmd *Command) (err error) {
	return nil
}

func AssetListCommandHandler(cmd *Command) (err error) {
	var assets hypertemplates.HyperFiles
	website := hypertexting.Website{URI: &url.URL{Path: cmd.Options.Get("config")}}
	if err = website.Load(false, nil); err != nil {
		return err
	} // load website + theme assets
	if website.Pages, err = hypertexting.GetPages(&website, website.ContentDir()); err != nil {
		return err
	}
	if err = hypertexting.LoadPages(context.TODO(), &website, website.Pages, nil, false); err != nil {
		return err
	} // load page assets
	assets = append(assets, website.Assets...)
	for _, page := range website.Pages {
		assets = append(assets, page.Assets...)
	}
	AssetListCommandOutput(cmd, assets)
	return nil
}
func AssetListCommandOutput(cmd *Command, assets hypertemplates.HyperFiles) {
	indexcol := 3
	sourcecol := cmd.ColumnWidth(assets, "Source", 1)
	sizecol := cmd.ColumnWidth(assets, "Size", 1)
	pathcol := cmd.ColumnWidth(assets, "Path", 2) // add padding for / prefix
	fmt.Printf("%-*s %-*s %-*v %-*s\n", indexcol, "No", sourcecol, "Source", sizecol, "Size", pathcol, "Path")
	fmt.Printf("%-*s %-*s %-*v %-*s\n", indexcol, "--", sourcecol, "------", sizecol, "----", pathcol, "----")
	for i := range assets {
		fmt.Printf(
			"%-*v %-*s %-*v %-*s\n",
			indexcol, i+1,
			sourcecol, assets[i].RelPath(),
			sizecol, assets[i].Size(),
			pathcol, fmt.Sprintf("/%s", assets[i].SubPath()),
		)
	}
}

func AssetLookupCommand() *Command {
	return &Command{
		Config: CommandConfig{
			Name:        "lookup",
			Description: "Lookup a static asset by its path.",
			Usage:       "hyperctl asset lookup [options] <path>",
			Validator:   AssetLookupCommandValidator,
			Handler:     AssetLookupCommandHandler,
		},
		Options: CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Required:    true,
				Secret:      false,
			},
		},
	}
}

func AssetLookupCommandValidator(cmd *Command) (err error) {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expected exactly 1 argument, got %v (%s)", len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}

func AssetLookupCommandHandler(cmd *Command) (err error) {
	assetPath := cmd.Args[0]
	website := hypertexting.Website{URI: &url.URL{Path: cmd.Options.Get("config")}}
	if err = website.Load(false, nil); err != nil {
		return err
	}
	assetFilePath := website.GetAssetURI(assetPath)
	if assetFilePath != nil {
		fmt.Printf("asset \"%s\" source \"%s\" found\n", assetPath, assetFilePath.Path)
	}
	return nil
}
