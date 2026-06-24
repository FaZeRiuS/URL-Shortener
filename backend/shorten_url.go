package main

import (
	"crypto/rand"
	"math/big"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func shortenURL() string {
	b := make([]byte, 6) // Створюємо заготовку (масив) відразу на 6 байт
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		b[i] = letters[n.Int64()] // Заповнюємо кожну комірку прямо в пам'яті
	}
	return string(b) // Конвертуємо весь масив у рядок один єдиний раз
}
