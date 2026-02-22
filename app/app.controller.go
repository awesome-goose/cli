package app

import (
	"fmt"

	"github.com/awesome-goose/goose/output"
	"github.com/awesome-goose/goose/types"
)

type AppController struct {
	appService *AppService `inject:""`
	log        types.Log   `inject:""`
}

// Version shows CLI version information
func (c *AppController) Version(body *VersionDto) types.Output {
	c.log.Info("AppController: Version called")
	return output.Box("Goose CLI", []string{
		fmt.Sprintf("Version: %s", c.appService.GetVersion()),
		"A tool for scaffolding Goose applications",
		"",
		"Commands:",
		"  goose app --name=<name> --template=<api|cli|web>  Create a new application",
		"  goose g module --name=<name> --type=<plain|resource>  Generate a module",
		"  goose version  Show version information",
	})
}

// App creates a new Goose application
func (c *AppController) App(body *AppDto) types.Output {
	c.log.Info("AppController: App called", map[string]any{
		"name":     body.Name,
		"template": body.Template,
		"path":     body.Path,
	})

	// Validate inputs
	if body.Name == "" {
		return output.ConsoleError("Error: --name flag is required")
	}
	if body.Template == "" {
		return output.ConsoleError("Error: --template flag is required (api, cli, or web)")
	}

	// Validate template
	validTemplates := []string{"api", "cli", "web"}
	isValid := false
	for _, t := range validTemplates {
		if body.Template == t {
			isValid = true
			break
		}
	}
	if !isValid {
		return output.ConsoleError(fmt.Sprintf("Error: invalid template '%s'. Valid templates: api, cli, web", body.Template))
	}

	// Create the application
	result, err := c.appService.CreateApp(body.Name, body.Template, body.Path)
	if err != nil {
		return output.ConsoleError(fmt.Sprintf("Error: %v", err))
	}

	return output.Box("Application Created", result)
}
