package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var homeTemplate *template.Template

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Forbidden HTTP Method", http.StatusForbidden)
	}

	homeTemplate = template.Must(template.ParseFiles("template/index.gohtml"))

	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(w, nil); err != nil {
		log.Fatal("template execution error:", err)
	}
}

func shortURL(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Forbidden HTTP Method", http.StatusForbidden)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Fatal("error parsing form", err)
	}
	_, _ = fmt.Fprintln(w, r.PostForm["full_url"])
}

func main() {

	// Static server
	staticServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticServer))

	// GET method
	http.HandleFunc("/", home)
	// Post method
	http.HandleFunc("/shortURL", shortURL)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
