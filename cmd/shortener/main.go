package main

import (
	"Ivan-Vorobev/shortener/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	router.Post("/", handler.CreateUrl)
	router.Get("/{slug}", handler.ReturnUrl)

	err := http.ListenAndServe(`:8080`, router)
	if err != nil {
		panic(err)
	}
}
