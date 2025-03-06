package main

import (
	"os"

	"github.com/set2105/whatTheTimeBot/internal/config"
	"github.com/set2105/whatTheTimeBot/internal/wttbot"
)

func main() {
	key := os.Getenv("TG_BOT_KEY")
	if key == "" {
		panic("TG_BOT_KEY is not set")
	}
	storage, err := config.InitConfigStorage("", "./config")
	if err != nil {
		panic(err)
	}
	disp := wttbot.InitConfigDispatcher(storage)
	bot, err := wttbot.Init(key, disp)
	if err != nil {
		panic(err)
	}
	for true {
		bot.Start()
	}
}
