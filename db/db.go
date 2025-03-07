package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"currency-service/models"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	dbTemp, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	if err := dbTemp.Ping(); err != nil {
		log.Fatal("Ошибка проверки соединения с БД:", err)
	}

	db = dbTemp

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS exchange_rates (
			id INT AUTO_INCREMENT PRIMARY KEY,
			cur_id INT,
			date DATETIME,
			cur_abbreviation VARCHAR(10),
			cur_scale INT,
			cur_name VARCHAR(255),
			cur_official_rate DECIMAL(10,4),
			UNIQUE KEY uniq_cur_date (cur_id, date)
		);
	`)
	if err != nil {
		log.Fatal("Ошибка создания таблицы exchange_rates:", err)
	}
}

func GetDB() *sql.DB {
	return db
}

func SaveRates(rates []models.ExchangeRate) {
	for _, rate := range rates {
		_, err := db.Exec(`
            INSERT IGNORE INTO exchange_rates (cur_id, date, cur_abbreviation, cur_scale, cur_name, cur_official_rate)
            VALUES (?, ?, ?, ?, ?, ?)`,
			rate.CurID, rate.Date, rate.CurAbbreviation, rate.CurScale, rate.CurName, rate.CurOfficialRate,
		)

		if err != nil {
			log.Println("Ошибка сохранения курса:", err)
		}
	}
}
