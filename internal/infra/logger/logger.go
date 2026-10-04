package logger

import "github.com/sirupsen/logrus"

var LogrusLogger *logrus.Logger

func Init() {
	LogrusLogger = InitLog()
}
