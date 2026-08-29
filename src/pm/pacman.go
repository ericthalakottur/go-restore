package pm

import "os/exec"

type PacmanPackageManager struct{}

func (p *PacmanPackageManager) Install(packageName string) error {
	cmd := exec.Command("pacman", "-S", packageName)
	return cmd.Run()
}

func (p *PacmanPackageManager) Search(packageName string) error {
	cmd := exec.Command("pacman", "-Ss", packageName)
	return cmd.Run()
}

func (p *PacmanPackageManager) Update() error {
	cmd := exec.Command("pacman", "-Syu")
	return cmd.Run()
}
