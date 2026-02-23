package generate

import (
	"fmt"
	"strings"

	"github.com/awesome-goose/goose/io/output"
	"github.com/awesome-goose/goose/types"
)

type GenerateController struct {
	generateService *GenerateService `inject:""`
	log             types.Log        `inject:""`
}

// Module generates a new module in the current Goose application
func (c *GenerateController) Module(body *ModuleDto) types.Output {
	c.log.Info("GenerateController: Module called", map[string]any{
		"name":     body.Name,
		"type":     body.Type,
		"template": body.Template,
	})

	// Validate inputs
	if body.Name == "" {
		return output.ConsoleError("Error: --name flag is required")
	}

	// Set default type
	moduleType := body.Type
	if moduleType == "" {
		moduleType = "plain"
	}

	// Validate module type
	validTypes := []string{"plain", "resource"}
	isValidType := false
	for _, t := range validTypes {
		if moduleType == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return output.ConsoleError(fmt.Sprintf("Error: invalid module type '%s'. Valid types: plain, resource", moduleType))
	}

	// Generate the module
	result, err := c.generateService.GenerateModule(body.Name, moduleType, body.Template)
	if err != nil {
		return output.ConsoleError(fmt.Sprintf("Error: %v", err))
	}

	return output.Box(fmt.Sprintf("Module '%s' Generated", strings.Title(body.Name)), result)
}
