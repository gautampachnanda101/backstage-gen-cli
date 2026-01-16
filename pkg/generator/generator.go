package generator

import (
	"fmt"
	"strings"

	"github.com/gautampachnanda101/backstage-gen-cli/pkg/detector"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/llm"
)

type Generator struct {
	config    *Config
	llmClient *llm.Client
}

type Config struct {
	Organization OrganizationConfig
	Defaults     DefaultsConfig
	LLM          *llm.Config
}

type OrganizationConfig struct {
	Name      string
	Namespace string
}

type DefaultsConfig struct {
	Owner       string
	System      string
	Lifecycle   string
	Annotations map[string]string
	Tags        []string
}

type Catalog struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   Metadata               `yaml:"metadata"`
	Spec       map[string]interface{} `yaml:"spec"`
}

type Metadata struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Tags        []string          `yaml:"tags,omitempty"`
}

func New(config *Config) *Generator {
	if config == nil {
		config = &Config{
			Organization: OrganizationConfig{Namespace: "default"},
			Defaults:     DefaultsConfig{Owner: "platform-team", Lifecycle: "production"},
		}
	}

	var llmClient *llm.Client
	if config.LLM != nil {
		llmClient = llm.NewClient(config.LLM)
	}

	return &Generator{
		config:    config,
		llmClient: llmClient,
	}
}

func (g *Generator) Generate(info *detector.RepositoryInfo, template string) (*Catalog, error) {
	return g.GenerateWithLLM(info, template, true)
}

func (g *Generator) GenerateWithLLM(info *detector.RepositoryInfo, template string, useLLM bool) (*Catalog, error) {
	switch template {
	case "service", "library", "website":
		return g.generateComponent(info, template, useLLM)
	case "resource":
		return g.generateResource(info, useLLM)
	default:
		return g.generateComponent(info, "service", useLLM)
	}
}

func (g *Generator) generateComponent(info *detector.RepositoryInfo, compType string, useLLM bool) (*Catalog, error) {
	catalog := &Catalog{
		APIVersion: "backstage.io/v1alpha1",
		Kind:       "Component",
		Metadata: Metadata{
			Name:        g.normalizeName(info.Name),
			Description: info.Description,
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Tags:        []string{},
		},
		Spec: map[string]interface{}{
			"type":      compType,
			"lifecycle": g.config.Defaults.Lifecycle,
			"owner":     g.getOwner(info),
		},
	}

	if g.config.Defaults.System != "" {
		catalog.Spec["system"] = g.config.Defaults.System
	}

	g.addTechnologyLabels(catalog, info)
	g.addAnnotations(catalog, info)
	catalog.Metadata.Tags = g.generateTags(info)

	// Use LLM to enhance with domain-specific tags and annotations
	if useLLM && g.llmClient != nil {
		if err := g.enhanceWithLLM(catalog, info); err == nil {
			// LLM enhancement succeeded
		}
		// Silently ignore LLM errors and continue with basic generation
	}

	return catalog, nil
}

func (g *Generator) generateResource(info *detector.RepositoryInfo, useLLM bool) (*Catalog, error) {
	catalog := &Catalog{
		APIVersion: "backstage.io/v1alpha1",
		Kind:       "Resource",
		Metadata: Metadata{
			Name:        g.normalizeName(info.Name),
			Description: info.Description,
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Tags:        []string{},
		},
		Spec: map[string]interface{}{
			"type":  "infrastructure",
			"owner": g.getOwner(info),
		},
	}

	g.addTechnologyLabels(catalog, info)
	g.addAnnotations(catalog, info)
	catalog.Metadata.Tags = g.generateTags(info)

	if useLLM && g.llmClient != nil {
		if err := g.enhanceWithLLM(catalog, info); err == nil {
			// LLM enhancement succeeded
		}
	}

	return catalog, nil
}

