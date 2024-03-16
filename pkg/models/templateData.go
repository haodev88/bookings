package models

type TempldateData struct {
	StringMap map[string]string
	intMap map[string]int
	FloatMap map[string]float32
	Data map[string]interface{}
	CSRFToken string
	Flash string
	Warning string
	Error string
}
