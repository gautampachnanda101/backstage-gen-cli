package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	appconfig "github.com/gautampachnanda101/backstage-gen-cli/pkg/config"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/detector"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/llm"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/output"
	"github.com/spf13/cobra"
)

var jsonOutput bool

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "🔍 Inspect repository and show detected information",
	Long: `Inspect repository and show detected information

Analyzes your repository to detect:
  • System information (OS, platform, architecture)
  • Technology stack (languages, frameworks, build tools)
  • Infrastructure components (Docker, Kubernetes, etc.)
  • Git repository information

Examples:
  # Inspect with formatted output
  backstage-gen-cli inspect

  # Output as JSON for scripting
  backstage-gen-cli inspect --json

  # Use in CI/CD pipelines
  backstage-gen-cli inspect --json | jq '.languages[]'`,
	RunE: runInspect,
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
}

func runInspect(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Try to load config for LLM support
	appConfig, err := appconfig.Load()
	if err != nil {
		output.PrintDebug("Failed to load config: %v\n", err)
	}

	var llmClient *llm.Client
	if appConfig != nil && appConfig.LLM != nil {
		llmClient = llm.NewClient(appConfig.LLM)
		if llmClient.IsAvailable() {
			output.PrintDebug("LLM client is available\n")
		} else {
			output.PrintDebug("LLM client is NOT available\n")
		}
	} else {
		output.PrintDebug("No LLM config found\n")
	}

	// Use LLM-enhanced detector if available
	var det *detector.Detector
	if llmClient != nil && llmClient.IsAvailable() {
		det = detector.NewWithLLM(cwd, llmClient)
		output.PrintDebug("Using LLM-enhanced detector\n")
	} else {
		det = detector.New(cwd)
		output.PrintDebug("Using standard detector\n")
	}

	info, err := det.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect: %w", err)
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	// In quiet mode, just exit successfully
	if output.IsQuiet() {
		return nil
	}

	// Header
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	magenta := color.New(color.FgMagenta)

	fmt.Println()
	cyan.Println("╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║          Repository Inspection Report                     ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// System Information
	yellow.Println("⚙️  System Information")
	fmt.Printf("  ├─ Platform: %s\n", green.Sprint(info.Platform))
	fmt.Printf("  ├─ OS: %s\n", green.Sprint(info.OS))
	fmt.Printf("  └─ Architecture: %s\n", green.Sprint(info.Arch))
	fmt.Println()

	// Basic Information
	yellow.Println("📦 Repository Information")
	fmt.Printf("  ├─ Name: %s\n", green.Sprint(info.Name))
	fmt.Printf("  ├─ Type: %s\n", magenta.Sprint(info.Type))
	fmt.Printf("  ├─ Path: %s\n", info.Path)
	if info.Description != "" {
		fmt.Printf("  └─ Description: %s\n", info.Description)
	}
	fmt.Println()

	// Technology Stack
	if len(info.Languages) > 0 {
		yellow.Println("💻 Technology Stack")
		for i, lang := range info.Languages {
			if i == len(info.Languages)-1 && len(info.Frameworks) == 0 && len(info.BuildTools) == 0 {
				fmt.Printf("  └─ Language: %s\n", green.Sprint(lang))
			} else {
				fmt.Printf("  ├─ Language: %s\n", green.Sprint(lang))
			}
		}

		if len(info.Frameworks) > 0 {
			for i, fw := range info.Frameworks {
				if i == len(info.Frameworks)-1 && len(info.BuildTools) == 0 {
					fmt.Printf("  └─ Framework: %s\n", green.Sprint(fw))
				} else {
					fmt.Printf("  ├─ Framework: %s\n", green.Sprint(fw))
				}
			}
		}

		if len(info.BuildTools) > 0 {
			for i, tool := range info.BuildTools {
				if i == len(info.BuildTools)-1 {
					fmt.Printf("  └─ Build Tool: %s\n", green.Sprint(tool))
				} else {
					fmt.Printf("  ├─ Build Tool: %s\n", green.Sprint(tool))
				}
			}
		}
		fmt.Println()
	}

	// Infrastructure
	infraItems := []string{}
	if info.HasDocker {
		infraItems = append(infraItems, "🐳 Docker")
	}
	if info.HasKubernetes {
		infraItems = append(infraItems, "☸️  Kubernetes")
	}
	if info.HasHelm {
		infraItems = append(infraItems, "⎈  Helm")
	}
	if info.HasTerraform {
		infraItems = append(infraItems, "🏗️  Terraform")
	}

	if len(infraItems) > 0 {
		yellow.Println("🏗️  Infrastructure")
		for i, item := range infraItems {
			if i == len(infraItems)-1 {
				fmt.Printf("  └─ %s\n", green.Sprint(item))
			} else {
				fmt.Printf("  ├─ %s\n", green.Sprint(item))
			}
		}
		fmt.Println()
	}

	// Git Information
	if info.GitRemote != "" {
		yellow.Println("🔗 Git Information")
		fmt.Printf("  └─ Remote: %s\n", green.Sprint(info.GitRemote))
		fmt.Println()
	}

	cyan.Println("✓ Inspection complete!")
	fmt.Println()
	return nil
}
