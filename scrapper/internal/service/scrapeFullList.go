package service

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"

	"github.com/sergiocltn/apartment-scrapper/internal/model"
	"github.com/sergiocltn/apartment-scrapper/internal/provider"
	"github.com/sergiocltn/apartment-scrapper/internal/provider/scrapper"
	"github.com/sergiocltn/apartment-scrapper/internal/repository"
)

type Service struct {
	ApartmentRepo repository.ApartmentRepository
}

func NewService(apartmentRepo repository.ApartmentRepository) *Service {
	return &Service{
		ApartmentRepo: apartmentRepo,
	}
}

func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (s *Service) ScrapeFullList() error {
	var firstID string
	page := 1

	for {
		scrapedList, err := scrapper.ScrapeList(page)
		if err != nil {
			return fmt.Errorf("failed to scrape list on page %d: %v", page, err)
		}

		if firstID != "" && slices.Contains(scrapedList.PropertyIDs, firstID) {
			provider.InfoLogger.Println("Found duplicate first ID, stopping scrape")
			break
		}

		if firstID == "" && len(scrapedList.PropertyIDs) > 0 {
			firstID = scrapedList.PropertyIDs[0]
			provider.InfoLogger.Printf("First ID: %s", firstID)
		}

		provider.InfoLogger.Printf("Analyzing page: %d", page)

		for _, apartmentID := range scrapedList.PropertyIDs {
			exists, err := s.ApartmentRepo.Exists(apartmentID)
			if err != nil {
				provider.ErrorLogger.Printf("Error checking existence of %s: %v", apartmentID, err)
				continue
			}
			if exists {
				provider.InfoLogger.Printf("Already saved: %s", apartmentID)
				continue
			}

			sleep(rand.Intn(5000) + 5000)

			scrapedData, err := scrapper.ScrapeApartment(apartmentID)
			if err != nil {
				provider.ErrorLogger.Printf("Failed to scrape apartment %s: %v", apartmentID, err)
				continue
			}

			apartment := model.Apartment{
				ID:                apartmentID,
				Title:             scrapedData.Title,
				PricePerSqm:       scrapedData.PriceFeatures.PricePerSqm,
				PropertyPrice:     scrapedData.PriceFeatures.PropertyPrice,
				CommunityFees:     scrapedData.PriceFeatures.CommunityFees,
				ApartmentStatus:   scrapedData.Details.ApartmentStatus,
				Building:          scrapedData.Details.Building,
				BasicFeatures:     scrapedData.Details.BasicFeatures,
				EnergyCertificate: scrapedData.Details.EnergyCertificate,
				Location:          strings.Join(scrapedData.Location, ", "),
				Description:       scrapedData.Description,
				CreatedAt:         time.Now(),
			}

			if err := s.ApartmentRepo.Save(apartment); err != nil {
				provider.ErrorLogger.Printf("Failed to save apartment %s: %v", apartmentID, err)
				continue
			}
		}

		provider.InfoLogger.Printf("Page %d scrapping finished", page)
		page++
		sleep(60000)
	}

	return nil
}
