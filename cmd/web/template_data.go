package main

import "github/jordani-alpuche/test1/internal/data"

type TemplateData struct {
	FormErrors map[string]string
	FormData   map[string]string
	Brand   []data.BrandData 
	Category   []data.CategoryData
	Product	[]data.ProductData
	CurrentPage string

	
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		FormErrors: map[string]string{},
		FormData:   map[string]string{},
	}
}