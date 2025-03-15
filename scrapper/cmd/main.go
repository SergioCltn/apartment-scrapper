package main

import (
	"math/rand"
	"time"

	"github.com/sergiocltn/apartment-scrapper/internal/config"
	"github.com/sergiocltn/apartment-scrapper/internal/provider"
	"github.com/sergiocltn/apartment-scrapper/internal/repository"
	"github.com/sergiocltn/apartment-scrapper/internal/service"
)

func main() {
	provider.InitLogger()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	if err := config.InitSQLite(); err != nil {
		provider.ErrorLogger.Fatalf("InitSqlite failed: %v", err)
	}
	db := config.DB

	apartmentRepo := *repository.NewApartmentRepository(db)
	if err := apartmentRepo.Initialize(); err != nil {
		provider.ErrorLogger.Fatalf("Failed to initialize apartment repository: %v", err)
	}

	svc := service.NewService(apartmentRepo)

	if err := svc.ScrapeFullList(); err != nil {
		provider.ErrorLogger.Fatalf("Scraping failed: %v", err)
	}

	provider.InfoLogger.Println("Scraper completed successfully")
}
