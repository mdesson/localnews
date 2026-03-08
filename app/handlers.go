package app

import "net/http"

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
	data := map[string]any{
		"Source": a.Sources[0], // TODO: multi source support
	}

	if err := a.templates.ExecuteTemplate(w, "articles.html", data); err != nil {
		a.l.Error("error loading articles", "error", err)
	}
}
