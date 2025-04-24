package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	// data.CSRFField = template.HTML(csrf.TemplateField(r)) 
	data.CurrentPage="/"
	// data.CSRFToken = csrf.Token(r)
	

	brandCount,err := app.brandInfo.CountAllBrands()
	if err != nil {
		app.logger.Error("failed to get brand count", "error", err, "url", r.URL.Path, "method", r.Method)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	categoryCount,err := app.categoryInfo.CountAllCategories()
	if err != nil {
		app.logger.Error("failed to get category count", "error", err, "url", r.URL.Path, "method", r.Method)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	productCount,err := app.productInfo.CountAllProducts()
	if err != nil {
		app.logger.Error("failed to get product count", "error", err, "url", r.URL.Path, "method", r.Method)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userCount,err := app.userInfo.CountAllActiveUsers()
	if err != nil {
		app.logger.Error("failed to get user count", "error", err, "url", r.URL.Path, "method", r.Method)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data.UserCount = userCount
	data.ProductCount = productCount
	data.CategoryCount = categoryCount
	data.BrandCount = brandCount

	err = app.render(w,r, http.StatusOK, "index.tmpl", data)

	if err != nil {
		app.logger.Error("failed to render home page", "template", "index.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}