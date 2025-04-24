package main

import (
	"fmt"
	"github/jordani-alpuche/test2/internal/data"
	"github/jordani-alpuche/test2/internal/validator"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
)

//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
/*******************************************************************************************************************************************************
*																Product Handlers																	   *
********************************************************************************************************************************************************/

func (app *application) GETProducts(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.CSRFField = csrf.TemplateField(r)
	data.CurrentPage="/products"

	products, err := app.productInfo.GET(0)
	if err != nil {
		app.logger.Error("failed to fetch Products", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	
	data.Product = products // 👈 Add Products to TemplateData


	err = app.render(w, r, http.StatusOK, "productlist.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render Products page", "template", "productlist.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) createForm(w http.ResponseWriter, r *http.Request) {

	data := NewTemplateData()
	data.CSRFField = csrf.TemplateField(r)
	data.CurrentPage="/product"
	data.CurrentPageType="create"

		// --- Fetch Brands ---
		brands, err := app.brandInfo.GET(0) // Assuming GET(0) fetches all brands
		if err != nil {
			app.logger.Error("failed to fetch brands for product form", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Add brands to the template data
		data.Brand = brands // Make sure TemplateData has a 'Brand' field
	
		// --- Fetch Categories ---
		categories, err := app.categoryInfo.GET(0) // Assuming GET(0) fetches all categories
		if err != nil {
			app.logger.Error("failed to fetch categories for product form", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Add categories to the template data
		data.Category = categories // Make sure TemplateData has a 'Category' field


	err = app.render(w, r, http.StatusOK, "addupdateproduct.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render products page", "template", "addupdateproduct.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}



func (app *application) createProducts(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	productName := r.PostForm.Get("ProductName")
	productDescription := r.PostForm.Get("ProductDescription")
	productPrice := r.PostForm.Get("ProductPrice")
	productCategoryID := r.PostForm.Get("ProductCategoryID")
	productBrandID := r.PostForm.Get("ProductBrandID")
	productQTY := r.PostForm.Get("ProductQTY")
	productStatus := r.PostForm.Get("ProductStatus")
	productPurchasedFrom := r.PostForm.Get("ProductPurchasedFrom")

	price, _ := strconv.ParseFloat(productPrice, 64) // or use your `number()` function
	categoryID, _ := strconv.Atoi(productCategoryID)
	brandID, _ := strconv.Atoi(productBrandID)
	qty, _ := strconv.Atoi(productQTY)

	// Call the AI model to generate a product tag
	 productTag, err := app.generateProductTag(productName, productDescription)
	if err != nil {
		app.logger.Error("failed to generate product tag", "error", err)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Printf("\nGenerated Product Tag: %s\n", productTag)
	
	product := &data.ProductData{		
		ProductName:  productName,
		ProductDescription:  productDescription,
		ProductPrice:  price,
		ProductCategoryID:  categoryID,
		ProductBrandID:  brandID,
		ProductQTY:  qty,
		ProductStatus:  productStatus,
		ProductPurchasedFrom:  productPurchasedFrom,
		ProductTag: productTag, // Add the generated product tag
	}

	// validate data
	v := validator.NewValidator()
	data.ValidateProduct(v, product)
	// 
	// Check for validation errors
	if !v.ValidData() {
		data := NewTemplateData()
		data.CSRFField = csrf.TemplateField(r)
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"ProductName":    productName,
			"ProductDescription":   productDescription,
			"ProductPrice":   productPrice,
			"ProductCategoryID":   productCategoryID,
			"ProductBrandID":   productBrandID,
			"ProductQTY":   productQTY,
			"ProductStatus":   productStatus,
			"ProductPurchasedFrom":   productPurchasedFrom,
		}

		        // --- RE-FETCH Brands and Categories ---
				brands, err := app.brandInfo.GET(0)
				if err != nil {
					app.logger.Error("failed to fetch brands for product form re-render", "error", err)
					// Decide how critical this is. If you can't show the form without brands, return 500.
					// If you can show it (with an empty dropdown), just log and continue.
					// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					// return
				} else {
					data.Brand = brands // Add brands to the template data
				}
		
				categories, err := app.categoryInfo.GET(0)
				if err != nil {
					app.logger.Error("failed to fetch categories for product form re-render", "error", err)
					 // Decide how critical this is.
					// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					// return
				} else {
					 data.Category = categories // Add categories to the template data
				}
				// --- End RE-FETCH ---
		

		err = app.render(w, r, http.StatusUnprocessableEntity, "addupdateproduct.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render Product Form", "template", "addupdateproduct.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.productInfo.POST(product)

	if err != nil {
		app.logger.Error("failed to insert product", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func (app *application) productItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // 👈 get ID from URL path like /product-item/25
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	products, err := app.productInfo.GET(id)
	
	if err != nil || len(products) == 0 {
		app.logger.Error("product not found", "id", id, "error", err)
	
		data := NewTemplateData()
		data.CSRFField = csrf.TemplateField(r)
		data.FormData = map[string]string{
			"Message": "The product you're looking for doesn't exist.",
		}
		
		// Render custom 404 page
		err = app.render(w, r, http.StatusNotFound, "error-404.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render 404 page", "error", err)
			http.Error(w, "Page not found", http.StatusNotFound)
		}
		return
	}

	product := products[0]


	data := NewTemplateData()
	data.CSRFField = csrf.TemplateField(r)
	data.CurrentPage="/products"
	data.FormData = map[string]string{
		"ProductName":         product.ProductName,
		"ProductDescription":  product.ProductDescription,
		"ProductPrice":        fmt.Sprintf("%.2f", product.ProductPrice),
		"ProductCategoryID":   strconv.Itoa(product.ProductCategoryID),
		"ProductBrandID":      strconv.Itoa(product.ProductBrandID),
		"ProductQTY":          strconv.Itoa(product.ProductQTY),
		"ProductStatus":       product.ProductStatus,
		"ProductPurchasedFrom": product.ProductPurchasedFrom,
	}

	// --- RE-FETCH Brands and Categories ---
	brands, err := app.brandInfo.GET(product.ProductBrandID)
	if err != nil {
		app.logger.Error("failed to fetch brands for product form re-render", "error", err)
		// Decide how critical this is. If you can't show the form without brands, return 500.
		// If you can show it (with an empty dropdown), just log and continue.
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// return
	} else {
		data.Brand = brands // Add brands to the template data
	}

	categories, err := app.categoryInfo.GET(product.ProductCategoryID)
	if err != nil {
		app.logger.Error("failed to fetch categories for product form re-render", "error", err)
			// Decide how critical this is.
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// return
	} else {
			data.Category = categories // Add categories to the template data
	}
	// --- End RE-FETCH ---
			

	fmt.Printf("Data: %+v\n", data)

	err = app.render(w, r, http.StatusOK, "product-details.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render viewer", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) updateForm(w http.ResponseWriter, r *http.Request){
	idStr := r.PathValue("id") 
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	products, err := app.productInfo.GET(id)
	
	if err != nil || len(products) == 0 {
		app.logger.Error("product not found", "id", id, "error", err)
	
		data := NewTemplateData()
		data.CSRFField = csrf.TemplateField(r)
		data.FormData = map[string]string{
			"Message": "The product you're looking for doesn't exist.",
		}
		
		err = app.render(w, r, http.StatusNotFound, "error-404.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render 404 page", "error", err)
			http.Error(w, "Page not found", http.StatusNotFound)
		}
		return
	}

	product := products[0]
	data := NewTemplateData()
	data.CSRFField = csrf.TemplateField(r)
	data.CurrentPage="/products"
	data.CurrentPageType = "update"
	data.FormData = map[string]string{
		"ID":                 strconv.FormatInt(product.ID, 10), 
		"ProductName":         product.ProductName,
		"ProductDescription":  product.ProductDescription,
		"ProductPrice":        fmt.Sprintf("%.2f", product.ProductPrice),
		"ProductCategoryID":   strconv.Itoa(product.ProductCategoryID),
		"ProductBrandID":      strconv.Itoa(product.ProductBrandID),
		"ProductQTY":          strconv.Itoa(product.ProductQTY),	
		"ProductStatus":       product.ProductStatus,
		"ProductPurchasedFrom": product.ProductPurchasedFrom,

}

		brands, err := app.brandInfo.GET(0) // Assuming GET(0) fetches all brands
		if err != nil {
			app.logger.Error("failed to fetch brands for product form", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Add brands to the template data
		data.Brand = brands // Make sure TemplateData has a 'Brand' field

		// --- Fetch Categories ---
		categories, err := app.categoryInfo.GET(0) // Assuming GET(0) fetches all categories
		if err != nil {
			app.logger.Error("failed to fetch categories for product form", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Add categories to the template data
		data.Category = categories // Make sure TemplateData has a 'Category' field

	err = app.render(w, r, http.StatusOK, "addupdateproduct.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render viewer", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}



func (app *application) updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // 👈 get ID from URL path like /product-item/25
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	err = r.ParseForm()
	if err != nil {
		// Handle error appropriately - maybe a client error?
		app.logger.Error("failed to parse form", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Bad Request", http.StatusBadRequest) // Or InternalServerError depending on context
		return
	}

	productName := r.PostForm.Get("ProductName")
	productDescription := r.PostForm.Get("ProductDescription")
	productPrice := r.PostForm.Get("ProductPrice")
	productCategoryID := r.PostForm.Get("ProductCategoryID")
	productBrandID := r.PostForm.Get("ProductBrandID")
	productQTY := r.PostForm.Get("ProductQTY")
	productStatus := r.PostForm.Get("ProductStatus")
	productPurchasedFrom := r.PostForm.Get("ProductPurchasedFrom")

	price, _ := strconv.ParseFloat(productPrice, 64) // or use your `number()` function
	categoryID, _ := strconv.Atoi(productCategoryID)
	brandID, _ := strconv.Atoi(productBrandID)
	qty, _ := strconv.Atoi(productQTY)

	product := &data.ProductData{		
		ProductName:  productName,
		ProductDescription:  productDescription,
		ProductPrice:  price,
		ProductCategoryID:  categoryID,
		ProductBrandID:  brandID,
		ProductQTY:  qty,
		ProductStatus:  productStatus,
		ProductPurchasedFrom:  productPurchasedFrom,
	}

		// validate data
		v := validator.NewValidator()
		data.ValidateProduct(v, product)
		// 
		// Check for validation errors
		if !v.ValidData() {
			data := NewTemplateData()
			data.CSRFField = csrf.TemplateField(r)
			data.CurrentPage="/products"
			data.CurrentPageType = "update"
			data.FormErrors = v.Errors
			data.FormData = map[string]string{
				"ID":                 idStr,
				"ProductName":    productName,
				"ProductDescription":   productDescription,
				"ProductPrice":   productPrice,
				"ProductCategoryID":   productCategoryID,
				"ProductBrandID":   productBrandID,
				"ProductQTY":   productQTY,
				"ProductStatus":   productStatus,
				"ProductPurchasedFrom":   productPurchasedFrom,
			}


			// --- RE-FETCH Brands and Categories ---
			brands, err := app.brandInfo.GET(0)
			if err != nil {
				app.logger.Error("failed to fetch brands for product form re-render", "error", err)
				// Decide how critical this is. If you can't show the form without brands, return 500.
				// If you can show it (with an empty dropdown), just log and continue.
				// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				// return
			} else {
				data.Brand = brands // Add brands to the template data
			}
	
			categories, err := app.categoryInfo.GET(0)
			if err != nil {
				app.logger.Error("failed to fetch categories for product form re-render", "error", err)
					// Decide how critical this is.
				// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				// return
			} else {
					data.Category = categories // Add categories to the template data
			}
			// --- End RE-FETCH ---
	
			err = app.render(w, r, http.StatusUnprocessableEntity, "addupdateproduct.tmpl", data)
			if err != nil {
				app.logger.Error("failed to render Product Form", "template", "addupdateproduct.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}

	err = app.productInfo.PUT(id, product)

	if err != nil {
		app.logger.Error("failed to update product", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/products", http.StatusSeeOther)	
}

func (app *application) deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // 👈 get ID from URL path like /product-item/25
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	err = app.productInfo.DELETE(id)
	if err != nil {
		app.logger.Error("failed to delete product", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/products", http.StatusSeeOther)	
}