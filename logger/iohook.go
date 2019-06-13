package hook

import (
	"github.com/sirupsen/logrus"
	"io"
)

type IOHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

func (hook *IOHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	_, err = hook.Writer.Write([]byte(line))
	return err
}

func (hook *IOHook) Levels() []logrus.Level {
	return hook.LogLevels
}
