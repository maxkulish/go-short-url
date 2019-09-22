package main

import (
	"go-short-url/cache"
	"go-short-url/config"
	"go-short-url/controller"
	"html/template"
	"log"
	"net/http"
)

var homeTemplate *template.Template

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Forbidden HTTP Method", http.StatusMethodNotAllowed)
	}

	homeTemplate = template.Must(template.ParseFiles("template/index.gohtml"))

	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(w, nil); err != nil {
		log.Fatal("template execution error:", err)
	}
}

func main() {

	newCache := cache.NewCache()
	userURL := controller.NewUserURL(newCache)
	shortURL := controller.NewShortURL(newCache)

	// Static server
	staticServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticServer))

	// GET method
	http.HandleFunc("/", homeHandler)
	// Post method
	http.HandleFunc("/createURL", userURL.CreateURL)

	// GET /short/:shortURL
	http.HandleFunc("/s/", shortURL.HandleShortURL)

	log.Fatal(http.ListenAndServe(config.HostURL, nil))
}
