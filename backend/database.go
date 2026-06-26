package main

import (
	"database/sql"
	"fmt"
	"time"

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

func (a *App) flushClicksWorker() {
	buffer := make(map[string]int)
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {

		select {
		case code := <-a.clicksChan:
			buffer[code]++
		case <-ticker.C:
			for code, count := range buffer {
				_, err := a.db.Exec(`UPDATE urls SET clicks = clicks + ? WHERE code = ?`, count, code)
				if err != nil {
					fmt.Println("Помилка фонового оновлення кліків:", err)
				}
			}
			buffer = make(map[string]int)
		}
	}
}
