package logger_setup

import (
	"annotater/internal/config"
	"errors"
	"net"
	"os"

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
		if errors.Is(err, &netErr) { // error with gives connection refused
			entry.Level = logrus.ErrorLevel
		}

	}
	return nil
}

func Setuplog(conf *config.Config) *logrus.Logger {

	log := logrus.New()
	useFile := conf.Logger.UseFile
	if useFile {
		log.Printf("using file %s\n", conf.OutputFilePath)
		f, err := os.OpenFile(conf.OutputFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Printf("Failed to create logfile %s:%s, defaulting to stderr", conf.OutputFilePath, err.Error())
			useFile = false
		}

		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stderr)
	}

	easyFormatter := &easy.Formatter{
		TimestampFormat: conf.TimestampFormat,
		LogFormat:       conf.LogFormat,
	}

	logLevel, _ := logrus.ParseLevel(conf.LogLevel)
	log.SetFormatter(easyFormatter)
	/*if conf.OutputFormat == "text" {
		log.SetFormatter(&logrus.TextFormatter{})
	}
	if conf.OutputFormat == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	}*/
	log.SetReportCaller(true)
	log.AddHook(&DatabaseRefusedConnHook{})
	log.SetLevel(logLevel)
	return log
}
