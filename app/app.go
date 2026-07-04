package app

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/mdesson/localnews/source"
	"github.com/pemistahl/lingua-go"
)

const (
	DEFAULT_PORT int = 8080
)

type App struct {
	l                *slog.Logger
	Port             int
	UserAgent        string
	Sources          []*source.Source
	LastUpdated      time.Time
	templates        *template.Template
	languageDetector lingua.LanguageDetector
}

// Start kicks off a background process that refreshes the article list every minute and starts the web server.
func (a *App) Start(staticFolder embed.FS) {
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

	// config web server
	http.HandleFunc("/", a.handleIndex)
	http.HandleFunc("/articles", a.handleArticles)
	http.Handle("/static/", http.FileServer(http.FS(staticFolder)))

	// start web server
	a.l.Info("starting server", "port", a.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.Port), nil); err != nil {
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
			// Fetch the article, log out errors
			err := s.FetchArticles(a.languageDetector, a.UserAgent)

			if err != nil {
				if errors.Is(err, source.ErrorNoFeedURL) {
					a.l.Warn("source has no feed", "source", s.Name)
				}
				a.l.Error("error updating source", "source", s.Name, "error", err.Error())
			} else {
				if len(s.Articles) == 0 {
					a.l.Warn("empty feed", "source", s.Name)
				}
				a.l.Info("updated source", "source", s.Name)
			}

			// Set the error status on it so it won't display on frontend
			s.ErrorOnUpdate = err != nil
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
func NewApp(sourcesBytes []byte, templatesFolder embed.FS) (*App, error) {
	// init logger
	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHandler)

	// select port or fall back to HTTPS default
	port := DEFAULT_PORT
	if portStr, ok := os.LookupEnv("PORT"); ok {
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
		port = portInt
	}

	// set user agent
	userAgent := ""
	if ua, ok := os.LookupEnv("USER_AGENT"); ok {
		userAgent = ua
		logger.Info("setting custom user agent", "userAgent", userAgent)
	}

	// load sources from embedded json files
	sources := make([]*source.Source, 0)
	if err := json.Unmarshal(sourcesBytes, &sources); err != nil {
		return nil, err
	}
	logger.Info("loaded sources from embedded json")

	// load templates
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"isEnglish": func(language source.Language) bool {
			return (language & source.LanguageEnglish) != 0
		},
		"isFrench": func(language source.Language) bool {
			return (language & source.LanguageFrench) != 0
		},
		"isBilingual": func(language source.Language) bool {
			return language == (source.LanguageFrench | source.LanguageEnglish)
		},
		"isSelected": func(sourceID string, selectedSources []string) bool {
			return slices.Contains(selectedSources, sourceID)
		},
	}

	templates := template.Must(
		template.New("").Funcs(funcMap).ParseFS(templatesFolder, "templates/*.html"),
	)
	// init language detector
	detector := lingua.NewLanguageDetectorBuilder().FromLanguages(lingua.English, lingua.French).WithPreloadedLanguageModels().Build()

	// initialize app and update sources
	app := &App{Port: port, UserAgent: userAgent, Sources: sources, l: logger, templates: templates, languageDetector: detector}
	app.l.Info("starting update")
	app.Update()
	app.l.Info("finished update")

	return app, nil
}
