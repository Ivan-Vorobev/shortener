package handler

import (
	"Ivan-Vorobev/shortener/internal/config"
	"Ivan-Vorobev/shortener/internal/model"
	"Ivan-Vorobev/shortener/internal/repository"
	"Ivan-Vorobev/shortener/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type InURL struct {
	URL string `json:"url"`
}

type OutURL struct {
	Result string `json:"result"`
}

type LinkHandler struct {
	configuration    config.Configuration
	shortLinkService service.ShortLinkService
}

func NewLinkHandler(config *config.Configuration, service *service.ShortLinkService) *LinkHandler {
	return &LinkHandler{
		configuration:    *config,
		shortLinkService: *service,
	}
}

func (l *LinkHandler) ReturnFullURL(res http.ResponseWriter, req *http.Request) {
	shortLink := model.NewShortLink(strings.TrimLeft(req.URL.Path, "/"))
	link, err := l.shortLinkService.Get(shortLink)

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

func (l *LinkHandler) CreateShortURL(res http.ResponseWriter, req *http.Request) {
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

	shortLink, err := l.shortLinkService.Create(link)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("%s%s", l.configuration.BaseURL, shortLink.String())))
}

func (l *LinkHandler) CreateAPIShortURL(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, "Invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	input := InURL{}

	if err := json.Unmarshal(body, &input); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	link, err := model.NewLink(input.URL)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	shortLink, err := l.shortLinkService.Create(link)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	result := OutURL{
		Result: fmt.Sprintf("%s%s", l.configuration.BaseURL, shortLink.String()),
	}

	responseData, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(responseData)
}
