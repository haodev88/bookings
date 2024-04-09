package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/haodev88/bookings/internal/config"
	"github.com/haodev88/bookings/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var app *config.AppConfig
var pathToTemplates = "./templates"

func NewRenderer(appConfig *config.AppConfig)  {
	app = appConfig
}

func addDefaultData(data *models.TemplateData, r *http.Request) *models.TemplateData {
	data.Flash     = app.Session.PopString(r.Context(), "flash")
	data.Error     = app.Session.PopString(r.Context(), "error")
	data.Warning   = app.Session.PopString(r.Context(), "warning")
	data.CSRFToken = nosurf.Token(r)
	return data
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, data *models.TemplateData) error {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc,_ = CreateTemplateCache()
	}

	t, ok:= tc[tmpl]
	if !ok {
		// log.Fatal(err)
		// log.Fatal("could not get template cache")
		return errors.New("Can't get template cache")
	}

	data = addDefaultData(data, r)
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, data)
	_,err := buf.WriteTo(w)
	if err!= nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}
	return nil
}

/** Create template cache **/
func CreateTemplateCache()(map[string]*template.Template, error)  {
	myCache := map[string]*template.Template{}
	pages, err:= filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
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
		matches,err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err!= nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts,err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err!= nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}