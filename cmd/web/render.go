package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

// func (app *application) render(w http.ResponseWriter, status int, page string, data *TemplateData) error {
// 	// app.logger.Info("rendering template", "template", page, "message", messagevalue)
// 	buf := bufferPool.Get().(*bytes.Buffer)
// 	buf.Reset()
// 	defer bufferPool.Put(buf)

// 	ts, ok := app.templateCache[page]
// 	if !ok {
// 		err := fmt.Errorf("template %s does not exist", page)
// 		app.logger.Error("template does not exist", "template", page, "error", err)
// 		return err
// 	}
// 	err := ts.Execute(buf, data)

// 	if err != nil {
// 		err = fmt.Errorf("failed to render template %s: %w", page, err)
// 		app.logger.Error("failed to render template", "template", page, "error", err)
// 		return err
// 	}

// 	w.WriteHeader(status)
// 	_, err = buf.WriteTo(w)
// 	if err != nil {
// 		err = fmt.Errorf("failed to write template to response: %w", err)
// 		app.logger.Error("failed to write template to response", "error", err)
// 		return err
// 	}

// 	return nil
// }

func (app *application) render(w http.ResponseWriter, status int, page string, data *TemplateData) error {
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufferPool.Put(buf)

    ts, ok := app.templateCache[page]
    if !ok {
        err := fmt.Errorf("template %s does not exist", page)
        app.logger.Error("template does not exist", "template", page, "error", err)
        return err
    }
    
    // Execute the "base" template (defined in layout.tmpl)
    err := ts.ExecuteTemplate(buf, "base", data)
    
    if err != nil {
        err = fmt.Errorf("failed to render template %s: %w", page, err)
        app.logger.Error("failed to render template", "template", page, "error", err)
        return err
    }

    w.WriteHeader(status)
    _, err = buf.WriteTo(w)
    if err != nil {
        err = fmt.Errorf("failed to write template to response: %w", err)
        app.logger.Error("failed to write template to response", "error", err)
        return err
    }

    return nil
}