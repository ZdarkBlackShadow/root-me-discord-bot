package config

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() *zap.Logger {
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.MkdirAll(logDir, 0755)
	}


	consoleEnabled := os.Getenv("CONSOLE_LOG_ENABLED") != "false"

	logLevelStr := os.Getenv("LOG_LEVEL")
	var level zapcore.Level
	switch strings.ToUpper(logLevelStr) {
	case "DEBUG":
		level = zap.DebugLevel
	case "INFO":
		level = zap.InfoLevel
	case "WARN":
		level = zap.WarnLevel
	case "ERROR":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	var cores []zapcore.Core

	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	logFile, err := os.OpenFile("logs/bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		writer := zapcore.AddSync(logFile)
		cores = append(cores, zapcore.NewCore(fileEncoder, writer, level))
	} else {
		fmt.Printf("Can't open log file : %v\n", err)
	}

	if consoleEnabled {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
	}

	combinedCore := zapcore.NewTee(cores...)

	return zap.New(combinedCore, zap.AddCaller())
}
