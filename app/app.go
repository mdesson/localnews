package app

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/mdesson/localnews/source"
)

type App struct {
	l           *slog.Logger
	Sources     []*source.Source
	LastUpdated time.Time
}

// Update fetches all sources' articles, as well as lastUpdated.
//
// Will not fail if there is an error fetching, it will log it out. this is to prevent one or two sources from stopping
// the entire update from finishing.
func (app *App) Update() {
	app.l.Info("updating sources")
	// fetch the articles
	var wg sync.WaitGroup
	for _, s := range app.Sources {
		go func() {
			wg.Add(1)
			defer wg.Done()
			if err := s.FetchArticles(); err != nil {
				app.l.Error("error updating source", "source", s.Name, "error", err.Error())
			} else {
				app.l.Info("updated source", "source", s.Name)
			}
		}()
	}
	wg.Wait()
	app.LastUpdated = time.Now()
	app.l.Info("updated sources")
}

// NewApp creates a new app. This includes fetching articles
//
// Note that fetching articles will not cause it to fail, since a single failing source should not cause every update to
// to fail.
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

	// initialize app and update sources
	app := &App{Sources: sources, l: logger}
	app.l.Info("starting update")
	app.Update()
	app.l.Info("finished update")

	return app, nil
}
