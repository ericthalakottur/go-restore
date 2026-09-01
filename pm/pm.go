package pm

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

func GetPackageManager() (PackageManger, error) {
	packageManager := AptPackageManager{}
	return packageManager, nil
}
