package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func initDB() *sql.DB {
	var err error
	database, err := sql.Open("sqlite", "file:shortener.db?_journal_mode=WAL")
	if err != nil {
		panic(err)
	}
	err = database.Ping()
	if err != nil {
		panic(err)
	}

	_, err = database.Exec(`CREATE TABLE IF NOT EXISTS urls(
		code TEXT PRIMARY KEY,
		original_url TEXT NOT NULL,
		clicks INTEGER DEFAULT 0
		)`)
	if err != nil {
		panic(err)
	}
	return database
}
