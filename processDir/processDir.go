package processdir

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/earmuffjam/smart-scan/trash"
	"github.com/earmuffjam/smart-scan/types"
	"github.com/earmuffjam/smart-scan/util"
)

// WalkDir ...
// walk through directories and sub-directories to find list of files. Moves directories to trash or skips them as defined in the env variable
func WalkDir(rootDir string) (map[int64][]types.File, error) {
	util.Info("Scanning directory: %s", rootDir)
	filesGroupedBySize := make(map[int64][]types.File, 0)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			util.Error("error accessing: %s. Details: %+v", path, err)
			return nil
		}

		if d.IsDir() {
			if shouldIgnoreDir(d.Name()) {
				util.Debug("skipping directory %s", d.Name())
				return fs.SkipDir
			}
			if shouldTrashDir(d.Name()) {
				trash.MoveToTrash(path)
				util.Debug("removed directory %s as requested by env vars", d.Name())
				return fs.SkipDir
			}
		}

		if !d.Type().IsRegular() {
			util.Info("skipping non-regular file: %s", path)
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			util.Error("failed to get info for file: %+v", err)
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

		util.Info("found file: %s (%d bytes)", file.Path, file.Size)
		return nil
	})

	if err != nil {
		util.Error("unable to walk directory: %+v", err)
		return filesGroupedBySize, err
	}

	return filesGroupedBySize, nil
}

// BuildHashFile ...
// verifies file checksums and creates hash for each file
func BuildHashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		util.Error("unable to build hash for file. error: %+v", err)
		return "", err
	}
	defer file.Close()

	hashVal := sha256.New()
	if _, err := io.Copy(hashVal, file); err != nil {
		util.Error("unable to copy hashValues. error: %+v", nil)
		return "", err
	}
	return hex.EncodeToString(hashVal.Sum(nil)), nil
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
