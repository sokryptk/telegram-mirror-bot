package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"strconv"
	"telegram-mirror-bot/utils/dryve"
	"telegram-mirror-bot/utils/parse"
)

func ListFiles(ctx *ext.Context) error {
	args := ctx.Args()
	if len(args) == 1 {
		return nil
	}

	files, err := dryve.List(args[1])
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(ctx.Bot, err.Error(), nil)
		if err != nil {
			return err
		}
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