package logger

import (
	"fmt"
	"github.com/multiplay/go-slack/chat"
	"github.com/multiplay/go-slack/lrhook"
	"github.com/multiplay/go-slack/webhook"
	"github.com/sirupsen/logrus"
)

type SlackConfig struct {
	user        string
	webHookURL  string
	channel     string
	machineName string
}

func (sc SlackConfig) NewFatalLoggerHook() *lrhook.Hook {
	cfg := lrhook.Config{
		MinLevel: logrus.FatalLevel,
		Async:    true,
		Message: chat.Message{
			Text:      sc.BaseMassage(),
			Channel:   sc.channel,
			Username:  "avv fatal logger",
			IconEmoji: ":avvfatal:",
			AsUser:    true,
		},
	}
	return lrhook.New(cfg, sc.webHookURL)
}

func (sc SlackConfig) BaseMassage() string {
	return fmt.Sprintf("<@%s> こちらはavvコマンドです\nマシン名:%s で実行していたお仕事がおわりました", sc.user, sc.machineName)
}

func (sc SlackConfig) PostMessage(text string) {
	hook := webhook.New(sc.webHookURL)
	m := &chat.Message{
		Text:      sc.BaseMassage() + "\n" + text,
		Channel:   sc.channel,
		Username:  "avv",
		IconEmoji: ":avv:",
		AsUser:    true,
	}
	res, err := m.Send(hook)
	if err != nil || !res.Ok() {
		logrus.WithError(err).Error("Failed post message to slack")
	} else {
		logrus.Info("Post Message to ", sc.channel)
	}
}
