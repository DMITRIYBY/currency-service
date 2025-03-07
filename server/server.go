package server

import (
	"encoding/json"
	"log"
	"net/http"

	"currency-service/db"
	"currency-service/models"
)

// Тут получаем все записи
func GetAllRates(w http.ResponseWriter, r *http.Request) {
	rows, err := db.GetDB().Query("SELECT cur_id, date, cur_abbreviation, cur_scale, cur_name, cur_official_rate FROM exchange_rates")
	if err != nil {
		http.Error(w, "Ошибка запроса", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Карта для группировки по дате
	ratesByDate := make(map[string][]models.ExchangeRate)

	// Чтение всех строк
	for rows.Next() {
		var rate models.ExchangeRate
		if err := rows.Scan(&rate.CurID, &rate.Date, &rate.CurAbbreviation, &rate.CurScale, &rate.CurName, &rate.CurOfficialRate); err != nil {
			http.Error(w, "Ошибка обработки данных", http.StatusInternalServerError)
			return
		}

		dateStr := rate.Date.Format("2006-01-02")

		// Добавляем курс в группу по дате
		ratesByDate[dateStr] = append(ratesByDate[dateStr], rate)
	}

	if len(ratesByDate) == 0 {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Записи не найдены",
		})
		return
	}

	// Возвращаем сгруппированные данные
	json.NewEncoder(w).Encode(ratesByDate)
}

// А тут запись за конкретную дату
func GetRatesByDate(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "Укажите дату в формате YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	rows, err := db.GetDB().Query("SELECT cur_id, date, cur_abbreviation, cur_scale, cur_name, cur_official_rate FROM exchange_rates WHERE date = ?", date)
	if err != nil {
		http.Error(w, "Ошибка запроса", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rates []models.ExchangeRate
	for rows.Next() {
		var rate models.ExchangeRate
		if err := rows.Scan(&rate.CurID, &rate.Date, &rate.CurAbbreviation, &rate.CurScale, &rate.CurName, &rate.CurOfficialRate); err != nil {
			http.Error(w, "Ошибка обработки данных", http.StatusInternalServerError)
			return
		}
		rates = append(rates, rate)
	}

	if len(rates) == 0 {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Записи за данную дату не найдены",
		})
		return
	}

	json.NewEncoder(w).Encode(rates)
}

func StartServer() {
	http.HandleFunc("/rates", GetAllRates)
	http.HandleFunc("/rate", GetRatesByDate)

	log.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
