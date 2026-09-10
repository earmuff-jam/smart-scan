package trash

import (
	"fmt"
	"os/exec"
	"path/filepath"

	log "github.com/earmuffjam/smart-scan/log"
)

func moveToTrash(file string) error {
	filePath, err := filepath.Abs(file)
	if err != nil {
		log.Debug("unable to resolve file path: %+v", err)
		return fmt.Errorf("unable to resolve file path: %w", err)
	}

	script := fmt.Sprintf(
		`tell application "Finder" to delete POSIX file "%s"`,
		filePath,
	)

	cmd := exec.Command("osascript", "-e", script)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Debug("unable to move file to trash: %v: %s", err, output)
		return fmt.Errorf("unable to move file to trash: %w: %s", err, output)
	}

	return nil
}
