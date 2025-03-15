package scrapper

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sergiocltn/apartment-scrapper/internal/config"
	"github.com/sergiocltn/apartment-scrapper/internal/provider"
)

type PriceFeatures struct {
	PropertyPrice string `json:"propertyPrice,omitempty"`
	PricePerSqm   string `json:"pricePerSqm,omitempty"`
	CommunityFees string `json:"communityFees,omitempty"`
}

type Details struct {
	ApartmentStatus   string `json:"apartmentStatus,omitempty"`
	BasicFeatures     string `json:"basicFeatures,omitempty"`
	Building          string `json:"building,omitempty"`
	EnergyCertificate string `json:"energyCertificate,omitempty"`
}

type ScrapedData struct {
	Title         string        `json:"title,omitempty"`
	Description   string        `json:"description,omitempty"`
	Details       Details       `json:"details,omitempty"`
	DetailInfoTag string        `json:"detailInfoTag,omitempty"`
	TitleMinor    string        `json:"titleMinor,omitempty"`
	InfoFeatures  []string      `json:"infoFeatures,omitempty"`
	PriceFeatures PriceFeatures `json:"priceFeatures"`
	Location      []string      `json:"location,omitempty"`
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0",
}

func getRandomUserAgent() string {
	rand.New(rand.NewSource(time.Now().Unix()))
	return userAgents[rand.Intn(len(userAgents))]
}

func ScrapeApartment(apartmentId string) (ScrapedData, error) {
	cfg := config.GetConfig()
	url := fmt.Sprintf("https://www.idealista.com/inmueble/%s", apartmentId)

	client := &http.Client{
		Timeout: cfg.ScraperTimeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ScrapedData{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return ScrapedData{}, fmt.Errorf("failed to scrape Idealista: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ScrapedData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ScrapedData{}, fmt.Errorf("failed to parse HTML: %v", err)
	}

	data := ScrapedData{}

	description := strings.TrimSpace(doc.Find(".comment .adCommentsLanguage p").First().Text())
	data.Description = description

	title := strings.TrimSpace(doc.Find(".main-info__title-main").First().Text())
	data.Title = title

	titleMinor := strings.TrimSpace(doc.Find(".main-info__title-minor").First().Text())
	data.TitleMinor = titleMinor

	detailInfoTag := strings.TrimSpace(doc.Find(".detail-info-tags .tag").First().Text())
	data.DetailInfoTag = detailInfoTag

	doc.Find(".info-features span").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			data.InfoFeatures = append(data.InfoFeatures, text)
		}
	})

	doc.Find("#details .details-property-h2").Each(func(_ int, s *goquery.Selection) {
		sectionTitle := strings.TrimSpace(s.Text())
		if sectionTitle == "" {
			return
		}

		value := ""
		s.Next().Filter(".details-property_features").Find("ul li").Each(func(i int, li *goquery.Selection) {
			text := strings.TrimSpace(li.Text())
			if text != "" {
				if i == 0 {
					value = text
				} else {
					value += "; " + text
				}
			}
		})

		switch sectionTitle {
		case "Situación de la vivienda":
			data.Details.ApartmentStatus = value
		case "Características básicas":
			data.Details.BasicFeatures = value
		case "Edificio":
			data.Details.Building = value
		case "Certificado energético":
			data.Details.EnergyCertificate = value
		}
	})

	doc.Find("#headerMap ul li.header-map-list").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			data.Location = append(data.Location, text)
		}
	})

	data.PriceFeatures = PriceFeatures{}
	doc.Find(".price-features__container p").Each(func(_ int, s *goquery.Selection) {
		spans := s.Find(".flex-feature-details")
		switch spans.Length() {
		case 2:
			label := strings.TrimSpace(spans.Eq(0).Text())
			value := strings.TrimSpace(spans.Eq(1).Text())
			if label == "" || value == "" {
				return
			}
			switch {
			case strings.Contains(label, "Precio del inmueble"):
				data.PriceFeatures.PropertyPrice = value
			case strings.Contains(label, "Precio por m²"):
				data.PriceFeatures.PricePerSqm = value
			}
		case 1:
			text := strings.TrimSpace(spans.Eq(0).Text())
			if strings.HasPrefix(text, "Gastos de comunidad") {
				communityFees := strings.TrimSpace(strings.Replace(text, "Gastos de comunidad", "", 1))
				data.PriceFeatures.CommunityFees = communityFees
			}
		}
	})

	provider.InfoLogger.Printf("Scraped apartment with ID: %s", apartmentId)
	return data, nil
}
