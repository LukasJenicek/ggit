package helpers

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetProjectRootFolder() (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("get root folder: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
