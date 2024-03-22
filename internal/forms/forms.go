package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

// form creates a custom form struct, embed a url.values object
type Form struct {
	url.Values
	Errors erros
}

// New intitializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		erros(map[string][]string{}),
	}
}

// Has check if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x:= f.Get(field)
	if x == "" {
		f.Errors.Add(field,"This field can not be blank")
		return false
	}
	return true
}

// check valid error
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// check require field with mutiple
func (f *Form) Required(fields ...string)  {
	for _,field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field can not be blank")
		}
	}
}

// check minlength
func (f *Form) MinLength(field string, length int)  {
	x:= f.Get(field)
	if (len(x) < length) {
		f.Errors.Add(field, fmt.Sprintf("This field must be at lastest %d characters", length))
	}
}

func (f *Form) IsEmail(field string) {
	var isEmail bool
	isEmail = govalidator.IsEmail(f.Get(field))
	if !isEmail {
		f.Errors.Add(field, "Invalidate email address")
	}
}