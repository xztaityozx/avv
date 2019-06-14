package logger

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"io/ioutil"
	"path/filepath"
	"time"
)

func NewLogger(base string, config SlackConfig) *logrus.Logger {
	logDir := filepath.Join(base, "log")
	log := logrus.New()

	// Hook to log logDir
	path := filepath.Join(logDir, time.Now().Format("2006-01-02-15-04-05")+".log")
	logrus.Info("LogFile: ", path)
	fileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename: path,
		MaxAge:   28,
		MaxSize:  500,
		Level:    logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{
			ForceColors:     true,
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		},
	})
	if err != nil {
		logrus.Fatal(err)
	}

	log.SetOutput(ioutil.Discard)
	// init logrus System
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
	log.AddHook(&IOHook{
		LogLevels: []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel},
		Writer:    colorable.NewColorableStderr(),
	})

	log.AddHook(fileHook)

	// Slack Hook
	slackHook := config.NewFatalLoggerHook()
	log.AddHook(slackHook)

	return log
}
