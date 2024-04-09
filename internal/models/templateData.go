package models

import "github.com/haodev88/bookings/internal/forms"

type TemplateData struct {
	StringMap map[string]string
	intMap map[string]int
	FloatMap map[string]float32
	Data map[string]interface{}
	CSRFToken string
	Flash string
	Warning string
	Error string
	Form *forms.Form
}
