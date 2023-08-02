package main

import (
	"log"
	"os"
)

type loggerInterface interface {
	Log(string)
}

type logger struct {
}

func newLogger() *logger {
	l := &logger{}
	return l
}

func (l *logger) Log(message string) {
	file, err := os.OpenFile("logfile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open logfile.log: ", err)
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println(message)
}
