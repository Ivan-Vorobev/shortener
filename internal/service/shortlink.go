package service

import (
	"Ivan-Vorobev/shortener/internal/model"
)

type ShortLinkRepository interface {
	Create(link model.Link) (model.ShortLink, error)
	Get(shortLink model.ShortLink) (model.Link, error)
}

func NewShortLinkService(repo ShortLinkRepository) *ShortLinkService {
	return &ShortLinkService{
		repo: repo,
	}
}

type ShortLinkService struct {
	repo ShortLinkRepository
}

func (s *ShortLinkService) Create(link model.Link) (model.ShortLink, error) {
	return s.repo.Create(link)
}

func (s *ShortLinkService) Get(shortLink model.ShortLink) (model.Link, error) {
	return s.repo.Get(shortLink)
}
