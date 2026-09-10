// trash_linux.go
package trash

import (
	"fmt"
	"os/exec"

	log "github.com/earmuffjam/smart-scan/log"
)

func moveToTrash(file string) error {
	cmd := exec.Command("gio", "trash", file)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Debug("unable to move file to trash: %v: %s", err, output)
		return fmt.Errorf("unable to move file to trash: %w: %s", err, output)
	}

	return nil
}
