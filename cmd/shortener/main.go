package main

import (
	"Ivan-Vorobev/shortener/internal/config"
	"Ivan-Vorobev/shortener/internal/handler"
	"Ivan-Vorobev/shortener/internal/logger"
	"Ivan-Vorobev/shortener/internal/utils"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	shutdown := utils.NewShutdown()
	conf := config.LoadConfiguration()
	logger, err := logger.NewLogger()

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		shutdown.Close()
		logger.Sync()
	}()

	router, err := handler.NewRouter(conf, logger, shutdown)

	if err != nil {
		logger.Fatal("Failed to create router", zap.Error(err))
	}

	logger.Info(
		"Starting server",
		zap.String("addr", conf.ServerAddress),
	)

	err = http.ListenAndServe(conf.ServerAddress, router)
	if err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
