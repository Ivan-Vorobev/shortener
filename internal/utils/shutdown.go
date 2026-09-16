package utils

import (
	"errors"
	"fmt"
	"io"
)

type closer struct {
	name   string
	closer io.Closer
}

type Shutdown struct {
	closers []closer
}

func NewShutdown() *Shutdown {
	return &Shutdown{
		closers: []closer{},
	}
}

func (s *Shutdown) Add(name string, c io.Closer) {
	s.closers = append(s.closers, closer{
		name:   name,
		closer: c,
	})
}

func (s *Shutdown) Close() error {
	var errs []error

	for _, item := range s.closers {
		if err := item.closer.Close(); err != nil {
			errs = append(
				errs,
				fmt.Errorf("%s: %w", item.name, err),
			)
		}
	}

	return errors.Join(errs...)
}
