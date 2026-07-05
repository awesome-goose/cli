package generator

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var generatorTitleCaser = cases.Title(language.English)

//go:embed all:templates/apps
var appTemplates embed.FS

//go:embed all:templates/modules
var moduleTemplates embed.FS

// frontendsDirName is the template subdirectory holding per-framework
// frontend variants; only the chosen framework's tree is generated.
const frontendsDirName = "frontends"

// AppGenerator generates a new application from templates
type AppGenerator struct {
	Name       string
	Template   string
	Framework  string
	OutputPath string
	ModulePath string
	data       map[string]string
}

// NewAppGenerator creates a new app generator
func NewAppGenerator(name, tmpl, outputPath string) *AppGenerator {
	modulePath := fmt.Sprintf("github.com/awesome-goose/%s", strings.ToLower(name))

	return &AppGenerator{
		Name:       name,
		Template:   tmpl,
		OutputPath: outputPath,
		ModulePath: modulePath,
		data: map[string]string{
			"AppName":      name,
			"AppNameLower": strings.ToLower(name),
			"AppNameSnake": toSnakeCase(name),
			"ModulePath":   modulePath,
			"Framework":    "",
		},
	}
}

// WithFramework selects the frontend framework variant to generate
func (g *AppGenerator) WithFramework(framework string) *AppGenerator {
	g.Framework = framework
	g.data["Framework"] = framework
	return g
}

// Generate creates the application structure
func (g *AppGenerator) Generate() error {
	templateDir := fmt.Sprintf("templates/apps/%s", g.Template)
	variantsRoot := templateDir + "/" + frontendsDirName

	// Pass 1: shared files, skipping the frontend variant tree
	err := fs.WalkDir(appTemplates, templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root template directory
		if path == templateDir {
			return nil
		}

		if path == variantsRoot {
			return fs.SkipDir
		}

		return g.emit(path, strings.TrimPrefix(path, templateDir+"/"), d)
	})
	if err != nil || g.Framework == "" {
		return err
	}

	// Pass 2: the chosen frontend variant, remapped under frontend/
	variantDir := variantsRoot + "/" + g.Framework

	return fs.WalkDir(appTemplates, variantDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == variantDir {
			return nil
		}

		return g.emit(path, "frontend/"+strings.TrimPrefix(path, variantDir+"/"), d)
	})
}

// emit writes a single template entry to the output tree. Files ending in
// .tmpl are rendered with the generator data; all other files are copied
// verbatim so frontend sources keep their own {{ }} syntax.
func (g *AppGenerator) emit(srcPath, relPath string, d fs.DirEntry) error {
	outputRelPath := strings.TrimSuffix(relPath, ".tmpl")
	outputPath := filepath.Join(g.OutputPath, outputRelPath)

	if d.IsDir() {
		return os.MkdirAll(outputPath, 0755)
	}

	content, err := appTemplates.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", srcPath, err)
	}

	out := content
	if strings.HasSuffix(relPath, ".tmpl") {
		processed, err := g.processTemplate(string(content))
		if err != nil {
			return fmt.Errorf("failed to process template %s: %w", srcPath, err)
		}
		out = []byte(processed)
	}

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(outputPath, out, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", outputPath, err)
	}

	fmt.Printf("  Created: %s\n", outputRelPath)
	return nil
}

func (g *AppGenerator) processTemplate(content string) (string, error) {
	tmpl, err := template.New("").Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g.data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ModuleGenerator generates a new module from templates
type ModuleGenerator struct {
	Name       string
	ModuleType string
	Template   string
	AppPath    string
	data       map[string]string
}

// NewModuleGenerator creates a new module generator
func NewModuleGenerator(name, moduleType, template, appPath string) *ModuleGenerator {
	return &ModuleGenerator{
		Name:       name,
		ModuleType: moduleType,
		Template:   template,
		AppPath:    appPath,
		data: map[string]string{
			"ModuleName":      toPascalCase(name),
			"ModuleNameLower": strings.ToLower(name),
			"ModuleNameSnake": toSnakeCase(name),
		},
	}
}

// Generate creates the module structure
func (g *ModuleGenerator) Generate() error {
	// Resolve the app's module path from go.mod so templates can import
	// the generated package
	modulePath, err := readModulePath(g.AppPath)
	if err != nil {
		return err
	}
	g.data["ModulePath"] = modulePath

	templateDir := fmt.Sprintf("templates/modules/%s/%s", g.ModuleType, g.Template)
	outputDir := filepath.Join(g.AppPath, "app", strings.ToLower(g.Name))

	err = fs.WalkDir(moduleTemplates, templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root template directory
		if path == templateDir {
			return nil
		}

		// Calculate relative path
		relPath := strings.TrimPrefix(path, templateDir+"/")

		// Remove .tmpl extension and replace generic names
		outputRelPath := strings.TrimSuffix(relPath, ".tmpl")
		outputRelPath = strings.ReplaceAll(outputRelPath, "module.go", fmt.Sprintf("%s.module.go", strings.ToLower(g.Name)))
		outputRelPath = strings.ReplaceAll(outputRelPath, "controller.go", fmt.Sprintf("%s.controller.go", strings.ToLower(g.Name)))
		outputRelPath = strings.ReplaceAll(outputRelPath, "service.go", fmt.Sprintf("%s.service.go", strings.ToLower(g.Name)))
		outputRelPath = strings.ReplaceAll(outputRelPath, "routes.go", fmt.Sprintf("%s.routes.go", strings.ToLower(g.Name)))
		outputRelPath = strings.ReplaceAll(outputRelPath, "dtos.go", fmt.Sprintf("%s.dtos.go", strings.ToLower(g.Name)))
		outputRelPath = strings.ReplaceAll(outputRelPath, "entity.go", fmt.Sprintf("%s.entity.go", strings.ToLower(g.Name)))

		outputPath := filepath.Join(outputDir, outputRelPath)

		if d.IsDir() {
			return os.MkdirAll(outputPath, 0755)
		}

		// Read template content
		content, err := moduleTemplates.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}

		// Process template
		processed, err := g.processTemplate(string(content))
		if err != nil {
			return fmt.Errorf("failed to process template %s: %w", path, err)
		}

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Write file
		if err := os.WriteFile(outputPath, []byte(processed), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", outputPath, err)
		}

		fmt.Printf("  Created: app/%s/%s\n", strings.ToLower(g.Name), outputRelPath)
		return nil
	})

	if err != nil {
		return err
	}

	// Update the app module to import the new module
	return g.updateAppModule()
}

