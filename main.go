package main

import (
	"os"

	log "github.com/earmuffjam/smart-scan/log"
	"github.com/earmuffjam/smart-scan/process"
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
	log.Info("Hunting for unwanted folders")
	files, err := process.WalkDirectory(rootDir)
	if err != nil {
		log.Debug("unable to walk files within %s directory. details: %+v", rootDir, err)
		return
	}

	// remove unwanted folders
	removedFoldersCount, err := process.RemoveUnwantedFolders(files)
	if err != nil {
		log.Debug("unable to remove unwanted folders. details: %+v", err)
		return
	}
	log.Info("Removed %d unwanted folder(s)", removedFoldersCount)

	// remove matching unwanted parent folders. Eg, test for test.zip
	removeOwnerFolders, err := process.RemoveOwnerFolders(files)
	if err != nil {
		log.Debug("unable to remove owner folders. details: %+v", err)
		return
	}
	log.Info("Removed %d owner folder(s)", removeOwnerFolders)

	log.Info("Hunting for duplicate files")
	filesMap, err := process.WalkFiles(rootDir)
	if err != nil {
		log.Debug("unable to walk files. details: %+v", err)
		return
	}

	if len(filesMap) == 0 {
		log.Debug("no files detected to process")
		return
	}

	removedFileGroup, err := process.RemoveDuplicate(filesMap)
	if err != nil {
		log.Debug("failed to remove duplicate files. details: %+v", err)
		return
	}
	log.Info("Removed files within %d duplicate groups", removedFileGroup)

	log.Info("Hunting for unwanted files")
	removedFiles, err := process.RemoveUnwantedFiles(filesMap)
	if err != nil {
		log.Debug("failed to remove unwanted files. details: %+v", err)
		return
	}

	log.Info("Removed %d unwanted file(s)", removedFiles)

}
