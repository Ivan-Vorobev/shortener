package logger

import (
	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	// создаём новую конфигурацию логера
	zl, err := zap.NewDevelopment(zap.IncreaseLevel(zap.InfoLevel))

	if err != nil {
		return nil, err
	}

	return zl, nil
}
