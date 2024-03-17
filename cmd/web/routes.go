package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/haodev88/bookings/internal/config"
	"github.com/haodev88/bookings/internal/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler  {
	mux:= chi.NewRouter()
	mux.Use(middleware.Recoverer)
	// mux.Use(WriteToConsole)
	mux.Use(SessionLoad)
	mux.Use(Nosurf)
	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	mux.Get("/generals-quarters", http.HandlerFunc(handlers.Repo.Generals))
	mux.Get("/majors-suite", http.HandlerFunc(handlers.Repo.Majors))
	mux.Get("/search-availabitily", http.HandlerFunc(handlers.Repo.Availabitily))
	mux.Post("/search-availabitily", http.HandlerFunc(handlers.Repo.PostAvailabitily))
	mux.Get("/make-reservation", http.HandlerFunc(handlers.Repo.Reservation))
	mux.Get("/contact", http.HandlerFunc(handlers.Repo.Contact))

	// load static file in forder static
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	// end load static file

	return mux
}
