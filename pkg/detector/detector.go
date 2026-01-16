package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Detector struct {
	rootPath string
}

type RepositoryInfo struct {
	Name          string            `json:"name"`
	Path          string            `json:"path"`
	Type          string            `json:"type"`
	Description   string            `json:"description"`
	GitRemote     string            `json:"git_remote,omitempty"`
	DefaultOwner  string            `json:"default_owner,omitempty"`
	OS            string            `json:"os"`
	Platform      string            `json:"platform"`
	Arch          string            `json:"arch"`
	Languages     []string          `json:"languages"`
	Frameworks    []string          `json:"frameworks"`
	BuildTools    []string          `json:"build_tools"`
	HasDocker     bool              `json:"has_docker"`
	HasKubernetes bool              `json:"has_kubernetes"`
	HasHelm       bool              `json:"has_helm"`
	HasTerraform  bool              `json:"has_terraform"`
	Dependencies  []Dependency      `json:"dependencies,omitempty"`
	KeyFiles      []string          `json:"key_files,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
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
		OS:           runtime.GOOS,
		Platform:     getPlatformName(),
		Arch:         runtime.GOARCH,
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
		"Go":         {"go.mod", "go.sum", "main.go"},
		"Python":     {"requirements.txt", "setup.py", "pyproject.toml", "Pipfile"},
		"JavaScript": {"package.json", "yarn.lock"},
		"TypeScript": {"tsconfig.json", "package.json"},
		"Java":       {"pom.xml", "build.gradle", "build.gradle.kts"},
		"Rust":       {"Cargo.toml", "Cargo.lock"},
		"Ruby":       {"Gemfile", "Rakefile", ".ruby-version"},
		"PHP":        {"composer.json", "composer.lock"},
		"C#":         {".csproj", ".sln", "paket.dependencies"},
		"C++":        {"CMakeLists.txt", "Makefile", ".cpp"},
		"Swift":      {"Package.swift", ".swift"},
		"Kotlin":     {"build.gradle.kts", ".kt"},
		"Scala":      {"build.sbt", ".scala"},
		"Elixir":     {"mix.exs", ".ex"},
		"Haskell":    {"stack.yaml", "cabal.project", ".hs"},
		"Clojure":    {"project.clj", "deps.edn"},
		"Perl":       {"cpanfile", ".pl"},
		"Lua":        {".lua", "rockspec"},
		"R":          {"DESCRIPTION", ".R"},
		"Dart":       {"pubspec.yaml", ".dart"},
		"Shell":      {".sh", "Dockerfile"},
	}

	for lang, markers := range languageMarkers {
		for _, marker := range markers {
			if d.fileExists(marker) || d.hasFilesWithExtension(marker) {
				if !contains(info.Languages, lang) {
					info.Languages = append(info.Languages, lang)
					if marker != "" && !strings.HasPrefix(marker, ".") {
						info.KeyFiles = append(info.KeyFiles, marker)
					}
				}
				break
			}
		}
	}
}

func (d *Detector) detectFrameworks(info *RepositoryInfo) {
	// JavaScript/TypeScript frameworks
	if d.fileExists("package.json") {
		content, _ := os.ReadFile(filepath.Join(d.rootPath, "package.json"))
		text := string(content)
		frameworks := map[string]string{
			"\"react\"":   "React",
			"\"vue\"":     "Vue",
			"\"angular\"": "Angular",
			"\"express\"": "Express",
			"\"next\"":    "Next.js",
			"\"nuxt\"":    "Nuxt",
			"\"svelte\"":  "Svelte",
			"\"nestjs\"":  "NestJS",
			"\"gatsby\"":  "Gatsby",
			"\"remix\"":   "Remix",
		}
		for pattern, name := range frameworks {
			if strings.Contains(text, pattern) && !contains(info.Frameworks, name) {
				info.Frameworks = append(info.Frameworks, name)
			}
		}
	}

	// Python frameworks
	if d.fileExists("requirements.txt") || d.fileExists("pyproject.toml") {
		content, _ := os.ReadFile(filepath.Join(d.rootPath, "requirements.txt"))
		text := string(content)
		frameworks := map[string]string{
			"django":    "Django",
			"flask":     "Flask",
			"fastapi":   "FastAPI",
			"tornado":   "Tornado",
			"pyramid":   "Pyramid",
			"streamlit": "Streamlit",
		}
		for pattern, name := range frameworks {
			if strings.Contains(strings.ToLower(text), pattern) && !contains(info.Frameworks, name) {
				info.Frameworks = append(info.Frameworks, name)
			}
		}
	}

	// Java frameworks
	if d.fileExists("pom.xml") || d.fileExists("build.gradle") {
		content, _ := os.ReadFile(filepath.Join(d.rootPath, "pom.xml"))
		text := string(content)
		frameworks := map[string]string{
			"spring-boot":      "Spring Boot",
			"spring-framework": "Spring",
			"quarkus":          "Quarkus",
			"micronaut":        "Micronaut",
			"vertx":            "Vert.x",
		}
		for pattern, name := range frameworks {
			if strings.Contains(strings.ToLower(text), pattern) && !contains(info.Frameworks, name) {
				info.Frameworks = append(info.Frameworks, name)
			}
		}
	}

	// Ruby frameworks
	if d.fileExists("Gemfile") {
		content, _ := os.ReadFile(filepath.Join(d.rootPath, "Gemfile"))
		text := string(content)
		if strings.Contains(text, "rails") && !contains(info.Frameworks, "Rails") {
			info.Frameworks = append(info.Frameworks, "Rails")
		}
		if strings.Contains(text, "sinatra") && !contains(info.Frameworks, "Sinatra") {
			info.Frameworks = append(info.Frameworks, "Sinatra")
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
	for _, name := range []string{"README.md", "readme.md", "README", "readme.txt"} {
		path := filepath.Join(d.rootPath, name)
		if content, err := os.ReadFile(path); err == nil {
			// Clean the entire content first to handle multi-line HTML
			cleaned := cleanHTMLTags(string(content))
			lines := strings.Split(cleaned, "\n")

			var descriptionLines []string
			inDescription := false

			// Try to find a meaningful description
			for _, line := range lines {
				line = strings.TrimSpace(line)

				// Skip empty lines
				if line == "" {
					if inDescription && len(descriptionLines) > 0 {
						// Empty line after we started collecting description - might be end
						break
					}
					continue
				}

				// Skip markdown code blocks
				if strings.HasPrefix(line, "```") {
					break
				}

				// Skip markdown headings (but the content after ## is good to start collecting)
				if strings.HasPrefix(line, "###") || strings.HasPrefix(line, "####") {
					continue
				}
				if strings.HasPrefix(line, "# ") {
					// Main heading - content after this is what we want
					inDescription = true
					continue
				}
				if strings.HasPrefix(line, "## ") {
					// Sub-heading - stop collecting if we already have content
					if len(descriptionLines) > 0 {
						break
					}
					continue
				}

				// Skip badge/shield lines
				if strings.Contains(line, "badge") || strings.Contains(line, "shields.io") {
					continue
				}
				// Skip separator lines
				if strings.Trim(line, "-=_*") == "" {
					continue
				}

				// Collect substantial lines (at least 20 chars)
				if len(line) > 20 {
					descriptionLines = append(descriptionLines, line)
					inDescription = true

					// Stop if we have enough content (2-3 sentences)
					combined := strings.Join(descriptionLines, " ")
					if len(combined) > 150 || len(descriptionLines) >= 3 {
						break
					}
				}
			}

			// Combine collected lines
			if len(descriptionLines) > 0 {
				info.Description = strings.Join(descriptionLines, " ")
				if len(info.Description) > 250 {
					info.Description = info.Description[:247] + "..."
				}
				return
			}
		}
	}
	info.Description = fmt.Sprintf("A %s application", info.Type)
}

// cleanHTMLTags removes HTML tags from a string
func cleanHTMLTags(s string) string {
	// Remove multi-line HTML blocks (div, img, etc.)
	// This regex handles tags that span multiple lines
	htmlBlockRe := regexp.MustCompile(`<[^>]*>`)
	s = htmlBlockRe.ReplaceAllString(s, "")

	// Remove markdown links but keep the text
	// [text](url) -> text
	markdownLinkRe := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	s = markdownLinkRe.ReplaceAllString(s, "$1")

	// Remove markdown images
	// ![alt](url) -> ""
	markdownImgRe := regexp.MustCompile(`!\[([^\]]*)\]\([^)]+\)`)
	s = markdownImgRe.ReplaceAllString(s, "")

	// Remove standalone URLs
	urlRe := regexp.MustCompile(`https?://[^\s]+`)
	s = urlRe.ReplaceAllString(s, "")

	// Clean up multiple spaces and newlines
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	s = regexp.MustCompile(`\n+`).ReplaceAllString(s, "\n")

	return strings.TrimSpace(s)
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
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

func getPlatformName() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	case "freebsd":
		return "FreeBSD"
	default:
		return runtime.GOOS
	}
}
