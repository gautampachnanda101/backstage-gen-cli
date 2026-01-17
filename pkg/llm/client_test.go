package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	// Test with nil config
	client := NewClient(nil)
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.config.Provider != "litellm" {
		t.Errorf("expected default provider litellm, got %s", client.config.Provider)
	}
	if client.config.BaseURL != "http://localhost:4000" {
		t.Errorf("expected default base_url http://localhost:4000, got %s", client.config.BaseURL)
	}
	if client.config.Timeout != 30 {
		t.Errorf("expected default timeout 30, got %d", client.config.Timeout)
	}
}

func TestNewClientWithConfig(t *testing.T) {
	config := &Config{
		Provider: "ollama",
		BaseURL:  "http://localhost:11434",
		Model:    "llama2",
		APIKey:   "test-key",
		Timeout:  60,
	}

	client := NewClient(config)
	if client.config.Provider != "ollama" {
		t.Errorf("expected provider ollama, got %s", client.config.Provider)
	}
	if client.config.BaseURL != "http://localhost:11434" {
		t.Errorf("expected base_url http://localhost:11434, got %s", client.config.BaseURL)
	}
	if client.config.Model != "llama2" {
		t.Errorf("expected model llama2, got %s", client.config.Model)
	}
	if client.config.APIKey != "test-key" {
		t.Errorf("expected api_key test-key, got %s", client.config.APIKey)
	}
	if client.config.Timeout != 60 {
		t.Errorf("expected timeout 60, got %d", client.config.Timeout)
	}
}

func TestNewClientWithZeroTimeout(t *testing.T) {
	config := &Config{
		Provider: "ollama",
		BaseURL:  "http://localhost:11434",
		Model:    "llama2",
		Timeout:  0,
	}

	client := NewClient(config)
	if client.config.Timeout != 30 {
		t.Errorf("expected default timeout 30 for zero value, got %d", client.config.Timeout)
	}
}

func TestGetEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		baseURL  string
		expected string
	}{
		{
			name:     "ollama",
			provider: "ollama",
			baseURL:  "http://localhost:11434",
			expected: "http://localhost:11434/api/chat",
		},
		{
			name:     "litellm",
			provider: "litellm",
			baseURL:  "http://localhost:4000",
			expected: "http://localhost:4000/chat/completions",
		},
		{
			name:     "openrouter",
			provider: "openrouter",
			baseURL:  "https://openrouter.ai",
			expected: "https://openrouter.ai/api/v1/chat/completions",
		},
		{
			name:     "unknown provider defaults to litellm format",
			provider: "custom",
			baseURL:  "http://custom.api",
			expected: "http://custom.api/chat/completions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(&Config{
				Provider: tt.provider,
				BaseURL:  tt.baseURL,
			})
			endpoint := client.getEndpoint()
			if endpoint != tt.expected {
				t.Errorf("expected endpoint %s, got %s", tt.expected, endpoint)
			}
		})
	}
}

func TestChatWithMockServer(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Return mock response (LiteLLM format)
		response := ChatResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "Hello from mock LLM!"}},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(&Config{
		Provider: "litellm",
		BaseURL:  server.URL,
		Model:    "test-model",
		Timeout:  5,
	})

	result, err := client.Chat("Test prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello from mock LLM!" {
		t.Errorf("expected 'Hello from mock LLM!', got '%s'", result)
	}
}

func TestChatWithOllamaFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return Ollama format response
		response := ChatResponse{
			Message: struct {
				Content string `json:"content"`
			}{Content: "Ollama response"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(&Config{
		Provider: "ollama",
		BaseURL:  server.URL,
		Model:    "llama2",
		Timeout:  5,
	})

	result, err := client.Chat("Test prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Ollama response" {
		t.Errorf("expected 'Ollama response', got '%s'", result)
	}
}

// Note: JSON extraction is tested indirectly through AnalyzeLanguages and AnalyzeDomain

func TestAnalyzeLanguagesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ChatResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: `{"primary_language": "Go", "secondary_languages": ["JavaScript"], "confidence": "high", "reasoning": "go.mod found"}`}},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(&Config{
		Provider: "litellm",
		BaseURL:  server.URL,
		Model:    "test-model",
		Timeout:  5,
	})

	analysis, err := client.AnalyzeLanguages("test file info")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if analysis.PrimaryLanguage != "Go" {
		t.Errorf("expected primary language Go, got %s", analysis.PrimaryLanguage)
	}
	if analysis.Confidence != "high" {
		t.Errorf("expected confidence high, got %s", analysis.Confidence)
	}
	if len(analysis.SecondaryLanguages) != 1 || analysis.SecondaryLanguages[0] != "JavaScript" {
		t.Errorf("expected secondary languages [JavaScript], got %v", analysis.SecondaryLanguages)
	}
}

func TestIsAvailableWithMockServer(t *testing.T) {
	// Server that responds successfully
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(&Config{
		Provider: "litellm",
		BaseURL:  server.URL,
		Timeout:  3,
	})

	if !client.IsAvailable() {
		t.Error("expected IsAvailable to return true for responding server")
	}
}

func TestIsAvailableWithDownServer(t *testing.T) {
	client := NewClient(&Config{
		Provider: "litellm",
		BaseURL:  "http://localhost:59999", // Port that's unlikely to be in use
		Timeout:  1,
	})

	if client.IsAvailable() {
		t.Error("expected IsAvailable to return false for non-responding server")
	}
}

func TestConfigStruct(t *testing.T) {
	config := &Config{
		Provider: "openrouter",
		BaseURL:  "https://openrouter.ai",
		Model:    "mistralai/mistral-7b",
		APIKey:   "sk-test-key",
		Timeout:  45,
	}

	if config.Provider != "openrouter" {
		t.Errorf("expected provider openrouter, got %s", config.Provider)
	}
	if config.BaseURL != "https://openrouter.ai" {
		t.Errorf("expected base_url https://openrouter.ai, got %s", config.BaseURL)
	}
	if config.Model != "mistralai/mistral-7b" {
		t.Errorf("expected model mistralai/mistral-7b, got %s", config.Model)
	}
	if config.APIKey != "sk-test-key" {
		t.Errorf("expected api_key sk-test-key, got %s", config.APIKey)
	}
	if config.Timeout != 45 {
		t.Errorf("expected timeout 45, got %d", config.Timeout)
	}
}

func TestChatRequestJSON(t *testing.T) {
	req := ChatRequest{
		Model: "test-model",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		Stream: false,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	if !strings.Contains(string(data), `"model":"test-model"`) {
		t.Error("expected model in JSON")
	}
	if !strings.Contains(string(data), `"role":"user"`) {
		t.Error("expected role in JSON")
	}
	if !strings.Contains(string(data), `"content":"Hello"`) {
		t.Error("expected content in JSON")
	}
}
