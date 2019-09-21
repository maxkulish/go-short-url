package controller

import (
	"fmt"
	"go-short-url/cache"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	length  = 12
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ123456789"
	scheme  = "http"
	domain  = "localhost:3000"
	route   = "short"
)

type UserURL struct {
	RawURL   string
	FullURL  *url.URL
	ShortURL string
	Cache    *cache.Cache
}

func NewUserURL(c *cache.Cache) *UserURL {
	return &UserURL{
		Cache: c,
	}
}

type Alert struct {
	Visible   bool
	Message   string
	FinalLink string
}

// POST /createURL
func (uu *UserURL) CreateURL(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Forbidden HTTP Method", http.StatusMethodNotAllowed)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Fatal("error parsing form", err)
	}

	uu.RawURL = r.PostForm["full_url"][0]

	if err := uu.NormalizeURL(uu.RawURL); err != nil {
		http.Error(w, "invalid URL format", http.StatusInternalServerError)
	}

	if err := uu.CreateShortURL(); err != nil {
		http.Error(w, "error creating short URL", http.StatusInternalServerError)
	}

	resTemplate, err := template.ParseFiles("template/index.gohtml")
	if err != nil {
		panic(err)
	}
	alert := Alert{
		Visible:   true,
		FinalLink: fmt.Sprintf("%s://%s/%s/%s", scheme, domain, route, uu.ShortURL),
	}

	uu.Cache.Add(uu.ShortURL, uu.FullURL.String())

	w.Header().Set("Content-Type", "text/html")
	if err := resTemplate.Execute(w, alert); err != nil {
		log.Fatal("template execution error:", err)
	}

}

//NormalizeURL takes inputURL from user and delete all tabs, spaces
func (uu *UserURL) NormalizeURL(inputURL string) error {

	uu.RawURL = strings.TrimSpace(inputURL)
	parse, err := url.Parse(uu.RawURL)
	if err != nil {
		log.Println("url parsing error")
		return nil
	}

	uu.FullURL = parse

	return nil
}

// CreateShortURL generate short URL
// Example: Kiu61yaLBDdL
func (uu *UserURL) CreateShortURL() error {

	var seederRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seederRand.Intn(len(charset))]
	}

	uu.ShortURL = string(b)
	return nil
}

func (uu *UserURL) SaveToCache() {

}
