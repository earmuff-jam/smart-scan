package service

import (
	"errors"
	"fmt"

	log "github.com/earmuffjam/smart-scan/log"
	"github.com/earmuffjam/smart-scan/process"
	"github.com/earmuffjam/smart-scan/trash"
	"github.com/earmuffjam/smart-scan/types"
)

// RemoveUnwantedFiles ...
// removes unwanted files slated for removal in env variables
func RemoveUnwantedFiles(groupedFiles map[int64][]types.File) error {
	removedCount := 0
	for _, files := range groupedFiles {
		if len(files) < 2 {
			continue
		}

		for _, file := range files {
			if !file.ValidatePrefix() {
				continue
			}

			log.Debug("moving file %s to trash.", file.Path)
			if err := trash.MoveToTrash(file.Path); err != nil {
				log.Debug(
					"unable to move file %s to trash. details: %+v",
					file.Path,
					err,
				)
				return errors.New("failed to remove file")
			}
			removedCount++
		}
	}
	log.Info("Removed %d file(s) from selected directory", removedCount)
	return nil
}

// DetectAndDeleteDuplicatesFiles ...
// used to detect and remove duplicate values
func DetectAndDeleteDuplicatesFiles(filesMap map[int64][]types.File) error {
	groupDuplicateFiles, err := detectDuplicates(filesMap)
	if err != nil {
		log.Debug("unable to detect duplicates. details: %+v", err)
		return err
	}

	log.Info("Found %d duplicate files. Processing ...", len(groupDuplicateFiles))

	if len(groupDuplicateFiles) == 0 {
		log.Debug("no duplicate files found.")
		return nil
	}

	err = moveDuplicateFilesToTrash(groupDuplicateFiles)
	if err != nil {
		log.Debug("failed to remove duplicate files")
		return err
	}

	return nil
}

// detectDuplicates ...
// defines a function that detects all duplicates
func detectDuplicates(groupedFiles map[int64][]types.File) (map[string][]types.File, error) {
	groupedByFileHash := make(map[string][]types.File)
	for _, files := range groupedFiles {
		if len(files) < 2 {
			continue
		}

		for _, file := range files {
			hash, err := process.BuildHashFile(file.Path)
			if err != nil {
				log.Debug(
					"failed to hash %s. details: %+v", file.Path, err)
				continue

			}
			// hash file with matching size
			groupedByFileHash[hash] = append(groupedByFileHash[hash], file)
		}
	}
	return groupedByFileHash, nil
}

// moveDuplicateFilesToTrash ...
// defines a function to move duplicates to trash
func moveDuplicateFilesToTrash(groupedByHash map[string][]types.File) error {

	for hash, files := range groupedByHash {
		if len(files) < 2 {
			continue
		}
		log.Debug("duplicate file group: %s", hash)

		// keep the first file
		for _, file := range files[1:] {
			log.Debug("moving file %s to trash.", file.Path)
			if err := trash.MoveToTrash(file.Path); err != nil {
				errorMsg := fmt.Sprintf("unable to move file %s to trash. details: %+v", file.Path, err)
				log.Debug("%s", errorMsg)
				return errors.New(errorMsg)
			}
		}
	}
	return nil
}
