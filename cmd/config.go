package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "⚙️  Manage configuration",
	Long: `Manage configuration for backstage-gen CLI

Configuration file can be placed at:
  • ~/.backstage-gen-cli.yaml
  • ~/.config/backstage-gen-cli/config.yaml
  • .backstage-gen-cli.yaml (current directory)

The configuration includes:
  • LLM provider settings for AI-powered suggestions
  • Organization defaults
  • Component defaults`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create example configuration file",
	Long: `Creates an example configuration file in your home directory

The file will be created at ~/.backstage-gen.yaml and includes
examples for different LLM providers including Ollama and OpenRouter.`,
	RunE: runConfigInit,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".backstage-gen-cli.yaml")

	if _, err := os.Stat(configPath); err == nil {
		green := color.New(color.FgYellow)
		green.Printf("⚠️  Config file already exists: %s\n", configPath)
		fmt.Print("Overwrite? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	if err := config.CreateExampleConfig(configPath); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	green := color.New(color.FgGreen, color.Bold)
	fmt.Println()
	green.Println("✓ Configuration File Created!")
	fmt.Printf("  📄 Location: %s\n", color.CyanString(configPath))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Edit config: %s\n", color.CyanString("nano "+configPath))
	fmt.Println("  2. Configure your LLM provider (Ollama, OpenRouter, etc.)")
	fmt.Println("  3. Set organization defaults")
	fmt.Println()
	fmt.Println("💡 Tip: Install Ollama for free local AI:")
	fmt.Printf("   %s\n", color.CyanString("https://ollama.ai"))
	fmt.Println()

	return nil
}
