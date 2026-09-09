package main

import (
	"os"

	processdir "github.com/earmuffjam/smart-scan/processDir"
	"github.com/earmuffjam/smart-scan/service"
	"github.com/earmuffjam/smart-scan/util"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		util.Error("unable to load env file. details: %+v", err)
		return
	}
	util.Info("Application smart-scan running ...")

	dirName := os.Args
	if len(dirName) < 2 {
		util.Debug("Invalid usage. Usage: scan <directory>")
		return
	}

	rootDir := os.Args[1]
	filesMap, err := processdir.WalkDir(rootDir)
	if err != nil {
		util.Debug("unable to process selected directories")
		return
	}

	if len(filesMap) == 0 {
		util.Debug("no files detected to process")
		return
	}

	groupDuplicateFiles, err := service.DetectDuplicates(filesMap)
	if err != nil {
		util.Debug("unable to detect duplicates. details: %+v", err)
		return
	}

	if len(groupDuplicateFiles) == 0 {
		util.Info("no duplicate files found.")
		return
	}

	err = service.MoveToTrash(groupDuplicateFiles)
}
