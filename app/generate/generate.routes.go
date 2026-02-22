package generate

import (
	"github.com/awesome-goose/goose/modules/router"
)

var ROUTES = router.ForRoutes(
	router.Cli("generate/module", []any{GenerateController{}, "Module"}),
	router.Cli("g/module", []any{GenerateController{}, "Module"}),
)
