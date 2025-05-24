package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"quotes/logger"
	"quotes/storage"
	"strconv"

	"github.com/gorilla/mux"
)

func Add(s *storage.JSONStorage, r *http.Request, log *logger.Logger) error {
	defer r.Body.Close()

	var quote storage.Quote
	err := json.NewDecoder(r.Body).Decode(&quote)
	if err != nil {
		return fmt.Errorf("Не удалось декодировать JSON из запроса: %w", err)
	}

	s.Add(quote)

	log.Info("Добавление новой цитаты прошло успешно (Author: " + quote.Author + "; Text: " + quote.Quote + ")")

	return nil
}

func GetQuotes(s *storage.JSONStorage, log *logger.Logger, r *http.Request) ([]storage.Quote, error) {
	quotes, err := s.GetQuotes()
	if err != nil {
		return quotes, err
	}

	log.Info("Получение всех цитат прошло успешно")

	params := r.URL.Query()
	author := params.Get("author")

	var response []storage.Quote

	if author == "" {
		return quotes, nil
	} else {
		for _, quote := range quotes {
			if author == quote.Author {
				response = append(response, quote)
			}
		}
		return response, nil
	}
}

func GetRandom(s *storage.JSONStorage, log *logger.Logger) (storage.Quote, error) {
	quotes, err := s.GetQuotes()
	if err != nil {
		return storage.Quote{}, err
	}

	randomQuote := quotes[rand.Intn(len(quotes))]

	log.Info("Получение случайной цитаты прошло успешно")

	return randomQuote, nil
}

func Delete(s *storage.JSONStorage, log *logger.Logger, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return fmt.Errorf("Неверный формат ID: %v", err)
	}

	if err = s.DeleteQuoteID(id); err != nil {
		return fmt.Errorf("Ошибка при удалении цитаты: %v", err)
	}

	return nil
}
