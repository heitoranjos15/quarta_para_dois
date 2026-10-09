package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v4/middleware"
	"golang.org/x/net/http2/h2c"
)

func Logger(next http.Handler) http.Handler {
	return middleware.Logger(next)
}

func Recoverer(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}

func RealIP(next http.Handler) http.Handler {
	return middleware.RealIP(next)
}

func RequestID(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return middleware.Timeout(timeout)
}

func Compress(level int) func(http.Handler) http.Handler {
	return middleware.Compress(level)
}

func GetHeader(next http.Handler) http.Handler {
	return middleware.GetHead(next)
}

func StripSlashes(next http.Handler) http.Handler {
	return middleware.StripSlashes(next)
}

func Heartbeat(path string) func(http.Handler) http.Handler {
	return middleware.Heartbeat(path)
}

func SetHeader(key, value string) func(http.Handler) http.Handler {
	return middleware.SetHeader(key, value)
}

func AllowContentType(types ...string) func(http.Handler) http.Handler {
	return middleware.AllowContentType(types...)
}

func AllowContentEncoding(types ...string) func(http.Handler) http.Handler {
	return middleware.AllowContentEncoding(types...)
}

func NoCache(next http.Handler) http.Handler {
	return middleware.NoCache(next)
}

func H2C(next http.Handler) http.Handler {
	return h2c.NewHandler(next, nil)
}