package main

import (
	"Ivan-Vorobev/shortener/internal/config"
	"Ivan-Vorobev/shortener/internal/handler"
	"net/http"
)

func main() {
	conf := config.LoadConfiguration()

	router := handler.NewRouter(conf)

	err := http.ListenAndServe(conf.ServerAddress, router)
	if err != nil {
		panic(err)
	}
}
