package handler

import (
	"Ivan-Vorobev/shortener/internal/model"
	"Ivan-Vorobev/shortener/internal/repository"
	"Ivan-Vorobev/shortener/internal/service"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func CreateMainHandler() func(http.ResponseWriter, *http.Request) {
	shortLinkService := service.NewShortLinkService(repository.NewShortLinkRepository())
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			if req.URL.Path != "/" {
				http.Error(res, fmt.Sprintf("Method POST not allow for url: %s", req.URL.Path), http.StatusBadRequest)
				return
			}

			createUrl(res, req, shortLinkService)
			return
		}

		if req.Method == http.MethodGet && len(req.URL.Path) > 1 {
			returnUrl(res, req, shortLinkService)
			return
		}

		http.Error(res, "Method not allowed", http.StatusBadRequest)
	}
}

func returnUrl(res http.ResponseWriter, req *http.Request, shortLinkService *service.ShortLinkService) {
	shortLink := model.NewShortLink(strings.TrimLeft(req.URL.Path, "/"))
	link, err := shortLinkService.Get(shortLink)

	if err != nil {
		if errors.Is(err, repository.ErrorNotFound) {
			http.Error(res, "Short link not found", http.StatusNotFound)
			return
		}

		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	linkValue, err := link.Get()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Location", linkValue)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func createUrl(res http.ResponseWriter, req *http.Request, shortLinkService *service.ShortLinkService) {
	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	link, err := model.NewLink(string(body))

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	shortLink, err := shortLinkService.Create(link)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("http://%s%s", req.Host, shortLink.String())))
}
