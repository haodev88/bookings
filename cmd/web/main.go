package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/haodev88/bookings/internal/config"
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
	err:= Run()
	fmt.Println("Running with port", PORT_NUM)
	srv:= &http.Server{
		Addr: PORT_NUM,
		Handler: routes(&app),
	}
	err= srv.ListenAndServe()
	log.Fatal(err)
}

func Run() error {
	// Register gob
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = true

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

	tc,err := render.CreateTemplateCache()
	if err!=nil {
		log.Fatal("can not create template cache")
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false
	app.Session  = session

	/** call handeler **/
	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)
	helpers.NewHelper(&app)

	/** render template **/
	render.NewRenderer(&app)
	return nil
}
