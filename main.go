package main

import (
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type PageData struct {
	Name         string
	ShortenedURL string
	Error        string
}

var urlMap = make(map[string]string)
var baseURL = "http://localhost:8080/"

func generateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			originalURL := r.FormValue("url")
			if originalURL == "" {
				tmpl.ExecuteTemplate(w, "index.html", PageData{
					Name:  "URL Shortener",
					Error: "URL cannot be empty!",
				})
				return
			}
			shortURL := generateShortURL()
			urlMap[shortURL] = originalURL
			tmpl.ExecuteTemplate(w, "index.html", PageData{
				Name:         "URL Shortener",
				ShortenedURL: baseURL + shortURL,
			})
		} else {
			tmpl.ExecuteTemplate(w, "index.html", PageData{
				Name: "URL Shortener",
			})
		}
	})

	router.HandleFunc("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		shortURL := strings.TrimPrefix(r.URL.Path, "/")
		originalURL, exists := urlMap[shortURL]
		if !exists {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, originalURL, http.StatusFound)
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Starting website at localhost:8080")

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("An error occurred:", err)
	}
}
