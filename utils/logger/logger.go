package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

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
