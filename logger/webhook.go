package logger

import (
	"fmt"
	"github.com/multiplay/go-slack/chat"
	"github.com/multiplay/go-slack/lrhook"
	"github.com/multiplay/go-slack/webhook"
	"github.com/sirupsen/logrus"
)

type SlackConfig struct {
	User        string
	WebHookURL  string
	Channel     string
	MachineName string
}

func (sc SlackConfig) NewFatalLoggerHook() *lrhook.Hook {
	cfg := lrhook.Config{
		MinLevel: logrus.FatalLevel,
		Async:    true,
		Message: chat.Message{
			Text:      sc.BaseMassage(),
			Channel:   sc.Channel,
			Username:  "avv fatal logger",
			IconEmoji: ":avvfatal:",
			AsUser:    true,
		},
	}
	return lrhook.New(cfg, sc.WebHookURL)
}

func (sc SlackConfig) BaseMassage() string {
<<<<<<< HEAD:logger/webhook.go
	return fmt.Sprintf("<@%s> こちらはavvコマンドです\nマシン名:%s で実行していたお仕事がおわりました", sc.User, sc.MachineName)
=======
	return fmt.Sprintf("<@%s> こちらはavvコマンドです\nマシン名:%s で実行していたお仕事がおわりました", sc.User, config.MachineName)
>>>>>>> master:cmd/webhook.go
}

func (sc SlackConfig) PostMessage(text string) {
	hook := webhook.New(sc.WebHookURL)
	m := &chat.Message{
		Text:      sc.BaseMassage() + "\n" + text,
		Channel:   sc.Channel,
		Username:  "avv",
		IconEmoji: ":avv:",
		AsUser:    true,
	}
	res, err := m.Send(hook)
	if err != nil || !res.Ok() {
		logrus.WithError(err).Error("Failed post message to slack")
	} else {
		logrus.Info("Post Message to ", sc.Channel)
	}
}
