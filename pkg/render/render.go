package render

import (
	"bytes"
	"fmt"
	"github.com/haodev88/bookings/pkg/config"
	"github.com/haodev88/bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig
func NewTemplates(appConfig *config.AppConfig)  {
	app = appConfig
}

func addDefaultData(data *models.TempldateData) *models.TempldateData {
	return data
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data *models.TempldateData)  {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc,_ = CreateTemplateCache()
	}

	t, ok:= tc[tmpl]
	if !ok {
		// log.Fatal(err)
		log.Fatal("could not get template cache")
	}

	data = addDefaultData(data)
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, data)
	_,err := buf.WriteTo(w)
	if err!= nil {
		fmt.Println("Error writing template to browser", err)
	}
}

/** Create template cache **/
func CreateTemplateCache()(map[string]*template.Template, error)  {
	myCache := map[string]*template.Template{}
	pages, err:= filepath.Glob("./templates/*.page.tmpl")
	if err!= nil {
		return myCache, err
	}

	for _, page:=range pages {
		name := filepath.Base(page)
		// fmt.Println("Page is currently", page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err!= nil {
			return myCache, err
		}
		matches,err := filepath.Glob("./templates/*.layout.tmpl")
		if err!= nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts,err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err!= nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}