package wttbot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

type WhatTheTimeBot struct {
	bot              *bot.Bot
	configDispatcher *ConfigDispatcher
}

func (wtt *WhatTheTimeBot) Start() {
	go wtt.configDispatcher.StartClear()
	wtt.bot.Start(context.Background())
}

func newWTT(token string) (*bot.Bot, error) {
	opts := []bot.Option{}
	wttBot, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}
	return wttBot, nil
}

func Init(token string, configDispatcher *ConfigDispatcher) (*WhatTheTimeBot, error) {
	if configDispatcher == nil {
		return nil, fmt.Errorf("configDispatcher is nil")
	}
	wttBot, err := newWTT(token)
	if err != nil {
		return nil, err
	}
	b := WhatTheTimeBot{bot: wttBot, configDispatcher: configDispatcher}
	b.RegisterHandlers()
	return &b, nil
}
