package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/awesome-goose/goose-cli/generator"
	test "github.com/awesome-goose/goose/testing"
)

func TestGenerator(t *testing.T) {
	test.NewSuiteRunner(t, &GeneratorSuite{}).Run()
}

type GeneratorSuite struct {
	test.Suite
	tmpDir string
}

func (s *GeneratorSuite) SetupTest() {
	tmpDir, err := os.MkdirTemp("", "goose-generator-test-*")
	if err != nil {
		s.T.Require(err).ToBeNil()
	}
	s.tmpDir = tmpDir
}

func (s *GeneratorSuite) TeardownTest() {
	if s.tmpDir != "" {
		os.RemoveAll(s.tmpDir)
	}
}

func (s *GeneratorSuite) createMockApp(appType string) string {
	appPath := filepath.Join(s.tmpDir, "mockapp-"+test.NewFixture().String(6))
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

// ==================== NewAppGenerator Tests ====================

func (s *GeneratorSuite) TestNewAppGenerator_CreatesGenerator() {
	gen := generator.NewAppGenerator("TestApp", "api", s.tmpDir)

	s.T.Expect(gen).Not().ToBeNil()
	s.T.Expect(gen.Name).ToEqual("TestApp")
	s.T.Expect(gen.Template).ToEqual("api")
	s.T.Expect(gen.OutputPath).ToEqual(s.tmpDir)
}

func (s *GeneratorSuite) TestNewAppGenerator_SetsModulePath() {
	gen := generator.NewAppGenerator("MyApp", "api", s.tmpDir)

	s.T.Expect(gen.ModulePath).ToContainString("myapp")
	s.T.Expect(gen.ModulePath).ToContainString("github.com/awesome-goose/")
}

// ==================== AppGenerator.Generate Tests ====================

func (s *GeneratorSuite) TestAppGenerator_Generate_ApiTemplate() {
	outputPath := filepath.Join(s.tmpDir, "testapi")
	gen := generator.NewAppGenerator("TestApi", "api", outputPath)

	err := gen.Generate()

	s.T.Expect(err).ToBeNil()
	s.T.Expect(dirExists(outputPath)).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "main.go"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "go.mod"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, ".env"))).ToBeTrue()
	s.T.Expect(dirExists(filepath.Join(outputPath, "app"))).ToBeTrue()
}

func (s *GeneratorSuite) TestAppGenerator_Generate_CliTemplate() {
	outputPath := filepath.Join(s.tmpDir, "testcli")
	gen := generator.NewAppGenerator("TestCli", "cli", outputPath)

	err := gen.Generate()

	s.T.Expect(err).ToBeNil()
	s.T.Expect(dirExists(outputPath)).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "main.go"))).ToBeTrue()
}

func (s *GeneratorSuite) TestAppGenerator_Generate_WebTemplate() {
	outputPath := filepath.Join(s.tmpDir, "testweb")
	gen := generator.NewAppGenerator("TestWeb", "web", outputPath)

	err := gen.Generate()

	s.T.Expect(err).ToBeNil()
	s.T.Expect(dirExists(outputPath)).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "main.go"))).ToBeTrue()
}

func (s *GeneratorSuite) TestAppGenerator_Generate_MultiTemplate() {
	outputPath := filepath.Join(s.tmpDir, "testmulti")
	gen := generator.NewAppGenerator("TestMulti", "multi", outputPath)

	err := gen.Generate()

	s.T.Expect(err).ToBeNil()
	s.T.Expect(dirExists(outputPath)).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "main.go"))).ToBeTrue()
}

func (s *GeneratorSuite) TestNewAppGenerator_WithFramework() {
	gen := generator.NewAppGenerator("TestSpa", "spa", s.tmpDir).WithFramework("react")

	s.T.Expect(gen.Framework).ToEqual("react")
}

func (s *GeneratorSuite) generateSpaApp(framework string) string {
	outputPath := filepath.Join(s.tmpDir, "testspa-"+framework)
	gen := generator.NewAppGenerator("TestSpa", "spa", outputPath).WithFramework(framework)

	err := gen.Generate()
	s.T.Expect(err).ToBeNil()

	return outputPath
}

