package logging

import (
	"log"
	"os"
)

func SetupLogging() (*os.File, error) {
	logFile, err := os.OpenFile("bookstore.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(logFile)
	return logFile, nil
}
