package main

import (
	"Ivan-Vorobev/shortener/internal/config"
	"Ivan-Vorobev/shortener/internal/handler"
	"Ivan-Vorobev/shortener/internal/logger"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	conf := config.LoadConfiguration()

	log, err := logger.NewLogger()

	if err != nil {
		panic(err)
	}

	defer log.Sync()

	router := handler.NewRouter(conf, log)

	log.Info(
		"Starting server",
		zap.String("addr", conf.ServerAddress),
	)

	err = http.ListenAndServe(conf.ServerAddress, router)
	if err != nil {
		panic(err)
	}
}
