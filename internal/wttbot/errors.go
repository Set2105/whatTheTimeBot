package wttbot

import (
	"context"

	"github.com/go-telegram/bot"
)

func SendError(chatId, threadId int, err error, ctx context.Context, b *bot.Bot) {
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:          chatId,
			MessageThreadID: threadId,
			Text:            err.Error(),
		})
	}
}
