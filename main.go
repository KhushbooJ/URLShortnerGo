package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}

var urlDB = make(map[string]URL)

func generateShortURL(originalUrl string) string {
	hash := md5.New()
	hash.Write(([]byte(originalUrl)))
	hashValue := fmt.Sprintf("%x", hash.Sum(nil))[:6] // Take the first 6 characters of the hash
	return hashValue
}

func saveURL(originalUrl string) string {
	shortUrl := generateShortURL(originalUrl)
	url := URL{
		ID:          shortUrl,
		OriginalURL: originalUrl,
		ShortURL:    shortUrl,
		CreatedAt:   time.Now(),
	}
	urlDB[shortUrl] = url
	return url.ShortURL
}

func getOriginalURL(id string) (URL, error) {
	url, isAvail := urlDB[id]
	if !isAvail {
		return URL{}, errors.New("URL not found")
	}
	return url, nil

}

func shortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	shortURL := saveURL(data.URL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{
		ShortURL: shortURL,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getOriginalURL(id)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	fmt.Println("Welcome to the URL Shortener Service!")

	http.HandleFunc("/shorten", shortUrlHandler)
	http.HandleFunc("/redirect/", redirectHandler)

	//start server on 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
}
