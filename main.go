package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"log"
	"telegram-mirror-bot/telegram"
	"telegram-mirror-bot/utils/config"
)

func main() {
	b , err := gotgbot.NewBot(config.C.Bot)
	if err != nil {
		log.Fatal(err)
	}

	updater := ext.NewUpdater(b, nil)
	dispatcher := updater.Dispatcher

	dispatcher.AddHandler(handlers.NewCommand("mirror", telegram.Mirror))

	if err := updater.StartPolling(b, &ext.PollingOpts{Clean: true}); err != nil {
		log.Fatal(err)
	}

	log.Println("Mirror Bot Up")
	updater.Idle()
}