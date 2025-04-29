package main

import (
	"github/jordani-alpuche/test2/internal/data"
	"html/template"
)

type TemplateData struct {
	FormErrors map[string]string
	FormData   map[string]string
	Brand   []data.BrandData 
	Category   []data.CategoryData
	Product	[]data.ProductData
	User	[]data.UsersData
	CurrentPage string
	CurrentPageType string
	PasswordFieldName string
	AlertMessage string // To hold general messages like "Invalid credentials"
	AlertType    string // e.g., "alert-danger", "alert-success"
	CSRFField     template.HTML

	// New fields for counts
	ProductCount  int
	CategoryCount int
	BrandCount    int
	UserCount    int
	

	
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		FormErrors: map[string]string{},
		FormData:   map[string]string{},
	}
}