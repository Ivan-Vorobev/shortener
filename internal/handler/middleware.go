package handler

import (
	"Ivan-Vorobev/shortener/internal/config"
	"context"
	"net/http"
)

const (
	CtxConfigKey = "config"
)

func NewConfMiddleware(configuration *config.Configuration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), CtxConfigKey, configuration)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
