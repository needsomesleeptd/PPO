package logger_setup

import (
	"annotater/internal/config"
	"errors"
	"net"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var UnableToDecodeUserReqF = "unable to decode user req:%v"
var UnableToGetUserifF = "cannot get userID from jwt %v in middleware"

type DatabaseRefusedConnHook struct {
}

func (hook *DatabaseRefusedConnHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.WarnLevel, logrus.InfoLevel}
}

func (hook *DatabaseRefusedConnHook) Fire(entry *logrus.Entry) error {
	var netErr net.OpError
	if err, ok := entry.Data["error"].(error); ok {
		if errors.Is(err, &netErr) || errors.Is(err, syscall.ECONNREFUSED) { // error with gives connection refused
			entry.Level = logrus.ErrorLevel
		}

	}
	return nil
}

func Setuplog(conf *config.Config) *logrus.Logger {

	logger := logrus.New()
	useFile := conf.Logger.UseFile
	if useFile {
		logger.Printf("using file %s\n", conf.OutputFilePath)
		f, err := os.OpenFile(conf.OutputFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			logger.Printf("Failed to create logfile %s:%s, defaulting to stderr", conf.OutputFilePath, err.Error())
			useFile = false
		}

		logger.SetOutput(f)
	} else {
		logger.SetOutput(os.Stderr)
	}

	easyFormatter := &easy.Formatter{
		TimestampFormat: conf.TimestampFormat,
		LogFormat:       conf.LogFormat,
	}

	logLevel, _ := logrus.ParseLevel(conf.LogLevel)
	logger.SetFormatter(easyFormatter)
	if conf.OutputFormat == "text" {
		logger.SetFormatter(&logrus.TextFormatter{
			QuoteEmptyFields: true,
		})
	}
	if conf.OutputFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	logger.SetReportCaller(true)
	logger.AddHook(&DatabaseRefusedConnHook{})
	logger.SetLevel(logLevel)
	return logger
}
