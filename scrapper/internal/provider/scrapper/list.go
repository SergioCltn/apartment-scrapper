package scrapper

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sergiocltn/apartment-scrapper/internal/provider"
)

type ScrapedListData struct {
	PropertyIDs []string `json:"propertyIds"`
}

func ScrapeList(page int) (ScrapedListData, error) {
	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:135.0) Gecko/20100101 Firefox/135.0",
		"Accept-Language": "en-US,en;q=0.5",
	}

	url := "https://www.idealista.com/areas/venta-viviendas/con-precio-hasta_190000,precio-desde_80000/pagina-" + strconv.Itoa(page) + "?shape=%28%28%28%7Bj%7C%7EE%7EhfZ%7CH_%40%7EEaAdEeDtL%7BLjDwElFyIxDyHlAwFmAgHgEkJqCuEdCnO%7BFsCwNwFcVwGuLuDyFsBiGsCqJgGiN%7DKcHwGwIiJoHiGsEsCoCqAoYqPwPkKaFuD%7DHgHoOoOwBuCoH%7DMoAcDoH_OgGaQoAuEmFaPYaBcCgGqCuDqCsC_FuCyDaBeEq%40aF_%40cC%3FeEn%40sEbBmFrCiBpA%7DHhHmA%60BcH%7CLiBvEmAtE_DhIoAfGaAhIe%40%7ENXfHr%40xHpCnNxD%7ENzHdUvBrDxDzIlAtDdEhJvBvFzAbCxDhJxK%7CLlM%60Q%60FhHjDbDbCtCfGtFfEtCzAn%40lRdEdErApCN%7EJrBjK%3FrEN%60McCpCqAlFcBpC_%40dE%3FzFpAzTzJvIbDnMxGzMdFxKbChGp%40%60H%5EpH%3F%29%29%29&ordenado-por=fecha-publicacion-desc"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ScrapedListData{}, fmt.Errorf("failed to create request: %v", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return ScrapedListData{}, fmt.Errorf("failed to scrape Idealista: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ScrapedListData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ScrapedListData{}, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var propertyIDs []string
	doc.Find(".items-container.items-list article.item").Each(func(_ int, s *goquery.Selection) {
		id, exists := s.Attr("data-element-id")
		if exists {
			propertyIDs = append(propertyIDs, id)
		}
	})

	data := ScrapedListData{PropertyIDs: propertyIDs}
	provider.InfoLogger.Printf("Scraped %d property IDs from page %d", len(propertyIDs), page)
	return data, nil
}
