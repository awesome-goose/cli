package generate

type ModuleDto struct {
	Name     string `flag:"name"`
	Type     string `flag:"type"`
	Template string `flag:"template"`
}
