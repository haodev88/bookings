package forms

type erros map[string][]string

// Add adds an error message for given field
func (e erros) Add(field, message string)  {
	e[field] = append(e[field], message)
}

// Get returns the first error message
func (e erros) Get(field string) string {
	es:=e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
