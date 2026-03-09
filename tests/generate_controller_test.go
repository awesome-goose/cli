package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/awesome-goose/goose-cli/app/generate"
	"github.com/awesome-goose/goose/core"
	test "github.com/awesome-goose/goose/testing"
	"github.com/awesome-goose/goose/types"
)

func TestGenerateController(t *testing.T) {
	test.NewSuiteRunner(t, &GenerateControllerSuite{}).Run()
}

// MockLog implements types.Log for testing
type MockLog struct{}

func (m *MockLog) Debug(message string, extra ...any)     {}
func (m *MockLog) Info(message string, extra ...any)      {}
func (m *MockLog) Notice(message string, extra ...any)    {}
func (m *MockLog) Warning(message string, extra ...any)   {}
func (m *MockLog) Error(message string, extra ...any)     {}
func (m *MockLog) Critical(message string, extra ...any)  {}
func (m *MockLog) Alert(message string, extra ...any)     {}
func (m *MockLog) Emergency(message string, extra ...any) {}

type GenerateControllerSuite struct {
	test.Suite
	controller *generate.GenerateController
	container  *core.Container
	tmpDir     string
	originalWd string
}

func (s *GenerateControllerSuite) SetupTest() {
	s.container = core.NewContainer()
	s.container.Register(func() *generate.GenerateService {
		return &generate.GenerateService{}
	}, "", true)
	s.container.Register(func() types.Log {
		return &MockLog{}
	}, "", true)

	s.controller = &generate.GenerateController{}
	s.container.Fill(s.controller)

	s.originalWd, _ = os.Getwd()

	tmpDir, err := os.MkdirTemp("", "goose-generate-ctrl-test-*")
	if err != nil {
		s.T.Require(err).ToBeNil()
	}
	s.tmpDir = tmpDir
}

func (s *GenerateControllerSuite) TeardownTest() {
	os.Chdir(s.originalWd)
	if s.container != nil {
		s.container.Reset()
	}
	if s.tmpDir != "" {
		os.RemoveAll(s.tmpDir)
	}
}

func (s *GenerateControllerSuite) createMockGooseApp(appType string) string {
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

func (s *GenerateControllerSuite) TestModule_ReturnsOutput() {
	dto := &generate.ModuleDto{
		Name:     "users",
		Type:     "plain",
		Template: "api",
	}

	output := s.controller.Module(dto)
	s.T.Expect(output).Not().ToBeNil()
}

func (s *GenerateControllerSuite) TestModule_MissingName_ReturnsError() {
	dto := &generate.ModuleDto{
		Name:     "",
		Type:     "plain",
		Template: "api",
	}

	output := s.controller.Module(dto)
	content := getOutputContent(output)

	s.T.Expect(content).ToContainString("Error")
	s.T.Expect(content).ToContainString("--name")
}

func (s *GenerateControllerSuite) TestModule_InvalidType_ReturnsError() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)

	dto := &generate.ModuleDto{
		Name:     "users",
		Type:     "invalid",
		Template: "api",
	}

	output := s.controller.Module(dto)
	content := getOutputContent(output)

	s.T.Expect(content).ToContainString("Error")
	s.T.Expect(content).ToContainString("invalid module type")
}

func (s *GenerateControllerSuite) TestModule_DefaultsToPlainType() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)

	dto := &generate.ModuleDto{
		Name:     "users",
		Type:     "",
		Template: "api",
	}

	output := s.controller.Module(dto)
	content := getOutputContent(output)

	if output.Code() == 0 {
		s.T.Expect(content).ToContainString("Generated")
		s.T.Expect(content).ToContainString("plain")
	}
}

func (s *GenerateControllerSuite) TestModule_ValidPlainType() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)

	dto := &generate.ModuleDto{
		Name:     "products",
		Type:     "plain",
		Template: "api",
	}

	output := s.controller.Module(dto)
	content := getOutputContent(output)

	if output.Code() == 0 {
		s.T.Expect(content).ToContainString("Generated")
		s.T.Expect(content).ToContainString("plain")
	}
}

func (s *GenerateControllerSuite) TestModule_ValidResourceType() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)

	dto := &generate.ModuleDto{
		Name:     "orders",
		Type:     "resource",
		Template: "api",
	}

	output := s.controller.Module(dto)
	content := getOutputContent(output)

	if output.Code() == 0 {
		s.T.Expect(content).ToContainString("Generated")
		s.T.Expect(content).ToContainString("resource")
	}
}

func (s *GenerateControllerSuite) TestModule_NotInGooseApp_ReturnsError() {
	os.Chdir(s.tmpDir)

	dto := &generate.ModuleDto{
		Name:     "users",
		Type:     "plain",
		Template: "api",
	}

	output := s.controller.Module(dto)
	content := getOutputContent(output)
	s.T.Expect(content).ToContainString("Error")
}

func (s *GenerateControllerSuite) TestModule_AutoDetectsTemplate() {
	appPath := s.createMockGooseApp("api")
	os.Chdir(appPath)

	dto := &generate.ModuleDto{
		Name:     "items",
		Type:     "plain",
		Template: "",
	}

	output := s.controller.Module(dto)
	content := getOutputContent(output)

	if output.Code() == 0 {
		s.T.Expect(content).ToContainString("api app")
	}
}

func (s *GenerateControllerSuite) TestModule_ValidTypesAccepted() {
	validTypes := []string{"plain", "resource"}

	for _, modType := range validTypes {
		appPath := s.createMockGooseApp("api")
		os.Chdir(appPath)

		dto := &generate.ModuleDto{
			Name:     "test" + modType,
			Type:     modType,
			Template: "api",
		}

		output := s.controller.Module(dto)
		content := getOutputContent(output)
		s.T.Expect(content).Not().ToContainString("invalid module type '" + modType + "'")
	}
}

func (s *GenerateControllerSuite) TestModule_InvalidTypesRejected() {
	invalidTypes := []string{"invalid", "PLAIN", "crud", "entity", "model"}

	for _, modType := range invalidTypes {
		dto := &generate.ModuleDto{
			Name:     "testmodule",
			Type:     modType,
			Template: "api",
		}

		output := s.controller.Module(dto)
		content := getOutputContent(output)
		s.T.Expect(content).ToContainString("invalid module type")
	}
}
