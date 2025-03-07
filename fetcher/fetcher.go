package fetcher

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"currency-service/models"
)

const NBRB_API_URL = "https://api.nbrb.by/exrates/rates?periodicity=0"

func FetchRates() []models.ExchangeRate {
	resp, err := http.Get(NBRB_API_URL)
	if err != nil {
		log.Println("Ошибка при запросе к API НБРБ:", err)
		return nil
	}
	defer resp.Body.Close()

	var apiRates []struct {
		CurID           int     `json:"Cur_ID"`
		Date            string  `json:"Date"`
		CurAbbreviation string  `json:"Cur_Abbreviation"`
		CurScale        int     `json:"Cur_Scale"`
		CurName         string  `json:"Cur_Name"`
		CurOfficialRate float64 `json:"Cur_OfficialRate"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiRates); err != nil {
		log.Println("Ошибка парсинга JSON:", err)
		return nil
	}

	var rates []models.ExchangeRate
	for _, r := range apiRates {
		parsedDate, _ := time.Parse("2006-01-02T15:04:05", r.Date)
		rates = append(rates, models.ExchangeRate{
			CurID:           r.CurID,
			Date:            parsedDate,
			CurAbbreviation: r.CurAbbreviation,
			CurScale:        r.CurScale,
			CurName:         r.CurName,
			CurOfficialRate: r.CurOfficialRate,
		})
	}

	return rates
}
