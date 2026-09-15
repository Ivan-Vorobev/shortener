package repository

import (
	"Ivan-Vorobev/shortener/internal/model"
	"bufio"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"sync"
)

const HashSize = 8

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	ErrorNotFound = errors.New("short link not found")
)

func NewMemoryShortLinkRepository(fileStoragePath string) (*ShortLinkRepository, error) {
	shortLinks := make(map[model.Link]model.ShortLink)
	links := make(map[model.ShortLink]model.Link)

	if fileStoragePath != "" {
		file, err := os.OpenFile(fileStoragePath, os.O_RDONLY, os.ModePerm)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				var linkStorage model.LinkStorageRow
				line := scanner.Text()

				if err := json.Unmarshal([]byte(line), &linkStorage); err != nil {
					return nil, err
				}

				shortLinks[linkStorage.OriginalURL] = linkStorage.ShortURL
				links[linkStorage.ShortURL] = linkStorage.OriginalURL
			}
		}
	}

	return &ShortLinkRepository{
		shortLinks:  shortLinks,
		links:       links,
		storagePath: fileStoragePath,
		mutex:       sync.RWMutex{},
	}, nil
}

type ShortLinkRepository struct {
	shortLinks  map[model.Link]model.ShortLink
	links       map[model.ShortLink]model.Link
	storagePath string
	mutex       sync.RWMutex
}

func (s *ShortLinkRepository) Get(shortLink model.ShortLink) (model.Link, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	link, ok := s.links[shortLink]
	if !ok {
		return model.Link{}, ErrorNotFound
	}

	return link, nil
}

func (s *ShortLinkRepository) Create(link model.Link) (shortLink model.ShortLink, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if shortLink, ok := s.shortLinks[link]; ok {
		return shortLink, nil
	}

	for {
		hash, err := generateHash()
		if err != nil {
			return model.ShortLink{}, fmt.Errorf("generateHsh: %w", err)
		}

		shortLink = model.NewShortLink(hash)

		if _, ok := s.links[shortLink]; ok {
			continue
		}

		if err := s.saveToStorage(link, shortLink); err != nil {
			return model.ShortLink{}, err
		}

		s.shortLinks[link] = shortLink
		s.links[shortLink] = link
		break
	}

	return shortLink, nil
}

func (s *ShortLinkRepository) saveToStorage(link model.Link, shortLink model.ShortLink) error {
	file, err := os.OpenFile(s.storagePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error append data to storage (%s): %w", s.storagePath, err)
	}
	defer file.Close()

	linkStorage := model.LinkStorageRow{
		UUID:        strconv.Itoa(len(s.shortLinks) + 1),
		ShortURL:    shortLink,
		OriginalURL: link,
	}
	storageData, err := json.Marshal(linkStorage)
	if err != nil {
		return err
	}

	storageData = append(storageData, '\n')

	if _, err := file.Write(storageData); err != nil {
		return err
	}

	return nil
}

func generateHash() (string, error) {
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
