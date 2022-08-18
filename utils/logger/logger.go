package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

// New constructs a Sugared Logger that writes to stdout and
// provides human-readable timestamps.
func New(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]any{
		"service": service,
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
