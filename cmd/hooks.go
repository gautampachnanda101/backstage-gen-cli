package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var hooksCmd = &cobra.Command{
	Use:   "hooks [install|uninstall]",
	Short: "Manage git hooks",
	Args:  cobra.ExactArgs(1),
	RunE:  runHooks,
}

func init() {
	rootCmd.AddCommand(hooksCmd)
}

func runHooks(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "install":
		return installHook()
	case "uninstall":
		return uninstallHook()
	default:
		return fmt.Errorf("unknown action: %s", args[0])
	}
}

func installHook() error {
	gitDir, err := findGitDir()
	if err != nil {
		return err
	}

	hooksDir := filepath.Join(gitDir, "hooks")
	hookPath := filepath.Join(hooksDir, "pre-commit")

	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return err
	}

	hookScript := `#!/bin/sh
# backstage-gen pre-commit hook
echo "Running backstage-gen validation..."
if [ -f "catalog-info.yaml" ]; then
    backstage-gen lint
    if [ $? -ne 0 ]; then
        echo "❌ Validation failed"
        exit 1
    fi
fi
echo "✓ Validation passed"
exit 0
`

	if err := os.WriteFile(hookPath, []byte(hookScript), 0755); err != nil {
		return err
	}

	fmt.Println(color.GreenString("✓") + " Hook installed: " + hookPath)
	return nil
}

func uninstallHook() error {
	gitDir, err := findGitDir()
	if err != nil {
		return err
	}

	hookPath := filepath.Join(gitDir, "hooks", "pre-commit")
	
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		fmt.Println("Hook not installed")
		return nil
	}

	if err := os.Remove(hookPath); err != nil {
		return err
	}

	fmt.Println(color.GreenString("✓") + " Hook uninstalled")
	return nil
}

func findGitDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := cwd
	for {
		gitDir := filepath.Join(dir, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			return gitDir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf(".git directory not found")
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