func (s *GeneratorSuite) TestAppGenerator_Generate_SpaTemplateReact() {
	outputPath := s.generateSpaApp("react")

	s.T.Expect(fileExists(filepath.Join(outputPath, "main.go"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "go.mod"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "Makefile"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "package.json"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "vite.config.js"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "src", "App.jsx"))).ToBeTrue()
}

func (s *GeneratorSuite) TestAppGenerator_Generate_SpaTemplateVue() {
	outputPath := s.generateSpaApp("vue")

	s.T.Expect(fileExists(filepath.Join(outputPath, "Makefile"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "package.json"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "src", "App.vue"))).ToBeTrue()
}

func (s *GeneratorSuite) TestAppGenerator_Generate_SpaTemplateSvelte() {
	outputPath := s.generateSpaApp("svelte")

	s.T.Expect(fileExists(filepath.Join(outputPath, "Makefile"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "package.json"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "src", "App.svelte"))).ToBeTrue()
}

func (s *GeneratorSuite) TestAppGenerator_Generate_SpaTemplateNg() {
	outputPath := s.generateSpaApp("ng")

	s.T.Expect(fileExists(filepath.Join(outputPath, "Makefile"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "package.json"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "angular.json"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "proxy.conf.json"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "src", "app", "app.component.ts"))).ToBeTrue()
}

func (s *GeneratorSuite) TestAppGenerator_Generate_Spa_VerbatimFrontendKeepsBraces() {
	outputPath := s.generateSpaApp("vue")

	content, err := os.ReadFile(filepath.Join(outputPath, "frontend", "src", "App.vue"))
	s.T.Expect(err).ToBeNil()

	// Verbatim copy must preserve Vue's own {{ }} interpolation
	s.T.Expect(string(content)).ToContainString("{{ info.name }}")
}

func (s *GeneratorSuite) TestAppGenerator_Generate_Spa_TemplatedFilesRendered() {
	outputPath := s.generateSpaApp("react")

	makefile, err := os.ReadFile(filepath.Join(outputPath, "Makefile"))
	s.T.Expect(err).ToBeNil()
	s.T.Expect(string(makefile)).ToContainString("FRAMEWORK    := react")
	s.T.Expect(string(makefile)).Not().ToContainString("{{.")

	mainGo, err := os.ReadFile(filepath.Join(outputPath, "main.go"))
	s.T.Expect(err).ToBeNil()
	s.T.Expect(string(mainGo)).ToContainString("platforms/spa")
	s.T.Expect(string(mainGo)).Not().ToContainString("{{.")

	pkgJSON, err := os.ReadFile(filepath.Join(outputPath, "frontend", "package.json"))
	s.T.Expect(err).ToBeNil()
	s.T.Expect(string(pkgJSON)).ToContainString("testspa-frontend")
}

func (s *GeneratorSuite) TestAppGenerator_Generate_Spa_SkipsOtherFrameworks() {
	outputPath := s.generateSpaApp("react")

	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "angular.json"))).ToBeFalse()
	s.T.Expect(fileExists(filepath.Join(outputPath, "frontend", "src", "App.vue"))).ToBeFalse()
	s.T.Expect(dirExists(filepath.Join(outputPath, "frontends"))).ToBeFalse()
}

func (s *GeneratorSuite) TestAppGenerator_Generate_CreatesAppModule() {
	outputPath := filepath.Join(s.tmpDir, "testapp1")
	gen := generator.NewAppGenerator("TestApp", "api", outputPath)

	gen.Generate()

	s.T.Expect(fileExists(filepath.Join(outputPath, "app", "app.module.go"))).ToBeTrue()
}

func (s *GeneratorSuite) TestAppGenerator_Generate_CreatesAppController() {
	outputPath := filepath.Join(s.tmpDir, "testapp2")
	gen := generator.NewAppGenerator("TestApp", "api", outputPath)

	gen.Generate()

	s.T.Expect(fileExists(filepath.Join(outputPath, "app", "app.controller.go"))).ToBeTrue()
}

func (s *GeneratorSuite) TestAppGenerator_Generate_TemplateVariablesReplaced() {
	outputPath := filepath.Join(s.tmpDir, "testapp3")
	gen := generator.NewAppGenerator("MyTestApp", "api", outputPath)

	gen.Generate()

	content, _ := os.ReadFile(filepath.Join(outputPath, "main.go"))
	contentStr := string(content)

	s.T.Expect(contentStr).ToContainString("mytestapp")
	s.T.Expect(contentStr).Not().ToContainString("{{")
	s.T.Expect(contentStr).Not().ToContainString("}}")
}

