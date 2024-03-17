package handlers

import (
	"fmt"
	"github.com/haodev88/bookings/internal/config"
	"github.com/haodev88/bookings/internal/models"
	"github.com/haodev88/bookings/internal/render"
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
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TempldateData{

	})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request)  {
	stringMap:= make(map[string]string)
	//remoteIp:= m.App.Session.GetString(r.Context() ,"remote_ip")
	//log.Println(remoteIp)

	stringMap["test"] = "Data Model from about page"
	stringMap["name"] = "Hao tran here"
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TempldateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Reservation (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TempldateData{

	})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request)  {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TempldateData{

	})
}

func (m *Repository) Majors (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TempldateData{

	})
}

func (m *Repository) Availabitily (w http.ResponseWriter, r *http.Request) {
	var stringMap = make(map[string]string)
	stringMap["title"] = "Search for Availability"
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TempldateData{
		StringMap: stringMap,
	})
}

func (m *Repository) PostAvailabitily(w http.ResponseWriter, r *http.Request) {
	startDate := r.Form.Get("start")
	endDate   := r.Form.Get("end")
	_,_ = w.Write([]byte(fmt.Sprintf("Start data is %s and end date is %s", startDate, endDate)))
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request)  {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TempldateData{

	})
}