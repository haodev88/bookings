package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/haodev88/bookings/pkg/config"
	"github.com/haodev88/bookings/pkg/handlers"
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

	// load static file in forder static
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	// end load static file

	return mux
}
