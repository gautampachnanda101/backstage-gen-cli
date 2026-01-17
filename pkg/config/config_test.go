package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultConfig(t *testing.T) {
	config := getDefaultConfig()

	if config == nil {
		t.Fatal("expected non-nil config")
	}

	// Check LLM defaults
	if config.LLM == nil {
		t.Fatal("expected non-nil LLM config")
	}
	if config.LLM.Provider != "litellm" {
		t.Errorf("expected provider litellm, got %s", config.LLM.Provider)
	}
	if config.LLM.BaseURL != "http://localhost:4000" {
		t.Errorf("expected base_url http://localhost:4000, got %s", config.LLM.BaseURL)
	}
	if config.LLM.Model != "ollama/llama2" {
		t.Errorf("expected model ollama/llama2, got %s", config.LLM.Model)
	}
	if config.LLM.Timeout != 30 {
		t.Errorf("expected timeout 30, got %d", config.LLM.Timeout)
	}

	// Check organization defaults
	if config.Organization.Namespace != "default" {
		t.Errorf("expected namespace default, got %s", config.Organization.Namespace)
	}

	// Check defaults
	if config.Defaults.Owner != "platform-team" {
		t.Errorf("expected owner platform-team, got %s", config.Defaults.Owner)
	}
	if config.Defaults.Lifecycle != "production" {
		t.Errorf("expected lifecycle production, got %s", config.Defaults.Lifecycle)
	}
}

func TestLoadFromFile(t *testing.T) {
	// Create a temporary config file
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configContent := `
llm:
  provider: ollama
  base_url: http://localhost:11434
  model: llama2
  timeout: 60

organization:
  name: test-org
  namespace: test-ns

defaults:
  owner: test-team
  system: test-system
  lifecycle: experimental
  tags:
    - test
    - internal
`

	configPath := filepath.Join(tmpDir, "test-config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	config, err := loadFromFile(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Check LLM config
	if config.LLM.Provider != "ollama" {
		t.Errorf("expected provider ollama, got %s", config.LLM.Provider)
	}
	if config.LLM.BaseURL != "http://localhost:11434" {
		t.Errorf("expected base_url http://localhost:11434, got %s", config.LLM.BaseURL)
	}
	if config.LLM.Model != "llama2" {
		t.Errorf("expected model llama2, got %s", config.LLM.Model)
	}
	if config.LLM.Timeout != 60 {
		t.Errorf("expected timeout 60, got %d", config.LLM.Timeout)
	}

	// Check organization config
	if config.Organization.Name != "test-org" {
		t.Errorf("expected org name test-org, got %s", config.Organization.Name)
	}
	if config.Organization.Namespace != "test-ns" {
		t.Errorf("expected namespace test-ns, got %s", config.Organization.Namespace)
	}

	// Check defaults
	if config.Defaults.Owner != "test-team" {
		t.Errorf("expected owner test-team, got %s", config.Defaults.Owner)
	}
	if config.Defaults.System != "test-system" {
		t.Errorf("expected system test-system, got %s", config.Defaults.System)
	}
	if config.Defaults.Lifecycle != "experimental" {
		t.Errorf("expected lifecycle experimental, got %s", config.Defaults.Lifecycle)
	}
	if len(config.Defaults.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(config.Defaults.Tags))
	}
}

func TestLoadFromFileWithDefaults(t *testing.T) {
	// Create a minimal config file
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Minimal config - should apply defaults
	configContent := `
organization:
  name: my-org
`

	configPath := filepath.Join(tmpDir, "minimal-config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	config, err := loadFromFile(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Check that defaults are applied
	if config.Organization.Namespace != "default" {
		t.Errorf("expected namespace default, got %s", config.Organization.Namespace)
	}
	if config.Defaults.Owner != "platform-team" {
		t.Errorf("expected owner platform-team, got %s", config.Defaults.Owner)
	}
	if config.Defaults.Lifecycle != "production" {
		t.Errorf("expected lifecycle production, got %s", config.Defaults.Lifecycle)
	}
}

func TestLoadFromFileInvalid(t *testing.T) {
	// Create an invalid config file
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configContent := `
this is not valid yaml
  - invalid
    nesting: [broken
`

	configPath := filepath.Join(tmpDir, "invalid-config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err = loadFromFile(configPath)
	if err == nil {
		t.Error("expected error loading invalid config, got nil")
	}
}

func TestLoadFromFileNotFound(t *testing.T) {
	_, err := loadFromFile("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error loading nonexistent file, got nil")
	}
}

func TestCreateExampleConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "example-config.yaml")
	if err := CreateExampleConfig(configPath); err != nil {
		t.Fatalf("failed to create example config: %v", err)
	}

	// Check file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("expected config file to exist")
	}

	// Check file content is valid YAML that can be loaded
	config, err := loadFromFile(configPath)
	if err != nil {
		t.Fatalf("example config is not valid: %v", err)
	}

	// Check some expected values
	if config.LLM.Provider != "litellm" {
		t.Errorf("expected provider litellm in example, got %s", config.LLM.Provider)
	}
}

func TestAppConfigStruct(t *testing.T) {
	config := &AppConfig{
		Organization: OrganizationConfig{
			Name:      "test-org",
			Namespace: "test-ns",
		},
		Defaults: DefaultsConfig{
			Owner:     "team-a",
			System:    "system-b",
			Lifecycle: "production",
			Annotations: map[string]string{
				"custom.io/key": "value",
			},
			Tags: []string{"tag1", "tag2"},
		},
	}

	if config.Organization.Name != "test-org" {
		t.Errorf("expected org name test-org, got %s", config.Organization.Name)
	}
	if config.Defaults.Annotations["custom.io/key"] != "value" {
		t.Error("expected annotation value")
	}
	if len(config.Defaults.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(config.Defaults.Tags))
	}
}
