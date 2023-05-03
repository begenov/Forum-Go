package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(msg string)
	Error(msg string, err error)
}

type logger struct {
	info *log.Logger
	e    *log.Logger
}

func NewLogger(l *log.Logger) Logger {
	_, err := os.Create("logs.txt")
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return &logger{
		info: log.New(file, "INFO: ", log.Ldate|log.Ltime),
		e:    log.New(file, "ERROR: ", log.Ldate|log.Ltime),
	}
}

func (l *logger) Info(msg string) {
	l.info.Println(msg)
}

func (l *logger) Error(msg string, err error) {
	l.e.Println(msg, err)
}
