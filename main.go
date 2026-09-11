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

	log.Info("Processing Folders")
	rootDir := os.Args[1]
	filesAndFolders, err := process.WalkDirectory(rootDir)
	if err != nil {
		log.Debug("unable to walk files within %s directory. details: %+v", rootDir, err)
		return
	}

	// remove unwanted folders
	removedFoldersCount, err := process.RemoveUnwantedFolders(filesAndFolders)
	if err != nil {
		log.Debug("unable to remove unwanted folders. details: %+v", err)
		return
	}
	log.Info("Removed %d unwanted folder(s)", removedFoldersCount)

	// remove matching unwanted parent folders. Eg, test for test.zip
	removeOwnerFolders, err := process.RemoveOwnerFolders(filesAndFolders)
	if err != nil {
		log.Debug("unable to remove owner folders. details: %+v", err)
		return
	}
	log.Info("Removed %d owner folder(s)", removeOwnerFolders)

	log.Info("Processing files")
	indiviualFiles, err := process.WalkFiles(rootDir)
	if err != nil {
		log.Debug("unable to walk files. details: %+v", err)
		return
	}

	if len(indiviualFiles) == 0 {
		log.Debug("no files detected to process")
		return
	}

	filesMap := process.GroupFilesBySize(indiviualFiles)

	log.Info("Removing duplicate files")
	removedDuplicateFilesCount, err := process.RemoveDuplicate(filesMap)
	if err != nil {
		log.Debug("failed to remove duplicate files. details: %+v", err)
		return
	}
	log.Info("Removed %d duplicate files", removedDuplicateFilesCount)

	log.Info("Removing unwanted files")
	removedFiles, err := process.RemoveUnwantedFiles(indiviualFiles)
	if err != nil {
		log.Debug("failed to remove unwanted files. details: %+v", err)
		return
	}

	log.Info("Removed %d unwanted file(s)", removedFiles)

}
