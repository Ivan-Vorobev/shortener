package handler

import (
	"Ivan-Vorobev/shortener/internal/config"
	"Ivan-Vorobev/shortener/internal/repository"
	"Ivan-Vorobev/shortener/internal/service"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(configuration *config.Configuration, log *zap.Logger) *chi.Mux {
	shortLinkRepository, err := repository.NewMemoryShortLinkRepository(configuration.FileStoragePath)
	if err != nil {
		panic(err)
	}

	shortLinkService := service.NewShortLinkService(shortLinkRepository)
	linkHandler := NewLinkHandler(configuration, shortLinkService)

	router := chi.NewRouter()
	router.Use(LoggingMiddleware(log))
	router.Use(CompressMiddleware)
	router.Post("/", linkHandler.CreateShortURL)
	router.Post("/api/shorten", linkHandler.CreateAPIShortURL)
	router.Get("/{slug}", linkHandler.ReturnFullURL)

	return router
}
