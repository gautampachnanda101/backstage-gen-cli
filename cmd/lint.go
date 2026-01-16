package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gautampachnanda101/backstage-gen-cli/pkg/validator"
	"github.com/spf13/cobra"
)

var (
	lintFile     string
	strictMode   bool
	outputFormat string
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "✅ Validate and lint Backstage catalog files",
	Long: `Validate and lint Backstage catalog files

Validates catalog-info.yaml files against:
  • Backstage schema requirements
  • Required metadata fields
  • Best practices and conventions
  • YAML syntax and structure

Examples:
  # Validate default catalog file
  backstage-gen-cli lint

  # Validate specific file
  backstage-gen-cli lint -f custom-catalog.yaml

  # Strict mode (warnings as errors)
  backstage-gen-cli lint --strict

  # Use in CI/CD
  backstage-gen-cli lint --strict --format json`,
	RunE: runLint,
}

func init() {
	rootCmd.AddCommand(lintCmd)
	lintCmd.Flags().StringVarP(&lintFile, "file", "f", "catalog-info.yaml", "File to validate")
	lintCmd.Flags().BoolVar(&strictMode, "strict", false, "Treat warnings as errors")
	lintCmd.Flags().StringVar(&outputFormat, "format", "text", "Output format")
}

func runLint(cmd *cobra.Command, args []string) error {
	val := validator.New(&validator.Config{Strict: strictMode})

	result, err := val.ValidateFile(lintFile)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	hasErrors := false
	if len(result.Errors) > 0 {
		fmt.Printf("\n%s %s\n", color.RedString("✗"), color.RedString(result.File))
		for _, e := range result.Errors {
			fmt.Printf("  %s %s\n", color.RedString("ERROR:"), e.Message)
			if e.Field != "" {
				fmt.Printf("    Field: %s\n", e.Field)
			}
		}
		hasErrors = true
	}

	if len(result.Warnings) > 0 {
		if len(result.Errors) == 0 {
			fmt.Printf("\n%s %s\n", color.YellowString("⚠"), result.File)
		}
		for _, w := range result.Warnings {
			fmt.Printf("  %s %s\n", color.YellowString("WARNING:"), w.Message)
			if w.Field != "" {
				fmt.Printf("    Field: %s\n", w.Field)
			}
		}
		if strictMode {
			hasErrors = true
		}
	}

	if !hasErrors {
		fmt.Println(color.GreenString("\n✓ All validations passed"))
		return nil
	}

	return fmt.Errorf("validation failed")
}
