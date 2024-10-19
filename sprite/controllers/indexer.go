package controllers

import (
	"context"
	"os"
	"path/filepath"

	"github.com/konstfish/aquarium/common/logging"
)

var Sprites map[string][]string

func InitSprites(rootDir string) {
	Sprites = GetSprites(rootDir)
}

func GetSprites(rootDir string) map[string][]string {
	parentFolders, err := os.ReadDir(rootDir)
	if err != nil {
		logging.Error(context.TODO(), "Failed to read directory", err.Error())
		os.Exit(1)
	}

	var folders map[string][]string = make(map[string][]string)

	for _, folder := range parentFolders {
		if folder.IsDir() {
			folderPath := filepath.Join(rootDir, folder.Name())

			files, err := os.ReadDir(folderPath)
			if err != nil {
				logging.Error(context.TODO(), err.Error())
				os.Exit(1)
			}

			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".txt" {
					folders[folder.Name()] = append(folders[folder.Name()], file.Name())
				}
			}
		}
	}

	return folders
}
