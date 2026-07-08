package main

import (
	"Ivan-Vorobev/shortener/internal/handler"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handler.CreateMainHandler())

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
