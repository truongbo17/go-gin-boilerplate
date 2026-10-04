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

type Formatter struct{}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	logLine := fmt.Sprintf("[%s] %s.%s: %s",
		entry.Time.Format(time.DateTime),
		config.EnvConfig.Env,
		entry.Level.String(),
		entry.Message,
	)

	if len(entry.Data) > 0 {
		logLine += fmt.Sprintf(" fields=%v", entry.Data)
	}

	logLine += "\n"
	return []byte(logLine), nil
}

func InitLog() *logrus.Logger {
	logger := logrus.New()
	currentDate := time.Now()

	logger.SetLevel(logrus.TraceLevel)
	logger.SetFormatter(&Formatter{})

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
