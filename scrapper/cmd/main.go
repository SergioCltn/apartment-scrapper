package main

import (
	"math/rand"
	"time"

	"github.com/sergiocltn/apartment-scrapper/internal/config"
	"github.com/sergiocltn/apartment-scrapper/internal/provider"
	"github.com/sergiocltn/apartment-scrapper/internal/service"
)

func runScraper() error {
	return service.ScrapeFullList()
}

func main() {
	provider.InitLogger()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	if err := config.InitSQLite(); err != nil {
		provider.ErrorLogger.Fatalf("InitSqlite failed: %v", err)
	}

	defer config.CloseDB()
	if err := runScraper(); err != nil {
	}
	provider.InfoLogger.Println("Scraper completed successfully")
}
