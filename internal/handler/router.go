package handler

import (
	"Ivan-Vorobev/shortener/internal/config"

	"github.com/go-chi/chi/v5"
)

func NewRouter(configuration *config.Configuration) *chi.Mux {
	router := chi.NewRouter()
	router.Use(NewConfMiddleware(configuration))
	router.Post("/", CreateUrl)
	router.Get("/{slug}", ReturnUrl)

	return router
}
