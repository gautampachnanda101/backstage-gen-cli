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
	Short: "🪝 Manage git hooks",
	Long: `Manage git hooks for automated validation

Installs or removes a pre-commit hook that automatically
validates catalog-info.yaml before each commit.

Commands:
  install      Install the pre-commit hook
  uninstall    Remove the pre-commit hook

Examples:
  # Install hook
  backstage-gen-cli hooks install

  # Remove hook
  backstage-gen-cli hooks uninstall

The hook will run 'backstage-gen-cli lint' before each commit.`,
	Args: cobra.ExactArgs(1),
	RunE: runHooks,
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
# backstage-gen-cli pre-commit hook
echo "Running backstage-gen-cli validation..."
if [ -f "catalog-info.yaml" ]; then
    backstage-gen-cli lint
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

	green := color.New(color.FgGreen, color.Bold)
	fmt.Println()
	green.Println("✔ Git Hook Installed Successfully!")
	fmt.Printf("  Location: %s\n", color.CyanString(hookPath))
	fmt.Println("\n  The hook will now validate catalog-info.yaml before each commit.")
	fmt.Println()
	return nil
}

func uninstallHook() error {
	gitDir, err := findGitDir()
	if err != nil {
		return err
	}

	hookPath := filepath.Join(gitDir, "hooks", "pre-commit")

	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		fmt.Println(color.YellowString("⚠ Hook not currently installed"))
		return nil
	}

	if err := os.Remove(hookPath); err != nil {
		return err
	}

	green := color.New(color.FgGreen, color.Bold)
	fmt.Println()
	green.Println("✔ Git Hook Uninstalled Successfully!")
	fmt.Printf("  Removed: %s\n", color.CyanString(hookPath))
	fmt.Println()
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
