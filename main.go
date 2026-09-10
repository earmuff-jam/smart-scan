package main

import (
	"os"

	log "github.com/earmuffjam/smart-scan/log"
	"github.com/earmuffjam/smart-scan/process"
	"github.com/earmuffjam/smart-scan/service"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Error("unable to load env file. details: %+v", err)
		return
	}

	log.Info("Application smart-scan running ...")

	dirName := os.Args
	if len(dirName) < 2 {
		log.Debug("Invalid usage. Usage: scan <directory>")
		return
	}

	rootDir := os.Args[1]
	err := process.RemoveUnwantedFolders(rootDir)
	if err != nil {
		log.Debug("unable to filter unwanted archives. details: %+v", err)
		return
	}

	filesMap, err := process.WalkDir(rootDir)
	if err != nil {
		log.Debug("unable to process selected directories")
		return
	}

	if len(filesMap) == 0 {
		log.Debug("no files detected to process")
		return
	}

	err = service.DetectAndDeleteDuplicatesFiles(filesMap)
	if err != nil {
		log.Debug("failed to remove duplicate files. details: %+v", err)
		return
	}

	err = service.RemoveUnwantedFiles(filesMap)
	if err != nil {
		log.Debug("failed to remove unwanted files. details: %+v", err)
		return
	}

}