func (s *GeneratorSuite) TestAppGenerator_Generate_GoModContainsModulePath() {
	outputPath := filepath.Join(s.tmpDir, "testapp4")
	gen := generator.NewAppGenerator("MyApp", "api", outputPath)

	gen.Generate()

	content, _ := os.ReadFile(filepath.Join(outputPath, "go.mod"))
	contentStr := string(content)

	s.T.Expect(contentStr).ToContainString("module")
	s.T.Expect(contentStr).ToContainString("myapp")
}

// ==================== NewModuleGenerator Tests ====================

func (s *GeneratorSuite) TestNewModuleGenerator_CreatesGenerator() {
	gen := generator.NewModuleGenerator("users", "plain", "api", s.tmpDir)

	s.T.Expect(gen).Not().ToBeNil()
	s.T.Expect(gen.Name).ToEqual("users")
	s.T.Expect(gen.ModuleType).ToEqual("plain")
	s.T.Expect(gen.Template).ToEqual("api")
	s.T.Expect(gen.AppPath).ToEqual(s.tmpDir)
}

// ==================== ModuleGenerator.Generate Tests ====================

func (s *GeneratorSuite) TestModuleGenerator_Generate_PlainModule() {
	appPath := s.createMockApp("api")
	gen := generator.NewModuleGenerator("users", "plain", "api", appPath)

	err := gen.Generate()

	s.T.Expect(err).ToBeNil()

	moduleDir := filepath.Join(appPath, "app", "users")
	s.T.Expect(dirExists(moduleDir)).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(moduleDir, "users.module.go"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(moduleDir, "users.controller.go"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(moduleDir, "users.service.go"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(moduleDir, "users.routes.go"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(moduleDir, "users.dtos.go"))).ToBeTrue()
}

func (s *GeneratorSuite) TestModuleGenerator_Generate_ResourceModule() {
	appPath := s.createMockApp("api")
	gen := generator.NewModuleGenerator("products", "resource", "api", appPath)

	err := gen.Generate()

	s.T.Expect(err).ToBeNil()

	moduleDir := filepath.Join(appPath, "app", "products")
	s.T.Expect(dirExists(moduleDir)).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(moduleDir, "products.entity.go"))).ToBeTrue()
}

func (s *GeneratorSuite) TestModuleGenerator_Generate_CliTemplate() {
	appPath := s.createMockApp("cli")
	gen := generator.NewModuleGenerator("commands", "plain", "cli", appPath)

	err := gen.Generate()

	s.T.Expect(err).ToBeNil()

	moduleDir := filepath.Join(appPath, "app", "commands")
	s.T.Expect(dirExists(moduleDir)).ToBeTrue()
}

func (s *GeneratorSuite) TestModuleGenerator_Generate_WebTemplate() {
	appPath := s.createMockApp("web")
	gen := generator.NewModuleGenerator("pages", "plain", "web", appPath)

	err := gen.Generate()

	s.T.Expect(err).ToBeNil()

	moduleDir := filepath.Join(appPath, "app", "pages")
	s.T.Expect(dirExists(moduleDir)).ToBeTrue()
}

func (s *GeneratorSuite) TestModuleGenerator_Generate_UpdatesAppModule() {
	appPath := s.createMockApp("api")
	gen := generator.NewModuleGenerator("accounts", "plain", "api", appPath)

	err := gen.Generate()

	s.T.Expect(err).ToBeNil()

	content, _ := os.ReadFile(filepath.Join(appPath, "app", "app.module.go"))
	contentStr := string(content)
	s.T.Expect(contentStr).ToContainString("accounts")
}

// ==================== DetectAppType Tests ====================

func (s *GeneratorSuite) TestDetectAppType_ApiApp() {
	appPath := s.createMockApp("api")

	appType, err := generator.DetectAppType(appPath)

	s.T.Expect(err).ToBeNil()
	s.T.Expect(appType).ToEqual("api")
}

func (s *GeneratorSuite) TestDetectAppType_CliApp() {
	appPath := s.createMockApp("cli")

	appType, err := generator.DetectAppType(appPath)

	s.T.Expect(err).ToBeNil()
	s.T.Expect(appType).ToEqual("cli")
}

