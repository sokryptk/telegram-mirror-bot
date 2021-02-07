package config

import (
	"gopkg.in/ini.v1"
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
	Bot string `ini:"bot"`
	Root string `ini:"root"`
	Index string `ini:"index"`
	Dir string `ini:"dir"`
	Interval string `ini:"interval"`
}

var C = &Configuration{
	Dir: "",
	Interval: "2s",
}

const (
	configFileName = "config.ini"
)

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


		C.Interval = os.Getenv("INT")
		C.Bot = bot
		C.Root = root
		C.Index = os.Getenv("INDEX") //We will check if the index exists in the the upload function
		C.Dir = dir
		return
	}

	err := ini.MapTo(C, configFileName)
	if err != nil {
		log.Fatal(err)
	}
}