package telegram

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"os"
	"path/filepath"
	"strings"
	"telegram-mirror-bot/utils/aria"
	"telegram-mirror-bot/utils/aria/ariaStatus"
	"telegram-mirror-bot/utils/dryve"
	"telegram-mirror-bot/utils/erst"
	"telegram-mirror-bot/utils/parse"
	"time"
)

func Mirror(ctx *ext.Context) error {
	args := ctx.Args()
	if len(args) == 1 {
		_, _ = ctx.EffectiveMessage.Reply(ctx.Bot, erst.NoUriToDownload, &gotgbot.SendMessageOpts{DisableNotification: true})
		return nil
	}

	gid, er := aria.Download(args[1])
	if er != nil {
		_, _ = ctx.EffectiveMessage.Reply(ctx.Bot, er.Error(), &gotgbot.SendMessageOpts{DisableNotification: true})
		return nil
	}

	m, _ := ctx.EffectiveMessage.Reply(ctx.Bot, "Downloading...", nil)

	go func(gid string, bot *gotgbot.Bot, message *gotgbot.Message) {
	DownloadLoop:
		for {
			si := aria.GetStatus(gid)
			tmi := parse.TelegramMirrorInfo{}
			tmi.ParseStatus(si)

			switch si.Status {
			case ariaStatus.ACTIVE:
				_, err := m.EditText(bot, tmi.FormatInfo(), &gotgbot.EditMessageTextOpts{ParseMode: "html"})
				if err != nil {
					log.Println(err)
				}
				time.Sleep(2 * time.Second)
			case ariaStatus.COMPLETE:
				//Probably a magnet
				if len(si.FollowedBy) > 0 {
					gid = si.FollowedBy[0]
					continue
				}

				_, _ = m.EditText(bot, fmt.Sprintf("Uploading %s", tmi.FileName), nil)
				break DownloadLoop
			}
		}

		st := aria.GetStatus(gid)
		basePath := strings.Replace(st.Files[0].Path, fmt.Sprintf("%s/", st.Dir), "", 1)
		//Concluding that it indeed is a file.
		// Since we would have had something else as the result from
		// filepath.Dir(basePath)
		if filepath.Dir(basePath) == "." {
			f, err := dryve.UploadFile(st.Files[0].Path)
			if err != nil {
				_, err = m.EditText(bot, err.Error(), nil)
			}

			parsedLink := parse.ConvertLinks(filepath.Base(st.Files[0].Path), dryve.ParseMediaToUsableFormat(*f))
			_, err = m.EditText(bot, parsedLink, &gotgbot.EditMessageTextOpts{ParseMode: "HTML"})
			if err != nil {
				log.Println(err)
			}


			if err := os.Remove(st.Files[0].Path); err != nil {
				log.Println(err)
			}

		} else {
			folderName := strings.Split(basePath, string(filepath.Separator))[1]
			folderPath := fmt.Sprintf("%s%s%s", st.Dir, string(filepath.Separator), folderName)
			folder, err := dryve.UploadFolder(folderPath)

			if err != nil {
				_, _ = ctx.EffectiveMessage.Reply(ctx.Bot, err.Error(), nil)
			}

			parsedLink := parse.ConvertLinks(folderName, dryve.ParseMediaToUsableFormat(*folder, true))
			_, err = m.EditText(bot, parsedLink, &gotgbot.EditMessageTextOpts{ParseMode: "HTML"})
			if err != nil {
				log.Println(err)
			}

			if err := os.RemoveAll(folderPath); err != nil {
				log.Println(err)
			}
		}

	}(gid, ctx.Bot, m)

	return nil
}
