package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	det := New("/tmp/test")
	if det == nil {
		t.Fatal("expected non-nil detector")
	}
	if det.rootPath != "/tmp/test" {
		t.Errorf("expected rootPath /tmp/test, got %s", det.rootPath)
	}
	if det.llmClient != nil {
		t.Error("expected nil llmClient for standard detector")
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "detector-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	det := New(tmpDir)

	// Test existing file
	if !det.fileExists("test.txt") {
		t.Error("expected test.txt to exist")
	}

	// Test non-existing file
	if det.fileExists("nonexistent.txt") {
		t.Error("expected nonexistent.txt to not exist")
	}
}

func TestDetectLanguagesGo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "detector-test-go")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod file
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	det := New(tmpDir)
	info := &RepositoryInfo{}
	det.detectLanguages(info)

	if len(info.Languages) != 1 {
		t.Fatalf("expected 1 language, got %d", len(info.Languages))
	}
	if info.Languages[0] != "Go" {
		t.Errorf("expected Go, got %s", info.Languages[0])
	}
}

func TestDetectLanguagesPython(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "detector-test-python")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create requirements.txt file
	reqFile := filepath.Join(tmpDir, "requirements.txt")
	if err := os.WriteFile(reqFile, []byte("flask==2.0.0"), 0644); err != nil {
		t.Fatalf("failed to create requirements.txt: %v", err)
	}

	det := New(tmpDir)
	info := &RepositoryInfo{}
	det.detectLanguages(info)

	if len(info.Languages) != 1 {
		t.Fatalf("expected 1 language, got %d", len(info.Languages))
	}
	if info.Languages[0] != "Python" {
		t.Errorf("expected Python, got %s", info.Languages[0])
	}
}

func TestDetectLanguagesRust(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "detector-test-rust")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create Cargo.toml file
	cargoFile := filepath.Join(tmpDir, "Cargo.toml")
	if err := os.WriteFile(cargoFile, []byte("[package]\nname = \"test\""), 0644); err != nil {
		t.Fatalf("failed to create Cargo.toml: %v", err)
	}

	det := New(tmpDir)
	info := &RepositoryInfo{}
	det.detectLanguages(info)

	if len(info.Languages) != 1 {
		t.Fatalf("expected 1 language, got %d", len(info.Languages))
	}
	if info.Languages[0] != "Rust" {
		t.Errorf("expected Rust, got %s", info.Languages[0])
	}
}

func TestDetectLanguagesJavaScript(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "detector-test-js")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create package.json file
	pkgFile := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(pkgFile, []byte(`{"name": "test"}`), 0644); err != nil {
		t.Fatalf("failed to create package.json: %v", err)
	}

	det := New(tmpDir)
	info := &RepositoryInfo{}
	det.detectLanguages(info)

	if len(info.Languages) != 1 {
		t.Fatalf("expected 1 language, got %d", len(info.Languages))
	}
	if info.Languages[0] != "JavaScript" {
		t.Errorf("expected JavaScript, got %s", info.Languages[0])
	}
}

func TestDetectLanguagesTypeScript(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "detector-test-ts")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create tsconfig.json file
	tsFile := filepath.Join(tmpDir, "tsconfig.json")
	if err := os.WriteFile(tsFile, []byte(`{}`), 0644); err != nil {
		t.Fatalf("failed to create tsconfig.json: %v", err)
	}

	det := New(tmpDir)
	info := &RepositoryInfo{}
	det.detectLanguages(info)

	if len(info.Languages) != 1 {
		t.Fatalf("expected 1 language, got %d", len(info.Languages))
	}
	if info.Languages[0] != "TypeScript" {
		t.Errorf("expected TypeScript, got %s", info.Languages[0])
	}
}

func TestDetectLanguagesJava(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "detector-test-java")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create pom.xml file
	pomFile := filepath.Join(tmpDir, "pom.xml")
	if err := os.WriteFile(pomFile, []byte(`<project></project>`), 0644); err != nil {
		t.Fatalf("failed to create pom.xml: %v", err)
	}

	det := New(tmpDir)
	info := &RepositoryInfo{}
	det.detectLanguages(info)

	if len(info.Languages) != 1 {
		t.Fatalf("expected 1 language, got %d", len(info.Languages))
	}
	if info.Languages[0] != "Java" {
		t.Errorf("expected Java, got %s", info.Languages[0])
	}
}

func TestDetectBuildTools(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "detector-test-build")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create Makefile
	makeFile := filepath.Join(tmpDir, "Makefile")
	if err := os.WriteFile(makeFile, []byte("all:\n\techo test"), 0644); err != nil {
		t.Fatalf("failed to create Makefile: %v", err)
	}

	det := New(tmpDir)
	info := &RepositoryInfo{}
	det.detectBuildTools(info)

	found := false
	for _, tool := range info.BuildTools {
		if tool == "Make" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected Make in build tools, got %v", info.BuildTools)
	}
}

func TestDetectInfrastructureDocker(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "detector-test-docker")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create Dockerfile
	dockerfile := filepath.Join(tmpDir, "Dockerfile")
	if err := os.WriteFile(dockerfile, []byte("FROM alpine"), 0644); err != nil {
		t.Fatalf("failed to create Dockerfile: %v", err)
	}

	det := New(tmpDir)
	info := &RepositoryInfo{}
	det.detectInfrastructure(info)

	if !info.HasDocker {
		t.Error("expected HasDocker to be true")
	}
}

