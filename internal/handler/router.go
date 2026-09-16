package handler

import (
	"Ivan-Vorobev/shortener/internal/config"
	"Ivan-Vorobev/shortener/internal/repository"
	"Ivan-Vorobev/shortener/internal/service"
	"Ivan-Vorobev/shortener/internal/utils"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(configuration *config.Configuration, log *zap.Logger, shutdown *utils.Shutdown) (*chi.Mux, error) {
	shortLinkRepository, err := repository.NewFileShortLinkRepository(configuration.FileStoragePath)

	if err != nil {
		return nil, err
	}

	shutdown.Add("FileShortLinkRepository", shortLinkRepository)
	shortLinkService := service.NewShortLinkService(shortLinkRepository)
	linkHandler := NewLinkHandler(configuration, shortLinkService)

	router := chi.NewRouter()
	router.Use(LoggingMiddleware(log))
	router.Use(CompressMiddleware)
	router.Post("/", linkHandler.CreateShortURL)
	router.Post("/api/shorten", linkHandler.CreateAPIShortURL)
	router.Get("/{slug}", linkHandler.ReturnFullURL)

	return router, nil
}
