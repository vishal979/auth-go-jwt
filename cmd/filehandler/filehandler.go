package filehandler

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Open sets log output to a file
func Open() *os.File {
	logsFile, err := os.OpenFile("logs/server.logs", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("error opening logs file")
	}
	log.SetOutput(logsFile)
	log.Info("log file opened successfully")
	return logsFile
}
