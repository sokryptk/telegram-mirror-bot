package cache

import (
	"telegram-mirror-bot/utils/aria/ariaStatus"
	"telegram-mirror-bot/utils/parse"
)

type cache map[int64][]parse.TelegramMirrorInfo

var db = make(cache)

func RemoveGIDFromChat(chatID int64, gid string)  {
	for i, k := range db[chatID] {
		if k.Gid == gid {
			db[chatID] = append(db[chatID][:i], db[chatID][i:]...)
		}
	}
}

func Set(chatID int64, tmi parse.TelegramMirrorInfo) []parse.TelegramMirrorInfo {
	if i, ok := isInfoInCache(tmi, db[chatID]); ok {
		db[chatID][i] = tmi
		if tmi.Status == ariaStatus.DONE || tmi.Status == ariaStatus.REMOVED {
			db[chatID] = append(db[chatID][:i], db[chatID][i+1:]...)
		}
	} else {
		db[chatID] = append(db[chatID], tmi)
	}


	return db[chatID]
}

func isInfoInCache(tmi parse.TelegramMirrorInfo, tmis []parse.TelegramMirrorInfo) (int, bool) {
	for i, t := range tmis {
		if t.Gid == tmi.Gid {
			return i, true
		}
	}

	return 0, false
}