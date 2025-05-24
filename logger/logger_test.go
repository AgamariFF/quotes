package logger_test

import (
	"os"
	"strings"
	"testing"

	"quotes/logger"
)

func TestLogger(t *testing.T) {
	log, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Не удалось создать логгер: %v", err)
	}
	defer log.Close()

	// Тест 1: Логирование Info
	infoMessage := "Информация..."
	log.Info(infoMessage)

	content, err := os.ReadFile("log.log")
	if err != nil {
		t.Fatalf("Не удалось прочитать файл: %v", err)
	}
	if !strings.Contains(string(content), infoMessage) {
		t.Errorf("Сообщение '%s' не найдено в логах: %s", infoMessage, string(content))
	}

	// Тест 2: Логирование Error
	errorMessage := "Ошибка..."
	log.Error(errorMessage)

	content, err = os.ReadFile("log.log")
	if err != nil {
		t.Fatalf("Не удалось прочитать файл: %v", err)
	}
	if !strings.Contains(string(content), errorMessage) {
		t.Errorf("Сообщение '%s' не найдено в логах: %s", errorMessage, string(content))
	}

	// Тест 3: Закрытие файла
	log.Close()
	err = log.File.Close()
	if err == nil {
		t.Error("Ожидалась ошибка при повторном закрытии файла")
	}
}
