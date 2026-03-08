package app

import (
	"net/http"

	"github.com/mdesson/localnews/source"
)

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Strings": translations(r),
		"Sources": a.Sources,
	}

	if err := a.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		a.l.Error("error loading index", "error", err)
	}
}

func (a *App) handleArticles(w http.ResponseWriter, r *http.Request) {
	// get user language
	userLanguage := source.UserLanguage(r)

	// get all relevant articles
	var articles []source.Article
	for _, s := range a.Sources {
		for _, article := range s.Articles {
			// TODO: This conditional will remove stuff that might be customer selected
			if userLanguage == article.Language {
				articles = append(articles, article)
			}
		}
	}

	// TODO: sort by date. go into debugger and see what's set on the gofeed item for each source in a.Sources

	data := map[string]any{
		"Articles": articles,
	}

	if err := a.templates.ExecuteTemplate(w, "articles.html", data); err != nil {
		a.l.Error("error loading articles", "error", err)
	}
}
