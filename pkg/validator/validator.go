package validator

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Validator struct {
	config *Config
}

type Config struct {
	Strict         bool
	RequiredFields []string
}

type ValidationResult struct {
	File     string
	Errors   []ValidationError
	Warnings []ValidationWarning
}

type ValidationError struct {
	Field   string
	Line    int
	Message string
	Rule    string
}

type ValidationWarning struct {
	Field   string
	Line    int
	Message string
	Rule    string
}

type CatalogEntity struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   map[string]interface{} `yaml:"metadata"`
	Spec       map[string]interface{} `yaml:"spec"`
}

func New(config *Config) *Validator {
	if config == nil {
		config = &Config{
			RequiredFields: []string{
				"metadata.name",
				"metadata.description",
				"spec.owner",
			},
		}
	}
	return &Validator{config: config}
}

func (v *Validator) ValidateFile(path string) (ValidationResult, error) {
	result := ValidationResult{
		File:     path,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return result, fmt.Errorf("failed to read file: %w", err)
	}

	var entity CatalogEntity
	if err := yaml.Unmarshal(data, &entity); err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Message: fmt.Sprintf("Invalid YAML: %v", err),
			Rule:    "yaml-syntax",
		})
		return result, nil
	}

	v.validateSchema(&result, &entity)
	v.validateRequiredFields(&result, &entity)
	v.validateFieldFormats(&result, &entity)
	v.checkBestPractices(&result, &entity)

	return result, nil
}

func (v *Validator) validateSchema(result *ValidationResult, entity *CatalogEntity) {
	if entity.APIVersion == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "apiVersion",
			Message: "apiVersion is required",
			Rule:    "required-field",
		})
	}

	validKinds := []string{"Component", "API", "Resource", "System", "Domain"}
	if entity.Kind == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "kind",
			Message: "kind is required",
			Rule:    "required-field",
		})
	} else if !contains(validKinds, entity.Kind) {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "kind",
			Message: fmt.Sprintf("kind must be one of: %s", strings.Join(validKinds, ", ")),
			Rule:    "invalid-kind",
		})
	}

	if entity.Metadata == nil {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "metadata",
			Message: "metadata is required",
			Rule:    "required-field",
		})
	}

	if entity.Spec == nil {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "spec",
			Message: "spec is required",
			Rule:    "required-field",
		})
	}
}

func (v *Validator) validateRequiredFields(result *ValidationResult, entity *CatalogEntity) {
	if name, ok := entity.Metadata["name"].(string); !ok || name == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "metadata.name",
			Message: "metadata.name is required",
			Rule:    "required-field",
		})
	}

	if owner, ok := entity.Spec["owner"].(string); !ok || owner == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "spec.owner",
			Message: "spec.owner is required",
			Rule:    "required-field",
		})
	}

	if entity.Kind == "Component" {
		if compType, ok := entity.Spec["type"].(string); !ok || compType == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "spec.type",
				Message: "spec.type is required for Component",
				Rule:    "required-field",
			})
		}
	}
}

func (v *Validator) validateFieldFormats(result *ValidationResult, entity *CatalogEntity) {
	if name, ok := entity.Metadata["name"].(string); ok {
		namePattern := regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
		if !namePattern.MatchString(name) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "metadata.name",
				Message: "metadata.name must be lowercase alphanumeric with hyphens",
				Rule:    "name-format",
			})
		}
	}

	if owner, ok := entity.Spec["owner"].(string); ok {
		ownerPattern := regexp.MustCompile(`^[a-z0-9-]+$`)
		if !ownerPattern.MatchString(owner) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "spec.owner",
				Message: "spec.owner must be lowercase alphanumeric with hyphens",
				Rule:    "owner-format",
			})
		}
	}
}

func (v *Validator) checkBestPractices(result *ValidationResult, entity *CatalogEntity) {
	if desc, ok := entity.Metadata["description"].(string); !ok || desc == "" {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "metadata.description",
			Message: "description is recommended",
			Rule:    "best-practice",
		})
	}

	if tags, ok := entity.Metadata["tags"].([]interface{}); !ok || len(tags) == 0 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "metadata.tags",
			Message: "tags improve discoverability",
			Rule:    "best-practice",
		})
	}
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ValidationResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
