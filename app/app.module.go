package app

import (
	"github.com/awesome-goose/goose-cli/app/generate"
	"github.com/awesome-goose/goose/types"
)

type AppModule struct{}

func (m *AppModule) Imports() []types.Module {
	return []types.Module{
		&generate.GenerateModule{},

		ROUTES,
	}
}

func (m *AppModule) Exports() []any {
	return []any{}
}

func (m *AppModule) Declarations() []any {
	return []any{
		&AppService{},
	}
}
