package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/mdesson/localnews/source"
)

const (
	PORT = 8080
)

type App struct {
	l           *slog.Logger
	Sources     []*source.Source
	LastUpdated time.Time
	templates   *template.Template
}

func (a *App) Start() {
	// update news every minute
	ticker := time.NewTicker(time.Minute)
	done := make(chan bool)
	defer func() {
		done <- true
	}()

	go func() {
		for {
			select {
			case <-done:
				a.l.Info("stopping update process")
				return
			case <-ticker.C:
				a.Update()
			}
		}
	}()

	// start web server
	http.HandleFunc("/", a.handleIndex)
	a.l.Info("starting server", "port", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil); err != nil {
		a.l.Error("server failed", "error", err)
		os.Exit(1)
	}
}

// Update fetches all sources' articles, as well as lastUpdated.
//
// Will not fail if there is an error fetching, it will log it out. this is to prevent one or two sources from stopping
// the entire update from finishing.
func (a *App) Update() {
	a.l.Info("updating sources")
	// fetch the articles
	var wg sync.WaitGroup
	for _, s := range a.Sources {
		go func() {
			wg.Add(1)
			defer wg.Done()
			if err := s.FetchArticles(); err != nil {
				a.l.Error("error updating source", "source", s.Name, "error", err.Error())
			} else {
				a.l.Info("updated source", "source", s.Name)
			}
		}()
	}
	wg.Wait()
	a.LastUpdated = time.Now()
	a.l.Info("updated sources")
}

// NewApp creates a new app. This includes fetching articles
//
// Note that fetching articles will not cause it to fail, since a single failing source should not cause every update to
// fail.
func NewApp(sourcesFile string) (*App, error) {
	// init logger
	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHandler)

	// load sources from json
	sourcesBytes, err := os.ReadFile(sourcesFile)
	if err != nil {
		return nil, err
	}

	sources := make([]*source.Source, 0)
	if err := json.Unmarshal(sourcesBytes, &sources); err != nil {
		return nil, err
	}
	logger.Info("loaded sources from json", "sourcesFile", sourcesFile)

	// load templates
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	templates := template.Must(
		template.New("").Funcs(funcMap).ParseGlob("templates/*.html"),
	)
	// initialize app and update sources
	app := &App{Sources: sources, l: logger, templates: templates}
	app.l.Info("starting update")
	app.Update()
	app.l.Info("finished update")

	return app, nil
}
