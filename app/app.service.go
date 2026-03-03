package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/awesome-goose/goose-cli/generator"
)

const VERSION = "0.0.0"

type AppService struct{}

func (s *AppService) GetVersion() string {
	return VERSION
}

func (s *AppService) CreateApp(name, template, path string) ([]string, error) {
	// Determine output path
	outputPath := path
	if outputPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
		outputPath = cwd
	}

	// Create full output path
	fullPath := filepath.Join(outputPath, name)

	// Check if directory already exists
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		return nil, fmt.Errorf("directory '%s' already exists", fullPath)
	}

	// Generate the application
	gen := generator.NewAppGenerator(name, template, fullPath)
	if err := gen.Generate(); err != nil {
		return nil, fmt.Errorf("failed to generate application: %w", err)
	}

	// Build result messages
	result := []string{
		fmt.Sprintf("Created %s application '%s'", template, name),
		fmt.Sprintf("Location: %s", fullPath),
		"",
		"Next steps:",
		fmt.Sprintf("  cd %s", name),
		"  go mod tidy",
	}

	switch template {
	case "api", "web":
		result = append(result, "  # Update database configuration in app/app.module.go")
		result = append(result, "  go run main.go")
	case "multi":
		result = append(result,
			"",
			"Run modes:",
			"  go run main.go                # Start API (port 8080) + Web (port 3000) servers",
			"  go run main.go cli <command>  # Run CLI commands",
		)
	default:
		result = append(result, "  go run main.go")
	}

	return result, nil
}
