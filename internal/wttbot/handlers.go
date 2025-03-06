package wttbot

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/set2105/whatTheTimeBot/internal/loc"
)

func (wtt *WhatTheTimeBot) RegisterHandlers() {
	wtt.bot.RegisterHandler(bot.HandlerTypeMessageText, "/add", bot.MatchTypePrefix, wtt.AddHandler)
	wtt.bot.RegisterHandler(bot.HandlerTypeMessageText, "/del", bot.MatchTypePrefix, wtt.DeleteHandler)
	wtt.bot.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypePrefix, wtt.ListHandler)
	wtt.bot.RegisterHandler(bot.HandlerTypeMessageText, "/time", bot.MatchTypePrefix, wtt.TimeHandler)
	wtt.bot.RegisterHandler(bot.HandlerTypeMessageText, "/datetime", bot.MatchTypePrefix, wtt.DateTimeHandler)
	wtt.bot.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypePrefix, wtt.HelpHandler)
}

func (wtt *WhatTheTimeBot) AddHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	threadId := update.Message.MessageThreadID

	args, err := parseArgs(update.Message.Text, 2)
	if err != nil {
		SendError(int(chatId), threadId, err, ctx, b)
		return
	}
	locName := args[0]
	locLocale := args[1]
	ltm, err := wtt.configDispatcher.Get(int(chatId))
	if err != nil {
		SendError(int(chatId), threadId, err, ctx, b)
		return
	}

	newLoc, err := loc.Init(locName, locLocale)
	if err != nil {
		SendError(int(chatId), threadId, err, ctx, b)
		return
	}

	ltm.Add(newLoc)
	err = wtt.configDispatcher.Save(int(chatId))
	if err != nil {
		SendError(int(chatId), threadId, err, ctx, b)
		return
	}
}

func (wtt *WhatTheTimeBot) DeleteHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	threadId := update.Message.MessageThreadID

	args, err := parseArgs(update.Message.Text, 1)
	if err != nil {
		SendError(int(chatId), threadId, err, ctx, b)
		return
	}
	locName := args[0]
	ltm, err := wtt.configDispatcher.Get(int(chatId))
	if err != nil {
		SendError(int(chatId), threadId, err, ctx, b)
		return
	}
	ltm.Delete(locName)
	err = wtt.configDispatcher.Save(int(chatId))
	if err != nil {
		SendError(int(chatId), threadId, err, ctx, b)
		return
	}
}

func (wtt *WhatTheTimeBot) ListHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	threadId := update.Message.MessageThreadID

	ltm, err := wtt.configDispatcher.Get(int(chatId))
	if err != nil {
		SendError(int(chatId), threadId, err, ctx, b)
		return
	}
	lts := ltm.Values()
	msg := "Locales:\n"
	for i, lt := range lts {
		msg += fmt.Sprintf("%d)%s:%s\n", i+1, lt.Name, lt.LocationName)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:          chatId,
		MessageThreadID: threadId,
		Text:            msg,
	})
}

func (wtt *WhatTheTimeBot) TimeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	threadId := update.Message.MessageThreadID

	ltm, err := wtt.configDispatcher.Get(int(chatId))
	if err != nil {
		SendError(int(chatId), threadId, err, ctx, b)
		return
	}

	args, err := parseVariousArgs(update.Message.Text, 0, 2)
	if err != nil {
		SendError(int(chatId), threadId, err, ctx, b)
		return
	}

	switch len(args) {
	case 0:
		ltm.Time = time.Now()
	case 2:
		locName := ""
		l, err := ltm.GetByPosition(args[1])
		if err == nil {
			locName = l.LocationName
		} else {
			locName = args[1]
		}
		t, err := ParseDateTime("", args[0], locName)
		if err != nil {
			SendError(int(chatId), threadId, err, ctx, b)
			return
		}
		ltm.Time = *t
	}
	lts := ltm.Values()
	msg := ""
	for _, lt := range lts {
		msg += fmt.Sprintf("%s: %s\n", lt.Name, lt.TimeString())
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:          chatId,
		MessageThreadID: threadId,
		Text:            msg,
	})
}

func (wtt *WhatTheTimeBot) DateTimeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	threadId := update.Message.MessageThreadID

	ltm, err := wtt.configDispatcher.Get(int(chatId))
	if err != nil {
		SendError(int(chatId), threadId, err, ctx, b)
		return
	}

	args, err := parseVariousArgs(update.Message.Text, 0, 3)
	if err != nil {
		SendError(int(chatId), threadId, err, ctx, b)
		return
	}

	switch len(args) {
	case 0:
		ltm.Time = time.Now()
	case 3:
		locName := ""
		l, err := ltm.GetByPosition(args[2])
		if err == nil {
			locName = l.LocationName
		} else {
			locName = args[2]
		}
		t, err := ParseDateTime(args[0], args[1], locName)
		if err != nil {
			SendError(int(chatId), threadId, err, ctx, b)
			return
		}
		ltm.Time = *t
	}
	lts := ltm.Values()
	msg := ""
	for _, lt := range lts {
		msg += fmt.Sprintf("%s:  %s  %s\n", lt.Name, lt.DateString(), lt.TimeString())
	}
	msg += ""

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:          chatId,
		MessageThreadID: threadId,
		Text:            msg,
	})
}

func (wtt *WhatTheTimeBot) HelpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	threadId := update.Message.MessageThreadID

	msg := `
<pre>
/add [name] [locale]
	add new locale to locales
	examples
		/add Tbilisi Asia/Tbilisi
		/add "Tbilisi the capital of Georgia" Asia/Tbilisi
-------------------------------------------------------------
/del [name]
	delete locale from locales
	example
		/del Tbilisi
-------------------------------------------------------------
/list
	list locales
-------------------------------------------------------------
/time
	send current time in locales
-------------------------------------------------------------
/time [HH:MM] [locale|locale id from list]
	send hours and minutes in locales
	examples
		/time 12:00 Asia/Tbilisi
		/time 12:00 1
-------------------------------------------------------------
/datetime 
	send current date and time  in locales
-------------------------------------------------------------
/datetime [MM-DD-YYYY] [HH:MM] [locale|locale id from list]
	send date and time in locales
	example
		/datetime 05-04-2025 16:05 Asia/Tbilisi
		/datetime 05-04-2025 16:05 1
	</pre>`
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:          chatId,
		MessageThreadID: threadId,
		Text:            msg,
		ParseMode:       models.ParseModeHTML,
	})
}
