package logger

import (
	"log"
	"os"
)

func CreateFile() *os.File {
	logFile, errFile := os.Create(os.Getenv("LOG_FILENAME"))

	if errFile != nil {
		log.Fatalf("Error creating file: %v", errFile)
		return nil
	}
	return logFile
}
