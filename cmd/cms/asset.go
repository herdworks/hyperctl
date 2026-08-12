package cms

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	cli "github.com/herdworks/cli-go"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
)

func AssetCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "asset",
			Description: "Manage website assets.",
			Usage:       "hyperctl cms asset [command] [options]",
			Handler:     NewCommandHandler,
		},
		Options: cli.CommandOptions{},
		SubCommands: cli.Commands{
			AssetAddCommand(),
			AssetListCommand(),
			AssetLookupCommand(),
		},
	}
}

func AssetAddCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "add",
			Description: "Add a static asset to a website.",
			Usage:       "hyperctl cms asset add [options] <asset> <path>",
			Validator:   AssetAddCommandValidator,
			Handler:     AssetAddCommandHandler,
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
		},
	}
}

func AssetAddCommandValidator(cmd *cli.Command) (err error) {
	requiredArgs := 2
	if len(cmd.Args) != requiredArgs {
		return fmt.Errorf("expected %v argument(s), got %v (%s)", requiredArgs, len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}

func AssetAddCommandHandler(cmd *cli.Command) (err error) {
	sourcePath := cmd.Args[0]
	assetPath := cmd.Args[1]
	website := hypertexting.Website{URI: &href.Href{URL: url.URL{Path: cmd.Options.Get("config")}}}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	if err = website.CopyAsset(sourcePath, assetPath, nil); err != nil {
		return err
	}
	return nil
}

func AssetListCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "ls",
			Description: "List all static assets.",
			Usage:       "hyperctl cms asset ls [options]",
			Validator:   AssetListCommandValidator,
			Handler:     AssetListCommandHandler,
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
		},
	}
}

func AssetListCommandValidator(cmd *cli.Command) (err error) {
	return nil
}

func AssetListCommandHandler(cmd *cli.Command) (err error) {
	var assets AssetListRows
	website := hypertexting.Website{URI: &href.Href{URL: url.URL{Path: cmd.Options.Get("config")}}}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	pages, err := hypertexting.GetPages(&website, website.ContentDirPath())
	if err != nil {
		return err
	}
	if err = hypertexting.LoadPages(context.TODO(), &website, pages, nil, false, nil); err != nil {
		return err
	}
	assetEntries, err := hypertexting.AssetEntries(&website.Theme, &website, pages)
	if err != nil {
		return err
	}
	if assets, err = assetListRows(&website, assetEntries); err != nil {
		return err
	}
	AssetListCommandOutput(cmd, assets)
	return nil
}

func assetListRows(website *hypertexting.Website, assets hypertexting.Assets) (AssetListRows, error) {
	rows := AssetListRows{}
	for _, asset := range assets {
		stat, err := os.Stat(asset.SourcePath)
		if err != nil {
			return nil, err
		}
		row := AssetListRow{
			Source: website.RelativePath(asset.SourcePath),
			Size:   stat.Size(),
			Path:   asset.PublicPath,
		}
		rows = append(rows, row)
	}
	return rows, nil
}

type AssetListRow struct {
	Source string
	Size   int64
	Path   string
}

type AssetListRows []AssetListRow

func (assets AssetListRows) ColumnValues(field string) (values []string) {
	values = []string{}
	for i := range assets {
		switch field {
		case "Source", "source":
			values = append(values, assets[i].Source)
		case "Size", "size":
			values = append(values, fmt.Sprintf("%d", assets[i].Size))
		case "Path", "path":
			values = append(values, assets[i].Path)
		}
	}
	return values
}

func AssetListCommandOutput(cmd *cli.Command, assets AssetListRows) {
	indexcol := 3
	sourcecol := cmd.ColumnWidth(assets, "Source", 1)
	sizecol := cmd.ColumnWidth(assets, "Size", 1)
	pathcol := cmd.ColumnWidth(assets, "Path", 1)
	fmt.Printf("%-*s %-*s %-*v %-*s\n", indexcol, "No", sourcecol, "Source", sizecol, "Size", pathcol, "Path")
	fmt.Printf("%-*s %-*s %-*v %-*s\n", indexcol, "--", sourcecol, "------", sizecol, "----", pathcol, "----")
	for i := range assets {
		fmt.Printf(
			"%-*v %-*s %-*v %-*s\n",
			indexcol, i+1,
			sourcecol, assets[i].Source,
			sizecol, assets[i].Size,
			pathcol, assets[i].Path,
		)
	}
}

func AssetLookupCommand() *cli.Command {
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "lookup",
			Description: "Lookup a static asset by its path.",
			Usage:       "hyperctl cms asset lookup [options] <path>",
			Validator:   AssetLookupCommandValidator,
			Handler:     AssetLookupCommandHandler,
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
		},
	}
}

func AssetLookupCommandValidator(cmd *cli.Command) (err error) {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expected exactly 1 argument, got %v (%s)", len(cmd.Args), strings.Join(cmd.Args, ", "))
	}
	return nil
}

func AssetLookupCommandHandler(cmd *cli.Command) (err error) {
	assetPath := cmd.Args[0]
	website := hypertexting.Website{URI: &href.Href{URL: url.URL{Path: cmd.Options.Get("config")}}}
	if err = website.Load(false, nil, nil); err != nil {
		return err
	}
	pages, err := hypertexting.GetPages(&website, website.ContentDirPath())
	if err != nil {
		return err
	}
	if err = hypertexting.LoadPages(context.TODO(), &website, pages, nil, false, nil); err != nil {
		return err
	}
	candidates, err := hypertexting.AssetCandidates(&website.Theme, &website, pages)
	if err != nil {
		return err
	}
	matches := assetLookupMatches(candidates, assetPath)
	fmt.Printf("asset \"%s\" lookup\n", assetPath)
	if len(matches) == 0 {
		fmt.Printf("asset \"%s\" not found\n", assetPath)
		return nil
	}
	for _, match := range matches {
		fmt.Printf("asset \"%s\" candidate \"%s\"\n", assetPath, website.RelativePath(match.SourcePath))
	}
	winner := matches[len(matches)-1]
	fmt.Printf("asset \"%s\" source \"%s\"\n", assetPath, website.RelativePath(winner.SourcePath))
	return nil
}

func assetLookupMatches(candidates hypertexting.Assets, query string) hypertexting.Assets {
	normalized := strings.TrimPrefix(query, "/")
	matches := hypertexting.Assets{}
	for _, asset := range candidates {
		if asset.RelPath == normalized {
			matches = append(matches, asset)
		}
	}
	return matches
}
