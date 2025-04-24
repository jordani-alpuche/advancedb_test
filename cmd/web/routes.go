package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	// mux.HandleFunc("GET /{$}", app.home)

	// Protected Routes - Apply authentication middleware to a separate ServeMux
	protectedMux := http.NewServeMux()

	///////////////////Admin Routes ///////////////////////

	protectedMux.Handle("GET /user", app.requireRole("admin")(http.HandlerFunc(app.createUserForm)))
	protectedMux.Handle("POST /user", app.requireRole("admin")(http.HandlerFunc(app.createUser)))
	protectedMux.Handle("GET /users", app.requireRole("admin")(http.HandlerFunc(app.GETUsers)))
	protectedMux.Handle(("GET /edit-user/{id}"), app.requireRole("admin")(http.HandlerFunc(app.updateUserForm)))
	protectedMux.Handle("POST /edit-user/{id}", app.requireRole("admin")(http.HandlerFunc(app.updateUser)))
	
	protectedMux.Handle("GET /delete-product/{id}", app.requireRole("admin")(http.HandlerFunc(app.deleteProduct)))
	protectedMux.Handle("GET /delete-user/{id}", app.requireRole("admin")(http.HandlerFunc(app.deleteUser)))
	protectedMux.Handle("GET /delete-brand/{id}", app.requireRole("admin")(http.HandlerFunc(app.deleteBrand)))
	protectedMux.Handle("GET /delete-category/{id}", app.requireRole("admin")(http.HandlerFunc(app.deleteCategory)))



	///////Products Routes ///////////
	protectedMux.HandleFunc("GET /{$}", app.home)
	protectedMux.HandleFunc("GET /product", app.createForm)
	protectedMux.HandleFunc("POST /product", app.createProducts)
	protectedMux.HandleFunc("GET /products", app.GETProducts)
	protectedMux.HandleFunc("GET /product-item/{id}", app.productItem)
	protectedMux.HandleFunc(("GET /edit-product/{id}"), app.updateForm)
	protectedMux.HandleFunc("POST /edit-product/{id}", app.updateProduct)


	/////////Brands Routes /////////////////////////
	protectedMux.HandleFunc("GET /brand", app.createBrandForm)
	protectedMux.HandleFunc("POST /brand", app.createBrand)
	protectedMux.HandleFunc("GET /brands", app.GETBrands)
	protectedMux.HandleFunc("GET /brand-item/{id}", app.brandItem)
	protectedMux.HandleFunc(("GET /edit-brand/{id}"), app.updateBrandForm)
	protectedMux.HandleFunc("POST /edit-brand/{id}", app.updateBrand)

	/////////Categories Routes /////////////////////////
	protectedMux.HandleFunc("GET /category", app.createCategoryForm)
	protectedMux.HandleFunc("POST /category", app.createCategory)
	protectedMux.HandleFunc("GET /categories", app.GETCategories)
	protectedMux.HandleFunc("GET /category-item/{id}", app.categoryItem)
	protectedMux.HandleFunc(("GET /edit-category/{id}"), app.updateCategoryForm)
	protectedMux.HandleFunc("POST /edit-category/{id}", app.updateCategory)



	// Apply the authentication middleware to the protected routes
	mux.Handle("/", app.authenticate(protectedMux))

	///////////Login Routes - These should be accessible without authentication
	mux.HandleFunc("GET /login", app.LoginForm)
	mux.HandleFunc("POST /login", app.login)
	mux.HandleFunc("GET /logout", app.logout)
	mux.HandleFunc("GET /signup", app.signupForm)
	mux.HandleFunc("POST /signup", app.signupUser)
	

	return app.rateLimitMiddleware(app.loggingMiddleware(mux)) // Apply rate limiting

}