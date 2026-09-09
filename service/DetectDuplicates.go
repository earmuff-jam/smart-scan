package service

import (
	"errors"
	"fmt"

	processdir "github.com/earmuffjam/smart-scan/processDir"
	"github.com/earmuffjam/smart-scan/trash"
	"github.com/earmuffjam/smart-scan/types"
	"github.com/earmuffjam/smart-scan/util"
)

// DetectDuplicates ...
// defines a function that detects all duplicates
func DetectDuplicates(groupedFiles map[int64][]types.File) (map[string][]types.File, error) {
	groupedByFileHash := make(map[string][]types.File)
	for _, files := range groupedFiles {
		if len(files) < 2 {
			continue
		}

		for _, file := range files {
			hash, err := processdir.BuildHashFile(file.Path)
			if err != nil {
				util.Error(
					"failed to hash %s. details: %+v", file.Path, err)
				continue

			}
			// hash file with matching size
			groupedByFileHash[hash] = append(groupedByFileHash[hash], file)
		}
	}
	return groupedByFileHash, nil
}

// MoveToTrash ...
// defines a function to move duplicates to trash
func MoveToTrash(groupedByHash map[string][]types.File) error {

	for hash, files := range groupedByHash {
		if len(files) < 2 {
			continue
		}
		util.Info("duplicate group: %s", hash)

		// keep the first file
		for _, file := range files[1:] {
			util.Debug("moving file %s to trash.", file.Path)
			if err := moveToTrash(file.Path); err != nil {
				errorMsg := fmt.Sprintf("unable to move file %s to trash. details: %+v", file.Path, err)
				util.Error(errorMsg)
				return errors.New(errorMsg)
			}
		}
	}
	return nil
}

func moveToTrash(file string) error {

	err := trash.MoveToTrash(file)
	if err != nil {
		util.Error("unable to move file to trash. details: %+v", err)
		return err
	}
	return nil
}
