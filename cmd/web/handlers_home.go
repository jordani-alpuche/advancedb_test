package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.CurrentPage="/"

	err := app.render(w, http.StatusOK, "index.tmpl", data)

	if err != nil {
		app.logger.Error("failed to render home page", "template", "index.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}