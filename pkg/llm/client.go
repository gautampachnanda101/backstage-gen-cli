package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type Client struct {
	config *Config
}

type Config struct {
	Provider string `yaml:"provider"` // "litellm", "ollama", "openrouter"
	BaseURL  string `yaml:"base_url"`
	Model    string `yaml:"model"`
	APIKey   string `yaml:"api_key,omitempty"`
	Timeout  int    `yaml:"timeout"` // seconds
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices,omitempty"`
}

func NewClient(config *Config) *Client {
	if config == nil {
		// Default to LiteLLM
		config = &Config{
			Provider: "litellm",
			BaseURL:  "http://localhost:4000",
			Model:    "ollama/llama2",
			Timeout:  30,
		}
	}
	if config.Timeout == 0 {
		config.Timeout = 30
	}
	return &Client{config: config}
}

func (c *Client) GenerateCatalogSuggestions(repoInfo, currentCatalog string) (string, error) {
	prompt := fmt.Sprintf(`You are an expert in Backstage catalog creation. Based on the repository information below, suggest improvements or generate a complete catalog-info.yaml.

Repository Information:
%s

Current Catalog (if any):
%s

Please provide:
1. A brief description of the repository
2. Suggested tags (3-5 relevant tags)
3. Any additional metadata or annotations that would be helpful
4. Lifecycle stage recommendation (production, experimental, deprecated)

Keep your response concise and focused on actionable suggestions.`, repoInfo, currentCatalog)

	return c.Chat(prompt)
}

// AnalyzeDomain analyzes the repository to understand its domain and suggest relevant tags
func (c *Client) AnalyzeDomain(repoInfo string) (DomainAnalysis, error) {
	prompt := fmt.Sprintf(`Analyze this repository and determine its domain/purpose. Provide your analysis in JSON format.

Repository Information:
%s

Respond with ONLY valid JSON (no markdown, no extra text) in this exact format:
{
  "domain": "web-framework|api|database|cli-tool|library|frontend|backend|devops|data-processing|ml|other",
  "purpose": "brief one-sentence description",
  "tags": ["tag1", "tag2", "tag3"],
  "annotations": {
    "custom.io/key": "value"
  }
}

Tags should be lowercase, hyphenated, and domain-specific (e.g., "rest-api", "microservice", "monitoring", "ci-cd").
Only include annotations if they're genuinely useful for this specific domain.`, repoInfo)

	response, err := c.Chat(prompt)
	if err != nil {
		return DomainAnalysis{}, err
	}

	// Clean up response - remove markdown code blocks if present
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var analysis DomainAnalysis
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return DomainAnalysis{}, fmt.Errorf("failed to parse LLM response: %w (response: %s)", err, response)
	}

	return analysis, nil
}

// AnalyzeLanguages uses LLM to intelligently detect primary and secondary languages
func (c *Client) AnalyzeLanguages(fileInfo string) (LanguageAnalysis, error) {
	prompt := fmt.Sprintf(`Analyze these file statistics and output JSON.

%s

Rules:
- Primary language = main code language (not config files)
- Cargo.toml → Rust
- go.mod → Go
- package.json with .ts → TypeScript
- package.json with .js → JavaScript
- pom.xml or build.gradle → Java
- requirements.txt or pyproject.toml → Python
- Gemfile → Ruby

Output ONLY valid JSON, nothing else:
{"primary_language":"NAME","secondary_languages":[],"confidence":"high","reasoning":"why"}

Example for Rust project:
{"primary_language":"Rust","secondary_languages":["Python"],"confidence":"high","reasoning":"Cargo.toml found"}

Your JSON response:`, fileInfo)

	response, err := c.Chat(prompt)
	if err != nil {
		return LanguageAnalysis{}, err
	}

	// Extract JSON from response - try multiple strategies
	response = strings.TrimSpace(response)

	// Strategy 1: Try to find a JSON object with "primary_language" key
	// This handles cases where models echo back questions or add explanatory text
	var analysis LanguageAnalysis

	// Find all potential JSON objects by looking for { and matching }
	for i := len(response) - 1; i >= 0; i-- {
		if response[i] == '}' {
			// Found end of potential JSON, search backwards for matching {
			depth := 1
			for j := i - 1; j >= 0 && depth > 0; j-- {
				if response[j] == '}' {
					depth++
				} else if response[j] == '{' {
					depth--
					if depth == 0 {
						// Try to parse this JSON object
						jsonStr := response[j : i+1]
						if err := json.Unmarshal([]byte(jsonStr), &analysis); err == nil {
							// Check if it has the expected field
							if analysis.PrimaryLanguage != "" {
								return analysis, nil
							}
						}
					}
				}
			}
		}
	}

	// Fallback: try first { to last }
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start == -1 || end == -1 || start > end {
		return LanguageAnalysis{}, fmt.Errorf("no valid JSON found in response: %s", response)
	}

	jsonStr := response[start : end+1]

	if err := json.Unmarshal([]byte(jsonStr), &analysis); err != nil {
		return LanguageAnalysis{}, fmt.Errorf("failed to parse JSON: %w (extracted: %s)", err, jsonStr)
	}

	return analysis, nil
}

