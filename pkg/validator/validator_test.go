package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	v := New(nil)
	if v == nil {
		t.Fatal("expected non-nil validator")
	}
	if v.config == nil {
		t.Fatal("expected non-nil config")
	}
	if len(v.config.RequiredFields) != 3 {
		t.Errorf("expected 3 required fields, got %d", len(v.config.RequiredFields))
	}
}

func TestNewWithConfig(t *testing.T) {
	config := &Config{
		Strict: true,
		RequiredFields: []string{
			"metadata.name",
		},
	}

	v := New(config)
	if !v.config.Strict {
		t.Error("expected strict mode to be true")
	}
	if len(v.config.RequiredFields) != 1 {
		t.Errorf("expected 1 required field, got %d", len(v.config.RequiredFields))
	}
}

func TestValidateFileValidCatalog(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	validCatalog := `apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: my-service
  description: A test service
  tags:
    - go
    - api
spec:
  type: service
  lifecycle: production
  owner: platform-team
`

	catalogPath := filepath.Join(tmpDir, "catalog-info.yaml")
	if err := os.WriteFile(catalogPath, []byte(validCatalog), 0644); err != nil {
		t.Fatalf("failed to write catalog file: %v", err)
	}

	v := New(nil)
	result, err := v.ValidateFile(catalogPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.HasErrors() {
		t.Errorf("expected no errors, got %v", result.Errors)
	}
	if result.HasWarnings() {
		t.Logf("warnings: %v", result.Warnings)
	}
}

func TestValidateFileMissingAPIVersion(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	invalidCatalog := `kind: Component
metadata:
  name: my-service
spec:
  type: service
  owner: platform-team
`

	catalogPath := filepath.Join(tmpDir, "catalog-info.yaml")
	if err := os.WriteFile(catalogPath, []byte(invalidCatalog), 0644); err != nil {
		t.Fatalf("failed to write catalog file: %v", err)
	}

	v := New(nil)
	result, err := v.ValidateFile(catalogPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.HasErrors() {
		t.Error("expected errors for missing apiVersion")
	}

	foundError := false
	for _, e := range result.Errors {
		if e.Field == "apiVersion" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected error for apiVersion field")
	}
}

func TestValidateFileMissingKind(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	invalidCatalog := `apiVersion: backstage.io/v1alpha1
metadata:
  name: my-service
spec:
  type: service
  owner: platform-team
`

	catalogPath := filepath.Join(tmpDir, "catalog-info.yaml")
	if err := os.WriteFile(catalogPath, []byte(invalidCatalog), 0644); err != nil {
		t.Fatalf("failed to write catalog file: %v", err)
	}

	v := New(nil)
	result, err := v.ValidateFile(catalogPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.HasErrors() {
		t.Error("expected errors for missing kind")
	}
}

func TestValidateFileInvalidKind(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	invalidCatalog := `apiVersion: backstage.io/v1alpha1
kind: InvalidKind
metadata:
  name: my-service
spec:
  type: service
  owner: platform-team
`

	catalogPath := filepath.Join(tmpDir, "catalog-info.yaml")
	if err := os.WriteFile(catalogPath, []byte(invalidCatalog), 0644); err != nil {
		t.Fatalf("failed to write catalog file: %v", err)
	}

	v := New(nil)
	result, err := v.ValidateFile(catalogPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.HasErrors() {
		t.Error("expected errors for invalid kind")
	}

	foundError := false
	for _, e := range result.Errors {
		if e.Rule == "invalid-kind" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected invalid-kind error")
	}
}

func TestValidateFileMissingMetadata(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	invalidCatalog := `apiVersion: backstage.io/v1alpha1
kind: Component
spec:
  type: service
  owner: platform-team
`

	catalogPath := filepath.Join(tmpDir, "catalog-info.yaml")
	if err := os.WriteFile(catalogPath, []byte(invalidCatalog), 0644); err != nil {
		t.Fatalf("failed to write catalog file: %v", err)
	}

	v := New(nil)
	result, err := v.ValidateFile(catalogPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.HasErrors() {
		t.Error("expected errors for missing metadata")
	}
}

func TestValidateFileMissingSpec(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	invalidCatalog := `apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: my-service
`

	catalogPath := filepath.Join(tmpDir, "catalog-info.yaml")
	if err := os.WriteFile(catalogPath, []byte(invalidCatalog), 0644); err != nil {
		t.Fatalf("failed to write catalog file: %v", err)
	}

	v := New(nil)
	result, err := v.ValidateFile(catalogPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.HasErrors() {
		t.Error("expected errors for missing spec")
	}
}

func TestValidateFileInvalidName(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	invalidCatalog := `apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: My_Invalid_Name
spec:
  type: service
  owner: platform-team
`

	catalogPath := filepath.Join(tmpDir, "catalog-info.yaml")
	if err := os.WriteFile(catalogPath, []byte(invalidCatalog), 0644); err != nil {
		t.Fatalf("failed to write catalog file: %v", err)
	}

	v := New(nil)
	result, err := v.ValidateFile(catalogPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foundError := false
	for _, e := range result.Errors {
		if e.Rule == "name-format" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected name-format error")
	}
}

func TestValidateFileInvalidYAML(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	invalidYAML := `
this is not valid yaml
  - broken
    nesting: [incomplete
`

	catalogPath := filepath.Join(tmpDir, "catalog-info.yaml")
	if err := os.WriteFile(catalogPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("failed to write catalog file: %v", err)
	}

	v := New(nil)
	result, err := v.ValidateFile(catalogPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.HasErrors() {
		t.Error("expected errors for invalid YAML")
	}

	foundError := false
	for _, e := range result.Errors {
		if e.Rule == "yaml-syntax" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected yaml-syntax error")
	}
}

func TestValidateFileNotFound(t *testing.T) {
	v := New(nil)
	_, err := v.ValidateFile("/nonexistent/path/catalog.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestValidateFileMissingDescription(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	catalogWithoutDesc := `apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: my-service
spec:
  type: service
  owner: platform-team
`

	catalogPath := filepath.Join(tmpDir, "catalog-info.yaml")
	if err := os.WriteFile(catalogPath, []byte(catalogWithoutDesc), 0644); err != nil {
		t.Fatalf("failed to write catalog file: %v", err)
	}

	v := New(nil)
	result, err := v.ValidateFile(catalogPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.HasWarnings() {
		t.Error("expected warning for missing description")
	}

	foundWarning := false
	for _, w := range result.Warnings {
		if w.Field == "metadata.description" {
			foundWarning = true
			break
		}
	}
	if !foundWarning {
		t.Error("expected description warning")
	}
}

func TestValidateFileMissingTags(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	catalogWithoutTags := `apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: my-service
  description: A test service
spec:
  type: service
  owner: platform-team
`

	catalogPath := filepath.Join(tmpDir, "catalog-info.yaml")
	if err := os.WriteFile(catalogPath, []byte(catalogWithoutTags), 0644); err != nil {
		t.Fatalf("failed to write catalog file: %v", err)
	}

	v := New(nil)
	result, err := v.ValidateFile(catalogPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foundWarning := false
	for _, w := range result.Warnings {
		if w.Field == "metadata.tags" {
			foundWarning = true
			break
		}
	}
	if !foundWarning {
		t.Error("expected tags warning")
	}
}

func TestValidateFileComponentMissingType(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	catalogWithoutType := `apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: my-service
spec:
  owner: platform-team
`

	catalogPath := filepath.Join(tmpDir, "catalog-info.yaml")
	if err := os.WriteFile(catalogPath, []byte(catalogWithoutType), 0644); err != nil {
		t.Fatalf("failed to write catalog file: %v", err)
	}

	v := New(nil)
	result, err := v.ValidateFile(catalogPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foundError := false
	for _, e := range result.Errors {
		if e.Field == "spec.type" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected spec.type error for Component")
	}
}

func TestValidationResultHasErrors(t *testing.T) {
	result := &ValidationResult{
		Errors: []ValidationError{
			{Message: "test error"},
		},
	}

	if !result.HasErrors() {
		t.Error("expected HasErrors to return true")
	}

	emptyResult := &ValidationResult{}
	if emptyResult.HasErrors() {
		t.Error("expected HasErrors to return false for empty result")
	}
}

func TestValidationResultHasWarnings(t *testing.T) {
	result := &ValidationResult{
		Warnings: []ValidationWarning{
			{Message: "test warning"},
		},
	}

	if !result.HasWarnings() {
		t.Error("expected HasWarnings to return true")
	}

	emptyResult := &ValidationResult{}
	if emptyResult.HasWarnings() {
		t.Error("expected HasWarnings to return false for empty result")
	}
}

func TestContains(t *testing.T) {
	slice := []string{"Component", "API", "Resource"}

	if !contains(slice, "Component") {
		t.Error("expected contains to return true for Component")
	}
	if !contains(slice, "API") {
		t.Error("expected contains to return true for API")
	}
	if contains(slice, "Invalid") {
		t.Error("expected contains to return false for Invalid")
	}
}

func TestValidKinds(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	validKinds := []string{"Component", "API", "Resource", "System", "Domain"}

	for _, kind := range validKinds {
		t.Run(kind, func(t *testing.T) {
			catalog := `apiVersion: backstage.io/v1alpha1
kind: ` + kind + `
metadata:
  name: test-entity
spec:
  owner: test-team
`
			if kind == "Component" {
				catalog += "  type: service\n"
			}

			catalogPath := filepath.Join(tmpDir, "catalog-"+kind+".yaml")
			if err := os.WriteFile(catalogPath, []byte(catalog), 0644); err != nil {
				t.Fatalf("failed to write catalog file: %v", err)
			}

			v := New(nil)
			result, err := v.ValidateFile(catalogPath)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check no "invalid-kind" error
			for _, e := range result.Errors {
				if e.Rule == "invalid-kind" {
					t.Errorf("kind %s should be valid, got invalid-kind error", kind)
				}
			}
		})
	}
}
