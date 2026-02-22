package generate

type ModuleDto struct {
	Name     string `query:"name"`
	Type     string `query:"type"`
	Template string `query:"template"`
}
