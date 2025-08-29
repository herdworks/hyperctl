package cli

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	hypertemplates "github.com/herdworks/hypertemplates-go/lib"
	hypertexting "github.com/herdworks/hypertexting-go/lib"
)

func ServerCommand() *Command {
	id := uuid.New()
	return &Command{
		Config: CommandConfig{
			Name:        "server",
			Description: "Run a local HTTP server.",
			Usage:       "hyperctl server [options] [content]",
			Validator:   ServerCommandValidator,
			Handler:     ServerCommandHandler,
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

func ServerCommandValidator(cmd *Command) (err error) {
	return nil
}

func ServerCommandHandler(cmd *Command) (err error) {
	port, err := strconv.Atoi(cmd.Options.Get("port"))
	if err != nil {
		port = 8080
	}
	server := Server{Hostname: cmd.Options.Get("host"), Port: port}
	server.URL, err = url.Parse(fmt.Sprintf("http://%s:%d/", server.Hostname, server.Port))
	server.ConfigFile = cmd.Options.Get("config")
	server.Interval, err = strconv.Atoi(cmd.Options.Get("interval"))
	if err != nil {
		server.Interval = 2
	}
	server.Duration = time.Duration(server.Interval) * time.Second
	server.WatchedFiles = map[string]fs.FileInfo{}
	website := hypertexting.Website{URI: &url.URL{Path: cmd.Options.Get("config")}}
	server.Build = &hypertexting.HyperBuild{
		URI:     &url.URL{Path: filepath.Join(website.BuildsDir(), cmd.Options.Get("id"))},
		ID:      cmd.Options.Get("id"),
		Website: &website,
		Pages:   &website.Pages,
	}
	server.WebsiteConfig = &hypertexting.WebsiteConfig{
		ContentDirectory: cmd.Options.Get("content-dir"),
		DataDirectory:    cmd.Options.Get("data-dir"),
		StaticDirectory:  cmd.Options.Get("static-dir"),
		ThemeSettings:    cmd.Options.Get("theme"),
	}
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

	ticker := time.NewTicker(server.Duration)
	defer ticker.Stop()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		ever := true
		for ever {
			select {
			case <-ticker.C:
				count, err := ServerCommandWatcher(&server)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if count > 0 {
					var plural string
					if count > 1 {
						plural = "s"
					}
					fmt.Printf("\n%v change%v detected, refreshing build...\n\n", count, plural)
					err = server.Build.Execute(context.TODO(), server.URL, server.WebsiteConfig)
					if err != nil {
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

func ServerCommandWatcher(server *Server) (changecount int, err error) {
	changecount = 0
	now := time.Now()
	dir := filepath.Dir(server.ConfigFile)
	for _, p := range hypertemplates.GetPaths(dir, true) {
		if strings.HasPrefix(p, ".") {
			continue
		}
		if strings.HasPrefix(p, "builds/") {
			continue
		}
		stat, err := os.Stat(p)
		if err != nil {
			return 0, err
		}
		if _, watched := server.WatchedFiles[p]; !watched {
			if changecount == 0 {
				fmt.Println()
			}
			fmt.Printf("file %s added\n", p)
			server.WatchedFiles[p] = stat
			changecount++
			continue
		}
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
