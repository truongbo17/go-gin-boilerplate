package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Open(env string) (*logrus.Logger, io.Closer) {
	logger := logrus.New()
	currentDate := time.Now()

	logger.SetLevel(logrus.InfoLevel)
	if env == config.DebugMode || env == config.LocalMode {
		logger.SetLevel(logrus.DebugLevel)
	}
	logger.SetFormatter(&logrus.JSONFormatter{})

	logFile := &lumberjack.Logger{
		Filename: fmt.Sprintf(config.PathLog, currentDate.Format(time.DateOnly)),
		MaxSize:  10,
		Compress: false,
	}

	writers := []io.Writer{logFile}
	if env != config.ReleaseMode {
		writers = append(writers, os.Stdout)
	}

	logger.SetOutput(io.MultiWriter(writers...))

	fmt.Println("Success init logger with Logrus")

	return logger, logFile
}
