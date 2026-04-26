package dev

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"

	cli "github.com/herdworks/cli-go"
	hypertemplates "github.com/herdworks/hypertemplates-go"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
)

func ServerCommand() *cli.Command {
	id := uuid.New()
	return &cli.Command{
		Config: cli.CommandConfig{
			Name:        "server",
			Description: "Live-reloading development server.",
			Usage:       "hyperctl dev server [options] [content]",
			Validator:   ServerCommandValidator,
			Handler:     ServerCommandHandler,
		},
		Options: cli.CommandOptions{
			{
				Name:        "config",
				ShortName:   "c",
				EnvName:     "HYPER_CONFIG",
				Description: "Path to the website configuration file.",
				Required:    true,
				Secret:      false,
			},
			{
				Name:        "id",
				ShortName:   "i",
				EnvName:     "HYPER_BUILD_ID",
				Description: "Build ID.",
				Required:    false,
				Secret:      false,
				Default:     id.String(),
			},
			{
				Name:        "interval",
				ShortName:   "I",
				EnvName:     "HYPER_REFRESH",
				Description: "Live reload interval seconds.",
				Required:    false,
				Secret:      false,
				Default:     "2",
			},
			{
				Name:        "host",
				ShortName:   "H",
				EnvName:     "HYPER_HOST",
				Description: "Development server host.",
				Required:    false,
				Secret:      false,
				Default:     "localhost",
			},
			{
				Name:        "port",
				ShortName:   "P",
				EnvName:     "HYPER_PORT",
				Description: "Development server port.",
				Required:    false,
				Secret:      false,
				Default:     "8080",
			},
			{
				Name:        "content-dir",
				EnvName:     "HYPER_CONTENT_DIR",
				Description: "Relative path to the content directory to use for the build (overwrites site.config.content_dir).",
			},
			{
				Name:        "data-dir",
				EnvName:     "HYPER_DATA_DIR",
				Description: "Relative path to the data directory to use for the build (overwrites site.config.data_dir).",
			},
			{
				Name:        "static-dir",
				EnvName:     "HYPER_STATIC_DIR",
				Description: "Relative path to the static directory to use for the build (overwrites site.config.static_dir).",
			},
			{
				Name:        "theme",
				EnvName:     "HYPER_THEME",
				Description: "Relative path to theme config to use for the build (overwrites site.config.theme).",
			},
		},
	}
}

func ServerCommandValidator(cmd *cli.Command) (err error) {
	return nil
}

// TODO: this thing was only meant to be a temporary prototype; need refactoring
func ServerCommandHandler(cmd *cli.Command) (err error) {
	var server Server = Server{Hostname: cmd.Options.Get("host")}
	if server.Port, err = strconv.Atoi(cmd.Options.Get("port")); err != nil {
		return fmt.Errorf("invalid port: %s", cmd.Options.Get("port"))
	}
	if server.URL, err = url.Parse(fmt.Sprintf("http://%s:%d/", server.Hostname, server.Port)); err != nil {
		return fmt.Errorf("invalid server URL: %v", err)
	}
	if server.Interval, err = strconv.Atoi(cmd.Options.Get("interval")); err != nil {
		server.Interval = 2
	}
	server.ConfigFile = cmd.Options.Get("config")
	server.Duration = time.Duration(server.Interval) * time.Second
	server.WatchedFiles = map[string]fs.FileInfo{}
	server.WebsiteConfig = &hypertexting.WebsiteConfig{
		ContentDirectory: cmd.Options.Get("content-dir"),
		DataDirectory:    cmd.Options.Get("data-dir"),
		StaticDirectory:  cmd.Options.Get("static-dir"),
		ThemePath:        cmd.Options.Get("theme"),
	}
	websiteURI, err := href.Parse(cmd.Options.Get("config"))
	if err != nil {
		return err
	}
	website := &hypertexting.Website{URI: websiteURI}
	if err = website.Load(false, server.WebsiteConfig); err != nil {
		return err
	}
	server.Build = &hypertexting.HyperBuild{
		URI:     &url.URL{Path: filepath.Join(website.BuildsDir(), cmd.Options.Get("id"))},
		ID:      cmd.Options.Get("id"),
		Website: website,
		Pages:   &website.Pages,
	}

	// start the server
	go func() {
		fmt.Printf("Listening on %s...\n", server.Address())
		ss := http.FileServer(http.Dir(server.Build.URI.Path))
		http.Handle("/", ss)
		err = http.ListenAndServe(server.Address(), requestLogger(http.DefaultServeMux))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	// start the file watcher
	var ticker *time.Ticker = time.NewTicker(server.Duration)
	defer ticker.Stop()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		var count int
		var ever bool = true
		for ever {
			select {
			case <-ticker.C:
				if count, err = ServerCommandWatcher(&server, website); err != nil {
					fmt.Println(err)
					continue
				}
				if count > 0 {
					var plural string
					if count > 1 {
						plural = "s"
					}
					fmt.Printf("\n%v change%v detected, refreshing build...\n\n", count, plural)
					server.Build.Website.Loaded = false // force reload of website config
					if err = server.Build.Website.Load(false, server.WebsiteConfig); err != nil {
						fmt.Println(err)
						continue
					}
					if err = server.Build.Execute(context.TODO(), server.URL, nil); err != nil {
						fmt.Println(err)
					}
					fmt.Println()
				}
				continue
			case <-done:
				ticker.Stop()
				fmt.Println("Shutting down...")
				return
			}
		}
	}()
	<-done
	return nil
}

func ServerCommandWatcher(server *Server, website *hypertexting.Website) (changecount int, err error) {
	changecount = 0
	now := time.Now()
	dir := filepath.Dir(server.ConfigFile)
	buildpath := path.Join(dir, website.Config.BuildsDirectory)
	for _, p := range hypertemplates.GetPaths(dir, true) {
		var stat fs.FileInfo
		if strings.HasPrefix(p, ".") || strings.HasPrefix(p, buildpath) {
			continue
		}
		if stat, err = os.Stat(p); err != nil {
			return 0, err
		}
		// new file
		if _, watched := server.WatchedFiles[p]; !watched {
			if changecount == 0 {
				fmt.Println()
			}
			fmt.Printf("file %s added\n", p)
			server.WatchedFiles[p] = stat
			changecount++
			continue
		}
		// modified file
		age := uint64(now.Sub(stat.ModTime()).Milliseconds())
		if age < uint64(server.Interval*1000) {
			if changecount == 0 {
				fmt.Println()
			}
			fmt.Printf("file %s updated\n", p)
			changecount++
			continue
		}
	}
	return changecount, nil
}

type Server struct {
	ConfigFile    string
	WebsiteConfig *hypertexting.WebsiteConfig
	Build         *hypertexting.HyperBuild
	Interval      int
	Duration      time.Duration
	WatchedFiles  map[string]fs.FileInfo
	Hostname      string
	Port          int
	URL           *url.URL
}

func (s *Server) Address() string {
	return fmt.Sprintf("%s:%d", s.Hostname, s.Port)
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%-5s %s (%s)\n", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
