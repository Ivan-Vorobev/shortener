package model

import (
	"errors"
	"strings"
)

type ShortLink struct {
	value string
}

func NewShortLink(shortLink string) ShortLink {
	return ShortLink{
		value: strings.TrimLeft(shortLink, "/"),
	}
}

func (s ShortLink) GetHash() (string, error) {
	if s.value == "" {
		return "", errors.New("Short link is empty")
	}

	return s.value, nil
}

func (s ShortLink) String() string {
	return "/" + s.value
}

func NewLink(link string) (Link, error) {
	link = strings.TrimSpace(link)

	if link == "" {
		return Link{}, errors.New("link is empty")
	}

	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		return Link{}, errors.New("invalid link")
	}

	return Link{
		value: link,
	}, nil
}

type Link struct {
	value string
}

func (l Link) Get() (string, error) {
	if l.value == "" {
		return "", errors.New("Link is empty")
	}

	return l.value, nil
}
