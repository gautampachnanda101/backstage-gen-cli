package generator

import (
	"testing"

	"github.com/gautampachnanda101/backstage-gen-cli/pkg/detector"
)

func TestNew(t *testing.T) {
	gen := New(nil)
	if gen == nil {
		t.Fatal("expected non-nil generator")
	}
	if gen.config.Organization.Namespace != "default" {
		t.Errorf("expected default namespace, got %s", gen.config.Organization.Namespace)
	}
	if gen.config.Defaults.Owner != "platform-team" {
		t.Errorf("expected default owner platform-team, got %s", gen.config.Defaults.Owner)
	}
}

func TestNewWithConfig(t *testing.T) {
	config := &Config{
		Organization: OrganizationConfig{
			Name:      "test-org",
			Namespace: "test-ns",
		},
		Defaults: DefaultsConfig{
			Owner:     "test-team",
			Lifecycle: "experimental",
		},
	}

	gen := New(config)
	if gen.config.Organization.Name != "test-org" {
		t.Errorf("expected org name test-org, got %s", gen.config.Organization.Name)
	}
	if gen.config.Organization.Namespace != "test-ns" {
		t.Errorf("expected namespace test-ns, got %s", gen.config.Organization.Namespace)
	}
}

func TestGenerateService(t *testing.T) {
	gen := New(&Config{
		Defaults: DefaultsConfig{
			Owner:     "platform-team",
			Lifecycle: "production",
		},
	})

	info := &detector.RepositoryInfo{
		Name:        "my-service",
		Description: "A test service",
		Languages:   []string{"Go"},
		HasDocker:   true,
	}

	catalog, err := gen.Generate(info, "service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if catalog.APIVersion != "backstage.io/v1alpha1" {
		t.Errorf("expected apiVersion backstage.io/v1alpha1, got %s", catalog.APIVersion)
	}
	if catalog.Kind != "Component" {
		t.Errorf("expected kind Component, got %s", catalog.Kind)
	}
	if catalog.Metadata.Name != "my-service" {
		t.Errorf("expected name my-service, got %s", catalog.Metadata.Name)
	}
	if catalog.Spec["type"] != "service" {
		t.Errorf("expected type service, got %v", catalog.Spec["type"])
	}
	if catalog.Spec["owner"] != "platform-team" {
		t.Errorf("expected owner platform-team, got %v", catalog.Spec["owner"])
	}
}

func TestGenerateLibrary(t *testing.T) {
	gen := New(&Config{
		Defaults: DefaultsConfig{
			Owner:     "libs-team",
			Lifecycle: "production",
		},
	})

	info := &detector.RepositoryInfo{
		Name:        "utils-library",
		Description: "A utility library",
		Languages:   []string{"TypeScript"},
	}

	catalog, err := gen.Generate(info, "library")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if catalog.Spec["type"] != "library" {
		t.Errorf("expected type library, got %v", catalog.Spec["type"])
	}
}

func TestGenerateWebsite(t *testing.T) {
	gen := New(&Config{
		Defaults: DefaultsConfig{
			Owner:     "frontend-team",
			Lifecycle: "production",
		},
	})

	info := &detector.RepositoryInfo{
		Name:        "my-website",
		Description: "A website",
		Languages:   []string{"JavaScript"},
		Frameworks:  []string{"React"},
	}

	catalog, err := gen.Generate(info, "website")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if catalog.Spec["type"] != "website" {
		t.Errorf("expected type website, got %v", catalog.Spec["type"])
	}
}

func TestGenerateResource(t *testing.T) {
	gen := New(&Config{
		Defaults: DefaultsConfig{
			Owner:     "infra-team",
			Lifecycle: "production",
		},
	})

	info := &detector.RepositoryInfo{
		Name:         "my-database",
		Description:  "A database resource",
		HasTerraform: true,
	}

	catalog, err := gen.Generate(info, "resource")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if catalog.Kind != "Resource" {
		t.Errorf("expected kind Resource, got %s", catalog.Kind)
	}
	if catalog.Spec["type"] != "infrastructure" {
		t.Errorf("expected type infrastructure, got %v", catalog.Spec["type"])
	}
}

func TestGenerateWithSystem(t *testing.T) {
	gen := New(&Config{
		Defaults: DefaultsConfig{
			Owner:     "platform-team",
			Lifecycle: "production",
			System:    "my-system",
		},
	})

	info := &detector.RepositoryInfo{
		Name: "my-service",
	}

	catalog, err := gen.Generate(info, "service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if catalog.Spec["system"] != "my-system" {
		t.Errorf("expected system my-system, got %v", catalog.Spec["system"])
	}
}

