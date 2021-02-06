package telegram

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"strconv"
	"strings"
	"telegram-mirror-bot/utils/dryve"
	"telegram-mirror-bot/utils/parse"
)

func ListFiles(ctx *ext.Context) error {
	args := ctx.Args()
	if len(args) == 1 {
		return nil
	}

	query := strings.Join(args[1:], " ")
	files, err := dryve.List(query)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(ctx.Bot, err.Error(), nil)
		if err != nil {
			return err
		}

		return err
	}

	if len(files) == 0 {
		_, err := ctx.EffectiveMessage.Reply(ctx.Bot, fmt.Sprintf("No results found for %s", query), nil)
		if err != nil {
			log.Println(err)
		}

		return err
	}

	var finMessage string
	for _, f := range files {
		//Check if folder or file and adjust accordingly
		if f.MimeType == dryve.FolderMimeType {
			finMessage += parse.ConvertLinks(f.Name, dryve.ParseMediaToUsableFormat(*f, true))
		} else {
			finMessage += parse.ConvertLinks(f.Name, dryve.ParseMediaToUsableFormat(*f, false), parse.BytesToHumanReadable(strconv.Itoa(int(f.Size))))
		}

		// Add a newLine after every link
		finMessage += "\n"
	}

	_, err = ctx.EffectiveMessage.Reply(ctx.Bot, finMessage, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	if err != nil {
		log.Println(err)
	}

	return nil
}