package main

import (
	"fmt"
	"github.com/haodev88/bookings/internal/helpers"
	"github.com/justinas/nosurf"
	"log"
	"net/http"
)

func WriteToConsole(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("hit the page")
		next.ServeHTTP(writer, request)
	})
}

// Nosurf add CSRF protection all POST request
func Nosurf(next http.Handler) http.Handler {
	log.Println("call no surf")
	csrfHandle:= nosurf.New(next)
	csrfHandle.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandle
}

// SessionLoad and saves the session every request
func SessionLoad(next http.Handler) http.Handler {
	log.Println("Session Load Ok")
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !helpers.IsAuthenticated(request) {
			session.Put(request.Context(), "error", "Log in first !")
			http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
		}
		next.ServeHTTP(writer, request)
	})
}