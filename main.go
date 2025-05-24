package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"quotes/handlers"
	"quotes/logger"
	"quotes/storage"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func WaitClose(log *logger.Logger) chan struct{} {
	stop := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Нажмите Enter для завершения работы программы")
		scanner.Scan()

		log.Info("Получен сигнал завершения работы")
		close(stop)
	}()

	return stop
}

func loadEnv() (map[string]string, error) {
	env := make(map[string]string)

	file, err := os.Open(".env")
	if err != nil {
		return nil, fmt.Errorf("Не удалось открыть .env файл: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, "=")
		env[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return env, nil
}

func main() {
	var env map[string]string
	env, err := loadEnv()
	if err != nil {
		fmt.Println(err)
		env["JSONPATH"] = "./storage/quotes.json"
		env["PORT"] = "8080"
	}

	log, err := logger.NewLogger()
	if err != nil {
		fmt.Printf("Ошибка инициализации логгера: %v\n", err)
		return
	}
	defer log.Close()

	log.Info("Запуск сервера")

	storage, err := storage.CreateJSONStorage(env["JSONPATH"], log)
	if err != nil {
		log.Error(fmt.Sprintf("Не удалось инициализировавть хранилище: %v", err))
		return
	}

	rand.Seed(time.Now().UnixNano())

	stop := WaitClose(log)

	defer func() {
		if err = storage.Save(env["JSONPATH"], log); err != nil {
			log.Error(fmt.Sprintf("Не удалось сохранить данные: %v", err))
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/quotes", handlers.HandlerQuotesPost(storage, log)).Methods("POST")
	r.HandleFunc("/quotes", handlers.HandlerQuotesGet(storage, log)).Methods("GET")
	r.HandleFunc("/quotes/random", handlers.HandlerQuotesRandomGet(storage, log)).Methods("GET")
	r.HandleFunc("/quotes/{id}", handlers.HandlerQuotesDelete(storage, log)).Methods("DELETE")

	go func() {
		if err := http.ListenAndServe(":"+env["PORT"], r); err != nil {
			log.Error(fmt.Sprintf("Ошибка при запуске сервера: %v", err))
		}
	}()

	<-stop
	log.Info("Завершение работы программы")
}
