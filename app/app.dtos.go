package app

type VersionDto struct{}

type AppDto struct {
	Name      string `flag:"name"`
	Template  string `flag:"template"`
	Framework string `flag:"framework"`
	Path      string `flag:"path"`
}
