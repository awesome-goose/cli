package app

type VersionDto struct{}

type AppDto struct {
	Name     string `query:"name"`
	Template string `query:"template"`
	Path     string `query:"path"`
}
