package io

import (
	"github.com/wiedzmin/goodies"
	"go.uber.org/zap"
)

var (
	rawLog *zap.Logger
	log    *zap.SugaredLogger
)

// GetLog ... get sugared logger
func GetLog() *zap.SugaredLogger {
	if log == nil {
		rawLog, log = goodies.InitLogging(false, false)
	}
	return log
}

// GetRawLog ... get raw logger
func GetRawLog() *zap.Logger {
	if rawLog == nil {
		rawLog, log = goodies.InitLogging(false, false)
	}
	return rawLog
}
