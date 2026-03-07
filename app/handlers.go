package app

import "net/http"

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":     "Vaudreil-Soulanges News",
		"Subheader": "Get your local news in one place",
		"Source":    a.Sources[0], // TODO: multi source support
	}

	if err := a.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		a.l.Error("error loading index", "error", err)
	}
}
