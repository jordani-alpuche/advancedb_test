package main

import (
	_ "fmt"
	"github/jordani-alpuche/test1/internal/data"
	"github/jordani-alpuche/test1/internal/validator"
	"net/http"
	"strconv"
)

/*******************************************************************************************************************************************************
* Brand Handlers                                                                                                                                                                                                                                                                                                                                                                                    *
********************************************************************************************************************************************************/

func (app *application) GETBrands(w http.ResponseWriter, r *http.Request) {
        data := NewTemplateData()
        data.CurrentPage = "/brands" // Set CurrentPage here

        brands, err := app.brandInfo.GET(0)
        if err != nil {
                app.logger.Error("failed to fetch brands", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }

        data.Brand = brands

        err = app.render(w, http.StatusOK, "brandlist.tmpl", data)
        if err != nil {
                app.logger.Error("failed to render brands page", "template", "brandlist.tmpl", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
}

func (app *application) createBrandForm(w http.ResponseWriter, r *http.Request) {
        data := NewTemplateData()
        data.CurrentPage = "/brand" // Set CurrentPage here

        err := app.render(w, http.StatusOK, "addbrand.tmpl", data)
        if err != nil {
                app.logger.Error("failed to render brands page", "template", "addbrand.tmpl", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
}

func (app *application) createBrand(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
                app.logger.Error("failed to parse form", "error", err)
                http.Error(w, "Bad Request", http.StatusBadRequest)
                return
        }

        brandName := r.PostForm.Get("BrandName")
        brandDescription := r.PostForm.Get("BrandDescription")

        brand := &data.BrandData{
                BrandName:        brandName,
                BrandDescription: brandDescription,
        }

        // validate data
        v := validator.NewValidator()
        data.ValidateBrands(v, brand)
        //
        // Check for validation errors
        if !v.ValidData() {
                data := NewTemplateData()
                data.CurrentPage = "/brand" // Set CurrentPage here
                data.FormErrors = v.Errors
                data.FormData = map[string]string{
                        "BrandName":        brandName,
                        "BrandDescription": brandDescription,
                }

                err := app.render(w, http.StatusUnprocessableEntity, "addbrand.tmpl", data)
                if err != nil {
                        app.logger.Error("failed to render brand Form", "template", "addbrand.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
                        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                        return
                }
                return
        }

        err = app.brandInfo.POST(brand)

        if err != nil {
                app.logger.Error("failed to insert brand", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
        http.Redirect(w, r, "/brands", http.StatusSeeOther)
}

func (app *application) brandItem(w http.ResponseWriter, r *http.Request) {
        idStr := r.PathValue("id") // 👈 get ID from URL path like /brand-item/25
        id, err := strconv.Atoi(idStr)
        if err != nil || id < 1 {
                http.NotFound(w, r)
                return
        }

        brands, err := app.brandInfo.GET(id)

        if err != nil || len(brands) == 0 {
                app.logger.Error("brand not found", "id", id, "error", err)

                data := NewTemplateData()
                data.FormData = map[string]string{
                        "Message": "The brand you're looking for doesn't exist.",
                }

                // Render custom 404 page
                err = app.render(w, http.StatusNotFound, "error-404.tmpl", data)
                if err != nil {
                        app.logger.Error("failed to render 404 page", "error", err)
                        http.Error(w, "Page not found", http.StatusNotFound)
                }
                return
        }

        brand := brands[0]

        data := NewTemplateData()
		data.CurrentPage =  "/brands"
        data.FormData = map[string]string{
                "BrandName":        brand.BrandName,
                "BrandDescription": brand.BrandDescription,
        }

        err = app.render(w, http.StatusOK, "brand-details.tmpl", data)
        if err != nil {
                app.logger.Error("failed to render viewer", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
}

func (app *application) updateBrandForm(w http.ResponseWriter, r *http.Request) {
        idStr := r.PathValue("id")
        id, err := strconv.Atoi(idStr)
        if err != nil || id < 1 {
                http.NotFound(w, r)
                return
        }

        brands, err := app.brandInfo.GET(id)

        if err != nil || len(brands) == 0 {
                app.logger.Error("brand not found", "id", id, "error", err)

                data := NewTemplateData()
                data.FormData = map[string]string{
                        "Message": "The Brand you're looking for doesn't exist.",
                }

                err = app.render(w, http.StatusNotFound, "error-404.tmpl", data)
                if err != nil {
                        app.logger.Error("failed to render 404 page", "error", err)
                        http.Error(w, "Page not found", http.StatusNotFound)
                }
                return
        }

        brand := brands[0]
        data := NewTemplateData()
        data.CurrentPage = "/brands" //Set CurrentPage here.
        data.FormData = map[string]string{
                "ID":               strconv.FormatInt(brand.ID, 10),
                "BrandName":        brand.BrandName,
                "BrandDescription": brand.BrandDescription,
        }

        err = app.render(w, http.StatusOK, "editbrand.tmpl", data)
        if err != nil {
                app.logger.Error("failed to render viewer", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
}

func (app *application) updateBrand(w http.ResponseWriter, r *http.Request) {
        idStr := r.PathValue("id") // 👈 get ID from URL path like /brand-item/25
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

        brandName := r.PostForm.Get("BrandName")
        brandDescription := r.PostForm.Get("BrandDescription")

        brand := &data.BrandData{
                BrandName:        brandName,
                BrandDescription: brandDescription,
        }

        // validate data
        v := validator.NewValidator()
        data.ValidateBrands(v, brand)
        //
        // Check for validation errors
        if !v.ValidData() {
                data := NewTemplateData()
                // data.CurrentPage = "/brand" // Set CurrentPage here.
                data.FormErrors = v.Errors
                data.FormData = map[string]string{
                        "ID":               idStr,
                        "BrandName":        brandName,
                        "BrandDescription": brandDescription,
                }

                err := app.render(w, http.StatusUnprocessableEntity, "editbrand.tmpl", data)
                if err != nil {
                        app.logger.Error("failed to render Brand Form", "template", "editbrand.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
                        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                        return
                }
                return
        }

        err = app.brandInfo.PUT(id, brand)

        if err != nil {
                app.logger.Error("failed to update brand", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
        http.Redirect(w, r, "/brands", http.StatusSeeOther)
}

func (app *application) deleteBrand(w http.ResponseWriter, r *http.Request) {
        idStr := r.PathValue("id") // 👈 get ID from URL path like /brand-item/25
        id, err := strconv.Atoi(idStr)
        if err != nil || id < 1 {
                http.NotFound(w, r)
                return
        }

        err = app.brandInfo.DELETE(id)
        if err != nil {
                app.logger.Error("failed to delete brand", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
        http.Redirect(w, r, "/brands", http.StatusSeeOther)
}