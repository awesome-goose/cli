package app

type VersionDto struct{}

type AppDto struct {
	Name     string `flag:"name"`
	Template string `flag:"template"`
	Path     string `flag:"path"`
}
