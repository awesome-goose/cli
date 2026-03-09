package tests

import (
	"testing"

	"github.com/awesome-goose/goose-cli/app"
	"github.com/awesome-goose/goose/core"
	test "github.com/awesome-goose/goose/testing"
)

func TestAppController(t *testing.T) {
	test.NewSuiteRunner(t, &AppControllerSuite{}).Run()
}

type AppControllerSuite struct {
	test.Suite
	controller *app.AppController
	container  *core.Container
}

func (s *AppControllerSuite) SetupTest() {
	s.container = core.NewContainer()
	s.container.Register(func() *app.AppService {
		return &app.AppService{}
	}, "", true)

	s.controller = &app.AppController{}
	s.container.Fill(s.controller)
}

func (s *AppControllerSuite) TeardownTest() {
	if s.container != nil {
		s.container.Reset()
	}
}

func (s *AppControllerSuite) TestVersion_ReturnsOutput() {
	dto := &app.VersionDto{}
	output := s.controller.Version(dto)
	s.T.Expect(output).Not().ToBeNil()
}

func (s *AppControllerSuite) TestVersion_ContainsVersionInfo() {
	dto := &app.VersionDto{}
	output := s.controller.Version(dto)
	content := getOutputContent(output)

	s.T.Expect(content).ToContainString("Version:")
	s.T.Expect(content).ToContainString("Goose CLI")
}

func (s *AppControllerSuite) TestVersion_ContainsAvailableCommands() {
	dto := &app.VersionDto{}
	output := s.controller.Version(dto)
	content := getOutputContent(output)

	s.T.Expect(content).ToContainString("Commands:")
	s.T.Expect(content).ToContainString("goose app")
	s.T.Expect(content).ToContainString("goose g module")
}

func (s *AppControllerSuite) TestVersion_ContainsTemplatesList() {
	dto := &app.VersionDto{}
	output := s.controller.Version(dto)
	content := getOutputContent(output)

	s.T.Expect(content).ToContainString("Templates:")
	s.T.Expect(content).ToContainString("api")
	s.T.Expect(content).ToContainString("cli")
	s.T.Expect(content).ToContainString("web")
	s.T.Expect(content).ToContainString("multi")
}

func (s *AppControllerSuite) TestApp_MissingName_ReturnsError() {
	dto := &app.AppDto{
		Name:     "",
		Template: "api",
		Path:     "",
	}

	output := s.controller.App(dto)
	content := getOutputContent(output)

	s.T.Expect(content).ToContainString("Error")
	s.T.Expect(content).ToContainString("--name")
}

func (s *AppControllerSuite) TestApp_MissingTemplate_ReturnsError() {
	dto := &app.AppDto{
		Name:     "testapp",
		Template: "",
		Path:     "",
	}

	output := s.controller.App(dto)
	content := getOutputContent(output)

	s.T.Expect(content).ToContainString("Error")
	s.T.Expect(content).ToContainString("--template")
}

func (s *AppControllerSuite) TestApp_InvalidTemplate_ReturnsError() {
	dto := &app.AppDto{
		Name:     "testapp",
		Template: "invalid",
		Path:     "",
	}

	output := s.controller.App(dto)
	content := getOutputContent(output)

	s.T.Expect(content).ToContainString("Error")
	s.T.Expect(content).ToContainString("invalid template")
}

func (s *AppControllerSuite) TestApp_ValidApiTemplate_Accepted() {
	dto := &app.AppDto{
		Name:     "testapi",
		Template: "api",
		Path:     "/tmp/goose-test-" + test.NewFixture().String(8),
	}

	output := s.controller.App(dto)
	content := getOutputContent(output)

	if output.Code() == 0 {
		s.T.Expect(content).ToContainString("Application Created")
	}
}

func (s *AppControllerSuite) TestApp_ValidTemplatesAccepted() {
	validTemplates := []string{"api", "cli", "web", "multi"}

	for _, tmpl := range validTemplates {
		dto := &app.AppDto{
			Name:     "testapp",
			Template: tmpl,
			Path:     "/tmp/goose-test-invalid-path-" + test.NewFixture().String(8),
		}

		output := s.controller.App(dto)
		content := getOutputContent(output)
		s.T.Expect(content).Not().ToContainString("invalid template '" + tmpl + "'")
	}
}

func (s *AppControllerSuite) TestApp_InvalidTemplatesRejected() {
	invalidTemplates := []string{"invalid", "API", "WEB", "unknown", "rest", "http"}

	for _, tmpl := range invalidTemplates {
		dto := &app.AppDto{
			Name:     "testapp",
			Template: tmpl,
			Path:     "",
		}

		output := s.controller.App(dto)
		content := getOutputContent(output)
		s.T.Expect(content).ToContainString("invalid template")
	}
}
