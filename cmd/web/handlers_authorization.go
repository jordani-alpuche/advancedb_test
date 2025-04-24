package main

import (
	"net/http"
)

func (app *application) forbidden(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.AlertMessage = "You are not authorized to access this page."
	data.AlertType = "alert-danger"
	app.render(w, r, http.StatusForbidden, "error-403.tmpl", data) // Create an error template
}