func TestDetectInfrastructureKubernetes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "detector-test-k8s")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create k8s directory with deployment.yaml
	k8sDir := filepath.Join(tmpDir, "k8s")
	if err := os.MkdirAll(k8sDir, 0755); err != nil {
		t.Fatalf("failed to create k8s dir: %v", err)
	}
	deployment := filepath.Join(k8sDir, "deployment.yaml")
	if err := os.WriteFile(deployment, []byte("apiVersion: apps/v1\nkind: Deployment"), 0644); err != nil {
		t.Fatalf("failed to create deployment.yaml: %v", err)
	}

	det := New(tmpDir)
	info := &RepositoryInfo{}
	det.detectInfrastructure(info)

	if !info.HasKubernetes {
		t.Error("expected HasKubernetes to be true")
	}
}

func TestDetectInfrastructureHelm(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "detector-test-helm")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create Chart.yaml
	chartFile := filepath.Join(tmpDir, "Chart.yaml")
	if err := os.WriteFile(chartFile, []byte("name: test-chart"), 0644); err != nil {
		t.Fatalf("failed to create Chart.yaml: %v", err)
	}

	det := New(tmpDir)
	info := &RepositoryInfo{}
	det.detectInfrastructure(info)

	if !info.HasHelm {
		t.Error("expected HasHelm to be true")
	}
}

func TestDetectInfrastructureTerraform(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "detector-test-tf")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create main.tf
	tfFile := filepath.Join(tmpDir, "main.tf")
	if err := os.WriteFile(tfFile, []byte("resource \"null\" \"test\" {}"), 0644); err != nil {
		t.Fatalf("failed to create main.tf: %v", err)
	}

	det := New(tmpDir)
	info := &RepositoryInfo{}
	det.detectInfrastructure(info)

	if !info.HasTerraform {
		t.Error("expected HasTerraform to be true")
	}
}

func TestDetectType(t *testing.T) {
	t.Run("resource with kubernetes", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "detector-test")
		defer os.RemoveAll(tmpDir)
		det := New(tmpDir)
		info := &RepositoryInfo{HasKubernetes: true}
		det.detectType(info)
		if info.Type != "resource" {
			t.Errorf("expected type resource, got %s", info.Type)
		}
	})

	t.Run("resource with terraform", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "detector-test")
		defer os.RemoveAll(tmpDir)
		det := New(tmpDir)
		info := &RepositoryInfo{HasTerraform: true}
		det.detectType(info)
		if info.Type != "resource" {
			t.Errorf("expected type resource, got %s", info.Type)
		}
	})

	t.Run("resource with helm", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "detector-test")
		defer os.RemoveAll(tmpDir)
		det := New(tmpDir)
		info := &RepositoryInfo{HasHelm: true}
		det.detectType(info)
		if info.Type != "resource" {
			t.Errorf("expected type resource, got %s", info.Type)
		}
	})

	t.Run("website with index.html", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "detector-test")
		defer os.RemoveAll(tmpDir)
		// Create index.html to trigger website detection
		os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte("<html></html>"), 0644)
		det := New(tmpDir)
		info := &RepositoryInfo{}
		det.detectType(info)
		if info.Type != "website" {
			t.Errorf("expected type website, got %s", info.Type)
		}
	})

	t.Run("library with setup.py", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "detector-test")
		defer os.RemoveAll(tmpDir)
		// Create setup.py without main.py to trigger library detection
		os.WriteFile(filepath.Join(tmpDir, "setup.py"), []byte("setup()"), 0644)
		det := New(tmpDir)
		info := &RepositoryInfo{}
		det.detectType(info)
		if info.Type != "library" {
			t.Errorf("expected type library, got %s", info.Type)
		}
	})

	t.Run("service with docker", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "detector-test")
		defer os.RemoveAll(tmpDir)
		det := New(tmpDir)
		info := &RepositoryInfo{HasDocker: true}
		det.detectType(info)
		if info.Type != "service" {
			t.Errorf("expected type service, got %s", info.Type)
		}
	})

	t.Run("default to service", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "detector-test")
		defer os.RemoveAll(tmpDir)
		det := New(tmpDir)
		info := &RepositoryInfo{}
		det.detectType(info)
		if info.Type != "service" {
			t.Errorf("expected type service, got %s", info.Type)
		}
	})
}

// Note: Description cleaning is tested indirectly through Detect()

func TestContains(t *testing.T) {
	slice := []string{"Go", "Python", "Rust"}

	if !contains(slice, "Go") {
		t.Error("expected contains to return true for Go")
	}
	if !contains(slice, "Python") {
		t.Error("expected contains to return true for Python")
	}
	if contains(slice, "Java") {
		t.Error("expected contains to return false for Java")
	}
}

func TestDetectMultipleLanguages(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "detector-test-multi")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create multiple language markers
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	pkgJson := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(pkgJson, []byte(`{"name": "test"}`), 0644); err != nil {
		t.Fatalf("failed to create package.json: %v", err)
	}

	det := New(tmpDir)
	info := &RepositoryInfo{}
	det.detectLanguages(info)

	if len(info.Languages) < 2 {
		t.Fatalf("expected at least 2 languages, got %d", len(info.Languages))
	}

	hasGo := false
	hasJS := false
	for _, lang := range info.Languages {
		if lang == "Go" {
			hasGo = true
		}
		if lang == "JavaScript" {
			hasJS = true
		}
	}
	if !hasGo {
		t.Error("expected Go in languages")
	}
	if !hasJS {
		t.Error("expected JavaScript in languages")
	}
}
