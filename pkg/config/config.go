package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gautampachnanda101/backstage-gen-cli/pkg/llm"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	LLM          *llm.Config        `yaml:"llm"`
	Organization OrganizationConfig `yaml:"organization"`
	Defaults     DefaultsConfig     `yaml:"defaults"`
}

type OrganizationConfig struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type DefaultsConfig struct {
	Owner     string `yaml:"owner"`
	System    string `yaml:"system"`
	Lifecycle string `yaml:"lifecycle"`
}

func Load() (*AppConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return getDefaultConfig(), nil
	}

	configPaths := []string{
		filepath.Join(homeDir, ".backstage-gen.yaml"),
		filepath.Join(homeDir, ".backstage-gen.yml"),
		filepath.Join(homeDir, ".config", "backstage-gen", "config.yaml"),
		".backstage-gen.yaml",
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return loadFromFile(path)
		}
	}

	return getDefaultConfig(), nil
}

func loadFromFile(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults
	if config.Organization.Namespace == "" {
		config.Organization.Namespace = "default"
	}
	if config.Defaults.Owner == "" {
		config.Defaults.Owner = "platform-team"
	}
	if config.Defaults.Lifecycle == "" {
		config.Defaults.Lifecycle = "production"
	}

	return &config, nil
}

func getDefaultConfig() *AppConfig {
	return &AppConfig{
		LLM: &llm.Config{
			Provider: "litellm",
			BaseURL:  "http://localhost:4000",
			Model:    "ollama/llama2",
			Timeout:  30,
		},
		Organization: OrganizationConfig{
			Namespace: "default",
		},
		Defaults: DefaultsConfig{
			Owner:     "platform-team",
			Lifecycle: "production",
		},
	}
}

func CreateExampleConfig(path string) error {
	example := `# Backstage Gen CLI Configuration
# Place this file in your home directory as ~/.backstage-gen.yaml

# LLM Configuration (for AI-powered suggestions)
# Using LiteLLM with Docker (recommended)
llm:
  provider: litellm
  base_url: http://localhost:4000
  model: ollama/llama2
  timeout: 30

# Organization settings
organization:
  name: my-org
  namespace: default

# Default values for generated catalogs
defaults:
  owner: platform-team
  system: ""
  lifecycle: production

# ═══════════════════════════════════════════════════════════
# LiteLLM Configuration Examples
# ═══════════════════════════════════════════════════════════
# LiteLLM supports 100+ LLM providers through a unified API
# Start LiteLLM: docker run -d --name litellm -p 4000:4000 ghcr.io/berriai/litellm:main-latest

# Using Ollama (local, free) via LiteLLM:
# llm:
#   provider: litellm
#   base_url: http://localhost:4000
#   model: ollama/llama2

# Using OpenAI via LiteLLM:
# llm:
#   provider: litellm
#   base_url: http://localhost:4000
#   model: gpt-3.5-turbo
#   api_key: your-openai-key

# Using Anthropic Claude via LiteLLM:
# llm:
#   provider: litellm
#   base_url: http://localhost:4000
#   model: claude-2
#   api_key: your-anthropic-key

# Using Direct Providers (without LiteLLM):

# Direct Ollama:
# llm:
#   provider: ollama
#   base_url: http://localhost:11434
#   model: llama2

# OpenRouter (hosted, free tier available):
# llm:
#   provider: openrouter
#   base_url: https://openrouter.ai
#   model: mistralai/mistral-7b-instruct:free
#   api_key: your-api-key
`

	return os.WriteFile(path, []byte(example), 0644)
}
