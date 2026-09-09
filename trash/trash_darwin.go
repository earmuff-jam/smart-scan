//go:build darwin

package trash

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func moveToTrash(file string) error {
	filePath, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf("unable to resolve file path: %w", err)
	}

	script := fmt.Sprintf(
		`tell application "Finder" to delete POSIX file %q`,
		filePath,
	)

	cmd := exec.Command("osascript", "-e", script)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to move file to trash: %w: %s", err, output)
	}

	return nil
}
