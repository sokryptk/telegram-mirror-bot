package cache

type cache map[int64][]string

var DB cache

func RemoveGIDFromChat(chatID int64, gid string) error {
	for i, k := range DB[chatID] {
		if k == gid {
			DB[chatID] = append(DB[chatID][:i], DB[chatID][i:]...)
			return nil
		}
	}

	return nil
}

func AddGIDToChat(chatID int64, gid string) error {
	(DB)[chatID] = append(DB[chatID], gid)

	return nil
}