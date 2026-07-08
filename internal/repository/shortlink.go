package repository

import (
	"Ivan-Vorobev/shortener/internal/model"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

const HashSize = 8

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	ErrorNotFound = errors.New("short link not found")
	ErrorCreate   = errors.New("short link could not be created")
)

var (
	shortLinks = make(map[model.Link]model.ShortLink)
	links      = make(map[model.ShortLink]model.Link)
)

func NewShortLinkRepository() *ShortLinkRepository {
	return &ShortLinkRepository{
		shortLinks: shortLinks,
		links:      links,
	}
}

type ShortLinkRepository struct {
	shortLinks map[model.Link]model.ShortLink
	links      map[model.ShortLink]model.Link
}

func (s *ShortLinkRepository) Get(shortLink model.ShortLink) (model.Link, error) {
	link, ok := s.links[shortLink]
	if !ok {
		return model.Link{}, ErrorNotFound
	}

	return link, nil
}

func (s *ShortLinkRepository) Create(link model.Link) (model.ShortLink, error) {
	if shortLink, ok := s.shortLinks[link]; ok {
		return shortLink, nil
	}

	for {
		hash, err := s.generateHash()
		if err != nil {
			return model.ShortLink{}, fmt.Errorf("ShortLink.generateHsh: %w", err)
		}

		shortLink := model.NewShortLink(hash)

		if _, ok := s.links[shortLink]; ok {
			continue
		}

		s.shortLinks[link] = shortLink
		s.links[shortLink] = link
		break
	}

	if res, ok := s.shortLinks[link]; ok {
		return res, nil
	}

	return model.ShortLink{}, ErrorCreate
}

func (s *ShortLinkRepository) generateHash() (string, error) {
	length := HashSize
	code := make([]byte, length)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}

		code[i] = alphabet[n.Int64()]
	}

	return string(code), nil
}
