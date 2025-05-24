package storage_test

import (
	"encoding/json"
	"os"
	"quotes/logger"
	"quotes/storage"
	"testing"
)

func TestCreateJSONStorage(t *testing.T) {
	log, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Не удалось создать логгер: %v", err)
	}
	defer os.Remove("log.log")

	// Тест 1: Файл не существует
	tempFile, err := os.CreateTemp("", "test_JSON.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	os.Remove(tempFile.Name())

	s, err := storage.CreateJSONStorage(tempFile.Name(), log)
	if err != nil {
		t.Fatalf("CreateJSONStorage вернула ошибку: %v", err)
	}
	if len(s.Quotes) != 0 || s.IdCounter != 1 {
		t.Errorf("Ожидалось пустое хранилище с IdCounter=1, получено: %+v", s)
	}

	// Тест 2: Пустой файл
	tempFile, err = os.CreateTemp("", "test_JSON.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tempFile.Name())

	s, err = storage.CreateJSONStorage(tempFile.Name(), log)
	if err != nil {
		t.Fatalf("CreateJSONStorage вернула ошибку: %v", err)
	}
	if len(s.Quotes) != 0 || s.IdCounter != 1 {
		t.Errorf("Ожидалось пустое хранилище с IdCounter=1, получено: %+v", s)
	}

	// Тест 3: Файл с данными
	tempFile, err = os.CreateTemp("", "test_JSON.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tempFile.Name())

	testData := []storage.QuoteStore{
		{ID: 1, Quote: "Quote 1", Author: "Author 1"},
		{ID: 2, Quote: "Quote 2", Author: "Author 2"},
	}
	fileData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Не удалось сериализовать тестовые данные: %v", err)
	}
	err = os.WriteFile(tempFile.Name(), fileData, 0644)
	if err != nil {
		t.Fatalf("Не удалось записать тестовые данные в файл: %v", err)
	}

	s, err = storage.CreateJSONStorage(tempFile.Name(), log)
	if err != nil {
		t.Fatalf("CreateJSONStorage вернула ошибку: %v", err)
	}
	if len(s.Quotes) != 2 || s.IdCounter != 3 {
		t.Errorf("Ожидалось хранилище с 2 цитатами и IdCounter=3, получено: %+v", s)
	}

}

func TestSave(t *testing.T) {
	log, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Не удалось создать логгер: %v", err)
	}
	defer os.Remove("log.log")

	tempFile, err := os.CreateTemp("", "test_JSON.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tempFile.Name())

	s := &storage.JSONStorage{
		Quotes: []storage.QuoteStore{
			{ID: 1, Quote: "Quote 1", Author: "Author 1"},
			{ID: 2, Quote: "Quote 2", Author: "Author 2"},
		},
		IdCounter: 3,
	}

	err = s.Save(tempFile.Name(), log)
	if err != nil {
		t.Fatalf("Save вернула ошибку: %v", err)
	}

	fileData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Не удалось прочитать файл: %v", err)
	}

	var savedQuotes []storage.QuoteStore
	err = json.Unmarshal(fileData, &savedQuotes)
	if err != nil {
		t.Fatalf("Не удалось десериализовать данные из файла: %v", err)
	}

	if len(savedQuotes) != len(s.Quotes) {
		t.Errorf("Ожидалось %d цитат, получено: %d", len(s.Quotes), len(savedQuotes))
	}

	for i, quote := range savedQuotes {
		if quote.ID != s.Quotes[i].ID || quote.Quote != s.Quotes[i].Quote || quote.Author != s.Quotes[i].Author {
			t.Errorf("Цитата %d не совпадает. Ожидалось %+v, получено %+v", i, s.Quotes[i], quote)
		}
	}

	// Тест 2: Ошибка записи в файл
	invalidPath := "/invalid/path/test.json"
	err = s.Save(invalidPath, log)
	if err == nil {
		t.Error("Ожидалась ошибка при записи в недоступный путь")
	}
}
