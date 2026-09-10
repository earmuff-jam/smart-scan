package process

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/earmuffjam/smart-scan/log"
	"github.com/earmuffjam/smart-scan/types"
)

// WalkDirectory ...
//
// defines a function that walks through a provided directory to return files within.
func WalkDirectory(rootDir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Debug("unable to walk provided directory %s. details: %+v", rootDir, err)
			return err
		}

		// ignore the main root directory
		if path == rootDir {
			return nil
		}

		if shouldIgnoreDir(d.Name()) {
			log.Debug("Ignore directory %s as requested by env variables.", d.Name())
			return fs.SkipDir
		}

		files = append(files, path)
		return nil
	})

	log.Debug("found %d item(s) in %s directory.", len(files), rootDir)
	return files, err
}

// WalkFiles ...
//
// defines a function that walks through a provided directory to return a list of files
func WalkFiles(rootDir string) ([]types.File, error) {
	files := make([]types.File, 0)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Debug("error accessing: %s. Details: %+v", path, err)
			return nil
		}

		if d.IsDir() {
			if shouldIgnoreDir(d.Name()) {
				log.Debug("Ignore directory %s as requested by env variables.", d.Name())
				return fs.SkipDir
			}
			return nil
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

		files = append(files, file)
		log.Debug("found file: %s (%d bytes)", file.Path, file.Size)

		return nil
	})

	return files, err
}

// GroupFilesBySize ...
//
// defines a function that groups files by their size
func GroupFilesBySize(files []types.File) map[int64][]types.File {
	filesGroupedBySize := make(map[int64][]types.File)

	for _, file := range files {
		filesGroupedBySize[file.Size] = append(
			filesGroupedBySize[file.Size],
			file,
		)
	}

	return filesGroupedBySize
}

func shouldIgnoreDir(directoryName string) bool {
	ignoreDirs := os.Getenv("FOLDERS_TO_IGNORE")

	for dir := range strings.SplitSeq(ignoreDirs, ",") {
		if strings.TrimSpace(dir) == directoryName {
			return true
		}
	}

	return false
}
