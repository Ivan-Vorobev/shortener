package repository

import (
	"Ivan-Vorobev/shortener/internal/model"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func NewFileShortLinkRepository(fileStoragePath string) (*FileShortLinkRepository, error) {
	inMemoryShortLinkRepository := NewInMemoryShortLinkRepository()

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

				inMemoryShortLinkRepository.shortLinks[linkStorage.OriginalURL] = linkStorage.ShortURL
				inMemoryShortLinkRepository.links[linkStorage.ShortURL] = linkStorage.OriginalURL
			}
		}
	}

	file, err := os.OpenFile(fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error append data to storage (%s): %w", fileStoragePath, err)
	}

	return &FileShortLinkRepository{
		InMemoryShortLinkRepository: inMemoryShortLinkRepository,
		storage:                     file,
	}, nil
}

type FileShortLinkRepository struct {
	*InMemoryShortLinkRepository
	storage *os.File
}

func (s *FileShortLinkRepository) Create(link model.Link) (shortLink model.ShortLink, err error) {
	shortLink, err = s.InMemoryShortLinkRepository.Create(link)

	if err != nil {
		return shortLink, err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if link, ok := s.links[shortLink]; ok {
		err = s.saveToStorage(link, shortLink)
	}

	return shortLink, err
}

func (s *FileShortLinkRepository) Close() (err error) {
	return s.storage.Close()
}

func (s *FileShortLinkRepository) saveToStorage(link model.Link, shortLink model.ShortLink) error {
	linkStorage := model.LinkStorageRow{
		UUID:        strconv.Itoa(len(s.InMemoryShortLinkRepository.shortLinks) + 1),
		ShortURL:    shortLink,
		OriginalURL: link,
	}
	storageData, err := json.Marshal(linkStorage)
	if err != nil {
		return err
	}

	storageData = append(storageData, '\n')

	if _, err := s.storage.Write(storageData); err != nil {
		return err
	}

	return nil
}
