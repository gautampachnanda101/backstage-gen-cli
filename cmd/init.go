package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/config"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/llm"
	"github.com/spf13/cobra"
)

var (
	skipLLM bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "🚀 Initialize backstage-gen with prerequisites validation",
	Long: `Initialize backstage-gen CLI tool

This command:
  • Validates all prerequisites (Docker, etc.)
  • Sets up LiteLLM with Docker for AI-powered suggestions
  • Creates example configuration file
  • Verifies the setup is working

Prerequisites:
  • Docker (optional, for LiteLLM AI features)`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&skipLLM, "skip-llm", false, "Skip LiteLLM setup")
}

func runInit(cmd *cobra.Command, args []string) error {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	red := color.New(color.FgRed, color.Bold)

	fmt.Println()
	cyan.Println("╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║          Backstage-Gen Initialization                     ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Step 1: Check prerequisites
	yellow.Println("📋 Step 1: Checking Prerequisites")
	fmt.Println()

	dockerOK := llm.CheckDockerRunning()
	if dockerOK {
		green.Println("  ✓ Docker is installed and running")
	} else {
		if skipLLM {
			yellow.Println("  ⚠️  Docker not found (skipping LiteLLM setup)")
		} else {
			red.Println("  ✗ Docker is not installed or not running")
			fmt.Println()
			fmt.Println("  To use AI-powered features, please:")
			fmt.Printf("    1. Install Docker: %s\n", color.CyanString("https://www.docker.com/get-started"))
			fmt.Println("    2. Start Docker")
			fmt.Printf("    3. Run: %s\n", color.CyanString("backstage-gen-cli init"))
			fmt.Println()
			fmt.Printf("  Or run with %s to skip LLM setup\n", color.CyanString("--skip-llm"))
			fmt.Println()
			return fmt.Errorf("docker prerequisite not met")
		}
	}

	fmt.Println()

	// Step 2: Create configuration file
	yellow.Println("📝 Step 2: Creating Configuration")
	fmt.Println()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".backstage-gen.yaml")

	if _, err := os.Stat(configPath); err == nil {
		yellow.Printf("  Config file already exists: %s\n", configPath)
	} else {
		if err := config.CreateExampleConfig(configPath); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
		green.Println("  ✓ Created configuration file")
		fmt.Printf("    Location: %s\n", color.CyanString(configPath))
	}

	fmt.Println()

	// Step 3: Setup LiteLLM
	if dockerOK && !skipLLM {
		yellow.Println("🤖 Step 3: Setting up LiteLLM")
		fmt.Println()

		if llm.CheckLiteLLMRunning() {
			green.Println("  ✓ LiteLLM is already running")
		} else {
			fmt.Println("  Starting LiteLLM Docker container...")
			if err := llm.StartLiteLLM(); err != nil {
				yellow.Printf("  ⚠️  Failed to start LiteLLM: %v\n", err)
				fmt.Println()
				fmt.Println("  You can start it manually:")
				fmt.Printf("    %s\n", color.CyanString("docker run -d --name litellm -p 4000:4000 ghcr.io/berriai/litellm:main-latest"))
			} else {
				green.Println("  ✓ LiteLLM container started successfully")
				fmt.Println("    Waiting for LiteLLM to be ready...")

				// Wait a bit for container to start
				for i := 0; i < 10; i++ {
					if llm.CheckLiteLLMRunning() {
						break
					}
					fmt.Print(".")
					// time.Sleep(1 * time.Second)
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}

	// Step 4: Verify setup
	yellow.Println("✅ Step 4: Verification")
	fmt.Println()

	appConfig, err := config.Load()
	if err == nil {
		green.Println("  ✓ Configuration loaded successfully")

		if !skipLLM && dockerOK {
			llmClient := llm.NewClient(appConfig.LLM)
			if llmClient.IsAvailable() {
				green.Println("  ✓ LLM provider is available")
			} else {
				yellow.Println("  ⚠️  LLM provider not responding yet")
				fmt.Println("    It may take a minute for LiteLLM to start")
			}
		}
	}

	fmt.Println()

	// Summary
	cyan.Println("╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║          Setup Complete!                                   ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Println("🎉 You're all set! Here's what you can do next:")
	fmt.Println()
	fmt.Println("  1. Inspect a repository:")
	fmt.Printf("     %s\n", color.CyanString("backstage-gen inspect"))
	fmt.Println()
	fmt.Println("  2. Generate catalog with AI wizard:")
	fmt.Printf("     %s\n", color.CyanString("backstage-gen generate -i"))
	fmt.Println()
	fmt.Println("  3. Generate without wizard:")
	fmt.Printf("     %s\n", color.CyanString("backstage-gen generate"))
	fmt.Println()

	if skipLLM || !dockerOK {
		fmt.Println("💡 Tips:")
		fmt.Println("  • Install Docker to enable AI-powered suggestions")
		fmt.Printf("  • Run %s after installing Docker\n", color.CyanString("backstage-gen init"))
		fmt.Println()
	} else {
		fmt.Println("💡 Tips:")
		fmt.Printf("  • Edit config: %s\n", color.CyanString("nano "+configPath))
		fmt.Println("  • LiteLLM supports 100+ models via different providers")
		fmt.Println("  • Use interactive mode (-i) for best results")
		fmt.Println()
	}

	return nil
}
