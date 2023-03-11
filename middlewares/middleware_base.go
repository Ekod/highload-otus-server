package middlewares

import "go.uber.org/zap"

type Middleware struct {
	Logger *zap.SugaredLogger
}
