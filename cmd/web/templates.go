package main

import (
	"html/template"
	"path/filepath"
)

func newTemplateCache() (map[string]*template.Template, error) {
    cache := map[string]*template.Template{}
    
    // Define the paths for layout (base) and pages
    layoutPath := "./ui/html/layout.tmpl" 
    pagesPath := "./ui/html/*.tmpl"
    
    // Get all page templates 
    pages, err := filepath.Glob(pagesPath)
    if err != nil {
        return nil, err
    }
    
    // Process each page template
    for _, page := range pages {
        // Get the filename (e.g., "index.tmpl")
        name := filepath.Base(page)

        // fmt.Printf("\nProcessing template: %s\n", name)

        if name == "layout.tmpl" {
            // Skip the layout.tmpl itself
            continue
        }

        // Create a new template set
        var ts *template.Template

        if name == "signin.tmpl" {
            // Special case: Signin page does NOT use layout
            ts, err = template.ParseFiles(page)
            if err != nil {
                return nil, err
            }
        } else {
            // Normal pages: use layout + page
            ts, err = template.New(name).ParseFiles(layoutPath, page)
            if err != nil {
                return nil, err
            }
        } 
    // Add template to cache
        cache[name] = ts
    }
    
    return cache, nil
}