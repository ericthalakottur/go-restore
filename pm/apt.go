package pm

import (
	"os/exec"
	"strings"
)

type AptPackageManager struct{}

func (apt AptPackageManager) ListPackages() ([]Package, error) {
	cmd := exec.Command("apt", "list", "--manual-installed")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []Package
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) > 1 {
			packages = append(packages, Package{
				Name:    parts[0],
				Version: parts[1],
			})
		}
	}
	return packages, nil
}
