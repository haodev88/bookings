package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/haodev88/bookings/internal/config"
	"github.com/haodev88/bookings/internal/handlers"
	"github.com/haodev88/bookings/internal/models"
	"github.com/haodev88/bookings/internal/render"
	"log"
	"net/http"
	"time"
)

const PORT_NUM =  ":8080"
var app config.AppConfig
var session *scs.SessionManager

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

	/** render template **/
	render.NewTemplates(&app)
	return nil
}