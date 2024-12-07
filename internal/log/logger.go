package log

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

// InitializeLogger initializes the global logger
func InitializeLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Logger = logger.Sugar()
}

// Infof logs a message at InfoLevel
func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

// Infof logs a message at InfoLevel
func Infoln(args ...interface{}) {
	Logger.Infoln(args...)
}

// Error logs a message at ErrorLevel
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Fatal logs a message at ErrorLevel
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}
