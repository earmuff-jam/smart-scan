package process

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/earmuffjam/smart-scan/log"
	"github.com/earmuffjam/smart-scan/trash"
)

// RemoveUnwantedFolders ...
//
// removes folders matching the basename of configured file types.
func RemoveUnwantedFolders(rootDir string) error {
	removeMatchingFolders := os.Getenv("REMOVE_MATCHING_FOLDERS")

	return filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.TrimPrefix(filepath.Ext(d.Name()), ".")
		if !strings.Contains(removeMatchingFolders, ext) {
			return nil
		}

		baseName := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
		parentDir := filepath.Dir(path)

		entries, err := os.ReadDir(parentDir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if !entry.IsDir() || entry.Name() != baseName {
				continue
			}

			matchingPath := filepath.Join(parentDir, entry.Name())

			log.Info(
				"found %s. removing associated directory: %s",
				path,
				matchingPath,
			)

			if err := trash.MoveToTrash(matchingPath); err != nil {
				log.Debug(
					"unable to move directory to trash. details: %+v",
					err,
				)
				return err
			}

			break
		}

		return nil
	})

}
