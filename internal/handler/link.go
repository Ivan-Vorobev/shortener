package handler

import (
	"Ivan-Vorobev/shortener/internal/config"
	"Ivan-Vorobev/shortener/internal/model"
	"Ivan-Vorobev/shortener/internal/repository"
	"Ivan-Vorobev/shortener/internal/service"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ReturnUrl(res http.ResponseWriter, req *http.Request) {
	shortLinkService := service.NewShortLinkService()
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

func CreateUrl(res http.ResponseWriter, req *http.Request) {
	conf, ok := req.Context().Value(CtxConfigKey).(*config.Configuration)
	if !ok {
		http.Error(res, "configuration not set", http.StatusInternalServerError)
		return
	}

	shortLinkService := service.NewShortLinkService()
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
	res.Write([]byte(fmt.Sprintf("%s%s", conf.BaseURL, shortLink.String())))
}
