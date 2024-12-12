package middleware

import (
	"net/http"

	cache "github.com/victorspringer/http-cache"
)

type CacheMiddleware struct {
	cacheClient *cache.Client
}

func NewCacheMiddleware(client *cache.Client) *CacheMiddleware {
	return &CacheMiddleware{
		cacheClient: client,
	}
}

func (m *CacheMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return m.cacheClient.Middleware(next).(http.HandlerFunc)
}