func TestNormalizeName(t *testing.T) {
	gen := New(nil)

	tests := []struct {
		input    string
		expected string
	}{
		{"my-service", "my-service"},
		{"My_Service", "my-service"},
		{"My Service", "my-service"},
		{"MyService123", "myservice123"},
		{"my-service-v2", "my-service-v2"},
		{"Service@#$%", "service"},
		{"UPPERCASE", "uppercase"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := gen.normalizeName(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestAddTechnologyLabels(t *testing.T) {
	gen := New(nil)

	info := &detector.RepositoryInfo{
		Languages:     []string{"Go", "JavaScript"},
		Frameworks:    []string{"React"},
		HasDocker:     true,
		HasKubernetes: true,
	}

	catalog := &Catalog{
		Metadata: Metadata{
			Labels: make(map[string]string),
		},
	}

	gen.addTechnologyLabels(catalog, info)

	if catalog.Metadata.Labels["language"] != "go" {
		t.Errorf("expected language go, got %s", catalog.Metadata.Labels["language"])
	}
	if catalog.Metadata.Labels["framework"] != "react" {
		t.Errorf("expected framework react, got %s", catalog.Metadata.Labels["framework"])
	}
	if catalog.Metadata.Labels["docker"] != "true" {
		t.Errorf("expected docker true, got %s", catalog.Metadata.Labels["docker"])
	}
	if catalog.Metadata.Labels["kubernetes"] != "true" {
		t.Errorf("expected kubernetes true, got %s", catalog.Metadata.Labels["kubernetes"])
	}
}

func TestAddAnnotations(t *testing.T) {
	gen := New(&Config{
		Defaults: DefaultsConfig{
			Annotations: map[string]string{
				"custom.io/key": "value",
			},
		},
	})

	info := &detector.RepositoryInfo{
		GitRemote: "git@github.com:owner/repo.git",
	}

	catalog := &Catalog{
		Metadata: Metadata{
			Annotations: make(map[string]string),
		},
	}

	gen.addAnnotations(catalog, info)

	if catalog.Metadata.Annotations["backstage.io/source-location"] != "url:git@github.com:owner/repo.git" {
		t.Error("expected source-location annotation")
	}
	if catalog.Metadata.Annotations["backstage.io/techdocs-ref"] != "dir:." {
		t.Error("expected techdocs-ref annotation")
	}
	if catalog.Metadata.Annotations["custom.io/key"] != "value" {
		t.Error("expected custom annotation")
	}
	if catalog.Metadata.Annotations["github.com/project-slug"] != "owner/repo" {
		t.Errorf("expected github slug owner/repo, got %s", catalog.Metadata.Annotations["github.com/project-slug"])
	}
}

func TestExtractGitHubSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "SSH format",
			input:    "git@github.com:owner/repo.git",
			expected: "owner/repo",
		},
		{
			name:     "HTTPS format",
			input:    "https://github.com/owner/repo.git",
			expected: "owner/repo",
		},
		{
			name:     "SSH without .git",
			input:    "git@github.com:owner/repo",
			expected: "owner/repo",
		},
		{
			name:     "HTTPS without .git",
			input:    "https://github.com/owner/repo",
			expected: "owner/repo",
		},
		{
			name:     "non-GitHub URL",
			input:    "git@gitlab.com:owner/repo.git",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractGitHubSlug(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGenerateTags(t *testing.T) {
	gen := New(&Config{
		Defaults: DefaultsConfig{
			Tags: []string{"internal"},
		},
	})

	info := &detector.RepositoryInfo{
		Languages:     []string{"Go"},
		Frameworks:    []string{"Gin"},
		BuildTools:    []string{"Make"},
		HasDocker:     true,
		HasKubernetes: true,
	}

	tags := gen.generateTags(info)

	expectedTags := []string{"go", "gin", "docker", "containerized", "kubernetes", "k8s", "cloud-native", "make", "internal"}
	for _, expected := range expectedTags {
		found := false
		for _, tag := range tags {
			if tag == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected tag %s not found in %v", expected, tags)
		}
	}
}

func TestUniqueStrings(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b", "d"}
	result := uniqueStrings(input)

	if len(result) != 4 {
		t.Errorf("expected 4 unique items, got %d", len(result))
	}

	expected := map[string]bool{"a": true, "b": true, "c": true, "d": true}
	for _, item := range result {
		if !expected[item] {
			t.Errorf("unexpected item %s", item)
		}
	}
}

func TestGetOwner(t *testing.T) {
	gen := New(&Config{
		Defaults: DefaultsConfig{
			Owner: "default-team",
		},
	})

	// Test with default owner
	info := &detector.RepositoryInfo{}
	owner := gen.getOwner(info)
	if owner != "default-team" {
		t.Errorf("expected default-team, got %s", owner)
	}

	// Test with info owner
	info.DefaultOwner = "codeowners-team"
	owner = gen.getOwner(info)
	if owner != "codeowners-team" {
		t.Errorf("expected codeowners-team, got %s", owner)
	}
}

func TestCatalogStruct(t *testing.T) {
	catalog := &Catalog{
		APIVersion: "backstage.io/v1alpha1",
		Kind:       "Component",
		Metadata: Metadata{
			Name:        "test-service",
			Description: "A test service",
			Labels:      map[string]string{"language": "go"},
			Annotations: map[string]string{"custom": "value"},
			Tags:        []string{"go", "api"},
		},
		Spec: map[string]interface{}{
			"type":      "service",
			"lifecycle": "production",
			"owner":     "platform-team",
		},
	}

	if catalog.APIVersion != "backstage.io/v1alpha1" {
		t.Error("unexpected apiVersion")
	}
	if catalog.Kind != "Component" {
		t.Error("unexpected kind")
	}
	if catalog.Metadata.Name != "test-service" {
		t.Error("unexpected name")
	}
	if len(catalog.Metadata.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(catalog.Metadata.Tags))
	}
}
