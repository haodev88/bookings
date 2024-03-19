package main

import (
	"fmt"
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

	mux.Get("/test", test)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/search-availabitily", handlers.Repo.Availabitily)
	mux.Post("/search-availabitily", handlers.Repo.PostAvailabitily)
	mux.Post("/search-availabitily-json", handlers.Repo.AvailabitilyJson)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/contact", handlers.Repo.Contact)

	// load static file in forder static
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	// end load static file

	return mux
}

func test(w http.ResponseWriter, r *http.Request)  {
	_,_ = w.Write([]byte(fmt.Sprintf("This is test function %s", "test")))
}
