package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"quotes/logger"
	"quotes/services"
	"quotes/storage"
)

func HandlerQuotesPost(s *storage.JSONStorage, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := services.Add(s, r, log); err != nil {
			log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Цитата успешно добавлена"))
	}
}

func HandlerQuotesGet(s *storage.JSONStorage, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var quotes []storage.Quote
		var err error

		quotes, err = services.GetQuotes(s, log, r)
		if err != nil {
			log.Error(fmt.Sprintf("Ошибка при получении цитат: %v", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(quotes); err != nil {
			log.Error(fmt.Sprintf("Ошибка при записи ответа: %v", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func HandlerQuotesRandomGet(s *storage.JSONStorage, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		quote, err := services.GetRandom(s, log)
		if err != nil {
			log.Error(fmt.Sprintf("Ошибка при получении рандомной цитаты: %v", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(quote); err != nil {
			log.Error(fmt.Sprintf("Ошибка при записи ответа: %v", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func HandlerQuotesDelete(s *storage.JSONStorage, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := services.Delete(s, log, r); err != nil {
			log.Error(err.Error())
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Цитата успешно удалена"))
	}
}
