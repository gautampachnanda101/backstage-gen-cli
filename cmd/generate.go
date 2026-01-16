package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/config"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/detector"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/generator"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/llm"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/wizard"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	outputFile   string
	templateName string
	interactive  bool
	dryRun       bool
	force        bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "📝 Generate catalog-info.yaml from repository inspection",
	Long: `Generate catalog-info.yaml from repository inspection

Inspects the current repository to detect:
  • Technology stack (languages, frameworks)
  • Build tools and dependencies
  • Infrastructure (Docker, Kubernetes, Helm)
  • Git information and metadata

Then generates a properly formatted catalog-info.yaml file.

Interactive mode uses AI-powered suggestions if LLM is configured.

Examples:
  # Generate with interactive wizard
  backstage-gen-cli generate -i

  # Generate with auto-detection only
  backstage-gen-cli generate

  # Preview without writing file
  backstage-gen-cli generate --dry-run

  # Force overwrite existing file
  backstage-gen-cli generate --force -i

  # Specify custom output path
  backstage-gen-cli generate -o custom-catalog.yaml`,
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "catalog-info.yaml", "Output file path")
	generateCmd.Flags().StringVarP(&templateName, "template", "t", "", "Template to use")
	generateCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode")
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be generated")
	generateCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if !force && !dryRun {
		if _, err := os.Stat(outputFile); err == nil {
			return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
		}
	}

	// Load configuration
	appConfig, err := config.Load()
	if err != nil {
		fmt.Printf("Warning: Failed to load config, using defaults: %v\n", err)
		appConfig = &config.AppConfig{}
	}

	// Initialize LLM client if configured
	var llmClient *llm.Client
	if appConfig.LLM != nil {
		llmClient = llm.NewClient(appConfig.LLM)
		if !llmClient.IsAvailable() && verboseMode {
			yellow := color.New(color.FgYellow)
			yellow.Println("⚠️  LLM provider not available. Continuing without AI suggestions.")
			yellow.Printf("   Tip: Make sure %s is running or update ~/.backstage-gen.yaml\n", appConfig.LLM.Provider)
			fmt.Println()
		}
	}

	if verboseMode {
		cyan := color.New(color.FgCyan, color.Bold)
		cyan.Printf("🔍 Analyzing repository at: %s\n", cwd)
	}

	// Detect repository information
	det := detector.New(cwd)
	info, err := det.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect repository information: %w", err)
	}

	if verboseMode {
		yellow := color.New(color.FgYellow, color.Bold)
		fmt.Println("\nDetected:")
		yellow.Printf("  ├─ Name: %s\n", info.Name)
		yellow.Printf("  ├─ Type: %s\n", info.Type)
		yellow.Printf("  └─ Languages: %v\n", info.Languages)
		fmt.Println()
	}

	// Prepare generator config
	genConfig := &generator.Config{
		Organization: generator.OrganizationConfig{
			Name:      appConfig.Organization.Name,
			Namespace: appConfig.Organization.Namespace,
		},
		Defaults: generator.DefaultsConfig{
			Owner:       appConfig.Defaults.Owner,
			System:      appConfig.Defaults.System,
			Lifecycle:   appConfig.Defaults.Lifecycle,
			Annotations: appConfig.Defaults.Annotations,
			Tags:        appConfig.Defaults.Tags,
		},
		LLM: appConfig.LLM,
	}

	if genConfig.Organization.Namespace == "" {
		genConfig.Organization.Namespace = "default"
	}
	if genConfig.Defaults.Owner == "" {
		genConfig.Defaults.Owner = "platform-team"
	}
	if genConfig.Defaults.Lifecycle == "" {
		genConfig.Defaults.Lifecycle = "production"
	}

	var catalog *generator.Catalog

	// Run interactive wizard if requested
	if interactive {
		wiz := wizard.New(llmClient)
		wizConfig, err := wiz.Run(info, genConfig)
		if err != nil {
			return fmt.Errorf("wizard failed: %w", err)
		}

		// Update info with wizard responses
		info.Name = wizConfig.Name
		info.Description = wizConfig.Description
		info.Type = wizConfig.Type
		genConfig.Defaults.Owner = wizConfig.Owner
		genConfig.Defaults.Lifecycle = wizConfig.Lifecycle
		genConfig.Defaults.System = wizConfig.System

		gen := generator.New(genConfig)
		catalog, err = gen.Generate(info, wizConfig.Type)
		if err != nil {
			return fmt.Errorf("failed to generate catalog: %w", err)
		}

		// Update tags from wizard
		if len(wizConfig.Tags) > 0 {
			catalog.Metadata.Tags = wizConfig.Tags
		}
	} else {
		// Auto-generate without wizard
		gen := generator.New(genConfig)
		template := templateName
		if template == "" {
			template = info.Type
		}
		catalog, err = gen.Generate(info, template)
		if err != nil {
			return fmt.Errorf("failed to generate catalog: %w", err)
		}
	}

	data, err := yaml.Marshal(catalog)
	if err != nil {
		return fmt.Errorf("failed to marshal catalog: %w", err)
	}

	if dryRun {
		cyan := color.New(color.FgCyan, color.Bold)
		yellow := color.New(color.FgYellow, color.Bold)
		fmt.Println()
		cyan.Println("╔════════════════════════════════════════════════════════════╗")
		cyan.Println("║          Dry Run - Preview Generated Catalog              ║")
		cyan.Println("╚════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println(string(data))
		yellow.Println("💡 This is a preview. Use without --dry-run to save the file.")
		fmt.Println()
		return nil
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	green := color.New(color.FgGreen, color.Bold)
	absPath, _ := filepath.Abs(outputFile)
	fmt.Println()
	green.Println("✓ Catalog Generated Successfully!")
	fmt.Printf("  📄 File: %s\n", color.CyanString(absPath))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Review: %s\n", color.CyanString("cat "+outputFile))
	fmt.Printf("  2. Validate: %s\n", color.CyanString("backstage-gen-cli lint"))
	fmt.Printf("  3. Commit: %s\n", color.CyanString("git add "+outputFile))
	fmt.Println()
	return nil
}
