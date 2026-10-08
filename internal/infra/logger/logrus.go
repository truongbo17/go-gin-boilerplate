package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

func InitLog() *logrus.Logger {
	logger := logrus.New()
	currentDate := time.Now()

	logger.SetLevel(logrus.InfoLevel)
	if config.EnvConfig.App.Env == config.DebugMode || config.EnvConfig.App.Env == config.LocalMode {
		logger.SetLevel(logrus.DebugLevel)
	}
	logger.SetFormatter(&logrus.JSONFormatter{})

	logFile := &lumberjack.Logger{
		Filename: fmt.Sprintf(config.PathLog, currentDate.Format(time.DateOnly)),
		MaxSize:  10,
		Compress: false,
	}

	writers := []io.Writer{logFile}
	if config.EnvConfig.Env != config.ReleaseMode {
		writers = append(writers, os.Stdout)
	}

	logger.SetOutput(io.MultiWriter(writers...))

	fmt.Println("Success init logger with Logrus")

	return logger
}
