package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Info("Failed to log to file, using default stderr")
	} else {
		log.SetOutput(file)
	}
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.StandardLogger().Level)
}

func LogError(message string, err error) {
	log.Error(fmt.Sprintf("%s. Error: %s", message, err.Error()))
}

func LogErrorMessage(message string) {
	log.Error(fmt.Sprintf("%s.", message))
}

func LogInfo(message string) {
	log.Info(fmt.Sprintf("%s.", message))
}
