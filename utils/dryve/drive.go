package dryve

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"google.golang.org/api/drive/v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"telegram-mirror-bot/utils/config"
)

const (
	FolderMimeType = "application/vnd.google-apps.folder"
)

func ParseMediaToUsableFormat(media drive.File, isFolder... bool) string {
	driveDomain := "https://drive.google.com"
	if len(isFolder) > 0 && isFolder[0] {
		return fmt.Sprintf("%s/drive/folders/%s", driveDomain, media.Id)
	}

	return fmt.Sprintf("%s/uc?id=%s", driveDomain, media.Id)
}

func createFile(name string, mimeType string, content io.Reader, parent... string) (*drive.File, error) {
	if len(parent) == 0 {
		parent = append(parent, config.C.Root)
	}

	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  parent,
	}
	file, err := S.Files.Create(f).Media(content).SupportsTeamDrives(true).Do()

	if err != nil {
		return nil, err
	}

	return file, nil
}

func createFolder(folderName string, parentID string) *drive.File {
	folder := &drive.File{
		Name: folderName,
		Parents: []string{
			parentID,
		},
		MimeType: FolderMimeType,
	}


	mkdir, _ := S.Files.Create(folder).SupportsTeamDrives(true).Do()

	return mkdir
}


func UploadFile(fpath string, folderID... string) (*drive.File, error) {
	var uploadFile *drive.File
	file, err := os.OpenFile(fpath, os.O_RDONLY, os.ModePerm)
	//noinspection ALL
	defer file.Close()

	if err != nil {
		return nil, fmt.Errorf("cancelled file")
	}

	mimeType, _ := mimetype.DetectFile(fpath)

	if len(folderID) == 0 {
		uploadFile, err = createFile(filepath.Base(fpath), mimeType.String(), file)
	} else {
		uploadFile, err = createFile(filepath.Base(fpath), mimeType.String(), file, folderID[0])
	}

	if err != nil {
		return nil, err
	}

	return uploadFile, nil
}

func UploadFolder(folderPath string) (*drive.File, error) {
	parentFolder := createFolder(filepath.Base(folderPath), config.C.Root)

	folderMap := map[string]string{
		filepath.Base(folderPath) : parentFolder.Id,
	}


	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if folderPath == path {
			return nil
		}

		if info.IsDir() {
			folderID := createFolder(filepath.Base(path), parentFolder.Id)

			folderMap[filepath.Base(path)] = folderID.Id

			return nil
		}


		if folderMap[filepath.Base(filepath.Dir(path))] != "" {
			fmt.Println(path)
			_, err  := UploadFile(path, folderMap[filepath.Base(filepath.Dir(path))])
			if err != nil {
				return err
			}
			log.Println("Uploaded")
		}

		return nil
	})


	if err != nil {
		return nil, err
	}

	return parentFolder, err
}