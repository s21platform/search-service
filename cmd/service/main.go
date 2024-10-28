package main

import "github.com/s21platform/search-service/internal/config"

func main() {
	// чтение конфига
	cfg := config.MustLoad()
	_ = cfg
}
