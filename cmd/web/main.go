package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/haodev88/bookings/internal/config"
	"github.com/haodev88/bookings/internal/driver"
	"github.com/haodev88/bookings/internal/handlers"
	"github.com/haodev88/bookings/internal/helpers"
	"github.com/haodev88/bookings/internal/models"
	"github.com/haodev88/bookings/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const PORT_NUM =  ":8080"
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err:= Run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)

	// Listen for template email
	fmt.Println("Listen for sending email")
	listenForMail()

	// Create route
	fmt.Println("Running with port", PORT_NUM)
	srv:= &http.Server{
		Addr: PORT_NUM,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func Run() (*driver.DB, error) {
	// Register gob
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailchan:= make(chan models.MailData)
	app.MailChan = mailchan

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist  = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure   = app.InProduction
	app.Session  = session

	// connect to database
	log.Println("Connecting to database")
	db, err:= driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=root")
	if err != nil {
		log.Fatal("Can not connect to database! Dying...")
	}

	log.Println("Connected to database!")


	tc,err := render.CreateTemplateCache()
	if err!=nil {
		log.Fatal("can not create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	/** call handeler **/
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)
	helpers.NewHelper(&app)

	/** render template **/
	render.NewRenderer(&app)
	return db, nil
}
