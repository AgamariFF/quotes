package services_test

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"quotes/logger"
	"quotes/services"
	"quotes/storage"
	"testing"

	"github.com/gorilla/mux"
)

func TestAdd(t *testing.T) {
	log, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Не удалось создать логгер: %v", err)
	}
	defer os.Remove("log.log")

	s, err := storage.CreateJSONStorage("temp_JSON.json", log)
	if err != nil {
		t.Fatalf("Не удалось инициализировать хранилище: %v", err)
	}
	defer os.Remove("temp_JSON.json")

	inputQuote := storage.Quote{
		Quote:  "Simple quote",
		Author: "Simple author",
	}

	jsonData, err := json.Marshal(inputQuote)
	if err != nil {
		t.Fatalf("Не удалось создать JSON: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	err = services.Add(s, req, log)
	if err != nil {
		t.Fatalf("Add вернула ошибку: %v", err)
	}

	if len(s.Quotes) != 1 {
		t.Fatalf("Ожидалась одна цитата, получено: %d", len(s.Quotes))
	}
	addedQuote := s.Quotes[0]
	if addedQuote.Quote != inputQuote.Quote || addedQuote.Author != inputQuote.Author {
		t.Errorf("Некорректная цитата в хранилище. Ожидалось: %+v, получено: %+v", inputQuote, addedQuote)
	}
}

func TestGetQuotes(t *testing.T) {
	log, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Не удалось создать логгер: %v", err)
	}
	defer os.Remove("log.log")

	s, err := storage.CreateJSONStorage("temp_JSON.json", log)
	if err != nil {
		t.Fatalf("Не удалось инициализировать хранилище: %v", err)
	}
	defer os.Remove("temp_JSON.json")

	quotes := []storage.Quote{
		{Quote: "Quote 1", Author: "Author 1"},
		{Quote: "Quote 2", Author: "Author 2"},
		{Quote: "Quote 3", Author: "Author 1"},
	}
	for _, quote := range quotes {
		s.Add(quote)
	}

	// Тест 1: Получение всех цитат
	req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
	allQuotes, err := services.GetQuotes(s, log, req)
	if err != nil {
		t.Fatalf("GetQuotes вернула ошибку: %v", err)
	}
	if len(allQuotes) != len(quotes) {
		t.Errorf("Ожидалось %d цитат, получено: %d", len(quotes), len(allQuotes))
	}

	// Тест 2: Фильтрация цитат по автору
	req = httptest.NewRequest(http.MethodGet, "/quotes?author=Author+1", nil)
	filteredQuotes, err := services.GetQuotes(s, log, req)
	if err != nil {
		t.Fatalf("GetQuotes вернула ошибку: %v", err)
	}
	if len(filteredQuotes) != 2 {
		t.Errorf("Ожидалось 2 цитаты от 'Author 1', получено: %d", len(filteredQuotes))
	}
	for _, quote := range filteredQuotes {
		if quote.Author != "Author 1" {
			t.Errorf("Получена цитата с неправильным автором: %+v", quote)
		}
	}

	// Тест 3: Фильтрация цитат по несуществующему автору
	req = httptest.NewRequest(http.MethodGet, "/quotes?author=A", nil)
	filteredQuotes, err = services.GetQuotes(s, log, req)
	if err != nil {
		t.Fatalf("GetQuotes вернула ошибку: %v", err)
	}
	if len(filteredQuotes) != 0 {
		t.Errorf("Ожидалось 0 цитат, получено: %d", len(filteredQuotes))
	}
}

func TestGetRandom(t *testing.T) {
	// Создаем логгер
	log, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Не удалось создать логгер: %v", err)
	}
	defer os.Remove("log.log")

	s, err := storage.CreateJSONStorage("temp_JSON.json", log)
	if err != nil {
		t.Fatalf("Не удалось инициализировать хранилище: %v", err)
	}
	defer os.Remove("temp_JSON.json")

	quotes := []storage.Quote{
		{Quote: "Quote 1", Author: "Author 1"},
		{Quote: "Quote 2", Author: "Author 2"},
		{Quote: "Quote 3", Author: "Author 1"},
	}
	for _, quote := range quotes {
		s.Add(quote)
	}

	rand.Seed(3)

	// Тест 1: Получение случайной цитаты
	randomQuote, err := services.GetRandom(s, log)
	if err != nil {
		t.Fatalf("GetRandom вернула ошибку: %v", err)
	}

	found := false
	for _, quote := range quotes {
		if quote.Quote == randomQuote.Quote && quote.Author == randomQuote.Author {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Получена цитата, которая отсутствует в хранилище: %+v", randomQuote)
	}

	// Тест 2: Пустое хранилище
	emptyStorage, err := storage.CreateJSONStorage("empty_JSON.json", log)
	if err != nil {
		t.Fatalf("Не удалось инициализировать пустое хранилище: %v", err)
	}
	defer os.Remove("empty_JSON.json")

	_, err = services.GetRandom(emptyStorage, log)
	if err == nil {
		t.Error("Ожидалась ошибка при получении случайной цитаты из пустого хранилища")
	}
}

func TestDelete(t *testing.T) {
	log, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Не удалось создать логгер: %v", err)
	}
	defer os.Remove("log.log")

	s, err := storage.CreateJSONStorage("temp_JSON.json", log)
	if err != nil {
		t.Fatalf("Не удалось инициализировать хранилище: %v", err)
	}
	defer os.Remove("temp_JSON.json")

	quotes := []storage.Quote{
		{Quote: "Quote 1", Author: "Author 1"},
		{Quote: "Quote 2", Author: "Author 2"},
		{Quote: "Quote 3", Author: "Author 1"},
	}
	for _, quote := range quotes {
		s.Add(quote)
	}

	// Тест 1: Удаление существующей цитаты
	req := httptest.NewRequest(http.MethodDelete, "/quotes/1", nil)
	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	err = services.Delete(s, log, req)
	if err != nil {
		t.Fatalf("Delete вернула ошибку: %v", err)
	}

	if len(s.Quotes) != len(quotes)-1 {
		t.Errorf("Ожидалось %d цитат, получено: %d", len(quotes)-1, len(s.Quotes))
	}

	// Тест 2: Удаление несуществующей цитаты
	req = httptest.NewRequest(http.MethodDelete, "/quotes/999", nil)
	vars = map[string]string{"id": "999"}
	req = mux.SetURLVars(req, vars)

	err = services.Delete(s, log, req)
	if err == nil {
		t.Error("Ожидалась ошибка при удалении несуществующей цитаты")
	}

	// Тест 3: Некорректный формат ID
	req = httptest.NewRequest(http.MethodDelete, "/quotes/invalid", nil)
	vars = map[string]string{"id": "invalid"}
	req = mux.SetURLVars(req, vars)

	err = services.Delete(s, log, req)
	if err == nil {
		t.Error("Ожидалась ошибка при некорректном формате ID")
	}
}
