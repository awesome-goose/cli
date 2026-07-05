package generate

import (
	"fmt"
	"os"
	"strings"

	"github.com/awesome-goose/goose-cli/generator"
)

type GenerateService struct{}

func (s *GenerateService) GenerateModule(name, moduleType, template string) ([]string, error) {
	// Get current working directory as app path
	appPath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	// Verify we're in a goose app directory
	mainPath := appPath + "/main.go"
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("not in a Goose application directory (main.go not found)")
	}

	appModulePath := appPath + "/app/app.module.go"
	if _, err := os.Stat(appModulePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("not in a Goose application directory (app/app.module.go not found)")
	}

	// Auto-detect template if not specified
	tmpl := template
	if tmpl == "" {
		detected, err := generator.DetectAppType(appPath)
		if err != nil {
			return nil, fmt.Errorf("failed to auto-detect app type: %w\nPlease specify --template flag", err)
		}
		tmpl = detected
	}

	// Validate template
	validTemplates := []string{"api", "cli", "web", "spa"}
	isValidTemplate := false
	for _, t := range validTemplates {
		if tmpl == t {
			isValidTemplate = true
			break
		}
	}
	if !isValidTemplate {
		return nil, fmt.Errorf("invalid template '%s'. Valid templates: api, cli, web, spa", tmpl)
	}

	// SPA backend modules are JSON API modules — reuse the api module templates
	moduleTemplate := tmpl
	if tmpl == "spa" {
		moduleTemplate = "api"
	}

	// Check if module already exists
	moduleDir := appPath + "/app/" + strings.ToLower(name)
	if _, err := os.Stat(moduleDir); !os.IsNotExist(err) {
		return nil, fmt.Errorf("module '%s' already exists at %s", name, moduleDir)
	}

	// Generate the module
	gen := generator.NewModuleGenerator(name, moduleType, moduleTemplate, appPath)
	if err := gen.Generate(); err != nil {
		return nil, fmt.Errorf("failed to generate module: %w", err)
	}

	// Build result messages
	result := []string{
		fmt.Sprintf("Generated %s module for %s app", moduleType, tmpl),
		fmt.Sprintf("Location: app/%s/", strings.ToLower(name)),
		"",
		"The module has been automatically registered in app/app.module.go",
		fmt.Sprintf("You can now customize the generated files in app/%s/", strings.ToLower(name)),
	}

	return result, nil
}