func (g *ModuleGenerator) processTemplate(content string) (string, error) {
	tmpl, err := template.New("").Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g.data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *ModuleGenerator) updateAppModule() error {
	appModulePath := filepath.Join(g.AppPath, "app", "app.module.go")

	content, err := os.ReadFile(appModulePath)
	if err != nil {
		return fmt.Errorf("failed to read app.module.go: %w", err)
	}

	moduleNameLower := strings.ToLower(g.Name)
	moduleName := toPascalCase(g.Name)

	basePath, err := readModulePath(g.AppPath)
	if err != nil {
		return err
	}

	// Add import if not exists
	importLine := fmt.Sprintf(`"%s/app/%s"`, basePath, moduleNameLower)
	contentStr := string(content)

	if !strings.Contains(contentStr, importLine) {
		// Find import block and add the new import
		importRegex := regexp.MustCompile(`import\s*\(([^)]+)\)`)
		contentStr = importRegex.ReplaceAllStringFunc(contentStr, func(match string) string {
			// Add the new import before the closing parenthesis
			return strings.TrimSuffix(match, ")") + fmt.Sprintf("\t%s\n)", importLine)
		})
	}

	// Add module registration if not exists
	moduleRegistration := fmt.Sprintf("&%s.%sModule{}", moduleNameLower, moduleName)
	if !strings.Contains(contentStr, moduleRegistration) {
		// Find the Imports() method and add the module registration
		importsRegex := regexp.MustCompile(`(func\s*\([^)]+\)\s*Imports\s*\(\)\s*\[\]types\.Module\s*\{\s*return\s*\[\]types\.Module\s*\{)`)
		contentStr = importsRegex.ReplaceAllStringFunc(contentStr, func(match string) string {
			return match + fmt.Sprintf("\n\t\t%s,", moduleRegistration)
		})
	}

	if err := os.WriteFile(appModulePath, []byte(contentStr), 0644); err != nil {
		return fmt.Errorf("failed to write app.module.go: %w", err)
	}

	fmt.Printf("  Updated: app/app.module.go\n")
	return nil
}

// readModulePath extracts the module path from an app's go.mod
func readModulePath(appPath string) (string, error) {
	goModPath := filepath.Join(appPath, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	modulePathRegex := regexp.MustCompile(`module\s+(\S+)`)
	matches := modulePathRegex.FindSubmatch(goModContent)
	if len(matches) < 2 {
		return "", fmt.Errorf("failed to find module path in go.mod")
	}

	return string(matches[1]), nil
}

// Helper functions
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func toPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(string(word[0])))
			result.WriteString(strings.ToLower(word[1:]))
		}
	}

	if result.Len() == 0 {
		return generatorTitleCaser.String(strings.ToLower(s))
	}

	return result.String()
}

// DetectAppType tries to detect the app type from an existing application
func DetectAppType(appPath string) (string, error) {
	mainPath := filepath.Join(appPath, "main.go")
	content, err := os.ReadFile(mainPath)
	if err != nil {
		return "", fmt.Errorf("failed to read main.go: %w", err)
	}

	contentStr := string(content)

	if strings.Contains(contentStr, "platforms/spa") {
		return "spa", nil
	}
	if strings.Contains(contentStr, "platforms/api") {
		return "api", nil
	}
	if strings.Contains(contentStr, "platforms/cli") {
		return "cli", nil
	}
	if strings.Contains(contentStr, "platforms/web") {
		return "web", nil
	}

	return "", fmt.Errorf("could not detect app type from main.go")
}
