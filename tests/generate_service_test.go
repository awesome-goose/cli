package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/awesome-goose/goose-cli/app/generate"
	test "github.com/awesome-goose/goose/testing"
)

func TestGenerateService(t *testing.T) {
	test.NewSuiteRunner(t, &GenerateServiceSuite{}).Run()
}

type GenerateServiceSuite struct {
	test.Suite
	service    *generate.GenerateService
	tmpDir     string
	originalWd string
}

func (s *GenerateServiceSuite) SetupTest() {
	s.service = &generate.GenerateService{}
	s.originalWd, _ = os.Getwd()

	tmpDir, err := os.MkdirTemp("", "goose-generate-test-*")
	if err != nil {
		s.T.Require(err).ToBeNil()
	}
	s.tmpDir = tmpDir
}

func (s *GenerateServiceSuite) TeardownTest() {
	os.Chdir(s.originalWd)
	if s.tmpDir != "" {
		os.RemoveAll(s.tmpDir)
	}
}

func (s *GenerateServiceSuite) createMockGooseApp(appType string) string {
	appPath := filepath.Join(s.tmpDir, "testapp-"+test.NewFixture().String(6))
	os.MkdirAll(filepath.Join(appPath, "app"), 0755)

	mainContent := "package main\nimport \"github.com/awesome-goose/goose/platforms/" + appType + "\"\nfunc main() { " + appType + ".Run() }"
	os.WriteFile(filepath.Join(appPath, "main.go"), []byte(mainContent), 0644)

	moduleContent := `package app

import "github.com/awesome-goose/goose/types"

type AppModule struct{}

func (m *AppModule) Imports() []types.Module {
	return []types.Module{}
}

func (m *AppModule) Declarations() []any {
	return []any{}
}

func (m *AppModule) Exports() []any {
	return []any{}
}
`
	os.WriteFile(filepath.Join(appPath, "app", "app.module.go"), []byte(moduleContent), 0644)

	goModContent := "module github.com/awesome-goose/testapp\n\ngo 1.21\n"
	os.WriteFile(filepath.Join(appPath, "go.mod"), []byte(goModContent), 0644)

	return appPath
}

func (s *GenerateServiceSuite) TestGenerateModule_NotInGooseApp_ReturnsError() {
	os.Chdir(s.tmpDir)

	result, err := s.service.GenerateModule("users", "plain", "api")

	s.T.Expect(err).Not().ToBeNil()
	s.T.Expect(result).ToBeNil()
	s.T.Expect(err.Error()).ToContainString("main.go not found")
}

func (s *GenerateServiceSuite) TestGenerateModule_MissingAppModule_ReturnsError() {
	appPath := filepath.Join(s.tmpDir, "nomodule")
	os.MkdirAll(appPath, 0755)
	os.WriteFile(filepath.Join(appPath, "main.go"), []byte("package main\nfunc main(){}"), 0644)
	os.Chdir(appPath)

	result, err := s.service.GenerateModule("users", "plain", "api")

	s.T.Expect(err).Not().ToBeNil()
	s.T.Expect(result).ToBeNil()
	s.T.Expect(err.Error()).ToContainString("app.module.go not found")
}

func (s *GenerateServiceSuite) TestGenerateModule_InvalidTemplate_ReturnsError() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)

	result, err := s.service.GenerateModule("users", "plain", "invalid")

	s.T.Expect(err).Not().ToBeNil()
	s.T.Expect(result).ToBeNil()
	s.T.Expect(err.Error()).ToContainString("invalid template")
}

func (s *GenerateServiceSuite) TestGenerateModule_ModuleAlreadyExists_ReturnsError() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)
	os.MkdirAll(filepath.Join(appPath, "app", "users"), 0755)

	result, err := s.service.GenerateModule("users", "plain", "api")

	s.T.Expect(err).Not().ToBeNil()
	s.T.Expect(result).ToBeNil()
	s.T.Expect(err.Error()).ToContainString("already exists")
}

func (s *GenerateServiceSuite) TestGenerateModule_ValidPlainModule_Success() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)

	result, err := s.service.GenerateModule("users", "plain", "api")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()
	s.T.Expect(len(result)).ToBeGreaterThan(0)

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("Generated plain module")
	s.T.Expect(dirExists(filepath.Join(appPath, "app", "users"))).ToBeTrue()
}

func (s *GenerateServiceSuite) TestGenerateModule_ValidResourceModule_Success() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)

	result, err := s.service.GenerateModule("products", "resource", "api")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("Generated resource module")
	s.T.Expect(dirExists(filepath.Join(appPath, "app", "products"))).ToBeTrue()
}

func (s *GenerateServiceSuite) TestGenerateModule_AutoDetectsApiTemplate() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)

	result, err := s.service.GenerateModule("orders", "plain", "")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("api app")
}

func (s *GenerateServiceSuite) TestGenerateModule_AutoDetectsCliTemplate() {
	appPath := s.createMockGooseApp("cli")
	os.Chdir(appPath)

	result, err := s.service.GenerateModule("commands", "plain", "")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("cli app")
}

func (s *GenerateServiceSuite) TestGenerateModule_AutoDetectsWebTemplate() {
	appPath := s.createMockGooseApp("web")
	os.Chdir(appPath)

	result, err := s.service.GenerateModule("pages", "plain", "")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("web app")
}

func (s *GenerateServiceSuite) TestGenerateModule_SpaTemplateUsesApiModules() {
	appPath := s.createMockGooseApp("spa")
	os.Chdir(appPath)

	result, err := s.service.GenerateModule("users", "plain", "spa")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("spa app")

	moduleDir := filepath.Join(appPath, "app", "users")
	s.T.Expect(fileExists(filepath.Join(moduleDir, "users.module.go"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(moduleDir, "users.controller.go"))).ToBeTrue()
}

func (s *GenerateServiceSuite) TestGenerateModule_AutoDetectsSpaTemplate() {
	appPath := s.createMockGooseApp("spa")
	os.Chdir(appPath)

	result, err := s.service.GenerateModule("orders", "plain", "")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("spa app")
}

func (s *GenerateServiceSuite) TestGenerateModule_GeneratesRequiredFiles() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)

	s.service.GenerateModule("items", "plain", "api")

	moduleDir := filepath.Join(appPath, "app", "items")
	s.T.Expect(fileExists(filepath.Join(moduleDir, "items.module.go"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(moduleDir, "items.controller.go"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(moduleDir, "items.service.go"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(moduleDir, "items.routes.go"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(moduleDir, "items.dtos.go"))).ToBeTrue()
}

func (s *GenerateServiceSuite) TestGenerateModule_ResourceGeneratesEntityFile() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)

	s.service.GenerateModule("customers", "resource", "api")

	moduleDir := filepath.Join(appPath, "app", "customers")
	s.T.Expect(fileExists(filepath.Join(moduleDir, "customers.entity.go"))).ToBeTrue()
}

func (s *GenerateServiceSuite) TestGenerateModule_UpdatesAppModule() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)

	s.service.GenerateModule("accounts", "plain", "api")

	content, _ := os.ReadFile(filepath.Join(appPath, "app", "app.module.go"))
	contentStr := string(content)
	s.T.Expect(contentStr).ToContainString("accounts")
}
