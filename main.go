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
		util.Error("unable to load .env: %+v", err)
		return
	}
	util.Info("Application smart-scan running ...")

	dirName := os.Args
	if len(dirName) < 2 {
		util.Error("Invalid usage. Usage: scan <directory>")
		return
	}

	rootDir := os.Args[1]
	groupedFiles, err := processdir.WalkDir(rootDir)
	if err != nil {
		util.Error("unable to process selected directories")
		return
	}

	groupByHash, err := service.DetectDuplicates(groupedFiles)
	if err != nil {
		util.Error("unable to detect duplicates. details: %+v", err)
		return
	}

	err = service.MoveToTrash(groupByHash)
}
