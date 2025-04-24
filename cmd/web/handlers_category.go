package main

import (
	_ "fmt"
	"github/jordani-alpuche/test2/internal/data"
	"github/jordani-alpuche/test2/internal/validator"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
)

/*******************************************************************************************************************************************************
*																category Handlers																	   *
********************************************************************************************************************************************************/

func (app *application) GETCategories(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.CurrentPage="/categories"

	categories, err := app.categoryInfo.GET(0)
	if err != nil {
		app.logger.Error("failed to fetch categories", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	
	data.Category = categories 
	data.CSRFField = csrf.TemplateField(r)


	err = app.render(w, r, http.StatusOK, "categorylist.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render category page", "template", "categorylist.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) createCategoryForm(w http.ResponseWriter, r *http.Request) {

	data := NewTemplateData()
	data.CurrentPage="/category"
	data.CurrentPageType="create"
	data.CSRFField = csrf.TemplateField(r)

	err:= app.render(w, r, http.StatusOK, "addupdatecategory.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render category page", "template", "addupdatecategory.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}



func (app *application) createCategory(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	categoryName := r.PostForm.Get("CategoryName")
	categoryDescription := r.PostForm.Get("CategoryDescription")
	categoryCode := r.PostForm.Get("CategoryCode")
	

	category := &data.CategoryData{		
		CategoryName: categoryName,
		CategoryDescription: categoryDescription,
		CategoryCode: categoryCode,
	}

	// validate data
	v := validator.NewValidator()
	data.ValidateCategory(v, category)
	// 
	// Check for validation errors
	if !v.ValidData() {
		data := NewTemplateData()
		data.CSRFField = csrf.TemplateField(r)
		data.FormErrors = v.Errors
		data.CurrentPage="/category"
		data.CurrentPageType="create"
		data.FormData = map[string]string{
			"CategoryName":    categoryName,
			"CategoryDescription":   categoryDescription,
			"CategoryCode":    categoryCode,
		}

		err := app.render(w, r, http.StatusUnprocessableEntity, "addupdatecategory.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render Category Form", "template", "addupdatecategory.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.categoryInfo.POST(category)

	if err != nil {
		app.logger.Error("failed to insert category", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func (app *application) categoryItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // 👈 get ID from URL path like /category-item/25
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	categories, err := app.categoryInfo.GET(id)
	
	if err != nil || len(categories) == 0 {
		app.logger.Error("category not found", "id", id, "error", err)
	
		data := NewTemplateData()
		data.CSRFField = csrf.TemplateField(r)
		data.FormData = map[string]string{
			"Message": "The category you're looking for doesn't exist.",
		}
		
		// Render custom 404 page
		err = app.render(w, r, http.StatusNotFound, "error-404.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render 404 page", "error", err)
			http.Error(w, "Page not found", http.StatusNotFound)
		}
		return
	}

	category := categories[0]


	data := NewTemplateData()
	data.CSRFField = csrf.TemplateField(r)
	data.CurrentPage="/categories"
	data.FormData = map[string]string{
		"CategoryName":         category.CategoryName,
		"CategoryDescription":  category.CategoryDescription,
		"CategoryCode":         category.CategoryCode,
	}


	err = app.render(w, r, http.StatusOK, "category-details.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render viewer", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) updateCategoryForm(w http.ResponseWriter, r *http.Request){
	idStr := r.PathValue("id") 
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	categories, err := app.categoryInfo.GET(id)
	
	if err != nil || len(categories) == 0 {
		app.logger.Error("category not found", "id", id, "error", err)
	
		data := NewTemplateData()
		data.CSRFField = csrf.TemplateField(r)
		data.FormData = map[string]string{
			"Message": "The Category you're looking for doesn't exist.",
		}
		
		err = app.render(w, r, http.StatusNotFound, "error-404.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render 404 page", "error", err)
			http.Error(w, "Page not found", http.StatusNotFound)
		}
		return
	}

	category := categories[0]
	data := NewTemplateData()
	data.CurrentPage="/categories"
	data.CurrentPageType="update"
	data.CSRFField = csrf.TemplateField(r)
	data.FormData = map[string]string{
		"ID":                 strconv.FormatInt(category.ID, 10), 
		"CategoryName":         category.CategoryName,
		"CategoryDescription":  category.CategoryDescription,
		"CategoryCode":         category.CategoryCode,

}

	err = app.render(w, r, http.StatusOK, "addupdatecategory.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render viewer", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}



func (app *application) updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // 👈 get ID from URL path like /category-item/25
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

	
	categoryName := r.PostForm.Get("CategoryName")
	categoryDescription := r.PostForm.Get("CategoryDescription")
	categoryCode := r.PostForm.Get("CategoryCode")


	category := &data.CategoryData{		
		CategoryName: categoryName,
		CategoryDescription: categoryDescription,
		CategoryCode: categoryCode,
	}

		// validate data
		v := validator.NewValidator()
		data.ValidateCategory(v, category)
		// 
		// Check for validation errors
		if !v.ValidData() {
			data := NewTemplateData()
			data.CSRFField = csrf.TemplateField(r)
			data.FormErrors = v.Errors
			data.CurrentPage="/categories"
			data.CurrentPageType="update"
			data.FormData = map[string]string{
				"ID":                 idStr,
				"CategoryName":         categoryName,
				"CategoryDescription":  categoryDescription,
				"CategoryCode":         categoryCode,
			}
	
			err := app.render(w, r, http.StatusUnprocessableEntity, "addupdatecategory.tmpl", data)
			if err != nil {
				app.logger.Error("failed to render Category Form", "template", "addupdatecategory.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}

	err = app.categoryInfo.PUT(id, category)

	if err != nil {
		app.logger.Error("failed to update category", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)	
}

func (app *application) deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // 👈 get ID from URL path like /category-item/25
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	err = app.categoryInfo.DELETE(id)
	if err != nil {
		app.logger.Error("failed to delete category", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)	
}