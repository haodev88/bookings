package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/haodev88/bookings/internal/config"
	"testing"
)

func TestRoutes(t *testing.T)  {
	var app config.AppConfig
	mux:= routes(&app)
	switch v:= mux.(type) {
		case *chi.Mux:
			// do something here
		default:
			t.Error(fmt.Sprintf("Type is not *chi.mux, type is %t", v))
	}
}