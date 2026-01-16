package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yourusername/backstage-gen/pkg/detector"
	"github.com/yourusername/backstage-gen/pkg/generator"
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
	Short: "Generate catalog-info.yaml from repository inspection",
	Long: `Inspects the current repository to detect technology stack, dependencies,
and organizational patterns, then generates a catalog-info.yaml file.`,
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

	if verboseMode {
		fmt.Printf("Analyzing repository at: %s\n", cwd)
	}

	if !force && !dryRun {
		if _, err := os.Stat(outputFile); err == nil {
			return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
		}
	}

	det := detector.New(cwd)
	info, err := det.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect repository information: %w", err)
	}

	if verboseMode {
		fmt.Println("\nDetected:")
		fmt.Printf("  Name: %s\n", info.Name)
		fmt.Printf("  Type: %s\n", info.Type)
		fmt.Printf("  Languages: %v\n", info.Languages)
	}

	config := &generator.Config{
		Organization: generator.OrganizationConfig{
			Name:      viper.GetString("organization.name"),
			Namespace: viper.GetString("organization.namespace"),
		},
		Defaults: generator.DefaultsConfig{
			Owner:     viper.GetString("defaults.owner"),
			System:    viper.GetString("defaults.system"),
			Lifecycle: viper.GetString("defaults.lifecycle"),
		},
	}

	if config.Organization.Namespace == "" {
		config.Organization.Namespace = "default"
	}
	if config.Defaults.Owner == "" {
		config.Defaults.Owner = "platform-team"
	}
	if config.Defaults.Lifecycle == "" {
		config.Defaults.Lifecycle = "production"
	}

	gen := generator.New(config)
	template := templateName
	if template == "" {
		template = info.Type
	}

	catalog, err := gen.Generate(info, template)
	if err != nil {
		return fmt.Errorf("failed to generate catalog: %w", err)
	}

	data, err := yaml.Marshal(catalog)
	if err != nil {
		return fmt.Errorf("failed to marshal catalog: %w", err)
	}

	if dryRun {
		fmt.Println(color.YellowString("=== Dry Run ==="))
		fmt.Println(string(data))
		fmt.Println(color.YellowString("=== End Dry Run ==="))
		return nil
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	absPath, _ := filepath.Abs(outputFile)
	fmt.Println(color.GreenString("✓") + " Generated: " + absPath)
	return nil
}
