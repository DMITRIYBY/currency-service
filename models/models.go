package models

import "time"

type ExchangeRate struct {
	ID              int       `json:"id"`
	CurID           int       `json:"cur_id"`
	Date            time.Time `json:"date"`
	CurAbbreviation string    `json:"cur_abbreviation"`
	CurScale        int       `json:"cur_scale"`
	CurName         string    `json:"cur_name"`
	CurOfficialRate float64   `json:"cur_official_rate"`
}
