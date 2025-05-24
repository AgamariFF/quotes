package logger

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	File     *os.File
	mu       sync.Mutex
}

func NewLogger() (*Logger, error) {
	file, err := os.OpenFile("log.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		InfoLog:  log.New(file, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(file, "ERROR\t", log.Ldate|log.Ltime),
		File:     file,
	}, nil
}

func (l *Logger) Info(str string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.InfoLog.Println(str)
}

func (l *Logger) Error(str string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ErrorLog.Println(str)
}

func (l *Logger) Close() {
	if l.File != nil {
		l.File.Close()
	}
}
