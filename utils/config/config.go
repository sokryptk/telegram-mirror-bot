package config

import (
	"log"
	"os"
	"telegram-mirror-bot/utils/erst"
)

/*
Bot : Telegram Bot Token
Root : Drive Folder ID To Upload
Index : Google Drive File Index Url
Dir : Download file directory
 */
type Configuration struct {
	Bot string `json:"bot"`
	Root string `json:"root"`
	Index string `json:"index"`
	Dir string `json:"dir"`
	Interval string `json:"interval"`
}

var C = &Configuration{}

func init() {
	if bot, isEnv := os.LookupEnv("BOT"); isEnv {
		root, hasRoot := os.LookupEnv("ROOT")
		if !hasRoot {
			log.Fatal(erst.GenMissingFlagErr("ROOT", "Drive Root Folder ID"))
		}

		dir, hasDir := os.LookupEnv("DIR")
		if !hasDir {
			log.Fatal(erst.GenMissingFlagErr("DIR", "DIR to download to"))
		}


		interval, hasInterval := os.LookupEnv("INT")
		if !hasInterval {
			//default interval is 2seconds
			C.Interval = "2s"
		}

		C.Interval = interval
		C.Bot = bot
		C.Root = root
		C.Index = os.Getenv("INDEX") //We will check if the index exists in the the upload function
		C.Dir = dir
	}
}