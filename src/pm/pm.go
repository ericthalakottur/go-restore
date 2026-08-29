package pm

type PackageManger interface {
	Install(packageName string) error
	Search(packageName string) []string
	Update() error
}
