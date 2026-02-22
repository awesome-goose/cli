package app

import (
	"github.com/awesome-goose/goose/modules/router"
)

var (
	ROUTES = router.ForRoutes(
		// Default command shows version/help
		router.Cli("/", []any{AppController{}, "Version"}),
		router.Cli("version", []any{AppController{}, "Version"}),

		// App creation command
		router.Cli("app", []any{AppController{}, "App"}),
	)
)
