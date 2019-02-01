package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

//default setting
var StdoutLogger = &logrus.Logger{
	Out:       os.Stdout,
	Formatter: &logrus.JSONFormatter{PrettyPrint: true},
	Level:     logrus.DebugLevel,
}

func InitLogrusLogger(logLevel string, prettyLog bool) {
	StdoutLogger.SetOutput(os.Stdout)
	StdoutLogger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: prettyLog})
	SetLogLevel(StdoutLogger, logLevel)
}
