package middleware

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerKey string

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.With().Str("method", r.Method).Str("path", r.URL.Path).Logger()
		logger.Info().Msg("request received")
		ctx := context.WithValue(r.Context(), LoggerKey("logger"), logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetLogger(ctx context.Context) *zerolog.Logger {
	logger, ok := ctx.Value(LoggerKey("logger")).(zerolog.Logger)
	if !ok {
		return &zerolog.Logger{}
	}
	return &logger
}
