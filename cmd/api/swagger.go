package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func mountSwagger(r chi.Router) {
	r.Get("/docs/*", httpSwagger.WrapHandler)
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusFound)
	})
	r.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/doc.json", http.StatusFound)
	})
}
