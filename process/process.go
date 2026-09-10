package process

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/earmuffjam/smart-scan/log"
	"github.com/earmuffjam/smart-scan/trash"
	"github.com/earmuffjam/smart-scan/types"
)

// BuildHashFile ...
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

// WalkDir ...
// walk through directories and sub-directories to find list of files.
// Moves directories to trash or skips them as defined in the env variable
func WalkDir(rootDir string) (map[int64][]types.File, error) {
	log.Debug("Scanning directory: %s", rootDir)
	filesGroupedBySize := make(map[int64][]types.File, 0)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Debug("error accessing: %s. Details: %+v", path, err)
			return nil
		}

		if d.IsDir() {
			if shouldIgnoreDir(d.Name()) {
				log.Debug("skipping directory %s", d.Name())
				return fs.SkipDir
			}

			if shouldTrashDir(d.Name()) {
				log.Info("Found %s directory set to remove. Processing ...", d.Name())
				trash.MoveToTrash(path)
				return fs.SkipDir
			}
		}

		if !d.Type().IsRegular() {
			log.Debug("skipping non-regular file: %s", path)
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			log.Debug("failed to get info for file: %+v", err)
			return nil
		}

		file := types.File{
			Path: path,
			Size: fileInfo.Size(),
		}

		filesGroupedBySize[file.Size] = append(
			filesGroupedBySize[file.Size],
			file,
		)

		log.Debug("found file: %s (%d bytes)", file.Path, file.Size)
		return nil
	})

	if err != nil {
		log.Debug("unable to walk directory: %+v", err)
		return filesGroupedBySize, err
	}

	return filesGroupedBySize, nil
}

func shouldTrashDir(directoryName string) bool {
	removeDirs := os.Getenv("REMOVE_DIRS")

	for dir := range strings.SplitSeq(removeDirs, ",") {
		if strings.TrimSpace(dir) == directoryName {
			return true
		}
	}

	return false
}

func shouldIgnoreDir(directoryName string) bool {
	ignoreDirs := os.Getenv("IGNORE_DIRS")

	for dir := range strings.SplitSeq(ignoreDirs, ",") {
		if strings.TrimSpace(dir) == directoryName {
			return true
		}
	}

	return false
}
