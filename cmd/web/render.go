package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/csrf"
)

var bufferPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

// func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data *TemplateData) error {
// 	// Get a buffer from the pool
// 	buf := bufferPool.Get().(*bytes.Buffer)
// 	buf.Reset()
// 	// Ensure the buffer is returned to the pool even if errors occur
// 	defer bufferPool.Put(buf)

// 	// Retrieve the appropriate template set from the cache
// 	ts, ok := app.templateCache[page]
// 	if !ok {
// 		err := fmt.Errorf("template %s does not exist", page)
// 		app.logger.Error("template does not exist", "template", page, "error", err.Error())
// 		return err
// 	}

// 	// Inject CSRF token
// 	data.CSRFToken = csrf.Token(r)


// 	var err error
// 	// Conditionally execute the correct template
// 	if page == "signin.tmpl" {
// 		err = ts.ExecuteTemplate(buf, "signin.tmpl", data)
// 	} else {
// 		err = ts.ExecuteTemplate(buf, "base", data)
// 	}

// 	if err != nil {
// 		renderErr := fmt.Errorf("failed to render template %s: %w", page, err)
// 		app.logger.Error("failed to render template", "template", page, "error", renderErr.Error())
// 		return renderErr
// 	}

// 	// Set headers and write response
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	w.WriteHeader(status)

// 	_, err = buf.WriteTo(w)
// 	if err != nil {
// 		writeErr := fmt.Errorf("failed to write template to response: %w", err)
// 		app.logger.Error("failed to write template to response", "error", writeErr.Error())
// 		return writeErr
// 	}

// 	return nil
// }

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data *TemplateData) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("template %s does not exist", page)
		app.logger.Error("template not found", "template", page, "error", err)
		return err
	}

	

	// Inject CSRF token
	data.CSRFField  = csrf.TemplateField(r)


	var err error
	if page == "signin.tmpl" {
		err = ts.ExecuteTemplate(buf, "signin.tmpl", data)
	}else if page == "signup.tmpl" {
		err = ts.ExecuteTemplate(buf, "signup.tmpl", data)
	}else {
		err = ts.ExecuteTemplate(buf, "base", data)
	}

	if err != nil {
		app.logger.Error("template render failed", "template", page, "error", err)
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	_, err = buf.WriteTo(w)
	if err != nil {
		app.logger.Error("write to response failed", "error", err)
		return err
	}

	return nil
}
