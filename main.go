package main

import (
	"currency-service/db"
	"currency-service/fetcher"
	"currency-service/server"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
)

func connectToDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	var dbConn *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		dbConn, err = sql.Open("mysql", dsn)
		if err == nil {
			err = dbConn.Ping()
			if err == nil {
				fmt.Println("Успешное подключение к БД")
				return dbConn, nil
			}
		}
		log.Println("Ошибка подключения к БД, повтор через 5 секунд...")
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
}

func collectAndSaveRates() {
	fmt.Println("Сбор данных о курсах валют...")
	rates := fetcher.FetchRates()
	if rates != nil {
		db.SaveRates(rates)
		fmt.Println("Курсы валют успешно обновлены в БД.")
	} else {
		fmt.Println("Ошибка: Не удалось получить курсы валют.")
	}
}

func main() {
	db.InitDB()

	dbConn, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Фетчинг запускается в 01:00 ночи каждые сутки
	c := cron.New()
	_, err = c.AddFunc("0 1 * * *", collectAndSaveRates)
	if err != nil {
		log.Fatal("Ошибка добавления cron-задачи:", err)
	}
	c.Start()

	server.StartServer()
	fmt.Println("Сервер запущен на порту 8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Сервис работает!")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
