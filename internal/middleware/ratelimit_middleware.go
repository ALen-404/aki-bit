package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

type RateLimitMiddleware struct {
}

func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{}
}

func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return httprate.Limit(
		10,
		10*time.Second,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "some specific response here", http.StatusTooManyRequests)
		}),
	)(next).(http.HandlerFunc)
}
