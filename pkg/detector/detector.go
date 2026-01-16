package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Detector struct {
	rootPath string
}

type RepositoryInfo struct {
	Name          string              `json:"name"`
	Path          string              `json:"path"`
	Type          string              `json:"type"`
	Description   string              `json:"description"`
	GitRemote     string              `json:"git_remote,omitempty"`
	DefaultOwner  string              `json:"default_owner,omitempty"`
	Languages     []string            `json:"languages"`
	Frameworks    []string            `json:"frameworks"`
	BuildTools    []string            `json:"build_tools"`
	HasDocker     bool                `json:"has_docker"`
	HasKubernetes bool                `json:"has_kubernetes"`
	HasHelm       bool                `json:"has_helm"`
	HasTerraform  bool                `json:"has_terraform"`
	Dependencies  []Dependency        `json:"dependencies,omitempty"`
	KeyFiles      []string            `json:"key_files,omitempty"`
	Metadata      map[string]string   `json:"metadata,omitempty"`
}

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Type    string `json:"type"`
}

func New(rootPath string) *Detector {
	return &Detector{rootPath: rootPath}
}

func (d *Detector) Detect() (*RepositoryInfo, error) {
	info := &RepositoryInfo{
		Path:         d.rootPath,
		Name:         filepath.Base(d.rootPath),
		Languages:    []string{},
		Frameworks:   []string{},
		BuildTools:   []string{},
		Dependencies: []Dependency{},
		KeyFiles:     []string{},
		Metadata:     make(map[string]string),
	}

	d.detectGitInfo(info)
	d.detectLanguages(info)
	d.detectFrameworks(info)
	d.detectBuildTools(info)
	d.detectInfrastructure(info)
	d.detectType(info)
	d.detectDescription(info)
	d.detectOwner(info)

	return info, nil
}

func (d *Detector) detectGitInfo(info *RepositoryInfo) {
	repo, err := git.PlainOpen(d.rootPath)
	if err != nil {
		return
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return
	}

	if len(remotes) > 0 {
		config := remotes[0].Config()
		if len(config.URLs) > 0 {
			info.GitRemote = config.URLs[0]
			if name := extractRepoName(config.URLs[0]); name != "" {
				info.Name = name
			}
		}
	}
}

func (d *Detector) detectLanguages(info *RepositoryInfo) {
	languageMarkers := map[string][]string{
		"Go":         {"go.mod", "go.sum"},
		"Python":     {"requirements.txt", "setup.py", "pyproject.toml"},
		"JavaScript": {"package.json"},
		"TypeScript": {"tsconfig.json"},
		"Java":       {"pom.xml", "build.gradle"},
		"Rust":       {"Cargo.toml"},
		"Ruby":       {"Gemfile"},
		"PHP":        {"composer.json"},
		"C#":         {".csproj", ".sln"},
	}

	for lang, markers := range languageMarkers {
		for _, marker := range markers {
			if d.fileExists(marker) {
				info.Languages = append(info.Languages, lang)
				info.KeyFiles = append(info.KeyFiles, marker)
				break
			}
		}
	}
}

func (d *Detector) detectFrameworks(info *RepositoryInfo) {
	if d.fileExists("package.json") {
		content, _ := os.ReadFile(filepath.Join(d.rootPath, "package.json"))
		text := string(content)
		if strings.Contains(text, "\"react\"") {
			info.Frameworks = append(info.Frameworks, "React")
		}
		if strings.Contains(text, "\"vue\"") {
			info.Frameworks = append(info.Frameworks, "Vue")
		}
		if strings.Contains(text, "\"express\"") {
			info.Frameworks = append(info.Frameworks, "Express")
		}
	}

	if d.fileExists("requirements.txt") {
		content, _ := os.ReadFile(filepath.Join(d.rootPath, "requirements.txt"))
		text := string(content)
		if strings.Contains(text, "django") {
			info.Frameworks = append(info.Frameworks, "Django")
		}
		if strings.Contains(text, "flask") {
			info.Frameworks = append(info.Frameworks, "Flask")
		}
	}
}

func (d *Detector) detectBuildTools(info *RepositoryInfo) {
	buildToolMarkers := map[string]string{
		"Maven":  "pom.xml",
		"Gradle": "build.gradle",
		"npm":    "package.json",
		"pip":    "requirements.txt",
		"cargo":  "Cargo.toml",
		"Make":   "Makefile",
	}

	for tool, marker := range buildToolMarkers {
		if d.fileExists(marker) {
			info.BuildTools = append(info.BuildTools, tool)
		}
	}
}

func (d *Detector) detectInfrastructure(info *RepositoryInfo) {
	if d.fileExists("Dockerfile") || d.fileExists("docker-compose.yml") {
		info.HasDocker = true
		info.KeyFiles = append(info.KeyFiles, "Dockerfile")
	}

	if d.dirExists("k8s") || d.dirExists("kubernetes") {
		info.HasKubernetes = true
		info.KeyFiles = append(info.KeyFiles, "k8s/")
	}

	if d.fileExists("Chart.yaml") {
		info.HasHelm = true
	}

	if d.hasFilesWithExtension(".tf") {
		info.HasTerraform = true
	}
}

func (d *Detector) detectType(info *RepositoryInfo) {
	if info.HasKubernetes || info.HasHelm || info.HasTerraform {
		info.Type = "resource"
		return
	}

	if d.isLibrary(info) {
		info.Type = "library"
		return
	}

	if d.isFrontend(info) {
		info.Type = "website"
		return
	}

	info.Type = "service"
}

func (d *Detector) detectDescription(info *RepositoryInfo) {
	for _, name := range []string{"README.md", "readme.md"} {
		path := filepath.Join(d.rootPath, name)
		if content, err := os.ReadFile(path); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") && len(line) > 10 {
					info.Description = line
					if len(info.Description) > 200 {
						info.Description = info.Description[:197] + "..."
					}
					return
				}
			}
		}
	}
	info.Description = fmt.Sprintf("A %s application", info.Type)
}

func (d *Detector) detectOwner(info *RepositoryInfo) {
	for _, path := range []string{".github/CODEOWNERS", "CODEOWNERS"} {
		fullPath := filepath.Join(d.rootPath, path)
		if content, err := os.ReadFile(fullPath); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "*") {
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						info.DefaultOwner = strings.TrimPrefix(parts[1], "@")
						return
					}
				}
			}
		}
	}
}

func (d *Detector) isLibrary(info *RepositoryInfo) bool {
	if info.HasDocker || info.HasKubernetes {
		return false
	}
	return d.fileExists("setup.py") && !d.fileExists("main.py")
}

func (d *Detector) isFrontend(info *RepositoryInfo) bool {
	frontendIndicators := []string{"public/index.html", "index.html"}
	for _, indicator := range frontendIndicators {
		if d.fileExists(indicator) {
			return true
		}
	}
	return false
}

func (d *Detector) fileExists(path string) bool {
	fullPath := filepath.Join(d.rootPath, path)
	info, err := os.Stat(fullPath)
	return err == nil && !info.IsDir()
}

func (d *Detector) dirExists(path string) bool {
	fullPath := filepath.Join(d.rootPath, path)
	info, err := os.Stat(fullPath)
	return err == nil && info.IsDir()
}

func (d *Detector) hasFilesWithExtension(ext string) bool {
	found := false
	filepath.Walk(d.rootPath, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ext) {
			found = true
			return filepath.SkipDir
		}
		return nil
	})
	return found
}

func extractRepoName(url string) string {
	re := regexp.MustCompile(`([^/]+?)(\.git)?$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return strings.TrimSuffix(matches[1], ".git")
	}
	return ""
}
