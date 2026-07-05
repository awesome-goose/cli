package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/awesome-goose/goose-cli/app"
	test "github.com/awesome-goose/goose/testing"
)

func TestAppService(t *testing.T) {
	test.NewSuiteRunner(t, &AppServiceSuite{}).Run()
}

type AppServiceSuite struct {
	test.Suite
	service *app.AppService
	tmpDir  string
}

func (s *AppServiceSuite) SetupTest() {
	s.service = &app.AppService{}
	tmpDir, err := os.MkdirTemp("", "goose-cli-test-*")
	if err != nil {
		s.T.Require(err).ToBeNil()
	}
	s.tmpDir = tmpDir
}

func (s *AppServiceSuite) TeardownTest() {
	if s.tmpDir != "" {
		os.RemoveAll(s.tmpDir)
	}
}

func (s *AppServiceSuite) TestGetVersion_ReturnsCorrectVersion() {
	version := s.service.GetVersion()
	s.T.Expect(version).Not().ToBeEmpty()
	s.T.Expect(version).ToEqual("0.0.9")
}

func (s *AppServiceSuite) TestGetVersion_ReturnsConsistentValue() {
	version1 := s.service.GetVersion()
	version2 := s.service.GetVersion()
	s.T.Expect(version1).ToEqual(version2)
}

func (s *AppServiceSuite) TestCreateApp_ApiTemplate_Success() {
	appName := "testapi"
	result, err := s.service.CreateApp(appName, "api", s.tmpDir, "")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()
	s.T.Expect(len(result)).ToBeGreaterThan(0)

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("Created api application")
	s.T.Expect(resultStr).ToContainString(appName)

	appPath := filepath.Join(s.tmpDir, appName)
	s.T.Expect(dirExists(appPath)).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(appPath, "main.go"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(appPath, "go.mod"))).ToBeTrue()
	s.T.Expect(dirExists(filepath.Join(appPath, "app"))).ToBeTrue()
}

func (s *AppServiceSuite) TestCreateApp_CliTemplate_Success() {
	appName := "testcli"
	result, err := s.service.CreateApp(appName, "cli", s.tmpDir, "")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("Created cli application")

	appPath := filepath.Join(s.tmpDir, appName)
	s.T.Expect(dirExists(appPath)).ToBeTrue()
}

func (s *AppServiceSuite) TestCreateApp_WebTemplate_Success() {
	appName := "testweb"
	result, err := s.service.CreateApp(appName, "web", s.tmpDir, "")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("Created web application")

	appPath := filepath.Join(s.tmpDir, appName)
	s.T.Expect(dirExists(appPath)).ToBeTrue()
}

func (s *AppServiceSuite) TestCreateApp_MultiTemplate_Success() {
	appName := "testmulti"
	result, err := s.service.CreateApp(appName, "multi", s.tmpDir, "")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("Created multi application")
	s.T.Expect(resultStr).ToContainString("Run modes:")

	appPath := filepath.Join(s.tmpDir, appName)
	s.T.Expect(dirExists(appPath)).ToBeTrue()
}

func (s *AppServiceSuite) TestCreateApp_SpaTemplate_Success() {
	appName := "testspa"
	result, err := s.service.CreateApp(appName, "spa", s.tmpDir, "react")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("Created spa application")
	s.T.Expect(resultStr).ToContainString("make install")
	s.T.Expect(resultStr).ToContainString("make dist")
	s.T.Expect(resultStr).ToContainString("react")

	appPath := filepath.Join(s.tmpDir, appName)
	s.T.Expect(dirExists(appPath)).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(appPath, "Makefile"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(appPath, "frontend", "package.json"))).ToBeTrue()
}

func (s *AppServiceSuite) TestCreateApp_DirectoryAlreadyExists_ReturnsError() {
	appName := "existing"
	appPath := filepath.Join(s.tmpDir, appName)
	os.MkdirAll(appPath, 0755)

	result, err := s.service.CreateApp(appName, "api", s.tmpDir, "")

	s.T.Expect(err).Not().ToBeNil()
	s.T.Expect(result).ToBeNil()
	s.T.Expect(err.Error()).ToContainString("already exists")
}

func (s *AppServiceSuite) TestCreateApp_EmptyPath_UsesCurrentDirectory() {
	originalDir, _ := os.Getwd()
	os.Chdir(s.tmpDir)
	defer os.Chdir(originalDir)

	appName := "testapp"
	result, err := s.service.CreateApp(appName, "api", "", "")

	s.T.Expect(err).ToBeNil()
	s.T.Expect(result).Not().ToBeNil()

	appPath := filepath.Join(s.tmpDir, appName)
	s.T.Expect(dirExists(appPath)).ToBeTrue()
}

func (s *AppServiceSuite) TestCreateApp_NextStepsIncluded() {
	appName := "testappsteps"
	result, err := s.service.CreateApp(appName, "api", s.tmpDir, "")

	s.T.Expect(err).ToBeNil()

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("Next steps:")
	s.T.Expect(resultStr).ToContainString("cd " + appName)
	s.T.Expect(resultStr).ToContainString("go mod tidy")
	s.T.Expect(resultStr).ToContainString("go run main.go")
}

func (s *AppServiceSuite) TestCreateApp_LocationIncluded() {
	appName := "testapp2"
	result, err := s.service.CreateApp(appName, "api", s.tmpDir, "")

	s.T.Expect(err).ToBeNil()

	resultStr := joinStrings(result)
	s.T.Expect(resultStr).ToContainString("Location:")
}
