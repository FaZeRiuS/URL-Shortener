package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	_ "modernc.org/sqlite"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (a *App) shortenerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	var req ShortenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to read JSON body", http.StatusInternalServerError)
		return
	}

	body := req.URL

	parsedURL, err := url.ParseRequestURI(body)
	if err != nil || parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid URL Format, must be a valid domain name (http or https)"})
		return
	}

	var code string
	for {
		code = shortenURL()

		_, dbErr := a.db.Exec(`INSERT INTO urls(code, original_url) VALUES(?, ?)`, code, body)
		if dbErr == nil {
			break
		}

		if strings.Contains(dbErr.Error(), "UNIQUE") {
			continue
		}

		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to save to database"})
		return
	}

	a.cache.Store(code, body)

	w.Header().Set("Location", "http://localhost:8080/"+code)
	resp := ShortenResponse{ShortenedURL: "http://localhost:8080/" + code}
	writeJSON(w, http.StatusCreated, resp)

}

func (a *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]

	if r.URL.Path == "/" {
		err := a.tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	} else if len(code) != 6 {
		http.Error(w, "Invalid Code Format", http.StatusNotFound)
		return
	}

	var originalUrl string

	if value, ok := a.cache.Load(code); ok {
		originalUrl = value.(string)
	} else {
		dbErr := a.db.QueryRow(`SELECT original_url FROM urls WHERE code = ?`, code).Scan(&originalUrl)
		if dbErr == sql.ErrNoRows {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		} else if dbErr != nil {
			fmt.Println(dbErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		a.cache.Store(code, originalUrl)
	}

	a.clicksChan <- code

	http.Redirect(w, r, originalUrl, http.StatusSeeOther)
}

func (a *App) statsHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[7:]

	if r.Method != "GET" {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "Method Not Allowed"})
		return
	}

	if len(code) != 6 {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "URL not found"})
		return
	}

	var clicksCount int

	var originalURL string

	err := a.db.QueryRow(`SELECT clicks, original_url FROM urls WHERE code = ?`, code).Scan(&clicksCount, &originalURL)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "URL not found"})
		return
	} else if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Unexpected error"})
		return
	}

	writeJSON(w, http.StatusOK, StatsResponse{Count: clicksCount, OriginalURL: originalURL})
}
