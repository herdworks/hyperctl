package cli

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/herdworks/clic"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
)

func DevServerCommand() clic.Command {
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "server",
			Description: "Live-reloading development server.",
			Usage: []string{
				"hyperctl dev server [options] [content]",
			},
		},
		Spec: clic.CommandSpec{
			Key: "server", Help: clic.DefaultHelp,
			Options: serverOptions(),
			Handler: serverHandler,
		},
	}
}
func serverOptions() clic.Options {
	return clic.Options{
		"config":      {Meta: clic.OptionMeta{Name: "config", ShortName: "c", EnvName: "HYPER_CONFIG", Description: "Path to the website configuration file.", Default: "site.yaml"}, Spec: clic.OptionSpec{Key: "config", Kind: clic.StringOption, Required: true}},
		"id":          {Meta: clic.OptionMeta{Name: "id", ShortName: "i", EnvName: "HYPER_BUILD_ID", Description: "Build ID.", Default: uuid.New().String()}, Spec: clic.OptionSpec{Key: "id", Kind: clic.StringOption}},
		"interval":    {Meta: clic.OptionMeta{Name: "interval", ShortName: "I", EnvName: "HYPER_REFRESH", Description: "Live reload interval seconds.", Default: "2"}, Spec: clic.OptionSpec{Key: "interval", Kind: clic.StringOption}},
		"host":        {Meta: clic.OptionMeta{Name: "host", ShortName: "H", EnvName: "HYPER_HOST", Description: "Development server host.", Default: "localhost"}, Spec: clic.OptionSpec{Key: "host", Kind: clic.StringOption}},
		"port":        {Meta: clic.OptionMeta{Name: "port", ShortName: "P", EnvName: "HYPER_PORT", Description: "Development server port.", Default: "8080"}, Spec: clic.OptionSpec{Key: "port", Kind: clic.StringOption}},
		"content-dir": {Meta: clic.OptionMeta{Name: "content-dir", EnvName: "HYPER_CONTENT_DIR", Description: "Relative path to the content directory to use for the build (overwrites site.config.content_dir)."}, Spec: clic.OptionSpec{Key: "content-dir", Kind: clic.StringOption}},
		"data-dir":    {Meta: clic.OptionMeta{Name: "data-dir", EnvName: "HYPER_DATA_DIR", Description: "Relative path to the data directory to use for the build (overwrites site.config.data_dir)."}, Spec: clic.OptionSpec{Key: "data-dir", Kind: clic.StringOption}},
		"static-dir":  {Meta: clic.OptionMeta{Name: "static-dir", EnvName: "HYPER_STATIC_DIR", Description: "Relative path to the static directory to use for the build (overwrites site.config.static_dir)."}, Spec: clic.OptionSpec{Key: "static-dir", Kind: clic.StringOption}},
		"theme":       {Meta: clic.OptionMeta{Name: "theme", EnvName: "HYPER_THEME", Description: "Relative path to theme config to use for the build (overwrites site.config.theme)."}, Spec: clic.OptionSpec{Key: "theme", Kind: clic.StringOption}},
		"data-file":   {Meta: clic.OptionMeta{Name: "data-file", Description: "Build-root-relative or absolute path for build data JSON, for example build.json."}, Spec: clic.OptionSpec{Key: "data-file", Kind: clic.StringOption}},
		"log-file":    {Meta: clic.OptionMeta{Name: "log-file", Description: "Build-root-relative or absolute path for build logs, for example build.log."}, Spec: clic.OptionSpec{Key: "log-file", Kind: clic.StringOption}},
	}
}
func serverHandler(execution *clic.Execution) error {
	options := execution.Options
	port, err := strconv.Atoi(options["port"].StringValue())
	if err != nil {
		return fmt.Errorf("invalid port: %s", options["port"].StringValue())
	}
	base, err := url.Parse(fmt.Sprintf("http://%s:%d/", options["host"].StringValue(), port))
	if err != nil {
		return fmt.Errorf("invalid server URL: %v", err)
	}
	interval, err := strconv.Atoi(options["interval"].StringValue())
	if err != nil {
		interval = 2
	}
	config := htxt.BuildConfig{WebsiteConfig: htxt.WebsiteConfig{ContentDirectory: options["content-dir"].StringValue(), DataDirectory: options["data-dir"].StringValue(), StaticDirectory: options["static-dir"].StringValue(), ThemePath: options["theme"].StringValue()}, BaseURL: base.String(), DataFile: options["data-file"].StringValue(), LogFile: options["log-file"].StringValue(), Verbose: true}
	uri, err := href.Parse(options["config"].StringValue())
	if err != nil {
		return err
	}
	website := &htxt.Website{URI: uri}
	if err = website.Load(false, &config.WebsiteConfig, nil); err != nil {
		return err
	}
	build := &htxt.Build{URI: &url.URL{Path: filepath.Join(website.BuildsDirPath(), options["id"].StringValue())}, ID: options["id"].StringValue(), Website: website, Config: config, Preview: true}
	if err = build.Execute(execution.Context); err != nil {
		return err
	}
	address := fmt.Sprintf("%s:%d", options["host"].StringValue(), port)
	if _, err = fmt.Fprintf(execution.Stdout, "Listening on %s...\n", address); err != nil {
		return err
	}
	server := &http.Server{Addr: address, Handler: http.FileServer(http.Dir(build.URI.Path)), ReadHeaderTimeout: time.Duration(interval) * time.Second}
	return server.ListenAndServe()
}
