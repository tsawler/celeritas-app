package main

import (
	"github.com/tsawler/celeritas"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.get("/", a.Handlers.Home)

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))
	
	// routes from celeritas
	a.App.Routes.Mount("/celeritas", celeritas.Routes())

	return a.App.Routes
}
