package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/haodev88/bookings/internal/config"
	"github.com/haodev88/bookings/internal/driver"
	"github.com/haodev88/bookings/internal/forms"
	"github.com/haodev88/bookings/internal/helpers"
	"github.com/haodev88/bookings/internal/models"
	"github.com/haodev88/bookings/internal/render"
	"github.com/haodev88/bookings/internal/repository"
	"github.com/haodev88/bookings/internal/repository/dbrepo"
	"net/http"
	"strconv"
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

func NewHandler(r *Repository)  {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request)  {
	m.DB.AllUsers()
	_= render.Template(w, r, "home.page.tmpl", &models.TempldateData{

	})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request)  {
	stringMap:= make(map[string]string)
	// remoteIp:= m.App.Session.GetString(r.Context() ,"remote_ip")
	//log.Println(remoteIp)

	stringMap["test"] = "Data Model from about page"
	stringMap["name"] = "Hao tran here"
	_= render.Template(w, r, "about.page.tmpl", &models.TempldateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Reservation (w http.ResponseWriter, r *http.Request) {
	// var emptyReservation models.Reservation
	res,ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	room, err:= m.DB.GetRoomById(res.RoomID)
	if err!=nil {
		helpers.ServerError(w, err)
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
	_= render.Template(w, r, "make-reservation.page.tmpl", &models.TempldateData {
		Form:forms.New(nil),
		Data: data,
		StringMap: stringMap,
	})
}

func (m *Repository) PostReservation (w http.ResponseWriter, r *http.Request) {

	reservation, ok:= m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can't get session"))
		return
	}


	err:=r.ParseForm()
	if err!=nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName =  r.Form.Get("first_name")
	reservation.LastName =  r.Form.Get("last_name")
	reservation.Phone =  r.Form.Get("phone")
	reservation.Email =  r.Form.Get("email")

	form:=forms.New(r.PostForm)
	// form.Has("first_name")
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name",3)
	form.IsEmail("email")

	if !form.Valid() {
		data:=make(map[string]interface{})
		data["reservation"] = reservation
		_ = render.Template(w, r, "make-reservation.page.tmpl", &models.TempldateData {
			Form: form,
			Data: data,
		})
		return
	}
	var newReservationId int
	newReservationId,err = m.DB.InsertReservation(reservation)

	if err!=nil {
		helpers.ServerError(w, err)
		return
	}

	// insert room restriction
	restriction:= models.RoomRestriction{
		StartDate: reservation.StartDate,
		EndDate: reservation.EndDate,
		RoomID: reservation.RoomID,
		ReservationID: newReservationId,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary display reservation Summary page
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

	sd:= reservation.StartDate.Format("2006-01-02")
	ed:= reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	_= render.Template(w, r, "reservation-summary.page.tmpl", &models.TempldateData{
		Data: data,
		StringMap: stringMap,
	})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request)  {
	_= render.Template(w, r, "generals.page.tmpl", &models.TempldateData{

	})
}

func (m *Repository) Majors (w http.ResponseWriter, r *http.Request) {
	_= render.Template(w, r, "majors.page.tmpl", &models.TempldateData{

	})
}

func (m *Repository) Availability (w http.ResponseWriter, r *http.Request) {
	var stringMap = make(map[string]string)
	stringMap["title"] = "Search for Availability"
	_= render.Template(w, r, "search-availability.page.tmpl", &models.TempldateData{
		StringMap: stringMap,
	})
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end   := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err!=nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err!=nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
	}

	var data = make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate: endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)
	_= render.Template(w, r, "choose-room.page.tmpl", &models.TempldateData{
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
	roomID, err:= strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res,ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
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
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomID = RoomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request)  {
	_= render.Template(w, r, "contact.page.tmpl", &models.TempldateData{

	})
}