package main

import (
	"os"

	log "github.com/earmuffjam/smart-scan/log"
	processdir "github.com/earmuffjam/smart-scan/processDir"
	"github.com/earmuffjam/smart-scan/service"
	"github.com/earmuffjam/smart-scan/types"
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
	filesMap, err := processdir.WalkDir(rootDir)
	if err != nil {
		log.Debug("unable to process selected directories")
		return
	}

	if len(filesMap) == 0 {
		log.Debug("no files detected to process")
		return
	}

	err = DetectAndDeleteDuplicatesFiles(filesMap)
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

// DetectAndDeleteDuplicatesFiles ...
// used to detect and remove duplicate values
func DetectAndDeleteDuplicatesFiles(filesMap map[int64][]types.File) error {
	groupDuplicateFiles, err := service.DetectDuplicates(filesMap)
	if err != nil {
		log.Debug("unable to detect duplicates. details: %+v", err)
		return err
	}

	log.Info("Found %d duplicate files. Processing ...", len(groupDuplicateFiles))

	if len(groupDuplicateFiles) == 0 {
		log.Debug("no duplicate files found.")
		return nil
	}

	err = service.MoveToTrash(groupDuplicateFiles)
	if err != nil {
		log.Debug("failed to remove duplicate files")
		return err
	}

	return nil
}
