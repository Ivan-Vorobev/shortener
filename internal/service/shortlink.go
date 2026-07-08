package service

import (
	"Ivan-Vorobev/shortener/internal/model"
	"Ivan-Vorobev/shortener/internal/repository"
)

func NewShortLinkService(repository *repository.ShortLinkRepository) *ShortLinkService {
	return &ShortLinkService{
		repo: repository,
	}
}

type ShortLinkService struct {
	repo *repository.ShortLinkRepository
}

func (s *ShortLinkService) Create(link model.Link) (model.ShortLink, error) {
	return s.repo.Create(link)
}

func (s *ShortLinkService) Get(shortLink model.ShortLink) (model.Link, error) {
	return s.repo.Get(shortLink)
}
