package handler

import (
	"Ivan-Vorobev/shortener/internal/config"
	"Ivan-Vorobev/shortener/internal/model"
	"Ivan-Vorobev/shortener/internal/repository"
	"Ivan-Vorobev/shortener/internal/service"

	"github.com/go-chi/chi/v5"
)

func NewRouter(configuration *config.Configuration) *chi.Mux {
	shortLinks := make(map[model.Link]model.ShortLink)
	links := make(map[model.ShortLink]model.Link)

	shortLinkRepository := repository.NewMemoryShortLinkRepository(shortLinks, links)
	shortLinkService := service.NewShortLinkService(shortLinkRepository)
	linkHandler := NewLinkHandler(configuration, shortLinkService)

	router := chi.NewRouter()
	router.Post("/", linkHandler.CreateShortURL)
	router.Get("/{slug}", linkHandler.ReturnFullURL)

	return router
}
