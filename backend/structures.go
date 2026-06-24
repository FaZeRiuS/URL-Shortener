package main

import (
	"database/sql"
	"html/template"
)

type App struct {
	db   *sql.DB
	tmpl *template.Template
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortenedURL string `json:"short_url"`
}

type StatsResponse struct {
	Count       int    `json:"count"`
	OriginalURL string `json:"original_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
