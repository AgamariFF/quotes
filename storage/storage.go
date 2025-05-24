package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"quotes/logger"
	"sync"
)

type JSONStorage struct {
	Quotes    []QuoteStore
	IdCounter int
	mute      sync.Mutex
}

func CreateJSONStorage(filename string, log *logger.Logger) (*JSONStorage, error) {
	var storage JSONStorage

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(filename)
			if err != nil {
				return &storage, fmt.Errorf("Не удалось создать файл: %w", err)
			}
			log.Info("Файл хранилища отсутствовал и был успешно создан")
			storage.Quotes = []QuoteStore{}
			storage.IdCounter = 1
			return &storage, nil
		}
	} else {
		log.Info("Файл хранилища успшено открыт")
	}

	if len(data) == 0 {
		log.Info("Файл пустой, инициализация пустого хранилища")
		storage.Quotes = []QuoteStore{}
		storage.IdCounter = 1
		return &storage, nil
	}

	if err = json.Unmarshal(data, &storage.Quotes); err != nil {
		return &storage, fmt.Errorf("Не удалось десериализовать данные: %w", err)
	}

	storage.IdCounter = len(storage.Quotes) + 1

	log.Info("Инициализация хранилища прошла успешно")
	return &storage, nil
}

func (storage *JSONStorage) Save(filename string, log *logger.Logger) error {
	storage.mute.Lock()
	defer storage.mute.Unlock()
	data, err := json.Marshal(storage.Quotes)
	if err != nil {
		return err
	}
	log.Info("Сериализация данных прошла успешно")

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}
	log.Info("Сохранение данных прошло успешно")

	return nil
}

func (storage *JSONStorage) Add(quote Quote) {
	storage.mute.Lock()
	defer storage.mute.Unlock()

	quoteStore := QuoteStore{
		Quote:  quote.Quote,
		Author: quote.Author,
		ID:     storage.IdCounter,
	}

	storage.Quotes = append(storage.Quotes, quoteStore)

	storage.IdCounter++
}

func (storage *JSONStorage) GetQuotes() ([]Quote, error) {
	storage.mute.Lock()
	defer storage.mute.Unlock()

	Quotes := make([]Quote, len(storage.Quotes))

	for i, quoteStore := range storage.Quotes {
		Quotes[i] = Quote{
			Quote:  quoteStore.Quote,
			Author: quoteStore.Author,
		}
	}

	if len(Quotes) == 0 {
		return Quotes, fmt.Errorf("Отсутствуют цитаты")
	}

	return Quotes, nil
}

func (storage *JSONStorage) DeleteQuoteID(id int) error {
	storage.mute.Lock()
	defer storage.mute.Unlock()

	for i, quote := range storage.Quotes {
		if quote.ID == id {
			storage.Quotes = append(storage.Quotes[:i], storage.Quotes[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Цитата с указанным ID %d не найдена", id)
}
