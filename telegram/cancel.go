package telegram

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"strings"
	"telegram-mirror-bot/utils/aria"
	"telegram-mirror-bot/utils/erst"
)

func CancelMirror(ctx *ext.Context) error {
	args := ctx.Args()
	if len(args) == 1 {
		return nil
	}

	var cancelledGids []string
	for _, gids := range args[1:] {
		if err := aria.Cancel(gids); err != nil {
			log.Println(err)
			continue
		}

		cancelledGids = append(cancelledGids, gids)
	}

	if len(cancelledGids) == 0 {
		_, err := ctx.EffectiveMessage.Reply(ctx.Bot, erst.NoGidFound, nil)
		if err != nil {
			log.Println(err)
		}

		return nil
	}

	downloadText := fmt.Sprintf("Download(s) : %s were cancelled sucessfully.", strings.Join(cancelledGids, ","))
	_, err := ctx.EffectiveMessage.Reply(ctx.Bot, downloadText, nil)
	if err != nil {
		log.Println(err)
	}

	return nil
}