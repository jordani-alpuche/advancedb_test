package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)


	///////Products Routes ///////////

	mux.HandleFunc("GET /product", app.createForm)
	mux.HandleFunc("POST /product", app.createProducts)
	
	mux.HandleFunc("GET /products", app.GETProducts)
	mux.HandleFunc("GET /product-item/{id}", app.productItem)

	mux.HandleFunc(("GET /edit-product/{id}"), app.updateForm)
	mux.HandleFunc("POST /edit-product/{id}", app.updateProduct)
	mux.HandleFunc("GET /delete-product/{id}", app.deleteProduct)
	
	/////////Brands Routes /////////////////////////
	mux.HandleFunc("GET /brand", app.createBrandForm)
	mux.HandleFunc("POST /brand", app.createBrand)
	
	mux.HandleFunc("GET /brands", app.GETBrands)
	mux.HandleFunc("GET /brand-item/{id}", app.brandItem)

	mux.HandleFunc(("GET /edit-brand/{id}"), app.updateBrandForm)
	mux.HandleFunc("POST /edit-brand/{id}", app.updateBrand)
	mux.HandleFunc("GET /delete-brand/{id}", app.deleteBrand)

	/////////Categories Routes /////////////////////////
	mux.HandleFunc("GET /category", app.createCategoryForm)
	mux.HandleFunc("POST /category", app.createCategory)
	
	mux.HandleFunc("GET /categories", app.GETCategories)
	mux.HandleFunc("GET /category-item/{id}", app.categoryItem)

	mux.HandleFunc(("GET /edit-category/{id}"), app.updateCategoryForm)
	mux.HandleFunc("POST /edit-category/{id}", app.updateCategory)
	mux.HandleFunc("GET /delete-category/{id}", app.deleteCategory)


	return app.loggingMiddleware(mux)
}
