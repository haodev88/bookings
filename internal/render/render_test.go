package render

import (
	"github.com/haodev88/bookings/internal/models"
	"net/http"
	"testing"
)


func TestAddDefaultData(t *testing.T)  {
	var td models.TempldateData
	r,err:=getSession()
	if err!=nil {
		t.Error(err)
	}
	session.Put(r.Context(), "Flash", "123")
	result:= addDefaultData(&td, r)
	if result.Flash == "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	var ww myWriter
	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TempldateData{})
	if err != nil {
		t.Error("error writing template to browser", err)
	}

	err = RenderTemplate(&ww, r, "non-existent.page.tmpl", &models.TempldateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}


func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error) {
	r,err := http.NewRequest("GET","/someurl", nil)
	if err!=nil {
		return r, nil
	}
	ctx:=r.Context()
	ctx,_ = session.Load(ctx, r.Header.Get("X-session"))
	r = r.WithContext(ctx)
	return r, nil
}