package app

import (
	"fmt"

	"github.com/awesome-goose/goose/io/output"
	"github.com/awesome-goose/goose/types"
)

type AppController struct {
	appService *AppService `inject:""`
	log        types.Log   `inject:""`
}

// Version shows CLI version information
func (c *AppController) Version(body *VersionDto) types.Output {
	return output.Box("Goose CLI", []string{
		fmt.Sprintf("Version: %s", c.appService.GetVersion()),
		"A tool for scaffolding Goose applications",
		"",
		"Commands:",
		"  goose app --name=<name> --template=<api|cli|web|multi|spa> [--framework=<react|vue|svelte|ng>]  Create a new application",
		"  goose g module --name=<name> --type=<plain|resource>  Generate a module",
		"  goose version  Show version information",
		"",
		"Templates:",
		"  api   - REST API application",
		"  cli   - Command-line interface application",
		"  web   - Web application with HTML templates",
		"  multi - Multi-platform app (API + Web + CLI in one)",
		"  spa   - Single-page app: JSON API + frontend served as one service (requires --framework)",
	})
}

// App creates a new Goose application
func (c *AppController) App(body *AppDto) types.Output {
	// Validate inputs
	if body.Name == "" {
		return output.ConsoleError("Error: --name flag is required")
	}
	if body.Template == "" {
		return output.ConsoleError("Error: --template flag is required (api, cli, web, multi, or spa)")
	}

	// Validate template
	validTemplates := []string{"api", "cli", "web", "multi", "spa"}
	isValid := false
	for _, t := range validTemplates {
		if body.Template == t {
			isValid = true
			break
		}
	}
	if !isValid {
		return output.ConsoleError(fmt.Sprintf("Error: invalid template '%s'. Valid templates: api, cli, web, multi, spa", body.Template))
	}

	// Validate framework
	if body.Template == "spa" {
		if body.Framework == "" {
			return output.ConsoleError("Error: --framework flag is required for the spa template (react, vue, svelte, or ng)")
		}
		validFrameworks := []string{"react", "vue", "svelte", "ng"}
		isValidFramework := false
		for _, f := range validFrameworks {
			if body.Framework == f {
				isValidFramework = true
				break
			}
		}
		if !isValidFramework {
			return output.ConsoleError(fmt.Sprintf("Error: invalid framework '%s'. Valid frameworks: react, vue, svelte, ng", body.Framework))
		}
	} else if body.Framework != "" {
		return output.ConsoleError("Error: --framework is only valid with --template=spa")
	}

	// Create the application
	result, err := c.appService.CreateApp(body.Name, body.Template, body.Path, body.Framework)
	if err != nil {
		return output.ConsoleError(fmt.Sprintf("Error: %v", err))
	}

	return output.Box("Application Created", result)
}
