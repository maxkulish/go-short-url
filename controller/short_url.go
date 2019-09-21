package controller

import (
	"go-short-url/cache"
	"net/http"
	"net/url"
	"strings"
)

type ShortURL struct {
	FullURL  *url.URL
	ShortKey string
	Cache    *cache.Cache
}

func NewShortURL(c *cache.Cache) *ShortURL {
	return &ShortURL{
		Cache: c,
	}
}

// GET /short/:shortKey
func (su *ShortURL) HandleShortURL(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	su.ShortKey = strings.TrimPrefix(path, "/short/")

	if bytes := su.Cache.Get(su.ShortKey); bytes == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else {
		http.Redirect(w, r, string(bytes), http.StatusFound)
	}

}