type DomainAnalysis struct {
	Domain      string            `json:"domain"`
	Purpose     string            `json:"purpose"`
	Tags        []string          `json:"tags"`
	Annotations map[string]string `json:"annotations"`
}

// LanguageAnalysis represents LLM-based language detection results
type LanguageAnalysis struct {
	PrimaryLanguage    string   `json:"primary_language"`
	SecondaryLanguages []string `json:"secondary_languages"`
	Confidence         string   `json:"confidence"`
	Reasoning          string   `json:"reasoning"`
}

func (c *Client) Chat(prompt string) (string, error) {
	endpoint := c.getEndpoint()

	reqBody := ChatRequest{
		Model: c.config.Model,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	client := &http.Client{
		Timeout: time.Duration(c.config.Timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Handle different response formats
	if chatResp.Message.Content != "" {
		return chatResp.Message.Content, nil
	}
	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no content in response")
}

func (c *Client) getEndpoint() string {
	switch c.config.Provider {
	case "litellm":
		return c.config.BaseURL + "/chat/completions"
	case "ollama":
		return c.config.BaseURL + "/api/chat"
	case "openrouter":
		return c.config.BaseURL + "/api/v1/chat/completions"
	default:
		return c.config.BaseURL + "/chat/completions"
	}
}

func (c *Client) IsAvailable() bool {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	var checkURL string
	switch c.config.Provider {
	case "ollama":
		checkURL = c.config.BaseURL + "/api/tags"
	case "litellm":
		checkURL = c.config.BaseURL + "/health"
	default:
		checkURL = c.config.BaseURL
	}

	resp, err := client.Get(checkURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// CheckDockerRunning checks if Docker is running
func CheckDockerRunning() bool {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	return err == nil
}

// CheckLiteLLMRunning checks if LiteLLM container is running
func CheckLiteLLMRunning() bool {
	cmd := exec.Command("docker", "ps", "--filter", "name=litellm", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "litellm")
}

// StartLiteLLM starts the LiteLLM Docker container
func StartLiteLLM() error {
	// Check if container exists but is stopped
	checkCmd := exec.Command("docker", "ps", "-a", "--filter", "name=litellm", "--format", "{{.Names}}")
	output, _ := checkCmd.Output()

	if strings.Contains(string(output), "litellm") {
		// Container exists, just start it
		cmd := exec.Command("docker", "start", "litellm")
		return cmd.Run()
	}

	// Create and start new container
	cmd := exec.Command("docker", "run", "-d",
		"--name", "litellm",
		"-p", "4000:4000",
		"-e", "LITELLM_LOG=INFO",
		"ghcr.io/berriai/litellm:main-latest",
		"--port", "4000",
	)
	return cmd.Run()
}

// StopLiteLLM stops the LiteLLM Docker container
func StopLiteLLM() error {
	cmd := exec.Command("docker", "stop", "litellm")
	return cmd.Run()
}
