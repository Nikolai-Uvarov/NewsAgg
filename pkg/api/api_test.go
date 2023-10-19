package api

import (
	"NewsAgg/pkg/db/dbmock"
	"NewsAgg/pkg/db/obj"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_posts(t *testing.T) {
	// Создаём чистый объект API для теста.
	dbase := dbmock.New()

	dbase.SavePost(obj.Post{})

	api := New(dbase)
	// Создаём HTTP-запрос.
	req := httptest.NewRequest(http.MethodGet, "/news/1", nil)
	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	api.r.ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Читаем тело ответа.
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Раскодируем JSON в массив новостей.
	var data []obj.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Проверяем, что в массиве ровно один элемент.
	const wantLen = 1
	if len(data) != wantLen {
		t.Fatalf("получено %d записей, ожидалось %d", len(data), wantLen)
	}
}

func TestAPI_postByID(t *testing.T) {
	// Создаём чистый объект API для теста.
	dbase := dbmock.New()

	dbase.SavePost(obj.Post{})

	api := New(dbase)
	// Создаём HTTP-запрос.
	req := httptest.NewRequest(http.MethodGet, "/news?postID=1", nil)
	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	api.r.ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Читаем тело ответа.
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Раскодируем JSON в объект.
	var data obj.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Проверяем, что получили нужный объект.
	const wantID = 1
	if data.ID != wantID {
		t.Fatalf("получена новость с  %d , ожидалась %d", data.ID, wantID)
	}
}
