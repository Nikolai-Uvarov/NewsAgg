package main

import (
	"NewsAgg/pkg/api"
	"NewsAgg/pkg/db/dbmock"
	"NewsAgg/pkg/rss"
	"log"
	"net/http"
)

func main() {

	// Инициализация БД в памяти.
	db := dbmock.New()

	//Запуск воркера, прослушивающего и сохраняющего в БД RSS-ленты
	rss.Collect(db)

	// Создание объекта API, использующего БД в памяти.
	api := api.New(db)

	// Запуск сетевой службы и HTTP-сервера
	// на всех локальных IP-адресах на порту 8080.
	err := http.ListenAndServe(":8080", api.Router())
	if err != nil {
		log.Fatal(err)
	}

}
