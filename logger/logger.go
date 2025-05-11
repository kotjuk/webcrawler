package logger

import (
	"log"
	"os"
)

var (
	logger *log.Logger
)

func Init() {
	file, err := os.OpenFile("crawler.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(file, "", log.Ldate|log.Ltime)
}

func Log(message string) {
	logger.Println(message)
}
