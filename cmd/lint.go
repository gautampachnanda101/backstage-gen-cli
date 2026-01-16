package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/yourusername/backstage-gen/pkg/validator"
)

var (
	lintFile     string
	strictMode   bool
	outputFormat string
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Validate and lint Backstage catalog files",
	Long:  `Validates catalog-info.yaml against Backstage schema and best practices.`,
	RunE:  runLint,
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
