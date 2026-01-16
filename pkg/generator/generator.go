package generator

import (
	"fmt"
	"strings"
	"time"

	"github.com/gautampachnanda101/backstage-gen-cli/pkg/detector"
)

type Generator struct {
	config *Config
}

type Config struct {
	Organization OrganizationConfig
	Defaults     DefaultsConfig
}

type OrganizationConfig struct {
	Name      string
	Namespace string
}

type DefaultsConfig struct {
	Owner     string
	System    string
	Lifecycle string
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
	return &Generator{config: config}
}

func (g *Generator) Generate(info *detector.RepositoryInfo, template string) (*Catalog, error) {
	switch template {
	case "service", "library", "website":
		return g.generateComponent(info, template)
	case "resource":
		return g.generateResource(info)
	default:
		return g.generateComponent(info, "service")
	}
}

func (g *Generator) generateComponent(info *detector.RepositoryInfo, compType string) (*Catalog, error) {
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

	return catalog, nil
}

func (g *Generator) generateResource(info *detector.RepositoryInfo) (*Catalog, error) {
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

	return catalog, nil
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
	if info.GitRemote != "" {
		catalog.Metadata.Annotations["backstage.io/source-location"] =
			fmt.Sprintf("url:%s", info.GitRemote)
	}
	catalog.Metadata.Annotations["generated-by"] = "backstage-gen"
	catalog.Metadata.Annotations["generated-at"] = time.Now().Format(time.RFC3339)
}

func (g *Generator) generateTags(info *detector.RepositoryInfo) []string {
	tags := []string{}

	for _, lang := range info.Languages {
		tags = append(tags, strings.ToLower(lang))
	}

	for _, fw := range info.Frameworks {
		tags = append(tags, strings.ToLower(fw))
	}

	if info.HasDocker {
		tags = append(tags, "docker")
	}
	if info.HasKubernetes {
		tags = append(tags, "kubernetes")
	}

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
