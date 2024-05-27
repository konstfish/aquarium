package controllers

import (
	"log"
	"os"
	"path/filepath"
)

var Sprites map[string][]string

func InitSprites(rootDir string) {
	Sprites = GetSprites(rootDir)
}

func GetSprites(rootDir string) map[string][]string {
	parentFolders, err := os.ReadDir(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	var folders map[string][]string = make(map[string][]string)

	for _, folder := range parentFolders {
		if folder.IsDir() {
			folderPath := filepath.Join(rootDir, folder.Name())

			files, err := os.ReadDir(folderPath)
			if err != nil {
				log.Fatal(err)
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
