package pm

import (
	"fmt"
	"os/exec"
)

type Package struct {
	Name    string
	Version string
}

type PackageManger interface {
	ListPackages() ([]Package, error)
	// Install(packageName string) error
	// Search(packageName string) []string
	// Update() error
}

func checkIfExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func GetPackageManager() (PackageManger, error) {
	var packageManager PackageManger
	var err error
	if checkIfExists("apt") {
		packageManager = AptPackageManager{}
	}
	if packageManager == nil {
		err = fmt.Errorf("No package manager found")
	}
	return packageManager, err
}