func (s *GeneratorSuite) TestDetectAppType_WebApp() {
	appPath := s.createMockApp("web")

	appType, err := generator.DetectAppType(appPath)

	s.T.Expect(err).ToBeNil()
	s.T.Expect(appType).ToEqual("web")
}

func (s *GeneratorSuite) TestDetectAppType_SpaApp() {
	appPath := s.createMockApp("spa")

	appType, err := generator.DetectAppType(appPath)

	s.T.Expect(err).ToBeNil()
	s.T.Expect(appType).ToEqual("spa")
}

func (s *GeneratorSuite) TestDetectAppType_NoMainGo_ReturnsError() {
	emptyDir := filepath.Join(s.tmpDir, "empty")
	os.MkdirAll(emptyDir, 0755)

	_, err := generator.DetectAppType(emptyDir)

	s.T.Expect(err).Not().ToBeNil()
	s.T.Expect(err.Error()).ToContainString("main.go")
}

func (s *GeneratorSuite) TestDetectAppType_UnknownType_ReturnsError() {
	unknownApp := filepath.Join(s.tmpDir, "unknown")
	os.MkdirAll(unknownApp, 0755)

	mainContent := "package main\nfunc main() {}"
	os.WriteFile(filepath.Join(unknownApp, "main.go"), []byte(mainContent), 0644)

	_, err := generator.DetectAppType(unknownApp)

	s.T.Expect(err).Not().ToBeNil()
	s.T.Expect(err.Error()).ToContainString("could not detect")
}

// ==================== Integration Tests ====================

func (s *GeneratorSuite) TestIntegration_CreateAppThenModule() {
	appPath := filepath.Join(s.tmpDir, "fullapp")
	appGen := generator.NewAppGenerator("FullApp", "api", appPath)
	err := appGen.Generate()

	s.T.Expect(err).ToBeNil()
	s.T.Expect(dirExists(appPath)).ToBeTrue()

	modGen := generator.NewModuleGenerator("users", "plain", "api", appPath)
	err = modGen.Generate()

	s.T.Expect(err).ToBeNil()
	s.T.Expect(dirExists(filepath.Join(appPath, "app", "users"))).ToBeTrue()

	modGen2 := generator.NewModuleGenerator("products", "resource", "api", appPath)
	err = modGen2.Generate()

	s.T.Expect(err).ToBeNil()
	s.T.Expect(dirExists(filepath.Join(appPath, "app", "products"))).ToBeTrue()

	content, _ := os.ReadFile(filepath.Join(appPath, "app", "app.module.go"))
	contentStr := string(content)

	s.T.Expect(contentStr).ToContainString("users")
	s.T.Expect(contentStr).ToContainString("products")
}

func (s *GeneratorSuite) TestIntegration_GeneratedFilesAreValidGo() {
	outputPath := filepath.Join(s.tmpDir, "validgo")
	gen := generator.NewAppGenerator("ValidGo", "api", outputPath)

	err := gen.Generate()

	s.T.Expect(err).ToBeNil()

	goFiles := []string{
		"main.go",
		"app/app.module.go",
		"app/app.controller.go",
		"app/app.service.go",
		"app/app.routes.go",
	}

	for _, file := range goFiles {
		content, err := os.ReadFile(filepath.Join(outputPath, file))
		if err != nil {
			continue
		}
		contentStr := string(content)

		s.T.Expect(contentStr).ToContainString("package")
		s.T.Expect(contentStr).Not().ToContainString("{{.")
		s.T.Expect(contentStr).Not().ToContainString(".tmpl")
	}
}

func (s *GeneratorSuite) TestIntegration_MultipleModuleSameApp() {
	appPath := s.createMockApp("api")

	modules := []string{"users", "posts", "comments", "tags"}

	for _, modName := range modules {
		gen := generator.NewModuleGenerator(modName, "plain", "api", appPath)
		err := gen.Generate()

		s.T.Expect(err).ToBeNil()
		s.T.Expect(dirExists(filepath.Join(appPath, "app", modName))).ToBeTrue()
	}

	content, _ := os.ReadFile(filepath.Join(appPath, "app", "app.module.go"))
	contentStr := string(content)

	for _, modName := range modules {
		s.T.Expect(contentStr).ToContainString(modName)
	}
}
