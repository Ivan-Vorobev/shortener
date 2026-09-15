package model

import (
	"encoding/json"
	"fmt"
)

type LinkStorageRow struct {
	UUID        string    `json:"uuid"`
	ShortURL    ShortLink `json:"short_url"`
	OriginalURL Link      `json:"original_url"`
}

type linkStorageRowJSON struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (l LinkStorageRow) MarshalJSON() ([]byte, error) {
	shortURL, err := l.ShortURL.GetHash()
	if err != nil {
		return nil, fmt.Errorf("short url: %w", err)
	}

	originalURL, err := l.OriginalURL.Get()
	if err != nil {
		return nil, fmt.Errorf("original url: %w", err)
	}

	return json.Marshal(linkStorageRowJSON{
		UUID:        l.UUID,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	})
}

func (l *LinkStorageRow) UnmarshalJSON(data []byte) error {
	var row linkStorageRowJSON
	if err := json.Unmarshal(data, &row); err != nil {
		return err
	}

	shortURL := NewShortLink(row.ShortURL)
	if _, err := shortURL.GetHash(); err != nil {
		return fmt.Errorf("short url: %w", err)
	}

	originalURL, err := NewLink(row.OriginalURL)
	if err != nil {
		return fmt.Errorf("original url: %w", err)
	}

	l.UUID = row.UUID
	l.ShortURL = shortURL
	l.OriginalURL = originalURL

	return nil
}
