package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/haodev88/bookings/internal/config"
	"github.com/haodev88/bookings/internal/driver"
	"github.com/haodev88/bookings/internal/forms"
	"github.com/haodev88/bookings/internal/helpers"
	"github.com/haodev88/bookings/internal/models"
	"github.com/haodev88/bookings/internal/render"
	"github.com/haodev88/bookings/internal/repository"
	"github.com/haodev88/bookings/internal/repository/dbrepo"
	"net/http"
)

type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}


var Repo *Repository

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository{
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandler(r *Repository)  {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request)  {
	m.DB.AllUsers()
	_= render.RenderTemplate(w, r, "home.page.tmpl", &models.TempldateData{

	})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request)  {
	stringMap:= make(map[string]string)
	//remoteIp:= m.App.Session.GetString(r.Context() ,"remote_ip")
	//log.Println(remoteIp)

	stringMap["test"] = "Data Model from about page"
	stringMap["name"] = "Hao tran here"
	_= render.RenderTemplate(w, r, "about.page.tmpl", &models.TempldateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Reservation (w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data:= make(map[string]interface{})
	data["reservation"] = emptyReservation
	_= render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TempldateData {
		Form:forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostReservation (w http.ResponseWriter, r *http.Request) {
	err:=r.ParseForm()
	if err!=nil {
		helpers.ServerError(w, err)
		return
	}

	reservation:= models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName: r.Form.Get("last_name"),
		Phone: r.Form.Get("phone"),
		Email: r.Form.Get("email"),
	}

	form:=forms.New(r.PostForm)
	// form.Has("first_name")
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name",3)
	form.IsEmail("email")

	if !form.Valid() {
		data:=make(map[string]interface{})
		data["reservation"] = reservation
		_ = render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TempldateData {
			Form: form,
			Data: data,
		})
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation,ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't reservation from session")
		http.Redirect(w, r,"/", http.StatusTemporaryRedirect)
	}
	m.App.Session.Remove(r.Context(), "reservation")
	var data map[string]interface{}
	data = make(map[string]interface{})
	data["reservation"] = reservation
	_= render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TempldateData{
		Data: data,
	})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request)  {
	_= render.RenderTemplate(w, r, "generals.page.tmpl", &models.TempldateData{

	})
}

func (m *Repository) Majors (w http.ResponseWriter, r *http.Request) {
	_= render.RenderTemplate(w, r, "majors.page.tmpl", &models.TempldateData{

	})
}

func (m *Repository) Availability (w http.ResponseWriter, r *http.Request) {
	var stringMap = make(map[string]string)
	stringMap["title"] = "Search for Availability"
	_= render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TempldateData{
		StringMap: stringMap,
	})
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	startDate := r.Form.Get("start")
	endDate   := r.Form.Get("end")
	_,_ = w.Write([]byte(fmt.Sprintf("Start data is %s and end date is %s", startDate, endDate)))
}

type jsonResponse struct {
	Ok bool `json:"ok"`
	Message string `json:"message"`
}

func (m * Repository) AvailabilityJson (w http.ResponseWriter, r *http.Request) {
	resp:= jsonResponse{
		Ok: false,
		Message: "AvailabitilyJson",
	}
	out,err := json.MarshalIndent(resp, "", "")
	if err!=nil {
		helpers.ServerError(w, err)
		return
	}
	// log.Println(string(out))
	w.Header().Set("Content-Type","application/json")
	_,_ = w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request)  {
	_= render.RenderTemplate(w, r, "contact.page.tmpl", &models.TempldateData{

	})
}
