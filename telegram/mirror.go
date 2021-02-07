package telegram

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"telegram-mirror-bot/cache"
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
			chatTmis := cache.Set(ctx.EffectiveChat.Id, tmi)

			switch si.Status {
			case ariaStatus.ACTIVE:
				_, err := m.EditText(bot, parse.Format(chatTmis), &gotgbot.EditMessageTextOpts{ParseMode: "html"})
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

				break DownloadLoop
			case ariaStatus.REMOVED:
				_, err := ctx.EffectiveMessage.Reply(ctx.Bot, "Download was cancelled", nil)
				if err != nil {
					log.Println(err)
				}

				_, err = m.Delete(ctx.Bot)
				if err != nil {
					log.Println(err)
				}

				//Cleanup
				basePath := strings.Replace(si.Files[0].Path, fmt.Sprintf("%s/", si.Dir), "", 1)
				if filepath.Dir(basePath) == "." {
					if err := os.Remove(si.Files[0].Path); err != nil {
						log.Println(err)
					}
				} else {
					folderName := strings.Split(basePath, string(filepath.Separator))[1]
					folderPath := fmt.Sprintf("%s%s%s", si.Dir, string(filepath.Separator), folderName)

					if err := os.RemoveAll(folderPath); err != nil {
						log.Println(err)
					}
				}
				//Cleanup end

				return
			}
		}

		st := aria.GetStatus(gid)
		tmi := parse.TelegramMirrorInfo{}
		tmi.ParseStatus(st)
		chatTmis := cache.Set(ctx.EffectiveChat.Id, tmi)

		_, err := m.EditText(bot, parse.Format(chatTmis), &gotgbot.EditMessageTextOpts{ParseMode: "html"})
		if err != nil {
			log.Println(err)
		}

		basePath := strings.Replace(st.Files[0].Path, fmt.Sprintf("%s/", st.Dir), "", 1)
		//Concluding that it indeed is a file.
		// Since we would have had something else as the result from
		// filepath.Dir(basePath)
		if filepath.Dir(basePath) == "." {
			f, err := dryve.UploadFile(st.Files[0].Path)
			if err != nil {
				_, err = m.EditText(bot, err.Error(), nil)
			}

			chatTmis = cache.Set(ctx.EffectiveChat.Id, parse.TelegramMirrorInfo{Status: "done"})
			_, err = m.EditText(bot, parse.Format(chatTmis), &gotgbot.EditMessageTextOpts{ParseMode: "html"})
			if err != nil {
				log.Println(err)
			}

			parsedLink := parse.ConvertLinks(filepath.Base(st.Files[0].Path), dryve.ParseMediaToUsableFormat(*f), parse.BytesToHumanReadable(strconv.Itoa(int(f.Size))))
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

			chatTmis = cache.Set(ctx.EffectiveChat.Id, parse.TelegramMirrorInfo{Status: "done"})
			_, err = m.EditText(bot, parse.Format(chatTmis), &gotgbot.EditMessageTextOpts{ParseMode: "html"})
			if err != nil {
				log.Println(err)
			}

			parsedLink := parse.ConvertLinks(folderName, dryve.ParseMediaToUsableFormat(*folder, true), parse.BytesToHumanReadable(strconv.Itoa(int(folder.Size))))
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
