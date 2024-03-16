package handlers

import (
	"github.com/haodev88/bookings/pkg/config"
	"github.com/haodev88/bookings/pkg/models"
	"github.com/haodev88/bookings/pkg/render"
	"net/http"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository{
	return &Repository{
		App: a,
	}
}

func NewHandler(r *Repository)  {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request)  {
	render.RenderTemplate(w, "home.page.tmpl", &models.TempldateData{

	})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request)  {
	stringMap:= make(map[string]string)
	stringMap["test"] = "Data Model from about page"
	stringMap["name"] = "Hao tran here"
	render.RenderTemplate(w, "about.page.tmpl", &models.TempldateData{
		StringMap: stringMap,
	})
}
