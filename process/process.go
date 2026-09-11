package process

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	log "github.com/earmuffjam/smart-scan/log"
	"github.com/earmuffjam/smart-scan/trash"
	"github.com/earmuffjam/smart-scan/types"
)

// BuildHashFile ...
//
// verifies file checksums and creates hash for each file
func BuildHashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Debug("unable to build hash for file. error: %+v", err)
		return "", err
	}
	defer file.Close()

	hashVal := sha256.New()
	if _, err := io.Copy(hashVal, file); err != nil {
		log.Debug("unable to copy hashValues. error: %+v", nil)
		return "", err
	}
	return hex.EncodeToString(hashVal.Sum(nil)), nil
}

// RemoveOwnerFolders ...
//
// defines a function that is used to remove owner folders
func RemoveOwnerFolders(files []string) (int, error) {
	removedCount := 0
	removeMatchingOwnerFolders := os.Getenv("REMOVE_OWNER_FOLDERS")

	for _, path := range files {
		name := filepath.Base(path)
		ext := strings.TrimPrefix(filepath.Ext(name), ".")

		folders := strings.Split(removeMatchingOwnerFolders, ", ")

		if len(ext) == 0 || !slices.Contains(folders, ext) {
			continue
		}

		baseName := strings.TrimSuffix(name, filepath.Ext(name))
		parentDir := filepath.Dir(path)

		entries, err := os.ReadDir(parentDir)
		if err != nil {
			log.Debug("unable to read provided directory. details: %+v", err)
			return removedCount, err
		}

		for _, entry := range entries {
			if !entry.IsDir() || entry.Name() != baseName {
				continue
			}

			matchingPath := filepath.Join(parentDir, entry.Name())

			log.Info(
				"found folder %s with %s extension. removing associated owner directory: %s",
				path,
				matchingPath,
			)

			if err := trash.MoveToTrash(matchingPath); err != nil {
				log.Debug(
					"unable to move directory to trash. details: %+v",
					err,
				)
				return removedCount, err
			}
			removedCount++
			break
		}
	}
	return removedCount, nil
}

// RemoveUnwantedFolders ...
//
// defines a function that is used to remove unwanted folders
func RemoveUnwantedFolders(files []string) (int, error) {
	removedCount := 0
	removeMatchingFolders := os.Getenv("FOLDERS_TO_REMOVE")

	for _, path := range files {
		name := filepath.Base(path)
		folders := strings.Split(removeMatchingFolders, ", ")

		if !slices.Contains(folders, name) {
			continue
		}

		if err := trash.MoveToTrash(path); err != nil {
			log.Debug(
				"unable to move directory to trash. details: %+v",
				err,
			)
			return removedCount, err
		}
		removedCount++
	}
	return removedCount, nil
}

// RemoveUnwantedFiles ...
//
// removes unwanted files slated for removal in env variables
func RemoveUnwantedFiles(files []types.File) (int, error) {
	removedCount := 0
	removeFileTypes := os.Getenv("REMOVE_FILE_TYPES")

	for _, file := range files {
		ext := strings.TrimPrefix(filepath.Ext(file.Path), ".")
		files := strings.Split(removeFileTypes, ", ")

		if len(ext) == 0 || !slices.Contains(files, ext) {
			continue
		}

		log.Debug("moving file %s to trash.", file.Path)

		if err := trash.MoveToTrash(file.Path); err != nil {
			log.Debug(
				"unable to move file %s to trash. details: %+v",
				file.Path,
				err,
			)
			return removedCount, errors.New("failed to remove file")
		}

		removedCount++
	}

	return removedCount, nil
}

// RemoveDuplicate ...
//
// used to detect and remove duplicate files
func RemoveDuplicate(filesMap map[int64][]types.File) (int, error) {
	groupDuplicateFiles, err := detectDuplicates(filesMap)
	if err != nil {
		log.Debug("unable to detect duplicates. details: %+v", err)
		return 0, err
	}

	if len(groupDuplicateFiles) == 0 {
		log.Debug("no duplicate files found.")
		return 0, nil
	}

	removedDuplicateFilesCount, err := moveDuplicateFilesToTrash(groupDuplicateFiles)
	if err != nil {
		log.Debug("failed to remove duplicate files")
		return 0, err
	}

	return removedDuplicateFilesCount, nil
}

func detectDuplicates(groupedFiles map[int64][]types.File) (map[string][]types.File, error) {
	groupedByFileHash := make(map[string][]types.File)
	for _, files := range groupedFiles {
		if len(files) < 2 {
			continue
		}

		for _, file := range files {
			hash, err := BuildHashFile(file.Path)
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

func moveDuplicateFilesToTrash(groupedByHash map[string][]types.File) (int, error) {
	removedFilesCount := 0
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
				return removedFilesCount, errors.New(errorMsg)
			}
			removedFilesCount++
		}
	}
	return removedFilesCount, nil
}
