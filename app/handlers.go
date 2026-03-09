package app

import (
	"net/http"
	"slices"
	"strings"

	"github.com/goodsign/monday"
	"github.com/mdesson/localnews/source"
)

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	var selected []string
	if c, err := r.Cookie("sources"); err == nil {
		selected = strings.Split(c.Value, ",")
	} else {
		userLang := source.UserLanguage(r)
		for _, s := range a.Sources {
			selected = append(selected, s.CheckboxIDs(userLang)...)
		}
	}

	data := map[string]any{
		"Strings":  translations(r),
		"Sources":  a.Sources,
		"Selected": selected,
	}

	if err := a.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		a.l.Error("error loading index", "error", err)
	}
}

func (a *App) handleArticles(w http.ResponseWriter, r *http.Request) {
	// get all relevant articles
	var articles []source.Article
	if c, err := r.Cookie("sources"); err == nil {
		// cookie is set, get articles relevant to cookie
		selectedSources := strings.Split(c.Value, ",")
		for _, s := range a.Sources {
			for _, article := range s.Articles {
				if slices.Contains(selectedSources, article.SelectedSourceID) {
					articles = append(articles, article)
				}
			}
		}
	} else {
		for _, s := range a.Sources {
			// not custom languages set, filter on user's language
			userLanguage := source.UserLanguage(r)
			for _, article := range s.Articles {
				if (userLanguage & article.Language) != 0 {
					articles = append(articles, article)
				}
			}
		}
	}

	// sort by date
	slices.SortFunc(articles, func(a, b source.Article) int {
		if a.PublishedParsed == nil || b.PublishedParsed == nil {
			return 0
		}
		return b.PublishedParsed.Compare(*a.PublishedParsed)
	})

	// format date
	layout := monday.DefaultFormatFrCADateTime
	locale := monday.LocaleEnUS
	if source.UserLanguage(r) == source.LanguageEnglish {
		layout = monday.DefaultFormatEnGBDateTime
		locale = monday.LocaleEnGB
	}
	for i, _ := range articles {
		articles[i].Published = monday.Format(*articles[i].PublishedParsed, layout, monday.Locale(locale))
	}

	data := map[string]any{
		"Articles": articles,
	}

	if err := a.templates.ExecuteTemplate(w, "articles.html", data); err != nil {
		a.l.Error("error loading articles", "error", err)
	}
}
