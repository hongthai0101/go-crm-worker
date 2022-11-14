package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func InitializeLogger() {
	zapLog, _ := zap.NewDevelopment()
	defer zapLog.Sync()
	Logger = zapLog.Sugar()
}
