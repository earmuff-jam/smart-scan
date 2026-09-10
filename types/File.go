package types

import (
	"os"
	"path/filepath"
	"strings"
)

// File ...
// defines the struct for file
type File struct {
	Path string
	Size int64
}

// ValidatePrefix ...
// checks whether the file extension is listed in the env var
func (f File) ValidatePrefix() bool {
	fileTypes := strings.Split(os.Getenv("REMOVE_FILE_TYPES"), ",")
	ext := strings.TrimPrefix(filepath.Ext(f.Path), ".")

	for _, fileType := range fileTypes {
		if strings.TrimSpace(fileType) == ext {
			return true
		}
	}

	return false
}