func (g *Generator) enhanceWithLLM(catalog *Catalog, info *detector.RepositoryInfo) error {
	repoSummary := fmt.Sprintf(`Name: %s
Description: %s
Languages: %v
Frameworks: %v
Build Tools: %v
Has Docker: %v
Has Kubernetes: %v`,
		info.Name, info.Description, info.Languages, info.Frameworks,
		info.BuildTools, info.HasDocker, info.HasKubernetes)

	analysis, err := g.llmClient.AnalyzeDomain(repoSummary)
	if err != nil {
		return err
	}

	// Merge LLM-suggested tags (remove duplicates)
	existingTags := make(map[string]bool)
	for _, tag := range catalog.Metadata.Tags {
		existingTags[tag] = true
	}

	for _, tag := range analysis.Tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag != "" && !existingTags[tag] {
			catalog.Metadata.Tags = append(catalog.Metadata.Tags, tag)
			existingTags[tag] = true
		}
	}

	// Add domain tag if meaningful
	if analysis.Domain != "" && analysis.Domain != "other" {
		domainTag := strings.ToLower(analysis.Domain)
		if !existingTags[domainTag] {
			catalog.Metadata.Tags = append(catalog.Metadata.Tags, domainTag)
		}
	}

	// Add LLM-suggested annotations (don't override existing standard ones)
	for key, value := range analysis.Annotations {
		if _, exists := catalog.Metadata.Annotations[key]; !exists {
			catalog.Metadata.Annotations[key] = value
		}
	}

	// Use LLM purpose if description is empty or generic
	if catalog.Metadata.Description == "" || len(catalog.Metadata.Description) < 20 {
		if analysis.Purpose != "" {
			catalog.Metadata.Description = analysis.Purpose
		}
	}

	return nil
}

func (g *Generator) normalizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, " ", "-")

	var result strings.Builder
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func (g *Generator) getOwner(info *detector.RepositoryInfo) string {
	if info.DefaultOwner != "" {
		return info.DefaultOwner
	}
	return g.config.Defaults.Owner
}

func (g *Generator) addTechnologyLabels(catalog *Catalog, info *detector.RepositoryInfo) {
	if len(info.Languages) > 0 {
		catalog.Metadata.Labels["language"] = strings.ToLower(info.Languages[0])
	}
	if len(info.Frameworks) > 0 {
		catalog.Metadata.Labels["framework"] = strings.ToLower(info.Frameworks[0])
	}
	if info.HasDocker {
		catalog.Metadata.Labels["docker"] = "true"
	}
	if info.HasKubernetes {
		catalog.Metadata.Labels["kubernetes"] = "true"
	}
}

func (g *Generator) addAnnotations(catalog *Catalog, info *detector.RepositoryInfo) {
	// Standard Backstage annotations
	if info.GitRemote != "" {
		catalog.Metadata.Annotations["backstage.io/source-location"] =
			fmt.Sprintf("url:%s", info.GitRemote)

		// Add GitHub project slug if it's a GitHub repo
		if strings.Contains(info.GitRemote, "github.com") {
			// Extract owner/repo from git remote
			slug := extractGitHubSlug(info.GitRemote)
			if slug != "" {
				catalog.Metadata.Annotations["github.com/project-slug"] = slug
			}
		}
	}

	// TechDocs reference (default to repo root)
	catalog.Metadata.Annotations["backstage.io/techdocs-ref"] = "dir:."

	// Add custom annotations from config
	for key, value := range g.config.Defaults.Annotations {
		catalog.Metadata.Annotations[key] = value
	}
}

func extractGitHubSlug(remote string) string {
	// Handle both HTTPS and SSH formats
	// https://github.com/owner/repo.git -> owner/repo
	// git@github.com:owner/repo.git -> owner/repo
	remote = strings.TrimSuffix(remote, ".git")

	if strings.Contains(remote, "github.com/") {
		parts := strings.Split(remote, "github.com/")
		if len(parts) == 2 {
			return parts[1]
		}
	} else if strings.Contains(remote, "github.com:") {
		parts := strings.Split(remote, "github.com:")
		if len(parts) == 2 {
			return parts[1]
		}
	}

	return ""
}

func (g *Generator) generateTags(info *detector.RepositoryInfo) []string {
	tags := []string{}

	// Add language tags
	for _, lang := range info.Languages {
		tags = append(tags, strings.ToLower(lang))
	}

	// Add framework tags
	for _, fw := range info.Frameworks {
		tags = append(tags, strings.ToLower(fw))
	}

	// Add infrastructure tags
	if info.HasDocker {
		tags = append(tags, "docker", "containerized")
	}
	if info.HasKubernetes {
		tags = append(tags, "kubernetes", "k8s", "cloud-native")
	}

	// Add build tool tags
	for _, tool := range info.BuildTools {
		tags = append(tags, strings.ToLower(tool))
	}

	// Add custom tags from config
	tags = append(tags, g.config.Defaults.Tags...)

	return uniqueStrings(tags)
}

func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
