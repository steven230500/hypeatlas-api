package http

import (
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	allowed := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowed == "" {
		// por defecto: local y dominio principal
		allowed = "http://localhost:3000,https://hypeatlas.app"
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(allowed, ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	r.Use(c.Handler)

	return r
}
