package generate

import "github.com/awesome-goose/goose/types"

type GenerateModule struct{}

func (m *GenerateModule) Imports() []types.Module {
	return []types.Module{
		ROUTES,
	}
}

func (m *GenerateModule) Exports() []any {
	return []any{}
}

func (m *GenerateModule) Declarations() []any {
	return []any{
		&GenerateService{},
	}
}
