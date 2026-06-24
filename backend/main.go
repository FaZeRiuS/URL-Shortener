package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	var err error

	database := initDB()
	defer database.Close()

	parsedTemplates, err := template.ParseFiles("templates/index.html")
	if err != nil {
		panic(err)
	}

	app := &App{
		db:   database,
		tmpl: parsedTemplates,
	}

	http.HandleFunc("/shorten", app.shortenerHandler)

	http.HandleFunc("/", app.redirectHandler)

	http.HandleFunc("/stats/", app.statsHandler)

	fmt.Println("Сервер запускається на порту :8080...")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
