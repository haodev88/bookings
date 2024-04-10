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
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}


var Repo *Repository

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository{
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Repository{
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingsRepo(a),
	}
}

func NewHandler(r *Repository)  {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request)  {
	m.DB.AllUsers()
	_= render.Template(w, r, "home.page.tmpl", &models.TemplateData{

	})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request)  {
	stringMap:= make(map[string]string)
	// remoteIp:= m.App.Session.GetString(r.Context() ,"remote_ip")
	//log.Println(remoteIp)

	stringMap["test"] = "Data Model from about page"
	stringMap["name"] = "Hao tran here"
	_= render.Template(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Reservation (w http.ResponseWriter, r *http.Request) {
	// var emptyReservation models.Reservation
	res,ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err:= m.DB.GetRoomById(res.RoomID)
	if err!=nil {
		m.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)

	sd:= res.StartDate.Format("2006-01-02")
	ed:= res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data:= make(map[string]interface{})
	data["reservation"] = res
	_= render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData {
		Form:forms.New(nil),
		Data: data,
		StringMap: stringMap,
	})
}

func (m *Repository) PostReservation (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}


	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "my own error message", http.StatusSeeOther)
		_= render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	htmlMassage:= fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br />
		Dear %s:, <br />
		This is confirm your reservation from %s to %s
	`,
		reservation.FirstName,
		reservation.StartDate.Format("2006-01-02"),
		reservation.EndDate.Format("2006-01-02"),
	)

	// send notifications - first to guest
	msg:= models.MailData{
		To: reservation.Email,
		From: "me@here.com",
		Subject: "Reservation confirm",
		Content: htmlMassage,
		Template: "basic.html",
	}
	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary display reservation Summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation,ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok {
		// m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't reservation from session")
		http.Redirect(w, r,"/", http.StatusTemporaryRedirect)
	}
	m.App.Session.Remove(r.Context(), "reservation")
	var data map[string]interface{}
	data = make(map[string]interface{})
	data["reservation"] = reservation

	sd:= reservation.StartDate.Format("2006-01-02")
	ed:= reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	_= render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
		StringMap: stringMap,
	})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request)  {
	_= render.Template(w, r, "generals.page.tmpl", &models.TemplateData{

	})
}

func (m *Repository) Majors (w http.ResponseWriter, r *http.Request) {
	_= render.Template(w, r, "majors.page.tmpl", &models.TemplateData{

	})
}

func (m *Repository) Availability (w http.ResponseWriter, r *http.Request) {
	var stringMap = make(map[string]string)
	stringMap["title"] = "Search for Availability"
	_= render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)
	_= render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
	// _,_ = w.Write([]byte(fmt.Sprintf("Start data is %s and end date is %s", start, end)))

}

type jsonResponse struct {
	Ok bool `json:"ok"`
	Message string `json:"message"`
	RoomID string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}

func (m * Repository) AvailabilityJson (w http.ResponseWriter, r *http.Request) {

	sd:= r.Form.Get("start")
	ed:= r.Form.Get("end")

	layout:= "2006-01-02"
	startDate,_ := time.Parse(layout, sd)
	endDate,_ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	avaliable,_:= m.DB.SearchAvailabilityByDatesByRoomId(startDate, endDate, roomID)
	resp:= jsonResponse{
		Ok: avaliable,
		Message: "",
		StartDate: sd,
		EndDate: ed,
		RoomID: strconv.Itoa(roomID),
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

// ChooseRoom display available room
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request)  {
	exploded := strings.Split(r.RequestURI, "/")
	// roomID, err:= strconv.Atoi(chi.URLParam(r, "id"))
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res,ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request)  {
	RoomID,_ := strconv.Atoi(r.URL.Query().Get("id"))
	sd:= r.URL.Query().Get("s")
	ed:= r.URL.Query().Get("e")

	layout:= "2006-01-02"
	startDate,_ :=time.Parse(layout, sd)
	endDate,_ := time.Parse(layout, ed)

	var res models.Reservation
	room, err := m.DB.GetRoomById(RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get room from db!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomID = RoomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request)  {
	_= render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowLogin handlers logging the user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request)  {
	_= m.App.Session.RenewToken(r.Context())
	err:= r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	email:= r.Form.Get("email")
	password:= r.Form.Get("password")

	form:= forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		_= render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id,_,err:= m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// Logout logs a User out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request)  {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request)  {
	_= render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{

	})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request)  {
	_= render.Template(w, r, "contact.page.tmpl", &models.TemplateData{

	})
}