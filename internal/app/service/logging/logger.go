package logging

import (
	"go.uber.org/zap"
	"log"
	"sync"
)

var once sync.Once
var Sugar *zap.SugaredLogger

func init() {
	once.Do(func() {
		Sugar = newConsoleLogger().Sugar()
	})
}

// newConsoleLogger
func newConsoleLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot initialize zap")
	}

	return logger
}
