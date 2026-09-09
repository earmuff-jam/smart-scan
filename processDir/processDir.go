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
// defines a function to traverse directories
// removes directory if found in env vars
func WalkDir(rootDir string) (map[int64][]types.File, error) {
	util.Info("Scanning directory: %s", rootDir)
	filesGroupedBySize := make(map[int64][]types.File)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			util.Error("error accessing: %s. Details: %+v", path, err)
			return nil
		}

		if d.IsDir() {
			if shouldTrashDir(d.Name()) {
				trash.MoveToTrash(path)
				util.Debug("removed directory %s as requested by env vars", d.Name())
			}
			util.Debug("skipping directory: %+v", d)
			return nil
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
		return nil, err
	}

	return filesGroupedBySize, nil
}

// BuildHashFile ...
// defines a function to verify file checksums
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
