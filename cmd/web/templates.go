package main

import (
	"html/template"
	"path/filepath"
)

// func newTemplateCache() (map[string]*template.Template, error) {
// 	cache := map[string]*template.Template{}
// 	pages, err := filepath.Glob("./ui/html/*.tmpl")
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, page := range pages {
// 		fileName := filepath.Base(page)

// 		ts, err := template.ParseFiles(page)
// 		if err != nil {
// 			return nil, err
// 		}
// 		cache[fileName] = ts
// 	}

// 	return cache, nil
// }

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
        
        // Skip the layout template when processing individual pages
        if name == "layout.tmpl" {
            continue
        }
        
        // Create a template set with the page name
        ts, err := template.New(name).ParseFiles(layoutPath)
        if err != nil {
            return nil, err
        }
        
        // Add the page template to the set
        ts, err = ts.ParseFiles(page)
        if err != nil {
            return nil, err
        }
        
        // Add template to cache
        cache[name] = ts
    }
    
    return cache, nil
}