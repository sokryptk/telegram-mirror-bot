package parse

import (
	"fmt"
	"github.com/zyxar/argo/rpc"
	"path/filepath"
	"strconv"
	"strings"
)

type TelegramMirrorInfo struct {
	Completed string `json:"completed"`
	Total     string `json:"total"`
	Speed     string `json:"speed"`
	ETA       string `json:"eta"`
	FileName  string `json:"file_name"`
}

func (tmi *TelegramMirrorInfo) ParseStatus(rpc rpc.StatusInfo) {
	tmi.Completed = BytesToHumanReadable(rpc.CompletedLength)
	tmi.Total = BytesToHumanReadable(rpc.TotalLength)
	tmi.Speed = BytesToHumanReadable(rpc.DownloadSpeed, true)
	tmi.ETA = CalculateETA(rpc)

	if rpc.Files[0].Path == "" || (strings.Contains(rpc.Files[0].Path, "METADATA")){
		tmi.FileName = "Downloading"
		return
	}

	//Will be a dir
	if len(rpc.Files) > 0 {
		folderBasePath := strings.Replace(rpc.Files[0].Path, rpc.Dir, "", 1)
		tmi.FileName = strings.Split(folderBasePath, string(filepath.Separator))[1]
	} else {
		tmi.FileName = filepath.Base(rpc.Files[0].Path)
	}
}

func BytesToHumanReadable(bytes string, speed ...bool) (hrf string) {
	bytesInt, _  := strconv.Atoi(bytes)

	switch {
	case bytesInt > (2 << 19):
		hrf = fmt.Sprintf("%dmb", bytesInt/(2 << 19))
	case bytesInt > (2 << 9):
		hrf = fmt.Sprintf("%dkb", bytesInt/(2 << 9))
	default:
		hrf = fmt.Sprintf("%db", bytesInt)
	}

	if len(speed) != 0 && speed[0] {
		hrf = fmt.Sprintf("%sps", hrf)
	}

	return
}

func SecondsToHumanReadable(seconds int) (hrf string) {
	switch {
	case seconds > (3600):
		hrf = fmt.Sprintf("%dh", seconds/3600)
	case seconds > (60):
		hrf = fmt.Sprintf("%dm", seconds/60)
	default:
		hrf = fmt.Sprintf("%ds", seconds)
	}

	return
}

func CalculateETA(rpc rpc.StatusInfo) string {
	tl, _ := strconv.Atoi(rpc.TotalLength)
	cl, _ := strconv.Atoi(rpc.CompletedLength)
	ds, _ := strconv.Atoi(rpc.DownloadSpeed)

	var eta int

	if ds != 0 {
		eta = (tl - cl)/ds
	} else {
		eta = 0
	}

	return SecondsToHumanReadable(eta)
}


